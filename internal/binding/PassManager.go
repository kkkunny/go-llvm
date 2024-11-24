package binding

/*
#include "PassManager.h"
*/
import "C"

func LLVMOptModule(ir LLVMModuleRef, level LLVMOptimizationLevel) {
	C.LLVMOptModule(ir.c, level)
}
