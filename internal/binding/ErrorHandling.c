#include "ErrorHandling.h"
#include "llvm-c/ErrorHandling.h"
#include "llvm-c/Core.h"
#include <stddef.h>

static void goFatalTrampoline(const char *msg) {
    goLLVMFatalErrorHandler((char *)msg);
}

void llvmInstallGoFatalErrorHandler(void) {
    LLVMInstallFatalErrorHandler(goFatalTrampoline);
}

static void goDiagnosticTrampoline(LLVMDiagnosticInfoRef DI, void *Ctx) {
    uintptr_t id = (uintptr_t)Ctx;
    int severity = (int)LLVMGetDiagInfoSeverity(DI);
    char *msg = LLVMGetDiagInfoDescription(DI);
    goLLVMDiagnosticHandler(id, severity, msg);
    LLVMDisposeMessage(msg);
}

void llvmInstallGoDiagnosticHandler(LLVMContextRef C, uintptr_t id) {
    LLVMContextSetDiagnosticHandler(C, goDiagnosticTrampoline, (void *)id);
}

void llvmClearGoDiagnosticHandler(LLVMContextRef C) {
    LLVMContextSetDiagnosticHandler(C, NULL, NULL);
}
