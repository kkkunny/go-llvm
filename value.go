package llvm

import (
	"github.com/kkkunny/go-llvm/internal/binding"
)

// AnyValue 值句柄的非泛型视图
type AnyValue interface {
	Ref() binding.LLVMValueRef
	Dyn() Value[DynT]
	Alive() bool
	Context() *Context
	Lifetime() *Lifetime
	String() string
	Name() string
	SetName(name string)
	IsNil() bool
	IsConstant() bool
}

// ValueRef 种类安全的值引用：Value[T] 与全部值角色（IntConst/Phi[T]/...）均实现。
// 用于泛型方法/函数参数，使泛型调用可直接推断 T。
type ValueRef[T Kind] interface {
	AsValue() Value[T]
}

// Value 种类级泛型值句柄；类别专属操作见角色包装（如 IntConst.SignedValue）
type Value[T Kind] struct {
	ref  binding.LLVMValueRef
	ctx  *Context
	life *Lifetime
}

// Ref 返回底层句柄（供 llvm/* 子包桥接使用）
func (v Value[T]) Ref() binding.LLVMValueRef { return v.ref }

// Dyn 擦除类型参数
func (v Value[T]) Dyn() Value[DynT] {
	return Value[DynT]{ref: v.ref, ctx: v.ctx, life: v.life}
}

// AsValue 返回自身（实现 ValueRef[T]；值角色经内嵌继承）
func (v Value[T]) AsValue() Value[T] { return v }

// Alive 值是否可用（非空句柄且所属 Context 与生命周期令牌均存活）
func (v Value[T]) Alive() bool {
	return !v.ref.IsNil() && v.ctx != nil && v.ctx.Alive() && v.life != nil && v.life.Alive()
}

// Context 返回所属上下文
func (v Value[T]) Context() *Context { return v.ctx }

// Lifetime 返回所属生命周期令牌
func (v Value[T]) Lifetime() *Lifetime { return v.life }

// IsNil 是否为空句柄
func (v Value[T]) IsNil() bool { return v.ref.IsNil() }

// Check 值操作前置校验：句柄非零值且 Context/生命周期均存活。
// 供 llvm/* 子包的值角色方法统一调用（角色经内嵌 Value[T] 自动继承）。
func (v Value[T]) Check(op string) {
	if v.ref.IsNil() {
		errPanic(ErrInvalidArg, op, "nil value handle")
	}
	if v.ctx == nil || !v.ctx.Alive() {
		errPanic(ErrUseAfterFree, op, "context is closed")
	}
	if v.life == nil || !v.life.Alive() {
		errPanic(ErrUseAfterFree, op, "value is freed")
	}
}

// String 值的 IR 文本表示
func (v Value[T]) String() string {
	if v.ref.IsNil() {
		return "<nil>"
	}
	v.Check("llvm.Value.String")
	return binding.LLVMPrintValueToString(v.ref)
}

// Name 值名称
func (v Value[T]) Name() string {
	v.Check("llvm.Value.Name")
	return binding.LLVMGetValueName(v.ref)
}

// SetName 设置值名称
func (v Value[T]) SetName(name string) {
	v.Check("llvm.Value.SetName")
	binding.LLVMSetValueName(v.ref, name)
}

// IsConstant 是否常量
func (v Value[T]) IsConstant() bool {
	v.Check("llvm.Value.IsConstant")
	return binding.LLVMIsConstant(v.ref)
}

// Type 值的类型（依赖类型：T 与值种类一致）
func (v Value[T]) Type() Type[T] {
	v.Check("llvm.Value.Type")
	return Type[T]{ref: binding.LLVMTypeOf(v.ref), ctx: v.ctx}
}

// As 运行时校验种类后转换类型参数；目标是 DynT 时始终成功
func (v Value[T]) As[U Kind]() (Value[U], error) {
	v.Check("llvm.Value.As")
	if !kindMatches[U](binding.LLVMTypeOf(v.ref)) {
		return Value[U]{}, &Error{
			Reason: ErrTypeMismatch,
			Op:     "llvm.Value.As",
			Msg:    "value kind mismatch: have " + kindName(kindOfType(binding.LLVMTypeOf(v.ref))) + ", want " + kindName(kindOf[U]()),
		}
	}
	return Value[U]{ref: v.ref, ctx: v.ctx, life: v.life}, nil
}

// MustAs As 的 panic 版本（程序员错误）
func (v Value[T]) MustAs[U Kind]() Value[U] {
	res, err := v.As[U]()
	if err != nil {
		panic(err)
	}
	return res
}

// NewValue 由底层句柄构建值（供 llvm/* 子包桥接使用）
func NewValue[T Kind](ctx *Context, life *Lifetime, ref binding.LLVMValueRef) Value[T] {
	return Value[T]{ref: ref, ctx: ctx, life: life}
}

// ValueOf 由底层句柄构建擦除值（指令遍历等动态来源）
func ValueOf(ctx *Context, life *Lifetime, ref binding.LLVMValueRef) Value[DynT] {
	return Value[DynT]{ref: ref, ctx: ctx, life: life}
}
