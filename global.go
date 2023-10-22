package llvm

import (
	"github.com/samber/lo"

	"github.com/kkkunny/go-llvm/internal/binding"
)

type Global interface {
	binding() binding.LLVMValueRef
	Module() Module
}

type Function binding.LLVMValueRef

func (m Module) NewFunction(name string, t FunctionType) Function {
	return Function(binding.LLVMAddFunction(m.binding(), name, t.binding()))
}

func (m Module) GetFunction(name string) Function {
	return Function(binding.LLVMGetNamedFunction(m.binding(), name))
}

func (c Function) String() string {
	return binding.LLVMPrintValueToString(c.binding())
}

func (c Function) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(c)
}

func (c Function) Type() Type {
	return lookupType(binding.LLVMTypeOf(c.binding()))
}

func (g Function) Module() Module {
	return Module(binding.LLVMGetGlobalParent(g.binding()))
}

func (f Function) Blocks() []Block {
	return lo.Map(binding.LLVMGetBasicBlocks(f.binding()), func(item binding.LLVMBasicBlockRef, index int) Block {
		return Block(item)
	})
}

func (f Function) EntryBlock() Block {
	return Block(binding.LLVMGetEntryBasicBlock(f.binding()))
}

type GlobalValue binding.LLVMValueRef

func (m Module) NewGlobal(name string, t Type) GlobalValue {
	return GlobalValue(binding.LLVMAddGlobal(m.binding(), t.binding(), name))
}

func (m Module) GetGlobal(name string) GlobalValue {
	return GlobalValue(binding.LLVMGetNamedGlobal(m.binding(), name))
}

func (v GlobalValue) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

func (v GlobalValue) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(v)
}

func (v GlobalValue) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (g GlobalValue) Module() Module {
	return Module(binding.LLVMGetGlobalParent(g.binding()))
}

func (g GlobalValue) ValueType() Type {
	return lookupType(binding.LLVMGlobalGetValueType(g.binding()))
}
