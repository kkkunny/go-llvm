package llvm

import (
	"runtime"

	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/errs"
)

// ByteOrder 目标字节序
type ByteOrder int32

// ByteOrder 取值对应 LLVM 目标字节序（binding.LLVMByteOrdering）。
// 决定多字节整数与浮点在内存中的字节排列顺序。
const (
	LittleEndian ByteOrder = ByteOrder(binding.LLVMLittleEndian) // 小端：低位字节在低地址
	BigEndian    ByteOrder = ByteOrder(binding.LLVMBigEndian)    // 大端：高位字节在低地址
)

// DataLayout 目标数据布局查询句柄；由 NewDataLayout 或 TargetMachine 创建，用毕 Close
type DataLayout struct {
	ref    binding.LLVMTargetDataRef
	closed bool
}

// NewDataLayout 由布局字符串创建数据布局（拥有句柄）。
//
// 空字符串得到的是 LLVM 默认布局（例如 i64 ABI 对齐为 4），与宿主目标无关；
// 需要宿主目标的真实布局请用
// [github.com/kkkunny/go-llvm/target.TargetMachine.DataLayout]。
// 非法布局字符串 panic [ErrInvalidArg]（经 DataLayout::parse 解析，不会触发 LLVM
// report_fatal_error 退出进程）。
func NewDataLayout(layout string) *DataLayout {
	ref, errMsg := binding.LLVMGoCreateTargetData(layout)
	if !ref.IsNil() {
		d := &DataLayout{ref: ref}
		runtime.SetFinalizer(d, (*DataLayout).finalize)
		return d
	}
	if errMsg == "" {
		errMsg = "unknown error"
	}
	errs.Panicf(ErrInvalidArg, "llvm.NewDataLayout", "invalid data layout string %q: %s", layout, errMsg)
	return nil
}

// DataLayoutOf 由底层句柄构建数据布局（供 llvm/target 桥接使用，接管所有权）
func DataLayoutOf(ref binding.LLVMTargetDataRef) *DataLayout {
	d := &DataLayout{ref: ref}
	runtime.SetFinalizer(d, (*DataLayout).finalize)
	return d
}

// finalize GC 兜底：忘记 Close 时释放底层句柄
func (d *DataLayout) finalize() {
	if d.closed {
		return
	}
	d.closed = true
	binding.LLVMDisposeTargetData(d.ref)
}

// check 前置校验：句柄非空且未释放
func (d *DataLayout) check(op string) {
	if d.ref.IsNil() {
		errs.Panicf(ErrInvalidArg, op, "nil data layout")
	}
	if d.closed {
		errs.Panicf(ErrUseAfterFree, op, "data layout is closed")
	}
}

// checkType 前置校验类型参数
func (d *DataLayout) checkType(op string, t AnyType) {
	d.check(op)
	if t == nil || t.IsNil() {
		errs.Panicf(ErrInvalidArg, op, "nil type")
	}
}

// checkValue 前置校验值参数
func (d *DataLayout) checkValue(op string, v AnyValue) {
	d.check(op)
	if v == nil || v.IsNil() {
		errs.Panicf(ErrInvalidArg, op, "nil value")
	}
}

// Close 释放句柄；二次调用返回 ErrClosed
func (d *DataLayout) Close() error {
	if d.closed {
		return &Error{Reason: ErrClosed, Op: "llvm.DataLayout.Close", Msg: "data layout already closed"}
	}
	d.closed = true
	runtime.SetFinalizer(d, nil)
	binding.LLVMDisposeTargetData(d.ref)
	return nil
}

// String 布局字符串
func (d *DataLayout) String() string {
	d.check("llvm.DataLayout.String")
	s := binding.LLVMCopyStringRepOfTargetData(d.ref)
	runtime.KeepAlive(d)
	return s
}

// ByteOrder 字节序
func (d *DataLayout) ByteOrder() ByteOrder {
	d.check("llvm.DataLayout.ByteOrder")
	v := ByteOrder(binding.LLVMByteOrder(d.ref))
	runtime.KeepAlive(d)
	return v
}

// PointerSize 指针字节数
func (d *DataLayout) PointerSize() uint32 {
	d.check("llvm.DataLayout.PointerSize")
	v := binding.LLVMPointerSize(d.ref)
	runtime.KeepAlive(d)
	return v
}

