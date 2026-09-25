package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func memModule(t *testing.T) (*llvm.Context, *Module, *Builder, Function) {
	t.Helper()
	ctx := llvm.NewContext()
	m := NewModule(ctx, "mem")
	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	b := NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	return ctx, m, b, fn
}

func TestBuilderMemoryGolden(t *testing.T) {
	ctx, m, b, fn := memModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	p := fn.ParamAs[llvm.IntT](0)

	slot := b.Alloca(i32, "slot")
	slot.SetAlign(8)
	b.Store(p, slot)
	st := b.Store(p, slot)
	st.SetAlign(4)
	if got := st.Align(); got != 4 {
		t.Fatalf("store align = %d, want 4", got)
	}
	loaded := b.Load(slot, i32, "loaded")
	b.Ret(loaded)

	got := m.String()
	for _, want := range []string{
		"%slot = alloca i32, align 8",
		"store i32 %0, ptr %slot",
		"store i32 %0, ptr %slot, align 4",
		"%loaded = load i32, ptr %slot",
		"ret i32 %loaded",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("module output missing %q:\n%s", want, got)
		}
	}
}

func TestBuilderGEPGolden(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "gep")
	defer m.Close()

	i32 := ctx.Int(32)
	i64 := ctx.Int(64)
	arrTy := ctx.Array(i32, 4)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{ctx.Ptr(0)}, false))
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(fn.NewBlock("entry"))

	p := fn.ParamAs[llvm.PtrT](0)
	idx := []llvm.ValueRef[llvm.IntT]{ctx.ConstInt(i64, 0), ctx.ConstInt(i64, 1)}
	elem := b.GEP(arrTy, p, idx, "elem")
	b.InBoundsGEP(i32, elem, []llvm.ValueRef[llvm.IntT]{ctx.ConstInt(i64, 0)}, "inner")
	b.RetVoid()

	got := m.String()
	for _, want := range []string{
		"%elem = getelementptr [4 x i32], ptr %0, i64 0, i64 1",
		"getelementptr inbounds i32, ptr %elem, i64 0",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("module output missing %q:\n%s", want, got)
		}
	}
}

func TestBuilderCastsGolden(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "cast")
	defer m.Close()

	i8 := ctx.Int(8)
	i32 := ctx.Int(32)
	i64 := ctx.Int(64)
	f32 := ctx.Float(llvm.FloatSingle)
	f64 := ctx.Float(llvm.FloatDouble)

	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{i8, i32, f32, f64, ctx.Ptr(0)}, false))
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(fn.NewBlock("entry"))

	small := fn.ParamAs[llvm.IntT](0)
	mid := fn.ParamAs[llvm.IntT](1)
	single := fn.ParamAs[llvm.FloatT](2)
	dbl := fn.ParamAs[llvm.FloatT](3)
	ptr := fn.ParamAs[llvm.PtrT](4)

	b.Trunc(mid, i8, "trunc")
	b.ZExt(small, i64, "zext")
	b.SExt(small, i64, "sext")
	b.FPTrunc(dbl, f32, "fptrunc")
	b.FPExt(single, f64, "fpext")
	b.FPToUI(single, i32, "fptoui")
	b.FPToSI(single, i32, "fptosi")
	b.UIToFP(mid, f64, "uitofp")
	b.SIToFP(mid, f64, "sitofp")
	b.PtrToInt(ptr, i64, "ptrtoint")
	b.IntToPtr(mid, ctx.Ptr(0), "inttoptr")
	b.BitCast(mid, f32, "bitcast")
	b.RetVoid()

	got := m.String()
	for _, want := range []string{
		"%trunc = trunc i32 %1 to i8",
		"%zext = zext i8 %0 to i64",
		"%sext = sext i8 %0 to i64",
		"%fptrunc = fptrunc double %3 to float",
		"%fpext = fpext float %2 to double",
		"%fptoui = fptoui float %2 to i32",
		"%fptosi = fptosi float %2 to i32",
		"%uitofp = uitofp i32 %1 to double",
		"%sitofp = sitofp i32 %1 to double",
		"%ptrtoint = ptrtoint ptr %4 to i64",
		"%inttoptr = inttoptr i32 %1 to ptr",
		"%bitcast = bitcast i32 %1 to float",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("module output missing %q:\n%s", want, got)
		}
	}
}

