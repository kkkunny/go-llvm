#include "execution.h"

extern void goFuncCallChannel(int funcIdx, void **ret_ptr_ptr, void **param_ptr_array_ptr);

void callGoFunc(int funcIdx, void **ret_ptr_ptr, void **param_ptr_array_ptr){
    goFuncCallChannel(funcIdx, ret_ptr_ptr, param_ptr_array_ptr);
}
