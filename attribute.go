package llvm

import (
	"fmt"

	"github.com/kkkunny/go-llvm/internal/binding"
)

// AttrIndex 属性位置：返回值、函数本体或参数（参数从 0 起）
type AttrIndex int32

// AttrIndex 常量对应 LLVM 属性位置索引（binding.LLVMAttributeReturnIndex 等）。
// 参数位置从 0 起（见 AttrParam），AttrReturn 对应返回值，AttrFunction 对应函数本体。
const (
	AttrReturn   AttrIndex = AttrIndex(binding.LLVMAttributeReturnIndex)   // 返回值属性位置（LLVM 索引 0）
	AttrFunction AttrIndex = AttrIndex(binding.LLVMAttributeFunctionIndex) // 函数本体属性位置（LLVM 索引 -1）
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

// CallConv 取值对应 LLVM 调用约定（binding.LLVMCallConv），决定参数与返回值的传递方式
// 以及寄存器保存规则；目标相关约定（x86/ARM/AVR/AMDGPU 等）只在相应后端上有效。
const (
	CallConvC             = CallConv(binding.LLVMCCallConv)             // C 调用约定（LLVM 默认）
	CallConvFast          = CallConv(binding.LLVMFastCallConv)          // fastcc：快速调用约定，尽量通过寄存器传递参数与返回值
	CallConvCold          = CallConv(binding.LLVMColdCallConv)          // coldcc：冷路径调用约定，尽量减少被调用方对寄存器的破坏
	CallConvGHC           = CallConv(binding.LLVMGHCCallConv)           // GHC 调用约定（Haskell GHC 运行时）
	CallConvHiPE          = CallConv(binding.LLVMHiPECallConv)          // HiPE 调用约定（Erlang HiPE 运行时）
	CallConvAnyReg        = CallConv(binding.LLVMAnyRegCallConv)        // anyregcc：允许值保存在任意寄存器中（供 patchpoint 等动态调用使用）
	CallConvPreserveMost  = CallConv(binding.LLVMPreserveMostCallConv)  // preserve_mostcc：保留绝大多数寄存器，除少量临时寄存器（如 X86-64 的 R11）外仅允许破坏传参与返回值所需寄存器
	CallConvPreserveAll   = CallConv(binding.LLVMPreserveAllCallConv)   // preserve_allcc：除返回值所需外保留所有寄存器
	CallConvSwift         = CallConv(binding.LLVMSwiftCallConv)         // swiftcc：Swift 语言调用约定
	CallConvCXXFastTLS    = CallConv(binding.LLVMCXXFASTTLSCallConv)    // cxx_fast_tlscc：C++ 快速 TLS 访问函数（如 __tls_get_addr）的调用约定
	CallConvX86StdCall    = CallConv(binding.LLVMX86StdcallCallConv)    // x86_stdcallcc：X86 stdcall，被调用方清栈
	CallConvX86FastCall   = CallConv(binding.LLVMX86FastcallCallConv)   // x86_fastcallcc：X86 fastcall，前两个参数经 ECX/EDX 传递
	CallConvARMAPCS       = CallConv(binding.LLVMARMAPCSCallConv)       // arm_apcscc：ARM APCS 调用约定（旧版）
	CallConvARMAAPCS      = CallConv(binding.LLVMARMAAPCSCallConv)      // arm_aapcscc：ARM AAPCS 过程调用标准
	CallConvARMAAPCSVFP   = CallConv(binding.LLVMARMAAPCSVFPCallConv)   // arm_aapcs_vfpcc：ARM AAPCS VFP，浮点参数经 VFP 寄存器传递
	CallConvMSP430INTR    = CallConv(binding.LLVMMSP430INTRCallConv)    // msp430_intrcc：MSP430 中断处理调用约定
	CallConvX86ThisCall   = CallConv(binding.LLVMX86ThisCallCallConv)   // x86_thiscallcc：X86 thiscall，this 指针经 ECX 传递
	CallConvPTXKernel     = CallConv(binding.LLVMPTXKernelCallConv)     // ptx_kernel：PTX 内核入口调用约定
	CallConvPTXDevice     = CallConv(binding.LLVMPTXDeviceCallConv)     // ptx_device：PTX 设备函数调用约定
	CallConvSPIRFunc      = CallConv(binding.LLVMSPIRFUNCCallConv)      // spir_func：SPIR 函数调用约定（OpenCL 设备函数）
	CallConvSPIRKernel    = CallConv(binding.LLVMSPIRKERNELCallConv)    // spir_kernel：SPIR 内核调用约定（OpenCL 内核入口）
	CallConvIntelOCLBI    = CallConv(binding.LLVMIntelOCLBICallConv)    // intel_ocl_bicc：Intel OpenCL 内置函数调用约定
	CallConvX8664SysV     = CallConv(binding.LLVMX8664SysVCallConv)     // x86_64_sysvcc：x86-64 System V ABI 调用约定
	CallConvWin64         = CallConv(binding.LLVMWin64CallConv)         // win64cc：Windows x64 调用约定
	CallConvX86VectorCall = CallConv(binding.LLVMX86VectorCallCallConv) // x86_vectorcallcc：X86 向量调用约定，向量参数经向量寄存器传递
	CallConvHHVM          = CallConv(binding.LLVMHHVMCallConv)          // hhvmcc：HHVM 调用约定（已废弃的占位值）
	CallConvHHVMC         = CallConv(binding.LLVMHHVMCCallConv)         // hhvm_ccc：HHVM C 调用约定（已废弃的占位值）
	CallConvX86INTR       = CallConv(binding.LLVMX86INTRCallConv)       // x86_intrcc：X86 硬件中断处理调用约定
	CallConvAVRINTR       = CallConv(binding.LLVMAVRINTRCallConv)       // avr_intrcc：AVR 中断处理调用约定
	CallConvAVRSignal     = CallConv(binding.LLVMAVRSIGNALCallConv)     // avr_signalcc：AVR 信号处理调用约定
	CallConvAVRBuiltin    = CallConv(binding.LLVMAVRBUILTINCallConv)    // AVR 运行时库内置函数调用约定（保留寄存器的优化约定）
	CallConvAMDGPUVS      = CallConv(binding.LLVMAMDGPUVSCallConv)      // amdgpu_vs：AMDGPU 顶点着色器
	CallConvAMDGPUGS      = CallConv(binding.LLVMAMDGPUGSCallConv)      // amdgpu_gs：AMDGPU 几何着色器
	CallConvAMDGPUPS      = CallConv(binding.LLVMAMDGPUPSCallConv)      // amdgpu_ps：AMDGPU 像素（片元）着色器
	CallConvAMDGPUCS      = CallConv(binding.LLVMAMDGPUCSCallConv)      // amdgpu_cs：AMDGPU 计算着色器
	CallConvAMDGPUKernel  = CallConv(binding.LLVMAMDGPUKERNELCallConv)  // amdgpu_kernel：AMDGPU 内核入口
	CallConvX86RegCall    = CallConv(binding.LLVMX86RegCallCallConv)    // x86_regcallcc：X86 register 调用约定（Intel regcall，尽量用寄存器传参）
	CallConvAMDGPUHS      = CallConv(binding.LLVMAMDGPUHSCallConv)      // amdgpu_hs：AMDGPU 曲面细分控制（hull）着色器
	CallConvMSP430Builtin = CallConv(binding.LLVMMSP430BUILTINCallConv) // MSP430 运行时库内置函数调用约定（使用额外寄存器的优化约定）
	CallConvAMDGPULS      = CallConv(binding.LLVMAMDGPULSCallConv)      // amdgpu_ls：AMDGPU 局部着色器（启用曲面细分时的 AMDPAL 顶点着色器）
	CallConvAMDGPUES      = CallConv(binding.LLVMAMDGPUESCallConv)      // amdgpu_es：AMDGPU 导出着色器（曲面细分时等价于评估着色器，否则为顶点着色器）
)
