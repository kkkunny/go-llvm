package target

import (
	"testing"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/errs"
)

func TestInitAndHostQueries(t *testing.T) {
	InitNative()

	native, err := NativeTarget()
	if err != nil {
		t.Fatal(err)
	}
	if !native.HasTargetMachine() {
		t.Fatalf("native target %q should have a target machine", native.Name())
	}
	if !native.HasAsmBackend() {
		t.Fatalf("native target %q should have an asm backend", native.Name())
	}
	if native.Name() == "" || native.Description() == "" {
		t.Fatal("native target name/description should be non-empty")
	}
	if triple := DefaultTriple(); triple == "" || NormalizeTriple(triple) == "" {
		t.Fatal("host triple should be non-empty")
	}
	if HostCPUName() == "" {
		t.Fatal("host cpu name should be non-empty")
	}
}

func TestInitAllAndLookup(t *testing.T) {
	InitAll()
	Init(X86)

	if got, ok := FromName("x86-64"); !ok || got.Name() != "x86-64" {
		t.Fatalf("FromName(x86-64) = %q, %v", got.Name(), ok)
	}
	if _, ok := FromName("definitely-not-a-target"); ok {
		t.Fatal("unknown target name should not resolve")
	}
	if _, err := FromTriple("x86_64-unknown-linux-gnu"); err != nil {
		t.Fatalf("FromTriple(x86_64-unknown-linux-gnu) = %v", err)
	}
	if _, err := FromTriple("not-a-triple"); err == nil || err.(*llvm.Error).Reason != llvm.ErrNotFound {
		t.Fatalf("invalid triple should return ErrNotFound, got %v", err)
	}
}

func TestInitUnknownArch(t *testing.T) {
	if err := errs.Catch(func() { Init(Arch(255)) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("unknown arch should panic ErrInvalidArg, got %v", err)
	}
}

// TestTargetRefAndHasJIT 验证 Ref 暴露的底层句柄指向全局唯一的目标描述
func TestTargetRefAndHasJIT(t *testing.T) {
	InitNative()

	native, err := NativeTarget()
	if err != nil {
		t.Fatal(err)
	}
	if native.Ref().IsNil() {
		t.Fatal("Ref should return a non-nil handle")
	}
	// 同一三元组重复查询应得到同一目标描述句柄
	again, err := FromTriple(DefaultTriple())
	if err != nil {
		t.Fatal(err)
	}
	if native.Ref() != again.Ref() {
		t.Fatalf("Ref handles should be identical for %q", DefaultTriple())
	}
	// 宿主目标是 LLJIT 的基础，必须支持 JIT
	if !native.HasJIT() {
		t.Fatalf("native target %q should support JIT", native.Name())
	}
}

// TestZeroTargetFloor 零值 Target 不得进入 LLVM（FromName 失败后的零值常见）。
func TestZeroTargetFloor(t *testing.T) {
	var zero Target
	if err := errs.Catch(func() { _ = zero.Name() }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("zero Target.Name should panic ErrInvalidArg, got %v", err)
	}
	if err := errs.Catch(func() {
		_, _ = NewTargetMachine(zero, DefaultTriple(), "", "", OptNone, RelocDefault, CodeModelDefault)
	}); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("NewTargetMachine(zero) should panic ErrInvalidArg, got %v", err)
	}
}
