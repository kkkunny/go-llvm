package binding

/*
#include "Core.h"
#include "OptimizationLevel.h"
*/
import "C"

type LLVMOptimizationLevel *C.char

var (
	LLVMOptimizationLevelO0 LLVMOptimizationLevel = C.LLVMOptimizationLevelO0
	LLVMOptimizationLevelO1 LLVMOptimizationLevel = C.LLVMOptimizationLevelO1
	LLVMOptimizationLevelO2 LLVMOptimizationLevel = C.LLVMOptimizationLevelO2
	LLVMOptimizationLevelO3 LLVMOptimizationLevel = C.LLVMOptimizationLevelO3
	LLVMOptimizationLevelOz LLVMOptimizationLevel = C.LLVMOptimizationLevelOz
	LLVMOptimizationLevelOs LLVMOptimizationLevel = C.LLVMOptimizationLevelOs
)
