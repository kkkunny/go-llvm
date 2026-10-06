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
void LLVMSetDSOLocal(LLVMValueRef v, LLVMBool Local);
LLVMBool LLVMIsDSOLocal(LLVMValueRef v);
void LLVMGoEmitError(LLVMContextRef c, const char *msg);

// LLVMGoConstSExt/LLVMGoConstZExt: ConstantExpr::getSExt/getZExt (missing from LLVM-C).
LLVMValueRef LLVMGoConstSExt(LLVMValueRef v, LLVMTypeRef ty);
LLVMValueRef LLVMGoConstZExt(LLVMValueRef v, LLVMTypeRef ty);

// Checked APInt accessors: return 1 and store the value when it fits in 64 bits,
// return 0 otherwise (the plain LLVM-C getters assert on out-of-range values).
int LLVMGoConstIntGetSExtValue(LLVMValueRef v, long long *out);
int LLVMGoConstIntGetZExtValue(LLVMValueRef v, unsigned long long *out);

// LLVMGoCreateTargetData: DataLayout::parse (non-fatal). On failure returns NULL and
// stores a malloc'd message in *ErrMsg (free with LLVMDisposeMessage).
LLVMTargetDataRef LLVMGoCreateTargetData(const char *StringRep, char **ErrMsg);

#ifdef __cplusplus
}
#endif

#endif
