package binding

/*
#include "llvm-c/ExecutionEngine.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

type (
	LLVMGenericValueRef       struct{ c C.LLVMGenericValueRef }
	LLVMExecutionEngineRef    struct{ c C.LLVMExecutionEngineRef }
	LLVMMCJITMemoryManagerRef struct{ c C.LLVMMCJITMemoryManagerRef }
	LLVMMCJITCompilerOptions  struct {
		OptLevel           LLVMCodeGenOptLevel
		CodeModel          LLVMCodeModel
		NoFramePointerElim bool
		EnableFastISel     bool
		MCJMM              LLVMMCJITMemoryManagerRef
	}
)

func (options *LLVMMCJITCompilerOptions) fromC(c C.struct_LLVMMCJITCompilerOptions) {
	options.OptLevel = LLVMCodeGenOptLevel(c.OptLevel)
	options.CodeModel = LLVMCodeModel(c.CodeModel)
	options.NoFramePointerElim = llvmBool2bool(c.NoFramePointerElim)
	options.EnableFastISel = llvmBool2bool(c.EnableFastISel)
	options.MCJMM = LLVMMCJITMemoryManagerRef{c: c.MCJMM}
}

func (options LLVMMCJITCompilerOptions) c() C.struct_LLVMMCJITCompilerOptions {
	var c C.struct_LLVMMCJITCompilerOptions
	c.OptLevel = C.unsigned(options.OptLevel)
	c.CodeModel = C.LLVMCodeModel(options.CodeModel)
	c.NoFramePointerElim = bool2LLVMBool(options.NoFramePointerElim)
	c.EnableFastISel = bool2LLVMBool(options.EnableFastISel)
	c.MCJMM = options.MCJMM.c
	return c
}

func LLVMCreateGenericValueOfInt(ty LLVMTypeRef, n uint64, isSigned bool) LLVMGenericValueRef {
	return LLVMGenericValueRef{c: C.LLVMCreateGenericValueOfInt(ty.c, C.ulonglong(n), bool2LLVMBool(isSigned))}
}

func LLVMCreateGenericValueOfPointer(p unsafe.Pointer) LLVMGenericValueRef {
	return LLVMGenericValueRef{c: C.LLVMCreateGenericValueOfPointer(p)}
}

func LLVMCreateGenericValueOfFloat(ty LLVMTypeRef, n float64) LLVMGenericValueRef {
	return LLVMGenericValueRef{c: C.LLVMCreateGenericValueOfFloat(ty.c, C.double(n))}
}

func LLVMGenericValueIntWidth(genValRef LLVMGenericValueRef) uint32 {
	return uint32(C.LLVMGenericValueIntWidth(genValRef.c))
}

func LLVMGenericValueToInt(genVal LLVMGenericValueRef, isSigned bool) uint64 {
	return uint64(C.LLVMGenericValueToInt(genVal.c, bool2LLVMBool(isSigned)))
}

func LLVMGenericValueToPointer(genVal LLVMGenericValueRef) unsafe.Pointer {
	return C.LLVMGenericValueToPointer(genVal.c)
}

func LLVMGenericValueToFloat(tyRef LLVMTypeRef, genVal LLVMGenericValueRef) float64 {
	return float64(C.LLVMGenericValueToFloat(tyRef.c, genVal.c))
}

func LLVMDisposeGenericValue(genVal LLVMGenericValueRef) {
	C.LLVMDisposeGenericValue(genVal.c)
}

func LLVMCreateExecutionEngineForModule(m LLVMModuleRef) (LLVMExecutionEngineRef, error) {
	var outee LLVMExecutionEngineRef
	err := llvmError2Error(func(outError **C.char) C.LLVMBool {
		return C.LLVMCreateExecutionEngineForModule(&outee.c, m.c, outError)
	})
	return outee, err
}

func LLVMCreateInterpreterForModule(m LLVMModuleRef) (LLVMExecutionEngineRef, error) {
	var outInterp LLVMExecutionEngineRef
	err := llvmError2Error(func(outError **C.char) C.LLVMBool {
		return C.LLVMCreateInterpreterForModule(&outInterp.c, m.c, outError)
	})
	return outInterp, err
}

func LLVMCreateJITCompilerForModule(m LLVMModuleRef, optLevel uint32) (LLVMExecutionEngineRef, error) {
	var outJIT LLVMExecutionEngineRef
	err := llvmError2Error(func(outError **C.char) C.LLVMBool {
		return C.LLVMCreateJITCompilerForModule(&outJIT.c, m.c, C.unsigned(optLevel), outError)
	})
	return outJIT, err
}

func LLVMInitializeMCJITCompilerOptions(options *LLVMMCJITCompilerOptions) {
	coptions := options.c()
	C.LLVMInitializeMCJITCompilerOptions(&coptions, C.size_t(unsafe.Sizeof(C.struct_LLVMMCJITCompilerOptions{})))
	options.fromC(coptions)
}

// LLVMCreateMCJITCompilerForModule Create an MCJIT execution engine for a module, with the given options.
// It is the responsibility of the caller to ensure that all fields in Options up to the given SizeOfOptions are initialized.
// It is correct to pass a smaller value of SizeOfOptions that omits some fields.
// The canonical way of using this is:
//
// LLVMMCJITCompilerOptions options;
// LLVMInitializeMCJITCompilerOptions(&options, sizeof(options));
// ... fill in those options you care about
// LLVMCreateMCJITCompilerForModule(&jit, mod, &options, sizeof(options), &error);
//
// Note that this is also correct, though possibly suboptimal:
//
// LLVMCreateMCJITCompilerForModule(&jit, mod, 0, 0, &error);
func LLVMCreateMCJITCompilerForModule(m LLVMModuleRef, options LLVMMCJITCompilerOptions) (LLVMExecutionEngineRef, error) {
	var outJIT LLVMExecutionEngineRef
	coptions := options.c()
	err := llvmError2Error(func(outError **C.char) C.LLVMBool {
		return C.LLVMCreateMCJITCompilerForModule(&outJIT.c, m.c, &coptions, C.size_t(unsafe.Sizeof(C.struct_LLVMMCJITCompilerOptions{})), outError)
	})
	return outJIT, err
}

func LLVMDisposeExecutionEngine(ee LLVMExecutionEngineRef) {
	C.LLVMDisposeExecutionEngine(ee.c)
}

func LLVMRunStaticConstructors(ee LLVMExecutionEngineRef) {
	C.LLVMRunStaticConstructors(ee.c)
}

func LLVMRunStaticDestructors(ee LLVMExecutionEngineRef) {
	C.LLVMRunStaticDestructors(ee.c)
}

func LLVMRunFunctionAsMain(ee LLVMExecutionEngineRef, f LLVMValueRef, argc uint32, argv *string, envp *string) int32 {
	var argvp, envpp **C.char
	if argv != nil {
		cargv := C.CString(*argv)
		defer C.free(unsafe.Pointer(cargv))
		argvp = &cargv
	}
	if envp != nil {
		cenvp := C.CString(*envp)
		defer C.free(unsafe.Pointer(cenvp))
		envpp = &cenvp
	}
	return int32(C.LLVMRunFunctionAsMain(ee.c, f.c, C.unsigned(argc), argvp, envpp))
}

func LLVMRunFunction(ee LLVMExecutionEngineRef, f LLVMValueRef, args []LLVMGenericValueRef) LLVMGenericValueRef {
	ptr, length := slice2Ptr[LLVMGenericValueRef, C.LLVMGenericValueRef](args)
	return LLVMGenericValueRef{c: C.LLVMRunFunction(ee.c, f.c, length, ptr)}
}

func LLVMFreeMachineCodeForFunction(ee LLVMExecutionEngineRef, f LLVMValueRef) {
	C.LLVMFreeMachineCodeForFunction(ee.c, f.c)
}

func LLVMAddModule(ee LLVMExecutionEngineRef, m LLVMModuleRef) {
	C.LLVMAddModule(ee.c, m.c)
}

func LLVMRemoveModule(ee LLVMExecutionEngineRef, m LLVMModuleRef) (LLVMModuleRef, error) {
	var outMod LLVMModuleRef
	err := llvmError2Error(func(outError **C.char) C.LLVMBool {
		return C.LLVMRemoveModule(ee.c, m.c, &outMod.c, outError)
	})
	return outMod, err
}

func LLVMFindFunction(ee LLVMExecutionEngineRef, name string) (LLVMValueRef, bool) {
	var outFn LLVMValueRef
	v := string2CString(name, func(name *C.char) bool {
		return llvmBool2bool(C.LLVMFindFunction(ee.c, name, &outFn.c))
	})
	return outFn, v
}

func LLVMRecompileAndRelinkFunction(ee LLVMExecutionEngineRef, fn LLVMValueRef) unsafe.Pointer {
	return C.LLVMRecompileAndRelinkFunction(ee.c, fn.c)
}

func LLVMGetExecutionEngineTargetData(ee LLVMExecutionEngineRef) LLVMTargetDataRef {
	return LLVMTargetDataRef{c: C.LLVMGetExecutionEngineTargetData(ee.c)}
}

func LLVMGetExecutionEngineTargetMachine(ee LLVMExecutionEngineRef) LLVMTargetMachineRef {
	return LLVMTargetMachineRef{c: C.LLVMGetExecutionEngineTargetMachine(ee.c)}
}

func LLVMAddGlobalMapping(ee LLVMExecutionEngineRef, global LLVMValueRef, addr unsafe.Pointer) {
	C.LLVMAddGlobalMapping(ee.c, global.c, addr)
}

func LLVMGetPointerToGlobal(ee LLVMExecutionEngineRef, global LLVMValueRef) unsafe.Pointer {
	return C.LLVMGetPointerToGlobal(ee.c, global.c)
}

func LLVMGetGlobalValueAddress(ee LLVMExecutionEngineRef, name string) uintptr {
	return string2CString(name, func(name *C.char) uintptr {
		return uintptr(C.LLVMGetGlobalValueAddress(ee.c, name))
	})
}

func LLVMGetFunctionAddress(ee LLVMExecutionEngineRef, name string) uintptr {
	return string2CString(name, func(name *C.char) uintptr {
		return uintptr(C.LLVMGetFunctionAddress(ee.c, name))
	})
}

func LLVMExecutionEngineGetErrMsg(ee LLVMExecutionEngineRef) error {
	return llvmError2Error(func(outError **C.char) C.LLVMBool {
		return C.LLVMExecutionEngineGetErrMsg(ee.c, outError)
	})
}
