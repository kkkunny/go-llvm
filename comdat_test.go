package llvm

import "testing"

func TestComdatSelectionKindValues(t *testing.T) {
	// 与 llvm-c/Comdat.h 的 LLVMComdatSelectionKind 枚举一一对应
	want := []struct {
		name string
		got  ComdatSelectionKind
		val  int
	}{
		{"Any", ComdatAny, 0},
		{"ExactMatch", ComdatExactMatch, 1},
		{"Largest", ComdatLargest, 2},
		{"NoDeduplicate", ComdatNoDeduplicate, 3},
		{"SameSize", ComdatSameSize, 4},
	}
	seen := make(map[int]string, len(want))
	for _, c := range want {
		if int(c.got) != c.val {
			t.Errorf("Comdat%s = %d, want %d", c.name, int(c.got), c.val)
		}
		if prev, ok := seen[int(c.got)]; ok {
			t.Errorf("Comdat%s 与 Comdat%s 取值重复 (%d)", c.name, prev, int(c.got))
		}
		seen[int(c.got)] = c.name
	}
}
