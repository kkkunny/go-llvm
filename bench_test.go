package llvm

import "testing"

func BenchmarkConstArray(b *testing.B) {
	ctx := NewContext()
	defer ctx.Close()
	i32 := ctx.Int(32)
	elems := make([]AnyValue, 8)
	for i := range elems {
		elems[i] = ctx.ConstInt(i32, uint64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.ConstArray(i32, elems...)
	}
}

func BenchmarkTypeOfGo(b *testing.B) {
	ctx := NewContext()
	defer ctx.Close()
	type S struct {
		A int32
		B float64
		C *int64
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := TypeOf[S](ctx); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFnSignatureOfGo(b *testing.B) {
	ctx := NewContext()
	defer ctx.Close()
	type F func(int32, float64, *int64) int64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, err := FnSignatureOf[F](ctx); err != nil {
			b.Fatal(err)
		}
	}
}
