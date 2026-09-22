package binding

// LLVMOptimizationLevel 优化级别管线名称（与 PassBuilder 的 OptimizationLevel 对应）
type LLVMOptimizationLevel string

const (
	LLVMOptimizationLevelO0 LLVMOptimizationLevel = "O0"
	LLVMOptimizationLevelO1 LLVMOptimizationLevel = "O1"
	LLVMOptimizationLevelO2 LLVMOptimizationLevel = "O2"
	LLVMOptimizationLevelO3 LLVMOptimizationLevel = "O3"
	LLVMOptimizationLevelOz LLVMOptimizationLevel = "Oz"
	LLVMOptimizationLevelOs LLVMOptimizationLevel = "Os"
)
