# AGENTS.md

Go bindings to a **system-installed LLVM**. Public API is split across `llvm` and its
sub-packages; all cgo lives in `internal/binding`.

## Hard prerequisites

- **This ref targets LLVM 22** (`MIN/MAX_SUPPORT_MAJOR_VERSION = 22` in the Makefile); older majors are out of scope here (see the README support table). If `go build ./...` fails with cgo errors like `could not determine what C.X refers to`, you are compiling against the wrong LLVM major — check `llvm-config --version`.
- **Go 1.27+** is required: the public API uses generic methods (`func (v Value[T]) As[U Kind]()`), which are only available on concrete types. Do not move generic methods into interfaces — Go does not allow it, and interfaces cannot be satisfied by generic methods.
- `#cgo` flags are checked in at `internal/binding/cgo.go`: static multi-candidate `-I`/`-L` dirs (`/usr/lib/llvm-NN` first, then `/usr`, `/usr/local`, `/usr/lib64`) with an unversioned `-lLLVM`. Consumers need no extra setup for standard layouts; this line lists `llvm-22`, and `make config` regenerates the flags when bumping the LLVM line.
- CI runs on every push and pull request via `.github/workflows/ci.yml`: a `gofmt` lint check plus a `default`/`release` test matrix on `ubuntu-24.04`, with LLVM 22 installed from apt.llvm.org. Reproduce the checks locally: `go build ./...`, `go vet ./...`, `go test ./...` **and `go test -tags=llvm_release ./...`** (`make test` / `make test-release` wrap the last two). Default builds run the full check stack, `llvm_release` is the trust build (see *Checks and build modes* below). Golden IR tests are regenerated with `go test ./ir -run TestGolden -update` and must match in both modes. Benchmarks: `make bench` / `make bench-release`, or `go test -run '^$' -bench . -benchmem ./...` (see `bench_test.go`, `ir/bench_test.go`, `jit/bench_test.go`).
- Panic-expected misuse tests are gated by `requireDebug(t)` and skip under `llvm_release`; crash-class floor tests must pass in both modes.

## cgo / Makefile quirks

- The Makefile is the maintainer tool for tracking new LLVM lines: `make config` (or `make config MAJOR_VERSION=NN` to pick the toolchain) rewrites `internal/binding/cgo.go` from the detected `llvm-config`, keeping the multi-candidate layout (`/usr/lib/llvm-NN` first, then `/usr`, `/usr/local`, `/usr/lib64`) and the unversioned `-lLLVM`. Review the `git diff` before committing; no `package main` artifact is produced.
- Consumers on non-standard prefixes override via env: `CGO_CFLAGS` / `CGO_CXXFLAGS` / `CGO_LDFLAGS` from `llvm-config-NN` (global, so they also reach `internal/binding`'s own compilation; a `#cgo` file in the consumer's main package cannot fix its include paths).
- Go callbacks into C use `//export` plus a tiny `.c`/`.h` trampoline (see `ErrorHandling.h`), never `reflect`-built function pointers. The JIT interop bridge (`internal/binding/bridge.c`) follows the same rule: `llvmBridgeGoChannel` → `goLLVMBridgeDispatch` for native→Go, `llvmBridgeCall` for Go→native through per-signature IR adapters.
- ORC ownership: `LLVMOrcCreateNewThreadSafeModule` takes the module, and the ThreadSafeContext takes the LLVMContext. `jit.LLJIT.AddIRModule` therefore calls `Module.Disown()` + `Context.Disown()` (both invalidate Go handles immediately) and never disposes them itself; `LLJIT.Close()` frees them. `AddObjectFile` likewise consumes the `MemoryBuffer` (`Disown`).
- C++ shims are compiled with `-fexceptions` (see `cgo.go`) and must catch exceptions at the `extern "C"` boundary, returning an error message instead of letting them cross into cgo.

## Layout / architecture

