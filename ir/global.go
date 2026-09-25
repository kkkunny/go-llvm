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
	return llvm.ValueOf(g.Context(), g.Lifetime(), ref), true
}

// SetInitializer 设置初始化器
func (g Global) SetInitializer(v llvm.AnyValue) {
	const op = "ir.Global.SetInitializer"
	g.Context().CheckValues(op, v)
	binding.LLVMSetInitializer(g.Ref(), v.Ref())
}

// IsConstant 全局是否为常量（有意遮蔽内嵌 Value.IsConstant：语义为全局常量标志）
func (g Global) IsConstant() bool {
	return binding.LLVMIsGlobalConstant(g.Ref())
}

// SetConstant 设置是否常量
func (g Global) SetConstant(c bool) {
	binding.LLVMSetGlobalConstant(g.Ref(), c)
}

// Align 对齐字节数
func (g Global) Align() uint32 {
	return binding.LLVMGetAlignment(g.Ref())
}

// SetAlign 设置对齐字节数
func (g Global) SetAlign(n uint32) {
	const op = "ir.Global.SetAlign"
	preAlign(op, n)
	binding.LLVMSetAlignment(g.Ref(), n)
}

// Linkage 链接类型
func (g Global) Linkage() llvm.Linkage {
	return llvm.Linkage(binding.LLVMGetLinkage(g.Ref()))
}

// SetLinkage 设置链接类型
func (g Global) SetLinkage(l llvm.Linkage) {
	binding.LLVMSetLinkage(g.Ref(), binding.LLVMLinkage(l))
}

// Section 段名
func (g Global) Section() string {
	return binding.LLVMGetSection(g.Ref())
}

// SetSection 设置段名
func (g Global) SetSection(s string) {
	binding.LLVMSetSection(g.Ref(), s)
}

// Visibility 符号可见性
func (g Global) Visibility() llvm.Visibility {
	return llvm.Visibility(binding.LLVMGetVisibility(g.Ref()))
}

// SetVisibility 设置符号可见性
func (g Global) SetVisibility(v llvm.Visibility) {
	binding.LLVMSetVisibility(g.Ref(), binding.LLVMVisibility(v))
}

// DLLStorageClass DLL 存储类
func (g Global) DLLStorageClass() llvm.DLLStorageClass {
	return llvm.DLLStorageClass(binding.LLVMGetDLLStorageClass(g.Ref()))
}

// SetDLLStorageClass 设置 DLL 存储类
func (g Global) SetDLLStorageClass(c llvm.DLLStorageClass) {
	binding.LLVMSetDLLStorageClass(g.Ref(), binding.LLVMDLLStorageClass(c))
}

// UnnamedAddr 匿名地址语义
func (g Global) UnnamedAddr() llvm.UnnamedAddr {
	return llvm.UnnamedAddr(binding.LLVMGetUnnamedAddress(g.Ref()))
}

// SetUnnamedAddr 设置匿名地址语义
func (g Global) SetUnnamedAddr(a llvm.UnnamedAddr) {
	binding.LLVMSetUnnamedAddress(g.Ref(), binding.LLVMUnnamedAddr(a))
}

// ThreadLocalMode 线程局部模式
func (g Global) ThreadLocalMode() llvm.ThreadLocalMode {
	return llvm.ThreadLocalMode(binding.LLVMGetThreadLocalMode(g.Ref()))
}

// SetThreadLocalMode 设置线程局部模式
func (g Global) SetThreadLocalMode(m llvm.ThreadLocalMode) {
	binding.LLVMSetThreadLocalMode(g.Ref(), binding.LLVMThreadLocalMode(m))
}
