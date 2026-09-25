package llvm

import (
	"io"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
)

// Context LLVM 上下文，也是资源所有权的根：Close 时级联关闭所有登记的子资源
type Context struct {
	ref      binding.LLVMContextRef
	life     *Lifetime
	diagID   uint64 // 诊断回调注册表键（全局唯一）
	mu       sync.Mutex
	nextID   uint64
	closers  []owned
	disowned bool

	cacheMu   sync.Mutex
	typeCache map[reflect.Type]AnyType // Go 类型 → LLVM 类型映射缓存
	fnCache   map[reflect.Type]FnType  // Go 函数签名 → LLVM 函数类型缓存
}

// owned 登记项；id 用于注销时精确定位，避免依赖接口值的可比性
type owned struct {
	id uint64
	c  io.Closer
}

// diagIDSeq 诊断回调注册表键的自增序列
var diagIDSeq atomic.Uint64

// NewContext 创建上下文。
// 默认安装诊断回调（A3）：LLVM 默认 handler 对 error 会退出/abort 进程，安装后
// 转为日志输出，进程得以存活并由上层继续报错；用户可用 SetDiagnosticHandler 覆盖。
func NewContext() *Context {
	ctx := &Context{ref: binding.LLVMContextCreate(), life: NewLifetime(), diagID: diagIDSeq.Add(1)}
	ctx.SetDiagnosticHandler(defaultDiagnosticHandler)
	return ctx
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
	// 先级联关闭子资源（此时 Context 仍存活），再失效令牌并释放底层上下文。
	// 调试层报告未显式 Close 的子资源（C1：资源审计）。
	if checks.Debug && len(closers) > 0 {
		logf("context closed with %d child resource(s) not explicitly closed; cascading", len(closers))
	}
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

// SyncScopeID 按名查询 sync scope ID。
//
// 常见作用域名（如 "system"、"singlethread"）由 LLVM 固定识别，但对应 ID 的具体数值
// 属 LLVM 内部约定，不同版本可能不同，勿硬编码；未注册的名字由 LLVM 分配新 ID。
func (ctx *Context) SyncScopeID(name string) uint32 {
	ctx.CheckAlive("llvm.Context.SyncScopeID")
	return binding.LLVMGetSyncScopeID(ctx.ref, name)
}

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

// ===== Go 类型映射缓存（per-Context，随 Context 释放） =====

// lookupGoType 查询 Go 类型映射缓存
func (ctx *Context) lookupGoType(t reflect.Type) (AnyType, bool) {
	ctx.cacheMu.Lock()
	defer ctx.cacheMu.Unlock()
	ty, ok := ctx.typeCache[t]
	return ty, ok
}

// storeGoType 写入 Go 类型映射缓存
func (ctx *Context) storeGoType(t reflect.Type, ty AnyType) {
	ctx.cacheMu.Lock()
	defer ctx.cacheMu.Unlock()
	if ctx.typeCache == nil {
		ctx.typeCache = make(map[reflect.Type]AnyType)
	}
	ctx.typeCache[t] = ty
}

// lookupFnSig 查询 Go 函数签名映射缓存
func (ctx *Context) lookupFnSig(t reflect.Type) (FnType, bool) {
	ctx.cacheMu.Lock()
	defer ctx.cacheMu.Unlock()
	sig, ok := ctx.fnCache[t]
	return sig, ok
}

// storeFnSig 写入 Go 函数签名映射缓存
func (ctx *Context) storeFnSig(t reflect.Type, sig FnType) {
	ctx.cacheMu.Lock()
	defer ctx.cacheMu.Unlock()
	if ctx.fnCache == nil {
		ctx.fnCache = make(map[reflect.Type]FnType)
	}
	ctx.fnCache[t] = sig
}
