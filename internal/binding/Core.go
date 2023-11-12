package binding

/*
#include "llvm-c/Core.h"
#include "Core.h"
*/
import "C"

type LLVMOpcode int32

const (
	LLVMRet LLVMOpcode = 1 + iota
	LLVMBr
	LLVMSwitch
	LLVMIndirectBr
	LLVMInvoke
	LLVMUnreachable LLVMOpcode = 2 + iota
	LLVMAdd
	LLVMFAdd
	LLVMSub
	LLVMFSub
	LLVMMul
	LLVMFMul
	LLVMUDiv
	LLVMSDiv
	LLVMFDiv
	LLVMURem
	LLVMSRem
	LLVMFRem
	LLVMShl
	LLVMLShr
	LLVMAShr
	LLVMAnd
	LLVMOr
	LLVMXor
	LLVMAlloca
	LLVMLoad
	LLVMStore
	LLVMGetElementPtr
	LLVMTrunc
	LLVMZExt
	LLVMSExt
	LLVMFPToUI
	LLVMFPToSI
	LLVMUIToFP
	LLVMSIToFP
	LLVMFPTrunc
	LLVMFPExt
	LLVMPtrToInt
	LLVMIntToPtr
	LLVMBitCast
	LLVMICmp
	LLVMFCmp
	LLVMPHI
	LLVMCall
	LLVMSelect
	LLVMUserOp1
	LLVMUserOp2
	LLVMVAArg
	LLVMExtractElement
	LLVMInsertElement
	LLVMShuffleVector
	LLVMExtractValue
	LLVMInsertValue
	LLVMFence
	LLVMAtomicCmpXchg
	LLVMAtomicRMW
	LLVMResume
	LLVMLandingPad
	LLVMAddrSpaceCast
	LLVMCleanupRet
	LLVMCatchRet
	LLVMCatchPad
	LLVMCleanupPad
	LLVMCatchSwitch
	LLVMFNeg
	LLVMCallBr
	LLVMFreeze
)

type LLVMTypeKind int32

const (
	// LLVMVoidTypeKind type with no size
	LLVMVoidTypeKind LLVMTypeKind = iota
	// LLVMHalfTypeKind 16 bit floating point type
	LLVMHalfTypeKind
	// LLVMFloatTypeKind 32 bit floating point type
	LLVMFloatTypeKind
	// LLVMDoubleTypeKind 64 bit floating point type
	LLVMDoubleTypeKind
	// LLVMX86_FP80TypeKind 80 bit floating point type (X87)
	LLVMX86_FP80TypeKind
	// LLVMFP128TypeKind 128 bit floating point type (112-bit mantissa)
	LLVMFP128TypeKind
	// LLVMPPC_FP128TypeKind 128 bit floating point type (two 64-bits)
	LLVMPPC_FP128TypeKind
	// LLVMLabelTypeKind Labels
	LLVMLabelTypeKind
	// LLVMIntegerTypeKind Arbitrary bit width integers
	LLVMIntegerTypeKind
	// LLVMFunctionTypeKind Functions
	LLVMFunctionTypeKind
	// LLVMStructTypeKind Structures
	LLVMStructTypeKind
	// LLVMArrayTypeKind Arrays
	LLVMArrayTypeKind
	// LLVMPointerTypeKind Pointers
	LLVMPointerTypeKind
	// LLVMVectorTypeKind Fixed width SIMD vector type
	LLVMVectorTypeKind
	// LLVMMetadataTypeKind Metadata
	LLVMMetadataTypeKind
	// LLVMX86_MMXTypeKind X86 MMX
	LLVMX86_MMXTypeKind
	// LLVMTokenTypeKind Tokens
	LLVMTokenTypeKind
	// LLVMScalableVectorTypeKind Scalable SIMD vector type
	LLVMScalableVectorTypeKind
	// LLVMBFloatTypeKind 16 bit brain floating point type
	LLVMBFloatTypeKind
	// LLVMX86_AMXTypeKind X86 AMX
	LLVMX86_AMXTypeKind
)

type LLVMLinkage int32

const (
	// LLVMExternalLinkage Externally visible function
	LLVMExternalLinkage LLVMLinkage = iota
	LLVMAvailableExternallyLinkage
	// LLVMLinkOnceAnyLinkage Keep one copy of function when linking (inline)
	LLVMLinkOnceAnyLinkage
	// LLVMLinkOnceODRLinkage Same, but only replaced by something equivalent.
	LLVMLinkOnceODRLinkage
	// LLVMLinkOnceODRAutoHideLinkage Obsolete
	LLVMLinkOnceODRAutoHideLinkage
	// LLVMWeakAnyLinkage Keep one copy of function when linking (weak)
	LLVMWeakAnyLinkage
	// LLVMWeakODRLinkage Same, but only replaced by something equivalent.
	LLVMWeakODRLinkage
	// LLVMAppendingLinkage Special purpose, only applies to global arrays
	LLVMAppendingLinkage
	// LLVMInternalLinkage Rename collisions when linking (static functions)
	LLVMInternalLinkage
	// LLVMPrivateLinkage Like Internal, but omit from symbol table
	LLVMPrivateLinkage
	// LLVMDLLImportLinkage Obsolete
	LLVMDLLImportLinkage
	// LLVMDLLExportLinkage Obsolete
	LLVMDLLExportLinkage
	// LLVMExternalWeakLinkage ExternalWeak linkage description
	LLVMExternalWeakLinkage
	// LLVMGhostLinkage Obsolete
	LLVMGhostLinkage
	// LLVMCommonLinkage Tentative definitions
	LLVMCommonLinkage
	// LLVMLinkerPrivateLinkage Like Private, but linker removes.
	LLVMLinkerPrivateLinkage
	// LLVMLinkerPrivateWeakLinkage Like LinkerPrivate, but is weak.
	LLVMLinkerPrivateWeakLinkage
)

