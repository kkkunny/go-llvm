# AGENTS.md

Go bindings to a **system-installed LLVM**. Public API is split across `llvm` and its
sub-packages; all cgo lives in `internal/binding`.

## Hard prerequisites

- **master targets LLVM 22** (`MIN/MAX_SUPPORT_MAJOR_VERSION = 22` in the Makefile). Older lines are pinned with git tags (currently `@llvm21`) and carry their own flags/Makefile numbers. If `go build ./...` fails with cgo errors like `could not determine what C.X refers to`, you are compiling this ref against the wrong LLVM major — check `llvm-config --version`.
- **Go 1.27+** is required: the public API uses generic methods (`func (v Value[T]) As[U Kind]()`), which are only available on concrete types. Do not move generic methods into interfaces — Go does not allow it, and interfaces cannot be satisfied by generic methods.
- `#cgo` flags are checked in at `internal/binding/cgo.go`: static multi-candidate `-I`/`-L` dirs (`/usr/lib/llvm-NN` first, then `/usr`, `/usr/local`, `/usr/lib64`) with an unversioned `-lLLVM`. Consumers need no extra setup for standard layouts; version dirs are per-line (master lists `llvm-22`).
- `go build ./...`, `go vet ./...` and `go test ./...` are the verification steps. Golden IR tests can be regenerated with `go test ./ir -run TestGolden -update`.

## cgo / Makefile quirks

- The Makefile is now only an escape hatch for non-standard prefixes (README, "Non-standard LLVM prefixes"): `make config` emits `llvm_config.go` with `package main` for the **consumer's** main package. **Never run `make config` in this repo root** — it creates a `package main` file inside the library and breaks the build. Select the toolchain with `make config MAJOR_VERSION=NN`.
- Same escape hatch via env: `CGO_CFLAGS` / `CGO_CXXFLAGS` / `CGO_LDFLAGS` from `llvm-config-NN`.
- Go callbacks into C use `//export` plus a tiny `.c`/`.h` trampoline (see `ErrorHandling.c`), never `reflect`-built function pointers.

## Layout / architecture

| Package | Responsibility | Must not |
|---|---|---|
| `llvm` | `Kind`/`Type[T]`/`Value[T]`, constants, `Context`, errors, lifetime, Go type mapping, `DataLayout` | import sub-packages |
| `llvm/ir` | `Module`/`Function`/`Block`/`Global`/`Builder`, instruction roles, `Verify`/print | touch target/JIT |
| `llvm/target` (P1) | targets, target machines, codegen | execute |
| `llvm/jit` (P1) | ORC LLJIT, symbol mapping, Go interop bridge | AOT codegen |
| `llvm/pass` (P2) | optimization pipelines | define passes |
| `internal/binding` | 1:1 cgo wrappers over LLVM-C + C++ shims | expose high-level API |

Dependency direction is strictly one-way: `llvm` ← `llvm/target` ← `llvm/ir` ← `llvm/pass`/`llvm/jit`.

### Core conventions

- **All LLVM enum constants** in `internal/binding` are declared `= C.Name` (e.g. `LLVMRet LLVMOpcode = C.LLVMRet`) so they bind to the local headers automatically. Never hand-copy numeric enum tables; new enum members appear by mapping the header name. Public enums in `llvm` (`Linkage`, `IntPred`, ...) forward these values.
- **Kind-level generics only**: `T` in `Value[T]`/`Type[T]` distinguishes categories (`IntT`, `FloatT`, `PtrT`, ...), never bit widths or element types. Widths, element types and function signatures are pre-checked at runtime.
- **Category-specific operations live on role wrappers**, because Go cannot add methods to an instantiated generic type (`Value[IntT]`). `Value[T]`/`Type[T]` only carry generic operations (`String`/`Name`/`Type`/`As`/`Dyn`/`Alive`/`Context`/`Ref`/`IsNil`). Roles: `IntType.Bits()`, `StructType.SetBody()`, `Alloca.SetAlign()`, `Phi.AddIncoming()`, `IntConst.SignedValue()`, ...
- **`TypeRef[T]` / `ValueRef[T]`**: kind-safe reference interfaces implemented by both `Type[T]`/`Value[T]` and all role wrappers. Use them as parameter types in generic methods/functions so calls like `ctx.ConstNull(ctx.Int(32))` infer `T` from either form.
- **`AnyType` / `AnyValue`**: non-generic views for heterogeneous collections (call args, mixed instructions). They deliberately exclude `Type()`/`DynType()`-style methods that would differ per instantiation; use `Dyn()`/`DynType()` for erasure and `As[U]()`/`MustAs[U]()` to recover a kind. Dynamic sources (parsed IR, instruction iteration) produce `Value[DynT]`.
- **Unknown kinds/opcodes degrade** to `Value[DynT]` instead of panicking.
- **Errors**: recoverable runtime failures return `error`; programmer errors `panic(*llvm.Error)` with `Reason`/`Op`/`Msg`, recoverable via `llvm.Catch`. Data-driven unsupported cases (`TypeOf[string]`, variadic Go funcs) return `ErrUnsupported`.
- **Pre-checks before cgo**: every `ir.Builder` method calls `pre`/`preSameType`/`preBlock`/`preAlign` first (positioned builder, same context, live handles, matching operand types, power-of-two alignment) so LLVM never sees invalid IR and never aborts.
- **Lifetime**: `Context` is the ownership root (`Own` registers sub-resources, `Close` cascades in reverse). `Module`/`Builder` implement `io.Closer`; `Value`/`Type`/`Block`/`Func` never expose `Free` — they carry a `Lifetime` token and are checked on every operation. Second `Close` returns `ErrClosed`.
- **Comments**: `internal/binding` in English, root and sub-packages in Chinese; match the file you edit.
- Commits use Conventional Commits, commonly with Chinese descriptions.

## Binding additions

Add new bindings in the same `/* #include ... */` + `import "C"` style, mapping the exact
C API name. Prefer wrapping the local header declaration over re-declaring it. C++ shims
(`Core.cpp`, `PassManager.cpp`) exist only for APIs missing from LLVM-C.
