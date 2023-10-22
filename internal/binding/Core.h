#ifndef GOLLVM_BINDINGS_CORE_H
#define GOLLVM_BINDINGS_CORE_H

#include "llvm-c/Core.h"
#ifdef __cplusplus
#include "llvm/Support/CBindingWrapping.h"
#endif

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

LLVMValueRef LLVMConstAggregateZero(LLVMTypeRef ty);

#ifdef __cplusplus
}
#endif

#endif
