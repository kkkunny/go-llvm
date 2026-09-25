package llvm

import (
	"github.com/kkkunny/go-llvm/internal/binding"
)

// foldPre 折叠操作前置校验：同上下文、同类型。
// 与 ConstArray/ConstVector 一致不加调试开关——LLVM 对类型不符的折叠会 assert。
func (c IntConst) foldPre(op string, r ValueRef[IntT]) {
	c.Check(op)
	rv := r.AsValue()
	rv.Check(op)
	if c.Context() != rv.Context() {
		errPanic(ErrCrossContext, op, "operand belongs to another context")
	}
	if !c.Type().Equal(rv.Type()) {
		errPanic(ErrTypeMismatch, op, "operand types differ: %s vs %s", c.Type(), rv.Type())
	}
}

func (c IntConst) wrapFold(ref binding.LLVMValueRef) IntConst {
	return IntConst{newValue[IntT](c.Context(), c.Lifetime(), ref)}
}

// Add 常量加法
func (c IntConst) Add(r ValueRef[IntT]) IntConst {
	c.foldPre("llvm.IntConst.Add", r)
	return c.wrapFold(binding.LLVMConstAdd(c.Ref(), r.AsValue().Ref()))
}

// NSWAdd 常量加法（no signed wrap）
func (c IntConst) NSWAdd(r ValueRef[IntT]) IntConst {
	c.foldPre("llvm.IntConst.NSWAdd", r)
	return c.wrapFold(binding.LLVMConstNSWAdd(c.Ref(), r.AsValue().Ref()))
}

// NUWAdd 常量加法（no unsigned wrap）
func (c IntConst) NUWAdd(r ValueRef[IntT]) IntConst {
	c.foldPre("llvm.IntConst.NUWAdd", r)
	return c.wrapFold(binding.LLVMConstNUWAdd(c.Ref(), r.AsValue().Ref()))
}

// Sub 常量减法
func (c IntConst) Sub(r ValueRef[IntT]) IntConst {
	c.foldPre("llvm.IntConst.Sub", r)
	return c.wrapFold(binding.LLVMConstSub(c.Ref(), r.AsValue().Ref()))
}

// NSWSub 常量减法（no signed wrap）
func (c IntConst) NSWSub(r ValueRef[IntT]) IntConst {
	c.foldPre("llvm.IntConst.NSWSub", r)
	return c.wrapFold(binding.LLVMConstNSWSub(c.Ref(), r.AsValue().Ref()))
}

// NUWSub 常量减法（no unsigned wrap）
func (c IntConst) NUWSub(r ValueRef[IntT]) IntConst {
	c.foldPre("llvm.IntConst.NUWSub", r)
	return c.wrapFold(binding.LLVMConstNUWSub(c.Ref(), r.AsValue().Ref()))
}

// Xor 常量按位异或
func (c IntConst) Xor(r ValueRef[IntT]) IntConst {
	c.foldPre("llvm.IntConst.Xor", r)
	return c.wrapFold(binding.LLVMConstXor(c.Ref(), r.AsValue().Ref()))
}

// Neg 常量取负
func (c IntConst) Neg() IntConst {
	c.Check("llvm.IntConst.Neg")
	return c.wrapFold(binding.LLVMConstNeg(c.Ref()))
}

// NSWNeg 常量取负（no signed wrap）
func (c IntConst) NSWNeg() IntConst {
	c.Check("llvm.IntConst.NSWNeg")
	return c.wrapFold(binding.LLVMConstNSWNeg(c.Ref()))
}

// NUWNeg 常量取负（no unsigned wrap）
func (c IntConst) NUWNeg() IntConst {
	c.Check("llvm.IntConst.NUWNeg")
	return c.wrapFold(binding.LLVMConstNUWNeg(c.Ref()))
}

// Not 常量按位取反
func (c IntConst) Not() IntConst {
	c.Check("llvm.IntConst.Not")
	return c.wrapFold(binding.LLVMConstNot(c.Ref()))
}

// Cast 常量整数位宽转换（trunc 或扩展皆可，语义同 LLVMConstTruncOrBitCast）
func (c IntConst) Cast(to IntType) IntConst {
	const op = "llvm.IntConst.Cast"
	c.Check(op)
	to.Check(op)
	if c.Context() != to.Context() {
		errPanic(ErrCrossContext, op, "target type belongs to another context")
	}
	return IntConst{newValue[IntT](to.Context(), to.Context().life, binding.LLVMConstTruncOrBitCast(c.Ref(), to.Ref()))}
}

// AllOnes 该类型的全 1 常量（iN 为 -1）
func (t IntType) AllOnes() IntConst {
	t.Check("llvm.IntType.AllOnes")
	return IntConst{newValue[IntT](t.Context(), t.Context().life, binding.LLVMConstAllOnes(t.Ref()))}
}
