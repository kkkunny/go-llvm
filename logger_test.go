package llvm

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

// TestLoggerAndDefaultDiagnosticHandler 验证日志替换与默认诊断回调的分级输出。
// 默认回调把 LLVM 默认 handler 的进程退出变为日志；这里直接调用回调本身，
// 不触发 LLVM 内部诊断（避免进程退出）。
func TestLoggerAndDefaultDiagnosticHandler(t *testing.T) {
	var buf bytes.Buffer
	old := logger
	SetLogger(log.New(&buf, "", 0))
	defer SetLogger(old)

	defaultDiagnosticHandler(DiagnosticError, "boom")
	if got := buf.String(); !strings.Contains(got, "llvm error: boom") {
		t.Fatalf("错误级别日志缺失: %q", got)
	}
	defaultDiagnosticHandler(DiagnosticWarning, "careful")
	if got := buf.String(); !strings.Contains(got, "llvm warning: careful") {
		t.Fatalf("警告级别日志缺失: %q", got)
	}

	// remark/note 属提示信息，默认回调不输出
	n := buf.Len()
	defaultDiagnosticHandler(DiagnosticRemark, "hint")
	defaultDiagnosticHandler(DiagnosticNote, "note")
	if buf.Len() != n {
		t.Fatalf("remark/note 不应输出日志: %q", buf.String()[n:])
	}

	// SetLogger(nil) 静默
	SetLogger(nil)
	defaultDiagnosticHandler(DiagnosticError, "silent")
	if buf.Len() != n {
		t.Fatalf("SetLogger(nil) 后不应再输出: %q", buf.String()[n:])
	}
}
