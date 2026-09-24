package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/checks"
)

func TestBuilderCallGolden(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "call")
	defer m.Close()

	i32 := ctx.Int(32)
	addSig := ctx.Fn(i32, []llvm.AnyType{i32, i32}, false)
	m.NewFunction("add", addSig)

	fn := m.NewFunction("caller", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(fn.NewBlock("entry"))

	add, _ := m.GetFunction("add")
	x := fn.ParamAs[llvm.IntT](0)
	y := fn.ParamAs[llvm.IntT](1)
	res := b.Call[llvm.IntT](add, []llvm.AnyValue{x.Dyn(), y.Dyn()}, "r")
	if res.ArgCount() != 2 {
		t.Fatalf("ArgCount() = %d", res.ArgCount())
	}
	if got := res.Arg(0).String(); got != "i32 %0" {
		t.Fatalf("Arg(0) = %q", got)
	}
	if callee, ok := res.CalledFunction(); !ok || callee.Name() != "add" {
		t.Fatalf("CalledFunction() = %v, %v", callee.Name(), ok)
	}
	b.Ret(res)

	got := m.String()
	for _, want := range []string{
		"%r = call i32 @add(i32 %0, i32 %1)",
		"ret i32 %r",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("module output missing %q:\n%s", want, got)
		}
	}
}

func TestBuilderCallChecks(t *testing.T) {
	requireDebug(t)
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "callcheck")
	defer m.Close()

	i32 := ctx.Int(32)
	i64 := ctx.Int(64)
	m.NewFunction("add", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	fn := m.NewFunction("caller", ctx.Fn(ctx.Void(), nil, false))
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(fn.NewBlock("entry"))

	add, _ := m.GetFunction("add")
	err := llvm.Catch(func() {
		b.Call[llvm.IntT](add, []llvm.AnyValue{ctx.ConstInt(i32, 1)}, "tooFew")
	})
	if err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("arity mismatch should panic ErrTypeMismatch, got %v", err)
	}

	err = llvm.Catch(func() {
		b.Call[llvm.IntT](add, []llvm.AnyValue{ctx.ConstInt(i64, 1), ctx.ConstInt(i32, 2)}, "badType")
	})
	if err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("arg type mismatch should panic ErrTypeMismatch, got %v", err)
	}
}

func TestBuilderCallVoid(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "callvoid")
	defer m.Close()

	i32 := ctx.Int(32)
	m.NewFunction("sink", ctx.Fn(ctx.Void(), []llvm.AnyType{i32}, false))
	fn := m.NewFunction("caller", ctx.Fn(ctx.Void(), nil, false))
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(fn.NewBlock("entry"))

	sink, _ := m.GetFunction("sink")
	b.Call[llvm.VoidT](sink, []llvm.AnyValue{ctx.ConstInt(i32, 7)}, "")
	b.RetVoid()

	if got := m.String(); !strings.Contains(got, "call void @sink(i32 7)") {
		t.Fatalf("module output:\n%s", got)
	}
}

func TestBuilderPHI(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "phi")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	entry := fn.NewBlock("entry")
	loop := fn.NewBlock("loop")
	exit := fn.NewBlock("exit")

	b := NewBuilder(ctx)
	defer b.Close()

	b.MoveToEnd(entry)
	b.Br(loop)

	b.MoveToEnd(loop)
	phi := b.PHI(i32, "p")
	phi.AddIncoming(
		Incoming[llvm.IntT]{Value: fn.ParamAs[llvm.IntT](0), Block: entry},
		Incoming[llvm.IntT]{Value: ctx.ConstInt(i32, 1).Value, Block: loop},
	)
	b.Br(exit)

	b.MoveToEnd(exit)
	b.Ret(phi)

	if phi.Count() != 2 {
		t.Fatalf("phi incoming count = %d", phi.Count())
	}
	in := phi.IncomingAt(1)
	if in.Value.String() != "i32 1" || in.Block.Name() != "loop" {
		t.Fatalf("IncomingAt(1) = %v/%v", in.Value, in.Block.Name())
	}

	got := m.String()
	want := "%p = phi i32 [ %0, %entry ], [ 1, %loop ]"
	if !strings.Contains(got, want) {
		t.Fatalf("module output missing %q:\n%s", want, got)
	}

	if checks.Debug {
		err := llvm.Catch(func() {
			phi.AddIncoming(Incoming[llvm.IntT]{Value: ctx.ConstInt(ctx.Int(64), 1).Value, Block: entry})
		})
		if err == nil || err.Reason != llvm.ErrTypeMismatch {
			t.Fatalf("phi type mismatch should panic ErrTypeMismatch, got %v", err)
		}
	}
}

func TestBuilderKindChecks(t *testing.T) {
	requireDebug(t)
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "kindcheck")
	defer m.Close()

	i32 := ctx.Int(32)
	st := ctx.Struct([]llvm.AnyType{i32, ctx.Int(64)}, false)
	m.NewFunction("add", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	fn := m.NewFunction("caller", ctx.Fn(ctx.Void(), []llvm.AnyType{st, ctx.Ptr(0)}, false))
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(fn.NewBlock("entry"))

	add, _ := m.GetFunction("add")
	agg := fn.ParamAs[llvm.StructT](0)
	args := []llvm.AnyValue{ctx.ConstInt(i32, 1), ctx.ConstInt(i32, 2)}

	err := llvm.Catch(func() { b.Call[llvm.FloatT](add, args, "") })
	if err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("Call[FloatT] on i32-returning fn should panic ErrTypeMismatch, got %v", err)
	}

	err = llvm.Catch(func() {
		b.CallIndirect[llvm.FloatT](fn.ParamAs[llvm.PtrT](1), ctx.Fn(i32, nil, false), nil, "")
	})
	if err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("CallIndirect[FloatT] on i32-returning sig should panic ErrTypeMismatch, got %v", err)
	}

	err = llvm.Catch(func() { b.ExtractValue[llvm.FloatT](agg, []uint32{0}, "") })
	if err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("ExtractValue[FloatT] on i32 element should panic ErrTypeMismatch, got %v", err)
	}

	err = llvm.Catch(func() { b.ExtractValue[llvm.IntT](agg, []uint32{9}, "") })
	if err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("out-of-range index should panic ErrInvalidArg, got %v", err)
	}

	err = llvm.Catch(func() { b.ExtractValue[llvm.IntT](ctx.ConstInt(i32, 1), []uint32{0}, "") })
	if err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("indexing scalar should panic ErrInvalidArg, got %v", err)
	}
}

func TestBuilderAggregate(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "agg")
	defer m.Close()

	i32 := ctx.Int(32)
	i64 := ctx.Int(64)
	st := ctx.Struct([]llvm.AnyType{i32, i64}, false)

	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{st, i32}, false))
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(fn.NewBlock("entry"))

	agg := fn.ParamAs[llvm.StructT](0)
	val := fn.ParamAs[llvm.IntT](1)
	inserted := b.InsertValue(agg, val, []uint32{0}, "ins")
	extracted := b.ExtractValue[llvm.IntT](inserted, []uint32{0}, "ext")
	b.Ret(extracted)

	got := m.String()
	for _, want := range []string{
		"%ins = insertvalue { i32, i64 } %0, i32 %1, 0",
		"%ext = extractvalue { i32, i64 } %ins, 0",
		"ret i32 %ext",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("module output missing %q:\n%s", want, got)
		}
	}
}
