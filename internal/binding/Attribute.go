package binding

/*
#include "llvm-c/Core.h"
*/
import "C"

// LLVMCallConv 调用约定
type LLVMCallConv int32

const (
	LLVMCCallConv             LLVMCallConv = C.LLVMCCallConv
	LLVMFastCallConv          LLVMCallConv = C.LLVMFastCallConv
	LLVMColdCallConv          LLVMCallConv = C.LLVMColdCallConv
	LLVMGHCCallConv           LLVMCallConv = C.LLVMGHCCallConv
	LLVMHiPECallConv          LLVMCallConv = C.LLVMHiPECallConv
	LLVMAnyRegCallConv        LLVMCallConv = C.LLVMAnyRegCallConv
	LLVMPreserveMostCallConv  LLVMCallConv = C.LLVMPreserveMostCallConv
	LLVMPreserveAllCallConv   LLVMCallConv = C.LLVMPreserveAllCallConv
	LLVMSwiftCallConv         LLVMCallConv = C.LLVMSwiftCallConv
	LLVMCXXFASTTLSCallConv    LLVMCallConv = C.LLVMCXXFASTTLSCallConv
	LLVMX86StdcallCallConv    LLVMCallConv = C.LLVMX86StdcallCallConv
	LLVMX86FastcallCallConv   LLVMCallConv = C.LLVMX86FastcallCallConv
	LLVMARMAPCSCallConv       LLVMCallConv = C.LLVMARMAPCSCallConv
	LLVMARMAAPCSCallConv      LLVMCallConv = C.LLVMARMAAPCSCallConv
	LLVMARMAAPCSVFPCallConv   LLVMCallConv = C.LLVMARMAAPCSVFPCallConv
	LLVMMSP430INTRCallConv    LLVMCallConv = C.LLVMMSP430INTRCallConv
	LLVMX86ThisCallCallConv   LLVMCallConv = C.LLVMX86ThisCallCallConv
	LLVMPTXKernelCallConv     LLVMCallConv = C.LLVMPTXKernelCallConv
	LLVMPTXDeviceCallConv     LLVMCallConv = C.LLVMPTXDeviceCallConv
	LLVMSPIRFUNCCallConv      LLVMCallConv = C.LLVMSPIRFUNCCallConv
	LLVMSPIRKERNELCallConv    LLVMCallConv = C.LLVMSPIRKERNELCallConv
	LLVMIntelOCLBICallConv    LLVMCallConv = C.LLVMIntelOCLBICallConv
	LLVMX8664SysVCallConv     LLVMCallConv = C.LLVMX8664SysVCallConv
	LLVMWin64CallConv         LLVMCallConv = C.LLVMWin64CallConv
	LLVMX86VectorCallCallConv LLVMCallConv = C.LLVMX86VectorCallCallConv
	LLVMHHVMCallConv          LLVMCallConv = C.LLVMHHVMCallConv
	LLVMHHVMCCallConv         LLVMCallConv = C.LLVMHHVMCCallConv
	LLVMX86INTRCallConv       LLVMCallConv = C.LLVMX86INTRCallConv
	LLVMAVRINTRCallConv       LLVMCallConv = C.LLVMAVRINTRCallConv
	LLVMAVRSIGNALCallConv     LLVMCallConv = C.LLVMAVRSIGNALCallConv
	LLVMAVRBUILTINCallConv    LLVMCallConv = C.LLVMAVRBUILTINCallConv
	LLVMAMDGPUVSCallConv      LLVMCallConv = C.LLVMAMDGPUVSCallConv
	LLVMAMDGPUGSCallConv      LLVMCallConv = C.LLVMAMDGPUGSCallConv
	LLVMAMDGPUPSCallConv      LLVMCallConv = C.LLVMAMDGPUPSCallConv
	LLVMAMDGPUCSCallConv      LLVMCallConv = C.LLVMAMDGPUCSCallConv
	LLVMAMDGPUKERNELCallConv  LLVMCallConv = C.LLVMAMDGPUKERNELCallConv
	LLVMX86RegCallCallConv    LLVMCallConv = C.LLVMX86RegCallCallConv
	LLVMAMDGPUHSCallConv      LLVMCallConv = C.LLVMAMDGPUHSCallConv
	LLVMMSP430BUILTINCallConv LLVMCallConv = C.LLVMMSP430BUILTINCallConv
	LLVMAMDGPULSCallConv      LLVMCallConv = C.LLVMAMDGPULSCallConv
	LLVMAMDGPUESCallConv      LLVMCallConv = C.LLVMAMDGPUESCallConv
)

// LLVMGetLastEnumAttributeKind Return the last attribute kind ID.
func LLVMGetLastEnumAttributeKind() uint32 {
	return uint32(C.LLVMGetLastEnumAttributeKind())
}

// LLVMGetEnumAttributeKind Get the unique id corresponding to the enum attribute passed as argument.
func LLVMGetEnumAttributeKind(a LLVMAttributeRef) uint32 {
	return uint32(C.LLVMGetEnumAttributeKind(a.c))
}

