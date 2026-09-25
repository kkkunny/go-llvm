package ir

import (
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestAsRoles(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "t")
	defer m.Close()

	i32 := ctx.Int(32)
	callee := m.NewFunction("callee", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()

	a := fn.ParamAs[llvm.IntT](0)
	slot := b.Alloca(i32, "slot")
	b.Store(a, slot)
	ld := b.Load[llvm.IntT](slot, i32, "v")
	sum := b.Add(ld, ld, "sum")
	call := b.Call[llvm.IntT](callee, []llvm.AnyValue{sum}, "c")
	phi := b.PHI[llvm.IntT](i32, "p")
	b.Ret(call)

	if _, ok := AsLoad[llvm.IntT](ld); !ok {
		t.Fatalf("load -> AsLoad failed")
	}
	if _, ok := AsStore(ld); ok {
		t.Fatalf("load should not convert to store")
	}
	if ldInst, ok := AsLoad[llvm.IntT](ld); !ok || ldInst.IsVolatile() {
		t.Fatalf("AsLoad result unexpected")
	}
	if _, ok := AsAlloca(slot); !ok {
		t.Fatalf("alloca -> AsAlloca failed")
	}
	if c, ok := AsCall[llvm.IntT](call); !ok || c.ArgCount() != 1 {
		t.Fatalf("call -> AsCall failed")
	}
	if _, ok := AsLoad[llvm.IntT](call); ok {
		t.Fatalf("call should not convert to load")
	}
	if _, ok := AsPhi[llvm.IntT](phi); !ok {
		t.Fatalf("phi -> AsPhi failed")
	}
	if _, ok := AsCall[llvm.IntT](sum); ok {
		t.Fatalf("add should not convert to call")
	}
}

// TestAsRoleConversions 覆盖其余 As* 角色转换的正向路径与拒绝路径。
func TestAsRoleConversions(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "roles")
	defer m.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	pers := m.NewFunction("pers", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	g := m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	h := m.NewFunction("h", ctx.Fn(ctx.Void(), nil, false))
	ti := m.NewGlobal("ti", ptr)

	b := NewBuilder(ctx)
	defer b.Close()

	// ---- invoke / switch / landingpad ----
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn.SetPersonality(pers)
	entry := fn.NewBlock("entry")
	cont := fn.NewBlock("cont")
	lpad := fn.NewBlock("lpad")
	def := fn.NewBlock("def")
	case1 := fn.NewBlock("case1")

	b.MoveToEnd(entry)
	iv := b.Invoke[llvm.IntT](g, []llvm.AnyValue{fn.ParamAs[llvm.IntT](0)}, cont, lpad, "iv")
	b.MoveToEnd(cont)
	sw := b.Switch(fn.ParamAs[llvm.IntT](0), def)
	sw.AddCase(ctx.ConstInt(i32, 1), case1)

	b.MoveToEnd(lpad)
	lp := b.LandingPad(ctx.Struct([]llvm.AnyType{ptr, i32}, false), "lp")
	lp.SetCleanup(true)
	b.Resume(lp)

	b.MoveToEnd(def)
	b.Ret(ctx.ConstInt(i32, 0).Value)
	b.MoveToEnd(case1)
	b.Ret(ctx.ConstInt(i32, 1).Value)

	if got, ok := AsInvoke[llvm.IntT](iv); !ok || got.ArgCount() != 1 || got.NormalBlock() != cont {
		t.Fatalf("invoke -> AsInvoke failed: %v %v", got, ok)
	}
	if got, ok := AsSwitch(sw); !ok || got.Count() != 1 {
		t.Fatalf("switch -> AsSwitch failed: %v %v", got, ok)
	}
	if got, ok := AsLandingPad[llvm.StructT](lp); !ok || got.ClauseCount() != 0 || !got.IsCleanup() {
		t.Fatalf("landingpad -> AsLandingPad failed: %v %v", got, ok)
	}
	if _, ok := AsInvoke[llvm.IntT](sw); ok {
		t.Fatal("switch should not convert to invoke")
	}

	// ---- atomicrmw / cmpxchg / fence / store ----
	af := m.NewFunction("atomics", ctx.Fn(ctx.Void(), []llvm.AnyType{ptr}, false))
	aentry := af.NewBlock("entry")
	b.MoveToEnd(aentry)
	p := af.ParamAs[llvm.PtrT](0)
	one := ctx.ConstInt(i32, 1).Value
	rm := b.AtomicRMW(llvm.RMWAdd, p, one, llvm.AtomicMonotonic, false, "rm")
	cx := b.CmpXchg(p, one, ctx.ConstInt(i32, 2).Value, llvm.AtomicAcquire, llvm.AtomicMonotonic, false, "cx")
	fe := b.Fence(llvm.AtomicAcquire, false)
	st := b.Store(one, p)
	b.RetVoid()

	if got, ok := AsAtomicRMW[llvm.IntT](rm); !ok || got.Op() != llvm.RMWAdd {
		t.Fatalf("atomicrmw -> AsAtomicRMW failed: %v %v", got, ok)
	}
	if got, ok := AsCmpXchg(cx); !ok || got.SuccessOrdering() != llvm.AtomicAcquire {
		t.Fatalf("cmpxchg -> AsCmpXchg failed: %v %v", got, ok)
	}
	if got, ok := AsFence(fe); !ok || got.Ordering() != llvm.AtomicAcquire {
		t.Fatalf("fence -> AsFence failed: %v %v", got, ok)
	}
	if got, ok := AsStore(st); !ok || got.IsVolatile() {
		t.Fatalf("store -> AsStore failed: %v %v", got, ok)
	}
	if _, ok := AsFence(rm); ok {
		t.Fatal("atomicrmw should not convert to fence")
	}

	// ---- catchswitch / funclet pad ----
	ef := m.NewFunction("funclets", ctx.Fn(ctx.Void(), nil, false))
	ef.SetPersonality(pers)
	e2 := ef.NewBlock("entry")
	c2 := ef.NewBlock("cont")
	disp := ef.NewBlock("dispatch")
	hnd := ef.NewBlock("hnd")
	done := ef.NewBlock("done")

	b.MoveToEnd(e2)
	b.Invoke[llvm.VoidT](h, nil, c2, disp, "")
	b.MoveToEnd(c2)
	b.RetVoid()
	b.MoveToEnd(disp)
	cs := b.CatchSwitch(nil, Block{}, "cs")
	cs.AddHandler(hnd)
	b.MoveToEnd(hnd)
	cp := b.CatchPad(cs, []llvm.AnyValue{ti}, "cp")
	b.CatchRet(cp, done)
	b.MoveToEnd(done)
	b.RetVoid()

	if got, ok := AsCatchSwitch(cs); !ok || got.HandlerCount() != 1 {
		t.Fatalf("catchswitch -> AsCatchSwitch failed: %v %v", got, ok)
	}
	if got, ok := AsFuncletPad(cp); !ok || got.ArgCount() != 1 || got.ParentCatchSwitch().Name() != "cs" {
		t.Fatalf("catchpad -> AsFuncletPad failed: %v %v", got, ok)
	}
	if _, ok := AsFuncletPad(cs); ok {
		t.Fatal("catchswitch should not convert to funclet pad")
	}
	if _, ok := AsCatchSwitch(cp); ok {
		t.Fatal("catchpad should not convert to catchswitch")
	}

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
}

// TestTerminatorPrecheck 非条件终结指令/越界后继在调试构建下必须 panic。
func TestTerminatorPrecheck(t *testing.T) {
	requireDebug(t)
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "term")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{i32}, false))
	entry := fn.NewBlock("entry")
	then := fn.NewBlock("then")
	els := fn.NewBlock("else")
	b := NewBuilderAt(entry)
	defer b.Close()

	cond := b.ICmp(llvm.IntEQ, fn.ParamAs[llvm.IntT](0), ctx.ConstInt(i32, 0), "c")
	br := b.Br(then) // 无条件 br：只有 1 个后继
	if IsConditional(br) {
		t.Fatal("unconditional br should not be conditional")
	}
	if n := SuccessorCount(br); n != 1 {
		t.Fatalf("successor count = %d, want 1", n)
	}
	// 注：SuccessorCount 对非终结指令是 LLVM-C 未定义行为（实测挂起），不做断言。
	if err := llvm.Catch(func() { Condition(br) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("Condition on unconditional br should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { SetCondition(br, cond) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("SetCondition on unconditional br should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { Successor(br, 1) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("successor out of range should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { SetSuccessor(br, 1, then) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("SetSuccessor out of range should panic ErrInvalidArg, got %v", err)
	}
	b.MoveToEnd(then)
	b.RetVoid()
	b.MoveToEnd(els)
	b.RetVoid()
}

func TestTerminatorOps(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "t")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{i32}, false))
	entry := fn.NewBlock("entry")
	then := fn.NewBlock("then")
	els := fn.NewBlock("else")
	b := NewBuilderAt(entry)
	defer b.Close()

	cond := b.ICmp(llvm.IntNE, fn.ParamAs[llvm.IntT](0), ctx.ConstInt(i32, 0), "c")
	br := b.CondBr(cond, then, els)

	if !IsConditional(br) {
		t.Fatalf("cond br should be conditional")
	}
	if SuccessorCount(br) != 2 {
		t.Fatalf("successor count = %d, want 2", SuccessorCount(br))
	}
	SetSuccessor(br, 1, then)
	if got := Successor(br, 1); got.Name() != "then" {
		t.Fatalf("successor 1 = %s, want then", got.Name())
	}
	one := ctx.ConstInt(i32, 1)
	SetCondition(br, one)
	if got := Condition(br); got.String() != one.String() {
		t.Fatalf("condition = %s, want %s", got, one)
	}
	if b2 := b.Br(then); IsConditional(b2) {
		t.Fatalf("unconditional br should not be conditional")
	}
}
