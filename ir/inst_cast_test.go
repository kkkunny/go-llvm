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
