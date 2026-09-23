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

func TestBuilderFunclets(t *testing.T) {
	ctx, m, b := buildEHModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	pers := m.NewFunction("pers", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	h := m.NewFunction("h", ctx.Fn(ctx.Void(), nil, false))
	ti := m.NewGlobal("ti", ptr)

	// ---- catchswitch -> catchpad -> catchret ----
	fn := m.NewFunction("funclets", ctx.Fn(ctx.Void(), nil, false))
	fn.SetPersonality(pers)
	entry := fn.NewBlock("entry")
	cont := fn.NewBlock("cont")
	dispatch := fn.NewBlock("dispatch")
	hnd := fn.NewBlock("hnd")
	done := fn.NewBlock("done")

	b.MoveToEnd(entry)
	b.Invoke[llvm.VoidT](h, nil, cont, dispatch, "") // EH pad 需 unwind 边进入

	b.MoveToEnd(cont)
	b.RetVoid()

	b.MoveToEnd(dispatch)                   // EH pad 不得在 entry 块
	cs := b.CatchSwitch(nil, Block{}, "cs") // parent=nil 即 within none；unwindTo 零块即 unwind to caller
	cs.AddHandler(hnd)
	if cs.HandlerCount() != 1 {
		t.Fatalf("handler count = %d, want 1", cs.HandlerCount())
	}
	if cs.HandlerAt(0) != hnd {
		t.Fatalf("handler 0 = %s, want hnd", cs.HandlerAt(0).Name())
	}

	b.MoveToEnd(hnd)
	cp := b.CatchPad(cs, []llvm.AnyValue{ti}, "cp")
	if got := cp.ParentCatchSwitch(); got.Name() != "cs" {
		t.Fatalf("parent catchswitch = %s, want cs", got.Name())
	}
	if cp.ArgCount() != 1 {
		t.Fatalf("catchpad arg count = %d, want 1", cp.ArgCount())
	}
	if got := cp.Arg(0); got.Name() != "ti" {
		t.Fatalf("catchpad arg 0 = %s, want ti", got.Name())
	}
	b.CatchRet(cp, done)

	b.MoveToEnd(done)
	b.RetVoid()

	// ---- cleanuppad -> cleanupret ----
	kn := m.NewFunction("cleanup", ctx.Fn(ctx.Void(), nil, false))
	kn.SetPersonality(pers)
	e2 := kn.NewBlock("entry")
	c2 := kn.NewBlock("cont")
	clean := kn.NewBlock("clean")

	b.MoveToEnd(e2)
	b.Invoke[llvm.VoidT](h, nil, c2, clean, "")
	b.MoveToEnd(c2)
	b.RetVoid()

	b.MoveToEnd(clean)
	clp := b.CleanupPad(nil, nil, "clp")
	if clp.ArgCount() != 0 {
		t.Fatalf("cleanuppad arg count = %d, want 0", clp.ArgCount())
	}
	b.CleanupRet(clp, Block{}) // unwind to caller

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := m.String()
	if !strings.Contains(got, "catchswitch within none [label %hnd] unwind to caller") {
		t.Fatalf("missing catchswitch:\n%s", got)
	}
	if !strings.Contains(got, "catchpad within %cs [ptr @ti]") {
		t.Fatalf("missing catchpad:\n%s", got)
	}
	if !strings.Contains(got, "catchret from %cp to label %done") {
		t.Fatalf("missing catchret:\n%s", got)
	}
	if !strings.Contains(got, "cleanuppad within none []") {
		t.Fatalf("missing cleanuppad:\n%s", got)
	}
	if !strings.Contains(got, "cleanupret from %clp unwind to caller") {
		t.Fatalf("missing cleanupret:\n%s", got)
	}
}

func TestBuilderFuncletPrecheck(t *testing.T) {
	ctx, m, b := buildEHModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))
	entry := fn.NewBlock("entry")
	b.MoveToEnd(entry)

	cs := b.CatchSwitch(nil, Block{}, "cs")

	// handler 下标越界
	if err := llvm.Catch(func() { cs.HandlerAt(3) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("handler out of range should panic ErrInvalidArg, got %v", err)
	}
	// funclet pad 实参下标越界
	pad := b.CatchPad(cs, nil, "cp")
	if err := llvm.Catch(func() { pad.Arg(0) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("pad arg out of range should panic ErrInvalidArg, got %v", err)
	}
	b.Unreachable()
}
