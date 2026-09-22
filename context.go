package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// GlobalContext 全局共享上下文（LLVM22起不再使用C全局上下文，改为进程内自有上下文）
var GlobalContext = NewContext()

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
