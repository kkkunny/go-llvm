package llvm

import (
	"iter"

	"github.com/kkkunny/go-llvm/internal/binding"
)

// FloatKind 浮点类型种类
type FloatKind binding.LLVMTypeKind

const (
	FloatHalf     = FloatKind(binding.LLVMHalfTypeKind)
	FloatBFloat   = FloatKind(binding.LLVMBFloatTypeKind)
	FloatSingle   = FloatKind(binding.LLVMFloatTypeKind)
	FloatDouble   = FloatKind(binding.LLVMDoubleTypeKind)
	FloatX86FP80  = FloatKind(binding.LLVMX86_FP80TypeKind)
	FloatFP128    = FloatKind(binding.LLVMFP128TypeKind)
	FloatPPCFP128 = FloatKind(binding.LLVMPPC_FP128TypeKind)
)

// AnyType 类型句柄的非泛型视图
type AnyType interface {
	Ref() binding.LLVMTypeRef
	DynType() Type[DynT]
	Context() *Context
	String() string
	IsSized() bool
	IsNil() bool
	Equal(other AnyType) bool
}

// TypeRef 种类安全的类型引用：Type[T] 与全部类型角色（IntType/FnType/...）均实现。
// 用于泛型方法/函数参数，使 ctx.ConstNull(ctx.Int(32)) 这类调用可直接推断 T。
type TypeRef[T Kind] interface {
	AsType() Type[T]
}

// Type 种类级泛型类型句柄；类别专属操作见各类型角色（如 IntType.Bits）
type Type[T Kind] struct {
	ref binding.LLVMTypeRef
	ctx *Context
}

// Ref 返回底层句柄（供 llvm/* 子包桥接使用）
func (t Type[T]) Ref() binding.LLVMTypeRef { return t.ref }

// DynType 擦除类型参数
func (t Type[T]) DynType() Type[DynT] { return Type[DynT]{ref: t.ref, ctx: t.ctx} }

// AsType 返回自身（实现 TypeRef[T]；类型角色经内嵌继承）
func (t Type[T]) AsType() Type[T] { return t }

// Context 返回所属上下文
func (t Type[T]) Context() *Context { return t.ctx }

// IsNil 是否为空句柄
func (t Type[T]) IsNil() bool { return t.ref.IsNil() }

// Check 类型操作前置校验：句柄非零值且 Context 存活。
// 供 llvm/* 子包的类型角色方法统一调用（角色经内嵌 Type[T] 自动继承）。
func (t Type[T]) Check(op string) {
	if t.ctx == nil || t.ref.IsNil() {
		errPanic(ErrInvalidArg, op, "nil type handle")
	}
	t.ctx.CheckAlive(op)
}

// String 类型的 IR 文本表示
func (t Type[T]) String() string {
	if t.ref.IsNil() {
		return "<nil>"
	}
	t.Check("llvm.Type.String")
	return binding.LLVMPrintTypeToString(t.ref)
}

// IsSized 类型是否具有确定大小
func (t Type[T]) IsSized() bool {
	t.Check("llvm.Type.IsSized")
	return binding.LLVMTypeIsSized(t.ref)
}

// Equal 与另一类型是否同一 LLVM 类型
func (t Type[T]) Equal(other AnyType) bool {
	if other == nil {
		return false
	}
	t.Check("llvm.Type.Equal")
	return t.ref.Equal(other.Ref())
}

// As 运行时校验种类后转换类型参数；目标是 DynT 时始终成功
func (t Type[T]) As[U Kind]() (Type[U], error) {
	t.Check("llvm.Type.As")
	if !kindMatches[U](t.ref) {
		return Type[U]{}, &Error{
			Reason: ErrTypeMismatch,
			Op:     "llvm.Type.As",
			Msg:    "type kind mismatch: have " + kindName(kindOfType(t.ref)) + ", want " + kindName(kindOf[U]()),
		}
	}
	return Type[U]{ref: t.ref, ctx: t.ctx}, nil
}

// MustAs As 的 panic 版本（程序员错误）
func (t Type[T]) MustAs[U Kind]() Type[U] {
	res, err := t.As[U]()
	if err != nil {
		panic(err)
	}
	return res
}

// NewType 由底层句柄构建类型（供 llvm/* 子包桥接使用）
func NewType[T Kind](ctx *Context, ref binding.LLVMTypeRef) Type[T] {
	return Type[T]{ref: ref, ctx: ctx}
}

// TypeOfRef 由底层句柄构建擦除类型（动态来源）
func TypeOfRef(ctx *Context, ref binding.LLVMTypeRef) Type[DynT] {
	return Type[DynT]{ref: ref, ctx: ctx}
}

