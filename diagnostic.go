package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// DiagnosticSeverity 诊断严重级别
type DiagnosticSeverity binding.LLVMDiagnosticSeverity

const (
	DiagnosticError   = DiagnosticSeverity(binding.LLVMDSError)
	DiagnosticWarning = DiagnosticSeverity(binding.LLVMDSWarning)
	DiagnosticRemark  = DiagnosticSeverity(binding.LLVMDSRemark)
	DiagnosticNote    = DiagnosticSeverity(binding.LLVMDSNote)
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
