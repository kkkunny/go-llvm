package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// ===== 命名元数据 / 模块 flag =====

// AddNamedMetadataOperand 向命名元数据节点追加操作数（不存在则创建）
func (m *Module) AddNamedMetadataOperand(name string, md llvm.Metadata) {
	const op = "ir.Module.AddNamedMetadataOperand"
	m.Check(op)
	md.Check(op)
	if md.Context() != m.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "metadata belongs to another context")
	}
	binding.LLVMAddNamedMetadataOperand(m.ref, name, md.Value().Ref())
}

// NamedMetadataCount 命名元数据节点数
func (m *Module) NamedMetadataCount(name string) uint32 {
	m.Check("ir.Module.NamedMetadataCount")
	return binding.LLVMGetNamedMetadataNumOperands(m.ref, name)
}

// NamedMetadataOperands 命名元数据节点的操作数
func (m *Module) NamedMetadataOperands(name string) []llvm.Metadata {
	m.Check("ir.Module.NamedMetadataOperands")
	vals := binding.LLVMGetNamedMetadataOperands(m.ref, name)
	if len(vals) == 0 {
		return nil
	}
	operands := make([]llvm.Metadata, len(vals))
	for i, v := range vals {
		operands[i] = llvm.MetadataOf(m.ctx, binding.LLVMValueAsMetadata(v))
	}
	return operands
}

// AddModuleFlag 追加模块级 flag（键已存在时不覆盖）
func (m *Module) AddModuleFlag(behavior llvm.ModuleFlagBehavior, key string, val llvm.Metadata) {
	const op = "ir.Module.AddModuleFlag"
	m.Check(op)
	val.Check(op)
	if val.Context() != m.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "metadata belongs to another context")
	}
	binding.LLVMAddModuleFlag(m.ref, binding.LLVMModuleFlagBehavior(behavior), key, val.Ref())
}

// ModuleFlag 模块级 flag（未设置返回 false）
func (m *Module) ModuleFlag(key string) (llvm.Metadata, bool) {
	m.Check("ir.Module.ModuleFlag")
	ref := binding.LLVMGetModuleFlag(m.ref, key)
	if ref.IsNil() {
		return llvm.Metadata{}, false
	}
	return llvm.MetadataOf(m.ctx, ref), true
}

// ===== 指令元数据附着 =====

// AttachMetadata 给指令附着指定种类的元数据（如 !dbg、!prof；LLVM 要求为 MDNode）
func AttachMetadata(inst llvm.AnyValue, kind string, md llvm.Metadata) {
	const op = "ir.AttachMetadata"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead instruction")
	}
	md.Check(op)
	if inst.Context() != md.Context() {
		llvm.Panicf(llvm.ErrCrossContext, op, "metadata belongs to another context")
	}
	if !md.IsNode() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "instruction metadata must be an MDNode")
	}
	kindID := binding.LLVMGetMDKindIDInContext(inst.Context().Ref(), kind)
	binding.LLVMSetMetadata(inst.Ref(), kindID, md.Value().Ref())
}

// InstMetadata 读取指令上指定种类的元数据（不存在返回 false）
func InstMetadata(inst llvm.AnyValue, kind string) (llvm.Metadata, bool) {
	const op = "ir.InstMetadata"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead instruction")
	}
	kindID := binding.LLVMGetMDKindIDInContext(inst.Context().Ref(), kind)
	ref := binding.LLVMGetMetadata(inst.Ref(), kindID)
	if ref.IsNil() {
		return llvm.Metadata{}, false
	}
	return llvm.MetadataOf(inst.Context(), binding.LLVMValueAsMetadata(ref)), true
}

// ===== blockaddress =====

// BlockAddress 构造 blockaddress 常量（指向函数内某基本块）
func BlockAddress(fn Function, blk Block) llvm.Value[llvm.PtrT] {
	const op = "ir.BlockAddress"
	fn.Check(op)
	blk.Check(op)
	if fn.Context() != blk.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "block belongs to another context")
	}
	ref := binding.LLVMBlockAddress(fn.Ref(), blk.ref)
	return llvm.NewValue[llvm.PtrT](fn.Context(), fn.Lifetime(), ref)
}

// BlockAddressFunction blockaddress 常量对应的函数
func BlockAddressFunction(ba llvm.AnyValue) Function {
	const op = "ir.BlockAddressFunction"
	if ba == nil || !ba.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead blockaddress")
	}
	ref := binding.LLVMGetBlockAddressFunction(ba.Ref())
	return Function{Value: llvm.NewValue[llvm.FnT](ba.Context(), ba.Lifetime(), ref)}
}

// BlockAddressBlock blockaddress 常量对应的基本块
func BlockAddressBlock(ba llvm.AnyValue) Block {
	const op = "ir.BlockAddressBlock"
	if ba == nil || !ba.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead blockaddress")
	}
	ref := binding.LLVMGetBlockAddressBasicBlock(ba.Ref())
	return wrapBlock(ba.Context(), ba.Lifetime(), ref)
}
