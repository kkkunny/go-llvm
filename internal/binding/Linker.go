package binding

/*
#include "llvm-c/Linker.h"
*/
import "C"
import "errors"

// LLVMLinkModules Links the source module into the destination module. The source module is destroyed.
// The return value is true if an error occurred, false otherwise.
// Use the diagnostic handler to get any diagnostic message.
func LLVMLinkModules(dest, src LLVMModuleRef) error {
	if llvmBool2bool(C.LLVMLinkModules2(dest.c, src.c)) {
		return errors.New("linking failed")
	}
	return nil
}
