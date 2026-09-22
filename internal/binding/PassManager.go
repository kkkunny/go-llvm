package binding

/*
#include "PassManager.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

// LLVMOptModule 运行默认优化管线；失败返回 C++ 异常消息
func LLVMOptModule(ir LLVMModuleRef, level LLVMOptimizationLevel) error {
	clevel := C.CString(string(level))
	defer C.free(unsafe.Pointer(clevel))
	var outError *C.char
	if C.LLVMOptModule(ir.c, clevel, &outError) != 0 {
		defer C.free(unsafe.Pointer(outError))
		return errors.New(C.GoString(outError))
	}
	return nil
}
