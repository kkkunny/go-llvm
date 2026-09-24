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

// emptyCString 静态空字符串缓冲（Name 等"可为空"的 C 参数用）。
// LLVM 对这类参数只读取不保留；用静态缓冲避免每次调用的 C 堆分配与两次额外 cgo 穿越。
var emptyCString = C.CString("")

// string2CString passes a temporary C string to f and frees it afterwards.
// 空串直接使用静态缓冲；非空串仍走 CString/free 保证 NUL 终止与生命周期安全。
func string2CString[T any](v string, f func(v *C.char) T) T {
	if len(v) == 0 {
		return f(emptyCString)
	}
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
