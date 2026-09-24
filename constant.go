package llvm

import (
	"github.com/kkkunny/go-llvm/internal/binding"
)

// IntConst 整数常量角色
type IntConst struct{ Value[IntT] }

// SignedValue 有符号值
func (c IntConst) SignedValue() int64 {
	c.Check("llvm.IntConst.SignedValue")
	return binding.LLVMConstIntGetSExtValue(c.ref)
}

// UnsignedValue 无符号值
func (c IntConst) UnsignedValue() uint64 {
	c.Check("llvm.IntConst.UnsignedValue")
	return binding.LLVMConstIntGetZExtValue(c.ref)
}

// IsNegative 是否有符号语义下为负
func (c IntConst) IsNegative() bool {
	c.Check("llvm.IntConst.IsNegative")
	return c.SignedValue() < 0
}

// FloatConst 浮点常量角色
type FloatConst struct{ Value[FloatT] }

// FloatValue 浮点值（方法名避开内嵌字段 Value）
func (c FloatConst) FloatValue() float64 {
	c.Check("llvm.FloatConst.FloatValue")
	v, _ := binding.LLVMConstRealGetDouble(c.ref)
	return v
}

// ConstInt 构造整数常量（值按无符号截断；负数请用 ConstSInt）
func (ctx *Context) ConstInt(t IntType, v uint64) IntConst {
	ctx.CheckType("llvm.Context.ConstInt", t)
	return IntConst{newValue[IntT](ctx, ctx.life, binding.LLVMConstInt(t.ref, v, false))}
}

// ConstSInt 构造有符号整数常量（负数与超宽值按符号扩展）
func (ctx *Context) ConstSInt(t IntType, v int64) IntConst {
	ctx.CheckType("llvm.Context.ConstSInt", t)
	return IntConst{newValue[IntT](ctx, ctx.life, binding.LLVMConstInt(t.ref, uint64(v), true))}
}

// Const 该类型的整数常量（类型导向糖：i32.Const(5)）
func (t IntType) Const(v uint64) IntConst { return t.ctx.ConstInt(t, v) }

// ConstS 该类型的有符号整数常量（类型导向糖：i32.ConstS(-1)）
func (t IntType) ConstS(v int64) IntConst { return t.ctx.ConstSInt(t, v) }

// Const 该类型的浮点常量（类型导向糖：f64.Const(3.14)）
func (t FloatType) Const(v float64) FloatConst { return t.ctx.ConstFloat(t, v) }

// ConstIntOfString 按进制解析字符串构造整数常量
func (ctx *Context) ConstIntOfString(t IntType, s string, radix uint8) IntConst {
	ctx.CheckType("llvm.Context.ConstIntOfString", t)
	return IntConst{newValue[IntT](ctx, ctx.life, binding.LLVMConstIntOfString(t.ref, s, radix))}
}

// ConstBool 构造布尔常量
func (ctx *Context) ConstBool(v bool) Value[IntT] {
	var n uint64
	if v {
		n = 1
	}
	return newValue[IntT](ctx, ctx.life, binding.LLVMConstInt(ctx.Bool().ref, n, false))
}

// ConstFloat 构造浮点常量
func (ctx *Context) ConstFloat(t FloatType, v float64) FloatConst {
	ctx.CheckType("llvm.Context.ConstFloat", t)
	return FloatConst{newValue[FloatT](ctx, ctx.life, binding.LLVMConstReal(t.ref, v))}
}

// ConstNull 构造指定类型的 null 常量（泛型方法，接受类型角色或裸 Type[T]）
func (ctx *Context) ConstNull[T Kind](t TypeRef[T]) Value[T] {
	tt := t.AsType()
	ctx.CheckType("llvm.Context.ConstNull", tt)
	return newValue[T](ctx, ctx.life, binding.LLVMConstNull(tt.ref))
}

// ConstZero 构造指定类型的零值常量（泛型方法）；聚合类型得到 zeroinitializer
func (ctx *Context) ConstZero[T Kind](t TypeRef[T]) Value[T] {
	tt := t.AsType()
	ctx.CheckType("llvm.Context.ConstZero", tt)
	var ref binding.LLVMValueRef
	switch kindOfType(tt.ref).(type) {
	case StructT, ArrayT, VecT:
		ref = binding.LLVMConstAggregateZero(tt.ref)
	default:
		ref = binding.LLVMConstNull(tt.ref)
	}
	return newValue[T](ctx, ctx.life, ref)
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
	ctx.CheckAlive("llvm.Context.ConstString")
	ref := binding.LLVMConstStringInContext(ctx.ref, s, !nullTerminate)
	return newValue[ArrayT](ctx, ctx.life, ref)
}

// ConstArray 构造数组常量；元素类型/归属不符则 panic
func (ctx *Context) ConstArray(elem AnyType, elems ...AnyValue) Value[ArrayT] {
	const op = "llvm.Context.ConstArray"
	ctx.CheckType(op, elem)
	ctx.CheckValues(op, elems...)
	elemRef := elem.Ref()
	for _, e := range elems {
		if ref := binding.LLVMTypeOf(e.Ref()); !elemRef.Equal(ref) {
			errPanic(ErrTypeMismatch, op,
				"element type %s does not match array element type %s", TypeOfRef(ctx, ref), elem)
		}
	}
	ref := binding.LLVMConstArray(elemRef, AnyValuesToRefs(elems))
	return newValue[ArrayT](ctx, ctx.life, ref)
}

