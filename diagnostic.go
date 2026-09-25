package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// DiagnosticSeverity 诊断严重级别
type DiagnosticSeverity binding.LLVMDiagnosticSeverity

// DiagnosticSeverity 取值对应 LLVM 诊断严重级别（binding.LLVMDiagnosticSeverity）。
// [Context.SetDiagnosticHandler] 回调的 severity 参数即为其中一种。
const (
	DiagnosticError   = DiagnosticSeverity(binding.LLVMDSError)   // 错误：报告编译错误
	DiagnosticWarning = DiagnosticSeverity(binding.LLVMDSWarning) // 警告：报告可疑但不中断编译的问题
	DiagnosticRemark  = DiagnosticSeverity(binding.LLVMDSRemark)  // 备注：报告优化等补充信息（如 -Rpass 输出）
	DiagnosticNote    = DiagnosticSeverity(binding.LLVMDSNote)    // 附注：附加在前一条诊断之后的补充信息
)

// SetDiagnosticHandler 安装上下文诊断回调（覆盖 LLVM 默认的 stderr 输出）。
// 回调可能由任意 goroutine 在 LLVM 内部触发，不得 panic；不调用时保持默认行为。
func (ctx *Context) SetDiagnosticHandler(handler func(severity DiagnosticSeverity, msg string)) {
	const op = "llvm.Context.SetDiagnosticHandler"
	ctx.CheckAlive(op)
	if handler == nil {
		errPanic(ErrInvalidArg, op, "nil diagnostic handler")
	}
	binding.LLVMContextSetDiagnosticHandlerGo(ctx.ref, ctx.diagID, func(severity binding.LLVMDiagnosticSeverity, msg string) {
		handler(DiagnosticSeverity(severity), msg)
	})
}

// ClearDiagnosticHandler 移除诊断回调，恢复 LLVM 默认行为
func (ctx *Context) ClearDiagnosticHandler() {
	ctx.CheckAlive("llvm.Context.ClearDiagnosticHandler")
	binding.LLVMContextClearDiagnosticHandlerGo(ctx.ref, ctx.diagID)
}
