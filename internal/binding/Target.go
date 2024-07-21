package binding

/*
#include "llvm-c/Target.h"
*/
import "C"
import "errors"

type LLVMTargetDataRef struct{ c C.LLVMTargetDataRef }

type LLVMByteOrdering int32

const (
	LLVMBigEndian    LLVMByteOrdering = C.LLVMBigEndian
	LLVMLittleEndian LLVMByteOrdering = C.LLVMLittleEndian
)

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

func LLVMInitializeAArch64TargetInfo()     { C.LLVMInitializeAArch64TargetInfo() }
func LLVMInitializeAMDGPUTargetInfo()      { C.LLVMInitializeAMDGPUTargetInfo() }
func LLVMInitializeARMTargetInfo()         { C.LLVMInitializeARMTargetInfo() }
func LLVMInitializeAVRTargetInfo()         { C.LLVMInitializeAVRTargetInfo() }
func LLVMInitializeBPFTargetInfo()         { C.LLVMInitializeBPFTargetInfo() }
func LLVMInitializeHexagonTargetInfo()     { C.LLVMInitializeHexagonTargetInfo() }
func LLVMInitializeLanaiTargetInfo()       { C.LLVMInitializeLanaiTargetInfo() }
func LLVMInitializeLoongArchTargetInfo()   { C.LLVMInitializeLoongArchTargetInfo() }
func LLVMInitializeMipsTargetInfo()        { C.LLVMInitializeMipsTargetInfo() }
func LLVMInitializeMSP430TargetInfo()      { C.LLVMInitializeMSP430TargetInfo() }
func LLVMInitializeNVPTXTargetInfo()       { C.LLVMInitializeNVPTXTargetInfo() }
func LLVMInitializePowerPCTargetInfo()     { C.LLVMInitializePowerPCTargetInfo() }
func LLVMInitializeRISCVTargetInfo()       { C.LLVMInitializeRISCVTargetInfo() }
func LLVMInitializeSparcTargetInfo()       { C.LLVMInitializeSparcTargetInfo() }
func LLVMInitializeSystemZTargetInfo()     { C.LLVMInitializeSystemZTargetInfo() }
func LLVMInitializeVETargetInfo()          { C.LLVMInitializeVETargetInfo() }
func LLVMInitializeWebAssemblyTargetInfo() { C.LLVMInitializeWebAssemblyTargetInfo() }
func LLVMInitializeX86TargetInfo()         { C.LLVMInitializeX86TargetInfo() }
func LLVMInitializeXCoreTargetInfo()       { C.LLVMInitializeXCoreTargetInfo() }

func LLVMInitializeAArch64Target()     { C.LLVMInitializeAArch64Target() }
func LLVMInitializeAMDGPUTarget()      { C.LLVMInitializeAMDGPUTarget() }
func LLVMInitializeARMTarget()         { C.LLVMInitializeARMTarget() }
func LLVMInitializeAVRTarget()         { C.LLVMInitializeAVRTarget() }
func LLVMInitializeBPFTarget()         { C.LLVMInitializeBPFTarget() }
func LLVMInitializeHexagonTarget()     { C.LLVMInitializeHexagonTarget() }
func LLVMInitializeLanaiTarget()       { C.LLVMInitializeLanaiTarget() }
func LLVMInitializeLoongArchTarget()   { C.LLVMInitializeLoongArchTarget() }
func LLVMInitializeMipsTarget()        { C.LLVMInitializeMipsTarget() }
func LLVMInitializeMSP430Target()      { C.LLVMInitializeMSP430Target() }
func LLVMInitializeNVPTXTarget()       { C.LLVMInitializeNVPTXTarget() }
func LLVMInitializePowerPCTarget()     { C.LLVMInitializePowerPCTarget() }
func LLVMInitializeRISCVTarget()       { C.LLVMInitializeRISCVTarget() }
func LLVMInitializeSparcTarget()       { C.LLVMInitializeSparcTarget() }
func LLVMInitializeSystemZTarget()     { C.LLVMInitializeSystemZTarget() }
func LLVMInitializeVETarget()          { C.LLVMInitializeVETarget() }
func LLVMInitializeWebAssemblyTarget() { C.LLVMInitializeWebAssemblyTarget() }
func LLVMInitializeX86Target()         { C.LLVMInitializeX86Target() }
func LLVMInitializeXCoreTarget()       { C.LLVMInitializeXCoreTarget() }

