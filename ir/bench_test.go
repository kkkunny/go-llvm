package ir

import (
	"testing"

	"github.com/kkkunny/go-llvm"
)

// 说明：LLVM IRBuilder 会对全常量操作数做常量折叠，基准必须用函数参数等非常量操作数，
// 否则测到的是折叠路径而不是指令构建路径。

func BenchmarkBuilderBinop(b *testing.B) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "bench_binop")
	defer m.Close()
	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	bld := NewBuilder(ctx)
	defer bld.Close()
	bld.MoveToEnd(fn.NewBlock("entry"))
	x, y := fn.ParamAs[llvm.IntT](0), fn.ParamAs[llvm.IntT](1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bld.Add(x, y, "")
	}
}

func BenchmarkBuilderCall(b *testing.B) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "bench_call")
	defer m.Close()
	i32 := ctx.Int(32)
	params := []llvm.AnyType{i32, i32, i32, i32}
	decl := m.NewFunction("callee", ctx.Fn(i32, params, false))
	fn := m.NewFunction("caller", ctx.Fn(i32, params, false))
	bld := NewBuilder(ctx)
	defer bld.Close()
	bld.MoveToEnd(fn.NewBlock("entry"))
	args := []llvm.AnyValue{fn.Param(0), fn.Param(1), fn.Param(2), fn.Param(3)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bld.Call[llvm.IntT](decl.Value, args, "")
	}
}

func BenchmarkBlockAllInsts(b *testing.B) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "bench_allinsts")
	defer m.Close()
	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	blk := fn.NewBlock("entry")
	bld := NewBuilder(ctx)
	defer bld.Close()
	bld.MoveToEnd(blk)
	x, y := fn.ParamAs[llvm.IntT](0), fn.ParamAs[llvm.IntT](1)
	for i := 0; i < 64; i++ {
		bld.Add(x, y, "")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n := 0
		for range blk.AllInsts() {
			n++
		}
		if n != 64 {
			b.Fatalf("unexpected instruction count: %d", n)
		}
	}
}

func BenchmarkBlockInsts(b *testing.B) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "bench_insts")
	defer m.Close()
	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	blk := fn.NewBlock("entry")
	bld := NewBuilder(ctx)
	defer bld.Close()
	bld.MoveToEnd(blk)
	x, y := fn.ParamAs[llvm.IntT](0), fn.ParamAs[llvm.IntT](1)
	for i := 0; i < 64; i++ {
		bld.Add(x, y, "")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if n := len(blk.Insts()); n != 64 {
			b.Fatalf("unexpected instruction count: %d", n)
		}
	}
}
