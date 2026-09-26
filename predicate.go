package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// IntPred 整数比较谓词
type IntPred binding.LLVMIntPredicate

// IntPred 取值对应 LLVM 整数比较谓词（binding.LLVMIntPredicate），用于 icmp 指令。
// EQ/NE 与符号无关，U 前缀为无符号比较，S 前缀为有符号比较。
const (
	IntEQ  = IntPred(binding.LLVMIntEQ)  // 等于（与符号无关）
	IntNE  = IntPred(binding.LLVMIntNE)  // 不等于（与符号无关）
	IntUGT = IntPred(binding.LLVMIntUGT) // 无符号大于
	IntUGE = IntPred(binding.LLVMIntUGE) // 无符号大于等于
	IntULT = IntPred(binding.LLVMIntULT) // 无符号小于
	IntULE = IntPred(binding.LLVMIntULE) // 无符号小于等于
	IntSGT = IntPred(binding.LLVMIntSGT) // 有符号大于
	IntSGE = IntPred(binding.LLVMIntSGE) // 有符号大于等于
	IntSLT = IntPred(binding.LLVMIntSLT) // 有符号小于
	IntSLE = IntPred(binding.LLVMIntSLE) // 有符号小于等于
)

// FloatPred 浮点比较谓词
type FloatPred binding.LLVMRealPredicate

// FloatPred 取值对应 LLVM 浮点比较谓词（binding.LLVMRealPredicate），用于 fcmp 指令。
// O 前缀表示有序（两个操作数均非 NaN），U 前缀表示无序（任一操作数为 NaN 即为真）。
const (
	FloatFalse = FloatPred(binding.LLVMRealPredicateFalse) // 恒假
	FloatOEQ   = FloatPred(binding.LLVMRealOEQ)            // 有序且相等
	FloatOGT   = FloatPred(binding.LLVMRealOGT)            // 有序大于
	FloatOGE   = FloatPred(binding.LLVMRealOGE)            // 有序大于等于
	FloatOLT   = FloatPred(binding.LLVMRealOLT)            // 有序小于
	FloatOLE   = FloatPred(binding.LLVMRealOLE)            // 有序小于等于
	FloatONE   = FloatPred(binding.LLVMRealONE)            // 有序且不等
	FloatORD   = FloatPred(binding.LLVMRealORD)            // 有序（两个操作数均非 NaN）
	FloatUNO   = FloatPred(binding.LLVMRealUNO)            // 无序（任一操作数为 NaN）
	FloatUEQ   = FloatPred(binding.LLVMRealUEQ)            // 无序或相等
	FloatUGT   = FloatPred(binding.LLVMRealUGT)            // 无序或大于
	FloatUGE   = FloatPred(binding.LLVMRealUGE)            // 无序或大于等于
	FloatULT   = FloatPred(binding.LLVMRealULT)            // 无序或小于
	FloatULE   = FloatPred(binding.LLVMRealULE)            // 无序或小于等于
	FloatUNE   = FloatPred(binding.LLVMRealUNE)            // 无序或不等
	FloatTrue  = FloatPred(binding.LLVMRealPredicateTrue)  // 恒真
)
