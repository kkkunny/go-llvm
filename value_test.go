package llvm

import (
	"strings"
	"testing"
)

func TestValueTypeAndAs(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	v := ctx.ConstInt(ctx.Int(32), 42, true)
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

	v := ctx.ConstInt(ctx.Int(32), 1, true)
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
	v := NewValue[IntT](ctx, life, ctx.ConstInt(ctx.Int(32), 1, true).Ref())
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
	ref := ctx.ConstInt(ctx.Int(32), 7, true).Ref()
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