| Package | Responsibility | Must not |
|---|---|---|
| `llvm` | `Kind`/`Type[T]`/`Value[T]`, constants, `Context`, errors, lifetime, Go type mapping, `DataLayout`, `MemoryBuffer` | import sub-packages |
| `llvm/ir` | `Module`/`Function`/`Block`/`Global`/`Builder`, instruction roles, `Verify`/print/parse/bitcode | touch JIT/execution |
| `llvm/target` | targets, target machines, codegen (`EmitToFile(m *ir.Module)`) | execute |
| `llvm/jit` | ORC LLJIT, symbol mapping, Go interop bridge | AOT codegen |
| `llvm/pass` | optimization pipelines | define passes |
| `internal/binding` | 1:1 cgo wrappers over LLVM-C + C++ shims | expose high-level API |
| `internal/checks` | check-mode switch (`Debug` const via `llvm_release` build tag) | import project packages |

Dependency direction is strictly one-way: `llvm` ← `llvm/ir` ← `llvm/target` ← `llvm/jit`,
and `llvm/ir` ← `llvm/pass`. `llvm/target` may import `llvm/ir` because codegen and
`ApplyTo` conveniences need `*ir.Module`; `llvm/ir` never imports target/JIT.

### Core conventions

- **All LLVM enum constants** in `internal/binding` are declared `= C.Name` (e.g. `LLVMRet LLVMOpcode = C.LLVMRet`) so they bind to the local headers automatically. Never hand-copy numeric enum tables; new enum members appear by mapping the header name. Public enums in `llvm` (`Linkage`, `IntPred`, ...) forward these values.
- **Kind-level generics only**: `T` in `Value[T]`/`Type[T]` distinguishes categories (`IntT`, `FloatT`, `PtrT`, ...), never bit widths or element types. Widths, element types and function signatures are pre-checked at runtime.
- **Category-specific operations live on role wrappers**, because Go cannot add methods to an instantiated generic type (`Value[IntT]`). `Value[T]`/`Type[T]` only carry generic operations (`String`/`Name`/`Type`/`As`/`Dyn`/`Alive`/`Context`/`Ref`/`IsNil`). Roles: `IntType.Bits()`, `StructType.SetBody()`, `Alloca.SetAlign()`, `Phi.AddIncoming()`, `IntConst.SignedValue()`, ...
- **Roles embed `Value[T]`/`Type[T]` anonymously** (`type Alloca struct{ llvm.Value[llvm.PtrT] }`), so they inherit `ValueRef[T]`/`AnyValue` for free. Do not add a `v` field or a `Value()` accessor; the generic view is `AsValue()`. A role may deliberately shadow a promoted method when the semantics differ (`Global.IsConstant` = global constant flag vs `Value.IsConstant` = constant expression) — document the shadowing in a comment.
- **Handle wrapping goes through `ir/wrap.go`** (`wrapBlock`): a handle's lifetime always comes from its owning `Module` token, never from the call site's local state. Values/types are cast directly with `llvm.NewValue`/`llvm.ValueOf`/`llvm.TypeOfRef`; `llvm.AnyValuesToRefs` converts heterogeneous value lists for binding calls.
- **Uniform lifetime checks / choke points**: the crash-class floor lives at the handle-delivery choke points. `Value.Ref()` / `Type.Ref()` validate nil/context/lifetime and are the only way role methods on `Value[T]`/`Type[T]`-based roles may obtain raw handles — such role methods therefore carry **no** prologue `Check` (do not read `.ref` directly; precise error location comes from the Go stack, and Builder errors still carry `recent builder ops`). Their fast path must stay inlinable: keep the happy-path condition in `Ref` itself and only the panicking slow path in `checkFloor`. `RawRef()` is the unchecked accessor reserved for pre-check internals (`ir.preVal`), where the explicit `checkVal` owns validation. Container handles without a checked `Ref` (`Block`, `Module`, `Comdat`, `Metadata`, `Attribute`) keep their explicit `Check(op)` at method entry. `Context.CheckAlive`/`CheckType`/`CheckTypes`/`CheckValues` are the shared argument/context entry points exported for sub-packages; `Module.Check` is reused by `jit`/`target` for module parameters.
- **`TypeRef[T]` / `ValueRef[T]`**: kind-safe reference interfaces implemented by both `Type[T]`/`Value[T]` and all role wrappers. Use them as parameter types in generic methods/functions so calls like `ctx.ConstNull(ctx.Int(32))` infer `T` from either form.
- **`AnyType` / `AnyValue`**: non-generic views for heterogeneous collections (call args, mixed instructions). They omit methods whose signature would differ per instantiation (`Type()`, `As`); erase with `Dyn()`/`DynType()` and recover a kind via `As[U]()`/`MustAs[U]()` on the generic `Value[DynT]`/`Type[DynT]`. Dynamic sources (parsed IR, instruction iteration) produce `Value[DynT]`.
- **Unknown kinds/opcodes degrade** to `Value[DynT]` instead of panicking.
- **Errors**: recoverable runtime failures return `error`; programmer errors `panic(*llvm.Error)` with `Reason`/`Op`/`Msg`, recoverable via `llvm.Catch` (or `llvm.Try[T]` when a value is produced). Data-driven unsupported cases (`TypeOf[string]`, variadic Go funcs) return `ErrUnsupported`. `internal/binding` returns plain `error` (it must not import the root package); sub-packages wrap it with `llvm.WrapError(reason, op, err)` and use `ErrCodeGen`/`ErrJIT`/`ErrIO`/`ErrParse` for target/JIT/IO/parse failures.
- **Constants**: `ConstInt` truncates like an unsigned value; `ConstSInt` sign-extends (use it for negative or over-wide values). Type-directed sugar: `IntType.Const/ConstS`, `FloatType.Const`.
- **Performance**: the pre-check layer uses the `preVal` value type (`core`/`coreAny`) instead of `AnyValue` interfaces so checks never box `Value[T]`; `Value.ty` caches the type handle in debug builds only (release pays nothing for it) so type/arg checks cost no cgo. `Builder` owns a scratch ref buffer (`valueRefs`/`refs`) reused for call args and GEP indices in release; the debug build allocates fresh per use so buffer-lifetime misuse is impossible. The binding layer passes a static empty C string for `Name` when the name is empty (`string2CString`), avoiding two cgo crossings plus a C malloc/free per instruction; keep this fast path when adding build bindings. Go type/signature mapping is memoized per `Context` (`typeCache`/`fnCache`) with liveness checked before lookup. The JIT bridge pools its slot arrays; `Func[F]`/`MapFunc[F]` calls still cost reflect + one cgo round trip each.
- **Traversal**: prefer the lazy `iter.Seq` APIs (`Block.AllInsts`, `Function.AllBlocks`/`AllParams`, `StructType.AllElems`) over the slice variants in hot paths; slice APIs stay for compatibility.
- **Checks and build modes (三层校验)**: checks are layered, and every new check must be placed by consequence, not by taste.
  1. *Compile-time kind safety* (`Value[T]`/`Type[T]`/`ValueRef`) — always on, prefer this when expressible.
  2. *Crash-class floor* — always on, **pure Go only** (`ref.IsNil()`, `Lifetime.Alive()`, ctx pointer compare, `closed` flags, insert-point tracking). These turn cgo SIGSEGV/UAF into catchable `panic(*llvm.Error)`. Never remove one; never add cgo to this layer.
  3. *Semantic contracts* — guard with `if checks.Debug { ... }` (`internal/checks`): operand/argument type equality, power-of-two alignment, index bounds, atomic orderings, call arity, result kinds, goroutine-owner sampling, boundary `Verify`, and JIT signature verification (`Func[F]`/`MapFunc[F]` are checked against signatures recorded by `AddIRModule`). In `-tags=llvm_release` these are dead-code-eliminated, so misuse falls back to LLVM asserts + `Module.Verify()` (same contract as C/Rust/inkwell).
  `ir.Builder` methods still start with `pre`/`preSameType`/`preBlock` (plus the package-level `preAlign`) — positioned builder, same context, live handles first; then the gated semantics; role methods on value/type roles rely on `Ref()`'s floor (see *Uniform lifetime checks / choke points*). Positioning methods use `preAlive` plus their own operand checks (`MoveToEnd`: `preBlockOwn`; `MoveBefore`: `pre` on the instruction). Debug builds add pending diagnostics (default handler installed in `NewContext`), `requireDebug(t)` for misuse tests, and resource/leak reports in `Context.Close`.
