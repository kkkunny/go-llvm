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

// llvmBool2bool converts LLVMBool to bool
func llvmBool2bool(v C.LLVMBool) bool {
	return v == 1
}

// bool2LLVMBool converts bool to LLVMBool
func bool2LLVMBool(v bool) C.LLVMBool {
	if v {
		return 1
	}
	return 0
}

// slice2Ptr converts a slice to a C pointer and element count
func slice2Ptr[T, F any](v []T) (*F, C.unsigned) {
	var ptr *F
	if len(v) > 0 {
		ptr = (*F)(unsafe.Pointer(&v[0]))
	}
	return ptr, C.unsigned(len(v))
}

// string2CString passes a temporary C string to f and frees it afterwards
func string2CString[T any](v string, f func(v *C.char) T) T {
	cstring := C.CString(v)
	defer C.free(unsafe.Pointer(cstring))
	return f(cstring)
}

// llvmError2Error converts the out-error pattern of LLVM-C APIs to a Go error
func llvmError2Error(f func(outError **C.char) C.LLVMBool) error {
	var outError *C.char
	if llvmBool2bool(f(&outError)) {
		defer LLVMDisposeMessage(outError)
		return errors.New(C.GoString(outError))
	}
	return nil
}

// FuncPtr wraps a C function pointer; T must be a function type.
type FuncPtr[T any] struct {
	ptr unsafe.Pointer
}

// NewFuncPtr wraps a C function pointer; f must point to a function, not a lambda or method.
func NewFuncPtr[T any](f unsafe.Pointer) FuncPtr[T] {
	return FuncPtr[T]{ptr: f}
}

// Func returns the typed function pointer.
func (f FuncPtr[T]) Func() T {
	return *(*T)(unsafe.Pointer(&f.ptr))
}
