package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestNewFuncGoSignature(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "gofunc")
	defer m.Close()

	f, err := m.NewFunc[func(int32, int32) int32]("add")
	if err != nil {
		t.Fatal(err)
	}
	if got := f.GoType().String(); got != "func(int32, int32) int32" {
		t.Fatalf("GoType() = %q", got)
	}
	if got := f.Function().Signature().String(); got != "i32 (i32, i32)" {
		t.Fatalf("Signature() = %q", got)
	}

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(f.Function().NewBlock("entry"))
	sum := b.Add(f.Function().ParamAs[llvm.IntT](0), f.Function().ParamAs[llvm.IntT](1), "sum")
	b.Ret(sum)

	got := m.String()
	for _, want := range []string{
		"define i32 @add(i32 %0, i32 %1)",
		"%sum = add i32 %0, %1",
		"ret i32 %sum",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("module output missing %q:\n%s", want, got)
		}
	}

	// void 函数
	v, err := m.NewFunc[func()]("nothing")
	if err != nil {
		t.Fatal(err)
	}
	if got := v.Function().Signature().String(); got != "void ()" {
		t.Fatalf("void signature = %q", got)
	}

	// 不支持的类型
	if _, err := m.NewFunc[func(string)]("bad"); err == nil || err.(*llvm.Error).Reason != llvm.ErrUnsupported {
		t.Fatalf("unsupported signature should return ErrUnsupported, got %v", err)
	}
	if _, err := m.NewFunc[int32]("notfunc"); err == nil || err.(*llvm.Error).Reason != llvm.ErrUnsupported {
		t.Fatalf("non-func should return ErrUnsupported, got %v", err)
	}
}
