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
	case binding.LLVMFunctionValueKind:
		return Function(ref)
	case binding.LLVMInstructionValueKind:
		return lookupInstruction(ref).(Value)
	default:
		panic(fmt.Errorf("unknown enum value `%d`", binding.LLVMGetValueKind(ref)))
	}
}
