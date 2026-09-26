package llvm

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kkkunny/go-llvm/internal/binding"
)

func TestMemoryBufferRoundTrip(t *testing.T) {
	buf := NewMemoryBuffer([]byte("hello"), "test")

	if !buf.Alive() {
		t.Fatal("new buffer should be alive")
	}
	if got := string(buf.Bytes()); got != "hello" {
		t.Fatalf("Bytes() = %q", got)
	}
	if buf.Len() != 5 {
		t.Fatalf("Len() = %d", buf.Len())
	}
	if err := buf.Close(); err != nil {
		t.Fatalf("first close should succeed, got %v", err)
	}
	if err := buf.Close(); err == nil || err.(*Error).Reason != ErrClosed {
		t.Fatalf("double close should return ErrClosed, got %v", err)
	}
}

func TestReadFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "in.txt")
	if err := os.WriteFile(path, []byte("abc"), 0o644); err != nil {
		t.Fatal(err)
	}
	buf, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	defer buf.Close()
	if got := string(buf.Bytes()); got != "abc" {
		t.Fatalf("Bytes() = %q", got)
	}

	if _, err := ReadFile(filepath.Join(t.TempDir(), "missing")); err == nil || err.(*Error).Reason != ErrIO {
		t.Fatalf("missing file should return ErrIO, got %v", err)
	}
}

func TestMemoryBufferRefAndFinalize(t *testing.T) {
	buf := NewMemoryBuffer([]byte("xyz"), "test")
	if buf.Ref().IsNil() {
		t.Fatal("Ref() 不应为空句柄")
	}
	buf.finalize() // 直接触发 GC 兜底路径
	if !buf.closed || buf.Alive() {
		t.Fatal("finalize 后应标记 closed 且不再 Alive")
	}
	buf.finalize() // 二次调用提前返回，不得重复释放
	if err := Catch(func() { buf.Bytes() }); err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("finalize 后使用应 panic ErrUseAfterFree, got %v", err)
	}
	if err := buf.Close(); err == nil || err.(*Error).Reason != ErrClosed {
		t.Fatalf("finalize 后 Close 应为 ErrClosed, got %v", err)
	}
}

func TestMemoryBufferDisown(t *testing.T) {
	buf := NewMemoryBuffer([]byte("xyz"), "test")
	ref := buf.Ref()
	buf.Disown()
	if buf.Alive() {
		t.Fatal("Disown 后句柄应失效")
	}
	buf.Disown() // 幂等
	if err := buf.Close(); err == nil || err.(*Error).Reason != ErrClosed {
		t.Fatalf("Disown 后 Close 应为 ErrClosed, got %v", err)
	}
	// 测试中扮演接管方，负责释放底层缓冲
	binding.LLVMDisposeMemoryBuffer(ref)
}

func TestMemoryBufferAfterClose(t *testing.T) {
	buf := NewMemoryBuffer([]byte("x"), "test")
	_ = buf.Close()
	if buf.Alive() {
		t.Fatal("closed buffer should not be alive")
	}
	if err := Catch(func() { buf.Bytes() }); err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("Bytes after close should panic ErrUseAfterFree, got %v", err)
	}
}
