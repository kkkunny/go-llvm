package pass_test

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
	"github.com/kkkunny/go-llvm/pass"
)

// buildAddModule 构造 f() { return 1 + 2 }
func buildAddModule(t *testing.T, name string) (*llvm.Context, *ir.Module, ir.Function) {
	t.Helper()
	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, name)
	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, nil, false))
	entry := fn.NewBlock("entry")
	b := ir.NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)
	sum := b.Add(ctx.ConstInt(i32, 1).Value, ctx.ConstInt(i32, 2).Value, "sum")
	b.Ret(sum)
	return ctx, m, fn
}

func TestAutoOpt(t *testing.T) {
	ctx, m, _ := buildAddModule(t, "opt")
	defer ctx.Close()
	defer m.Close()

	if err := pass.AutoOpt(m, pass.O2); err != nil {
		t.Fatalf("AutoOpt: %v", err)
	}
	if got := m.String(); !strings.Contains(got, "ret i32 3") {
		t.Fatalf("expected constant fold after O2:\n%s", got)
	}
}

func TestRunPassesOnFunction(t *testing.T) {
	ctx, m, fn := buildAddModule(t, "opt-fn")
	defer ctx.Close()
	defer m.Close()

	if err := pass.RunPassesOnFunction(fn, "instcombine"); err != nil {
		t.Fatalf("RunPassesOnFunction: %v", err)
	}
	if got := m.String(); !strings.Contains(got, "ret i32 3") {
		t.Fatalf("expected constant fold after instcombine:\n%s", got)
	}
}

func TestRunPassesErrors(t *testing.T) {
	ctx, m, _ := buildAddModule(t, "opt-err")
	defer ctx.Close()
	defer m.Close()

	if err := pass.RunPasses(m, "definitely-not-a-pipeline"); err == nil {
		t.Fatal("unknown pipeline should fail")
	} else if e, ok := err.(*llvm.Error); !ok || e.Reason != llvm.ErrPass {
		t.Fatalf("unknown pipeline should return ErrPass, got %v", err)
	}
	if err := pass.AutoOpt(m, pass.Level("O9")); err == nil {
		t.Fatal("unknown level should fail")
	} else if e, ok := err.(*llvm.Error); !ok || e.Reason != llvm.ErrPass {
		t.Fatalf("unknown level should return ErrPass, got %v", err)
	}
}

func TestRunPassesOptions(t *testing.T) {
	ctx, m, _ := buildAddModule(t, "opt-opts")
	defer ctx.Close()
	defer m.Close()

	err := pass.RunPasses(m, "default<O1>",
		pass.VerifyEach(true),
		pass.DebugLogging(false),
		pass.LoopVectorization(true),
		pass.SLPVectorization(true),
		pass.LoopUnrolling(false),
		pass.LoopInterleaving(false),
		pass.ForgetAllSCEVInLoopUnroll(true),
		pass.InlinerThreshold(10),
	)
	if err != nil {
		t.Fatalf("RunPasses with options: %v", err)
	}
}
