package binding

/*
#include "llvm-c/Core.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

// covert LLVMBool to bool
func llvmBool2bool(v C.LLVMBool) bool {
	return v == 1
}

// covert LLVMBool to bool
func bool2LLVMBool(v bool) C.LLVMBool {
	if v {
		return 1
	} else {
		return 0
	}
}

// covert slice to c pointer
func slice2Ptr[T, F any](v []T) (*F, C.unsigned) {
	var ptr *F
	if len(v) > 0 {
		ptr = (*F)(unsafe.Pointer(&v[0]))
	}
	return ptr, C.unsigned(len(v))
}

// covert string to c char *
func string2CString[T any](v string, f func(v *C.char) T) T {
	cstring := C.CString(v)
	defer C.free(unsafe.Pointer(cstring))
	return f(cstring)
}

func llvmError2Error(f func(outError **C.char) C.LLVMBool) error {
	var outError *C.char
	if llvmBool2bool(f(&outError)) {
		defer LLVMDisposeMessage(outError)
		return errors.New(C.GoString(outError))
	}
	return nil
}

type FuncPtr[T any] struct {
	ptr unsafe.Pointer
}

// NewFuncPtr
// 必须是函数指针，不能是lambda和方法
func NewFuncPtr[T any](f unsafe.Pointer) FuncPtr[T] {
	return FuncPtr[T]{ptr: f}
}

func (f FuncPtr[T]) Func() T {
	return *(*T)(unsafe.Pointer(&f.ptr))
}
