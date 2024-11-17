package llvm

import "C"
import (
	"reflect"

	"github.com/kkkunny/go-llvm/internal/binding"
)

var globalFatalErrorHandler func(string)

func fatalErrorHandler(error *C.char) {
	if globalFatalErrorHandler != nil {
		globalFatalErrorHandler(C.GoString(error))
	}
}

func RegisterFatalErrorHandler(handler func(string)) {
	globalFatalErrorHandler = handler
	if globalFatalErrorHandler == nil {
		binding.LLVMInstallFatalErrorHandler(binding.NewFuncPtr[binding.LLVMFatalErrorHandler](reflect.ValueOf(fatalErrorHandler).UnsafePointer()))
	}
}

func ResetFatalErrorHandler() {
	globalFatalErrorHandler = nil
	binding.LLVMResetFatalErrorHandler()
}

func EnablePrettyStackTrace() {
	binding.LLVMEnablePrettyStackTrace()
}
