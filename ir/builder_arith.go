package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// ===== 整数算术/位运算 =====

// Add 插入 add（回绕语义）
func (b *Builder) Add(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.Add", l, r, name, binding.LLVMBuildAdd)
}

// AddNSW 插入 add nsw
func (b *Builder) AddNSW(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.AddNSW", l, r, name, binding.LLVMBuildNSWAdd)
}

// AddNUW 插入 add nuw
func (b *Builder) AddNUW(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.AddNUW", l, r, name, binding.LLVMBuildNUWAdd)
}

// Sub 插入 sub
func (b *Builder) Sub(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.Sub", l, r, name, binding.LLVMBuildSub)
}

// SubNSW 插入 sub nsw
func (b *Builder) SubNSW(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.SubNSW", l, r, name, binding.LLVMBuildNSWSub)
}

// SubNUW 插入 sub nuw
func (b *Builder) SubNUW(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.SubNUW", l, r, name, binding.LLVMBuildNUWSub)
}

// Mul 插入 mul
func (b *Builder) Mul(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.Mul", l, r, name, binding.LLVMBuildMul)
}

// MulNSW 插入 mul nsw
func (b *Builder) MulNSW(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.MulNSW", l, r, name, binding.LLVMBuildNSWMul)
}

// MulNUW 插入 mul nuw
func (b *Builder) MulNUW(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.MulNUW", l, r, name, binding.LLVMBuildNUWMul)
}

// SDiv 插入有符号除法
func (b *Builder) SDiv(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.SDiv", l, r, name, binding.LLVMBuildSDiv)
}

// UDiv 插入无符号除法
func (b *Builder) UDiv(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.UDiv", l, r, name, binding.LLVMBuildUDiv)
}

// SRem 插入有符号取余
func (b *Builder) SRem(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.SRem", l, r, name, binding.LLVMBuildSRem)
}

// URem 插入无符号取余
func (b *Builder) URem(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.URem", l, r, name, binding.LLVMBuildURem)
}

// Shl 插入左移
func (b *Builder) Shl(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.Shl", l, r, name, binding.LLVMBuildShl)
}

// LShr 插入逻辑右移
func (b *Builder) LShr(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.LShr", l, r, name, binding.LLVMBuildLShr)
}

// AShr 插入算术右移
func (b *Builder) AShr(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.AShr", l, r, name, binding.LLVMBuildAShr)
}

// And 插入按位与
func (b *Builder) And(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.And", l, r, name, binding.LLVMBuildAnd)
}

// Or 插入按位或
func (b *Builder) Or(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.Or", l, r, name, binding.LLVMBuildOr)
}

// Xor 插入按位异或
func (b *Builder) Xor(l, r llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	return b.intBinop("ir.Builder.Xor", l, r, name, binding.LLVMBuildXor)
}

// Neg 插入整数取负
func (b *Builder) Neg(v llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	const op = "ir.Builder.Neg"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	return wrapValue[llvm.IntT](b.ctx, b.inserted.life, binding.LLVMBuildNeg(b.ref, vv.Ref(), name))
}

// Not 插入按位取反
func (b *Builder) Not(v llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.IntT] {
	const op = "ir.Builder.Not"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	return wrapValue[llvm.IntT](b.ctx, b.inserted.life, binding.LLVMBuildNot(b.ref, vv.Ref(), name))
}

// intBinop 整数二元指令公共实现
func (b *Builder) intBinop(op string, l, r llvm.ValueRef[llvm.IntT], name string, build func(binding.LLVMBuilderRef, binding.LLVMValueRef, binding.LLVMValueRef, string) binding.LLVMValueRef) llvm.Value[llvm.IntT] {
	lv, rv := l.AsValue(), r.AsValue()
	b.preSameType(op, lv.Dyn(), rv.Dyn())
	return wrapValue[llvm.IntT](b.ctx, b.inserted.life, build(b.ref, lv.Ref(), rv.Ref(), name))
}

// ===== 浮点算术 =====

// FAdd 插入浮点加法
func (b *Builder) FAdd(l, r llvm.ValueRef[llvm.FloatT], name string) llvm.Value[llvm.FloatT] {
	return b.floatBinop("ir.Builder.FAdd", l, r, name, binding.LLVMBuildFAdd)
}

