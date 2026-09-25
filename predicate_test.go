package llvm

import "testing"

// TestIntPredValues 与 llvm-c/Core.h 的 LLVMIntPredicate 枚举对齐（从 32 起编号）。
func TestIntPredValues(t *testing.T) {
	want := []struct {
		name string
		got  IntPred
		val  int
	}{
		{"EQ", IntEQ, 32},
		{"NE", IntNE, 33},
		{"UGT", IntUGT, 34},
		{"UGE", IntUGE, 35},
		{"ULT", IntULT, 36},
		{"ULE", IntULE, 37},
		{"SGT", IntSGT, 38},
		{"SGE", IntSGE, 39},
		{"SLT", IntSLT, 40},
		{"SLE", IntSLE, 41},
	}
	for _, c := range want {
		if int(c.got) != c.val {
			t.Errorf("Int%s = %d, want %d", c.name, int(c.got), c.val)
		}
	}
}

// TestFloatPredValues 与 llvm-c/Core.h 的 LLVMRealPredicate 枚举对齐（0..15）。
func TestFloatPredValues(t *testing.T) {
	want := []struct {
		name string
		got  FloatPred
		val  int
	}{
		{"False", FloatFalse, 0},
		{"OEQ", FloatOEQ, 1},
		{"OGT", FloatOGT, 2},
		{"OGE", FloatOGE, 3},
		{"OLT", FloatOLT, 4},
		{"OLE", FloatOLE, 5},
		{"ONE", FloatONE, 6},
		{"ORD", FloatORD, 7},
		{"UNO", FloatUNO, 8},
		{"UEQ", FloatUEQ, 9},
		{"UGT", FloatUGT, 10},
		{"UGE", FloatUGE, 11},
		{"ULT", FloatULT, 12},
		{"ULE", FloatULE, 13},
		{"UNE", FloatUNE, 14},
		{"True", FloatTrue, 15},
	}
	for _, c := range want {
		if int(c.got) != c.val {
			t.Errorf("Float%s = %d, want %d", c.name, int(c.got), c.val)
		}
	}
}
