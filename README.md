# go-llvm

This library provides bindings to a system-installed LLVM.

Currently supported:

| Local LLVM | How to use |
|---|---|
| 22 (latest, default branch) | `go get github.com/kkkunny/go-llvm` |
| 21 | `go get github.com/kkkunny/go-llvm@llvm21` |

Notes:

* These are non-semver tags, so `go.mod` records them as pseudo-versions.
* Older LLVM lines (20 and earlier) are no longer provided; they can be pinned to
  historic commits if needed.

## Usage

Install a supported LLVM together with its development headers (for example
`llvm-22-dev`), then:

```shell
go get github.com/kkkunny/go-llvm
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

### Non-standard LLVM prefixes

`internal/binding/cgo.go` ships `#cgo` flags for the common Linux layouts
(`/usr/lib/llvm-NN`, `/usr`, `/usr/local`, `/usr/lib64`) with an unversioned
`-lLLVM`. If your LLVM lives somewhere else, override via environment
variables in your build:

```shell
export CGO_CFLAGS="$(llvm-config --cflags)"
export CGO_CXXFLAGS="$(llvm-config --cxxflags)"
export CGO_LDFLAGS="$(llvm-config --ldflags --libs)"
```

or generate an override file in **your own** main package (never in this
repository's root):

```shell
curl -O https://raw.githubusercontent.com/kkkunny/go-llvm/master/Makefile
make config             # or pin the toolchain: make config MAJOR_VERSION=22
```

## Updating for a new LLVM release

1. Read the new release notes at
   `https://releases.llvm.org/N.1.0/docs/ReleaseNotes.html`, section
   **Changes to the C API** (removals, deprecations, behavior changes).
2. Cross-check every C symbol this repo references against the local headers:
   `rg -o 'C\.[A-Za-z_]\w*' --glob '*.go'`, then verify each name with
   `grep -rw NAME /usr/include/llvm-c/`.
3. Freeze the previous line first, then bump the four version spots on master:
   `internal/binding/cgo.go` candidate dirs, `Makefile`
   `MIN/MAX_SUPPORT_MAJOR_VERSION`, the support table above, and this README.

   ```shell
   git branch llvm-NN master     # keep the old line reachable
   git tag -a llvmNN -m "LLVM NN line"
   ```
4. Enum constants need no changes: they are declared as `C.Name` and bind to
   whatever the local headers define. Unknown value kinds / opcodes degrade to
   generic fallback wrappers instead of panicking.
5. Bind newly added C APIs only when needed.
