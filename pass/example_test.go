package pass_test

import (
	"fmt"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
	"github.com/kkkunny/go-llvm/pass"
)

// ExampleAutoOpt 演示按优化级别运行默认管线：O2 会把 f() { return 1 + 2 }
// 常量折叠为 ret i32 3，函数只剩 1 条指令。
func ExampleAutoOpt() {
	ctx := llvm.NewContext()
	defer ctx.Close()

	m := ir.NewModule(ctx, "opt")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, nil, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	sum := b.Add(ctx.ConstInt(i32, 1), ctx.ConstInt(i32, 2), "sum")
	b.Ret(sum)
	if err := b.Close(); err != nil {
		panic(err)
	}

	// default<O2>：标准优化级别；可选 O0/O1/O2/O3/Os/Oz
	if err := pass.AutoOpt(m, pass.O2); err != nil {
		panic(err)
	}

	insts := 0
	for blk := range fn.AllBlocks() {
		for range blk.AllInsts() {
			insts++
		}
	}
	fmt.Printf("ok, insts=%d\n", insts)
	// Output:
	// ok, insts=1
}

// ExampleRunPassesOnFunction 演示在单个函数上运行管线（opt -passes 语法）：
// instcombine 做指令合并，把 g(a) { return a + 0 } 化简为直接返回参数。
func ExampleRunPassesOnFunction() {
	ctx := llvm.NewContext()
	defer ctx.Close()

	m := ir.NewModule(ctx, "opt-fn")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	sum := b.Add(fn.ParamAs[llvm.IntT](0), ctx.ConstInt(i32, 0), "sum")
	b.Ret(sum)
	if err := b.Close(); err != nil {
		panic(err)
	}

	if err := pass.RunPassesOnFunction(fn, "instcombine"); err != nil {
		panic(err)
	}

	insts := 0
	for blk := range fn.AllBlocks() {
		for range blk.AllInsts() {
			insts++
		}
	}
	fmt.Printf("ok, insts=%d\n", insts)
	// Output:
	// ok, insts=1
}
