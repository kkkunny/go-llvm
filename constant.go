package llvm

import (
	"fmt"
	"strings"

	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/errs"
)

// IntConst 整数常量角色
type IntConst struct{ Value[IntT] }

// SignedValue 有符号值
func (c IntConst) SignedValue() int64 {
	return binding.LLVMConstIntGetSExtValue(c.Ref())
}

// UnsignedValue 无符号值
func (c IntConst) UnsignedValue() uint64 {
	return binding.LLVMConstIntGetZExtValue(c.Ref())
}

// IsNegative 是否有符号语义下为负
func (c IntConst) IsNegative() bool {
	return c.SignedValue() < 0
}

// FloatConst 浮点常量角色
type FloatConst struct{ Value[FloatT] }

// FloatValue 浮点值（方法名避开内嵌字段 Value）
func (c FloatConst) FloatValue() float64 {
	v, _ := binding.LLVMConstRealGetDouble(c.Ref())
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

// ConstBool 构造布尔常量（i1，返回 Value[IntT]；LLVM 无独立 bool 类型，bool 即 i1）
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

// Undef 该类型的 undef 常量（角色经内嵌 Type[T] 自动继承）
func (t Type[T]) Undef() Value[T] {
	t.Check("llvm.Type.Undef")
	return newValue[T](t.ctx, t.ctx.life, binding.LLVMGetUndef(t.ref))
}

// Poison 该类型的 poison 常量（角色经内嵌 Type[T] 自动继承）
func (t Type[T]) Poison() Value[T] {
	t.Check("llvm.Type.Poison")
	return newValue[T](t.ctx, t.ctx.life, binding.LLVMGetPoison(t.ref))
}

// ConstString 构造字符串常量；nullTerminate 为 true 时末尾附加 \00
func (ctx *Context) ConstString(s string, nullTerminate bool) Value[ArrayT] {
	ctx.CheckAlive("llvm.Context.ConstString")
	ref := binding.LLVMConstStringInContext(ctx.ref, s, !nullTerminate)
	return newValue[ArrayT](ctx, ctx.life, ref)
}

// ConstArray 构造数组常量；元素非常量或类型/归属不符则 panic
func (ctx *Context) ConstArray(elem AnyType, elems ...AnyValue) Value[ArrayT] {
	const op = "llvm.Context.ConstArray"
	ctx.CheckType(op, elem)
	ctx.CheckValues(op, elems...)
	elemRef := elem.Ref()
	for i, e := range elems {
		checkConstOperand(op, i, e)
		if ref := binding.LLVMTypeOf(e.Ref()); !elemRef.Equal(ref) {
			errs.Panicf(ErrTypeMismatch, op,
				"element type %s does not match array element type %s", TypeOfRef(ctx, ref), elem)
		}
	}
	ref := binding.LLVMConstArray(elemRef, AnyValuesToRefs(elems))
	return newValue[ArrayT](ctx, ctx.life, ref)
}

// ConstVector 构造向量常量；元素非常量或类型/归属不符则 panic
func (ctx *Context) ConstVector(elem AnyType, elems ...AnyValue) Value[VecT] {
	const op = "llvm.Context.ConstVector"
	ctx.CheckType(op, elem)
	ctx.CheckValues(op, elems...)
	elemRef := elem.Ref()
	for i, e := range elems {
		checkConstOperand(op, i, e)
		if ref := binding.LLVMTypeOf(e.Ref()); !elemRef.Equal(ref) {
			errs.Panicf(ErrTypeMismatch, op,
				"element %d type %s does not match vector element type %s", i, TypeOfRef(ctx, ref), elem)
		}
	}
	ref := binding.LLVMConstVector(AnyValuesToRefs(elems))
	return newValue[VecT](ctx, ctx.life, ref)
}

// ConstStruct 构造字面量结构体常量；元素非常量则 panic
func (ctx *Context) ConstStruct(packed bool, elems ...AnyValue) Value[StructT] {
	const op = "llvm.Context.ConstStruct"
	ctx.CheckValues(op, elems...)
	for i, e := range elems {
		checkConstOperand(op, i, e)
	}
	ref := binding.LLVMConstStructInContext(ctx.ref, AnyValuesToRefs(elems), packed)
	return newValue[StructT](ctx, ctx.life, ref)
}

// ConstNamedStruct 构造命名结构体常量；元素非常量或个数/类型不符则 panic
func (ctx *Context) ConstNamedStruct(t StructType, elems ...AnyValue) Value[StructT] {
	const op = "llvm.Context.ConstNamedStruct"
	ctx.CheckType(op, t)
	ctx.CheckValues(op, elems...)
	if got, want := len(elems), int(binding.LLVMCountStructElementTypes(t.ref)); got != want {
		errs.Panicf(ErrTypeMismatch, op, "expect %d elements, got %d", want, got)
	}
	for i, e := range elems {
		checkConstOperand(op, i, e)
		want := binding.LLVMStructGetTypeAtIndex(t.ref, uint32(i))
		if ref := binding.LLVMTypeOf(e.Ref()); !want.Equal(ref) {
			errs.Panicf(ErrTypeMismatch, op,
				"element %d type %s does not match field type %s", i, TypeOfRef(ctx, ref), TypeOfRef(ctx, want))
		}
	}
	ref := binding.LLVMConstNamedStruct(t.ref, AnyValuesToRefs(elems))
	return newValue[StructT](ctx, ctx.life, ref)
}

// ConstGEP 构造常量 GEP 表达式；elem 为源元素类型，base 与索引都必须是常量
func (ctx *Context) ConstGEP(elem AnyType, base ValueRef[PtrT], inBounds bool, idx ...ValueRef[IntT]) Value[PtrT] {
	const op = "llvm.Context.ConstGEP"
	ctx.CheckType(op, elem)
	baseV := base.AsValue()
	ctx.checkValueOwn(op, baseV.ref, baseV.life, baseV.ctx)
	checkConstOperand(op, -1, baseV)
	idxRefs := make([]binding.LLVMValueRef, len(idx))
	for i, x := range idx {
		v := x.AsValue()
		ctx.checkValueOwn(op, v.ref, v.life, v.ctx)
		checkConstOperand(op, i, v)
		idxRefs[i] = v.ref
	}
	var ref binding.LLVMValueRef
	if inBounds {
		ref = binding.LLVMConstInBoundsGEP2(elem.Ref(), baseV.ref, idxRefs)
	} else {
		ref = binding.LLVMConstGEP2(elem.Ref(), baseV.ref, idxRefs)
	}
	return newValue[PtrT](ctx, ctx.life, ref)
}

// SizeOf 构造类型分配大小的 i64 常量表达式（语义对应 C API 的 LLVMSizeOf）。
//
// 结果是目标无关的常量表达式，在代码生成期按值所在模块的数据布局折叠为分配大小
// （getTypeAllocSize，byte）；t 须在代码生成前 sized，否则该表达式无法折叠。
// t 为 nil 或属于其他 Context 时 panic。
func (ctx *Context) SizeOf(t AnyType) Value[IntT] {
	const op = "llvm.Context.SizeOf"
	ctx.CheckType(op, t)
	return newValue[IntT](ctx, ctx.life, binding.LLVMSizeOf(t.Ref()))
}

// AlignOf 构造类型对齐的 i64 常量表达式（语义对应 C API 的 LLVMAlignOf）。
//
// 结果是目标无关的常量表达式，在代码生成期按值所在模块的数据布局折叠为 ABI 对齐
// （getTypeABIAlignment，byte）；t 须在代码生成前 sized，否则该表达式无法折叠。
// t 为 nil 或属于其他 Context 时 panic。
func (ctx *Context) AlignOf(t AnyType) Value[IntT] {
	const op = "llvm.Context.AlignOf"
	ctx.CheckType(op, t)
	return newValue[IntT](ctx, ctx.life, binding.LLVMAlignOf(t.Ref()))
}

// ConstIntToPtr 构造 inttoptr 常量表达式
func (ctx *Context) ConstIntToPtr(v ValueRef[IntT], to PtrType) Value[PtrT] {
	const op = "llvm.Context.ConstIntToPtr"
	vv := v.AsValue()
	ctx.checkValueOwn(op, vv.ref, vv.life, vv.ctx)
	checkConstOperand(op, -1, vv)
	ctx.CheckType(op, to)
	ref := binding.LLVMConstIntToPtr(vv.Ref(), to.Ref())
	return newValue[PtrT](ctx, ctx.life, ref)
}

// ConstBitCast 构造 bitcast 常量表达式（源与目标须位宽一致，见 LLVM 常量转换规则）
func (ctx *Context) ConstBitCast[U Kind](v AnyValue, to TypeRef[U]) Value[U] {
	const op = "llvm.Context.ConstBitCast"
	ctx.CheckValues(op, v)
	checkConstOperand(op, -1, v)
	tt := to.AsType()
	ctx.CheckType(op, tt)
	ref := binding.LLVMConstBitCast(v.Ref(), tt.Ref())
	return newValue[U](ctx, ctx.life, ref)
}

// ConstPointerCast 构造指针 cast 常量表达式（不透明指针下等价于 bitcast，地址空间须相同）
func (ctx *Context) ConstPointerCast(v ValueRef[PtrT], to PtrType) Value[PtrT] {
	const op = "llvm.Context.ConstPointerCast"
	vv := v.AsValue()
	ctx.checkValueOwn(op, vv.ref, vv.life, vv.ctx)
	checkConstOperand(op, -1, vv)
	ctx.CheckType(op, to)
	ref := binding.LLVMConstPointerCast(vv.Ref(), to.Ref())
	return newValue[PtrT](ctx, ctx.life, ref)
}

// ===== 内部辅助 =====

// checkConstOperand 校验常量构造的操作数是常量表达式；非常量则 panic。
// 索引 i 为负表示单一操作数（如 GEP 的 base），消息里只说"operand"。
// 常量构造器把操作数原样交给 LLVM：非常量会构造出"常量里包含指令"的畸形值，
// 之后 Module.Verify 的报错指向那条指令本身而非构造调用处，定位成本极高，故在源头提前失败。
// 与 ConstArray/ConstNamedStruct 的类型校验一致：冷路径，恒定校验（不加调试开关）。
func checkConstOperand(op string, i int, v AnyValue) {
	ref := v.Ref()
	if binding.LLVMIsConstant(ref) {
		return
	}
	desc := "operand"
	if i >= 0 {
		desc = fmt.Sprintf("operand %d", i)
	}
	if name := v.Name(); name != "" {
		desc += " @" + name
	}
	errs.Panicf(ErrInvalidArg, op, "%s is not a constant: %s", desc, strings.TrimSpace(binding.LLVMPrintValueToString(ref)))
}

// CheckValues 校验值归属同一 Context 且存活
func (ctx *Context) CheckValues(op string, vs ...AnyValue) {
	for _, v := range vs {
		if v == nil {
			errs.Panicf(ErrInvalidArg, op, "nil value")
		}
		ctx.checkValueOwn(op, v.Ref(), v.Lifetime(), v.Context())
	}
}

// checkValueOwn 校验具体值句柄归属本 Context 且存活（无装箱）
func (ctx *Context) checkValueOwn(op string, ref binding.LLVMValueRef, life *Lifetime, vctx *Context) {
	if ref.IsNil() {
		errs.Panicf(ErrInvalidArg, op, "nil value")
	}
	if !ctx.life.Alive() {
		errs.Panicf(ErrUseAfterFree, op, "context is closed")
	}
	if life == nil || !life.Alive() {
		errs.Panicf(ErrUseAfterFree, op, "value is freed")
	}
	if vctx != ctx {
		errs.Panicf(ErrCrossContext, op, "value belongs to another context")
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
