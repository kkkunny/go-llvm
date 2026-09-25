package ir

import (
	"iter"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Alias 全局别名角色
type Alias struct {
	llvm.Value[llvm.PtrT]
}

// Aliasee 别名指向的值
func (a Alias) Aliasee() llvm.Value[llvm.DynT] {
	return llvm.ValueOf(a.Context(), a.Lifetime(), binding.LLVMAliasGetAliasee(a.Ref()))
}

// SetAliasee 设置别名指向的值
func (a Alias) SetAliasee(v llvm.AnyValue) {
	const op = "ir.Alias.SetAliasee"
	a.Context().CheckValues(op, v)
	binding.LLVMAliasSetAliasee(a.Ref(), v.Ref())
}

// NewAlias 创建全局别名（valueTy 为被指值类型，地址空间 0）
func (m *Module) NewAlias(name string, valueTy llvm.AnyType, aliasee llvm.AnyValue) Alias {
	const op = "ir.Module.NewAlias"
	m.Check(op)
	m.ctx.CheckType(op, valueTy)
	m.ctx.CheckValues(op, aliasee)
	ref := binding.LLVMAddAlias2(m.ref, valueTy.Ref(), 0, aliasee.Ref(), name)
	return Alias{Value: llvm.NewValue[llvm.PtrT](m.ctx, m.life, ref)}
}

// NewIFunc 创建全局 IFunc（resolver 为解析函数，地址空间 0）
func (m *Module) NewIFunc(name string, t llvm.AnyType, resolver llvm.AnyValue) Global {
	const op = "ir.Module.NewIFunc"
	m.Check(op)
	m.ctx.CheckType(op, t)
	m.ctx.CheckValues(op, resolver)
	ref := binding.LLVMAddGlobalIFunc(m.ref, name, t.Ref(), 0, resolver.Ref())
	return Global{Value: llvm.NewValue[llvm.PtrT](m.ctx, m.life, ref)}
}

// GetAlias 按名查找别名
func (m *Module) GetAlias(name string) (Alias, bool) {
	m.Check("ir.Module.GetAlias")
	ref := binding.LLVMGetNamedGlobalAlias(m.ref, name)
	if ref.IsNil() {
		return Alias{}, false
	}
	return Alias{Value: llvm.NewValue[llvm.PtrT](m.ctx, m.life, ref)}, true
}

// AllAliases 模块内全部别名的惰性遍历
func (m *Module) AllAliases() iter.Seq[Alias] {
	return func(yield func(Alias) bool) {
		m.Check("ir.Module.AllAliases")
		for ref := binding.LLVMGetFirstGlobalAlias(m.ref); !ref.IsNil(); ref = binding.LLVMGetNextGlobalAlias(ref) {
			if !yield(Alias{Value: llvm.NewValue[llvm.PtrT](m.ctx, m.life, ref)}) {
				return
			}
		}
	}
}
