package ir

import (
	"iter"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
	"github.com/kkkunny/go-llvm/internal/errs"
)

// OperandCount 指令/常量的操作数个数（call 的被调方是最后一个操作数，见 LLVM-C 约定）
func OperandCount(inst llvm.AnyValue) uint32 {
	const op = "ir.OperandCount"
	if inst == nil || !inst.Alive() {
		errs.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	return uint32(binding.LLVMGetNumOperands(inst.Ref()))
}

// OperandAt 第 i 个操作数；越界校验仅调试层
func OperandAt(inst llvm.AnyValue, i uint32) llvm.Value[llvm.DynT] {
	const op = "ir.OperandAt"
	if inst == nil || !inst.Alive() {
		errs.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	if checks.Debug && i >= OperandCount(inst) {
		errs.Panicf(llvm.ErrInvalidArg, op, "operand index %d out of range", i)
	}
	return llvm.ValueOf(inst.Context(), inst.Lifetime(), binding.LLVMGetOperand(inst.Ref(), i))
}

// SetOperand 替换第 i 个操作数；越界/跨上下文校验按三层约定，类型一致仅调试层
func SetOperand(inst llvm.AnyValue, i uint32, v llvm.AnyValue) {
	const op = "ir.SetOperand"
	if inst == nil || !inst.Alive() {
		errs.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	if v == nil || !v.Alive() {
		errs.Panicf(llvm.ErrInvalidArg, op, "nil or dead operand")
	}
	if inst.Context() != v.Context() {
		errs.Panicf(llvm.ErrCrossContext, op, "operand belongs to another context")
	}
	if checks.Debug {
		if i >= OperandCount(inst) {
			errs.Panicf(llvm.ErrInvalidArg, op, "operand index %d out of range", i)
		}
		want := binding.LLVMTypeOf(binding.LLVMGetOperand(inst.Ref(), i))
		if got := binding.LLVMTypeOf(v.Ref()); !want.Equal(got) {
			errs.Panicf(llvm.ErrTypeMismatch, op, "operand type %s does not match replaced operand type %s",
				llvm.TypeOfRef(inst.Context(), got), llvm.TypeOfRef(inst.Context(), want))
		}
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
	u.check("ir.Use.User")
	return llvm.ValueOf(u.ctx, u.life, binding.LLVMGetUser(u.ref))
}

// UsedValue 该使用记录指向的值
func (u Use) UsedValue() llvm.Value[llvm.DynT] {
	u.check("ir.Use.UsedValue")
	return llvm.ValueOf(u.ctx, u.life, binding.LLVMGetUsedValue(u.ref))
}

// check 前置校验（崩溃类地板：句柄/上下文/生命周期）
func (u Use) check(op string) {
	if u.ref.IsNil() {
		errs.Panicf(llvm.ErrInvalidArg, op, "nil use handle")
	}
	if u.ctx == nil || !u.ctx.Alive() || u.life == nil || !u.life.Alive() {
		errs.Panicf(llvm.ErrUseAfterFree, op, "use handle is dead")
	}
}

// Uses 值的使用记录遍历（谁在用我）
func Uses(v llvm.AnyValue) iter.Seq[Use] {
	return func(yield func(Use) bool) {
		const op = "ir.Uses"
		if v == nil || !v.Alive() {
			errs.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
		}
		for ref := binding.LLVMGetFirstUse(v.Ref()); !ref.IsNil(); ref = binding.LLVMGetNextUse(ref) {
			if !yield(Use{ref: ref, ctx: v.Context(), life: v.Lifetime()}) {
				return
			}
		}
	}
}

// ReplaceAllUses 把 old 的全部使用替换为 new（RAUW）；两者须同上下文，类型一致仅调试层
// （上游 Value::doRAUW 对类型不一致会断言）
func ReplaceAllUses(old, new llvm.AnyValue) {
	const op = "ir.ReplaceAllUses"
	if old == nil || !old.Alive() {
		errs.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	if new == nil || !new.Alive() {
		errs.Panicf(llvm.ErrInvalidArg, op, "nil or dead replacement")
	}
	if old.Context() != new.Context() {
		errs.Panicf(llvm.ErrCrossContext, op, "replacement belongs to another context")
	}
	if checks.Debug {
		oldTy := binding.LLVMTypeOf(old.Ref())
		newTy := binding.LLVMTypeOf(new.Ref())
		if !oldTy.Equal(newTy) {
			errs.Panicf(llvm.ErrTypeMismatch, op, "replacement type %s differs from original type %s",
				llvm.TypeOfRef(old.Context(), newTy), llvm.TypeOfRef(old.Context(), oldTy))
		}
	}
	binding.LLVMReplaceAllUsesWith(old.Ref(), new.Ref())
}
