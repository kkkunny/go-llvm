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
| 22 | `go get github.com/kkkunny/go-llvm` |

Notes:

* Requires **Go 1.27+** (the API uses generic methods).
* Only LLVM 22 is supported; LLVM 21 and earlier are no longer provided. They can be
  pinned to historic commits if needed.
* The common Linux, macOS (Homebrew) and FreeBSD install layouts work with no extra
  setup; custom prefixes need a one-time step — see
  [Non-standard LLVM prefixes](#non-standard-llvm-prefixes).

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

The standard Linux/macOS/FreeBSD layouts are covered out of the box; for other
install prefixes see [Non-standard LLVM prefixes](#non-standard-llvm-prefixes).

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

`internal/binding/cgo.go` carries a portable per-OS candidate list (the **A layer**),
so the common install layouts work out of the box — no environment variables needed:

| Platform | Layout covered | Typical install |
|---|---|---|
| Linux | `/usr/lib/llvm-22`, `/usr/include/llvm-22`, `/usr/include/llvm-c-22`, `/usr/include`, `/usr/local`, `/usr/lib64` | Debian/Ubuntu `llvm-22-dev` (apt.llvm.org); Arch |
| macOS (Homebrew) | `/opt/homebrew/opt/llvm@22` (Apple Silicon), `/usr/local/opt/llvm@22` (Intel) | `brew install llvm@22` |
| FreeBSD | `/usr/local/llvm22` | `pkg install llvm22` (not confirmed on real hardware yet) |

The Debian/Ubuntu layouts are exercised by CI and Arch has been verified locally; the Fedora and
FreeBSD candidates are listed for convenience but have not been confirmed on real hardware yet —
if they miss, the generator (B layer) below is the fallback.

`-lLLVM` is intentionally unversioned: it resolves against the first matching library
(Debian/Ubuntu `llvm-22-dev` ships both `libLLVM-22.so` and `libLLVM.so`). If your LLVM
lives elsewhere (custom `--prefix`, a non-standard multi-version toolchain, a store
layout), pick one of the two paths below.

#### In a local checkout (B layer): generate machine-specific flags

```shell
# from the repository root (./internal/binding is a cwd-relative package pattern)
go generate ./internal/binding                        # or: make config
LLVM_CONFIG=/path/to/llvm-config go generate ./internal/binding
LLVM_PREFIX=/path/to/prefix go generate ./internal/binding

# from any other directory inside this module: use the full package path
go run github.com/kkkunny/go-llvm/internal/cmd/llvmconfig
```

`internal/cmd/llvmconfig` probes `$LLVM_CONFIG` → `$LLVM_PREFIX/bin/llvm-config` →
`llvm-config-22` → `llvm-config` on `PATH`, queries `--includedir`/`--libdir`/`--libs`/
`--system-libs`, and rewrites `internal/binding/cgo.go` with those flags in **both**
`CFLAGS` and `CXXFLAGS` (the C++ shims include LLVM headers too). The result is
**machine-specific**: do not commit it unless everyone shares the same layout — keep the
portable candidate list in version control. `go run ./internal/cmd/llvmconfig --check`
only prints what the generator detected, without writing anything.

#### As a dependency from the module cache: override the cgo flags

The module cache is read-only and managed by the `go` command; the generator refuses to
write there and tells you so. Consumers who `go get` this library and have no local
checkout must export the flags themselves:

```shell
export CGO_CFLAGS="$(llvm-config-22 --cflags)"
# --cxxflags ends with -fno-exceptions, which would break the C++ shims; filter it out.
export CGO_CXXFLAGS="$(llvm-config-22 --cxxflags | sed 's/-fno-exceptions//g')"
export CGO_LDFLAGS="$(llvm-config-22 --ldflags --libs)"
```

`CGO_CFLAGS`/`CGO_CXXFLAGS` are global, so they also reach `internal/binding`'s own
compilation; a `#cgo` file in your own main package would not. Instead of environment
variables you can `replace` the module with a local checkout and generate there, then copy
`internal/binding/cgo.go` into the `vendor/` copy — `go mod vendor` overwrites `vendor/`,
and a vendor tree has no `go.mod`, so the generator cannot run inside it.

#### Troubleshooting

| Symptom | Likely cause | Fix |
|---|---|---|
| `fatal error: llvm-c/Core.h: No such file or directory` | include path missed by the candidates | generate flags (B layer) or set `CGO_CFLAGS` **and** `CGO_CXXFLAGS` |
| `could not determine what C.X refers to` | same — cgo compiled without the LLVM headers | same |
| `cannot find -lLLVM` / undefined `LLVM*` symbols | library path (or library name) missed | generate flags (B layer) or set `CGO_LDFLAGS`; the generator may emit `-lLLVM-22` instead of `-lLLVM` — treat a miss on either the same way |
| `panic: llvm.NewContext: linked LLVM library is …` (`ErrVersionMismatch`, debug builds) | runtime library major ≠ compile-time headers major | install the dev package matching the library, or regenerate the flags for that version |

Quick self-check inside a checkout: `go run ./internal/cmd/llvmconfig --check`. Without a
checkout: `llvm-config --includedir --libdir --libs` (use `llvm-config-22` when several
LLVM majors are installed).

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
`default`/`release` test matrix on `ubuntu-24.04` with LLVM 22 installed from
apt.llvm.org (the A-layer Linux layout), a macOS job for the Homebrew layout, and a
custom-prefix job that regenerates the cgo flags with the B-layer generator and builds
against them.

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
   `Makefile` `MIN/MAX_SUPPORT_MAJOR_VERSION`, the support table above, this README,
   and the per-OS candidate paths in `internal/binding/cgo.go`. Check the flags against
   the local toolchain with `LLVM_CONFIG=/path/to/llvm-config-NN go generate
   ./internal/binding` (or `go run ./internal/cmd/llvmconfig --check`); the generated
   file is machine-specific, so keep the portable candidates in the commit.

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
