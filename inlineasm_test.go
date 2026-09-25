package llvm

import (
	"strings"
	"testing"
)

func TestInlineAsm(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	sig := ctx.Fn(ctx.Void(), nil, false)
	asm := ctx.InlineAsm(sig, "nop", "", true, false, InlineAsmATT, false)
	if asm.IsNil() {
		t.Fatal("InlineAsm 句柄不应为空")
	}
	// LLVM 22 的内联汇编常量为不透明指针（函数指针）类型
	if got := asm.Type().String(); got != "ptr" {
		t.Fatalf("InlineAsm Type() = %q, want ptr", got)
	}
	got := asm.String()
	for _, want := range []string{"sideeffect", `"nop"`, `""`} {
		if !strings.Contains(got, want) {
			t.Fatalf("InlineAsm String() 缺少 %q: %s", want, got)
		}
	}

	// Intel 方言 / 栈对齐 / canThrow 的打印差异
	intel := ctx.InlineAsm(sig, "nop", "~{dirflag}", false, true, InlineAsmIntel, true)
	got = intel.String()
	for _, want := range []string{"inteldialect", "alignstack", "unwind", `"nop"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("Intel InlineAsm String() 缺少 %q: %s", want, got)
		}
	}

	// 非函数类型（nil 类型角色）应 panic ErrInvalidArg
	if err := Catch(func() { ctx.InlineAsm(FnType{}, "nop", "", false, false, InlineAsmATT, false) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("nil 函数类型应 panic ErrInvalidArg, got %v", err)
	}
}
