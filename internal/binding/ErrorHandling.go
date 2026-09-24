package binding

/*
#include "ErrorHandling.h"
#include "llvm-c/ErrorHandling.h"
#include <stdlib.h>
*/
import "C"
import "sync"

type LLVMFatalErrorHandler func(*C.char)

var goFatalHandler func(string)

//export goLLVMFatalErrorHandler
func goLLVMFatalErrorHandler(msg *C.char) {
	if goFatalHandler != nil {
		goFatalHandler(C.GoString(msg))
	}
}

// diagnosticHandlers 诊断回调注册表：Context id → Go 回调（可能被任意 goroutine 触发）
var diagnosticHandlers sync.Map

//export goLLVMDiagnosticHandler
func goLLVMDiagnosticHandler(id C.uintptr_t, severity C.int, msg *C.char) {
	if h, ok := diagnosticHandlers.Load(uint64(id)); ok {
		h.(func(LLVMDiagnosticSeverity, string))(LLVMDiagnosticSeverity(severity), C.GoString(msg))
	}
}

// LLVMContextSetDiagnosticHandlerGo 安装 Go 侧诊断回调；id 由调用方保证唯一
func LLVMContextSetDiagnosticHandlerGo(c LLVMContextRef, id uint64, handler func(LLVMDiagnosticSeverity, string)) {
	diagnosticHandlers.Store(id, handler)
	C.llvmInstallGoDiagnosticHandler(c.c, C.uintptr_t(id))
}

// LLVMContextClearDiagnosticHandlerGo 移除诊断回调
func LLVMContextClearDiagnosticHandlerGo(c LLVMContextRef, id uint64) {
	diagnosticHandlers.Delete(id)
	C.llvmClearGoDiagnosticHandler(c.c)
}

// LLVMInstallFatalErrorHandler Install a fatal error handler. By default, if LLVM detects a fatal error, it will call exit(1).
// This may not be appropriate in many contexts.
// For example, doing exit(1) will bypass many crash reporting/tracing system tools.
// This function allows you to install a callback that will be invoked prior to the call to exit(1).
func LLVMInstallFatalErrorHandler(handler FuncPtr[LLVMFatalErrorHandler]) {
	C.LLVMInstallFatalErrorHandler((C.LLVMFatalErrorHandler)(handler.ptr))
}

// LLVMInstallFatalErrorHandlerGo 安装 Go 侧 fatal error 回调
func LLVMInstallFatalErrorHandlerGo(handler func(string)) {
	goFatalHandler = handler
	C.llvmInstallGoFatalErrorHandler()
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

// LLVMContextHasDiagnosticHandler 是否安装了诊断回调
func LLVMContextHasDiagnosticHandler(c LLVMContextRef) bool {
	return C.LLVMContextGetDiagnosticHandler(c.c) != nil
}
