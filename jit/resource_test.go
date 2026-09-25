package jit

import (
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestResourceTrackerUnload(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	ctx, m := retModule(t, "rt", 42)
	rt := j.NewResourceTracker()
	if err := rt.AddIRModule(m); err != nil {
		t.Fatal(err)
	}
	if m.Lifetime().Alive() || ctx.Alive() {
		t.Fatal("module and context should be invalidated after ownership transfer")
	}

	if p, err := j.Lookup("answer"); err != nil || p == nil {
		t.Fatalf("Lookup before Remove: %v", err)
	}
	if err := rt.Remove(); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if _, err := j.Lookup("answer"); err == nil || err.(*llvm.Error).Reason != llvm.ErrNotFound {
		t.Fatalf("symbol should be gone after Remove, got %v", err)
	}
}

func TestResourceTrackerTransferAndClose(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	rt1 := j.NewResourceTracker()
	rt2 := j.NewResourceTracker()
	rt1.TransferTo(rt2)
	if err := rt1.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if err := rt1.Close(); err == nil {
		t.Fatalf("double close should error")
	}
	if err := rt2.Close(); err != nil {
		t.Fatalf("Close rt2: %v", err)
	}
}

func TestClearSymbols(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	_, m := retModule(t, "clr", 1)
	if err := j.AddIRModule(m); err != nil {
		t.Fatal(err)
	}
	if _, err := j.Lookup("answer"); err != nil {
		t.Fatalf("Lookup before clear: %v", err)
	}
	if err := j.ClearSymbols(); err != nil {
		t.Fatalf("ClearSymbols: %v", err)
	}
	if _, err := j.Lookup("answer"); err == nil || err.(*llvm.Error).Reason != llvm.ErrNotFound {
		t.Fatalf("symbol should be gone after ClearSymbols, got %v", err)
	}

	// 清空后 dylib 仍可复用
	_, m2 := retModule(t, "clr2", 2)
	if err := j.AddIRModule(m2); err != nil {
		t.Fatalf("re-add after clear: %v", err)
	}
	if _, err := j.Lookup("answer"); err != nil {
		t.Fatalf("Lookup after re-add: %v", err)
	}
}
