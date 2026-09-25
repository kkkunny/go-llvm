# go-llvm

Go bindings to a **system-installed LLVM**, with a type-safe, highly-encapsulated API
inspired by [inkwell](https://github.com/TheDan64/inkwell).

Currently supported:

| Local LLVM | How to use |
|---|---|
| 22 (latest, default branch) | `go get github.com/kkkunny/go-llvm` |
| 21 | `go get github.com/kkkunny/go-llvm@llvm21` |

Notes:

* Requires **Go 1.27+** (the API uses generic methods).
* These are non-semver tags, so `go.mod` records them as pseudo-versions.
* Older LLVM lines (20 and earlier) are no longer provided; they can be pinned to
  historic commits if needed.

## Packages

| Package | Responsibility |
|---|---|
| `llvm` | Core vocabulary: `Kind`, `Type[T]`, `Value[T]`, constants, `Context`, errors, lifetime, Go type mapping, `DataLayout`, `MemoryBuffer` |
| `llvm/ir` | IR construction: `Module`, `Function`, `Block`, `Builder`, instructions, `Verify`/print/parse/bitcode |
| `llvm/target` | Target machines and code generation (`EmitToFile`/`Emit` for OBJ/ASM) |
| `llvm/jit` | ORC LLJIT execution engine, `Func[F]`/`MapFunc[F]`/`MapSymbol`/`RunMain` |
| `llvm/pass` | Optimization pipelines (`RunPasses`/`AutoOpt`) |

All cgo lives in `internal/binding`.

## Usage

Install a supported LLVM together with its development headers (for example
`llvm-22-dev`), then:

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

### Type safety and errors

* Values and types are kind-parameterized: `llvm.Value[llvm.IntT]`, `llvm.Type[llvm.FnT]`, etc.
  Category-specific operations live on role wrappers (`IntType.Bits()`, `Alloca.SetAlign()`, `Phi.AddIncoming()`).
  Role wrappers implement `ValueRef[T]`/`TypeRef[T]`, so they can be passed directly wherever a typed value is
  expected — no unwrapping needed (`builder.Add(i32.Const(1), i32.Const(2), "sum")`).
* Constants come in unsigned-truncating (`ConstInt`) and sign-extending (`ConstSInt`) flavors, with
  type-directed sugar: `i32.Const(5)`, `i32.ConstS(-1)`, `f64.Const(3.14)`.
* Dynamic sources (parsed IR, instruction iteration) yield `llvm.Value[llvm.DynT]`; recover the kind with
  the generic method `As[U]()` / `MustAs[U]()`.
* Recoverable failures return `error`; programmer errors `panic(*llvm.Error)`, which `llvm.Catch` converts
  back to an error and `llvm.Try` converts while also returning a value. Checks are layered:
  * **崩溃类地板（任何构建都开，纯 Go）**：nil 句柄、已释放（`ErrUseAfterFree`）、跨 `Context`、
    已关闭的 Context/Module/Builder。缺少它们时 cgo 内的非法句柄是 SIGSEGV —— `recover` 接不住、
    进程直接死且没有 Go 堆栈；有了它们，所有误用都退化为可 `Catch` 的 Go panic。
  * **语义契约（仅调试构建）**：操作数类型一致、对齐为 2 的幂、索引越界、原子序合法性、
    调用实参数量/类型、结果种类。这些在 `-tags=llvm_release` 信任构建下**编译期整体消除**，
    误用交由 LLVM assert + `Module.Verify()` 兜底（与 C/Rust/inkwell 的 release 语义一致）。
  * 调试构建另含：单 goroutine 契约采样检查、链接/代码生成/JIT 边界自动 `Verify`、
    `panic` 消息附带最近操作现场、Context 关闭时的未释放资源报告、LLVM 诊断回调默认安装
    （把默认 handler 的进程退出变为日志）。
* `Context`/`Module`/`Builder` implement `io.Closer`. Values are owned by their context/module;
  value/type role methods deliver handles through the `Value.Ref()`/`Type.Ref()` choke point, which
  enforces the crash-class floor in one place (no per-method `Check` duplication; `RawRef()` is the
  unchecked accessor used only by the pre-check layer). `Module.Check`/`Block.Check` guard their own
  handles. Every operation self-checks lifetime, so use-after-free and double-close surface as
  `ErrUseAfterFree`/`ErrClosed` panics instead of touching dangling handles. `MemoryBuffer`/`DataLayout`
  additionally carry a GC finalizer as a safety net, so a forgotten `Close` leaks until the next GC
  instead of forever.
* Traversal APIs come in slice and lazy flavors: `Block.Insts()` / `Block.AllInsts()`,
  `Function.Blocks()` / `AllBlocks()` / `AllParams()`, `StructType.Elems()` / `AllElems()`.
  The `All*` variants are `iter.Seq` and allocate nothing (`for inst := range blk.AllInsts()`).

### Code generation

```go
target.InitNative()
native, _ := target.NativeTarget()
tm, _ := target.NewTargetMachine(native, target.DefaultTriple(), target.HostCPUName(), target.HostCPUFeatures(),
	target.OptDefault, target.RelocPIC, target.CodeModelDefault)
defer tm.Close()

tm.ApplyTo(module)                                   // 写入 triple + data layout
_ = tm.EmitToFile(module, "main.o", target.ObjectFile)
asm, _ := tm.Emit(module, target.AsmFile)          // 或产出到内存缓冲
defer asm.Close()
```

### JIT execution

```go
target.InitNative()
j, _ := jit.NewLLJIT()
defer j.Close()

_ = j.AddIRModule(module)                          // 模块与 Context 所有权移交 JIT
fib, _ := j.Func[func(int32) int32]("fib")         // 真实 Go 函数值
fmt.Println(fib(10))

_ = j.MapFunc("host_cb", func(x int32) int32 { return x * 2 })  // 宿主函数注册为 JIT 符号
p, _ := j.Lookup("some_symbol")                    // unsafe.Pointer 低层入口
code, _ := j.RunMain([]string{"prog"})             // 按 main(argc, argv, envp) 调用
```

`Func[F]` / `MapFunc[F]` calls go through `reflect` + a fixed-signature C channel, so each call costs
roughly 0.1–0.5 µs and one small allocation (slot arrays are pooled). Fetch the Go function value once
and keep it. For hot loops, look the symbol up with `Lookup` and call the address from your own cgo
binding — note that a bare `unsafe.Pointer` cannot be called from pure Go without cgo or an
assembly trampoline, so the escape hatch requires a cgo-enabled caller.

### Non-standard LLVM prefixes

`internal/binding/cgo.go` ships `#cgo` flags for the common Linux layouts
(`/usr/lib/llvm-NN`, `/usr`, `/usr/local`, `/usr/lib64`) with an unversioned
`-lLLVM`. If your LLVM lives somewhere else, override via environment
variables in your build:

```shell
export CGO_CFLAGS="$(llvm-config --cflags)"
export CGO_CXXFLAGS="$(llvm-config --cxxflags)"
export CGO_LDFLAGS="$(llvm-config --ldflags --libs)"
```

`CGO_CFLAGS`/`CGO_CXXFLAGS` are global, so they also apply to `internal/binding`'s
own compilation; a `#cgo` file in your own main package would not.

## Development

```shell
go build ./...
go vet ./...
go test ./...                            # 调试构建（默认）：三层校验全开
go test -tags=llvm_release ./...         # 信任构建：语义契约/调试增强编译期消除
go test ./ir -run TestGolden -update     # 重新生成 golden IR（两种构建下结果必须一致）
make test test-release bench bench-release
```

语义契约类负向测试（期望 panic 的误用测试）通过 `requireDebug(t)` 挂在调试构建，
`llvm_release` 矩阵下自动跳过；崩溃类地板测试两种构建都必须通过。

## Updating for a new LLVM release

1. Read the new release notes at
   `https://releases.llvm.org/N.1.0/docs/ReleaseNotes.html`, section
   **Changes to the C API** (removals, deprecations, behavior changes).
2. Cross-check every C symbol this repo references against the local headers:
   `rg -o 'C\.[A-Za-z_]\w*' --glob '*.go'`, then verify each name with
   `grep -rw NAME /usr/include/llvm-c/`.
3. Freeze the previous line first, then bump the version spots on master:
   `Makefile` `MIN/MAX_SUPPORT_MAJOR_VERSION`, the support table above, and
   this README. Regenerate `internal/binding/cgo.go` with
   `make config MAJOR_VERSION=NN` and review the diff.

   ```shell
   git branch llvm-NN master     # keep the old line reachable
   git tag -a llvmNN -m "LLVM NN line"
   ```
4. Enum constants need no changes: they are declared as `C.Name` and bind to
   whatever the local headers define. Unknown value kinds / opcodes degrade to
   generic fallback wrappers instead of panicking.
5. Bind newly added C APIs only when needed.
