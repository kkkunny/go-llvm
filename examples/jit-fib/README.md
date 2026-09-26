# jit-fib

A runnable go-llvm JIT example:

1. builds a recursive `fib` function directly in LLVM IR (`target.InitNative`
   + `jit.NewLLJIT`), and calls it from Go via `Func[func(int32) int32]`;
2. registers a host Go callback `host_double` with `MapFunc`, then calls it
   from JIT code through an automatically generated IR wrapper
   (`double_fib(n) = host_double(fib(n))`).

## Run

Requires LLVM 22 (same as the library) and Go 1.27+:

```shell
go run ./examples/jit-fib
```

## Expected output

```text
fib(10) = 55
double_fib(10) = host_double(fib(10)) = 110
```

The example never panics on the normal path; errors (target/JIT setup, module
verification) are reported by panicking with the underlying error.
