package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

// TestGlobalLinkage 覆盖全局变量链接类型的读写与 IR 打印。
func TestGlobalLinkage(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "global")
	defer m.Close()

	i32 := ctx.Int(32)
	g := m.NewGlobal("g", i32)
	if got := g.Linkage(); got != llvm.LinkageExternal {
		t.Fatalf("default linkage = %v, want external", got)
	}

	g.SetLinkage(llvm.LinkageInternal)
	if got := g.Linkage(); got != llvm.LinkageInternal {
		t.Fatalf("linkage after set = %v, want internal", got)
	}

	g2 := m.NewGlobal("g2", i32)
	g2.SetLinkage(llvm.LinkageWeakODR)
	if got := g2.Linkage(); got != llvm.LinkageWeakODR {
		t.Fatalf("g2 linkage = %v, want weak_odr", got)
	}

	out := m.String()
	for _, want := range []string{"@g = internal global i32", "@g2 = weak_odr global i32"} {
		if !strings.Contains(out, want) {
			t.Fatalf("IR missing %q:\n%s", want, out)
		}
	}
}
