//go:build llvm15 || llvm16 || llvm17

package binding

/*
#include "llvm-c/Core.h"
#include "Core.h"
*/
import "C"

// Deprecated
func LLVMConstAnd(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstAnd(lHSConstant.c, rHSConstant.c)}
}

// Deprecated
func LLVMConstOr(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstOr(lHSConstant.c, rHSConstant.c)}
}

// Deprecated
func LLVMConstLShr(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstLShr(lHSConstant.c, rHSConstant.c)}
}

// Deprecated
func LLVMConstAShr(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstAShr(lHSConstant.c, rHSConstant.c)}
}

// Deprecated
func LLVMConstZExt(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstZExt(constantVal.c, toType.c)}
}

// Deprecated
func LLVMConstSExt(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstSExt(constantVal.c, toType.c)}
}

// Deprecated
func LLVMConstFPTrunc(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstFPTrunc(constantVal.c, toType.c)}
}

// Deprecated
func LLVMConstFPExt(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstFPExt(constantVal.c, toType.c)}
}

// Deprecated
func LLVMConstFPToUI(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstFPToUI(constantVal.c, toType.c)}
}

// Deprecated
func LLVMConstFPToSI(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstFPToSI(constantVal.c, toType.c)}
}

// Deprecated
func LLVMConstUIToFP(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstUIToFP(constantVal.c, toType.c)}
}

// Deprecated
func LLVMConstSIToFP(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstSIToFP(constantVal.c, toType.c)}
}

// Deprecated
func LLVMConstFPCast(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstFPCast(constantVal.c, toType.c)}
}

// Deprecated
func LLVMConstIntCast(constantVal LLVMValueRef, toType LLVMTypeRef, isSigned bool) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstIntCast(constantVal.c, toType.c, bool2LLVMBool(isSigned))}
}

// Deprecated
func LLVMConstSExtOrBitCast(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstSExtOrBitCast(constantVal.c, toType.c)}
}

// Deprecated
func LLVMConstZExtOrBitCast(constantVal LLVMValueRef, toType LLVMTypeRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstZExtOrBitCast(constantVal.c, toType.c)}
}
