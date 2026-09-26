package target

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Arch 目标架构
type Arch uint8

// Arch 取值对应 LLVM 内置的目标架构（LLVM target backend）；各架构可用的组件
// （info/target/MC/asm printer/asm parser/disassembler）由包内初始化表决定。
const (
	AArch64     Arch = iota // ARM 64 位架构
	AMDGPU                  // AMD GPU 架构
	ARM                     // ARM 32 位架构
	AVR                     // Atmel AVR 8 位微控制器架构
	BPF                     // eBPF 架构（Linux 内核可编程字节码）
	Hexagon                 // Qualcomm Hexagon DSP 架构
	Lanai                   // Lanai 架构（Google 内部使用的 32 位处理器）
	LoongArch               // 龙芯 LoongArch 架构
	Mips                    // MIPS 架构
	MSP430                  // 德州仪器 MSP430 16 位微控制器架构
	NVPTX                   // NVIDIA PTX 虚拟指令集（CUDA GPU 目标）
	PowerPC                 // PowerPC 架构
	RISCV                   // RISC-V 架构
	Sparc                   // Sun SPARC 架构
	SystemZ                 // IBM z/Architecture 大型机架构
	VE                      // NEC SX-Aurora TSUBASA 向量引擎架构
	WebAssembly             // WebAssembly 架构
	X86                     // Intel x86 / x86-64 架构
	XCore                   // XMOS XCore 架构
)

// archInit 单架构初始化函数表；nil 表示该组件在 binding 中不可用
type archInit struct {
	info         func()
	target       func()
	mc           func()
	asmPrinter   func()
	asmParser    func()
	disassembler func()
}

