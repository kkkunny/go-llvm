package llvm

import "C"
import (
	"strings"

	"github.com/kkkunny/go-llvm/internal/binding"
)

func InitializeAllTargetInfos() {
	binding.LLVMInitializeAllTargetInfos()
}

func InitializeAllTargets() {
	binding.LLVMInitializeAllTargets()
}

func InitializeAllTargetMCs() {
	binding.LLVMInitializeAllTargetMCs()
}

func InitializeAllAsmPrinters() {
	binding.LLVMInitializeAllAsmPrinters()
}

func InitializeAllAsmParsers() {
	binding.LLVMInitializeAllAsmParsers()
}

func InitializeAllDisassemblers() {
	binding.LLVMInitializeAllDisassemblers()
}

func InitializeNativeTarget() error {
	return binding.LLVMInitializeNativeTarget()
}

func InitializeNativeAsmParser() error {
	return binding.LLVMInitializeNativeAsmParser()
}

func InitializeNativeAsmPrinter() error {
	return binding.LLVMInitializeNativeAsmPrinter()
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
		return nil, nil
	}
	machine := binding.LLVMCreateTargetMachine(t, triple, cpu, feature, binding.LLVMCodeGenLevelNone, binding.LLVMRelocDefault, binding.LLVMCodeModelDefault)
	layout := binding.LLVMCreateTargetDataLayout(machine)
	return &Target{
		dataLayout: dataLayout(layout),
		machine:    machine,
	}, nil
}

func (m Module) SetTarget(t *Target) {
	binding.LLVMSetModuleDataLayout(m.binding(), t.dataLayout.binding())
	binding.LLVMSetTarget(m.binding(), t.Triple())
}

func (m Module) GetTarget(cpu, feature string) *Target {
	triple := binding.LLVMGetTarget(m.binding())
	target, _ := NewTargetFromTriple(triple, cpu, feature)
	return target
}

func (t Target) getTargetRef() binding.LLVMTargetRef {
	return binding.LLVMGetTargetMachineTarget(t.machine)
}

func (t Target) String() string {
	return t.Name()
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
	return strings.Contains(t.Name(), "linux")
}

func (t Target) IsDarwin() bool {
	return strings.Contains(t.Name(), "darwin")
}

func (t Target) IsWindows() bool {
	return strings.Contains(t.Name(), "windows")
}
