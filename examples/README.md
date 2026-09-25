# Examples

Runnable examples for go-llvm. All of them need LLVM 22 (same as the library)
and Go 1.27+, and can be run from the repository root:

| Example | Description | Command |
|---|---|---|
| [hello](hello/) | Build, verify, and print the IR of the smallest module (`1 + 2`). | `go run ./examples/hello` |
| [jit-fib](jit-fib/) | Build a recursive `fib` in IR, JIT it, and call it from Go; also register a host Go callback with `MapFunc`. | `go run ./examples/jit-fib` |
| [codegen](codegen/) | AOT-compile a module to an object file (`.o`) and assembly text (`.s`) with `EmitToFile`. | `go run ./examples/codegen` |
| [opt](opt/) | Run the `default<O2>` pipeline with `pass.AutoOpt` and print the instruction count before and after. | `go run ./examples/opt` |
| [kaleidoscope](kaleidoscope/) | A full Kaleidoscope language front-end (tutorial chapters 1–7) on top of go-llvm's JIT. | `go run ./examples/kaleidoscope` |

Each example has its own README with the exact expected output; the
kaleidoscope one also documents its REPL flags (`-e`, `-dl`, `-dp`, `-dc`).
