package binding

/*
#include "llvm-c/Error.h"
*/
import "C"

// LLVMErrorRef Opaque reference to an error instance. Null serves as the 'success' value.
type LLVMErrorRef struct{ c C.LLVMErrorRef }
type LLVMErrorTypeId struct{ c C.LLVMErrorTypeId }

// LLVMGetErrorTypeId Returns the type id for the given error instance, which must be a failure value (i.e. non-null).
func LLVMGetErrorTypeId(err LLVMErrorRef) LLVMErrorTypeId {
	return LLVMErrorTypeId{c: C.LLVMGetErrorTypeId(err.c)}
}

// LLVMConsumeError Dispose of the given error without handling it.
// This operation consumes the error, and the given LLVMErrorRef value is not usable once this call returns.
// Note: This method *only* needs to be called if the error is not being passed to some other consuming operation, e.g. LLVMGetErrorMessage.
func LLVMConsumeError(err LLVMErrorRef) {
	C.LLVMConsumeError(err.c)
}

// LLVMGetErrorMessage Returns the given string's error message. This operation consumes the error, and the given LLVMErrorRef value is not usable once this call returns.
func LLVMGetErrorMessage(err LLVMErrorRef) string {
	cptr := C.LLVMGetErrorMessage(err.c)
	defer C.LLVMDisposeErrorMessage(cptr)
	return C.GoString(cptr)
}

// LLVMGetStringErrorTypeId Returns the type id for llvm StringError.
func LLVMGetStringErrorTypeId() LLVMErrorTypeId {
	return LLVMErrorTypeId{c: C.LLVMGetStringErrorTypeId()}
}

// LLVMCreateStringError Create a StringError.
func LLVMCreateStringError(errMsg string) LLVMErrorRef {
	return string2CString(errMsg, func(errMsg *C.char) LLVMErrorRef {
		return LLVMErrorRef{c: C.LLVMCreateStringError(errMsg)}
	})
}
