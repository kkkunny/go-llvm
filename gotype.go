package llvm

import (
	"errors"
	"reflect"
	"strconv"
)

// TypeOf 将 Go 类型映射为 LLVM 类型。
//
// 映射规则：bool→i1；intN/uintN→iN；int/uint→按平台位宽；float32→float；float64→double；
// *T/unsafe.Pointer→ptr；[N]T→数组；全字段可映射的 struct→字面量结构体。
// 不支持的类型（string/slice/map/chan/func/interface/complex/含不可映射字段的结构体）返回 ErrUnsupported。
func TypeOf[T any](ctx *Context) (Type[DynT], error) {
	return TypeOfGo(ctx, reflect.TypeOf((*T)(nil)).Elem())
}

// TypeOfGo TypeOf 的 reflect.Type 版本
func TypeOfGo(ctx *Context, t reflect.Type) (Type[DynT], error) {
	ty, err := typeOfGo(ctx, t)
	if err != nil {
		return Type[DynT]{}, err
	}
	return ty.DynType(), nil
}

func typeOfGo(ctx *Context, t reflect.Type) (AnyType, error) {
	switch t.Kind() {
	case reflect.Bool:
		return ctx.Bool(), nil
	case reflect.Int, reflect.Uint:
		return ctx.Int(uint32(strconv.IntSize)), nil
	case reflect.Int8:
		return ctx.Int(8), nil
	case reflect.Int16:
		return ctx.Int(16), nil
	case reflect.Int32:
		return ctx.Int(32), nil
	case reflect.Int64:
		return ctx.Int(64), nil
	case reflect.Uint8:
		return ctx.Int(8), nil
	case reflect.Uint16:
		return ctx.Int(16), nil
	case reflect.Uint32:
		return ctx.Int(32), nil
	case reflect.Uint64:
		return ctx.Int(64), nil
	case reflect.Uintptr:
		return ctx.Int(uint32(strconv.IntSize)), nil
	case reflect.Float32:
		return ctx.Float(FloatSingle), nil
	case reflect.Float64:
		return ctx.Float(FloatDouble), nil
	case reflect.Pointer, reflect.UnsafePointer:
		return ctx.Ptr(0), nil
	case reflect.Array:
		elem, err := typeOfGo(ctx, t.Elem())
		if err != nil {
			return nil, err
		}
		return ctx.Array(elem, uint64(t.Len())), nil
	case reflect.Struct:
		fields := make([]AnyType, t.NumField())
		for i := range fields {
			ft, err := typeOfGo(ctx, t.Field(i).Type)
			if err != nil {
				var inner *Error
				if errors.As(err, &inner) {
					return nil, &Error{
						Reason: inner.Reason,
						Op:     "llvm.TypeOfGo",
						Msg:    "struct field " + t.Field(i).Name + ": " + inner.Msg,
						cause:  err,
					}
				}
				return nil, &Error{
					Reason: ErrUnsupported,
					Op:     "llvm.TypeOfGo",
					Msg:    "struct field " + t.Field(i).Name + ": " + err.Error(),
					cause:  err,
				}
			}
			fields[i] = ft
		}
		return ctx.Struct(fields, false), nil
	default:
		return nil, &Error{
			Reason: ErrUnsupported,
			Op:     "llvm.TypeOfGo",
			Msg:    "unsupported Go type: " + t.String(),
		}
	}
}

// ConstOf 将 Go 值映射为 LLVM 常量。
//
// 标量映射为整数/浮点常量；nil 指针映射为 null；非 nil 指针映射为 inttoptr(uintptr)；
// struct/array 递归构造。不支持的类型返回 ErrUnsupported。
func ConstOf[T any](ctx *Context, v T) (Value[DynT], error) {
	return ConstOfGo(ctx, reflect.ValueOf(v))
}

