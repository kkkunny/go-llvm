package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// ===== 内存序预检 =====

// preOrderingFence fence 合法内存序：acquire/release/acq_rel/seq_cst
func preOrderingFence(op string, o llvm.AtomicOrdering) {
	switch o {
	case llvm.AtomicAcquire, llvm.AtomicRelease, llvm.AtomicAcquireRelease, llvm.AtomicSequentiallyConsistent:
	default:
		llvm.Panicf(llvm.ErrInvalidArg, op, "invalid fence ordering %d", o)
	}
}

// preOrderingRMW atomicrmw/成功比较交换合法内存序：monotonic 起（not_atomic/unordered 非法）
func preOrderingRMW(op string, o llvm.AtomicOrdering) {
	switch o {
	case llvm.AtomicMonotonic, llvm.AtomicAcquire, llvm.AtomicRelease, llvm.AtomicAcquireRelease, llvm.AtomicSequentiallyConsistent:
	default:
		llvm.Panicf(llvm.ErrInvalidArg, op, "invalid atomic ordering %d", o)
	}
}

// orderingRank 内存序强度序（用于失败序不得强于成功序）
func orderingRank(o llvm.AtomicOrdering) int {
	switch o {
	case llvm.AtomicNotAtomic:
		return 0
	case llvm.AtomicUnordered:
		return 1
	case llvm.AtomicMonotonic:
		return 2
	case llvm.AtomicAcquire, llvm.AtomicRelease:
		return 3
	case llvm.AtomicAcquireRelease:
		return 4
	default:
		return 5 // seq_cst
	}
}

// preOrderingCmpXchg 比较交换内存序预检：失败序不得为 release/acq_rel 且不得强于成功序
func preOrderingCmpXchg(op string, success, failure llvm.AtomicOrdering) {
	preOrderingRMW(op, success)
	switch failure {
	case llvm.AtomicUnordered, llvm.AtomicMonotonic, llvm.AtomicAcquire, llvm.AtomicSequentiallyConsistent:
	default:
		llvm.Panicf(llvm.ErrInvalidArg, op, "invalid cmpxchg failure ordering %d", failure)
	}
	if orderingRank(failure) > orderingRank(success) {
		llvm.Panicf(llvm.ErrInvalidArg, op, "failure ordering %d is stronger than success ordering %d", failure, success)
	}
}

// preOrderingLoad load 合法内存序：unordered/monotonic/acquire/seq_cst（及 not_atomic 复位）
func preOrderingLoad(op string, o llvm.AtomicOrdering) {
	switch o {
	case llvm.AtomicNotAtomic, llvm.AtomicUnordered, llvm.AtomicMonotonic, llvm.AtomicAcquire, llvm.AtomicSequentiallyConsistent:
	default:
		llvm.Panicf(llvm.ErrInvalidArg, op, "invalid load ordering %d", o)
	}
}

// preOrderingStore store 合法内存序：unordered/monotonic/release/seq_cst（及 not_atomic 复位）
func preOrderingStore(op string, o llvm.AtomicOrdering) {
	switch o {
	case llvm.AtomicNotAtomic, llvm.AtomicUnordered, llvm.AtomicMonotonic, llvm.AtomicRelease, llvm.AtomicSequentiallyConsistent:
	default:
		llvm.Panicf(llvm.ErrInvalidArg, op, "invalid store ordering %d", o)
	}
}

// ===== 指令角色 =====

// Fence fence 指令角色（内嵌 Value[VoidT]；void 值指令不可命名）
type Fence struct {
	llvm.Value[llvm.VoidT]
}

// Ordering 内存序
func (f Fence) Ordering() llvm.AtomicOrdering {
	f.Check("ir.Fence.Ordering")
	return llvm.AtomicOrdering(binding.LLVMGetOrdering(f.Ref()))
}

// SetOrdering 设置内存序
func (f Fence) SetOrdering(o llvm.AtomicOrdering) {
	const op = "ir.Fence.SetOrdering"
	f.Check(op)
	preOrderingFence(op, o)
	binding.LLVMSetOrdering(f.Ref(), binding.LLVMAtomicOrdering(o))
}

