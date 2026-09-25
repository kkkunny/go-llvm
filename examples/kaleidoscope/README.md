# Kaleidoscope

用 [go-llvm](../../) 实现 LLVM 官方教程
[My First Language Frontend](https://llvm.org/docs/tutorial/MyFirstLanguageFrontend/index.html)
的 Kaleidoscope 语言（功能覆盖到第 7 章：JIT 与优化），参考
[inkwell 的同名示例](https://github.com/TheDan64/inkwell/tree/master/examples/kaleidoscope)。

## 运行

需要已安装 LLVM 22（与库的要求一致），然后：

```shell
go run ./examples/kaleidoscope
```

也可以直接求值单个表达式后退出：

```shell
go run ./examples/kaleidoscope -e '1 + 2 * 2'
```

## 示例会话

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

`for` 表达式的值恒为 0，示例里用 `putchard` 打印星号观察循环次数。

## 命令行参数

| 参数 | 含义 |
|---|---|
| `-e <expr>` | 求值单个表达式后退出 |
| `-dl` | 显示词法分析结果（token 流） |
| `-dp` | 显示语法分析结果（AST） |
| `-dc` | 显示生成的 LLVM IR |

例如：

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

## 语言特性（对应教程章节）

| 教程章节 | 特性 | 示例 |
|---|---|---|
| 1 | 词法分析与 REPL | `1 + 1` |
| 2 | AST 与优先级解析 | `1 + 2 * 3` |
| 3 | 数字与二元运算代码生成 | `10 / 4` |
| 4 | 函数、`extern`、调用、`if/then/else` | `def fib(n) if n < 2 then n else fib(n - 1) + fib(n - 2)` |
| 5 | `for/in` 循环、自定义 `binary`/`unary` 算符、赋值 | `def binary^ 40 (a, b) a * b` |
| 7 | 顶层表达式 JIT 求值 + 优化管线 | `fib(40)` |

`var a = 1, b in ...` 作用域同样支持。

## 实现说明

| 文件 | 职责 |
|---|---|
| `lexer.go` | 词法分析 |
| `parser.go` | 递归下降 + 优先级爬升解析 |
| `compiler.go` | AST → `ir.Module`（Alloca/Load/Store、PHI、CondBr、自定义算符调用……） |
| `session.go` | 会话状态：常驻 `jit.LLJIT`、历史定义、优化管线 |
| `main.go` | 命令行与 REPL |

- 每次求值新建 `Context` + `ir.Module`，把历史定义和新函数一起重编译（对齐 inkwell），
  经 `pass.RunPasses("instcombine,reassociate,gvn,simplifycfg,mem2reg")` 优化后，
  用 `jit.ResourceTracker` 加入 JIT，求值完立即 `Remove()`（对应教程 LangImpl07 的
  `addModule`/`removeModule`）。
- `extern sin(x)` 之类声明由 `jit.LLJIT.AddProcessSymbols()` 解析到宿主进程符号；
  `putchard`/`printd` 通过 `jit.LLJIT.MapFunc` 注册为 Go 回调。
- 与 inkwell 的一处有意差异：同名 `def` 重定义会替换历史定义，而 inkwell 会在模块里
  残留旧函数。
- 未通过 `binary <op> <prec>` 声明的算符不能出现在二元位置（对齐官方教程；inkwell 会给
  未注册算符默认优先级）。

## 测试

```shell
go test ./examples/kaleidoscope
go test -tags=llvm_release ./examples/kaleidoscope
```
