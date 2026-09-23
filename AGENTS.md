# AGENTS.md

Go bindings to a **system-installed LLVM**. Public API is split across `llvm` and its
sub-packages; all cgo lives in `internal/binding`.

## Hard prerequisites

- **master targets LLVM 22** (`MIN/MAX_SUPPORT_MAJOR_VERSION = 22` in the Makefile). Older lines are pinned with git tags (currently `@llvm21`) and carry their own flags/Makefile numbers. If `go build ./...` fails with cgo errors like `could not determine what C.X refers to`, you are compiling this ref against the wrong LLVM major — check `llvm-config --version`.
- **Go 1.27+** is required: the public API uses generic methods (`func (v Value[T]) As[U Kind]()`), which are only available on concrete types. Do not move generic methods into interfaces — Go does not allow it, and interfaces cannot be satisfied by generic methods.
- `#cgo` flags are checked in at `internal/binding/cgo.go`: static multi-candidate `-I`/`-L` dirs (`/usr/lib/llvm-NN` first, then `/usr`, `/usr/local`, `/usr/lib64`) with an unversioned `-lLLVM`. Consumers need no extra setup for standard layouts; version dirs are per-line (master lists `llvm-22`).
- `go build ./...`, `go vet ./...` and `go test ./...` are the verification steps. Golden IR tests can be regenerated with `go test ./ir -run TestGolden -update`. Benchmarks: `go test -run '^$' -bench . -benchmem ./...` (see `bench_test.go`, `ir/bench_test.go`, `jit/bench_test.go`).

## cgo / Makefile quirks

- The Makefile is now only an escape hatch for non-standard prefixes (README, "Non-standard LLVM prefixes"): `make config` emits `llvm_config.go` with `package main` for the **consumer's** main package. **Never run `make config` in this repo root** — it creates a `package main` file inside the library and breaks the build. Select the toolchain with `make config MAJOR_VERSION=NN`.
- Same escape hatch via env: `CGO_CFLAGS` / `CGO_CXXFLAGS` / `CGO_LDFLAGS` from `llvm-config-NN`.
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
| `llvm/pass` (P2) | optimization pipelines | define passes |
| `internal/binding` | 1:1 cgo wrappers over LLVM-C + C++ shims | expose high-level API |

Dependency direction is strictly one-way: `llvm` ← `llvm/ir` ← `llvm/target` ← `llvm/jit`,
and `llvm/ir` ← `llvm/pass`. `llvm/target` may import `llvm/ir` because codegen and
`ApplyTo` conveniences need `*ir.Module`; `llvm/ir` never imports target/JIT.

### Core conventions

