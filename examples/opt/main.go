// Command opt 演示 go-llvm 的模块优化：构造一个含冗余指令的函数，
// 运行 default<O2> 默认管线，打印优化前后的指令数。
//
// 运行：
//
//	go run ./examples/opt
package main

import (
	"fmt"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
	"github.com/kkkunny/go-llvm/pass"
)

func main() {
	ctx := llvm.NewContext()
	defer ctx.Close()

	module := ir.NewModule(ctx, "opt")
	defer module.Close()

	// i32 compute()：常量先写进 alloca 再读回，最后套上 +0 与 *1 恒等运算
	i32 := ctx.Int(32)
	compute := module.NewFunction("compute", ctx.Fn(i32, nil, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(compute.NewBlock("entry"))

	slot := b.Alloca(i32, "slot")
	b.Store(i32.Const(1), slot) // 先写 1，随后被覆盖：死存储
	b.Store(i32.Const(2), slot)

	x := b.Load(slot, i32, "x").Value
	a := b.Add(x, i32.Const(0), "a") // x + 0
	v := b.Mul(a, i32.Const(1), "v") // a * 1
	b.Ret(v)

	if err := b.Close(); err != nil {
		panic(err)
	}
	if err := module.Verify(); err != nil {
		panic(err)
	}

	// AutoOpt 运行 default<O*> 默认管线：mem2reg 提升 alloca，
	// instcombine 折叠恒等运算并删除死存储
	before := countInsts(module)
	if err := pass.AutoOpt(module, pass.O2); err != nil {
		panic(err)
	}
	after := countInsts(module)
	fmt.Printf("insts: %d -> %d\n", before, after)
}

// countInsts 统计模块内全部函数定义的指令数（声明没有基本块，计 0）
func countInsts(m *ir.Module) int {
	n := 0
	for fn := range m.AllFunctions() {
		for blk := range fn.AllBlocks() {
			for range blk.AllInsts() {
				n++
			}
		}
	}
	return n
}
