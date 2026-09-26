package llvm

import "testing"

// TestAtomicOrderingValues 与 llvm-c/Core.h 的 LLVMAtomicOrdering 枚举对齐；
// 注意 Acquire 从 4 起（3 为 C++ 侧保留值）。
func TestAtomicOrderingValues(t *testing.T) {
	want := []struct {
		name string
		got  AtomicOrdering
		val  int
	}{
		{"NotAtomic", AtomicNotAtomic, 0},
		{"Unordered", AtomicUnordered, 1},
		{"Monotonic", AtomicMonotonic, 2},
		{"Acquire", AtomicAcquire, 4},
		{"Release", AtomicRelease, 5},
		{"AcquireRelease", AtomicAcquireRelease, 6},
		{"SequentiallyConsistent", AtomicSequentiallyConsistent, 7},
	}
	for _, c := range want {
		if int(c.got) != c.val {
			t.Errorf("Atomic%s = %d, want %d", c.name, int(c.got), c.val)
		}
	}
}

// TestRMWOpValues 与 llvm-c/Core.h 的 LLVMAtomicRMWBinOp 枚举对齐（0..20）。
func TestRMWOpValues(t *testing.T) {
	want := []struct {
		name string
		got  RMWOp
		val  int
	}{
		{"Xchg", RMWXchg, 0},
		{"Add", RMWAdd, 1},
		{"Sub", RMWSub, 2},
		{"And", RMWAnd, 3},
		{"Nand", RMWNand, 4},
		{"Or", RMWOr, 5},
		{"Xor", RMWXor, 6},
		{"Max", RMWMax, 7},
		{"Min", RMWMin, 8},
		{"UMax", RMWUMax, 9},
		{"UMin", RMWUMin, 10},
		{"FAdd", RMWFAdd, 11},
		{"FSub", RMWFSub, 12},
		{"FMax", RMWFMax, 13},
		{"FMin", RMWFMin, 14},
		{"UIncWrap", RMWUIncWrap, 15},
		{"UDecWrap", RMWUDecWrap, 16},
		{"USubCond", RMWUSubCond, 17},
		{"USubSat", RMWUSubSat, 18},
		{"FMaximum", RMWFMaximum, 19},
		{"FMinimum", RMWFMinimum, 20},
	}
	seen := make(map[int]string, len(want))
	for _, c := range want {
		if int(c.got) != c.val {
			t.Errorf("RMW%s = %d, want %d", c.name, int(c.got), c.val)
		}
		if prev, ok := seen[int(c.got)]; ok {
			t.Errorf("RMW%s 与 RMW%s 取值重复 (%d)", c.name, prev, int(c.got))
		}
		seen[int(c.got)] = c.name
	}
}
