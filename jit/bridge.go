package jit

import (
	"math"
	"reflect"
	"sync"
	"unsafe"

	"github.com/kkkunny/go-llvm/internal/binding"
)

// bridgeRegistry 注册的宿主 Go 函数（native→Go 方向）；下标即通道 idx
var (
	bridgeMu       sync.RWMutex
	bridgeRegistry []reflect.Value
)

func init() {
	binding.SetBridgeDispatch(dispatchToGo)
}

// registerGoFunc 注册宿主 Go 函数并返回通道 idx
func registerGoFunc(f reflect.Value) int64 {
	bridgeMu.Lock()
	defer bridgeMu.Unlock()
	bridgeRegistry = append(bridgeRegistry, f)
	return int64(len(bridgeRegistry) - 1)
}

// dispatchToGo 固定签名通道的 Go 侧终点：解箱 slots → reflect 调用 → 返回装箱值
func dispatchToGo(idx int64, slots []uint64) uint64 {
	bridgeMu.RLock()
	var f reflect.Value
	if idx >= 0 && int(idx) < len(bridgeRegistry) {
		f = bridgeRegistry[idx]
	}
	bridgeMu.RUnlock()
	if !f.IsValid() {
		return 0
	}
	ft := f.Type()
	args := make([]reflect.Value, ft.NumIn())
	for i := range args {
		args[i] = unboxValue(ft.In(i), slots[i])
	}
	out := f.Call(args)
	if len(out) == 0 {
		return 0
	}
	return boxValue(out[0])
}

// boxValue Go 值 → 槽位
func boxValue(v reflect.Value) uint64 {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return 1
		}
		return 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint()
	case reflect.Float32:
		return uint64(math.Float32bits(float32(v.Float())))
	case reflect.Float64:
		return math.Float64bits(v.Float())
	case reflect.Pointer, reflect.UnsafePointer:
		if v.IsNil() {
			return 0
		}
		return uint64(v.Pointer())
	}
	return 0
}

// unboxValue 槽位 → Go 值
func unboxValue(t reflect.Type, slot uint64) reflect.Value {
	out := reflect.New(t).Elem()
	switch t.Kind() {
	case reflect.Bool:
		out.SetBool(slot != 0)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out.SetInt(int64(slot))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		out.SetUint(slot)
	case reflect.Float32:
		out.SetFloat(float64(math.Float32frombits(uint32(slot))))
	case reflect.Float64:
		out.SetFloat(math.Float64frombits(slot))
	case reflect.Pointer, reflect.UnsafePointer:
		out.SetPointer(slotToPtr(slot))
	}
	return out
}

// slotToPtr 槽位位模式 → 指针（避免 uintptr→unsafe.Pointer 转换触发 vet 告警）
func slotToPtr(slot uint64) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&slot))
}
