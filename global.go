package llvm

import (
	"github.com/samber/lo"

	"github.com/kkkunny/go-llvm/internal/binding"
)

type Global interface {
	binding() binding.LLVMValueRef
	global()
}

type Function binding.LLVMValueRef

func (m Module) NewFunction(name string, t FunctionType) Function {
	return Function(binding.LLVMAddFunction(m.binding(), name, t.binding()))
}

func (m Module) GetFunction(name string) (Function, bool) {
	f := Function(binding.LLVMGetNamedFunction(m.binding(), name))
	if f.binding().IsNil() {
		return Function{}, false
	}
	return f, true
}

func (Function) constant() {}

func (Function) global() {}

func (c Function) String() string {
	return binding.LLVMPrintValueToString(c.binding())
}

func (c Function) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(c)
}

func (c Function) Type() Type {
	return lookupType(binding.LLVMTypeOf(c.binding()))
}

func (c Function) Name() string {
	return binding.LLVMGetValueName(c.binding())
}

func (c Function) SetName(name string) {
	binding.LLVMSetValueName(c.binding(), name)
}

func (f Function) FirstBlock() Block {
	return Block(binding.LLVMGetFirstBasicBlock(f.binding()))
}

func (f Function) LastBlock() Block {
	return Block(binding.LLVMGetLastBasicBlock(f.binding()))
}

func (f Function) Blocks() []Block {
	return lo.Map(binding.LLVMGetBasicBlocks(f.binding()), func(item binding.LLVMBasicBlockRef, index int) Block {
		return Block(item)
	})
}

func (f Function) EntryBlock() Block {
	return Block(binding.LLVMGetEntryBasicBlock(f.binding()))
}

func (f Function) CountParams() uint {
	return uint(binding.LLVMCountParams(f.binding()))
}

func (f Function) Params() []Param {
	return lo.Map(binding.LLVMGetParams(f.binding()), func(item binding.LLVMValueRef, index int) Param {
		return Param(item)
	})
}

func (f Function) GetParam(i uint) Param {
	return Param(binding.LLVMGetParam(f.binding(), uint32(i)))
}

func (f Function) FirstParam() Param {
	return Param(binding.LLVMGetFirstParam(f.binding()))
}

func (f Function) LastParam() Param {
	return Param(binding.LLVMGetLastParam(f.binding()))
}

func (f Function) Verify() bool {
	return !binding.LLVMVerifyFunction(f.binding(), binding.LLVMReturnStatusAction)
}

func (f Function) FunctionType() FunctionType {
	return lookupType(binding.LLVMGetFunctionType(f.binding())).(FunctionType)
}

func (f Function) VerifyWithCFG(only bool) {
	if !only {
		binding.LLVMViewFunctionCFG(f.binding())
	} else {
		binding.LLVMViewFunctionCFGOnly(f.binding())
	}
}

func (g Function) Linkage() Linkage {
	return Linkage(binding.LLVMGetLinkage(g.binding()))
}

func (g Function) SetLinkage(linkage Linkage) {
	binding.LLVMSetLinkage(g.binding(), binding.LLVMLinkage(linkage))
}

type GlobalValue binding.LLVMValueRef

func (m Module) NewGlobal(name string, t Type) GlobalValue {
	return GlobalValue(binding.LLVMAddGlobal(m.binding(), t.binding(), name))
}

func (m Module) DelGlobal(g GlobalValue) {
	binding.LLVMDeleteGlobal(g.binding())
}

func (m Module) GetGlobal(name string) (GlobalValue, bool) {
	v := GlobalValue(binding.LLVMGetNamedGlobal(m.binding(), name))
	if v.binding().IsNil() {
		return GlobalValue{}, false
	}
	return v, true
}

func (v GlobalValue) constant() {}

func (v GlobalValue) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

func (v GlobalValue) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(v)
}

func (v GlobalValue) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (g GlobalValue) ValueType() Type {
	return lookupType(binding.LLVMGlobalGetValueType(g.binding()))
}

func (g GlobalValue) SetAlign(align uint) {
	binding.LLVMSetAlignment(g.binding(), uint32(align))
}

func (g GlobalValue) IsDeclaration() bool {
	return binding.LLVMIsDeclaration(g.binding())
}

type Linkage binding.LLVMLinkage

