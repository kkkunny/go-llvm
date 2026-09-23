package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestFunctionPersonality(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "eh")
	defer m.Close()

	i32 := ctx.Int(32)
	pers := m.NewFunction("pers", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))

	fn.SetPersonality(pers)
	got, ok := fn.Personality()
	if !ok {
		t.Fatal("personality should be set")
	}
	if got.Name() != "pers" {
		t.Fatalf("personality = %s, want pers", got.Name())
	}
	if !strings.Contains(m.String(), "personality ptr @pers") {
		t.Fatalf("module should contain personality:\n%s", m.String())
	}

	if err := llvm.Catch(func() { fn.SetPersonality(llvm.Value[llvm.FnT]{}) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil personality should panic ErrInvalidArg, got %v", err)
	}
}
