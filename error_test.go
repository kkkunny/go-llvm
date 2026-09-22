package llvm

import (
	"strings"
	"testing"
)

func TestCatch(t *testing.T) {
	t.Run("收敛 *Error", func(t *testing.T) {
		err := Catch(func() {
			errPanic(ErrInvalidArg, "llvm.Test", "bad arg %d", 7)
		})
		if err == nil {
			t.Fatal("want error, got nil")
		}
		if err.Reason != ErrInvalidArg || err.Op != "llvm.Test" || err.Msg != "bad arg 7" {
			t.Fatalf("unexpected error: %+v", err)
		}
	})

	t.Run("无 panic 返回 nil", func(t *testing.T) {
		if err := Catch(func() {}); err != nil {
			t.Fatalf("want nil, got %v", err)
		}
	})

	t.Run("非 *Error 原样重抛", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("want re-panic, got none")
			}
			if s, ok := r.(string); !ok || s != "boom" {
				t.Fatalf("unexpected panic value: %v", r)
			}
		}()
		Catch(func() { panic("boom") })
	})
}

func TestMust(t *testing.T) {
	if got := Must(42, nil); got != 42 {
		t.Fatalf("want 42, got %d", got)
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("want panic, got none")
		}
		if err, ok := r.(*Error); !ok || err.Reason != ErrNotFound {
			t.Fatalf("unexpected panic value: %v", r)
		}
	}()
	Must(0, &Error{Reason: ErrNotFound, Op: "llvm.Test", Msg: "missing"})
}

func TestErrorMessage(t *testing.T) {
	err := &Error{Reason: ErrTypeMismatch, Op: "llvm.Builder.Add", Msg: "operand kinds differ"}
	want := "llvm.Builder.Add: operand kinds differ"
	if got := err.Error(); got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
	if !strings.Contains(err.Error(), "operand kinds differ") {
		t.Fatal("message lost")
	}
}
