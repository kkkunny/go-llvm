package binding

/*
#include "llvm-c/Core.h"
#include "Core.h"
*/
import "C"
import (
	"unsafe"
)

type LLVMOpcode int32

const (
	LLVMRet            LLVMOpcode = C.LLVMRet
	LLVMBr             LLVMOpcode = C.LLVMBr
	LLVMSwitch         LLVMOpcode = C.LLVMSwitch
	LLVMIndirectBr     LLVMOpcode = C.LLVMIndirectBr
	LLVMInvoke         LLVMOpcode = C.LLVMInvoke
	LLVMUnreachable    LLVMOpcode = C.LLVMUnreachable
	LLVMAdd            LLVMOpcode = C.LLVMAdd
	LLVMFAdd           LLVMOpcode = C.LLVMFAdd
	LLVMSub            LLVMOpcode = C.LLVMSub
	LLVMFSub           LLVMOpcode = C.LLVMFSub
	LLVMMul            LLVMOpcode = C.LLVMMul
	LLVMFMul           LLVMOpcode = C.LLVMFMul
	LLVMUDiv           LLVMOpcode = C.LLVMUDiv
	LLVMSDiv           LLVMOpcode = C.LLVMSDiv
	LLVMFDiv           LLVMOpcode = C.LLVMFDiv
	LLVMURem           LLVMOpcode = C.LLVMURem
	LLVMSRem           LLVMOpcode = C.LLVMSRem
	LLVMFRem           LLVMOpcode = C.LLVMFRem
	LLVMShl            LLVMOpcode = C.LLVMShl
	LLVMLShr           LLVMOpcode = C.LLVMLShr
	LLVMAShr           LLVMOpcode = C.LLVMAShr
	LLVMAnd            LLVMOpcode = C.LLVMAnd
	LLVMOr             LLVMOpcode = C.LLVMOr
	LLVMXor            LLVMOpcode = C.LLVMXor
	LLVMAlloca         LLVMOpcode = C.LLVMAlloca
	LLVMLoad           LLVMOpcode = C.LLVMLoad
	LLVMStore          LLVMOpcode = C.LLVMStore
	LLVMGetElementPtr  LLVMOpcode = C.LLVMGetElementPtr
	LLVMTrunc          LLVMOpcode = C.LLVMTrunc
	LLVMZExt           LLVMOpcode = C.LLVMZExt
	LLVMSExt           LLVMOpcode = C.LLVMSExt
	LLVMFPToUI         LLVMOpcode = C.LLVMFPToUI
	LLVMFPToSI         LLVMOpcode = C.LLVMFPToSI
	LLVMUIToFP         LLVMOpcode = C.LLVMUIToFP
	LLVMSIToFP         LLVMOpcode = C.LLVMSIToFP
	LLVMFPTrunc        LLVMOpcode = C.LLVMFPTrunc
	LLVMFPExt          LLVMOpcode = C.LLVMFPExt
	LLVMPtrToInt       LLVMOpcode = C.LLVMPtrToInt
	LLVMIntToPtr       LLVMOpcode = C.LLVMIntToPtr
	LLVMBitCast        LLVMOpcode = C.LLVMBitCast
	LLVMICmp           LLVMOpcode = C.LLVMICmp
	LLVMFCmp           LLVMOpcode = C.LLVMFCmp
	LLVMPHI            LLVMOpcode = C.LLVMPHI
	LLVMCall           LLVMOpcode = C.LLVMCall
	LLVMSelect         LLVMOpcode = C.LLVMSelect
	LLVMUserOp1        LLVMOpcode = C.LLVMUserOp1
	LLVMUserOp2        LLVMOpcode = C.LLVMUserOp2
	LLVMVAArg          LLVMOpcode = C.LLVMVAArg
	LLVMExtractElement LLVMOpcode = C.LLVMExtractElement
	LLVMInsertElement  LLVMOpcode = C.LLVMInsertElement
	LLVMShuffleVector  LLVMOpcode = C.LLVMShuffleVector
	LLVMExtractValue   LLVMOpcode = C.LLVMExtractValue
	LLVMInsertValue    LLVMOpcode = C.LLVMInsertValue
	LLVMFence          LLVMOpcode = C.LLVMFence
	LLVMAtomicCmpXchg  LLVMOpcode = C.LLVMAtomicCmpXchg
	LLVMAtomicRMW      LLVMOpcode = C.LLVMAtomicRMW
	LLVMResume         LLVMOpcode = C.LLVMResume
	LLVMLandingPad     LLVMOpcode = C.LLVMLandingPad
	LLVMAddrSpaceCast  LLVMOpcode = C.LLVMAddrSpaceCast
	LLVMCleanupRet     LLVMOpcode = C.LLVMCleanupRet
	LLVMCatchRet       LLVMOpcode = C.LLVMCatchRet
	LLVMCatchPad       LLVMOpcode = C.LLVMCatchPad
	LLVMCleanupPad     LLVMOpcode = C.LLVMCleanupPad
	LLVMCatchSwitch    LLVMOpcode = C.LLVMCatchSwitch
	LLVMFNeg           LLVMOpcode = C.LLVMFNeg
	LLVMCallBr         LLVMOpcode = C.LLVMCallBr
	LLVMFreeze         LLVMOpcode = C.LLVMFreeze
	LLVMPtrToAddr      LLVMOpcode = C.LLVMPtrToAddr
)

type LLVMTypeKind int32

const (
	// LLVMVoidTypeKind type with no size
	LLVMVoidTypeKind LLVMTypeKind = C.LLVMVoidTypeKind
	// LLVMHalfTypeKind 16 bit floating point type
	LLVMHalfTypeKind LLVMTypeKind = C.LLVMHalfTypeKind
	// LLVMFloatTypeKind 32 bit floating point type
	LLVMFloatTypeKind LLVMTypeKind = C.LLVMFloatTypeKind
	// LLVMDoubleTypeKind 64 bit floating point type
	LLVMDoubleTypeKind LLVMTypeKind = C.LLVMDoubleTypeKind
	// LLVMX86_FP80TypeKind 80 bit floating point type (X87)
	LLVMX86_FP80TypeKind LLVMTypeKind = C.LLVMX86_FP80TypeKind
	// LLVMFP128TypeKind 128 bit floating point type (112-bit mantissa)
	LLVMFP128TypeKind LLVMTypeKind = C.LLVMFP128TypeKind
	// LLVMPPC_FP128TypeKind 128 bit floating point type (two 64-bits)
	LLVMPPC_FP128TypeKind LLVMTypeKind = C.LLVMPPC_FP128TypeKind
	// LLVMLabelTypeKind Labels
	LLVMLabelTypeKind LLVMTypeKind = C.LLVMLabelTypeKind
	// LLVMIntegerTypeKind Arbitrary bit width integers
	LLVMIntegerTypeKind LLVMTypeKind = C.LLVMIntegerTypeKind
	// LLVMFunctionTypeKind Functions
	LLVMFunctionTypeKind LLVMTypeKind = C.LLVMFunctionTypeKind
	// LLVMStructTypeKind Structures
	LLVMStructTypeKind LLVMTypeKind = C.LLVMStructTypeKind
	// LLVMArrayTypeKind Arrays
	LLVMArrayTypeKind LLVMTypeKind = C.LLVMArrayTypeKind
	// LLVMPointerTypeKind Pointers
	LLVMPointerTypeKind LLVMTypeKind = C.LLVMPointerTypeKind
	// LLVMVectorTypeKind Fixed width SIMD vector type
	LLVMVectorTypeKind LLVMTypeKind = C.LLVMVectorTypeKind
	// LLVMMetadataTypeKind Metadata
	LLVMMetadataTypeKind LLVMTypeKind = C.LLVMMetadataTypeKind
	// LLVMTokenTypeKind Tokens
	LLVMTokenTypeKind LLVMTypeKind = C.LLVMTokenTypeKind
	// LLVMScalableVectorTypeKind Scalable SIMD vector type
	LLVMScalableVectorTypeKind LLVMTypeKind = C.LLVMScalableVectorTypeKind
	// LLVMBFloatTypeKind 16 bit brain floating point type
	LLVMBFloatTypeKind LLVMTypeKind = C.LLVMBFloatTypeKind
	// LLVMX86_AMXTypeKind X86 AMX
	LLVMX86_AMXTypeKind   LLVMTypeKind = C.LLVMX86_AMXTypeKind
	LLVMTargetExtTypeKind LLVMTypeKind = C.LLVMTargetExtTypeKind
)

type LLVMLinkage int32

const (
	// LLVMExternalLinkage Externally visible function
	LLVMExternalLinkage            LLVMLinkage = C.LLVMExternalLinkage
	LLVMAvailableExternallyLinkage LLVMLinkage = C.LLVMAvailableExternallyLinkage
	// LLVMLinkOnceAnyLinkage Keep one copy of function when linking (inline)
	LLVMLinkOnceAnyLinkage LLVMLinkage = C.LLVMLinkOnceAnyLinkage
	// LLVMLinkOnceODRLinkage Same, but only replaced by something equivalent.
	LLVMLinkOnceODRLinkage LLVMLinkage = C.LLVMLinkOnceODRLinkage
	// LLVMLinkOnceODRAutoHideLinkage Obsolete
	LLVMLinkOnceODRAutoHideLinkage LLVMLinkage = C.LLVMLinkOnceODRAutoHideLinkage
	// LLVMWeakAnyLinkage Keep one copy of function when linking (weak)
	LLVMWeakAnyLinkage LLVMLinkage = C.LLVMWeakAnyLinkage
	// LLVMWeakODRLinkage Same, but only replaced by something equivalent.
	LLVMWeakODRLinkage LLVMLinkage = C.LLVMWeakODRLinkage
	// LLVMAppendingLinkage Special purpose, only applies to global arrays
	LLVMAppendingLinkage LLVMLinkage = C.LLVMAppendingLinkage
	// LLVMInternalLinkage Rename collisions when linking (static functions)
	LLVMInternalLinkage LLVMLinkage = C.LLVMInternalLinkage
	// LLVMPrivateLinkage Like Internal, but omit from symbol table
	LLVMPrivateLinkage LLVMLinkage = C.LLVMPrivateLinkage
	// LLVMDLLImportLinkage Obsolete
	LLVMDLLImportLinkage LLVMLinkage = C.LLVMDLLImportLinkage
	// LLVMDLLExportLinkage Obsolete
	LLVMDLLExportLinkage LLVMLinkage = C.LLVMDLLExportLinkage
	// LLVMExternalWeakLinkage ExternalWeak linkage description
	LLVMExternalWeakLinkage LLVMLinkage = C.LLVMExternalWeakLinkage
	// LLVMGhostLinkage Obsolete
	LLVMGhostLinkage LLVMLinkage = C.LLVMGhostLinkage
	// LLVMCommonLinkage Tentative definitions
	LLVMCommonLinkage LLVMLinkage = C.LLVMCommonLinkage
	// LLVMLinkerPrivateLinkage Like Private, but linker removes.
	LLVMLinkerPrivateLinkage LLVMLinkage = C.LLVMLinkerPrivateLinkage
	// LLVMLinkerPrivateWeakLinkage Like LinkerPrivate, but is weak.
	LLVMLinkerPrivateWeakLinkage LLVMLinkage = C.LLVMLinkerPrivateWeakLinkage
)

type LLVMVisibility int32

const (
	// LLVMDefaultVisibility The GV is visible
	LLVMDefaultVisibility LLVMVisibility = C.LLVMDefaultVisibility
	// LLVMHiddenVisibility The GV is hidden
	LLVMHiddenVisibility LLVMVisibility = C.LLVMHiddenVisibility
	// LLVMProtectedVisibility The GV is protected
	LLVMProtectedVisibility LLVMVisibility = C.LLVMProtectedVisibility
)

type LLVMUnnamedAddr int32

const (
	// LLVMNoUnnamedAddr Address of the GV is significant.
	LLVMNoUnnamedAddr LLVMUnnamedAddr = C.LLVMNoUnnamedAddr
	// LLVMLocalUnnamedAddr Address of the GV is locally insignificant.
	LLVMLocalUnnamedAddr LLVMUnnamedAddr = C.LLVMLocalUnnamedAddr
	// LLVMGlobalUnnamedAddr Address of the GV is globally insignificant.
	LLVMGlobalUnnamedAddr LLVMUnnamedAddr = C.LLVMGlobalUnnamedAddr
)

type LLVMValueKind int32

const (
	LLVMArgumentValueKind              LLVMValueKind = C.LLVMArgumentValueKind
	LLVMBasicBlockValueKind            LLVMValueKind = C.LLVMBasicBlockValueKind
	LLVMMemoryUseValueKind             LLVMValueKind = C.LLVMMemoryUseValueKind
	LLVMMemoryDefValueKind             LLVMValueKind = C.LLVMMemoryDefValueKind
	LLVMMemoryPhiValueKind             LLVMValueKind = C.LLVMMemoryPhiValueKind
	LLVMFunctionValueKind              LLVMValueKind = C.LLVMFunctionValueKind
	LLVMGlobalAliasValueKind           LLVMValueKind = C.LLVMGlobalAliasValueKind
	LLVMGlobalIFuncValueKind           LLVMValueKind = C.LLVMGlobalIFuncValueKind
	LLVMGlobalVariableValueKind        LLVMValueKind = C.LLVMGlobalVariableValueKind
	LLVMBlockAddressValueKind          LLVMValueKind = C.LLVMBlockAddressValueKind
	LLVMConstantExprValueKind          LLVMValueKind = C.LLVMConstantExprValueKind
	LLVMConstantArrayValueKind         LLVMValueKind = C.LLVMConstantArrayValueKind
	LLVMConstantStructValueKind        LLVMValueKind = C.LLVMConstantStructValueKind
	LLVMConstantVectorValueKind        LLVMValueKind = C.LLVMConstantVectorValueKind
	LLVMUndefValueValueKind            LLVMValueKind = C.LLVMUndefValueValueKind
	LLVMConstantAggregateZeroValueKind LLVMValueKind = C.LLVMConstantAggregateZeroValueKind
	LLVMConstantDataArrayValueKind     LLVMValueKind = C.LLVMConstantDataArrayValueKind
	LLVMConstantDataVectorValueKind    LLVMValueKind = C.LLVMConstantDataVectorValueKind
	LLVMConstantIntValueKind           LLVMValueKind = C.LLVMConstantIntValueKind
	LLVMConstantFPValueKind            LLVMValueKind = C.LLVMConstantFPValueKind
	LLVMConstantPointerNullValueKind   LLVMValueKind = C.LLVMConstantPointerNullValueKind
	LLVMConstantTokenNoneValueKind     LLVMValueKind = C.LLVMConstantTokenNoneValueKind
	LLVMMetadataAsValueValueKind       LLVMValueKind = C.LLVMMetadataAsValueValueKind
	LLVMInlineAsmValueKind             LLVMValueKind = C.LLVMInlineAsmValueKind
	LLVMInstructionValueKind           LLVMValueKind = C.LLVMInstructionValueKind
	LLVMPoisonValueValueKind           LLVMValueKind = C.LLVMPoisonValueValueKind
	LLVMConstantTargetNoneValueKind    LLVMValueKind = C.LLVMConstantTargetNoneValueKind
	LLVMConstantPtrAuthValueKind       LLVMValueKind = C.LLVMConstantPtrAuthValueKind
)

