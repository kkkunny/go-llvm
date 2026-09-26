package jit

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/ir"
)

// slotsPool 复用桥调用的槽位数组，避免每次 Go→JIT 调用都堆分配 128B
var slotsPool = sync.Pool{New: func() any { return new([binding.BridgeMaxSlots]uint64) }}

// Func 按名字取出 JIT 函数并包装为真实 Go 函数值（reflect.MakeFunc）。
// 调用闭包把实参装箱到固定槽位，经 C 桥调用按签名缓存的 IR 适配器。
func (j *LLJIT) Func[F any](name string) (F, error) {
	var zero F
	const op = "jit.LLJIT.Func"
	j.check(op)
	ft := reflect.TypeOf((*F)(nil)).Elem()
	if ft.Kind() != reflect.Func {
		return zero, errBridge("not a function type: " + ft.String())
	}
	j.checkSymbolSig(op, name, ft)
	e, err := j.adapterFor(ft)
	if err != nil {
		return zero, err
	}
	addr, err := j.Lookup(name)
	if err != nil {
		return zero, err
	}

	fn := reflect.MakeFunc(ft, func(args []reflect.Value) []reflect.Value {
		sp := slotsPool.Get().(*[binding.BridgeMaxSlots]uint64)
		slots := sp[:]
		for i, a := range args {
			slots[i] = boxValue(a)
		}
		raw := binding.BridgeCall(e.adapter, addr, slots)
		for i := range slots {
			slots[i] = 0
		}
		slotsPool.Put(sp)
		if ft.NumOut() == 0 {
			return nil
		}
		return []reflect.Value{unboxValue(ft.Out(0), raw)}
	})
	return fn.Interface().(F), nil
}

// MapFunc 把宿主 Go 函数注册为 JIT 符号 name（native→Go 方向）。
// 生成真实签名的 IR 包装体：装箱参数 → callGoChannel(idx, slots) → 解箱返回值。
func (j *LLJIT) MapFunc[F any](name string, f F) error {
	const op = "jit.LLJIT.MapFunc"
	j.check(op)
	ft := reflect.TypeOf((*F)(nil)).Elem()
	if ft.Kind() != reflect.Func {
		return errBridge("not a function type: " + ft.String())
	}
	j.checkSymbolSig(op, name, ft)
	if _, _, err := checkBridgeFunc(ft); err != nil {
		return err
	}
	if err := j.ensureGoChannel(); err != nil {
		return err
	}
	idx := registerGoFunc(reflect.ValueOf(f))
	return j.compileWrapper(name, ft, idx)
}

// ensureGoChannel 把 C 通道函数挂接为 JIT 符号 callGoChannel（仅一次）
func (j *LLJIT) ensureGoChannel() error {
	j.channelOnce.Do(func() {
		j.channelErr = j.MapSymbol("callGoChannel", binding.BridgeGoChannelAddr())
	})
	if j.channelErr != nil {
		return llvm.WrapError(llvm.ErrJIT, "jit.LLJIT.ensureGoChannel", j.channelErr)
	}
	return nil
}

// compileWrapper 生成 native→Go 包装体：
//
//	<name>(params...) -> ret   // 真实 native ABI 签名
//
// 函数体把参数装箱到 slots，调用 callGoChannel(idx, slots)，再解箱返回值
func (j *LLJIT) compileWrapper(name string, ft reflect.Type, idx int64) error {
	const op = "jit.compileWrapper"
	params, ret, err := checkBridgeFunc(ft)
	if err != nil {
		return err
	}

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "bridge_wrapper")
	sig, _, err := llvm.FnSignatureOfGo(ctx, ft)
	if err != nil {
		ctx.Close()
		return llvm.WrapError(llvm.ErrUnsupported, op, err)
	}
	fn := m.NewFunction(name, sig)
	fn.SetLinkage(llvm.LinkageExternal)

	i64 := ctx.Int(64)
	ptr := ctx.Ptr(0)
	channelSig := ctx.Fn(i64, []llvm.AnyType{i64, ptr}, false)
	channel := m.NewFunction("callGoChannel", channelSig)

	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	slotArr := ctx.Array(i64, uint64(binding.BridgeMaxSlots))
	slots := b.Alloca(slotArr, "").Value
	for i, k := range params {
		raw := packSlot(b, ctx, fn.Param(uint(i)).Dyn(), i64, k)
		dst := b.GEP(slotArr, slots, []llvm.ValueRef[llvm.IntT]{ctx.ConstInt(i64, 0), ctx.ConstInt(i64, uint64(i))}, "")
		b.Store(raw, dst)
	}

	res := b.Call[llvm.IntT](channel.Value, []llvm.AnyValue{ctx.ConstInt(i64, uint64(idx)), slots}, "").Value
	if ret == slotUnsupported {
		b.RetVoid()
	} else {
		b.Ret(unpackSlot(b, ctx, res, sig.Return(), ret))
	}
	if err := b.Close(); err != nil {
		ctx.Close()
		return llvm.WrapError(llvm.ErrInternal, op, err)
	}
	if err := m.Verify(); err != nil {
		ctx.Close()
		return llvm.WrapError(llvm.ErrInternal, op, err)
	}
	return j.AddIRModule(m)
}

// RunMain 以 C 的 main(argc, argv, envp) 约定调用 JIT 中的 main，返回退出码
func (j *LLJIT) RunMain(args []string) (int32, error) {
	const op = "jit.LLJIT.RunMain"
	j.check(op)
	mainType := reflect.TypeOf((func(int32, unsafe.Pointer, unsafe.Pointer) int32)(nil))
	e, err := j.adapterFor(mainType)
	if err != nil {
		return 0, err
	}
	addr, err := j.Lookup("main")
	if err != nil {
		return 0, err
	}

	argv := binding.NewCStringArray(args)
	defer argv.Free()
	var slots [binding.BridgeMaxSlots]uint64
	slots[0] = uint64(len(args))
	slots[1] = uint64(uintptr(unsafe.Pointer(argv.Ptr())))
	slots[2] = 0
	raw := binding.BridgeCall(e.adapter, addr, slots[:])
	return int32(raw), nil
}
