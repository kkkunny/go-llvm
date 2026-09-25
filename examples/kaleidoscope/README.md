# Kaleidoscope

An implementation of the Kaleidoscope language from the official LLVM tutorial
[My First Language Frontend](https://llvm.org/docs/tutorial/MyFirstLanguageFrontend/index.html)
on top of [go-llvm](../../), covering the tutorial up to chapter 7 (JIT and
optimization), modeled on
[the inkwell example of the same name](https://github.com/TheDan64/inkwell/tree/master/examples/kaleidoscope).

## Running

Requires LLVM 22 (same requirement as the library), then:

```shell
go run ./examples/kaleidoscope
```

You can also evaluate a single expression and exit:

```shell
go run ./examples/kaleidoscope -e '1 + 2 * 2'
```

## Sample session

```text
Kaleidoscope REPL（输入 exit 或 quit 退出）
?> 1 + 1
=> 2
?> var a = 5, b = 10 in a * b
=> 50
?> def fib(n) if n < 2 then n else fib(n - 1) + fib(n - 2)
?> fib(40)
=> 102334155
?> extern putchard(x)
?> for i = 1, i < 5 in putchard(42)
*****=> 0
?> exit
```

(The REPL banner is printed in Chinese, meaning "enter exit or quit to
leave".)

The value of a `for` expression is always 0; the sample uses `putchard` to
print asterisks so that the number of iterations is visible.

## Command-line flags

| Flag | Meaning |
|---|---|
| `-e <expr>` | Evaluate a single expression and exit |
| `-dl` | Dump the lexer output (token stream) |
| `-dp` | Dump the parser output (AST) |
| `-dc` | Dump the generated LLVM IR |

For example:

```shell
$ go run ./examples/kaleidoscope -dl -dp -dc -e '1 + 2 * 2'
-> tokens: [Number(1) Op('+') Number(2) Op('*') Number(2) EOF]
-> expression: (1 + (2 * 2))
-> IR:
define double @__anon_expr.1() {
entry:
  ret double 5.000000e+00
}

=> 5
```

## Language features (tutorial chapters)

| Tutorial chapter | Feature | Example |
|---|---|---|
| 1 | Lexing and the REPL | `1 + 1` |
| 2 | AST and precedence parsing | `1 + 2 * 3` |
| 3 | Codegen for numbers and binary operations | `10 / 4` |
| 4 | Functions, `extern`, calls, `if/then/else` | `def fib(n) if n < 2 then n else fib(n - 1) + fib(n - 2)` |
| 5 | `for/in` loops, custom `binary`/`unary` operators, assignment | `def binary^ 40 (a, b) a * b` |
| 7 | Top-level expression JIT evaluation + optimization pipeline | `fib(40)` |

Scoping with `var a = 1, b in ...` is supported as well.

## Implementation notes

| File | Responsibility |
|---|---|
| `lexer.go` | Lexing |
| `parser.go` | Recursive descent + precedence climbing |
| `compiler.go` | AST → `ir.Module` (Alloca/Load/Store, PHI, CondBr, custom operator calls, ...) |
| `session.go` | Session state: the long-lived `jit.LLJIT`, historical definitions, optimization pipeline |
| `main.go` | Command line and REPL |

- Every evaluation creates a fresh `Context` + `ir.Module` and recompiles the
  historical definitions together with the new function (matching inkwell);
  after `pass.RunPasses("instcombine,reassociate,gvn,simplifycfg,mem2reg")`,
  the module is added to the JIT via `jit.ResourceTracker` and `Remove()`d as
  soon as evaluation finishes (corresponding to `addModule`/`removeModule` in
  the tutorial's LangImpl07).
- Declarations such as `extern sin(x)` are resolved to host process symbols by
  `jit.LLJIT.AddProcessSymbols()`; `putchard`/`printd` are registered as Go
  callbacks through `jit.LLJIT.MapFunc`.
- One deliberate difference from inkwell: redefining a `def` with the same
  name replaces the historical definition, whereas inkwell leaves the old
  function behind in the module.
- Operators not declared with `binary <op> <prec>` cannot appear in binary
  position (matching the official tutorial; inkwell gives unregistered
  operators a default precedence).

## Tests

```shell
go test ./examples/kaleidoscope
go test -tags=llvm_release ./examples/kaleidoscope
```
