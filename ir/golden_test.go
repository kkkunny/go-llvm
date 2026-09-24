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

// TestGoldenAtomic 覆盖 fence/atomicrmw/cmpxchg 与 load/store 原子形态
func TestGoldenAtomic(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "atomic")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("atomics", ctx.Fn(ctx.Void(), []llvm.AnyType{ctx.Ptr(0)}, false))
	entry := fn.NewBlock("entry")

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	p := fn.ParamAs[llvm.PtrT](0)
	one := ctx.ConstInt(i32, 1).Value

	b.Fence(llvm.AtomicSequentiallyConsistent, false)
	rm := b.AtomicRMW(llvm.RMWAdd, p, one, llvm.AtomicMonotonic, false, "old")
	rm.SetAlign(4)
	b.Store(one, p).SetOrdering(llvm.AtomicRelease)
	ld := b.Load(p, i32, "cur")
	ld.SetOrdering(llvm.AtomicAcquire)
	b.CmpXchg(p, ld, ctx.ConstInt(i32, 2).Value, llvm.AtomicAcquire, llvm.AtomicMonotonic, true, "cx")
	b.RetVoid()

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	checkGolden(t, m.String(), "atomic.ll")
}

// TestGoldenVec 覆盖 insertelement/extractelement/shufflevector 与 ConstVector
func TestGoldenVec(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "vec")
	defer m.Close()

	i32 := ctx.Int(32)
	vty := ctx.Vec(i32, 4)
	fn := m.NewFunction("vec_demo", ctx.Fn(i32, []llvm.AnyType{vty}, false))
	entry := fn.NewBlock("entry")

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	// 以参数为基向量，避免全常量操作数被 IRBuilder 常量折叠
	vec := fn.ParamAs[llvm.VecT](0)
	ins := b.InsertElement(vec, ctx.ConstInt(i32, 42).Value, ctx.ConstInt(i32, 1).Value, "ins")
	mask := ctx.ConstVector(i32,
		ctx.ConstInt(i32, 0).Value, ctx.ConstInt(i32, 0).Value,
		ctx.ConstInt(i32, 0).Value, ctx.ConstInt(i32, 0).Value)
	sh := b.ShuffleVector(ins, ins, mask, "sh")
	ex := b.ExtractElement[llvm.IntT](sh, ctx.ConstInt(i32, 0).Value, "ex")
	b.Ret(ex)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	checkGolden(t, m.String(), "vec.ll")
}

// TestGoldenAttrsMeta 覆盖属性/调用约定、元数据/模块 flag、comdat 与 ctor 的端到端打印
func TestGoldenAttrsMeta(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "attrs_meta")
	defer m.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)

	g := m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{ptr}, false))
	g.AddAttr(llvm.AttrFunction, ctx.EnumAttr(llvm.AttrNoInline, 0))
	g.AddAttr(llvm.AttrFunction, ctx.StringAttr("my-attr", "v1"))
	g.AddAttr(llvm.AttrParam(0), ctx.AlignAttr(8))
	g.AddAttr(llvm.AttrReturn, ctx.EnumAttr(llvm.AttrNoUndef, 0))
	g.SetCallConv(llvm.CallConvFast)

	ctor := m.NewFunction("ctor", ctx.Fn(ctx.Void(), nil, false))
	entry := ctor.NewBlock("entry")
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)
	v := b.Call[llvm.IntT](g, []llvm.AnyValue{ctx.ConstNull(ptr)}, "v")
	AttachMetadata(v, "my.kind", ctx.MDNode(ctx.MDString("tag")))
	b.RetVoid()

	m.AppendCtor(ctor, 65535)
	m.AddModuleFlag(llvm.ModuleFlagOverride, "my.flag", ctx.MDString("v"))
	m.AddNamedMetadataOperand("my.md", ctx.MDNode(ctx.MDString("x")))
	c := m.GetOrInsertComdat("mycomdat")
	c.SetSelectionKind(llvm.ComdatNoDeduplicate)
	glob := m.NewGlobal("glob", i32)
	glob.SetInitializer(ctx.ConstInt(i32, 0).Value)
	glob.SetComdat(c)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	checkGolden(t, m.String(), "attrs_meta.ll")
}
