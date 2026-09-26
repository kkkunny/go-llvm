package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestParseIRString(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m, err := ParseIRString(ctx, `
		define i32 @f() {
		entry:
			ret i32 42
		}`)
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()
	fn, ok := m.GetFunction("f")
	if !ok || fn.OnlyDecl() {
		t.Fatalf("f should be defined: ok=%v onlyDecl=%v", ok, fn.OnlyDecl())
	}
}

func TestPtrAddAndBuilderAt(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "ptradd")
	defer m.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	fn := m.NewFunction("f", ctx.Fn(ptr, []llvm.AnyType{ptr, i32}, false))
	b := NewBuilderAt(fn.NewBlock("entry"))
	defer b.Close()

	p := fn.ParamAs[llvm.PtrT](0)
	off := fn.ParamAs[llvm.IntT](1)
	b.Ret(b.PtrAdd(p, off, "p2"))

	if err := m.Verify(); err != nil {
		t.Fatal(err)
	}
	if got := m.String(); !strings.Contains(got, "getelementptr i8, ptr %0, i32 %1") {
		t.Fatalf("module output:\n%s", got)
	}
}

func TestNewBuilderAtChecks(t *testing.T) {
	err := llvm.Catch(func() { NewBuilderAt(Block{}) })
	if err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil block should panic ErrInvalidArg, got %v", err)
	}
}
