package llvm

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/kkkunny/go-llvm/internal/binding"
)

type Type interface {
	fmt.Stringer
	binding() binding.LLVMTypeRef
	Context() Context
	IsSized() bool
}

type AggregateType interface {
	Type
	aggregate()
}

func lookupType(ref binding.LLVMTypeRef) Type {
	switch binding.LLVMGetTypeKind(ref) {
	case binding.LLVMVoidTypeKind:
		return VoidType(ref)
	case binding.LLVMHalfTypeKind, binding.LLVMFloatTypeKind, binding.LLVMDoubleTypeKind, binding.LLVMBFloatTypeKind,
		binding.LLVMX86_FP80TypeKind, binding.LLVMFP128TypeKind, binding.LLVMPPC_FP128TypeKind:
		return FloatType(ref)
	case binding.LLVMIntegerTypeKind:
		return IntegerType(ref)
	case binding.LLVMFunctionTypeKind:
		return FunctionType(ref)
	case binding.LLVMStructTypeKind:
		return StructType(ref)
	case binding.LLVMArrayTypeKind:
		return ArrayType(ref)
	case binding.LLVMPointerTypeKind:
		return PointerType(ref)
	default:
		panic(fmt.Errorf("unknown enum value `%d`", binding.LLVMGetTypeKind(ref)))
	}
}

type VoidType binding.LLVMTypeRef

func (ctx Context) VoidType() VoidType {
	return VoidType(binding.LLVMVoidTypeInContext(ctx.binding()))
}

func (t VoidType) String() string {
	return binding.LLVMPrintTypeToString(t.binding())
}

func (t VoidType) binding() binding.LLVMTypeRef {
	return binding.LLVMTypeRef(t)
}

func (t VoidType) Context() Context {
	return Context(binding.LLVMGetTypeContext(t.binding()))
}

func (t VoidType) IsSized() bool {
	return true
}

type FloatTypeKind binding.LLVMTypeKind

const (
	FloatTypeKindHalf     = FloatTypeKind(binding.LLVMHalfTypeKind)
	FloatTypeKindBFloat   = FloatTypeKind(binding.LLVMBFloatTypeKind)
	FloatTypeKindFloat    = FloatTypeKind(binding.LLVMFloatTypeKind)
	FloatTypeKindDouble   = FloatTypeKind(binding.LLVMDoubleTypeKind)
	FloatTypeKindX86FP80  = FloatTypeKind(binding.LLVMX86_FP80TypeKind)
	FloatTypeKindFP128    = FloatTypeKind(binding.LLVMFP128TypeKind)
	FloatTypeKindPPCFP128 = FloatTypeKind(binding.LLVMPPC_FP128TypeKind)
)

type FloatType binding.LLVMTypeRef

func (ctx Context) FloatType(kind FloatTypeKind) FloatType {
	switch kind {
	case FloatTypeKindHalf:
		return FloatType(binding.LLVMHalfTypeInContext(ctx.binding()))
	case FloatTypeKindBFloat:
		return FloatType(binding.LLVMBFloatTypeInContext(ctx.binding()))
	case FloatTypeKindFloat:
		return FloatType(binding.LLVMFloatTypeInContext(ctx.binding()))
	case FloatTypeKindDouble:
		return FloatType(binding.LLVMDoubleTypeInContext(ctx.binding()))
	case FloatTypeKindX86FP80:
		return FloatType(binding.LLVMX86FP80TypeInContext(ctx.binding()))
	case FloatTypeKindFP128:
		return FloatType(binding.LLVMFP128TypeInContext(ctx.binding()))
	case FloatTypeKindPPCFP128:
		return FloatType(binding.LLVMPPCFP128TypeInContext(ctx.binding()))
	default:
		panic("unreachable")
	}
}

func (t FloatType) String() string {
	return binding.LLVMPrintTypeToString(t.binding())
}

func (t FloatType) binding() binding.LLVMTypeRef {
	return binding.LLVMTypeRef(t)
}

func (t FloatType) Context() Context {
	return Context(binding.LLVMGetTypeContext(t.binding()))
}

func (t FloatType) Kind() FloatTypeKind {
	return FloatTypeKind(binding.LLVMGetTypeKind(binding.LLVMTypeRef(t)))
}

func (t FloatType) IsSized() bool {
	return true
}

type IntegerType binding.LLVMTypeRef

func (ctx Context) IntegerType(bits uint32) IntegerType {
	return IntegerType(binding.LLVMIntTypeInContext(ctx.binding(), bits))
}

func (ctx Context) IntPtrType(t *Target) IntegerType {
	return IntegerType(binding.LLVMIntPtrTypeInContext(ctx.binding(), t.dataLayout.binding()))
}

func (t IntegerType) String() string {
	return binding.LLVMPrintTypeToString(t.binding())
}

func (t IntegerType) binding() binding.LLVMTypeRef {
	return binding.LLVMTypeRef(t)
}

func (t IntegerType) Context() Context {
	return Context(binding.LLVMGetTypeContext(t.binding()))
}

func (t IntegerType) IsSized() bool {
	return true
}

func (t IntegerType) Bits() uint32 {
	return binding.LLVMGetIntTypeWidth(t.binding())
}

type FunctionType binding.LLVMTypeRef

func (ctx Context) FunctionType(isVarArg bool, ret Type, param ...Type) FunctionType {
	var ps []binding.LLVMTypeRef
	if len(param) > 0 {
		ps = lo.Map(param, func(e Type, index int) binding.LLVMTypeRef {
			return e.binding()
		})
	}
	return FunctionType(binding.LLVMFunctionType(ret.binding(), ps, isVarArg))
}

