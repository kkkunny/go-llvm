package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// Metadata 元数据句柄（由 Context 唯一化，无独立释放）
type Metadata struct {
	ref binding.LLVMMetadataRef
	ctx *Context
}

// MetadataOf 由底层句柄铸造元数据（供 llvm/* 子包桥接使用）
func MetadataOf(ctx *Context, ref binding.LLVMMetadataRef) Metadata {
	return Metadata{ref: ref, ctx: ctx}
}

// Ref 返回底层句柄（供 llvm/* 子包桥接使用）
func (m Metadata) Ref() binding.LLVMMetadataRef { return m.ref }

// Context 返回所属上下文
func (m Metadata) Context() *Context { return m.ctx }

// IsNil 是否为空句柄
func (m Metadata) IsNil() bool { return m.ref.IsNil() }

// Check 元数据可用性前置校验
func (m Metadata) Check(op string) {
	if m.ref.IsNil() || m.ctx == nil {
		errPanic(ErrInvalidArg, op, "nil metadata")
	}
	m.ctx.CheckAlive(op)
}

// Value 元数据的值形式（Value[MetaT]）
func (m Metadata) Value() Value[MetaT] {
	m.Check("llvm.Metadata.Value")
	return newValue[MetaT](m.ctx, m.ctx.life, binding.LLVMMetadataAsValue(m.ctx.ref, m.ref))
}

// IsString 是否 MDString
func (m Metadata) IsString() bool {
	m.Check("llvm.Metadata.IsString")
	return !binding.LLVMIsAMDString(m.Value().Ref()).IsNil()
}

// IsNode 是否 MDNode
func (m Metadata) IsNode() bool {
	m.Check("llvm.Metadata.IsNode")
	return !binding.LLVMIsAMDNode(m.Value().Ref()).IsNil()
}

// IsValueAsMetadata 是否值的元数据包装（ValueAsMetadata）
func (m Metadata) IsValueAsMetadata() bool {
	m.Check("llvm.Metadata.IsValueAsMetadata")
	return !binding.LLVMIsAValueAsMetadata(m.Value().Ref()).IsNil()
}

// StringValue MDString 的字符串（非 MDString panic）
func (m Metadata) StringValue() string {
	const op = "llvm.Metadata.StringValue"
	m.Check(op)
	if !m.IsString() {
		errPanic(ErrTypeMismatch, op, "metadata is not an MDString")
	}
	s, _ := binding.LLVMGetMDString(m.Value().Ref())
	return s
}

// Operands MDNode 的操作数（非 MDNode panic）
func (m Metadata) Operands() []Metadata {
	const op = "llvm.Metadata.Operands"
	m.Check(op)
	if !m.IsNode() {
		errPanic(ErrTypeMismatch, op, "metadata is not an MDNode")
	}
	vals := binding.LLVMGetMDNodeOperands(m.Value().Ref())
	if len(vals) == 0 {
		return nil
	}
	operands := make([]Metadata, len(vals))
	for i, v := range vals {
		operands[i] = Metadata{ref: binding.LLVMValueAsMetadata(v), ctx: m.ctx}
	}
	return operands
}

// String 元数据的 IR 文本形式
func (m Metadata) String() string {
	return m.Value().String()
}

// MDString 构造 MDString
func (ctx *Context) MDString(s string) Metadata {
	ctx.CheckAlive("llvm.Context.MDString")
	return Metadata{ref: binding.LLVMMDStringInContext2(ctx.ref, s), ctx: ctx}
}

// MDNode 构造 MDNode；元素须属于同一 Context
func (ctx *Context) MDNode(elems ...Metadata) Metadata {
	const op = "llvm.Context.MDNode"
	ctx.CheckAlive(op)
	refs := make([]binding.LLVMMetadataRef, len(elems))
	for i, e := range elems {
		e.Check(op)
		if e.ctx != ctx {
			errPanic(ErrCrossContext, op, "metadata belongs to another context")
		}
		refs[i] = e.ref
	}
	return Metadata{ref: binding.LLVMMDNodeInContext2(ctx.ref, refs), ctx: ctx}
}

// ValueAsMetadata 把值包装为元数据（如 MDNode 的常量操作数）
func (ctx *Context) ValueAsMetadata(v AnyValue) Metadata {
	const op = "llvm.Context.ValueAsMetadata"
	ctx.CheckValues(op, v)
	return Metadata{ref: binding.LLVMValueAsMetadata(v.Ref()), ctx: ctx}
}

// ModuleFlagBehavior 模块级 flag 的合并行为
type ModuleFlagBehavior binding.LLVMModuleFlagBehavior

const (
	ModuleFlagError        = ModuleFlagBehavior(binding.LLVMModuleFlagBehaviorError)
	ModuleFlagWarning      = ModuleFlagBehavior(binding.LLVMModuleFlagBehaviorWarning)
	ModuleFlagRequire      = ModuleFlagBehavior(binding.LLVMModuleFlagBehaviorRequire)
	ModuleFlagOverride     = ModuleFlagBehavior(binding.LLVMModuleFlagBehaviorOverride)
	ModuleFlagAppend       = ModuleFlagBehavior(binding.LLVMModuleFlagBehaviorAppend)
	ModuleFlagAppendUnique = ModuleFlagBehavior(binding.LLVMModuleFlagBehaviorAppendUnique)
)
