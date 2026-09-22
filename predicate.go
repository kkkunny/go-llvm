package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// IntPred 整数比较谓词
type IntPred binding.LLVMIntPredicate

const (
	IntEQ  = IntPred(binding.LLVMIntEQ)
	IntNE  = IntPred(binding.LLVMIntNE)
	IntUGT = IntPred(binding.LLVMIntUGT)
	IntUGE = IntPred(binding.LLVMIntUGE)
	IntULT = IntPred(binding.LLVMIntULT)
	IntULE = IntPred(binding.LLVMIntULE)
	IntSGT = IntPred(binding.LLVMIntSGT)
	IntSGE = IntPred(binding.LLVMIntSGE)
	IntSLT = IntPred(binding.LLVMIntSLT)
	IntSLE = IntPred(binding.LLVMIntSLE)
)

// FloatPred 浮点比较谓词
type FloatPred binding.LLVMRealPredicate

const (
	FloatFalse = FloatPred(binding.LLVMRealPredicateFalse)
	FloatOEQ   = FloatPred(binding.LLVMRealOEQ)
	FloatOGT   = FloatPred(binding.LLVMRealOGT)
	FloatOGE   = FloatPred(binding.LLVMRealOGE)
	FloatOLT   = FloatPred(binding.LLVMRealOLT)
	FloatOLE   = FloatPred(binding.LLVMRealOLE)
	FloatONE   = FloatPred(binding.LLVMRealONE)
	FloatORD   = FloatPred(binding.LLVMRealORD)
	FloatUNO   = FloatPred(binding.LLVMRealUNO)
	FloatUEQ   = FloatPred(binding.LLVMRealUEQ)
	FloatUGT   = FloatPred(binding.LLVMRealUGT)
	FloatUGE   = FloatPred(binding.LLVMRealUGE)
	FloatULT   = FloatPred(binding.LLVMRealULT)
	FloatULE   = FloatPred(binding.LLVMRealULE)
	FloatUNE   = FloatPred(binding.LLVMRealUNE)
	FloatTrue  = FloatPred(binding.LLVMRealPredicateTrue)
)
