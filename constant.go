package llvm

import (
	"fmt"
	"strconv"

	"github.com/samber/lo"

	"github.com/kkkunny/go-llvm/internal/binding"
)

type Constant interface {
	Value
	constant()
}

func lookupConstant(ref binding.LLVMValueRef) Constant {
	if ref.IsNil() {
		return nil
	}

	switch constKind := binding.LLVMGetValueKind(ref); constKind {
	case binding.LLVMGlobalVariableValueKind:
		return GlobalValue(ref)
	case binding.LLVMConstantIntValueKind:
		return ConstInteger(ref)
	case binding.LLVMConstantFPValueKind:
		return ConstFloat(ref)
	case binding.LLVMConstantArrayValueKind:
		return ConstArray(ref)
	case binding.LLVMConstantStructValueKind:
		return ConstStruct(ref)
	case binding.LLVMConstantPointerNullValueKind:
		return ConstPointer(ref)
	case binding.LLVMConstantAggregateZeroValueKind:
		switch typeKind := binding.LLVMGetTypeKind(binding.LLVMTypeOf(ref)); typeKind {
		case binding.LLVMArrayTypeKind:
			return ConstArray(ref)
		case binding.LLVMStructTypeKind:
			return ConstStruct(ref)
		default:
			panic(fmt.Errorf("unknown type `%d`", typeKind))
		}
	case binding.LLVMConstantExprValueKind:
		switch opcode := binding.LLVMGetConstOpcode(ref); opcode {
		case binding.LLVMExtractElement:
			return ConstExtractElement(ref)
		case binding.LLVMGetElementPtr:
			return ConstGetElementPtr(ref)
		default:
			panic(fmt.Errorf("unknown opcode `%d`", opcode))
		}
	default:
		panic(fmt.Errorf("unknown constant `%d`", constKind))
	}
}

func (ctx Context) ConstNull(t Type) Constant {
	return lookupConstant(binding.LLVMConstNull(t.binding()))
}

func (ctx Context) ConstAggregateZero(t AggregateType) Constant {
	return lookupConstant(binding.LLVMConstAggregateZero(t.binding()))
}

func (ctx Context) ConstZero(t Type) Constant {
	return ctx.ConstNull(t)
}

type ConstInteger binding.LLVMValueRef

func (ctx Context) ConstIntegerFromString(t IntegerType, s string, radix uint8) ConstInteger {
	return ConstInteger(binding.LLVMConstIntOfString(t.binding(), s, radix))
}

func (ctx Context) ConstInteger(t IntegerType, v int64) ConstInteger {
	return ctx.ConstIntegerFromString(t, strconv.FormatInt(v, 10), 10)
}

func (ctx Context) ConstBoolean(v bool) ConstInteger {
	if v {
		return ctx.ConstInteger(ctx.BooleanType(), 1)
	} else {
		return ctx.ConstInteger(ctx.BooleanType(), 0)
	}
}

func (c ConstInteger) String() string {
	return binding.LLVMPrintValueToString(c.binding())
}

func (c ConstInteger) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(c)
}

func (c ConstInteger) Type() Type {
	return lookupType(binding.LLVMTypeOf(c.binding()))
}

func (ConstInteger) constant() {}

func (c ConstInteger) SignedValue() int64 {
	return binding.LLVMConstIntGetSExtValue(c.binding())
}

func (c ConstInteger) UnsignedValue() uint64 {
	return binding.LLVMConstIntGetZExtValue(c.binding())
}

type ConstFloat binding.LLVMValueRef

func (ctx Context) ConstFloatFromString(t FloatType, s string) ConstFloat {
	return ConstFloat(binding.LLVMConstRealOfString(t.binding(), s))
}

func (ctx Context) ConstFloat(t FloatType, v float64) ConstFloat {
	return ConstFloat(binding.LLVMConstReal(t.binding(), v))
}

func (c ConstFloat) String() string {
	return binding.LLVMPrintValueToString(c.binding())
}

func (c ConstFloat) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(c)
}

func (c ConstFloat) Type() Type {
	return lookupType(binding.LLVMTypeOf(c.binding()))
}

func (ConstFloat) constant() {}

func (c ConstFloat) Value() float64 {
	v, _ := binding.LLVMConstRealGetDouble(c.binding())
	return v
}

