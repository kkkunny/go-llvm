//go:build llvm15 || llvm16

package binding

/*
#include "llvm-c/Core.h"
#include "Core.h"
*/
import "C"

// Deprecated
func LLVMConstSelect(constantCondition, constantIfTrue, constantIfFalse LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMConstSelect(constantCondition.c, constantIfTrue.c, constantIfFalse.c)}
}