type LLVMIntPredicate int32

const (
	// LLVMIntEQ equal
	LLVMIntEQ LLVMIntPredicate = C.LLVMIntEQ
	// LLVMIntNE not equal
	LLVMIntNE LLVMIntPredicate = C.LLVMIntNE
	// LLVMIntUGT unsigned greater than
	LLVMIntUGT LLVMIntPredicate = C.LLVMIntUGT
	// LLVMIntUGE unsigned greater or equal
	LLVMIntUGE LLVMIntPredicate = C.LLVMIntUGE
	// LLVMIntULT unsigned less than
	LLVMIntULT LLVMIntPredicate = C.LLVMIntULT
	// LLVMIntULE unsigned less or equal
	LLVMIntULE LLVMIntPredicate = C.LLVMIntULE
	// LLVMIntSGT signed greater than
	LLVMIntSGT LLVMIntPredicate = C.LLVMIntSGT
	// LLVMIntSGE signed greater or equal
	LLVMIntSGE LLVMIntPredicate = C.LLVMIntSGE
	// LLVMIntSLT signed less than
	LLVMIntSLT LLVMIntPredicate = C.LLVMIntSLT
	// LLVMIntSLE signed less or equal
	LLVMIntSLE LLVMIntPredicate = C.LLVMIntSLE
)

type LLVMRealPredicate int32

const (
	// LLVMRealPredicateFalse Always false (always folded)
	LLVMRealPredicateFalse LLVMRealPredicate = C.LLVMRealPredicateFalse
	// LLVMRealOEQ True if ordered and equal
	LLVMRealOEQ LLVMRealPredicate = C.LLVMRealOEQ
	// LLVMRealOGT True if ordered and greater than
	LLVMRealOGT LLVMRealPredicate = C.LLVMRealOGT
	// LLVMRealOGE True if ordered and greater than or equal
	LLVMRealOGE LLVMRealPredicate = C.LLVMRealOGE
	// LLVMRealOLT True if ordered and less than
	LLVMRealOLT LLVMRealPredicate = C.LLVMRealOLT
	// LLVMRealOLE True if ordered and less than or equal
	LLVMRealOLE LLVMRealPredicate = C.LLVMRealOLE
	// LLVMRealONE True if ordered and operands are unequal
	LLVMRealONE LLVMRealPredicate = C.LLVMRealONE
	// LLVMRealORD True if ordered (no nans)
	LLVMRealORD LLVMRealPredicate = C.LLVMRealORD
	// LLVMRealUNO True if unordered: isnan(X) | isnan(Y)
	LLVMRealUNO LLVMRealPredicate = C.LLVMRealUNO
	// LLVMRealUEQ True if unordered or equal
	LLVMRealUEQ LLVMRealPredicate = C.LLVMRealUEQ
	// LLVMRealUGT True if unordered or greater than
	LLVMRealUGT LLVMRealPredicate = C.LLVMRealUGT
	// LLVMRealUGE True if unordered, greater than, or equal
	LLVMRealUGE LLVMRealPredicate = C.LLVMRealUGE
	// LLVMRealULT True if unordered or less than
	LLVMRealULT LLVMRealPredicate = C.LLVMRealULT
	// LLVMRealULE True if unordered, less than, or equal
	LLVMRealULE LLVMRealPredicate = C.LLVMRealULE
	// LLVMRealUNE True if unordered or not equal
	LLVMRealUNE LLVMRealPredicate = C.LLVMRealUNE
	// LLVMRealPredicateTrue Always true (always folded)
	LLVMRealPredicateTrue LLVMRealPredicate = C.LLVMRealPredicateTrue
)

type LLVMThreadLocalMode int32

const (
	LLVMNotThreadLocal         LLVMThreadLocalMode = C.LLVMNotThreadLocal
	LLVMGeneralDynamicTLSModel LLVMThreadLocalMode = C.LLVMGeneralDynamicTLSModel
	LLVMLocalDynamicTLSModel   LLVMThreadLocalMode = C.LLVMLocalDynamicTLSModel
	LLVMInitialExecTLSModel    LLVMThreadLocalMode = C.LLVMInitialExecTLSModel
	LLVMLocalExecTLSModel      LLVMThreadLocalMode = C.LLVMLocalExecTLSModel
)

type LLVMDiagnosticSeverity int32

const (
	LLVMDSError   LLVMDiagnosticSeverity = C.LLVMDSError
	LLVMDSWarning LLVMDiagnosticSeverity = C.LLVMDSWarning
	LLVMDSRemark  LLVMDiagnosticSeverity = C.LLVMDSRemark
	LLVMDSNote    LLVMDiagnosticSeverity = C.LLVMDSNote
)

// LLVMAtomicOrdering is the ordering of a fence/atomic load-store instruction.
type LLVMAtomicOrdering int32

const (
	LLVMAtomicOrderingNotAtomic              LLVMAtomicOrdering = C.LLVMAtomicOrderingNotAtomic
	LLVMAtomicOrderingUnordered              LLVMAtomicOrdering = C.LLVMAtomicOrderingUnordered
	LLVMAtomicOrderingMonotonic              LLVMAtomicOrdering = C.LLVMAtomicOrderingMonotonic
	LLVMAtomicOrderingAcquire                LLVMAtomicOrdering = C.LLVMAtomicOrderingAcquire
	LLVMAtomicOrderingRelease                LLVMAtomicOrdering = C.LLVMAtomicOrderingRelease
	LLVMAtomicOrderingAcquireRelease         LLVMAtomicOrdering = C.LLVMAtomicOrderingAcquireRelease
	LLVMAtomicOrderingSequentiallyConsistent LLVMAtomicOrdering = C.LLVMAtomicOrderingSequentiallyConsistent
)

// LLVMAtomicRMWBinOp is the operation of an atomicrmw instruction.
type LLVMAtomicRMWBinOp int32

const (
	LLVMAtomicRMWBinOpXchg     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpXchg
	LLVMAtomicRMWBinOpAdd      LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpAdd
	LLVMAtomicRMWBinOpSub      LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpSub
	LLVMAtomicRMWBinOpAnd      LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpAnd
	LLVMAtomicRMWBinOpNand     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpNand
	LLVMAtomicRMWBinOpOr       LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpOr
	LLVMAtomicRMWBinOpXor      LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpXor
	LLVMAtomicRMWBinOpMax      LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpMax
	LLVMAtomicRMWBinOpMin      LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpMin
	LLVMAtomicRMWBinOpUMax     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpUMax
	LLVMAtomicRMWBinOpUMin     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpUMin
	LLVMAtomicRMWBinOpFAdd     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpFAdd
	LLVMAtomicRMWBinOpFSub     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpFSub
	LLVMAtomicRMWBinOpFMax     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpFMax
	LLVMAtomicRMWBinOpFMin     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpFMin
	LLVMAtomicRMWBinOpUIncWrap LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpUIncWrap
	LLVMAtomicRMWBinOpUDecWrap LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpUDecWrap
	LLVMAtomicRMWBinOpUSubCond LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpUSubCond
	LLVMAtomicRMWBinOpUSubSat  LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpUSubSat
	LLVMAtomicRMWBinOpFMaximum LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpFMaximum
	LLVMAtomicRMWBinOpFMinimum LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpFMinimum
)

type LLVMAttributeIndex int32

const (
	LLVMAttributeReturnIndex   LLVMAttributeIndex = C.LLVMAttributeReturnIndex
	LLVMAttributeFunctionIndex LLVMAttributeIndex = C.LLVMAttributeFunctionIndex
)

func LLVMDisposeMessage(message *C.char) {
	C.LLVMDisposeMessage(message)
}

type LLVMDiagnosticHandler func(LLVMDiagnosticInfoRef, unsafe.Pointer)

// LLVMContextCreate Create a new context.
// Every call to this function should be paired with a call to LLVMContextDispose() or the context will leak memory.
func LLVMContextCreate() LLVMContextRef {
	return LLVMContextRef{c: C.LLVMContextCreate()}
}

// LLVMContextSetDiagnosticHandler Set the diagnostic handler for this context.
func LLVMContextSetDiagnosticHandler(c LLVMContextRef, handler FuncPtr[LLVMDiagnosticHandler], diagnosticContext unsafe.Pointer) {
	C.LLVMContextSetDiagnosticHandler(c.c, (C.LLVMDiagnosticHandler)(handler.ptr), diagnosticContext)
}

// LLVMContextGetDiagnosticHandler Get the diagnostic handler of this context.
func LLVMContextGetDiagnosticHandler(c LLVMContextRef) FuncPtr[LLVMDiagnosticHandler] {
	return NewFuncPtr[LLVMDiagnosticHandler](unsafe.Pointer(C.LLVMContextGetDiagnosticHandler(c.c)))
}

// LLVMContextGetDiagnosticContext Get the diagnostic context of this context.
func LLVMContextGetDiagnosticContext(c LLVMContextRef) unsafe.Pointer {
	return C.LLVMContextGetDiagnosticContext(c.c)
}

// LLVMContextDispose Destroy a context instance.
// This should be called for every call to LLVMContextCreate() or memory will be leaked.
func LLVMContextDispose(c LLVMContextRef) {
	C.LLVMContextDispose(c.c)
}

// LLVMGetDiagInfoDescription Return a string representation of the DiagnosticInfo.
func LLVMGetDiagInfoDescription(di LLVMDiagnosticInfoRef) string {
	cstring := C.LLVMGetDiagInfoDescription(di.c)
	defer LLVMDisposeMessage(cstring)
	return C.GoString(cstring)
}

// LLVMGetDiagInfoSeverity Return an enum LLVMDiagnosticSeverity.
func LLVMGetDiagInfoSeverity(di LLVMDiagnosticInfoRef) LLVMDiagnosticSeverity {
	return LLVMDiagnosticSeverity(C.LLVMGetDiagInfoSeverity(di.c))
}

// LLVMGetEnumAttributeKindForName Return an unique id given the name of a enum attribute, or 0 if no attribute by that name exists.
// http://llvm.org/docs/LangRef.html#parameter-attributes
// http://llvm.org/docs/LangRef.html#function-attributes
func LLVMGetEnumAttributeKindForName(name string) uint32 {
	return string2CString(name, func(cname *C.char) uint32 {
		return uint32(C.LLVMGetEnumAttributeKindForName(cname, C.size_t(len(name))))
	})
}

// LLVMCreateEnumAttribute Create an enum attribute.
func LLVMCreateEnumAttribute(c LLVMContextRef, kindID uint32, val uint64) LLVMAttributeRef {
	return LLVMAttributeRef{c: C.LLVMCreateEnumAttribute(c.c, C.unsigned(kindID), C.uint64_t(val))}
}

// LLVMGetTypeByName Obtain a Type from a context by its registered name.
func LLVMGetTypeByName(c LLVMContextRef, name string) LLVMTypeRef {
	return string2CString(name, func(name *C.char) LLVMTypeRef {
		return LLVMTypeRef{c: C.LLVMGetTypeByName2(c.c, name)}
	})
}

// LLVMModuleCreateWithNameInContext Create a new, empty module in a specific context.
// Every invocation should be paired with LLVMDisposeModule() or memory will be leaked.
func LLVMModuleCreateWithNameInContext(moduleID string, c LLVMContextRef) LLVMModuleRef {
	return string2CString(moduleID, func(moduleID *C.char) LLVMModuleRef {
		return LLVMModuleRef{c: C.LLVMModuleCreateWithNameInContext(moduleID, c.c)}
	})
}

// LLVMCloneModule Return an exact copy of the specified module.
func LLVMCloneModule(m LLVMModuleRef) LLVMModuleRef {
	return LLVMModuleRef{c: C.LLVMCloneModule(m.c)}
}

// LLVMDisposeModule Destroy a module instance.
// This must be called for every created module or memory will be leaked.
func LLVMDisposeModule(m LLVMModuleRef) {
	C.LLVMDisposeModule(m.c)
}

// LLVMGetSourceFileName Obtain the module's original source file name.
func LLVMGetSourceFileName(m LLVMModuleRef) string {
	var length C.size_t
	return C.GoString(C.LLVMGetSourceFileName(m.c, &length))
}

// LLVMSetSourceFileName Set the original source file name of a module to a string Name with length Len.
func LLVMSetSourceFileName(m LLVMModuleRef, name string) {
	string2CString(name, func(cname *C.char) int {
		C.LLVMSetSourceFileName(m.c, cname, C.size_t(len(name)))
		return 1
	})
}

// LLVMGetDataLayoutStr Obtain the data layout for a module.
func LLVMGetDataLayoutStr(m LLVMModuleRef) string {
	return C.GoString(C.LLVMGetDataLayoutStr(m.c))
}

// LLVMSetDataLayout Set the data layout for a module.
func LLVMSetDataLayout(m LLVMModuleRef, name string) {
	string2CString(name, func(name *C.char) int {
		C.LLVMSetDataLayout(m.c, name)
		return 1
	})
}

// LLVMGetTarget Obtain the target triple for a module.
func LLVMGetTarget(m LLVMModuleRef) string {
	return C.GoString(C.LLVMGetTarget(m.c))
}

// LLVMSetTarget Set the target triple for a module.
func LLVMSetTarget(m LLVMModuleRef, triple string) {
	string2CString(triple, func(triple *C.char) int {
		C.LLVMSetTarget(m.c, triple)
		return 1
	})
}

// LLVMPrintModuleToString Return a string representation of the module.
func LLVMPrintModuleToString(m LLVMModuleRef) string {
	cstring := C.LLVMPrintModuleToString(m.c)
	defer LLVMDisposeMessage(cstring)
	return C.GoString(cstring)
}

// LLVMPrintModuleToFile Print a module to a file.
func LLVMPrintModuleToFile(m LLVMModuleRef, filename string) error {
	return string2CString(filename, func(filename *C.char) error {
		return llvmError2Error(func(errstr **C.char) C.LLVMBool {
			return C.LLVMPrintModuleToFile(m.c, filename, errstr)
		})
	})
}

// LLVMGetModuleContext Obtain the context to which this module is associated.
// @see Module::getContext()
func LLVMGetModuleContext(m LLVMModuleRef) LLVMContextRef {
	return LLVMContextRef{c: C.LLVMGetModuleContext(m.c)}
}

// LLVMAddFunction Add a function to a module under a specified name.
// @see llvm::Function::Create()
func LLVMAddFunction(m LLVMModuleRef, name string, functionTy LLVMTypeRef) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMAddFunction(m.c, name, functionTy.c)}
	})
}

