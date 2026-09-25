package ir

import (
	"strconv"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Op 指令操作码；转发 binding.LLVMOpcode（仿 llvm.Linkage 先例）
type Op binding.LLVMOpcode

const (
	OpRet            = Op(binding.LLVMRet)
	OpBr             = Op(binding.LLVMBr)
	OpSwitch         = Op(binding.LLVMSwitch)
	OpIndirectBr     = Op(binding.LLVMIndirectBr)
	OpInvoke         = Op(binding.LLVMInvoke)
	OpUnreachable    = Op(binding.LLVMUnreachable)
	OpFNeg           = Op(binding.LLVMFNeg)
	OpAdd            = Op(binding.LLVMAdd)
	OpFAdd           = Op(binding.LLVMFAdd)
	OpSub            = Op(binding.LLVMSub)
	OpFSub           = Op(binding.LLVMFSub)
	OpMul            = Op(binding.LLVMMul)
	OpFMul           = Op(binding.LLVMFMul)
	OpUDiv           = Op(binding.LLVMUDiv)
	OpSDiv           = Op(binding.LLVMSDiv)
	OpFDiv           = Op(binding.LLVMFDiv)
	OpURem           = Op(binding.LLVMURem)
	OpSRem           = Op(binding.LLVMSRem)
	OpFRem           = Op(binding.LLVMFRem)
	OpShl            = Op(binding.LLVMShl)
	OpLShr           = Op(binding.LLVMLShr)
	OpAShr           = Op(binding.LLVMAShr)
	OpAnd            = Op(binding.LLVMAnd)
	OpOr             = Op(binding.LLVMOr)
	OpXor            = Op(binding.LLVMXor)
	OpAlloca         = Op(binding.LLVMAlloca)
	OpLoad           = Op(binding.LLVMLoad)
	OpStore          = Op(binding.LLVMStore)
	OpGetElementPtr  = Op(binding.LLVMGetElementPtr)
	OpTrunc          = Op(binding.LLVMTrunc)
	OpZExt           = Op(binding.LLVMZExt)
	OpSExt           = Op(binding.LLVMSExt)
	OpFPToUI         = Op(binding.LLVMFPToUI)
	OpFPToSI         = Op(binding.LLVMFPToSI)
	OpUIToFP         = Op(binding.LLVMUIToFP)
	OpSIToFP         = Op(binding.LLVMSIToFP)
	OpFPTrunc        = Op(binding.LLVMFPTrunc)
	OpFPExt          = Op(binding.LLVMFPExt)
	OpPtrToInt       = Op(binding.LLVMPtrToInt)
	OpIntToPtr       = Op(binding.LLVMIntToPtr)
	OpBitCast        = Op(binding.LLVMBitCast)
	OpAddrSpaceCast  = Op(binding.LLVMAddrSpaceCast)
	OpICmp           = Op(binding.LLVMICmp)
	OpFCmp           = Op(binding.LLVMFCmp)
	OpPHI            = Op(binding.LLVMPHI)
	OpCall           = Op(binding.LLVMCall)
	OpSelect         = Op(binding.LLVMSelect)
	OpUserOp1        = Op(binding.LLVMUserOp1)
	OpUserOp2        = Op(binding.LLVMUserOp2)
	OpVAArg          = Op(binding.LLVMVAArg)
	OpExtractElement = Op(binding.LLVMExtractElement)
	OpInsertElement  = Op(binding.LLVMInsertElement)
	OpShuffleVector  = Op(binding.LLVMShuffleVector)
	OpExtractValue   = Op(binding.LLVMExtractValue)
	OpInsertValue    = Op(binding.LLVMInsertValue)
	OpFence          = Op(binding.LLVMFence)
	OpAtomicCmpXchg  = Op(binding.LLVMAtomicCmpXchg)
	OpAtomicRMW      = Op(binding.LLVMAtomicRMW)
	OpResume         = Op(binding.LLVMResume)
	OpLandingPad     = Op(binding.LLVMLandingPad)
	OpCleanupRet     = Op(binding.LLVMCleanupRet)
	OpCatchRet       = Op(binding.LLVMCatchRet)
	OpCatchPad       = Op(binding.LLVMCatchPad)
	OpCleanupPad     = Op(binding.LLVMCleanupPad)
	OpCatchSwitch    = Op(binding.LLVMCatchSwitch)
	OpCallBr         = Op(binding.LLVMCallBr)
	OpFreeze         = Op(binding.LLVMFreeze)
	OpPtrToAddr      = Op(binding.LLVMPtrToAddr)
)

// opNames 操作码短名（String 与调试输出用；未知码输出 op<N>）
var opNames = map[Op]string{
	OpRet: "ret", OpBr: "br", OpSwitch: "switch", OpIndirectBr: "indirectbr",
	OpInvoke: "invoke", OpUnreachable: "unreachable", OpFNeg: "fneg",
	OpAdd: "add", OpFAdd: "fadd", OpSub: "sub", OpFSub: "fsub",
	OpMul: "mul", OpFMul: "fmul", OpUDiv: "udiv", OpSDiv: "sdiv", OpFDiv: "fdiv",
	OpURem: "urem", OpSRem: "srem", OpFRem: "frem", OpShl: "shl",
	OpLShr: "lshr", OpAShr: "ashr", OpAnd: "and", OpOr: "or", OpXor: "xor",
	OpAlloca: "alloca", OpLoad: "load", OpStore: "store", OpGetElementPtr: "getelementptr",
	OpTrunc: "trunc", OpZExt: "zext", OpSExt: "sext", OpFPToUI: "fptoui",
	OpFPToSI: "fptosi", OpUIToFP: "uitofp", OpSIToFP: "sitofp", OpFPTrunc: "fptrunc",
	OpFPExt: "fpext", OpPtrToInt: "ptrtoint", OpIntToPtr: "inttoptr",
	OpBitCast: "bitcast", OpAddrSpaceCast: "addrspacecast", OpICmp: "icmp",
	OpFCmp: "fcmp", OpPHI: "phi", OpCall: "call", OpSelect: "select",
	OpUserOp1: "userop1", OpUserOp2: "userop2", OpVAArg: "va_arg",
	OpExtractElement: "extractelement", OpInsertElement: "insertelement",
	OpShuffleVector: "shufflevector", OpExtractValue: "extractvalue",
	OpInsertValue: "insertvalue", OpFence: "fence", OpAtomicCmpXchg: "cmpxchg",
	OpAtomicRMW: "atomicrmw", OpResume: "resume", OpLandingPad: "landingpad",
	OpCleanupRet: "cleanupret", OpCatchRet: "catchret", OpCatchPad: "catchpad",
	OpCleanupPad: "cleanuppad", OpCatchSwitch: "catchswitch", OpCallBr: "callbr",
	OpFreeze: "freeze", OpPtrToAddr: "ptrtoaddr",
}

// String 操作码名（小写，与 LLVM IR 文本一致）
func (o Op) String() string {
	if s, ok := opNames[o]; ok {
		return s
	}
	return "op" + strconv.Itoa(int(o))
}

// OpOf 指令的操作码；非指令值（常量/参数/全局等）返回 false
func OpOf(inst llvm.AnyValue) (Op, bool) {
	const op = "ir.OpOf"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	ref := inst.Ref()
	if binding.LLVMGetValueKind(ref) != binding.LLVMInstructionValueKind {
		return 0, false
	}
	return Op(binding.LLVMGetInstructionOpcode(ref)), true
}
