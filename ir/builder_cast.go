package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Trunc 整数截断
func (b *Builder) Trunc(v llvm.ValueRef[llvm.IntT], to llvm.IntType, name string) llvm.Value[llvm.IntT] {
	const op = "ir.Builder.Trunc"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	preCastTo(op, b, to)
	return wrapValue[llvm.IntT](b.ctx, b.inserted.life, binding.LLVMBuildTrunc(b.ref, vv.Ref(), to.Ref(), name))
}

// ZExt 整数零扩展
func (b *Builder) ZExt(v llvm.ValueRef[llvm.IntT], to llvm.IntType, name string) llvm.Value[llvm.IntT] {
	const op = "ir.Builder.ZExt"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	preCastTo(op, b, to)
	return wrapValue[llvm.IntT](b.ctx, b.inserted.life, binding.LLVMBuildZExt(b.ref, vv.Ref(), to.Ref(), name))
}

// SExt 整数符号扩展
func (b *Builder) SExt(v llvm.ValueRef[llvm.IntT], to llvm.IntType, name string) llvm.Value[llvm.IntT] {
	const op = "ir.Builder.SExt"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	preCastTo(op, b, to)
	return wrapValue[llvm.IntT](b.ctx, b.inserted.life, binding.LLVMBuildSExt(b.ref, vv.Ref(), to.Ref(), name))
}

// FPTrunc 浮点精度截断
func (b *Builder) FPTrunc(v llvm.ValueRef[llvm.FloatT], to llvm.FloatType, name string) llvm.Value[llvm.FloatT] {
	const op = "ir.Builder.FPTrunc"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	preCastTo(op, b, to)
	return wrapValue[llvm.FloatT](b.ctx, b.inserted.life, binding.LLVMBuildFPTrunc(b.ref, vv.Ref(), to.Ref(), name))
}

// FPExt 浮点精度扩展
func (b *Builder) FPExt(v llvm.ValueRef[llvm.FloatT], to llvm.FloatType, name string) llvm.Value[llvm.FloatT] {
	const op = "ir.Builder.FPExt"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	preCastTo(op, b, to)
	return wrapValue[llvm.FloatT](b.ctx, b.inserted.life, binding.LLVMBuildFPExt(b.ref, vv.Ref(), to.Ref(), name))
}

// FPToUI 浮点转无符号整数
func (b *Builder) FPToUI(v llvm.ValueRef[llvm.FloatT], to llvm.IntType, name string) llvm.Value[llvm.IntT] {
	const op = "ir.Builder.FPToUI"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	preCastTo(op, b, to)
	return wrapValue[llvm.IntT](b.ctx, b.inserted.life, binding.LLVMBuildFPToUI(b.ref, vv.Ref(), to.Ref(), name))
}

// FPToSI 浮点转有符号整数
func (b *Builder) FPToSI(v llvm.ValueRef[llvm.FloatT], to llvm.IntType, name string) llvm.Value[llvm.IntT] {
	const op = "ir.Builder.FPToSI"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	preCastTo(op, b, to)
	return wrapValue[llvm.IntT](b.ctx, b.inserted.life, binding.LLVMBuildFPToSI(b.ref, vv.Ref(), to.Ref(), name))
}

// UIToFP 无符号整数转浮点
func (b *Builder) UIToFP(v llvm.ValueRef[llvm.IntT], to llvm.FloatType, name string) llvm.Value[llvm.FloatT] {
	const op = "ir.Builder.UIToFP"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	preCastTo(op, b, to)
	return wrapValue[llvm.FloatT](b.ctx, b.inserted.life, binding.LLVMBuildUIToFP(b.ref, vv.Ref(), to.Ref(), name))
}

// SIToFP 有符号整数转浮点
func (b *Builder) SIToFP(v llvm.ValueRef[llvm.IntT], to llvm.FloatType, name string) llvm.Value[llvm.FloatT] {
	const op = "ir.Builder.SIToFP"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	preCastTo(op, b, to)
	return wrapValue[llvm.FloatT](b.ctx, b.inserted.life, binding.LLVMBuildSIToFP(b.ref, vv.Ref(), to.Ref(), name))
}

// PtrToInt 指针转整数
func (b *Builder) PtrToInt(v llvm.ValueRef[llvm.PtrT], to llvm.IntType, name string) llvm.Value[llvm.IntT] {
	const op = "ir.Builder.PtrToInt"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	preCastTo(op, b, to)
	return wrapValue[llvm.IntT](b.ctx, b.inserted.life, binding.LLVMBuildPtrToInt(b.ref, vv.Ref(), to.Ref(), name))
}

// IntToPtr 整数转指针
func (b *Builder) IntToPtr(v llvm.ValueRef[llvm.IntT], to llvm.PtrType, name string) llvm.Value[llvm.PtrT] {
	const op = "ir.Builder.IntToPtr"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	preCastTo(op, b, to)
	return wrapValue[llvm.PtrT](b.ctx, b.inserted.life, binding.LLVMBuildIntToPtr(b.ref, vv.Ref(), to.Ref(), name))
}

// BitCast 位重解释转换（泛型方法：结果种类由调用方断言）
func (b *Builder) BitCast[U llvm.Kind](v llvm.AnyValue, to llvm.TypeRef[U], name string) llvm.Value[U] {
	const op = "ir.Builder.BitCast"
	tt := to.AsType()
	b.pre(op, v)
	preCastTo(op, b, tt)
	return wrapValue[U](b.ctx, b.inserted.life, binding.LLVMBuildBitCast(b.ref, v.Ref(), tt.Ref(), name))
}

// preCastTo 预检目标类型归属同一 Context
func preCastTo[T llvm.Kind](op string, b *Builder, to llvm.TypeRef[T]) {
	if to.AsType().Context() != b.ctx {
		errPanic(llvm.ErrCrossContext, op, "target type belongs to another context")
	}
}
