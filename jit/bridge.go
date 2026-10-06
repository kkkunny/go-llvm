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

// registerGoFunc 注册宿主 Go 函数并返回通道 idx（复用已注销的空槽）
func registerGoFunc(f reflect.Value) int64 {
	bridgeMu.Lock()
	defer bridgeMu.Unlock()
	for i, v := range bridgeRegistry {
		if !v.IsValid() {
			bridgeRegistry[i] = f
			return int64(i)
		}
	}
	bridgeRegistry = append(bridgeRegistry, f)
	return int64(len(bridgeRegistry) - 1)
}

// unregisterGoFunc 注销宿主 Go 函数（编译包装体失败时回收；槽位可被后续注册复用）
func unregisterGoFunc(idx int64) {
	bridgeMu.Lock()
	defer bridgeMu.Unlock()
	if idx >= 0 && int(idx) < len(bridgeRegistry) {
		bridgeRegistry[idx] = reflect.Value{}
	}
}

// bridgeArgsPool 复用回调参数切片，避免每次 native→Go 调用分配
var bridgeArgsPool = sync.Pool{New: func() any {
	s := make([]reflect.Value, 0, 8)
	return &s
}}

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
	bufp := bridgeArgsPool.Get().(*[]reflect.Value)
	args := (*bufp)[:0]
	for i := 0; i < ft.NumIn(); i++ {
		args = append(args, unboxValue(ft.In(i), slots[i]))
	}
	out := f.Call(args)
	clear(args)
	*bufp = args[:0]
	bridgeArgsPool.Put(bufp)
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

// 常见标量类型句柄：unboxValue 用它们判断能否走无分配快路径
var (
	typeBool          = reflect.TypeOf(false)
	typeInt           = reflect.TypeOf(int(0))
	typeInt32         = reflect.TypeOf(int32(0))
	typeInt64         = reflect.TypeOf(int64(0))
	typeUint          = reflect.TypeOf(uint(0))
	typeUint32        = reflect.TypeOf(uint32(0))
	typeUint64        = reflect.TypeOf(uint64(0))
	typeUintptr       = reflect.TypeOf(uintptr(0))
	typeFloat32       = reflect.TypeOf(float32(0))
	typeFloat64       = reflect.TypeOf(float64(0))
	typeUnsafePointer = reflect.TypeOf(unsafe.Pointer(nil))
)

// unboxValue 槽位 → Go 值。常见标量/指针走 reflect.ValueOf/NewAt 快路径（无 reflect.New
// 分配）；具名类型与其余种类回退到 reflect.New(t).Elem()。
func unboxValue(t reflect.Type, slot uint64) reflect.Value {
	switch t {
	case typeBool:
		return reflect.ValueOf(slot != 0)
	case typeInt:
		return reflect.ValueOf(int(slot))
	case typeInt32:
		return reflect.ValueOf(int32(slot))
	case typeInt64:
		return reflect.ValueOf(int64(slot))
	case typeUint:
		return reflect.ValueOf(uint(slot))
	case typeUint32:
		return reflect.ValueOf(uint32(slot))
	case typeUint64:
		return reflect.ValueOf(slot)
	case typeUintptr:
		return reflect.ValueOf(uintptr(slot))
	case typeFloat32:
		return reflect.ValueOf(math.Float32frombits(uint32(slot)))
	case typeFloat64:
		return reflect.ValueOf(math.Float64frombits(slot))
	case typeUnsafePointer:
		return reflect.ValueOf(slotToPtr(slot))
	}
	if t.Kind() == reflect.Pointer && t.Name() == "" {
		if slot == 0 {
			return reflect.Zero(t)
		}
		return reflect.NewAt(t.Elem(), slotToPtr(slot))
	}
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
