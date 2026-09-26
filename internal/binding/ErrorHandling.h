#ifndef GOLLVM_BINDINGS_ERRORHANDLING_H
#define GOLLVM_BINDINGS_ERRORHANDLING_H

#include "llvm-c/Core.h"
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

extern void goLLVMFatalErrorHandler(char *);

void llvmInstallGoFatalErrorHandler(void);

extern void goLLVMDiagnosticHandler(uintptr_t id, int severity, char *msg);

void llvmInstallGoDiagnosticHandler(LLVMContextRef C, uintptr_t id);

void llvmClearGoDiagnosticHandler(LLVMContextRef C);

#ifdef __cplusplus
}
#endif

#endif
