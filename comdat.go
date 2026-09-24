package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// ComdatSelectionKind Comdat 冲突解决方式
type ComdatSelectionKind binding.LLVMComdatSelectionKind

const (
	ComdatAny           = ComdatSelectionKind(binding.LLVMAnyComdatSelectionKind)
	ComdatExactMatch    = ComdatSelectionKind(binding.LLVMExactMatchComdatSelectionKind)
	ComdatLargest       = ComdatSelectionKind(binding.LLVMLargestComdatSelectionKind)
	ComdatNoDeduplicate = ComdatSelectionKind(binding.LLVMNoDeduplicateComdatSelectionKind)
	ComdatSameSize      = ComdatSelectionKind(binding.LLVMSameSizeComdatSelectionKind)
)