// PointerSizeForAS 指定地址空间的指针字节数
func (d *DataLayout) PointerSizeForAS(addrspace uint32) uint32 {
	d.check("llvm.DataLayout.PointerSizeForAS")
	v := binding.LLVMPointerSizeForAS(d.ref, addrspace)
	runtime.KeepAlive(d)
	return v
}

// IntPtrType 指针等宽整数类型
func (d *DataLayout) IntPtrType(ctx *Context) IntType {
	d.check("llvm.DataLayout.IntPtrType")
	ctx.CheckAlive("llvm.DataLayout.IntPtrType")
	t := IntType{Type[IntT]{ref: binding.LLVMIntPtrTypeInContext(ctx.ref, d.ref), ctx: ctx}}
	runtime.KeepAlive(d)
	return t
}

// IntPtrTypeForAS 指定地址空间指针等宽整数类型
func (d *DataLayout) IntPtrTypeForAS(ctx *Context, addrspace uint32) IntType {
	d.check("llvm.DataLayout.IntPtrTypeForAS")
	ctx.CheckAlive("llvm.DataLayout.IntPtrTypeForAS")
	t := IntType{Type[IntT]{ref: binding.LLVMIntPtrTypeForASInContext(ctx.ref, d.ref, addrspace), ctx: ctx}}
	runtime.KeepAlive(d)
	return t
}

// SizeOfTypeInBits 类型位宽（bit）
func (d *DataLayout) SizeOfTypeInBits(t AnyType) uint64 {
	d.checkType("llvm.DataLayout.SizeOfTypeInBits", t)
	v := binding.LLVMSizeOfTypeInBits(d.ref, t.Ref())
	runtime.KeepAlive(d)
	return v
}

// StoreSizeOfType 实际存储大小（byte）
func (d *DataLayout) StoreSizeOfType(t AnyType) uint64 {
	d.checkType("llvm.DataLayout.StoreSizeOfType", t)
	v := binding.LLVMStoreSizeOfType(d.ref, t.Ref())
	runtime.KeepAlive(d)
	return v
}

// ABISizeOfType ABI 大小（byte）
func (d *DataLayout) ABISizeOfType(t AnyType) uint64 {
	d.checkType("llvm.DataLayout.ABISizeOfType", t)
	v := binding.LLVMABISizeOfType(d.ref, t.Ref())
	runtime.KeepAlive(d)
	return v
}

// ABIAlignOfType ABI 对齐（byte）
func (d *DataLayout) ABIAlignOfType(t AnyType) uint32 {
	d.checkType("llvm.DataLayout.ABIAlignOfType", t)
	v := binding.LLVMABIAlignmentOfType(d.ref, t.Ref())
	runtime.KeepAlive(d)
	return v
}

// CallFrameAlignOfType 调用栈帧对齐（byte）
func (d *DataLayout) CallFrameAlignOfType(t AnyType) uint32 {
	d.checkType("llvm.DataLayout.CallFrameAlignOfType", t)
	v := binding.LLVMCallFrameAlignmentOfType(d.ref, t.Ref())
	runtime.KeepAlive(d)
	return v
}

// PrefAlignOfType 编译器推荐对齐（byte）
func (d *DataLayout) PrefAlignOfType(t AnyType) uint32 {
	d.checkType("llvm.DataLayout.PrefAlignOfType", t)
	v := binding.LLVMPreferredAlignmentOfType(d.ref, t.Ref())
	runtime.KeepAlive(d)
	return v
}

// PrefAlignOfGlobal 全局变量推荐对齐（byte）
func (d *DataLayout) PrefAlignOfGlobal(g AnyValue) uint32 {
	d.checkValue("llvm.DataLayout.PrefAlignOfGlobal", g)
	v := binding.LLVMPreferredAlignmentOfGlobal(d.ref, g.Ref())
	runtime.KeepAlive(d)
	return v
}

// ElementAtOffset 包含指定字节偏移的结构体元素下标
func (d *DataLayout) ElementAtOffset(st AnyType, offset uint64) uint32 {
	d.checkType("llvm.DataLayout.ElementAtOffset", st)
	v := binding.LLVMElementAtOffset(d.ref, st.Ref(), offset)
	runtime.KeepAlive(d)
	return v
}

// OffsetOfElement 指定结构体元素的字节偏移
func (d *DataLayout) OffsetOfElement(st AnyType, i uint32) uint64 {
	d.checkType("llvm.DataLayout.OffsetOfElement", st)
	v := binding.LLVMOffsetOfElement(d.ref, st.Ref(), i)
	runtime.KeepAlive(d)
	return v
}
