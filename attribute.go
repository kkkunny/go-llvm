package llvm

import (
	"fmt"

	"github.com/kkkunny/go-llvm/internal/binding"
)

// AttrIndex 属性位置：返回值、函数本体或参数（参数从 0 起）
type AttrIndex int32

const (
	AttrReturn   AttrIndex = AttrIndex(binding.LLVMAttributeReturnIndex)
	AttrFunction AttrIndex = AttrIndex(binding.LLVMAttributeFunctionIndex)
)

// AttrParam 第 i 个参数（i 从 0 起）的属性位置
func AttrParam(i uint) AttrIndex { return AttrIndex(i + 1) }

// AttributeKind 枚举属性种类 ID（由 LLVM 按属性名全局唯一分配）
type AttributeKind uint32

// AttributeKindForName 按名称查询枚举属性 ID；LLVM 不认识该名称时返回 0
func AttributeKindForName(name string) AttributeKind {
	return AttributeKind(binding.LLVMGetEnumAttributeKindForName(name))
}

// MustAttributeKind 按名称查询枚举属性 ID；未知名称 panic ErrNotFound
func MustAttributeKind(name string) AttributeKind {
	k := AttributeKindForName(name)
	if k == 0 {
		errPanic(ErrNotFound, "llvm.MustAttributeKind", "unknown attribute %q", name)
	}
	return k
}

// attrKindNames 常用枚举属性名称表（同时用于反查打印）
var attrKindNames = []string{
	"alwaysinline", "noinline", "noreturn", "nounwind", "cold", "hot",
	"optnone", "optsize", "minsize", "readnone", "readonly", "writeonly",
	"memory", "willreturn", "mustprogress", "convergent", "nocallback",
	"nofree", "nosync", "noduplicate", "returns_twice", "uwtable", "naked",
	"builtin", "noalias", "nonnull", "noundef", "inreg", "signext", "zeroext",
	"align", "dereferenceable", "dereferenceable_or_null", "allocsize",
	"swifterror", "swiftself", "nofpclass", "captures", "noprofile",
	"skipprofile", "speculatable", "strictfp", "sanitize_address",
	"sanitize_thread", "sanitize_memory", "ssp", "sspstrong", "sspreq",
	"byval", "sret", "byref", "preallocated",
}

// 常用枚举属性 kind
var (
	AttrAlwaysInline    = MustAttributeKind("alwaysinline")
	AttrNoInline        = MustAttributeKind("noinline")
	AttrNoReturn        = MustAttributeKind("noreturn")
	AttrNoUnwind        = MustAttributeKind("nounwind")
	AttrCold            = MustAttributeKind("cold")
	AttrHot             = MustAttributeKind("hot")
	AttrOptNone         = MustAttributeKind("optnone")
	AttrOptSize         = MustAttributeKind("optsize")
	AttrMinSize         = MustAttributeKind("minsize")
	AttrReadNone        = MustAttributeKind("readnone")
	AttrReadOnly        = MustAttributeKind("readonly")
	AttrWriteOnly       = MustAttributeKind("writeonly")
	AttrMemory          = MustAttributeKind("memory")
	AttrWillReturn      = MustAttributeKind("willreturn")
	AttrMustProgress    = MustAttributeKind("mustprogress")
	AttrConvergent      = MustAttributeKind("convergent")
	AttrNoCallback      = MustAttributeKind("nocallback")
	AttrNoFree          = MustAttributeKind("nofree")
	AttrNoSync          = MustAttributeKind("nosync")
	AttrNoDuplicate     = MustAttributeKind("noduplicate")
	AttrReturnsTwice    = MustAttributeKind("returns_twice")
	AttrUWTable         = MustAttributeKind("uwtable")
	AttrNaked           = MustAttributeKind("naked")
	AttrBuiltin         = MustAttributeKind("builtin")
	AttrNoAlias         = MustAttributeKind("noalias")
	AttrNonNull         = MustAttributeKind("nonnull")
	AttrNoUndef         = MustAttributeKind("noundef")
	AttrInReg           = MustAttributeKind("inreg")
	AttrSignExt         = MustAttributeKind("signext")
	AttrZeroExt         = MustAttributeKind("zeroext")
	AttrAlign           = MustAttributeKind("align")
	AttrDereferenceable = MustAttributeKind("dereferenceable")
	AttrAllocSize       = MustAttributeKind("allocsize")
	AttrByVal           = MustAttributeKind("byval")
	AttrSRet            = MustAttributeKind("sret")
)