const (
	ExternalLinkage            = Linkage(binding.LLVMExternalLinkage)
	AvailableExternallyLinkage = Linkage(binding.LLVMAvailableExternallyLinkage)
	LinkOnceAnyLinkage         = Linkage(binding.LLVMLinkOnceAnyLinkage)
	LinkOnceODRLinkage         = Linkage(binding.LLVMLinkOnceODRLinkage)
	LinkOnceODRAutoHideLinkage = Linkage(binding.LLVMLinkOnceODRAutoHideLinkage)
	WeakAnyLinkage             = Linkage(binding.LLVMWeakAnyLinkage)
	WeakODRLinkage             = Linkage(binding.LLVMWeakODRLinkage)
	AppendingLinkage           = Linkage(binding.LLVMAppendingLinkage)
	InternalLinkage            = Linkage(binding.LLVMInternalLinkage)
	PrivateLinkage             = Linkage(binding.LLVMPrivateLinkage)
	DLLImportLinkage           = Linkage(binding.LLVMDLLImportLinkage)
	DLLExportLinkage           = Linkage(binding.LLVMDLLExportLinkage)
	ExternalWeakLinkage        = Linkage(binding.LLVMExternalWeakLinkage)
	GhostLinkage               = Linkage(binding.LLVMGhostLinkage)
	CommonLinkage              = Linkage(binding.LLVMCommonLinkage)
	LinkerPrivateLinkage       = Linkage(binding.LLVMLinkerPrivateLinkage)
	LinkerPrivateWeakLinkage   = Linkage(binding.LLVMLinkerPrivateWeakLinkage)
)

func (g GlobalValue) Linkage() Linkage {
	return Linkage(binding.LLVMGetLinkage(g.binding()))
}

func (g GlobalValue) SetLinkage(linkage Linkage) {
	binding.LLVMSetLinkage(g.binding(), binding.LLVMLinkage(linkage))
}

type Visibility binding.LLVMVisibility

const (
	DefaultVisibility   = Visibility(binding.LLVMDefaultVisibility)
	HiddenVisibility    = Visibility(binding.LLVMHiddenVisibility)
	ProtectedVisibility = Visibility(binding.LLVMProtectedVisibility)
)

func (g GlobalValue) Visibility() Visibility {
	return Visibility(binding.LLVMGetVisibility(g.binding()))
}

func (GlobalValue) global() {}

func (g GlobalValue) SetVisibility(visibility Visibility) {
	binding.LLVMSetVisibility(g.binding(), binding.LLVMVisibility(visibility))
}

type UnnamedAddr binding.LLVMUnnamedAddr

const (
	NoUnnamedAddr     = UnnamedAddr(binding.LLVMNoUnnamedAddr)
	LocalUnnamedAddr  = UnnamedAddr(binding.LLVMLocalUnnamedAddr)
	GlobalUnnamedAddr = UnnamedAddr(binding.LLVMGlobalUnnamedAddr)
)

func (g GlobalValue) UnnamedAddress() UnnamedAddr {
	return UnnamedAddr(binding.LLVMGetUnnamedAddress(g.binding()))
}

func (g GlobalValue) SetUnnamedAddress(unnamedAddr UnnamedAddr) {
	binding.LLVMSetUnnamedAddress(g.binding(), binding.LLVMUnnamedAddr(unnamedAddr))
}

func (g GlobalValue) GetInitializer() (Constant, bool) {
	init := binding.LLVMGetInitializer(g.binding())
	if init.IsNil() {
		return nil, false
	}
	return lookupConstant(init), true
}

func (g GlobalValue) SetInitializer(v Constant) {
	binding.LLVMSetInitializer(g.binding(), v.binding())
}

func (g GlobalValue) IsThreadLocal() bool {
	return binding.LLVMIsThreadLocal(g.binding())
}

func (g GlobalValue) SetThreadLocal(isThreadLocal bool) {
	binding.LLVMSetThreadLocal(g.binding(), isThreadLocal)
}

func (g GlobalValue) IsGlobalConstant() bool {
	return binding.LLVMIsGlobalConstant(g.binding())
}

func (g GlobalValue) SetGlobalConstant(isConstant bool) {
	binding.LLVMSetGlobalConstant(g.binding(), isConstant)
}

type ThreadLocalMode binding.LLVMThreadLocalMode

const (
	NotThreadLocal         = ThreadLocalMode(binding.LLVMNotThreadLocal)
	GeneralDynamicTLSModel = ThreadLocalMode(binding.LLVMGeneralDynamicTLSModel)
	LocalDynamicTLSModel   = ThreadLocalMode(binding.LLVMLocalDynamicTLSModel)
	InitialExecTLSModel    = ThreadLocalMode(binding.LLVMInitialExecTLSModel)
	LocalExecTLSModel      = ThreadLocalMode(binding.LLVMLocalExecTLSModel)
)

func (g GlobalValue) ThreadLocalMode() ThreadLocalMode {
	return ThreadLocalMode(binding.LLVMGetThreadLocalMode(g.binding()))
}

func (g GlobalValue) SetThreadLocalMode(mode ThreadLocalMode) {
	binding.LLVMSetThreadLocalMode(g.binding(), binding.LLVMThreadLocalMode(mode))
}

func (g GlobalValue) IsExternallyInitialized() bool {
	return binding.LLVMIsExternallyInitialized(g.binding())
}

func (g GlobalValue) SetExternallyInitialized(isConstant bool) {
	binding.LLVMSetExternallyInitialized(g.binding(), isConstant)
}
