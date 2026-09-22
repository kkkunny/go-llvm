package binding

/*
#include "llvm-c/IRReader.h"
*/
import "C"

// LLVMParseIRInContext Read LLVM IR from a memory buffer and convert it into an LLVM module.
// Wraps the non-deprecated LLVMParseIRInContext2; the memory buffer is not consumed.
func LLVMParseIRInContext(c LLVMContextRef, mem LLVMMemoryBufferRef) (LLVMModuleRef, error) {
	var out LLVMModuleRef
	err := llvmError2Error(func(errstr **C.char) C.LLVMBool {
		return C.LLVMParseIRInContext2(c.c, mem.c, &out.c, errstr)
	})
	if err != nil {
		return LLVMModuleRef{}, err
	}
	return out, nil
}
