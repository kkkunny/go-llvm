package llvm

import (
	"github.com/kkkunny/go-llvm/internal/binding"
)

// IntConst 整数常量角色
type IntConst struct{ Value[IntT] }

// SignedValue 有符号值
func (c IntConst) SignedValue() int64 { return binding.LLVMConstIntGetSExtValue(c.ref) }

// UnsignedValue 无符号值
func (c IntConst) UnsignedValue() uint64 { return binding.LLVMConstIntGetZExtValue(c.ref) }

// IsNegative 是否有符号语义下为负
func (c IntConst) IsNegative() bool { return c.SignedValue() < 0 }

// FloatConst 浮点常量角色
type FloatConst struct{ Value[FloatT] }

// FloatValue 浮点值（方法名避开内嵌字段 Value）
func (c FloatConst) FloatValue() float64 {
	v, _ := binding.LLVMConstRealGetDouble(c.ref)
	return v
}

// ConstInt 构造整数常量
func (ctx *Context) ConstInt(t IntType, v uint64, signed bool) IntConst {
	ctx.checkType("llvm.Context.ConstInt", t)
	return IntConst{Value[IntT]{ref: binding.LLVMConstInt(t.ref, v, signed), ctx: ctx, life: ctx.life}}
}

// ConstIntOfString 按进制解析字符串构造整数常量
func (ctx *Context) ConstIntOfString(t IntType, s string, radix uint8) IntConst {
	ctx.checkType("llvm.Context.ConstIntOfString", t)
	return IntConst{Value[IntT]{ref: binding.LLVMConstIntOfString(t.ref, s, radix), ctx: ctx, life: ctx.life}}
}

// ConstBool 构造布尔常量
func (ctx *Context) ConstBool(v bool) Value[IntT] {
	var n uint64
	if v {
		n = 1
	}
	return Value[IntT]{ref: binding.LLVMConstInt(ctx.Bool().ref, n, false), ctx: ctx, life: ctx.life}
}

// ConstFloat 构造浮点常量
func (ctx *Context) ConstFloat(t FloatType, v float64) FloatConst {
	ctx.checkType("llvm.Context.ConstFloat", t)
	return FloatConst{Value[FloatT]{ref: binding.LLVMConstReal(t.ref, v), ctx: ctx, life: ctx.life}}
}

// ConstNull 构造指定类型的 null 常量（泛型方法，接受类型角色或裸 Type[T]）
func (ctx *Context) ConstNull[T Kind](t TypeRef[T]) Value[T] {
	tt := t.AsType()
	ctx.checkType("llvm.Context.ConstNull", tt)
	return Value[T]{ref: binding.LLVMConstNull(tt.ref), ctx: ctx, life: ctx.life}
}

// ConstZero 构造指定类型的零值常量（泛型方法）；聚合类型得到 zeroinitializer
func (ctx *Context) ConstZero[T Kind](t TypeRef[T]) Value[T] {
	tt := t.AsType()
	ctx.checkType("llvm.Context.ConstZero", tt)
	var ref binding.LLVMValueRef
	switch kindOfType(tt.ref).(type) {
	case StructT, ArrayT, VecT:
		ref = binding.LLVMConstAggregateZero(tt.ref)
	default:
		ref = binding.LLVMConstNull(tt.ref)
	}
	return Value[T]{ref: ref, ctx: ctx, life: ctx.life}
}

// Null 该类型的 null 常量（角色经内嵌 Type[T] 自动继承）
func (t Type[T]) Null() Value[T] {
	return t.ctx.ConstNull(t)
}

// Zero 该类型的零值常量（角色经内嵌 Type[T] 自动继承）
func (t Type[T]) Zero() Value[T] {
	return t.ctx.ConstZero(t)
}

// ConstString 构造字符串常量；nullTerminate 为 true 时末尾附加 \00
func (ctx *Context) ConstString(s string, nullTerminate bool) Value[ArrayT] {
	ctx.checkAlive("llvm.Context.ConstString")
	ref := binding.LLVMConstStringInContext(ctx.ref, s, !nullTerminate)
	return Value[ArrayT]{ref: ref, ctx: ctx, life: ctx.life}
}

// ConstArray 构造数组常量；元素类型/归属不符则 panic
func (ctx *Context) ConstArray(elem AnyType, elems ...AnyValue) Value[ArrayT] {
	ctx.checkType("llvm.Context.ConstArray", elem)
	ctx.checkValues("llvm.Context.ConstArray", elems...)
	for _, e := range elems {
		if !elem.Equal(TypeOfRef(ctx, binding.LLVMTypeOf(e.Ref()))) {
			errPanic(ErrTypeMismatch, "llvm.Context.ConstArray",
				"element type %s does not match array element type %s", TypeOfRef(ctx, binding.LLVMTypeOf(e.Ref())), elem)
		}
	}
	ref := binding.LLVMConstArray(elem.Ref(), anyValuesToRefs(elems))
	return Value[ArrayT]{ref: ref, ctx: ctx, life: ctx.life}
}

