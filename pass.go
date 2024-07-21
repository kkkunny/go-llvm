package llvm

import (
	"strings"

	"github.com/kkkunny/go-llvm/internal/binding"
)

type PassOption binding.LLVMPassBuilderOptionsRef

func NewPassOption() PassOption {
	return PassOption(binding.LLVMCreatePassBuilderOptions())
}

func (o PassOption) binding() binding.LLVMPassBuilderOptionsRef {
	return binding.LLVMPassBuilderOptionsRef(o)
}

func (o PassOption) SetVerifyEach(v bool) {
	binding.LLVMPassBuilderOptionsSetVerifyEach(o.binding(), v)
}

func (o PassOption) SetDebugLogging(v bool) {
	binding.LLVMPassBuilderOptionsSetDebugLogging(o.binding(), v)
}

func (o PassOption) SetLoopInterleaving(v bool) {
	binding.LLVMPassBuilderOptionsSetLoopInterleaving(o.binding(), v)
}

func (o PassOption) SetLoopVectorization(v bool) {
	binding.LLVMPassBuilderOptionsSetLoopVectorization(o.binding(), v)
}

func (o PassOption) SetSLPVectorization(v bool) {
	binding.LLVMPassBuilderOptionsSetSLPVectorization(o.binding(), v)
}

func (o PassOption) SetLoopUnrolling(v bool) {
	binding.LLVMPassBuilderOptionsSetLoopUnrolling(o.binding(), v)
}

func (o PassOption) SetForgetAllSCEVInLoopUnroll(v bool) {
	binding.LLVMPassBuilderOptionsSetForgetAllSCEVInLoopUnroll(o.binding(), v)
}

func (o PassOption) SetLicmMssaOptCap(v uint32) {
	binding.LLVMPassBuilderOptionsSetLicmMssaOptCap(o.binding(), v)
}

func (o PassOption) SetLicmMssaNoAccForPromotionCap(v uint32) {
	binding.LLVMPassBuilderOptionsSetLicmMssaNoAccForPromotionCap(o.binding(), v)
}

func (o PassOption) SetCallGraphProfile(v bool) {
	binding.LLVMPassBuilderOptionsSetCallGraphProfile(o.binding(), v)
}

func (o PassOption) SetMergeFunctions(v bool) {
	binding.LLVMPassBuilderOptionsSetMergeFunctions(o.binding(), v)
}

func (o PassOption) Free() {
	binding.LLVMDisposePassBuilderOptions(o.binding())
}

// RunPasses https://llvm.org/docs/Passes.html
func (m Module) RunPasses(option PassOption, pass ...string) error {
	target, err := m.GetTarget()
	if err != nil {
		return err
	}
	machine := newMachine(target.targetInfo, target.triple, "generic", "", CodeOptLevelNone, RelocModeDefault, CodeModelDefault)
	defer machine.Free()
	return binding.LLVMRunPasses(m.binding(), strings.Join(pass, ","), machine.binding(), option.binding())
}
