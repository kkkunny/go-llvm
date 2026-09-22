package llvm

import (
	"strings"
	"testing"
)

func TestDataLayoutQueries(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()
	dl := NewDataLayout("e-p:64:64-i64:64-n8:16:32:64")
	defer dl.Close()

	if got := dl.ByteOrder(); got != LittleEndian {
		t.Fatalf("ByteOrder() = %v", got)
	}
	if got := dl.PointerSize(); got != 8 {
		t.Fatalf("PointerSize() = %d", got)
	}
	if got := dl.PointerSizeForAS(0); got != 8 {
		t.Fatalf("PointerSizeForAS(0) = %d", got)
	}

	i64 := ctx.Int(64)
	if got := dl.SizeOfTypeInBits(i64); got != 64 {
		t.Fatalf("SizeOfTypeInBits(i64) = %d", got)
	}
	if got := dl.StoreSizeOfType(i64); got != 8 {
		t.Fatalf("StoreSizeOfType(i64) = %d", got)
	}
	if got := dl.ABISizeOfType(i64); got != 8 {
		t.Fatalf("ABISizeOfType(i64) = %d", got)
	}
	if got := dl.ABIAlignOfType(i64); got != 8 {
		t.Fatalf("ABIAlignOfType(i64) = %d", got)
	}
	if got := dl.PrefAlignOfType(i64); got != 8 {
		t.Fatalf("PrefAlignOfType(i64) = %d", got)
	}
	if got := dl.CallFrameAlignOfType(i64); got == 0 {
		t.Fatal("CallFrameAlignOfType(i64) should be non-zero")
	}
	if !dl.IntPtrType(ctx).Equal(ctx.Int(64)) {
		t.Fatalf("IntPtrType() = %s", dl.IntPtrType(ctx))
	}
	if !dl.IntPtrTypeForAS(ctx, 0).Equal(ctx.Int(64)) {
		t.Fatalf("IntPtrTypeForAS() = %s", dl.IntPtrTypeForAS(ctx, 0))
	}
	if got := dl.String(); !strings.Contains(got, "p:64:64") || !strings.Contains(got, "i64:64") {
		t.Fatalf("String() = %q", got)
	}
}

func TestDataLayoutStructQueries(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()
	dl := NewDataLayout("e-p:64:64-i64:64-n8:16:32:64")
	defer dl.Close()

	st := ctx.Struct([]AnyType{ctx.Int(8), ctx.Int(64)}, false)
	if got := dl.OffsetOfElement(st, 1); got != 8 {
		t.Fatalf("OffsetOfElement(st, 1) = %d", got)
	}
	if got := dl.ElementAtOffset(st, 8); got != 1 {
		t.Fatalf("ElementAtOffset(st, 8) = %d", got)
	}
	if got := dl.SizeOfTypeInBits(st); got != 128 {
		t.Fatalf("SizeOfTypeInBits(st) = %d", got)
	}
	if got := dl.ABISizeOfType(st); got != 16 {
		t.Fatalf("ABISizeOfType(st) = %d", got)
	}
}

func TestDataLayoutClose(t *testing.T) {
	dl := NewDataLayout("e-p:64:64")
	_ = dl.Close()
	if err := dl.Close(); err == nil || err.(*Error).Reason != ErrClosed {
		t.Fatalf("double close should return ErrClosed, got %v", err)
	}
	if err := Catch(func() { dl.PointerSize() }); err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("use after close should panic ErrUseAfterFree, got %v", err)
	}
}