// ConstOfGo ConstOf 的 reflect.Value 版本
func ConstOfGo(ctx *Context, v reflect.Value) (Value[DynT], error) {
	c, err := constOfGo(ctx, v)
	if err != nil {
		return Value[DynT]{}, err
	}
	return c.Dyn(), nil
}

func constOfGo(ctx *Context, v reflect.Value) (AnyValue, error) {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return ctx.ConstBool(true), nil
		}
		return ctx.ConstBool(false), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		t, err := typeOfGo(ctx, v.Type())
		if err != nil {
			return nil, err
		}
		return ctx.ConstInt(AsIntType(t), uint64(v.Int()), true).Value, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		t, err := typeOfGo(ctx, v.Type())
		if err != nil {
			return nil, err
		}
		return ctx.ConstInt(AsIntType(t), v.Uint(), false).Value, nil
	case reflect.Float32, reflect.Float64:
		t, err := typeOfGo(ctx, v.Type())
		if err != nil {
			return nil, err
		}
		return ctx.ConstFloat(AsFloatType(t), v.Float()).Value, nil
	case reflect.Pointer, reflect.UnsafePointer:
		ptrTy := ctx.Ptr(0)
		if v.IsNil() {
			return ctx.ConstNull(ptrTy), nil
		}
		addr := ctx.ConstInt(ctx.Int(uint32(strconv.IntSize)), uint64(v.Pointer()), false)
		return ctx.ConstIntToPtr(addr, ptrTy), nil
	case reflect.Array:
		elems := make([]AnyValue, v.Len())
		for i := range elems {
			e, err := constOfGo(ctx, v.Index(i))
			if err != nil {
				return nil, err
			}
			elems[i] = e
		}
		elemTy, err := typeOfGo(ctx, v.Type().Elem())
		if err != nil {
			return nil, err
		}
		return ctx.ConstArray(elemTy, elems...), nil
	case reflect.Struct:
		elems := make([]AnyValue, v.NumField())
		for i := range elems {
			e, err := constOfGo(ctx, v.Field(i))
			if err != nil {
				return nil, err
			}
			elems[i] = e
		}
		return ctx.ConstStruct(false, elems...), nil
	default:
		return nil, &Error{
			Reason: ErrUnsupported,
			Op:     "llvm.ConstOfGo",
			Msg:    "unsupported Go value: " + v.Type().String(),
		}
	}
}

// FnSignatureOf 将 Go 函数类型映射为 LLVM 函数类型。
// 返回签名、Go 函数类型（供 JIT 桥使用）或 ErrUnsupported。
func FnSignatureOf[F any](ctx *Context) (FnType, reflect.Type, error) {
	return FnSignatureOfGo(ctx, reflect.TypeOf((*F)(nil)).Elem())
}

// FnSignatureOfGo FnSignatureOf 的 reflect.Type 版本
func FnSignatureOfGo(ctx *Context, ft reflect.Type) (FnType, reflect.Type, error) {
	if ft.Kind() != reflect.Func {
		return FnType{}, nil, &Error{Reason: ErrUnsupported, Op: "llvm.FnSignatureOf", Msg: "not a function type: " + ft.String()}
	}
	if ft.IsVariadic() {
		return FnType{}, nil, &Error{Reason: ErrUnsupported, Op: "llvm.FnSignatureOf", Msg: "variadic Go functions are not supported"}
	}
	if ft.NumOut() > 1 {
		return FnType{}, nil, &Error{Reason: ErrUnsupported, Op: "llvm.FnSignatureOf", Msg: "functions with multiple results are not supported"}
	}

	var ret AnyType = ctx.Void()
	if ft.NumOut() == 1 {
		r, err := typeOfGo(ctx, ft.Out(0))
		if err != nil {
			return FnType{}, nil, err
		}
		ret = r
	}
	params := make([]AnyType, ft.NumIn())
	for i := range params {
		p, err := typeOfGo(ctx, ft.In(i))
		if err != nil {
			return FnType{}, nil, err
		}
		params[i] = p
	}
	return ctx.Fn(ret, params, false), ft, nil
}
