package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestOperandBundles(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "t")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("callee", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	caller := m.NewFunction("caller", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	blk := caller.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()

	a := caller.ParamAs[llvm.IntT](0)
	tag := NewOperandBundle(ctx, "funclet", []llvm.AnyValue{a})
	if tag.Tag() != "funclet" {
		t.Fatalf("tag = %s", tag.Tag())
	}

	// invoke 带捆绑（entry 以 invoke 终结），call 放在 then 块
	then := caller.NewBlock("then")
	catch := caller.NewBlock("catch")
	b.InvokeWithBundles[llvm.IntT](fn, []llvm.AnyValue{a}, []OperandBundle{tag}, then, catch, "inv")
	b.MoveToEnd(then)

	call := b.CallWithBundles[llvm.IntT](fn, []llvm.AnyValue{a}, []OperandBundle{tag}, "c")
	if call.ArgCount() != 1 {
		t.Fatalf("args = %d", call.ArgCount())
	}
	b.Ret(call)

	out := m.String()
	if !strings.Contains(out, `"funclet"(`) {
		t.Fatalf("IR missing operand bundle:\n%s", out)
	}

	if err := tag.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if err := tag.Close(); err == nil {
		t.Fatalf("double close should error")
	}

	// 已关闭捆绑不得再用于构建（崩溃类地板：始终校验）
	if err := llvm.Catch(func() { _ = tag.Tag() }); err == nil || err.Reason != llvm.ErrUseAfterFree {
		t.Fatalf("Tag on closed bundle should panic ErrUseAfterFree, got %v", err)
	}
	if err := llvm.Catch(func() {
		b.CallWithBundles[llvm.IntT](fn, []llvm.AnyValue{a}, []OperandBundle{tag}, "")
	}); err == nil || err.Reason != llvm.ErrUseAfterFree {
		t.Fatalf("CallWithBundles with closed bundle should panic ErrUseAfterFree, got %v", err)
	}
}
