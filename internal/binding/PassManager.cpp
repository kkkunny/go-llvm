#include "PassManager.h"
#include <llvm/Passes/PassBuilder.h>
#include <cstdlib>
#include <cstring>
#include <exception>
#include <string>

using namespace llvm;

static OptimizationLevel ParseOptimizationLevel(const char *level){
    std::string s(level);
    if (s == "O0"){
        return OptimizationLevel::O0;
    } else if (s == "O1"){
        return OptimizationLevel::O1;
    } else if (s == "O2"){
        return OptimizationLevel::O2;
    } else if (s == "O3"){
        return OptimizationLevel::O3;
    } else if (s == "Oz"){
        return OptimizationLevel::Oz;
    } else if (s == "Os"){
        return OptimizationLevel::Os;
    } else {
        return OptimizationLevel();
    }
}

static char *CopyError(const char *msg){
    size_t len = strlen(msg) + 1;
    char *out = static_cast<char *>(malloc(len));
    if (out != nullptr){
        memcpy(out, msg, len);
    }
    return out;
}

int LLVMOptModule(LLVMModuleRef IR, const char *level, char **outError){
    try {
        ModuleAnalysisManager MAM;
        LoopAnalysisManager LAM;
        FunctionAnalysisManager FAM;
        CGSCCAnalysisManager CGAM;

        PassBuilder PB;

        PB.registerModuleAnalyses(MAM);
        PB.registerCGSCCAnalyses(CGAM);
        PB.registerFunctionAnalyses(FAM);
        PB.registerLoopAnalyses(LAM);
        PB.crossRegisterProxies(LAM, FAM, CGAM, MAM);

        auto MPM = PB.buildPerModuleDefaultPipeline(ParseOptimizationLevel(level));
        MPM.run(*unwrap(IR), MAM);
        return 0;
    } catch (const std::exception &e){
        if (outError != nullptr){
            *outError = CopyError(e.what());
        }
        return 1;
    } catch (...){
        if (outError != nullptr){
            *outError = CopyError("unknown C++ exception in LLVMOptModule");
        }
        return 1;
    }
}
