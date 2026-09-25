package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
)

// Alloca alloca 指令角色（内嵌 Value[PtrT]，自动实现 llvm.ValueRef/AnyValue）
type Alloca struct {
	llvm.Value[llvm.PtrT]
}

// SetAlign 设置分配对齐
func (a Alloca) SetAlign(n uint32) {
	const op = "ir.Alloca.SetAlign"
	preAlign(op, n)
	binding.LLVMSetAlignment(a.Ref(), n)
}

// Align 分配对齐
func (a Alloca) Align() uint32 {
	return binding.LLVMGetAlignment(a.Ref())
}

// Load 加载指令角色（内嵌 Value[T]）
type Load[T llvm.Kind] struct {
	llvm.Value[T]
}

// SetAlign 设置加载对齐
func (l Load[T]) SetAlign(n uint32) {
	const op = "ir.Load.SetAlign"
	preAlign(op, n)
	binding.LLVMSetAlignment(l.Ref(), n)
}

// Align 加载对齐
func (l Load[T]) Align() uint32 {
	return binding.LLVMGetAlignment(l.Ref())
}

// SetVolatile 设置 volatile 访问
func (l Load[T]) SetVolatile(v bool) {
	binding.LLVMSetVolatile(l.Ref(), v)
}

// IsVolatile 是否 volatile 访问
func (l Load[T]) IsVolatile() bool {
	return binding.LLVMGetVolatile(l.Ref())
}

// SetOrdering 设置原子内存序
func (l Load[T]) SetOrdering(o llvm.AtomicOrdering) {
	const op = "ir.Load.SetOrdering"
	preOrderingLoad(op, o)
	binding.LLVMSetOrdering(l.Ref(), binding.LLVMAtomicOrdering(o))
}

// Ordering 原子内存序（非原子访问为 AtomicNotAtomic）
func (l Load[T]) Ordering() llvm.AtomicOrdering {
	return llvm.AtomicOrdering(binding.LLVMGetOrdering(l.Ref()))
}

// Store 存储指令角色（内嵌 Value[VoidT]）
type Store struct {
	llvm.Value[llvm.VoidT]
}

// SetAlign 设置存储对齐
func (s Store) SetAlign(n uint32) {
	const op = "ir.Store.SetAlign"
	preAlign(op, n)
	binding.LLVMSetAlignment(s.Ref(), n)
}

// Align 存储对齐
func (s Store) Align() uint32 {
	return binding.LLVMGetAlignment(s.Ref())
}

// SetVolatile 设置 volatile 访问
func (s Store) SetVolatile(v bool) {
	binding.LLVMSetVolatile(s.Ref(), v)
}

// IsVolatile 是否 volatile 访问
func (s Store) IsVolatile() bool {
	return binding.LLVMGetVolatile(s.Ref())
}

// SetOrdering 设置原子内存序
func (s Store) SetOrdering(o llvm.AtomicOrdering) {
	const op = "ir.Store.SetOrdering"
	preOrderingStore(op, o)
	binding.LLVMSetOrdering(s.Ref(), binding.LLVMAtomicOrdering(o))
}

// Ordering 原子内存序（非原子访问为 AtomicNotAtomic）
func (s Store) Ordering() llvm.AtomicOrdering {
	return llvm.AtomicOrdering(binding.LLVMGetOrdering(s.Ref()))
}

// ===== 内存指令 =====

// Alloca 在栈上分配 t 类型的空间
func (b *Builder) Alloca(t llvm.AnyType, name string) Alloca {
	const op = "ir.Builder.Alloca"
	b.pre(op)
	if t == nil {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil type")
	}
	if t.Context() != b.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "type belongs to another context")
	}
	ref := binding.LLVMBuildAlloca(b.ref, t.Ref(), name)
	return Alloca{Value: llvm.NewValue[llvm.PtrT](b.ctx, b.inserted.life, ref)}
}

// Load 从指针加载 U 类型的值（泛型方法：结果种类由调用方断言并预检）
func (b *Builder) Load[U llvm.Kind](p llvm.ValueRef[llvm.PtrT], t llvm.TypeRef[U], name string) Load[U] {
	const op = "ir.Builder.Load"
	pv, tt := p.AsValue(), t.AsType()
	b.pre(op, core(pv))
	if tt.Context() != b.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "type belongs to another context")
	}
	ref := binding.LLVMBuildLoad(b.ref, tt.Ref(), pv.Ref(), name)
	return Load[U]{Value: llvm.NewValue[U](b.ctx, b.inserted.life, ref)}
}

// Store 将值写入指针
func (b *Builder) Store(v llvm.AnyValue, p llvm.ValueRef[llvm.PtrT]) Store {
	const op = "ir.Builder.Store"
	pv := p.AsValue()
	b.pre(op, coreAny(v), core(pv))
	ref := binding.LLVMBuildStore(b.ref, v.Ref(), pv.Ref())
	return Store{Value: llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)}
}

// GEP 插入 getelementptr
func (b *Builder) GEP(elem llvm.AnyType, p llvm.ValueRef[llvm.PtrT], idx []llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.PtrT] {
	return b.gep("ir.Builder.GEP", elem, p, idx, name, false)
}

// InBoundsGEP 插入 getelementptr inbounds
func (b *Builder) InBoundsGEP(elem llvm.AnyType, p llvm.ValueRef[llvm.PtrT], idx []llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.PtrT] {
	return b.gep("ir.Builder.InBoundsGEP", elem, p, idx, name, true)
}

