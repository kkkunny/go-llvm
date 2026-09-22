package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Builder IR 构建器
type Builder struct {
	ref      binding.LLVMBuilderRef
	ctx      *llvm.Context
	inserted *Block
	unown    func()
	closed   bool
}

// NewBuilder 创建构建器并登记到 Context 生命周期
func NewBuilder(ctx *llvm.Context) *Builder {
	if !ctx.Alive() {
		errPanic(llvm.ErrUseAfterFree, "ir.NewBuilder", "context is closed")
	}
	b := &Builder{ref: binding.LLVMCreateBuilderInContext(ctx.Ref()), ctx: ctx}
	b.unown = ctx.Own(b)
	return b
}

// Close 释放构建器；二次调用返回 ErrClosed
func (b *Builder) Close() error {
	if b.closed {
		return &llvm.Error{Reason: llvm.ErrClosed, Op: "ir.Builder.Close", Msg: "builder already closed"}
	}
	b.closed = true
	b.inserted = nil
	b.unown()
	binding.LLVMDisposeBuilder(b.ref)
	return nil
}

// Context 返回所属上下文
func (b *Builder) Context() *llvm.Context { return b.ctx }

// MoveToEnd 将插入点移到基本块末尾
func (b *Builder) MoveToEnd(blk Block) {
	const op = "ir.Builder.MoveToEnd"
	b.preAlive(op)
	b.preBlockOwn(op, blk)
	b.inserted = &blk
	binding.LLVMPositionBuilderAtEnd(b.ref, blk.ref)
}

// MoveBefore 将插入点移到指令之前
func (b *Builder) MoveBefore(inst llvm.AnyValue) {
	const op = "ir.Builder.MoveBefore"
	b.preAlive(op)
	b.pre(op, inst)
	ref := binding.LLVMGetInstructionParent(inst.Ref())
	blk := wrapBlock(b.ctx, inst.Lifetime(), ref)
	b.inserted = &blk
	binding.LLVMPositionBuilderBefore(b.ref, inst.Ref())
}

// CurrentBlock 当前插入块
func (b *Builder) CurrentBlock() (Block, bool) {
	ref := binding.LLVMGetInsertBlock(b.ref)
	if ref.IsNil() {
		return Block{}, false
	}
	if b.inserted != nil {
		return wrapBlock(b.ctx, b.inserted.life, ref), true
	}
	return Block{}, false
}

// ===== 统一预检 =====

// preAlive 校验 Builder 未关闭（不要求已定位，供定位类方法使用）
func (b *Builder) preAlive(op string) {
	if b.closed {
		errPanic(llvm.ErrClosed, op, "builder already closed")
	}
}

// pre 所有 Builder 方法入口统一调用；v 可为 nil（无操作数指令）
func (b *Builder) pre(op string, vs ...llvm.AnyValue) {
	b.preAlive(op)
	if b.inserted == nil || binding.LLVMGetInsertBlock(b.ref).IsNil() {
		errPanic(llvm.ErrInvalidArg, op, "builder is not positioned at any block")
	}
	for _, v := range vs {
		if v == nil {
			continue
		}
		if v.IsNil() {
			errPanic(llvm.ErrInvalidArg, op, "nil operand")
		}
		if !v.Alive() {
			errPanic(llvm.ErrUseAfterFree, op, "operand is freed")
		}
		if v.Context() != b.ctx {
			errPanic(llvm.ErrCrossContext, op, "operand belongs to another context")
		}
	}
}

// preSameType 预检并要求两操作数类型一致
func (b *Builder) preSameType(op string, l, r llvm.AnyValue) {
	b.pre(op, l, r)
	if !l.Dyn().Type().Equal(r.Dyn().Type()) {
		errPanic(llvm.ErrTypeMismatch, op, "operand types differ: %s vs %s", l.Dyn().Type(), r.Dyn().Type())
	}
}

// preAlign 预检对齐值为 2 的幂
func preAlign(op string, n uint32) {
	if n == 0 || n&(n-1) != 0 {
		errPanic(llvm.ErrInvalidArg, op, "alignment %d is not a power of two", n)
	}
}

// errPanic 转发到 root 包的 panic 构造（保持错误类型统一）
func errPanic(reason llvm.ErrKind, op, format string, args ...any) {
	llvm.Panicf(reason, op, format, args...)
}

// preBlockOwn 预检基本块句柄本身（nil/跨 Context/已释放），不要求 Builder 已定位
func (b *Builder) preBlockOwn(op string, blk Block) {
	if blk.ref.IsNil() {
		errPanic(llvm.ErrInvalidArg, op, "nil block")
	}
	if blk.ctx != b.ctx {
		errPanic(llvm.ErrCrossContext, op, "block belongs to another context")
	}
	if !blk.life.Alive() {
		errPanic(llvm.ErrUseAfterFree, op, "block is freed")
	}
}

// preBlock 预检基本块归属
func (b *Builder) preBlock(op string, blk Block) {
	b.pre(op)
	b.preBlockOwn(op, blk)
}
