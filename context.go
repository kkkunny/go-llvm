package llvm

import (
	"io"
	"sync"

	"github.com/kkkunny/go-llvm/internal/binding"
)

// Context LLVM 上下文，也是资源所有权的根：Close 时级联关闭所有登记的子资源
type Context struct {
	ref      binding.LLVMContextRef
	life     *Lifetime
	mu       sync.Mutex
	nextID   uint64
	closers  []owned
	disowned bool
}

// owned 登记项；id 用于注销时精确定位，避免依赖接口值的可比性
type owned struct {
	id uint64
	c  io.Closer
}

// NewContext 创建上下文
func NewContext() *Context {
	return &Context{ref: binding.LLVMContextCreate(), life: NewLifetime()}
}

// Own 登记子资源（供 llvm/* 子包使用）；Context.Close 时按逆序级联 Close。
// 返回注销函数：资源自行关闭或移交所有权后应调用以解除登记（幂等）。
func (ctx *Context) Own(c io.Closer) func() {
	ctx.CheckAlive("llvm.Context.Own")
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	id := ctx.nextID
	ctx.nextID++
	ctx.closers = append(ctx.closers, owned{id: id, c: c})
	return func() {
		ctx.mu.Lock()
		defer ctx.mu.Unlock()
		for i := range ctx.closers {
			if ctx.closers[i].id == id {
				ctx.closers = append(ctx.closers[:i], ctx.closers[i+1:]...)
				return
			}
		}
	}
}

// Close 级联关闭子资源后释放 Context；二次调用返回 ErrClosed
func (ctx *Context) Close() error {
	if ctx.disowned {
		return &Error{Reason: ErrClosed, Op: "llvm.Context.Close", Msg: "context ownership has been transferred"}
	}
	if !ctx.life.Alive() {
		return &Error{Reason: ErrClosed, Op: "llvm.Context.Close", Msg: "context already closed"}
	}
	ctx.mu.Lock()
	closers := ctx.closers
	ctx.closers = nil
	ctx.mu.Unlock()
	// 先级联关闭子资源（此时 Context 仍存活），再失效令牌并释放底层上下文
	for i := len(closers) - 1; i >= 0; i-- {
		_ = closers[i].c.Close()
	}
	ctx.life.Kill()
	binding.LLVMContextDispose(ctx.ref)
	return nil
}

// Disown 立即解除与外部接管方（如 JIT ThreadSafeModule）的所有权：级联关闭未移交的子资源
// （Builder 等），但不释放底层上下文（由接管方释放）；其下所有 Go 侧句柄立即失效。
func (ctx *Context) Disown() {
	ctx.mu.Lock()
	if ctx.disowned || !ctx.life.Alive() {
		ctx.mu.Unlock()
		return
	}
	ctx.disowned = true
	closers := ctx.closers
	ctx.closers = nil
	ctx.mu.Unlock()
	for i := len(closers) - 1; i >= 0; i-- {
		_ = closers[i].c.Close()
	}
	ctx.life.Kill()
}

// Ref 返回底层句柄（供 llvm/* 子包桥接使用）
func (ctx *Context) Ref() binding.LLVMContextRef { return ctx.ref }

// Lifetime 返回上下文生命周期令牌
func (ctx *Context) Lifetime() *Lifetime { return ctx.life }

// Alive 上下文是否存活
func (ctx *Context) Alive() bool { return ctx.life.Alive() }

// CheckAlive 校验上下文存活；类型/常量/值构造入口统一调用（供 llvm/* 子包使用）
func (ctx *Context) CheckAlive(op string) {
	if !ctx.life.Alive() {
		errPanic(ErrUseAfterFree, op, "context is closed")
	}
}
