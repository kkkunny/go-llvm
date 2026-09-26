package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// ComdatSelectionKind Comdat 冲突解决方式
type ComdatSelectionKind binding.LLVMComdatSelectionKind

// ComdatSelectionKind 取值对应 LLVM COMDAT 选择类型（binding.LLVMComdatSelectionKind）。
// 链接器据此在多个同名 COMDAT 组之间取舍：除 nodeduplicate 外最多保留一个组，其余组的成员一并丢弃。
const (
	ComdatAny           = ComdatSelectionKind(binding.LLVMAnyComdatSelectionKind)           // 任选：链接器任选一个同名组，不检查内容
	ComdatExactMatch    = ComdatSelectionKind(binding.LLVMExactMatchComdatSelectionKind)    // 内容一致：各组数据必须相同，链接器从中任选一个
	ComdatLargest       = ComdatSelectionKind(binding.LLVMLargestComdatSelectionKind)       // 保留最大：链接器选择体积最大的组
	ComdatNoDeduplicate = ComdatSelectionKind(binding.LLVMNoDeduplicateComdatSelectionKind) // 不去重：禁止链接器合并同名组
	ComdatSameSize      = ComdatSelectionKind(binding.LLVMSameSizeComdatSelectionKind)      // 大小一致：各组数据大小必须相同，链接器从中任选一个
)
