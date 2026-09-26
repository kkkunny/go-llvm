package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
)

// vecKindTy 预检值的底层类型确为向量，返回该类型；非向量 panic（语义契约，仅调试层调用）。
// 防止 LLVMGetElementType 收到非顺序类型（C API 不校验）
func (b *Builder) vecKindTy(op string, v preVal) binding.LLVMTypeRef {
	ref := typeOfVal(v.ref, v.ty)
	switch binding.LLVMGetTypeKind(ref) {
	case binding.LLVMVectorTypeKind, binding.LLVMScalableVectorTypeKind:
		return ref
	default:
		llvm.Panicf(llvm.ErrTypeMismatch, op, "expect a vector, got %s", typeRefString(b.ctx, ref))
		return binding.LLVMTypeRef{}
	}
}

// ExtractElement 提取向量第 idx 个元素；结果种类 U 与元素类型比对，不符 panic（仅调试层）
func (b *Builder) ExtractElement[U llvm.Kind](vec llvm.ValueRef[llvm.VecT], idx llvm.ValueRef[llvm.IntT], name string) llvm.Value[U] {
	const op = "ir.Builder.ExtractElement"
	vv, iv := vec.AsValue(), idx.AsValue()
	b.pre(op, core(vv), core(iv))
	if checks.Debug {
		elemTy := binding.LLVMGetElementType(b.vecKindTy(op, core(vv)))
		checkKind[U](op, b.ctx, elemTy)
	}
	ref := binding.LLVMBuildExtractElement(b.ref, vv.Ref(), iv.Ref(), name)
	return llvm.NewValue[U](b.ctx, b.inserted.life, ref)
}

// InsertElement 将元素插入向量第 idx 位；elem 须与向量元素类型一致（仅调试层）
func (b *Builder) InsertElement(vec llvm.ValueRef[llvm.VecT], elem llvm.AnyValue, idx llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.VecT] {
	const op = "ir.Builder.InsertElement"
	vv, iv := vec.AsValue(), idx.AsValue()
	b.pre(op, core(vv), coreAny(elem), core(iv))
	if checks.Debug {
		elemTy := binding.LLVMGetElementType(b.vecKindTy(op, core(vv)))
		if ref := anyValType(elem); !elemTy.Equal(ref) {
			llvm.Panicf(llvm.ErrTypeMismatch, op, "element type %s does not match vector element type %s",
				typeString(b.ctx, elem.Ref()), typeRefString(b.ctx, elemTy))
		}
	}
	ref := binding.LLVMBuildInsertElement(b.ref, vv.Ref(), elem.Ref(), iv.Ref(), name)
	return llvm.NewValue[llvm.VecT](b.ctx, b.inserted.life, ref)
}

// ShuffleVector 按 mask 重排两个同型向量；mask 为 <N x i32> 常量向量（校验仅调试层）
func (b *Builder) ShuffleVector(v1, v2 llvm.ValueRef[llvm.VecT], mask llvm.ValueRef[llvm.VecT], name string) llvm.Value[llvm.VecT] {
	const op = "ir.Builder.ShuffleVector"
	a, c, mv := v1.AsValue(), v2.AsValue(), mask.AsValue()
	b.preSameType(op, core(a), core(c))
	b.checkVal(op, core(mv))
	if checks.Debug {
		maskElem := binding.LLVMGetElementType(b.vecKindTy(op, core(mv)))
		i32Ref := b.ctx.Int(32).Ref()
		if binding.LLVMGetTypeKind(maskElem) != binding.LLVMIntegerTypeKind || !maskElem.Equal(i32Ref) {
			llvm.Panicf(llvm.ErrTypeMismatch, op, "mask element must be i32, got %s", typeRefString(b.ctx, maskElem))
		}
	}
	ref := binding.LLVMBuildShuffleVector(b.ref, a.Ref(), c.Ref(), mv.Ref(), name)
	return llvm.NewValue[llvm.VecT](b.ctx, b.inserted.life, ref)
}
