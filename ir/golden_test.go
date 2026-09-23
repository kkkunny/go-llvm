package ir

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/kkkunny/go-llvm"
)

var updateGolden = flag.Bool("update", false, "更新 golden 文件")

func checkGolden(t *testing.T, got, name string) {
	t.Helper()
	goldenPath := filepath.Join("testdata", "golden", name)
	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(goldenPath, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden (可用 -update 生成): %v", err)
	}
	if got != string(want) {
		t.Fatalf("module output differs from golden %s:\n--- got ---\n%s\n--- want ---\n%s", name, got, want)
	}
}

// TestGoldenMain 端到端构建 README 示例模块并与 golden 文件比对
func TestGoldenMain(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "main")
	defer m.Close()

	i32 := ctx.Int(32)
	main, err := m.NewFunc[func() int32]("main")
	if err != nil {
		t.Fatal(err)
	}
	entry := main.Function().NewBlock("entry")

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	// int32 result = add(1, 2); return result == 3 ? 0 : 1
	one := ctx.ConstInt(i32, 1).Value
	two := ctx.ConstInt(i32, 2).Value
	sum := b.Add(one, two, "sum")
	isThree := b.ICmp(llvm.IntEQ, sum, ctx.ConstInt(i32, 3).Value, "is_three")
	ret := b.Select(isThree, ctx.ConstInt(i32, 0).Value, ctx.ConstInt(i32, 1).Value, "ret")
	b.Ret(ret)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}

	checkGolden(t, m.String(), "main.ll")
}

// TestGoldenLoop 覆盖 PHI/Switch/内存指令的端到端场景
func TestGoldenLoop(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "loop")
	defer m.Close()

	i32 := ctx.Int(32)
	i64 := ctx.Int(64)
	boolTy := ctx.Bool()
	fn, err := m.NewFunc[func(int32) int32]("sum_to")
	if err != nil {
		t.Fatal(err)
	}
	entry := fn.Function().NewBlock("entry")
	loop := fn.Function().NewBlock("loop")
	exit := fn.Function().NewBlock("exit")

	b := NewBuilder(ctx)
	defer b.Close()

	b.MoveToEnd(entry)
	acc := b.Alloca(i32, "acc")
	b.Store(ctx.ConstInt(i32, 0).Value, acc)
	b.Br(loop)

	b.MoveToEnd(loop)
	i := b.PHI(i32, "i")
	cur := b.Load(acc, i32, "cur")
	i.AddIncoming(Incoming[llvm.IntT]{Value: ctx.ConstInt(i32, 0).Value, Block: entry})
	next := b.Add(i, ctx.ConstInt(i32, 1).Value, "next")
	b.Store(next, acc)
	done := b.ICmp(llvm.IntSGE, next, fn.Function().ParamAs[llvm.IntT](0), "done")
	b.CondBr(done, exit, loop)
	i.AddIncoming(Incoming[llvm.IntT]{Value: next, Block: loop})
	_ = cur
	_ = boolTy

	b.MoveToEnd(exit)
	result := b.Load(acc, i32, "result")
	b.Ret(result)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}

	checkGolden(t, m.String(), "loop.ll")
	_ = i64
}

// TestGoldenEH 覆盖 Itanium 式（invoke/landingpad/resume）与 funclet 式（catchswitch/catchpad/catchret/cleanuppad/cleanupret）
func TestGoldenEH(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "eh")
	defer m.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	pers := m.NewFunction("pers", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	g := m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	h := m.NewFunction("h", ctx.Fn(ctx.Void(), nil, false))
	ti := m.NewGlobal("ti", ptr)

	b := NewBuilder(ctx)
	defer b.Close()

	// ---- itanium ----
	f1 := m.NewFunction("itanium", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	f1.SetPersonality(pers)
	e1 := f1.NewBlock("entry")
	c1 := f1.NewBlock("cont")
	l1 := f1.NewBlock("lpad")

	b.MoveToEnd(e1)
	v := b.Invoke[llvm.IntT](g, []llvm.AnyValue{f1.ParamAs[llvm.IntT](0)}, c1, l1, "v")
	b.MoveToEnd(c1)
	b.Ret(v)
	b.MoveToEnd(l1)
	lp := b.LandingPad(ctx.Struct([]llvm.AnyType{ptr, i32}, false), "lp")
	lp.AddClause(ti)
	lp.SetCleanup(true)
	b.Resume(lp)

	// ---- funclets ----
	f2 := m.NewFunction("funclets", ctx.Fn(ctx.Void(), nil, false))
	f2.SetPersonality(pers)
	e2 := f2.NewBlock("entry")
	c2 := f2.NewBlock("cont")
	d2 := f2.NewBlock("dispatch")
	hn := f2.NewBlock("hnd")
	dn := f2.NewBlock("done")

	b.MoveToEnd(e2)
	b.Invoke[llvm.VoidT](h, nil, c2, d2, "")
	b.MoveToEnd(c2)
	b.RetVoid()
	b.MoveToEnd(d2)
	cs := b.CatchSwitch(nil, Block{}, "cs")
	cs.AddHandler(hn)
	b.MoveToEnd(hn)
	cp := b.CatchPad(cs, []llvm.AnyValue{ti}, "cp")
	b.CatchRet(cp, dn)
	b.MoveToEnd(dn)
	b.RetVoid()

	// ---- cleanup ----
	f3 := m.NewFunction("cleanup", ctx.Fn(ctx.Void(), nil, false))
	f3.SetPersonality(pers)
	e3 := f3.NewBlock("entry")
	c3 := f3.NewBlock("cont")
	k3 := f3.NewBlock("clean")

	b.MoveToEnd(e3)
	b.Invoke[llvm.VoidT](h, nil, c3, k3, "")
	b.MoveToEnd(c3)
	b.RetVoid()
	b.MoveToEnd(k3)
	clp := b.CleanupPad(nil, nil, "clp")
	b.CleanupRet(clp, Block{})

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	checkGolden(t, m.String(), "eh.ll")
}