// LLVMGetNamedFunction Obtain a Function value from a Module by its name.
// The returned value corresponds to a llvm::Function value.
// @see llvm::Module::getFunction()
func LLVMGetNamedFunction(m LLVMModuleRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMGetNamedFunction(m.c, name)}
	})
}

// LLVMGetFirstFunction Obtain the first Function in a Module.
// @see llvm::Module::begin()
func LLVMGetFirstFunction(m LLVMModuleRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetFirstFunction(m.c)}
}

// LLVMGetNextFunction Advance to the next Function in a Module.
// @see llvm::Module::iterator::operator++()
func LLVMGetNextFunction(fn LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetNextFunction(fn.c)}
}

// LLVMGetNamedFunctionWithLength Obtain a Function value from a Module by its name.
// The returned value corresponds to a llvm::Function value.
// @see llvm::Module::getFunction()
func LLVMGetNamedFunctionWithLength(m LLVMModuleRef, name string, length uint) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMGetNamedFunctionWithLength(m.c, name, C.size_t(length))}
	})
}

// LLVMGetTypeKind Obtain the enumerated type of a Type instance.
// @see llvm::Type:getTypeID()
func LLVMGetTypeKind(ty LLVMTypeRef) LLVMTypeKind {
	return LLVMTypeKind(C.LLVMGetTypeKind(ty.c))
}

// LLVMTypeIsSized Whether the type has a known size.
// Things that don't have a size are abstract types, labels, and void.a
// @see llvm::Type::isSized()
func LLVMTypeIsSized(ty LLVMTypeRef) bool {
	return llvmBool2bool(C.LLVMTypeIsSized(ty.c))
}

// LLVMGetTypeContext Obtain the context to which this type instance is associated.
// @see llvm::Type::getContext()
func LLVMGetTypeContext(ty LLVMTypeRef) LLVMContextRef {
	return LLVMContextRef{c: C.LLVMGetTypeContext(ty.c)}
}

// LLVMPrintTypeToString Return a string representation of the type.
// @see llvm::Type::print()
func LLVMPrintTypeToString(val LLVMTypeRef) string {
	cstring := C.LLVMPrintTypeToString(val.c)
	defer LLVMDisposeMessage(cstring)
	return C.GoString(cstring)
}

// LLVMIntTypeInContext Obtain an integer type from a context with specified bit width.
func LLVMIntTypeInContext(c LLVMContextRef, numBits uint32) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMIntTypeInContext(c.c, C.unsigned(numBits))}
}

// LLVMGetIntTypeWidth Obtain an integer type from the global context with a specified bit width.
func LLVMGetIntTypeWidth(integerTy LLVMTypeRef) uint32 {
	return uint32(C.LLVMGetIntTypeWidth(integerTy.c))
}

// LLVMHalfTypeInContext Obtain a 16-bit floating point type from a context.
func LLVMHalfTypeInContext(c LLVMContextRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMHalfTypeInContext(c.c)}
}

// LLVMBFloatTypeInContext Obtain a 16-bit brain floating point type from a context.
func LLVMBFloatTypeInContext(c LLVMContextRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMBFloatTypeInContext(c.c)}
}

// LLVMFloatTypeInContext Obtain a 32-bit floating point type from a context.
func LLVMFloatTypeInContext(c LLVMContextRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMFloatTypeInContext(c.c)}
}

// LLVMDoubleTypeInContext Obtain a 64-bit floating point type from a context.
func LLVMDoubleTypeInContext(c LLVMContextRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMDoubleTypeInContext(c.c)}
}

// LLVMX86FP80TypeInContext Obtain a 80-bit floating point type (X87) from a context.
func LLVMX86FP80TypeInContext(c LLVMContextRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMX86FP80TypeInContext(c.c)}
}

// LLVMFP128TypeInContext Obtain a 128-bit floating point type (112-bit mantissa) from a context.
func LLVMFP128TypeInContext(c LLVMContextRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMFP128TypeInContext(c.c)}
}

// LLVMPPCFP128TypeInContext Obtain a 128-bit floating point type (two 64-bits) from a context.
func LLVMPPCFP128TypeInContext(c LLVMContextRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMPPCFP128TypeInContext(c.c)}
}

// LLVMFunctionType Obtain a function type consisting of a specified signature.
// The function is defined as a tuple of a return Type, a list of parameter types, and whether the function is variadic.
func LLVMFunctionType(returnType LLVMTypeRef, paramTypes []LLVMTypeRef, isVarArg bool) LLVMTypeRef {
	ptr, length := slice2Ptr[LLVMTypeRef, C.LLVMTypeRef](paramTypes)
	return LLVMTypeRef{c: C.LLVMFunctionType(returnType.c, ptr, length, bool2LLVMBool(isVarArg))}
}

// LLVMIsFunctionVarArg Returns whether a function type is variadic.
func LLVMIsFunctionVarArg(functionTy LLVMTypeRef) bool {
	return llvmBool2bool(C.LLVMIsFunctionVarArg(functionTy.c))
}

// LLVMGetReturnType Obtain the Type this function Type returns.
func LLVMGetReturnType(functionTy LLVMTypeRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMGetReturnType(functionTy.c)}
}

// LLVMCountParamTypes Obtain the number of parameters this function accepts.
func LLVMCountParamTypes(functionTy LLVMTypeRef) uint32 {
	return uint32(C.LLVMCountParamTypes(functionTy.c))
}

// LLVMGetParamTypes Obtain the types of a function's parameters.
func LLVMGetParamTypes(functionTy LLVMTypeRef) []LLVMTypeRef {
	length := LLVMCountParamTypes(functionTy)
	dest := make([]LLVMTypeRef, length)
	ptr, _ := slice2Ptr[LLVMTypeRef, C.LLVMTypeRef](dest)
	C.LLVMGetParamTypes(functionTy.c, ptr)
	return dest
}

// LLVMStructTypeInContext Create a new structure type in a context.
// A structure is specified by a list of inner elements/types and whether these can be packed together.
// @see llvm::StructType::create()
func LLVMStructTypeInContext(c LLVMContextRef, elementTypes []LLVMTypeRef, packed bool) LLVMTypeRef {
	ptr, length := slice2Ptr[LLVMTypeRef, C.LLVMTypeRef](elementTypes)
	return LLVMTypeRef{c: C.LLVMStructTypeInContext(c.c, ptr, length, bool2LLVMBool(packed))}
}

// LLVMStructCreateNamed Create an empty structure in a context having a specified name.
// @see llvm::StructType::create()
func LLVMStructCreateNamed(c LLVMContextRef, name string) LLVMTypeRef {
	return string2CString(name, func(name *C.char) LLVMTypeRef {
		return LLVMTypeRef{c: C.LLVMStructCreateNamed(c.c, name)}
	})
}

// LLVMGetStructName Obtain the name of a structure.
// @see llvm::StructType::getName()
func LLVMGetStructName(ty LLVMTypeRef) string {
	return C.GoString(C.LLVMGetStructName(ty.c))
}

// LLVMStructSetBody Set the contents of a structure type.
// @see llvm::StructType::setBody()
func LLVMStructSetBody(structTy LLVMTypeRef, elementTypes []LLVMTypeRef, packed bool) {
	ptr, length := slice2Ptr[LLVMTypeRef, C.LLVMTypeRef](elementTypes)
	C.LLVMStructSetBody(structTy.c, ptr, length, bool2LLVMBool(packed))
}

// LLVMCountStructElementTypes Get the number of elements defined inside the structure.
// @see llvm::StructType::getNumElements()
func LLVMCountStructElementTypes(structTy LLVMTypeRef) uint32 {
	return uint32(C.LLVMCountStructElementTypes(structTy.c))
}

// LLVMGetStructElementTypes Get the elements within a structure.
func LLVMGetStructElementTypes(structTy LLVMTypeRef) []LLVMTypeRef {
	length := LLVMCountStructElementTypes(structTy)
	dest := make([]LLVMTypeRef, length)
	ptr, _ := slice2Ptr[LLVMTypeRef, C.LLVMTypeRef](dest)
	C.LLVMGetStructElementTypes(structTy.c, ptr)
	return dest
}

// LLVMStructGetTypeAtIndex Get the type of the element at a given index in the structure.
// @see llvm::StructType::getTypeAtIndex()
func LLVMStructGetTypeAtIndex(structTy LLVMTypeRef, i uint32) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMStructGetTypeAtIndex(structTy.c, C.unsigned(i))}
}

// LLVMIsPackedStruct Determine whether a structure is packed.
// @see llvm::StructType::isPacked()
func LLVMIsPackedStruct(structTy LLVMTypeRef) bool {
	return llvmBool2bool(C.LLVMIsPackedStruct(structTy.c))
}

// LLVMIsOpaqueStruct Determine whether a structure is opaque.
// @see llvm::StructType::isOpaque()
func LLVMIsOpaqueStruct(structTy LLVMTypeRef) bool {
	return llvmBool2bool(C.LLVMIsOpaqueStruct(structTy.c))
}

// LLVMGetElementType Obtain the element type of an array or vector type.
// @see llvm::SequentialType::getElementType()
func LLVMGetElementType(ty LLVMTypeRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMGetElementType(ty.c)}
}

// LLVMArrayType Create a fixed size array type that refers to a specific type.
// The created type will exist in the context that its element type exists in.
// @see llvm::ArrayType::get()
func LLVMArrayType(elementType LLVMTypeRef, elementCount uint32) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMArrayType(elementType.c, C.unsigned(elementCount))}
}

// LLVMGetArrayLength Obtain the length of an array type.
// This only works on types that represent arrays.
// @see llvm::ArrayType::getNumElements()
func LLVMGetArrayLength(arrayTy LLVMTypeRef) uint32 {
	return uint32(C.LLVMGetArrayLength(arrayTy.c))
}

// LLVMArrayType2 Create a fixed size array type that refers to a specific type.
// The created type will exist in the context that its element type exists in.
// @see llvm::ArrayType::get()
func LLVMArrayType2(elementType LLVMTypeRef, elementCount uint64) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMArrayType2(elementType.c, C.uint64_t(elementCount))}
}

// LLVMGetArrayLength2 Obtain the length of an array type.
// This only works on types that represent arrays.
// @see llvm::ArrayType::getNumElements()
func LLVMGetArrayLength2(arrayTy LLVMTypeRef) uint64 {
	return uint64(C.LLVMGetArrayLength2(arrayTy.c))
}

// LLVMVectorType Create a vector type that contains a defined type and has a specific number of elements.
// The created type will exist in the context thats its element type exists in.
// @see llvm::VectorType::get()
func LLVMVectorType(elementType LLVMTypeRef, elementCount uint32) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMVectorType(elementType.c, C.unsigned(elementCount))}
}

// LLVMGetVectorSize Obtain the number of elements in a vector type.
// This only works on types that represent vectors.
// @see llvm::VectorType::getNumElements()
func LLVMGetVectorSize(vectorTy LLVMTypeRef) uint32 {
	return uint32(C.LLVMGetVectorSize(vectorTy.c))
}

// LLVMPointerType Create a pointer type that points to a defined type.
// The created type will exist in the context that its pointee type exists in.
// @see llvm::PointerType::get()
func LLVMPointerType(elementType LLVMTypeRef, addressSpace uint32) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMPointerType(elementType.c, C.unsigned(addressSpace))}
}

// LLVMPointerTypeIsOpaque Determine whether a pointer is opaque.
// True if this is an instance of an opaque PointerType.
// @see llvm::Type::isOpaquePointerTy()
func LLVMPointerTypeIsOpaque(ty LLVMTypeRef) bool {
	return llvmBool2bool(C.LLVMPointerTypeIsOpaque(ty.c))
}

// LLVMPointerTypeInContext Create an opaque pointer type in a context.
// @see llvm::PointerType::get()
func LLVMPointerTypeInContext(c LLVMContextRef, addressSpace uint32) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMPointerTypeInContext(c.c, C.unsigned(addressSpace))}
}

// LLVMGetPointerAddressSpace Obtain the address space of a pointer type.
// This only works on types that represent pointers.
// @see llvm::PointerType::getAddressSpace()
func LLVMGetPointerAddressSpace(pointerTy LLVMTypeRef) uint32 {
	return uint32(C.LLVMGetPointerAddressSpace(pointerTy.c))
}

// LLVMVoidTypeInContext Create a void type in a context.
func LLVMVoidTypeInContext(c LLVMContextRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMVoidTypeInContext(c.c)}
}

// LLVMTypeOf Obtain the type of a value.
func LLVMTypeOf(val LLVMValueRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMTypeOf(val.c)}
}

// LLVMGetValueKind Obtain the enumerated type of a Value instance.
// @see llvm::Value::getValueID()
func LLVMGetValueKind(val LLVMValueRef) LLVMValueKind {
	return LLVMValueKind(C.LLVMGetValueKind(val.c))
}

// LLVMGetValueName Obtain the string name of a value.
func LLVMGetValueName(val LLVMValueRef) string {
	var length C.size_t
	s := C.LLVMGetValueName2(val.c, &length)
	return C.GoStringN(s, C.int(length))
}

// LLVMSetValueName Set the string name of a value.
func LLVMSetValueName(val LLVMValueRef, name string) {
	string2CString(name, func(cname *C.char) bool {
		C.LLVMSetValueName2(val.c, cname, C.size_t(len(name)))
		return false
	})
}

// LLVMDumpValue Dump a representation of a value to stderr.
func LLVMDumpValue(val LLVMValueRef) {
	C.LLVMDumpValue(val.c)
}

// LLVMPrintValueToString Return a string representation of the value.
func LLVMPrintValueToString(val LLVMValueRef) string {
	cstring := C.LLVMPrintValueToString(val.c)
	defer LLVMDisposeMessage(cstring)
	return C.GoString(cstring)
}

// LLVMIsConstant Determine whether the specified value instance is constant.
func LLVMIsConstant(val LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsConstant(val.c))
}

// LLVMConstNull Obtain a constant value referring to the null instance of a type.
// @see llvm::Constant::getNullValue()
func LLVMConstNull(ty LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstNull(ty.c)}
}

// LLVMIsNull Determine whether a value instance is null.
// @see llvm::Constant::isNullValue()
func LLVMIsNull(val LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsNull(val.c))
}

func LLVMConstAggregateZero(ty LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstAggregateZero(ty.c)}
}

// LLVMConstPointerNull Obtain a constant that is a constant pointer pointing to NULL for a specified type.
func LLVMConstPointerNull(ty LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstPointerNull(ty.c)}
}

// LLVMConstAllOnes Obtain a constant value for an integer type with all bits set.
func LLVMConstAllOnes(ty LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstAllOnes(ty.c)}
}

// LLVMConstInt Obtain a constant value for an integer type.
// The returned value corresponds to a llvm::ConstantInt.
// @see llvm::ConstantInt::get()
func LLVMConstInt(intTy LLVMTypeRef, n uint64, signExtend bool) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstInt(intTy.c, C.ulonglong(n), bool2LLVMBool(signExtend))}
}

