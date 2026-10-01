package llvm

import (
	"strings"
	"testing"
)

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
