package ir

import (
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestOperandsAndReplace(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "t")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()

	a := fn.ParamAs[llvm.IntT](0)
	c := fn.ParamAs[llvm.IntT](1)
	add := b.Add(a, c, "x")
	sub := b.Sub(add, c, "y")
	b.Ret(sub)

	if n := OperandCount(add); n != 2 {
		t.Fatalf("OperandCount = %d, want 2", n)
	}
	if op0 := OperandAt(add, 0); op0.String() != a.String() {
		t.Fatalf("operand 0 = %v, want %v", op0, a)
	}

	n := 0
	for range Operands(sub) {
		n++
	}
	if n != 2 {
		t.Fatalf("Operands count = %d, want 2", n)
	}

	uses := 0
	for range Uses(add) {
		uses++
	}
	if uses != 1 {
		t.Fatalf("Uses(add) = %d, want 1", uses)
	}

	ReplaceAllUses(add, a)
	if got := OperandAt(sub, 0).String(); got != a.String() {
		t.Fatalf("sub operand 0 = %s, want %s", got, a.String())
	}
	uses = 0
	for range Uses(add) {
		uses++
	}
	if uses != 0 {
		t.Fatalf("Uses(add) after RAUW = %d, want 0", uses)
	}
}

func TestOperandAPIsAfterModuleClose(t *testing.T) {
	requireDebug(t)
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "t")
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()
	one := ctx.ConstInt(ctx.Int(32), 1)
	add := b.Add(one, one, "x")
	b.RetVoid()
	_ = m.Close()

	defer func() {
		if recover() == nil {
			t.Fatalf("OperandCount on dead value should panic")
		}
	}()
	OperandCount(add)
}
