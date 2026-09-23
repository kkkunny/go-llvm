package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Invoke invoke 指令角色（内嵌 Value[T]）；实参操作同 Call（LLVM 侧同为 CallBase），
// 另有正常/异常出口块访问
type Invoke[T llvm.Kind] struct {
	llvm.Value[T]
}

// NormalBlock 正常出口块
func (c Invoke[T]) NormalBlock() Block {
	c.Check("ir.Invoke.NormalBlock")
	return wrapBlock(c.Context(), c.Lifetime(), binding.LLVMGetSuccessor(c.Ref(), 0))
}

// UnwindBlock 异常出口块
func (c Invoke[T]) UnwindBlock() Block {
	c.Check("ir.Invoke.UnwindBlock")
	return wrapBlock(c.Context(), c.Lifetime(), binding.LLVMGetSuccessor(c.Ref(), 1))
}

// ArgCount 实参个数
func (c Invoke[T]) ArgCount() uint32 {
	c.Check("ir.Invoke.ArgCount")
	return binding.LLVMGetNumArgOperands(c.Ref())
}

// Arg 第 i 个实参（擦除种类）
func (c Invoke[T]) Arg(i uint32) llvm.Value[llvm.DynT] {
	const op = "ir.Invoke.Arg"
	c.Check(op)
	if i >= c.ArgCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "argument index %d out of range", i)
	}
	return llvm.ValueOf(c.Context(), c.Lifetime(), binding.LLVMGetOperand(c.Ref(), i))
}

// SetArg 替换第 i 个实参
func (c Invoke[T]) SetArg(i uint32, v llvm.AnyValue) {
	const op = "ir.Invoke.SetArg"
	c.Check(op)
	if i >= c.ArgCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "argument index %d out of range", i)
	}
	c.Context().CheckValues(op, v)
	binding.LLVMSetOperand(c.Ref(), i, v.Ref())
}

// CalledFunction 被调用函数（非间接 invoke 时）
func (c Invoke[T]) CalledFunction() (llvm.Value[llvm.FnT], bool) {
	c.Check("ir.Invoke.CalledFunction")
	ref := binding.LLVMGetCalledValue(c.Ref())
	if ref.IsNil() {
		return llvm.Value[llvm.FnT]{}, false
	}
	return llvm.NewValue[llvm.FnT](c.Context(), c.Lifetime(), ref), true
}

// LandingPad landingpad 指令角色（内嵌 Value[T]，T 通常为 StructT）
type LandingPad[T llvm.Kind] struct {
	llvm.Value[T]
}

// AddClause 追加 catch/filter 子句（catch：类型信息全局；filter：常量数组）。
// 注：经 Dyn() 走 Value[DynT].IsConstant，避免 Global 角色遮蔽语义（全局常量标志）
func (l LandingPad[T]) AddClause(v llvm.AnyValue) {
	const op = "ir.LandingPad.AddClause"
	l.Check(op)
	l.Context().CheckValues(op, v)
	if !v.Dyn().IsConstant() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "clause must be a constant")
	}
	binding.LLVMAddClause(l.Ref(), v.Ref())
}

// ClauseCount 子句数量
func (l LandingPad[T]) ClauseCount() uint32 {
	l.Check("ir.LandingPad.ClauseCount")
	return binding.LLVMGetNumClauses(l.Ref())
}

// Clause 第 i 条子句（擦除种类）
func (l LandingPad[T]) Clause(i uint32) llvm.Value[llvm.DynT] {
	const op = "ir.LandingPad.Clause"
	l.Check(op)
	if i >= l.ClauseCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "clause index %d out of range", i)
	}
	return llvm.ValueOf(l.Context(), l.Lifetime(), binding.LLVMGetClause(l.Ref(), i))
}

// SetCleanup 设置 cleanup 标志
func (l LandingPad[T]) SetCleanup(v bool) {
	l.Check("ir.LandingPad.SetCleanup")
	binding.LLVMSetCleanup(l.Ref(), v)
}

// IsCleanup 是否 cleanup
func (l LandingPad[T]) IsCleanup() bool {
	l.Check("ir.LandingPad.IsCleanup")
	return binding.LLVMIsCleanup(l.Ref())
}

// CatchSwitch catchswitch 指令角色（内嵌 Value[TokenT]）
type CatchSwitch struct {
	llvm.Value[llvm.TokenT]
}

// AddHandler 追加处理器入口块（须与 catchswitch 同函数）
func (s CatchSwitch) AddHandler(blk Block) {
	const op = "ir.CatchSwitch.AddHandler"
	s.Check(op)
	blk.Check(op)
	binding.LLVMAddHandler(s.Ref(), blk.ref)
}

