# codegen

A runnable go-llvm ahead-of-time code generation example. It initializes the
native target with `target.InitNative`, creates a `target.TargetMachine`,
copies the target triple and data layout into the module with `ApplyTo`, and
compiles the module both to an object file (`.o`) and to assembly text (`.s`)
with `EmitToFile`. The artifacts go to a directory created by `os.MkdirTemp`,
which is removed when the example exits.

## Run

Requires LLVM 22 (same as the library) and Go 1.27+:

```shell
go run ./examples/codegen
```

## Expected output

For each artifact the example prints its full path and its size in bytes. The
directory has a random suffix and the sizes depend on the host target, so the
output looks like this:

```text
object: /tmp/go-llvm-codegen-1234567890/add.o (752 bytes)
assembly: /tmp/go-llvm-codegen-1234567890/add.s (213 bytes)
```

The example panics if code generation fails or if either artifact is empty; on
the normal path it never panics.
