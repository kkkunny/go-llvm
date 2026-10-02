// Compile-time version gate for the LLVM major this build targets (see AGENTS.md,
// "Version tags"). cgo.go is rewritten by the llvmconfig generator, so the gate lives
// in its own file and must stay here.
//
// Untagged builds target the latest supported major; -tags llvmNN pins major NN. A
// header/tag mismatch fails here, at compile time, in every build mode (including
// -tags=llvm_release); the runtime header/library check in NewContext stays debug-only.

package binding

func init() {
	if LLVM_VERSION_MAJOR != 23 {
		panic("go-llvm: this build targets LLVM 23; install LLVM 23, or pass the -tags llvmNN matching your LLVM major (see the README support table)")
	}
}
