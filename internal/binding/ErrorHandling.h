#ifndef GOLLVM_BINDINGS_ERRORHANDLING_H
#define GOLLVM_BINDINGS_ERRORHANDLING_H

#ifdef __cplusplus
extern "C" {
#endif

extern void goLLVMFatalErrorHandler(char *);

void llvmInstallGoFatalErrorHandler(void);

#ifdef __cplusplus
}
#endif

#endif
