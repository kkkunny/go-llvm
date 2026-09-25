# hello

The smallest runnable go-llvm program: build a `main` function that returns
`1 + 2`, verify the module, and print its LLVM IR. It is the runnable version
of the quick-start snippet in the repository root
[`example_test.go`](../../example_test.go).

## Run

Requires LLVM 22 (same as the library) and Go 1.27+:

```shell
go run ./examples/hello
```

## Expected output

LLVM constant-folds the add, so the emitted function is `ret i32 3`:

```llvm
; ModuleID = 'main'
source_filename = "main"

define i32 @main() {
entry:
  ret i32 3
}
```

The example never panics on the normal path; `Verify` failures panic with the underlying `*llvm.Error`.
