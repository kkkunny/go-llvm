package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// ===== 属性位置共享助手（函数/调用点通用） =====

// addAttr 在 idx 位置追加属性；callSite 为真走调用点 API（LLVM 的 Function 与 CallBase 属性 API 分开）
func addAttr(op string, ctx *llvm.Context, ref binding.LLVMValueRef, idx llvm.AttrIndex, a llvm.Attribute, callSite bool) {
	a.Check(op)
	if a.Context() != ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "attribute belongs to another context")
	}
	if callSite {
		binding.LLVMAddCallSiteAttribute(ref, binding.LLVMAttributeIndex(idx), a.Ref())
		return
	}
	binding.LLVMAddAttributeAtIndex(ref, binding.LLVMAttributeIndex(idx), a.Ref())
}

// attrCount idx 位置属性个数
func attrCount(ref binding.LLVMValueRef, idx llvm.AttrIndex, callSite bool) uint32 {
	if callSite {
		return binding.LLVMGetCallSiteAttributeCount(ref, binding.LLVMAttributeIndex(idx))
	}
	return binding.LLVMGetAttributeCountAtIndex(ref, binding.LLVMAttributeIndex(idx))
}

// attrsAt idx 位置全部属性
func attrsAt(ctx *llvm.Context, ref binding.LLVMValueRef, idx llvm.AttrIndex, callSite bool) []llvm.Attribute {
	var refs []binding.LLVMAttributeRef
	if callSite {
		refs = binding.LLVMGetCallSiteAttributes(ref, binding.LLVMAttributeIndex(idx))
	} else {
		refs = binding.LLVMGetAttributesAtIndex(ref, binding.LLVMAttributeIndex(idx))
	}
	if len(refs) == 0 {
		return nil
	}
	attrs := make([]llvm.Attribute, len(refs))
	for i, r := range refs {
		attrs[i] = llvm.AttributeOf(ctx, r)
	}
	return attrs
}

// enumAttrAt idx 位置的枚举属性（不存在返回 false）
func enumAttrAt(ctx *llvm.Context, ref binding.LLVMValueRef, idx llvm.AttrIndex, kind llvm.AttributeKind, callSite bool) (llvm.Attribute, bool) {
	var r binding.LLVMAttributeRef
	if callSite {
		r = binding.LLVMGetCallSiteEnumAttribute(ref, binding.LLVMAttributeIndex(idx), uint32(kind))
	} else {
		r = binding.LLVMGetEnumAttributeAtIndex(ref, binding.LLVMAttributeIndex(idx), uint32(kind))
	}
	if r.IsNil() {
		return llvm.Attribute{}, false
	}
	return llvm.AttributeOf(ctx, r), true
}

// stringAttrAt idx 位置的字符串属性（不存在返回 false）
func stringAttrAt(ctx *llvm.Context, ref binding.LLVMValueRef, idx llvm.AttrIndex, kind string, callSite bool) (llvm.Attribute, bool) {
	var r binding.LLVMAttributeRef
	if callSite {
		r = binding.LLVMGetCallSiteStringAttribute(ref, binding.LLVMAttributeIndex(idx), kind)
	} else {
		r = binding.LLVMGetStringAttributeAtIndex(ref, binding.LLVMAttributeIndex(idx), kind)
	}
	if r.IsNil() {
		return llvm.Attribute{}, false
	}
	return llvm.AttributeOf(ctx, r), true
}

// removeEnumAttrAt 移除 idx 位置的枚举属性
func removeEnumAttrAt(ref binding.LLVMValueRef, idx llvm.AttrIndex, kind llvm.AttributeKind, callSite bool) {
	if callSite {
		binding.LLVMRemoveCallSiteEnumAttribute(ref, binding.LLVMAttributeIndex(idx), uint32(kind))
		return
	}
	binding.LLVMRemoveEnumAttributeAtIndex(ref, binding.LLVMAttributeIndex(idx), uint32(kind))
}

// removeStringAttrAt 移除 idx 位置的字符串属性
func removeStringAttrAt(ref binding.LLVMValueRef, idx llvm.AttrIndex, kind string, callSite bool) {
	if callSite {
		binding.LLVMRemoveCallSiteStringAttribute(ref, binding.LLVMAttributeIndex(idx), kind)
		return
	}
	binding.LLVMRemoveStringAttributeAtIndex(ref, binding.LLVMAttributeIndex(idx), kind)
}

// ===== Function =====

