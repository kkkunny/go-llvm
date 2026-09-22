package ir

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

// stripModuleID 去掉 `; ModuleID = ...` 行：解析往返时该行会变为缓冲名或被 bitcode 丢弃
func stripModuleID(s string) string {
	var b strings.Builder
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(line, "; ModuleID =") {
			continue
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}

func buildRetModule(t *testing.T, name string) (*llvm.Context, *Module) {
	t.Helper()
	ctx := llvm.NewContext()
	m := NewModule(ctx, name)
	i32 := ctx.Int(32)
	fn := m.NewFunction("main", ctx.Fn(i32, nil, false))
	b := NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(ctx.ConstInt(i32, 0, false))
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	if err := m.Verify(); err != nil {
		t.Fatal(err)
	}
	return ctx, m
}

func TestParseIRRoundTrip(t *testing.T) {
	ctx, m := buildRetModule(t, "roundtrip")
	defer ctx.Close()
	defer m.Close()

	text := m.String()
	buf := llvm.NewMemoryBuffer([]byte(text), "roundtrip.ll")
	defer buf.Close()
	parsed, err := ParseIR(ctx, buf)
	if err != nil {
		t.Fatal(err)
	}
	defer parsed.Close()

	// 解析后 ModuleID 变为缓冲名，其余 IR 必须逐字一致
	if stripModuleID(parsed.String()) != stripModuleID(text) {
		t.Fatalf("round trip mismatch:\n%s\n---\n%s", parsed.String(), text)
	}
	if !buf.Alive() {
		t.Fatal("ParseIR must not consume the memory buffer")
	}
}

func TestParseIRInvalid(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	buf := llvm.NewMemoryBuffer([]byte("this is not IR"), "bad.ll")
	defer buf.Close()

	_, err := ParseIR(ctx, buf)
	if err == nil || err.(*llvm.Error).Reason != llvm.ErrParse {
		t.Fatalf("invalid IR should return ErrParse, got %v", err)
	}
}

func TestBitcodeRoundTrip(t *testing.T) {
	ctx, m := buildRetModule(t, "bitcode")
	defer ctx.Close()
	defer m.Close()

	path := filepath.Join(t.TempDir(), "m.bc")
	if err := m.WriteBitcode(path); err != nil {
		t.Fatal(err)
	}
	buf, err := llvm.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	defer buf.Close()
	parsed, err := ParseBitcode(ctx, buf)
	if err != nil {
		t.Fatal(err)
	}
	defer parsed.Close()
	if stripModuleID(parsed.String()) != stripModuleID(m.String()) {
		t.Fatalf("bitcode round trip mismatch:\n%s\n---\n%s", parsed.String(), m.String())
	}

	bc := m.Bitcode()
	defer bc.Close()
	parsed2, err := ParseBitcode(ctx, bc)
	if err != nil {
		t.Fatal(err)
	}
	defer parsed2.Close()
	if stripModuleID(parsed2.String()) != stripModuleID(m.String()) {
		t.Fatalf("memory bitcode round trip mismatch")
	}

	bad := llvm.NewMemoryBuffer([]byte("junk"), "bad.bc")
	defer bad.Close()
	if _, err := ParseBitcode(ctx, bad); err == nil || err.(*llvm.Error).Reason != llvm.ErrParse {
		t.Fatalf("invalid bitcode should return ErrParse, got %v", err)
	}
}

func TestModuleWriteToFile(t *testing.T) {
	ctx, m := buildRetModule(t, "writefile")
	defer ctx.Close()
	defer m.Close()

	path := filepath.Join(t.TempDir(), "m.ll")
	if err := m.WriteToFile(path); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "define i32 @main()") {
		t.Fatalf("written IR:\n%s", data)
	}

	if err := m.WriteToFile(filepath.Join(t.TempDir(), "no/such/dir/m.ll")); err == nil || err.(*llvm.Error).Reason != llvm.ErrIO {
		t.Fatalf("bad path should return ErrIO, got %v", err)
	}
}

func TestModuleDataLayout(t *testing.T) {
	ctx, m := buildRetModule(t, "dl")
	defer ctx.Close()
	defer m.Close()

	dl := m.DataLayout()
	defer dl.Close()
	if dl.PointerSize() == 0 {
		t.Fatal("module data layout pointer size should be non-zero")
	}

	g := m.NewGlobal("g", ctx.Int(32))
	g.SetAlign(4)
	if got := dl.PrefAlignOfGlobal(g); got != 4 {
		t.Fatalf("PrefAlignOfGlobal(g) = %d", got)
	}
}
