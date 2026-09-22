#ifndef GOLLVM_BINDINGS_BRIDGE_H
#define GOLLVM_BINDINGS_BRIDGE_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

// goLLVMBridgeDispatch 是 //export 导出的 Go 函数，由 C 侧通道转发
extern uint64_t goLLVMBridgeDispatch(int64_t idx, uint64_t *slots);

// llvmBridgeCall 以固定 C ABI 调用 JIT 生成的适配器：
//   uint64_t adapter(void *fn, uint64_t *slots)
// fn 为目标函数地址，slots[0..] 为实参槽，返回值即适配器返回
uint64_t llvmBridgeCall(void *adapter, void *fn, uint64_t *slots);

// llvmBridgeGoChannel 是 JIT 代码调用的固定签名通道：
//   uint64_t callGoChannel(int64_t idx, uint64_t *slots)
// 转发到 goLLVMBridgeDispatch
uint64_t llvmBridgeGoChannel(int64_t idx, uint64_t *slots);

// llvmBridgeGoChannelAddr 返回 llvmBridgeGoChannel 的函数地址，
// 供 ORC DefineAbsoluteSymbols 挂接 callGoChannel 符号
void *llvmBridgeGoChannelAddr(void);

#ifdef __cplusplus
}
#endif

#endif
