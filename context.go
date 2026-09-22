package llvm

import (
	"io"
	"sync"

	"github.com/kkkunny/go-llvm/internal/binding"
)

// Context LLVM 上下文，也是资源所有权的根：Close 时级联关闭所有登记的子资源
type Context struct {
	ref     binding.LLVMContextRef
	life    *Lifetime
	mu      sync.Mutex
	closers []io.Closer
}

// NewContext 创建上下文
func NewContext() *Context {
	return &Context{ref: binding.LLVMContextCreate(), life: NewLifetime()}
}

// Own 登记子资源（供 llvm/* 子包使用）；Context.Close 时按逆序级联 Close
func (ctx *Context) Own(c io.Closer) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.closers = append(ctx.closers, c)
}

// Close 级联关闭子资源后释放 Context；二次调用返回 ErrClosed
func (ctx *Context) Close() error {
	if !ctx.life.Alive() {
		return &Error{Reason: ErrClosed, Op: "llvm.Context.Close", Msg: "context already closed"}
	}
	ctx.life.Kill()
	ctx.mu.Lock()
	closers := ctx.closers
	ctx.closers = nil
	ctx.mu.Unlock()
	for i := len(closers) - 1; i >= 0; i-- {
		_ = closers[i].Close()
	}
	binding.LLVMContextDispose(ctx.ref)
	return nil
}

// Ref 返回底层句柄（供 llvm/* 子包桥接使用）
func (ctx *Context) Ref() binding.LLVMContextRef { return ctx.ref }

// Lifetime 返回上下文生命周期令牌
func (ctx *Context) Lifetime() *Lifetime { return ctx.life }

// Alive 上下文是否存活
func (ctx *Context) Alive() bool { return ctx.life.Alive() }
