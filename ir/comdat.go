package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Comdat comdat 角色（Module 拥有，随模块失效）
type Comdat struct {
	ref  binding.LLVMComdatRef
	ctx  *llvm.Context
	life *llvm.Lifetime
	name string
}

// Ref 返回底层句柄（供包内桥接使用）
func (c Comdat) Ref() binding.LLVMComdatRef { return c.ref }

// Context 返回所属上下文
func (c Comdat) Context() *llvm.Context { return c.ctx }

// Lifetime 返回所属生命周期令牌
func (c Comdat) Lifetime() *llvm.Lifetime { return c.life }

// Check comdat 可用性前置校验
func (c Comdat) Check(op string) {
	if c.ref.IsNil() || c.ctx == nil {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil comdat")
	}
	c.ctx.CheckAlive(op)
	if c.life == nil || !c.life.Alive() {
		llvm.Panicf(llvm.ErrUseAfterFree, op, "comdat is freed")
	}
}

// Name comdat 名称（经 Global.Comdat 取得时可能为空）
func (c Comdat) Name() string {
	c.Check("ir.Comdat.Name")
	return c.name
}

// SelectionKind 冲突解决方式
func (c Comdat) SelectionKind() llvm.ComdatSelectionKind {
	c.Check("ir.Comdat.SelectionKind")
	return llvm.ComdatSelectionKind(binding.LLVMGetComdatSelectionKind(c.ref))
}

// SetSelectionKind 设置冲突解决方式
func (c Comdat) SetSelectionKind(k llvm.ComdatSelectionKind) {
	c.Check("ir.Comdat.SetSelectionKind")
	binding.LLVMSetComdatSelectionKind(c.ref, binding.LLVMComdatSelectionKind(k))
}

// GetOrInsertComdat 取得模块中指定名称的 comdat，不存在则创建
func (m *Module) GetOrInsertComdat(name string) Comdat {
	const op = "ir.Module.GetOrInsertComdat"
	m.Check(op)
	if name == "" {
		llvm.Panicf(llvm.ErrInvalidArg, op, "empty comdat name")
	}
	ref := binding.LLVMGetOrInsertComdat(m.ref, name)
	return Comdat{ref: ref, ctx: m.ctx, life: m.life, name: name}
}

// Comdat 全局变量所属 comdat（未设置返回 false）
func (g Global) Comdat() (Comdat, bool) {
	g.Check("ir.Global.Comdat")
	ref := binding.LLVMGetComdat(g.Ref())
	if ref.IsNil() {
		return Comdat{}, false
	}
	return Comdat{ref: ref, ctx: g.Context(), life: g.Lifetime()}, true
}

// SetComdat 给全局变量设置 comdat
func (g Global) SetComdat(c Comdat) {
	const op = "ir.Global.SetComdat"
	g.Check(op)
	c.Check(op)
	if c.ctx != g.Context() {
		llvm.Panicf(llvm.ErrCrossContext, op, "comdat belongs to another context")
	}
	binding.LLVMSetComdat(g.Ref(), c.ref)
}
