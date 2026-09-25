package llvm_test

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
)

// TestPrefAlignOfGlobal 经 ir 构建全局变量，验证 DataLayout 对全局推荐对齐的查询。
// 根包内部测试不得 import ir（依赖方向 llvm ← ir），故放在外部测试包。
func TestPrefAlignOfGlobal(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := ir.NewModule(ctx, "align")
	defer m.Close()

	dl := llvm.NewDataLayout("e-p:64:64-i64:64-i32:32")
	defer dl.Close()

	g := m.NewGlobal("g", ctx.Int(32))
	if got := dl.PrefAlignOfGlobal(g); got != 4 {
		t.Fatalf("PrefAlignOfGlobal(i32) = %d, want 4", got)
	}
	g64 := m.NewGlobal("g64", ctx.Int(64))
	if got := dl.PrefAlignOfGlobal(g64); got != 8 {
		t.Fatalf("PrefAlignOfGlobal(i64) = %d, want 8", got)
	}
}

// TestValueSetName 经 ir 构建全局变量，验证 Value.SetName/Name 的名称往返。
func TestValueSetName(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := ir.NewModule(ctx, "name")
	defer m.Close()

	g := m.NewGlobal("before", ctx.Int(32))
	g.SetName("after")
	if got := g.Name(); got != "after" {
		t.Fatalf("SetName 后 Name() = %q, want after", got)
	}
	if got := m.String(); !strings.Contains(got, "@after") {
		t.Fatalf("模块 IR 中缺少重命名后的全局变量:\n%s", got)
	}
}