// LLVMGetEnumAttributeValue Get the enum attribute's value. 0 is returned if none exists.
func LLVMGetEnumAttributeValue(a LLVMAttributeRef) uint64 {
	return uint64(C.LLVMGetEnumAttributeValue(a.c))
}

// LLVMCreateTypeAttribute Create a type attribute.
func LLVMCreateTypeAttribute(c LLVMContextRef, kindID uint32, ty LLVMTypeRef) LLVMAttributeRef {
	return LLVMAttributeRef{c: C.LLVMCreateTypeAttribute(c.c, C.unsigned(kindID), ty.c)}
}

// LLVMGetTypeAttributeValue Get the type attribute's value.
func LLVMGetTypeAttributeValue(a LLVMAttributeRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMGetTypeAttributeValue(a.c)}
}

// LLVMCreateStringAttribute Create a string attribute.
func LLVMCreateStringAttribute(c LLVMContextRef, k, v string) LLVMAttributeRef {
	klen, vlen := C.unsigned(len(k)), C.unsigned(len(v))
	return string2CString(k, func(ck *C.char) LLVMAttributeRef {
		return string2CString(v, func(cv *C.char) LLVMAttributeRef {
			return LLVMAttributeRef{c: C.LLVMCreateStringAttribute(c.c, ck, klen, cv, vlen)}
		})
	})
}

// LLVMGetStringAttributeKind Get the string attribute's kind.
func LLVMGetStringAttributeKind(a LLVMAttributeRef) string {
	var length C.unsigned
	ptr := C.LLVMGetStringAttributeKind(a.c, &length)
	if ptr == nil {
		return ""
	}
	return C.GoStringN(ptr, C.int(length))
}

// LLVMGetStringAttributeValue Get the string attribute's value.
func LLVMGetStringAttributeValue(a LLVMAttributeRef) string {
	var length C.unsigned
	ptr := C.LLVMGetStringAttributeValue(a.c, &length)
	if ptr == nil {
		return ""
	}
	return C.GoStringN(ptr, C.int(length))
}

// LLVMIsEnumAttribute Check whether the attribute is an enum attribute.
func LLVMIsEnumAttribute(a LLVMAttributeRef) bool {
	return llvmBool2bool(C.LLVMIsEnumAttribute(a.c))
}

// LLVMIsStringAttribute Check whether the attribute is a string attribute.
func LLVMIsStringAttribute(a LLVMAttributeRef) bool {
	return llvmBool2bool(C.LLVMIsStringAttribute(a.c))
}

// LLVMIsTypeAttribute Check whether the attribute is a type attribute.
func LLVMIsTypeAttribute(a LLVMAttributeRef) bool {
	return llvmBool2bool(C.LLVMIsTypeAttribute(a.c))
}

// LLVMGetAttributeCountAtIndex Obtain the number of attributes at the given index.
func LLVMGetAttributeCountAtIndex(f LLVMValueRef, idx LLVMAttributeIndex) uint32 {
	return uint32(C.LLVMGetAttributeCountAtIndex(f.c, C.LLVMAttributeIndex(idx)))
}

// LLVMGetAttributesAtIndex Obtain the attributes at the given index.
func LLVMGetAttributesAtIndex(f LLVMValueRef, idx LLVMAttributeIndex) []LLVMAttributeRef {
	n := int(LLVMGetAttributeCountAtIndex(f, idx))
	if n == 0 {
		return nil
	}
	attrs := make([]C.LLVMAttributeRef, n)
	C.LLVMGetAttributesAtIndex(f.c, C.LLVMAttributeIndex(idx), &attrs[0])
	refs := make([]LLVMAttributeRef, n)
	for i, a := range attrs {
		refs[i] = LLVMAttributeRef{c: a}
	}
	return refs
}

// LLVMGetEnumAttributeAtIndex Obtain the enum attribute at the given index.
func LLVMGetEnumAttributeAtIndex(f LLVMValueRef, idx LLVMAttributeIndex, kindID uint32) LLVMAttributeRef {
	return LLVMAttributeRef{c: C.LLVMGetEnumAttributeAtIndex(f.c, C.LLVMAttributeIndex(idx), C.unsigned(kindID))}
}

// LLVMGetStringAttributeAtIndex Obtain the string attribute at the given index.
func LLVMGetStringAttributeAtIndex(f LLVMValueRef, idx LLVMAttributeIndex, k string) LLVMAttributeRef {
	klen := C.unsigned(len(k))
	return string2CString(k, func(ck *C.char) LLVMAttributeRef {
		return LLVMAttributeRef{c: C.LLVMGetStringAttributeAtIndex(f.c, C.LLVMAttributeIndex(idx), ck, klen)}
	})
}

