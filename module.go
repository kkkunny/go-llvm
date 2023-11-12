package llvm

import (
	"errors"
	"strings"

	"github.com/kkkunny/go-llvm/internal/binding"
)

type Module binding.LLVMModuleRef

func (ctx Context) NewModule(name string) Module {
	return Module(binding.LLVMModuleCreateWithNameInContext(name, ctx.binding()))
}

func (m Module) binding() binding.LLVMModuleRef {
	return binding.LLVMModuleRef(m)
}

func (m Module) Free() {
	binding.LLVMDisposeModule(m.binding())
}

func (m Module) String() string {
	return binding.LLVMPrintModuleToString(m.binding())
}

func (m Module) Clone() Module {
	return Module(binding.LLVMCloneModule(m.binding()))
}

func (m Module) Context() Context {
	return Context(binding.LLVMGetModuleContext(m.binding()))
}

func (m Module) GetSource() string {
	return binding.LLVMGetSourceFileName(m.binding())
}

func (m Module) SetSource(source string) {
	binding.LLVMSetSourceFileName(m.binding(), source)
}

func (m Module) RunPasses(target *Target, option PassOption, pass ...string) error {
	return binding.LLVMRunPasses(m.binding(), strings.Join(pass, ","), target.machine, option.binding())
}

func (m Module) Verify() error {
	msg, fail := binding.LLVMVerifyModule(m.binding(), binding.LLVMReturnStatusAction)
	if fail {
		return errors.New(msg)
	}
	return nil
}
