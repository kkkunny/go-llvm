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

func (f Function) FirstBlock() (Block, bool) {
	ref := binding.LLVMGetFirstBasicBlock(f.binding())
	if ref.IsNil() {
		return Block{}, false
	}
	return Block(ref), true
}

func (f Function) LastBlock() (Block, bool) {
	ref := binding.LLVMGetLastBasicBlock(f.binding())
	if ref.IsNil() {
		return Block{}, false
	}
	return Block(ref), true
}

func (f Function) Blocks() []Block {
	return lo.Map(binding.LLVMGetBasicBlocks(f.binding()), func(item binding.LLVMBasicBlockRef, index int) Block {
		return Block(item)
	})
}

func (f Function) ForeachBlock(cb func(block Block)) {
	for inst, ok := f.FirstBlock(); ok; inst, ok = inst.Next() {
		cb(inst)
	}
}

func (f Function) EntryBlock() (Block, bool) {
	ref := binding.LLVMGetEntryBasicBlock(f.binding())
	if ref.IsNil() {
		return Block{}, false
	}
	return Block(ref), true
}

func (f Function) OnlyDecl() bool {
	_, ok := f.EntryBlock()
	return !ok
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

type FuncAttribute uint8

const (
	FuncAttributeNoReturn     FuncAttribute = iota // 函数不会返回
	FuncAttributeInlineHint                        // 自动内联
	FuncAttributeAlwaysInline                      // 必须内联
	FuncAttributeNoInline                          // 禁止内联
	// FuncAttributeAllocKind 内存分配类型
	// 1-alloc 2-realloc 1|2-alloc,realloc 4-free 8-uninitialized 16-zeroed 32-aligned
	FuncAttributeAllocKind
)

func (f Function) AddAttribute(attr FuncAttribute, attrValue ...uint) {
	ctx := binding.LLVMGetTypeContext(f.FunctionType().binding())
	var kind uint32
	switch attr {
	case FuncAttributeNoReturn:
		kind = binding.LLVMGetEnumAttributeKindForName("noreturn")
	case FuncAttributeInlineHint:
		kind = binding.LLVMGetEnumAttributeKindForName("inlinehint")
	case FuncAttributeAlwaysInline:
		kind = binding.LLVMGetEnumAttributeKindForName("alwaysinline")
	case FuncAttributeNoInline:
		kind = binding.LLVMGetEnumAttributeKindForName("noinline")
	case FuncAttributeAllocKind:
		kind = binding.LLVMGetEnumAttributeKindForName("allockind")
	default:
		panic("unreachable")
	}
	var val uint64
	if len(attrValue) > 0 {
		val = uint64(attrValue[0])
	}
	bindAttr := binding.LLVMCreateEnumAttribute(ctx, kind, val)
	binding.LLVMAddAttributeAtIndex(f.binding(), binding.LLVMAttributeFunctionIndex, bindAttr)
}

func (f Function) IsDSOLocal() bool {
	return binding.LLVMIsDSOLocal(f.binding())
}

func (f Function) SetDSOLocal(local bool) {
	binding.LLVMSetDSOLocal(f.binding(), local)
}

type GlobalValue binding.LLVMValueRef

func (m Module) NewGlobal(name string, t Type, v Constant) GlobalValue {
	gv := GlobalValue(binding.LLVMAddGlobal(m.binding(), t.binding(), name))
	if v != nil {
		gv.SetInitializer(v)
	}
	return gv
}

func (m Module) NewConstant(name string, v Constant) GlobalValue {
	gv := m.NewGlobal(name, v.Type(), v)
	gv.SetGlobalConstant(true)
	return gv
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

func (g GlobalValue) SetAlign(align uint32) {
	binding.LLVMSetAlignment(g.binding(), uint32(align))
}

func (g GlobalValue) GetAlign() uint32 {
	return binding.LLVMGetAlignment(g.binding())
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

func (g GlobalValue) IsDSOLocal() bool {
	return binding.LLVMIsDSOLocal(g.binding())
}

func (g GlobalValue) SetDSOLocal(local bool) {
	binding.LLVMSetDSOLocal(g.binding(), local)
}
