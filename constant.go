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
	switch constKind := binding.LLVMGetValueKind(ref); constKind {
	case binding.LLVMConstantIntValueKind:
		return Integer(ref)
	case binding.LLVMConstantFPValueKind:
		return Float(ref)
	case binding.LLVMConstantArrayValueKind:
		return Array(ref)
	case binding.LLVMConstantStructValueKind:
		return Struct(ref)
	case binding.LLVMConstantPointerNullValueKind:
		return Pointer(ref)
	case binding.LLVMConstantAggregateZeroValueKind:
		switch typeKind := binding.LLVMGetTypeKind(binding.LLVMTypeOf(ref)); typeKind {
		case binding.LLVMArrayTypeKind:
			return Array(ref)
		case binding.LLVMStructTypeKind:
			return Struct(ref)
		default:
			panic(fmt.Errorf("unknown enum value `%d`", typeKind))
		}
	default:
		panic(fmt.Errorf("unknown enum value `%d`", constKind))
	}
}

func (ctx Context) ConstNull(t Type) Constant {
	return lookupConstant(binding.LLVMConstNull(t.binding()))
}

func (ctx Context) ConstAggregateZero(t AggregateType) Constant {
	return lookupConstant(binding.LLVMConstAggregateZero(t.binding()))
}

type Integer binding.LLVMValueRef

func (ctx Context) ConstIntegerFromString(t IntegerType, s string, radix uint8) Integer {
	return Integer(binding.LLVMConstIntOfString(t.binding(), s, radix))
}

func (ctx Context) ConstInteger(t IntegerType, v int64) Integer {
	return ctx.ConstIntegerFromString(t, strconv.FormatInt(v, 10), 10)
}

func (c Integer) String() string {
	return binding.LLVMPrintValueToString(c.binding())
}

func (c Integer) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(c)
}

func (c Integer) Type() Type {
	return lookupType(binding.LLVMTypeOf(c.binding()))
}

func (Integer) constant() {}

type Float binding.LLVMValueRef

func (ctx Context) ConstFloatFromString(t FloatType, s string) Float {
	return Float(binding.LLVMConstRealOfString(t.binding(), s))
}

func (ctx Context) ConstFloat(t FloatType, v float64) Float {
	return Float(binding.LLVMConstReal(t.binding(), v))
}

func (c Float) String() string {
	return binding.LLVMPrintValueToString(c.binding())
}

func (c Float) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(c)
}

func (c Float) Type() Type {
	return lookupType(binding.LLVMTypeOf(c.binding()))
}

func (Float) constant() {}

func (c Float) Value() float64 {
	v, _ := binding.LLVMConstRealGetDouble(c.binding())
	return v
}

type Array binding.LLVMValueRef

func (ctx Context) ConstArray(et Type, elem ...Constant) Array {
	var es []binding.LLVMValueRef
	if len(elem) > 0 {
		es = lo.Map(elem, func(item Constant, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return Array(binding.LLVMConstArray(et.binding(), es))
}

func (ctx Context) ConstString(s string) Array {
	return Array(binding.LLVMConstStringInContext(ctx.binding(), s, false))
}

func (c Array) String() string {
	return binding.LLVMPrintValueToString(c.binding())
}

func (c Array) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(c)
}

func (c Array) Type() Type {
	return lookupType(binding.LLVMTypeOf(c.binding()))
}

func (Array) constant() {}

func (c Array) GetElem(i uint) Constant {
	return lookupConstant(binding.LLVMGetAggregateElement(c.binding(), uint32(i)))
}

type Struct binding.LLVMValueRef

func (ctx Context) ConstStruct(packed bool, elem ...Constant) Struct {
	var es []binding.LLVMValueRef
	if len(elem) > 0 {
		es = lo.Map(elem, func(item Constant, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return Struct(binding.LLVMConstStructInContext(ctx.binding(), es, packed))
}

func (ctx Context) ConstNamedStruct(t StructType, elem ...Constant) Struct {
	var es []binding.LLVMValueRef
	if len(elem) > 0 {
		es = lo.Map(elem, func(item Constant, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return Struct(binding.LLVMConstNamedStruct(t.binding(), es))
}

func (c Struct) String() string {
	return binding.LLVMPrintValueToString(c.binding())
}

func (c Struct) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(c)
}

func (c Struct) Type() Type {
	return lookupType(binding.LLVMTypeOf(c.binding()))
}

func (Struct) constant() {}

func (c Struct) GetElem(i uint) Constant {
	return lookupConstant(binding.LLVMGetAggregateElement(c.binding(), uint32(i)))
}

type Pointer binding.LLVMValueRef

func (ctx Context) ConstPointer(t Type) Pointer {
	return Pointer(binding.LLVMConstPointerNull(t.binding()))
}

func (c Pointer) String() string {
	return binding.LLVMPrintValueToString(c.binding())
}

func (c Pointer) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(c)
}

func (c Pointer) Type() Type {
	return lookupType(binding.LLVMTypeOf(c.binding()))
}

func (Pointer) constant() {}

func (ctx Context) ConstGEP(t Type, v Constant, indice ...Constant) Constant {
	indices := lo.Map(indice, func(item Constant, _ int) binding.LLVMValueRef {
		return item.binding()
	})
	return lookupConstant(binding.LLVMConstGEP(t.binding(), v.binding(), indices))
}

func (ctx Context) ConstInBoundsGEP(t Type, v Constant, indice ...Constant) Constant {
	indices := lo.Map(indice, func(item Constant, _ int) binding.LLVMValueRef {
		return item.binding()
	})
	return lookupConstant(binding.LLVMConstInBoundsGEP(t.binding(), v.binding(), indices))
}

func (ctx Context) ConstExtractElement(v, index Constant) Constant {
	return lookupConstant(binding.LLVMConstExtractElement(v.binding(), index.binding()))
}
