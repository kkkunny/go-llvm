package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// InlineAsmDialect 内联汇编方言
type InlineAsmDialect binding.LLVMInlineAsmDialect

const (
	InlineAsmATT   = InlineAsmDialect(binding.LLVMInlineAsmDialectATT)
	InlineAsmIntel = InlineAsmDialect(binding.LLVMInlineAsmDialectIntel)
)

// InlineAsm 构造内联汇编常量（函数类型的值，可经 Builder.Call 调用）
func (ctx *Context) InlineAsm(fn FnType, asm, constraints string, sideEffects, alignStack bool, dialect InlineAsmDialect, canThrow bool) Value[FnT] {
	const op = "llvm.Context.InlineAsm"
	ctx.CheckType(op, fn)
	ref := binding.LLVMGetInlineAsm(fn.Ref(), asm, constraints, sideEffects, alignStack, binding.LLVMInlineAsmDialect(dialect), canThrow)
	return newValue[FnT](ctx, ctx.life, ref)
}
