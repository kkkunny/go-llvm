package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Switch switch 终结指令角色
type Switch struct {
	v llvm.Value[llvm.IntT]
}

// Value 返回底层泛型值
func (s Switch) Value() llvm.Value[llvm.IntT] { return s.v }

// AsValue 实现 llvm.ValueRef[IntT]
func (s Switch) AsValue() llvm.Value[llvm.IntT] { return s.v }

// CondType switch 条件操作数的类型
func (s Switch) CondType() llvm.Type[llvm.IntT] {
	return llvm.NewType[llvm.IntT](s.v.Context(), binding.LLVMTypeOf(binding.LLVMGetOperand(s.v.Ref(), 0)))
}

// AddCase 追加 case；条件类型不符 panic
func (s Switch) AddCase(cond llvm.ValueRef[llvm.IntT], blk Block) {
	cv := cond.AsValue()
	if !cv.Type().Equal(s.CondType()) {
		errPanic(llvm.ErrTypeMismatch, "ir.Switch.AddCase", "case type %s differs from switch type %s", cv.Type(), s.CondType())
	}
	if blk.ref.IsNil() {
		errPanic(llvm.ErrInvalidArg, "ir.Switch.AddCase", "nil block")
	}
	binding.LLVMAddCase(s.v.Ref(), cv.Ref(), blk.ref)
}

// Count case 个数（后继数减去默认分支）
func (s Switch) Count() uint32 {
	n := binding.LLVMGetNumSuccessors(s.v.Ref())
	if n == 0 {
		return 0
	}
	return n - 1
}

// DefaultBlock 默认分支
func (s Switch) DefaultBlock() Block {
	ref := binding.LLVMGetSwitchDefaultDest(s.v.Ref())
	return Block{ref: ref, ctx: s.v.Context(), life: s.v.Lifetime()}
}

// CaseBlock 第 i 个 case 的目标块（i 从 0 开始）
func (s Switch) CaseBlock(i uint32) Block {
	ref := binding.LLVMGetSuccessor(s.v.Ref(), i+1)
	return Block{ref: ref, ctx: s.v.Context(), life: s.v.Lifetime()}
}

// CaseValue 第 i 个 case 的常量（i 从 0 开始）
func (s Switch) CaseValue(i uint32) llvm.Value[llvm.DynT] {
	ref := binding.LLVMGetSwitchCaseValue(s.v.Ref(), i+1)
	return llvm.ValueOf(s.v.Context(), s.v.Lifetime(), ref)
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
	b.pre(op, v)
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
	b.pre(op, cv.Dyn())
	b.preBlock(op, then)
	b.preBlock(op, els)
	if bits := llvm.AsIntType(cv.Type()).Bits(); bits != 1 {
		errPanic(llvm.ErrTypeMismatch, op, "condition must be i1, got i%d", bits)
	}
	ref := binding.LLVMBuildCondBr(b.ref, cv.Ref(), then.ref, els.ref)
	return llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)
}

// Switch 插入 switch 终结指令
func (b *Builder) Switch(v llvm.ValueRef[llvm.IntT], def Block) Switch {
	const op = "ir.Builder.Switch"
	vv := v.AsValue()
	b.pre(op, vv.Dyn())
	b.preBlock(op, def)
	ref := binding.LLVMBuildSwitch(b.ref, vv.Ref(), def.ref, 0)
	return Switch{v: llvm.NewValue[llvm.IntT](b.ctx, b.inserted.life, ref)}
}

// Unreachable 插入 unreachable
func (b *Builder) Unreachable() llvm.Value[llvm.VoidT] {
	const op = "ir.Builder.Unreachable"
	b.pre(op)
	ref := binding.LLVMBuildUnreachable(b.ref)
	return llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)
}
