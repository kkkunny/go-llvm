// Package pass runs optimization pipelines through LLVM's PassBuilder; it does not define
// passes itself.
//
// [RunPasses] and [RunPassesOnFunction] take pipelines in the opt -passes syntax, such as
// "default<O2>" or "function(instcombine)"; [AutoOpt] runs the default pipeline for an
// optimization [Level]. [WithTargetMachine] lets a pipeline see the target machine, which
// enables target-dependent optimizations and cost models. The dependency direction is one-way:
// this package depends on [github.com/kkkunny/go-llvm/ir] and
// [github.com/kkkunny/go-llvm/target], never the reverse.
package pass

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/ir"
	"github.com/kkkunny/go-llvm/target"
)

// Level 优化级别
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

// Option 管线运行选项，由本包的构造器创建（[VerifyEach]、[DebugLogging]、[WithTargetMachine] 等）。
// 底层类型故意不导出，调用方只能经构造器组合选项。
type Option func(c *runConfig)

// runConfig 运行期配置：PassBuilder 选项 + 目标机器
type runConfig struct {
	opts binding.LLVMPassBuilderOptionsRef
	tm   *target.TargetMachine
}

// newRunConfig 创建运行期配置（含默认 PassBuilder 选项），用毕须释放 opts
func newRunConfig() *runConfig {
	return &runConfig{opts: binding.LLVMCreatePassBuilderOptions()}
}

// WithTargetMachine 让管线使用目标机器：default<O2> 等据此做目标相关优化与成本模型
// （等效于 opt 的 -mcpu 行为）。传 nil 表示不使用，与不传该选项一致。
func WithTargetMachine(tm *target.TargetMachine) Option {
	return func(c *runConfig) { c.tm = tm }
}

// VerifyEach 每趟 pass 后运行 verifier
func VerifyEach(v bool) Option {
	return func(c *runConfig) { binding.LLVMPassBuilderOptionsSetVerifyEach(c.opts, v) }
}

// DebugLogging 输出 pass 调试日志
func DebugLogging(v bool) Option {
	return func(c *runConfig) { binding.LLVMPassBuilderOptionsSetDebugLogging(c.opts, v) }
}

// LoopInterleaving 启用循环交错
func LoopInterleaving(v bool) Option {
	return func(c *runConfig) { binding.LLVMPassBuilderOptionsSetLoopInterleaving(c.opts, v) }
}

// LoopVectorization 启用循环向量化
func LoopVectorization(v bool) Option {
	return func(c *runConfig) { binding.LLVMPassBuilderOptionsSetLoopVectorization(c.opts, v) }
}

// SLPVectorization 启用 SLP 向量化
func SLPVectorization(v bool) Option {
	return func(c *runConfig) { binding.LLVMPassBuilderOptionsSetSLPVectorization(c.opts, v) }
}

// LoopUnrolling 启用循环展开
func LoopUnrolling(v bool) Option {
	return func(c *runConfig) { binding.LLVMPassBuilderOptionsSetLoopUnrolling(c.opts, v) }
}

// ForgetAllSCEVInLoopUnroll 循环展开后丢弃全部 SCEV
func ForgetAllSCEVInLoopUnroll(v bool) Option {
	return func(c *runConfig) {
		binding.LLVMPassBuilderOptionsSetForgetAllSCEVInLoopUnroll(c.opts, v)
	}
}

// InlinerThreshold 内联阈值
func InlinerThreshold(n int32) Option {
	return func(c *runConfig) { binding.LLVMPassBuilderOptionsSetInlinerThreshold(c.opts, n) }
}

// applyOptions 逐项应用选项并返回目标机器句柄（未指定时为空的 LLVMTargetMachineRef）
func (c *runConfig) applyOptions(op string, opts []Option) binding.LLVMTargetMachineRef {
	for _, opt := range opts {
		opt(c)
	}
	if c.tm == nil {
		return binding.LLVMTargetMachineRef{}
	}
	c.tm.Check(op)
	return c.tm.Ref()
}

// RunPasses 在模块上运行 passes 管线（opt -passes 语法，如 "default<O2>"、"function(instcombine)"）
func RunPasses(m *ir.Module, pipeline string, opts ...Option) error {
	const op = "pass.RunPasses"
	m.Check(op)
	cfg := newRunConfig()
	defer binding.LLVMDisposePassBuilderOptions(cfg.opts)
	tm := cfg.applyOptions(op, opts)
	if err := binding.LLVMRunPasses(m.Ref(), pipeline, tm, cfg.opts); err != nil {
		return llvm.WrapError(llvm.ErrPass, op, err)
	}
	return nil
}

// RunPassesOnFunction 在单个函数上运行 passes 管线
func RunPassesOnFunction(f ir.Function, pipeline string, opts ...Option) error {
	const op = "pass.RunPassesOnFunction"
	f.Check(op)
	cfg := newRunConfig()
	defer binding.LLVMDisposePassBuilderOptions(cfg.opts)
	tm := cfg.applyOptions(op, opts)
	if err := binding.LLVMRunPassesOnFunction(f.Ref(), pipeline, tm, cfg.opts); err != nil {
		return llvm.WrapError(llvm.ErrPass, op, err)
	}
	return nil
}

// AutoOpt 运行 default<level> 默认管线
func AutoOpt(m *ir.Module, level Level, opts ...Option) error {
	return RunPasses(m, "default<"+string(level)+">", opts...)
}