- **Builder return types**: return a role wrapper only when the instruction has role-specific operations (`Alloca`/`Load`/`Store`/`Call`/`Invoke`/`Phi`/`Switch`/`LandingPad`/`CatchSwitch`/`FuncletPad`/`Fence`/`AtomicRMW`/`CmpXchg`); everything else returns the plain `Value[T]`.
- **Concurrency**: `Context` (ownership registry), `Lifetime` (atomic), the `LLJIT` adapter cache and the bridge registry are lock-protected; all other handles (`Module`/`Builder`/`Value`/`Type`/`Block`/`TargetMachine`/`DataLayout`/`MemoryBuffer`) are not goroutine-safe and must be used from a single goroutine. Debug builds sample the owning goroutine in `Builder`/`Module` operations and panic on cross-goroutine use (the race detector cannot see C-side state). `LLJIT.Func`/`MapFunc`/`Lookup` may be called concurrently, but `Close` must be serialized by the caller.
- **Lifetime**: `Context` is the ownership root (`Own` returns an unregister func, `Close` cascades in reverse). `Module`/`Builder` implement `io.Closer`; `Value`/`Type`/`Block`/`GoFunc` never expose `Free` — they carry a `Lifetime` token and are checked on every operation. Second `Close` returns `ErrClosed`. `Module.Disown()`/`Context.Disown()`/`MemoryBuffer.Disown()` transfer ownership to an external owner (JIT) and **immediately** invalidate Go-side handles (no return value); `Context`-independent resources (`TargetMachine`, `MemoryBuffer`, `LLJIT`) are their own roots and are not registered via `Own`. `MemoryBuffer`/`DataLayout` carry a GC finalizer as a leak safety net — `Close`/`Disown` must clear it (`runtime.SetFinalizer(x, nil)`) before releasing the handle.
- **Comments**: Package docs in English; symbol comments in Chinese for the root and sub-packages; `internal/binding` in English.
- Commits use Conventional Commits, commonly with Chinese descriptions.

