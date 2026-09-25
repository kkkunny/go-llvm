package llvm

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm/internal/binding"
)

func TestValueTypeAndAs(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	v := ctx.ConstSInt(ctx.Int(32), 42)
	if got := v.Type().String(); got != "i32" {
		t.Fatalf("Type() = %q, want i32", got)
	}
	if !strings.Contains(v.String(), "i32 42") {
		t.Fatalf("String() = %q", v.String())
	}

	got := v.MustAs[IntT]()
	if AsIntType(got.Type()).Bits() != 32 {
		t.Fatalf("MustAs[IntT] = %v", got)
	}
	if _, err := v.Dyn().As[FloatT](); err == nil || err.(*Error).Reason != ErrTypeMismatch {
		t.Fatalf("want ErrTypeMismatch, got %v", err)
	}
	if _, err := v.As[DynT](); err != nil {
		t.Fatalf("As[DynT] should always succeed: %v", err)
	}
}

func TestValueName(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	v := ctx.ConstSInt(ctx.Int(32), 1)
	if v.Name() != "" {
		t.Fatalf("constant name should be empty, got %q", v.Name())
	}
	if v.IsNil() {
		t.Fatal("constant should not be nil")
	}
	if !v.IsConstant() {
		t.Fatal("constant should report IsConstant")
	}
}

func TestValueAlive(t *testing.T) {
	ctx := NewContext()
	life := NewLifetime()
	v := NewValue[IntT](ctx, life, ctx.ConstSInt(ctx.Int(32), 1).Ref())
	if !v.Alive() {
		t.Fatal("value should be alive")
	}
	life.Kill()
	if v.Alive() {
		t.Fatal("value should be dead after lifetime kill")
	}
	ctx.Close()
	if v.Alive() {
		t.Fatal("value should be dead after context close")
	}
}

func TestValueOfDispatch(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	life := NewLifetime()
	ref := ctx.ConstSInt(ctx.Int(32), 7).Ref()
	v := ValueOf(ctx, life, ref)
	if got := v.Type().String(); got != "i32" {
		t.Fatalf("ValueOf Type() = %q, want i32", got)
	}
	if !strings.Contains(v.String(), "i32 7") {
		t.Fatalf("ValueOf String() = %q", v.String())
	}

	if v := ValueOf(ctx, life, ctx.ConstFloat(ctx.Float(FloatDouble), 1.5).Ref()); v.Type().String() != "double" {
		t.Fatalf("ValueOf float Type() = %q", v.Type().String())
	}
}

func TestCheckValuesFloor(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()
	i32 := ctx.Int(32)

	// 跨 Context 值
	other := NewContext()
	defer other.Close()
	if err := Catch(func() {
		ctx.CheckValues("llvm.Test.CheckValues", other.ConstSInt(other.Int(32), 1))
	}); err == nil || err.Reason != ErrCrossContext {
		t.Fatalf("跨 Context 值应 panic ErrCrossContext, got %v", err)
	}

	// nil 句柄：公开 API 的 Ref 会先行拦截，这里直接走底层入口验证地板
	if err := Catch(func() {
		ctx.checkValueOwn("llvm.Test.CheckValues", binding.LLVMValueRef{}, nil, ctx)
	}); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("nil 句柄应 panic ErrInvalidArg, got %v", err)
	}

	// 已关闭 Context：绕过 Ref 的重复校验，直接走底层入口
	closed := NewContext()
	cv := closed.ConstSInt(closed.Int(32), 1)
	_ = closed.Close()
	if err := Catch(func() {
		closed.checkValueOwn("llvm.Test.CheckValues", cv.RawRef(), cv.Lifetime(), cv.Context())
	}); err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("已关闭 Context 应 panic ErrUseAfterFree, got %v", err)
	}

	// 生命周期令牌已结束
	life := NewLifetime()
	lv := NewValue[IntT](ctx, life, i32.Const(1).RawRef())
	life.Kill()
	if err := Catch(func() {
		ctx.checkValueOwn("llvm.Test.CheckValues", lv.RawRef(), life, lv.Context())
	}); err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("生命周期结束后应 panic ErrUseAfterFree, got %v", err)
	}
}

func TestValueCrashFloor(t *testing.T) {
	// nil 句柄
	var nilVal Value[DynT]
	if err := Catch(func() { nilVal.Ref() }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("nil 值 Ref() 应 panic ErrInvalidArg, got %v", err)
	}

	ctx := NewContext()
	defer ctx.Close()
	ref := ctx.ConstSInt(ctx.Int(32), 1).RawRef()

	// ctx 为 nil
	if err := Catch(func() { (Value[IntT]{ref: ref}).Ref() }); err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("ctx 为 nil 的值 Ref() 应 panic ErrUseAfterFree, got %v", err)
	}

	// 生命周期令牌已结束
	life := NewLifetime()
	v := NewValue[IntT](ctx, life, ref)
	life.Kill()
	if err := Catch(func() { v.Ref() }); err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("生命周期结束后 Ref() 应 panic ErrUseAfterFree, got %v", err)
	}

	// MustAs 种类不符
	if err := Catch(func() { ctx.ConstSInt(ctx.Int(32), 1).MustAs[FloatT]() }); err == nil || err.Reason != ErrTypeMismatch {
		t.Fatalf("MustAs 种类不符应 panic ErrTypeMismatch, got %v", err)
	}
}

func TestValueNil(t *testing.T) {
	var v Value[DynT]
	if !v.IsNil() {
		t.Fatal("zero value should be nil")
	}
	if v.String() != "<nil>" {
		t.Fatalf("nil value String() = %q", v.String())
	}
	if v.Alive() {
		t.Fatal("nil value should not be alive")
	}
}
