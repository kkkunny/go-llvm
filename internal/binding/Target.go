package binding

/*
#include "llvm-c/Target.h"
*/
import "C"
import "errors"

type LLVMTargetDataRef struct{ c C.LLVMTargetDataRef }

func (ref LLVMTargetDataRef) IsNil() bool { return ref.c == nil }

// LLVMInitializeAllTargetInfos The main program should call this function if it wants access to all available targets that LLVM is configured to support.
func LLVMInitializeAllTargetInfos() {
	C.LLVMInitializeAllTargetInfos()
}

// LLVMInitializeAllTargets The main program should call this function if it wants to link in all available targets that LLVM is configured to support.
func LLVMInitializeAllTargets() {
	C.LLVMInitializeAllTargets()
}

// LLVMInitializeAllTargetMCs The main program should call this function if it wants access to all available target MC that LLVM is configured to support.
func LLVMInitializeAllTargetMCs() {
	C.LLVMInitializeAllTargetMCs()
}

// LLVMInitializeAllAsmPrinters The main program should call this function if it wants all asm printers that LLVM is configured to support, to make them available via the TargetRegistry.
func LLVMInitializeAllAsmPrinters() {
	C.LLVMInitializeAllAsmPrinters()
}

// LLVMInitializeAllAsmParsers The main program should call this function if it wants all asm parsers that LLVM is configured to support, to make them available via the TargetRegistry.
func LLVMInitializeAllAsmParsers() {
	C.LLVMInitializeAllAsmParsers()
}

// LLVMInitializeAllDisassemblers The main program should call this function if it wants all disassemblers that LLVM is configured to support, to make them available via the TargetRegistry.
func LLVMInitializeAllDisassemblers() {
	C.LLVMInitializeAllDisassemblers()
}

// LLVMInitializeNativeTarget The main program should call this function to initialize the native target corresponding to the host.  This is useful for JIT applications to ensure that the target gets linked in correctly.
func LLVMInitializeNativeTarget() error {
	if llvmBool2bool(C.LLVMInitializeNativeTarget()) {
		return errors.New("failed to initialize native target")
	}
	return nil
}

// LLVMInitializeNativeAsmParser The main program should call this function to initialize the parser for the native target corresponding to the host.
func LLVMInitializeNativeAsmParser() error {
	if llvmBool2bool(C.LLVMInitializeNativeAsmParser()) {
		return errors.New("failed to initialize native asm parser")
	}
	return nil
}

// LLVMInitializeNativeAsmPrinter The main program should call this function to initialize the printer for the native target corresponding to the host.
func LLVMInitializeNativeAsmPrinter() error {
	if llvmBool2bool(C.LLVMInitializeNativeAsmPrinter()) {
		return errors.New("failed to initialize native asm printer")
	}
	return nil
}

// LLVMInitializeNativeDisassembler The main program should call this function to initialize the disassembler for the native target corresponding to the host.
func LLVMInitializeNativeDisassembler() error {
	if llvmBool2bool(C.LLVMInitializeNativeDisassembler()) {
		return errors.New("failed to initialize native disassembler")
	}
	return nil
}

// LLVMGetModuleDataLayout Obtain the data layout for a module.
func LLVMGetModuleDataLayout(m LLVMModuleRef) LLVMTargetDataRef {
	return LLVMTargetDataRef{c: C.LLVMGetModuleDataLayout(m.c)}
}

// LLVMSetModuleDataLayout Set the data layout for a module.
func LLVMSetModuleDataLayout(m LLVMModuleRef, dl LLVMTargetDataRef) {
	C.LLVMSetModuleDataLayout(m.c, dl.c)
}

// LLVMCreateTargetData Creates target data from a target layout string.
func LLVMCreateTargetData(stringRep string) LLVMTargetDataRef {
	return string2CString(stringRep, func(stringRep *C.char) LLVMTargetDataRef {
		return LLVMTargetDataRef{c: C.LLVMCreateTargetData(stringRep)}
	})
}

// LLVMDisposeTargetData Deallocates a TargetData.
func LLVMDisposeTargetData(td LLVMTargetDataRef) {
	C.LLVMDisposeTargetData(td.c)
}

// LLVMCopyStringRepOfTargetData Converts target data to a target layout string.
func LLVMCopyStringRepOfTargetData(td LLVMTargetDataRef) string {
	cstring := C.LLVMCopyStringRepOfTargetData(td.c)
	defer LLVMDisposeMessage(cstring)
	return C.GoString(cstring)
}

