package binding

/*
#include "llvm-c/Comdat.h"
*/
import "C"

// LLVMComdatSelectionKind is the comdat conflict resolution kind.
type LLVMComdatSelectionKind int32

const (
	LLVMAnyComdatSelectionKind           LLVMComdatSelectionKind = C.LLVMAnyComdatSelectionKind
	LLVMExactMatchComdatSelectionKind    LLVMComdatSelectionKind = C.LLVMExactMatchComdatSelectionKind
	LLVMLargestComdatSelectionKind       LLVMComdatSelectionKind = C.LLVMLargestComdatSelectionKind
	LLVMNoDeduplicateComdatSelectionKind LLVMComdatSelectionKind = C.LLVMNoDeduplicateComdatSelectionKind
	LLVMSameSizeComdatSelectionKind      LLVMComdatSelectionKind = C.LLVMSameSizeComdatSelectionKind
)

// LLVMGetOrInsertComdat Return the Comdat in the module with the specified name; creates it if it does not exist.
func LLVMGetOrInsertComdat(m LLVMModuleRef, name string) LLVMComdatRef {
	return string2CString(name, func(name *C.char) LLVMComdatRef {
		return LLVMComdatRef{c: C.LLVMGetOrInsertComdat(m.c, name)}
	})
}

// LLVMGetComdat Get the Comdat assigned to the given global object; returns nil if none is set.
func LLVMGetComdat(v LLVMValueRef) LLVMComdatRef {
	return LLVMComdatRef{c: C.LLVMGetComdat(v.c)}
}

// LLVMSetComdat Assign the Comdat to the given global object.
func LLVMSetComdat(v LLVMValueRef, c LLVMComdatRef) {
	C.LLVMSetComdat(v.c, c.c)
}

// LLVMGetComdatSelectionKind Get the conflict resolution selection kind for the Comdat.
func LLVMGetComdatSelectionKind(c LLVMComdatRef) LLVMComdatSelectionKind {
	return LLVMComdatSelectionKind(C.LLVMGetComdatSelectionKind(c.c))
}

// LLVMSetComdatSelectionKind Set the conflict resolution selection kind for the Comdat.
func LLVMSetComdatSelectionKind(c LLVMComdatRef, kind LLVMComdatSelectionKind) {
	C.LLVMSetComdatSelectionKind(c.c, C.LLVMComdatSelectionKind(kind))
}