// ConstVector 构造向量常量；元素类型/归属不符则 panic
func (ctx *Context) ConstVector(elem AnyType, elems ...AnyValue) Value[VecT] {
	const op = "llvm.Context.ConstVector"
	ctx.CheckType(op, elem)
	ctx.CheckValues(op, elems...)
	elemRef := elem.Ref()
	for i, e := range elems {
		if ref := binding.LLVMTypeOf(e.Ref()); !elemRef.Equal(ref) {
			errPanic(ErrTypeMismatch, op,
				"element %d type %s does not match vector element type %s", i, TypeOfRef(ctx, ref), elem)
		}
	}
	ref := binding.LLVMConstVector(AnyValuesToRefs(elems))
	return newValue[VecT](ctx, ctx.life, ref)
}

// ConstStruct 构造字面量结构体常量
func (ctx *Context) ConstStruct(packed bool, elems ...AnyValue) Value[StructT] {
	ctx.CheckValues("llvm.Context.ConstStruct", elems...)
	ref := binding.LLVMConstStructInContext(ctx.ref, AnyValuesToRefs(elems), packed)
	return newValue[StructT](ctx, ctx.life, ref)
}

// ConstNamedStruct 构造命名结构体常量；元素个数/类型不符则 panic
func (ctx *Context) ConstNamedStruct(t StructType, elems ...AnyValue) Value[StructT] {
	const op = "llvm.Context.ConstNamedStruct"
	ctx.CheckType(op, t)
	ctx.CheckValues(op, elems...)
	if got, want := len(elems), int(binding.LLVMCountStructElementTypes(t.ref)); got != want {
		errPanic(ErrTypeMismatch, op, "expect %d elements, got %d", want, got)
	}
	for i, e := range elems {
		want := binding.LLVMStructGetTypeAtIndex(t.ref, uint32(i))
		if ref := binding.LLVMTypeOf(e.Ref()); !want.Equal(ref) {
			errPanic(ErrTypeMismatch, op,
				"element %d type %s does not match field type %s", i, TypeOfRef(ctx, ref), TypeOfRef(ctx, want))
		}
	}
	ref := binding.LLVMConstNamedStruct(t.ref, AnyValuesToRefs(elems))
	return newValue[StructT](ctx, ctx.life, ref)
}

// ConstGEP 构造常量 GEP 表达式；elem 为源元素类型，base 必须是指针值
func (ctx *Context) ConstGEP(elem AnyType, base ValueRef[PtrT], inBounds bool, idx ...ValueRef[IntT]) Value[PtrT] {
	const op = "llvm.Context.ConstGEP"
	ctx.CheckType(op, elem)
	baseV := base.AsValue()
	ctx.checkValueOwn(op, baseV.ref, baseV.life, baseV.ctx)
	idxRefs := make([]binding.LLVMValueRef, len(idx))
	for i, x := range idx {
		v := x.AsValue()
		ctx.checkValueOwn(op, v.ref, v.life, v.ctx)
		idxRefs[i] = v.ref
	}
	var ref binding.LLVMValueRef
	if inBounds {
		ref = binding.LLVMConstInBoundsGEP(elem.Ref(), baseV.ref, idxRefs)
	} else {
		ref = binding.LLVMConstGEP(elem.Ref(), baseV.ref, idxRefs)
	}
	return newValue[PtrT](ctx, ctx.life, ref)
}

// ConstIntToPtr 构造 inttoptr 常量表达式
func (ctx *Context) ConstIntToPtr(v ValueRef[IntT], to PtrType) Value[PtrT] {
	const op = "llvm.Context.ConstIntToPtr"
	vv := v.AsValue()
	ctx.checkValueOwn(op, vv.ref, vv.life, vv.ctx)
	ctx.CheckType(op, to)
	ref := binding.LLVMConstIntToPtr(vv.Ref(), to.Ref())
	return newValue[PtrT](ctx, ctx.life, ref)
}

// ===== 内部辅助 =====

// CheckValues 校验值归属同一 Context 且存活
func (ctx *Context) CheckValues(op string, vs ...AnyValue) {
	for _, v := range vs {
		if v == nil {
			errPanic(ErrInvalidArg, op, "nil value")
		}
		ctx.checkValueOwn(op, v.Ref(), v.Lifetime(), v.Context())
	}
}

// checkValueOwn 校验具体值句柄归属本 Context 且存活（无装箱）
func (ctx *Context) checkValueOwn(op string, ref binding.LLVMValueRef, life *Lifetime, vctx *Context) {
	if ref.IsNil() {
		errPanic(ErrInvalidArg, op, "nil value")
	}
	if !ctx.life.Alive() {
		errPanic(ErrUseAfterFree, op, "context is closed")
	}
	if life == nil || !life.Alive() {
		errPanic(ErrUseAfterFree, op, "value is freed")
	}
	if vctx != ctx {
		errPanic(ErrCrossContext, op, "value belongs to another context")
	}
}

// AnyValuesToRefs 将值列表转换为底层句柄列表（供 llvm/* 子包桥接使用）
func AnyValuesToRefs(values []AnyValue) []binding.LLVMValueRef {
	refs := make([]binding.LLVMValueRef, len(values))
	for i, v := range values {
		refs[i] = v.Ref()
	}
	return refs
}
