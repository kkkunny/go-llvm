package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// AtomicOrdering 原子内存序
type AtomicOrdering binding.LLVMAtomicOrdering

// AtomicOrdering 取值对应 LLVM 原子内存序（binding.LLVMAtomicOrdering）。
// AtomicNotAtomic 表示非原子访问，其余取值描述原子访问的顺序约束，除 Acquire 与 Release
// 互不可比外，强度随枚举顺序递增。
const (
	AtomicNotAtomic              = AtomicOrdering(binding.LLVMAtomicOrderingNotAtomic)              // 非原子访问，不提供原子性
	AtomicUnordered              = AtomicOrdering(binding.LLVMAtomicOrderingUnordered)              // 无序原子访问：保证原子性但不保证顺序，最弱的原子序
	AtomicMonotonic              = AtomicOrdering(binding.LLVMAtomicOrderingMonotonic)              // 单调原子访问：同一地址上的操作存在一致顺序，但不建立同步
	AtomicAcquire                = AtomicOrdering(binding.LLVMAtomicOrderingAcquire)                // 获取语义：之后的读写不能重排到本操作之前（用于加载）
	AtomicRelease                = AtomicOrdering(binding.LLVMAtomicOrderingRelease)                // 释放语义：之前的读写不能重排到本操作之后（用于存储）
	AtomicAcquireRelease         = AtomicOrdering(binding.LLVMAtomicOrderingAcquireRelease)         // 获取-释放语义：同时具备获取与释放约束（用于读改写与栅栏）
	AtomicSequentiallyConsistent = AtomicOrdering(binding.LLVMAtomicOrderingSequentiallyConsistent) // 顺序一致：在获取-释放基础上，所有此类操作间存在全局全序
)

// RMWOp 原子读改写操作
type RMWOp binding.LLVMAtomicRMWBinOp

// RMWOp 取值对应 LLVM 原子读改写二元操作（binding.LLVMAtomicRMWBinOp）。
// 所有操作都原子地更新目标内存并返回修改前的旧值；Max/Min 按有符号解释，
// UMax/UMin 与 U 前缀系列按无符号解释，F 前缀系列按浮点解释。
const (
	RMWXchg     = RMWOp(binding.LLVMAtomicRMWBinOpXchg)     // 原子交换：写入操作数并返回旧值
	RMWAdd      = RMWOp(binding.LLVMAtomicRMWBinOpAdd)      // 原子加法
	RMWSub      = RMWOp(binding.LLVMAtomicRMWBinOpSub)      // 原子减法
	RMWAnd      = RMWOp(binding.LLVMAtomicRMWBinOpAnd)      // 原子按位与
	RMWNand     = RMWOp(binding.LLVMAtomicRMWBinOpNand)     // 原子按位与非（~(旧值 & 操作数)）
	RMWOr       = RMWOp(binding.LLVMAtomicRMWBinOpOr)       // 原子按位或
	RMWXor      = RMWOp(binding.LLVMAtomicRMWBinOpXor)      // 原子按位异或
	RMWMax      = RMWOp(binding.LLVMAtomicRMWBinOpMax)      // 有符号最大值
	RMWMin      = RMWOp(binding.LLVMAtomicRMWBinOpMin)      // 有符号最小值
	RMWUMax     = RMWOp(binding.LLVMAtomicRMWBinOpUMax)     // 无符号最大值
	RMWUMin     = RMWOp(binding.LLVMAtomicRMWBinOpUMin)     // 无符号最小值
	RMWFAdd     = RMWOp(binding.LLVMAtomicRMWBinOpFAdd)     // 浮点加法
	RMWFSub     = RMWOp(binding.LLVMAtomicRMWBinOpFSub)     // 浮点减法
	RMWFMax     = RMWOp(binding.LLVMAtomicRMWBinOpFMax)     // 浮点最大值（maxnum 语义：忽略单个 qNaN 操作数；两操作数均为 NaN 或存在 sNaN 时结果为 NaN）
	RMWFMin     = RMWOp(binding.LLVMAtomicRMWBinOpFMin)     // 浮点最小值（minnum 语义：忽略单个 qNaN 操作数；两操作数均为 NaN 或存在 sNaN 时结果为 NaN）
	RMWUIncWrap = RMWOp(binding.LLVMAtomicRMWBinOpUIncWrap) // 无符号自增一：旧值不小于操作数（视为上限）时回绕为 0
	RMWUDecWrap = RMWOp(binding.LLVMAtomicRMWBinOpUDecWrap) // 无符号自减一：旧值为 0 或大于操作数（视为下限）时取操作数
	RMWUSubCond = RMWOp(binding.LLVMAtomicRMWBinOpUSubCond) // 无符号条件减法：仅当不发生下溢时相减，否则保持旧值
	RMWUSubSat  = RMWOp(binding.LLVMAtomicRMWBinOpUSubSat)  // 无符号饱和减法：下溢时结果为 0
	RMWFMaximum = RMWOp(binding.LLVMAtomicRMWBinOpFMaximum) // 浮点最大值（maximum 语义：任一操作数为 NaN 时结果为 NaN）
	RMWFMinimum = RMWOp(binding.LLVMAtomicRMWBinOpFMinimum) // 浮点最小值（minimum 语义：任一操作数为 NaN 时结果为 NaN）
)
