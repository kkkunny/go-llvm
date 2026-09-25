package jit

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
	"github.com/kkkunny/go-llvm/ir"
)

// ResourceTracker ORC 资源跟踪器：跟踪一批已加入 JIT 的符号，支持整批卸载。
// 持有 ORC 引用计数；Close 释放引用（不卸载），Remove 卸载全部已跟踪符号并释放引用。
// 并发规则同 LLJIT：多 goroutine 可用，Close/Remove 由调用方串行化。
type ResourceTracker struct {
	ref    binding.LLVMOrcResourceTrackerRef
	j      *LLJIT
	closed bool
}

// NewResourceTracker 在主 JITDylib 上创建资源跟踪器
func (j *LLJIT) NewResourceTracker() ResourceTracker {
	const op = "jit.LLJIT.NewResourceTracker"
	j.check(op)
	jd := binding.LLVMOrcLLJITGetMainJITDylib(j.ref)
	return ResourceTracker{ref: binding.LLVMOrcJITDylibCreateResourceTracker(jd), j: j}
}

// ClearSymbols 卸载主 JITDylib 的全部符号定义（等价逐个 tracker remove，含默认 tracker）；
// 之后可继续添加模块。对应 JITDylib::clear()。
func (j *LLJIT) ClearSymbols() error {
	const op = "jit.LLJIT.ClearSymbols"
	j.check(op)
	if err := binding.LLVMOrcJITDylibClear(binding.LLVMOrcLLJITGetMainJITDylib(j.ref)); err != nil {
		return llvm.WrapError(llvm.ErrJIT, op, err)
	}
	return nil
}

// check 前置校验（崩溃类地板）
func (rt *ResourceTracker) check(op string) {
	if rt == nil || rt.ref.IsNil() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil resource tracker")
	}
	if rt.closed {
		llvm.Panicf(llvm.ErrUseAfterFree, op, "resource tracker is closed")
	}
	rt.j.check(op)
}

// AddIRModule 把模块加入 JIT 并纳入本跟踪器；所有权移交（模块与其 Context 的 Go 句柄即刻失效）
func (rt *ResourceTracker) AddIRModule(mod *ir.Module) error {
	const op = "jit.ResourceTracker.AddIRModule"
	rt.check(op)
	mod.Check(op)
	if checks.Debug {
		rt.j.recordModuleSigs(mod)
		if err := mod.Verify(); err != nil {
			llvm.Panicf(llvm.ErrVerify, op, "module verification failed before JIT: %s", err)
		}
	}
	tsm := consumeModule(mod)
	if err := binding.LLVMOrcLLJITAddLLVMIRModuleWithRT(rt.j.ref, rt.ref, tsm); err != nil {
		return llvm.WrapError(llvm.ErrJIT, op, err)
	}
	return nil
}

// AddObjectFile 把目标文件缓冲纳入本跟踪器；所有权移交
func (rt *ResourceTracker) AddObjectFile(buf *llvm.MemoryBuffer) error {
	const op = "jit.ResourceTracker.AddObjectFile"
	rt.check(op)
	if buf == nil || !buf.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or closed memory buffer")
	}
	b := buf.Ref()
	buf.Disown()
	if err := binding.LLVMOrcLLJITAddObjectFileWithRT(rt.j.ref, rt.ref, b); err != nil {
		return llvm.WrapError(llvm.ErrJIT, op, err)
	}
	return nil
}

// Remove 卸载全部已跟踪符号并释放引用；之后本跟踪器不可再用
func (rt *ResourceTracker) Remove() error {
	const op = "jit.ResourceTracker.Remove"
	rt.check(op)
	err := binding.LLVMOrcResourceTrackerRemove(rt.ref)
	rt.closed = true
	binding.LLVMOrcReleaseResourceTracker(rt.ref)
	if err != nil {
		return llvm.WrapError(llvm.ErrJIT, op, err)
	}
	return nil
}

// TransferTo 把全部已跟踪符号移交给另一跟踪器（C API 无错误返回）
func (rt *ResourceTracker) TransferTo(dst ResourceTracker) {
	const op = "jit.ResourceTracker.TransferTo"
	rt.check(op)
	dst.check(op)
	binding.LLVMOrcResourceTrackerTransferTo(rt.ref, dst.ref)
}

// Close 释放引用（不卸载）；二次调用返回 ErrClosed
func (rt *ResourceTracker) Close() error {
	if rt.closed {
		return &llvm.Error{Reason: llvm.ErrClosed, Op: "jit.ResourceTracker.Close", Msg: "resource tracker already closed"}
	}
	rt.closed = true
	binding.LLVMOrcReleaseResourceTracker(rt.ref)
	return nil
}
