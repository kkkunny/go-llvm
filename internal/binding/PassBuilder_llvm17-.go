//go:build llvm17

package binding

/*
#include "llvm-c/Transforms/PassBuilder.h"
*/
import "C"

func LLVMPassBuilderOptionsSetInlinerThreshold(options LLVMPassBuilderOptionsRef, threshold int32) {
	C.LLVMPassBuilderOptionsSetInlinerThreshold(options.c, C.int(threshold))
}