// IsSingleThread 是否 syncscope("singlethread")
func (f Fence) IsSingleThread() bool {
	f.Check("ir.Fence.IsSingleThread")
	return binding.LLVMIsAtomicSingleThread(f.Ref())
}

// SetSingleThread 设置 syncscope("singlethread")
func (f Fence) SetSingleThread(v bool) {
	f.Check("ir.Fence.SetSingleThread")
	binding.LLVMSetAtomicSingleThread(f.Ref(), v)
}

// AtomicRMW atomicrmw 指令角色（内嵌 Value[T]，T 为操作数种类，返回旧值）
type AtomicRMW[T llvm.Kind] struct {
	llvm.Value[T]
}

// Op 读改写操作
func (r AtomicRMW[T]) Op() llvm.RMWOp {
	r.Check("ir.AtomicRMW.Op")
	return llvm.RMWOp(binding.LLVMGetAtomicRMWBinOp(r.Ref()))
}

// SetOp 设置读改写操作
func (r AtomicRMW[T]) SetOp(op llvm.RMWOp) {
	r.Check("ir.AtomicRMW.SetOp")
	binding.LLVMSetAtomicRMWBinOp(r.Ref(), binding.LLVMAtomicRMWBinOp(op))
}

// Ordering 内存序
func (r AtomicRMW[T]) Ordering() llvm.AtomicOrdering {
	r.Check("ir.AtomicRMW.Ordering")
	return llvm.AtomicOrdering(binding.LLVMGetOrdering(r.Ref()))
}

// SetOrdering 设置内存序
func (r AtomicRMW[T]) SetOrdering(o llvm.AtomicOrdering) {
	const op = "ir.AtomicRMW.SetOrdering"
	r.Check(op)
	preOrderingRMW(op, o)
	binding.LLVMSetOrdering(r.Ref(), binding.LLVMAtomicOrdering(o))
}

// IsVolatile 是否 volatile
func (r AtomicRMW[T]) IsVolatile() bool {
	r.Check("ir.AtomicRMW.IsVolatile")
	return binding.LLVMGetVolatile(r.Ref())
}

// SetVolatile 设置 volatile
func (r AtomicRMW[T]) SetVolatile(v bool) {
	r.Check("ir.AtomicRMW.SetVolatile")
	binding.LLVMSetVolatile(r.Ref(), v)
}

// Align 对齐字节数
func (r AtomicRMW[T]) Align() uint32 {
	r.Check("ir.AtomicRMW.Align")
	return binding.LLVMGetAlignment(r.Ref())
}

// SetAlign 设置对齐字节数
func (r AtomicRMW[T]) SetAlign(n uint32) {
	const op = "ir.AtomicRMW.SetAlign"
	r.Check(op)
	preAlign(op, n)
	binding.LLVMSetAlignment(r.Ref(), n)
}

// CmpXchg cmpxchg 指令角色（内嵌 Value[StructT]，结果为 {旧值, i1 成功标志}）
type CmpXchg struct {
	llvm.Value[llvm.StructT]
}

// SuccessOrdering 成功序
func (x CmpXchg) SuccessOrdering() llvm.AtomicOrdering {
	x.Check("ir.CmpXchg.SuccessOrdering")
	return llvm.AtomicOrdering(binding.LLVMGetCmpXchgSuccessOrdering(x.Ref()))
}

// SetSuccessOrdering 设置成功序
func (x CmpXchg) SetSuccessOrdering(o llvm.AtomicOrdering) {
	const op = "ir.CmpXchg.SetSuccessOrdering"
	x.Check(op)
	preOrderingRMW(op, o)
	binding.LLVMSetCmpXchgSuccessOrdering(x.Ref(), binding.LLVMAtomicOrdering(o))
}

// FailureOrdering 失败序
func (x CmpXchg) FailureOrdering() llvm.AtomicOrdering {
	x.Check("ir.CmpXchg.FailureOrdering")
	return llvm.AtomicOrdering(binding.LLVMGetCmpXchgFailureOrdering(x.Ref()))
}

