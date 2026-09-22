package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Module LLVM 模块；Close 后其下所有值失效
type Module struct {
	ref      binding.LLVMModuleRef
	ctx      *llvm.Context
	life     *llvm.Lifetime
	unown    func()
	closed   bool
	disowned bool
}

// NewModule 创建模块并登记到 Context 生命周期
func NewModule(ctx *llvm.Context, name string) *Module {
	m := &Module{
		ref:  binding.LLVMModuleCreateWithNameInContext(name, ctx.Ref()),
		ctx:  ctx,
		life: llvm.NewLifetime(),
	}
	m.unown = ctx.Own(m)
	return m
}

// Close 释放模块；其下值随之失效，二次调用返回 ErrClosed
func (m *Module) Close() error {
	if m.closed {
		return &llvm.Error{Reason: llvm.ErrClosed, Op: "ir.Module.Close", Msg: "module already closed"}
	}
	if m.disowned {
		return &llvm.Error{Reason: llvm.ErrClosed, Op: "ir.Module.Close", Msg: "module ownership has been transferred"}
	}
	m.closed = true
	m.unown()
	m.life.Kill()
	binding.LLVMDisposeModule(m.ref)
	return nil
}

// Disown 解除与 Context 的级联所有权，把模块移交给外部接管方（如 JIT）。
// 返回 release：接管方释放底层模块时调用它使 Go 侧句柄失效；此后 Close 返回 ErrClosed。
func (m *Module) Disown() func() {
	if m.disowned || m.closed {
		return func() {}
	}
	m.disowned = true
	m.unown()
	return m.life.Kill
}

// Context 返回所属上下文
func (m *Module) Context() *llvm.Context { return m.ctx }

// Lifetime 返回模块生命周期令牌
func (m *Module) Lifetime() *llvm.Lifetime { return m.life }

// Ref 返回底层句柄（供包内桥接使用）
func (m *Module) Ref() binding.LLVMModuleRef { return m.ref }

// String 模块 IR 文本
func (m *Module) String() string { return binding.LLVMPrintModuleToString(m.ref) }

// Source 模块源文件名
func (m *Module) Source() string { return binding.LLVMGetSourceFileName(m.ref) }

// SetSource 设置模块源文件名
func (m *Module) SetSource(source string) { binding.LLVMSetSourceFileName(m.ref, source) }

// TargetTriple 目标三元组
func (m *Module) TargetTriple() string { return binding.LLVMGetTarget(m.ref) }

// SetTargetTriple 设置目标三元组
func (m *Module) SetTargetTriple(triple string) { binding.LLVMSetTarget(m.ref, triple) }

// SetDataLayout 设置数据布局
func (m *Module) SetDataLayout(layout string) { binding.LLVMSetDataLayout(m.ref, layout) }

// Verify 校验模块；失败返回 ErrVerify 与完整诊断
func (m *Module) Verify() error {
	msg, failed := binding.LLVMVerifyModule(m.ref, binding.LLVMReturnStatusAction)
	if failed {
		return &llvm.Error{Reason: llvm.ErrVerify, Op: "ir.Module.Verify", Msg: msg}
	}
	return nil
}

// Clone 深拷贝模块
func (m *Module) Clone() *Module {
	clone := &Module{
		ref:  binding.LLVMCloneModule(m.ref),
		ctx:  m.ctx,
		life: llvm.NewLifetime(),
	}
	clone.unown = m.ctx.Own(clone)
	return clone
}

// NewFunction 按函数类型声明/定义函数
func (m *Module) NewFunction(name string, t llvm.FnType) Function {
	ref := binding.LLVMAddFunction(m.ref, name, t.Ref())
	return Function{Value: wrapValue[llvm.FnT](m.ctx, m.life, ref)}
}

// GetFunction 按名称查找函数
func (m *Module) GetFunction(name string) (Function, bool) {
	ref := binding.LLVMGetNamedFunction(m.ref, name)
	if ref.IsNil() {
		return Function{}, false
	}
	return Function{Value: wrapValue[llvm.FnT](m.ctx, m.life, ref)}, true
}

// NewGlobal 声明全局变量（无初始化器）
func (m *Module) NewGlobal(name string, t llvm.AnyType) Global {
	ref := binding.LLVMAddGlobal(m.ref, t.Ref(), name)
	return Global{Value: wrapValue[llvm.PtrT](m.ctx, m.life, ref)}
}

// NewConstant 声明常量全局变量
func (m *Module) NewConstant(name string, v llvm.AnyValue) Global {
	g := m.NewGlobal(name, llvm.TypeOfRef(m.ctx, binding.LLVMTypeOf(v.Ref())))
	g.SetInitializer(v)
	g.SetConstant(true)
	return g
}

// GetGlobal 按名称查找全局变量
func (m *Module) GetGlobal(name string) (Global, bool) {
	ref := binding.LLVMGetNamedGlobal(m.ref, name)
	if ref.IsNil() {
		return Global{}, false
	}
	return Global{Value: wrapValue[llvm.PtrT](m.ctx, m.life, ref)}, true
}

// DelGlobal 删除全局变量
func (m *Module) DelGlobal(g Global) {
	binding.LLVMDeleteGlobal(g.Ref())
}
