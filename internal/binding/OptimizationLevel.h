#ifndef GOLLVM_BINDINGS_OPTIMIZATIONLEVEL_H
#define GOLLVM_BINDINGS_OPTIMIZATIONLEVEL_H

#include "PassManager.h"
#include "llvm-c/TargetMachine.h"
#ifdef __cplusplus
#include "llvm/Support/CBindingWrapping.h"
#endif

#ifdef __cplusplus
extern "C" {
#endif

const char *LLVMOptimizationLevelO0 = "O0";
const char *LLVMOptimizationLevelO1 = "O1";
const char *LLVMOptimizationLevelO2 = "O2";
const char *LLVMOptimizationLevelO3 = "O3";
const char *LLVMOptimizationLevelOz = "Oz";
const char *LLVMOptimizationLevelOs = "Os";

#ifdef __cplusplus
}
#endif

#endif
