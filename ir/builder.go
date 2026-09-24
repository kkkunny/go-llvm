package ir

import (
	"strings"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
)

// Builder IR 构建器
type Builder struct {
	ref      binding.LLVMBuilderRef
	ctx      *llvm.Context
	inserted *Block
	unown    func()
	closed   bool
	refs     []binding.LLVMValueRef // 句柄转换 scratch（单 goroutine 使用）
	i8       llvm.IntType           // 惰性缓存的 i8 类型

	ownerGID uint64   // 调试层：owner goroutine id（0=未记录）
	ops      uint64   // 调试层：采样计数
	recent   []string // 调试层：最近 op 记录（A1 现场 dump）
}

// recentOps 调试层记录/展示的最近操作数
const recentOps = 8

// NewBuilder 创建构建器并登记到 Context 生命周期
func NewBuilder(ctx *llvm.Context) *Builder {
	if !ctx.Alive() {
		llvm.Panicf(llvm.ErrUseAfterFree, "ir.NewBuilder", "context is closed")
	}
	b := &Builder{ref: binding.LLVMCreateBuilderInContext(ctx.Ref()), ctx: ctx}
	b.unown = ctx.Own(b)
	return b
}

// NewBuilderAt 创建构建器并把插入点定位到 blk 末尾
func NewBuilderAt(blk Block) *Builder {
	blk.Check("ir.NewBuilderAt")
	b := NewBuilder(blk.ctx)
	b.MoveToEnd(blk)
	return b
}

// i8Type 惰性缓存的 i8 类型（PtrAdd 等使用）
func (b *Builder) i8Type() llvm.IntType {
	if b.i8.IsNil() {
		b.i8 = b.ctx.Int(8)
	}
	return b.i8
}

// Close 释放构建器；二次调用返回 ErrClosed
func (b *Builder) Close() error {
	if b.closed {
		return &llvm.Error{Reason: llvm.ErrClosed, Op: "ir.Builder.Close", Msg: "builder already closed"}
	}
	checkOwner("ir.Builder.Close", &b.ownerGID, &b.ops, "builder", true)
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
	checkOwner(op, &b.ownerGID, &b.ops, "builder", true)
	b.preBlockOwn(op, blk)
	b.inserted = &blk
	binding.LLVMPositionBuilderAtEnd(b.ref, blk.ref)
}

