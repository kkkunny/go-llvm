package llvm

import (
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
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
	ty   binding.LLVMTypeRef // 调试层类型缓存（构造时查询一次）；release 恒为空
	ctx  *Context
	life *Lifetime
}

// RawRef 返回底层句柄且不做校验；仅供预检查询（如 ir 的 preVal），业务路径请用 Ref
func (v Value[T]) RawRef() binding.LLVMValueRef { return v.ref }

// Ref 返回底层句柄：崩溃类地板的统一咽喉——nil / Context 已关 / 已释放在此 panic。
// 快路径保持可内联（慢路径单独走 checkFloor），热路径调用只付指针比较与 atomic load。
func (v Value[T]) Ref() binding.LLVMValueRef {
	if v.ref.IsNil() || v.ctx == nil || !v.ctx.Alive() || v.life == nil || !v.life.Alive() {
		v.checkFloor("llvm.Value.Ref")
	}
	return v.ref
}

// checkFloor 崩溃类地板校验（任何构建都开，纯 Go）
func (v Value[T]) checkFloor(op string) {
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

// RawType 返回构造时缓存的底层类型句柄（不做校验）；release 构建恒为空句柄。
// 供 llvm/* 子包预检在调试层复用，避免重复的 cgo 类型查询。
func (v Value[T]) RawType() binding.LLVMTypeRef { return v.ty }

// Dyn 擦除类型参数
func (v Value[T]) Dyn() Value[DynT] {
	return Value[DynT]{ref: v.ref, ty: v.ty, ctx: v.ctx, life: v.life}
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

// Check 值操作前置校验（显式入口，供需要指定 op 的场景）；与 Ref 同一套地板校验。
// 供 llvm/* 子包的值角色方法统一调用（角色经内嵌 Value[T] 自动继承）。
func (v Value[T]) Check(op string) { v.checkFloor(op) }

// String 值的 IR 文本表示
func (v Value[T]) String() string {
	if v.ref.IsNil() {
		return "<nil>"
	}
	v.checkFloor("llvm.Value.String")
	return binding.LLVMPrintValueToString(v.ref)
}

// Name 值名称
func (v Value[T]) Name() string {
	return binding.LLVMGetValueName(v.Ref())
}

// SetName 设置值名称
func (v Value[T]) SetName(name string) {
	binding.LLVMSetValueName(v.Ref(), name)
}

// IsConstant 是否常量
func (v Value[T]) IsConstant() bool {
	return binding.LLVMIsConstant(v.Ref())
}

// Type 值的类型（依赖类型：T 与值种类一致）
func (v Value[T]) Type() Type[T] {
	ref := v.Ref()
	ty := v.ty
	if ty.IsNil() {
		ty = binding.LLVMTypeOf(ref)
	}
	return Type[T]{ref: ty, ctx: v.ctx}
}

// As 运行时校验种类后转换类型参数；目标是 DynT 时始终成功
func (v Value[T]) As[U Kind]() (Value[U], error) {
	ref := v.Ref()
	ty := v.ty
	if ty.IsNil() {
		ty = binding.LLVMTypeOf(ref)
	}
	if !kindMatches[U](ty) {
		return Value[U]{}, &Error{
			Reason: ErrTypeMismatch,
			Op:     "llvm.Value.As",
			Msg:    "value kind mismatch: have " + kindName(kindOfType(ty)) + ", want " + kindName(kindOf[U]()),
		}
	}
	return Value[U]{ref: ref, ty: v.ty, ctx: v.ctx, life: v.life}, nil
}

// MustAs As 的 panic 版本（程序员错误）
func (v Value[T]) MustAs[U Kind]() Value[U] {
	res, err := v.As[U]()
	if err != nil {
		panic(err)
	}
	return res
}

// newValue 内部构造值并预取类型句柄（仅调试层；release 构建零额外开销）
func newValue[T Kind](ctx *Context, life *Lifetime, ref binding.LLVMValueRef) Value[T] {
	v := Value[T]{ref: ref, ctx: ctx, life: life}
	if checks.Debug && !ref.IsNil() {
		v.ty = binding.LLVMTypeOf(ref)
	}
	return v
}

// NewValue 由底层句柄构建值（供 llvm/* 子包桥接使用）
func NewValue[T Kind](ctx *Context, life *Lifetime, ref binding.LLVMValueRef) Value[T] {
	return newValue[T](ctx, life, ref)
}

// ValueOf 由底层句柄构建擦除值（指令遍历等动态来源）
func ValueOf(ctx *Context, life *Lifetime, ref binding.LLVMValueRef) Value[DynT] {
	return newValue[DynT](ctx, life, ref)
}
