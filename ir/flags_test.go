package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestFastMathFlags(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "t")
	defer m.Close()

	f64 := ctx.Float(llvm.FloatDouble)
	fn := m.NewFunction("f", ctx.Fn(f64, []llvm.AnyType{f64}, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()

	a := fn.ParamAs[llvm.FloatT](0)
	fadd := b.FAdd(a, a, "x")

	if !CanFastMath(fadd) {
		t.Fatalf("fadd should accept fast-math flags")
	}
	if got := FastMathOf(fadd); got != FastMathNone {
		t.Fatalf("initial flags = %v, want none", got)
	}
	SetFastMath(fadd, FastMathNoNaNs|FastMathAllowReciprocal)
	if got := FastMathOf(fadd); got&FastMathNoNaNs == 0 || got&FastMathAllowReciprocal == 0 {
		t.Fatalf("flags = %v", got)
	}
	b.Ret(fadd)

	out := m.String()
	for _, want := range []string{"fadd", "nnan", "arcp"} {
		if !strings.Contains(out, want) {
			t.Fatalf("IR missing %q:\n%s", want, out)
		}
	}
}

func TestGEPNoWrap(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "t")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Ptr(0), []llvm.AnyType{ctx.Ptr(0)}, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()

	p := fn.ParamAs[llvm.PtrT](0)
	g := b.GEPWithFlags(i32, p, []llvm.ValueRef[llvm.IntT]{ctx.ConstInt(i32, 1)}, NoWrapInBounds|NoWrapNUW, "g")
	if got := GEPNoWrapOf(g); got&NoWrapInBounds == 0 || got&NoWrapNUW == 0 {
		t.Fatalf("flags = %v", got)
	}
	SetGEPNoWrap(g, NoWrapInBounds)
	// inbounds 蕴含 nusw，LLVM 会规范化 flags
	if got := GEPNoWrapOf(g); got != NoWrapInBounds|NoWrapNUSW {
		t.Fatalf("flags after set = %v, want inbounds|nusw", got)
	}
	b.Ret(g)
}

func TestInstFlagsAndTailCall(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "t")
	defer m.Close()

	i32 := ctx.Int(32)
	callee := m.NewFunction("callee", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	caller := m.NewFunction("caller", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	blk := caller.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()

	a := caller.ParamAs[llvm.IntT](0)
	add := b.Add(a, a, "x")
	SetNSW(add, true)
	SetNUW(add, true)
	sd := b.SDiv(a, a, "q")
	SetExact(sd, true)
	ext := b.ZExt(a, ctx.Int(64), "e")
	SetNNeg(ext, true)

	call := b.Call[llvm.IntT](callee, []llvm.AnyValue{a}, "c")
	call.SetTailCallKind(TailCallMust)
	if !call.IsTailCall() {
		t.Fatalf("musttail should imply tail")
	}
	call.SetParamAlign(0, 16)
	b.Ret(call)

	out := m.String()
	for _, want := range []string{"nuw", "nsw", "exact", "nneg", "musttail", "align 16"} {
		if !strings.Contains(out, want) {
			t.Fatalf("IR missing %q:\n%s", want, out)
		}
	}
}