func (t FunctionType) String() string {
	return binding.LLVMPrintTypeToString(t.binding())
}

func (t FunctionType) binding() binding.LLVMTypeRef {
	return binding.LLVMTypeRef(t)
}

func (t FunctionType) Context() Context {
	return Context(binding.LLVMGetTypeContext(t.binding()))
}

func (t FunctionType) IsSized() bool {
	return true
}

func (t FunctionType) IsVarArg() bool {
	return binding.LLVMIsFunctionVarArg(t.binding())
}

func (t FunctionType) ReturnType() Type {
	return lookupType(binding.LLVMGetReturnType(t.binding()))
}

func (t FunctionType) CountParams() uint32 {
	return binding.LLVMCountParamTypes(t.binding())
}

func (t FunctionType) Params() []Type {
	return lo.Map(binding.LLVMGetParamTypes(t.binding()), func(e binding.LLVMTypeRef, index int) Type {
		return lookupType(e)
	})
}

type StructType binding.LLVMTypeRef

func (ctx Context) StructType(packed bool, elems ...Type) StructType {
	var es []binding.LLVMTypeRef
	if len(elems) > 0 {
		es = lo.Map(elems, func(e Type, index int) binding.LLVMTypeRef {
			return e.binding()
		})
	}
	return StructType(binding.LLVMStructTypeInContext(ctx.binding(), es, packed))
}

func (ctx Context) NamedStructType(name string, packed bool, elems ...Type) StructType {
	st := StructType(binding.LLVMStructCreateNamed(ctx.binding(), name))
	st.SetElems(packed, elems...)
	return st
}

func (ctx Context) GetTypeByName(name string) *StructType {
	st := binding.LLVMGetTypeByName(ctx.binding(), name)
	if st.IsNil() {
		return nil
	}
	t := StructType(st)
	return &t
}

func (t StructType) String() string {
	return binding.LLVMPrintTypeToString(t.binding())
}

func (t StructType) binding() binding.LLVMTypeRef {
	return binding.LLVMTypeRef(t)
}

func (t StructType) Context() Context {
	return Context(binding.LLVMGetTypeContext(t.binding()))
}

func (t StructType) IsSized() bool {
	return binding.LLVMTypeIsSized(t.binding())
}

func (t StructType) Name() string {
	return binding.LLVMGetStructName(t.binding())
}

func (t StructType) CountElems() uint32 {
	return binding.LLVMCountStructElementTypes(t.binding())
}

func (t StructType) Elems() []Type {
	return lo.Map(binding.LLVMGetStructElementTypes(t.binding()), func(e binding.LLVMTypeRef, index int) Type {
		return lookupType(e)
	})
}

func (t StructType) GetElem(i uint32) Type {
	return lookupType(binding.LLVMStructGetTypeAtIndex(t.binding(), i))
}

func (t StructType) IsPacked() bool {
	return binding.LLVMIsPackedStruct(t.binding())
}

func (t StructType) IsOpaque() bool {
	return binding.LLVMIsOpaqueStruct(t.binding())
}

func (t StructType) SetElems(packed bool, elems ...Type) {
	var es []binding.LLVMTypeRef
	if len(elems) > 0 {
		es = lo.Map(elems, func(e Type, index int) binding.LLVMTypeRef {
			return e.binding()
		})
	}
	binding.LLVMStructSetBody(t.binding(), es, packed)
}

func (StructType) aggregate() {}

type ArrayType binding.LLVMTypeRef

func (ctx Context) ArrayType(elem Type, cap uint32) ArrayType {
	return ArrayType(binding.LLVMArrayType(elem.binding(), cap))
}

func (t ArrayType) String() string {
	return binding.LLVMPrintTypeToString(t.binding())
}

func (t ArrayType) binding() binding.LLVMTypeRef {
	return binding.LLVMTypeRef(t)
}

func (t ArrayType) Context() Context {
	return Context(binding.LLVMGetTypeContext(t.binding()))
}

func (t ArrayType) IsSized() bool {
	return binding.LLVMTypeIsSized(t.binding())
}

func (t ArrayType) Capacity() uint32 {
	return binding.LLVMGetArrayLength(t.binding())
}

func (t ArrayType) Element() Type {
	return lookupType(binding.LLVMGetElementType(t.binding()))
}

func (ArrayType) aggregate() {}

type PointerType binding.LLVMTypeRef

func (ctx Context) PointerType(elem Type) PointerType {
	// TODO: unknown AddressSpace
	return PointerType(binding.LLVMPointerType(elem.binding(), 0))
}

func (ctx Context) OpaquePointerType() PointerType {
	// TODO: unknown AddressSpace
	return PointerType(binding.LLVMPointerTypeInContext(ctx.binding(), 0))
}

func (t PointerType) String() string {
	return binding.LLVMPrintTypeToString(t.binding())
}

func (t PointerType) binding() binding.LLVMTypeRef {
	return binding.LLVMTypeRef(t)
}

func (t PointerType) Context() Context {
	return Context(binding.LLVMGetTypeContext(t.binding()))
}

func (t PointerType) IsSized() bool {
	return true
}

func (t PointerType) IsOpaque() bool {
	return binding.LLVMPointerTypeIsOpaque(t.binding())
}

func (t PointerType) AddressSpace() uint32 {
	return binding.LLVMGetPointerAddressSpace(t.binding())
}
