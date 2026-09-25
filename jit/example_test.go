package jit_test

import (
	"fmt"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
	"github.com/kkkunny/go-llvm/jit"
	"github.com/kkkunny/go-llvm/target"
)

// newJIT 初始化宿主目标并创建 LLJIT，供示例复用
func newJIT() *jit.LLJIT {
	if err := target.InitNative(); err != nil {
		panic(err)
	}
	j, err := jit.NewLLJIT()
	if err != nil {
		panic(err)
	}
	return j
}

// ExampleLLJIT 演示最简 JIT 流程：把模块交给 LLJIT，按名字取出函数并调用。
// 模块（连同其 Context）的所有权在 AddIRModule 时移交给 JIT。
func ExampleLLJIT() {
	j := newJIT()
	defer j.Close()

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "add")
	i32 := ctx.Int(32)
	fn := m.NewFunction("add", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(b.Add(fn.ParamAs[llvm.IntT](0), fn.ParamAs[llvm.IntT](1), "sum"))
	if err := b.Close(); err != nil {
		panic(err)
	}

	// AddIRModule 后模块与 Context 归 JIT 所有，Go 侧句柄立即失效
	if err := j.AddIRModule(m); err != nil {
		panic(err)
	}

	// Func 把 JIT 符号包装成真实 Go 函数值（reflect.MakeFunc + 桥接调用）
	add, err := j.Func[func(int32, int32) int32]("add")
	if err != nil {
		panic(err)
	}
	fmt.Println(add(20, 22))
	// Output:
	// 42
}

// ExampleLLJIT_MapFunc 演示宿主回调：MapFunc 把 Go 函数注册为 JIT 符号，
// JIT 代码经自动生成的 IR 包装体调用它（native → Go 方向）。
func ExampleLLJIT_MapFunc() {
	j := newJIT()
	defer j.Close()

	// 宿主回调即 JIT 代码中该符号的实现
	hostMulAdd := func(a, b int32) int32 { return a*b + 1 }
	if err := j.MapFunc("host_mul_add", hostMulAdd); err != nil {
		panic(err)
	}

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "caller")
	i32 := ctx.Int(32)
	decl := m.NewFunction("host_mul_add", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	fn := m.NewFunction("caller", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(b.Call[llvm.IntT](
		decl.Value,
		[]llvm.AnyValue{fn.ParamAs[llvm.IntT](0).Dyn(), fn.ParamAs[llvm.IntT](1).Dyn()},
		"",
	).Value)
	if err := b.Close(); err != nil {
		panic(err)
	}
	if err := j.AddIRModule(m); err != nil {
		panic(err)
	}

	caller, err := j.Func[func(int32, int32) int32]("caller")
	if err != nil {
		panic(err)
	}
	fmt.Println(caller(6, 7)) // 6*7 + 1
	// Output:
	// 43
}

// ExampleLLJIT_RunMain 演示以 C 的 main(argc, argv, envp) 约定调用 JIT 中的
// main：RunMain 负责构造 argv 并返回 main 的退出码。
func ExampleLLJIT_RunMain() {
	j := newJIT()
	defer j.Close()

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "main")
	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	fn := m.NewFunction("main", ctx.Fn(i32, []llvm.AnyType{i32, ptr, ptr}, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(ctx.ConstInt(i32, 7))
	if err := b.Close(); err != nil {
		panic(err)
	}
	if err := j.AddIRModule(m); err != nil {
		panic(err)
	}

	code, err := j.RunMain([]string{"prog"})
	if err != nil {
		panic(err)
	}
	fmt.Println(code)
	// Output:
	// 7
}