// ===== 类型角色转换（种类不符 panic，程序员错误） =====

func asRole[T Kind](op string, t AnyType) Type[T] {
	if t == nil {
		errPanic(ErrInvalidArg, op, "nil type")
	}
	return Must(t.DynType().As[T]())
}

// AsVoidType 提取 void 类型角色
func AsVoidType(t AnyType) VoidType { return VoidType{asRole[VoidT]("llvm.AsVoidType", t)} }

// AsIntType 提取整数类型角色
func AsIntType(t AnyType) IntType { return IntType{asRole[IntT]("llvm.AsIntType", t)} }

// AsFloatType 提取浮点类型角色
func AsFloatType(t AnyType) FloatType { return FloatType{asRole[FloatT]("llvm.AsFloatType", t)} }

// AsPtrType 提取指针类型角色
func AsPtrType(t AnyType) PtrType { return PtrType{asRole[PtrT]("llvm.AsPtrType", t)} }

// AsStructType 提取结构体类型角色
func AsStructType(t AnyType) StructType { return StructType{asRole[StructT]("llvm.AsStructType", t)} }

// AsArrayType 提取数组类型角色
func AsArrayType(t AnyType) ArrayType { return ArrayType{asRole[ArrayT]("llvm.AsArrayType", t)} }

// AsVecType 提取向量类型角色
func AsVecType(t AnyType) VecType { return VecType{asRole[VecT]("llvm.AsVecType", t)} }

// AsFnType 提取函数类型角色
func AsFnType(t AnyType) FnType { return FnType{asRole[FnT]("llvm.AsFnType", t)} }

// ===== 类型角色 =====

// VoidType void 类型
type VoidType struct{ Type[VoidT] }

// IntType 整数类型
type IntType struct{ Type[IntT] }

// FloatType 浮点类型
type FloatType struct{ Type[FloatT] }

// PtrType 指针类型（LLVM 22 为不透明指针）
type PtrType struct{ Type[PtrT] }

// StructType 结构体类型
type StructType struct{ Type[StructT] }

// ArrayType 数组类型
type ArrayType struct{ Type[ArrayT] }

// VecType 向量类型
type VecType struct{ Type[VecT] }

// FnType 函数类型
type FnType struct{ Type[FnT] }

// Bits 整数位宽
func (t IntType) Bits() uint32 {
	t.Check("llvm.IntType.Bits")
	return binding.LLVMGetIntTypeWidth(t.ref)
}

// Kind 浮点种类
func (t FloatType) Kind() FloatKind {
	t.Check("llvm.FloatType.Kind")
	return FloatKind(binding.LLVMGetTypeKind(t.ref))
}

// Addrspace 指针地址空间
func (t PtrType) Addrspace() uint32 {
	t.Check("llvm.PtrType.Addrspace")
	return binding.LLVMGetPointerAddressSpace(t.ref)
}

// IsOpaque 指针是否不透明
func (t PtrType) IsOpaque() bool {
	t.Check("llvm.PtrType.IsOpaque")
	return binding.LLVMPointerTypeIsOpaque(t.ref)
}

// Name 结构体名称（字面量结构体为空）
func (t StructType) Name() string {
	t.Check("llvm.StructType.Name")
	return binding.LLVMGetStructName(t.ref)
}

// Elems 结构体元素类型；需要零分配遍历时用 AllElems
func (t StructType) Elems() []AnyType {
	t.Check("llvm.StructType.Elems")
	refs := binding.LLVMGetStructElementTypes(t.ref)
	elems := make([]AnyType, len(refs))
	for i, ref := range refs {
		elems[i] = TypeOfRef(t.ctx, ref)
	}
	return elems
}

// AllElems 惰性遍历结构体元素类型（range 友好，无切片分配）
func (t StructType) AllElems() iter.Seq[AnyType] {
	return func(yield func(AnyType) bool) {
		t.Check("llvm.StructType.AllElems")
		n := binding.LLVMCountStructElementTypes(t.ref)
		for i := uint32(0); i < n; i++ {
			if !yield(TypeOfRef(t.ctx, binding.LLVMStructGetTypeAtIndex(t.ref, i))) {
				return
			}
		}
	}
}

// Elem 第 i 个结构体元素类型
func (t StructType) Elem(i uint32) AnyType {
	t.Check("llvm.StructType.Elem")
	return TypeOfRef(t.ctx, binding.LLVMStructGetTypeAtIndex(t.ref, i))
}

// SetBody 设置结构体成员
func (t StructType) SetBody(elems []AnyType, packed bool) {
	t.Check("llvm.StructType.SetBody")
	t.ctx.CheckTypes("llvm.StructType.SetBody", elems)
	binding.LLVMStructSetBody(t.ref, anyTypesToRefs(elems), packed)
}

