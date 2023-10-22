package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

type Builder binding.LLVMBuilderRef

func (ctx Context) NewBuilder() Builder {
	return Builder(binding.LLVMCreateBuilderInContext(ctx.binding()))
}

func (b Builder) binding() binding.LLVMBuilderRef {
	return binding.LLVMBuilderRef(b)
}

func (b Builder) Free() {
	binding.LLVMDisposeBuilder(b.binding())
}

func (b Builder) MoveToAfter(block Block) {
	binding.LLVMPositionBuilderAtEnd(b.binding(), block.binding())
}

func (b Builder) CurrentBlock() Block {
	return Block(binding.LLVMGetInsertBlock(b.binding()))
}