// MoveBefore 将插入点移到指令之前
func (b *Builder) MoveBefore(inst llvm.AnyValue) {
	const op = "ir.Builder.MoveBefore"
	b.preAlive(op)
	checkOwner(op, &b.ownerGID, &b.ops, "builder", true)
	b.checkVal(op, coreAny(inst))
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

// preVal 值句柄的校验视图。
// 预检若直接接收 llvm.AnyValue，每次传参都会把 Value[T] 装箱为接口并堆分配；
// 内部改用本结构（纯值类型）消除该开销。
// ty 为调试层缓存的操作数类型句柄（见 llvm.Value.RawType），release 下恒为空。
type preVal struct {
	ref  binding.LLVMValueRef
	ty   binding.LLVMTypeRef
	ctx  *llvm.Context
	life *llvm.Lifetime
	ok   bool
}

// core 由具体值构造校验视图（无装箱）
func core[T llvm.Kind](v llvm.Value[T]) preVal {
	return preVal{ref: v.Ref(), ty: rawType(v), ctx: v.Context(), life: v.Lifetime(), ok: true}
}

// rawTypeer 值角色可提供调试层缓存的类型句柄，避免预检查询 cgo
type rawTypeer interface{ RawType() binding.LLVMTypeRef }

// rawType 取缓存类型；release 构建不查询（返回空句柄）
func rawType[T llvm.Kind](v llvm.Value[T]) binding.LLVMTypeRef {
	if checks.Debug {
		return v.RawType()
	}
	return binding.LLVMTypeRef{}
}

// coreAny 由任意值接口构造校验视图（不产生新的装箱）
func coreAny(v llvm.AnyValue) preVal {
	if v == nil {
		return preVal{}
	}
	p := preVal{ref: v.Ref(), ctx: v.Context(), life: v.Lifetime(), ok: true}
	if checks.Debug {
		if rt, ok := v.(rawTypeer); ok {
			p.ty = rt.RawType()
		}
	}
	return p
}

// valType 操作数的类型句柄：优先用调试层缓存，缺失时回退查询（仅调试层调用）
func (b *Builder) valType(p preVal) binding.LLVMTypeRef {
	return typeOfVal(p.ref, p.ty)
}

// typeOfVal 句柄类型：优先缓存，缺失回退查询（仅调试层调用）
func typeOfVal(ref binding.LLVMValueRef, cached binding.LLVMTypeRef) binding.LLVMTypeRef {
	if !cached.IsNil() {
		return cached
	}
	return binding.LLVMTypeOf(ref)
}

// anyValType 任意值接口的类型句柄：优先调试层缓存，缺失回退查询（仅调试层调用）
func anyValType(v llvm.AnyValue) binding.LLVMTypeRef {
	if rt, ok := v.(rawTypeer); ok {
		if ty := rt.RawType(); !ty.IsNil() {
			return ty
		}
	}
	return binding.LLVMTypeOf(v.Ref())
}

// typeString 底层值句柄的类型文本（仅错误路径使用）
func typeString(ctx *llvm.Context, ref binding.LLVMValueRef) string {
	return llvm.TypeOfRef(ctx, binding.LLVMTypeOf(ref)).String()
}

// typeRefString 底层类型句柄的文本（仅错误路径使用）
func typeRefString(ctx *llvm.Context, ref binding.LLVMTypeRef) string {
	return llvm.TypeOfRef(ctx, ref).String()
}

// preAlive 校验 Builder 未关闭（不要求已定位，供定位类方法使用）
func (b *Builder) preAlive(op string) {
	if b.closed {
		b.panicf(llvm.ErrClosed, op, "builder already closed")
	}
	checkOwner(op, &b.ownerGID, &b.ops, "builder", false)
}

// panicf 构造并 panic *Error；调试层附加最近操作现场（A1）
func (b *Builder) panicf(reason llvm.ErrKind, op, format string, args ...any) {
	if checks.Debug && len(b.recent) > 0 {
		format += " (recent builder ops: %s)"
		args = append(args, strings.Join(b.recent, ", "))
	}
	llvm.Panicf(reason, op, format, args...)
}

// recordOp 调试层记录最近 op（环形，A1 现场 dump 用）
func (b *Builder) recordOp(op string) {
	if len(b.recent) == recentOps {
		copy(b.recent, b.recent[1:])
		b.recent = b.recent[:recentOps-1]
	}
	b.recent = append(b.recent, op)
}

// prePosition 校验 Builder 已定位到某个基本块（纯 Go：插入点由 Go 侧跟踪，
// 块/模块已释放时其生命周期令牌已失效）
func (b *Builder) prePosition(op string) {
	if b.inserted == nil {
		b.panicf(llvm.ErrInvalidArg, op, "builder is not positioned at any block")
	}
	if !b.inserted.life.Alive() {
		b.panicf(llvm.ErrUseAfterFree, op, "insert block is freed")
	}
}

// checkVal 校验单个操作数；供 pre 与批量校验路径复用（避免构造临时切片）
func (b *Builder) checkVal(op string, v preVal) {
	if !v.ok || v.ref.IsNil() {
		b.panicf(llvm.ErrInvalidArg, op, "nil operand")
	}
	if v.ctx == nil || !v.ctx.Alive() {
		b.panicf(llvm.ErrUseAfterFree, op, "operand context is closed")
	}
	if v.life == nil || !v.life.Alive() {
		b.panicf(llvm.ErrUseAfterFree, op, "operand is freed")
	}
	if v.ctx != b.ctx {
		b.panicf(llvm.ErrCrossContext, op, "operand belongs to another context")
	}
}

// pre 所有 Builder 方法入口统一调用；vs 可为空（无操作数指令）
func (b *Builder) pre(op string, vs ...preVal) {
	if checks.Debug {
		b.recordOp(op)
	}
	b.preAlive(op)
	b.prePosition(op)
	for _, v := range vs {
		b.checkVal(op, v)
	}
}

// preSameType 预检并要求两操作数类型一致（语义契约，仅调试层）
func (b *Builder) preSameType(op string, l, r preVal) {
	b.pre(op, l, r)
	if !checks.Debug {
		return
	}
	if !b.valType(l).Equal(b.valType(r)) {
		b.panicf(llvm.ErrTypeMismatch, op, "operand types differ: %s vs %s", typeString(b.ctx, l.ref), typeString(b.ctx, r.ref))
	}
}

// preAlign 预检对齐值为 2 的幂（语义契约，仅调试层）
func preAlign(op string, n uint32) {
	if !checks.Debug {
		return
	}
	if n == 0 || n&(n-1) != 0 {
		llvm.Panicf(llvm.ErrInvalidArg, op, "alignment %d is not a power of two", n)
	}
}

// preBlockOwn 预检基本块句柄本身（nil/跨 Context/已释放），不要求 Builder 已定位
func (b *Builder) preBlockOwn(op string, blk Block) {
	if blk.ref.IsNil() {
		b.panicf(llvm.ErrInvalidArg, op, "nil block")
	}
	if blk.ctx != b.ctx {
		b.panicf(llvm.ErrCrossContext, op, "block belongs to another context")
	}
	if !blk.life.Alive() {
		b.panicf(llvm.ErrUseAfterFree, op, "block is freed")
	}
}

// preBlock 预检基本块归属
func (b *Builder) preBlock(op string, blk Block) {
	b.pre(op)
	b.preBlockOwn(op, blk)
}

// valueRefs 复用 Builder 的 scratch 缓冲把值列表转换为底层句柄列表。
// 前提：Builder 单 goroutine 使用，且返回值只在紧随其后的 binding 调用内消费。
// 调试层不复用缓冲（B4）：把上述约定违规从"潜在 UB"变为"不可能发生"。
func (b *Builder) valueRefs(vs []llvm.AnyValue) []binding.LLVMValueRef {
	if checks.Debug {
		refs := make([]binding.LLVMValueRef, len(vs))
		for i, v := range vs {
			refs[i] = v.Ref()
		}
		return refs
	}
	if cap(b.refs) < len(vs) {
		b.refs = make([]binding.LLVMValueRef, len(vs))
	}
	refs := b.refs[:len(vs)]
	for i, v := range vs {
		refs[i] = v.Ref()
	}
	b.refs = refs
	return refs
}
