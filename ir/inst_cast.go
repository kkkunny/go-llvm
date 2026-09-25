package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
)

// asInst 按操作码把任意值铸成指令角色：opcode 不符返回 false（不 panic）。
// T 为角色的结果种类，调试层与实际结果类型比对（语义契约）。
func asInst[T llvm.Kind, R any](op string, v llvm.AnyValue, want Op, wrap func(llvm.Value[T]) R) (R, bool) {
	var zero R
	got, ok := OpOf(v)
	if !ok || got != want {
		return zero, false
	}
	vv := llvm.NewValue[T](v.Context(), v.Lifetime(), v.Ref())
	if checks.Debug {
		checkKind[T](op, v.Context(), binding.LLVMTypeOf(v.Ref()))
	}
	return wrap(vv), true
}

// AsLoad 转为 Load 角色（T 为加载结果种类）
func AsLoad[T llvm.Kind](v llvm.AnyValue) (Load[T], bool) {
	return asInst[T, Load[T]]("ir.AsLoad", v, OpLoad, func(x llvm.Value[T]) Load[T] { return Load[T]{x} })
}

// AsStore 转为 Store 角色
func AsStore(v llvm.AnyValue) (Store, bool) {
	return asInst[llvm.VoidT, Store]("ir.AsStore", v, OpStore, func(x llvm.Value[llvm.VoidT]) Store { return Store{x} })
}

// AsAlloca 转为 Alloca 角色
func AsAlloca(v llvm.AnyValue) (Alloca, bool) {
	return asInst[llvm.PtrT, Alloca]("ir.AsAlloca", v, OpAlloca, func(x llvm.Value[llvm.PtrT]) Alloca { return Alloca{x} })
}

// AsCall 转为 Call 角色（T 为返回种类）
func AsCall[T llvm.Kind](v llvm.AnyValue) (Call[T], bool) {
	return asInst[T, Call[T]]("ir.AsCall", v, OpCall, func(x llvm.Value[T]) Call[T] { return Call[T]{x} })
}

// AsInvoke 转为 Invoke 角色（T 为返回种类）
func AsInvoke[T llvm.Kind](v llvm.AnyValue) (Invoke[T], bool) {
	return asInst[T, Invoke[T]]("ir.AsInvoke", v, OpInvoke, func(x llvm.Value[T]) Invoke[T] { return Invoke[T]{x} })
}

// AsPhi 转为 Phi 角色（T 为结果种类）
func AsPhi[T llvm.Kind](v llvm.AnyValue) (Phi[T], bool) {
	return asInst[T, Phi[T]]("ir.AsPhi", v, OpPHI, func(x llvm.Value[T]) Phi[T] { return Phi[T]{x} })
}

// AsSwitch 转为 Switch 角色（switch 指令底层类型为 void）
func AsSwitch(v llvm.AnyValue) (Switch, bool) {
	return asInst[llvm.VoidT, Switch]("ir.AsSwitch", v, OpSwitch, func(x llvm.Value[llvm.VoidT]) Switch { return Switch{x} })
}

// AsAtomicRMW 转为 AtomicRMW 角色（T 为操作数种类）
func AsAtomicRMW[T llvm.Kind](v llvm.AnyValue) (AtomicRMW[T], bool) {
	return asInst[T, AtomicRMW[T]]("ir.AsAtomicRMW", v, OpAtomicRMW, func(x llvm.Value[T]) AtomicRMW[T] { return AtomicRMW[T]{x} })
}

// AsCmpXchg 转为 CmpXchg 角色
func AsCmpXchg(v llvm.AnyValue) (CmpXchg, bool) {
	return asInst[llvm.StructT, CmpXchg]("ir.AsCmpXchg", v, OpAtomicCmpXchg, func(x llvm.Value[llvm.StructT]) CmpXchg { return CmpXchg{x} })
}

// AsFence 转为 Fence 角色
func AsFence(v llvm.AnyValue) (Fence, bool) {
	return asInst[llvm.VoidT, Fence]("ir.AsFence", v, OpFence, func(x llvm.Value[llvm.VoidT]) Fence { return Fence{x} })
}

// AsLandingPad 转为 LandingPad 角色（T 为结果种类）
func AsLandingPad[T llvm.Kind](v llvm.AnyValue) (LandingPad[T], bool) {
	return asInst[T, LandingPad[T]]("ir.AsLandingPad", v, OpLandingPad, func(x llvm.Value[T]) LandingPad[T] { return LandingPad[T]{x} })
}

