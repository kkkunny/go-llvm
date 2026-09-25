package llvm

import "testing"

func TestIntConstFolding(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()
	i32 := ctx.Int(32)

	a := i32.Const(10)
	b := i32.Const(3)
	if got := a.Add(b).UnsignedValue(); got != 13 {
		t.Fatalf("Add = %d", got)
	}
	if got := a.Sub(b).UnsignedValue(); got != 7 {
		t.Fatalf("Sub = %d", got)
	}
	if got := a.Xor(b).UnsignedValue(); got != 9 {
		t.Fatalf("Xor = %d", got)
	}
	if got := a.Neg().SignedValue(); got != -10 {
		t.Fatalf("Neg = %d", got)
	}
	if got, want := a.Not().UnsignedValue(), uint64(0xffffffff^10); got != want {
		t.Fatalf("Not = %d, want %d", got, want)
	}
	if got := a.NSWAdd(b).UnsignedValue(); got != 13 {
		t.Fatalf("NSWAdd = %d", got)
	}
	if got := a.NUWAdd(b).UnsignedValue(); got != 13 {
		t.Fatalf("NUWAdd = %d", got)
	}
	if got := a.NSWSub(b).UnsignedValue(); got != 7 {
		t.Fatalf("NSWSub = %d", got)
	}
	if got := a.NSWNeg().SignedValue(); got != -10 {
		t.Fatalf("NSWNeg = %d", got)
	}
	if got := a.NUWNeg().UnsignedValue(); got != uint64(0xffffffff-10+1) {
		t.Fatalf("NUWNeg = %d", got)
	}

	ones := i32.AllOnes()
	if ones.UnsignedValue() != 0xffffffff {
		t.Fatalf("AllOnes = %x", ones.UnsignedValue())
	}

	// cast 折叠：i32 -> i16
	i16 := ctx.Int(16)
	narrow := i32.Const(0x12345678).Cast(i16)
	if narrow.UnsignedValue() != 0x5678 {
		t.Fatalf("Cast = %x", narrow.UnsignedValue())
	}
	// 扩展方向同样可用
	wide := narrow.Cast(i32)
	if wide.UnsignedValue() != 0x5678 {
		t.Fatalf("Cast wide = %x", wide.UnsignedValue())
	}
}

func TestConstFoldMismatchPanic(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	i32 := ctx.Int(32)
	i64 := ctx.Int(64)

	err := Catch(func() {
		i32.Const(1).Add(i64.Const(1))
	})
	if err == nil || err.Reason != ErrTypeMismatch {
		t.Fatalf("want ErrTypeMismatch, got %v", err)
	}
}
