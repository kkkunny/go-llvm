package binding

/*
#include "llvm-c/Analysis.h"
*/
import "C"

type LLVMVerifierFailureAction int32

const (
	// LLVMAbortProcessAction verifier will print to stderr and abort()
	LLVMAbortProcessAction LLVMVerifierFailureAction = iota
	// LLVMPrintMessageAction verifier will print to stderr and return true
	LLVMPrintMessageAction
	// LLVMReturnStatusAction verifier will just return true
	LLVMReturnStatusAction
)

// LLVMVerifyModule Verifies that a module is valid, taking the specified action if not.
// Optionally returns a human-readable description of any invalid constructs.
func LLVMVerifyModule(m LLVMModuleRef, action LLVMVerifierFailureAction) (string, bool) {
	var outMessage *C.char
	res := llvmBool2bool(C.LLVMVerifyModule(m.c, C.LLVMVerifierFailureAction(action), &outMessage))
	defer LLVMDisposeMessage(outMessage)
	return C.GoString(outMessage), res
}

// LLVMVerifyFunction Verifies that a single function is valid, taking the specified action.
// Useful for debugging.
func LLVMVerifyFunction(fn LLVMValueRef, action LLVMVerifierFailureAction) bool {
	return llvmBool2bool(C.LLVMVerifyFunction(fn.c, C.LLVMVerifierFailureAction(action)))
}

// LLVMViewFunctionCFG Open up a ghostview window that displays the CFG of the current function.
// Useful for debugging.
func LLVMViewFunctionCFG(fn LLVMValueRef) {
	C.LLVMViewFunctionCFG(fn.c)
}

// LLVMViewFunctionCFGOnly Open up a ghostview window that displays the CFG of the current function.
// Useful for debugging.
func LLVMViewFunctionCFGOnly(fn LLVMValueRef) {
	C.LLVMViewFunctionCFGOnly(fn.c)
}
