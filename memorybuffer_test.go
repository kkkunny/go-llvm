package llvm

import (
	"os"
	"path/filepath"
	"testing"
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
