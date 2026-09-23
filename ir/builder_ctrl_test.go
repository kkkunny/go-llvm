package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func buildAddModule(t *testing.T) (*llvm.Context, *Module, *Builder) {
	t.Helper()
	ctx := llvm.NewContext()
	m := NewModule(ctx, "ctrl")
	i32 := ctx.Int(32)
	fn := m.NewFunction("add", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	blk := fn.NewBlock("entry")
	b := NewBuilder(ctx)
	b.MoveToEnd(blk)
	return ctx, m, b
}

func TestBuilderPrecheck(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	b := NewBuilder(ctx)
	defer b.Close()

	if err := llvm.Catch(func() { b.RetVoid() }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("unpositioned builder should panic ErrInvalidArg, got %v", err)
	}

	m := NewModule(ctx, "pre")
	defer m.Close()
	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, nil, false))
	b.MoveToEnd(fn.NewBlock("entry"))

	ctx2 := llvm.NewContext()
	defer ctx2.Close()
	other := ctx2.ConstInt(ctx2.Int(32), 1)
	if err := llvm.Catch(func() { b.Ret(other) }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("cross-context operand should panic ErrCrossContext, got %v", err)
	}

	life := llvm.NewLifetime()
	dead := llvm.NewValue[llvm.IntT](ctx, life, ctx.ConstInt(i32, 1).Ref())
	life.Kill()
	if err := llvm.Catch(func() { b.Ret(dead) }); err == nil || err.Reason != llvm.ErrUseAfterFree {
		t.Fatalf("dead operand should panic ErrUseAfterFree, got %v", err)
	}
}

func TestBuilderMovePrecheck(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "move-pre")
	defer m.Close()
	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, nil, false))

	b := NewBuilder(ctx)
	defer b.Close()

	if err := llvm.Catch(func() { b.MoveToEnd(Block{}) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil block should panic ErrInvalidArg, got %v", err)
	}

	ctx2 := llvm.NewContext()
	defer ctx2.Close()
	m2 := NewModule(ctx2, "other")
	fn2 := m2.NewFunction("g", ctx2.Fn(ctx2.Int(32), nil, false))
	blk2 := fn2.NewBlock("entry")
	if err := llvm.Catch(func() { b.MoveToEnd(blk2) }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("cross-context block should panic ErrCrossContext, got %v", err)
	}

	blk := fn.NewBlock("entry")
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	if err := llvm.Catch(func() { b.MoveToEnd(blk) }); err == nil || err.Reason != llvm.ErrClosed {
		t.Fatalf("closed builder MoveToEnd should panic ErrClosed, got %v", err)
	}
}

func TestBuilderMoveBeforeLifetime(t *testing.T) {
	ctx, m, b := buildAddModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	fn, _ := m.GetFunction("add")
	entry := fn.NewBlock("entry")
	b.MoveToEnd(entry)
	ret := b.Ret(fn.ParamAs[llvm.IntT](0))

	b.MoveBefore(ret)
	b.RetVoid()
	got := m.String()
	if !strings.Contains(got, "ret void") {
		t.Fatalf("MoveBefore should position before ret:\n%s", got)
	}
	if !strings.Contains(got, "ret i32 %0") {
		t.Fatalf("original ret should remain:\n%s", got)
	}
}

func TestBuilderClosePrecheck(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	b := NewBuilder(ctx)
	_ = b.Close()
	if err := llvm.Catch(func() { b.RetVoid() }); err == nil || err.Reason != llvm.ErrClosed {
		t.Fatalf("closed builder should panic ErrClosed, got %v", err)
	}
	if err := b.Close(); err == nil || err.(*llvm.Error).Reason != llvm.ErrClosed {
		t.Fatalf("double close should return ErrClosed, got %v", err)
	}
}

func TestBuilderRet(t *testing.T) {
	ctx, m, b := buildAddModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	fn, _ := m.GetFunction("add")
	b.Ret(fn.ParamAs[llvm.IntT](0))

	got := m.String()
	for _, want := range []string{
		"define i32 @add(i32 %0, i32 %1)",
		"ret i32 %0",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("module output missing %q:\n%s", want, got)
		}
	}
}

func TestBuilderRetVoid(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "void")
	defer m.Close()
	fn := m.NewFunction("nothing", ctx.Fn(ctx.Void(), nil, false))
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(fn.NewBlock("entry"))
	b.RetVoid()

	if got := m.String(); !strings.Contains(got, "define void @nothing()") || !strings.Contains(got, "ret void") {
		t.Fatalf("module output:\n%s", got)
	}
}

