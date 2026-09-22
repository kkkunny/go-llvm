package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Call 调用指令角色（内嵌 Value[T]，自动实现 llvm.ValueRef/AnyValue）
type Call[T llvm.Kind] struct {
	llvm.Value[T]
}

// ArgCount 实参个数
func (c Call[T]) ArgCount() uint32 { return binding.LLVMGetNumArgOperands(c.Ref()) }

// Arg 第 i 个实参（擦除种类）
func (c Call[T]) Arg(i uint32) llvm.Value[llvm.DynT] {
	ref := binding.LLVMGetOperand(c.Ref(), i)
	return wrapDyn(c.Context(), c.Lifetime(), ref)
}

// SetArg 替换第 i 个实参
func (c Call[T]) SetArg(i uint32, v llvm.AnyValue) {
	if i >= c.ArgCount() {
		errPanic(llvm.ErrInvalidArg, "ir.Call.SetArg", "argument index %d out of range", i)
	}
	binding.LLVMSetOperand(c.Ref(), i, v.Ref())
}

// CalledFunction 被调用函数（非间接调用时）
func (c Call[T]) CalledFunction() (llvm.Value[llvm.FnT], bool) {
	ref := binding.LLVMGetCalledValue(c.Ref())
	if ref.IsNil() {
		return llvm.Value[llvm.FnT]{}, false
	}
	return wrapValue[llvm.FnT](c.Context(), c.Lifetime(), ref), true
}

// Phi PHI 节点角色（内嵌 Value[T]）
type Phi[T llvm.Kind] struct {
	llvm.Value[T]
}

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
		if !in.Value.Type().Equal(p.Type()) {
			errPanic(llvm.ErrTypeMismatch, "ir.Phi.AddIncoming", "incoming type %s differs from phi type %s", in.Value.Type(), p.Type())
		}
		values[i] = in.Value.Ref()
		blocks[i] = in.Block.ref
	}
	binding.LLVMAddIncoming(p.Ref(), values, blocks)
}

// Count PHI 输入条数
func (p Phi[T]) Count() uint32 { return binding.LLVMCountIncoming(p.Ref()) }

// IncomingAt 第 i 条 PHI 输入
func (p Phi[T]) IncomingAt(i uint32) Incoming[T] {
	val := binding.LLVMGetIncomingValue(p.Ref(), i)
	blk := binding.LLVMGetIncomingBlock(p.Ref(), i)
	return Incoming[T]{
		Value: wrapValue[T](p.Context(), p.Lifetime(), val),
		Block: wrapBlock(p.Context(), p.Lifetime(), blk),
	}
}