// LLVMPointerSize Returns the pointer size in bytes for a target.
func LLVMPointerSize(td LLVMTargetDataRef) uint32 {
	return uint32(C.LLVMPointerSize(td.c))
}

// LLVMPointerSizeForAS Returns the pointer size in bytes for a target for a specified address space.
func LLVMPointerSizeForAS(td LLVMTargetDataRef, as uint32) uint32 {
	return uint32(C.LLVMPointerSizeForAS(td.c, C.unsigned(as)))
}

// LLVMIntPtrType Returns the integer type that is the same size as a pointer on a target.
func LLVMIntPtrType(td LLVMTargetDataRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMIntPtrType(td.c)}
}

// LLVMIntPtrTypeForAS Returns the integer type that is the same size as a pointer on a target.
func LLVMIntPtrTypeForAS(td LLVMTargetDataRef, as uint32) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMIntPtrTypeForAS(td.c, C.unsigned(as))}
}

// LLVMIntPtrTypeInContext Returns the integer type that is the same size as a pointer on a target.
func LLVMIntPtrTypeInContext(c LLVMContextRef, td LLVMTargetDataRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMIntPtrTypeInContext(c.c, td.c)}
}

// LLVMIntPtrTypeForASInContext Returns the integer type that is the same size as a pointer on a target.
func LLVMIntPtrTypeForASInContext(c LLVMContextRef, td LLVMTargetDataRef, as uint32) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMIntPtrTypeForASInContext(c.c, td.c, C.unsigned(as))}
}

// LLVMSizeOfTypeInBits Computes the size of a type in bytes for a target.
func LLVMSizeOfTypeInBits(td LLVMTargetDataRef, ty LLVMTypeRef) uint64 {
	return uint64(C.LLVMSizeOfTypeInBits(td.c, ty.c))
}

// LLVMStoreSizeOfType Computes the storage size of a type in bytes for a target.
func LLVMStoreSizeOfType(td LLVMTargetDataRef, ty LLVMTypeRef) uint64 {
	return uint64(C.LLVMStoreSizeOfType(td.c, ty.c))
}

// LLVMABISizeOfType Computes the ABI size of a type in bytes for a target.
func LLVMABISizeOfType(td LLVMTargetDataRef, ty LLVMTypeRef) uint64 {
	return uint64(C.LLVMABISizeOfType(td.c, ty.c))
}

// LLVMABIAlignmentOfType Computes the ABI alignment of a type in bytes for a target.
func LLVMABIAlignmentOfType(td LLVMTargetDataRef, ty LLVMTypeRef) uint32 {
	return uint32(C.LLVMABIAlignmentOfType(td.c, ty.c))
}

// LLVMCallFrameAlignmentOfType Computes the call frame alignment of a type in bytes for a target.
func LLVMCallFrameAlignmentOfType(td LLVMTargetDataRef, ty LLVMTypeRef) uint32 {
	return uint32(C.LLVMCallFrameAlignmentOfType(td.c, ty.c))
}

// LLVMPreferredAlignmentOfType Computes the preferred alignment of a type in bytes for a target.
func LLVMPreferredAlignmentOfType(td LLVMTargetDataRef, ty LLVMTypeRef) uint32 {
	return uint32(C.LLVMPreferredAlignmentOfType(td.c, ty.c))
}

// LLVMPreferredAlignmentOfGlobal Computes the preferred alignment of a global variable in bytes for a target.
func LLVMPreferredAlignmentOfGlobal(td LLVMTargetDataRef, globalVar LLVMValueRef) uint32 {
	return uint32(C.LLVMPreferredAlignmentOfGlobal(td.c, globalVar.c))
}

// LLVMElementAtOffset Computes the structure element that contains the byte offset for a target.
func LLVMElementAtOffset(td LLVMTargetDataRef, structTy LLVMTypeRef, offset uint64) uint32 {
	return uint32(C.LLVMElementAtOffset(td.c, structTy.c, C.ulonglong(offset)))
}

// LLVMOffsetOfElement Computes the byte offset of the indexed struct element for a target.
func LLVMOffsetOfElement(td LLVMTargetDataRef, structTy LLVMTypeRef, element uint32) uint64 {
	return uint64(C.LLVMOffsetOfElement(td.c, structTy.c, C.unsigned(element)))
}
