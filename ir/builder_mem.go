package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Alloca alloca 指令角色
type Alloca struct {
	v llvm.Value[llvm.PtrT]
}

// Value 返回底层泛型值
func (a Alloca) Value() llvm.Value[llvm.PtrT] { return a.v }

// SetAlign 设置分配对齐
func (a Alloca) SetAlign(n uint32) {
	preAlign("ir.Alloca.SetAlign", n)
	binding.LLVMSetAlignment(a.v.Ref(), n)
}

// Align 分配对齐
func (a Alloca) Align() uint32 { return binding.LLVMGetAlignment(a.v.Ref()) }

// Load 加载指令角色
type Load[T llvm.Kind] struct {
	v llvm.Value[T]
}

// Value 返回底层泛型值
func (l Load[T]) Value() llvm.Value[T] { return l.v }

// SetAlign 设置加载对齐
func (l Load[T]) SetAlign(n uint32) {
	preAlign("ir.Load.SetAlign", n)
	binding.LLVMSetAlignment(l.v.Ref(), n)
}

// Align 加载对齐
func (l Load[T]) Align() uint32 { return binding.LLVMGetAlignment(l.v.Ref()) }

// Store 存储指令角色
type Store struct {
	v llvm.Value[llvm.VoidT]
}

// Value 返回底层泛型值
func (s Store) Value() llvm.Value[llvm.VoidT] { return s.v }

// SetAlign 设置存储对齐
func (s Store) SetAlign(n uint32) {
	preAlign("ir.Store.SetAlign", n)
	binding.LLVMSetAlignment(s.v.Ref(), n)
}

// Align 存储对齐
func (s Store) Align() uint32 { return binding.LLVMGetAlignment(s.v.Ref()) }

// ===== 内存指令 =====

// Alloca 在栈上分配 t 类型的空间
func (b *Builder) Alloca(t llvm.AnyType, name string) Alloca {
	const op = "ir.Builder.Alloca"
	b.pre(op)
	if t == nil {
		errPanic(llvm.ErrInvalidArg, op, "nil type")
	}
	if t.Context() != b.ctx {
		errPanic(llvm.ErrCrossContext, op, "type belongs to another context")
	}
	ref := binding.LLVMBuildAlloca(b.ref, t.Ref(), name)
	return Alloca{v: llvm.NewValue[llvm.PtrT](b.ctx, b.inserted.life, ref)}
}

// Load 从指针加载 U 类型的值（泛型方法：结果种类由调用方断言并预检）
func (b *Builder) Load[U llvm.Kind](p llvm.ValueRef[llvm.PtrT], t llvm.TypeRef[U], name string) Load[U] {
	const op = "ir.Builder.Load"
	pv, tt := p.AsValue(), t.AsType()
	b.pre(op, pv.Dyn())
	if tt.Context() != b.ctx {
		errPanic(llvm.ErrCrossContext, op, "type belongs to another context")
	}
	ref := binding.LLVMBuildLoad(b.ref, tt.Ref(), pv.Ref(), name)
	return Load[U]{v: llvm.NewValue[U](b.ctx, b.inserted.life, ref)}
}

// Store 将值写入指针
func (b *Builder) Store(v llvm.AnyValue, p llvm.ValueRef[llvm.PtrT]) Store {
	const op = "ir.Builder.Store"
	pv := p.AsValue()
	b.pre(op, v, pv.Dyn())
	ref := binding.LLVMBuildStore(b.ref, v.Ref(), pv.Ref())
	return Store{v: llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)}
}

// GEP 插入 getelementptr
func (b *Builder) GEP(elem llvm.AnyType, p llvm.ValueRef[llvm.PtrT], idx []llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.PtrT] {
	return b.gep("ir.Builder.GEP", elem, p, idx, name, false)
}

// InBoundsGEP 插入 getelementptr inbounds
func (b *Builder) InBoundsGEP(elem llvm.AnyType, p llvm.ValueRef[llvm.PtrT], idx []llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.PtrT] {
	return b.gep("ir.Builder.InBoundsGEP", elem, p, idx, name, true)
}

func (b *Builder) gep(op string, elem llvm.AnyType, p llvm.ValueRef[llvm.PtrT], idx []llvm.ValueRef[llvm.IntT], name string, inBounds bool) llvm.Value[llvm.PtrT] {
	pv := p.AsValue()
	b.pre(op, pv.Dyn())
	if elem == nil {
		errPanic(llvm.ErrInvalidArg, op, "nil element type")
	}
	if elem.Context() != b.ctx {
		errPanic(llvm.ErrCrossContext, op, "element type belongs to another context")
	}
	refs := make([]binding.LLVMValueRef, len(idx))
	for i, x := range idx {
		v := x.AsValue()
		b.pre(op, v.Dyn())
		refs[i] = v.Ref()
	}
	var ref binding.LLVMValueRef
	if inBounds {
		ref = binding.LLVMBuildInBoundsGEP(b.ref, elem.Ref(), pv.Ref(), refs, name)
	} else {
		ref = binding.LLVMBuildGEP(b.ref, elem.Ref(), pv.Ref(), refs, name)
	}
	return llvm.NewValue[llvm.PtrT](b.ctx, b.inserted.life, ref)
}

