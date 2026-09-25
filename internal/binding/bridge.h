#ifndef GOLLVM_BINDINGS_BRIDGE_H
#define GOLLVM_BINDINGS_BRIDGE_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

// goLLVMBridgeDispatch is the Go function exported via //export; the C-side channel forwards to it.
extern uint64_t goLLVMBridgeDispatch(int64_t idx, uint64_t *slots);

// llvmBridgeCall calls a JIT-generated adapter with the fixed C ABI:
//   uint64_t adapter(void *fn, uint64_t *slots)
// fn is the target function address, slots[0..] are the argument slots, and the
// returned value is the adapter's return value.
uint64_t llvmBridgeCall(void *adapter, void *fn, uint64_t *slots);

// llvmBridgeGoChannel is the fixed-signature channel called by JIT code:
//   uint64_t callGoChannel(int64_t idx, uint64_t *slots)
// It forwards to goLLVMBridgeDispatch.
uint64_t llvmBridgeGoChannel(int64_t idx, uint64_t *slots);

// llvmBridgeGoChannelAddr returns the function address of llvmBridgeGoChannel,
// for ORC DefineAbsoluteSymbols to hook up the callGoChannel symbol.
void *llvmBridgeGoChannelAddr(void);

#ifdef __cplusplus
}
#endif

#endif
