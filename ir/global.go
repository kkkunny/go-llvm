package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Global 全局变量角色（内嵌 Value[PtrT]，自动实现 llvm.ValueRef/AnyValue）
type Global struct {
	llvm.Value[llvm.PtrT]
}

// ValueType 全局变量内容类型
func (g Global) ValueType() llvm.AnyType {
	ref := binding.LLVMGlobalGetValueType(g.Ref())
	return llvm.TypeOfRef(g.Context(), ref)
}

// Initializer 初始化器
func (g Global) Initializer() (llvm.Value[llvm.DynT], bool) {
	ref := binding.LLVMGetInitializer(g.Ref())
	if ref.IsNil() {
		return llvm.Value[llvm.DynT]{}, false
	}
	return wrapDyn(g.Context(), g.Lifetime(), ref), true
}

// SetInitializer 设置初始化器
func (g Global) SetInitializer(v llvm.AnyValue) {
	binding.LLVMSetInitializer(g.Ref(), v.Ref())
}

// IsConstant 全局是否为常量（有意遮蔽内嵌 Value.IsConstant：语义为全局常量标志）
func (g Global) IsConstant() bool { return binding.LLVMIsGlobalConstant(g.Ref()) }

// SetConstant 设置是否常量
func (g Global) SetConstant(c bool) { binding.LLVMSetGlobalConstant(g.Ref(), c) }

// Align 对齐字节数
func (g Global) Align() uint32 { return binding.LLVMGetAlignment(g.Ref()) }

// SetAlign 设置对齐字节数
func (g Global) SetAlign(n uint32) { binding.LLVMSetAlignment(g.Ref(), n) }

// Linkage 链接类型
func (g Global) Linkage() llvm.Linkage { return llvm.Linkage(binding.LLVMGetLinkage(g.Ref())) }

// SetLinkage 设置链接类型
func (g Global) SetLinkage(l llvm.Linkage) { binding.LLVMSetLinkage(g.Ref(), binding.LLVMLinkage(l)) }
