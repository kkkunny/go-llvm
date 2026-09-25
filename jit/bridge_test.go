package jit

import (
	"testing"
	"unsafe"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
)

// addModule 定义 `i32 add(i32, i32)`
func addModule(t *testing.T) *ir.Module {
	t.Helper()
	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "add")
	i32 := ctx.Int(32)
	fn := m.NewFunction("add", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(b.Add(fn.ParamAs[llvm.IntT](0), fn.ParamAs[llvm.IntT](1), ""))
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	return m
}

func TestLLJITFunc(t *testing.T) {
	j := newJIT(t)
	defer j.Close()
	if err := j.AddIRModule(addModule(t)); err != nil {
		t.Fatal(err)
	}

	add, err := j.Func[func(int32, int32) int32]("add")
	if err != nil {
		t.Fatal(err)
	}
	if got := add(20, 22); got != 42 {
		t.Fatalf("add(20, 22) = %d", got)
	}

	if _, err := j.Func[func(int32, int32) int32]("missing"); err == nil || err.(*llvm.Error).Reason != llvm.ErrNotFound {
		t.Fatalf("missing function should return ErrNotFound, got %v", err)
	}
	if _, err := j.Func[func(string) string]("add"); err == nil || err.(*llvm.Error).Reason != llvm.ErrUnsupported {
		t.Fatalf("unsupported signature should return ErrUnsupported, got %v", err)
	}
}

func TestLLJITMapFunc(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	host := func(a int32, b int32) int32 { return a*b + 1 }
	if err := j.MapFunc("host_mul_add", host); err != nil {
		t.Fatal(err)
	}

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "caller")
	i32 := ctx.Int(32)
	decl := m.NewFunction("host_mul_add", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	fn := m.NewFunction("caller", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(b.Call[llvm.IntT](decl.Value, []llvm.AnyValue{fn.ParamAs[llvm.IntT](0).Dyn(), fn.ParamAs[llvm.IntT](1).Dyn()}, "").Value)
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	if err := m.Verify(); err != nil {
		t.Fatal(err)
	}
	if err := j.AddIRModule(m); err != nil {
		t.Fatal(err)
	}

	caller, err := j.Func[func(int32, int32) int32]("caller")
	if err != nil {
		t.Fatal(err)
	}
	if got := caller(6, 7); got != 43 {
		t.Fatalf("caller(6, 7) = %d", got)
	}
}

func TestLLJITFuncFloatsAndPointers(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "mixed")
	f64 := ctx.Float(llvm.FloatDouble)
	fn := m.NewFunction("scale", ctx.Fn(f64, []llvm.AnyType{f64, f64}, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(b.FMul(fn.ParamAs[llvm.FloatT](0), fn.ParamAs[llvm.FloatT](1), ""))
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	if err := j.AddIRModule(m); err != nil {
		t.Fatal(err)
	}
	scale, err := j.Func[func(float64, float64) float64]("scale")
	if err != nil {
		t.Fatal(err)
	}
	if got := scale(1.5, 4); got != 6 {
		t.Fatalf("scale(1.5, 4) = %v", got)
	}
}

// TestSlotToPtr 槽位位模式与指针互转：非零槽位保持地址，零槽位得到 nil
func TestSlotToPtr(t *testing.T) {
	v := 42
	p := unsafe.Pointer(&v)
	if got := slotToPtr(uint64(uintptr(p))); got != p {
		t.Fatalf("slotToPtr round trip: got %p, want %p", got, p)
	}
	if got := slotToPtr(0); got != nil {
		t.Fatalf("slotToPtr(0) = %p, want nil", got)
	}
}

// TestLLJITPointerBridge 指针经桥的完整往返：Go 实参 → JIT → Go 回调 → JIT → 返回值，
// 覆盖槽位与指针互转及指针身份的保持
func TestLLJITPointerBridge(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	// native→Go 方向：回调把收到的指针原样返回
	if err := j.MapFunc("go_ptr_id", func(p unsafe.Pointer) unsafe.Pointer { return p }); err != nil {
		t.Fatal(err)
	}

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "ptr_bridge")
	ptr := ctx.Ptr(0)
	decl := m.NewFunction("go_ptr_id", ctx.Fn(ptr, []llvm.AnyType{ptr}, false))
	fn := m.NewFunction("echo_ptr", ctx.Fn(ptr, []llvm.AnyType{ptr}, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(b.Call[llvm.PtrT](decl.Value, []llvm.AnyValue{fn.ParamAs[llvm.PtrT](0).Dyn()}, "").Value)
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	if err := j.AddIRModule(m); err != nil {
		t.Fatal(err)
	}

	echo, err := j.Func[func(unsafe.Pointer) unsafe.Pointer]("echo_ptr")
	if err != nil {
		t.Fatal(err)
	}

	v := int32(1234)
	got := echo(unsafe.Pointer(&v))
	if got != unsafe.Pointer(&v) {
		t.Fatalf("echo_ptr(&v) = %p, want %p", got, unsafe.Pointer(&v))
	}
	if *(*int32)(got) != 1234 {
		t.Fatalf("echo_ptr(&v) deref = %d, want 1234", *(*int32)(got))
	}

	// 空指针往返仍为 nil
	if got := echo(nil); got != nil {
		t.Fatalf("echo_ptr(nil) = %p, want nil", got)
	}
}

func TestLLJITRunMain(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "main")
	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	fn := m.NewFunction("main", ctx.Fn(i32, []llvm.AnyType{i32, ptr, ptr}, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	argc := fn.ParamAs[llvm.IntT](0)
	b.Ret(b.Add(argc, ctx.ConstInt(i32, 3), ""))
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	if err := j.AddIRModule(m); err != nil {
		t.Fatal(err)
	}

	code, err := j.RunMain([]string{"prog", "a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	if code != 6 { // argc=3，加 3
		t.Fatalf("RunMain exit code = %d", code)
	}
}
