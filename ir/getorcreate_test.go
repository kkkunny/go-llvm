package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestGetOrCreateFunction(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "goc-fn")
	defer m.Close()

	sig := ctx.Fn(ctx.Int(32), nil, false)
	fn, created := m.GetOrCreateFunction("f", sig)
	if !created {
		t.Fatal("first GetOrCreateFunction should create")
	}
	again, created := m.GetOrCreateFunction("f", sig)
	if created {
		t.Fatal("second GetOrCreateFunction should reuse")
	}
	if !again.Ref().Equal(fn.Ref()) {
		t.Fatal("reused function handle should be the same LLVM value")
	}

	// 模块里只有一个 @f：没有静默改名的 @f.1
	if got := m.String(); strings.Contains(got, "@f.1") {
		t.Fatalf("GetOrCreateFunction must not create @f.1:\n%s", got)
	}

	// 同名不同签名是程序员错误：panic 而非静默复用/改名
	other := ctx.Fn(ctx.Int(64), nil, false)
	if err := llvm.Catch(func() { m.GetOrCreateFunction("f", other) }); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("signature mismatch should panic ErrTypeMismatch, got %v", err)
	}
}

func TestGetOrCreateGlobal(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "goc-global")
	defer m.Close()

	i32 := ctx.Int(32)
	g, created := m.GetOrCreateGlobal("g", i32)
	if !created {
		t.Fatal("first GetOrCreateGlobal should create")
	}
	again, created := m.GetOrCreateGlobal("g", i32)
	if created {
		t.Fatal("second GetOrCreateGlobal should reuse")
	}
	if !again.Ref().Equal(g.Ref()) {
		t.Fatal("reused global handle should be the same LLVM value")
	}
	if !g.ValueType().Equal(i32) {
		t.Fatalf("global value type = %s, want i32", g.ValueType())
	}
	if got := m.String(); strings.Contains(got, "@g.1") {
		t.Fatalf("GetOrCreateGlobal must not create @g.1:\n%s", got)
	}

	if err := llvm.Catch(func() { m.GetOrCreateGlobal("g", ctx.Int(64)) }); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("type mismatch should panic ErrTypeMismatch, got %v", err)
	}
}
