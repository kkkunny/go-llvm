package llvm

import (
	"github.com/kkkunny/go-llvm/internal/binding"
)

// IntConst 整数常量角色
type IntConst struct{ Value[IntT] }

// SignedValue 有符号值
func (c IntConst) SignedValue() int64 { return binding.LLVMConstIntGetSExtValue(c.ref) }

// UnsignedValue 无符号值
func (c IntConst) UnsignedValue() uint64 { return binding.LLVMConstIntGetZExtValue(c.ref) }

// IsNegative 是否有符号语义下为负
func (c IntConst) IsNegative() bool { return c.SignedValue() < 0 }

// FloatConst 浮点常量角色
type FloatConst struct{ Value[FloatT] }

// FloatValue 浮点值（方法名避开内嵌字段 Value）
func (c FloatConst) FloatValue() float64 {
	v, _ := binding.LLVMConstRealGetDouble(c.ref)
	return v
}

// ConstInt 构造整数常量
func (ctx *Context) ConstInt(t IntType, v uint64, signed bool) IntConst {
	ctx.checkType("llvm.Context.ConstInt", t)
	return IntConst{Value[IntT]{
		ref:  binding.LLVMConstInt(t.ref, v, signed),
		ctx:  ctx,
		life: ctx.life,
	}}
}

// ConstBool 构造布尔常量
func (ctx *Context) ConstBool(v bool) Value[IntT] {
	var n uint64
	if v {
		n = 1
	}
	return Value[IntT]{
		ref:  binding.LLVMConstInt(ctx.Bool().ref, n, false),
		ctx:  ctx,
		life: ctx.life,
	}
}

// ConstFloat 构造浮点常量
func (ctx *Context) ConstFloat(t FloatType, v float64) FloatConst {
	ctx.checkType("llvm.Context.ConstFloat", t)
	return FloatConst{Value[FloatT]{
		ref:  binding.LLVMConstReal(t.ref, v),
		ctx:  ctx,
		life: ctx.life,
	}}
}