// LLVMConstIntOfString Obtain a constant value for an integer parsed from a string.
// A similar API, LLVMConstIntOfStringAndSize is also available. If the string's length is available, it is preferred to call that function instead.
func LLVMConstIntOfString(intTy LLVMTypeRef, text string, radix uint8) LLVMValueRef {
	return string2CString(text, func(v *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMConstIntOfString(intTy.c, v, C.uint8_t(radix))}
	})
}

// LLVMConstReal Obtain a constant value referring to a double floating point value.
func LLVMConstReal(realTy LLVMTypeRef, n float64) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstReal(realTy.c, C.double(n))}
}

// LLVMConstRealOfString Obtain a constant for a floating point value parsed from a string.
// A similar API, LLVMConstRealOfStringAndSize is also available. It should be used if the input string's length is known.
func LLVMConstRealOfString(realTy LLVMTypeRef, text string) LLVMValueRef {
	return string2CString(text, func(v *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMConstRealOfString(realTy.c, v)}
	})
}

// LLVMConstIntGetZExtValue Obtain the zero extended value for an integer constant value.
func LLVMConstIntGetZExtValue(constantVal LLVMValueRef) uint64 {
	return uint64(C.LLVMConstIntGetZExtValue(constantVal.c))
}

// LLVMConstIntGetSExtValue Obtain the sign extended value for an integer constant value.
func LLVMConstIntGetSExtValue(constantVal LLVMValueRef) int64 {
	return int64(C.LLVMConstIntGetSExtValue(constantVal.c))
}

// LLVMConstRealGetDouble Obtain the double value for an floating point constant value.
func LLVMConstRealGetDouble(constantVal LLVMValueRef) (float64, bool) {
	var li C.LLVMBool
	v := float64(C.LLVMConstRealGetDouble(constantVal.c, &li))
	return v, llvmBool2bool(li)
}

// LLVMConstStringInContext Create a ConstantDataSequential and initialize it with a string.
// @see llvm::ConstantDataArray::getString()
func LLVMConstStringInContext(c LLVMContextRef, str string, dontNullTerminate bool) LLVMValueRef {
	return string2CString(str, func(v *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMConstStringInContext(c.c, v, C.unsigned(len(str)), bool2LLVMBool(dontNullTerminate))}
	})
}

// LLVMIsConstantString Returns true if the specified constant is an array of i8.
// @see ConstantDataSequential::getAsString()
func LLVMIsConstantString(c LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsConstantString(c.c))
}

// LLVMGetAsString Get the given constant data sequential as a string.
// @see ConstantDataSequential::getAsString()
func LLVMGetAsString(c LLVMValueRef) string {
	var length C.size_t
	return C.GoString(C.LLVMGetAsString(c.c, &length))
}

// LLVMConstStructInContext Create an anonymous ConstantStruct with the specified values.
// @see llvm::ConstantStruct::getAnon()
func LLVMConstStructInContext(c LLVMContextRef, constantVals []LLVMValueRef, packed bool) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](constantVals)
	return LLVMValueRef{c: C.LLVMConstStructInContext(c.c, ptr, length, bool2LLVMBool(packed))}
}

// LLVMConstArray Create a ConstantArray from values.
// @see llvm::ConstantArray::get()
func LLVMConstArray(elementTy LLVMTypeRef, constantVals []LLVMValueRef) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](constantVals)
	return LLVMValueRef{c: C.LLVMConstArray(elementTy.c, ptr, length)}
}

// LLVMConstNamedStruct Create a non-anonymous ConstantStruct from values.
// @see llvm::ConstantStruct::get()
func LLVMConstNamedStruct(structTy LLVMTypeRef, constantVals []LLVMValueRef) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](constantVals)
	return LLVMValueRef{c: C.LLVMConstNamedStruct(structTy.c, ptr, length)}
}

// LLVMConstVector Create a ConstantVector from values.
func LLVMConstVector(scalarConstantVals []LLVMValueRef) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](scalarConstantVals)
	return LLVMValueRef{c: C.LLVMConstVector(ptr, length)}
}

// LLVMGetAggregateElement Get element of a constant aggregate (struct, array or vector) at the specified index. Returns null if the index is out of range, or it's not possible to determine the element (e.g., because the constant is a constant expression.)
// @see llvm::Constant::getAggregateElement()
func LLVMGetAggregateElement(c LLVMValueRef, idx uint32) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetAggregateElement(c.c, C.unsigned(idx))}
}

// LLVMGetElementAsConstant an element at specified index as a constant.
// @see ConstantDataSequential::getElementAsConstant()
// Deprecated: Use LLVMGetAggregateElement instead
var LLVMGetElementAsConstant = LLVMGetAggregateElement

func LLVMGetConstOpcode(constantVal LLVMValueRef) LLVMOpcode {
	return LLVMOpcode(C.LLVMGetConstOpcode(constantVal.c))
}

func LLVMAlignOf(ty LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMAlignOf(ty.c)}
}

func LLVMSizeOf(ty LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMSizeOf(ty.c)}
}

func LLVMConstNeg(constantVal LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstNeg(constantVal.c)}
}

func LLVMConstNSWNeg(constantVal LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstNSWNeg(constantVal.c)}
}

func LLVMConstNot(constantVal LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstNot(constantVal.c)}
}

func LLVMConstAdd(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstAdd(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstNSWAdd(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstNSWAdd(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstNUWAdd(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstNUWAdd(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstSub(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstSub(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstNSWSub(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstNSWSub(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstNUWSub(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstNUWSub(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstXor(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstXor(lHSConstant.c, rHSConstant.c)}
}

// LLVMConstGEP2 Create a constant getelementptr expression with the given element type.
func LLVMConstGEP2(ty LLVMTypeRef, constantVal LLVMValueRef, constantIndices []LLVMValueRef) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](constantIndices)
	return LLVMValueRef{c: C.LLVMConstGEP2(ty.c, constantVal.c, ptr, length)}
}

// LLVMConstInBoundsGEP2 Create a constant inbounds getelementptr expression with the given element type.
func LLVMConstInBoundsGEP2(ty LLVMTypeRef, constantVal LLVMValueRef, constantIndices []LLVMValueRef) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](constantIndices)
	return LLVMValueRef{c: C.LLVMConstInBoundsGEP2(ty.c, constantVal.c, ptr, length)}
}

func LLVMConstTrunc(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstTrunc(constantVal.c, toType.c)}
}

func LLVMConstPtrToInt(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstPtrToInt(constantVal.c, toType.c)}
}

func LLVMConstIntToPtr(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstIntToPtr(constantVal.c, toType.c)}
}

func LLVMConstBitCast(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstBitCast(constantVal.c, toType.c)}
}

func LLVMConstAddrSpaceCast(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstAddrSpaceCast(constantVal.c, toType.c)}
}

func LLVMConstTruncOrBitCast(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstTruncOrBitCast(constantVal.c, toType.c)}
}

func LLVMConstPointerCast(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstPointerCast(constantVal.c, toType.c)}
}

func LLVMConstExtractElement(vectorConstant, indexConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstExtractElement(vectorConstant.c, indexConstant.c)}
}

func LLVMConstInsertElement(vectorConstant, elementValueConstant, indexConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstInsertElement(vectorConstant.c, elementValueConstant.c, indexConstant.c)}
}

func LLVMConstShuffleVector(vectorAConstant, vectorBConstant, maskConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstShuffleVector(vectorAConstant.c, vectorBConstant.c, maskConstant.c)}
}

// LLVMConstInlineAsm
// Deprecated: Use LLVMGetInlineAsm instead.
func LLVMConstInlineAsm(ty LLVMTypeRef, asmString, constraints string, hasSideEffects bool, isAlignStack bool) LLVMValueRef {
	return string2CString(asmString, func(asmString *C.char) LLVMValueRef {
		return string2CString(constraints, func(constraints *C.char) LLVMValueRef {
			return LLVMValueRef{c: C.LLVMConstInlineAsm(ty.c, asmString, constraints, bool2LLVMBool(hasSideEffects), bool2LLVMBool(isAlignStack))}
		})
	})
}

func LLVMGetGlobalParent(global LLVMValueRef) LLVMModuleRef {
	return LLVMModuleRef{c: C.LLVMGetGlobalParent(global.c)}
}

func LLVMIsDeclaration(global LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsDeclaration(global.c))
}

func LLVMGetLinkage(global LLVMValueRef) LLVMLinkage {
	return LLVMLinkage(C.LLVMGetLinkage(global.c))
}

func LLVMSetLinkage(global LLVMValueRef, linkage LLVMLinkage) {
	C.LLVMSetLinkage(global.c, C.LLVMLinkage(linkage))
}

func LLVMGetVisibility(global LLVMValueRef) LLVMVisibility {
	return LLVMVisibility(C.LLVMGetVisibility(global.c))
}

func LLVMSetVisibility(global LLVMValueRef, viz LLVMVisibility) {
	C.LLVMSetVisibility(global.c, C.LLVMVisibility(viz))
}

func LLVMGetUnnamedAddress(global LLVMValueRef) LLVMUnnamedAddr {
	return LLVMUnnamedAddr(C.LLVMGetUnnamedAddress(global.c))
}

func LLVMSetUnnamedAddress(global LLVMValueRef, unnamedAddr LLVMUnnamedAddr) {
	C.LLVMSetUnnamedAddress(global.c, C.LLVMUnnamedAddr(unnamedAddr))
}

// LLVMGlobalGetValueType Returns the "value type" of a global value.  This differs from the formal type of a global value which is always a pointer type.
// @see llvm::GlobalValue::getValueType()
func LLVMGlobalGetValueType(global LLVMValueRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMGlobalGetValueType(global.c)}
}

// LLVMGetAlignment Obtain the preferred alignment of the value.
// @see llvm::AllocaInst::getAlignment()
// @see llvm::LoadInst::getAlignment()
// @see llvm::StoreInst::getAlignment()
// @see llvm::AtomicRMWInst::setAlignment()
// @see llvm::AtomicCmpXchgInst::setAlignment()
// @see llvm::GlobalValue::getAlignment()
func LLVMGetAlignment(v LLVMValueRef) uint32 {
	return uint32(C.LLVMGetAlignment(v.c))
}

// LLVMSetAlignment Set the preferred alignment of the value.
// @see llvm::AllocaInst::setAlignment()
// @see llvm::LoadInst::setAlignment()
// @see llvm::StoreInst::setAlignment()
// @see llvm::AtomicRMWInst::setAlignment()
// @see llvm::AtomicCmpXchgInst::setAlignment()
// @see llvm::GlobalValue::setAlignment()
func LLVMSetAlignment(v LLVMValueRef, bytes uint32) {
	C.LLVMSetAlignment(v.c, C.unsigned(bytes))
}

func LLVMAddGlobal(m LLVMModuleRef, ty LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMAddGlobal(m.c, ty.c, name)}
	})
}

func LLVMAddGlobalInAddressSpace(m LLVMModuleRef, ty LLVMTypeRef, name string, addressSpace uint32) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMAddGlobalInAddressSpace(m.c, ty.c, name, C.unsigned(addressSpace))}
	})
}

func LLVMGetNamedGlobal(m LLVMModuleRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMGetNamedGlobal(m.c, name)}
	})
}

func LLVMGetNamedGlobalWithLength(m LLVMModuleRef, name string, length uint) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMGetNamedGlobalWithLength(m.c, name, C.size_t(length))}
	})
}

func LLVMDeleteGlobal(globalVar LLVMValueRef) {
	C.LLVMDeleteGlobal(globalVar.c)
}

func LLVMGetInitializer(globalVar LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetInitializer(globalVar.c)}
}

func LLVMSetInitializer(globalVar LLVMValueRef, constantVal LLVMValueRef) {
	C.LLVMSetInitializer(globalVar.c, constantVal.c)
}

func LLVMIsThreadLocal(globalVar LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsThreadLocal(globalVar.c))
}

func LLVMSetThreadLocal(globalVar LLVMValueRef, isThreadLocal bool) {
	C.LLVMSetThreadLocal(globalVar.c, bool2LLVMBool(isThreadLocal))
}

func LLVMIsDSOLocal(globalVar LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsDSOLocal(globalVar.c))
}

func LLVMSetDSOLocal(globalVar LLVMValueRef, local bool) {
	C.LLVMSetDSOLocal(globalVar.c, bool2LLVMBool(local))
}

func LLVMIsGlobalConstant(globalVar LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsGlobalConstant(globalVar.c))
}

func LLVMSetGlobalConstant(globalVar LLVMValueRef, isConstant bool) {
	C.LLVMSetGlobalConstant(globalVar.c, bool2LLVMBool(isConstant))
}

func LLVMGetThreadLocalMode(globalVar LLVMValueRef) LLVMThreadLocalMode {
	return LLVMThreadLocalMode(C.LLVMGetThreadLocalMode(globalVar.c))
}

func LLVMSetThreadLocalMode(globalVar LLVMValueRef, mode LLVMThreadLocalMode) {
	C.LLVMSetThreadLocalMode(globalVar.c, C.LLVMThreadLocalMode(mode))
}

func LLVMIsExternallyInitialized(globalVar LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsExternallyInitialized(globalVar.c))
}

func LLVMSetExternallyInitialized(globalVar LLVMValueRef, isExtInit bool) {
	C.LLVMSetExternallyInitialized(globalVar.c, bool2LLVMBool(isExtInit))
}

// LLVMAddAttributeAtIndex Add an attribute to a function.
func LLVMAddAttributeAtIndex(f LLVMValueRef, idx LLVMAttributeIndex, a LLVMAttributeRef) {
	C.LLVMAddAttributeAtIndex(f.c, C.LLVMAttributeIndex(idx), a.c)
}

// LLVMCountParams Obtain the number of parameters in a function.
func LLVMCountParams(fn LLVMValueRef) uint32 {
	return uint32(C.LLVMCountParams(fn.c))
}

// LLVMGetParams Obtain the parameters in a function.
func LLVMGetParams(fn LLVMValueRef) []LLVMValueRef {
	length := LLVMCountParams(fn)
	params := make([]LLVMValueRef, length)
	ptr, _ := slice2Ptr[LLVMValueRef, C.LLVMValueRef](params)
	C.LLVMGetParams(fn.c, ptr)
	return params
}

// LLVMGetParam Obtain the parameter at the specified index.
// Parameters are indexed from 0.
func LLVMGetParam(fn LLVMValueRef, index uint32) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetParam(fn.c, C.unsigned(index))}
}

// LLVMGetParamParent Obtain the function to which this argument belongs.
// Unlike other functions in this group, this one takes an LLVMValueRef that corresponds to a llvm::Attribute.
// The returned LLVMValueRef is the llvm::Function to which this argument belongs.
func LLVMGetParamParent(fn LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetParamParent(fn.c)}
}

// LLVMGetFirstParam Obtain the first parameter to a function.
func LLVMGetFirstParam(fn LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetFirstParam(fn.c)}
}

// LLVMGetLastParam Obtain the last parameter to a function.
func LLVMGetLastParam(fn LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetLastParam(fn.c)}
}

