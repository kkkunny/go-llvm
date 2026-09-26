package ir

import (
	"strings"
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

func TestSetOperandAndUses(t *testing.T) {
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

	// use-list 读取：UsedValue 与 User 必须指向真实的两端
	uses := 0
	for u := range Uses(add) {
		if got := u.UsedValue().String(); got != add.String() {
			t.Fatalf("UsedValue = %s, want %s", got, add)
		}
		if got := u.User().String(); got != sub.String() {
			t.Fatalf("User = %s, want %s", got, sub)
		}
		uses++
	}
	if uses != 1 {
		t.Fatalf("Uses(add) = %d, want 1", uses)
	}

	// SetOperand 定向替换第 0 个操作数
	SetOperand(sub, 0, c)
	if got := OperandAt(sub, 0).String(); got != c.String() {
		t.Fatalf("sub operand 0 = %s, want %s", got, c)
	}
	uses = 0
	for range Uses(add) {
		uses++
	}
	if uses != 0 {
		t.Fatalf("Uses(add) after SetOperand = %d, want 0", uses)
	}
	if out := m.String(); !strings.Contains(out, "%y = sub i32 %1, %1") {
		t.Fatalf("IR missing rewritten operand:\n%s", out)
	}

	// 提前 break 后遍历仍可重来
	n := 0
	for range Operands(sub) {
		n++
		break
	}
	if n != 1 {
		t.Fatalf("early break yielded %d operands", n)
	}
	if got := countSeq(Operands(sub)); got != 2 {
		t.Fatalf("Operands(sub) = %d, want 2", got)
	}
	n = 0
	for range Uses(c) {
		n++
		break
	}
	if n != 1 {
		t.Fatalf("early break yielded %d uses", n)
	}

	// nil 替换值/操作数（崩溃类地板：始终校验）
	if err := llvm.Catch(func() { SetOperand(sub, 0, nil) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil operand should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { ReplaceAllUses(sub, nil) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil replacement should panic ErrInvalidArg, got %v", err)
	}
}

// TestOperandIndexPrecheck 越界操作数下标在调试构建下必须 panic。
func TestOperandIndexPrecheck(t *testing.T) {
	requireDebug(t)
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
	b.Ret(add)

	if err := llvm.Catch(func() { OperandAt(add, 5) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("operand index out of range should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { SetOperand(add, 5, c) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("SetOperand index out of range should panic ErrInvalidArg, got %v", err)
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
