package jit

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
	"github.com/kkkunny/go-llvm/ir"
)

// LLJIT ORC 延迟编译 JIT 实例；独立所有权根（不经 Context.Own），用毕 Close。
// 移交给它的模块与其 Context 由 LLJIT 负责释放。
// 并发：Func/MapFunc/Lookup 等方法的适配器缓存与注册表受锁保护；
// Close 与其他方法之间的并发调用仍由调用方自行串行化。
type LLJIT struct {
	ref         binding.LLVMOrcLLJITRef
	closed      bool
	mu          sync.Mutex // 保护 adapters 缓存
	adapters    map[reflect.Type]*adapterEntry
	channelOnce sync.Once
	channelErr  error
}

// NewLLJIT 创建面向宿主的目标 JIT；失败返回 ErrJIT
func NewLLJIT() (*LLJIT, error) {
	const op = "jit.NewLLJIT"
	jtmb, err := binding.LLVMOrcJITTargetMachineBuilderDetectHost()
	if err != nil {
		return nil, llvm.WrapError(llvm.ErrJIT, op, err)
	}
	builder := binding.LLVMOrcCreateLLJITBuilder()
	binding.LLVMOrcLLJITBuilderSetJITTargetMachineBuilder(builder, jtmb)
	ref, err := binding.LLVMOrcCreateLLJIT(builder)
	if err != nil {
		return nil, llvm.WrapError(llvm.ErrJIT, op, err)
	}
	return &LLJIT{ref: ref}, nil
}

// check 前置校验：未释放
func (j *LLJIT) check(op string) {
	if j.closed {
		llvm.Panicf(llvm.ErrUseAfterFree, op, "LLJIT is closed")
	}
}

// Close 释放 JIT 及其托管的所有模块/上下文；二次调用返回 ErrClosed
func (j *LLJIT) Close() error {
	if j.closed {
		return &llvm.Error{Reason: llvm.ErrClosed, Op: "jit.LLJIT.Close", Msg: "LLJIT already closed"}
	}
	j.closed = true
	err := binding.LLVMOrcDisposeLLJIT(j.ref)
	if err != nil {
		return llvm.WrapError(llvm.ErrJIT, "jit.LLJIT.Close", err)
	}
	return nil
}

// Triple 目标三元组
func (j *LLJIT) Triple() string {
	j.check("jit.LLJIT.Triple")
	return binding.LLVMOrcLLJITGetTripleString(j.ref)
}

// DataLayoutStr 默认数据布局字符串
func (j *LLJIT) DataLayoutStr() string {
	j.check("jit.LLJIT.DataLayoutStr")
	return binding.LLVMOrcLLJITGetDataLayoutStr(j.ref)
}

// AddIRModule 把模块加入主 JITDylib；失败返回 ErrJIT。
// 所有权移交：模块及其整个 Context 归 JIT 所有，Go 侧句柄立即失效（此后使用会 panic）。
func (j *LLJIT) AddIRModule(mod *ir.Module) error {
	const op = "jit.LLJIT.AddIRModule"
	j.check(op)
	mod.Check(op)
	if checks.Debug {
		if err := mod.Verify(); err != nil {
			llvm.Panicf(llvm.ErrVerify, op, "module verification failed before JIT: %s", err)
		}
	}

	ctx := mod.Context()
	tsctx := binding.LLVMOrcCreateNewThreadSafeContextFromLLVMContext(ctx.Ref())
	mod.Disown()
	ctx.Disown()
	tsm := binding.LLVMOrcCreateNewThreadSafeModule(mod.Ref(), tsctx)
	binding.LLVMOrcDisposeThreadSafeContext(tsctx)
	// ORC 约定：调用后所有权无条件移交（失败时由 JIT 错误路径释放 TSM）
	if err := binding.LLVMOrcLLJITAddLLVMIRModule(j.ref, binding.LLVMOrcLLJITGetMainJITDylib(j.ref), tsm); err != nil {
		return llvm.WrapError(llvm.ErrJIT, op, err)
	}
	return nil
}

// AddObjectFile 把目标文件缓冲加入主 JITDylib；失败返回 ErrJIT。
// 所有权移交：调用后缓冲归 JIT 所有（失败时由 JIT 释放），Go 侧句柄失效。
func (j *LLJIT) AddObjectFile(buf *llvm.MemoryBuffer) error {
	const op = "jit.LLJIT.AddObjectFile"
	j.check(op)
	if buf == nil || !buf.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or closed memory buffer")
	}
	buf.Disown()
	if err := binding.LLVMOrcLLJITAddObjectFile(j.ref, binding.LLVMOrcLLJITGetMainJITDylib(j.ref), buf.Ref()); err != nil {
		return llvm.WrapError(llvm.ErrJIT, op, err)
	}
	return nil
}

// Lookup 在主 JITDylib 中查找符号（触发按需编译）；未找到返回 ErrNotFound
func (j *LLJIT) Lookup(name string) (unsafe.Pointer, error) {
	const op = "jit.LLJIT.Lookup"
	j.check(op)
	addr, err := binding.LLVMOrcLLJITLookup(j.ref, name)
	if err != nil {
		return nil, llvm.WrapError(llvm.ErrNotFound, op, err)
	}
	return addr, nil
}

// MapSymbol 把宿主地址定义为 JIT 符号（绝对符号）；失败返回 ErrJIT
func (j *LLJIT) MapSymbol(name string, p unsafe.Pointer) error {
	const op = "jit.LLJIT.MapSymbol"
	j.check(op)
	if p == nil {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil symbol address")
	}
	entry := binding.LLVMOrcLLJITMangleAndIntern(j.ref, name)
	mu := binding.LLVMOrcAbsoluteSymbols([]binding.LLVMOrcCSymbolMapPair{{
		Name: entry,
		Sym: binding.LLVMJITEvaluatedSymbol{
			Address: uint64(uintptr(p)),
			Flags: binding.LLVMJITSymbolFlags{
				GenericFlags: binding.LLVMJITSymbolGenericFlagsExported | binding.LLVMJITSymbolGenericFlagsCallable,
			},
		},
	}})
	if err := binding.LLVMOrcJITDylibDefine(binding.LLVMOrcLLJITGetMainJITDylib(j.ref), mu); err != nil {
		binding.LLVMOrcDisposeMaterializationUnit(mu)
		return llvm.WrapError(llvm.ErrJIT, op, err)
	}
	return nil
}
