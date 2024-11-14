//go:build llvm15 || llvm16 || llvm17 || llvm18

package binding

/*
#include "llvm-c/Core.h"
#include "Core.h"
*/
import "C"

// Deprecated
func LLVMConstFCmp(predicate LLVMRealPredicate, lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstFCmp(C.LLVMRealPredicate(predicate), lHSConstant.c, rHSConstant.c)}
}

// Deprecated
func LLVMConstICmp(predicate LLVMIntPredicate, lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstICmp(C.LLVMIntPredicate(predicate), lHSConstant.c, rHSConstant.c)}
}

// Deprecated
func LLVMConstShl(lHSConstant, rHSConstant LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstShl(lHSConstant.c, rHSConstant.c)}
}

func LLVMBuildNUWNeg(builder LLVMBuilderRef, v LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildNUWNeg(builder.c, v.c, name)}
	})
}

func LLVMConstNUWNeg(constantVal LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstNUWNeg(constantVal.c)}
}