// LLVMRemoveEnumAttributeAtIndex Remove the enum attribute at the given index.
func LLVMRemoveEnumAttributeAtIndex(f LLVMValueRef, idx LLVMAttributeIndex, kindID uint32) {
	C.LLVMRemoveEnumAttributeAtIndex(f.c, C.LLVMAttributeIndex(idx), C.unsigned(kindID))
}

// LLVMRemoveStringAttributeAtIndex Remove the string attribute at the given index.
func LLVMRemoveStringAttributeAtIndex(f LLVMValueRef, idx LLVMAttributeIndex, k string) {
	klen := C.unsigned(len(k))
	string2CString(k, func(ck *C.char) bool {
		C.LLVMRemoveStringAttributeAtIndex(f.c, C.LLVMAttributeIndex(idx), ck, klen)
		return false
	})
}

// LLVMAddCallSiteAttribute Add an attribute to a call/invoke instruction.
func LLVMAddCallSiteAttribute(c LLVMValueRef, idx LLVMAttributeIndex, a LLVMAttributeRef) {
	C.LLVMAddCallSiteAttribute(c.c, C.LLVMAttributeIndex(idx), a.c)
}

// LLVMGetCallSiteAttributeCount Obtain the number of attributes at the given index of a call site.
func LLVMGetCallSiteAttributeCount(c LLVMValueRef, idx LLVMAttributeIndex) uint32 {
	return uint32(C.LLVMGetCallSiteAttributeCount(c.c, C.LLVMAttributeIndex(idx)))
}

// LLVMGetCallSiteAttributes Obtain the attributes at the given index of a call site.
func LLVMGetCallSiteAttributes(c LLVMValueRef, idx LLVMAttributeIndex) []LLVMAttributeRef {
	n := int(LLVMGetCallSiteAttributeCount(c, idx))
	if n == 0 {
		return nil
	}
	attrs := make([]C.LLVMAttributeRef, n)
	C.LLVMGetCallSiteAttributes(c.c, C.LLVMAttributeIndex(idx), &attrs[0])
	refs := make([]LLVMAttributeRef, n)
	for i, a := range attrs {
		refs[i] = LLVMAttributeRef{c: a}
	}
	return refs
}

// LLVMGetCallSiteEnumAttribute Obtain the enum attribute at the given index of a call site.
func LLVMGetCallSiteEnumAttribute(c LLVMValueRef, idx LLVMAttributeIndex, kindID uint32) LLVMAttributeRef {
	return LLVMAttributeRef{c: C.LLVMGetCallSiteEnumAttribute(c.c, C.LLVMAttributeIndex(idx), C.unsigned(kindID))}
}

// LLVMGetCallSiteStringAttribute Obtain the string attribute at the given index of a call site.
func LLVMGetCallSiteStringAttribute(c LLVMValueRef, idx LLVMAttributeIndex, k string) LLVMAttributeRef {
	klen := C.unsigned(len(k))
	return string2CString(k, func(ck *C.char) LLVMAttributeRef {
		return LLVMAttributeRef{c: C.LLVMGetCallSiteStringAttribute(c.c, C.LLVMAttributeIndex(idx), ck, klen)}
	})
}

// LLVMRemoveCallSiteEnumAttribute Remove the enum attribute at the given index of a call site.
func LLVMRemoveCallSiteEnumAttribute(c LLVMValueRef, idx LLVMAttributeIndex, kindID uint32) {
	C.LLVMRemoveCallSiteEnumAttribute(c.c, C.LLVMAttributeIndex(idx), C.unsigned(kindID))
}

// LLVMRemoveCallSiteStringAttribute Remove the string attribute at the given index of a call site.
func LLVMRemoveCallSiteStringAttribute(c LLVMValueRef, idx LLVMAttributeIndex, k string) {
	klen := C.unsigned(len(k))
	string2CString(k, func(ck *C.char) bool {
		C.LLVMRemoveCallSiteStringAttribute(c.c, C.LLVMAttributeIndex(idx), ck, klen)
		return false
	})
}

// LLVMGetFunctionCallConv Obtain the calling convention of a function.
func LLVMGetFunctionCallConv(fn LLVMValueRef) LLVMCallConv {
	return LLVMCallConv(C.LLVMGetFunctionCallConv(fn.c))
}

// LLVMSetFunctionCallConv Set the calling convention of a function.
func LLVMSetFunctionCallConv(fn LLVMValueRef, cc LLVMCallConv) {
	C.LLVMSetFunctionCallConv(fn.c, C.unsigned(cc))
}

// LLVMGetInstructionCallConv Obtain the calling convention of a call/invoke instruction.
func LLVMGetInstructionCallConv(instr LLVMValueRef) LLVMCallConv {
	return LLVMCallConv(C.LLVMGetInstructionCallConv(instr.c))
}

// LLVMSetInstructionCallConv Set the calling convention of a call/invoke instruction.
func LLVMSetInstructionCallConv(instr LLVMValueRef, cc LLVMCallConv) {
	C.LLVMSetInstructionCallConv(instr.c, C.unsigned(cc))
}
