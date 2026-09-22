package jit

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"unsafe"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/ir"
)

func errBridge(msg string) error {
	return &llvm.Error{Reason: llvm.ErrUnsupported, Op: "jit.bridge", Msg: msg}
}

// hashType 函数签名的稳定短哈希（用于适配器/包装体符号名）
func hashType(ft reflect.Type) string {
	sum := sha256.Sum256([]byte(ft.String()))
	return hex.EncodeToString(sum[:8])
}

// adapterName 适配器符号名（按签名区分，包内唯一）
func adapterName(ft reflect.Type) string {
	return "__go_llvm_bridge_adapter_" + hashType(ft)
}

// bridgeSlotKind 槽位装箱/拆箱类别
type bridgeSlotKind uint8

const (
	slotUnsupported bridgeSlotKind = iota
	slotI1
	slotI32
	slotI64
	slotF32
	slotF64
	slotPtr
)

// bridgeSlots 桥首版支持的 Go 类型（其余返回 ErrUnsupported）
var bridgeSlots = map[reflect.Kind]bridgeSlotKind{
	reflect.Bool:          slotI1,
	reflect.Int32:         slotI32,
	reflect.Uint32:        slotI32,
	reflect.Int64:         slotI64,
	reflect.Uint64:        slotI64,
	reflect.Int:           slotI64,
	reflect.Uint:          slotI64,
	reflect.Float32:       slotF32,
	reflect.Float64:       slotF64,
	reflect.Pointer:       slotPtr,
	reflect.UnsafePointer: slotPtr,
}

// slotKindOf 判定 Go 类型可否过桥
func slotKindOf(t reflect.Type) bridgeSlotKind {
	k, ok := bridgeSlots[t.Kind()]
	if !ok {
		return slotUnsupported
	}
	if t.Kind() == reflect.Pointer && t.Elem() != nil {
		// 允许 *T；不支持指向不支持类型的指针无影响（不透明指针）
	}
	return k
}

// checkBridgeFunc 校验函数签名可过桥；返回参数/返回值槽类别
func checkBridgeFunc(ft reflect.Type) ([]bridgeSlotKind, bridgeSlotKind, error) {
	if ft.Kind() != reflect.Func {
		return nil, slotUnsupported, errBridge("not a function type: " + ft.String())
	}
	if ft.IsVariadic() {
		return nil, slotUnsupported, errBridge("variadic functions are not supported")
	}
	if ft.NumIn() > binding.BridgeMaxSlots-1 {
		return nil, slotUnsupported, errBridge("too many parameters (max 15)")
	}
	if ft.NumOut() > 1 {
		return nil, slotUnsupported, errBridge("multiple results are not supported")
	}
	params := make([]bridgeSlotKind, ft.NumIn())
	for i := range params {
		k := slotKindOf(ft.In(i))
		if k == slotUnsupported {
			return nil, slotUnsupported, errBridge("unsupported parameter type: " + ft.In(i).String())
		}
		params[i] = k
	}
	var ret = slotUnsupported
	if ft.NumOut() == 1 {
		ret = slotKindOf(ft.Out(0))
		if ret == slotUnsupported {
			return nil, slotUnsupported, errBridge("unsupported result type: " + ft.Out(0).String())
		}
	}
	return params, ret, nil
}

// adapterCache 按 Go 函数签名缓存的适配器/包装体
type adapterEntry struct {
	// Go→native 适配器地址；参数/返回值槽类别
	adapter unsafe.Pointer
	params  []bridgeSlotKind
	ret     bridgeSlotKind
}

func (j *LLJIT) adapterFor(ft reflect.Type) (*adapterEntry, error) {
	if e, ok := j.adapters[ft]; ok {
		return e, nil
	}
	params, ret, err := checkBridgeFunc(ft)
	if err != nil {
		return nil, err
	}
	e := &adapterEntry{params: params, ret: ret}
	if err := j.compileAdapter(ft, e); err != nil {
		return nil, err
	}
	if j.adapters == nil {
		j.adapters = make(map[reflect.Type]*adapterEntry)
	}
	j.adapters[ft] = e
	return e, nil
}

