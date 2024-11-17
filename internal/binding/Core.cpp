#include "Core.h"
#include "llvm/IR/Constants.h"
#include "llvm/IR/Function.h"

using namespace llvm;

LLVMValueRef LLVMConstAggregateZero(LLVMTypeRef ty) {
    return wrap(llvm::ConstantAggregateZero::get(unwrap(ty)));
}

LLVMTypeRef LLVMGetFunctionType(LLVMValueRef f) {
    return wrap(unwrap<Function>(f)->getFunctionType());
}

void LLVMSetDSOLocal(LLVMValueRef v, LLVMBool Local) {
    unwrap<Function>(v)->setDSOLocal(Local);
}

LLVMBool LLVMIsDSOLocal(LLVMValueRef v) {
    return unwrap<Function>(v)->isDSOLocal();
}
