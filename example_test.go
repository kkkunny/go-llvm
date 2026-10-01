package llvm_test

import (
	"errors"
	"fmt"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
)

// Example 包总览：建立 Context 与 Module，用 ir 子包构建一个简单函数，
// Verify 后打印模块 IR。
func Example() {
	ctx := llvm.NewContext()
	defer ctx.Close()

	module := ir.NewModule(ctx, "main")
	defer module.Close()

	// Go 签名直接映射为 LLVM 函数类型
	mainFn, err := module.NewFunc[func() int32]("main")
	if err != nil {
		panic(err)
	}

	b := ir.NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(mainFn.Function().NewBlock("entry"))

	// int32 result = add(1, 2); return result == 3 ? 0 : 1
	i32 := ctx.Int(32)
	sum := b.Add(i32.Const(1), i32.Const(2), "sum")
	ok := b.ICmp(llvm.IntEQ, sum, i32.Const(3), "ok")
	result := b.Select(ok, i32.Const(0), i32.Const(1), "result")
	b.Ret(result)

	// Verify 校验 IR 合法性；失败返回 *llvm.Error
	if err := module.Verify(); err != nil {
		panic(err)
	}
	fmt.Println(module)
	// Output:
	// ; ModuleID = 'main'
	// source_filename = "main"
	//
	// define i32 @main() {
	// entry:
	//   ret i32 0
	// }
}

// demoResource 实现 io.Closer，用于演示 Context.Own 登记的子资源。
type demoResource struct{ closed bool }

// Close 标记资源已关闭
func (r *demoResource) Close() error {
	r.closed = true
	return nil
}

// ExampleContext 演示 Context 的生命周期：Own 登记子资源，Close 时级联关闭；
// Alive 在关闭前为 true，关闭后为 false。
func ExampleContext() {
	ctx := llvm.NewContext()

	res := &demoResource{}
	ctx.Own(res) // 登记子资源：Context.Close 时按逆序级联关闭

	fmt.Println(ctx.Alive())

	ctx.Close()
	fmt.Println(ctx.Alive())
	fmt.Println(res.closed)
	// Output:
	// true
	// false
	// true
}

// ExampleTypeOf 演示 Go 类型到 LLVM 类型的映射。
func ExampleTypeOf() {
	ctx := llvm.NewContext()
	defer ctx.Close()

	// 标量：int32 → i32
	i32, err := llvm.TypeOf[int32](ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(i32)

	// 结构体：各字段类型递归映射为字面量结构体
	type point struct {
		X int32
		Y float64
	}
	p, err := llvm.TypeOf[point](ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(p)

	// 不支持的类型（切片）返回 ErrUnsupported 而不是 panic
	_, err = llvm.TypeOf[[]int32](ctx)
	unsupported := (*llvm.Error)(nil)
	fmt.Println(errors.As(err, &unsupported) && unsupported.Reason == llvm.ErrUnsupported)
	// Output:
	// i32
	// { i32, double }
	// true
}

// ExampleFnSignatureOf 演示 Go 函数签名到 LLVM 函数类型的映射。
func ExampleFnSignatureOf() {
	ctx := llvm.NewContext()
	defer ctx.Close()

	sig, _, err := llvm.FnSignatureOf[func(int32, float64) int32](ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(sig)
	fmt.Println(len(sig.Params()))
	fmt.Println(sig.Return())
	fmt.Println(sig.IsVarArg())
	// Output:
	// i32 (i32, double)
	// 2
	// i32
	// false
}