// HandlerCount 处理器数量
func (s CatchSwitch) HandlerCount() uint32 {
	s.Check("ir.CatchSwitch.HandlerCount")
	return binding.LLVMGetNumHandlers(s.Ref())
}

// HandlerAt 第 i 个处理器入口块
func (s CatchSwitch) HandlerAt(i uint32) Block {
	const op = "ir.CatchSwitch.HandlerAt"
	s.Check(op)
	if i >= s.HandlerCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "handler index %d out of range", i)
	}
	return wrapBlock(s.Context(), s.Lifetime(), binding.LLVMGetHandlers(s.Ref())[i])
}

// FuncletPad funclet pad 指令角色（catchpad/cleanuppad，内嵌 Value[TokenT]）
type FuncletPad struct {
	llvm.Value[llvm.TokenT]
}

// ArgCount 实参个数（LLVMGetNumOperands 含末位 parent pad，故减 1）
func (p FuncletPad) ArgCount() uint32 {
	p.Check("ir.FuncletPad.ArgCount")
	return uint32(binding.LLVMGetNumOperands(p.Ref())) - 1
}

// Arg 第 i 个实参（擦除种类）
func (p FuncletPad) Arg(i uint32) llvm.Value[llvm.DynT] {
	const op = "ir.FuncletPad.Arg"
	p.Check(op)
	if i >= p.ArgCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "argument index %d out of range", i)
	}
	return llvm.ValueOf(p.Context(), p.Lifetime(), binding.LLVMGetArgOperand(p.Ref(), i))
}

// SetArg 替换第 i 个实参
func (p FuncletPad) SetArg(i uint32, v llvm.AnyValue) {
	const op = "ir.FuncletPad.SetArg"
	p.Check(op)
	if i >= p.ArgCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "argument index %d out of range", i)
	}
	p.Context().CheckValues(op, v)
	binding.LLVMSetArgOperand(p.Ref(), i, v.Ref())
}

// ParentCatchSwitch 所属 catchswitch（仅 catchpad 有意义）
func (p FuncletPad) ParentCatchSwitch() CatchSwitch {
	p.Check("ir.FuncletPad.ParentCatchSwitch")
	return CatchSwitch{Value: llvm.NewValue[llvm.TokenT](p.Context(), p.Lifetime(), binding.LLVMGetParentCatchSwitch(p.Ref()))}
}

// ===== EH 构建方法 =====

// Invoke 插入 invoke 调用：then 为正常出口、unwind 为异常出口；
// 返回种类 U 在调用前与函数返回类型比对，不符 panic
func (b *Builder) Invoke[U llvm.Kind](fn llvm.ValueRef[llvm.FnT], args []llvm.AnyValue, then, unwind Block, name string) Invoke[U] {
	const op = "ir.Builder.Invoke"
	fv := fn.AsValue()
	b.pre(op, core(fv))
	b.preBlockOwn(op, then)
	b.preBlockOwn(op, unwind)
	sig := llvm.AsFnType(llvm.TypeOfRef(b.ctx, binding.LLVMGetFunctionType(fv.Ref())))
	checkKind[U](op, b.ctx, binding.LLVMGetReturnType(sig.Ref()))
	b.checkCallArgs(op, sig, args)
	ref := binding.LLVMBuildInvoke(b.ref, sig.Ref(), fv.Ref(), b.valueRefs(args), then.ref, unwind.ref, name)
	return Invoke[U]{Value: llvm.NewValue[U](b.ctx, b.inserted.life, ref)}
}

// InvokeIndirect 通过函数指针 invoke（不透明指针 + 签名）；返回种类 U 与签名返回类型比对
func (b *Builder) InvokeIndirect[U llvm.Kind](fnPtr llvm.ValueRef[llvm.PtrT], sig llvm.FnType, args []llvm.AnyValue, then, unwind Block, name string) Invoke[U] {
	const op = "ir.Builder.InvokeIndirect"
	pv := fnPtr.AsValue()
	b.pre(op, core(pv))
	b.preBlockOwn(op, then)
	b.preBlockOwn(op, unwind)
	if sig.Context() != b.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "signature belongs to another context")
	}
	checkKind[U](op, b.ctx, binding.LLVMGetReturnType(sig.Ref()))
	b.checkCallArgs(op, sig, args)
	ref := binding.LLVMBuildInvoke(b.ref, sig.Ref(), pv.Ref(), b.valueRefs(args), then.ref, unwind.ref, name)
	return Invoke[U]{Value: llvm.NewValue[U](b.ctx, b.inserted.life, ref)}
}

