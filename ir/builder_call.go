package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
)

// Call 插入函数调用；返回种类 U 在调用前与函数返回类型比对，不符 panic（语义契约，仅调试层）
func (b *Builder) Call[U llvm.Kind](fn llvm.ValueRef[llvm.FnT], args []llvm.AnyValue, name string) Call[U] {
	const op = "ir.Builder.Call"
	fv := fn.AsValue()
	b.pre(op, core(fv))
	sig := callSig(b.ctx, fv.Ref(), fv.RawType())
	if checks.Debug {
		checkKind[U](op, b.ctx, binding.LLVMGetReturnType(sig.Ref()))
	}
	ref := b.call(op, fv.Ref(), sig, args, name)
	return Call[U]{Value: llvm.NewValue[U](b.ctx, b.inserted.life, ref)}
}

// callSig 调用目标的函数类型：Function 走 LLVMGetFunctionType，
// InlineAsm 走 LLVMGetInlineAsmFunctionType（其值为 ptr 类型，不能反推签名）。
// calleeTy 为调用方缓存的类型句柄（可为空，空则查询）。
func callSig(ctx *llvm.Context, callee binding.LLVMValueRef, calleeTy binding.LLVMTypeRef) llvm.FnType {
	ty := calleeTy
	if ty.IsNil() {
		ty = binding.LLVMTypeOf(callee)
	}
	if binding.LLVMGetTypeKind(ty) == binding.LLVMFunctionTypeKind {
		return llvm.AsFnType(llvm.TypeOfRef(ctx, ty))
	}
	if binding.LLVMGetValueKind(callee) == binding.LLVMInlineAsmValueKind {
		return llvm.AsFnType(llvm.TypeOfRef(ctx, binding.LLVMGetInlineAsmFunctionType(callee)))
	}
	return llvm.AsFnType(llvm.TypeOfRef(ctx, binding.LLVMGetFunctionType(callee)))
}

// CallIndirect 通过函数指针调用（不透明指针 + 签名）；返回种类 U 在调用前与签名返回类型比对，不符 panic（仅调试层）
func (b *Builder) CallIndirect[U llvm.Kind](fnPtr llvm.ValueRef[llvm.PtrT], sig llvm.FnType, args []llvm.AnyValue, name string) Call[U] {
	const op = "ir.Builder.CallIndirect"
	pv := fnPtr.AsValue()
	b.pre(op, core(pv))
	if sig.Context() != b.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "signature belongs to another context")
	}
	if checks.Debug {
		checkKind[U](op, b.ctx, binding.LLVMGetReturnType(sig.Ref()))
	}
	ref := b.call(op, pv.Ref(), sig, args, name)
	return Call[U]{Value: llvm.NewValue[U](b.ctx, b.inserted.life, ref)}
}

// call 调用公共路径：实参预检 + 参数个数/类型校验后发指令（callee 已由调用方校验）
func (b *Builder) call(op string, callee binding.LLVMValueRef, sig llvm.FnType, args []llvm.AnyValue, name string) binding.LLVMValueRef {
	b.pre(op)
	b.checkCallArgs(op, sig, args)
	return binding.LLVMBuildCall(b.ref, sig.Ref(), callee, b.valueRefs(args), name)
}

// checkCallArgs 调用类指令（call/invoke）公共实参预检：崩溃类地板常开；个数/类型为语义契约，仅调试层
func (b *Builder) checkCallArgs(op string, sig llvm.FnType, args []llvm.AnyValue) {
	for _, a := range args {
		b.checkVal(op, coreAny(a))
	}
	if !checks.Debug {
		return
	}
	n := binding.LLVMCountParamTypes(sig.Ref())
	got := uint(len(args))
	if got < uint(n) {
		llvm.Panicf(llvm.ErrTypeMismatch, op, "expect at least %d arguments, got %d", n, got)
	}
	if got > uint(n) && !sig.IsVarArg() {
		llvm.Panicf(llvm.ErrTypeMismatch, op, "expect %d arguments, got %d", n, got)
	}
	if n > 0 {
		params := binding.LLVMGetParamTypes(sig.Ref())
		for i, p := range params {
			if !p.Equal(anyValType(args[i])) {
				llvm.Panicf(llvm.ErrTypeMismatch, op, "argument %d type %s does not match parameter type %s",
					i, typeString(b.ctx, args[i].Ref()), typeRefString(b.ctx, p))
			}
		}
	}
}

