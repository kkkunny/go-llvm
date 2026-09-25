# opt

A runnable go-llvm optimization example. It builds the deliberately wasteful
function below — a dead store, an `alloca` that is only loaded once, and the
`+ 0` / `* 1` identities:

```llvm
define i32 @compute() {
entry:
  %slot = alloca i32, align 4
  store i32 1, ptr %slot, align 4
  store i32 2, ptr %slot, align 4
  %x = load i32, ptr %slot, align 4
  %a = add i32 %x, 0
  %v = mul i32 %a, 1
  ret i32 %v
}
```

then runs the default `-O2` pipeline with `pass.AutoOpt` and prints the number
of instructions in the module before and after.

## Run

Requires LLVM 22 (same as the library) and Go 1.27+:

```shell
go run ./examples/opt
```

## Expected output

```text
insts: 7 -> 1
```

`mem2reg` promotes the alloca away (which also makes the dead store
disappear), and `instcombine` folds the `+ 0` / `* 1` identities, so all that
remains is a constant return (function attributes omitted here, the module
carries the usual ones):

```llvm
define noundef i32 @compute() local_unnamed_addr #0 {
entry:
  ret i32 2
}
```

The example panics if building the module or running the pipeline fails; on
the normal path it never panics.
