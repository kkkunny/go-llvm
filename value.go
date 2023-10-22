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
		panic(fmt.Errorf("unknown enum value `%d`", binding.LLVMGetValueKind(ref)))
	}
}

type Param binding.LLVMValueRef

func (v Param) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}
func (v Param) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(v)
}

func (v Param) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (p Param) Belong() Function {
	return Function(binding.LLVMGetParamParent(p.binding()))
}
