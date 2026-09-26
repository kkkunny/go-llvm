// This file holds portable #cgo flag candidates for the common Linux, macOS and FreeBSD
// LLVM layouts (the checked-in version). `go generate ./internal/binding` or `make config`
// rewrites it with machine-specific flags queried from the local llvm-config; keep the
// checked-in copy portable and review the diff before committing a generated one.
//
// CFLAGS and CXXFLAGS must always carry the same -I paths: the C++ shims include LLVM
// headers too, and a C-only include path was the root cause of a past CI failure.
// Nonexistent -I/-L directories are ignored by the compiler and linker, so several
// candidate layouts can be listed at once.

package binding

/*
#cgo CFLAGS: -D_GNU_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D__STDC_LIMIT_MACROS
#cgo CXXFLAGS: -std=c++17 -fexceptions -D_GNU_SOURCE -D_GLIBCXX_USE_CXX11_ABI=1 -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D__STDC_LIMIT_MACROS

#cgo linux CFLAGS: -I/usr/lib/llvm-22/include -I/usr/include/llvm-22 -I/usr/include/llvm-c-22 -I/usr/include -I/usr/local/include
#cgo linux CXXFLAGS: -I/usr/lib/llvm-22/include -I/usr/include/llvm-22 -I/usr/include/llvm-c-22 -I/usr/include -I/usr/local/include
#cgo linux LDFLAGS: -L/usr/lib/llvm-22/lib -L/usr/lib64/llvm-22/lib -L/usr/lib64/llvm22/lib -L/usr/lib -L/usr/lib64 -L/usr/local/lib -lLLVM

#cgo darwin,arm64 CFLAGS: -I/opt/homebrew/opt/llvm@22/include
#cgo darwin,arm64 CXXFLAGS: -I/opt/homebrew/opt/llvm@22/include
#cgo darwin,arm64 LDFLAGS: -L/opt/homebrew/opt/llvm@22/lib -Wl,-search_paths_first -Wl,-headerpad_max_install_names -lLLVM -lz -lm

#cgo darwin,amd64 CFLAGS: -I/usr/local/opt/llvm@22/include
#cgo darwin,amd64 CXXFLAGS: -I/usr/local/opt/llvm@22/include
#cgo darwin,amd64 LDFLAGS: -L/usr/local/opt/llvm@22/lib -Wl,-search_paths_first -Wl,-headerpad_max_install_names -lLLVM -lz -lm

#cgo freebsd CFLAGS: -I/usr/local/llvm22/include -I/usr/local/llvm22/include/llvm-c
#cgo freebsd CXXFLAGS: -I/usr/local/llvm22/include -I/usr/local/llvm22/include/llvm-c
#cgo freebsd LDFLAGS: -L/usr/local/llvm22/lib -lLLVM
*/
import "C"