type LLVMVisibility int32

const (
	// LLVMDefaultVisibility The GV is visible
	LLVMDefaultVisibility LLVMVisibility = iota
	// LLVMHiddenVisibility The GV is hidden
	LLVMHiddenVisibility
	// LLVMProtectedVisibility The GV is protected
	LLVMProtectedVisibility
)

type LLVMUnnamedAddr int32

const (
	// LLVMNoUnnamedAddr Address of the GV is significant.
	LLVMNoUnnamedAddr LLVMUnnamedAddr = iota
	// LLVMLocalUnnamedAddr Address of the GV is locally insignificant.
	LLVMLocalUnnamedAddr
	// LLVMGlobalUnnamedAddr Address of the GV is globally insignificant.
	LLVMGlobalUnnamedAddr
)

type LLVMValueKind int32

const (
	LLVMArgumentValueKind LLVMValueKind = iota
	LLVMBasicBlockValueKind
	LLVMMemoryUseValueKind
	LLVMMemoryDefValueKind
	LLVMMemoryPhiValueKind
	LLVMFunctionValueKind
	LLVMGlobalAliasValueKind
	LLVMGlobalIFuncValueKind
	LLVMGlobalVariableValueKind
	LLVMBlockAddressValueKind
	LLVMConstantExprValueKind
	LLVMConstantArrayValueKind
	LLVMConstantStructValueKind
	LLVMConstantVectorValueKind
	LLVMUndefValueValueKind
	LLVMConstantAggregateZeroValueKind
	LLVMConstantDataArrayValueKind
	LLVMConstantDataVectorValueKind
	LLVMConstantIntValueKind
	LLVMConstantFPValueKind
	LLVMConstantPointerNullValueKind
	LLVMConstantTokenNoneValueKind
	LLVMMetadataAsValueValueKind
	LLVMInlineAsmValueKind
	LLVMInstructionValueKind
	LLVMPoisonValueValueKind
)

type LLVMIntPredicate int32

const (
	// LLVMIntEQ equal
	LLVMIntEQ LLVMTypeKind = 32 + iota
	// LLVMIntNE not equal
	LLVMIntNE
	// LLVMIntUGT unsigned greater than
	LLVMIntUGT
	// LLVMIntUGE unsigned greater or equal
	LLVMIntUGE
	// LLVMIntULT unsigned less than
	LLVMIntULT
	// LLVMIntULE unsigned less or equal
	LLVMIntULE
	// LLVMIntSGT signed greater than
	LLVMIntSGT
	// LLVMIntSGE signed greater or equal
	LLVMIntSGE
	// LLVMIntSLT signed less than
	LLVMIntSLT
	// LLVMIntSLE signed less or equal
	LLVMIntSLE
)

type LLVMRealPredicate int32

const (
	// LLVMRealPredicateFalse Always false (always folded)
	LLVMRealPredicateFalse LLVMTypeKind = iota
	// LLVMRealOEQ True if ordered and equal
	LLVMRealOEQ
	// LLVMRealOGT True if ordered and greater than
	LLVMRealOGT
	// LLVMRealOGE True if ordered and greater than or equal
	LLVMRealOGE
	// LLVMRealOLT True if ordered and less than
	LLVMRealOLT
	// LLVMRealOLE True if ordered and less than or equal
	LLVMRealOLE
	// LLVMRealONE True if ordered and operands are unequal
	LLVMRealONE
	// LLVMRealORD True if ordered (no nans)
	LLVMRealORD
	// LLVMRealUNO True if unordered: isnan(X) | isnan(Y)
	LLVMRealUNO
	// LLVMRealUEQ True if unordered or equal
	LLVMRealUEQ
	// LLVMRealUGT True if unordered or greater than
	LLVMRealUGT
	// LLVMRealUGE True if unordered, greater than, or equal
	LLVMRealUGE
	// LLVMRealULT True if unordered or less than
	LLVMRealULT
	// LLVMRealULE True if unordered, less than, or equal
	LLVMRealULE
	// LLVMRealUNE True if unordered or not equal
	LLVMRealUNE
	// LLVMRealPredicateTrue Always true (always folded)
	LLVMRealPredicateTrue
)

type LLVMThreadLocalMode int32

const (
	LLVMNotThreadLocal LLVMThreadLocalMode = iota
	LLVMGeneralDynamicTLSModel
	LLVMLocalDynamicTLSModel
	LLVMInitialExecTLSModel
	LLVMLocalExecTLSModel
)

