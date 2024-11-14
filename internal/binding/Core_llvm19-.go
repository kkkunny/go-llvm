//go:build !llvm15 && !llvm16 && !llvm17 && !llvm18

package binding

/*
#include "llvm-c/Core.h"
#include "Core.h"
*/
import "C"

func LLVMSetNUW(arithInst LLVMValueRef, hasNUW bool) {
	C.LLVMSetNUW(arithInst.c, bool2LLVMBool(hasNUW))
}

func LLVMBuildNUWNeg(builder LLVMBuilderRef, v LLVMValueRef, name string) LLVMValueRef {
	ref := LLVMBuildNeg(builder, v, name)
	LLVMSetNUW(ref, true)
	return ref
}

func LLVMConstNUWNeg(constantVal LLVMValueRef) LLVMValueRef {
	return LLVMConstNull(LLVMGlobalGetValueType(constantVal))
}
