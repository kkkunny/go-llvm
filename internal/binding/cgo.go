package binding

/*
#cgo CFLAGS: -D_GNU_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D__STDC_LIMIT_MACROS
#cgo CFLAGS: -I/usr/lib/llvm-22/include -I/usr/include -I/usr/local/include
#cgo CXXFLAGS: -std=c++17 -fexceptions -D_GNU_SOURCE -D_GLIBCXX_USE_CXX11_ABI=1 -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D__STDC_LIMIT_MACROS
#cgo LDFLAGS: -L/usr/lib/llvm-22/lib -L/usr/lib64 -L/usr/lib -L/usr/local/lib -lLLVM
*/
import "C"
