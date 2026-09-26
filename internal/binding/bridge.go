package binding

/*
#include "bridge.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// BridgeMaxSlots is the maximum number of slots in the fixed-signature channel (8 bytes per slot;
// the channel function's return value is passed back as the call's return).
const BridgeMaxSlots = 16

var bridgeDispatch func(idx int64, slots []uint64) uint64

//export goLLVMBridgeDispatch
func goLLVMBridgeDispatch(idx C.int64_t, slots *C.uint64_t) C.uint64_t {
	if bridgeDispatch == nil {
		return 0
	}
	return C.uint64_t(bridgeDispatch(int64(idx), unsafe.Slice((*uint64)(unsafe.Pointer(slots)), BridgeMaxSlots)))
}

// SetBridgeDispatch registers the Go-side dispatcher (used by the llvm/jit bridge).
func SetBridgeDispatch(fn func(idx int64, slots []uint64) uint64) {
	bridgeDispatch = fn
}

// BridgeCall calls a JIT-generated adapter with the fixed C ABI (slots must be at least BridgeMaxSlots long).
func BridgeCall(adapter, fn unsafe.Pointer, slots []uint64) uint64 {
	var ptr *C.uint64_t
	if len(slots) > 0 {
		ptr = (*C.uint64_t)(unsafe.Pointer(&slots[0]))
	}
	return uint64(C.llvmBridgeCall(adapter, fn, ptr))
}

// BridgeGoChannelAddr returns the C function address of the callGoChannel channel.
func BridgeGoChannelAddr() unsafe.Pointer {
	return C.llvmBridgeGoChannelAddr()
}

// CStringArray is an array of C strings (NULL-terminated); Free it when done.
type CStringArray struct {
	ptrs []*C.char
}

// NewCStringArray builds a NULL-terminated array of C strings.
func NewCStringArray(strs []string) *CStringArray {
	ptrs := make([]*C.char, len(strs)+1)
	for i, s := range strs {
		ptrs[i] = C.CString(s)
	}
	return &CStringArray{ptrs: ptrs}
}

// Ptr returns the array head pointer (NULL-terminated).
func (a *CStringArray) Ptr() **C.char {
	if len(a.ptrs) == 0 {
		return nil
	}
	return (**C.char)(unsafe.Pointer(&a.ptrs[0]))
}

// Free frees the array and the strings it contains.
func (a *CStringArray) Free() {
	for _, p := range a.ptrs {
		if p != nil {
			C.free(unsafe.Pointer(p))
		}
	}
	a.ptrs = nil
}
