package llvm

import (
	"unsafe"

	"github.com/kkkunny/go-llvm/internal/binding"
)

const TargetTriple = binding.LLVM_DEFAULT_TARGET_TRIPLE

// MajorVersion LLVM大版本号
const MajorVersion = binding.LLVM_VERSION_MAJOR

// Version LLVM版本号
const Version = binding.LLVM_VERSION_STRING

var (
	// CPUName cpu名称
	CPUName = binding.LLVMGetHostCPUName()
	// CPUFeatures cpu特性
	CPUFeatures = binding.LLVMGetHostCPUFeatures()
	// PointerSize 宿主 Go 进程的指针大小（字节）；目标指针大小请用 DataLayout（P1）
	PointerSize = unsafe.Sizeof(uintptr(0))
)
