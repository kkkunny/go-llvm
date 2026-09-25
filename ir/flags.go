package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
)

// FastMath 浮点指令 fast-math flags 位掩码
type FastMath binding.LLVMFastMathFlags

// FastMath 取值对应 LLVM fast-math flags（binding.LLVMFastMathFlags），可按位或组合，
// 用于浮点算术、比较与转换等指令（见 CanFastMath）。各标志放松浮点语义假设以启用更多
// 优化，违反假设时结果可能是 poison，如 FastMathNoNaNs/FastMathNoInfs。
const (
	FastMathNone            = FastMath(binding.LLVMFastMathNone)            // 无标志：不放松任何浮点语义
	FastMathAllowReassoc    = FastMath(binding.LLVMFastMathAllowReassoc)    // 允许重结合：浮点运算可按代数等价的方式重新结合
	FastMathNoNaNs          = FastMath(binding.LLVMFastMathNoNaNs)          // 假设无 NaN：操作数为 NaN 或结果本应为 NaN 时得到 poison
	FastMathNoInfs          = FastMath(binding.LLVMFastMathNoInfs)          // 假设无无穷大：操作数为 ±Inf 或结果本应为 ±Inf 时得到 poison
	FastMathNoSignedZeros   = FastMath(binding.LLVMFastMathNoSignedZeros)   // 忽略零的符号：±0.0 的符号位可被非确定地翻转
	FastMathAllowReciprocal = FastMath(binding.LLVMFastMathAllowReciprocal) // 允许用倒数近似除法：a/b 可改写为 a*(1/b)
	FastMathAllowContract   = FastMath(binding.LLVMFastMathAllowContract)   // 允许收缩：乘法与加法可融合为乘加（FMA），但不做重结合
	FastMathApproxFunc      = FastMath(binding.LLVMFastMathApproxFunc)      // 允许近似函数：sin/log/sqrt 等可替换为近似计算
	FastMathAll             = FastMath(binding.LLVMFastMathAll)             // 以上所有标志的组合（等价于 fast）
)

// CanFastMath 该指令是否可携带 fast-math flags（浮点算术/转换等）
func CanFastMath(inst llvm.AnyValue) bool {
	const op = "ir.CanFastMath"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	return binding.LLVMCanValueUseFastMathFlags(inst.Ref())
}

// FastMathOf 读取 fast-math flags（不可携带的指令 panic，仅调试层）
func FastMathOf(inst llvm.AnyValue) FastMath {
	const op = "ir.FastMathOf"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	if checks.Debug && !CanFastMath(inst) {
		llvm.Panicf(llvm.ErrInvalidArg, op, "instruction cannot carry fast-math flags")
	}
	return FastMath(binding.LLVMGetFastMathFlags(inst.Ref()))
}

// SetFastMath 设置 fast-math flags
func SetFastMath(inst llvm.AnyValue, f FastMath) {
	const op = "ir.SetFastMath"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	if checks.Debug && !CanFastMath(inst) {
		llvm.Panicf(llvm.ErrInvalidArg, op, "instruction cannot carry fast-math flags")
	}
	binding.LLVMSetFastMathFlags(inst.Ref(), binding.LLVMFastMathFlags(f))
}

// NoWrap GEP 无回绕 flags 位掩码
type NoWrap binding.LLVMGEPNoWrapFlags

// NoWrap 取值对应 LLVM GEP no-wrap flags（binding.LLVMGEPNoWrapFlags），可按位或组合。
// 违反任一保证时 GEP 的结果为 poison。
const (
	NoWrapNone     = NoWrap(0)                           // 无额外保证
	NoWrapInBounds = NoWrap(binding.LLVMGEPFlagInBounds) // inbounds：基址与所有中间地址都位于同一已分配对象内（蕴含 nusw 规则）
	NoWrapNUSW     = NoWrap(binding.LLVMGEPFlagNUSW)     // nusw（no unsigned signed wrap）：索引按有符号解释时截断/乘加不回绕
	NoWrapNUW      = NoWrap(binding.LLVMGEPFlagNUW)      // nuw（no unsigned wrap）：索引按无符号解释时截断/乘加不回绕
)

// GEPNoWrapOf 读取 GEP no-wrap flags
func GEPNoWrapOf(v llvm.AnyValue) NoWrap {
	const op = "ir.GEPNoWrapOf"
	if v == nil || !v.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	return NoWrap(binding.LLVMGEPGetNoWrapFlags(v.Ref()))
}

