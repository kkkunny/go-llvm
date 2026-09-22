#ifndef GOLLVM_BINDINGS_PASSMANAGER_H
#define GOLLVM_BINDINGS_PASSMANAGER_H

#include <llvm-c/Types.h>
#ifdef __cplusplus
#include "llvm/Support/CBindingWrapping.h"
#endif

#ifdef __cplusplus
extern "C" {
#endif

// LLVMOptModule 运行默认优化管线；成功返回 0，失败返回 1 并通过 outError 返回 malloc 的错误消息
int LLVMOptModule(LLVMModuleRef IR, const char *level, char **outError);

#ifdef __cplusplus
}
#endif

#endif