func LLVMDisposeMessage(message *C.char) {
	C.LLVMDisposeMessage(message)
}

// LLVMContextCreate Create a new context.
// Every call to this function should be paired with a call to LLVMContextDispose() or the context will leak memory.
func LLVMContextCreate() LLVMContextRef {
	return LLVMContextRef{c: C.LLVMContextCreate()}
}

// LLVMGetGlobalContext Obtain the global context instance.
func LLVMGetGlobalContext() LLVMContextRef {
	return LLVMContextRef{c: C.LLVMGetGlobalContext()}
}

// LLVMContextDispose Destroy a context instance.
// This should be called for every call to LLVMContextCreate() or memory will be leaked.
func LLVMContextDispose(c LLVMContextRef) {
	C.LLVMContextDispose(c.c)
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

func LLVMConstNUWNeg(constantVal LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstNUWNeg(constantVal.c)}
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

func LLVMConstMul(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstMul(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstNSWMul(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstNSWMul(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstNUWMul(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstNUWMul(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstAnd(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstAnd(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstOr(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstOr(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstXor(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstXor(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstICmp(predicate LLVMIntPredicate, lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstICmp(C.LLVMIntPredicate(predicate), lHSConstant.c, rHSConstant.c)}
}

func LLVMConstFCmp(predicate LLVMRealPredicate, lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstFCmp(C.LLVMRealPredicate(predicate), lHSConstant.c, rHSConstant.c)}
}

func LLVMConstShl(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstShl(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstLShr(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstLShr(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstAShr(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstAShr(lHSConstant.c, rHSConstant.c)}
}

func LLVMConstGEP(ty LLVMTypeRef, constantVal LLVMValueRef, constantIndices []LLVMValueRef) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](constantIndices)
	return LLVMValueRef{c: C.LLVMConstGEP2(ty.c, constantVal.c, ptr, length)}
}

func LLVMConstInBoundsGEP(ty LLVMTypeRef, constantVal LLVMValueRef, constantIndices []LLVMValueRef) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](constantIndices)
	return LLVMValueRef{c: C.LLVMConstInBoundsGEP2(ty.c, constantVal.c, ptr, length)}
}

func LLVMConstTrunc(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstTrunc(constantVal.c, toType.c)}
}

func LLVMConstSExt(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstSExt(constantVal.c, toType.c)}
}

func LLVMConstZExt(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstZExt(constantVal.c, toType.c)}
}

func LLVMConstFPTrunc(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstFPTrunc(constantVal.c, toType.c)}
}

func LLVMConstFPExt(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstFPExt(constantVal.c, toType.c)}
}

func LLVMConstUIToFP(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstUIToFP(constantVal.c, toType.c)}
}

func LLVMConstSIToFP(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstSIToFP(constantVal.c, toType.c)}
}

func LLVMConstFPToUI(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstFPToUI(constantVal.c, toType.c)}
}

func LLVMConstFPToSI(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstFPToSI(constantVal.c, toType.c)}
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

func LLVMConstZExtOrBitCast(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstZExtOrBitCast(constantVal.c, toType.c)}
}

func LLVMConstSExtOrBitCast(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstSExtOrBitCast(constantVal.c, toType.c)}
}

func LLVMConstTruncOrBitCast(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstTruncOrBitCast(constantVal.c, toType.c)}
}

func LLVMConstPointerCast(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstPointerCast(constantVal.c, toType.c)}
}

func LLVMConstIntCast(constantVal LLVMValueRef, toType LLVMTypeRef, isSigned bool) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstIntCast(constantVal.c, toType.c, bool2LLVMBool(isSigned))}
}

func LLVMConstFPCast(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstFPCast(constantVal.c, toType.c)}
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

func LLVMSetInitializer(globalVar LLVMValueRef, constantVal LLVMValueRef) {
	C.LLVMSetInitializer(globalVar.c, constantVal.c)
}

func LLVMIsThreadLocal(globalVar LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsThreadLocal(globalVar.c))
}

func LLVMSetThreadLocal(globalVar LLVMValueRef, isThreadLocal bool) {
	C.LLVMSetThreadLocal(globalVar.c, bool2LLVMBool(isThreadLocal))
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

// LLVMGetEntryBasicBlock Obtain the basic block that corresponds to the entry point of a function.
func LLVMGetEntryBasicBlock(fn LLVMValueRef) LLVMBasicBlockRef {
	return LLVMBasicBlockRef{c: C.LLVMGetEntryBasicBlock(fn.c)}
}

// LLVMAppendBasicBlockInContext Append a basic block to the end of a function.
func LLVMAppendBasicBlockInContext(c LLVMContextRef, fn LLVMValueRef, name string) LLVMBasicBlockRef {
	return string2CString(name, func(name *C.char) LLVMBasicBlockRef {
		return LLVMBasicBlockRef{c: C.LLVMAppendBasicBlockInContext(c.c, fn.c, name)}
	})
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

func LLVMBuildNUWNeg(builder LLVMBuilderRef, v LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildNUWNeg(builder.c, v.c, name)}
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