func LLVMInitializeAArch64TargetMC()     { C.LLVMInitializeAArch64TargetMC() }
func LLVMInitializeAMDGPUTargetMC()      { C.LLVMInitializeAMDGPUTargetMC() }
func LLVMInitializeARMTargetMC()         { C.LLVMInitializeARMTargetMC() }
func LLVMInitializeAVRTargetMC()         { C.LLVMInitializeAVRTargetMC() }
func LLVMInitializeBPFTargetMC()         { C.LLVMInitializeBPFTargetMC() }
func LLVMInitializeHexagonTargetMC()     { C.LLVMInitializeHexagonTargetMC() }
func LLVMInitializeLanaiTargetMC()       { C.LLVMInitializeLanaiTargetMC() }
func LLVMInitializeLoongArchTargetMC()   { C.LLVMInitializeLoongArchTargetMC() }
func LLVMInitializeMipsTargetMC()        { C.LLVMInitializeMipsTargetMC() }
func LLVMInitializeMSP430TargetMC()      { C.LLVMInitializeMSP430TargetMC() }
func LLVMInitializeNVPTXTargetMC()       { C.LLVMInitializeNVPTXTargetMC() }
func LLVMInitializePowerPCTargetMC()     { C.LLVMInitializePowerPCTargetMC() }
func LLVMInitializeRISCVTargetMC()       { C.LLVMInitializeRISCVTargetMC() }
func LLVMInitializeSparcTargetMC()       { C.LLVMInitializeSparcTargetMC() }
func LLVMInitializeSystemZTargetMC()     { C.LLVMInitializeSystemZTargetMC() }
func LLVMInitializeVETargetMC()          { C.LLVMInitializeVETargetMC() }
func LLVMInitializeWebAssemblyTargetMC() { C.LLVMInitializeWebAssemblyTargetMC() }
func LLVMInitializeX86TargetMC()         { C.LLVMInitializeX86TargetMC() }
func LLVMInitializeXCoreTargetMC()       { C.LLVMInitializeXCoreTargetMC() }

func LLVMInitializeAArch64AsmParser()     { C.LLVMInitializeAArch64AsmParser() }
func LLVMInitializeAMDGPUAsmParser()      { C.LLVMInitializeAMDGPUAsmParser() }
func LLVMInitializeARMAsmParser()         { C.LLVMInitializeARMAsmParser() }
func LLVMInitializeAVRAsmParser()         { C.LLVMInitializeAVRAsmParser() }
func LLVMInitializeBPFAsmParser()         { C.LLVMInitializeBPFAsmParser() }
func LLVMInitializeHexagonAsmParser()     { C.LLVMInitializeHexagonAsmParser() }
func LLVMInitializeLanaiAsmParser()       { C.LLVMInitializeLanaiAsmParser() }
func LLVMInitializeLoongArchAsmParser()   { C.LLVMInitializeLoongArchAsmParser() }
func LLVMInitializeMipsAsmParser()        { C.LLVMInitializeMipsAsmParser() }
func LLVMInitializeMSP430AsmParser()      { C.LLVMInitializeMSP430AsmParser() }
func LLVMInitializePowerPCAsmParser()     { C.LLVMInitializePowerPCAsmParser() }
func LLVMInitializeRISCVAsmParser()       { C.LLVMInitializeRISCVAsmParser() }
func LLVMInitializeSparcAsmParser()       { C.LLVMInitializeSparcAsmParser() }
func LLVMInitializeSystemZAsmParser()     { C.LLVMInitializeSystemZAsmParser() }
func LLVMInitializeVEAsmParser()          { C.LLVMInitializeVEAsmParser() }
func LLVMInitializeWebAssemblyAsmParser() { C.LLVMInitializeWebAssemblyAsmParser() }
func LLVMInitializeX86AsmParser()         { C.LLVMInitializeX86AsmParser() }