type ConstArray binding.LLVMValueRef

func (ctx Context) ConstArray(et Type, elem ...Constant) ConstArray {
	var es []binding.LLVMValueRef
	if len(elem) > 0 {
		es = lo.Map(elem, func(item Constant, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return ConstArray(binding.LLVMConstArray(et.binding(), es))
}

func (ctx Context) ConstString(s string) ConstArray {
	return ConstArray(binding.LLVMConstStringInContext(ctx.binding(), s, false))
}

func (c ConstArray) String() string {
	return binding.LLVMPrintValueToString(c.binding())
}

func (c ConstArray) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(c)
}

func (c ConstArray) Type() Type {
	return lookupType(binding.LLVMTypeOf(c.binding()))
}

func (ConstArray) constant() {}

func (c ConstArray) GetElem(i uint) Constant {
	return lookupConstant(binding.LLVMGetAggregateElement(c.binding(), uint32(i)))
}

type ConstStruct binding.LLVMValueRef

func (ctx Context) ConstStruct(packed bool, elem ...Constant) ConstStruct {
	var es []binding.LLVMValueRef
	if len(elem) > 0 {
		es = lo.Map(elem, func(item Constant, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return ConstStruct(binding.LLVMConstStructInContext(ctx.binding(), es, packed))
}

func (ctx Context) ConstNamedStruct(t StructType, elem ...Constant) ConstStruct {
	var es []binding.LLVMValueRef
	if len(elem) > 0 {
		es = lo.Map(elem, func(item Constant, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return ConstStruct(binding.LLVMConstNamedStruct(t.binding(), es))
}

func (c ConstStruct) String() string {
	return binding.LLVMPrintValueToString(c.binding())
}

func (c ConstStruct) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(c)
}

func (c ConstStruct) Type() Type {
	return lookupType(binding.LLVMTypeOf(c.binding()))
}

func (ConstStruct) constant() {}

func (c ConstStruct) GetElem(i uint) Constant {
	return lookupConstant(binding.LLVMGetAggregateElement(c.binding(), uint32(i)))
}

type ConstPointer binding.LLVMValueRef

func (ctx Context) ConstPointer(t Type) ConstPointer {
	return ConstPointer(binding.LLVMConstPointerNull(t.binding()))
}

func (c ConstPointer) String() string {
	return binding.LLVMPrintValueToString(c.binding())
}

func (c ConstPointer) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(c)
}

func (c ConstPointer) Type() Type {
	return lookupType(binding.LLVMTypeOf(c.binding()))
}

func (ConstPointer) constant() {}

type ConstGetElementPtr binding.LLVMValueRef

func (ctx Context) ConstGEP(t Type, v Constant, indice ...Constant) ConstGetElementPtr {
	indices := lo.Map(indice, func(item Constant, _ int) binding.LLVMValueRef {
		return item.binding()
	})
	return ConstGetElementPtr(binding.LLVMConstGEP(t.binding(), v.binding(), indices))
}

func (ctx Context) ConstInBoundsGEP(t Type, v Constant, indice ...Constant) ConstGetElementPtr {
	indices := lo.Map(indice, func(item Constant, _ int) binding.LLVMValueRef {
		return item.binding()
	})
	return ConstGetElementPtr(binding.LLVMConstInBoundsGEP(t.binding(), v.binding(), indices))
}

func (c ConstGetElementPtr) String() string {
	return binding.LLVMPrintValueToString(c.binding())
}

func (c ConstGetElementPtr) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(c)
}

func (c ConstGetElementPtr) Type() Type {
	return lookupType(binding.LLVMTypeOf(c.binding()))
}

func (ConstGetElementPtr) constant() {}

type ConstExtractElement binding.LLVMValueRef

func (ctx Context) ConstExtractElement(v, index Constant) ConstExtractElement {
	return ConstExtractElement(binding.LLVMConstExtractElement(v.binding(), index.binding()))
}

func (c ConstExtractElement) String() string {
	return binding.LLVMPrintValueToString(c.binding())
}

func (c ConstExtractElement) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(c)
}

func (c ConstExtractElement) Type() Type {
	return lookupType(binding.LLVMTypeOf(c.binding()))
}

func (ConstExtractElement) constant() {}
