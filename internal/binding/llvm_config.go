package binding

/*
#include "llvm/Config/llvm-config.h"
*/
import "C"

const Version string = C.LLVM_VERSION_STRING
