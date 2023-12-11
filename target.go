package llvm

import "C"
import (
	"errors"
	"strings"

	"github.com/kkkunny/go-llvm/internal/binding"
)

// Arch 架构
type Arch uint8

const (
	AArch64 Arch = iota
	AMDGPU
	ARM
	AVR
	BPF
	Hexagon
	Lanai
	LoongArch
	Mips
	MSP430
	NVPTX
	PowerPC
	RISCV
	Sparc
	SystemZ
	VE
	WebAssembly
	X86
	XCore
)

// ErrUnknownArch 未知的架构
var ErrUnknownArch = errors.New("unknown arch")

func InitializeAllTargetInfos() {
	binding.LLVMInitializeAllTargetInfos()
}

// InitializeTargetInfo 初始化目标信息
func InitializeTargetInfo(arch Arch) error {
	switch arch {
	case AArch64:
		binding.LLVMInitializeAArch64TargetInfo()
	case AMDGPU:
		binding.LLVMInitializeAMDGPUTargetInfo()
	case ARM:
		binding.LLVMInitializeARMTargetInfo()
	case AVR:
		binding.LLVMInitializeAVRTargetInfo()
	case BPF:
		binding.LLVMInitializeBPFTargetInfo()
	case Hexagon:
		binding.LLVMInitializeHexagonTargetInfo()
	case Lanai:
		binding.LLVMInitializeLanaiTargetInfo()
	case LoongArch:
		binding.LLVMInitializeLoongArchTargetInfo()
	case Mips:
		binding.LLVMInitializeMipsTargetInfo()
	case MSP430:
		binding.LLVMInitializeMSP430TargetInfo()
	case NVPTX:
		binding.LLVMInitializeNVPTXTargetInfo()
	case PowerPC:
		binding.LLVMInitializePowerPCTargetInfo()
	case RISCV:
		binding.LLVMInitializeRISCVTargetInfo()
	case Sparc:
		binding.LLVMInitializeSparcTargetInfo()
	case SystemZ:
		binding.LLVMInitializeSystemZTargetInfo()
	case VE:
		binding.LLVMInitializeVETargetInfo()
	case WebAssembly:
		binding.LLVMInitializeWebAssemblyTargetInfo()
	case X86:
		binding.LLVMInitializeX86TargetInfo()
	case XCore:
		binding.LLVMInitializeXCoreTargetInfo()
	default:
		return ErrUnknownArch
	}
	return nil
}

func InitializeAllTargets() {
	binding.LLVMInitializeAllTargets()
}

// InitializeTarget 初始化目标
func InitializeTarget(arch Arch) error {
	switch arch {
	case AArch64:
		binding.LLVMInitializeAArch64Target()
	case AMDGPU:
		binding.LLVMInitializeAMDGPUTarget()
	case ARM:
		binding.LLVMInitializeARMTarget()
	case AVR:
		binding.LLVMInitializeAVRTarget()
	case BPF:
		binding.LLVMInitializeBPFTarget()
	case Hexagon:
		binding.LLVMInitializeHexagonTarget()
	case Lanai:
		binding.LLVMInitializeLanaiTarget()
	case LoongArch:
		binding.LLVMInitializeLoongArchTarget()
	case Mips:
		binding.LLVMInitializeMipsTarget()
	case MSP430:
		binding.LLVMInitializeMSP430Target()
	case NVPTX:
		binding.LLVMInitializeNVPTXTarget()
	case PowerPC:
		binding.LLVMInitializePowerPCTarget()
	case RISCV:
		binding.LLVMInitializeRISCVTarget()
	case Sparc:
		binding.LLVMInitializeSparcTarget()
	case SystemZ:
		binding.LLVMInitializeSystemZTarget()
	case VE:
		binding.LLVMInitializeVETarget()
	case WebAssembly:
		binding.LLVMInitializeWebAssemblyTarget()
	case X86:
		binding.LLVMInitializeX86Target()
	case XCore:
		binding.LLVMInitializeXCoreTarget()
	default:
		return ErrUnknownArch
	}
	return nil
}

func InitializeAllTargetMCs() {
	binding.LLVMInitializeAllTargetMCs()
}