- **All LLVM enum constants** in `internal/binding` are declared `= C.Name` (e.g. `LLVMRet LLVMOpcode = C.LLVMRet`) so they bind to the local headers automatically. Never hand-copy numeric enum tables; new enum members appear by mapping the header name. Public enums in `llvm` (`Linkage`, `IntPred`, ...) forward these values.
- **Kind-level generics only**: `T` in `Value[T]`/`Type[T]` distinguishes categories (`IntT`, `FloatT`, `PtrT`, ...), never bit widths or element types. Widths, element types and function signatures are pre-checked at runtime.
- **Category-specific operations live on role wrappers**, because Go cannot add methods to an instantiated generic type (`Value[IntT]`). `Value[T]`/`Type[T]` only carry generic operations (`String`/`Name`/`Type`/`As`/`Dyn`/`Alive`/`Context`/`Ref`/`IsNil`). Roles: `IntType.Bits()`, `StructType.SetBody()`, `Alloca.SetAlign()`, `Phi.AddIncoming()`, `IntConst.SignedValue()`, ...
- **Roles embed `Value[T]`/`Type[T]` anonymously** (`type Alloca struct{ llvm.Value[llvm.PtrT] }`), so they inherit `ValueRef[T]`/`AnyValue` for free. Do not add a `v` field or a `Value()` accessor; the generic view is `AsValue()`. A role may deliberately shadow a promoted method when the semantics differ (`Global.IsConstant` = global constant flag vs `Value.IsConstant` = constant expression) — document the shadowing in a comment.
- **Handle wrapping goes through `ir/wrap.go`** (`wrapBlock`): a handle's lifetime always comes from its owning `Module` token, never from the call site's local state. Values/types are cast directly with `llvm.NewValue`/`llvm.ValueOf`/`llvm.TypeOfRef`; `llvm.AnyValuesToRefs` converts heterogeneous value lists for binding calls.
- **Uniform lifetime checks**: `Value.Check`/`Type.Check` (roles inherit them) plus `Context.CheckAlive`/`CheckType`/`CheckTypes`/`CheckValues` are the shared pre-check entry points exported for sub-packages. Every `ir` role method that calls binding directly must start with `Check`/`pre`; `Module.Check`/`Block.Check` guard their handles (nil/closed/ownership-transferred) and are reused by `jit`/`target` for module parameters.
- **`TypeRef[T]` / `ValueRef[T]`**: kind-safe reference interfaces implemented by both `Type[T]`/`Value[T]` and all role wrappers. Use them as parameter types in generic methods/functions so calls like `ctx.ConstNull(ctx.Int(32))` infer `T` from either form.
- **`AnyType` / `AnyValue`**: non-generic views for heterogeneous collections (call args, mixed instructions). They deliberately exclude `Type()`/`DynType()`-style methods that would differ per instantiation; use `Dyn()`/`DynType()` for erasure and `As[U]()`/`MustAs[U]()` to recover a kind. Dynamic sources (parsed IR, instruction iteration) produce `Value[DynT]`.
- **Unknown kinds/opcodes degrade** to `Value[DynT]` instead of panicking.
- **Errors**: recoverable runtime failures return `error`; programmer errors `panic(*llvm.Error)` with `Reason`/`Op`/`Msg`, recoverable via `llvm.Catch` (or `llvm.Try[T]` when a value is produced). Data-driven unsupported cases (`TypeOf[string]`, variadic Go funcs) return `ErrUnsupported`. `internal/binding` returns plain `error` (it must not import the root package); sub-packages wrap it with `llvm.WrapError(reason, op, err)` and use `ErrCodeGen`/`ErrJIT`/`ErrIO`/`ErrParse` for target/JIT/IO/parse failures.
- **Constants**: `ConstInt` truncates like an unsigned value; `ConstSInt` sign-extends (use it for negative or over-wide values). Type-directed sugar: `IntType.Const/ConstS`, `FloatType.Const`.
- **Performance**: the pre-check layer uses the `preVal` value type (`core`/`coreAny`) instead of `AnyValue` interfaces so checks never box `Value[T]`; `Builder` owns a scratch ref buffer (`valueRefs`/`refs`) reused for call args and GEP indices; `Builder.call` reads parameter types straight from binding. Checks that must ask LLVM (`LLVMTypeOf` for operand/argument type equality) are the accepted cgo floor — do not add cgo to the pre-check path without a benchmark. Go type/signature mapping is memoized per `Context` (`typeCache`/`fnCache`) with liveness checked before lookup. The JIT bridge pools its slot arrays; `Func[F]`/`MapFunc[F]` calls still cost reflect + one cgo round trip each.
- **Traversal**: prefer the lazy `iter.Seq` APIs (`Block.AllInsts`, `Function.AllBlocks`/`AllParams`, `StructType.AllElems`) over the slice variants in hot paths; slice APIs stay for compatibility.
- **Pre-checks before cgo**: every `ir.Builder` method calls `pre`/`preSameType`/`preBlock`/`preAlign` first (positioned builder, same context, live handles, matching operand types, power-of-two alignment) so LLVM never sees invalid IR and never aborts. Positioning methods use `preAlive` plus their own operand checks instead (`MoveToEnd`: `preBlockOwn`; `MoveBefore`: `pre` on the instruction), since they must work before the first `pre`.
- **Builder return types**: return a role wrapper only when the instruction has role-specific operations (`Alloca`/`Load`/`Store`/`Call`/`Invoke`/`Phi`/`Switch`/`LandingPad`/`CatchSwitch`/`FuncletPad`/`Fence`/`AtomicRMW`/`CmpXchg`); everything else returns the plain `Value[T]`.
- **Concurrency**: `Context` (ownership registry), `Lifetime` (atomic), the `LLJIT` adapter cache and the bridge registry are lock-protected; all other handles (`Module`/`Builder`/`Value`/`Type`/`Block`/`TargetMachine`/`DataLayout`/`MemoryBuffer`) are not goroutine-safe and must be used from a single goroutine. `LLJIT.Func`/`MapFunc`/`Lookup` may be called concurrently, but `Close` must be serialized by the caller.
- **Lifetime**: `Context` is the ownership root (`Own` returns an unregister func, `Close` cascades in reverse). `Module`/`Builder` implement `io.Closer`; `Value`/`Type`/`Block`/`GoFunc` never expose `Free` — they carry a `Lifetime` token and are checked on every operation. Second `Close` returns `ErrClosed`. `Module.Disown()`/`Context.Disown()`/`MemoryBuffer.Disown()` transfer ownership to an external owner (JIT) and **immediately** invalidate Go-side handles (no return value); `Context`-independent resources (`TargetMachine`, `MemoryBuffer`, `LLJIT`) are their own roots and are not registered via `Own`. `MemoryBuffer`/`DataLayout` carry a GC finalizer as a leak safety net — `Close`/`Disown` must clear it (`runtime.SetFinalizer(x, nil)`) before releasing the handle.
- **Comments**: `internal/binding` in English, root and sub-packages in Chinese; match the file you edit.
- Commits use Conventional Commits, commonly with Chinese descriptions.

## Binding additions

Add new bindings in the same `/* #include ... */` + `import "C"` style, mapping the exact
C API name. Prefer wrapping the local header declaration over re-declaring it. C++ shims
(`Core.cpp`, `PassManager.cpp`) exist only for APIs missing from LLVM-C.

Some bindings are ahead of the public API and wait for their P2 consumers:
`Linker.go` (`Module.Link`), `PassBuilder.go`/`PassManager.*`/`OptimizationLevel.go` (`llvm/pass`).
