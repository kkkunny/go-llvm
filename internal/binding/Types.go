package binding

/*
#include "llvm-c/Core.h"
*/
import "C"

type (
	// LLVMContextRef The top-level container for all LLVM global data. See the LLVMContext class.
	LLVMContextRef struct{ c C.LLVMContextRef }

	// LLVMModuleRef The top-level container for all other LLVM Intermediate Representation (IR) objects.
	LLVMModuleRef struct{ c C.LLVMModuleRef }

	// LLVMTypeRef Each value in the LLVM IR has a type, an LLVMTypeRef.
	LLVMTypeRef struct{ c C.LLVMTypeRef }

	// LLVMValueRef Represents an individual value in LLVM IR.
	LLVMValueRef struct{ c C.LLVMValueRef }

	// LLVMBasicBlockRef Represents a basic block of instructions in LLVM IR.
	LLVMBasicBlockRef struct{ c C.LLVMBasicBlockRef }

	// LLVMBuilderRef Represents an LLVM basic block builder.
	LLVMBuilderRef struct{ c C.LLVMBuilderRef }
)
