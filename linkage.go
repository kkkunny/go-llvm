package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// Linkage 链接类型
type Linkage binding.LLVMLinkage

const (
	LinkageExternal                = Linkage(binding.LLVMExternalLinkage)
	LinkageAvailableExternally     = Linkage(binding.LLVMAvailableExternallyLinkage)
	LinkageLinkOnceAny             = Linkage(binding.LLVMLinkOnceAnyLinkage)
	LinkageLinkOnceODR             = Linkage(binding.LLVMLinkOnceODRLinkage)
	LinkageLinkOnceODRAutoHide     = Linkage(binding.LLVMLinkOnceODRAutoHideLinkage)
	LinkageWeakAny                 = Linkage(binding.LLVMWeakAnyLinkage)
	LinkageWeakODR                 = Linkage(binding.LLVMWeakODRLinkage)
	LinkageAppending               = Linkage(binding.LLVMAppendingLinkage)
	LinkageInternal                = Linkage(binding.LLVMInternalLinkage)
	LinkagePrivate                 = Linkage(binding.LLVMPrivateLinkage)
	LinkageExternalWeak            = Linkage(binding.LLVMExternalWeakLinkage)
	LinkageCommon                  = Linkage(binding.LLVMCommonLinkage)
)

// Visibility 符号可见性
type Visibility binding.LLVMVisibility

const (
	VisibilityDefault = Visibility(binding.LLVMDefaultVisibility)
	VisibilityHidden  = Visibility(binding.LLVMHiddenVisibility)
	VisibilityProtected = Visibility(binding.LLVMProtectedVisibility)
)
