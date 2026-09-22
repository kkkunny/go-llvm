package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Call 插入函数调用；返回种类由调用方断言并在调用前预检
func (b *Builder) Call[U llvm.Kind](fn llvm.ValueRef[llvm.FnT], args []llvm.AnyValue, name string) Call[U] {
	const op = "ir.Builder.Call"
	fv := fn.AsValue()
	b.pre(op, fv.Dyn())
	sig := llvm.AsFnType(llvm.TypeOfRef(b.ctx, binding.LLVMGetFunctionType(fv.Ref())))
	ref := b.call(op, fv.Dyn(), sig, args, name)
	return Call[U]{Value: wrapValue[U](b.ctx, b.inserted.life, ref)}
}

// CallIndirect 通过函数指针调用（不透明指针 + 签名）；返回种类由调用方断言并在调用前预检
func (b *Builder) CallIndirect[U llvm.Kind](fnPtr llvm.ValueRef[llvm.PtrT], sig llvm.FnType, args []llvm.AnyValue, name string) Call[U] {
	const op = "ir.Builder.CallIndirect"
	pv := fnPtr.AsValue()
	b.pre(op, pv.Dyn())
	if sig.Context() != b.ctx {
		errPanic(llvm.ErrCrossContext, op, "signature belongs to another context")
	}
	ref := b.call(op, pv.Dyn(), sig, args, name)
	return Call[U]{Value: wrapValue[U](b.ctx, b.inserted.life, ref)}
}

// call 调用公共路径：实参个数/类型预检后发指令
func (b *Builder) call(op string, callee llvm.AnyValue, sig llvm.FnType, args []llvm.AnyValue, name string) binding.LLVMValueRef {
	b.pre(op, args...)
	params := sig.Params()
	if !sig.IsVarArg() && uint(len(args)) != uint(len(params)) {
		errPanic(llvm.ErrTypeMismatch, op, "expect %d arguments, got %d", len(params), len(args))
	}
	if uint(len(args)) < uint(len(params)) {
		errPanic(llvm.ErrTypeMismatch, op, "expect at least %d arguments, got %d", len(params), len(args))
	}
	for i, p := range params {
		if !p.Equal(args[i].Dyn().Type()) {
			errPanic(llvm.ErrTypeMismatch, op, "argument %d type %s does not match parameter type %s", i, args[i].Dyn().Type(), p)
		}
	}

	return binding.LLVMBuildCall(b.ref, sig.Ref(), callee.Ref(), anyValuesToRefs(args), name)
}

// PHI 插入 PHI 节点
func (b *Builder) PHI[T llvm.Kind](t llvm.TypeRef[T], name string) Phi[T] {
	const op = "ir.Builder.PHI"
	tt := t.AsType()
	b.pre(op)
	if tt.Context() != b.ctx {
		errPanic(llvm.ErrCrossContext, op, "type belongs to another context")
	}
	ref := binding.LLVMBuildPhi(b.ref, tt.Ref(), name)
	return Phi[T]{Value: wrapValue[T](b.ctx, b.inserted.life, ref)}
}

// ExtractValue 从聚合值提取第 indices 路径的元素；结果种类由调用方断言
func (b *Builder) ExtractValue[U llvm.Kind](agg llvm.AnyValue, indices []uint32, name string) llvm.Value[U] {
	const op = "ir.Builder.ExtractValue"
	b.pre(op, agg)
	if len(indices) == 0 {
		errPanic(llvm.ErrInvalidArg, op, "empty index path")
	}
	ref := binding.LLVMBuildExtractValue(b.ref, agg.Ref(), indices[0], name)
	for _, idx := range indices[1:] {
		ref = binding.LLVMBuildExtractValue(b.ref, ref, idx, name)
	}
	return wrapValue[U](b.ctx, b.inserted.life, ref)
}

// InsertValue 将值插入聚合值的第 indices 路径
func (b *Builder) InsertValue[T llvm.Kind](agg llvm.ValueRef[T], v llvm.AnyValue, indices []uint32, name string) llvm.Value[T] {
	const op = "ir.Builder.InsertValue"
	av := agg.AsValue()
	b.pre(op, av.Dyn(), v)
	if len(indices) == 0 {
		errPanic(llvm.ErrInvalidArg, op, "empty index path")
	}
	ref := binding.LLVMBuildInsertValue(b.ref, av.Ref(), v.Ref(), indices[0], name)
	for _, idx := range indices[1:] {
		ref = binding.LLVMBuildInsertValue(b.ref, ref, v.Ref(), idx, name)
	}
	return wrapValue[T](b.ctx, b.inserted.life, ref)
}

func anyValuesToRefs(values []llvm.AnyValue) []binding.LLVMValueRef {
	refs := make([]binding.LLVMValueRef, len(values))
	for i, v := range values {
		refs[i] = v.Ref()
	}
	return refs
}
