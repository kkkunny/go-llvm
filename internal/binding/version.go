// Compile-time version gate for the LLVM major this build targets (see AGENTS.md,
// "Version tags"). cgo.go is rewritten by the llvmconfig generator, so the gate lives
// in its own file and must stay here.
//
// Untagged builds target the latest supported major; -tags llvmNN pins major NN. A
// header/tag mismatch fails here, at compile time, in every build mode (including
// -tags=llvm_release); the runtime header/library check in NewContext stays debug-only.
//
// Only LLVM 23 is supported today, so the gate below is hardcoded; when a second major
// is supported, split this file per tag (version_llvmNN.go) as described in AGENTS.md.

package binding

/*
#include "llvm/Config/llvm-config.h"
#if LLVM_VERSION_MAJOR != 23
#error "go-llvm: this build targets LLVM 23; install LLVM 23 (see the README support table)"
#endif
*/
import "C"
