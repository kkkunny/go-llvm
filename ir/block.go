package ir

import (
	"iter"

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
	return Function{Value: llvm.NewValue[llvm.FnT](b.ctx, b.life, binding.LLVMGetBasicBlockParent(b.ref))}
}

// Insts 全部指令（擦除种类）；需要零分配遍历时用 AllInsts
func (b Block) Insts() []llvm.Value[llvm.DynT] {
	b.Check("ir.Block.Insts")
	var insts []llvm.Value[llvm.DynT]
	for ref := binding.LLVMGetFirstInstruction(b.ref); !ref.IsNil(); ref = binding.LLVMGetNextInstruction(ref) {
		insts = append(insts, llvm.ValueOf(b.ctx, b.life, ref))
	}
	return insts
}

// AllInsts 惰性遍历全部指令（range 友好，无切片分配）
func (b Block) AllInsts() iter.Seq[llvm.Value[llvm.DynT]] {
	return func(yield func(llvm.Value[llvm.DynT]) bool) {
		b.Check("ir.Block.AllInsts")
		for ref := binding.LLVMGetFirstInstruction(b.ref); !ref.IsNil(); ref = binding.LLVMGetNextInstruction(ref) {
			if !yield(llvm.ValueOf(b.ctx, b.life, ref)) {
				return
			}
		}
	}
}

// FirstInst 第一条指令
func (b Block) FirstInst() (llvm.Value[llvm.DynT], bool) {
	b.Check("ir.Block.FirstInst")
	ref := binding.LLVMGetFirstInstruction(b.ref)
	if ref.IsNil() {
		return llvm.Value[llvm.DynT]{}, false
	}
	return llvm.ValueOf(b.ctx, b.life, ref), true
}

// LastInst 最后一条指令
func (b Block) LastInst() (llvm.Value[llvm.DynT], bool) {
	b.Check("ir.Block.LastInst")
	ref := binding.LLVMGetLastInstruction(b.ref)
	if ref.IsNil() {
		return llvm.Value[llvm.DynT]{}, false
	}
	return llvm.ValueOf(b.ctx, b.life, ref), true
}

// Terminator 终结指令；块尚未被终结（尚无终结指令）时返回 false。
// 与 [Block.LastInst] 不同，这是"块已被终结"的可靠判断：void 调用等非终结指令不会误判。
// 返回指令可用 [OpOf]、[Successor] 等继续处理；块失效 panic
// [github.com/kkkunny/go-llvm.ErrUseAfterFree]。
func (b Block) Terminator() (llvm.Value[llvm.DynT], bool) {
	b.Check("ir.Block.Terminator")
	ref := binding.LLVMGetBasicBlockTerminator(b.ref)
	if ref.IsNil() {
		return llvm.Value[llvm.DynT]{}, false
	}
	return llvm.ValueOf(b.ctx, b.life, ref), true
}

// IsTerminating 块是否已有终结指令（ret/br/switch/indirectbr/invoke/unreachable/
// resume/cleanupret/catchret/catchswitch/callbr）；没有终结指令的块不能直接交给
// 代码生成/校验，需先补终结指令。
// 要看具体的终结指令用 [Block.Terminator]；块失效 panic
// [github.com/kkkunny/go-llvm.ErrUseAfterFree]。
func (b Block) IsTerminating() bool {
	b.Check("ir.Block.IsTerminating")
	return !binding.LLVMGetBasicBlockTerminator(b.ref).IsNil()
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
