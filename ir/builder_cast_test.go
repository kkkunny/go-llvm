package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestFreezeAndCasts(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "t")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()

	a := fn.ParamAs[llvm.IntT](0)
	fr := b.Freeze(a, "frozen")
	iv := b.IntCast(a, ctx.Int(16), true, "narrow")
	if llvm.AsIntType(iv.Type()).Bits() != 16 {
		t.Fatalf("IntCast bits = %d", llvm.AsIntType(iv.Type()).Bits())
	}
	b.Ret(fr)

	fn2 := m.NewFunction("g", ctx.Fn(ctx.Ptr(0), []llvm.AnyType{ctx.Ptr(1)}, false))
	b.MoveToEnd(fn2.NewBlock("entry"))
	asc := b.AddrSpaceCast[llvm.PtrT](fn2.ParamAs[llvm.PtrT](0), ctx.Ptr(0), "asc")
	b.Ret(asc)

	got := m.String()
	for _, want := range []string{"freeze i32", "trunc i32", "addrspacecast ptr addrspace(1)"} {
		if !strings.Contains(got, want) {
			t.Fatalf("IR missing %q:\n%s", want, got)
		}
	}
}

// TestCastTargetCrossContext 目标类型跨 Context 时必须 panic ErrCrossContext（崩溃类地板：始终校验）。
func TestCastTargetCrossContext(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "castctx")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	b := NewBuilderAt(fn.NewBlock("entry"))
	defer b.Close()
	a := fn.ParamAs[llvm.IntT](0)

	ctx2 := llvm.NewContext()
	defer ctx2.Close()
	if err := llvm.Catch(func() { b.Trunc(a, ctx2.Int(8), "") }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("foreign trunc target should panic ErrCrossContext, got %v", err)
	}
	if err := llvm.Catch(func() { b.BitCast(a, ctx2.Int(32), "") }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("foreign bitcast target should panic ErrCrossContext, got %v", err)
	}
	if err := llvm.Catch(func() { b.IntToPtr(a, ctx2.Ptr(0), "") }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("foreign inttoptr target should panic ErrCrossContext, got %v", err)
	}
	b.Ret(a)
}
