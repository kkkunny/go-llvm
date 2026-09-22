package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Call 调用指令角色
type Call[T llvm.Kind] struct {
	v llvm.Value[T]
}

// Value 返回底层泛型值
func (c Call[T]) Value() llvm.Value[T] { return c.v }

// AsValue 实现 llvm.ValueRef[T]
func (c Call[T]) AsValue() llvm.Value[T] { return c.v }

// ArgCount 实参个数
func (c Call[T]) ArgCount() uint32 { return binding.LLVMGetNumArgOperands(c.v.Ref()) }

// Arg 第 i 个实参（擦除种类）
func (c Call[T]) Arg(i uint32) llvm.Value[llvm.DynT] {
	ref := binding.LLVMGetOperand(c.v.Ref(), i)
	return llvm.ValueOf(c.v.Context(), c.v.Lifetime(), ref)
}

// SetArg 替换第 i 个实参
func (c Call[T]) SetArg(i uint32, v llvm.AnyValue) {
	if i >= c.ArgCount() {
		errPanic(llvm.ErrInvalidArg, "ir.Call.SetArg", "argument index %d out of range", i)
	}
	binding.LLVMSetOperand(c.v.Ref(), i, v.Ref())
}

// CalledFunction 被调用函数（非间接调用时）
func (c Call[T]) CalledFunction() (llvm.Value[llvm.FnT], bool) {
	ref := binding.LLVMGetCalledValue(c.v.Ref())
	if ref.IsNil() {
		return llvm.Value[llvm.FnT]{}, false
	}
	return llvm.NewValue[llvm.FnT](c.v.Context(), c.v.Lifetime(), ref), true
}

// Phi PHI 节点角色
type Phi[T llvm.Kind] struct {
	v llvm.Value[T]
}

// Value 返回底层泛型值
func (p Phi[T]) Value() llvm.Value[T] { return p.v }

// AsValue 实现 llvm.ValueRef[T]
func (p Phi[T]) AsValue() llvm.Value[T] { return p.v }

// Incoming 一条 PHI 输入
type Incoming[T llvm.Kind] struct {
	Value llvm.Value[T]
	Block Block
}

// AddIncoming 追加 PHI 输入
func (p Phi[T]) AddIncoming(incomings ...Incoming[T]) {
	values := make([]binding.LLVMValueRef, len(incomings))
	blocks := make([]binding.LLVMBasicBlockRef, len(incomings))
	for i, in := range incomings {
		if in.Block.ref.IsNil() {
			errPanic(llvm.ErrInvalidArg, "ir.Phi.AddIncoming", "nil incoming block")
		}
		if !in.Value.Type().Equal(p.v.Type()) {
			errPanic(llvm.ErrTypeMismatch, "ir.Phi.AddIncoming", "incoming type %s differs from phi type %s", in.Value.Type(), p.v.Type())
		}
		values[i] = in.Value.Ref()
		blocks[i] = in.Block.ref
	}
	binding.LLVMAddIncoming(p.v.Ref(), values, blocks)
}

// Count PHI 输入条数
func (p Phi[T]) Count() uint32 { return binding.LLVMCountIncoming(p.v.Ref()) }

// IncomingAt 第 i 条 PHI 输入
func (p Phi[T]) IncomingAt(i uint32) Incoming[T] {
	val := binding.LLVMGetIncomingValue(p.v.Ref(), i)
	blk := binding.LLVMGetIncomingBlock(p.v.Ref(), i)
	return Incoming[T]{
		Value: llvm.NewValue[T](p.v.Context(), p.v.Lifetime(), val),
		Block: Block{ref: blk, ctx: p.v.Context(), life: p.v.Lifetime()},
	}
}
