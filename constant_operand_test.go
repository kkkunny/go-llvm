package llvm_test

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/errs"
	"github.com/kkkunny/go-llvm/ir"
)

// newInstructionFixture 构造 f(i32) { entry: %slot = alloca i32; %ld = load i32, ptr %slot }，
// 返回非常量值 slot/ld，供"常量构造器拒绝指令操作数"等负向测试使用。
// 注意用 load 而不是 add：Builder 默认带常量折叠，常量相加的结果是常量而非指令。
// 根包内部测试不得 import ir（依赖方向 llvm ← ir），故放在外部测试包。
func newInstructionFixture(t *testing.T, name string) (*llvm.Context, *ir.Module, *ir.Builder, ir.Alloca, llvm.Value[llvm.IntT]) {
	t.Helper()
	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, name)
	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{i32}, false))
	b := ir.NewBuilderAt(fn.NewBlock("entry"))
	slot := b.Alloca(i32, "slot")
	ld := b.Load(slot, i32, "ld")
	return ctx, m, b, slot, ld.Value
}

// TestConstOperandRejectsInstructions 常量构造器遇到指令操作数时必须在源头 panic
// ErrInvalidArg（而不是构造出"常量里包含指令"的畸形 IR，等 Verify 指向那条指令）。
func TestConstOperandRejectsInstructions(t *testing.T) {
	ctx, m, b, slot, ld := newInstructionFixture(t, "const-operand")
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	i64 := ctx.Int(64)
	pair := ctx.NamedStruct("One")
	pair.SetBody([]llvm.AnyType{i32}, false)

	cases := []struct {
		name string
		op   string
		want string // 消息中应出现的违规值标识
		fn   func()
	}{
		{"ConstArray", "llvm.Context.ConstArray", "slot", func() { ctx.ConstArray(i32, slot.Value) }},
		{"ConstVector", "llvm.Context.ConstVector", "slot", func() { ctx.ConstVector(i32, slot.Value) }},
		{"ConstStruct", "llvm.Context.ConstStruct", "slot", func() { ctx.ConstStruct(false, slot.Value) }},
		{"ConstNamedStruct", "llvm.Context.ConstNamedStruct", "slot", func() { ctx.ConstNamedStruct(pair, slot.Value) }},
		{"ConstGEP-base", "llvm.Context.ConstGEP", "slot", func() {
			ctx.ConstGEP(i32, slot.Value, false, ctx.ConstInt(i64, 1).Value)
		}},
		{"ConstGEP-index", "llvm.Context.ConstGEP", "ld", func() {
			ctx.ConstGEP(i32, ctx.ConstNull(ctx.Ptr(0)), true, ld)
		}},
		{"ConstIntToPtr", "llvm.Context.ConstIntToPtr", "ld", func() { ctx.ConstIntToPtr(ld, ctx.Ptr(0)) }},
		{"ConstBitCast", "llvm.Context.ConstBitCast", "slot", func() { ctx.ConstBitCast(slot.Value, ctx.Ptr(0)) }},
		{"ConstPointerCast", "llvm.Context.ConstPointerCast", "slot", func() { ctx.ConstPointerCast(slot.Value, ctx.Ptr(0)) }},
		{"IntConst.Add", "llvm.IntConst.Add", "ld", func() { ctx.ConstInt(i32, 1).Add(ld) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := errs.Catch(tc.fn)
			if err == nil || err.Reason != llvm.ErrInvalidArg {
				t.Fatalf("want ErrInvalidArg, got %v", err)
			}
			if err.Op != tc.op {
				t.Fatalf("op = %q, want %q", err.Op, tc.op)
			}
			if !strings.Contains(err.Msg, "not a constant") || !strings.Contains(err.Msg, tc.want) {
				t.Fatalf("msg should name the offending value %q: %q", tc.want, err.Msg)
			}
		})
	}
}

// TestMustAsPanicMentionsValue MustAs 的 panic 消息必须带值的名字与 IR 文本，
// 便于在动态代码中定位转换点。
func TestMustAsPanicMentionsValue(t *testing.T) {
	ctx, m, b, slot, _ := newInstructionFixture(t, "must-as")
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	err := errs.Catch(func() { slot.Value.MustAs[llvm.IntT]() })
	if err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("want ErrTypeMismatch, got %v", err)
	}
	for _, want := range []string{"value kind mismatch", "have pointer", "want int", "slot", "alloca"} {
		if !strings.Contains(err.Msg, want) {
			t.Fatalf("msg missing %q: %q", want, err.Msg)
		}
	}
}
