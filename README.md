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
  back to an error and `llvm.Try` converts while also returning a value. All builder calls are pre-checked
  (positioned builder, same context, live handles, matching operand types, result kinds) before reaching cgo.
* `Context`/`Module`/`Builder` implement `io.Closer`. Values are owned by their context/module;
  every operation (builder calls plus module/block/instruction role methods) self-checks lifetime,
  so use-after-free and double-close surface as `ErrUseAfterFree`/`ErrClosed` panics instead of
  touching dangling handles. `MemoryBuffer`/`DataLayout` additionally carry a GC finalizer as a
  safety net, so a forgotten `Close` leaks until the next GC instead of forever.
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

or generate an override file in **your own** main package (never in this
repository's root):

```shell
curl -O https://raw.githubusercontent.com/kkkunny/go-llvm/master/Makefile
make config             # or pin the toolchain: make config MAJOR_VERSION=22
```

## Development

```shell
go build ./...
go vet ./...
go test ./...
go test ./ir -run TestGolden -update   # 重新生成 golden IR
```

## Updating for a new LLVM release

1. Read the new release notes at
   `https://releases.llvm.org/N.1.0/docs/ReleaseNotes.html`, section
   **Changes to the C API** (removals, deprecations, behavior changes).
2. Cross-check every C symbol this repo references against the local headers:
   `rg -o 'C\.[A-Za-z_]\w*' --glob '*.go'`, then verify each name with
   `grep -rw NAME /usr/include/llvm-c/`.
3. Freeze the previous line first, then bump the four version spots on master:
   `internal/binding/cgo.go` candidate dirs, `Makefile`
   `MIN/MAX_SUPPORT_MAJOR_VERSION`, the support table above, and this README.

   ```shell
   git branch llvm-NN master     # keep the old line reachable
   git tag -a llvmNN -m "LLVM NN line"
   ```
4. Enum constants need no changes: they are declared as `C.Name` and bind to
   whatever the local headers define. Unknown value kinds / opcodes degrade to
   generic fallback wrappers instead of panicking.
5. Bind newly added C APIs only when needed.
