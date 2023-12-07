package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

type Block binding.LLVMBasicBlockRef

func (f Function) NewBlock(name string) Block {
	return Block(binding.LLVMAppendBasicBlockInContext(binding.LLVMGetTypeContext(f.Type().binding()), f.binding(), name))
}

func (b Block) binding() binding.LLVMBasicBlockRef {
	return binding.LLVMBasicBlockRef(b)
}

func (b Block) Belong() Function {
	return Function(binding.LLVMGetBasicBlockParent(b.binding()))
}

func (b Block) Name() string {
	return binding.LLVMGetBasicBlockName(b.binding())
}

func (b Block) GetTerminator() Terminator {
	inst := binding.LLVMGetBasicBlockTerminator(b.binding())
	if inst.IsNil() {
		return nil
	}
	return lookupTerminator(inst)
}

func (b Block) IsTerminating() bool {
	return b.GetTerminator() != nil
}
