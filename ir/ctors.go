package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// AppendCtor 向 llvm.global_ctors 追加构造器（priority 越小越先执行，clang 惯例 65535）
func (m *Module) AppendCtor(fn Function, priority uint32) {
	m.appendGlobalList("ir.Module.AppendCtor", "llvm.global_ctors", fn, priority)
}

// AppendDtor 向 llvm.global_dtors 追加析构器（priority 越小越先执行）
func (m *Module) AppendDtor(fn Function, priority uint32) {
	m.appendGlobalList("ir.Module.AppendDtor", "llvm.global_dtors", fn, priority)
}

// appendGlobalList 读取既有 {i32, ptr, ptr} 数组并追加一项，重建 appending 全局。
// LLVM 不允许 GlobalVariable 改类型，故先删后建
func (m *Module) appendGlobalList(op, name string, fn Function, priority uint32) {
	m.Check(op)
	fn.Check(op)
	if fn.Context() != m.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "function belongs to another context")
	}
	i32, ptr := m.ctx.Int(32), m.ctx.Ptr(0)
	elemTy := m.ctx.Struct([]llvm.AnyType{i32, ptr, ptr}, false)

	var elems []llvm.AnyValue
	if g, ok := m.GetGlobal(name); ok {
		if init, has := g.Initializer(); has {
			n := binding.LLVMGetArrayLength2(binding.LLVMTypeOf(init.Ref()))
			for i := uint64(0); i < n; i++ {
				elems = append(elems, llvm.ValueOf(m.ctx, m.life, binding.LLVMGetAggregateElement(init.Ref(), uint32(i))))
			}
		}
		m.DelGlobal(g)
	}
	entry := m.ctx.ConstStruct(false, m.ctx.ConstInt(i32, uint64(priority)).Value, fn, m.ctx.ConstNull(ptr))
	elems = append(elems, entry)

	g := m.NewGlobal(name, m.ctx.Array(elemTy, uint64(len(elems))))
	g.SetInitializer(m.ctx.ConstArray(elemTy, elems...))
	g.SetLinkage(llvm.LinkageAppending)
}
