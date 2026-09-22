package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Block 基本块
type Block struct {
	ref  binding.LLVMBasicBlockRef
	ctx  *llvm.Context
	life *llvm.Lifetime
}

// Ref 返回底层句柄
func (b Block) Ref() binding.LLVMBasicBlockRef { return b.ref }

// Name 基本块名称
func (b Block) Name() string { return binding.LLVMGetBasicBlockName(b.ref) }

// SetName 设置基本块名称
func (b Block) SetName(name string) {
	binding.LLVMSetValueName(binding.LLVMBasicBlockAsValue(b.ref), name)
}

// Belong 所属函数
func (b Block) Belong() Function {
	return Function{v: llvm.NewValue[llvm.FnT](b.ctx, b.life, binding.LLVMGetBasicBlockParent(b.ref))}
}

// Insts 全部指令（擦除种类）
func (b Block) Insts() []llvm.Value[llvm.DynT] {
	var insts []llvm.Value[llvm.DynT]
	for ref := binding.LLVMGetFirstInstruction(b.ref); !ref.IsNil(); ref = binding.LLVMGetNextInstruction(ref) {
		insts = append(insts, llvm.ValueOf(b.ctx, b.life, ref))
	}
	return insts
}

// FirstInst 第一条指令
func (b Block) FirstInst() (llvm.Value[llvm.DynT], bool) {
	ref := binding.LLVMGetFirstInstruction(b.ref)
	if ref.IsNil() {
		return llvm.Value[llvm.DynT]{}, false
	}
	return llvm.ValueOf(b.ctx, b.life, ref), true
}

// LastInst 最后一条指令
func (b Block) LastInst() (llvm.Value[llvm.DynT], bool) {
	ref := binding.LLVMGetLastInstruction(b.ref)
	if ref.IsNil() {
		return llvm.Value[llvm.DynT]{}, false
	}
	return llvm.ValueOf(b.ctx, b.life, ref), true
}

// Next 下一个基本块
func (b Block) Next() (Block, bool) {
	ref := binding.LLVMGetNextBasicBlock(b.ref)
	if ref.IsNil() {
		return Block{}, false
	}
	return Block{ref: ref, ctx: b.ctx, life: b.life}, true
}

// Prev 上一个基本块
func (b Block) Prev() (Block, bool) {
	ref := binding.LLVMGetPreviousBasicBlock(b.ref)
	if ref.IsNil() {
		return Block{}, false
	}
	return Block{ref: ref, ctx: b.ctx, life: b.life}, true
}

// Empty 是否无指令
func (b Block) Empty() bool { return binding.LLVMGetFirstInstruction(b.ref).IsNil() }
