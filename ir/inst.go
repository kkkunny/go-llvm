package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
)

// Call 调用指令角色（内嵌 Value[T]，自动实现 llvm.ValueRef/AnyValue）
type Call[T llvm.Kind] struct {
	llvm.Value[T]
}

// ArgCount 实参个数
func (c Call[T]) ArgCount() uint32 {
	return binding.LLVMGetNumArgOperands(c.Ref())
}

// Arg 第 i 个实参（擦除种类）；越界校验仅调试层
func (c Call[T]) Arg(i uint32) llvm.Value[llvm.DynT] {
	const op = "ir.Call.Arg"
	if checks.Debug && i >= c.ArgCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "argument index %d out of range", i)
	}
	ref := binding.LLVMGetOperand(c.Ref(), i)
	return llvm.ValueOf(c.Context(), c.Lifetime(), ref)
}

// SetArg 替换第 i 个实参；越界校验仅调试层
func (c Call[T]) SetArg(i uint32, v llvm.AnyValue) {
	const op = "ir.Call.SetArg"
	if checks.Debug && i >= c.ArgCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "argument index %d out of range", i)
	}
	c.Context().CheckValues(op, v)
	binding.LLVMSetOperand(c.Ref(), i, v.Ref())
}

// CalledFunction 被调用函数（非间接调用时）
func (c Call[T]) CalledFunction() (llvm.Value[llvm.FnT], bool) {
	ref := binding.LLVMGetCalledValue(c.Ref())
	if ref.IsNil() {
		return llvm.Value[llvm.FnT]{}, false
	}
	return llvm.NewValue[llvm.FnT](c.Context(), c.Lifetime(), ref), true
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

// AddIncoming 追加 PHI 输入；值类型一致为语义契约，仅调试层
func (p Phi[T]) AddIncoming(incomings ...Incoming[T]) {
	const op = "ir.Phi.AddIncoming"
	values := make([]binding.LLVMValueRef, len(incomings))
	blocks := make([]binding.LLVMBasicBlockRef, len(incomings))
	for i, in := range incomings {
		in.Block.Check(op)
		if checks.Debug {
			phiTy := typeOfVal(p.Ref(), p.RawType())
			inTy := typeOfVal(in.Value.Ref(), in.Value.RawType())
			if !inTy.Equal(phiTy) {
				llvm.Panicf(llvm.ErrTypeMismatch, op, "incoming type %s differs from phi type %s",
					typeString(p.Context(), in.Value.Ref()), typeString(p.Context(), p.Ref()))
			}
		}
		values[i] = in.Value.Ref()
		blocks[i] = in.Block.ref
	}
	binding.LLVMAddIncoming(p.Ref(), values, blocks)
}

// Count PHI 输入条数
func (p Phi[T]) Count() uint32 {
	return binding.LLVMCountIncoming(p.Ref())
}

// IncomingAt 第 i 条 PHI 输入；越界校验仅调试层
func (p Phi[T]) IncomingAt(i uint32) Incoming[T] {
	const op = "ir.Phi.IncomingAt"
	if checks.Debug && i >= p.Count() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "incoming index %d out of range", i)
	}
	val := binding.LLVMGetIncomingValue(p.Ref(), i)
	blk := binding.LLVMGetIncomingBlock(p.Ref(), i)
	return Incoming[T]{
		Value: llvm.NewValue[T](p.Context(), p.Lifetime(), val),
		Block: wrapBlock(p.Context(), p.Lifetime(), blk),
	}
}
