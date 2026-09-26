package jit

import (
	"math"
	"sync"
	"testing"
	"unsafe"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
	"github.com/kkkunny/go-llvm/target"
)

func newJIT(t testing.TB) *LLJIT {
	t.Helper()
	if err := target.InitNative(); err != nil {
		t.Fatal(err)
	}
	j, err := NewLLJIT()
	if err != nil {
		t.Fatal(err)
	}
	return j
}

func retModule(t testing.TB, name string, v int64) (*llvm.Context, *ir.Module) {
	t.Helper()
	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, name)
	i32 := ctx.Int(32)
	fn := m.NewFunction("answer", ctx.Fn(i32, nil, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(ctx.ConstInt(i32, uint64(v)))
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	return ctx, m
}

func TestLLJITAddAndLookup(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	if j.Triple() == "" || j.DataLayoutStr() == "" {
		t.Fatalf("triple/data layout should be non-empty: %q, %q", j.Triple(), j.DataLayoutStr())
	}

	ctx, m := retModule(t, "lookup", 42)
	if err := j.AddIRModule(m); err != nil {
		t.Fatal(err)
	}
	if m.Lifetime().Alive() || ctx.Alive() {
		t.Fatal("module and context should be invalidated after ownership transfer")
	}

	p, err := j.Lookup("answer")
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("Lookup(answer) should resolve an address")
	}

	if _, err := j.Lookup("no_such_symbol"); err == nil || err.(*llvm.Error).Reason != llvm.ErrNotFound {
		t.Fatalf("missing symbol should return ErrNotFound, got %v", err)
	}
}

func TestLLJITMapSymbol(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	_, src := retModule(t, "src", 7)
	if err := j.AddIRModule(src); err != nil {
		t.Fatal(err)
	}
	addr, err := j.Lookup("answer")
	if err != nil {
		t.Fatal(err)
	}
	if err := j.MapSymbol("answer_alias", addr); err != nil {
		t.Fatal(err)
	}
	if err := llvm.Catch(func() { j.MapSymbol("nil_symbol", nil) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil symbol should panic ErrInvalidArg, got %v", err)
	}

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "alias_user")
	i32 := ctx.Int(32)
	decl := m.NewFunction("answer_alias", ctx.Fn(i32, nil, false))
	fn := m.NewFunction("call_alias", ctx.Fn(i32, nil, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(b.Call[llvm.IntT](decl.Value, nil, "").Value)
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	if err := m.Verify(); err != nil {
		t.Fatal(err)
	}
	if err := j.AddIRModule(m); err != nil {
		t.Fatal(err)
	}

	if _, err := j.Lookup("call_alias"); err != nil {
		t.Fatal(err)
	}
}

func TestLLJITAddProcessSymbols(t *testing.T) {
	j := newJIT(t)
	defer j.Close()
	if err := j.AddProcessSymbols(); err != nil {
		t.Fatal(err)
	}

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "process_symbols")
	f64 := ctx.Float(llvm.FloatDouble)
	sin := m.NewFunction("sin", ctx.Fn(f64, []llvm.AnyType{f64}, false))
	fn := m.NewFunction("neg_sin", ctx.Fn(f64, []llvm.AnyType{f64}, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	x := fn.ParamAs[llvm.FloatT](0)
	res := b.Call[llvm.FloatT](sin, []llvm.AnyValue{x.Dyn()}, "")
	b.Ret(b.FNeg(res.Value, ""))
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	if err := m.Verify(); err != nil {
		t.Fatal(err)
	}
	if err := j.AddIRModule(m); err != nil {
		t.Fatal(err)
	}

	negSin, err := j.Func[func(float64) float64]("neg_sin")
	if err != nil {
		t.Fatal(err)
	}
	if got := negSin(math.Pi / 2); math.Abs(got+1) > 1e-9 {
		t.Fatalf("neg_sin(pi/2) = %v, want -1", got)
	}
}

func TestLLJITAddObjectFile(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	ctx, m := retModule(t, "obj", 3)
	tm, err := target.NewTargetMachine(mustNativeTarget(t), target.DefaultTriple(), target.HostCPUName(), target.HostCPUFeatures(), target.OptNone, target.RelocPIC, target.CodeModelDefault)
	if err != nil {
		t.Fatal(err)
	}
	defer tm.Close()
	tm.ApplyTo(m)
	obj, err := tm.Emit(m, target.ObjectFile)
	if err != nil {
		t.Fatal(err)
	}
	if err := j.AddObjectFile(obj); err != nil {
		t.Fatal(err)
	}
	if obj.Alive() {
		t.Fatal("object buffer should be invalidated after ownership transfer")
	}
	if _, err := j.Lookup("answer"); err != nil {
		t.Fatal(err)
	}
	_ = ctx
}

func TestLLJITConcurrentFunc(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	_, m := retModule(t, "concurrent", 42)
	if err := j.AddIRModule(m); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f, err := j.Func[func() int32]("answer")
			if err != nil {
				t.Errorf("Func: %v", err)
				return
			}
			if got := f(); got != 42 {
				t.Errorf("answer() = %d, want 42", got)
			}
		}()
	}
	wg.Wait()
}

func mustNativeTarget(t testing.TB) target.Target {
	t.Helper()
	native, err := target.NativeTarget()
	if err != nil {
		t.Fatal(err)
	}
	return native
}

func TestLLJITChecks(t *testing.T) {
	j := newJIT(t)
	_, m := retModule(t, "checks", 1)
	if err := j.AddIRModule(m); err != nil {
		t.Fatal(err)
	}
	if err := llvm.Catch(func() { j.AddIRModule(m) }); err == nil || err.Reason != llvm.ErrUseAfterFree {
		t.Fatalf("re-adding transferred module should panic ErrUseAfterFree, got %v", err)
	}
	if err := j.Close(); err != nil {
		t.Fatal(err)
	}
	if err := j.Close(); err == nil || err.(*llvm.Error).Reason != llvm.ErrClosed {
		t.Fatalf("double close should return ErrClosed, got %v", err)
	}
	if err := llvm.Catch(func() { j.Triple() }); err == nil || err.Reason != llvm.ErrUseAfterFree {
		t.Fatalf("use after close should panic ErrUseAfterFree, got %v", err)
	}

	j2 := newJIT(t)
	defer j2.Close()
	if err := llvm.Catch(func() { j2.AddObjectFile(nil) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil object buffer should panic ErrInvalidArg, got %v", err)
	}
	_ = unsafe.Pointer(nil)
}