func TestBuilderMemIntrinsics(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "memintrin")
	defer m.Close()

	i8 := ctx.Int(8)
	i64 := ctx.Int(64)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{ctx.Ptr(0), ctx.Ptr(0), i64}, false))
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(fn.NewBlock("entry"))

	dst := fn.ParamAs[llvm.PtrT](0)
	src := fn.ParamAs[llvm.PtrT](1)
	n := fn.ParamAs[llvm.IntT](2)

	b.MemSet(dst, ctx.ConstInt(i8, 0), n, 4)
	b.MemCpy(dst, 4, src, 4, n)
	b.MemMove(dst, 8, src, 8, n)
	b.RetVoid()

	got := m.String()
	for _, want := range []string{"memset", "memcpy", "memmove", "i64 %2"} {
		if !strings.Contains(got, want) {
			t.Fatalf("module output missing %q:\n%s", want, got)
		}
	}
}

func TestBuilderMallocFree(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "malloc")
	defer m.Close()
	m.SetDataLayout("e-m:e-p:64:64-i64:64-n8:16:32:64-S128")

	i32 := ctx.Int(32)
	i64 := ctx.Int(64)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{i64}, false))
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(fn.NewBlock("entry"))

	n := fn.ParamAs[llvm.IntT](0)
	p1 := b.Malloc(i32, "one")
	p2 := b.MallocArray(i32, n, "many")
	b.Free(p1)
	b.Free(p2)
	b.RetVoid()

	got := m.String()
	// 注：LLVM C API 的 LLVMBuildMalloc 固定以 i32 表达 size（已用 C 程序核实）
	for _, want := range []string{"call ptr @malloc(i32", "call void @free(ptr", "@malloc", "@free"} {
		if !strings.Contains(got, want) {
			t.Fatalf("module output missing %q:\n%s", want, got)
		}
	}
}

func TestBuilderAlignPrecheck(t *testing.T) {
	requireDebug(t)
	ctx, m, b, _ := memModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	err := llvm.Catch(func() { b.Alloca(ctx.Int(32), "x").SetAlign(3) })
	if err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("non-power-of-two align should panic ErrInvalidArg, got %v", err)
	}

	err = llvm.Catch(func() {
		b.MemCpy(ctx.Ptr(0).Zero(), 3, ctx.Ptr(0).Zero(), 4, ctx.ConstInt(ctx.Int(64), 1).Value)
	})
	if err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("non-power-of-two dst align should panic ErrInvalidArg, got %v", err)
	}
}

// TestBuilderMemTypePrecheck 覆盖 Alloca/GEP 的类型前置校验（nil/跨 Context）。
func TestBuilderMemTypePrecheck(t *testing.T) {
	ctx, m, b, _ := memModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	// nil 类型（崩溃类地板）
	if err := llvm.Catch(func() { b.Alloca(nil, "") }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil alloca type should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { b.GEP(nil, ctx.Ptr(0).Zero(), nil, "") }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil GEP element type should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { b.Malloc(nil, "") }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil malloc type should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() {
		b.MallocArray(nil, ctx.ConstInt(ctx.Int(64), 1).Value, "")
	}); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil malloc-array type should panic ErrInvalidArg, got %v", err)
	}

	// 跨 Context 类型
	ctx2 := llvm.NewContext()
	defer ctx2.Close()
	if err := llvm.Catch(func() { b.Alloca(ctx2.Int(32), "") }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("foreign alloca type should panic ErrCrossContext, got %v", err)
	}
	if err := llvm.Catch(func() { b.GEP(ctx2.Int(32), ctx.Ptr(0).Zero(), nil, "") }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("foreign GEP element type should panic ErrCrossContext, got %v", err)
	}
	if err := llvm.Catch(func() { b.Malloc(ctx2.Int(32), "") }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("foreign malloc type should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { b.Load(ctx.Ptr(0).Zero(), ctx2.Int(32), "") }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("foreign load type should panic ErrCrossContext, got %v", err)
	}
}

func TestIsNullPtrDiff(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "t")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Int(1), []llvm.AnyType{ctx.Ptr(0), ctx.Ptr(0)}, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()

	p := fn.ParamAs[llvm.PtrT](0)
	q := fn.ParamAs[llvm.PtrT](1)
	n := b.IsNull(p, "isnull")
	nn := b.IsNotNull(p, "isnotnull")
	d := b.PtrDiff(i32, p, q, "diff")
	if d.IsNil() {
		t.Fatalf("PtrDiff produced nil")
	}
	b.Ret(nn)

	got := m.String()
	if !strings.Contains(got, "icmp eq ptr") {
		t.Fatalf("IR missing icmp eq ptr:\n%s", got)
	}
	if !strings.Contains(got, "icmp ne ptr") {
		t.Fatalf("IR missing icmp ne ptr:\n%s", got)
	}
	if !strings.Contains(got, "ptrtoint") {
		t.Fatalf("IR missing ptrtoint:\n%s", got)
	}
	_ = n
}