// PtrAdd 插入按字节偏移的指针加法（等价 GEP i8, ptr, off）
func (b *Builder) PtrAdd(p llvm.ValueRef[llvm.PtrT], off llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.PtrT] {
	return b.gep("ir.Builder.PtrAdd", b.i8Type(), p, []llvm.ValueRef[llvm.IntT]{off}, name, false)
}

// IsNull 判断值为 null（ptr/int 均可）
func (b *Builder) IsNull(v llvm.AnyValue, name string) llvm.Value[llvm.IntT] {
	const op = "ir.Builder.IsNull"
	b.pre(op, coreAny(v))
	return llvm.NewValue[llvm.IntT](b.ctx, b.inserted.life, binding.LLVMBuildIsNull(b.ref, v.Ref(), name))
}

// IsNotNull 判断值非 null
func (b *Builder) IsNotNull(v llvm.AnyValue, name string) llvm.Value[llvm.IntT] {
	const op = "ir.Builder.IsNotNull"
	b.pre(op, coreAny(v))
	return llvm.NewValue[llvm.IntT](b.ctx, b.inserted.life, binding.LLVMBuildIsNotNull(b.ref, v.Ref(), name))
}

// PtrDiff 两指针按 elem 元素大小的差值（结果种类由 LLVM 决定，为整数）
func (b *Builder) PtrDiff(elem llvm.AnyType, l, r llvm.ValueRef[llvm.PtrT], name string) llvm.Value[llvm.IntT] {
	const op = "ir.Builder.PtrDiff"
	lv, rv := l.AsValue(), r.AsValue()
	b.pre(op, core(lv), core(rv))
	b.ctx.CheckType(op, elem)
	return llvm.NewValue[llvm.IntT](b.ctx, b.inserted.life, binding.LLVMBuildPtrDiff(b.ref, elem.Ref(), lv.Ref(), rv.Ref(), name))
}

func (b *Builder) gep(op string, elem llvm.AnyType, p llvm.ValueRef[llvm.PtrT], idx []llvm.ValueRef[llvm.IntT], name string, inBounds bool) llvm.Value[llvm.PtrT] {
	pv := p.AsValue()
	b.pre(op, core(pv))
	if elem == nil {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil element type")
	}
	if elem.Context() != b.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "element type belongs to another context")
	}
	// 调试层不复用 scratch（B4）
	var refs []binding.LLVMValueRef
	if !checks.Debug {
		refs = b.refs[:0]
	}
	for _, x := range idx {
		v := x.AsValue()
		b.checkVal(op, core(v))
		refs = append(refs, v.RawRef())
	}
	if !checks.Debug {
		b.refs = refs
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
	b.pre(op, core(pv), core(vv), core(nv))
	preAlign(op, align)
	ref := binding.LLVMBuildMemSet(b.ref, pv.Ref(), vv.Ref(), nv.Ref(), align)
	return Call[llvm.VoidT]{Value: llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)}
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
	b.pre(op, core(dv), core(sv), core(nv))
	preAlign(op, dstAlign)
	preAlign(op, srcAlign)
	var ref binding.LLVMValueRef
	if move {
		ref = binding.LLVMBuildMemMove(b.ref, dv.Ref(), dstAlign, sv.Ref(), srcAlign, nv.Ref())
	} else {
		ref = binding.LLVMBuildMemCpy(b.ref, dv.Ref(), dstAlign, sv.Ref(), srcAlign, nv.Ref())
	}
	return Call[llvm.VoidT]{Value: llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)}
}

// Malloc 插入 malloc 调用
func (b *Builder) Malloc(t llvm.AnyType, name string) Call[llvm.PtrT] {
	const op = "ir.Builder.Malloc"
	b.pre(op)
	if t == nil || t.Context() != b.ctx {
		llvm.Panicf(llvm.ErrInvalidArg, op, "invalid type")
	}
	ref := binding.LLVMBuildMalloc(b.ref, t.Ref(), name)
	return Call[llvm.PtrT]{Value: llvm.NewValue[llvm.PtrT](b.ctx, b.inserted.life, ref)}
}

// MallocArray 插入数组 malloc 调用
func (b *Builder) MallocArray(elem llvm.AnyType, n llvm.ValueRef[llvm.IntT], name string) Call[llvm.PtrT] {
	const op = "ir.Builder.MallocArray"
	nv := n.AsValue()
	b.pre(op, core(nv))
	if elem == nil || elem.Context() != b.ctx {
		llvm.Panicf(llvm.ErrInvalidArg, op, "invalid type")
	}
	ref := binding.LLVMBuildArrayMalloc(b.ref, elem.Ref(), nv.Ref(), name)
	return Call[llvm.PtrT]{Value: llvm.NewValue[llvm.PtrT](b.ctx, b.inserted.life, ref)}
}

// Free 插入 free 调用
func (b *Builder) Free(p llvm.ValueRef[llvm.PtrT]) Call[llvm.VoidT] {
	const op = "ir.Builder.Free"
	pv := p.AsValue()
	b.pre(op, core(pv))
	ref := binding.LLVMBuildFree(b.ref, pv.Ref())
	return Call[llvm.VoidT]{Value: llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)}
}
