package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestModuleLink(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	dst := NewModule(ctx, "dst")
	defer dst.Close()
	src := NewModule(ctx, "src") // 被 link 消费，不 defer Close

	i32 := ctx.Int(32)
	g := src.NewFunction("g", ctx.Fn(i32, nil, false))
	entry := g.NewBlock("entry")
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)
	b.Ret(ctx.ConstInt(i32, 42).Value)

	if err := dst.Link(src); err != nil {
		t.Fatalf("link: %v", err)
	}
	if fn, ok := dst.GetFunction("g"); !ok || fn.OnlyDecl() {
		t.Fatalf("linked function missing or declaration: %v %v", ok, fn.OnlyDecl())
	}

	// 源模块句柄立即失效
	if err := llvm.Catch(func() { _ = src.String() }); err == nil || err.Reason != llvm.ErrUseAfterFree {
		t.Fatalf("consumed src should panic ErrUseAfterFree, got %v", err)
	}
	if err := src.Close(); err == nil {
		t.Fatal("consumed src Close should return ErrClosed")
	} else if e, ok := err.(*llvm.Error); !ok || e.Reason != llvm.ErrClosed {
		t.Fatalf("consumed src Close should return ErrClosed, got %v", err)
	}

	if err := dst.Verify(); err != nil {
		t.Fatalf("verify dst: %v", err)
	}
	if got := dst.String(); !strings.Contains(got, "define i32 @g()") || !strings.Contains(got, "ret i32 42") {
		t.Fatalf("linked module:\n%s", got)
	}

	// 自链接
	if err := llvm.Catch(func() { _ = dst.Link(dst) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("self link should panic ErrInvalidArg, got %v", err)
	}
	// 跨 Context 链接
	ctx2 := llvm.NewContext()
	defer ctx2.Close()
	foreign := NewModule(ctx2, "foreign")
	if err := llvm.Catch(func() { _ = dst.Link(foreign) }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("cross-context link should panic ErrCrossContext, got %v", err)
	}
}

func TestCtorsDtors(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "ctors")
	defer m.Close()

	voidFn := ctx.Fn(ctx.Void(), nil, false)
	c1 := m.NewFunction("ctor1", voidFn)
	c2 := m.NewFunction("ctor2", voidFn)
	b := NewBuilder(ctx)
	defer b.Close()
	for _, f := range []Function{c1, c2} {
		e := f.NewBlock("entry")
		b.MoveToEnd(e)
		b.RetVoid()
	}

	m.AppendCtor(c1, 65535)
	m.AppendCtor(c2, 1)
	m.AppendDtor(c2, 10)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := m.String()
	for _, want := range []string{
		"@llvm.global_ctors = appending global [2 x { i32, ptr, ptr }]",
		"@llvm.global_dtors = appending global [1 x { i32, ptr, ptr }]",
		"{ i32 65535, ptr @ctor1, ptr null }",
		"{ i32 1, ptr @ctor2, ptr null }",
		"{ i32 10, ptr @ctor2, ptr null }",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q:\n%s", want, got)
		}
	}
}
