package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

var GlobalContext = Context(binding.LLVMGetGlobalContext())

type Context binding.LLVMContextRef

func NewContext() Context {
	return Context(binding.LLVMContextCreate())
}

func (ctx Context) binding() binding.LLVMContextRef {
	return binding.LLVMContextRef(ctx)
}

func (ctx Context) Free() {
	binding.LLVMContextDispose(ctx.binding())
}