// LLVMGetNextParam Obtain the next parameter to a function.
// This takes an LLVMValueRef obtained from LLVMGetFirstParam() (which is actually a wrapped iterator) and obtains the next parameter from the underlying iterator.
func LLVMGetNextParam(arg LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetNextParam(arg.c)}
}

// LLVMGetPreviousParam Obtain the previous parameter to a function.
// This is the opposite of LLVMGetNextParam().
func LLVMGetPreviousParam(arg LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetPreviousParam(arg.c)}
}

// LLVMSetParamAlignment Set the alignment for a function parameter.
func LLVMSetParamAlignment(arg LLVMValueRef, align uint32) {
	C.LLVMSetParamAlignment(arg.c, C.unsigned(align))
}

// LLVMGetBasicBlockName Obtain the string name of a basic block.
func LLVMGetBasicBlockName(bb LLVMBasicBlockRef) string {
	return C.GoString(C.LLVMGetBasicBlockName(bb.c))
}

// LLVMGetBasicBlockParent Obtain the function to which a basic block belongs.
func LLVMGetBasicBlockParent(bb LLVMBasicBlockRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetBasicBlockParent(bb.c)}
}

// LLVMGetBasicBlockTerminator Obtain the terminator instruction for a basic block.
// If the basic block does not have a terminator (it is not well-formed if it doesn't), then NULL is returned.
// The returned LLVMValueRef corresponds to an llvm::Instruction.
func LLVMGetBasicBlockTerminator(bb LLVMBasicBlockRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetBasicBlockTerminator(bb.c)}
}

// LLVMCountBasicBlocks Obtain the number of basic blocks in a function.
func LLVMCountBasicBlocks(fn LLVMValueRef) uint32 {
	return uint32(C.LLVMCountBasicBlocks(fn.c))
}

// LLVMGetBasicBlocks Obtain all of the basic blocks in a function.
func LLVMGetBasicBlocks(fn LLVMValueRef) []LLVMBasicBlockRef {
	length := LLVMCountBasicBlocks(fn)
	basicBlocks := make([]LLVMBasicBlockRef, length)
	ptr, _ := slice2Ptr[LLVMBasicBlockRef, C.LLVMBasicBlockRef](basicBlocks)
	C.LLVMGetBasicBlocks(fn.c, ptr)
	return basicBlocks
}

// LLVMGetFirstBasicBlock Obtain the first basic block in a function.
// The returned basic block can be used as an iterator. You will likely eventually call into LLVMGetNextBasicBlock() with it.
func LLVMGetFirstBasicBlock(fn LLVMValueRef) LLVMBasicBlockRef {
	return LLVMBasicBlockRef{c: C.LLVMGetFirstBasicBlock(fn.c)}
}

// LLVMGetLastBasicBlock Obtain the last basic block in a function.
func LLVMGetLastBasicBlock(fn LLVMValueRef) LLVMBasicBlockRef {
	return LLVMBasicBlockRef{c: C.LLVMGetLastBasicBlock(fn.c)}
}

// LLVMGetNextBasicBlock Advance a basic block iterator.
func LLVMGetNextBasicBlock(bb LLVMBasicBlockRef) LLVMBasicBlockRef {
	return LLVMBasicBlockRef{c: C.LLVMGetNextBasicBlock(bb.c)}
}

// LLVMGetPreviousBasicBlock Go backwards in a basic block iterator.
func LLVMGetPreviousBasicBlock(bb LLVMBasicBlockRef) LLVMBasicBlockRef {
	return LLVMBasicBlockRef{c: C.LLVMGetPreviousBasicBlock(bb.c)}
}

// LLVMGetEntryBasicBlock Obtain the basic block that corresponds to the entry point of a function.
func LLVMGetEntryBasicBlock(fn LLVMValueRef) LLVMBasicBlockRef {
	return LLVMBasicBlockRef{c: C.LLVMGetEntryBasicBlock(fn.c)}
}

// LLVMBasicBlockAsValue Convert a basic block instance to a value type.
func LLVMBasicBlockAsValue(bb LLVMBasicBlockRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBasicBlockAsValue(bb.c)}
}

// LLVMValueAsBasicBlock Convert a value instance to a basic block, if it is one.
func LLVMValueAsBasicBlock(val LLVMValueRef) LLVMBasicBlockRef {
	return LLVMBasicBlockRef{c: C.LLVMValueAsBasicBlock(val.c)}
}

// LLVMAppendBasicBlockInContext Append a basic block to the end of a function.
func LLVMAppendBasicBlockInContext(c LLVMContextRef, fn LLVMValueRef, name string) LLVMBasicBlockRef {
	return string2CString(name, func(name *C.char) LLVMBasicBlockRef {
		return LLVMBasicBlockRef{c: C.LLVMAppendBasicBlockInContext(c.c, fn.c, name)}
	})
}

// LLVMDeleteBasicBlock Remove a basic block from a function and delete it.
// This deletes the basic block from its containing function and deletes the basic block itself.
func LLVMDeleteBasicBlock(bb LLVMBasicBlockRef) {
	C.LLVMDeleteBasicBlock(bb.c)
}

// LLVMRemoveBasicBlockFromParent Remove a basic block from a function.
// This deletes the basic block from its containing function but keep the basic block alive.
func LLVMRemoveBasicBlockFromParent(bb LLVMBasicBlockRef) {
	C.LLVMRemoveBasicBlockFromParent(bb.c)
}

// LLVMGetFirstInstruction Obtain the first instruction in a basic block.
// The returned LLVMValueRef corresponds to a llvm::Instruction instance.
func LLVMGetFirstInstruction(bb LLVMBasicBlockRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetFirstInstruction(bb.c)}
}

// LLVMGetLastInstruction Obtain the last instruction in a basic block.
// The returned LLVMValueRef corresponds to an LLVM:Instruction.
func LLVMGetLastInstruction(bb LLVMBasicBlockRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetLastInstruction(bb.c)}
}

// LLVMGetInstructionParent Obtain the basic block to which an instruction belongs.
func LLVMGetInstructionParent(inst LLVMValueRef) LLVMBasicBlockRef {
	return LLVMBasicBlockRef{c: C.LLVMGetInstructionParent(inst.c)}
}

// LLVMGetNextInstruction Obtain the instruction that occurs after the one specified.
// The next instruction will be from the same basic block.
// If this is the last instruction in a basic block, NULL will be returned.
func LLVMGetNextInstruction(inst LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetNextInstruction(inst.c)}
}

// LLVMGetPreviousInstruction Obtain the instruction that occurred before this one.
// If the instruction is the first instruction in a basic block, NULL will be returned.
func LLVMGetPreviousInstruction(inst LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetPreviousInstruction(inst.c)}
}

// LLVMInstructionRemoveFromParent Remove an instruction.
// The instruction specified is removed from its containing building block but is kept alive.
func LLVMInstructionRemoveFromParent(inst LLVMValueRef) {
	C.LLVMInstructionRemoveFromParent(inst.c)
}

// LLVMDeleteInstruction Delete an instruction.
// The instruction specified is deleted. It must have previously been removed from its containing building block.
func LLVMDeleteInstruction(inst LLVMValueRef) {
	C.LLVMDeleteInstruction(inst.c)
}

// LLVMGetInstructionOpcode Obtain the code opcode for an individual instruction.
func LLVMGetInstructionOpcode(inst LLVMValueRef) LLVMOpcode {
	return LLVMOpcode(C.LLVMGetInstructionOpcode(inst.c))
}

// LLVMGetICmpPredicate Obtain the predicate of an instruction.
// This is only valid for instructions that correspond to llvm::ICmpInst or llvm::ConstantExpr whose opcode is llvm::Instruction::ICmp.
func LLVMGetICmpPredicate(inst LLVMValueRef) LLVMIntPredicate {
	return LLVMIntPredicate(C.LLVMGetICmpPredicate(inst.c))
}

// LLVMGetFCmpPredicate Obtain the float predicate of an instruction.
// This is only valid for instructions that correspond to llvm::FCmpInst or llvm::ConstantExpr whose opcode is llvm::Instruction::FCmp.
func LLVMGetFCmpPredicate(inst LLVMValueRef) LLVMRealPredicate {
	return LLVMRealPredicate(C.LLVMGetFCmpPredicate(inst.c))
}

// LLVMInstructionClone Create a copy of 'this' instruction that is identical in all ways
// except the following:
//   - The instruction has no parent
//   - The instruction has no name
func LLVMInstructionClone(inst LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMInstructionClone(inst.c)}
}

// LLVMIsATerminatorInst Determine whether an instruction is a terminator.
// This routine is named to be compatible with historical functions that did this by querying the underlying C++ type.
func LLVMIsATerminatorInst(inst LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMIsATerminatorInst(inst.c)}
}

// LLVMGetCalledFunctionType Obtain the function type called by this instruction.
func LLVMGetCalledFunctionType(c LLVMValueRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMGetCalledFunctionType(c.c)}
}

func LLVMGetFunctionType(c LLVMValueRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMGetFunctionType(c.c)}
}

// LLVMAddIncoming Add an incoming value to the end of a PHI list.
func LLVMAddIncoming(phiNode LLVMValueRef, incomingValues []LLVMValueRef, incomingBlocks []LLVMBasicBlockRef) {
	ptr1, count := slice2Ptr[LLVMValueRef, C.LLVMValueRef](incomingValues)
	ptr2, count := slice2Ptr[LLVMBasicBlockRef, C.LLVMBasicBlockRef](incomingBlocks)
	C.LLVMAddIncoming(phiNode.c, ptr1, ptr2, count)
}

// LLVMCountIncoming Obtain the number of incoming basic blocks to a PHI node.
func LLVMCountIncoming(phiNode LLVMValueRef) uint32 {
	return uint32(C.LLVMCountIncoming(phiNode.c))
}

// LLVMGetIncomingValue Obtain an incoming value to a PHI node as an LLVMValueRef.
func LLVMGetIncomingValue(phiNode LLVMValueRef, index uint32) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetIncomingValue(phiNode.c, C.unsigned(index))}
}

// LLVMGetIncomingBlock Obtain an incoming value to a PHI node as an LLVMValueRef.
func LLVMGetIncomingBlock(phiNode LLVMValueRef, index uint32) LLVMBasicBlockRef {
	return LLVMBasicBlockRef{c: C.LLVMGetIncomingBlock(phiNode.c, C.unsigned(index))}
}

func LLVMCreateBuilderInContext(c LLVMContextRef) LLVMBuilderRef {
	return LLVMBuilderRef{c: C.LLVMCreateBuilderInContext(c.c)}
}

func LLVMPositionBuilderBefore(builder LLVMBuilderRef, instr LLVMValueRef) {
	C.LLVMPositionBuilderBefore(builder.c, instr.c)
}

func LLVMPositionBuilderAtEnd(builder LLVMBuilderRef, block LLVMBasicBlockRef) {
	C.LLVMPositionBuilderAtEnd(builder.c, block.c)
}

func LLVMGetInsertBlock(builder LLVMBuilderRef) LLVMBasicBlockRef {
	return LLVMBasicBlockRef{c: C.LLVMGetInsertBlock(builder.c)}
}

func LLVMDisposeBuilder(builder LLVMBuilderRef) {
	C.LLVMDisposeBuilder(builder.c)
}

func LLVMBuildRetVoid(builder LLVMBuilderRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildRetVoid(builder.c)}
}

func LLVMBuildRet(builder LLVMBuilderRef, v LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildRet(builder.c, v.c)}
}

func LLVMBuildBr(builder LLVMBuilderRef, dest LLVMBasicBlockRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildBr(builder.c, dest.c)}
}

func LLVMBuildCondBr(builder LLVMBuilderRef, ifv LLVMValueRef, thenb, elseb LLVMBasicBlockRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildCondBr(builder.c, ifv.c, thenb.c, elseb.c)}
}

func LLVMBuildSwitch(builder LLVMBuilderRef, v LLVMValueRef, elseBlock LLVMBasicBlockRef, numCases uint32) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildSwitch(builder.c, v.c, elseBlock.c, C.unsigned(numCases))}
}

func LLVMBuildInvoke(builder LLVMBuilderRef, ty LLVMTypeRef, fn LLVMValueRef, args []LLVMValueRef, then, catch LLVMBasicBlockRef, name string) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](args)
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildInvoke2(builder.c, ty.c, fn.c, ptr, length, then.c, catch.c, name)}
	})
}

func LLVMBuildResume(builder LLVMBuilderRef, ext LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildResume(builder.c, ext.c)}
}

func LLVMBuildLandingPad(builder LLVMBuilderRef, ty LLVMTypeRef, persFn LLVMValueRef, numClauses uint32, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildLandingPad(builder.c, ty.c, persFn.c, C.unsigned(numClauses), name)}
	})
}

func LLVMBuildCleanupRet(builder LLVMBuilderRef, catchPad LLVMValueRef, bb LLVMBasicBlockRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildCleanupRet(builder.c, catchPad.c, bb.c)}
}

func LLVMBuildCatchRet(builder LLVMBuilderRef, catchPad LLVMValueRef, bb LLVMBasicBlockRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildCatchRet(builder.c, catchPad.c, bb.c)}
}

func LLVMBuildCatchPad(builder LLVMBuilderRef, parentPad LLVMValueRef, args []LLVMValueRef, name string) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](args)
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildCatchPad(builder.c, parentPad.c, ptr, length, name)}
	})
}

func LLVMBuildCleanupPad(builder LLVMBuilderRef, parentPad LLVMValueRef, args []LLVMValueRef, name string) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](args)
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildCleanupPad(builder.c, parentPad.c, ptr, length, name)}
	})
}

func LLVMBuildCatchSwitch(builder LLVMBuilderRef, parentPad LLVMValueRef, unwindBB LLVMBasicBlockRef, numHandlers uint32, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildCatchSwitch(builder.c, parentPad.c, unwindBB.c, C.unsigned(numHandlers), name)}
	})
}

func LLVMAddCase(v LLVMValueRef, onVal LLVMValueRef, dest LLVMBasicBlockRef) {
	C.LLVMAddCase(v.c, onVal.c, dest.c)
}

// LLVMGetNumClauses Get the number of clauses on the landingpad instruction.
func LLVMGetNumClauses(landingPad LLVMValueRef) uint32 {
	return uint32(C.LLVMGetNumClauses(landingPad.c))
}

// LLVMGetClause Get the value of the clause at index Idx on the landingpad instruction.
func LLVMGetClause(landingPad LLVMValueRef, idx uint32) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetClause(landingPad.c, C.unsigned(idx))}
}

// LLVMAddClause Add a catch or filter clause to the landingpad instruction.
func LLVMAddClause(landingPad, clauseVal LLVMValueRef) {
	C.LLVMAddClause(landingPad.c, clauseVal.c)
}

// LLVMIsCleanup Get the 'cleanup' flag in the landingpad instruction.
func LLVMIsCleanup(landingPad LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsCleanup(landingPad.c))
}

