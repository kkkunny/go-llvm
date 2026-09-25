package ir

import (
	"iter"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func countSeq[T any](seq iter.Seq[T]) int {
	n := 0
	for range seq {
		n++
	}
	return n
}

func TestIterLazyTraversal(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "iter")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	blk := fn.NewBlock("entry")
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(blk)
	x, y := fn.ParamAs[llvm.IntT](0), fn.ParamAs[llvm.IntT](1)
	for i := 0; i < 3; i++ {
		b.Add(x, y, "")
	}
	b.Ret(x)

	if got, want := countSeq(blk.AllInsts()), len(blk.Insts()); got != want || got != 4 {
		t.Fatalf("AllInsts = %d, want %d (4)", got, want)
	}
	if got, want := countSeq(fn.AllBlocks()), len(fn.Blocks()); got != want || got != 1 {
		t.Fatalf("AllBlocks = %d, want %d (1)", got, want)
	}
	if got, want := countSeq(fn.AllParams()), len(fn.Params()); got != want || got != 2 {
		t.Fatalf("AllParams = %d, want %d (2)", got, want)
	}

	// 提前 break 后块仍可正常使用
	n := 0
	for range blk.AllInsts() {
		n++
		break
	}
	if n != 1 {
		t.Fatalf("early break yielded %d items", n)
	}
	if got := countSeq(blk.AllInsts()); got != 4 {
		t.Fatalf("traversal after early break = %d", got)
	}

	// AllBlocks/AllParams 提前 break
	if n := countSeq(fn.AllBlocks()); n != 1 {
		t.Fatalf("AllBlocks = %d", n)
	}
	n = 0
	for range fn.AllBlocks() {
		n++
		break
	}
	if n != 1 {
		t.Fatalf("AllBlocks early break yielded %d", n)
	}
	n = 0
	for range fn.AllParams() {
		n++
		break
	}
	if n != 1 {
		t.Fatalf("AllParams early break yielded %d", n)
	}

	// 模块级惰性遍历提前 break
	m.NewFunction("g", ctx.Fn(i32, nil, false))
	m.NewGlobal("gv", i32)
	n = 0
	for range m.AllFunctions() {
		n++
		break
	}
	if n != 1 {
		t.Fatalf("AllFunctions early break yielded %d", n)
	}
	n = 0
	for range m.AllGlobals() {
		n++
		break
	}
	if n != 1 {
		t.Fatalf("AllGlobals early break yielded %d", n)
	}
}
