package ir

import (
	"iter"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
)

// OperandCount 指令/常量的操作数个数（call 的被调方是最后一个操作数，见 LLVM-C 约定）
func OperandCount(inst llvm.AnyValue) uint32 {
	const op = "ir.OperandCount"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	return uint32(binding.LLVMGetNumOperands(inst.Ref()))
}

// OperandAt 第 i 个操作数；越界校验仅调试层
func OperandAt(inst llvm.AnyValue, i uint32) llvm.Value[llvm.DynT] {
	const op = "ir.OperandAt"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	if checks.Debug && i >= OperandCount(inst) {
		llvm.Panicf(llvm.ErrInvalidArg, op, "operand index %d out of range", i)
	}
	return llvm.ValueOf(inst.Context(), inst.Lifetime(), binding.LLVMGetOperand(inst.Ref(), i))
}

// SetOperand 替换第 i 个操作数；越界/跨上下文校验按三层约定
func SetOperand(inst llvm.AnyValue, i uint32, v llvm.AnyValue) {
	const op = "ir.SetOperand"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	if v == nil || !v.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead operand")
	}
	if inst.Context() != v.Context() {
		llvm.Panicf(llvm.ErrCrossContext, op, "operand belongs to another context")
	}
	if checks.Debug && i >= OperandCount(inst) {
		llvm.Panicf(llvm.ErrInvalidArg, op, "operand index %d out of range", i)
	}
	binding.LLVMSetOperand(inst.Ref(), i, v.Ref())
}

// Operands 操作数的惰性遍历
func Operands(inst llvm.AnyValue) iter.Seq[llvm.Value[llvm.DynT]] {
	return func(yield func(llvm.Value[llvm.DynT]) bool) {
		n := OperandCount(inst)
		for i := uint32(0); i < n; i++ {
			if !yield(OperandAt(inst, i)) {
				return
			}
		}
	}
}

// Use 值的一条使用记录（use-list 节点）。
// 仅提供读取：LLVM-C 未暴露单条 use 的替换（`Use::set`），定向替换请用 SetOperand。
type Use struct {
	ref  binding.LLVMUseRef
	ctx  *llvm.Context
	life *llvm.Lifetime
}

// User 使用该值的指令/常量
func (u Use) User() llvm.Value[llvm.DynT] {
	return llvm.ValueOf(u.ctx, u.life, binding.LLVMGetUser(u.ref))
}

// UsedValue 该使用记录指向的值
func (u Use) UsedValue() llvm.Value[llvm.DynT] {
	return llvm.ValueOf(u.ctx, u.life, binding.LLVMGetUsedValue(u.ref))
}

// Uses 值的使用记录遍历（谁在用我）
func Uses(v llvm.AnyValue) iter.Seq[Use] {
	return func(yield func(Use) bool) {
		const op = "ir.Uses"
		if v == nil || !v.Alive() {
			llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
		}
		for ref := binding.LLVMGetFirstUse(v.Ref()); !ref.IsNil(); ref = binding.LLVMGetNextUse(ref) {
			if !yield(Use{ref: ref, ctx: v.Context(), life: v.Lifetime()}) {
				return
			}
		}
	}
}

// ReplaceAllUses 把 old 的全部使用替换为 new（RAUW）；两者须同上下文
func ReplaceAllUses(old, new llvm.AnyValue) {
	const op = "ir.ReplaceAllUses"
	if old == nil || !old.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	if new == nil || !new.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead replacement")
	}
	if old.Context() != new.Context() {
		llvm.Panicf(llvm.ErrCrossContext, op, "replacement belongs to another context")
	}
	binding.LLVMReplaceAllUsesWith(old.Ref(), new.Ref())
}
