#include "PassManager.h"
#include <llvm/Passes/PassBuilder.h>
#include <string>

using namespace llvm;

DEFINE_SIMPLE_CONVERSION_FUNCTIONS(Module, LLVMModuleRef)

OptimizationLevel ParseOptimizationLevel(const char *level){
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

void LLVMOptModule(LLVMModuleRef IR, const char *level){
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

    auto MPM = PB.buildPerModuleDefaultPipeline(ParseOptimizationLevel(level), false);
    MPM.run(*unwrap(IR), MAM);
}
