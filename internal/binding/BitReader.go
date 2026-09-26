package binding

/*
#include "llvm-c/BitReader.h"
*/
import "C"

// LLVMParseBitcodeInContext Read bitcode from a memory buffer and convert it to an LLVM module.
// The memory buffer is not consumed; the caller keeps ownership.
func LLVMParseBitcodeInContext(c LLVMContextRef, mem LLVMMemoryBufferRef) (LLVMModuleRef, error) {
	var out LLVMModuleRef
	err := llvmError2Error(func(errstr **C.char) C.LLVMBool {
		return C.LLVMParseBitcodeInContext(c.c, mem.c, &out.c, errstr)
	})
	if err != nil {
		return LLVMModuleRef{}, err
	}
	return out, nil
}