// InitializeTargetMC 初始化目标机器
func InitializeTargetMC(arch Arch) error {
	switch arch {
	case AArch64:
		binding.LLVMInitializeAArch64TargetMC()
	case AMDGPU:
		binding.LLVMInitializeAMDGPUTargetMC()
	case ARM:
		binding.LLVMInitializeARMTargetMC()
	case AVR:
		binding.LLVMInitializeAVRTargetMC()
	case BPF:
		binding.LLVMInitializeBPFTargetMC()
	case Hexagon:
		binding.LLVMInitializeHexagonTargetMC()
	case Lanai:
		binding.LLVMInitializeLanaiTargetMC()
	case LoongArch:
		binding.LLVMInitializeLoongArchTargetMC()
	case Mips:
		binding.LLVMInitializeMipsTargetMC()
	case MSP430:
		binding.LLVMInitializeMSP430TargetMC()
	case NVPTX:
		binding.LLVMInitializeNVPTXTargetMC()
	case PowerPC:
		binding.LLVMInitializePowerPCTargetMC()
	case RISCV:
		binding.LLVMInitializeRISCVTargetMC()
	case Sparc:
		binding.LLVMInitializeSparcTargetMC()
	case SystemZ:
		binding.LLVMInitializeSystemZTargetMC()
	case VE:
		binding.LLVMInitializeVETargetMC()
	case WebAssembly:
		binding.LLVMInitializeWebAssemblyTargetMC()
	case X86:
		binding.LLVMInitializeX86TargetMC()
	case XCore:
		binding.LLVMInitializeXCoreTargetMC()
	default:
		return ErrUnknownArch
	}
	return nil
}

func InitializeNativeTarget() error {
	return binding.LLVMInitializeNativeTarget()
}

func InitializeAllAsmPrinters() {
	binding.LLVMInitializeAllAsmPrinters()
}

// InitializeAsmPrinter 初始化汇编输出器
func InitializeAsmPrinter(arch Arch) error {
	switch arch {
	case AArch64:
		binding.LLVMInitializeAArch64AsmPrinter()
	case AMDGPU:
		binding.LLVMInitializeAMDGPUAsmPrinter()
	case ARM:
		binding.LLVMInitializeARMAsmPrinter()
	case AVR:
		binding.LLVMInitializeAVRAsmPrinter()
	case BPF:
		binding.LLVMInitializeBPFAsmPrinter()
	case Hexagon:
		binding.LLVMInitializeHexagonAsmPrinter()
	case Lanai:
		binding.LLVMInitializeLanaiAsmPrinter()
	case LoongArch:
		binding.LLVMInitializeLoongArchAsmPrinter()
	case Mips:
		binding.LLVMInitializeMipsAsmPrinter()
	case MSP430:
		binding.LLVMInitializeMSP430AsmPrinter()
	case PowerPC:
		binding.LLVMInitializePowerPCAsmPrinter()
	case RISCV:
		binding.LLVMInitializeRISCVAsmPrinter()
	case Sparc:
		binding.LLVMInitializeSparcAsmPrinter()
	case SystemZ:
		binding.LLVMInitializeSystemZAsmPrinter()
	case VE:
		binding.LLVMInitializeVEAsmPrinter()
	case WebAssembly:
		binding.LLVMInitializeWebAssemblyAsmPrinter()
	case X86:
		binding.LLVMInitializeX86AsmPrinter()
	default:
		return ErrUnknownArch
	}
	return nil
}

func InitializeNativeAsmPrinter() error {
	return binding.LLVMInitializeNativeAsmPrinter()
}

func InitializeAllAsmParsers() {
	binding.LLVMInitializeAllAsmParsers()
}

// InitializeAsmParser 初始化汇编格式化器
func InitializeAsmParser(arch Arch) error {
	switch arch {
	case AArch64:
		binding.LLVMInitializeAArch64AsmParser()
	case AMDGPU:
		binding.LLVMInitializeAMDGPUAsmParser()
	case ARM:
		binding.LLVMInitializeARMAsmParser()
	case AVR:
		binding.LLVMInitializeAVRAsmParser()
	case BPF:
		binding.LLVMInitializeBPFAsmParser()
	case Hexagon:
		binding.LLVMInitializeHexagonAsmParser()
	case Lanai:
		binding.LLVMInitializeLanaiAsmParser()
	case LoongArch:
		binding.LLVMInitializeLoongArchAsmParser()
	case Mips:
		binding.LLVMInitializeMipsAsmParser()
	case MSP430:
		binding.LLVMInitializeMSP430AsmParser()
	case PowerPC:
		binding.LLVMInitializePowerPCAsmParser()
	case RISCV:
		binding.LLVMInitializeRISCVAsmParser()
	case Sparc:
		binding.LLVMInitializeSparcAsmParser()
	case SystemZ:
		binding.LLVMInitializeSystemZAsmParser()
	case VE:
		binding.LLVMInitializeVEAsmParser()
	case WebAssembly:
		binding.LLVMInitializeWebAssemblyAsmParser()
	case X86:
		binding.LLVMInitializeX86AsmParser()
	default:
		return ErrUnknownArch
	}
	return nil
}

func InitializeNativeAsmParser() error {
	return binding.LLVMInitializeNativeAsmParser()
}

