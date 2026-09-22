# AGENTS.md

Go bindings to a **system-installed LLVM**. Public API is the root `llvm` package; all cgo lives in `internal/binding`.

## Hard prerequisites

- **master targets LLVM 22** (`MIN/MAX_SUPPORT_MAJOR_VERSION = 22` in the Makefile). Older lines are pinned with git tags (currently `@llvm21`) and carry their own flags/Makefile numbers. If `go build ./...` fails with cgo errors like `could not determine what C.X refers to`, you are compiling this ref against the wrong LLVM major — check `llvm-config --version`.
- `#cgo` flags are checked in at `internal/binding/cgo.go`: static multi-candidate `-I`/`-L` dirs (`/usr/lib/llvm-NN` first, then `/usr`, `/usr/local`, `/usr/lib64`) with an unversioned `-lLLVM`. Consumers need no extra setup for standard layouts; version dirs are per-line (master lists `llvm-22`).
- No tests, CI, or lint config exist. `go build ./...` is the only verification step.

## cgo / Makefile quirks

- The Makefile is now only an escape hatch for non-standard prefixes (README, "Non-standard LLVM prefixes"): `make config` emits `llvm_config.go` with `package main` for the **consumer's** main package. **Never run `make config` in this repo root** — it creates a `package main` file inside package `llvm` and breaks the build. Select the toolchain with `make config MAJOR_VERSION=NN`.
- Same escape hatch via env: `CGO_CFLAGS` / `CGO_CXXFLAGS` / `CGO_LDFLAGS` from `llvm-config-NN`.

## Layout / architecture

- `internal/binding/`: near 1:1 cgo wrappers over the LLVM-C API. `Core.go` is the large hand-maintained one. `Core.cpp`, `PassManager.cpp` and their headers are C++ shims for APIs missing from LLVM-C (`PassManager.cpp` backs `LLVMOptModule`). Add new bindings in the same `/* #include ... */` + `import "C"` style.
- **All LLVM enum constants** in `internal/binding` are declared `= C.Name` (e.g. `LLVMRet LLVMOpcode = C.LLVMRet`) so they bind to the local headers automatically. Never hand-copy numeric enum tables; new enum members appear by mapping the header name.
- Root `llvm` package: wrapper types are plain casts of binding refs, e.g. `type Block binding.LLVMBasicBlockRef`, each with a `binding()` method.
  - `Value` is an interface; `lookupValue` (`value.go`) and `lookupInstruction` (`instruction.go`) dispatch on LLVM value kind / opcode. Adding an instruction requires a marker type (`_Add`, `_Alloca`, ...), a `case` in `lookupInstruction`, and a builder method. Unknown kinds/opcodes degrade to `fallbackValue` / `_Fallback` wrappers instead of panicking.
  - Instruction wrappers are generic: `valueInst[T]` / `valueInstWithAlign[T]` around `genericInst[T]` (`instruction.go`).
- JIT-to-Go callbacks (`ExecutionEngine.MapFunctionToGo`) exist only on amd64 (`execution_amd64.go` + `execution.c/h`); there is no fallback on other architectures.
- `internal/binding/llvm_config.go` is a checked-in wrapper for `llvm/Config/llvm-config.h` (version constants) — unrelated to the Makefile-generated root `llvm_config.go`.
- Version-line maintenance procedure lives in README, "Updating for a new LLVM release".

## Conventions

- Go 1.22; only runtime dependency is `github.com/samber/lo`.
- Comments: `internal/binding` in English, root package often in Chinese; match the file you edit.
- Commits use Conventional Commits, commonly with Chinese descriptions.