// FSub 插入浮点减法
func (b *Builder) FSub(l, r llvm.ValueRef[llvm.FloatT], name string) llvm.Value[llvm.FloatT] {
	return b.floatBinop("ir.Builder.FSub", l, r, name, binding.LLVMBuildFSub)
}

// FMul 插入浮点乘法
func (b *Builder) FMul(l, r llvm.ValueRef[llvm.FloatT], name string) llvm.Value[llvm.FloatT] {
	return b.floatBinop("ir.Builder.FMul", l, r, name, binding.LLVMBuildFMul)
}

// FDiv 插入浮点除法
func (b *Builder) FDiv(l, r llvm.ValueRef[llvm.FloatT], name string) llvm.Value[llvm.FloatT] {
	return b.floatBinop("ir.Builder.FDiv", l, r, name, binding.LLVMBuildFDiv)
}

// FRem 插入浮点取余
func (b *Builder) FRem(l, r llvm.ValueRef[llvm.FloatT], name string) llvm.Value[llvm.FloatT] {
	return b.floatBinop("ir.Builder.FRem", l, r, name, binding.LLVMBuildFRem)
}

// FNeg 插入浮点取负
func (b *Builder) FNeg(v llvm.ValueRef[llvm.FloatT], name string) llvm.Value[llvm.FloatT] {
	const op = "ir.Builder.FNeg"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	return wrapValue[llvm.FloatT](b.ctx, b.inserted.life, binding.LLVMBuildFNeg(b.ref, vv.Ref(), name))
}

// floatBinop 浮点二元指令公共实现
func (b *Builder) floatBinop(op string, l, r llvm.ValueRef[llvm.FloatT], name string, build func(binding.LLVMBuilderRef, binding.LLVMValueRef, binding.LLVMValueRef, string) binding.LLVMValueRef) llvm.Value[llvm.FloatT] {
	lv, rv := l.AsValue(), r.AsValue()
	b.preSameType(op, lv.Dyn(), rv.Dyn())
	return wrapValue[llvm.FloatT](b.ctx, b.inserted.life, build(b.ref, lv.Ref(), rv.Ref(), name))
}

// ===== 比较与选择 =====

// ICmp 插入整数/指针比较，返回 i1
func (b *Builder) ICmp(pred llvm.IntPred, l, r llvm.AnyValue, name string) llvm.Value[llvm.IntT] {
	const op = "ir.Builder.ICmp"
	b.preSameType(op, l, r)
	ref := binding.LLVMBuildICmp(b.ref, binding.LLVMIntPredicate(pred), l.Ref(), r.Ref(), name)
	return wrapValue[llvm.IntT](b.ctx, b.inserted.life, ref)
}

// FCmp 插入浮点比较，返回 i1
func (b *Builder) FCmp(pred llvm.FloatPred, l, r llvm.ValueRef[llvm.FloatT], name string) llvm.Value[llvm.IntT] {
	const op = "ir.Builder.FCmp"
	lv, rv := l.AsValue(), r.AsValue()
	b.preSameType(op, lv.Dyn(), rv.Dyn())
	ref := binding.LLVMBuildFCmp(b.ref, binding.LLVMRealPredicate(pred), lv.Ref(), rv.Ref(), name)
	return wrapValue[llvm.IntT](b.ctx, b.inserted.life, ref)
}

// Select 插入 select；x/y 种类一致由编译期保证
func (b *Builder) Select[T llvm.Kind](cond llvm.ValueRef[llvm.IntT], x, y llvm.ValueRef[T], name string) llvm.Value[T] {
	const op = "ir.Builder.Select"
	cv, xv, yv := cond.AsValue(), x.AsValue(), y.AsValue()
	b.pre(op, cv.Dyn(), xv.Dyn(), yv.Dyn())
	b.preSameType(op, xv.Dyn(), yv.Dyn())
	if bits := llvm.AsIntType(cv.Type()).Bits(); bits != 1 {
		llvm.Panicf(llvm.ErrTypeMismatch, op, "condition must be i1, got i%d", bits)
	}
	ref := binding.LLVMBuildSelect(b.ref, cv.Ref(), xv.Ref(), yv.Ref(), name)
	return wrapValue[T](b.ctx, b.inserted.life, ref)
}
