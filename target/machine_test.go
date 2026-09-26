package target

import (
	"bytes"
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

// TestTargetMachineCPUFeatures 验证建机时传入的 CPU/特性串可原样回读
func TestTargetMachineCPUFeatures(t *testing.T) {
	ctx, m, tm := newMachineModule(t)
	defer ctx.Close()
	defer m.Close()
	defer tm.Close()

	// newMachineModule 以宿主 CPU/特性建机，回读应一致
	if got, want := tm.CPU(), HostCPUName(); got != want {
		t.Fatalf("CPU() = %q, want %q", got, want)
	}
	if got, want := tm.Features(), HostCPUFeatures(); got != want {
		t.Fatalf("Features() = %q, want %q", got, want)
	}
}

// TestTargetMachineSetAsmVerbosity 验证汇编详细模式确实改变汇编产物，且可关闭恢复
func TestTargetMachineSetAsmVerbosity(t *testing.T) {
	ctx, m, tm := newMachineModule(t)
	defer ctx.Close()
	defer m.Close()
	defer tm.Close()

	plain, err := tm.Emit(m, AsmFile)
	if err != nil {
		t.Fatal(err)
	}
	defer plain.Close()

	// 打开详细汇编：产物应带有块注释等额外信息
	tm.SetAsmVerbosity(true)
	verbose, err := tm.Emit(m, AsmFile)
	if err != nil {
		t.Fatal(err)
	}
	defer verbose.Close()
	if bytes.Equal(plain.Bytes(), verbose.Bytes()) {
		t.Fatal("verbose asm output should differ from the default output")
	}

	// 关闭后应恢复默认产物
	tm.SetAsmVerbosity(false)
	restored, err := tm.Emit(m, AsmFile)
	if err != nil {
		t.Fatal(err)
	}
	defer restored.Close()
	if !bytes.Equal(plain.Bytes(), restored.Bytes()) {
		t.Fatal("asm output should return to the default after disabling verbosity")
	}
}

// TestTargetMachineEmitToFileError 输出路径不可写时返回 ErrCodeGen（而非崩溃或静默成功）
func TestTargetMachineEmitToFileError(t *testing.T) {
	ctx, m, tm := newMachineModule(t)
	defer ctx.Close()
	defer m.Close()
	defer tm.Close()

	// 父目录不存在：打开输出文件失败
	bad := filepath.Join(t.TempDir(), "no-such-dir", "m.s")
	if err := tm.EmitToFile(m, bad, AsmFile); err == nil || err.(*llvm.Error).Reason != llvm.ErrCodeGen {
		t.Fatalf("emit to unwritable path should return ErrCodeGen, got %v", err)
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
	// 释放后 CPU/Features/SetAsmVerbosity 均须在触碰句柄前拦截
	for name, fn := range map[string]func(){
		"CPU":             func() { tm.CPU() },
		"Features":        func() { tm.Features() },
		"SetAsmVerbosity": func() { tm.SetAsmVerbosity(true) },
	} {
		if err := llvm.Catch(fn); err == nil || err.Reason != llvm.ErrUseAfterFree {
			t.Fatalf("%s after close should panic ErrUseAfterFree, got %v", name, err)
		}
	}
}