// MemSet 插入 memset 调用
func (b *Builder) MemSet(p llvm.ValueRef[llvm.PtrT], val, n llvm.ValueRef[llvm.IntT], align uint32) Call[llvm.VoidT] {
	const op = "ir.Builder.MemSet"
	pv, vv, nv := p.AsValue(), val.AsValue(), n.AsValue()
	b.pre(op, pv.Dyn(), vv.Dyn(), nv.Dyn())
	preAlign(op, align)
	ref := binding.LLVMBuildMemSet(b.ref, pv.Ref(), vv.Ref(), nv.Ref(), align)
	return Call[llvm.VoidT]{v: llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)}
}

// MemCpy 插入 memcpy 调用
func (b *Builder) MemCpy(dst llvm.ValueRef[llvm.PtrT], dstAlign uint32, src llvm.ValueRef[llvm.PtrT], srcAlign uint32, n llvm.ValueRef[llvm.IntT]) Call[llvm.VoidT] {
	return b.memTransfer("ir.Builder.MemCpy", dst, dstAlign, src, srcAlign, n, false)
}

// MemMove 插入 memmove 调用
func (b *Builder) MemMove(dst llvm.ValueRef[llvm.PtrT], dstAlign uint32, src llvm.ValueRef[llvm.PtrT], srcAlign uint32, n llvm.ValueRef[llvm.IntT]) Call[llvm.VoidT] {
	return b.memTransfer("ir.Builder.MemMove", dst, dstAlign, src, srcAlign, n, true)
}

func (b *Builder) memTransfer(op string, dst llvm.ValueRef[llvm.PtrT], dstAlign uint32, src llvm.ValueRef[llvm.PtrT], srcAlign uint32, n llvm.ValueRef[llvm.IntT], move bool) Call[llvm.VoidT] {
	dv, sv, nv := dst.AsValue(), src.AsValue(), n.AsValue()
	b.pre(op, dv.Dyn(), sv.Dyn(), nv.Dyn())
	preAlign(op, dstAlign)
	preAlign(op, srcAlign)
	var ref binding.LLVMValueRef
	if move {
		ref = binding.LLVMBuildMemMove(b.ref, dv.Ref(), dstAlign, sv.Ref(), srcAlign, nv.Ref())
	} else {
		ref = binding.LLVMBuildMemCpy(b.ref, dv.Ref(), dstAlign, sv.Ref(), srcAlign, nv.Ref())
	}
	return Call[llvm.VoidT]{v: llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)}
}

// Malloc 插入 malloc 调用
func (b *Builder) Malloc(t llvm.AnyType, name string) Call[llvm.PtrT] {
	const op = "ir.Builder.Malloc"
	b.pre(op)
	if t == nil || t.Context() != b.ctx {
		errPanic(llvm.ErrInvalidArg, op, "invalid type")
	}
	ref := binding.LLVMBuildMalloc(b.ref, t.Ref(), name)
	return Call[llvm.PtrT]{v: llvm.NewValue[llvm.PtrT](b.ctx, b.inserted.life, ref)}
}

// MallocArray 插入数组 malloc 调用
func (b *Builder) MallocArray(elem llvm.AnyType, n llvm.ValueRef[llvm.IntT], name string) Call[llvm.PtrT] {
	const op = "ir.Builder.MallocArray"
	nv := n.AsValue()
	b.pre(op, nv.Dyn())
	if elem == nil || elem.Context() != b.ctx {
		errPanic(llvm.ErrInvalidArg, op, "invalid type")
	}
	ref := binding.LLVMBuildArrayMalloc(b.ref, elem.Ref(), nv.Ref(), name)
	return Call[llvm.PtrT]{v: llvm.NewValue[llvm.PtrT](b.ctx, b.inserted.life, ref)}
}

// Free 插入 free 调用
func (b *Builder) Free(p llvm.ValueRef[llvm.PtrT]) Call[llvm.VoidT] {
	const op = "ir.Builder.Free"
	pv := p.AsValue()
	b.pre(op, pv.Dyn())
	ref := binding.LLVMBuildFree(b.ref, pv.Ref())
	return Call[llvm.VoidT]{v: llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)}
}