func LLVMInitializeAArch64AsmPrinter()     { C.LLVMInitializeAArch64AsmPrinter() }
func LLVMInitializeAMDGPUAsmPrinter()      { C.LLVMInitializeAMDGPUAsmPrinter() }
func LLVMInitializeARMAsmPrinter()         { C.LLVMInitializeARMAsmPrinter() }
func LLVMInitializeAVRAsmPrinter()         { C.LLVMInitializeAVRAsmPrinter() }
func LLVMInitializeBPFAsmPrinter()         { C.LLVMInitializeBPFAsmPrinter() }
func LLVMInitializeHexagonAsmPrinter()     { C.LLVMInitializeHexagonAsmPrinter() }
func LLVMInitializeLanaiAsmPrinter()       { C.LLVMInitializeLanaiAsmPrinter() }
func LLVMInitializeLoongArchAsmPrinter()   { C.LLVMInitializeLoongArchAsmPrinter() }
func LLVMInitializeMipsAsmPrinter()        { C.LLVMInitializeMipsAsmPrinter() }
func LLVMInitializeMSP430AsmPrinter()      { C.LLVMInitializeMSP430AsmPrinter() }
func LLVMInitializePowerPCAsmPrinter()     { C.LLVMInitializePowerPCAsmPrinter() }
func LLVMInitializeRISCVAsmPrinter()       { C.LLVMInitializeRISCVAsmPrinter() }
func LLVMInitializeSparcAsmPrinter()       { C.LLVMInitializeSparcAsmPrinter() }
func LLVMInitializeSystemZAsmPrinter()     { C.LLVMInitializeSystemZAsmPrinter() }
func LLVMInitializeVEAsmPrinter()          { C.LLVMInitializeVEAsmPrinter() }
func LLVMInitializeWebAssemblyAsmPrinter() { C.LLVMInitializeWebAssemblyAsmPrinter() }
func LLVMInitializeX86AsmPrinter()         { C.LLVMInitializeX86AsmPrinter() }

func LLVMInitializeAArch64Disassembler()     { C.LLVMInitializeAArch64Disassembler() }
func LLVMInitializeAMDGPUDisassembler()      { C.LLVMInitializeAMDGPUDisassembler() }
func LLVMInitializeARMDisassembler()         { C.LLVMInitializeARMDisassembler() }
func LLVMInitializeAVRDisassembler()         { C.LLVMInitializeAVRDisassembler() }
func LLVMInitializeBPFDisassembler()         { C.LLVMInitializeBPFDisassembler() }
func LLVMInitializeHexagonDisassembler()     { C.LLVMInitializeHexagonDisassembler() }
func LLVMInitializeLanaiDisassembler()       { C.LLVMInitializeLanaiDisassembler() }
func LLVMInitializeLoongArchDisassembler()   { C.LLVMInitializeLoongArchDisassembler() }
func LLVMInitializeMipsDisassembler()        { C.LLVMInitializeMipsDisassembler() }
func LLVMInitializeMSP430Disassembler()      { C.LLVMInitializeMSP430Disassembler() }
func LLVMInitializePowerPCDisassembler()     { C.LLVMInitializePowerPCDisassembler() }
func LLVMInitializeRISCVDisassembler()       { C.LLVMInitializeRISCVDisassembler() }
func LLVMInitializeSparcDisassembler()       { C.LLVMInitializeSparcDisassembler() }
func LLVMInitializeSystemZDisassembler()     { C.LLVMInitializeSystemZDisassembler() }
func LLVMInitializeVEDisassembler()          { C.LLVMInitializeVEDisassembler() }
func LLVMInitializeWebAssemblyDisassembler() { C.LLVMInitializeWebAssemblyDisassembler() }
func LLVMInitializeX86Disassembler()         { C.LLVMInitializeX86Disassembler() }
func LLVMInitializeXCoreDisassembler()       { C.LLVMInitializeXCoreDisassembler() }

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

// LLVMByteOrder Returns the byte order of a target, either LLVMBigEndian or LLVMLittleEndian.
func LLVMByteOrder(td LLVMTargetDataRef) LLVMByteOrdering {
	return LLVMByteOrdering(C.LLVMByteOrder(td.c))
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
