# go-llvm

[![Go Reference](https://pkg.go.dev/badge/github.com/kkkunny/go-llvm.svg)](https://pkg.go.dev/github.com/kkkunny/go-llvm)
[![Go Report Card](https://goreportcard.com/badge/github.com/kkkunny/go-llvm)](https://goreportcard.com/report/github.com/kkkunny/go-llvm)
[![CI](https://github.com/kkkunny/go-llvm/actions/workflows/ci.yml/badge.svg)](https://github.com/kkkunny/go-llvm/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go 1.27+](https://img.shields.io/badge/Go-1.27%2B-00ADD8.svg?logo=go)](https://go.dev/)

English | [简体中文](README.zh-CN.md)

Go bindings to a **system-installed LLVM**, with a type-safe, highly-encapsulated API
inspired by [inkwell](https://github.com/TheDan64/inkwell).

## Table of contents

- [Supported LLVM versions](#supported-llvm-versions)
- [Packages](#packages)
- [Usage](#usage)
  - [Type safety and errors](#type-safety-and-errors)
  - [Code generation](#code-generation)
  - [JIT execution](#jit-execution)
  - [Kaleidoscope example](#kaleidoscope-example)
  - [Non-standard LLVM prefixes](#non-standard-llvm-prefixes)
- [Examples](#examples)
- [Documentation](#documentation)
- [Development](#development)
- [Updating for a new LLVM release](#updating-for-a-new-llvm-release)
- [License](#license)

## Supported LLVM versions

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

	// The Go signature maps directly to an LLVM function type.
	mainFn, err := module.NewFunc[func() int32]("main")
	if err != nil {
		panic(err)
	}

	builder := ir.NewBuilder(ctx)
	defer builder.Close()
	builder.MoveToEnd(mainFn.Function().NewBlock("entry"))

	// Type safety: Add only accepts i32 values, ICmp returns i1, and Select
	// requires both arms to have the same type.
	i32 := ctx.Int(32)
	sum := builder.Add(i32.Const(1), i32.Const(2), "sum")
	ok := builder.ICmp(llvm.IntEQ, sum, i32.Const(3), "ok")
	result := builder.Select(ok, i32.Const(0), i32.Const(1), "result")
	builder.Ret(result)

	if err := module.Verify(); err != nil {
		panic(err) // structured *llvm.Error with full diagnostics
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
  * **Crash-class floor (always on, pure Go)**: nil handles, freed handles (`ErrUseAfterFree`),
    cross-`Context` use, and closed `Context`/`Module`/`Builder` handles. Without it, an invalid handle
    inside cgo is a SIGSEGV — `recover` cannot catch it, the process dies without a Go stack. With it,
    every misuse degrades to a catchable Go panic.
  * **Semantic contracts (debug builds only)**: operand type agreement, power-of-two alignment,
    index bounds, atomic orderings, call arity/types, result kinds. Under `-tags=llvm_release` these are
    **compiled out entirely**; misuse falls back to LLVM asserts + `Module.Verify()` — the same contract
    as C/Rust/inkwell release builds.
  * Debug builds additionally provide: single-goroutine contract sampling, automatic `Verify` at
    link/codegen/JIT boundaries, `panic` messages carrying the recent operation context, an
    unreleased-resource report when a `Context` is closed, and a default LLVM diagnostic handler
    (turning the default handler's process exit into a log).
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

tm.ApplyTo(module)                                 // write triple + data layout
_ = tm.EmitToFile(module, "main.o", target.ObjectFile)
asm, _ := tm.Emit(module, target.AsmFile)          // or emit into an in-memory buffer
defer asm.Close()
```

### JIT execution

```go
target.InitNative()
j, _ := jit.NewLLJIT()
defer j.Close()

_ = j.AddProcessSymbols()                          // process symbols (libc/libm, ...) resolve extern declarations in JIT code
_ = j.AddIRModule(module)                          // ownership of the module and its Context moves to the JIT
fib, _ := j.Func[func(int32) int32]("fib")         // a real Go function value
fmt.Println(fib(10))

_ = j.MapFunc("host_cb", func(x int32) int32 { return x * 2 })  // register a host function as a JIT symbol
p, _ := j.Lookup("some_symbol")                    // low-level unsafe.Pointer entry point
code, _ := j.RunMain([]string{"prog"})             // called as main(argc, argv, envp)
```

`Func[F]` / `MapFunc[F]` calls go through `reflect` + a fixed-signature C channel, so each call costs
roughly 0.1–0.5 µs and one small allocation (slot arrays are pooled). Fetch the Go function value once
and keep it. Debug builds verify `F` against the signature the JIT module declared for the symbol
(before registering), so a mismatched signature panics as `ErrTypeMismatch` instead of running
undefined behavior; the trust build (`-tags=llvm_release`) skips this and trusts the caller.
For hot loops, look the symbol up with `Lookup` and call the address from your own cgo
binding — note that a bare `unsafe.Pointer` cannot be called from pure Go without cgo or an
assembly trampoline, so the escape hatch requires a cgo-enabled caller.

`AddProcessSymbols` installs a `DynamicLibrarySearchGenerator` on the main JITDylib, so IR that
only declares a function (e.g. `extern double sin(double)`) resolves against the host process at
materialization time. Call it before `AddIRModule`.

### Kaleidoscope example

[`examples/kaleidoscope`](examples/kaleidoscope) is a runnable implementation of the LLVM
tutorial's Kaleidoscope language (up to chapter 7: JIT and optimization), built entirely on
this library's public API — lexer, Pratt parser, `ir` codegen, `pass` optimization pipeline,
and `jit.LLJIT` + `ResourceTracker` execution. Run it with:

```shell
go run ./examples/kaleidoscope          # REPL
go run ./examples/kaleidoscope -e '1 + 2 * 2'
go run ./examples/kaleidoscope -dl -dp -dc -e '1 + 2 * 2'   # dump tokens/AST/IR
```

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

## Examples

Runnable examples live in [`examples/`](examples); see
[`examples/README.md`](examples/README.md) for the exact commands and expected output:

| Example | Description |
|---|---|
| [hello](examples/hello/) | Build, verify, and print the IR of the smallest module (`1 + 2`). |
| [jit-fib](examples/jit-fib/) | Build a recursive `fib` in IR, JIT it, and call it from Go; also register a host Go callback with `MapFunc`. |
| [codegen](examples/codegen/) | AOT-compile a module to an object file (`.o`) and assembly text (`.s`) with `EmitToFile`. |
| [opt](examples/opt/) | Run the `default<O2>` pipeline with `pass.AutoOpt` and print the instruction count before and after. |
| [kaleidoscope](examples/kaleidoscope/) | A full Kaleidoscope language front-end (tutorial chapters 1–7) on top of go-llvm's JIT. |

## Documentation

* API reference: [pkg.go.dev/github.com/kkkunny/go-llvm](https://pkg.go.dev/github.com/kkkunny/go-llvm)
* Contributor and architecture notes: [`AGENTS.md`](AGENTS.md)

## Development

```shell
go build ./...
go vet ./...
go test ./...                            # debug build (default): full check stack
go test -tags=llvm_release ./...         # trust build: semantic checks compiled out
go test ./ir -run TestGolden -update     # regenerate golden IR (must match in both builds)
make test test-release bench bench-release
```

CI runs on every push and pull request via
[`.github/workflows/ci.yml`](.github/workflows/ci.yml): a `gofmt` lint check plus a
`default`/`release` test matrix on `ubuntu-24.04`, with LLVM 22 installed from
apt.llvm.org.

Negative tests for semantic-contract misuse are gated by `requireDebug(t)` and run in
debug builds only (skipped under the `llvm_release` matrix); crash-class floor tests
must pass in both builds.

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

## License

Licensed under the Apache License, Version 2.0 — see [`LICENSE`](LICENSE).
