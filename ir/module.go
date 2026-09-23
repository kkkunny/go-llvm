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

// newModule 由底层句柄铸造模块并登记到 Context 生命周期
func newModule(ctx *llvm.Context, ref binding.LLVMModuleRef) *Module {
	m := &Module{
		ref:  ref,
		ctx:  ctx,
		life: llvm.NewLifetime(),
	}
	m.unown = ctx.Own(m)
	return m
}

// NewModule 创建模块并登记到 Context 生命周期
func NewModule(ctx *llvm.Context, name string) *Module {
	if !ctx.Alive() {
		llvm.Panicf(llvm.ErrUseAfterFree, "ir.NewModule", "context is closed")
	}
	return newModule(ctx, binding.LLVMModuleCreateWithNameInContext(name, ctx.Ref()))
}

// Check 模块可用性前置校验（nil/所有权已移交/已关闭）
func (m *Module) Check(op string) {
	if m == nil || m.ref.IsNil() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil module")
	}
	if m.disowned {
		llvm.Panicf(llvm.ErrUseAfterFree, op, "module ownership has been transferred")
	}
	if m.closed || !m.life.Alive() {
		llvm.Panicf(llvm.ErrUseAfterFree, op, "module is closed")
	}
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

// Disown 立即解除与 Context 的级联所有权，把模块移交给外部接管方（如 JIT）：
// Go 侧句柄立即失效（此后使用 panic），底层模块由接管方释放；此后 Close 返回 ErrClosed。
func (m *Module) Disown() {
	if m.disowned || m.closed {
		return
	}
	m.disowned = true
	m.unown()
	m.life.Kill()
}

// Context 返回所属上下文
func (m *Module) Context() *llvm.Context { return m.ctx }

// Lifetime 返回模块生命周期令牌
func (m *Module) Lifetime() *llvm.Lifetime { return m.life }

// Ref 返回底层句柄（供包内桥接使用）
func (m *Module) Ref() binding.LLVMModuleRef { return m.ref }

// String 模块 IR 文本
func (m *Module) String() string {
	m.Check("ir.Module.String")
	return binding.LLVMPrintModuleToString(m.ref)
}

// Source 模块源文件名
func (m *Module) Source() string {
	m.Check("ir.Module.Source")
	return binding.LLVMGetSourceFileName(m.ref)
}

// SetSource 设置模块源文件名
func (m *Module) SetSource(source string) {
	m.Check("ir.Module.SetSource")
	binding.LLVMSetSourceFileName(m.ref, source)
}

// TargetTriple 目标三元组
func (m *Module) TargetTriple() string {
	m.Check("ir.Module.TargetTriple")
	return binding.LLVMGetTarget(m.ref)
}

// SetTargetTriple 设置目标三元组
func (m *Module) SetTargetTriple(triple string) {
	m.Check("ir.Module.SetTargetTriple")
	binding.LLVMSetTarget(m.ref, triple)
}

// SetDataLayout 设置数据布局
func (m *Module) SetDataLayout(layout string) {
	m.Check("ir.Module.SetDataLayout")
	binding.LLVMSetDataLayout(m.ref, layout)
}

// Verify 校验模块；失败返回 ErrVerify 与完整诊断
func (m *Module) Verify() error {
	m.Check("ir.Module.Verify")
	msg, failed := binding.LLVMVerifyModule(m.ref, binding.LLVMReturnStatusAction)
	if failed {
		return &llvm.Error{Reason: llvm.ErrVerify, Op: "ir.Module.Verify", Msg: msg}
	}
	return nil
}

// Clone 深拷贝模块
func (m *Module) Clone() *Module {
	m.Check("ir.Module.Clone")
	return newModule(m.ctx, binding.LLVMCloneModule(m.ref))
}

// DataLayout 模块数据布局（owned 副本，用毕 Close）
func (m *Module) DataLayout() *llvm.DataLayout {
	m.Check("ir.Module.DataLayout")
	return llvm.NewDataLayout(binding.LLVMGetDataLayoutStr(m.ref))
}

// WriteToFile 将模块 IR 文本写入文件；失败返回 ErrIO
func (m *Module) WriteToFile(path string) error {
	m.Check("ir.Module.WriteToFile")
	if err := binding.LLVMPrintModuleToFile(m.ref, path); err != nil {
		return llvm.WrapError(llvm.ErrIO, "ir.Module.WriteToFile", err)
	}
	return nil
}

// WriteBitcode 将模块 bitcode 写入文件；失败返回 ErrIO
func (m *Module) WriteBitcode(path string) error {
	m.Check("ir.Module.WriteBitcode")
	if err := binding.LLVMWriteBitcodeToFile(m.ref, path); err != nil {
		return llvm.WrapError(llvm.ErrIO, "ir.Module.WriteBitcode", err)
	}
	return nil
}

// Bitcode 将模块序列化为 bitcode 内存缓冲
func (m *Module) Bitcode() *llvm.MemoryBuffer {
	m.Check("ir.Module.Bitcode")
	return llvm.MemoryBufferOf(binding.LLVMWriteBitcodeToMemoryBuffer(m.ref))
}

// NewFunction 按函数类型声明/定义函数
func (m *Module) NewFunction(name string, t llvm.FnType) Function {
	const op = "ir.Module.NewFunction"
	m.Check(op)
	t.Check(op)
	if t.Context() != m.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "function type belongs to another context")
	}
	ref := binding.LLVMAddFunction(m.ref, name, t.Ref())
	return Function{Value: llvm.NewValue[llvm.FnT](m.ctx, m.life, ref)}
}

// GetFunction 按名称查找函数
func (m *Module) GetFunction(name string) (Function, bool) {
	m.Check("ir.Module.GetFunction")
	ref := binding.LLVMGetNamedFunction(m.ref, name)
	if ref.IsNil() {
		return Function{}, false
	}
	return Function{Value: llvm.NewValue[llvm.FnT](m.ctx, m.life, ref)}, true
}

// NewGlobal 声明全局变量（无初始化器）
func (m *Module) NewGlobal(name string, t llvm.AnyType) Global {
	const op = "ir.Module.NewGlobal"
	m.Check(op)
	m.ctx.CheckType(op, t)
	ref := binding.LLVMAddGlobal(m.ref, t.Ref(), name)
	return Global{Value: llvm.NewValue[llvm.PtrT](m.ctx, m.life, ref)}
}

// NewGlobalConst 声明常量全局变量
func (m *Module) NewGlobalConst(name string, v llvm.AnyValue) Global {
	const op = "ir.Module.NewGlobalConst"
	m.Check(op)
	m.ctx.CheckValues(op, v)
	g := m.NewGlobal(name, llvm.TypeOfRef(m.ctx, binding.LLVMTypeOf(v.Ref())))
	g.SetInitializer(v)
	g.SetConstant(true)
	return g
}

// GetGlobal 按名称查找全局变量
func (m *Module) GetGlobal(name string) (Global, bool) {
	m.Check("ir.Module.GetGlobal")
	ref := binding.LLVMGetNamedGlobal(m.ref, name)
	if ref.IsNil() {
		return Global{}, false
	}
	return Global{Value: llvm.NewValue[llvm.PtrT](m.ctx, m.life, ref)}, true
}

// DelGlobal 删除全局变量
func (m *Module) DelGlobal(g Global) {
	const op = "ir.Module.DelGlobal"
	m.Check(op)
	m.ctx.CheckValues(op, g)
	binding.LLVMDeleteGlobal(g.Ref())
}
