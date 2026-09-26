package llvm

import (
	"strings"
	"testing"
)

func TestVectorConsts(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()
	i32 := ctx.Int(32)
	vt := ctx.Vec(i32, 4)
	v := ctx.ConstVector(i32, i32.Const(1).Value, i32.Const(2).Value, i32.Const(3).Value, i32.Const(4).Value)

	// ExtractElement：取出第 2 个元素（常量）
	got := ctx.ConstExtractElement(v, i32.Const(2))
	ei, err := got.As[IntT]()
	if err != nil {
		t.Fatalf("As[IntT]: %v", err)
	}
	if iv := (IntConst{ei}).UnsignedValue(); iv != 3 {
		t.Fatalf("extract = %d, want 3", iv)
	}

	// InsertElement：把第 0 个元素换成 9
	ins := ctx.ConstInsertElement(v, i32.Const(9), i32.Const(0))
	if !ins.Alive() {
		t.Fatalf("nil insert result")
	}

	// ShuffleVector：取 v 的前两个元素重复
	mask := ctx.ConstVector(i32, i32.Const(0).Value, i32.Const(1).Value, i32.Const(0).Value, i32.Const(1).Value)
	sh := ctx.ConstShuffleVector(v, vt.Zero(), mask)
	if !sh.Alive() {
		t.Fatalf("nil shuffle result")
	}
	if out := sh.String(); !strings.Contains(out, "<i32 1, i32 2, i32 1, i32 2>") {
		t.Fatalf("shuffle = %s", out)
	}

	// ConstGEP：常量基址 + 常量下标（已含元素类型，GEP2 语义）
	base := ctx.ConstNull(ctx.Ptr(0))
	gep := ctx.ConstGEP(i32, base, false, i32.Const(1))
	if gep.IsNil() {
		t.Fatalf("nil const gep")
	}
	ib := ctx.ConstGEP(i32, base, true, i32.Const(2))
	if ib.IsNil() {
		t.Fatalf("nil const inbounds gep")
	}
}
