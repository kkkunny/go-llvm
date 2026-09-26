package llvm

import "testing"

// TestLinkageValues 与 llvm-c/Core.h 的 LLVMLinkage 枚举逐项对齐。
// 公共包只导出稳定的链接类型；已废弃的 DLLImport/Ghost 等不导出。
func TestLinkageValues(t *testing.T) {
	want := []struct {
		name string
		got  Linkage
		val  int
	}{
		{"External", LinkageExternal, 0},
		{"AvailableExternally", LinkageAvailableExternally, 1},
		{"LinkOnceAny", LinkageLinkOnceAny, 2},
		{"LinkOnceODR", LinkageLinkOnceODR, 3},
		{"LinkOnceODRAutoHide", LinkageLinkOnceODRAutoHide, 4},
		{"WeakAny", LinkageWeakAny, 5},
		{"WeakODR", LinkageWeakODR, 6},
		{"Appending", LinkageAppending, 7},
		{"Internal", LinkageInternal, 8},
		{"Private", LinkagePrivate, 9},
		{"ExternalWeak", LinkageExternalWeak, 12},
		{"Common", LinkageCommon, 14},
	}
	seen := make(map[int]string, len(want))
	for _, c := range want {
		if int(c.got) != c.val {
			t.Errorf("Linkage%s = %d, want %d", c.name, int(c.got), c.val)
		}
		if prev, ok := seen[int(c.got)]; ok {
			t.Errorf("Linkage%s 与 Linkage%s 取值重复 (%d)", c.name, prev, int(c.got))
		}
		seen[int(c.got)] = c.name
	}
}

func TestVisibilityValues(t *testing.T) {
	if VisibilityDefault != 0 || VisibilityHidden != 1 || VisibilityProtected != 2 {
		t.Fatalf("Visibility = %d/%d/%d, want 0/1/2",
			VisibilityDefault, VisibilityHidden, VisibilityProtected)
	}
}

func TestDLLStorageClassValues(t *testing.T) {
	if DLLStorageDefault != 0 || DLLStorageImport != 1 || DLLStorageExport != 2 {
		t.Fatalf("DLLStorageClass = %d/%d/%d, want 0/1/2",
			DLLStorageDefault, DLLStorageImport, DLLStorageExport)
	}
}

func TestUnnamedAddrValues(t *testing.T) {
	if UnnamedAddrNone != 0 || UnnamedAddrLocal != 1 || UnnamedAddrGlobal != 2 {
		t.Fatalf("UnnamedAddr = %d/%d/%d, want 0/1/2",
			UnnamedAddrNone, UnnamedAddrLocal, UnnamedAddrGlobal)
	}
}

func TestThreadLocalModeValues(t *testing.T) {
	want := []struct {
		name string
		got  ThreadLocalMode
		val  int
	}{
		{"None", ThreadLocalNone, 0},
		{"General", ThreadLocalGeneral, 1},
		{"LocalDynamic", ThreadLocalLocalDynamic, 2},
		{"InitialExec", ThreadLocalInitialExec, 3},
		{"LocalExec", ThreadLocalLocalExec, 4},
	}
	for _, c := range want {
		if int(c.got) != c.val {
			t.Errorf("ThreadLocal%s = %d, want %d", c.name, int(c.got), c.val)
		}
	}
}
