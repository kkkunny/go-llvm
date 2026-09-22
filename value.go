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
	String() string
	Name() string
	SetName(name string)
	IsNil() bool
	IsConstant() bool
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

// Alive 值是否可用（所属 Context 与生命周期令牌均存活）
func (v Value[T]) Alive() bool {
	return v.ref.IsNil() || (v.ctx.Alive() && v.life.Alive())
}

// Context 返回所属上下文
func (v Value[T]) Context() *Context { return v.ctx }

// IsNil 是否为空句柄
func (v Value[T]) IsNil() bool { return v.ref.IsNil() }

// String 值的 IR 文本表示
func (v Value[T]) String() string {
	if v.ref.IsNil() {
		return "<nil>"
	}
	return binding.LLVMPrintValueToString(v.ref)
}

// Name 值名称
func (v Value[T]) Name() string {
	if v.ref.IsNil() {
		return ""
	}
	return binding.LLVMGetValueName(v.ref)
}

// SetName 设置值名称
func (v Value[T]) SetName(name string) {
	binding.LLVMSetValueName(v.ref, name)
}

// IsConstant 是否常量
func (v Value[T]) IsConstant() bool {
	if v.ref.IsNil() {
		return false
	}
	return binding.LLVMIsConstant(v.ref)
}

// Type 值的类型（依赖类型：T 与值种类一致）
func (v Value[T]) Type() Type[T] {
	return Type[T]{ref: binding.LLVMTypeOf(v.ref), ctx: v.ctx}
}

// As 运行时校验种类后转换类型参数；目标是 DynT 时始终成功
func (v Value[T]) As[U Kind]() (Value[U], error) {
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
