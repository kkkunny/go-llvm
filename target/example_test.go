package target_test

import (
	"fmt"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
	"github.com/kkkunny/go-llvm/target"
)

// ExampleInitNative 演示初始化宿主目标并查询其能力：InitNative 注册原生目标的
// info/目标/MC/汇编打印/解析/反汇编组件，之后即可通过 NativeTarget 取得宿主目标。
func ExampleInitNative() {
	// 初始化宿主目标；失败返回 ErrCodeGen
	if err := target.InitNative(); err != nil {
		panic(err)
	}

	native, err := target.NativeTarget()
	if err != nil {
		panic(err)
	}

	// 宿主目标通常具备 JIT、目标机器与汇编后端
	fmt.Println(native.HasJIT())
	fmt.Println(native.HasTargetMachine())
	fmt.Println(native.HasAsmBackend())
	// Output:
	// true
	// true
	// true
}

// ExampleTargetMachine_Emit 演示宿主目标的按需代码生成：创建 TargetMachine，
// 用 ApplyTo 把三元组与数据布局写入模块，再把模块编译为汇编文本的内存缓冲
// （只检查产出非空，不打印汇编内容，避免跨环境差异）。
func ExampleTargetMachine_Emit() {
	if err := target.InitNative(); err != nil {
		panic(err)
	}
	native, err := target.NativeTarget()
	if err != nil {
		panic(err)
	}

	ctx := llvm.NewContext()
	defer ctx.Close()

	module := ir.NewModule(ctx, "emit")
	defer module.Close()

	// i32 add(i32, i32)
	i32 := ctx.Int(32)
	fn := module.NewFunction("add", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(b.Add(fn.ParamAs[llvm.IntT](0), fn.ParamAs[llvm.IntT](1), "sum"))
	if err := b.Close(); err != nil {
		panic(err)
	}
	if err := module.Verify(); err != nil {
		panic(err)
	}

	// 目标机器是独立所有权根（不经 Context.Own），用毕 Close；
	// 三元组/CPU/特性串取宿主值，优化级别与重定位模式按需选择
	tm, err := target.NewTargetMachine(
		native,
		target.DefaultTriple(),
		target.HostCPUName(),
		target.HostCPUFeatures(),
		target.OptNone,
		target.RelocPIC,
		target.CodeModelDefault,
	)
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	// ApplyTo 把目标三元组与数据布局写入模块，Emit 产出汇编内存缓冲
	tm.ApplyTo(module)
	asm, err := tm.Emit(module, target.AsmFile)
	if err != nil {
		panic(err)
	}
	defer asm.Close()

	fmt.Println(asm.Len() > 0)
	// Output:
	// true
}
