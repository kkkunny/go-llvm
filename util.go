package llvm

import (
	"reflect"
	"unsafe"
)

// 类型是否是
func is[T any](v any) bool {
	_, ok := v.(T)
	return ok
}

// 是否匹配
func match[T comparable](v T, to ...T) bool {
	for _, t := range to {
		if v == t {
			return true
		}
	}
	return false
}

func unsafeGetPointer(v any) unsafe.Pointer {
	vv := reflect.ValueOf(v)
	switch vt := vv.Type(); vt.Kind() {
	case reflect.Pointer:
		vp := reflect.New(vt)
		vp.Elem().Set(vv)
		return vp.UnsafePointer()
	default:
		return reflect.ValueOf(vv).FieldByName("ptr").UnsafePointer()
	}
}
