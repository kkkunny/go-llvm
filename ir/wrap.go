package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// wrapBlock 由底层句柄铸造基本块。
// ir 内所有 Block 构造统一走此入口，lifetime 一律取所属 Module 的令牌，
// 避免同一块从不同入口铸造时携带不同 lifetime。
func wrapBlock(ctx *llvm.Context, life *llvm.Lifetime, ref binding.LLVMBasicBlockRef) Block {
	return Block{ref: ref, ctx: ctx, life: life}
}
