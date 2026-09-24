package llvm

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm/internal/binding"
)

func TestDiagnosticHandler(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	type diag struct {
		sev DiagnosticSeverity
		msg string
	}
	var got []diag
	ctx.SetDiagnosticHandler(func(sev DiagnosticSeverity, msg string) {
		got = append(got, diag{sev, msg})
	})
	binding.LLVMGoEmitError(ctx.Ref(), "boom")
	if len(got) != 1 || got[0].sev != DiagnosticError || !strings.Contains(got[0].msg, "boom") {
		t.Fatalf("diagnostics = %v", got)
	}

	// 清除后不再有回调（默认 handler 对 error 会终止进程，故此处不再触发诊断）
	ctx.ClearDiagnosticHandler()
	if binding.LLVMContextHasDiagnosticHandler(ctx.Ref()) {
		t.Fatal("diagnostic handler should be cleared")
	}

	// 重新安装后仍可收到诊断
	ctx.SetDiagnosticHandler(func(sev DiagnosticSeverity, msg string) {
		got = append(got, diag{sev, msg})
	})
	binding.LLVMGoEmitError(ctx.Ref(), "again")
	if len(got) != 2 || got[1].msg != "again" {
		t.Fatalf("diagnostics after reinstall = %v", got)
	}

	if err := Catch(func() { ctx.SetDiagnosticHandler(nil) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("nil handler should panic ErrInvalidArg, got %v", err)
	}
}