// attributeKindName kind ID 反查名称（仅覆盖 attrKindNames 表；未知返回 #ID）
func attributeKindName(k AttributeKind) string {
	if k == 0 {
		return "#0"
	}
	if name, ok := attrKindNameIndex[k]; ok {
		return name
	}
	return fmt.Sprintf("#%d", uint32(k))
}

// attrKindNameIndex 懒加载的反查表
var attrKindNameIndex = func() map[AttributeKind]string {
	m := make(map[AttributeKind]string, len(attrKindNames))
	for _, name := range attrKindNames {
		if k := AttributeKindForName(name); k != 0 {
			m[k] = name
		}
	}
	return m
}()

// Attribute 属性句柄（由 Context 唯一化，无独立释放）
type Attribute struct {
	ref binding.LLVMAttributeRef
	ctx *Context
}

// AttributeOf 由底层句柄铸造属性（供 llvm/* 子包桥接使用）
func AttributeOf(ctx *Context, ref binding.LLVMAttributeRef) Attribute {
	return Attribute{ref: ref, ctx: ctx}
}

// Ref 返回底层句柄（供 llvm/* 子包桥接使用）
func (a Attribute) Ref() binding.LLVMAttributeRef { return a.ref }

// Context 返回所属上下文
func (a Attribute) Context() *Context { return a.ctx }

// Check 属性可用性前置校验
func (a Attribute) Check(op string) {
	if a.ref.IsNil() || a.ctx == nil {
		errPanic(ErrInvalidArg, op, "nil attribute")
	}
	a.ctx.CheckAlive(op)
}

// IsEnum 是否枚举属性
func (a Attribute) IsEnum() bool {
	a.Check("llvm.Attribute.IsEnum")
	return binding.LLVMIsEnumAttribute(a.ref)
}

// IsString 是否字符串属性
func (a Attribute) IsString() bool {
	a.Check("llvm.Attribute.IsString")
	return binding.LLVMIsStringAttribute(a.ref)
}

// IsType 是否类型属性
func (a Attribute) IsType() bool {
	a.Check("llvm.Attribute.IsType")
	return binding.LLVMIsTypeAttribute(a.ref)
}

// EnumKind 枚举属性种类（非枚举属性 panic）
func (a Attribute) EnumKind() AttributeKind {
	const op = "llvm.Attribute.EnumKind"
	a.Check(op)
	if !binding.LLVMIsEnumAttribute(a.ref) {
		errPanic(ErrTypeMismatch, op, "attribute is not an enum attribute")
	}
	return AttributeKind(binding.LLVMGetEnumAttributeKind(a.ref))
}

// EnumValue 枚举属性值（非枚举属性 panic）
func (a Attribute) EnumValue() uint64 {
	const op = "llvm.Attribute.EnumValue"
	a.Check(op)
	if !binding.LLVMIsEnumAttribute(a.ref) {
		errPanic(ErrTypeMismatch, op, "attribute is not an enum attribute")
	}
	return binding.LLVMGetEnumAttributeValue(a.ref)
}

// TypeValue 类型属性的类型（非类型属性 panic）
func (a Attribute) TypeValue() AnyType {
	const op = "llvm.Attribute.TypeValue"
	a.Check(op)
	if !binding.LLVMIsTypeAttribute(a.ref) {
		errPanic(ErrTypeMismatch, op, "attribute is not a type attribute")
	}
	return TypeOfRef(a.ctx, binding.LLVMGetTypeAttributeValue(a.ref))
}

