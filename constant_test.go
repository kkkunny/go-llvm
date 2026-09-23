package llvm

import (
	"strings"
	"testing"
)

func TestConstInt(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	c := ctx.ConstSInt(ctx.Int(32), 42)
	if c.SignedValue() != 42 || c.UnsignedValue() != 42 || c.IsNegative() {
		t.Fatalf("ConstInt(42): %d/%d/%v", c.SignedValue(), c.UnsignedValue(), c.IsNegative())
	}

	neg := ctx.ConstSInt(ctx.Int(32), -1) // -1
	if neg.SignedValue() != -1 || !neg.IsNegative() {
		t.Fatalf("ConstInt(-1): %d/%v", neg.SignedValue(), neg.IsNegative())
	}

	wide := ctx.ConstSInt(ctx.Int(64), -1<<63) // INT64_MIN
	if wide.SignedValue() != -1<<63 {
		t.Fatalf("ConstInt(INT64_MIN): %d", wide.SignedValue())
	}

	if got := ctx.ConstBool(true).String(); got != "i1 true" {
		t.Fatalf("ConstBool(true) = %q", got)
	}
	if got := ctx.ConstBool(false).String(); got != "i1 false" {
		t.Fatalf("ConstBool(false) = %q", got)
	}

	if got := ctx.ConstIntOfString(ctx.Int(32), "ff", 16).String(); got != "i32 255" {
		t.Fatalf("ConstIntOfString = %q", got)
	}
}

func TestConstSugar(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()
	i32 := ctx.Int(32)
	if got := i32.Const(42).SignedValue(); got != 42 {
		t.Fatalf("i32.Const(42) = %d", got)
	}
	if got := i32.ConstS(-1).SignedValue(); got != -1 {
		t.Fatalf("i32.ConstS(-1) = %d", got)
	}
	f64 := ctx.Float(FloatDouble)
	if got := f64.Const(3.5).FloatValue(); got != 3.5 {
		t.Fatalf("f64.Const(3.5) = %v", got)
	}
}

func TestConstFloat(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	c := ctx.ConstFloat(ctx.Float(FloatDouble), 1.5)
	if c.FloatValue() != 1.5 {
		t.Fatalf("FloatValue = %v", c.FloatValue())
	}
	if got := c.String(); got != "double 1.500000e+00" {
		t.Fatalf("ConstFloat = %q", got)
	}
}

func TestConstNullZero(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	i32 := ctx.Int(32)
	if got := ctx.ConstNull[DynT](i32.DynType()).String(); got != "i32 0" {
		t.Fatalf("ConstNull = %q", got)
	}
	if got := ctx.ConstNull(i32).String(); got != "i32 0" {
		t.Fatalf("ConstNull generic method = %q", got)
	}
	if got := ctx.ConstNull(ctx.Ptr(0)).String(); got != "ptr null" {
		t.Fatalf("ConstNull ptr = %q", got)
	}

	st := ctx.Struct([]AnyType{i32, ctx.Int(64)}, false)
	if got := ctx.ConstZero(st).String(); !strings.Contains(got, "zeroinitializer") {
		t.Fatalf("ConstZero struct = %q", got)
	}
	if got := ctx.ConstZero(i32).String(); got != "i32 0" {
		t.Fatalf("ConstZero int = %q", got)
	}
}

func TestConstString(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	withNull := ctx.ConstString("hi", true)
	if got := withNull.Type().String(); got != "[3 x i8]" {
		t.Fatalf("ConstString with null terminator = %q", got)
	}
	if got := withNull.String(); !strings.Contains(got, `c"hi\00"`) {
		t.Fatalf("ConstString = %q", got)
	}

	noNull := ctx.ConstString("hi", false)
	if got := noNull.Type().String(); got != "[2 x i8]" {
		t.Fatalf("ConstString without null terminator = %q", got)
	}
}

func TestConstArrayStruct(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	i32 := ctx.Int(32)
	arr := ctx.ConstArray(i32, ctx.ConstInt(i32, 1), ctx.ConstInt(i32, 2))
	if got := arr.Type().String(); got != "[2 x i32]" {
		t.Fatalf("ConstArray type = %q", got)
	}
	if got := arr.String(); !strings.HasSuffix(got, "[i32 1, i32 2]") {
		t.Fatalf("ConstArray = %q", got)
	}

	st := ctx.ConstStruct(false, ctx.ConstInt(i32, 1), ctx.ConstFloat(ctx.Float(FloatDouble), 2.0))
	if got := st.Type().String(); got != "{ i32, double }" {
		t.Fatalf("ConstStruct type = %q", got)
	}
	if got := st.String(); !strings.HasSuffix(got, "{ i32 1, double 2.000000e+00 }") {
		t.Fatalf("ConstStruct = %q", got)
	}

	named := ctx.NamedStruct("Pair")
	named.SetBody([]AnyType{i32, i32}, false)
	ns := ctx.ConstNamedStruct(named, ctx.ConstInt(i32, 3), ctx.ConstInt(i32, 4))
	if got := ns.Type().String(); got != "%Pair = type { i32, i32 }" {
		t.Fatalf("ConstNamedStruct type = %q", got)
	}
	if got := ns.String(); !strings.HasSuffix(got, "{ i32 3, i32 4 }") {
		t.Fatalf("ConstNamedStruct = %q", got)
	}
}

func TestConstGEP(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	i32 := ctx.Int(32)
	i64 := ctx.Int(64)
	nullPtr := ctx.ConstNull(ctx.Ptr(0))

	gep := ctx.ConstGEP(i32, nullPtr, false, ctx.ConstInt(i64, 1))
	if got := gep.Type().String(); got != "ptr" {
		t.Fatalf("ConstGEP type = %q", got)
	}
	if got := gep.String(); !strings.Contains(got, "getelementptr (i32, ptr null, i64 1)") {
		t.Fatalf("ConstGEP = %q", got)
	}

	inbounds := ctx.ConstGEP(i32, nullPtr, true, ctx.ConstInt(i64, 2))
	if got := inbounds.String(); !strings.Contains(got, "getelementptr inbounds (i32, ptr null, i64 2)") {
		t.Fatalf("ConstInBoundsGEP = %q", got)
	}
}

func TestConstMismatchPanic(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	i32 := ctx.Int(32)
	f32 := ctx.Float(FloatSingle)

	err := Catch(func() {
		ctx.ConstArray(i32, ctx.ConstFloat(f32, 1.0))
	})
	if err == nil || err.Reason != ErrTypeMismatch {
		t.Fatalf("want ErrTypeMismatch, got %v", err)
	}

	ctx2 := NewContext()
	defer ctx2.Close()
	err = Catch(func() {
		ctx.ConstArray(i32, ctx2.ConstInt(i32, 1))
	})
	if err == nil || err.Reason != ErrCrossContext {
		t.Fatalf("want ErrCrossContext, got %v", err)
	}
}

func TestConstVector(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	i32 := ctx.Int(32)
	v := ctx.ConstVector(i32, ctx.ConstInt(i32, 1).Value, ctx.ConstInt(i32, 2).Value)
	if got := v.String(); !strings.Contains(got, "<i32 1, i32 2>") {
		t.Fatalf("const vector = %s", got)
	}

	// 元素类型不符
	if err := Catch(func() {
		ctx.ConstVector(i32, ctx.ConstFloat(ctx.Float(FloatDouble), 1).Value)
	}); err == nil || err.Reason != ErrTypeMismatch {
		t.Fatalf("elem type mismatch should panic ErrTypeMismatch, got %v", err)
	}
}
