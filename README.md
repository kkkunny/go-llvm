# go-llvm

This library provides bindings to a system-installed LLVM.

Currently supported:

* LLVM 15(check in linux/amd64).
* LLVM 17(check in windows/amd64).

## Usage

If you have a supported LLVM installation, you should be able to do a simple `go get`:

    go get github.com/kkkunny/go-llvm@llvm17