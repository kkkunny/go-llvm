// Package pass 提供优化管线执行入口（PassBuilder），不定义 pass 本身。
//
// 依赖方向：llvm/ir ← llvm/pass（pass 不反向被 ir 依赖）。
package pass

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/ir"
)

// Level 默认管线优化级别（对应 default<O*>）
type Level string

const (
	O0 Level = "O0"
	O1 Level = "O1"
	O2 Level = "O2"
	O3 Level = "O3"
	Oz Level = "Oz"
	Os Level = "Os"
)

// Option PassBuilder 选项
type Option func(o binding.LLVMPassBuilderOptionsRef)

// VerifyEach 每趟 pass 后运行 verifier
func VerifyEach(v bool) Option {
	return func(o binding.LLVMPassBuilderOptionsRef) { binding.LLVMPassBuilderOptionsSetVerifyEach(o, v) }
}

// DebugLogging 输出 pass 调试日志
func DebugLogging(v bool) Option {
	return func(o binding.LLVMPassBuilderOptionsRef) { binding.LLVMPassBuilderOptionsSetDebugLogging(o, v) }
}

// LoopInterleaving 启用循环交错
func LoopInterleaving(v bool) Option {
	return func(o binding.LLVMPassBuilderOptionsRef) { binding.LLVMPassBuilderOptionsSetLoopInterleaving(o, v) }
}

// LoopVectorization 启用循环向量化
func LoopVectorization(v bool) Option {
	return func(o binding.LLVMPassBuilderOptionsRef) { binding.LLVMPassBuilderOptionsSetLoopVectorization(o, v) }
}

// SLPVectorization 启用 SLP 向量化
func SLPVectorization(v bool) Option {
	return func(o binding.LLVMPassBuilderOptionsRef) { binding.LLVMPassBuilderOptionsSetSLPVectorization(o, v) }
}

// LoopUnrolling 启用循环展开
func LoopUnrolling(v bool) Option {
	return func(o binding.LLVMPassBuilderOptionsRef) { binding.LLVMPassBuilderOptionsSetLoopUnrolling(o, v) }
}

// ForgetAllSCEVInLoopUnroll 循环展开后丢弃全部 SCEV
func ForgetAllSCEVInLoopUnroll(v bool) Option {
	return func(o binding.LLVMPassBuilderOptionsRef) {
		binding.LLVMPassBuilderOptionsSetForgetAllSCEVInLoopUnroll(o, v)
	}
}

// InlinerThreshold 内联阈值
func InlinerThreshold(n int32) Option {
	return func(o binding.LLVMPassBuilderOptionsRef) { binding.LLVMPassBuilderOptionsSetInlinerThreshold(o, n) }
}

// RunPasses 在模块上运行 passes 管线（opt -passes 语法，如 "default<O2>"、"function(instcombine)"）
func RunPasses(m *ir.Module, pipeline string, opts ...Option) error {
	const op = "pass.RunPasses"
	m.Check(op)
	o := binding.LLVMCreatePassBuilderOptions()
	defer binding.LLVMDisposePassBuilderOptions(o)
	for _, opt := range opts {
		opt(o)
	}
	if err := binding.LLVMRunPasses(m.Ref(), pipeline, binding.LLVMTargetMachineRef{}, o); err != nil {
		return llvm.WrapError(llvm.ErrPass, op, err)
	}
	return nil
}

// RunPassesOnFunction 在单个函数上运行 passes 管线
func RunPassesOnFunction(f ir.Function, pipeline string, opts ...Option) error {
	const op = "pass.RunPassesOnFunction"
	f.Check(op)
	o := binding.LLVMCreatePassBuilderOptions()
	defer binding.LLVMDisposePassBuilderOptions(o)
	for _, opt := range opts {
		opt(o)
	}
	if err := binding.LLVMRunPassesOnFunction(f.Ref(), pipeline, binding.LLVMTargetMachineRef{}, o); err != nil {
		return llvm.WrapError(llvm.ErrPass, op, err)
	}
	return nil
}

// AutoOpt 运行 default<level> 默认管线
func AutoOpt(m *ir.Module, level Level, opts ...Option) error {
	return RunPasses(m, "default<"+string(level)+">", opts...)
}