## Doc comments

Doc-comment language follows the *Comments* rule in *Core conventions*: package docs in English,
symbol comments in Chinese for the root and sub-packages, `internal/binding` in English. The
rules below apply to both languages.

- The first sentence begins with the symbol name and ends with a period (`。` in Chinese symbol
  comments): `NewModule 创建模块并登记到 Context 生命周期。`.
- Cross-reference symbols as `[Symbol]` / `[Type.Method]`, and packages by full import path
  (`[github.com/kkkunny/go-llvm/ir.Module]`). Keep links few and resolvable — an unresolvable
  reference does not error, but it carries no information either.
- Document a `const (...)` block with one block-level comment describing the value domain and
  the LLVM mapping, plus a short inline comment per constant (the `error.go` style); do not give
  every constant its own doc paragraph.
- Method/function docs state semantics, units, ownership, and error/panic conditions
  (`panic(*llvm.Error)`, `[ErrUseAfterFree]`, ...). Do not write filler that merely restates the
  name (`Name 返回名字` is wrong).
- Name godoc examples `Example` (package overview), `ExampleContext` (function),
  `ExampleBuilder_Add` (method); a second example for the same symbol takes a suffix
  (`ExampleContext_second`). Every example needs a stable `// Output:` and must pass in both
  build modes.
- Normalize doc comments with `gofmt`; after touching them, `gofmt -l .` must produce no output.

## Binding additions

Add new bindings in the same `/* #include ... */` + `import "C"` style, mapping the exact
C API name. Prefer wrapping the local header declaration over re-declaring it. C++ shims
(`Core.cpp`) exist only for APIs missing from LLVM-C.

As-built extras: `Context.SetDiagnosticHandler`/`ClearDiagnosticHandler` (Go callback registry +
`ErrorHandling.c` trampoline; LLVM's default handler calls `exit(1)` on errors — `NewContext` now
installs a default Go handler that logs instead, so unhandled diagnostics no longer kill the
process), `Module.Link` (same-context pre-check; the source module is consumed by LLVM and its Go
handle dies immediately; debug builds verify both modules first), `Module.AppendCtor/AppendDtor`,
and `llvm/pass` (`RunPasses`/`AutoOpt`/`RunPassesOnFunction`, all through PassBuilder). Debug-only
boundary verification also guards `target.EmitToFile` and `jit.AddIRModule`.

The binding layer keeps its 1:1 LLVM-C mapping role: wrappers with no public consumer yet are allowed
to stay, and cleanup targets only non-mapping dead code (old shims, orphan helpers) rather than the
mapping itself; new bindings are added on demand.