// AddAttr 在 idx 位置（返回/函数/参数）追加属性
func (f Function) AddAttr(idx llvm.AttrIndex, a llvm.Attribute) {
	const op = "ir.Function.AddAttr"
	addAttr(op, f.Context(), f.Ref(), idx, a, false)
}

// AttrCount idx 位置的属性个数
func (f Function) AttrCount(idx llvm.AttrIndex) uint32 {
	return attrCount(f.Ref(), idx, false)
}

// Attrs idx 位置的全部属性
func (f Function) Attrs(idx llvm.AttrIndex) []llvm.Attribute {
	return attrsAt(f.Context(), f.Ref(), idx, false)
}

// EnumAttr idx 位置的枚举属性（不存在返回 false）
func (f Function) EnumAttr(idx llvm.AttrIndex, kind llvm.AttributeKind) (llvm.Attribute, bool) {
	return enumAttrAt(f.Context(), f.Ref(), idx, kind, false)
}

// StringAttr idx 位置的字符串属性（不存在返回 false）
func (f Function) StringAttr(idx llvm.AttrIndex, kind string) (llvm.Attribute, bool) {
	return stringAttrAt(f.Context(), f.Ref(), idx, kind, false)
}

// RemoveEnumAttr 移除 idx 位置的枚举属性
func (f Function) RemoveEnumAttr(idx llvm.AttrIndex, kind llvm.AttributeKind) {
	removeEnumAttrAt(f.Ref(), idx, kind, false)
}

// RemoveStringAttr 移除 idx 位置的字符串属性
func (f Function) RemoveStringAttr(idx llvm.AttrIndex, kind string) {
	removeStringAttrAt(f.Ref(), idx, kind, false)
}

// CallConv 函数调用约定
func (f Function) CallConv() llvm.CallConv {
	return llvm.CallConv(binding.LLVMGetFunctionCallConv(f.Ref()))
}

// SetCallConv 设置函数调用约定
func (f Function) SetCallConv(cc llvm.CallConv) {
	binding.LLVMSetFunctionCallConv(f.Ref(), binding.LLVMCallConv(cc))
}

// ===== Param =====

// Index 参数在函数中的序号（从 0 起）
func (p Param) Index() uint32 {
	const op = "ir.Param.Index"
	fn := p.Belong()
	n := uint32(fn.CountParams())
	for i := uint32(0); i < n; i++ {
		if binding.LLVMGetParam(fn.Ref(), i).Equal(p.Ref()) {
			return i
		}
	}
	llvm.Panicf(llvm.ErrNotFound, op, "parameter not found in parent function")
	return 0
}

// AddAttr 追加参数属性
func (p Param) AddAttr(a llvm.Attribute) {
	const op = "ir.Param.AddAttr"
	p.Check(op)
	addAttr(op, p.Context(), p.Belong().Ref(), llvm.AttrParam(uint(p.Index())), a, false)
}

// AttrCount 参数属性个数
func (p Param) AttrCount() uint32 {
	p.Check("ir.Param.AttrCount")
	return attrCount(p.Belong().Ref(), llvm.AttrParam(uint(p.Index())), false)
}

// Attrs 参数全部属性
func (p Param) Attrs() []llvm.Attribute {
	p.Check("ir.Param.Attrs")
	return attrsAt(p.Context(), p.Belong().Ref(), llvm.AttrParam(uint(p.Index())), false)
}

// EnumAttr 参数的枚举属性（不存在返回 false）
func (p Param) EnumAttr(kind llvm.AttributeKind) (llvm.Attribute, bool) {
	p.Check("ir.Param.EnumAttr")
	return enumAttrAt(p.Context(), p.Belong().Ref(), llvm.AttrParam(uint(p.Index())), kind, false)
}

// StringAttr 参数的字符串属性（不存在返回 false）
func (p Param) StringAttr(kind string) (llvm.Attribute, bool) {
	p.Check("ir.Param.StringAttr")
	return stringAttrAt(p.Context(), p.Belong().Ref(), llvm.AttrParam(uint(p.Index())), kind, false)
}

// RemoveEnumAttr 移除参数的枚举属性
func (p Param) RemoveEnumAttr(kind llvm.AttributeKind) {
	p.Check("ir.Param.RemoveEnumAttr")
	removeEnumAttrAt(p.Belong().Ref(), llvm.AttrParam(uint(p.Index())), kind, false)
}

// RemoveStringAttr 移除参数的字符串属性
func (p Param) RemoveStringAttr(kind string) {
	p.Check("ir.Param.RemoveStringAttr")
	removeStringAttrAt(p.Belong().Ref(), llvm.AttrParam(uint(p.Index())), kind, false)
}

