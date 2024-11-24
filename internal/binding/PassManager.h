#ifndef GOLLVM_BINDINGS_PASSMANAGER_H
#define GOLLVM_BINDINGS_PASSMANAGER_H

#include <llvm-c/Types.h>
#ifdef __cplusplus
#include "llvm/Support/CBindingWrapping.h"
#endif

#ifdef __cplusplus
extern "C" {
#endif

void LLVMOptModule(LLVMModuleRef IR, const char *level);

#ifdef __cplusplus
}
#endif

#endif
