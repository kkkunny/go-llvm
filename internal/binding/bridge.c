#include "bridge.h"

typedef uint64_t (*goLLVMBridgeAdapter)(void *, uint64_t *);

uint64_t llvmBridgeCall(void *adapter, void *fn, uint64_t *slots) {
    return ((goLLVMBridgeAdapter)adapter)(fn, slots);
}

uint64_t llvmBridgeGoChannel(int64_t idx, uint64_t *slots) {
    return goLLVMBridgeDispatch(idx, slots);
}

void *llvmBridgeGoChannelAddr(void) {
    return (void *)&llvmBridgeGoChannel;
}
