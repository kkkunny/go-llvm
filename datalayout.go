package llvm

import (
	"github.com/kkkunny/go-llvm/internal/binding"
)

// ByteOrder 目标字节序
type ByteOrder int32

const (
	LittleEndian ByteOrder = ByteOrder(binding.LLVMLittleEndian)
	BigEndian    ByteOrder = ByteOrder(binding.LLVMBigEndian)
)

// DataLayout 目标数据布局查询句柄；由 NewDataLayout 或 TargetMachine 创建，用毕 Close
type DataLayout struct {
	ref    binding.LLVMTargetDataRef
	closed bool
}

// NewDataLayout 由布局字符串创建数据布局（拥有句柄）
func NewDataLayout(layout string) *DataLayout {
	ref := binding.LLVMCreateTargetData(layout)
	if ref.IsNil() {
		errPanic(ErrInvalidArg, "llvm.NewDataLayout", "invalid data layout string %q", layout)
	}
	return &DataLayout{ref: ref}
}

// DataLayoutOf 由底层句柄构建数据布局（供 llvm/target 桥接使用，接管所有权）
func DataLayoutOf(ref binding.LLVMTargetDataRef) *DataLayout {
	return &DataLayout{ref: ref}
}

// check 前置校验：句柄非空且未释放
func (d *DataLayout) check(op string) {
	if d.ref.IsNil() {
		errPanic(ErrInvalidArg, op, "nil data layout")
	}
	if d.closed {
		errPanic(ErrUseAfterFree, op, "data layout is closed")
	}
}

// checkType 前置校验类型参数
func (d *DataLayout) checkType(op string, t AnyType) {
	d.check(op)
	if t == nil || t.IsNil() {
		errPanic(ErrInvalidArg, op, "nil type")
	}
}

// checkValue 前置校验值参数
func (d *DataLayout) checkValue(op string, v AnyValue) {
	d.check(op)
	if v == nil || v.IsNil() {
		errPanic(ErrInvalidArg, op, "nil value")
	}
}

// Close 释放句柄；二次调用返回 ErrClosed
func (d *DataLayout) Close() error {
	if d.closed {
		return &Error{Reason: ErrClosed, Op: "llvm.DataLayout.Close", Msg: "data layout already closed"}
	}
	d.closed = true
	binding.LLVMDisposeTargetData(d.ref)
	return nil
}

// String 布局字符串
func (d *DataLayout) String() string {
	d.check("llvm.DataLayout.String")
	return binding.LLVMCopyStringRepOfTargetData(d.ref)
}

// ByteOrder 字节序
func (d *DataLayout) ByteOrder() ByteOrder {
	d.check("llvm.DataLayout.ByteOrder")
	return ByteOrder(binding.LLVMByteOrder(d.ref))
}

// PointerSize 指针字节数
func (d *DataLayout) PointerSize() uint32 {
	d.check("llvm.DataLayout.PointerSize")
	return binding.LLVMPointerSize(d.ref)
}

// PointerSizeForAS 指定地址空间的指针字节数
func (d *DataLayout) PointerSizeForAS(addrspace uint32) uint32 {
	d.check("llvm.DataLayout.PointerSizeForAS")
	return binding.LLVMPointerSizeForAS(d.ref, addrspace)
}

// IntPtrType 指针等宽整数类型
func (d *DataLayout) IntPtrType(ctx *Context) IntType {
	d.check("llvm.DataLayout.IntPtrType")
	ctx.CheckAlive("llvm.DataLayout.IntPtrType")
	return IntType{Type[IntT]{ref: binding.LLVMIntPtrTypeInContext(ctx.ref, d.ref), ctx: ctx}}
}

// IntPtrTypeForAS 指定地址空间指针等宽整数类型
func (d *DataLayout) IntPtrTypeForAS(ctx *Context, addrspace uint32) IntType {
	d.check("llvm.DataLayout.IntPtrTypeForAS")
	ctx.CheckAlive("llvm.DataLayout.IntPtrTypeForAS")
	return IntType{Type[IntT]{ref: binding.LLVMIntPtrTypeForASInContext(ctx.ref, d.ref, addrspace), ctx: ctx}}
}

// SizeOfTypeInBits 类型位宽（bit）
func (d *DataLayout) SizeOfTypeInBits(t AnyType) uint64 {
	d.checkType("llvm.DataLayout.SizeOfTypeInBits", t)
	return binding.LLVMSizeOfTypeInBits(d.ref, t.Ref())
}

// StoreSizeOfType 实际存储大小（byte）
func (d *DataLayout) StoreSizeOfType(t AnyType) uint64 {
	d.checkType("llvm.DataLayout.StoreSizeOfType", t)
	return binding.LLVMStoreSizeOfType(d.ref, t.Ref())
}

// ABISizeOfType ABI 大小（byte）
func (d *DataLayout) ABISizeOfType(t AnyType) uint64 {
	d.checkType("llvm.DataLayout.ABISizeOfType", t)
	return binding.LLVMABISizeOfType(d.ref, t.Ref())
}

// ABIAlignOfType ABI 对齐（byte）
func (d *DataLayout) ABIAlignOfType(t AnyType) uint32 {
	d.checkType("llvm.DataLayout.ABIAlignOfType", t)
	return binding.LLVMABIAlignmentOfType(d.ref, t.Ref())
}

// CallFrameAlignOfType 调用栈帧对齐（byte）
func (d *DataLayout) CallFrameAlignOfType(t AnyType) uint32 {
	d.checkType("llvm.DataLayout.CallFrameAlignOfType", t)
	return binding.LLVMCallFrameAlignmentOfType(d.ref, t.Ref())
}

// PrefAlignOfType 编译器推荐对齐（byte）
func (d *DataLayout) PrefAlignOfType(t AnyType) uint32 {
	d.checkType("llvm.DataLayout.PrefAlignOfType", t)
	return binding.LLVMPreferredAlignmentOfType(d.ref, t.Ref())
}

// PrefAlignOfGlobal 全局变量推荐对齐（byte）
func (d *DataLayout) PrefAlignOfGlobal(g AnyValue) uint32 {
	d.checkValue("llvm.DataLayout.PrefAlignOfGlobal", g)
	return binding.LLVMPreferredAlignmentOfGlobal(d.ref, g.Ref())
}

// ElementAtOffset 包含指定字节偏移的结构体元素下标
func (d *DataLayout) ElementAtOffset(st AnyType, offset uint64) uint32 {
	d.checkType("llvm.DataLayout.ElementAtOffset", st)
	return binding.LLVMElementAtOffset(d.ref, st.Ref(), offset)
}

// OffsetOfElement 指定结构体元素的字节偏移
func (d *DataLayout) OffsetOfElement(st AnyType, i uint32) uint64 {
	d.checkType("llvm.DataLayout.OffsetOfElement", st)
	return binding.LLVMOffsetOfElement(d.ref, st.Ref(), i)
}