// StringKind 字符串属性的键（非字符串属性 panic）
func (a Attribute) StringKind() string {
	const op = "llvm.Attribute.StringKind"
	a.Check(op)
	if !binding.LLVMIsStringAttribute(a.ref) {
		errPanic(ErrTypeMismatch, op, "attribute is not a string attribute")
	}
	return binding.LLVMGetStringAttributeKind(a.ref)
}

// StringValue 字符串属性的值（非字符串属性 panic）
func (a Attribute) StringValue() string {
	const op = "llvm.Attribute.StringValue"
	a.Check(op)
	if !binding.LLVMIsStringAttribute(a.ref) {
		errPanic(ErrTypeMismatch, op, "attribute is not a string attribute")
	}
	return binding.LLVMGetStringAttributeValue(a.ref)
}

// String 属性的 IR 文本形式
func (a Attribute) String() string {
	a.Check("llvm.Attribute.String")
	switch {
	case binding.LLVMIsEnumAttribute(a.ref):
		k := AttributeKind(binding.LLVMGetEnumAttributeKind(a.ref))
		v := binding.LLVMGetEnumAttributeValue(a.ref)
		if v == 0 {
			return attributeKindName(k)
		}
		return fmt.Sprintf("%s %d", attributeKindName(k), v)
	case binding.LLVMIsTypeAttribute(a.ref):
		k := AttributeKind(binding.LLVMGetEnumAttributeKind(a.ref))
		return fmt.Sprintf("%s(%s)", attributeKindName(k), TypeOfRef(a.ctx, binding.LLVMGetTypeAttributeValue(a.ref)))
	default:
		return fmt.Sprintf("%q=%q", binding.LLVMGetStringAttributeKind(a.ref), binding.LLVMGetStringAttributeValue(a.ref))
	}
}

// EnumAttr 构造枚举属性
func (ctx *Context) EnumAttr(kind AttributeKind, val uint64) Attribute {
	ctx.CheckAlive("llvm.Context.EnumAttr")
	if kind == 0 {
		errPanic(ErrInvalidArg, "llvm.Context.EnumAttr", "unknown attribute kind")
	}
	return Attribute{ref: binding.LLVMCreateEnumAttribute(ctx.ref, uint32(kind), val), ctx: ctx}
}

// EnumAttrName 按属性名构造枚举属性；名称未知 panic ErrNotFound
func (ctx *Context) EnumAttrName(name string, val uint64) Attribute {
	return ctx.EnumAttr(MustAttributeKind(name), val)
}

// TypeAttr 构造类型属性（如 byval/sret）
func (ctx *Context) TypeAttr(kind AttributeKind, t AnyType) Attribute {
	const op = "llvm.Context.TypeAttr"
	ctx.CheckAlive(op)
	if kind == 0 {
		errPanic(ErrInvalidArg, op, "unknown attribute kind")
	}
	ctx.CheckType(op, t)
	return Attribute{ref: binding.LLVMCreateTypeAttribute(ctx.ref, uint32(kind), t.Ref()), ctx: ctx}
}

// StringAttr 构造字符串属性
func (ctx *Context) StringAttr(kind, val string) Attribute {
	ctx.CheckAlive("llvm.Context.StringAttr")
	return Attribute{ref: binding.LLVMCreateStringAttribute(ctx.ref, kind, val), ctx: ctx}
}

// AlignAttr 对齐属性
func (ctx *Context) AlignAttr(n uint32) Attribute { return ctx.EnumAttr(AttrAlign, uint64(n)) }

// DereferenceableAttr 可解引用字节数属性
func (ctx *Context) DereferenceableAttr(n uint64) Attribute {
	return ctx.EnumAttr(AttrDereferenceable, n)
}

// ByValAttr 按值传参属性
func (ctx *Context) ByValAttr(t AnyType) Attribute { return ctx.TypeAttr(AttrByVal, t) }

// SRetAttr 结构体返回属性
func (ctx *Context) SRetAttr(t AnyType) Attribute { return ctx.TypeAttr(AttrSRet, t) }

