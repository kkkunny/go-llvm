# go-llvm

[![Go Reference](https://pkg.go.dev/badge/github.com/kkkunny/go-llvm.svg)](https://pkg.go.dev/github.com/kkkunny/go-llvm)
[![Go Report Card](https://goreportcard.com/badge/github.com/kkkunny/go-llvm)](https://goreportcard.com/report/github.com/kkkunny/go-llvm)
[![CI](https://github.com/kkkunny/go-llvm/actions/workflows/ci.yml/badge.svg)](https://github.com/kkkunny/go-llvm/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go 1.27+](https://img.shields.io/badge/Go-1.27%2B-00ADD8.svg?logo=go)](https://go.dev/)

[English](README.md) | 简体中文

基于**系统安装的 LLVM** 的 Go 绑定，提供受 [inkwell](https://github.com/TheDan64/inkwell) 启发的类型安全、高封装 API。

## 目录

- [支持的 LLVM 版本](#支持的-llvm-版本)
- [包结构](#包结构)
- [快速上手](#快速上手)
  - [类型安全与错误处理](#类型安全与错误处理)
  - [代码生成](#代码生成)
  - [JIT 执行](#jit-执行)
  - [Kaleidoscope 示例](#kaleidoscope-示例)
  - [非标准 LLVM 前缀](#非标准-llvm-前缀)
- [示例](#示例)
- [文档](#文档)
- [开发](#开发)
- [升级到新的 LLVM 版本](#升级到新的-llvm-版本)
- [许可证](#许可证)

## 支持的 LLVM 版本

| 本地 LLVM | 使用方式 |
|---|---|
| 22（最新，默认分支） | `go get github.com/kkkunny/go-llvm` |
| 21 | `go get github.com/kkkunny/go-llvm@llvm21` |

注意：

* 需要 **Go 1.27+**（API 使用了泛型方法）。
* 这些 tag 不是语义化版本，`go.mod` 会将其记录为伪版本（pseudo-version）。
* 不再提供更早的 LLVM 版本线（20 及以前）；如有需要可固定到历史提交。

## 包结构

| 包 | 职责 |
|---|---|
| `llvm` | 核心词汇表：`Kind`、`Type[T]`、`Value[T]`、常量、`Context`、错误、生命周期、Go 类型映射、`DataLayout`、`MemoryBuffer` |
| `llvm/ir` | IR 构建：`Module`、`Function`、`Block`、`Builder`、指令、`Verify`/打印/解析/bitcode |
| `llvm/target` | 目标机器与代码生成（`EmitToFile`/`Emit` 产出 OBJ/ASM） |
| `llvm/jit` | ORC LLJIT 执行引擎，`Func[F]`/`MapFunc[F]`/`MapSymbol`/`RunMain` |
| `llvm/pass` | 优化管线（`RunPasses`/`AutoOpt`） |

所有 cgo 都位于 `internal/binding`。

## 快速上手

安装受支持的 LLVM 及其开发头文件（例如 `llvm-22-dev`），然后：

```shell
go get github.com/kkkunny/go-llvm
```

```go
package main

import (
	"fmt"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
)

func main() {
	ctx := llvm.NewContext()
	defer ctx.Close()

	module := ir.NewModule(ctx, "main")
	defer module.Close()

	// Go 签名直接映射为 LLVM 函数类型
	mainFn, err := module.NewFunc[func() int32]("main")
	if err != nil {
		panic(err)
	}

	builder := ir.NewBuilder(ctx)
	defer builder.Close()
	builder.MoveToEnd(mainFn.Function().NewBlock("entry"))

	// 类型安全：Add 只接受 i32 值，ICmp 返回 i1，Select 要求分支同类型
	i32 := ctx.Int(32)
	sum := builder.Add(i32.Const(1), i32.Const(2), "sum")
	ok := builder.ICmp(llvm.IntEQ, sum, i32.Const(3), "ok")
	result := builder.Select(ok, i32.Const(0), i32.Const(1), "result")
	builder.Ret(result)

	if err := module.Verify(); err != nil {
		panic(err) // 结构化 *llvm.Error，含完整诊断
	}
	fmt.Println(module)
}
```

### 类型安全与错误处理

* 值与类型按种类参数化：`llvm.Value[llvm.IntT]`、`llvm.Type[llvm.FnT]` 等。
  种类特有的操作位于角色包装类型上（`IntType.Bits()`、`Alloca.SetAlign()`、`Phi.AddIncoming()`）。
  角色包装类型实现了 `ValueRef[T]`/`TypeRef[T]`，因此可以直接传给任何期望带类型值的位置，
  无需解包（`builder.Add(i32.Const(1), i32.Const(2), "sum")`）。
* 常量分为无符号截断（`ConstInt`）与符号扩展（`ConstSInt`）两种，并提供类型导向的语法糖：
  `i32.Const(5)`、`i32.ConstS(-1)`、`f64.Const(3.14)`。
* 动态来源（解析出的 IR、指令迭代）产出 `llvm.Value[llvm.DynT]`；用泛型方法 `As[U]()` /
  `MustAs[U]()` 恢复其种类。
* 可恢复的失败返回 `error`；程序员错误 `panic(*llvm.Error)`，`llvm.Catch` 可将其转回 error，
  `llvm.Try` 在转回 error 的同时返回值。校验分三层：
  * **崩溃类地板（任何构建都开，纯 Go）**：nil 句柄、已释放（`ErrUseAfterFree`）、跨 `Context`、
    已关闭的 Context/Module/Builder。缺少它们时 cgo 内的非法句柄是 SIGSEGV —— `recover` 接不住、
    进程直接死且没有 Go 堆栈；有了它们，所有误用都退化为可 `Catch` 的 Go panic。
  * **语义契约（仅调试构建）**：操作数类型一致、对齐为 2 的幂、索引越界、原子序合法性、
    调用实参数量/类型、结果种类。这些在 `-tags=llvm_release` 信任构建下**编译期整体消除**，
    误用交由 LLVM assert + `Module.Verify()` 兜底（与 C/Rust/inkwell 的 release 语义一致）。
  * 调试构建另含：单 goroutine 契约采样检查、链接/代码生成/JIT 边界自动 `Verify`、
    `panic` 消息附带最近操作现场、Context 关闭时的未释放资源报告、LLVM 诊断回调默认安装
    （把默认 handler 的进程退出变为日志）。
* `Context`/`Module`/`Builder` 实现 `io.Closer`。值归其 context/module 所有；
  值/类型角色方法统一通过 `Value.Ref()`/`Type.Ref()` 这一咽喉点获取句柄，崩溃类地板只在此处
  实施（无需每个方法重复 `Check`；`RawRef()` 是仅供预检层使用的未校验访问器）。
  `Module.Check`/`Block.Check` 守卫各自的句柄。每个操作都会自检生命周期，因此 use-after-free
  与重复关闭会表现为 `ErrUseAfterFree`/`ErrClosed` panic，而不会触碰悬空句柄。
  `MemoryBuffer`/`DataLayout` 另带 GC finalizer 作为安全网，忘记 `Close` 只会泄漏到下一次 GC，
  而不是永久泄漏。
* 遍历 API 同时提供切片与惰性两种形式：`Block.Insts()` / `Block.AllInsts()`、
  `Function.Blocks()` / `AllBlocks()` / `AllParams()`、`StructType.Elems()` / `AllElems()`。
  `All*` 变体返回 `iter.Seq` 且不分配内存（`for inst := range blk.AllInsts()`）。

### 代码生成

```go
target.InitNative()
native, _ := target.NativeTarget()
tm, _ := target.NewTargetMachine(native, target.DefaultTriple(), target.HostCPUName(), target.HostCPUFeatures(),
	target.OptDefault, target.RelocPIC, target.CodeModelDefault)
defer tm.Close()

tm.ApplyTo(module)                                 // 写入 triple + data layout
_ = tm.EmitToFile(module, "main.o", target.ObjectFile)
asm, _ := tm.Emit(module, target.AsmFile)          // 或产出到内存缓冲
defer asm.Close()
```

### JIT 执行

```go
target.InitNative()
j, _ := jit.NewLLJIT()
defer j.Close()

_ = j.AddProcessSymbols()                          // 进程符号（libc/libm 等）可被 JIT 内 extern 声明解析
_ = j.AddIRModule(module)                          // 模块与 Context 所有权移交 JIT
fib, _ := j.Func[func(int32) int32]("fib")         // 真实 Go 函数值
fmt.Println(fib(10))

_ = j.MapFunc("host_cb", func(x int32) int32 { return x * 2 })  // 宿主函数注册为 JIT 符号
p, _ := j.Lookup("some_symbol")                    // unsafe.Pointer 低层入口
code, _ := j.RunMain([]string{"prog"})             // 按 main(argc, argv, envp) 调用
```

`Func[F]` / `MapFunc[F]` 调用经由 `reflect` + 固定签名 C 通道，因此每次调用约 0.1–0.5 µs，
并有一次小分配（slot 数组已池化）。请获取一次 Go 函数值后长期持有。调试构建会在注册前按 JIT
模块为该符号声明的签名校验 `F`，签名不匹配时以 `ErrTypeMismatch` panic，而不是执行未定义行为；
信任构建（`-tags=llvm_release`）跳过校验、信任调用方。
热循环场景请用 `Lookup` 取出符号地址，并在自己的 cgo 绑定中调用 —— 注意裸 `unsafe.Pointer`
无法从纯 Go 调用（除非借助 cgo 或汇编跳板），因此这一逃生舱需要启用 cgo 的调用方。

`AddProcessSymbols` 在主 JITDylib 上安装 `DynamicLibrarySearchGenerator`，因此只声明函数
（如 `extern double sin(double)`）的 IR 会在具体化（materialization）时解析到宿主进程。
请在 `AddIRModule` 之前调用。

### Kaleidoscope 示例

[`examples/kaleidoscope`](examples/kaleidoscope) 是 LLVM 教程 Kaleidoscope 语言（到第 7 章：
JIT 与优化）的可运行实现，完全构建在本库公开 API 之上 —— 词法分析器、Pratt 解析器、
`ir` 代码生成、`pass` 优化管线，以及 `jit.LLJIT` + `ResourceTracker` 执行。运行方式：

```shell
go run ./examples/kaleidoscope          # REPL
go run ./examples/kaleidoscope -e '1 + 2 * 2'
go run ./examples/kaleidoscope -dl -dp -dc -e '1 + 2 * 2'   # dump tokens/AST/IR
```

### 非标准 LLVM 前缀

`internal/binding/cgo.go` 内置了常见 Linux 布局（`/usr/lib/llvm-NN`、`/usr`、`/usr/local`、
`/usr/lib64`）的 `#cgo` 标志，并使用不带版本号的 `-lLLVM`。如果你的 LLVM 安装在其他位置，
可在构建时通过环境变量覆盖：

```shell
export CGO_CFLAGS="$(llvm-config --cflags)"
export CGO_CXXFLAGS="$(llvm-config --cxxflags)"
export CGO_LDFLAGS="$(llvm-config --ldflags --libs)"
```

`CGO_CFLAGS`/`CGO_CXXFLAGS` 是全局的，因此同样会作用于 `internal/binding` 自身的编译；
在自己 main 包里放一个 `#cgo` 文件则不会。

## 示例

可运行示例位于 [`examples/`](examples)；确切的命令与期望输出见
[`examples/README.md`](examples/README.md)：

| 示例 | 描述 |
|---|---|
| [hello](examples/hello/) | 构建、校验并打印最小模块（`1 + 2`）的 IR。 |
| [jit-fib](examples/jit-fib/) | 用 IR 构建递归 `fib`，JIT 后从 Go 调用；并用 `MapFunc` 注册宿主 Go 回调。 |
| [codegen](examples/codegen/) | 用 `EmitToFile` 将模块 AOT 编译为目标文件（`.o`）与汇编文本（`.s`）。 |
| [opt](examples/opt/) | 用 `pass.AutoOpt` 运行 `default<O2>` 管线，并打印优化前后的指令数。 |
| [kaleidoscope](examples/kaleidoscope/) | 构建在 go-llvm JIT 之上的完整 Kaleidoscope 语言前端（教程第 1–7 章）。 |

## 文档

* API 参考：[pkg.go.dev/github.com/kkkunny/go-llvm](https://pkg.go.dev/github.com/kkkunny/go-llvm)
* 贡献者与架构说明：[`AGENTS.md`](AGENTS.md)

## 开发

```shell
go build ./...
go vet ./...
go test ./...                            # 调试构建（默认）：三层校验全开
go test -tags=llvm_release ./...         # 信任构建：语义契约/调试增强编译期消除
go test ./ir -run TestGolden -update     # 重新生成 golden IR（两种构建下结果必须一致）
make test test-release bench bench-release
```

CI 在每次 push 与 pull request 时通过
[`.github/workflows/ci.yml`](.github/workflows/ci.yml) 运行：`gofmt` lint 检查，加上
`ubuntu-24.04` 上的 `default`/`release` 测试矩阵，LLVM 22 来自 apt.llvm.org。

语义契约类负向测试（期望 panic 的误用测试）通过 `requireDebug(t)` 挂在调试构建，
`llvm_release` 矩阵下自动跳过；崩溃类地板测试两种构建都必须通过。

## 升级到新的 LLVM 版本

1. 阅读新版本的发行说明：
   `https://releases.llvm.org/N.1.0/docs/ReleaseNotes.html` 的
   **Changes to the C API** 一节（移除、弃用、行为变化）。
2. 对照本地头文件核查本仓库引用的每个 C 符号：
   `rg -o 'C\.[A-Za-z_]\w*' --glob '*.go'`，再用
   `grep -rw NAME /usr/include/llvm-c/` 逐一确认。
3. 先冻结上一条版本线，再在 master 上更新版本相关位置：
   `Makefile` 的 `MIN/MAX_SUPPORT_MAJOR_VERSION`、上面的支持矩阵，以及
   本 README。用 `make config MAJOR_VERSION=NN` 重新生成
   `internal/binding/cgo.go` 并审查 diff。

   ```shell
   git branch llvm-NN master     # 保持旧版本线可达
   git tag -a llvmNN -m "LLVM NN line"
   ```
4. 枚举常量无需改动：它们声明为 `C.Name`，自动绑定到本地头文件定义的值。
   未知的值种类/操作码会退化为通用回退包装类型，而不是 panic。
5. 仅在需要时绑定新增的 C API。

## 许可证

基于 Apache License 2.0 许可 —— 见 [`LICENSE`](LICENSE)。
