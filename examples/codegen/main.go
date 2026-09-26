// Command codegen 演示 go-llvm 的 AOT 代码生成：把内存中的 IR 模块编译为
// 宿主平台的汇编文本（.s）与目标文件（.o），写入随机临时目录。
//
// 运行：
//
//	go run ./examples/codegen
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
	"github.com/kkkunny/go-llvm/target"
)

func main() {
	// 初始化宿主目标：注册原生目标的 info/target/MC/汇编打印等组件
	if err := target.InitNative(); err != nil {
		panic(err)
	}
	native, err := target.NativeTarget()
	if err != nil {
		panic(err)
	}

	ctx := llvm.NewContext()
	defer ctx.Close()

	module := ir.NewModule(ctx, "add")
	defer module.Close()

	// i32 add(i32 a, i32 b) { return a + b }
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
		target.OptDefault,
		target.RelocPIC,
		target.CodeModelDefault,
	)
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	// ApplyTo 把目标三元组与数据布局写入模块，是 EmitToFile 的前置步骤
	tm.ApplyTo(module)

	// 输出到随机临时目录，退出时整目录删除
	dir, err := os.MkdirTemp("", "go-llvm-codegen-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	for _, out := range []struct {
		kind string
		name string
		ft   target.FileType
	}{
		{"object", "add.o", target.ObjectFile},
		{"assembly", "add.s", target.AsmFile},
	} {
		path := filepath.Join(dir, out.name)
		if err := tm.EmitToFile(module, path, out.ft); err != nil {
			panic(err)
		}
		info, err := os.Stat(path)
		if err != nil {
			panic(err)
		}
		// 空文件说明代码生成没有产出，按异常处理
		if info.Size() == 0 {
			panic("empty output: " + path)
		}
		fmt.Printf("%s: %s (%d bytes)\n", out.kind, path, info.Size())
	}
}
