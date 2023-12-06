package binding

/*
#include "llvm-c/Transforms/PassBuilder.h"
*/
import "C"
import "errors"

// LLVMPassBuilderOptionsRef A set of options passed which are attached to the Pass Manager upon run.
// This corresponds to an llvm::LLVMPassBuilderOptions instance
// The details for how the different properties of this structure are used can be found in the source for LLVMRunPasses
type LLVMPassBuilderOptionsRef struct{ c C.LLVMPassBuilderOptionsRef }

// LLVMRunPasses Construct and run a set of passes over a module
// This function takes a string with the passes that should be used.
// The format of this string is the same as opt's -passes argument for the new pass manager.
// Individual passes may be specified, separated by commas.
// Full pipelines may also be invoked using `default<O3>` and friends.
// See opt for full reference of the Passes format.
func LLVMRunPasses(m LLVMModuleRef, passes string, tm LLVMTargetMachineRef, options LLVMPassBuilderOptionsRef) error {
	err := string2CString(passes, func(passes *C.char) LLVMErrorRef {
		return LLVMErrorRef{c: C.LLVMRunPasses(m.c, passes, tm.c, options.c)}
	})
	if err.c != nil {
		return errors.New(LLVMGetErrorMessage(err))
	}
	return nil
}

// LLVMCreatePassBuilderOptions Create a new set of options for a PassBuilder
// Ownership of the returned instance is given to the client, and they are responsible for it.
// The client should call LLVMDisposePassBuilderOptions to free the pass builder options.
func LLVMCreatePassBuilderOptions() LLVMPassBuilderOptionsRef {
	return LLVMPassBuilderOptionsRef{c: C.LLVMCreatePassBuilderOptions()}
}

// LLVMPassBuilderOptionsSetVerifyEach Toggle adding the VerifierPass for the PassBuilder, ensuring all functions inside the module is valid.
func LLVMPassBuilderOptionsSetVerifyEach(options LLVMPassBuilderOptionsRef, verifyEach bool) {
	C.LLVMPassBuilderOptionsSetVerifyEach(options.c, bool2LLVMBool(verifyEach))
}

// LLVMPassBuilderOptionsSetDebugLogging Toggle debug logging when running the PassBuilder
func LLVMPassBuilderOptionsSetDebugLogging(options LLVMPassBuilderOptionsRef, debugLogging bool) {
	C.LLVMPassBuilderOptionsSetDebugLogging(options.c, bool2LLVMBool(debugLogging))
}

func LLVMPassBuilderOptionsSetLoopInterleaving(options LLVMPassBuilderOptionsRef, loopInterleaving bool) {
	C.LLVMPassBuilderOptionsSetLoopInterleaving(options.c, bool2LLVMBool(loopInterleaving))
}

func LLVMPassBuilderOptionsSetLoopVectorization(options LLVMPassBuilderOptionsRef, loopVectorization bool) {
	C.LLVMPassBuilderOptionsSetLoopVectorization(options.c, bool2LLVMBool(loopVectorization))
}

func LLVMPassBuilderOptionsSetSLPVectorization(options LLVMPassBuilderOptionsRef, slpVectorization bool) {
	C.LLVMPassBuilderOptionsSetSLPVectorization(options.c, bool2LLVMBool(slpVectorization))
}

func LLVMPassBuilderOptionsSetLoopUnrolling(options LLVMPassBuilderOptionsRef, loopUnrolling bool) {
	C.LLVMPassBuilderOptionsSetLoopUnrolling(options.c, bool2LLVMBool(loopUnrolling))
}

func LLVMPassBuilderOptionsSetForgetAllSCEVInLoopUnroll(options LLVMPassBuilderOptionsRef, forgetAllSCEVInLoopUnroll bool) {
	C.LLVMPassBuilderOptionsSetForgetAllSCEVInLoopUnroll(options.c, bool2LLVMBool(forgetAllSCEVInLoopUnroll))
}

func LLVMPassBuilderOptionsSetLicmMssaOptCap(options LLVMPassBuilderOptionsRef, licmMssaOptCap uint32) {
	C.LLVMPassBuilderOptionsSetLicmMssaOptCap(options.c, C.unsigned(licmMssaOptCap))
}

func LLVMPassBuilderOptionsSetLicmMssaNoAccForPromotionCap(options LLVMPassBuilderOptionsRef, licmMssaNoAccForPromotionCap uint32) {
	C.LLVMPassBuilderOptionsSetLicmMssaNoAccForPromotionCap(options.c, C.unsigned(licmMssaNoAccForPromotionCap))
}

func LLVMPassBuilderOptionsSetCallGraphProfile(options LLVMPassBuilderOptionsRef, callGraphProfile bool) {
	C.LLVMPassBuilderOptionsSetCallGraphProfile(options.c, bool2LLVMBool(callGraphProfile))
}

func LLVMPassBuilderOptionsSetMergeFunctions(options LLVMPassBuilderOptionsRef, mergeFunctions bool) {
	C.LLVMPassBuilderOptionsSetMergeFunctions(options.c, bool2LLVMBool(mergeFunctions))
}

// LLVMDisposePassBuilderOptions Dispose of a heap-allocated PassBuilderOptions instance
func LLVMDisposePassBuilderOptions(options LLVMPassBuilderOptionsRef) {
	C.LLVMDisposePassBuilderOptions(options.c)
}
