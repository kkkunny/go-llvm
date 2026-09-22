package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Global 全局变量角色
type Global struct {
	v   llvm.Value[llvm.PtrT]
	mod *Module
}

// Value 返回底层泛型值（全局变量自身类型为 ptr）
func (g Global) Value() llvm.Value[llvm.PtrT] { return g.v }

// AsValue 实现 llvm.ValueRef[PtrT]
func (g Global) AsValue() llvm.Value[llvm.PtrT] { return g.v }

// Name 全局变量名
func (g Global) Name() string { return g.v.Name() }

// SetName 设置全局变量名
func (g Global) SetName(name string) { g.v.SetName(name) }

// ValueType 全局变量内容类型
func (g Global) ValueType() llvm.AnyType {
	ref := binding.LLVMGlobalGetValueType(g.v.Ref())
	return llvm.TypeOfRef(g.v.Context(), ref)
}

// Initializer 初始化器
func (g Global) Initializer() (llvm.Value[llvm.DynT], bool) {
	ref := binding.LLVMGetInitializer(g.v.Ref())
	if ref.IsNil() {
		return llvm.Value[llvm.DynT]{}, false
	}
	return llvm.ValueOf(g.v.Context(), g.mod.life, ref), true
}

// SetInitializer 设置初始化器
func (g Global) SetInitializer(v llvm.AnyValue) {
	binding.LLVMSetInitializer(g.v.Ref(), v.Ref())
}

// IsConstant 是否常量
func (g Global) IsConstant() bool { return binding.LLVMIsGlobalConstant(g.v.Ref()) }

// SetConstant 设置是否常量
func (g Global) SetConstant(c bool) { binding.LLVMSetGlobalConstant(g.v.Ref(), c) }

// Align 对齐字节数
func (g Global) Align() uint32 { return binding.LLVMGetAlignment(g.v.Ref()) }

// SetAlign 设置对齐字节数
func (g Global) SetAlign(n uint32) { binding.LLVMSetAlignment(g.v.Ref(), n) }

// Linkage 链接类型
func (g Global) Linkage() llvm.Linkage { return llvm.Linkage(binding.LLVMGetLinkage(g.v.Ref())) }

// SetLinkage 设置链接类型
func (g Global) SetLinkage(l llvm.Linkage) { binding.LLVMSetLinkage(g.v.Ref(), binding.LLVMLinkage(l)) }
