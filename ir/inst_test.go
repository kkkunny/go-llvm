package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

// TestCallArgAccessors 覆盖 Call 的实参读写、被调函数与 tail 标志角色方法。
func TestCallArgAccessors(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "inst")
	defer m.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	g := m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	b := NewBuilderAt(fn.NewBlock("entry"))
	defer b.Close()

	a := fn.ParamAs[llvm.IntT](0)
	call := b.Call[llvm.IntT](g, []llvm.AnyValue{a}, "c")
	if got := call.Arg(0); got.String() != a.String() {
		t.Fatalf("call arg 0 = %s, want %s", got, a)
	}
	call.SetArg(0, ctx.ConstInt(i32, 5).Value)
	if got := call.Arg(0).String(); got != "i32 5" {
		t.Fatalf("call arg 0 after SetArg = %s, want i32 5", got)
	}
	// 角色方法的 tail 标志读写
	call.SetTailCall(true)
	if !call.IsTailCall() {
		t.Fatal("call should be tail after SetTailCall(true)")
	}
	call.SetTailCall(false)
	if call.IsTailCall() {
		t.Fatal("call should not be tail after SetTailCall(false)")
	}
	if called, ok := call.CalledFunction(); !ok || called.Name() != "g" {
		t.Fatalf("called function = %v %v, want g", called, ok)
	}
	b.Ret(call)

	// 间接调用：CallIndirect 走指针 + 签名
	fp := m.NewFunction("fp", ctx.Fn(i32, []llvm.AnyType{ptr}, false))
	b.MoveToEnd(fp.NewBlock("entry"))
	sig := ctx.Fn(i32, []llvm.AnyType{i32}, false)
	ic := b.CallIndirect[llvm.IntT](fp.ParamAs[llvm.PtrT](0), sig, []llvm.AnyValue{ctx.ConstInt(i32, 1).Value}, "ic")
	b.Ret(ic)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	out := m.String()
	for _, want := range []string{
		"call i32 @g(i32 5)",
		"call i32 %0(i32 1)",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("module output missing %q:\n%s", want, out)
		}
	}
}

// TestCallArgPrecheck 越界实参下标在调试构建下必须 panic。
func TestCallArgPrecheck(t *testing.T) {
	requireDebug(t)
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "instpre")
	defer m.Close()

	i32 := ctx.Int(32)
	g := m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	b := NewBuilderAt(fn.NewBlock("entry"))
	defer b.Close()

	call := b.Call[llvm.IntT](g, []llvm.AnyValue{ctx.ConstInt(i32, 1).Value}, "c")
	if err := llvm.Catch(func() { call.Arg(3) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("Call.Arg out of range should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { call.SetArg(3, ctx.ConstInt(i32, 1).Value) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("Call.SetArg out of range should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { fn.Param(5) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("Function.Param out of range should panic ErrInvalidArg, got %v", err)
	}
	b.Ret(call)
}