// LLVMSetCleanup Set the 'cleanup' flag in the landingpad instruction.
func LLVMSetCleanup(landingPad LLVMValueRef, val bool) {
	C.LLVMSetCleanup(landingPad.c, bool2LLVMBool(val))
}

// LLVMAddHandler Add a destination to the catchswitch instruction.
func LLVMAddHandler(catchSwitch LLVMValueRef, dest LLVMBasicBlockRef) {
	C.LLVMAddHandler(catchSwitch.c, dest.c)
}

// LLVMGetNumHandlers Get the number of handlers on the catchswitch instruction.
func LLVMGetNumHandlers(catchSwitch LLVMValueRef) uint32 {
	return uint32(C.LLVMGetNumHandlers(catchSwitch.c))
}

// LLVMGetHandlers Obtain the basic blocks acting as handlers for a catchswitch instruction.
func LLVMGetHandlers(catchSwitch LLVMValueRef) []LLVMBasicBlockRef {
	n := int(LLVMGetNumHandlers(catchSwitch))
	if n == 0 {
		return nil
	}
	handlers := make([]C.LLVMBasicBlockRef, n)
	C.LLVMGetHandlers(catchSwitch.c, &handlers[0])
	refs := make([]LLVMBasicBlockRef, n)
	for i, h := range handlers {
		refs[i] = LLVMBasicBlockRef{c: h}
	}
	return refs
}

// LLVMGetArgOperand Get a funcletpad argument at the given index.
func LLVMGetArgOperand(funclet LLVMValueRef, i uint32) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetArgOperand(funclet.c, C.unsigned(i))}
}

// LLVMSetArgOperand Set a funcletpad argument at the given index.
func LLVMSetArgOperand(funclet LLVMValueRef, i uint32, value LLVMValueRef) {
	C.LLVMSetArgOperand(funclet.c, C.unsigned(i), value.c)
}

// LLVMGetParentCatchSwitch Get the parent catchswitch instruction of a catchpad instruction.
func LLVMGetParentCatchSwitch(catchPad LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetParentCatchSwitch(catchPad.c)}
}

// LLVMGetPersonalityFn Get the personality function attached to the function.
func LLVMGetPersonalityFn(fn LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetPersonalityFn(fn.c)}
}

// LLVMSetPersonalityFn Set the personality function attached to the function.
func LLVMSetPersonalityFn(fn, persFn LLVMValueRef) {
	C.LLVMSetPersonalityFn(fn.c, persFn.c)
}

// LLVMGetNumSuccessors Obtain the number of successors of a terminator instruction.
func LLVMGetNumSuccessors(term LLVMValueRef) uint32 {
	return uint32(C.LLVMGetNumSuccessors(term.c))
}

// LLVMGetOperand Obtain an operand at a specific index in a LLVM value.
func LLVMGetOperand(val LLVMValueRef, index uint32) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetOperand(val.c, C.unsigned(index))}
}

// LLVMSetOperand Set an operand at a specific index in a LLVM value.
func LLVMSetOperand(user LLVMValueRef, index uint32, val LLVMValueRef) {
	C.LLVMSetOperand(user.c, C.unsigned(index), val.c)
}

// LLVMGetNumArgOperands Obtain the number of arguments to a call instruction.
func LLVMGetNumArgOperands(instr LLVMValueRef) uint32 {
	return uint32(C.LLVMGetNumArgOperands(instr.c))
}

// LLVMGetCalledValue Get the argument to a call instruction.
func LLVMGetCalledValue(instr LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetCalledValue(instr.c)}
}

// LLVMBuildInsertValue Insert a value into an aggregate value's element.
func LLVMBuildInsertValue(builder LLVMBuilderRef, aggVal, eltVal LLVMValueRef, index uint32, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildInsertValue(builder.c, aggVal.c, eltVal.c, C.unsigned(index), name)}
	})
}

// LLVMGetNumOperands Return the number of operands for this value.
func LLVMGetNumOperands(val LLVMValueRef) int {
	return int(C.LLVMGetNumOperands(val.c))
}

// LLVMGetSuccessor Obtain the specified successor of a terminator instruction.
func LLVMGetSuccessor(term LLVMValueRef, i uint32) LLVMBasicBlockRef {
	return LLVMBasicBlockRef{c: C.LLVMGetSuccessor(term.c, C.unsigned(i))}
}

// LLVMGetSwitchDefaultDest Obtain the default destination basic block of a switch instruction.
func LLVMGetSwitchDefaultDest(switchInstr LLVMValueRef) LLVMBasicBlockRef {
	return LLVMBasicBlockRef{c: C.LLVMGetSwitchDefaultDest(switchInstr.c)}
}

// LLVMGetSwitchCaseValue Obtain the case value for a successor of a switch instruction.
// i corresponds to the successor index; the first successor is the default destination, so i must be greater than zero.
func LLVMGetSwitchCaseValue(switchInstr LLVMValueRef, i uint32) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetSwitchCaseValue(switchInstr.c, C.unsigned(i))}
}

func LLVMBuildUnreachable(builder LLVMBuilderRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildUnreachable(builder.c)}
}

func LLVMBuildAdd(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildAdd(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildNSWAdd(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildNSWAdd(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildNUWAdd(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildNUWAdd(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildFAdd(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFAdd(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildSub(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildSub(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildNSWSub(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildNSWSub(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildNUWSub(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildNUWSub(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildFSub(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFSub(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildMul(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildMul(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildNSWMul(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildNSWMul(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildNUWMul(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildNUWMul(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildFMul(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFMul(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildUDiv(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildUDiv(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildExactUDiv(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildExactUDiv(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildSDiv(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildSDiv(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildExactSDiv(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildExactSDiv(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildFDiv(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFDiv(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildURem(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildURem(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildSRem(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildSRem(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildFRem(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFRem(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildShl(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildShl(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildLShr(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildLShr(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildAShr(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildAShr(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildAnd(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildAnd(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildOr(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildOr(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildXor(builder LLVMBuilderRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildXor(builder.c, lHS.c, rHS.c, name)}
	})
}

func LLVMBuildBinOp(builder LLVMBuilderRef, op LLVMOpcode, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildBinOp(builder.c, C.LLVMOpcode(op), lHS.c, rHS.c, name)}
	})
}

func LLVMBuildNeg(builder LLVMBuilderRef, v LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildNeg(builder.c, v.c, name)}
	})
}

func LLVMBuildNSWNeg(builder LLVMBuilderRef, v LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildNSWNeg(builder.c, v.c, name)}
	})
}

func LLVMBuildFNeg(builder LLVMBuilderRef, v LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFNeg(builder.c, v.c, name)}
	})
}

func LLVMBuildNot(builder LLVMBuilderRef, v LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildNot(builder.c, v.c, name)}
	})
}

func LLVMBuildMalloc(builder LLVMBuilderRef, ty LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildMalloc(builder.c, ty.c, name)}
	})
}

func LLVMBuildArrayMalloc(builder LLVMBuilderRef, ty LLVMTypeRef, val LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildArrayMalloc(builder.c, ty.c, val.c, name)}
	})
}

// LLVMBuildMemSet Creates and inserts a memset to the specified pointer and the specified value.
func LLVMBuildMemSet(b LLVMBuilderRef, ptr, val, len LLVMValueRef, align uint32) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildMemSet(b.c, ptr.c, val.c, len.c, C.unsigned(align))}
}

// LLVMBuildMemCpy Creates and inserts a memcpy between the specified pointers.
func LLVMBuildMemCpy(b LLVMBuilderRef, dst LLVMValueRef, dstAlign uint32, src LLVMValueRef, srcAlign uint32, size LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildMemCpy(b.c, dst.c, C.unsigned(dstAlign), src.c, C.unsigned(srcAlign), size.c)}
}

// LLVMBuildMemMove Creates and inserts a memmove between the specified pointers.
func LLVMBuildMemMove(b LLVMBuilderRef, dst LLVMValueRef, dstAlign uint32, src LLVMValueRef, srcAlign uint32, size LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildMemMove(b.c, dst.c, C.unsigned(dstAlign), src.c, C.unsigned(srcAlign), size.c)}
}

func LLVMBuildAlloca(builder LLVMBuilderRef, ty LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildAlloca(builder.c, ty.c, name)}
	})
}

func LLVMBuildArrayAlloca(builder LLVMBuilderRef, ty LLVMTypeRef, val LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildArrayAlloca(builder.c, ty.c, val.c, name)}
	})
}

func LLVMBuildFree(builder LLVMBuilderRef, pointerVal LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildFree(builder.c, pointerVal.c)}
}

func LLVMBuildLoad(builder LLVMBuilderRef, ty LLVMTypeRef, pointerVal LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildLoad2(builder.c, ty.c, pointerVal.c, name)}
	})
}

func LLVMBuildStore(builder LLVMBuilderRef, val, ptr LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildStore(builder.c, val.c, ptr.c)}
}

func LLVMBuildGEP(builder LLVMBuilderRef, ty LLVMTypeRef, pointer LLVMValueRef, indices []LLVMValueRef, name string) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](indices)
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildGEP2(builder.c, ty.c, pointer.c, ptr, length, name)}
	})
}

func LLVMBuildInBoundsGEP(builder LLVMBuilderRef, ty LLVMTypeRef, pointer LLVMValueRef, indices []LLVMValueRef, name string) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](indices)
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildInBoundsGEP2(builder.c, ty.c, pointer.c, ptr, length, name)}
	})
}

func LLVMBuildStructGEP(builder LLVMBuilderRef, ty LLVMTypeRef, pointer LLVMValueRef, idx uint, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildStructGEP2(builder.c, ty.c, pointer.c, C.unsigned(idx), name)}
	})
}

func LLVMBuildGlobalString(builder LLVMBuilderRef, str string, name string) LLVMValueRef {
	return string2CString(str, func(str *C.char) LLVMValueRef {
		return string2CString(name, func(name *C.char) LLVMValueRef {
			return LLVMValueRef{c: C.LLVMBuildGlobalString(builder.c, str, name)}
		})
	})
}

func LLVMBuildGlobalStringPtr(builder LLVMBuilderRef, str string, name string) LLVMValueRef {
	return string2CString(str, func(str *C.char) LLVMValueRef {
		return string2CString(name, func(name *C.char) LLVMValueRef {
			return LLVMValueRef{c: C.LLVMBuildGlobalStringPtr(builder.c, str, name)}
		})
	})
}

func LLVMBuildTrunc(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildTrunc(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildZExt(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildZExt(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildSExt(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildSExt(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildFPToUI(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFPToUI(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildFPToSI(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFPToSI(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildUIToFP(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildUIToFP(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildSIToFP(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildSIToFP(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildFPTrunc(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFPTrunc(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildFPExt(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFPExt(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildPtrToInt(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildPtrToInt(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildIntToPtr(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildIntToPtr(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildBitCast(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildBitCast(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildAddrSpaceCast(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildAddrSpaceCast(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildZExtOrBitCast(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildZExtOrBitCast(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildSExtOrBitCast(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildSExtOrBitCast(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildTruncOrBitCast(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildTruncOrBitCast(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildCast(builder LLVMBuilderRef, op LLVMOpcode, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildCast(builder.c, C.LLVMOpcode(op), val.c, destTy.c, name)}
	})
}

func LLVMBuildPointerCast(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildPointerCast(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildIntCast(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, isSigned bool, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildIntCast2(builder.c, val.c, destTy.c, bool2LLVMBool(isSigned), name)}
	})
}

func LLVMGetCastOpcode(src LLVMValueRef, srcIsSigned bool, destTy LLVMTypeRef, destIsSigned bool) LLVMOpcode {
	return LLVMOpcode(C.LLVMGetCastOpcode(src.c, bool2LLVMBool(srcIsSigned), destTy.c, bool2LLVMBool(destIsSigned)))
}

func LLVMBuildFPCast(builder LLVMBuilderRef, val LLVMValueRef, destTy LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFPCast(builder.c, val.c, destTy.c, name)}
	})
}

func LLVMBuildICmp(builder LLVMBuilderRef, op LLVMIntPredicate, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildICmp(builder.c, C.LLVMIntPredicate(op), lHS.c, rHS.c, name)}
	})
}

func LLVMBuildFCmp(builder LLVMBuilderRef, op LLVMRealPredicate, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFCmp(builder.c, C.LLVMRealPredicate(op), lHS.c, rHS.c, name)}
	})
}

// LLVMGetVolatile Get the volatile flag of the memory access instruction.
func LLVMGetVolatile(inst LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMGetVolatile(inst.c))
}

// LLVMSetVolatile Set the volatile flag of the memory access instruction.
func LLVMSetVolatile(inst LLVMValueRef, isVolatile bool) {
	C.LLVMSetVolatile(inst.c, bool2LLVMBool(isVolatile))
}

// LLVMGetWeak Get the weak flag of the cmpxchg instruction.
func LLVMGetWeak(cmpXchgInst LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMGetWeak(cmpXchgInst.c))
}

// LLVMSetWeak Set the weak flag of the cmpxchg instruction.
func LLVMSetWeak(cmpXchgInst LLVMValueRef, isWeak bool) {
	C.LLVMSetWeak(cmpXchgInst.c, bool2LLVMBool(isWeak))
}

// LLVMGetOrdering Get the ordering of the memory access instruction.
func LLVMGetOrdering(inst LLVMValueRef) LLVMAtomicOrdering {
	return LLVMAtomicOrdering(C.LLVMGetOrdering(inst.c))
}

// LLVMSetOrdering Set the ordering of the memory access instruction.
func LLVMSetOrdering(inst LLVMValueRef, ordering LLVMAtomicOrdering) {
	C.LLVMSetOrdering(inst.c, C.LLVMAtomicOrdering(ordering))
}

// LLVMGetAtomicRMWBinOp Get the operation of the atomicrmw instruction.
func LLVMGetAtomicRMWBinOp(inst LLVMValueRef) LLVMAtomicRMWBinOp {
	return LLVMAtomicRMWBinOp(C.LLVMGetAtomicRMWBinOp(inst.c))
}

// LLVMSetAtomicRMWBinOp Set the operation of the atomicrmw instruction.
func LLVMSetAtomicRMWBinOp(inst LLVMValueRef, binOp LLVMAtomicRMWBinOp) {
	C.LLVMSetAtomicRMWBinOp(inst.c, C.LLVMAtomicRMWBinOp(binOp))
}

// LLVMGetCmpXchgSuccessOrdering Get the success ordering of the cmpxchg instruction.
func LLVMGetCmpXchgSuccessOrdering(cmpXchgInst LLVMValueRef) LLVMAtomicOrdering {
	return LLVMAtomicOrdering(C.LLVMGetCmpXchgSuccessOrdering(cmpXchgInst.c))
}

// LLVMSetCmpXchgSuccessOrdering Set the success ordering of the cmpxchg instruction.
func LLVMSetCmpXchgSuccessOrdering(cmpXchgInst LLVMValueRef, ordering LLVMAtomicOrdering) {
	C.LLVMSetCmpXchgSuccessOrdering(cmpXchgInst.c, C.LLVMAtomicOrdering(ordering))
}

// LLVMGetCmpXchgFailureOrdering Get the failure ordering of the cmpxchg instruction.
func LLVMGetCmpXchgFailureOrdering(cmpXchgInst LLVMValueRef) LLVMAtomicOrdering {
	return LLVMAtomicOrdering(C.LLVMGetCmpXchgFailureOrdering(cmpXchgInst.c))
}

// LLVMSetCmpXchgFailureOrdering Set the failure ordering of the cmpxchg instruction.
func LLVMSetCmpXchgFailureOrdering(cmpXchgInst LLVMValueRef, ordering LLVMAtomicOrdering) {
	C.LLVMSetCmpXchgFailureOrdering(cmpXchgInst.c, C.LLVMAtomicOrdering(ordering))
}

// LLVMIsAtomicSingleThread Get the singlethread flag of the atomic instruction.
func LLVMIsAtomicSingleThread(inst LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsAtomicSingleThread(inst.c))
}

// LLVMSetAtomicSingleThread Set the singlethread flag of the atomic instruction.
func LLVMSetAtomicSingleThread(inst LLVMValueRef, singleThread bool) {
	C.LLVMSetAtomicSingleThread(inst.c, bool2LLVMBool(singleThread))
}

// LLVMBuildFence Create a fence instruction. Note: fence is void-valued and must not be named.
func LLVMBuildFence(builder LLVMBuilderRef, ordering LLVMAtomicOrdering, singleThread bool, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFence(builder.c, C.LLVMAtomicOrdering(ordering), bool2LLVMBool(singleThread), name)}
	})
}

// LLVMBuildAtomicRMW Create an atomicrmw instruction. The C API takes no name;
// apply LLVMSetValueName afterwards if needed.
func LLVMBuildAtomicRMW(builder LLVMBuilderRef, op LLVMAtomicRMWBinOp, ptr, val LLVMValueRef, ordering LLVMAtomicOrdering, singleThread bool) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildAtomicRMW(builder.c, C.LLVMAtomicRMWBinOp(op), ptr.c, val.c, C.LLVMAtomicOrdering(ordering), bool2LLVMBool(singleThread))}
}

// LLVMBuildAtomicCmpXchg Create an atomic cmpxchg instruction. The C API takes no name;
// apply LLVMSetValueName afterwards if needed.
func LLVMBuildAtomicCmpXchg(builder LLVMBuilderRef, ptr, cmp, new LLVMValueRef, successOrdering, failureOrdering LLVMAtomicOrdering, singleThread bool) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildAtomicCmpXchg(builder.c, ptr.c, cmp.c, new.c, C.LLVMAtomicOrdering(successOrdering), C.LLVMAtomicOrdering(failureOrdering), bool2LLVMBool(singleThread))}
}

func LLVMBuildPhi(builder LLVMBuilderRef, ty LLVMTypeRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildPhi(builder.c, ty.c, name)}
	})
}

func LLVMBuildCall(builder LLVMBuilderRef, ty LLVMTypeRef, fn LLVMValueRef, args []LLVMValueRef, name string) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](args)
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildCall2(builder.c, ty.c, fn.c, ptr, length, name)}
	})
}

func LLVMBuildSelect(builder LLVMBuilderRef, ifv LLVMValueRef, thenv, elsev LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildSelect(builder.c, ifv.c, thenv.c, elsev.c, name)}
	})
}