func InitializeAllDisassemblers() {
	binding.LLVMInitializeAllDisassemblers()
}

// InitializeDisassembler 初始化目标反汇编器
func InitializeDisassembler(arch Arch) error {
	switch arch {
	case AArch64:
		binding.LLVMInitializeAArch64Disassembler()
	case AMDGPU:
		binding.LLVMInitializeAMDGPUDisassembler()
	case ARM:
		binding.LLVMInitializeARMDisassembler()
	case AVR:
		binding.LLVMInitializeAVRDisassembler()
	case BPF:
		binding.LLVMInitializeBPFDisassembler()
	case Hexagon:
		binding.LLVMInitializeHexagonDisassembler()
	case Lanai:
		binding.LLVMInitializeLanaiDisassembler()
	case LoongArch:
		binding.LLVMInitializeLoongArchDisassembler()
	case Mips:
		binding.LLVMInitializeMipsDisassembler()
	case MSP430:
		binding.LLVMInitializeMSP430Disassembler()
	case PowerPC:
		binding.LLVMInitializePowerPCDisassembler()
	case RISCV:
		binding.LLVMInitializeRISCVDisassembler()
	case Sparc:
		binding.LLVMInitializeSparcDisassembler()
	case SystemZ:
		binding.LLVMInitializeSystemZDisassembler()
	case VE:
		binding.LLVMInitializeVEDisassembler()
	case WebAssembly:
		binding.LLVMInitializeWebAssemblyDisassembler()
	case X86:
		binding.LLVMInitializeX86Disassembler()
	case XCore:
		binding.LLVMInitializeXCoreDisassembler()
	default:
		return ErrUnknownArch
	}
	return nil
}

func InitializeNativeDisassembler() error {
	return binding.LLVMInitializeNativeDisassembler()
}

type dataLayout binding.LLVMTargetDataRef

func (d dataLayout) binding() binding.LLVMTargetDataRef {
	return binding.LLVMTargetDataRef(d)
}

func (d dataLayout) Free() {
	binding.LLVMDisposeTargetData(d.binding())
}

func (d dataLayout) String() string {
	return binding.LLVMCopyStringRepOfTargetData(d.binding())
}

func (d dataLayout) PointerSize() uint {
	return uint(binding.LLVMPointerSize(d.binding()))
}

func (d dataLayout) GetSizeOfType(t Type) uint {
	return uint(binding.LLVMSizeOfTypeInBits(d.binding(), t.binding()))
}

func (d dataLayout) GetStoreSizeOfType(t Type) uint {
	return uint(binding.LLVMStoreSizeOfType(d.binding(), t.binding()))
}

func (d dataLayout) GetABISizeOfType(t Type) uint {
	return uint(binding.LLVMABISizeOfType(d.binding(), t.binding()))
}

func (d dataLayout) GetABIAlignOfType(t Type) uint {
	return uint(binding.LLVMABIAlignmentOfType(d.binding(), t.binding()))
}

func (d dataLayout) GetCallFrameAlignOfType(t Type) uint {
	return uint(binding.LLVMCallFrameAlignmentOfType(d.binding(), t.binding()))
}

func (d dataLayout) GetPrefAlignOfType(t Type) uint {
	return uint(binding.LLVMPreferredAlignmentOfType(d.binding(), t.binding()))
}

func (d dataLayout) GetPrefAlignOfGlobal(g GlobalValue) uint {
	return uint(binding.LLVMPreferredAlignmentOfGlobal(d.binding(), g.binding()))
}

func (d dataLayout) GetOffsetOfElem(st StructType, i uint) uint {
	return uint(binding.LLVMOffsetOfElement(d.binding(), st.binding(), uint32(i)))
}

type Target struct {
	dataLayout
	machine binding.LLVMTargetMachineRef
}

func NativeTarget() (*Target, error) {
	return NewTargetFromTriple(binding.LLVMGetDefaultTargetTriple(), CPUName, CPUFeatures)
}

func NewTargetFromTriple(triple string, cpu, feature string) (*Target, error) {
	t, err := binding.LLVMGetTargetFromTriple(triple)
	if err != nil {
		return nil, err
	}
	machine := binding.LLVMCreateTargetMachine(t, triple, cpu, feature, binding.LLVMCodeGenLevelNone, binding.LLVMRelocDefault, binding.LLVMCodeModelDefault)
	layout := binding.LLVMCreateTargetDataLayout(machine)
	return &Target{
		dataLayout: dataLayout(layout),
		machine:    machine,
	}, nil
}

func (m *Module) SetTarget(t *Target) {
	binding.LLVMSetModuleDataLayout(m.binding(), t.dataLayout.binding())
	binding.LLVMSetTarget(m.binding(), t.Triple())
	m.target = t
}

