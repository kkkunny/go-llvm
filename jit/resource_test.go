package jit

import (
	"testing"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/target"
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

// TestResourceTrackerAddObjectFile 验证目标文件缓冲可纳入跟踪器并执行，
// Remove 时随跟踪器一起卸载
func TestResourceTrackerAddObjectFile(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	ctx, m := retModule(t, "rt_obj", 7)
	defer ctx.Close()
	defer m.Close()

	tm, err := target.NewTargetMachine(mustNativeTarget(t), target.DefaultTriple(), target.HostCPUName(), target.HostCPUFeatures(), target.OptNone, target.RelocPIC, target.CodeModelDefault)
	if err != nil {
		t.Fatal(err)
	}
	defer tm.Close()
	tm.ApplyTo(m)
	obj, err := tm.Emit(m, target.ObjectFile)
	if err != nil {
		t.Fatal(err)
	}

	rt := j.NewResourceTracker()
	defer func() { _ = rt.Close() }()
	if err := rt.AddObjectFile(obj); err != nil {
		t.Fatalf("AddObjectFile: %v", err)
	}
	if obj.Alive() {
		t.Fatal("object buffer should be invalidated after ownership transfer")
	}

	answer, err := j.Func[func() int32]("answer")
	if err != nil {
		t.Fatal(err)
	}
	if got := answer(); got != 7 {
		t.Fatalf("answer() = %d, want 7", got)
	}

	// Remove 卸载跟踪器中的符号定义
	if err := rt.Remove(); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if _, err := j.Lookup("answer"); err == nil || err.(*llvm.Error).Reason != llvm.ErrNotFound {
		t.Fatalf("symbol should be gone after tracker Remove, got %v", err)
	}
}

// TestResourceTrackerAddObjectFileChecks 非法/损坏缓冲的地板检查与错误返回
func TestResourceTrackerAddObjectFileChecks(t *testing.T) {
	j := newJIT(t)
	defer j.Close()

	rt := j.NewResourceTracker()
	defer func() { _ = rt.Close() }()

	if err := llvm.Catch(func() { _ = rt.AddObjectFile(nil) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil buffer should panic ErrInvalidArg, got %v", err)
	}

	closed := llvm.NewMemoryBuffer([]byte("obj"), "closed.o")
	if err := closed.Close(); err != nil {
		t.Fatal(err)
	}
	if err := llvm.Catch(func() { _ = rt.AddObjectFile(closed) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("closed buffer should panic ErrInvalidArg, got %v", err)
	}

	// 内容不是目标文件：JIT 解析失败应返回 ErrJIT 而非崩溃
	bad := llvm.NewMemoryBuffer([]byte("definitely not an object file"), "bad.o")
	var addErr error
	if panicErr := llvm.Catch(func() { addErr = rt.AddObjectFile(bad) }); panicErr != nil {
		t.Fatalf("invalid object should return an error, got panic %v", panicErr)
	}
	if addErr == nil || addErr.(*llvm.Error).Reason != llvm.ErrJIT {
		t.Fatalf("invalid object should return ErrJIT, got %v", addErr)
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
