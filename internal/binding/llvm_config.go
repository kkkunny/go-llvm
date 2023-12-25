package binding

/*
#include "llvm/Config/llvm-config.h"
*/
import "C"

// LLVM_DEFAULT_TARGET_TRIPLE Target triple LLVM will generate code for by default
const LLVM_DEFAULT_TARGET_TRIPLE string = C.LLVM_DEFAULT_TARGET_TRIPLE

// LLVM_HOST_TRIPLE Host triple LLVM will be executed on
const LLVM_HOST_TRIPLE string = C.LLVM_HOST_TRIPLE

// LLVM_VERSION_MAJOR Major version of the LLVM API
const LLVM_VERSION_MAJOR int32 = C.LLVM_VERSION_MAJOR

// LLVM_VERSION_MINOR Minor version of the LLVM API
const LLVM_VERSION_MINOR int32 = C.LLVM_VERSION_MINOR

// LLVM_VERSION_PATCH Patch version of the LLVM API
const LLVM_VERSION_PATCH int32 = C.LLVM_VERSION_PATCH

// LLVM_VERSION_STRING LLVM version string
const LLVM_VERSION_STRING string = C.LLVM_VERSION_STRING
