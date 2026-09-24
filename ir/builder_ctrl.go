package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
)

// Switch switch 终结指令角色（内嵌 Value[IntT]）
type Switch struct {
	llvm.Value[llvm.IntT]
}

// CondType switch 条件操作数的类型
func (s Switch) CondType() llvm.Type[llvm.IntT] {
	s.Check("ir.Switch.CondType")
	return llvm.NewType[llvm.IntT](s.Context(), binding.LLVMTypeOf(binding.LLVMGetOperand(s.Ref(), 0)))
}

// AddCase 追加 case；条件类型不符 panic（语义契约，仅调试层）
func (s Switch) AddCase(cond llvm.ValueRef[llvm.IntT], blk Block) {
	const op = "ir.Switch.AddCase"
	s.Check(op)
	cv := cond.AsValue()
	cv.Check(op)
	blk.Check(op)
	if checks.Debug && !cv.Type().Equal(s.CondType()) {
		llvm.Panicf(llvm.ErrTypeMismatch, op, "case type %s differs from switch type %s", cv.Type(), s.CondType())
	}
	binding.LLVMAddCase(s.Ref(), cv.Ref(), blk.ref)
}

// Count case 个数（后继数减去默认分支）
func (s Switch) Count() uint32 {
	s.Check("ir.Switch.Count")
	n := binding.LLVMGetNumSuccessors(s.Ref())
	if n == 0 {
		return 0
	}
	return n - 1
}

// DefaultBlock 默认分支
func (s Switch) DefaultBlock() Block {
	s.Check("ir.Switch.DefaultBlock")
	ref := binding.LLVMGetSwitchDefaultDest(s.Ref())
	return wrapBlock(s.Context(), s.Lifetime(), ref)
}

// CaseBlock 第 i 个 case 的目标块（i 从 0 开始）；越界校验仅调试层
func (s Switch) CaseBlock(i uint32) Block {
	const op = "ir.Switch.CaseBlock"
	s.Check(op)
	if checks.Debug && i >= s.Count() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "case index %d out of range", i)
	}
	ref := binding.LLVMGetSuccessor(s.Ref(), i+1)
	return wrapBlock(s.Context(), s.Lifetime(), ref)
}

// CaseValue 第 i 个 case 的常量（i 从 0 开始）；越界校验仅调试层
func (s Switch) CaseValue(i uint32) llvm.Value[llvm.DynT] {
	const op = "ir.Switch.CaseValue"
	s.Check(op)
	if checks.Debug && i >= s.Count() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "case index %d out of range", i)
	}
	ref := binding.LLVMGetSwitchCaseValue(s.Ref(), i+1)
	return llvm.ValueOf(s.Context(), s.Lifetime(), ref)
}

// ===== Builder 终结指令 =====

// RetVoid 插入 ret void
func (b *Builder) RetVoid() llvm.Value[llvm.VoidT] {
	const op = "ir.Builder.RetVoid"
	b.pre(op)
	ref := binding.LLVMBuildRetVoid(b.ref)
	return llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)
}

// Ret 插入 ret v
func (b *Builder) Ret(v llvm.AnyValue) llvm.Value[llvm.VoidT] {
	const op = "ir.Builder.Ret"
	b.pre(op, coreAny(v))
	ref := binding.LLVMBuildRet(b.ref, v.Ref())
	return llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)
}

// Br 插入无条件跳转
func (b *Builder) Br(blk Block) llvm.Value[llvm.VoidT] {
	const op = "ir.Builder.Br"
	b.pre(op)
	b.preBlock(op, blk)
	ref := binding.LLVMBuildBr(b.ref, blk.ref)
	return llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)
}

// CondBr 插入条件跳转
func (b *Builder) CondBr(cond llvm.ValueRef[llvm.IntT], then, els Block) llvm.Value[llvm.VoidT] {
	const op = "ir.Builder.CondBr"
	cv := cond.AsValue()
	b.pre(op, core(cv))
	b.preBlock(op, then)
	b.preBlock(op, els)
	if checks.Debug {
		if bits := llvm.AsIntType(cv.Type()).Bits(); bits != 1 {
			llvm.Panicf(llvm.ErrTypeMismatch, op, "condition must be i1, got i%d", bits)
		}
	}
	ref := binding.LLVMBuildCondBr(b.ref, cv.Ref(), then.ref, els.ref)
	return llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)
}

// Switch 插入 switch 终结指令
func (b *Builder) Switch(v llvm.ValueRef[llvm.IntT], def Block) Switch {
	const op = "ir.Builder.Switch"
	vv := v.AsValue()
	b.pre(op, core(vv))
	b.preBlock(op, def)
	ref := binding.LLVMBuildSwitch(b.ref, vv.Ref(), def.ref, 0)
	return Switch{Value: llvm.NewValue[llvm.IntT](b.ctx, b.inserted.life, ref)}
}

// Unreachable 插入 unreachable
func (b *Builder) Unreachable() llvm.Value[llvm.VoidT] {
	const op = "ir.Builder.Unreachable"
	b.pre(op)
	ref := binding.LLVMBuildUnreachable(b.ref)
	return llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)
}