var archInits = map[Arch]archInit{
	AArch64:     {binding.LLVMInitializeAArch64TargetInfo, binding.LLVMInitializeAArch64Target, binding.LLVMInitializeAArch64TargetMC, binding.LLVMInitializeAArch64AsmPrinter, binding.LLVMInitializeAArch64AsmParser, binding.LLVMInitializeAArch64Disassembler},
	AMDGPU:      {binding.LLVMInitializeAMDGPUTargetInfo, binding.LLVMInitializeAMDGPUTarget, binding.LLVMInitializeAMDGPUTargetMC, binding.LLVMInitializeAMDGPUAsmPrinter, binding.LLVMInitializeAMDGPUAsmParser, binding.LLVMInitializeAMDGPUDisassembler},
	ARM:         {binding.LLVMInitializeARMTargetInfo, binding.LLVMInitializeARMTarget, binding.LLVMInitializeARMTargetMC, binding.LLVMInitializeARMAsmPrinter, binding.LLVMInitializeARMAsmParser, binding.LLVMInitializeARMDisassembler},
	AVR:         {binding.LLVMInitializeAVRTargetInfo, binding.LLVMInitializeAVRTarget, binding.LLVMInitializeAVRTargetMC, binding.LLVMInitializeAVRAsmPrinter, binding.LLVMInitializeAVRAsmParser, binding.LLVMInitializeAVRDisassembler},
	BPF:         {binding.LLVMInitializeBPFTargetInfo, binding.LLVMInitializeBPFTarget, binding.LLVMInitializeBPFTargetMC, binding.LLVMInitializeBPFAsmPrinter, binding.LLVMInitializeBPFAsmParser, binding.LLVMInitializeBPFDisassembler},
	Hexagon:     {binding.LLVMInitializeHexagonTargetInfo, binding.LLVMInitializeHexagonTarget, binding.LLVMInitializeHexagonTargetMC, binding.LLVMInitializeHexagonAsmPrinter, binding.LLVMInitializeHexagonAsmParser, binding.LLVMInitializeHexagonDisassembler},
	Lanai:       {binding.LLVMInitializeLanaiTargetInfo, binding.LLVMInitializeLanaiTarget, binding.LLVMInitializeLanaiTargetMC, binding.LLVMInitializeLanaiAsmPrinter, binding.LLVMInitializeLanaiAsmParser, binding.LLVMInitializeLanaiDisassembler},
	LoongArch:   {binding.LLVMInitializeLoongArchTargetInfo, binding.LLVMInitializeLoongArchTarget, binding.LLVMInitializeLoongArchTargetMC, binding.LLVMInitializeLoongArchAsmPrinter, binding.LLVMInitializeLoongArchAsmParser, binding.LLVMInitializeLoongArchDisassembler},
	Mips:        {binding.LLVMInitializeMipsTargetInfo, binding.LLVMInitializeMipsTarget, binding.LLVMInitializeMipsTargetMC, binding.LLVMInitializeMipsAsmPrinter, binding.LLVMInitializeMipsAsmParser, binding.LLVMInitializeMipsDisassembler},
	MSP430:      {binding.LLVMInitializeMSP430TargetInfo, binding.LLVMInitializeMSP430Target, binding.LLVMInitializeMSP430TargetMC, binding.LLVMInitializeMSP430AsmPrinter, binding.LLVMInitializeMSP430AsmParser, binding.LLVMInitializeMSP430Disassembler},
	NVPTX:       {binding.LLVMInitializeNVPTXTargetInfo, binding.LLVMInitializeNVPTXTarget, binding.LLVMInitializeNVPTXTargetMC, nil, nil, nil},
	PowerPC:     {binding.LLVMInitializePowerPCTargetInfo, binding.LLVMInitializePowerPCTarget, binding.LLVMInitializePowerPCTargetMC, binding.LLVMInitializePowerPCAsmPrinter, binding.LLVMInitializePowerPCAsmParser, binding.LLVMInitializePowerPCDisassembler},
	RISCV:       {binding.LLVMInitializeRISCVTargetInfo, binding.LLVMInitializeRISCVTarget, binding.LLVMInitializeRISCVTargetMC, binding.LLVMInitializeRISCVAsmPrinter, binding.LLVMInitializeRISCVAsmParser, binding.LLVMInitializeRISCVDisassembler},
	Sparc:       {binding.LLVMInitializeSparcTargetInfo, binding.LLVMInitializeSparcTarget, binding.LLVMInitializeSparcTargetMC, binding.LLVMInitializeSparcAsmPrinter, binding.LLVMInitializeSparcAsmParser, binding.LLVMInitializeSparcDisassembler},
	SystemZ:     {binding.LLVMInitializeSystemZTargetInfo, binding.LLVMInitializeSystemZTarget, binding.LLVMInitializeSystemZTargetMC, binding.LLVMInitializeSystemZAsmPrinter, binding.LLVMInitializeSystemZAsmParser, binding.LLVMInitializeSystemZDisassembler},
	VE:          {binding.LLVMInitializeVETargetInfo, binding.LLVMInitializeVETarget, binding.LLVMInitializeVETargetMC, binding.LLVMInitializeVEAsmPrinter, binding.LLVMInitializeVEAsmParser, binding.LLVMInitializeVEDisassembler},
	WebAssembly: {binding.LLVMInitializeWebAssemblyTargetInfo, binding.LLVMInitializeWebAssemblyTarget, binding.LLVMInitializeWebAssemblyTargetMC, binding.LLVMInitializeWebAssemblyAsmPrinter, binding.LLVMInitializeWebAssemblyAsmParser, binding.LLVMInitializeWebAssemblyDisassembler},
	X86:         {binding.LLVMInitializeX86TargetInfo, binding.LLVMInitializeX86Target, binding.LLVMInitializeX86TargetMC, binding.LLVMInitializeX86AsmPrinter, binding.LLVMInitializeX86AsmParser, binding.LLVMInitializeX86Disassembler},
	XCore:       {binding.LLVMInitializeXCoreTargetInfo, binding.LLVMInitializeXCoreTarget, binding.LLVMInitializeXCoreTargetMC, nil, nil, binding.LLVMInitializeXCoreDisassembler},
}

