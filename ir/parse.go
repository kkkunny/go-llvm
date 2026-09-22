package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// ParseIR 从内存缓冲解析 LLVM IR 文本；缓冲不被消费，仍由调用方 Close。
// 解析失败返回 ErrParse 与完整诊断。
func ParseIR(ctx *llvm.Context, buf *llvm.MemoryBuffer) (*Module, error) {
	const op = "ir.ParseIR"
	if !ctx.Alive() {
		errPanic(llvm.ErrUseAfterFree, op, "context is closed")
	}
	if buf == nil {
		errPanic(llvm.ErrInvalidArg, op, "nil memory buffer")
	}
	if !buf.Alive() {
		errPanic(llvm.ErrUseAfterFree, op, "memory buffer is closed")
	}
	ref, err := binding.LLVMParseIRInContext(ctx.Ref(), buf.Ref())
	if err != nil {
		return nil, llvm.WrapError(llvm.ErrParse, op, err)
	}
	return newModule(ctx, ref), nil
}

// ParseBitcode 从内存缓冲解析 bitcode；缓冲不被消费，仍由调用方 Close。
// 解析失败返回 ErrParse 与完整诊断。
func ParseBitcode(ctx *llvm.Context, buf *llvm.MemoryBuffer) (*Module, error) {
	const op = "ir.ParseBitcode"
	if !ctx.Alive() {
		errPanic(llvm.ErrUseAfterFree, op, "context is closed")
	}
	if buf == nil {
		errPanic(llvm.ErrInvalidArg, op, "nil memory buffer")
	}
	if !buf.Alive() {
		errPanic(llvm.ErrUseAfterFree, op, "memory buffer is closed")
	}
	ref, err := binding.LLVMParseBitcodeInContext(ctx.Ref(), buf.Ref())
	if err != nil {
		return nil, llvm.WrapError(llvm.ErrParse, op, err)
	}
	return newModule(ctx, ref), nil
}
