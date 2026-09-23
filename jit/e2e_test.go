package jit

import (
	"testing"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
	"github.com/kkkunny/go-llvm/target"
)

// fibModule 定义递归 fib（含 PHI 循环与自引用调用）
func fibModule(t *testing.T) *ir.Module {
	t.Helper()
	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "fib")
	i32 := ctx.Int(32)
	fn := m.NewFunction("fib", ctx.Fn(i32, []llvm.AnyType{i32}, false))

	entry := fn.NewBlock("entry")
	recur := fn.NewBlock("recur")
	base := fn.NewBlock("base")
	done := fn.NewBlock("done")

	b := ir.NewBuilder(ctx)
	n := fn.ParamAs[llvm.IntT](0)
	b.MoveToEnd(entry)
	b.CondBr(b.ICmp(llvm.IntSLT, n, ctx.ConstInt(i32, 2), ""), base, recur)

	b.MoveToEnd(base)
	b.Br(done)

	b.MoveToEnd(recur)
	a := b.Sub(n, ctx.ConstInt(i32, 1), "")
	bv := b.Sub(n, ctx.ConstInt(i32, 2), "")
	fa := b.Call[llvm.IntT](fn.Value, []llvm.AnyValue{a.Dyn()}, "")
	fb := b.Call[llvm.IntT](fn.Value, []llvm.AnyValue{bv.Dyn()}, "")
	sum := b.Add(fa.Value, fb.Value, "")
	b.Br(done)

	b.MoveToEnd(done)
	res := b.PHI(i32, "res")
	res.AddIncoming(
		ir.Incoming[llvm.IntT]{Value: n, Block: base},
		ir.Incoming[llvm.IntT]{Value: sum, Block: recur},
	)
	b.Ret(res.Value)

	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	if err := m.Verify(); err != nil {
		t.Fatal(err)
	}
	return m
}

func TestE2EFib(t *testing.T) {
	j := newJIT(t)
	defer j.Close()
	if err := j.AddIRModule(fibModule(t)); err != nil {
		t.Fatal(err)
	}

	fib, err := j.Func[func(int32) int32]("fib")
	if err != nil {
		t.Fatal(err)
	}
	want := []int32{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
	for n, w := range want {
		if got := fib(int32(n)); got != w {
			t.Fatalf("fib(%d) = %d, want %d", n, got, w)
		}
	}
}

func TestE2EGoCallback(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	calls := 0
	callback := func(x int32) int32 {
		calls++
		return x * 3
	}
	if err := j.MapFunc("go_callback", callback); err != nil {
		t.Fatal(err)
	}

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "cb")
	i32 := ctx.Int(32)
	decl := m.NewFunction("go_callback", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn := m.NewFunction("twice", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	x := fn.ParamAs[llvm.IntT](0)
	first := b.Call[llvm.IntT](decl.Value, []llvm.AnyValue{x.Dyn()}, "")
	second := b.Call[llvm.IntT](decl.Value, []llvm.AnyValue{first.Value.Dyn()}, "")
	b.Ret(second.Value)
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	if err := m.Verify(); err != nil {
		t.Fatal(err)
	}
	if err := j.AddIRModule(m); err != nil {
		t.Fatal(err)
	}

	twice, err := j.Func[func(int32) int32]("twice")
	if err != nil {
		t.Fatal(err)
	}
	if got := twice(5); got != 45 {
		t.Fatalf("twice(5) = %d", got)
	}
	if calls != 2 {
		t.Fatalf("callback invoked %d times, want 2", calls)
	}
}

// TestE2EObjectRoundTrip 验证 Emit(ObjectFile) 的产物可被 JIT 加载并调用
func TestE2EObjectRoundTrip(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "aot")
	i32 := ctx.Int(32)
	fn := m.NewFunction("square", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(b.Mul(fn.ParamAs[llvm.IntT](0), fn.ParamAs[llvm.IntT](0), ""))
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}

	tm, err := target.NewTargetMachine(mustNativeTarget(t), target.DefaultTriple(), target.HostCPUName(), target.HostCPUFeatures(), target.OptNone, target.RelocPIC, target.CodeModelDefault)
	if err != nil {
		t.Fatal(err)
	}
	defer tm.Close()
	tm.ApplyTo(m)
	if err := m.Verify(); err != nil {
		t.Fatal(err)
	}
	obj, err := tm.Emit(m, target.ObjectFile)
	if err != nil {
		t.Fatal(err)
	}
	if err := j.AddObjectFile(obj); err != nil {
		t.Fatal(err)
	}

	square, err := j.Func[func(int32) int32]("square")
	if err != nil {
		t.Fatal(err)
	}
	if got := square(9); got != 81 {
		t.Fatalf("square(9) = %d", got)
	}
}

// TestE2EIRRoundTrip 验证模块 IR 打印 → 解析 → JIT 执行的完整链路
func TestE2EIRRoundTrip(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	ctx := llvm.NewContext()
	defer ctx.Close()
	buf := llvm.NewMemoryBuffer([]byte(fibModule(t).String()), "fib.ll")
	defer buf.Close()
	m, err := ir.ParseIR(ctx, buf)
	if err != nil {
		t.Fatal(err)
	}
	if err := j.AddIRModule(m); err != nil {
		t.Fatal(err)
	}
	fib, err := j.Func[func(int32) int32]("fib")
	if err != nil {
		t.Fatal(err)
	}
	if got := fib(10); got != 55 {
		t.Fatalf("fib(10) = %d", got)
	}
}