// CallConv 调用约定
type CallConv binding.LLVMCallConv

const (
	CallConvC             = CallConv(binding.LLVMCCallConv)
	CallConvFast          = CallConv(binding.LLVMFastCallConv)
	CallConvCold          = CallConv(binding.LLVMColdCallConv)
	CallConvGHC           = CallConv(binding.LLVMGHCCallConv)
	CallConvHiPE          = CallConv(binding.LLVMHiPECallConv)
	CallConvAnyReg        = CallConv(binding.LLVMAnyRegCallConv)
	CallConvPreserveMost  = CallConv(binding.LLVMPreserveMostCallConv)
	CallConvPreserveAll   = CallConv(binding.LLVMPreserveAllCallConv)
	CallConvSwift         = CallConv(binding.LLVMSwiftCallConv)
	CallConvCXXFastTLS    = CallConv(binding.LLVMCXXFASTTLSCallConv)
	CallConvX86StdCall    = CallConv(binding.LLVMX86StdcallCallConv)
	CallConvX86FastCall   = CallConv(binding.LLVMX86FastcallCallConv)
	CallConvARMAPCS       = CallConv(binding.LLVMARMAPCSCallConv)
	CallConvARMAAPCS      = CallConv(binding.LLVMARMAAPCSCallConv)
	CallConvARMAAPCSVFP   = CallConv(binding.LLVMARMAAPCSVFPCallConv)
	CallConvMSP430INTR    = CallConv(binding.LLVMMSP430INTRCallConv)
	CallConvX86ThisCall   = CallConv(binding.LLVMX86ThisCallCallConv)
	CallConvPTXKernel     = CallConv(binding.LLVMPTXKernelCallConv)
	CallConvPTXDevice     = CallConv(binding.LLVMPTXDeviceCallConv)
	CallConvSPIRFunc      = CallConv(binding.LLVMSPIRFUNCCallConv)
	CallConvSPIRKernel    = CallConv(binding.LLVMSPIRKERNELCallConv)
	CallConvIntelOCLBI    = CallConv(binding.LLVMIntelOCLBICallConv)
	CallConvX8664SysV     = CallConv(binding.LLVMX8664SysVCallConv)
	CallConvWin64         = CallConv(binding.LLVMWin64CallConv)
	CallConvX86VectorCall = CallConv(binding.LLVMX86VectorCallCallConv)
	CallConvHHVM          = CallConv(binding.LLVMHHVMCallConv)
	CallConvHHVMC         = CallConv(binding.LLVMHHVMCCallConv)
	CallConvX86INTR       = CallConv(binding.LLVMX86INTRCallConv)
	CallConvAVRINTR       = CallConv(binding.LLVMAVRINTRCallConv)
	CallConvAVRSignal     = CallConv(binding.LLVMAVRSIGNALCallConv)
	CallConvAVRBuiltin    = CallConv(binding.LLVMAVRBUILTINCallConv)
	CallConvAMDGPUVS      = CallConv(binding.LLVMAMDGPUVSCallConv)
	CallConvAMDGPUGS      = CallConv(binding.LLVMAMDGPUGSCallConv)
	CallConvAMDGPUPS      = CallConv(binding.LLVMAMDGPUPSCallConv)
	CallConvAMDGPUCS      = CallConv(binding.LLVMAMDGPUCSCallConv)
	CallConvAMDGPUKernel  = CallConv(binding.LLVMAMDGPUKERNELCallConv)
	CallConvX86RegCall    = CallConv(binding.LLVMX86RegCallCallConv)
	CallConvAMDGPUHS      = CallConv(binding.LLVMAMDGPUHSCallConv)
	CallConvMSP430Builtin = CallConv(binding.LLVMMSP430BUILTINCallConv)
	CallConvAMDGPULS      = CallConv(binding.LLVMAMDGPULSCallConv)
	CallConvAMDGPUES      = CallConv(binding.LLVMAMDGPUESCallConv)
)
