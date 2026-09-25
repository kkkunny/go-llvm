// Package pass runs optimization pipelines through LLVM's PassBuilder; it does not define
// passes itself.
//
// [RunPasses] and [RunPassesOnFunction] take pipelines in the opt -passes syntax, such as
// "default<O2>" or "function(instcombine)"; [AutoOpt] runs the default pipeline for an
// optimization [Level]. The dependency direction is one-way: this package depends on
// [github.com/kkkunny/go-llvm/ir], never the reverse.
package pass

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/ir"
)

// Level 默认管线优化级别（对应 default<O*>）
type Level string

// Level 取值对应 PassBuilder 的优化级别（default<O*> 管线），决定默认管线包含的
// 优化 pass 及激进程度；Oz/Os 以代码体积为优先。
const (
	O0 Level = "O0" // 不优化
	O1 Level = "O1" // 轻度优化
	O2 Level = "O2" // 标准优化：常用的平衡级别
	O3 Level = "O3" // 激进优化
	Oz Level = "Oz" // 体积优先：以最小代码体积为目标，可能牺牲运行速度
	Os Level = "Os" // 体积优化：在兼顾性能的前提下减小代码体积
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
