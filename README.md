# go-llvm

This library provides bindings to a system-installed LLVM.

Currently supported:

* LLVM 15
* LLVM 16
* LLVM 17
* LLVM 18

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
make config EXPECT_VERSION=VERION OF LLVM
# eg.make config EXPECT_VERSION=15
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
