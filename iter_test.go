package llvm

import "testing"

func TestStructAllElems(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()
	i32 := ctx.Int(32)
	st := ctx.Struct([]AnyType{i32, ctx.Float(FloatDouble)}, false)

	var got []string
	for e := range st.AllElems() {
		got = append(got, e.String())
	}
	if len(got) != 2 || got[0] != "i32" || got[1] != "double" {
		t.Fatalf("AllElems = %v", got)
	}

	// 提前 break 不 panic，且元素仍可用
	for range st.AllElems() {
		break
	}
	if st.Elem(1).String() != "double" {
		t.Fatal("element access after early break failed")
	}
}

func TestTypeOfGoCacheHit(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()
	type S struct {
		A int32
		B float64
	}
	a, err := TypeOf[S](ctx)
	if err != nil {
		t.Fatal(err)
	}
	b, err := TypeOf[S](ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !a.Equal(b) {
		t.Fatal("cached mapping should return the same LLVM type")
	}
}

func TestTypeOfGoCacheLifetime(t *testing.T) {
	ctx := NewContext()
	type S struct{ A int32 }
	if _, err := TypeOf[S](ctx); err != nil {
		t.Fatal(err)
	}
	if err := ctx.Close(); err != nil {
		t.Fatal(err)
	}
	err := Catch(func() { _, _ = TypeOf[S](ctx) })
	if err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("cached mapping on closed context should panic ErrUseAfterFree, got %v", err)
	}
}