// ===== Call / Invoke =====

// AddAttr 在 idx 位置（返回/函数/参数）追加调用点属性
func (c Call[T]) AddAttr(idx llvm.AttrIndex, a llvm.Attribute) {
	const op = "ir.Call.AddAttr"
	addAttr(op, c.Context(), c.Ref(), idx, a, true)
}

// AttrCount idx 位置的属性个数
func (c Call[T]) AttrCount(idx llvm.AttrIndex) uint32 {
	return attrCount(c.Ref(), idx, true)
}

// Attrs idx 位置的全部属性
func (c Call[T]) Attrs(idx llvm.AttrIndex) []llvm.Attribute {
	return attrsAt(c.Context(), c.Ref(), idx, true)
}

// EnumAttr idx 位置的枚举属性（不存在返回 false）
func (c Call[T]) EnumAttr(idx llvm.AttrIndex, kind llvm.AttributeKind) (llvm.Attribute, bool) {
	return enumAttrAt(c.Context(), c.Ref(), idx, kind, true)
}

// StringAttr idx 位置的字符串属性（不存在返回 false）
func (c Call[T]) StringAttr(idx llvm.AttrIndex, kind string) (llvm.Attribute, bool) {
	return stringAttrAt(c.Context(), c.Ref(), idx, kind, true)
}

// RemoveEnumAttr 移除 idx 位置的枚举属性
func (c Call[T]) RemoveEnumAttr(idx llvm.AttrIndex, kind llvm.AttributeKind) {
	removeEnumAttrAt(c.Ref(), idx, kind, true)
}

// RemoveStringAttr 移除 idx 位置的字符串属性
func (c Call[T]) RemoveStringAttr(idx llvm.AttrIndex, kind string) {
	removeStringAttrAt(c.Ref(), idx, kind, true)
}

// CallConv 调用点调用约定
func (c Call[T]) CallConv() llvm.CallConv {
	return llvm.CallConv(binding.LLVMGetInstructionCallConv(c.Ref()))
}

// SetCallConv 设置调用点调用约定
func (c Call[T]) SetCallConv(cc llvm.CallConv) {
	binding.LLVMSetInstructionCallConv(c.Ref(), binding.LLVMCallConv(cc))
}

// AddAttr 在 idx 位置（返回/函数/参数）追加调用点属性
func (c Invoke[T]) AddAttr(idx llvm.AttrIndex, a llvm.Attribute) {
	const op = "ir.Invoke.AddAttr"
	addAttr(op, c.Context(), c.Ref(), idx, a, true)
}

// AttrCount idx 位置的属性个数
func (c Invoke[T]) AttrCount(idx llvm.AttrIndex) uint32 {
	return attrCount(c.Ref(), idx, true)
}

// Attrs idx 位置的全部属性
func (c Invoke[T]) Attrs(idx llvm.AttrIndex) []llvm.Attribute {
	return attrsAt(c.Context(), c.Ref(), idx, true)
}

// EnumAttr idx 位置的枚举属性（不存在返回 false）
func (c Invoke[T]) EnumAttr(idx llvm.AttrIndex, kind llvm.AttributeKind) (llvm.Attribute, bool) {
	return enumAttrAt(c.Context(), c.Ref(), idx, kind, true)
}

// StringAttr idx 位置的字符串属性（不存在返回 false）
func (c Invoke[T]) StringAttr(idx llvm.AttrIndex, kind string) (llvm.Attribute, bool) {
	return stringAttrAt(c.Context(), c.Ref(), idx, kind, true)
}

// RemoveEnumAttr 移除 idx 位置的枚举属性
func (c Invoke[T]) RemoveEnumAttr(idx llvm.AttrIndex, kind llvm.AttributeKind) {
	removeEnumAttrAt(c.Ref(), idx, kind, true)
}

// RemoveStringAttr 移除 idx 位置的字符串属性
func (c Invoke[T]) RemoveStringAttr(idx llvm.AttrIndex, kind string) {
	removeStringAttrAt(c.Ref(), idx, kind, true)
}

// CallConv 调用点调用约定
func (c Invoke[T]) CallConv() llvm.CallConv {
	return llvm.CallConv(binding.LLVMGetInstructionCallConv(c.Ref()))
}

// SetCallConv 设置调用点调用约定
func (c Invoke[T]) SetCallConv(cc llvm.CallConv) {
	binding.LLVMSetInstructionCallConv(c.Ref(), binding.LLVMCallConv(cc))
}