// SetGEPNoWrap 设置 GEP no-wrap flags
func SetGEPNoWrap(v llvm.AnyValue, f NoWrap) {
	const op = "ir.SetGEPNoWrap"
	if v == nil || !v.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	binding.LLVMGEPSetNoWrapFlags(v.Ref(), binding.LLVMGEPNoWrapFlags(f))
}

// SetNSW 设置算术指令 no-signed-wrap 标志
func SetNSW(inst llvm.AnyValue, v bool) {
	const op = "ir.SetNSW"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	binding.LLVMSetNSW(inst.Ref(), v)
}

// SetNUW 设置算术指令 no-unsigned-wrap 标志
func SetNUW(inst llvm.AnyValue, v bool) {
	const op = "ir.SetNUW"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	binding.LLVMSetNUW(inst.Ref(), v)
}

// SetExact 设置 div/rem/shift 的 exact 标志
func SetExact(inst llvm.AnyValue, v bool) {
	const op = "ir.SetExact"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	binding.LLVMSetExact(inst.Ref(), v)
}

// SetNNeg 设置 sext 类指令 non-negative 标志
func SetNNeg(inst llvm.AnyValue, v bool) {
	const op = "ir.SetNNeg"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	binding.LLVMSetNNeg(inst.Ref(), v)
}

// SetInBounds 设置 GEP inbounds 标志（与 SetGEPNoWrap 的 NoWrapInBounds 等价）
func SetInBounds(inst llvm.AnyValue, v bool) {
	const op = "ir.SetInBounds"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	binding.LLVMSetIsInBounds(inst.Ref(), v)
}

// TailCallKind tail-call 种类
type TailCallKind binding.LLVMTailCallKind

// TailCallKind 取值对应 LLVM tail-call 种类（binding.LLVMTailCallKind）。
const (
	TailCallNone = TailCallKind(binding.LLVMTailCallKindNone)     // 未指定：由后端自行决定是否做尾调用
	TailCallTail = TailCallKind(binding.LLVMTailCallKindTail)     // tail：提示后端做尾调用，但不保证
	TailCallMust = TailCallKind(binding.LLVMTailCallKindMustTail) // musttail：强制尾调用，不满足结构约束即为非法 IR
	TailCallNo   = TailCallKind(binding.LLVMTailCallKindNoTail)   // notail：禁止把该调用优化为尾调用
)

// SetTailCall 设置调用指令的 tail 标志
func SetTailCall(inst llvm.AnyValue, v bool) {
	const op = "ir.SetTailCall"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	binding.LLVMSetTailCall(inst.Ref(), v)
}

// IsTailCall 调用指令是否 tail 调用
func IsTailCall(inst llvm.AnyValue) bool {
	const op = "ir.IsTailCall"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	return binding.LLVMIsTailCall(inst.Ref())
}

// SetTailCallKind 设置 tail-call 种类。
// LLVM-C 无种类 getter（musttail/nottail 无法读回），需要精确值请解析 IR 文本。
func SetTailCallKind(inst llvm.AnyValue, k TailCallKind) {
	const op = "ir.SetTailCallKind"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	binding.LLVMSetTailCallKind(inst.Ref(), binding.LLVMTailCallKind(k))
}

// SetParamAlign 设置调用指令第 i 个实参的对齐（i 从 0 起）
func SetParamAlign(inst llvm.AnyValue, i uint32, align uint32) {
	const op = "ir.SetParamAlign"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	preAlign(op, align)
	if checks.Debug && i >= uint32(binding.LLVMGetNumArgOperands(inst.Ref())) {
		llvm.Panicf(llvm.ErrInvalidArg, op, "argument index %d out of range", i)
	}
	binding.LLVMSetInstrParamAlignment(inst.Ref(), i+1, align) // LLVM-C 索引 1-based，0 为返回值
}

// SyncScopeOf 读取原子指令的 sync scope ID。
//
// ID 由 LLVM 分配，常见作用域名（如 "system"、"singlethread"）的对应数值属 LLVM 内部
// 约定，勿硬编码；按名查询 ID 用 [llvm.Context.SyncScopeID]。
func SyncScopeOf(inst llvm.AnyValue) uint32 {
	const op = "ir.SyncScopeOf"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	return binding.LLVMGetAtomicSyncScopeID(inst.Ref())
}

// SetSyncScope 设置原子指令的 sync scope ID
func SetSyncScope(inst llvm.AnyValue, ssid uint32) {
	const op = "ir.SetSyncScope"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	binding.LLVMSetAtomicSyncScopeID(inst.Ref(), ssid)
}
