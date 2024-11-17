package llvm

/*
void goFuncCallChannel(int funcIdx, void **ret_ptr_ptr, void **param_ptr_array_ptr);
*/
import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"github.com/kkkunny/go-llvm/internal/binding"
)

var funcMapLock sync.RWMutex
var funcMap = make(map[int32]func(retPtrPtr *unsafe.Pointer, paramPtrArrayPtr *unsafe.Pointer))

//export goFuncCallChannel
func goFuncCallChannel(funcIdx C.int, retPtrPtr *unsafe.Pointer, paramPtrArrayPtr *unsafe.Pointer) {
	funcMapLock.RLock()
	defer funcMapLock.RUnlock()
	funcMap[int32(funcIdx)](retPtrPtr, paramPtrArrayPtr)
}

// MapFunctionToGo 映射函数到go函数
func (engine ExecutionEngine) MapFunctionToGo(name string, to any) error {
	// from check
	if name == "" {
		return errors.New("expect a name which is not empty")
	}
	f, ok := engine.module.GetFunction(name)
	if !ok {
		return fmt.Errorf("unknown function which named `%s`", name)
	}

	// to check
	ctx := engine.module.Context()
	ft := f.FunctionType()
	toVal := reflect.ValueOf(to)
	toFt := toVal.Type()
	if toFt.Kind() != reflect.Func {
		return errors.New("expect a function")
	} else if toFt.NumIn() != int(ft.CountParams()) {
		return fmt.Errorf("type mismatch, expected type: %s", ft.String())
	}
	var ftRetNum int
	if ft.ReturnType().String() != ctx.VoidType().String() {
		ftRetNum = 1
	}
	if toFt.NumOut() != ftRetNum {
		return fmt.Errorf("type mismatch, expected type: %s", ft.String())
	}

	// covert to
	toFn := func(retPtrPtr *unsafe.Pointer, paramPtrArrayPtr *unsafe.Pointer) {
		params := make([]reflect.Value, ft.CountParams())
		for i, _ := range ft.Params() {
			paramPtrPtr := unsafe.Pointer(uintptr(unsafe.Pointer(paramPtrArrayPtr)) + PointerSize*uintptr(i))
			paramPtr := *(*unsafe.Pointer)(paramPtrPtr)
			params[i] = llvmPtr2GoReflectValue(toFt.In(i), paramPtr)
		}
		res := toVal.Call(params)
		if len(res) != 0 {
			*retPtrPtr = goReflectValue2llvmPtr(res[0])
		}
	}

	// lock
	funcMapLock.Lock()
	defer funcMapLock.Unlock()

	// idx
	idx := int32(len(funcMap))
	if existIdx, ok := engine.funcMap[name]; ok {
		funcMap[existIdx] = toFn
		return nil
	} else {
		funcMap[idx] = toFn
	}

	// ctx
	builder := ctx.NewBuilder()

	// channel
	cf, ok := engine.GetFunction("goFuncCallChannel")
	if !ok {
		cf = engine.module.NewFunction("goFuncCallChannel", ctx.FunctionType(false, ctx.VoidType(), ctx.IntegerType(32), ctx.PointerType(ctx.PointerType(ctx.VoidType())), ctx.PointerType(ctx.PointerType(ctx.VoidType()))))
		binding.LLVMAddGlobalMapping(engine.binding(), cf.binding(), C.goFuncCallChannel)
	}

	builder.MoveToAfter(f.NewBlock("entry"))

	// ret
	var retPtrPtr Value
	var retPtrPtrParam Value
	if ft.ReturnType().String() != ctx.VoidType().String() {
		retPtrPtr = builder.CreateAlloca("", ctx.PointerType(ft.ReturnType()))
		retPtrPtrParam = builder.CreateBitCast("", retPtrPtr, ctx.PointerType(ctx.PointerType(ctx.VoidType())))
	} else {
		retPtrPtrParam = ctx.ConstNull(ctx.PointerType(ctx.PointerType(ctx.VoidType())))
	}

	// params
	var paramPtrArrayPtrParam Value
	if ft.CountParams() > 0 {
		at := ctx.ArrayType(ctx.PointerType(ctx.VoidType()), ft.CountParams())
		paramPtrArrayPtrParam = builder.CreateAlloca("", at)
		for i, pt := range ft.Params() {
			paramPtrPtr := builder.CreateInBoundsGEP("", at, paramPtrArrayPtrParam, ctx.ConstInteger(ctx.IntegerType(64), 0), ctx.ConstInteger(ctx.IntegerType(64), int64(i)))
			paramPtr := builder.CreateAlloca("", pt)
			builder.CreateStore(f.GetParam(uint(i)), paramPtr)
			paramPtr = Alloca(builder.CreateBitCast("", paramPtr, ctx.PointerType(ctx.VoidType())))
			builder.CreateStore(paramPtr, paramPtrPtr)
		}
	} else {
		paramPtrArrayPtrParam = ctx.ConstNull(ctx.PointerType(ctx.PointerType(ctx.VoidType())))
	}

	// call
	builder.CreateCall("", cf.FunctionType(), cf, ctx.ConstInteger(ctx.IntegerType(32), int64(idx)), retPtrPtrParam, paramPtrArrayPtrParam)

	// ret
	if ftRetNum != 0 {
		retPtr := builder.CreateLoad("", ctx.PointerType(ft.ReturnType()), retPtrPtr)
		var ret Value = builder.CreateLoad("", ft.ReturnType(), retPtr)
		builder.CreateRet(&ret)
	} else {
		builder.CreateRet(nil)
	}
	return nil
}

func llvmPtr2GoReflectValue(t reflect.Type, ptr unsafe.Pointer) reflect.Value {
	valPtr := reflect.NewAt(t, ptr)
	return valPtr.Elem()
}

func goReflectValue2llvmPtr(val reflect.Value) unsafe.Pointer {
	ptr := unsafeGetPointer(val.Interface())
	return ptr
}
