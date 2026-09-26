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
	LLVMBuilderRef     struct{ c C.LLVMBuilderRef }
	LLVMPassManagerRef struct{ c C.LLVMPassManagerRef }

	LLVMAttributeRef      struct{ c C.LLVMAttributeRef }
	LLVMDiagnosticInfoRef struct{ c C.LLVMDiagnosticInfoRef }
	LLVMMetadataRef       struct{ c C.LLVMMetadataRef }
	LLVMNamedMDNodeRef    struct{ c C.LLVMNamedMDNodeRef }
	LLVMComdatRef         struct{ c C.LLVMComdatRef }

	// LLVMMemoryBufferRef A memory buffer.
	LLVMMemoryBufferRef struct{ c C.LLVMMemoryBufferRef }

	// LLVMUseRef A reference to a use of a value.
	LLVMUseRef struct{ c C.LLVMUseRef }

	// LLVMOperandBundleRef A reference to an operand bundle.
	LLVMOperandBundleRef struct{ c C.LLVMOperandBundleRef }
)

func (ref LLVMContextRef) IsNil() bool        { return ref.c == nil }
func (ref LLVMModuleRef) IsNil() bool         { return ref.c == nil }
func (ref LLVMTypeRef) IsNil() bool           { return ref.c == nil }
func (ref LLVMValueRef) IsNil() bool          { return ref.c == nil }
func (ref LLVMBasicBlockRef) IsNil() bool     { return ref.c == nil }
func (ref LLVMBuilderRef) IsNil() bool        { return ref.c == nil }
func (ref LLVMAttributeRef) IsNil() bool      { return ref.c == nil }
func (ref LLVMDiagnosticInfoRef) IsNil() bool { return ref.c == nil }
func (ref LLVMMetadataRef) IsNil() bool       { return ref.c == nil }
func (ref LLVMNamedMDNodeRef) IsNil() bool    { return ref.c == nil }
func (ref LLVMComdatRef) IsNil() bool         { return ref.c == nil }
func (ref LLVMMemoryBufferRef) IsNil() bool   { return ref.c == nil }
func (ref LLVMUseRef) IsNil() bool            { return ref.c == nil }
func (ref LLVMOperandBundleRef) IsNil() bool  { return ref.c == nil }

// Equal compares handle equality.
func (ref LLVMContextRef) Equal(other LLVMContextRef) bool { return ref.c == other.c }

// Equal compares handle equality.
func (ref LLVMModuleRef) Equal(other LLVMModuleRef) bool { return ref.c == other.c }

// Equal compares handle equality.
func (ref LLVMTypeRef) Equal(other LLVMTypeRef) bool { return ref.c == other.c }

// Equal compares handle equality.
func (ref LLVMValueRef) Equal(other LLVMValueRef) bool { return ref.c == other.c }

// Equal compares handle equality.
func (ref LLVMBasicBlockRef) Equal(other LLVMBasicBlockRef) bool { return ref.c == other.c }

// Equal compares handle equality.
func (ref LLVMBuilderRef) Equal(other LLVMBuilderRef) bool { return ref.c == other.c }
