package binding

/*
#include "llvm-c/Comdat.h"
*/
import "C"

// LLVMComdatSelectionKind Comdat 冲突解决方式
type LLVMComdatSelectionKind int32

const (
	LLVMAnyComdatSelectionKind           LLVMComdatSelectionKind = C.LLVMAnyComdatSelectionKind
	LLVMExactMatchComdatSelectionKind    LLVMComdatSelectionKind = C.LLVMExactMatchComdatSelectionKind
	LLVMLargestComdatSelectionKind       LLVMComdatSelectionKind = C.LLVMLargestComdatSelectionKind
	LLVMNoDeduplicateComdatSelectionKind LLVMComdatSelectionKind = C.LLVMNoDeduplicateComdatSelectionKind
	LLVMSameSizeComdatSelectionKind      LLVMComdatSelectionKind = C.LLVMSameSizeComdatSelectionKind
)

// LLVMGetOrInsertComdat Return the Comdat in the module with the specified name; 不存在则创建。
func LLVMGetOrInsertComdat(m LLVMModuleRef, name string) LLVMComdatRef {
	return string2CString(name, func(name *C.char) LLVMComdatRef {
		return LLVMComdatRef{c: C.LLVMGetOrInsertComdat(m.c, name)}
	})
}

// LLVMGetComdat Get the Comdat assigned to the given global object; 未设置返回 nil。
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
