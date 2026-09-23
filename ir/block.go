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

// Check 基本块可用性前置校验（nil/所属 Context 或生命周期失效）
func (b Block) Check(op string) {
	if b.ref.IsNil() || b.ctx == nil {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil block")
	}
	b.ctx.CheckAlive(op)
	if b.life == nil || !b.life.Alive() {
		llvm.Panicf(llvm.ErrUseAfterFree, op, "block is freed")
	}
}

// Name 基本块名称
func (b Block) Name() string {
	b.Check("ir.Block.Name")
	return binding.LLVMGetBasicBlockName(b.ref)
}

// SetName 设置基本块名称
func (b Block) SetName(name string) {
	b.Check("ir.Block.SetName")
	binding.LLVMSetValueName(binding.LLVMBasicBlockAsValue(b.ref), name)
}

// Belong 所属函数
func (b Block) Belong() Function {
	b.Check("ir.Block.Belong")
	return Function{Value: wrapValue[llvm.FnT](b.ctx, b.life, binding.LLVMGetBasicBlockParent(b.ref))}
}

// Insts 全部指令（擦除种类）
func (b Block) Insts() []llvm.Value[llvm.DynT] {
	b.Check("ir.Block.Insts")
	var insts []llvm.Value[llvm.DynT]
	for ref := binding.LLVMGetFirstInstruction(b.ref); !ref.IsNil(); ref = binding.LLVMGetNextInstruction(ref) {
		insts = append(insts, wrapDyn(b.ctx, b.life, ref))
	}
	return insts
}

// FirstInst 第一条指令
func (b Block) FirstInst() (llvm.Value[llvm.DynT], bool) {
	b.Check("ir.Block.FirstInst")
	ref := binding.LLVMGetFirstInstruction(b.ref)
	if ref.IsNil() {
		return llvm.Value[llvm.DynT]{}, false
	}
	return wrapDyn(b.ctx, b.life, ref), true
}

// LastInst 最后一条指令
func (b Block) LastInst() (llvm.Value[llvm.DynT], bool) {
	b.Check("ir.Block.LastInst")
	ref := binding.LLVMGetLastInstruction(b.ref)
	if ref.IsNil() {
		return llvm.Value[llvm.DynT]{}, false
	}
	return wrapDyn(b.ctx, b.life, ref), true
}

// Next 下一个基本块
func (b Block) Next() (Block, bool) {
	b.Check("ir.Block.Next")
	ref := binding.LLVMGetNextBasicBlock(b.ref)
	if ref.IsNil() {
		return Block{}, false
	}
	return wrapBlock(b.ctx, b.life, ref), true
}

// Prev 上一个基本块
func (b Block) Prev() (Block, bool) {
	b.Check("ir.Block.Prev")
	ref := binding.LLVMGetPreviousBasicBlock(b.ref)
	if ref.IsNil() {
		return Block{}, false
	}
	return wrapBlock(b.ctx, b.life, ref), true
}

// Empty 是否无指令
func (b Block) Empty() bool {
	b.Check("ir.Block.Empty")
	return binding.LLVMGetFirstInstruction(b.ref).IsNil()
}
