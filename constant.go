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
	switch binding.LLVMGetValueKind(ref) {
	case binding.LLVMConstantIntValueKind:
		return Integer(ref)
	case binding.LLVMConstantFPValueKind:
		return Float(ref)
	case binding.LLVMConstantArrayValueKind:
		return Array(ref)
	case binding.LLVMConstantStructValueKind:
		return Struct(ref)
	default:
		panic(fmt.Errorf("unknown enum value `%d`", binding.LLVMGetValueKind(ref)))
	}
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

func (ctx Context) ConstArray(et Type, elems []Constant) Array {
	var es []binding.LLVMValueRef
	if len(elems) > 0 {
		es = lo.Map(elems, func(item Constant, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return Array(binding.LLVMConstArray(et.binding(), es))
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

func (ctx Context) ConstStruct(packed bool, elems []Constant) Struct {
	var es []binding.LLVMValueRef
	if len(elems) > 0 {
		es = lo.Map(elems, func(item Constant, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return Struct(binding.LLVMConstStructInContext(ctx.binding(), es, packed))
}

func (ctx Context) ConstNamedStruct(t StructType, elems []Constant) Struct {
	var es []binding.LLVMValueRef
	if len(elems) > 0 {
		es = lo.Map(elems, func(item Constant, index int) binding.LLVMValueRef {
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
