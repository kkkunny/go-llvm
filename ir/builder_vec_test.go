package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestBuilderVectorInsts(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "vec")
	defer m.Close()

	i32 := ctx.Int(32)
	vty := ctx.Vec(i32, 4)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{vty}, false))
	entry := fn.NewBlock("entry")

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	vec := fn.ParamAs[llvm.VecT](0)
	one := ctx.ConstInt(i32, 1).Value
	idx := ctx.ConstInt(i32, 2).Value

	ins := b.InsertElement(vec, one, idx, "ins")
	ex := b.ExtractElement[llvm.IntT](ins, idx, "ex")
	mask := ctx.ConstVector(i32,
		ctx.ConstInt(i32, 3).Value, ctx.ConstInt(i32, 2).Value,
		ctx.ConstInt(i32, 1).Value, ctx.ConstInt(i32, 0).Value)
	sh := b.ShuffleVector(ins, ins, mask, "sh")

	b.Ret(ex)
	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := m.String()
	for _, want := range []string{
		"%ins = insertelement <4 x i32> %0, i32 1, i32 2",
		"%ex = extractelement <4 x i32> %ins, i32 2",
		"shufflevector <4 x i32> %ins, <4 x i32> %ins, <4 x i32> <i32 3, i32 2, i32 1, i32 0>",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q:\n%s", want, got)
		}
	}
	_ = sh
}

func TestBuilderVectorPrecheck(t *testing.T) {
	requireDebug(t)
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "vec-pre")
	defer m.Close()

	i32 := ctx.Int(32)
	vty := ctx.Vec(i32, 2)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))
	entry := fn.NewBlock("entry")
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	zero := ctx.ConstZero(vty)
	idx := ctx.ConstInt(i32, 0).Value

	// 元素类型不符
	if err := llvm.Catch(func() {
		b.InsertElement(zero, ctx.ConstFloat(ctx.Float(llvm.FloatDouble), 1).Value, idx, "")
	}); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("elem type mismatch should panic ErrTypeMismatch, got %v", err)
	}
	// ExtractElement 种类不符
	if err := llvm.Catch(func() {
		b.ExtractElement[llvm.FloatT](zero, idx, "")
	}); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("kind mismatch should panic ErrTypeMismatch, got %v", err)
	}
	// ShuffleVector v1/v2 类型不符
	if err := llvm.Catch(func() {
		b.ShuffleVector(zero, ctx.ConstZero(ctx.Vec(ctx.Int(64), 2)), zero, "")
	}); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("shuffle type mismatch should panic ErrTypeMismatch, got %v", err)
	}
	b.RetVoid()
}