// SetFailureOrdering 设置失败序（不得强于成功序）
func (x CmpXchg) SetFailureOrdering(o llvm.AtomicOrdering) {
	const op = "ir.CmpXchg.SetFailureOrdering"
	x.Check(op)
	preOrderingCmpXchg(op, x.SuccessOrdering(), o)
	binding.LLVMSetCmpXchgFailureOrdering(x.Ref(), binding.LLVMAtomicOrdering(o))
}

// IsWeak 是否 weak
func (x CmpXchg) IsWeak() bool {
	x.Check("ir.CmpXchg.IsWeak")
	return binding.LLVMGetWeak(x.Ref())
}

// SetWeak 设置 weak
func (x CmpXchg) SetWeak(v bool) {
	x.Check("ir.CmpXchg.SetWeak")
	binding.LLVMSetWeak(x.Ref(), v)
}

// IsVolatile 是否 volatile
func (x CmpXchg) IsVolatile() bool {
	x.Check("ir.CmpXchg.IsVolatile")
	return binding.LLVMGetVolatile(x.Ref())
}

// SetVolatile 设置 volatile
func (x CmpXchg) SetVolatile(v bool) {
	x.Check("ir.CmpXchg.SetVolatile")
	binding.LLVMSetVolatile(x.Ref(), v)
}

// Align 对齐字节数
func (x CmpXchg) Align() uint32 {
	x.Check("ir.CmpXchg.Align")
	return binding.LLVMGetAlignment(x.Ref())
}

// SetAlign 设置对齐字节数
func (x CmpXchg) SetAlign(n uint32) {
	const op = "ir.CmpXchg.SetAlign"
	x.Check(op)
	preAlign(op, n)
	binding.LLVMSetAlignment(x.Ref(), n)
}

// ===== 构建方法 =====

// Fence 插入 fence（void 值指令不可命名，故无 name 参数）
func (b *Builder) Fence(order llvm.AtomicOrdering, singleThread bool) Fence {
	const op = "ir.Builder.Fence"
	b.pre(op)
	preOrderingFence(op, order)
	ref := binding.LLVMBuildFence(b.ref, binding.LLVMAtomicOrdering(order), singleThread, "")
	return Fence{Value: llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)}
}

// AtomicRMW 插入 atomicrmw，返回旧值（种类 T 与 val 一致）；
// C 构建 API 无 name，构建后经 SetValueName 命名
func (b *Builder) AtomicRMW[T llvm.Kind](op llvm.RMWOp, ptr llvm.ValueRef[llvm.PtrT], val llvm.ValueRef[T], order llvm.AtomicOrdering, singleThread bool, name string) AtomicRMW[T] {
	const opr = "ir.Builder.AtomicRMW"
	pv, vv := ptr.AsValue(), val.AsValue()
	b.pre(opr, core(pv), core(vv))
	preOrderingRMW(opr, order)
	ref := binding.LLVMBuildAtomicRMW(b.ref, binding.LLVMAtomicRMWBinOp(op), pv.Ref(), vv.Ref(), binding.LLVMAtomicOrdering(order), singleThread)
	binding.LLVMSetValueName(ref, name)
	return AtomicRMW[T]{Value: llvm.NewValue[T](b.ctx, b.inserted.life, ref)}
}

// CmpXchg 插入原子比较交换，返回 {旧值, i1 成功标志}；cmp/new 须同类型。
// C 构建 API 无 name，构建后经 SetValueName 命名
func (b *Builder) CmpXchg(ptr llvm.ValueRef[llvm.PtrT], cmp, new llvm.AnyValue, success, failure llvm.AtomicOrdering, singleThread bool, name string) CmpXchg {
	const op = "ir.Builder.CmpXchg"
	pv := ptr.AsValue()
	b.pre(op, core(pv), coreAny(cmp), coreAny(new))
	b.preSameType(op, coreAny(cmp), coreAny(new))
	preOrderingCmpXchg(op, success, failure)
	ref := binding.LLVMBuildAtomicCmpXchg(b.ref, pv.Ref(), cmp.Ref(), new.Ref(),
		binding.LLVMAtomicOrdering(success), binding.LLVMAtomicOrdering(failure), singleThread)
	binding.LLVMSetValueName(ref, name)
	return CmpXchg{Value: llvm.NewValue[llvm.StructT](b.ctx, b.inserted.life, ref)}
}
