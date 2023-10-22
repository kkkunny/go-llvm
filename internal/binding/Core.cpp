#include "Core.h"
#include "llvm/IR/Constants.h"

using namespace llvm;

LLVMValueRef LLVMConstAggregateZero(LLVMTypeRef ty) {
    return wrap(llvm::ConstantAggregateZero::get(unwrap(ty)));
}
