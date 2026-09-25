// Command hello 演示 go-llvm 的最小工作流：创建 Context/Module、构建函数、
// Verify 后打印模块 IR。
//
// 运行：
//
//	go run ./examples/hello
package main

import (
	"fmt"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
)

func main() {
	// Context 是所有 LLVM 实体的所有权根，Close 时级联释放登记过的子资源
	ctx := llvm.NewContext()
	defer ctx.Close()

	// Module 承载一个编译单元；名字会写入 IR 的 ModuleID 与 source_filename
	module := ir.NewModule(ctx, "main")
	defer module.Close()

	// Go 签名直接映射为 LLVM 函数类型：int32 main()
	mainFunc, err := module.NewFunc[func() int32]("main")
	if err != nil {
		panic(err)
	}

	// 构建器把指令插入指定基本块
	b := ir.NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(mainFunc.Function().NewBlock("entry"))

	// return 1 + 2；操作数都是常量，LLVM 会在构建期折叠，最终 IR 为 ret i32 3
	i32 := ctx.Int(32)
	sum := b.Add(i32.Const(1), i32.Const(2), "sum")
	b.Ret(sum)

	// Verify 校验 IR 合法性；失败返回 *llvm.Error（可用 llvm.Catch 收敛为 error）
	if err := module.Verify(); err != nil {
		panic(err)
	}

	// Module 实现了 fmt.Stringer，打印完整 IR 文本
	fmt.Println(module)
}
