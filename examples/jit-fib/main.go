// Command jit-fib 演示 go-llvm 的 JIT 能力：用 IR 构建递归 fib 函数，
// 经 Func 包装为 Go 函数值调用，并用 MapFunc 把宿主 Go 回调注册为 JIT 符号。
//
// 运行：
//
//	go run ./examples/jit-fib
package main

import (
	"fmt"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
	"github.com/kkkunny/go-llvm/jit"
	"github.com/kkkunny/go-llvm/target"
)

func main() {
	// JIT 使用宿主目标；未初始化时 NewLLJIT 探测不到目标机器
	if err := target.InitNative(); err != nil {
		panic(err)
	}
	j, err := jit.NewLLJIT()
	if err != nil {
		panic(err)
	}
	defer j.Close()

	// 把宿主 Go 函数注册为 JIT 符号 host_double：JIT 代码调用它会经自动
	// 生成的 IR 包装体回到 Go（native → Go 方向）
	hostDouble := func(x int32) int32 { return x * 2 }
	if err := j.MapFunc("host_double", hostDouble); err != nil {
		panic(err)
	}

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "fib")
	i32 := ctx.Int(32)
	fnTy := ctx.Fn(i32, []llvm.AnyType{i32}, false)

	// i32 fib(i32 n)：纯 SSA 的递归定义
	fib := m.NewFunction("fib", fnTy)
	hostDoubleFn := m.NewFunction("host_double", fnTy)

	b := ir.NewBuilder(ctx)
	entry := fib.NewBlock("entry")
	base := fib.NewBlock("base")
	recurse := fib.NewBlock("recurse")

	// n < 2 时跳 base，否则跳 recurse
	b.MoveToEnd(entry)
	n := fib.ParamAs[llvm.IntT](0)
	b.CondBr(b.ICmp(llvm.IntSLT, n, i32.Const(2), "cmp"), base, recurse)

	// base: return n
	b.MoveToEnd(base)
	b.Ret(n)

	// recurse: return fib(n-1) + fib(n-2)
	b.MoveToEnd(recurse)
	n1 := b.Sub(n, i32.Const(1), "n1")
	f1 := b.Call[llvm.IntT](fib, []llvm.AnyValue{n1.Dyn()}, "f1").Value
	n2 := b.Sub(n, i32.Const(2), "n2")
	f2 := b.Call[llvm.IntT](fib, []llvm.AnyValue{n2.Dyn()}, "f2").Value
	b.Ret(b.Add(f1, f2, "sum"))

	// i32 double_fib(i32 n)：JIT 内先调 fib，再调宿主回调 host_double
	doubleFib := m.NewFunction("double_fib", fnTy)
	b.MoveToEnd(doubleFib.NewBlock("entry"))
	arg := doubleFib.ParamAs[llvm.IntT](0)
	fibN := b.Call[llvm.IntT](fib, []llvm.AnyValue{arg.Dyn()}, "fib_n").Value
	doubled := b.Call[llvm.IntT](hostDoubleFn, []llvm.AnyValue{fibN.Dyn()}, "doubled").Value
	b.Ret(doubled)

	// 构建器先于模块释放；AddIRModule 之后模块与 Context 归 JIT 所有
	if err := b.Close(); err != nil {
		panic(err)
	}
	if err := m.Verify(); err != nil {
		panic(err)
	}
	if err := j.AddIRModule(m); err != nil {
		panic(err)
	}

	// Func 把 JIT 符号包装成真实 Go 函数值（按名字解析 + 桥接调用）
	fibFn, err := j.Func[func(int32) int32]("fib")
	if err != nil {
		panic(err)
	}
	doubleFibFn, err := j.Func[func(int32) int32]("double_fib")
	if err != nil {
		panic(err)
	}

	fmt.Printf("fib(10) = %d\n", fibFn(10))
	fmt.Printf("double_fib(10) = host_double(fib(10)) = %d\n", doubleFibFn(10))
}