func (m Module) GetTarget() (*Target, bool) {
	if m.target == nil {
		return nil, false
	}
	return m.target, true
}

func (t Target) getTargetRef() binding.LLVMTargetRef {
	return binding.LLVMGetTargetMachineTarget(t.machine)
}

func (t Target) String() string {
	return t.Triple()
}

func (t Target) Name() string {
	return binding.LLVMGetTargetName(t.getTargetRef())
}

func (t Target) Description() string {
	return binding.LLVMGetTargetDescription(t.getTargetRef())
}

func (t Target) HasJIT() bool {
	return binding.LLVMTargetHasJIT(t.getTargetRef())
}

func (t Target) HasTargetMachine() bool {
	return binding.LLVMTargetHasTargetMachine(t.getTargetRef())
}

func (t Target) HasAsmBackend() bool {
	return binding.LLVMTargetHasAsmBackend(t.getTargetRef())
}

func (t Target) Free() {
	t.dataLayout.Free()
	binding.LLVMDisposeTargetMachine(t.machine)
}

func (t Target) Triple() string {
	return binding.LLVMGetTargetMachineTriple(t.machine)
}

func (t Target) CPU() string {
	return binding.LLVMGetTargetMachineCPU(t.machine)
}

func (t Target) Feature() string {
	return binding.LLVMGetTargetMachineFeatureString(t.machine)
}

type CodeOptLevel binding.LLVMCodeGenOptLevel

const (
	CodeOptLevelNone       = CodeOptLevel(binding.LLVMCodeGenLevelNone)
	CodeOptLevelLess       = CodeOptLevel(binding.LLVMCodeGenLevelLess)
	CodeOptLevelDefault    = CodeOptLevel(binding.LLVMCodeGenLevelDefault)
	CodeOptLevelAggressive = CodeOptLevel(binding.LLVMCodeGenLevelAggressive)
)

type RelocMode binding.LLVMRelocMode

const (
	RelocModeDefault      = RelocMode(binding.LLVMRelocDefault)
	RelocModeStatic       = RelocMode(binding.LLVMRelocStatic)
	RelocModePIC          = RelocMode(binding.LLVMRelocPIC)
	RelocModeDynamicNoPic = RelocMode(binding.LLVMRelocDynamicNoPic)
	RelocModeROPI         = RelocMode(binding.LLVMRelocROPI)
	RelocModeRWPI         = RelocMode(binding.LLVMRelocRWPI)
	RelocModeROPI_RWPI    = RelocMode(binding.LLVMRelocROPI_RWPI)
)

type CodeModel binding.LLVMRelocMode

const (
	CodeModelDefault    = CodeModel(binding.LLVMCodeModelDefault)
	CodeModelJITDefault = CodeModel(binding.LLVMCodeModelJITDefault)
	CodeModelTiny       = CodeModel(binding.LLVMCodeModelTiny)
	CodeModelSmall      = CodeModel(binding.LLVMCodeModelSmall)
	CodeModelKernel     = CodeModel(binding.LLVMCodeModelKernel)
	CodeModelMedium     = CodeModel(binding.LLVMCodeModelMedium)
	CodeModelLarge      = CodeModel(binding.LLVMCodeModelLarge)
)

func (t Target) WriteASMToFile(m Module, file string, opt CodeOptLevel, reloc RelocMode, code CodeModel) error {
	machine := binding.LLVMCreateTargetMachine(t.getTargetRef(), t.Triple(), t.CPU(), t.Feature(), binding.LLVMCodeGenOptLevel(opt), binding.LLVMRelocMode(reloc), binding.LLVMCodeModel(code))
	defer binding.LLVMDisposeTargetMachine(machine)
	return binding.LLVMTargetMachineEmitToFile(machine, m.binding(), file, binding.LLVMAssemblyFile)
}

func (t Target) WriteOBJToFile(m Module, file string, opt CodeOptLevel, reloc RelocMode, code CodeModel) error {
	machine := binding.LLVMCreateTargetMachine(t.getTargetRef(), t.Triple(), t.CPU(), t.Feature(), binding.LLVMCodeGenOptLevel(opt), binding.LLVMRelocMode(reloc), binding.LLVMCodeModel(code))
	defer binding.LLVMDisposeTargetMachine(machine)
	return binding.LLVMTargetMachineEmitToFile(machine, m.binding(), file, binding.LLVMObjectFile)
}

func (t Target) IsLinux() bool {
	return strings.Contains(t.String(), "linux")
}

func (t Target) IsDarwin() bool {
	return strings.Contains(t.String(), "darwin")
}

func (t Target) IsWindows() bool {
	return strings.Contains(t.String(), "windows")
}
