package binding

/*
#include "llvm-c/BitWriter.h"
*/
import "C"
import "errors"

// LLVMWriteBitcodeToFile Write a module to the specified path.
func LLVMWriteBitcodeToFile(m LLVMModuleRef, path string) error {
	return string2CString(path, func(path *C.char) error {
		if C.LLVMWriteBitcodeToFile(m.c, path) != 0 {
			return errors.New("failed to write bitcode file")
		}
		return nil
	})
}

// LLVMWriteBitcodeToMemoryBuffer Write a module to a memory buffer owned by the caller.
func LLVMWriteBitcodeToMemoryBuffer(m LLVMModuleRef) LLVMMemoryBufferRef {
	return LLVMMemoryBufferRef{c: C.LLVMWriteBitcodeToMemoryBuffer(m.c)}
}
