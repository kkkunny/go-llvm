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

	// InsertValue 负向：空索引路径 / 元素类型不符 / 索引越界
	err = llvm.Catch(func() { b.InsertValue(agg, ctx.ConstInt(i32, 1), nil, "") })
	if err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("empty insert index path should panic ErrInvalidArg, got %v", err)
	}
	err = llvm.Catch(func() { b.InsertValue(agg, ctx.ConstInt(i32, 1), []uint32{1}, "") })
	if err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("insert i32 into i64 slot should panic ErrTypeMismatch, got %v", err)
	}
	err = llvm.Catch(func() { b.InsertValue(agg, ctx.ConstInt(ctx.Int(64), 1), []uint32{9}, "") })
	if err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("insert out-of-range index should panic ErrInvalidArg, got %v", err)
	}

	// 跨 Context 的间接调用签名（崩溃类地板：始终校验）
	ctx2 := llvm.NewContext()
	defer ctx2.Close()
	sig2 := ctx2.Fn(ctx2.Int(32), nil, false)
	err = llvm.Catch(func() {
		b.CallIndirect[llvm.IntT](fn.ParamAs[llvm.PtrT](1), sig2, nil, "")
	})
	if err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("foreign signature should panic ErrCrossContext, got %v", err)
	}
}

// TestBuilderAggregateIndexPaths 覆盖 InsertValue/ExtractValue 的数组/嵌套结构体元素类型路径。
func TestBuilderAggregateIndexPaths(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "aggpath")
	defer m.Close()

	i32 := ctx.Int(32)
	arr := ctx.Array(i32, 4)
	inner := ctx.Struct([]llvm.AnyType{i32, i32}, false)
	st := ctx.Struct([]llvm.AnyType{arr, inner}, false)
	b := NewBuilder(ctx)
	defer b.Close()

	// 数组聚合的单层插入/提取
	arrFn := m.NewFunction("arr", ctx.Fn(i32, []llvm.AnyType{arr}, false))
	b.MoveToEnd(arrFn.NewBlock("entry"))
	av := arrFn.ParamAs[llvm.ArrayT](0)
	ins := b.InsertValue(av, ctx.ConstInt(i32, 7), []uint32{1}, "ins")
	ex := b.ExtractValue[llvm.IntT](ins, []uint32{1}, "ex")
	b.Ret(ex)

	// 结构体 -> 数组 / 结构体 -> 嵌套结构体的多层提取
	stFn := m.NewFunction("st", ctx.Fn(i32, []llvm.AnyType{st}, false))
	b.MoveToEnd(stFn.NewBlock("entry"))
	sv := stFn.ParamAs[llvm.StructT](0)
	a := b.ExtractValue[llvm.IntT](sv, []uint32{0, 1}, "a")
	n := b.ExtractValue[llvm.IntT](sv, []uint32{1, 0}, "n")
	b.Ret(b.Add(a, n, "sum"))

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	out := m.String()
	for _, want := range []string{
		"%ins = insertvalue [4 x i32] %0, i32 7, 1",
		"%ex = extractvalue [4 x i32] %ins, 1",
		// 多层路径由单层指令链实现，LLVM 自动为后续层去重命名
		"%a = extractvalue { [4 x i32], { i32, i32 } } %0, 0",
		"%a1 = extractvalue [4 x i32] %a, 1",
		"%n = extractvalue { [4 x i32], { i32, i32 } } %0, 1",
		"%n2 = extractvalue { i32, i32 } %n, 0",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("module output missing %q:\n%s", want, out)
		}
	}
}

// TestAggregateVectorPathPrecheck 覆盖 elementTypeAt 的向量下标越界拒绝。
// 注：LLVM 的 extractvalue/insertvalue 不接受向量聚合（langref：operand must be aggregate type），
// 合法向量下标会让底层 LLVM 崩溃，故这里只验证越界路径在调试层被拦下。
func TestAggregateVectorPathPrecheck(t *testing.T) {
	requireDebug(t)
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "aggvec")
	defer m.Close()

	i32 := ctx.Int(32)
	st := ctx.Struct([]llvm.AnyType{i32, ctx.Vec(i32, 4)}, false)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{st}, false))
	b := NewBuilderAt(fn.NewBlock("entry"))
	defer b.Close()
	agg := fn.ParamAs[llvm.StructT](0)

	if err := llvm.Catch(func() { b.ExtractValue[llvm.IntT](agg, []uint32{1, 9}, "") }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("vector index out of range should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { b.InsertValue(agg, ctx.ConstInt(i32, 1), []uint32{1, 9}, "") }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("vector insert index out of range should panic ErrInvalidArg, got %v", err)
	}
	b.RetVoid()
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
