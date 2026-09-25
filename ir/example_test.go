package ir_test

import (
	"errors"
	"fmt"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
)

// ExampleNewModule 演示模块、函数、基本块与 Builder 的最小构建流程：
// Go 函数签名直接映射为 LLVM 函数类型，指令通过 Builder 插入基本块。
func ExampleNewModule() {
	ctx := llvm.NewContext()
	defer ctx.Close()

	module := ir.NewModule(ctx, "example")
	defer module.Close()

	// func(int32, int32) int32 映射为 LLVM 函数类型 i32 (i32, i32)
	add, err := module.NewFunc[func(int32, int32) int32]("add")
	if err != nil {
		panic(err)
	}
	fn := add.Function()

	fmt.Println(fn.Name())
	fmt.Println(fn.Signature())
	fmt.Println(fn.OnlyDecl()) // 尚无基本块，仍是声明

	entry := fn.NewBlock("entry")

	b := ir.NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	sum := b.Add(fn.ParamAs[llvm.IntT](0), fn.ParamAs[llvm.IntT](1), "sum")
	b.Ret(sum)

	fmt.Println(fn.OnlyDecl())
	fmt.Println(len(entry.Insts()))
	// Output:
	// add
	// i32 (i32, i32)
	// true
	// false
	// 2
}

// ExampleBuilder 演示算术指令的构建与种类恢复：Builder 的算术方法只接受
// 整数值（编译期种类安全）；动态来源（指令遍历等）得到的是擦除种类的
// Value[DynT]，需经 As 在运行时校验种类后恢复类型参数。
func ExampleBuilder() {
	ctx := llvm.NewContext()
	defer ctx.Close()

	module := ir.NewModule(ctx, "arith")
	defer module.Close()

	fn, err := module.NewFunc[func(int32, int32) int32]("arith")
	if err != nil {
		panic(err)
	}
	f := fn.Function()
	entry := f.NewBlock("entry")

	b := ir.NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	a, c := f.ParamAs[llvm.IntT](0), f.ParamAs[llvm.IntT](1)
	sum := b.Add(a, c, "sum")      // 回绕加法
	diff := b.Sub(sum, c, "diff")  // 减法
	prod := b.Mul(diff, a, "prod") // 乘法
	b.Ret(prod)

	if err := module.Verify(); err != nil {
		panic(err)
	}

	// 动态来源的值：As 校验种类后恢复为 llvm.Value[llvm.IntT]
	first, ok := entry.FirstInst()
	if !ok {
		panic("empty entry block")
	}
	recovered, err := first.As[llvm.IntT]()
	if err != nil {
		panic(err)
	}
	fmt.Println(recovered.Name())

	fmt.Println(module)
	// Output:
	// sum
	// ; ModuleID = 'arith'
	// source_filename = "arith"
	//
	// define i32 @arith(i32 %0, i32 %1) {
	// entry:
	//   %sum = add i32 %0, %1
	//   %diff = sub i32 %sum, %1
	//   %prod = mul i32 %diff, %0
	//   ret i32 %prod
	// }
}

// ExampleModule_Verify 演示模块校验：合法模块的 Verify 返回 nil；
// 非法模块（基本块缺少终结指令）返回 Reason 为 ErrVerify 的错误而非 panic。
func ExampleModule_Verify() {
	ctx := llvm.NewContext()
	defer ctx.Close()

	b := ir.NewBuilder(ctx)
	defer b.Close()

	// 合法模块：入口块以 ret 终结
	valid := ir.NewModule(ctx, "valid")
	fn, err := valid.NewFunc[func() int32]("f")
	if err != nil {
		panic(err)
	}
	entry := fn.Function().NewBlock("entry")
	b.MoveToEnd(entry)
	b.Ret(ctx.Int(32).Const(0))
	fmt.Println("valid module:", valid.Verify() == nil)
	valid.Close()

	// 非法模块：基本块缺少终结指令
	invalid := ir.NewModule(ctx, "invalid")
	bad, err := invalid.NewFunc[func() int32]("g")
	if err != nil {
		panic(err)
	}
	bad.Function().NewBlock("entry")
	defer invalid.Close()

	err = invalid.Verify()
	fmt.Println("Verify error:", err != nil)
	var verr *llvm.Error
	if !errors.As(err, &verr) {
		panic(err)
	}
	fmt.Println("error is ErrVerify:", verr.Reason == llvm.ErrVerify)
	fmt.Println("error op:", verr.Op)
	// Output:
	// valid module: true
	// Verify error: true
	// error is ErrVerify: true
	// error op: ir.Module.Verify
}

// ExampleParseIRString 演示从文本 IR 解析模块：解析出的模块与手工构建的
// 模块等价，可直接查询函数、基本块与指令。
func ExampleParseIRString() {
	ctx := llvm.NewContext()
	defer ctx.Close()

	const src = `define i32 @square(i32 %x) {
entry:
  %r = mul i32 %x, %x
  ret i32 %r
}
`
	module, err := ir.ParseIRString(ctx, src)
	if err != nil {
		panic(err)
	}
	defer module.Close()

	fn, ok := module.GetFunction("square")
	if !ok {
		panic("function square not found")
	}
	fmt.Println(fn.Name())
	fmt.Println(fn.Signature())
	fmt.Println(len(fn.Blocks()))
	fmt.Println(len(fn.Blocks()[0].Insts()))
	// Output:
	// square
	// i32 (i32)
	// 1
	// 2
}

// ExampleFunction_AllBlocks 演示惰性遍历：AllBlocks/AllInsts 均返回
// iter.Seq，range 友好且无切片分配，适合在大模块上统计与遍历。
func ExampleFunction_AllBlocks() {
	ctx := llvm.NewContext()
	defer ctx.Close()

	module := ir.NewModule(ctx, "blocks")
	defer module.Close()

	fn, err := module.NewFunc[func(int32) int32]("abs")
	if err != nil {
		panic(err)
	}
	f := fn.Function()
	entry := f.NewBlock("entry")
	then := f.NewBlock("then")
	exit := f.NewBlock("exit")

	b := ir.NewBuilder(ctx)
	defer b.Close()

	b.MoveToEnd(entry)
	neg := b.ICmp(llvm.IntSLT, f.ParamAs[llvm.IntT](0), ctx.Int(32).Const(0), "neg")
	b.CondBr(neg, then, exit)

	b.MoveToEnd(then)
	b.Ret(b.Sub(ctx.Int(32).Const(0), f.ParamAs[llvm.IntT](0), "sub"))

	b.MoveToEnd(exit)
	b.Ret(f.ParamAs[llvm.IntT](0))

	if err := module.Verify(); err != nil {
		panic(err)
	}

	blocks, insts := 0, 0
	for blk := range f.AllBlocks() {
		blocks++
		n := 0
		for range blk.AllInsts() {
			n++
			insts++
		}
		fmt.Printf("%s: %d\n", blk.Name(), n)
	}
	fmt.Printf("blocks=%d insts=%d\n", blocks, insts)
	// Output:
	// entry: 2
	// then: 2
	// exit: 1
	// blocks=3 insts=5
}
