package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/errs"
)

// OperandBundle 调用的操作数捆绑（如 funclet/clang.arc.attachedcall/statepoint）。
// 拥有 C 侧资源，用毕 Close；指令构建时只借用，构建后即可释放。
// 指针语义：Close 后所有持有者都会观察到失效（值类型副本会绕过 closed 检查）。
type OperandBundle struct {
	ref    binding.LLVMOperandBundleRef
	ctx    *llvm.Context
	closed bool
}

// NewOperandBundle 创建操作数捆绑（ctx 用于校验与归属）
func NewOperandBundle(ctx *llvm.Context, tag string, values []llvm.AnyValue) *OperandBundle {
	const op = "ir.NewOperandBundle"
	ctx.CheckAlive(op)
	ctx.CheckValues(op, values...)
	return &OperandBundle{ref: binding.LLVMCreateOperandBundle(tag, llvm.AnyValuesToRefs(values)), ctx: ctx}
}

// Tag 捆绑标签
func (b *OperandBundle) Tag() string {
	b.check("ir.OperandBundle.Tag")
	return binding.LLVMGetOperandBundleTag(b.ref)
}

// check 前置校验（崩溃类地板：nil/已关闭/上下文失效）
func (b *OperandBundle) check(op string) {
	if b == nil || b.ref.IsNil() || b.closed {
		errs.Panicf(llvm.ErrUseAfterFree, op, "operand bundle is closed")
	}
	if b.ctx == nil || !b.ctx.Alive() {
		errs.Panicf(llvm.ErrUseAfterFree, op, "operand bundle context is closed")
	}
}

// Close 释放捆绑；二次调用返回 ErrClosed
func (b *OperandBundle) Close() error {
	if b == nil || b.closed {
		return &llvm.Error{Reason: llvm.ErrClosed, Op: "ir.OperandBundle.Close", Msg: "operand bundle already closed"}
	}
	b.closed = true
	binding.LLVMDisposeOperandBundle(b.ref)
	return nil
}

// bundlesToRefs 取句柄列表（内部用；在指令构建时借用捆绑）
func bundlesToRefs(bundles []*OperandBundle) []binding.LLVMOperandBundleRef {
	refs := make([]binding.LLVMOperandBundleRef, len(bundles))
	for i, b := range bundles {
		b.check("ir.bundlesToRefs")
		refs[i] = b.ref
	}
	return refs
}
