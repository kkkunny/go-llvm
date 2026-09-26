package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
)

// Intrinsic LLVM 内建函数（运行时按名查找，inkwell Intrinsic::find 的对应设计）。
// 不做代码生成：声明取得后用 Builder.Call 调用。
type Intrinsic struct {
	id uint32
}

// FindIntrinsic 按名查找内建函数；不存在返回 false
func FindIntrinsic(name string) (Intrinsic, bool) {
	id := binding.LLVMLookupIntrinsicID(name)
	if id == 0 {
		return Intrinsic{}, false
	}
	return Intrinsic{id: id}, true
}

// Name 内建函数名
func (i Intrinsic) Name() string { return binding.LLVMIntrinsicGetName(i.id) }

// IsOverloaded 是否重载（参数类型决定具体声明）
func (i Intrinsic) IsOverloaded() bool { return binding.LLVMIntrinsicIsOverloaded(i.id) }

// Declaration 取得/创建模块内的声明；重载 intrinsic 未提供参数类型返回 false
// （LLVM 对空参数的重载查询会崩溃，此处前置拦截——崩溃类地板，任何构建都生效）
func (i Intrinsic) Declaration(m *Module, params []llvm.AnyType) (Function, bool) {
	const op = "ir.Intrinsic.Declaration"
	m.Check(op)
	if i.IsOverloaded() && len(params) == 0 {
		return Function{}, false
	}
	if checks.Debug {
		m.ctx.CheckTypes(op, params)
	}
	ref := binding.LLVMGetIntrinsicDeclaration(m.ref, i.id, llvm.AnyTypesToRefs(params))
	if ref.IsNil() {
		return Function{}, false
	}
	return Function{Value: llvm.NewValue[llvm.FnT](m.ctx, m.life, ref)}, true
}
