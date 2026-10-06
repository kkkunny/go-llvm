package ir

import (
	"runtime"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/errs"
)

// ParseIR 从内存缓冲解析 LLVM IR 文本；缓冲不被消费，仍由调用方 Close。
// 解析失败返回 ErrParse 与完整诊断。
func ParseIR(ctx *llvm.Context, buf *llvm.MemoryBuffer) (*Module, error) {
	const op = "ir.ParseIR"
	if !ctx.Alive() {
		errs.Panicf(llvm.ErrUseAfterFree, op, "context is closed")
	}
	if buf == nil {
		errs.Panicf(llvm.ErrInvalidArg, op, "nil memory buffer")
	}
	if !buf.Alive() {
		errs.Panicf(llvm.ErrUseAfterFree, op, "memory buffer is closed")
	}
	ref, err := binding.LLVMParseIRInContext(ctx.Ref(), buf.Ref())
	runtime.KeepAlive(buf)
	if err != nil {
		return nil, errs.WrapError(llvm.ErrParse, op, err)
	}
	return newModule(ctx, ref), nil
}

// ParseBitcode 从内存缓冲解析 bitcode；缓冲不被消费，仍由调用方 Close。
// 解析失败返回 ErrParse 与完整诊断。
func ParseBitcode(ctx *llvm.Context, buf *llvm.MemoryBuffer) (*Module, error) {
	const op = "ir.ParseBitcode"
	if !ctx.Alive() {
		errs.Panicf(llvm.ErrUseAfterFree, op, "context is closed")
	}
	if buf == nil {
		errs.Panicf(llvm.ErrInvalidArg, op, "nil memory buffer")
	}
	if !buf.Alive() {
		errs.Panicf(llvm.ErrUseAfterFree, op, "memory buffer is closed")
	}
	ref, err := binding.LLVMParseBitcodeInContext(ctx.Ref(), buf.Ref())
	runtime.KeepAlive(buf)
	if err != nil {
		return nil, errs.WrapError(llvm.ErrParse, op, err)
	}
	return newModule(ctx, ref), nil
}

// ParseIRString 从 IR 文本解析模块（便捷入口：内部创建并释放内存缓冲）
func ParseIRString(ctx *llvm.Context, s string) (*Module, error) {
	buf := llvm.NewMemoryBuffer([]byte(s), "<ir>")
	defer buf.Close()
	return ParseIR(ctx, buf)
}

// ParseBitcodeBytes 从 bitcode 字节解析模块（便捷入口：内部创建并释放内存缓冲）
func ParseBitcodeBytes(ctx *llvm.Context, data []byte) (*Module, error) {
	buf := llvm.NewMemoryBuffer(data, "<bitcode>")
	defer buf.Close()
	return ParseBitcode(ctx, buf)
}
