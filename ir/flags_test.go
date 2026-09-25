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
