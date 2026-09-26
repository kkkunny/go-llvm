package llvm

import (
	"github.com/kkkunny/go-llvm/internal/binding"
)

// ConstExtractElement 常量向量取元素（结果种类随元素类型，运行时才知道，返回 Dyn）
func (ctx *Context) ConstExtractElement(v ValueRef[VecT], idx ValueRef[IntT]) Value[DynT] {
	const op = "llvm.Context.ConstExtractElement"
	ctx.CheckAlive(op)
	vv, iv := v.AsValue(), idx.AsValue()
	ctx.CheckValues(op, vv, iv)
	return ValueOf(ctx, ctx.life, binding.LLVMConstExtractElement(vv.Ref(), iv.Ref()))
}

// ConstInsertElement 常量向量写元素（elem 类型须与向量元素一致，不符则由 LLVM 报错）
func (ctx *Context) ConstInsertElement(v ValueRef[VecT], elem AnyValue, idx ValueRef[IntT]) Value[VecT] {
	const op = "llvm.Context.ConstInsertElement"
	ctx.CheckAlive(op)
	vv, iv := v.AsValue(), idx.AsValue()
	ctx.CheckValues(op, vv, elem, iv)
	return newValue[VecT](ctx, ctx.life, binding.LLVMConstInsertElement(vv.Ref(), elem.Ref(), iv.Ref()))
}

// ConstShuffleVector 常量向量洗牌；mask 的元素为 32 位整数常量
func (ctx *Context) ConstShuffleVector(a, b ValueRef[VecT], mask ValueRef[VecT]) Value[VecT] {
	const op = "llvm.Context.ConstShuffleVector"
	ctx.CheckAlive(op)
	av, bv, mv := a.AsValue(), b.AsValue(), mask.AsValue()
	ctx.CheckValues(op, av, bv, mv)
	return newValue[VecT](ctx, ctx.life, binding.LLVMConstShuffleVector(av.Ref(), bv.Ref(), mv.Ref()))
}
