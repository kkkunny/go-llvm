package llvm

import (
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/errs"
)

// ConstExtractElement 常量向量取元素（结果种类随元素类型，运行时才知道，返回 Dyn）
func (ctx *Context) ConstExtractElement(v ValueRef[VecT], idx ValueRef[IntT]) Value[DynT] {
	const op = "llvm.Context.ConstExtractElement"
	ctx.CheckAlive(op)
	vv, iv := v.AsValue(), idx.AsValue()
	ctx.CheckValues(op, vv, iv)
	checkConstOperand(op, -1, vv)
	checkConstOperand(op, -1, iv)
	return ValueOf(ctx, ctx.life, binding.LLVMConstExtractElement(vv.Ref(), iv.Ref()))
}

// ConstInsertElement 常量向量写元素；elem 类型须与向量元素一致（不符或非常量 panic）
func (ctx *Context) ConstInsertElement(v ValueRef[VecT], elem AnyValue, idx ValueRef[IntT]) Value[VecT] {
	const op = "llvm.Context.ConstInsertElement"
	ctx.CheckAlive(op)
	vv, iv := v.AsValue(), idx.AsValue()
	ctx.CheckValues(op, vv, elem, iv)
	checkConstOperand(op, -1, vv)
	checkConstOperand(op, -1, elem)
	checkConstOperand(op, -1, iv)
	if want := MustVecType(vv.Type()).Elem(); !elem.Dyn().Type().Equal(want) {
		errs.Panicf(ErrTypeMismatch, op, "element type %s does not match vector element type %s", elem.Dyn().Type(), want)
	}
	return newValue[VecT](ctx, ctx.life, binding.LLVMConstInsertElement(vv.Ref(), elem.Ref(), iv.Ref()))
}

// ConstShuffleVector 常量向量洗牌；mask 为 32 位整数常量向量，a/b 类型须一致且均为常量
func (ctx *Context) ConstShuffleVector(a, b ValueRef[VecT], mask ValueRef[VecT]) Value[VecT] {
	const op = "llvm.Context.ConstShuffleVector"
	ctx.CheckAlive(op)
	av, bv, mv := a.AsValue(), b.AsValue(), mask.AsValue()
	ctx.CheckValues(op, av, bv, mv)
	checkConstOperand(op, -1, av)
	checkConstOperand(op, -1, bv)
	checkConstOperand(op, -1, mv)
	if !av.Type().Equal(bv.Type()) {
		errs.Panicf(ErrTypeMismatch, op, "vector types differ: %s vs %s", av.Type(), bv.Type())
	}
	if elem := MustVecType(mv.Type()).Elem(); !elem.Equal(ctx.Int(32)) {
		errs.Panicf(ErrTypeMismatch, op, "shuffle mask element type must be i32, got %s", elem)
	}
	return newValue[VecT](ctx, ctx.life, binding.LLVMConstShuffleVector(av.Ref(), bv.Ref(), mv.Ref()))
}
