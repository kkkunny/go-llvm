package llvm

import (
	"fmt"

	"github.com/kkkunny/go-llvm/internal/binding"
)

type Value interface {
	fmt.Stringer
	binding() binding.LLVMValueRef
	Type() Type
}

func lookupValue(ref binding.LLVMValueRef) Value {
	if ref.IsNil() {
		return nil
	}

	if binding.LLVMIsConstant(ref) {
		return lookupConstant(ref)
	}
	switch binding.LLVMGetValueKind(ref) {
	case binding.LLVMArgumentValueKind:
		return Param(ref)
	case binding.LLVMFunctionValueKind:
		return Function(ref)
	case binding.LLVMInstructionValueKind:
		return lookupInstruction(ref).(Value)
	default:
		return fallbackValue{ref}
	}
}

// fallbackValue 未知 kind 的通用降级值
type fallbackValue struct{ ref binding.LLVMValueRef }

func (v fallbackValue) String() string {
	return binding.LLVMPrintValueToString(v.ref)
}
func (v fallbackValue) binding() binding.LLVMValueRef {
	return v.ref
}
func (v fallbackValue) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.ref))
}

type Param binding.LLVMValueRef

func (p Param) String() string {
	return binding.LLVMPrintValueToString(p.binding())
}
func (p Param) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(p)
}

func (p Param) Type() Type {
	return lookupType(binding.LLVMTypeOf(p.binding()))
}

func (p Param) Belong() Function {
	return Function(binding.LLVMGetParamParent(p.binding()))
}

func (p Param) SetAlign(align uint32) {
	binding.LLVMSetParamAlignment(p.binding(), align)
}
