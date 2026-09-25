package binding

/*
#include "llvm-c/Core.h"
*/
import "C"

// LLVMModuleFlagBehavior is the merge behavior for module-level flags.
type LLVMModuleFlagBehavior int32

const (
	LLVMModuleFlagBehaviorError        LLVMModuleFlagBehavior = C.LLVMModuleFlagBehaviorError
	LLVMModuleFlagBehaviorWarning      LLVMModuleFlagBehavior = C.LLVMModuleFlagBehaviorWarning
	LLVMModuleFlagBehaviorRequire      LLVMModuleFlagBehavior = C.LLVMModuleFlagBehaviorRequire
	LLVMModuleFlagBehaviorOverride     LLVMModuleFlagBehavior = C.LLVMModuleFlagBehaviorOverride
	LLVMModuleFlagBehaviorAppend       LLVMModuleFlagBehavior = C.LLVMModuleFlagBehaviorAppend
	LLVMModuleFlagBehaviorAppendUnique LLVMModuleFlagBehavior = C.LLVMModuleFlagBehaviorAppendUnique
)

// LLVMMDStringInContext2 Create an MDString value from a given string value.
func LLVMMDStringInContext2(c LLVMContextRef, s string) LLVMMetadataRef {
	return string2CString(s, func(cs *C.char) LLVMMetadataRef {
		return LLVMMetadataRef{c: C.LLVMMDStringInContext2(c.c, cs, C.size_t(len(s)))}
	})
}

// LLVMMDNodeInContext2 Create an MDNode value with the given array of operands.
func LLVMMDNodeInContext2(c LLVMContextRef, mds []LLVMMetadataRef) LLVMMetadataRef {
	ptr, length := slice2Ptr[LLVMMetadataRef, C.LLVMMetadataRef](mds)
	return LLVMMetadataRef{c: C.LLVMMDNodeInContext2(c.c, ptr, C.size_t(length))}
}

// LLVMMetadataAsValue Obtain a Metadata as a Value.
func LLVMMetadataAsValue(c LLVMContextRef, md LLVMMetadataRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMMetadataAsValue(c.c, md.c)}
}

// LLVMValueAsMetadata Obtain a Value as a Metadata.
func LLVMValueAsMetadata(v LLVMValueRef) LLVMMetadataRef {
	return LLVMMetadataRef{c: C.LLVMValueAsMetadata(v.c)}
}

// LLVMGetMDString Obtain the underlying string from a MDString value; ok is false when the value is not an MDString.
func LLVMGetMDString(v LLVMValueRef) (string, bool) {
	var length C.unsigned
	ptr := C.LLVMGetMDString(v.c, &length)
	if ptr == nil {
		return "", false
	}
	return C.GoStringN(ptr, C.int(length)), true
}

// LLVMIsAMDNode returns non-nil for both MDNode and ValueAsMetadata (an upstream
// LLVM-C behavior); combine with LLVMIsAValueAsMetadata to distinguish real nodes.
func LLVMIsAMDNode(v LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMIsAMDNode(v.c)}
}

// LLVMIsAValueAsMetadata returns non-nil when the metadata-as-value is a ValueAsMetadata.
func LLVMIsAValueAsMetadata(v LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMIsAValueAsMetadata(v.c)}
}

// LLVMIsAMDString returns non-nil when the metadata-as-value is an MDString.
func LLVMIsAMDString(v LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMIsAMDString(v.c)}
}

// LLVMGetMDNodeNumOperands Obtain the number of operands from an MDNode value.
func LLVMGetMDNodeNumOperands(v LLVMValueRef) uint32 {
	return uint32(C.LLVMGetMDNodeNumOperands(v.c))
}

// LLVMGetMDNodeOperands Obtain the given MDNode's operands.
func LLVMGetMDNodeOperands(v LLVMValueRef) []LLVMValueRef {
	n := int(LLVMGetMDNodeNumOperands(v))
	if n == 0 {
		return nil
	}
	operands := make([]C.LLVMValueRef, n)
	C.LLVMGetMDNodeOperands(v.c, &operands[0])
	refs := make([]LLVMValueRef, n)
	for i, o := range operands {
		refs[i] = LLVMValueRef{c: o}
	}
	return refs
}

// LLVMGetNamedMetadataNumOperands Obtain the number of operands for named metadata in a module.
func LLVMGetNamedMetadataNumOperands(m LLVMModuleRef, name string) uint32 {
	return string2CString(name, func(name *C.char) uint32 {
		return uint32(C.LLVMGetNamedMetadataNumOperands(m.c, name))
	})
}

