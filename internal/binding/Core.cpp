#include "Core.h"
#include "llvm/IR/Constants.h"
#include "llvm/IR/DataLayout.h"
#include "llvm/IR/Function.h"
#include "llvm/IR/GlobalValue.h"
#include "llvm/IR/LLVMContext.h"
#include "llvm/Support/Error.h"
#include <cstring>

using namespace llvm;

// reportCurrentException emits the in-flight exception as a diagnostic and returns.
// Must only be called from inside a catch block (it rethrows to inspect the type).
static void reportCurrentException(LLVMContext *C) {
    try {
        throw;
    } catch (const std::exception &E) {
        if (C)
            C->emitError(E.what());
    } catch (...) {
        if (C)
            C->emitError("unknown C++ exception in a go-llvm shim");
    }
}

LLVMValueRef LLVMConstAggregateZero(LLVMTypeRef ty) {
    Type *Ty = unwrap(ty);
    try {
        return wrap(ConstantAggregateZero::get(Ty));
    } catch (...) {
        reportCurrentException(&Ty->getContext());
        return nullptr;
    }
}

LLVMTypeRef LLVMGetFunctionType(LLVMValueRef f) {
    Function *F = unwrap<Function>(f);
    try {
        return wrap(F->getFunctionType());
    } catch (...) {
        reportCurrentException(&F->getContext());
        return nullptr;
    }
}

void LLVMSetDSOLocal(LLVMValueRef v, LLVMBool Local) {
    GlobalValue *GV = unwrap<GlobalValue>(v);
    try {
        GV->setDSOLocal(Local);
    } catch (...) {
        reportCurrentException(&GV->getContext());
    }
}

LLVMBool LLVMIsDSOLocal(LLVMValueRef v) {
    GlobalValue *GV = unwrap<GlobalValue>(v);
    try {
        return GV->isDSOLocal();
    } catch (...) {
        reportCurrentException(&GV->getContext());
        return 0;
    }
}

void LLVMGoEmitError(LLVMContextRef c, const char *msg) {
    LLVMContext *C = unwrap(c);
    try {
        C->emitError(msg);
    } catch (...) {
        // Nothing sensible to do: emitting the error is already the fallback path.
    }
}

LLVMValueRef LLVMGoConstSExt(LLVMValueRef v, LLVMTypeRef ty) {
    Type *Ty = unwrap(ty);
    // LLVM 23 dropped sext/zext constant expressions, so only plain ConstantInt
    // operands can be extended; anything else returns null and the Go side panics.
    if (auto *CI = dyn_cast<ConstantInt>(unwrap(v)))
        return wrap(ConstantInt::get(Ty, CI->getValue().sext(Ty->getIntegerBitWidth())));
    return nullptr;
}

LLVMValueRef LLVMGoConstZExt(LLVMValueRef v, LLVMTypeRef ty) {
    Type *Ty = unwrap(ty);
    if (auto *CI = dyn_cast<ConstantInt>(unwrap(v)))
        return wrap(ConstantInt::get(Ty, CI->getValue().zext(Ty->getIntegerBitWidth())));
    return nullptr;
}

LLVMTargetDataRef LLVMGoCreateTargetData(const char *StringRep, char **ErrMsg) {
    if (ErrMsg)
        *ErrMsg = nullptr;
    Expected<DataLayout> Parsed = DataLayout::parse(StringRep);
    if (!Parsed) {
        if (ErrMsg) {
            std::string Msg = toString(Parsed.takeError());
            *ErrMsg = strdup(Msg.c_str());
        }
        return nullptr;
    }
    return wrap(new DataLayout(*Parsed));
}

int LLVMGoConstIntGetSExtValue(LLVMValueRef v, long long *out) {
    auto *CI = dyn_cast<ConstantInt>(unwrap(v));
    if (!CI)
        return 0;
    const APInt &V = CI->getValue();
    if (!V.isSignedIntN(64))
        return 0;
    *out = V.getSExtValue();
    return 1;
}

int LLVMGoConstIntGetZExtValue(LLVMValueRef v, unsigned long long *out) {
    auto *CI = dyn_cast<ConstantInt>(unwrap(v));
    if (!CI)
        return 0;
    const APInt &V = CI->getValue();
    if (V.getActiveBits() > 64)
        return 0;
    *out = V.getZExtValue();
    return 1;
}
