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

func TestTypeBridging(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	i32 := ctx.Int(32)
	if i32.RawRef() != i32.Ref() {
		t.Fatal("RawRef 应与 Ref 指向同一句柄")
	}
	nt := IntType{NewType[IntT](ctx, i32.RawRef())}
	if nt.Bits() != 32 {
		t.Fatalf("NewType[IntT] = %s", nt)
	}

	refs := AnyTypesToRefs([]AnyType{ctx.Int(32), ctx.Float(FloatDouble), ctx.Ptr(0)})
	if len(refs) != 3 || refs[0] != i32.RawRef() {
		t.Fatalf("AnyTypesToRefs = %v", refs)
	}
	if got := TypeOfRef(ctx, refs[1]).String(); got != "double" {
		t.Fatalf("TypeOfRef(refs[1]) = %q", got)
	}
}

func TestAsTypeRoles(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	if got := AsVoidType(ctx.Void()).String(); got != "void" {
		t.Fatalf("AsVoidType = %q", got)
	}
	if got := AsIntType(ctx.Int(32)).Bits(); got != 32 {
		t.Fatalf("AsIntType = %d", got)
	}
	if got := AsFloatType(ctx.Float(FloatSingle)).Kind(); got != FloatSingle {
		t.Fatalf("AsFloatType = %v", got)
	}
	if got := AsPtrType(ctx.Ptr(0)).Addrspace(); got != 0 {
		t.Fatalf("AsPtrType = %d", got)
	}
	st := ctx.Struct([]AnyType{ctx.Int(32)}, false)
	if got := AsStructType(st).Elems(); len(got) != 1 || got[0].String() != "i32" {
		t.Fatalf("AsStructType = %v", got)
	}
	if got := AsArrayType(ctx.Array(ctx.Int(8), 4)).Len(); got != 4 {
		t.Fatalf("AsArrayType = %d", got)
	}
	if got := AsVecType(ctx.Vec(ctx.Int(8), 2)).Len(); got != 2 {
		t.Fatalf("AsVecType = %d", got)
	}
	if AsFnType(ctx.Fn(ctx.Void(), nil, false)).IsVarArg() {
		t.Fatal("AsFnType 应保留非变参签名")
	}

	// nil 类型与种类不符都应 panic
	if err := Catch(func() { AsIntType(nil) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("AsIntType(nil) 应 panic ErrInvalidArg, got %v", err)
	}
	if err := Catch(func() { AsIntType(ctx.Float(FloatSingle)) }); err == nil || err.Reason != ErrTypeMismatch {
		t.Fatalf("AsIntType(float) 应 panic ErrTypeMismatch, got %v", err)
	}
}

func TestPtrTypeOpaque(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()
	if !ctx.Ptr(0).IsOpaque() {
		t.Fatal("LLVM 22 的指针类型应为不透明")
	}
}

func TestTypeNilFloor(t *testing.T) {
	var t0 Type[IntT]
	if got := t0.String(); got != "<nil>" {
		t.Fatalf("nil 类型 String() = %q", got)
	}
	if err := Catch(func() { t0.Ref() }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("nil 类型 Ref() 应 panic ErrInvalidArg, got %v", err)
	}
	if err := Catch(func() { t0.Check("llvm.Test.Type.Check") }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("nil 类型 Check() 应 panic ErrInvalidArg, got %v", err)
	}

	ctx := NewContext()
	defer ctx.Close()
	if err := Catch(func() { ctx.Float(FloatDouble).MustAs[IntT]() }); err == nil || err.Reason != ErrTypeMismatch {
		t.Fatalf("MustAs 种类不符应 panic ErrTypeMismatch, got %v", err)
	}
	// 已释放 Context 上的类型句柄走同一地板
	dead := NewContext()
	ty := dead.Int(32)
	_ = dead.Close()
	if err := Catch(func() { ty.Ref() }); err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("已释放 Context 的类型应 panic ErrUseAfterFree, got %v", err)
	}
}

func TestCheckTypeValidation(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	if err := Catch(func() { ctx.CheckType("llvm.Test.CheckType", nil) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("nil 类型应 panic ErrInvalidArg, got %v", err)
	}
	if err := Catch(func() { ctx.CheckType("llvm.Test.CheckType", Type[IntT]{}) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("nil 句柄应 panic ErrInvalidArg, got %v", err)
	}

	ctx2 := NewContext()
	defer ctx2.Close()
	if err := Catch(func() { ctx.CheckType("llvm.Test.CheckType", ctx2.Int(32)) }); err == nil || err.Reason != ErrCrossContext {
		t.Fatalf("跨 Context 类型应 panic ErrCrossContext, got %v", err)
	}
}

func TestFloatKinds(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	want := []struct {
		kind FloatKind
		str  string
	}{
		{FloatHalf, "half"},
		{FloatBFloat, "bfloat"},
		{FloatSingle, "float"},
		{FloatDouble, "double"},
		{FloatX86FP80, "x86_fp80"},
		{FloatFP128, "fp128"},
		{FloatPPCFP128, "ppc_fp128"},
	}
	for _, c := range want {
		got := ctx.Float(c.kind)
		if s := got.String(); s != c.str {
			t.Errorf("Float(%v) = %q, want %q", c.kind, s, c.str)
		}
		if k := got.Kind(); k != c.kind {
			t.Errorf("Float(%v).Kind() = %v", c.kind, k)
		}
	}
	if err := Catch(func() { ctx.Float(FloatKind(99)) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("未知浮点种类应 panic ErrInvalidArg, got %v", err)
	}
}

func TestKindOfTypeAndName(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	// metadata 值的内建类型应归类为 MetaT
	if _, err := ctx.MDString("x").Value().Type().DynType().As[MetaT](); err != nil {
		t.Fatalf("metadata 类型应归类为 MetaT: %v", err)
	}
	if _, err := ctx.Int(32).DynType().As[MetaT](); err == nil || err.(*Error).Reason != ErrTypeMismatch {
		t.Fatalf("i32 不应归类为 MetaT, got %v", err)
	}

	// kindName 错误文本辅助：全部种类
	cases := []struct {
		kind Kind
		want string
	}{
		{VoidT{}, "void"}, {IntT{}, "int"}, {FloatT{}, "float"}, {PtrT{}, "pointer"},
		{StructT{}, "struct"}, {ArrayT{}, "array"}, {VecT{}, "vector"}, {FnT{}, "function"},
		{LabelT{}, "label"}, {MetaT{}, "metadata"}, {TokenT{}, "token"}, {DynT{}, "dyn"},
	}
	for _, c := range cases {
		if got := kindName(c.kind); got != c.want {
			t.Errorf("kindName(%T) = %q, want %q", c.kind, got, c.want)
		}
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