// LLVMGetNamedMetadataOperands Obtain the named metadata operands for a module.
func LLVMGetNamedMetadataOperands(m LLVMModuleRef, name string) []LLVMValueRef {
	n := int(LLVMGetNamedMetadataNumOperands(m, name))
	if n == 0 {
		return nil
	}
	operands := make([]C.LLVMValueRef, n)
	string2CString(name, func(name *C.char) bool {
		C.LLVMGetNamedMetadataOperands(m.c, name, &operands[0])
		return false
	})
	refs := make([]LLVMValueRef, n)
	for i, o := range operands {
		refs[i] = LLVMValueRef{c: o}
	}
	return refs
}

// LLVMAddNamedMetadataOperand Add an operand to named metadata.
func LLVMAddNamedMetadataOperand(m LLVMModuleRef, name string, v LLVMValueRef) {
	string2CString(name, func(name *C.char) bool {
		C.LLVMAddNamedMetadataOperand(m.c, name, v.c)
		return false
	})
}

// LLVMAddModuleFlag Add a module-level flag to the module-level flags metadata if it doesn't already exist.
func LLVMAddModuleFlag(m LLVMModuleRef, behavior LLVMModuleFlagBehavior, key string, val LLVMMetadataRef) {
	klen := C.size_t(len(key))
	string2CString(key, func(ckey *C.char) bool {
		C.LLVMAddModuleFlag(m.c, C.LLVMModuleFlagBehavior(behavior), ckey, klen, val.c)
		return false
	})
}

// LLVMGetModuleFlag Get a module-level flag from the module-level flags metadata; returns nil if not set.
func LLVMGetModuleFlag(m LLVMModuleRef, key string) LLVMMetadataRef {
	klen := C.size_t(len(key))
	return string2CString(key, func(ckey *C.char) LLVMMetadataRef {
		return LLVMMetadataRef{c: C.LLVMGetModuleFlag(m.c, ckey, klen)}
	})
}

// LLVMGetMDKindIDInContext Map a metadata kind name to a ID unique within this context.
func LLVMGetMDKindIDInContext(c LLVMContextRef, name string) uint32 {
	nlen := C.unsigned(len(name))
	return string2CString(name, func(cname *C.char) uint32 {
		return uint32(C.LLVMGetMDKindIDInContext(c.c, cname, nlen))
	})
}

// LLVMSetMetadata Set the metadata of the specified kind on a value.
func LLVMSetMetadata(v LLVMValueRef, kindID uint32, node LLVMValueRef) {
	C.LLVMSetMetadata(v.c, C.unsigned(kindID), node.c)
}

// LLVMGetMetadata Get the metadata of the specified kind on a value; returns nil if absent.
func LLVMGetMetadata(v LLVMValueRef, kindID uint32) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetMetadata(v.c, C.unsigned(kindID))}
}

// LLVMBlockAddress Get the block address of a basic block in a function.
func LLVMBlockAddress(f LLVMValueRef, bb LLVMBasicBlockRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBlockAddress(f.c, bb.c)}
}

// LLVMGetBlockAddressFunction Get the function associated with a given BlockAddress constant value.
func LLVMGetBlockAddressFunction(blockAddr LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetBlockAddressFunction(blockAddr.c)}
}

// LLVMGetBlockAddressBasicBlock Get the basic block associated with a given BlockAddress constant value.
func LLVMGetBlockAddressBasicBlock(blockAddr LLVMValueRef) LLVMBasicBlockRef {
	return LLVMBasicBlockRef{c: C.LLVMGetBlockAddressBasicBlock(blockAddr.c)}
}

// LLVMInlineAsmDialect is the inline assembly dialect.
type LLVMInlineAsmDialect int32

const (
	LLVMInlineAsmDialectATT   LLVMInlineAsmDialect = C.LLVMInlineAsmDialectATT
	LLVMInlineAsmDialectIntel LLVMInlineAsmDialect = C.LLVMInlineAsmDialectIntel
)

// LLVMGetInlineAsm Create the specified uniqued inline asm string.
func LLVMGetInlineAsm(ty LLVMTypeRef, asmString, constraints string, hasSideEffects, isAlignStack bool, dialect LLVMInlineAsmDialect, canThrow bool) LLVMValueRef {
	return string2CString(asmString, func(casm *C.char) LLVMValueRef {
		return string2CString(constraints, func(cconstraints *C.char) LLVMValueRef {
			return LLVMValueRef{c: C.LLVMGetInlineAsm(ty.c, casm, C.size_t(len(asmString)), cconstraints, C.size_t(len(constraints)),
				bool2LLVMBool(hasSideEffects), bool2LLVMBool(isAlignStack), C.LLVMInlineAsmDialect(dialect), bool2LLVMBool(canThrow))}
		})
	})
}

// LLVMGetInlineAsmFunctionType Get the function type of an inline asm value.
func LLVMGetInlineAsmFunctionType(inlineAsmVal LLVMValueRef) LLVMTypeRef {
	return LLVMTypeRef{c: C.LLVMGetInlineAsmFunctionType(inlineAsmVal.c)}
}
