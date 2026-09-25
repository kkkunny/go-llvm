package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// InlineAsmDialect 内联汇编方言
type InlineAsmDialect binding.LLVMInlineAsmDialect

// InlineAsmDialect 取值对应 LLVM 内联汇编方言（binding.LLVMInlineAsmDialect）。
// 决定 [Context.InlineAsm] 的 asm 字符串按哪种汇编语法书写。
const (
	InlineAsmATT   = InlineAsmDialect(binding.LLVMInlineAsmDialectATT)   // AT&T 方言：寄存器带 % 前缀，操作数顺序为源、目标
	InlineAsmIntel = InlineAsmDialect(binding.LLVMInlineAsmDialectIntel) // Intel 方言：寄存器不带前缀，操作数顺序为目标、源
)

// InlineAsm 构造内联汇编常量（函数类型的值，可经 Builder.Call 调用）
func (ctx *Context) InlineAsm(fn FnType, asm, constraints string, sideEffects, alignStack bool, dialect InlineAsmDialect, canThrow bool) Value[FnT] {
	const op = "llvm.Context.InlineAsm"
	ctx.CheckType(op, fn)
	ref := binding.LLVMGetInlineAsm(fn.Ref(), asm, constraints, sideEffects, alignStack, binding.LLVMInlineAsmDialect(dialect), canThrow)
	return newValue[FnT](ctx, ctx.life, ref)
}
