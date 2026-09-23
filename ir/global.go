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
	g.Check("ir.Global.ValueType")
	ref := binding.LLVMGlobalGetValueType(g.Ref())
	return llvm.TypeOfRef(g.Context(), ref)
}

// Initializer 初始化器
func (g Global) Initializer() (llvm.Value[llvm.DynT], bool) {
	g.Check("ir.Global.Initializer")
	ref := binding.LLVMGetInitializer(g.Ref())
	if ref.IsNil() {
		return llvm.Value[llvm.DynT]{}, false
	}
	return wrapDyn(g.Context(), g.Lifetime(), ref), true
}

// SetInitializer 设置初始化器
func (g Global) SetInitializer(v llvm.AnyValue) {
	const op = "ir.Global.SetInitializer"
	g.Check(op)
	g.Context().CheckValues(op, v)
	binding.LLVMSetInitializer(g.Ref(), v.Ref())
}

// IsConstant 全局是否为常量（有意遮蔽内嵌 Value.IsConstant：语义为全局常量标志）
func (g Global) IsConstant() bool {
	g.Check("ir.Global.IsConstant")
	return binding.LLVMIsGlobalConstant(g.Ref())
}

// SetConstant 设置是否常量
func (g Global) SetConstant(c bool) {
	g.Check("ir.Global.SetConstant")
	binding.LLVMSetGlobalConstant(g.Ref(), c)
}

// Align 对齐字节数
func (g Global) Align() uint32 {
	g.Check("ir.Global.Align")
	return binding.LLVMGetAlignment(g.Ref())
}

// SetAlign 设置对齐字节数
func (g Global) SetAlign(n uint32) {
	const op = "ir.Global.SetAlign"
	g.Check(op)
	preAlign(op, n)
	binding.LLVMSetAlignment(g.Ref(), n)
}

// Linkage 链接类型
func (g Global) Linkage() llvm.Linkage {
	g.Check("ir.Global.Linkage")
	return llvm.Linkage(binding.LLVMGetLinkage(g.Ref()))
}

// SetLinkage 设置链接类型
func (g Global) SetLinkage(l llvm.Linkage) {
	g.Check("ir.Global.SetLinkage")
	binding.LLVMSetLinkage(g.Ref(), binding.LLVMLinkage(l))
}
