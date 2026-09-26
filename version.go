package llvm

import (
	"github.com/kkkunny/go-llvm/internal/binding"
)

// TargetTriple LLVM 默认目标三元组
const TargetTriple = binding.LLVM_DEFAULT_TARGET_TRIPLE

// MajorVersion LLVM大版本号
const MajorVersion = binding.LLVM_VERSION_MAJOR

// Version LLVM版本号
const Version = binding.LLVM_VERSION_STRING
