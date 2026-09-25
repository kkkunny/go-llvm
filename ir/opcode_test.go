package ir

import (
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestOpOf(t *testing.T) {
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

	// 非指令值（函数参数）：返回 false
	if _, ok := OpOf(a); ok {
		t.Fatalf("param should not be an instruction")
	}

	got := map[string]bool{}
	for inst := range blk.AllInsts() {
		op, ok := OpOf(inst)
		if !ok {
			t.Fatalf("inst %v is not recognized as instruction", inst)
		}
		got[op.String()] = true
	}
	if !got["add"] || !got["ret"] {
		t.Fatalf("opcodes = %v, want add & ret", got)
	}
}