// InitAll 初始化全部目标（infos/targets/MCs/asm printers/asm parsers/disassemblers）
func InitAll() {
	binding.LLVMInitializeAllTargetInfos()
	binding.LLVMInitializeAllTargets()
	binding.LLVMInitializeAllTargetMCs()
	binding.LLVMInitializeAllAsmPrinters()
	binding.LLVMInitializeAllAsmParsers()
	binding.LLVMInitializeAllDisassemblers()
}

// InitNative 初始化宿主目标；失败返回 ErrCodeGen
func InitNative() error {
	for _, f := range []struct {
		op string
		fn func() error
	}{
		{"target.InitNative", binding.LLVMInitializeNativeTarget},
		{"target.InitNativeAsmPrinter", binding.LLVMInitializeNativeAsmPrinter},
		{"target.InitNativeAsmParser", binding.LLVMInitializeNativeAsmParser},
		{"target.InitNativeDisassembler", binding.LLVMInitializeNativeDisassembler},
	} {
		if err := f.fn(); err != nil {
			return llvm.WrapError(llvm.ErrCodeGen, f.op, err)
		}
	}
	return nil
}

// Init 初始化指定架构的全部可用组件；未知架构 panic ErrInvalidArg
func Init(arch Arch) {
	init, ok := archInits[arch]
	if !ok {
		llvm.Panicf(llvm.ErrInvalidArg, "target.Init", "unknown arch %d", arch)
	}
	init.info()
	init.target()
	init.mc()
	if init.asmPrinter != nil {
		init.asmPrinter()
	}
	if init.asmParser != nil {
		init.asmParser()
	}
	if init.disassembler != nil {
		init.disassembler()
	}
}

// Target 目标描述（初始化后全局唯一，无需释放）
type Target struct{ ref binding.LLVMTargetRef }

// Ref 返回底层句柄（供 llvm/target 内部桥接使用）
func (t Target) Ref() binding.LLVMTargetRef { return t.ref }

// FromName 按目标名查找（如 "x86-64"）
func FromName(name string) (Target, bool) {
	ref := binding.LLVMGetTargetFromName(name)
	if ref.IsNil() {
		return Target{}, false
	}
	return Target{ref: ref}, true
}

// FromTriple 按三元组查找；未知三元组返回 ErrNotFound
func FromTriple(triple string) (Target, error) {
	ref, err := binding.LLVMGetTargetFromTriple(triple)
	if err != nil {
		return Target{}, llvm.WrapError(llvm.ErrNotFound, "target.FromTriple", err)
	}
	return Target{ref: ref}, nil
}

// NativeTarget 宿主目标
func NativeTarget() (Target, error) { return FromTriple(DefaultTriple()) }

// Name 目标名
func (t Target) Name() string { return binding.LLVMGetTargetName(t.ref) }

// Description 目标描述
func (t Target) Description() string { return binding.LLVMGetTargetDescription(t.ref) }

// HasJIT 是否支持 JIT
func (t Target) HasJIT() bool { return binding.LLVMTargetHasJIT(t.ref) }

// HasTargetMachine 是否支持目标机器
func (t Target) HasTargetMachine() bool { return binding.LLVMTargetHasTargetMachine(t.ref) }

// HasAsmBackend 是否有汇编后端
func (t Target) HasAsmBackend() bool { return binding.LLVMTargetHasAsmBackend(t.ref) }

// DefaultTriple 宿主三元组
func DefaultTriple() string { return binding.LLVMGetDefaultTargetTriple() }

// NormalizeTriple 规范化三元组
func NormalizeTriple(triple string) string { return binding.LLVMNormalizeTargetTriple(triple) }

// HostCPUName 宿主 CPU 名
func HostCPUName() string { return binding.LLVMGetHostCPUName() }

// HostCPUFeatures 宿主 CPU 特性串
func HostCPUFeatures() string { return binding.LLVMGetHostCPUFeatures() }
