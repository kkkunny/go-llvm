package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
)

// FastMath 浮点指令 fast-math flags 位掩码
type FastMath binding.LLVMFastMathFlags

const (
	FastMathNone            = FastMath(binding.LLVMFastMathNone)
	FastMathAllowReassoc    = FastMath(binding.LLVMFastMathAllowReassoc)
	FastMathNoNaNs          = FastMath(binding.LLVMFastMathNoNaNs)
	FastMathNoInfs          = FastMath(binding.LLVMFastMathNoInfs)
	FastMathNoSignedZeros   = FastMath(binding.LLVMFastMathNoSignedZeros)
	FastMathAllowReciprocal = FastMath(binding.LLVMFastMathAllowReciprocal)
	FastMathAllowContract   = FastMath(binding.LLVMFastMathAllowContract)
	FastMathApproxFunc      = FastMath(binding.LLVMFastMathApproxFunc)
	FastMathAll             = FastMath(binding.LLVMFastMathAll)
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

const (
	NoWrapNone     = NoWrap(0)
	NoWrapInBounds = NoWrap(binding.LLVMGEPFlagInBounds)
	NoWrapNUSW     = NoWrap(binding.LLVMGEPFlagNUSW)
	NoWrapNUW      = NoWrap(binding.LLVMGEPFlagNUW)
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

const (
	TailCallNone = TailCallKind(binding.LLVMTailCallKindNone)
	TailCallTail = TailCallKind(binding.LLVMTailCallKindTail)
	TailCallMust = TailCallKind(binding.LLVMTailCallKindMustTail)
	TailCallNo   = TailCallKind(binding.LLVMTailCallKindNoTail)
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

// SyncScopeOf 读取原子指令的 sync scope ID（0 = system）
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
