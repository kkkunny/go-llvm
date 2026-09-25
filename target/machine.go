package target

import (
	"errors"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
	"github.com/kkkunny/go-llvm/ir"
)

// OptLevel 代码生成优化级别
type OptLevel int32

// OptLevel 取值对应 LLVM 代码生成优化级别（binding.LLVMCodeGenOptLevel），
// 影响后端指令选择与调度等优化强度；IR 层优化由 pass 管线负责。
const (
	OptNone       OptLevel = OptLevel(binding.LLVMCodeGenLevelNone)       // 不优化（-O0）
	OptLess       OptLevel = OptLevel(binding.LLVMCodeGenLevelLess)       // 轻度优化（-O1）
	OptDefault    OptLevel = OptLevel(binding.LLVMCodeGenLevelDefault)    // 默认优化（-O2，-Os/-Oz 亦映射到本级）
	OptAggressive OptLevel = OptLevel(binding.LLVMCodeGenLevelAggressive) // 激进优化（-O3）
)

// RelocMode 重定位模式
type RelocMode int32

// RelocMode 取值对应 LLVM 重定位模式（binding.LLVMRelocMode），决定生成代码如何
// 寻址全局符号与函数；ROPI/RWPI 系列主要供 ARM 目标使用。
const (
	RelocDefault      RelocMode = RelocMode(binding.LLVMRelocDefault)      // 目标默认：由目标后端选择重定位模型
	RelocStatic       RelocMode = RelocMode(binding.LLVMRelocStatic)       // 静态：地址在链接时固定，生成非位置无关代码
	RelocPIC          RelocMode = RelocMode(binding.LLVMRelocPIC)          // 位置无关：代码可加载到任意地址（-fPIC）
	RelocDynamicNoPic RelocMode = RelocMode(binding.LLVMRelocDynamicNoPic) // 动态非 PIC：非位置无关，但可放入动态可执行文件（Darwin 的 -mdynamic-no-pic）
	RelocROPI         RelocMode = RelocMode(binding.LLVMRelocROPI)         // 只读位置无关（ROPI）：只读段位置无关，可写数据位于固定地址
	RelocRWPI         RelocMode = RelocMode(binding.LLVMRelocRWPI)         // 可写位置无关（RWPI）：可写数据经静态基址寄存器寻址，只读段位于固定地址
	RelocROPI_RWPI    RelocMode = RelocMode(binding.LLVMRelocROPI_RWPI)    // ROPI 与 RWPI 兼具：只读段与可写数据均位置无关
)

// CodeModel 代码模型
type CodeModel int32

// CodeModel 取值对应 LLVM 代码模型（binding.LLVMCodeModel），描述代码与数据的
// 大小及地址范围假设，影响寻址指令的选取。
const (
	CodeModelDefault    CodeModel = CodeModel(binding.LLVMCodeModelDefault)    // 目标默认：由目标后端选择代码模型
	CodeModelJITDefault CodeModel = CodeModel(binding.LLVMCodeModelJITDefault) // JIT 默认：由目标为 JIT 代码选择默认代码模型
	CodeModelTiny       CodeModel = CodeModel(binding.LLVMCodeModelTiny)       // 微模型：假设代码与数据位于 16 位地址空间内（仅部分目标支持）
	CodeModelSmall      CodeModel = CodeModel(binding.LLVMCodeModelSmall)      // 小模型：假设代码与数据位于前 2 GiB 地址空间内（多数目标的默认）
	CodeModelKernel     CodeModel = CodeModel(binding.LLVMCodeModelKernel)     // 内核模型：代码位于高地址区（如 x86-64 的负 2 GiB），供操作系统内核使用
	CodeModelMedium     CodeModel = CodeModel(binding.LLVMCodeModelMedium)     // 中模型：代码位于前 2 GiB，数据可位于任意地址
	CodeModelLarge      CodeModel = CodeModel(binding.LLVMCodeModelLarge)      // 大模型：对代码与数据的大小和地址不做假设
)

// FileType 产出文件类型
type FileType int32

// FileType 取值对应 LLVM 代码生成输出文件类型（binding.LLVMCodeGenFileType）。
const (
	AsmFile    FileType = FileType(binding.LLVMAssemblyFile) // 汇编文本文件（.s）
	ObjectFile FileType = FileType(binding.LLVMObjectFile)   // 目标文件（.o/.obj）
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

// ApplyTo 把目标三元组与数据布局写入模块（旧 Module.SetTarget 的替代）
func (m *TargetMachine) ApplyTo(mod *ir.Module) {
	const op = "target.TargetMachine.ApplyTo"
	m.check(op)
	mod.Check(op)
	mod.SetTargetTriple(m.Triple())
	dl := m.DataLayout()
	defer dl.Close()
	mod.SetDataLayout(dl.String())
}

// EmitToFile 将模块编译为汇编/目标文件；失败返回 ErrCodeGen
func (m *TargetMachine) EmitToFile(mod *ir.Module, path string, ft FileType) error {
	const op = "target.TargetMachine.EmitToFile"
	m.check(op)
	mod.Check(op)
	if checks.Debug {
		if err := mod.Verify(); err != nil {
			llvm.Panicf(llvm.ErrVerify, op, "module verification failed before codegen: %s", err)
		}
	}
	if err := binding.LLVMTargetMachineEmitToFile(m.ref, mod.Ref(), path, binding.LLVMCodeGenFileType(ft)); err != nil {
		return llvm.WrapError(llvm.ErrCodeGen, op, err)
	}
	return nil
}

// Emit 将模块编译为汇编/目标代码内存缓冲
func (m *TargetMachine) Emit(mod *ir.Module, ft FileType) (*llvm.MemoryBuffer, error) {
	const op = "target.TargetMachine.Emit"
	m.check(op)
	mod.Check(op)
	buf, err := binding.LLVMTargetMachineEmitToMemoryBuffer(m.ref, mod.Ref(), binding.LLVMCodeGenFileType(ft))
	if err != nil {
		return nil, llvm.WrapError(llvm.ErrCodeGen, op, err)
	}
	return llvm.MemoryBufferOf(buf), nil
}
