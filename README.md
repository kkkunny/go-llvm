# go-llvm

This library provides bindings to a system-installed LLVM.

Currently supported:

* LLVM 20

## Usage

First, you need to make sure that you have installed a supported version of LLVM.

And then

```shell
go get github.com/kkkunny/go-llvm
```

```shell
curl -O https://raw.githubusercontent.com/kkkunny/go-llvm/master/Makefile
```

```shell
make config
# or specify a version: make config LLVM_CONFIG_BIN=20
```

```go
package main

import (
	"os"

	"github.com/kkkunny/go-llvm"
)

func main() {
	ctx := llvm.NewContext()
	module := ctx.NewModule("main")
	builder := ctx.NewBuilder()

	mainFn := module.NewFunction("main", ctx.FunctionType(false, ctx.IntegerType(8)))
	mainFnEntry := mainFn.NewBlock("entry")
	builder.MoveToAfter(mainFnEntry)
	var ret llvm.Value = ctx.ConstInteger(ctx.IntegerType(8), 0)
	builder.CreateRet(&ret)

	_ = llvm.InitializeNativeTarget()
	_ = llvm.InitializeNativeAsmPrinter()

	jiter, err := llvm.NewJITCompiler(module, llvm.CodeOptLevelNone)
	if err != nil {
		panic(err)
	}
	os.Exit(int(jiter.RunMainFunction(mainFn, nil, nil)))
}
```
