package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Invoke invoke 指令角色（内嵌 Value[T]）；实参操作同 Call（LLVM 侧同为 CallBase），
// 另有正常/异常出口块访问
type Invoke[T llvm.Kind] struct {
	llvm.Value[T]
}

// NormalBlock 正常出口块
func (c Invoke[T]) NormalBlock() Block {
	c.Check("ir.Invoke.NormalBlock")
	return wrapBlock(c.Context(), c.Lifetime(), binding.LLVMGetSuccessor(c.Ref(), 0))
}

// UnwindBlock 异常出口块
func (c Invoke[T]) UnwindBlock() Block {
	c.Check("ir.Invoke.UnwindBlock")
	return wrapBlock(c.Context(), c.Lifetime(), binding.LLVMGetSuccessor(c.Ref(), 1))
}

// ArgCount 实参个数
func (c Invoke[T]) ArgCount() uint32 {
	c.Check("ir.Invoke.ArgCount")
	return binding.LLVMGetNumArgOperands(c.Ref())
}

// Arg 第 i 个实参（擦除种类）
func (c Invoke[T]) Arg(i uint32) llvm.Value[llvm.DynT] {
	const op = "ir.Invoke.Arg"
	c.Check(op)
	if i >= c.ArgCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "argument index %d out of range", i)
	}
	return llvm.ValueOf(c.Context(), c.Lifetime(), binding.LLVMGetOperand(c.Ref(), i))
}

// SetArg 替换第 i 个实参
func (c Invoke[T]) SetArg(i uint32, v llvm.AnyValue) {
	const op = "ir.Invoke.SetArg"
	c.Check(op)
	if i >= c.ArgCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "argument index %d out of range", i)
	}
	c.Context().CheckValues(op, v)
	binding.LLVMSetOperand(c.Ref(), i, v.Ref())
}

// CalledFunction 被调用函数（非间接 invoke 时）
func (c Invoke[T]) CalledFunction() (llvm.Value[llvm.FnT], bool) {
	c.Check("ir.Invoke.CalledFunction")
	ref := binding.LLVMGetCalledValue(c.Ref())
	if ref.IsNil() {
		return llvm.Value[llvm.FnT]{}, false
	}
	return llvm.NewValue[llvm.FnT](c.Context(), c.Lifetime(), ref), true
}

// ===== EH 构建方法 =====

// Invoke 插入 invoke 调用：then 为正常出口、unwind 为异常出口；
// 返回种类 U 在调用前与函数返回类型比对，不符 panic
func (b *Builder) Invoke[U llvm.Kind](fn llvm.ValueRef[llvm.FnT], args []llvm.AnyValue, then, unwind Block, name string) Invoke[U] {
	const op = "ir.Builder.Invoke"
	fv := fn.AsValue()
	b.pre(op, core(fv))
	b.preBlockOwn(op, then)
	b.preBlockOwn(op, unwind)
	sig := llvm.AsFnType(llvm.TypeOfRef(b.ctx, binding.LLVMGetFunctionType(fv.Ref())))
	checkKind[U](op, b.ctx, binding.LLVMGetReturnType(sig.Ref()))
	b.checkCallArgs(op, sig, args)
	ref := binding.LLVMBuildInvoke(b.ref, sig.Ref(), fv.Ref(), b.valueRefs(args), then.ref, unwind.ref, name)
	return Invoke[U]{Value: llvm.NewValue[U](b.ctx, b.inserted.life, ref)}
}

// InvokeIndirect 通过函数指针 invoke（不透明指针 + 签名）；返回种类 U 与签名返回类型比对
func (b *Builder) InvokeIndirect[U llvm.Kind](fnPtr llvm.ValueRef[llvm.PtrT], sig llvm.FnType, args []llvm.AnyValue, then, unwind Block, name string) Invoke[U] {
	const op = "ir.Builder.InvokeIndirect"
	pv := fnPtr.AsValue()
	b.pre(op, core(pv))
	b.preBlockOwn(op, then)
	b.preBlockOwn(op, unwind)
	if sig.Context() != b.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "signature belongs to another context")
	}
	checkKind[U](op, b.ctx, binding.LLVMGetReturnType(sig.Ref()))
	b.checkCallArgs(op, sig, args)
	ref := binding.LLVMBuildInvoke(b.ref, sig.Ref(), pv.Ref(), b.valueRefs(args), then.ref, unwind.ref, name)
	return Invoke[U]{Value: llvm.NewValue[U](b.ctx, b.inserted.life, ref)}
}