// LandingPad 插入 landingpad（t 须为首类聚合类型）；personality 经 Function.SetPersonality 设置。
// 注：LLVM 22 的 PersFn 构建参数已废弃，传零值
func (b *Builder) LandingPad[T llvm.Kind](t llvm.TypeRef[T], name string) LandingPad[T] {
	const op = "ir.Builder.LandingPad"
	tt := t.AsType()
	b.pre(op)
	if tt.Context() != b.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "type belongs to another context")
	}
	ref := binding.LLVMBuildLandingPad(b.ref, tt.Ref(), binding.LLVMValueRef{}, 0, name)
	return LandingPad[T]{Value: llvm.NewValue[T](b.ctx, b.inserted.life, ref)}
}

// Resume 以 landingpad 值恢复异常传播（void 指令，无 name）
func (b *Builder) Resume(exn llvm.AnyValue) llvm.Value[llvm.VoidT] {
	const op = "ir.Builder.Resume"
	b.pre(op, coreAny(exn))
	ref := binding.LLVMBuildResume(b.ref, exn.Ref())
	return llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)
}

// CatchSwitch 插入 catchswitch（终结指令）；parent 为 nil 表示 within none，
// unwindTo 为零块表示 unwind to caller。处理器用 CatchSwitch.AddHandler 追加
func (b *Builder) CatchSwitch(parent llvm.ValueRef[llvm.TokenT], unwindTo Block, name string) CatchSwitch {
	const op = "ir.Builder.CatchSwitch"
	b.pre(op)
	parentRef := b.preParentPad(op, parent)
	unwindRef := b.preOptBlock(op, unwindTo)
	ref := binding.LLVMBuildCatchSwitch(b.ref, parentRef, unwindRef, 0, name)
	return CatchSwitch{Value: llvm.NewValue[llvm.TokenT](b.ctx, b.inserted.life, ref)}
}

// CatchPad 插入 catchpad（须为块首指令）；parent 为 nil 表示 within none
func (b *Builder) CatchPad(parent llvm.ValueRef[llvm.TokenT], args []llvm.AnyValue, name string) FuncletPad {
	const op = "ir.Builder.CatchPad"
	b.pre(op)
	parentRef := b.preParentPad(op, parent)
	for _, a := range args {
		b.checkVal(op, coreAny(a))
	}
	ref := binding.LLVMBuildCatchPad(b.ref, parentRef, b.valueRefs(args), name)
	return FuncletPad{Value: llvm.NewValue[llvm.TokenT](b.ctx, b.inserted.life, ref)}
}

// CleanupPad 插入 cleanuppad（须为块首指令）；parent 为 nil 表示 within none
func (b *Builder) CleanupPad(parent llvm.ValueRef[llvm.TokenT], args []llvm.AnyValue, name string) FuncletPad {
	const op = "ir.Builder.CleanupPad"
	b.pre(op)
	parentRef := b.preParentPad(op, parent)
	for _, a := range args {
		b.checkVal(op, coreAny(a))
	}
	ref := binding.LLVMBuildCleanupPad(b.ref, parentRef, b.valueRefs(args), name)
	return FuncletPad{Value: llvm.NewValue[llvm.TokenT](b.ctx, b.inserted.life, ref)}
}

// CatchRet 从 catchpad 转移到目标块（void 终结指令，无 name）
func (b *Builder) CatchRet(pad FuncletPad, to Block) llvm.Value[llvm.VoidT] {
	const op = "ir.Builder.CatchRet"
	b.pre(op, core(pad.Value))
	b.preBlockOwn(op, to)
	ref := binding.LLVMBuildCatchRet(b.ref, pad.Ref(), to.ref)
	return llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)
}

// CleanupRet 从 cleanuppad 转移；unwindTo 为零块表示 unwind to caller（void 终结指令，无 name）
func (b *Builder) CleanupRet(pad FuncletPad, unwindTo Block) llvm.Value[llvm.VoidT] {
	const op = "ir.Builder.CleanupRet"
	b.pre(op, core(pad.Value))
	unwindRef := b.preOptBlock(op, unwindTo)
	ref := binding.LLVMBuildCleanupRet(b.ref, pad.Ref(), unwindRef)
	return llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)
}

// preParentPad 预检可选父 pad（nil = within none）；返回底层句柄（零值 = none）
func (b *Builder) preParentPad(op string, parent llvm.ValueRef[llvm.TokenT]) binding.LLVMValueRef {
	if parent == nil {
		return binding.LLVMValueRef{}
	}
	pv := parent.AsValue()
	b.checkVal(op, core(pv))
	return pv.Ref()
}

// preOptBlock 预检可选目标块（零块 = to/unwind to caller）；返回底层句柄（零值 = caller）
func (b *Builder) preOptBlock(op string, blk Block) binding.LLVMBasicBlockRef {
	if blk.ref.IsNil() {
		return binding.LLVMBasicBlockRef{}
	}
	b.preBlockOwn(op, blk)
	return blk.ref
}
