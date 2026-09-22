#include "ErrorHandling.h"
#include "llvm-c/ErrorHandling.h"

static void goFatalTrampoline(const char *msg) {
    goLLVMFatalErrorHandler((char *)msg);
}

void llvmInstallGoFatalErrorHandler(void) {
    LLVMInstallFatalErrorHandler(goFatalTrampoline);
}