func LLVMBuildExtractElement(builder LLVMBuilderRef, vecVal, index LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildExtractElement(builder.c, vecVal.c, index.c, name)}
	})
}

func LLVMBuildExtractValue(builder LLVMBuilderRef, aggVal LLVMValueRef, index uint32, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildExtractValue(builder.c, aggVal.c, C.unsigned(index), name)}
	})
}

// LLVMBuildInsertElement Insert a value into a vector element.
func LLVMBuildInsertElement(builder LLVMBuilderRef, vecVal, eltVal, index LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildInsertElement(builder.c, vecVal.c, eltVal.c, index.c, name)}
	})
}

// LLVMBuildShuffleVector Create a shufflevector instruction.
func LLVMBuildShuffleVector(builder LLVMBuilderRef, v1, v2, mask LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildShuffleVector(builder.c, v1.c, v2.c, mask.c, name)}
	})
}

func LLVMBuildIsNull(builder LLVMBuilderRef, val LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildIsNull(builder.c, val.c, name)}
	})
}

func LLVMBuildIsNotNull(builder LLVMBuilderRef, val LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildIsNotNull(builder.c, val.c, name)}
	})
}

func LLVMBuildPtrDiff(builder LLVMBuilderRef, elemTy LLVMTypeRef, lHS, rHS LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildPtrDiff2(builder.c, elemTy.c, lHS.c, rHS.c, name)}
	})
}

// LLVMCreatePassManager Constructs a new whole-module pass pipeline.
// This type of pipeline is suitable for link-time optimization and whole-module transformations.
func LLVMCreatePassManager() LLVMPassManagerRef {
	return LLVMPassManagerRef{c: C.LLVMCreatePassManager()}
}

// LLVMCreateFunctionPassManagerForModule Constructs a new function-by-function pass pipeline over the module provider.
// It does not take ownership of the module provider.
// This type of pipeline is suitable for code generation and JIT compilation tasks.
func LLVMCreateFunctionPassManagerForModule(m LLVMModuleRef) LLVMPassManagerRef {
	return LLVMPassManagerRef{c: C.LLVMCreateFunctionPassManagerForModule(m.c)}
}

// LLVMRunPassManager Initializes, executes on the provided module, and finalizes all of the passes scheduled in the pass manager.
// Returns true if any of the passes modified the module, false otherwise.
func LLVMRunPassManager(pm LLVMPassManagerRef, m LLVMModuleRef) bool {
	return llvmBool2bool(C.LLVMRunPassManager(pm.c, m.c))
}

// LLVMInitializeFunctionPassManager Initializes all of the function passes scheduled in the function pass manager.
// Returns true if any of the passes modified the module, false otherwise.
func LLVMInitializeFunctionPassManager(fpm LLVMPassManagerRef) bool {
	return llvmBool2bool(C.LLVMInitializeFunctionPassManager(fpm.c))
}

// LLVMRunFunctionPassManager Executes all of the function passes scheduled in the function pass manager on the provided function.
// Returns true if any of the passes modified the function, false otherwise.
func LLVMRunFunctionPassManager(fpm LLVMPassManagerRef, f LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMRunFunctionPassManager(fpm.c, f.c))
}

// LLVMFinalizeFunctionPassManager Finalizes all of the function passes scheduled in the function pass manager.
// Returns true if any of the passes modified the module, false otherwise.
func LLVMFinalizeFunctionPassManager(fpm LLVMPassManagerRef) bool {
	return llvmBool2bool(C.LLVMFinalizeFunctionPassManager(fpm.c))
}

// LLVMDisposePassManager Frees the memory of a pass pipeline.
// For function pipelines, does not free the module provider.
func LLVMDisposePassManager(pm LLVMPassManagerRef) {
	C.LLVMDisposePassManager(pm.c)
}

func LLVMSetNUW(arithInst LLVMValueRef, hasNUW bool) {
	C.LLVMSetNUW(arithInst.c, bool2LLVMBool(hasNUW))
}

func LLVMBuildNUWNeg(builder LLVMBuilderRef, v LLVMValueRef, name string) LLVMValueRef {
	ref := LLVMBuildNeg(builder, v, name)
	LLVMSetNUW(ref, true)
	return ref
}

// LLVMConstNUWNeg Obtain the negation of a constant.
// LLVM 22 将 C API LLVMConstNUWNeg 标记弃用（弃用提示词为 "Use LLVMConstNull instead."，
// 语义明显有误）；折叠后的整数常量无法携带 nuw 标志，ConstNeg 等价且无弃用警告。
func LLVMConstNUWNeg(constantVal LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstNeg(constantVal.c)}
}

// LLVMCreateMemoryBufferWithContentsOfFile Read a file into a memory buffer.
func LLVMCreateMemoryBufferWithContentsOfFile(path string) (LLVMMemoryBufferRef, error) {
	var buf LLVMMemoryBufferRef
	err := string2CString(path, func(cpath *C.char) error {
		return llvmError2Error(func(errstr **C.char) C.LLVMBool {
			return C.LLVMCreateMemoryBufferWithContentsOfFile(cpath, &buf.c, errstr)
		})
	})
	if err != nil {
		return LLVMMemoryBufferRef{}, err
	}
	return buf, nil
}

// LLVMCreateMemoryBufferWithMemoryRangeCopy Create a memory buffer from a memory range, copying the data.
func LLVMCreateMemoryBufferWithMemoryRangeCopy(data []byte, name string) LLVMMemoryBufferRef {
	return string2CString(name, func(name *C.char) LLVMMemoryBufferRef {
		var ptr *C.char
		if len(data) > 0 {
			ptr = (*C.char)(unsafe.Pointer(&data[0]))
		}
		return LLVMMemoryBufferRef{c: C.LLVMCreateMemoryBufferWithMemoryRangeCopy(ptr, C.size_t(len(data)), name)}
	})
}

// LLVMGetBufferStart Get the start of the buffer.
func LLVMGetBufferStart(buf LLVMMemoryBufferRef) []byte {
	ptr := C.LLVMGetBufferStart(buf.c)
	if ptr == nil {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(ptr)), int(C.LLVMGetBufferSize(buf.c)))
}

// LLVMGetBufferSize Get the size of the buffer.
func LLVMGetBufferSize(buf LLVMMemoryBufferRef) uint64 {
	return uint64(C.LLVMGetBufferSize(buf.c))
}

// LLVMDisposeMemoryBuffer Dispose of a memory buffer.
func LLVMDisposeMemoryBuffer(buf LLVMMemoryBufferRef) {
	C.LLVMDisposeMemoryBuffer(buf.c)
}

// LLVMGoEmitError 经 LLVMContext::emitError 发出诊断（LLVM-C 无对应 API，供诊断回调测试/调试用）
func LLVMGoEmitError(c LLVMContextRef, msg string) {
	string2CString(msg, func(cmsg *C.char) bool {
		C.LLVMGoEmitError(c.c, cmsg)
		return false
	})
}

// LLVMGetFirstUse Obtain the first use of a value.
func LLVMGetFirstUse(v LLVMValueRef) LLVMUseRef {
	return LLVMUseRef{c: C.LLVMGetFirstUse(v.c)}
}

// LLVMGetNextUse Obtain the next use following a use.
func LLVMGetNextUse(u LLVMUseRef) LLVMUseRef {
	return LLVMUseRef{c: C.LLVMGetNextUse(u.c)}
}

// LLVMGetUser Obtain the user value for a use.
func LLVMGetUser(u LLVMUseRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetUser(u.c)}
}

// LLVMGetUsedValue Obtain the value this use corresponds to.
func LLVMGetUsedValue(u LLVMUseRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetUsedValue(u.c)}
}

// LLVMReplaceAllUsesWith Replace all uses of a value with another one.
func LLVMReplaceAllUsesWith(oldVal, newVal LLVMValueRef) {
	C.LLVMReplaceAllUsesWith(oldVal.c, newVal.c)
}

// LLVMSetSuccessor Set the specified successor of the terminator instruction.
func LLVMSetSuccessor(term LLVMValueRef, i uint32, blk LLVMBasicBlockRef) {
	C.LLVMSetSuccessor(term.c, C.unsigned(i), blk.c)
}

// LLVMIsConditional Determine whether a terminator instruction is conditional.
func LLVMIsConditional(branch LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsConditional(branch.c))
}

// LLVMGetCondition Obtain the condition of the terminator instruction.
func LLVMGetCondition(branch LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetCondition(branch.c)}
}

// LLVMSetCondition Set the condition of the terminator instruction.
func LLVMSetCondition(branch, cond LLVMValueRef) {
	C.LLVMSetCondition(branch.c, cond.c)
}

// LLVMGetUndef Create an undef value of the given type.
func LLVMGetUndef(ty LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetUndef(ty.c)}
}

// LLVMGetPoison Create a poison value of the given type.
func LLVMGetPoison(ty LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetPoison(ty.c)}
}

// LLVMBuildFreeze Freeze the given value.
func LLVMBuildFreeze(b LLVMBuilderRef, v LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(cs *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFreeze(b.c, v.c, cs)}
	})
}

// LLVMGetFirstGlobal Obtain the first global variable in a module.
func LLVMGetFirstGlobal(m LLVMModuleRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetFirstGlobal(m.c)}
}

// LLVMGetNextGlobal Obtain the next global variable in a module.
func LLVMGetNextGlobal(g LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetNextGlobal(g.c)}
}

// LLVMLookupIntrinsicID Look up the ID for the specified intrinsic name. Returns 0 when not found.
func LLVMLookupIntrinsicID(name string) uint32 {
	return string2CString(name, func(cs *C.char) uint32 {
		return uint32(C.LLVMLookupIntrinsicID(cs, C.size_t(len(name))))
	})
}

// LLVMIntrinsicGetName Get the name of the specified intrinsic.
func LLVMIntrinsicGetName(id uint32) string {
	var length C.size_t
	s := C.LLVMIntrinsicGetName(C.unsigned(id), &length)
	return C.GoStringN(s, C.int(length))
}

// LLVMIntrinsicIsOverloaded Determine whether the specified intrinsic is overloaded.
func LLVMIntrinsicIsOverloaded(id uint32) bool {
	return llvmBool2bool(C.LLVMIntrinsicIsOverloaded(C.unsigned(id)))
}

// LLVMGetIntrinsicDeclaration Get the declaration of the specified intrinsic in the module.
func LLVMGetIntrinsicDeclaration(m LLVMModuleRef, id uint32, paramTypes []LLVMTypeRef) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMTypeRef, C.LLVMTypeRef](paramTypes)
	return LLVMValueRef{c: C.LLVMGetIntrinsicDeclaration(m.c, C.unsigned(id), ptr, C.size_t(length))}
}

// LLVMFastMathFlags Fast-math flags bitmask.
type LLVMFastMathFlags uint32

