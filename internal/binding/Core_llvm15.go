//go:build llvm15

package binding

/*
#include "llvm-c/Core.h"
#include "Core.h"
*/
import "C"

func LLVMConstFNeg(constantVal LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstFNeg(constantVal.c)}
}