// IsPacked 是否 packed
func (t StructType) IsPacked() bool {
	t.Check("llvm.StructType.IsPacked")
	return binding.LLVMIsPackedStruct(t.ref)
}

// IsOpaque 是否 opaque（未设置成员）
func (t StructType) IsOpaque() bool {
	t.Check("llvm.StructType.IsOpaque")
	return binding.LLVMIsOpaqueStruct(t.ref)
}

// Elem 数组元素类型
func (t ArrayType) Elem() AnyType {
	t.Check("llvm.ArrayType.Elem")
	return TypeOfRef(t.ctx, binding.LLVMGetElementType(t.ref))
}

// Len 数组长度
func (t ArrayType) Len() uint64 {
	t.Check("llvm.ArrayType.Len")
	return binding.LLVMGetArrayLength2(t.ref)
}

// Elem 向量元素类型
func (t VecType) Elem() AnyType {
	t.Check("llvm.VecType.Elem")
	return TypeOfRef(t.ctx, binding.LLVMGetElementType(t.ref))
}

// Len 向量元素个数
func (t VecType) Len() uint32 {
	t.Check("llvm.VecType.Len")
	return binding.LLVMGetVectorSize(t.ref)
}

// Return 函数返回类型
func (t FnType) Return() AnyType {
	t.Check("llvm.FnType.Return")
	return TypeOfRef(t.ctx, binding.LLVMGetReturnType(t.ref))
}

// Params 函数参数类型
func (t FnType) Params() []AnyType {
	t.Check("llvm.FnType.Params")
	refs := binding.LLVMGetParamTypes(t.ref)
	params := make([]AnyType, len(refs))
	for i, ref := range refs {
		params[i] = TypeOfRef(t.ctx, ref)
	}
	return params
}

// IsVarArg 是否变参
func (t FnType) IsVarArg() bool {
	t.Check("llvm.FnType.IsVarArg")
	return binding.LLVMIsFunctionVarArg(t.ref)
}

// ===== Context 类型构造器 =====

// Void void 类型
func (ctx *Context) Void() VoidType {
	ctx.CheckAlive("llvm.Context.Void")
	return VoidType{Type[VoidT]{ref: binding.LLVMVoidTypeInContext(ctx.ref), ctx: ctx}}
}

// Int 指定位宽的整数类型
func (ctx *Context) Int(bits uint32) IntType {
	ctx.CheckAlive("llvm.Context.Int")
	return IntType{Type[IntT]{ref: binding.LLVMIntTypeInContext(ctx.ref, bits), ctx: ctx}}
}

// Bool 布尔类型（i1）
func (ctx *Context) Bool() IntType { return ctx.Int(1) }

// Float 浮点类型
func (ctx *Context) Float(kind FloatKind) FloatType {
	ctx.CheckAlive("llvm.Context.Float")
	var ref binding.LLVMTypeRef
	switch kind {
	case FloatHalf:
		ref = binding.LLVMHalfTypeInContext(ctx.ref)
	case FloatBFloat:
		ref = binding.LLVMBFloatTypeInContext(ctx.ref)
	case FloatSingle:
		ref = binding.LLVMFloatTypeInContext(ctx.ref)
	case FloatDouble:
		ref = binding.LLVMDoubleTypeInContext(ctx.ref)
	case FloatX86FP80:
		ref = binding.LLVMX86FP80TypeInContext(ctx.ref)
	case FloatFP128:
		ref = binding.LLVMFP128TypeInContext(ctx.ref)
	case FloatPPCFP128:
		ref = binding.LLVMPPCFP128TypeInContext(ctx.ref)
	default:
		errPanic(ErrInvalidArg, "llvm.Context.Float", "unknown float kind %d", kind)
	}
	return FloatType{Type[FloatT]{ref: ref, ctx: ctx}}
}

// Ptr 指定地址空间的不透明指针类型
func (ctx *Context) Ptr(addrspace uint32) PtrType {
	ctx.CheckAlive("llvm.Context.Ptr")
	return PtrType{Type[PtrT]{ref: binding.LLVMPointerTypeInContext(ctx.ref, addrspace), ctx: ctx}}
}

// Struct 字面量结构体类型
func (ctx *Context) Struct(elems []AnyType, packed bool) StructType {
	ctx.CheckAlive("llvm.Context.Struct")
	ctx.CheckTypes("llvm.Context.Struct", elems)
	return StructType{Type[StructT]{ref: binding.LLVMStructTypeInContext(ctx.ref, anyTypesToRefs(elems), packed), ctx: ctx}}
}

