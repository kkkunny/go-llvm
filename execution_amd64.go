package llvm

/*
void go_func_call_channel(int funcIdx, void **ret_ptr_ptr, void **param_ptr_array_ptr);
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

//export go_func_call_channel
func go_func_call_channel(funcIdx C.int, retPtrPtr *unsafe.Pointer, paramPtrArrayPtr *unsafe.Pointer){
	funcMapLock.RLock()
	defer funcMapLock.RUnlock()
	funcMap[int32(funcIdx)](retPtrPtr, paramPtrArrayPtr)
}

// MapGlobalToC 映射全局值到c语言值
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
	} else if toFt.NumIn() != int(ft.CountParams()){
		return fmt.Errorf("type mismatch, expected type: %s", ft.String())
	}
	var ftRetNum int
	if ft.ReturnType().String() != ctx.VoidType().String(){
		ftRetNum = 1
	}
	if toFt.NumOut() != ftRetNum{
		return fmt.Errorf("type mismatch, expected type: %s", ft.String())
	}

	// covert to
	toFn := func(retPtrPtr *unsafe.Pointer, paramPtrArrayPtr *unsafe.Pointer){
		params := make([]reflect.Value, ft.CountParams())
		for i, pt := range ft.Params(){
			paramPtrPtr := unsafe.Pointer(uintptr(unsafe.Pointer(paramPtrArrayPtr))+8*uintptr(i))
			paramPtr := *(*unsafe.Pointer)(paramPtrPtr)
			params[i] = llvmPtr2GoReflectValue(pt, toFt.In(i), paramPtr)
		}
		res := toVal.Call(params)
		if len(res) != 0{
			*retPtrPtr = goReflectValue2llvmPtr(toFt.Out(0), ft.ReturnType(), res[0])
		}
	}

	// lock
	funcMapLock.Lock()
	defer funcMapLock.Unlock()

	// idx
	idx := int32(len(funcMap))
	if existIdx, ok := engine.funcMap[name]; ok{
		funcMap[existIdx] = toFn
		return nil
	}else{
		funcMap[idx] = toFn
	}

	// ctx
	builder := ctx.NewBuilder()

	// channel
	cf, ok := engine.GetFunction("go_func_call_channel")
	if !ok{
		cf = engine.module.NewFunction("go_func_call_channel", ctx.FunctionType(false, ctx.VoidType(), ctx.IntegerType(32), ctx.PointerType(ctx.PointerType(ctx.VoidType())), ctx.PointerType(ctx.PointerType(ctx.VoidType()))))
		binding.LLVMAddGlobalMapping(engine.binding(), cf.binding(), C.go_func_call_channel)
	}

	builder.MoveToAfter(f.NewBlock("entry"))

	// ret
	var retPtrPtr Value
	var retPtrPtrParam Value
	if ft.ReturnType().String() != ctx.VoidType().String(){
		retPtrPtr = builder.CreateAlloca("", ctx.PointerType(ft.ReturnType()))
		retPtrPtrParam = builder.CreateBitCast("", retPtrPtr, ctx.PointerType(ctx.PointerType(ctx.VoidType())))
	}else{
		retPtrPtrParam = ctx.ConstNull(ctx.PointerType(ctx.PointerType(ctx.VoidType())))
	}

	// params
	var paramPtrArrayPtrParam Value
	if ft.CountParams() > 0{
		at := ctx.ArrayType(ctx.PointerType(ctx.VoidType()), ft.CountParams())
		paramPtrArrayPtrParam = builder.CreateAlloca("", at)
		for i, pt := range ft.Params(){
			paramPtrPtr := builder.CreateInBoundsGEP("", at, paramPtrArrayPtrParam, ctx.ConstInteger(ctx.IntegerType(64), 0), ctx.ConstInteger(ctx.IntegerType(64), int64(i)))
			paramPtr := builder.CreateAlloca("", pt)
			builder.CreateStore(f.GetParam(uint(i)), paramPtr)
			paramPtr = Alloca(builder.CreateBitCast("", paramPtr, ctx.PointerType(ctx.VoidType())))
			builder.CreateStore(paramPtr, paramPtrPtr)
		}
	}else{
		paramPtrArrayPtrParam = ctx.ConstNull(ctx.PointerType(ctx.PointerType(ctx.VoidType())))
	}

	// call
	builder.CreateCall("", cf.FunctionType(), cf, ctx.ConstInteger(ctx.IntegerType(32), int64(idx)), retPtrPtrParam, paramPtrArrayPtrParam)

	// ret
	if ftRetNum != 0{
		retPtr := builder.CreateLoad("", ctx.PointerType(ft.ReturnType()), retPtrPtr)
		var ret Value = builder.CreateLoad("", ft.ReturnType(), retPtr)
		builder.CreateRet(&ret)
	}else{
		builder.CreateRet(nil)
	}
	return nil
}

func llvmPtr2GoReflectValue(f Type, t reflect.Type, ptr unsafe.Pointer)reflect.Value{
	switch {
	case is[IntegerType](f) && t.Kind() == reflect.Int:
		if f.(IntegerType).Bits() == 64{
			v := *(*int)(ptr)
			return reflect.ValueOf(v)
		}else{
			panic("unreachable")
		}
	case is[IntegerType](f) && t.Kind() == reflect.Int8:
		if f.(IntegerType).Bits() == 8{
			v := *(*int8)(ptr)
			return reflect.ValueOf(v)
		}else{
			panic("unreachable")
		}
	case is[IntegerType](f) && t.Kind() == reflect.Int16:
		if f.(IntegerType).Bits() == 16{
			v := *(*int16)(ptr)
			return reflect.ValueOf(v)
		}else{
			panic("unreachable")
		}
	case is[IntegerType](f) && t.Kind() == reflect.Int32:
		if f.(IntegerType).Bits() == 32{
			v := *(*int32)(ptr)
			return reflect.ValueOf(v)
		}else{
			panic("unreachable")
		}
	case is[IntegerType](f) && t.Kind() == reflect.Int64:
		if f.(IntegerType).Bits() == 64{
			v := *(*int64)(ptr)
			return reflect.ValueOf(v)
		}else{
			panic("unreachable")
		}
	case is[IntegerType](f) && t.Kind() == reflect.Uint:
		if f.(IntegerType).Bits() == 64{
			v := *(*uint)(ptr)
			return reflect.ValueOf(v)
		}else{
			panic("unreachable")
		}
	case is[IntegerType](f) && t.Kind() == reflect.Uint8:
		if f.(IntegerType).Bits() == 8{
			v := *(*uint8)(ptr)
			return reflect.ValueOf(v)
		}else{
			panic("unreachable")
		}
	case is[IntegerType](f) && t.Kind() == reflect.Uint16:
		if f.(IntegerType).Bits() == 16{
			v := *(*uint16)(ptr)
			return reflect.ValueOf(v)
		}else{
			panic("unreachable")
		}
	case is[IntegerType](f) && t.Kind() == reflect.Uint32:
		if f.(IntegerType).Bits() == 32{
			v := *(*uint32)(ptr)
			return reflect.ValueOf(v)
		}else{
			panic("unreachable")
		}
	case is[IntegerType](f) && t.Kind() == reflect.Uint64:
		if f.(IntegerType).Bits() == 64{
			v := *(*uint64)(ptr)
			return reflect.ValueOf(v)
		}else{
			panic("unreachable")
		}
	case is[IntegerType](f) && t.Kind() == reflect.Uintptr:
		if f.(IntegerType).Bits() == 64{
			v := *(*uintptr)(ptr)
			return reflect.ValueOf(v)
		}else{
			panic("unreachable")
		}
	case is[IntegerType](f) && t.Kind() == reflect.Bool:
		switch f.(IntegerType).Bits(){
		case 8:
			v := *(*int8)(ptr) != 0
			return reflect.ValueOf(v)
		case 16:
			v := *(*int16)(ptr) != 0
			return reflect.ValueOf(v)
		case 32:
			v := *(*int32)(ptr) != 0
			return reflect.ValueOf(v)
		case 64:
			v := *(*int64)(ptr) != 0
			return reflect.ValueOf(v)
		default:
			panic("unreachable")
		}
	case is[FloatType](f) && t.Kind() == reflect.Float32:
		if f.(IntegerType).Bits() == 32{
			v := *(*float32)(ptr)
			return reflect.ValueOf(v)
		}else{
			panic("unreachable")
		}
	case is[FloatType](f) && t.Kind() == reflect.Float64:
		if f.(IntegerType).Bits() == 64{
			v := *(*float64)(ptr)
			return reflect.ValueOf(v)
		}else{
			panic("unreachable")
		}
	case is[PointerType](f) && t.Kind() == reflect.UnsafePointer:
		v := *(*unsafe.Pointer)(ptr)
		return reflect.ValueOf(v)
	default:
		panic("unreachable")
	}
}

func goReflectValue2llvmPtr(f reflect.Type, t Type, val reflect.Value)unsafe.Pointer{
	switch {
	case match(f.Kind(), reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64) && is[IntegerType](t):
		v := val.Int()
		switch t.(IntegerType).Bits(){
		case 8:
			vv := int8(v)
			return unsafe.Pointer(&vv)
		case 16:
			vv := int16(v)
			return unsafe.Pointer(&vv)
		case 32:
			vv := int32(v)
			return unsafe.Pointer(&vv)
		case 64:
			return unsafe.Pointer(&v)
		default:
			panic("unreachable")
		}
	case match(f.Kind(), reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr) && is[IntegerType](t):
		v := val.Uint()
		switch t.(IntegerType).Bits(){
		case 8:
			vv := uint8(v)
			return unsafe.Pointer(&vv)
		case 16:
			vv := uint16(v)
			return unsafe.Pointer(&vv)
		case 32:
			vv := uint32(v)
			return unsafe.Pointer(&vv)
		case 64:
			return unsafe.Pointer(&v)
		default:
			panic("unreachable")
		}
	case match(f.Kind(), reflect.Bool) && is[IntegerType](t):
		var v int64
		if val.Bool(){
			v = 1
		}
		switch t.(IntegerType).Bits(){
		case 8:
			vv := int8(v)
			return unsafe.Pointer(&vv)
		case 16:
			vv := int16(v)
			return unsafe.Pointer(&vv)
		case 32:
			vv := int32(v)
			return unsafe.Pointer(&vv)
		case 64:
			return unsafe.Pointer(&v)
		default:
			panic("unreachable")
		}
	case match(f.Kind(), reflect.Float32, reflect.Float64) && is[IntegerType](t):
		v := val.Float()
		switch t.(FloatType).Kind(){
		case FloatTypeKindFloat:
			vv := float32(v)
			return unsafe.Pointer(&vv)
		case FloatTypeKindDouble:
			return unsafe.Pointer(&v)
		default:
			panic("unreachable")
		}
	case match(f.Kind(), reflect.UnsafePointer, reflect.Pointer) && is[PointerType](t):
		v := val.UnsafePointer()
		return unsafe.Pointer(&v)
	default:
		panic("unreachable")
	}
}
