package target

import (
	"errors"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/ir"
)

// OptLevel 代码生成优化级别
type OptLevel int32

const (
	OptNone       OptLevel = OptLevel(binding.LLVMCodeGenLevelNone)
	OptLess       OptLevel = OptLevel(binding.LLVMCodeGenLevelLess)
	OptDefault    OptLevel = OptLevel(binding.LLVMCodeGenLevelDefault)
	OptAggressive OptLevel = OptLevel(binding.LLVMCodeGenLevelAggressive)
)

// RelocMode 重定位模式
type RelocMode int32

const (
	RelocDefault      RelocMode = RelocMode(binding.LLVMRelocDefault)
	RelocStatic       RelocMode = RelocMode(binding.LLVMRelocStatic)
	RelocPIC          RelocMode = RelocMode(binding.LLVMRelocPIC)
	RelocDynamicNoPic RelocMode = RelocMode(binding.LLVMRelocDynamicNoPic)
	RelocROPI         RelocMode = RelocMode(binding.LLVMRelocROPI)
	RelocRWPI         RelocMode = RelocMode(binding.LLVMRelocRWPI)
	RelocROPI_RWPI    RelocMode = RelocMode(binding.LLVMRelocROPI_RWPI)
)

// CodeModel 代码模型
type CodeModel int32

const (
	CodeModelDefault    CodeModel = CodeModel(binding.LLVMCodeModelDefault)
	CodeModelJITDefault CodeModel = CodeModel(binding.LLVMCodeModelJITDefault)
	CodeModelTiny       CodeModel = CodeModel(binding.LLVMCodeModelTiny)
	CodeModelSmall      CodeModel = CodeModel(binding.LLVMCodeModelSmall)
	CodeModelKernel     CodeModel = CodeModel(binding.LLVMCodeModelKernel)
	CodeModelMedium     CodeModel = CodeModel(binding.LLVMCodeModelMedium)
	CodeModelLarge      CodeModel = CodeModel(binding.LLVMCodeModelLarge)
)

// FileType 产出文件类型
type FileType int32

const (
	AsmFile    FileType = FileType(binding.LLVMAssemblyFile)
	ObjectFile FileType = FileType(binding.LLVMObjectFile)
)

// TargetMachine 目标机器；独立所有权根（不经 Context.Own），用毕 Close
type TargetMachine struct {
	ref    binding.LLVMTargetMachineRef
	closed bool
}

// NewTargetMachine 创建目标机器；失败返回 ErrCodeGen
func NewTargetMachine(t Target, triple, cpu, features string, opt OptLevel, reloc RelocMode, cm CodeModel) (*TargetMachine, error) {
	ref := binding.LLVMCreateTargetMachine(t.ref, triple, cpu, features, binding.LLVMCodeGenOptLevel(opt), binding.LLVMRelocMode(reloc), binding.LLVMCodeModel(cm))
	if ref.IsNil() {
		return nil, llvm.WrapError(llvm.ErrCodeGen, "target.NewTargetMachine", errors.New("LLVMCreateTargetMachine returned null"))
	}
	return &TargetMachine{ref: ref}, nil
}

// check 前置校验：未释放
func (m *TargetMachine) check(op string) {
	if m.closed {
		llvm.Panicf(llvm.ErrUseAfterFree, op, "target machine is closed")
	}
}

// checkModule 前置校验模块可用
func checkModule(op string, mod *ir.Module) {
	if mod == nil || mod.Ref().IsNil() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil module")
	}
	if !mod.Lifetime().Alive() {
		llvm.Panicf(llvm.ErrUseAfterFree, op, "module is closed")
	}
}

// Close 释放目标机器；二次调用返回 ErrClosed
func (m *TargetMachine) Close() error {
	if m.closed {
		return &llvm.Error{Reason: llvm.ErrClosed, Op: "target.TargetMachine.Close", Msg: "target machine already closed"}
	}
	m.closed = true
	binding.LLVMDisposeTargetMachine(m.ref)
	return nil
}

// Triple 目标三元组
func (m *TargetMachine) Triple() string {
	m.check("target.TargetMachine.Triple")
	return binding.LLVMGetTargetMachineTriple(m.ref)
}

// CPU 目标 CPU
func (m *TargetMachine) CPU() string {
	m.check("target.TargetMachine.CPU")
	return binding.LLVMGetTargetMachineCPU(m.ref)
}

// Features 目标特性串
func (m *TargetMachine) Features() string {
	m.check("target.TargetMachine.Features")
	return binding.LLVMGetTargetMachineFeatureString(m.ref)
}

// Target 目标机器对应的目标
func (m *TargetMachine) Target() Target {
	m.check("target.TargetMachine.Target")
	return Target{ref: binding.LLVMGetTargetMachineTarget(m.ref)}
}

// DataLayout 目标数据布局（owned，用毕 Close）
func (m *TargetMachine) DataLayout() *llvm.DataLayout {
	m.check("target.TargetMachine.DataLayout")
	return llvm.DataLayoutOf(binding.LLVMCreateTargetDataLayout(m.ref))
}

// SetAsmVerbosity 设置汇编输出详细程度
func (m *TargetMachine) SetAsmVerbosity(verbose bool) {
	m.check("target.TargetMachine.SetAsmVerbosity")
	binding.LLVMSetTargetMachineAsmVerbosity(m.ref, verbose)
}

// SetTo 把目标三元组与数据布局写入模块（旧 Module.SetTarget 的替代）
func (m *TargetMachine) SetTo(mod *ir.Module) {
	const op = "target.TargetMachine.SetTo"
	m.check(op)
	checkModule(op, mod)
	mod.SetTargetTriple(m.Triple())
	dl := m.DataLayout()
	defer dl.Close()
	mod.SetDataLayout(dl.String())
}

// EmitToFile 将模块编译为汇编/目标文件；失败返回 ErrCodeGen
func (m *TargetMachine) EmitToFile(mod *ir.Module, path string, ft FileType) error {
	const op = "target.TargetMachine.EmitToFile"
	m.check(op)
	checkModule(op, mod)
	if err := binding.LLVMTargetMachineEmitToFile(m.ref, mod.Ref(), path, binding.LLVMCodeGenFileType(ft)); err != nil {
		return llvm.WrapError(llvm.ErrCodeGen, op, err)
	}
	return nil
}

// Emit 将模块编译为汇编/目标代码内存缓冲
func (m *TargetMachine) Emit(mod *ir.Module, ft FileType) (*llvm.MemoryBuffer, error) {
	const op = "target.TargetMachine.Emit"
	m.check(op)
	checkModule(op, mod)
	buf, err := binding.LLVMTargetMachineEmitToMemoryBuffer(m.ref, mod.Ref(), binding.LLVMCodeGenFileType(ft))
	if err != nil {
		return nil, llvm.WrapError(llvm.ErrCodeGen, op, err)
	}
	return llvm.MemoryBufferOf(buf), nil
}
