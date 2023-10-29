# go-llvm

This library provides bindings to a system-installed LLVM.

Currently supported:

* LLVM 15
* LLVM 16
* LLVM 17

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