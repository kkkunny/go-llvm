package ir

import (
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestGlobalAttrsAndAlias(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "t")
	defer m.Close()
	i32 := ctx.Int(32)

	g := m.NewGlobal("g", i32)
	g.SetSection(".mydata")
	g.SetVisibility(llvm.VisibilityHidden)
	g.SetDLLStorageClass(llvm.DLLStorageExport)
	g.SetUnnamedAddr(llvm.UnnamedAddrLocal)
	g.SetThreadLocalMode(llvm.ThreadLocalLocalExec)

	if g.Section() != ".mydata" || g.Visibility() != llvm.VisibilityHidden ||
		g.DLLStorageClass() != llvm.DLLStorageExport || g.UnnamedAddr() != llvm.UnnamedAddrLocal ||
		g.ThreadLocalMode() != llvm.ThreadLocalLocalExec {
		t.Fatalf("global attrs mismatch: %s", g)
	}

	fn := m.NewFunction("f", ctx.Fn(i32, nil, false))
	fn.SetGC("statepoint-example")
	fn.SetSection(".text.f")
	if fn.GC() != "statepoint-example" || fn.Section() != ".text.f" {
		t.Fatalf("fn attrs mismatch")
	}

	prefix := ctx.ConstString("P", false)
	fn.SetPrefixData(prefix)
	if fn.PrefixData().IsNil() {
		t.Fatalf("prefix data not set")
	}
	prologue := ctx.ConstString("Q", false)
	fn.SetPrologueData(prologue)
	if fn.PrologueData().IsNil() {
		t.Fatalf("prologue data not set")
	}

	a := m.NewAlias("alias_g", i32, g)
	if a.Aliasee().String() != g.String() {
		t.Fatalf("aliasee mismatch: %s vs %s", a.Aliasee(), g)
	}
	if _, ok := m.GetAlias("alias_g"); !ok {
		t.Fatalf("GetAlias failed")
	}
	n := 0
	for range m.AllAliases() {
		n++
	}
	if n != 1 {
		t.Fatalf("AllAliases = %d", n)
	}

	resolver := m.NewFunction("resolver", ctx.Fn(ctx.Ptr(0), nil, false))
	ifu := m.NewIFunc("ifunc_g", ctx.Fn(ctx.Ptr(0), nil, false), resolver)
	if ifu.IsNil() {
		t.Fatalf("IFunc not created")
	}
}
