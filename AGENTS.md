# AGENTS.md

Go bindings to a **system-installed LLVM**. Public API is the root `llvm` package; all cgo lives in `internal/binding`.

## Hard prerequisites

- **Only LLVM 20 is supported** (`MIN/MAX_SUPPORT_MAJOR_VERSION = 20` in the Makefile; commit "支持且仅支持llvm20"). If `go build ./...` fails with cgo errors like `could not determine what C.LLVMConstMul refers to`, check `llvm-config --version` first — it is an LLVM version mismatch, not a code bug. (This host only has LLVM 22, so the build currently fails that way.)
- No `#cgo` directives exist anywhere in this repo; they are supplied by the consumer.
- No tests, CI, or lint config exist. `go build ./...` is the only verification step.

## cgo / Makefile quirks

- The Makefile is shipped to consumer projects (see README); `make config` there emits `llvm_config.go` with `package main` plus the `#cgo` flags. **Never run `make config` in this repo root** — it creates a `package main` file inside package `llvm` and breaks the build.
- README's `make config LLVM_CONFIG_BIN=20` is stale and literally runs the command `20`; the Makefile selects the binary via `make config MAJOR_VERSION=20`.
- To link a scratch binary/test here instead of generating that file, set `CGO_CFLAGS` / `CGO_CXXFLAGS` / `CGO_LDFLAGS` from `llvm-config-20`.

## Layout / architecture

- `internal/binding/`: near 1:1 cgo wrappers over the LLVM-C API. `Core.go` is the large hand-maintained one. `Core.cpp`, `PassManager.cpp` and their headers are C++ shims for APIs missing from LLVM-C (`PassManager.cpp` backs `LLVMOptModule`). Add new bindings in the same `/* #include ... */` + `import "C"` style.
- Root `llvm` package: wrapper types are plain casts of binding refs, e.g. `type Block binding.LLVMBasicBlockRef`, each with a `binding()` method.
  - `Value` is an interface; `lookupValue` (value.go:15) and `lookupInstruction` (instruction.go) dispatch on LLVM value kind / opcode. Adding an instruction requires a marker type (`_Add`, `_Alloca`, ...), a `case` in `lookupInstruction`, and a builder method.
  - Instruction wrappers are generic: `valueInst[T]` / `valueInstWithAlign[T]` around `genericInst[T]` (instruction.go:174).
- JIT-to-Go callbacks (`ExecutionEngine.MapFunctionToGo`) exist only on amd64 (`execution_amd64.go` + `execution.c/h`); there is no fallback on other architectures.
- `internal/binding/llvm_config.go` is a checked-in wrapper for `llvm/Config/llvm-config.h` (version constants) — unrelated to the Makefile-generated root `llvm_config.go`.

## Conventions

- Go 1.22; only runtime dependency is `github.com/samber/lo`.
- Comments: `internal/binding` in English, root package often in Chinese; match the file you edit.
- Commits use Conventional Commits, commonly with Chinese descriptions.