// compileAdapter 生成 Go→native 适配器：
//
//	uint64_t adapter(void *fn, uint64_t *slots)
//
// 解箱 slots[0..n-1] 为参数，间接调用 fn，返回值装箱为 uint64 返回
func (j *LLJIT) compileAdapter(ft reflect.Type, e *adapterEntry) error {
	const op = "jit.compileAdapter"
	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "bridge_adapter")
	sig, _, err := llvm.FnSignatureOfGo(ctx, ft)
	if err != nil {
		ctx.Close()
		return llvm.WrapError(llvm.ErrUnsupported, op, err)
	}

	i64 := ctx.Int(64)
	ptr := ctx.Ptr(0)
	bridgeSig := ctx.Fn(i64, []llvm.AnyType{ptr, ptr}, false)
	fn := m.NewFunction(adapterName(ft), bridgeSig)
	fn.SetLinkage(llvm.LinkageExternal)

	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	fnPtr := fn.ParamAs[llvm.PtrT](0)
	slots := fn.ParamAs[llvm.PtrT](1)

	args := make([]llvm.AnyValue, len(e.params))
	for i, k := range e.params {
		slot := b.GEP(i64, slots, []llvm.ValueRef[llvm.IntT]{ctx.ConstInt(i64, uint64(i), false)}, "")
		raw := b.Load(slot, i64, "").Value
		args[i] = unpackSlot(b, ctx, raw, sig.Params()[i], k).Dyn()
	}

	call := b.CallIndirect[llvm.DynT](fnPtr, sig, args, "")
	var ret llvm.AnyValue = call.Value
	if e.ret != slotUnsupported {
		ret = packSlot(b, ctx, call.Value, i64, e.ret)
	}
	b.Ret(ret)
	if err := b.Close(); err != nil {
		ctx.Close()
		return llvm.WrapError(llvm.ErrInternal, op, err)
	}
	if err := m.Verify(); err != nil {
		ctx.Close()
		return llvm.WrapError(llvm.ErrInternal, op, err)
	}
	if err := j.AddIRModule(m); err != nil {
		return err
	}
	addr, err := j.Lookup(adapterName(ft))
	if err != nil {
		return err
	}
	e.adapter = addr
	return nil
}

// unpackSlot 槽位（i64）→ 目标类型值
func unpackSlot(b *ir.Builder, ctx *llvm.Context, raw llvm.Value[llvm.IntT], ty llvm.AnyType, k bridgeSlotKind) llvm.AnyValue {
	switch k {
	case slotI1:
		return b.Trunc(raw, llvm.AsIntType(ty), "").Dyn()
	case slotI32:
		return b.Trunc(raw, llvm.AsIntType(ty), "").Dyn()
	case slotI64:
		return raw.Dyn()
	case slotF32:
		i32 := b.Trunc(raw, ctx.Int(32), "")
		return b.BitCast(i32, llvm.AsFloatType(ty), "").Dyn()
	case slotF64:
		return b.BitCast(raw, llvm.AsFloatType(ty), "").Dyn()
	case slotPtr:
		return b.IntToPtr(raw, llvm.AsPtrType(ty), "").Dyn()
	}
	return raw.Dyn()
}

// packSlot 目标类型值 → 槽位（i64）
func packSlot(b *ir.Builder, ctx *llvm.Context, v llvm.Value[llvm.DynT], i64 llvm.IntType, k bridgeSlotKind) llvm.AnyValue {
	switch k {
	case slotI1:
		return b.ZExt(v.MustAs[llvm.IntT](), i64, "").Dyn()
	case slotI32:
		return b.ZExt(v.MustAs[llvm.IntT](), i64, "").Dyn()
	case slotI64:
		return v.Dyn()
	case slotF32:
		bits := b.BitCast(v.MustAs[llvm.FloatT](), ctx.Int(32), "")
		return b.ZExt(bits, i64, "").Dyn()
	case slotF64:
		return b.BitCast(v.MustAs[llvm.FloatT](), i64, "").Dyn()
	case slotPtr:
		return b.PtrToInt(v.MustAs[llvm.PtrT](), i64, "").Dyn()
	}
	return v.Dyn()
}
