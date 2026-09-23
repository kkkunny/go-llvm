package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// AtomicOrdering 原子内存序
type AtomicOrdering binding.LLVMAtomicOrdering

const (
	AtomicNotAtomic              = AtomicOrdering(binding.LLVMAtomicOrderingNotAtomic)
	AtomicUnordered              = AtomicOrdering(binding.LLVMAtomicOrderingUnordered)
	AtomicMonotonic              = AtomicOrdering(binding.LLVMAtomicOrderingMonotonic)
	AtomicAcquire                = AtomicOrdering(binding.LLVMAtomicOrderingAcquire)
	AtomicRelease                = AtomicOrdering(binding.LLVMAtomicOrderingRelease)
	AtomicAcquireRelease         = AtomicOrdering(binding.LLVMAtomicOrderingAcquireRelease)
	AtomicSequentiallyConsistent = AtomicOrdering(binding.LLVMAtomicOrderingSequentiallyConsistent)
)

// RMWOp 原子读改写操作
type RMWOp binding.LLVMAtomicRMWBinOp

const (
	RMWXchg     = RMWOp(binding.LLVMAtomicRMWBinOpXchg)
	RMWAdd      = RMWOp(binding.LLVMAtomicRMWBinOpAdd)
	RMWSub      = RMWOp(binding.LLVMAtomicRMWBinOpSub)
	RMWAnd      = RMWOp(binding.LLVMAtomicRMWBinOpAnd)
	RMWNand     = RMWOp(binding.LLVMAtomicRMWBinOpNand)
	RMWOr       = RMWOp(binding.LLVMAtomicRMWBinOpOr)
	RMWXor      = RMWOp(binding.LLVMAtomicRMWBinOpXor)
	RMWMax      = RMWOp(binding.LLVMAtomicRMWBinOpMax)
	RMWMin      = RMWOp(binding.LLVMAtomicRMWBinOpMin)
	RMWUMax     = RMWOp(binding.LLVMAtomicRMWBinOpUMax)
	RMWUMin     = RMWOp(binding.LLVMAtomicRMWBinOpUMin)
	RMWFAdd     = RMWOp(binding.LLVMAtomicRMWBinOpFAdd)
	RMWFSub     = RMWOp(binding.LLVMAtomicRMWBinOpFSub)
	RMWFMax     = RMWOp(binding.LLVMAtomicRMWBinOpFMax)
	RMWFMin     = RMWOp(binding.LLVMAtomicRMWBinOpFMin)
	RMWUIncWrap = RMWOp(binding.LLVMAtomicRMWBinOpUIncWrap)
	RMWUDecWrap = RMWOp(binding.LLVMAtomicRMWBinOpUDecWrap)
	RMWUSubCond = RMWOp(binding.LLVMAtomicRMWBinOpUSubCond)
	RMWUSubSat  = RMWOp(binding.LLVMAtomicRMWBinOpUSubSat)
	RMWFMaximum = RMWOp(binding.LLVMAtomicRMWBinOpFMaximum)
	RMWFMinimum = RMWOp(binding.LLVMAtomicRMWBinOpFMinimum)
)
