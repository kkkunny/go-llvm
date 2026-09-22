package binding

/*
#include "bridge.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// BridgeMaxSlots 固定签名通道的最大槽数（每槽 8 字节；返回值为通道函数返回）
const BridgeMaxSlots = 16

var bridgeDispatch func(idx int64, slots []uint64) uint64

//export goLLVMBridgeDispatch
func goLLVMBridgeDispatch(idx C.int64_t, slots *C.uint64_t) C.uint64_t {
	if bridgeDispatch == nil {
		return 0
	}
	return C.uint64_t(bridgeDispatch(int64(idx), unsafe.Slice((*uint64)(unsafe.Pointer(slots)), BridgeMaxSlots)))
}

// SetBridgeDispatch 注册 Go 侧派发器（供 llvm/jit 桥接使用）
func SetBridgeDispatch(fn func(idx int64, slots []uint64) uint64) {
	bridgeDispatch = fn
}

// BridgeCall 以固定 C ABI 调用 JIT 生成的适配器（slots 长度至少 BridgeMaxSlots）
func BridgeCall(adapter, fn unsafe.Pointer, slots []uint64) uint64 {
	var ptr *C.uint64_t
	if len(slots) > 0 {
		ptr = (*C.uint64_t)(unsafe.Pointer(&slots[0]))
	}
	return uint64(C.llvmBridgeCall(adapter, fn, ptr))
}

// BridgeGoChannelAddr 返回 callGoChannel 通道的 C 函数地址
func BridgeGoChannelAddr() unsafe.Pointer {
	return C.llvmBridgeGoChannelAddr()
}

// CStringArray C 字符串数组（NULL 结尾）；用毕 Free
type CStringArray struct {
	ptrs []*C.char
}

// NewCStringArray 构造以 NULL 结尾的 C 字符串数组
func NewCStringArray(strs []string) *CStringArray {
	ptrs := make([]*C.char, len(strs)+1)
	for i, s := range strs {
		ptrs[i] = C.CString(s)
	}
	return &CStringArray{ptrs: ptrs}
}

// Ptr 返回数组头指针（NULL 结尾）
func (a *CStringArray) Ptr() **C.char {
	if len(a.ptrs) == 0 {
		return nil
	}
	return (**C.char)(unsafe.Pointer(&a.ptrs[0]))
}

// Free 释放数组及其中字符串
func (a *CStringArray) Free() {
	for _, p := range a.ptrs {
		if p != nil {
			C.free(unsafe.Pointer(p))
		}
	}
	a.ptrs = nil
}
