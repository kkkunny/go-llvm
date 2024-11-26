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

func (b Block) RemoveFromBelong() {
	binding.LLVMRemoveBasicBlockFromParent(b.binding())
}

func (b Block) FirstInst() (Instruction, bool) {
	i := binding.LLVMGetFirstInstruction(b.binding())
	if i.IsNil() {
		return nil, false
	}
	return lookupInstruction(i), true
}

func (b Block) LastInst() (Instruction, bool) {
	i := binding.LLVMGetLastInstruction(b.binding())
	if i.IsNil() {
		return nil, false
	}
	return lookupInstruction(i), true
}

func (b Block) ForeachInst(cb func(inst Instruction)) {
	for inst, ok := b.FirstInst(); ok; inst, ok = inst.Next() {
		cb(inst)
	}
}

func (b Block) Instructions() []Instruction {
	var insts []Instruction
	b.ForeachInst(func(inst Instruction) {
		insts = append(insts, inst)
	})
	return insts
}

func (b Block) Next() (Block, bool) {
	ref := binding.LLVMGetNextBasicBlock(b.binding())
	if ref.IsNil() {
		return Block{}, false
	}
	return Block(ref), true
}

func (b Block) Prev() (Block, bool) {
	ref := binding.LLVMGetPreviousBasicBlock(b.binding())
	if ref.IsNil() {
		return Block{}, false
	}
	return Block(ref), true
}

func (b Block) Empty() bool {
	_, ok := b.FirstInst()
	return !ok
}