// AsFuncletPad 转为 FuncletPad 角色（catchpad/cleanuppad 共用）
func AsFuncletPad(v llvm.AnyValue) (FuncletPad, bool) {
	got, ok := OpOf(v)
	if !ok || (got != OpCatchPad && got != OpCleanupPad) {
		return FuncletPad{}, false
	}
	if checks.Debug {
		checkKind[llvm.TokenT]("ir.AsFuncletPad", v.Context(), binding.LLVMTypeOf(v.Ref()))
	}
	return FuncletPad{llvm.NewValue[llvm.TokenT](v.Context(), v.Lifetime(), v.Ref())}, true
}

// AsCatchSwitch 转为 CatchSwitch 角色
func AsCatchSwitch(v llvm.AnyValue) (CatchSwitch, bool) {
	return asInst[llvm.TokenT, CatchSwitch]("ir.AsCatchSwitch", v, OpCatchSwitch, func(x llvm.Value[llvm.TokenT]) CatchSwitch { return CatchSwitch{x} })
}

// ===== 终结指令操作 =====

// requireTerminator 校验值是存活的终结指令（崩溃类地板，两种构建均生效）：
// nil/已释放或非终结指令一律 panic ErrInvalidArg，避免 LLVM-C 对非终结指令的 UB（挂起/SIGSEGV）。
func requireTerminator(op string, term llvm.AnyValue) {
	if term == nil || !term.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead terminator")
	}
	if binding.LLVMIsATerminatorInst(term.Ref()).IsNil() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "not a terminator instruction")
	}
}

// SuccessorCount 终结指令的后继个数；非终结指令 panic（崩溃类地板，两种构建均生效）
func SuccessorCount(term llvm.AnyValue) uint32 {
	const op = "ir.SuccessorCount"
	requireTerminator(op, term)
	return binding.LLVMGetNumSuccessors(term.Ref())
}

// Successor 第 i 个后继块；越界校验仅调试层
func Successor(term llvm.AnyValue, i uint32) Block {
	const op = "ir.Successor"
	requireTerminator(op, term)
	if checks.Debug && i >= SuccessorCount(term) {
		llvm.Panicf(llvm.ErrInvalidArg, op, "successor index %d out of range", i)
	}
	return wrapBlock(term.Context(), term.Lifetime(), binding.LLVMGetSuccessor(term.Ref(), i))
}

// SetSuccessor 替换第 i 个后继块
func SetSuccessor(term llvm.AnyValue, i uint32, blk Block) {
	const op = "ir.SetSuccessor"
	requireTerminator(op, term)
	blk.Check(op)
	if term.Context() != blk.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "block belongs to another context")
	}
	if checks.Debug && i >= SuccessorCount(term) {
		llvm.Panicf(llvm.ErrInvalidArg, op, "successor index %d out of range", i)
	}
	binding.LLVMSetSuccessor(term.Ref(), i, blk.Ref())
}

// IsConditional 终结指令是否有条件（条件 br / switch）；非终结指令 panic（崩溃类地板，两种构建均生效）
func IsConditional(term llvm.AnyValue) bool {
	const op = "ir.IsConditional"
	requireTerminator(op, term)
	return binding.LLVMIsConditional(term.Ref())
}

// Condition 条件值（非条件终结指令 panic，仅调试层）
func Condition(term llvm.AnyValue) llvm.Value[llvm.DynT] {
	const op = "ir.Condition"
	requireTerminator(op, term)
	if checks.Debug && !IsConditional(term) {
		llvm.Panicf(llvm.ErrInvalidArg, op, "terminator is not conditional")
	}
	return llvm.ValueOf(term.Context(), term.Lifetime(), binding.LLVMGetCondition(term.Ref()))
}

// SetCondition 替换条件值
func SetCondition(term llvm.AnyValue, cond llvm.AnyValue) {
	const op = "ir.SetCondition"
	requireTerminator(op, term)
	if cond == nil || !cond.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead condition")
	}
	if term.Context() != cond.Context() {
		llvm.Panicf(llvm.ErrCrossContext, op, "condition belongs to another context")
	}
	if checks.Debug && !IsConditional(term) {
		llvm.Panicf(llvm.ErrInvalidArg, op, "terminator is not conditional")
	}
	binding.LLVMSetCondition(term.Ref(), cond.Ref())
}
