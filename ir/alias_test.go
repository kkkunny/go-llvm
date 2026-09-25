package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

// TestAliasSetAliasee 覆盖别名的改指、查找失败与提前退出遍历。
func TestAliasSetAliasee(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "alias")
	defer m.Close()

	i32 := ctx.Int(32)
	g1 := m.NewGlobal("g1", i32)
	g1.SetInitializer(ctx.ConstInt(i32, 1))
	g2 := m.NewGlobal("g2", i32)
	g2.SetInitializer(ctx.ConstInt(i32, 2))

	a := m.NewAlias("a", i32, g1)
	if got := a.Aliasee(); got.String() != g1.String() {
		t.Fatalf("aliasee = %s, want %s", got, g1)
	}

	// SetAliasee 改指后，读回与打印文本都必须跟随
	a.SetAliasee(g2)
	if got := a.Aliasee().String(); got != g2.String() {
		t.Fatalf("aliasee after set = %s, want %s", got, g2)
	}
	if out := m.String(); !strings.Contains(out, "@a = alias i32, ptr @g2") {
		t.Fatalf("IR missing alias target:\n%s", out)
	}

	// 未登记的别名
	if _, ok := m.GetAlias("no-such-alias"); ok {
		t.Fatal("GetAlias should return false for missing name")
	}
	// AllAliases 提前 break 后不得继续遍历
	n := 0
	for range m.AllAliases() {
		n++
		break
	}
	if n != 1 {
		t.Fatalf("early break yielded %d aliases", n)
	}

	// 跨 Context 的 aliasee（崩溃类地板：始终校验）
	ctx2 := llvm.NewContext()
	defer ctx2.Close()
	foreign := ctx2.ConstInt(ctx2.Int(32), 1)
	if err := llvm.Catch(func() { a.SetAliasee(foreign) }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("foreign aliasee should panic ErrCrossContext, got %v", err)
	}
}
