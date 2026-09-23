package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func buildAtomicModule(t *testing.T) (*llvm.Context, *Module, *Builder) {
	t.Helper()
	ctx := llvm.NewContext()
	m := NewModule(ctx, "atomic")
	b := NewBuilder(ctx)
	return ctx, m, b
}

func TestBuilderAtomics(t *testing.T) {
	ctx, m, b := buildAtomicModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{ctx.Ptr(0)}, false))
	entry := fn.NewBlock("entry")
	b.MoveToEnd(entry)
	p := fn.ParamAs[llvm.PtrT](0)
	one := ctx.ConstInt(i32, 1).Value

	fe := b.Fence(llvm.AtomicSequentiallyConsistent, false)
	if fe.Ordering() != llvm.AtomicSequentiallyConsistent {
		t.Fatalf("fence ordering = %v", fe.Ordering())
	}
	fe.SetOrdering(llvm.AtomicAcquire)
	if fe.Ordering() != llvm.AtomicAcquire {
		t.Fatalf("fence ordering after set = %v", fe.Ordering())
	}
	fe.SetSingleThread(true)
	if !fe.IsSingleThread() {
		t.Fatal("fence should be singlethread after set")
	}
	fe.SetSingleThread(false)

	rm := b.AtomicRMW(llvm.RMWAdd, p, one, llvm.AtomicMonotonic, false, "old")
	if rm.Ordering() != llvm.AtomicMonotonic || rm.Op() != llvm.RMWAdd {
		t.Fatalf("rmw = %v %v", rm.Ordering(), rm.Op())
	}
	rm.SetVolatile(true)
	rm.SetAlign(4)
	if !rm.IsVolatile() || rm.Align() != 4 {
		t.Fatalf("rmw volatile/align = %v %d", rm.IsVolatile(), rm.Align())
	}

	cx := b.CmpXchg(p, one, ctx.ConstInt(i32, 2).Value, llvm.AtomicAcquire, llvm.AtomicMonotonic, false, "cx")
	if cx.SuccessOrdering() != llvm.AtomicAcquire || cx.FailureOrdering() != llvm.AtomicMonotonic {
		t.Fatalf("cmpxchg orderings = %v %v", cx.SuccessOrdering(), cx.FailureOrdering())
	}
	if cx.IsWeak() {
		t.Fatal("cmpxchg should not be weak by default")
	}
	cx.SetWeak(true)
	if !cx.IsWeak() {
		t.Fatal("cmpxchg should be weak after set")
	}
	cx.SetWeak(false)
	cx.SetAlign(4)

	b.RetVoid()
	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := m.String()
	for _, want := range []string{
		"fence acquire",
		"atomicrmw volatile add ptr %0, i32 1 monotonic, align 4",
		"cmpxchg ptr %0, i32 1, i32 2 acquire monotonic, align 4",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q:\n%s", want, got)
		}
	}
}

func TestBuilderAtomicPrecheck(t *testing.T) {
	ctx, m, b := buildAtomicModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{ctx.Ptr(0)}, false))
	entry := fn.NewBlock("entry")
	b.MoveToEnd(entry)
	p := fn.ParamAs[llvm.PtrT](0)
	one := ctx.ConstInt(i32, 1).Value

	// fence 不允许 monotonic/not_atomic
	if err := llvm.Catch(func() { b.Fence(llvm.AtomicMonotonic, false) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("fence monotonic should panic ErrInvalidArg, got %v", err)
	}
	// atomicrmw 不允许 unordered/not_atomic
	if err := llvm.Catch(func() {
		b.AtomicRMW(llvm.RMWAdd, p, one, llvm.AtomicNotAtomic, false, "")
	}); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("rmw not_atomic should panic ErrInvalidArg, got %v", err)
	}
	// cmpxchg 失败序不得为 release/acq_rel
	if err := llvm.Catch(func() {
		b.CmpXchg(p, one, one, llvm.AtomicAcquire, llvm.AtomicRelease, false, "")
	}); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("cmpxchg failure release should panic ErrInvalidArg, got %v", err)
	}
	// cmpxchg 失败序不得强于成功序
	if err := llvm.Catch(func() {
		b.CmpXchg(p, one, one, llvm.AtomicMonotonic, llvm.AtomicAcquire, false, "")
	}); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("cmpxchg failure stronger should panic ErrInvalidArg, got %v", err)
	}
	// cmp/new 类型必须一致
	if err := llvm.Catch(func() {
		b.CmpXchg(p, one, ctx.ConstFloat(ctx.Float(llvm.FloatDouble), 1).Value, llvm.AtomicAcquire, llvm.AtomicMonotonic, false, "")
	}); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("cmpxchg type mismatch should panic ErrTypeMismatch, got %v", err)
	}
	b.RetVoid()
}

func TestBuilderLoadStoreAtomic(t *testing.T) {
	ctx, m, b := buildAtomicModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{ctx.Ptr(0)}, false))
	entry := fn.NewBlock("entry")
	b.MoveToEnd(entry)
	p := fn.ParamAs[llvm.PtrT](0)

	ld := b.Load(p, i32, "v")
	ld.SetVolatile(true)
	ld.SetOrdering(llvm.AtomicAcquire)
	ld.SetAlign(4)
	if !ld.IsVolatile() || ld.Ordering() != llvm.AtomicAcquire || ld.Align() != 4 {
		t.Fatalf("load flags = %v %v %d", ld.IsVolatile(), ld.Ordering(), ld.Align())
	}

	st := b.Store(ctx.ConstInt(i32, 1).Value, p)
	st.SetVolatile(true)
	st.SetOrdering(llvm.AtomicRelease)
	if !st.IsVolatile() || st.Ordering() != llvm.AtomicRelease {
		t.Fatalf("store flags = %v %v", st.IsVolatile(), st.Ordering())
	}

	b.RetVoid()
	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := m.String()
	if !strings.Contains(got, "load atomic volatile i32, ptr %0 acquire, align 4") {
		t.Fatalf("missing atomic load:\n%s", got)
	}
	if !strings.Contains(got, "store atomic volatile i32 1, ptr %0 release") {
		t.Fatalf("missing atomic store:\n%s", got)
	}

	// load 不得设 release/acq_rel；store 不得设 acquire/acq_rel
	if err := llvm.Catch(func() { ld.SetOrdering(llvm.AtomicRelease) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("load release should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { st.SetOrdering(llvm.AtomicAcquire) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("store acquire should panic ErrInvalidArg, got %v", err)
	}
}