// ConstStruct 构造字面量结构体常量
func (ctx *Context) ConstStruct(packed bool, elems ...AnyValue) Value[StructT] {
	ctx.checkValues("llvm.Context.ConstStruct", elems...)
	ref := binding.LLVMConstStructInContext(ctx.ref, anyValuesToRefs(elems), packed)
	return Value[StructT]{ref: ref, ctx: ctx, life: ctx.life}
}

// ConstNamedStruct 构造命名结构体常量；元素个数/类型不符则 panic
func (ctx *Context) ConstNamedStruct(t StructType, elems ...AnyValue) Value[StructT] {
	ctx.checkType("llvm.Context.ConstNamedStruct", t)
	ctx.checkValues("llvm.Context.ConstNamedStruct", elems...)
	if got, want := len(elems), len(t.Elems()); got != want {
		errPanic(ErrTypeMismatch, "llvm.Context.ConstNamedStruct", "expect %d elements, got %d", want, got)
	}
	for i, e := range elems {
		if !t.Elem(uint32(i)).Equal(TypeOfRef(ctx, binding.LLVMTypeOf(e.Ref()))) {
			errPanic(ErrTypeMismatch, "llvm.Context.ConstNamedStruct",
				"element %d type %s does not match field type %s", i, TypeOfRef(ctx, binding.LLVMTypeOf(e.Ref())), t.Elem(uint32(i)))
		}
	}
	ref := binding.LLVMConstNamedStruct(t.ref, anyValuesToRefs(elems))
	return Value[StructT]{ref: ref, ctx: ctx, life: ctx.life}
}

// ConstGEP 构造常量 GEP 表达式；elem 为源元素类型，base 必须是指针值
func (ctx *Context) ConstGEP(elem AnyType, base ValueRef[PtrT], inBounds bool, idx ...ValueRef[IntT]) Value[PtrT] {
	ctx.checkType("llvm.Context.ConstGEP", elem)
	baseV := base.AsValue()
	ctx.checkValues("llvm.Context.ConstGEP", baseV.Dyn())
	idxRefs := make([]binding.LLVMValueRef, len(idx))
	for i, x := range idx {
		v := x.AsValue()
		ctx.checkValues("llvm.Context.ConstGEP", v.Dyn())
		idxRefs[i] = v.ref
	}
	var ref binding.LLVMValueRef
	if inBounds {
		ref = binding.LLVMConstInBoundsGEP(elem.Ref(), baseV.ref, idxRefs)
	} else {
		ref = binding.LLVMConstGEP(elem.Ref(), baseV.ref, idxRefs)
	}
	return Value[PtrT]{ref: ref, ctx: ctx, life: ctx.life}
}

// ConstIntToPtr 构造 inttoptr 常量表达式
func (ctx *Context) ConstIntToPtr(v ValueRef[IntT], to PtrType) Value[PtrT] {
	vv := v.AsValue()
	ctx.checkValues("llvm.Context.ConstIntToPtr", vv.Dyn())
	ctx.checkType("llvm.Context.ConstIntToPtr", to)
	ref := binding.LLVMConstIntToPtr(vv.Ref(), to.Ref())
	return Value[PtrT]{ref: ref, ctx: ctx, life: ctx.life}
}

// ===== 内部辅助 =====

// checkValues 校验值归属同一 Context 且存活
func (ctx *Context) checkValues(op string, vs ...AnyValue) {
	for _, v := range vs {
		if v == nil || v.IsNil() {
			errPanic(ErrInvalidArg, op, "nil value")
		}
		if !v.Alive() {
			errPanic(ErrUseAfterFree, op, "value is freed")
		}
		if v.Context() != ctx {
			errPanic(ErrCrossContext, op, "value belongs to another context")
		}
	}
}

func anyValuesToRefs(values []AnyValue) []binding.LLVMValueRef {
	refs := make([]binding.LLVMValueRef, len(values))
	for i, v := range values {
		refs[i] = v.Ref()
	}
	return refs
}

func intValuesToRefs(values []Value[IntT]) []binding.LLVMValueRef {
	refs := make([]binding.LLVMValueRef, len(values))
	for i, v := range values {
		refs[i] = v.ref
	}
	return refs
}
