package llvm

import (
	"log"
	"os"
)

// logger go-llvm 诊断日志（默认 stderr）；SetLogger(nil) 可静默。
// 调试层告警（资源未显式释放等）与默认 LLVM 诊断回调经此输出。
var logger = log.New(os.Stderr, "[go-llvm] ", 0)

// SetLogger 替换诊断日志输出（nil 静默）
func SetLogger(l *log.Logger) { logger = l }

// logf 内部日志输出
func logf(format string, args ...any) {
	if logger != nil {
		logger.Printf(format, args...)
	}
}

// defaultDiagnosticHandler 默认诊断回调：把 LLVM 默认 handler 的进程退出/abort
// 变为可观测的日志（回调可能由任意 goroutine 触发，禁止 panic）。
func defaultDiagnosticHandler(severity DiagnosticSeverity, msg string) {
	switch severity {
	case DiagnosticError:
		logf("llvm error: %s", msg)
	case DiagnosticWarning:
		logf("llvm warning: %s", msg)
	default:
		logf("llvm diagnostic: %s", msg)
	}
}