const (
	LLVMFastMathAllowReassoc    LLVMFastMathFlags = C.LLVMFastMathAllowReassoc
	LLVMFastMathNoNaNs          LLVMFastMathFlags = C.LLVMFastMathNoNaNs
	LLVMFastMathNoInfs          LLVMFastMathFlags = C.LLVMFastMathNoInfs
	LLVMFastMathNoSignedZeros   LLVMFastMathFlags = C.LLVMFastMathNoSignedZeros
	LLVMFastMathAllowReciprocal LLVMFastMathFlags = C.LLVMFastMathAllowReciprocal
	LLVMFastMathAllowContract   LLVMFastMathFlags = C.LLVMFastMathAllowContract
	LLVMFastMathApproxFunc      LLVMFastMathFlags = C.LLVMFastMathApproxFunc
	LLVMFastMathNone            LLVMFastMathFlags = 0
	LLVMFastMathAll             LLVMFastMathFlags = LLVMFastMathAllowReassoc | LLVMFastMathNoNaNs |
		LLVMFastMathNoInfs | LLVMFastMathNoSignedZeros | LLVMFastMathAllowReciprocal |
		LLVMFastMathAllowContract | LLVMFastMathApproxFunc
)

// LLVMGetFastMathFlags Get the fast-math flags of an FP instruction.
func LLVMGetFastMathFlags(v LLVMValueRef) LLVMFastMathFlags {
	return LLVMFastMathFlags(C.LLVMGetFastMathFlags(v.c))
}

// LLVMSetFastMathFlags Set the fast-math flags of an FP instruction.
func LLVMSetFastMathFlags(v LLVMValueRef, f LLVMFastMathFlags) {
	C.LLVMSetFastMathFlags(v.c, C.LLVMFastMathFlags(f))
}

// LLVMCanValueUseFastMathFlags Whether the value can carry fast-math flags.
func LLVMCanValueUseFastMathFlags(v LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMCanValueUseFastMathFlags(v.c))
}

// LLVMGEPNoWrapFlags GEP no-wrap flags bitmask.
type LLVMGEPNoWrapFlags uint32

const (
	LLVMGEPFlagInBounds LLVMGEPNoWrapFlags = C.LLVMGEPFlagInBounds
	LLVMGEPFlagNUSW     LLVMGEPNoWrapFlags = C.LLVMGEPFlagNUSW
	LLVMGEPFlagNUW      LLVMGEPNoWrapFlags = C.LLVMGEPFlagNUW
)

// LLVMGEPGetNoWrapFlags Get the no-wrap flags of a GEP instruction.
func LLVMGEPGetNoWrapFlags(gep LLVMValueRef) LLVMGEPNoWrapFlags {
	return LLVMGEPNoWrapFlags(C.LLVMGEPGetNoWrapFlags(gep.c))
}

// LLVMGEPSetNoWrapFlags Set the no-wrap flags of a GEP instruction.
func LLVMGEPSetNoWrapFlags(gep LLVMValueRef, f LLVMGEPNoWrapFlags) {
	C.LLVMGEPSetNoWrapFlags(gep.c, C.LLVMGEPNoWrapFlags(f))
}

// LLVMBuildGEPWithNoWrapFlags Create a getelementptr instruction with no-wrap flags.
func LLVMBuildGEPWithNoWrapFlags(b LLVMBuilderRef, ty LLVMTypeRef, p LLVMValueRef, indices []LLVMValueRef, name string, f LLVMGEPNoWrapFlags) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](indices)
	return string2CString(name, func(cs *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildGEPWithNoWrapFlags(b.c, ty.c, p.c, ptr, C.unsigned(length), cs, C.LLVMGEPNoWrapFlags(f))}
	})
}

// LLVMSetNSW Set the no-signed-wrap flag of an arithmetic instruction.
func LLVMSetNSW(v LLVMValueRef, hasNSW bool) {
	C.LLVMSetNSW(v.c, bool2LLVMBool(hasNSW))
}

// LLVMSetExact Set the exact flag of a div/rem/shift instruction.
func LLVMSetExact(v LLVMValueRef, isExact bool) {
	C.LLVMSetExact(v.c, bool2LLVMBool(isExact))
}

// LLVMSetNNeg Set the non-negative flag of a sign-extending instruction.
func LLVMSetNNeg(v LLVMValueRef, isNonNeg bool) {
	C.LLVMSetNNeg(v.c, bool2LLVMBool(isNonNeg))
}

// LLVMSetIsInBounds Set the inbounds flag of a GEP instruction.
func LLVMSetIsInBounds(v LLVMValueRef, inBounds bool) {
	C.LLVMSetIsInBounds(v.c, bool2LLVMBool(inBounds))
}

// LLVMTailCallKind Tail call kind.
type LLVMTailCallKind int32

const (
	LLVMTailCallKindNone     LLVMTailCallKind = C.LLVMTailCallKindNone
	LLVMTailCallKindTail     LLVMTailCallKind = C.LLVMTailCallKindTail
	LLVMTailCallKindMustTail LLVMTailCallKind = C.LLVMTailCallKindMustTail
	LLVMTailCallKindNoTail   LLVMTailCallKind = C.LLVMTailCallKindNoTail
)

// LLVMSetTailCall Set whether the call is a tail call.
func LLVMSetTailCall(v LLVMValueRef, isTailCall bool) {
	C.LLVMSetTailCall(v.c, bool2LLVMBool(isTailCall))
}

// LLVMIsTailCall Whether the call is a tail call.
func LLVMIsTailCall(v LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsTailCall(v.c))
}

// LLVMSetTailCallKind Set the tail call kind of the call.
func LLVMSetTailCallKind(v LLVMValueRef, k LLVMTailCallKind) {
	C.LLVMSetTailCallKind(v.c, C.LLVMTailCallKind(k))
}

// LLVMSetInstrParamAlignment Set the alignment attribute of the call's parameter (1-based, 0 = return).
func LLVMSetInstrParamAlignment(v LLVMValueRef, index uint32, align uint32) {
	C.LLVMSetInstrParamAlignment(v.c, C.LLVMAttributeIndex(index), C.unsigned(align))
}

// LLVMGetSyncScopeID Get the sync scope ID for the given name. LLVM recognizes
// common names (e.g. "system", "singlethread") but assigns the IDs itself; do not
// hardcode specific values.
func LLVMGetSyncScopeID(c LLVMContextRef, name string) uint32 {
	return string2CString(name, func(cs *C.char) uint32 {
		return uint32(C.LLVMGetSyncScopeID(c.c, cs, C.size_t(len(name))))
	})
}

// LLVMGetAtomicSyncScopeID Get the sync scope ID of an atomic instruction.
func LLVMGetAtomicSyncScopeID(v LLVMValueRef) uint32 {
	return uint32(C.LLVMGetAtomicSyncScopeID(v.c))
}

// LLVMSetAtomicSyncScopeID Set the sync scope ID of an atomic instruction.
func LLVMSetAtomicSyncScopeID(v LLVMValueRef, ssid uint32) {
	C.LLVMSetAtomicSyncScopeID(v.c, C.unsigned(ssid))
}

// LLVMBuildFenceSyncScope Create a fence instruction with a sync scope.
func LLVMBuildFenceSyncScope(b LLVMBuilderRef, order LLVMAtomicOrdering, ssid uint32, name string) LLVMValueRef {
	return string2CString(name, func(cs *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFenceSyncScope(b.c, C.LLVMAtomicOrdering(order), C.unsigned(ssid), cs)}
	})
}

// LLVMBuildAtomicRMWSyncScope Create an atomicrmw instruction with a sync scope.
func LLVMBuildAtomicRMWSyncScope(b LLVMBuilderRef, op LLVMAtomicRMWBinOp, ptr, val LLVMValueRef, order LLVMAtomicOrdering, ssid uint32) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildAtomicRMWSyncScope(b.c, C.LLVMAtomicRMWBinOp(op), ptr.c, val.c, C.LLVMAtomicOrdering(order), C.unsigned(ssid))}
}

// LLVMBuildAtomicCmpXchgSyncScope Create a cmpxchg instruction with a sync scope.
func LLVMBuildAtomicCmpXchgSyncScope(b LLVMBuilderRef, ptr, cmp, newVal LLVMValueRef, successOrder, failureOrder LLVMAtomicOrdering, ssid uint32) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildAtomicCmpXchgSyncScope(b.c, ptr.c, cmp.c, newVal.c, C.LLVMAtomicOrdering(successOrder), C.LLVMAtomicOrdering(failureOrder), C.unsigned(ssid))}
}

// LLVMCreateOperandBundle Create an operand bundle with the given tag and inputs.
func LLVMCreateOperandBundle(tag string, inputs []LLVMValueRef) LLVMOperandBundleRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](inputs)
	return string2CString(tag, func(cs *C.char) LLVMOperandBundleRef {
		return LLVMOperandBundleRef{c: C.LLVMCreateOperandBundle(cs, C.size_t(len(tag)), ptr, C.unsigned(length))}
	})
}

// LLVMDisposeOperandBundle Dispose an operand bundle.
func LLVMDisposeOperandBundle(b LLVMOperandBundleRef) {
	C.LLVMDisposeOperandBundle(b.c)
}

// LLVMGetOperandBundleTag Get the tag of an operand bundle.
func LLVMGetOperandBundleTag(b LLVMOperandBundleRef) string {
	var length C.size_t
	s := C.LLVMGetOperandBundleTag(b.c, &length)
	return C.GoStringN(s, C.int(length))
}

// LLVMBuildCallWithOperandBundles Create a call instruction with operand bundles.
func LLVMBuildCallWithOperandBundles(b LLVMBuilderRef, ty LLVMTypeRef, fn LLVMValueRef, args []LLVMValueRef, bundles []LLVMOperandBundleRef, name string) LLVMValueRef {
	aptr, alen := slice2Ptr[LLVMValueRef, C.LLVMValueRef](args)
	bptr, blen := slice2Ptr[LLVMOperandBundleRef, C.LLVMOperandBundleRef](bundles)
	return string2CString(name, func(cs *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildCallWithOperandBundles(b.c, ty.c, fn.c, aptr, C.unsigned(alen), bptr, C.unsigned(blen), cs)}
	})
}

// LLVMBuildInvokeWithOperandBundles Create an invoke instruction with operand bundles.
func LLVMBuildInvokeWithOperandBundles(b LLVMBuilderRef, ty LLVMTypeRef, fn LLVMValueRef, args []LLVMValueRef, then, catch LLVMBasicBlockRef, bundles []LLVMOperandBundleRef, name string) LLVMValueRef {
	aptr, alen := slice2Ptr[LLVMValueRef, C.LLVMValueRef](args)
	bptr, blen := slice2Ptr[LLVMOperandBundleRef, C.LLVMOperandBundleRef](bundles)
	return string2CString(name, func(cs *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildInvokeWithOperandBundles(b.c, ty.c, fn.c, aptr, C.unsigned(alen), then.c, catch.c, bptr, C.unsigned(blen), cs)}
	})
}

// LLVMGetSection Get the section of the global value.
func LLVMGetSection(g LLVMValueRef) string {
	s := C.LLVMGetSection(g.c)
	if s == nil {
		return ""
	}
	return C.GoString(s)
}

// LLVMSetSection Set the section of the global value.
func LLVMSetSection(g LLVMValueRef, section string) {
	string2CString(section, func(cs *C.char) struct{} {
		C.LLVMSetSection(g.c, cs)
		return struct{}{}
	})
}

// LLVMDLLStorageClass DLL storage class.
type LLVMDLLStorageClass int32

const (
	LLVMDefaultStorageClass   LLVMDLLStorageClass = C.LLVMDefaultStorageClass
	LLVMDLLImportStorageClass LLVMDLLStorageClass = C.LLVMDLLImportStorageClass
	LLVMDLLExportStorageClass LLVMDLLStorageClass = C.LLVMDLLExportStorageClass
)

// LLVMGetDLLStorageClass Get the DLL storage class of the global value.
func LLVMGetDLLStorageClass(g LLVMValueRef) LLVMDLLStorageClass {
	return LLVMDLLStorageClass(C.LLVMGetDLLStorageClass(g.c))
}

// LLVMSetDLLStorageClass Set the DLL storage class of the global value.
func LLVMSetDLLStorageClass(g LLVMValueRef, class LLVMDLLStorageClass) {
	C.LLVMSetDLLStorageClass(g.c, C.LLVMDLLStorageClass(class))
}

// LLVMGetGC Get the GC name of the function.
func LLVMGetGC(fn LLVMValueRef) string {
	s := C.LLVMGetGC(fn.c)
	if s == nil {
		return ""
	}
	return C.GoString(s)
}

// LLVMSetGC Set the GC name of the function.
func LLVMSetGC(fn LLVMValueRef, name string) {
	string2CString(name, func(cs *C.char) struct{} {
		C.LLVMSetGC(fn.c, cs)
		return struct{}{}
	})
}

// LLVMGetPrefixData Get the prefix data of the function.
func LLVMGetPrefixData(fn LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetPrefixData(fn.c)}
}

// LLVMSetPrefixData Set the prefix data of the function.
func LLVMSetPrefixData(fn, prefixData LLVMValueRef) {
	C.LLVMSetPrefixData(fn.c, prefixData.c)
}

// LLVMGetPrologueData Get the prologue data of the function.
func LLVMGetPrologueData(fn LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetPrologueData(fn.c)}
}

// LLVMSetPrologueData Set the prologue data of the function.
func LLVMSetPrologueData(fn, prologueData LLVMValueRef) {
	C.LLVMSetPrologueData(fn.c, prologueData.c)
}

// LLVMAddAlias2 Create an alias with the given value type, address space and aliasee.
func LLVMAddAlias2(m LLVMModuleRef, valueTy LLVMTypeRef, addrSpace uint32, aliasee LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(cs *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMAddAlias2(m.c, valueTy.c, C.unsigned(addrSpace), aliasee.c, cs)}
	})
}

// LLVMAliasGetAliasee Get the aliasee of the alias.
func LLVMAliasGetAliasee(a LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMAliasGetAliasee(a.c)}
}

// LLVMAliasSetAliasee Set the aliasee of the alias.
func LLVMAliasSetAliasee(a, aliasee LLVMValueRef) {
	C.LLVMAliasSetAliasee(a.c, aliasee.c)
}

// LLVMAddGlobalIFunc Add a global indirect function to the module.
func LLVMAddGlobalIFunc(m LLVMModuleRef, name string, ty LLVMTypeRef, addrSpace uint32, resolver LLVMValueRef) LLVMValueRef {
	return string2CString(name, func(cs *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMAddGlobalIFunc(m.c, cs, C.size_t(len(name)), ty.c, C.unsigned(addrSpace), resolver.c)}
	})
}

// LLVMGetFirstGlobalAlias Obtain the first global alias in a module.
func LLVMGetFirstGlobalAlias(m LLVMModuleRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetFirstGlobalAlias(m.c)}
}

// LLVMGetNextGlobalAlias Obtain the next global alias in a module.
func LLVMGetNextGlobalAlias(a LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetNextGlobalAlias(a.c)}
}

// LLVMGetNamedGlobalAlias Obtain the global alias with the given name.
func LLVMGetNamedGlobalAlias(m LLVMModuleRef, name string) LLVMValueRef {
	return string2CString(name, func(cs *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMGetNamedGlobalAlias(m.c, cs, C.size_t(len(name)))}
	})
}
