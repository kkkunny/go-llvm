package llvm

import "C"
import (
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

type ExecutionEngine binding.LLVMExecutionEngineRef

func NewExecutionEngine(m Module) (ExecutionEngine, error) {
	b, err := binding.LLVMCreateExecutionEngineForModule(m.binding())
	return ExecutionEngine(b), err
}

func NewInterpreter(m Module) (ExecutionEngine, error) {
	b, err := binding.LLVMCreateInterpreterForModule(m.binding())
	return ExecutionEngine(b), err
}

func NewJITCompiler(m Module, opt CodeOptLevel) (ExecutionEngine, error) {
	b, err := binding.LLVMCreateJITCompilerForModule(m.binding(), uint32(opt))
	return ExecutionEngine(b), err
}

func DefaultMCJITCompiler(m Module) (ExecutionEngine, error) {
	var option binding.LLVMMCJITCompilerOptions
	binding.LLVMInitializeMCJITCompilerOptions(&option)
	b, err := binding.LLVMCreateMCJITCompilerForModule(m.binding(), option)
	return ExecutionEngine(b), err
}

func (engine ExecutionEngine) binding() binding.LLVMExecutionEngineRef {
	return binding.LLVMExecutionEngineRef(engine)
}

func (engine ExecutionEngine) Free() {
	binding.LLVMDisposeExecutionEngine(engine.binding())
}

func (engine ExecutionEngine) RunFunction(f Function, args ...ExecutionValue) ExecutionValue {
	as := lo.Map(args, func(item ExecutionValue, index int) binding.LLVMGenericValueRef {
		return item.binding()
	})
	return ExecutionValue(binding.LLVMRunFunction(engine.binding(), f.binding(), as))
}

func (engine ExecutionEngine) GetFunction(name string) *Function {
	v, fail := binding.LLVMFindFunction(engine.binding(), name)
	if fail {
		return nil
	}
	f := Function(v)
	return &f
}

func (engine ExecutionEngine) MapGlobal(g Global, to unsafe.Pointer) {
	binding.LLVMAddGlobalMapping(engine.binding(), g.binding(), to)
}
