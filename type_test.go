package llvm

import "testing"

func TestTypeConstruct(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	if got := ctx.Int(32).String(); got != "i32" {
		t.Fatalf("Int(32) = %q, want i32", got)
	}
	if got := ctx.Int(32).Bits(); got != 32 {
		t.Fatalf("Int(32).Bits() = %d, want 32", got)
	}
	if got := ctx.Bool().String(); got != "i1" {
		t.Fatalf("Bool() = %q, want i1", got)
	}
	if got := ctx.Void().String(); got != "void" {
		t.Fatalf("Void() = %q, want void", got)
	}
	if got := ctx.Float(FloatDouble).String(); got != "double" {
		t.Fatalf("Float(FloatDouble) = %q, want double", got)
	}
	if got := ctx.Float(FloatDouble).Kind(); got != FloatDouble {
		t.Fatalf("Kind() = %v, want FloatDouble", got)
	}
	if got := ctx.Ptr(0).String(); got != "ptr" {
		t.Fatalf("Ptr(0) = %q, want ptr", got)
	}
	if got := ctx.Ptr(0).Addrspace(); got != 0 {
		t.Fatalf("Addrspace() = %d, want 0", got)
	}
	if !ctx.Int(32).IsSized() {
		t.Fatal("i32 should be sized")
	}
	if ctx.Void().IsSized() {
		t.Fatal("void should not be sized")
	}
	if ctx.Int(32).IsNil() {
		t.Fatal("i32 should not be nil")
	}
}

func TestFunctionType(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	ft := ctx.Fn(ctx.Bool(), []AnyType{ctx.Int(32), ctx.Int(64)}, false)
	if ft.IsVarArg() {
		t.Fatal("should not be vararg")
	}
	if got := ft.Return().String(); got != "i1" {
		t.Fatalf("Return() = %q, want i1", got)
	}
	params := ft.Params()
	if len(params) != 2 || params[0].String() != "i32" || params[1].String() != "i64" {
		t.Fatalf("Params() = %v", params)
	}
	if got := ft.String(); got != "i1 (i32, i64)" {
		t.Fatalf("Fn = %q", got)
	}

	vararg := ctx.Fn(ctx.Void(), nil, true)
	if !vararg.IsVarArg() || len(vararg.Params()) != 0 {
		t.Fatal("vararg function type mismatch")
	}
}

func TestAggregateTypes(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	st := ctx.Struct([]AnyType{ctx.Int(32), ctx.Int(64)}, false)
	if got := st.String(); got != "{ i32, i64 }" {
		t.Fatalf("Struct = %q", got)
	}
	if st.IsPacked() || st.IsOpaque() {
		t.Fatal("literal struct should be sized and not packed")
	}
	if elems := st.Elems(); len(elems) != 2 || elems[1].String() != "i64" {
		t.Fatalf("Elems() = %v", elems)
	}

	named := ctx.NamedStruct("Pair")
	if !named.IsOpaque() || named.Name() != "Pair" {
		t.Fatal("NamedStruct should start opaque")
	}
	named.SetBody([]AnyType{ctx.Int(32), ctx.Int(32)}, true)
	if named.IsOpaque() || !named.IsPacked() || len(named.Elems()) != 2 {
		t.Fatalf("SetBody failed: %s", named)
	}

	arr := ctx.Array(ctx.Int(32), 4)
	if got := arr.String(); got != "[4 x i32]" {
		t.Fatalf("Array = %q", got)
	}
	if arr.Len() != 4 || arr.Elem().String() != "i32" {
		t.Fatalf("Array Elem/Len = %v/%d", arr.Elem(), arr.Len())
	}

	vec := ctx.Vec(ctx.Float(FloatSingle), 4)
	if got := vec.String(); got != "<4 x float>" {
		t.Fatalf("Vec = %q", got)
	}
	if vec.Len() != 4 || vec.Elem().String() != "float" {
		t.Fatalf("Vec Elem/Len = %v/%d", vec.Elem(), vec.Len())
	}
}

func TestTypeEqual(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	if !ctx.Int(32).Equal(ctx.Int(32)) {
		t.Fatal("same type should be equal")
	}
	if ctx.Int(32).Equal(ctx.Int(64)) {
		t.Fatal("different types should not be equal")
	}
	if ctx.Int(32).Equal(ctx.Float(FloatSingle)) {
		t.Fatal("int and float should not be equal")
	}
	if ctx.Int(32).Equal(nil) {
		t.Fatal("nil should not be equal")
	}
}

func TestTypeAsMismatch(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	if _, err := ctx.Int(32).DynType().As[FloatT](); err == nil || err.(*Error).Reason != ErrTypeMismatch {
		t.Fatalf("want ErrTypeMismatch, got %v", err)
	}
	got := ctx.Int(32).DynType().MustAs[IntT]()
	if AsIntType(got).Bits() != 32 {
		t.Fatalf("MustAs[IntT] = %v", got)
	}
	if _, err := ctx.Int(32).As[DynT](); err != nil {
		t.Fatalf("As[DynT] should always succeed: %v", err)
	}
}

func TestTypeCrossContextPanic(t *testing.T) {
	ctx1 := NewContext()
	defer ctx1.Close()
	ctx2 := NewContext()
	defer ctx2.Close()

	err := Catch(func() {
		ctx1.Struct([]AnyType{ctx2.Int(32)}, false)
	})
	if err == nil || err.Reason != ErrCrossContext {
		t.Fatalf("want ErrCrossContext, got %v", err)
	}

	err = Catch(func() {
		ctx1.Fn(ctx1.Int(32), []AnyType{ctx2.Int(32)}, false)
	})
	if err == nil || err.Reason != ErrCrossContext {
		t.Fatalf("want ErrCrossContext, got %v", err)
	}
}