func TestBuilderBranches(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "br")
	defer m.Close()

	i32 := ctx.Int(32)
	boolTy := ctx.Bool()
	fn := m.NewFunction("branchy", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	entry := fn.NewBlock("entry")
	thenBlk := fn.NewBlock("then")
	elseBlk := fn.NewBlock("else")

	b := NewBuilder(ctx)
	defer b.Close()

	b.MoveToEnd(entry)
	b.CondBr(ctx.ConstBool(true), thenBlk, elseBlk)

	b.MoveToEnd(thenBlk)
	b.Ret(ctx.ConstSInt(i32, -1))

	b.MoveToEnd(elseBlk)
	b.Br(thenBlk)

	got := m.String()
	for _, want := range []string{
		"br i1 true, label %then, label %else",
		"br label %then",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("module output missing %q:\n%s", want, got)
		}
	}

	// CondBr 条件必须 i1
	blk2 := fn.NewBlock("bad")
	b.MoveToEnd(blk2)
	err := llvm.Catch(func() {
		b.CondBr(ctx.ConstInt(i32, 1), thenBlk, elseBlk)
	})
	if err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("non-i1 condition should panic ErrTypeMismatch, got %v", err)
	}
	_ = boolTy
}

func TestBuilderSwitch(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "switch")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("sw", ctx.Fn(ctx.Void(), []llvm.AnyType{i32}, false))
	entry := fn.NewBlock("entry")
	defBlk := fn.NewBlock("default")
	case1 := fn.NewBlock("case1")
	case2 := fn.NewBlock("case2")

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	sw := b.Switch(fn.ParamAs[llvm.IntT](0), defBlk)
	sw.AddCase(ctx.ConstInt(i32, 1), case1)
	sw.AddCase(ctx.ConstInt(i32, 2), case2)

	if sw.Count() != 2 {
		t.Fatalf("switch case count = %d", sw.Count())
	}
	if sw.DefaultBlock().Name() != "default" {
		t.Fatalf("default = %q", sw.DefaultBlock().Name())
	}
	if sw.CaseBlock(0).Name() != "case1" || sw.CaseBlock(1).Name() != "case2" {
		t.Fatalf("case blocks = %q/%q", sw.CaseBlock(0).Name(), sw.CaseBlock(1).Name())
	}
	if got := sw.CaseValue(0).String(); got != "i32 1" {
		t.Fatalf("case value = %q", got)
	}

	b.MoveToEnd(defBlk)
	b.RetVoid()
	b.MoveToEnd(case1)
	b.Br(case2)
	b.MoveToEnd(case2)
	b.RetVoid()

	got := m.String()
	for _, want := range []string{
		"switch i32 %0, label %default [",
		"i32 1, label %case1",
		"i32 2, label %case2",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("module output missing %q:\n%s", want, got)
		}
	}

	// case 类型不符
	b.MoveToEnd(case1)
	err := llvm.Catch(func() { sw.AddCase(ctx.ConstInt(ctx.Int(64), 1), case2) })
	if err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("case type mismatch should panic, got %v", err)
	}
}

func TestBuilderUnreachableAndMoveBefore(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "unreach")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, nil, false))
	entry := fn.NewBlock("entry")

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)
	ret := b.Ret(ctx.ConstInt(i32, 0))
	b.Unreachable()

	if got := m.String(); !strings.Contains(got, "unreachable") {
		t.Fatalf("module output:\n%s", got)
	}

	// MoveBefore 在 ret 前插入指令
	b.MoveBefore(ret)
	b.Ret(ctx.ConstInt(i32, 1))

	if cur, ok := b.CurrentBlock(); !ok || cur.Name() != "entry" {
		t.Fatalf("CurrentBlock() = %v, %v", cur.Name(), ok)
	}
	insts := entry.Insts()
	if len(insts) != 3 {
		t.Fatalf("entry should have 3 instructions, got %d", len(insts))
	}
	if !strings.Contains(insts[0].String(), "ret i32 1") {
		t.Fatalf("instruction order wrong: %v", insts)
	}
	if !strings.Contains(insts[1].String(), "ret i32 0") {
		t.Fatalf("instruction order wrong: %v", insts)
	}
}
