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
	closed   bool
}

// NewBuilder 创建构建器并登记到 Context 生命周期
func NewBuilder(ctx *llvm.Context) *Builder {
	b := &Builder{ref: binding.LLVMCreateBuilderInContext(ctx.Ref()), ctx: ctx}
	ctx.Own(b)
	return b
}

// Close 释放构建器；二次调用返回 ErrClosed
func (b *Builder) Close() error {
	if b.closed {
		return &llvm.Error{Reason: llvm.ErrClosed, Op: "ir.Builder.Close", Msg: "builder already closed"}
	}
	b.closed = true
	b.inserted = nil
	binding.LLVMDisposeBuilder(b.ref)
	return nil
}

// Context 返回所属上下文
func (b *Builder) Context() *llvm.Context { return b.ctx }

// MoveToEnd 将插入点移到基本块末尾
func (b *Builder) MoveToEnd(blk Block) {
	b.inserted = &blk
	binding.LLVMPositionBuilderAtEnd(b.ref, blk.ref)
}

// MoveBefore 将插入点移到指令之前
func (b *Builder) MoveBefore(inst llvm.AnyValue) {
	ref := binding.LLVMGetInstructionParent(inst.Ref())
	b.inserted = &Block{ref: ref, ctx: b.ctx, life: inst.Lifetime()}
	binding.LLVMPositionBuilderBefore(b.ref, inst.Ref())
}

// CurrentBlock 当前插入块
func (b *Builder) CurrentBlock() (Block, bool) {
	ref := binding.LLVMGetInsertBlock(b.ref)
	if ref.IsNil() {
		return Block{}, false
	}
	return Block{ref: ref, ctx: b.ctx, life: b.inserted.life}, true
}

// ===== 统一预检 =====

// pre 所有 Builder 方法入口统一调用；v 可为 nil（无操作数指令）
func (b *Builder) pre(op string, vs ...llvm.AnyValue) {
	if b.closed {
		errPanic(llvm.ErrClosed, op, "builder already closed")
	}
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

// preBlock 预检基本块归属
func (b *Builder) preBlock(op string, blk Block) {
	b.pre(op)
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
