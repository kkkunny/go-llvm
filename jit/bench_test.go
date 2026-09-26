package jit

import (
	"testing"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
)

func BenchmarkFuncCall(b *testing.B) {
	j := newJIT(b)
	defer j.Close()

	ctx, m := retModule(b, "bench_func", 42)
	if err := j.AddIRModule(m); err != nil {
		b.Fatal(err)
	}
	_ = ctx
	f, err := j.Func[func() int32]("answer")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if f() != 42 {
			b.Fatal("bad result")
		}
	}
}

func BenchmarkMapFuncCallback(b *testing.B) {
	j := newJIT(b)
	defer j.Close()

	if err := j.MapFunc("host_cb", func(x int64) int64 { return x + 1 }); err != nil {
		b.Fatal(err)
	}

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "bench_cb")
	i64 := ctx.Int(64)
	decl := m.NewFunction("host_cb", ctx.Fn(i64, []llvm.AnyType{i64}, false))
	fn := m.NewFunction("call_cb", ctx.Fn(i64, []llvm.AnyType{i64}, false))
	bld := ir.NewBuilder(ctx)
	bld.MoveToEnd(fn.NewBlock("entry"))
	bld.Ret(bld.Call[llvm.IntT](decl.Value, []llvm.AnyValue{fn.Param(0)}, "").Value)
	if err := bld.Close(); err != nil {
		b.Fatal(err)
	}
	if err := m.Verify(); err != nil {
		b.Fatal(err)
	}
	if err := j.AddIRModule(m); err != nil {
		b.Fatal(err)
	}

	f, err := j.Func[func(int64) int64]("call_cb")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if f(int64(i)) != int64(i)+1 {
			b.Fatal("bad result")
		}
	}
}
