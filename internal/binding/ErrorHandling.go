package binding

/*
#include "llvm-c/ErrorHandling.h"
*/
import "C"

type LLVMFatalErrorHandler func(*C.char)

// LLVMInstallFatalErrorHandler Install a fatal error handler. By default, if LLVM detects a fatal error, it will call exit(1).
// This may not be appropriate in many contexts.
// For example, doing exit(1) will bypass many crash reporting/tracing system tools.
// This function allows you to install a callback that will be invoked prior to the call to exit(1).
func LLVMInstallFatalErrorHandler(handler FuncPtr[LLVMFatalErrorHandler]) {
	C.LLVMInstallFatalErrorHandler((C.LLVMFatalErrorHandler)(handler.ptr))
}

// LLVMResetFatalErrorHandler Reset the fatal error handler.
// This resets LLVM's fatal error handling behavior to the default.
func LLVMResetFatalErrorHandler() {
	C.LLVMResetFatalErrorHandler()
}

// LLVMEnablePrettyStackTrace Enable LLVM's built-in stack trace code.
// This intercepts the OS's crash signals and prints which component of LLVM you were in at the time if the crash.
func LLVMEnablePrettyStackTrace() {
	C.LLVMEnablePrettyStackTrace()
}
