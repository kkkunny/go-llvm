#ifndef GOLLVM_BINDINGS_CORE_H
#define GOLLVM_BINDINGS_CORE_H

#include "llvm-c/Core.h"
#include "llvm-c/Transforms/PassBuilder.h"
#ifdef __cplusplus
#include "llvm/Support/CBindingWrapping.h"
#endif

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

LLVMValueRef LLVMConstAggregateZero(LLVMTypeRef ty);
LLVMTypeRef LLVMGetFunctionType(LLVMValueRef f);

#ifdef __cplusplus
}
#endif

#endif
