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