// NamedStruct 命名结构体类型（初始为 opaque，需 SetBody 设置成员）
func (ctx *Context) NamedStruct(name string) StructType {
	ctx.CheckAlive("llvm.Context.NamedStruct")
	return StructType{Type[StructT]{ref: binding.LLVMStructCreateNamed(ctx.ref, name), ctx: ctx}}
}

// Array 数组类型
func (ctx *Context) Array(elem AnyType, n uint64) ArrayType {
	ctx.CheckAlive("llvm.Context.Array")
	ctx.CheckType("llvm.Context.Array", elem)
	return ArrayType{Type[ArrayT]{ref: binding.LLVMArrayType2(elem.Ref(), n), ctx: ctx}}
}

// Vec 向量类型
func (ctx *Context) Vec(elem AnyType, n uint32) VecType {
	ctx.CheckAlive("llvm.Context.Vec")
	ctx.CheckType("llvm.Context.Vec", elem)
	return VecType{Type[VecT]{ref: binding.LLVMVectorType(elem.Ref(), n), ctx: ctx}}
}

// Fn 函数类型
func (ctx *Context) Fn(ret AnyType, params []AnyType, vararg bool) FnType {
	ctx.CheckAlive("llvm.Context.Fn")
	ctx.CheckType("llvm.Context.Fn", ret)
	ctx.CheckTypes("llvm.Context.Fn", params)
	return FnType{Type[FnT]{ref: binding.LLVMFunctionType(ret.Ref(), anyTypesToRefs(params), vararg), ctx: ctx}}
}

// ===== 校验入口（供 llvm/* 子包统一调用） =====

// CheckType 校验类型归属同一 Context 且非空
func (ctx *Context) CheckType(op string, t AnyType) {
	ctx.CheckAlive(op)
	if t == nil || t.IsNil() {
		errPanic(ErrInvalidArg, op, "nil type")
	}
	if t.Context() != ctx {
		errPanic(ErrCrossContext, op, "type belongs to another context")
	}
}

// CheckTypes 校验类型列表归属同一 Context
func (ctx *Context) CheckTypes(op string, elems []AnyType) {
	for _, elem := range elems {
		ctx.CheckType(op, elem)
	}
}

func anyTypesToRefs(types []AnyType) []binding.LLVMTypeRef {
	refs := make([]binding.LLVMTypeRef, len(types))
	for i, t := range types {
		refs[i] = t.Ref()
	}
	return refs
}

// kindOfType 底层类型句柄的种类标记
func kindOfType(ref binding.LLVMTypeRef) Kind {
	switch binding.LLVMGetTypeKind(ref) {
	case binding.LLVMVoidTypeKind:
		return VoidT{}
	case binding.LLVMIntegerTypeKind:
		return IntT{}
	case binding.LLVMHalfTypeKind, binding.LLVMBFloatTypeKind, binding.LLVMFloatTypeKind,
		binding.LLVMDoubleTypeKind, binding.LLVMX86_FP80TypeKind, binding.LLVMFP128TypeKind,
		binding.LLVMPPC_FP128TypeKind:
		return FloatT{}
	case binding.LLVMPointerTypeKind:
		return PtrT{}
	case binding.LLVMStructTypeKind:
		return StructT{}
	case binding.LLVMArrayTypeKind:
		return ArrayT{}
	case binding.LLVMVectorTypeKind, binding.LLVMScalableVectorTypeKind:
		return VecT{}
	case binding.LLVMFunctionTypeKind:
		return FnT{}
	case binding.LLVMLabelTypeKind:
		return LabelT{}
	case binding.LLVMMetadataTypeKind:
		return MetaT{}
	case binding.LLVMTokenTypeKind:
		return TokenT{}
	default:
		return DynT{}
	}
}

// kindMatches 类型句柄是否匹配种类参数 U（U 为 DynT 时始终匹配）
func kindMatches[U Kind](ref binding.LLVMTypeRef) bool {
	if _, ok := any(kindOf[U]()).(DynT); ok {
		return true
	}
	return sameKind(kindOfType(ref), kindOf[U]())
}

func sameKind(a, b Kind) bool {
	return any(a) == any(b)
}

func kindName(k Kind) string {
	switch k.(type) {
	case VoidT:
		return "void"
	case IntT:
		return "int"
	case FloatT:
		return "float"
	case PtrT:
		return "pointer"
	case StructT:
		return "struct"
	case ArrayT:
		return "array"
	case VecT:
		return "vector"
	case FnT:
		return "function"
	case LabelT:
		return "label"
	case MetaT:
		return "metadata"
	case TokenT:
		return "token"
	default:
		return "dyn"
	}
}
