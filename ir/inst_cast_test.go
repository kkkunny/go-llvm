package ir

import (
	"fmt"
	"strings"
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

// TestTerminatorPrecheck 后继越界在调试构建下必须 panic；
// 非终结指令与无条件终结指令的常开拒绝见 TestSuccessorOpsRejectNonTerminator 与 TestSwitchConditionAccess。
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

	b.ICmp(llvm.IntEQ, fn.ParamAs[llvm.IntT](0), ctx.ConstInt(i32, 0), "c")
	br := b.Br(then) // 无条件 br：只有 1 个后继
	if IsConditional(br) {
		t.Fatal("unconditional br should not be conditional")
	}
	if n := SuccessorCount(br); n != 1 {
		t.Fatalf("successor count = %d, want 1", n)
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

// TestSuccessorOpsRejectNonTerminator 非终结指令传给后继访问器必须在两种构建模式下
// panic ErrInvalidArg，而不是把 UB 交给 LLVM-C（LLVMGetNumSuccessors 对非终结指令
// 实测挂起，LLVMGetCondition 可能 SIGSEGV）。
func TestSuccessorOpsRejectNonTerminator(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "termguard")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{i32}, false))
	entry := fn.NewBlock("entry")
	exit := fn.NewBlock("exit")
	b := NewBuilderAt(entry)
	defer b.Close()

	icmp := b.ICmp(llvm.IntEQ, fn.ParamAs[llvm.IntT](0), ctx.ConstInt(i32, 0), "c")
	if err := llvm.Catch(func() { SuccessorCount(icmp) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("SuccessorCount on non-terminator should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { Successor(icmp, 0) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("Successor on non-terminator should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { SetSuccessor(icmp, 0, exit) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("SetSuccessor on non-terminator should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { IsConditional(icmp) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("IsConditional on non-terminator should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { Condition(icmp) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("Condition on non-terminator should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { SetCondition(icmp, ctx.ConstInt(i32, 1)) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("SetCondition on non-terminator should panic ErrInvalidArg, got %v", err)
	}

	// 终结指令路径保持正确：br 有 1 个后继、ret 有 0 个后继
	br := b.Br(exit)
	if IsConditional(br) {
		t.Fatal("unconditional br should not be conditional")
	}
	if n := SuccessorCount(br); n != 1 {
		t.Fatalf("br successor count = %d, want 1", n)
	}
	b.MoveToEnd(exit)
	ret := b.RetVoid()
	if n := SuccessorCount(ret); n != 0 {
		t.Fatalf("ret successor count = %d, want 0", n)
	}
	if IsConditional(ret) {
		t.Fatal("ret should not be conditional")
	}
}

// TestSwitchConditionAccess 覆盖 switch 的条件访问：llvm-c/Core.h 明示
// LLVMIsConditional/LLVMGetCondition/LLVMSetCondition 只支持 BranchInst，
// switch 必须按操作码分流到操作数 0（LLVMGetOperand/LLVMSetOperand）。
// 本用例不设 requireDebug，两种构建模式都验证。
func TestSwitchConditionAccess(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "switchcond")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{i32}, false))
	entry := fn.NewBlock("entry")
	b := NewBuilderAt(entry)
	defer b.Close()

	cond := fn.ParamAs[llvm.IntT](0)

	// 0 / 1 / 2 case 的 switch 条件都必须可读可替换
	for n := 0; n <= 2; n++ {
		head := fn.NewBlock(fmt.Sprintf("head%d", n))
		def := fn.NewBlock(fmt.Sprintf("def%d", n))
		b.MoveToEnd(head)
		sw := b.Switch(cond, def)

		b.MoveToEnd(def)
		b.RetVoid()
		for i := 0; i < n; i++ {
			cb := fn.NewBlock(fmt.Sprintf("case%d_%d", n, i))
			sw.AddCase(ctx.ConstInt(i32, uint64(i+1)), cb)
			b.MoveToEnd(cb)
			b.RetVoid()
		}

		if !IsConditional(sw) {
			t.Errorf("%d case switch should be conditional", n)
		}
		if got := Condition(sw); got.String() != cond.String() {
			t.Errorf("%d case switch condition = %s, want %s", n, got, cond)
		}
		want := ctx.ConstInt(i32, uint64(10+n))
		SetCondition(sw, want)
		if got := Condition(sw); got.String() != want.String() {
			t.Errorf("%d case switch condition after SetCondition = %s, want %s", n, got, want)
		}
	}

	// 无条件 br 与 ret：Condition/SetCondition 在两种构建模式下都必须 panic ErrInvalidArg
	b.MoveToEnd(entry)
	b.RetVoid()
	brHead := fn.NewBlock("brhead")
	brDest := fn.NewBlock("brdest")
	b.MoveToEnd(brHead)
	br := b.Br(brDest)
	b.MoveToEnd(brDest)
	ret := b.RetVoid()
	for _, tc := range []struct {
		name string
		term llvm.AnyValue
	}{
		{"unconditional br", br},
		{"ret", ret},
	} {
		if err := llvm.Catch(func() { Condition(tc.term) }); err == nil || err.Reason != llvm.ErrInvalidArg {
			t.Errorf("%s: Condition should panic ErrInvalidArg, got %v", tc.name, err)
		}
		if err := llvm.Catch(func() { SetCondition(tc.term, cond) }); err == nil || err.Reason != llvm.ErrInvalidArg {
			t.Errorf("%s: SetCondition should panic ErrInvalidArg, got %v", tc.name, err)
		}
	}
	if IsConditional(br) {
		t.Error("unconditional br should not be conditional")
	}
	if IsConditional(ret) {
		t.Error("ret should not be conditional")
	}

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	out := m.String()
	for n := 0; n <= 2; n++ {
		want := fmt.Sprintf("switch i32 %d, label %%def%d", 10+n, n)
		if !strings.Contains(out, want) {
			t.Errorf("module output missing %q:\n%s", want, out)
		}
	}
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
