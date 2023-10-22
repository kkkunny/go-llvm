package llvm

import "C"
import "github.com/kkkunny/go-llvm/internal/binding"

const Version = binding.Version

var (
	CPUName     = binding.LLVMGetHostCPUName()
	CPUFeatures = binding.LLVMGetHostCPUFeatures()
)
