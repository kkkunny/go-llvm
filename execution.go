package llvm

import "C"
import (
	"os"
	"unsafe"

	"github.com/samber/lo"

	"github.com/kkkunny/go-llvm/internal/binding"
)

type ExecutionValue binding.LLVMGenericValueRef

func NewIntExecutionValue(t IntegerType, v uint64, isSigned bool) ExecutionValue {
	return ExecutionValue(binding.LLVMCreateGenericValueOfInt(t.binding(), v, isSigned))
}

func NewFloatExecutionValue(t FloatType, v float64) ExecutionValue {
	return ExecutionValue(binding.LLVMCreateGenericValueOfFloat(t.binding(), v))
}

func (v ExecutionValue) binding() binding.LLVMGenericValueRef {
	return binding.LLVMGenericValueRef(v)
}

func (v ExecutionValue) Free() {
	binding.LLVMDisposeGenericValue(v.binding())
}

func (v ExecutionValue) Integer(isSigned bool) uint64 {
	return binding.LLVMGenericValueToInt(v.binding(), isSigned)
}

func (v ExecutionValue) Float(t FloatType) float64 {
	return binding.LLVMGenericValueToFloat(t.binding(), v.binding())
}

type ExecutionEngine struct {
	bind   binding.LLVMExecutionEngineRef
	module Module

	funcMaps map[uint64]any
}

func newExecutionEngine(m Module, b binding.LLVMExecutionEngineRef) *ExecutionEngine {
	return &ExecutionEngine{
		bind:     b,
		module:   m,
		funcMaps: make(map[uint64]any),
	}
}

func NewExecutionEngine(m Module) (*ExecutionEngine, error) {
	b, err := binding.LLVMCreateExecutionEngineForModule(m.binding())
	return newExecutionEngine(m, b), err
}

func NewInterpreter(m Module) (*ExecutionEngine, error) {
	b, err := binding.LLVMCreateInterpreterForModule(m.binding())
	return newExecutionEngine(m, b), err
}

func NewJITCompiler(m Module, opt CodeOptLevel) (*ExecutionEngine, error) {
	b, err := binding.LLVMCreateJITCompilerForModule(m.binding(), uint32(opt))
	return newExecutionEngine(m, b), err
}

func DefaultMCJITCompiler(m Module) (*ExecutionEngine, error) {
	var option binding.LLVMMCJITCompilerOptions
	binding.LLVMInitializeMCJITCompilerOptions(&option)
	b, err := binding.LLVMCreateMCJITCompilerForModule(m.binding(), option)
	return newExecutionEngine(m, b), err
}

func (engine ExecutionEngine) binding() binding.LLVMExecutionEngineRef {
	return engine.bind
}

func (engine ExecutionEngine) Free() {
	binding.LLVMDisposeExecutionEngine(engine.binding())
}

func (engine ExecutionEngine) RunMainFunction(f Function, argv, envp []string) uint8 {
	return uint8(binding.LLVMRunFunctionAsMain(engine.binding(), f.binding(), argv, envp))
}

func (engine ExecutionEngine) RunMainFunctionWithParentEnv(f Function) uint8 {
	return engine.RunMainFunction(f, os.Args, os.Environ())
}

func (engine ExecutionEngine) RunFunction(f Function, args ...ExecutionValue) ExecutionValue {
	as := lo.Map(args, func(item ExecutionValue, index int) binding.LLVMGenericValueRef {
		return item.binding()
	})
	return ExecutionValue(binding.LLVMRunFunction(engine.binding(), f.binding(), as))
}

func (engine ExecutionEngine) GetFunction(name string) (Function, bool) {
	v, fail := binding.LLVMFindFunction(engine.binding(), name)
	if !fail {
		return Function(v), !v.IsNil()
	}
	return engine.module.GetFunction(name)
}

func (engine ExecutionEngine) GetFunctionRuntimePointer(f Function) (unsafe.Pointer, bool) {
	res := binding.LLVMGetPointerToGlobal(engine.bind, f.binding())
	return res, res != nil
}

func (engine ExecutionEngine) GetVariable(name string) (GlobalValue, bool) {
	return engine.module.GetGlobal(name)
}

func (engine ExecutionEngine) GetVariableRuntimePointer(v GlobalValue) (unsafe.Pointer, bool) {
	res := binding.LLVMGetPointerToGlobal(engine.bind, v.binding())
	return res, res != nil
}

// MapGlobalToC 映射全局值到c语言值
func (engine ExecutionEngine) MapGlobalToC(g Global, to unsafe.Pointer) {
	binding.LLVMAddGlobalMapping(engine.binding(), g.binding(), to)
}
