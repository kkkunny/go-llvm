package target

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
)

func newMachineModule(t *testing.T) (*llvm.Context, *ir.Module, *TargetMachine) {
	t.Helper()
	if err := InitNative(); err != nil {
		t.Fatal(err)
	}
	native, err := NativeTarget()
	if err != nil {
		t.Fatal(err)
	}
	tm, err := NewTargetMachine(native, DefaultTriple(), HostCPUName(), HostCPUFeatures(), OptDefault, RelocPIC, CodeModelDefault)
	if err != nil {
		t.Fatal(err)
	}

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "codegen")
	i32 := ctx.Int(32)
	fn := m.NewFunction("main", ctx.Fn(i32, nil, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(ctx.ConstInt(i32, 0))
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	tm.ApplyTo(m)
	if err := m.Verify(); err != nil {
		t.Fatal(err)
	}
	return ctx, m, tm
}

func TestTargetMachineApplyTo(t *testing.T) {
	ctx, m, tm := newMachineModule(t)
	defer ctx.Close()
	defer m.Close()
	defer tm.Close()

	if tm.Triple() != m.TargetTriple() {
		t.Fatalf("ApplyTo triple: module %q, machine %q", m.TargetTriple(), tm.Triple())
	}
	if tm.Target().Name() == "" {
		t.Fatal("Target() should be non-nil")
	}
	mdl := m.DataLayout()
	defer mdl.Close()
	tmdl := tm.DataLayout()
	defer tmdl.Close()
	if mdl.String() != tmdl.String() {
		t.Fatalf("ApplyTo data layout: module %q, machine %q", mdl.String(), tmdl.String())
	}
}

func TestTargetMachineEmit(t *testing.T) {
	ctx, m, tm := newMachineModule(t)
	defer ctx.Close()
	defer m.Close()
	defer tm.Close()

	asm, err := tm.Emit(m, AsmFile)
	if err != nil {
		t.Fatal(err)
	}
	defer asm.Close()
	if !strings.Contains(string(asm.Bytes()), "main") {
		t.Fatalf("asm output missing main:\n%s", asm.Bytes())
	}

	obj, err := tm.Emit(m, ObjectFile)
	if err != nil {
		t.Fatal(err)
	}
	defer obj.Close()
	if obj.Len() == 0 {
		t.Fatal("object output should be non-empty")
	}

	path := filepath.Join(t.TempDir(), "m.s")
	if err := tm.EmitToFile(m, path, AsmFile); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "main") {
		t.Fatalf("asm file missing main:\n%s", data)
	}
}

func TestTargetMachineChecks(t *testing.T) {
	ctx, m, tm := newMachineModule(t)
	defer ctx.Close()
	defer m.Close()

	closed := ir.NewModule(ctx, "closed")
	if err := closed.Close(); err != nil {
		t.Fatal(err)
	}
	if err := llvm.Catch(func() { tm.EmitToFile(closed, "x.s", AsmFile) }); err == nil || err.Reason != llvm.ErrUseAfterFree {
		t.Fatalf("emit closed module should panic ErrUseAfterFree, got %v", err)
	}

	if err := tm.Close(); err != nil {
		t.Fatal(err)
	}
	if err := tm.Close(); err == nil || err.(*llvm.Error).Reason != llvm.ErrClosed {
		t.Fatalf("double close should return ErrClosed, got %v", err)
	}
	if err := llvm.Catch(func() { tm.Triple() }); err == nil || err.Reason != llvm.ErrUseAfterFree {
		t.Fatalf("use after close should panic ErrUseAfterFree, got %v", err)
	}
}