// checkKind 校验底层类型句柄与种类参数 U 匹配；不符 panic
func checkKind[U llvm.Kind](op string, ctx *llvm.Context, ref binding.LLVMTypeRef) {
	if _, err := llvm.TypeOfRef(ctx, ref).As[U](); err != nil {
		msg := err.Error()
		if e, ok := err.(*llvm.Error); ok {
			msg = e.Msg
		}
		llvm.Panicf(llvm.ErrTypeMismatch, op, "%s", msg)
	}
}

// elementTypeAt 沿 indices 路径求聚合值的元素类型；越界或不可索引 panic
func elementTypeAt(op string, ctx *llvm.Context, agg binding.LLVMTypeRef, indices []uint32) binding.LLVMTypeRef {
	cur := agg
	for depth, idx := range indices {
		switch binding.LLVMGetTypeKind(cur) {
		case binding.LLVMStructTypeKind:
			if idx >= binding.LLVMCountStructElementTypes(cur) {
				llvm.Panicf(llvm.ErrInvalidArg, op, "index %d out of range at depth %d", idx, depth)
			}
			cur = binding.LLVMStructGetTypeAtIndex(cur, idx)
		case binding.LLVMArrayTypeKind:
			if uint64(idx) >= binding.LLVMGetArrayLength2(cur) {
				llvm.Panicf(llvm.ErrInvalidArg, op, "index %d out of range at depth %d", idx, depth)
			}
			cur = binding.LLVMGetElementType(cur)
		case binding.LLVMVectorTypeKind, binding.LLVMScalableVectorTypeKind:
			if idx >= binding.LLVMGetVectorSize(cur) {
				llvm.Panicf(llvm.ErrInvalidArg, op, "index %d out of range at depth %d", idx, depth)
			}
			cur = binding.LLVMGetElementType(cur)
		default:
			llvm.Panicf(llvm.ErrInvalidArg, op, "cannot index into %s at depth %d", typeRefString(ctx, cur), depth)
		}
	}
	return cur
}

// PHI 插入 PHI 节点
func (b *Builder) PHI[T llvm.Kind](t llvm.TypeRef[T], name string) Phi[T] {
	const op = "ir.Builder.PHI"
	tt := t.AsType()
	b.pre(op)
	if tt.Context() != b.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "type belongs to another context")
	}
	ref := binding.LLVMBuildPhi(b.ref, tt.Ref(), name)
	return Phi[T]{Value: llvm.NewValue[T](b.ctx, b.inserted.life, ref)}
}

// ExtractValue 从聚合值提取第 indices 路径的元素；结果种类 U 与元素类型比对，不符 panic（仅调试层）
func (b *Builder) ExtractValue[U llvm.Kind](agg llvm.AnyValue, indices []uint32, name string) llvm.Value[U] {
	const op = "ir.Builder.ExtractValue"
	b.pre(op, coreAny(agg))
	if len(indices) == 0 {
		llvm.Panicf(llvm.ErrInvalidArg, op, "empty index path")
	}
	if checks.Debug {
		elemTy := elementTypeAt(op, b.ctx, anyValType(agg), indices)
		checkKind[U](op, b.ctx, elemTy)
	}
	ref := binding.LLVMBuildExtractValue(b.ref, agg.Ref(), indices[0], name)
	for _, idx := range indices[1:] {
		ref = binding.LLVMBuildExtractValue(b.ref, ref, idx, name)
	}
	return llvm.NewValue[U](b.ctx, b.inserted.life, ref)
}

// InsertValue 将值插入聚合值的第 indices 路径；索引路径与元素类型校验仅调试层
func (b *Builder) InsertValue[T llvm.Kind](agg llvm.ValueRef[T], v llvm.AnyValue, indices []uint32, name string) llvm.Value[T] {
	const op = "ir.Builder.InsertValue"
	av := agg.AsValue()
	b.pre(op, core(av), coreAny(v))
	if len(indices) == 0 {
		llvm.Panicf(llvm.ErrInvalidArg, op, "empty index path")
	}
	if checks.Debug {
		elemTy := elementTypeAt(op, b.ctx, typeOfVal(av.Ref(), av.RawType()), indices)
		if !elemTy.Equal(anyValType(v)) {
			llvm.Panicf(llvm.ErrTypeMismatch, op, "value type %s does not match element type %s at path %v",
				typeString(b.ctx, v.Ref()), typeRefString(b.ctx, elemTy), indices)
		}
	}
	ref := binding.LLVMBuildInsertValue(b.ref, av.Ref(), v.Ref(), indices[0], name)
	for _, idx := range indices[1:] {
		ref = binding.LLVMBuildInsertValue(b.ref, ref, v.Ref(), idx, name)
	}
	return llvm.NewValue[T](b.ctx, b.inserted.life, ref)
}
