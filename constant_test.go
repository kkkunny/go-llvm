package llvm

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm/internal/errs"
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

func TestTypeNull(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	if got := ctx.Int(32).Null().String(); got != "i32 0" {
		t.Fatalf("IntType.Null() = %q", got)
	}
	if got := ctx.Ptr(0).Null().String(); got != "ptr null" {
		t.Fatalf("PtrType.Null() = %q", got)
	}
	if got := ctx.Int(32).DynType().Null().Type().String(); got != "i32" {
		t.Fatalf("Type[DynT].Null() 类型 = %q", got)
	}
	if got := ctx.Struct([]AnyType{ctx.Int(32), ctx.Int(64)}, false).Null().String(); !strings.Contains(got, "zeroinitializer") {
		t.Fatalf("StructType.Null() = %q", got)
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

func TestConstNamedStructMismatch(t *testing.T) {
	requireDebug(t)

	ctx := NewContext()
	defer ctx.Close()
	i32 := ctx.Int(32)
	named := ctx.NamedStruct("Pair")
	named.SetBody([]AnyType{i32, i32}, false)

	// 元素个数不符
	if err := errs.Catch(func() {
		ctx.ConstNamedStruct(named, ctx.ConstInt(i32, 1))
	}); err == nil || err.Reason != ErrTypeMismatch {
		t.Fatalf("元素个数不符应 panic ErrTypeMismatch, got %v", err)
	}
	// 元素类型不符
	if err := errs.Catch(func() {
		ctx.ConstNamedStruct(named, ctx.ConstInt(i32, 1), ctx.ConstFloat(ctx.Float(FloatDouble), 1))
	}); err == nil || err.Reason != ErrTypeMismatch {
		t.Fatalf("元素类型不符应 panic ErrTypeMismatch, got %v", err)
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

func TestSizeOfAlignOf(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	i64 := ctx.Int(64)
	size := ctx.SizeOf(i64)
	if got := size.Type().String(); got != "i64" {
		t.Fatalf("SizeOf(i64) type = %q, want i64", got)
	}
	if !size.IsConstant() {
		t.Fatalf("SizeOf(i64) 应为常量表达式, got %s", size.String())
	}
	if got := size.String(); !strings.Contains(got, "getelementptr") || !strings.Contains(got, "ptrtoint") {
		t.Fatalf("SizeOf(i64) = %q", got)
	}

	align := ctx.AlignOf(i64)
	if got := align.Type().String(); got != "i64" {
		t.Fatalf("AlignOf(i64) type = %q, want i64", got)
	}
	if !align.IsConstant() {
		t.Fatalf("AlignOf(i64) 应为常量表达式, got %s", align.String())
	}
	if got := align.String(); !strings.Contains(got, "getelementptr") || !strings.Contains(got, "ptrtoint") {
		t.Fatalf("AlignOf(i64) = %q", got)
	}
}

func TestSizeOfAlignOfChecks(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()
	ctx2 := NewContext()
	defer ctx2.Close()

	if err := errs.Catch(func() { ctx.SizeOf(nil) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("SizeOf(nil) 应 panic ErrInvalidArg, got %v", err)
	}
	if err := errs.Catch(func() { ctx.AlignOf(nil) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("AlignOf(nil) 应 panic ErrInvalidArg, got %v", err)
	}
	if err := errs.Catch(func() { ctx.SizeOf(Type[IntT]{}) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("SizeOf(nil 句柄) 应 panic ErrInvalidArg, got %v", err)
	}
	if err := errs.Catch(func() { ctx.AlignOf(Type[IntT]{}) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("AlignOf(nil 句柄) 应 panic ErrInvalidArg, got %v", err)
	}
	if err := errs.Catch(func() { ctx.SizeOf(ctx2.Int(32)) }); err == nil || err.Reason != ErrCrossContext {
		t.Fatalf("SizeOf(跨 Context) 应 panic ErrCrossContext, got %v", err)
	}
	if err := errs.Catch(func() { ctx.AlignOf(ctx2.Int(32)) }); err == nil || err.Reason != ErrCrossContext {
		t.Fatalf("AlignOf(跨 Context) 应 panic ErrCrossContext, got %v", err)
	}
}

func TestConstMismatchPanic(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	i32 := ctx.Int(32)
	f32 := ctx.Float(FloatSingle)

	err := errs.Catch(func() {
		ctx.ConstArray(i32, ctx.ConstFloat(f32, 1.0))
	})
	if err == nil || err.Reason != ErrTypeMismatch {
		t.Fatalf("want ErrTypeMismatch, got %v", err)
	}

	ctx2 := NewContext()
	defer ctx2.Close()
	err = errs.Catch(func() {
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
	if err := errs.Catch(func() {
		ctx.ConstVector(i32, ctx.ConstFloat(ctx.Float(FloatDouble), 1).Value)
	}); err == nil || err.Reason != ErrTypeMismatch {
		t.Fatalf("elem type mismatch should panic ErrTypeMismatch, got %v", err)
	}
}

func TestUndefPoison(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()
	i32 := ctx.Int(32)
	if u := i32.Undef(); u.String() != "i32 undef" {
		t.Fatalf("Undef = %s", u.String())
	}
	if p := i32.Poison(); p.String() != "i32 poison" {
		t.Fatalf("Poison = %s", p.String())
	}
}

func TestConstGuards(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	// 空向量常量：上游 ConstantVector::getImpl 断言 !V.empty()，NDEBUG 下解引用空指针
	if err := errs.Catch(func() { ctx.ConstVector(ctx.Int(32)) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("empty ConstVector should panic ErrInvalidArg, got %v", err)
	}
	// 非法 radix / 空串：上游 APInt 断言
	if err := errs.Catch(func() { ctx.ConstIntOfString(ctx.Int(32), "1", 3) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("bad radix should panic ErrInvalidArg, got %v", err)
	}
	if err := errs.Catch(func() { ctx.ConstIntOfString(ctx.Int(32), "", 10) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("empty integer string should panic ErrInvalidArg, got %v", err)
	}
	// 超宽值读取不再静默截断
	wide := ctx.ConstIntOfString(ctx.Int(128), "1"+strings.Repeat("0", 25), 16) // 1<<100
	if err := errs.Catch(func() { wide.UnsignedValue() }); err == nil || err.Reason != ErrUnsupported {
		t.Fatalf("over-wide UnsignedValue should panic ErrUnsupported, got %v", err)
	}
	if err := errs.Catch(func() { wide.SignedValue() }); err == nil || err.Reason != ErrUnsupported {
		t.Fatalf("over-wide SignedValue should panic ErrUnsupported, got %v", err)
	}
	// 加宽必须显式走 SExt/ZExt
	i16, i32 := ctx.Int(16), ctx.Int(32)
	if err := errs.Catch(func() { i16.Const(1).Cast(i32) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("widening Cast should panic ErrInvalidArg, got %v", err)
	}
	if err := errs.Catch(func() { i32.Const(1).SExt(i16) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("narrowing SExt should panic ErrInvalidArg, got %v", err)
	}
}
