package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestFunctionPersonality(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "eh")
	defer m.Close()

	i32 := ctx.Int(32)
	pers := m.NewFunction("pers", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))

	fn.SetPersonality(pers)
	got, ok := fn.Personality()
	if !ok {
		t.Fatal("personality should be set")
	}
	if got.Name() != "pers" {
		t.Fatalf("personality = %s, want pers", got.Name())
	}
	if !strings.Contains(m.String(), "personality ptr @pers") {
		t.Fatalf("module should contain personality:\n%s", m.String())
	}

	if err := llvm.Catch(func() { fn.SetPersonality(llvm.Value[llvm.FnT]{}) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil personality should panic ErrInvalidArg, got %v", err)
	}
}

func buildEHModule(t *testing.T) (*llvm.Context, *Module, *Builder) {
	t.Helper()
	ctx := llvm.NewContext()
	m := NewModule(ctx, "eh")
	b := NewBuilder(ctx)
	return ctx, m, b
}

func TestBuilderInvoke(t *testing.T) {
	ctx, m, b := buildEHModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	g := m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	pers := m.NewFunction("pers", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn.SetPersonality(pers)
	entry := fn.NewBlock("entry")
	cont := fn.NewBlock("cont")
	lpad := fn.NewBlock("lpad")

	b.MoveToEnd(entry)
	iv := b.Invoke[llvm.IntT](g, []llvm.AnyValue{fn.ParamAs[llvm.IntT](0)}, cont, lpad, "v")

	if iv.NormalBlock() != cont {
		t.Fatalf("normal block = %s, want cont", iv.NormalBlock().Name())
	}
	if iv.UnwindBlock() != lpad {
		t.Fatalf("unwind block = %s, want lpad", iv.UnwindBlock().Name())
	}
	if iv.ArgCount() != 1 {
		t.Fatalf("arg count = %d, want 1", iv.ArgCount())
	}
	if called, ok := iv.CalledFunction(); !ok || called.Name() != "g" {
		t.Fatalf("called function = %v %v, want g", called, ok)
	}

	b.MoveToEnd(cont)
	b.Ret(iv)

	b.MoveToEnd(lpad)
	padTy := ctx.Struct([]llvm.AnyType{ctx.Ptr(0), i32}, false)
	lp := b.LandingPad(padTy, "lp")
	lp.SetCleanup(true)
	b.Resume(lp)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}

	got := m.String()
	if !strings.Contains(got, "invoke i32 @g(i32 %0)") {
		t.Fatalf("missing invoke:\n%s", got)
	}
	if !strings.Contains(got, "to label %cont unwind label %lpad") {
		t.Fatalf("missing successors:\n%s", got)
	}
}

func TestBuilderInvokePrecheck(t *testing.T) {
	ctx, m, b := buildEHModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	g := m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	entry := fn.NewBlock("entry")
	cont := fn.NewBlock("cont")
	lpad := fn.NewBlock("lpad")
	b.MoveToEnd(entry)

	// 未定位
	b2 := NewBuilder(ctx)
	defer b2.Close()
	if err := llvm.Catch(func() {
		b2.Invoke[llvm.IntT](g, []llvm.AnyValue{fn.ParamAs[llvm.IntT](0)}, cont, lpad, "")
	}); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("unpositioned invoke should panic ErrInvalidArg, got %v", err)
	}

	// 返回种类不符
	if err := llvm.Catch(func() {
		b.Invoke[llvm.FloatT](g, []llvm.AnyValue{fn.ParamAs[llvm.IntT](0)}, cont, lpad, "")
	}); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("wrong return kind should panic ErrTypeMismatch, got %v", err)
	}

	// 实参个数不符
	if err := llvm.Catch(func() {
		b.Invoke[llvm.IntT](g, nil, cont, lpad, "")
	}); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("wrong arg count should panic ErrTypeMismatch, got %v", err)
	}
}

func TestBuilderLandingPadResume(t *testing.T) {
	ctx, m, b := buildEHModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	pers := m.NewFunction("pers", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	g := m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	ti := m.NewGlobal("ti", ptr)

	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn.SetPersonality(pers)
	entry := fn.NewBlock("entry")
	cont := fn.NewBlock("cont")
	lpad := fn.NewBlock("lpad")

	b.MoveToEnd(entry)
	iv := b.Invoke[llvm.IntT](g, []llvm.AnyValue{fn.ParamAs[llvm.IntT](0)}, cont, lpad, "v")

	b.MoveToEnd(cont)
	b.Ret(iv)

	b.MoveToEnd(lpad)
	padTy := ctx.Struct([]llvm.AnyType{ptr, i32}, false)
	lp := b.LandingPad(padTy, "lp")
	lp.AddClause(ti)
	lp.AddClause(ti)
	lp.SetCleanup(true)
	if lp.ClauseCount() != 2 {
		t.Fatalf("clause count = %d, want 2", lp.ClauseCount())
	}
	if got := lp.Clause(0); got.Name() != "ti" {
		t.Fatalf("clause 0 = %s, want ti", got.Name())
	}
	if !lp.IsCleanup() {
		t.Fatal("cleanup should be set")
	}
	lp.SetCleanup(false)
	if lp.IsCleanup() {
		t.Fatal("cleanup should be cleared")
	}
	b.Resume(lp)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := m.String()
	if !strings.Contains(got, "landingpad { ptr, i32 }") {
		t.Fatalf("missing landingpad:\n%s", got)
	}
	if !strings.Contains(got, "catch ptr @ti") {
		t.Fatalf("missing catch clause:\n%s", got)
	}
	if !strings.Contains(got, "resume { ptr, i32 } %lp") {
		t.Fatalf("missing resume:\n%s", got)
	}
}

func TestBuilderLandingPadPrecheck(t *testing.T) {
	ctx, m, b := buildEHModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{ptr}, false))
	entry := fn.NewBlock("entry")
	b.MoveToEnd(entry)

	padTy := ctx.Struct([]llvm.AnyType{ptr, i32}, false)
	lp := b.LandingPad(padTy, "lp")

	// 子句必须是常量
	if err := llvm.Catch(func() { lp.AddClause(fn.Param(0)) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("non-constant clause should panic ErrInvalidArg, got %v", err)
	}
	// 子句下标越界
	if err := llvm.Catch(func() { lp.Clause(9) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("clause out of range should panic ErrInvalidArg, got %v", err)
	}
	b.Unreachable()
}
