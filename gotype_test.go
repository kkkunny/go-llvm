package llvm

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestTypeOf(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	cases := []struct {
		name string
		got  func() (Type[DynT], error)
		want string
	}{
		{"bool", func() (Type[DynT], error) { return TypeOf[bool](ctx) }, "i1"},
		{"int8", func() (Type[DynT], error) { return TypeOf[int8](ctx) }, "i8"},
		{"int16", func() (Type[DynT], error) { return TypeOf[int16](ctx) }, "i16"},
		{"int32", func() (Type[DynT], error) { return TypeOf[int32](ctx) }, "i32"},
		{"int64", func() (Type[DynT], error) { return TypeOf[int64](ctx) }, "i64"},
		{"uint64", func() (Type[DynT], error) { return TypeOf[uint64](ctx) }, "i64"},
		{"float32", func() (Type[DynT], error) { return TypeOf[float32](ctx) }, "float"},
		{"float64", func() (Type[DynT], error) { return TypeOf[float64](ctx) }, "double"},
		{"ptr", func() (Type[DynT], error) { return TypeOf[*int32](ctx) }, "ptr"},
		{"unsafe.Pointer", func() (Type[DynT], error) { return TypeOf[unsafe.Pointer](ctx) }, "ptr"},
		{"array", func() (Type[DynT], error) { return TypeOf[[4]int32](ctx) }, "[4 x i32]"},
		{"struct", func() (Type[DynT], error) {
			return TypeOf[struct {
				A int32
				B float64
			}](ctx)
		}, "{ i32, double }"},
	}
	for _, c := range cases {
		ty, err := c.got()
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if got := ty.String(); got != c.want {
			t.Fatalf("%s: got %q, want %q", c.name, got, c.want)
		}
	}

	unsupported := []struct {
		name string
		got  func() (Type[DynT], error)
	}{
		{"string", func() (Type[DynT], error) { return TypeOf[string](ctx) }},
		{"slice", func() (Type[DynT], error) { return TypeOf[[]int32](ctx) }},
		{"map", func() (Type[DynT], error) { return TypeOf[map[int32]int32](ctx) }},
		{"chan", func() (Type[DynT], error) { return TypeOf[chan int32](ctx) }},
		{"func", func() (Type[DynT], error) { return TypeOf[func()](ctx) }},
		{"interface", func() (Type[DynT], error) { return TypeOf[any](ctx) }},
		{"complex", func() (Type[DynT], error) { return TypeOf[complex64](ctx) }},
		{"struct-with-string", func() (Type[DynT], error) {
			return TypeOf[struct{ S string }](ctx)
		}},
	}
	for _, c := range unsupported {
		_, err := c.got()
		if err == nil || err.(*Error).Reason != ErrUnsupported {
			t.Fatalf("%s: want ErrUnsupported, got %v", c.name, err)
		}
	}
}

func TestConstOf(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	v, err := ConstOf(ctx, int32(42))
	if err != nil || v.String() != "i32 42" {
		t.Fatalf("ConstOf(int32) = %v, %v", v, err)
	}
	v, err = ConstOf(ctx, uint8(255))
	if err != nil || v.String() != "i8 -1" { // LLVM 以有符号打印 i8 位模式
		t.Fatalf("ConstOf(uint8) = %v, %v", v, err)
	}
	v, err = ConstOf(ctx, uint32(4294967295))
	if err != nil || v.String() != "i32 -1" {
		t.Fatalf("ConstOf(uint32) = %v, %v", v, err)
	}
	v, err = ConstOf(ctx, float64(1.5))
	if err != nil || v.String() != "double 1.500000e+00" {
		t.Fatalf("ConstOf(float64) = %v, %v", v, err)
	}
	v, err = ConstOf(ctx, true)
	if err != nil || v.String() != "i1 true" {
		t.Fatalf("ConstOf(true) = %v, %v", v, err)
	}

	var nilPtr *int32
	v, err = ConstOf(ctx, nilPtr)
	if err != nil || v.String() != "ptr null" {
		t.Fatalf("ConstOf(nil ptr) = %v, %v", v, err)
	}

	x := int32(7)
	v, err = ConstOf(ctx, &x)
	if err != nil || v.Type().String() != "ptr" {
		t.Fatalf("ConstOf(ptr) = %v, %v", v, err)
	}

	v, err = ConstOf(ctx, struct {
		A int32
		B float64
	}{A: 1, B: 2})
	if err != nil || v.Type().String() != "{ i32, double }" {
		t.Fatalf("ConstOf(struct) = %v, %v", v, err)
	}

	v, err = ConstOf(ctx, [2]int32{1, 2})
	if err != nil || v.Type().String() != "[2 x i32]" {
		t.Fatalf("ConstOf(array) = %v, %v", v, err)
	}

	if _, err := ConstOf(ctx, "nope"); err == nil || err.(*Error).Reason != ErrUnsupported {
		t.Fatalf("ConstOf(string) should be unsupported, got %v", err)
	}
}

func TestFnSignatureOf(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	sig, goTy, err := FnSignatureOf[func(int32, float64) int32](ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got := sig.String(); got != "i32 (i32, double)" {
		t.Fatalf("signature = %q", got)
	}
	if goTy.Kind() != reflect.Func || goTy.String() != "func(int32, float64) int32" {
		t.Fatalf("go type = %v", goTy)
	}

	sig, _, err = FnSignatureOf[func()](ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got := sig.String(); got != "void ()" {
		t.Fatalf("void signature = %q", got)
	}

	if _, _, err := FnSignatureOf[int32](ctx); err == nil || err.(*Error).Reason != ErrUnsupported {
		t.Fatalf("non-func should be unsupported, got %v", err)
	}
	if _, _, err := FnSignatureOf[func(int32) (int32, int32)](ctx); err == nil || err.(*Error).Reason != ErrUnsupported {
		t.Fatalf("multi-result should be unsupported, got %v", err)
	}
	if _, _, err := FnSignatureOf[func(...int32)](ctx); err == nil || err.(*Error).Reason != ErrUnsupported {
		t.Fatalf("variadic should be unsupported, got %v", err)
	}
	if _, _, err := FnSignatureOf[func(string)](ctx); err == nil || err.(*Error).Reason != ErrUnsupported {
		t.Fatalf("unsupported param type should be unsupported, got %v", err)
	}
}
