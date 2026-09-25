package llvm

import (
	"fmt"
	"testing"
)

func TestAttrParam(t *testing.T) {
	// AttrReturn/AttrFunction 直接对应 LLVM-C 的 0/-1 特殊索引
	if AttrReturn != 0 || AttrFunction != -1 {
		t.Fatalf("AttrReturn=%d AttrFunction=%d", AttrReturn, AttrFunction)
	}
	if got := AttrParam(0); got != 1 {
		t.Fatalf("AttrParam(0) = %d, want 1", got)
	}
	if got := AttrParam(3); got != 4 {
		t.Fatalf("AttrParam(3) = %d, want 4", got)
	}
}

func TestAttributeKindName(t *testing.T) {
	if got := attributeKindName(0); got != "#0" {
		t.Fatalf("attributeKindName(0) = %q, want #0", got)
	}
	if got := attributeKindName(AttrNoInline); got != "noinline" {
		t.Fatalf("attributeKindName(AttrNoInline) = %q, want noinline", got)
	}
	// LLVM 认识但反查表未收录的属性名回退为 #ID 形式
	unlisted := AttributeKindForName("norecurse")
	if unlisted == 0 {
		t.Fatal("LLVM 应认识 norecurse 属性")
	}
	if got, want := attributeKindName(unlisted), fmt.Sprintf("#%d", uint32(unlisted)); got != want {
		t.Fatalf("attributeKindName(unlisted) = %q, want %q", got, want)
	}
}

func TestAttributeEnum(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	a := ctx.EnumAttr(AttrAlwaysInline, 0)
	cast := AttributeOf(ctx, a.Ref())
	if cast.Ref() != a.Ref() || cast.Context() != ctx {
		t.Fatal("AttributeOf 应保留底层句柄与所属 Context")
	}
	cast.Check("llvm.Test.Attribute.Check") // 显式校验入口

	if !a.IsEnum() || a.IsString() || a.IsType() {
		t.Fatalf("枚举属性判定错误: enum=%v string=%v type=%v", a.IsEnum(), a.IsString(), a.IsType())
	}
	if got := a.EnumKind(); got != AttrAlwaysInline {
		t.Fatalf("EnumKind() = %d, want %d", got, AttrAlwaysInline)
	}
	if got := a.EnumValue(); got != 0 {
		t.Fatalf("EnumValue() = %d, want 0", got)
	}
	if got := a.String(); got != "alwaysinline" {
		t.Fatalf("String() = %q, want alwaysinline", got)
	}

	aligned := ctx.AlignAttr(16)
	if got := aligned.EnumKind(); got != AttrAlign {
		t.Fatalf("AlignAttr kind = %d, want %d", got, AttrAlign)
	}
	if got := aligned.EnumValue(); got != 16 {
		t.Fatalf("AlignAttr value = %d, want 16", got)
	}
	if got := aligned.String(); got != "align 16" {
		t.Fatalf("AlignAttr String() = %q, want %q", got, "align 16")
	}

	if got := ctx.DereferenceableAttr(32).String(); got != "dereferenceable 32" {
		t.Fatalf("DereferenceableAttr String() = %q", got)
	}
	if got := ctx.EnumAttrName("noinline", 0).String(); got != "noinline" {
		t.Fatalf("EnumAttrName String() = %q", got)
	}
	// 带值的枚举属性打印为 "名称 值"
	if got := ctx.EnumAttr(AttrAlign, 8).String(); got != "align 8" {
		t.Fatalf("带值枚举属性 String() = %q", got)
	}
}

func TestAttributeStringAttr(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	sa := ctx.StringAttr("some-key", "some-value")
	if !sa.IsString() || sa.IsEnum() || sa.IsType() {
		t.Fatalf("字符串属性判定错误: enum=%v string=%v type=%v", sa.IsEnum(), sa.IsString(), sa.IsType())
	}
	if got := sa.StringKind(); got != "some-key" {
		t.Fatalf("StringKind() = %q", got)
	}
	if got := sa.StringValue(); got != "some-value" {
		t.Fatalf("StringValue() = %q", got)
	}
	if got := sa.String(); got != `"some-key"="some-value"` {
		t.Fatalf("String() = %q", got)
	}
}

func TestAttributeTypeAttr(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()
	i32 := ctx.Int(32)

	ta := ctx.TypeAttr(AttrByVal, i32)
	if !ta.IsType() || ta.IsEnum() || ta.IsString() {
		t.Fatalf("类型属性判定错误: enum=%v string=%v type=%v", ta.IsEnum(), ta.IsString(), ta.IsType())
	}
	if got := ta.TypeValue(); !got.Equal(i32) {
		t.Fatalf("TypeValue() = %s, want i32", got)
	}
	if got := ta.String(); got != "byval(i32)" {
		t.Fatalf("TypeAttr String() = %q, want byval(i32)", got)
	}
	if got := ctx.ByValAttr(ctx.Int(64)).String(); got != "byval(i64)" {
		t.Fatalf("ByValAttr String() = %q", got)
	}
	if got := ctx.SRetAttr(i32).String(); got != "sret(i32)" {
		t.Fatalf("SRetAttr String() = %q", got)
	}
}

// TestAttributeCrashFloor 崩溃类地板负向：nil 句柄、已释放 Context、跨 Context，
// 两构建模式都必须 panic（不得跳过）。
func TestAttributeCrashFloor(t *testing.T) {
	if err := Catch(func() { Attribute{}.IsEnum() }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("nil 属性应 panic ErrInvalidArg, got %v", err)
	}

	ctx := NewContext()
	ref := ctx.EnumAttr(AttrCold, 0).Ref()
	if err := Catch(func() { (Attribute{ref: ref}).IsEnum() }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("ctx 为 nil 的属性应 panic ErrInvalidArg, got %v", err)
	}

	attr := ctx.EnumAttr(AttrCold, 0)
	_ = ctx.Close()
	if err := Catch(func() { attr.IsEnum() }); err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("已释放 Context 上的属性应 panic ErrUseAfterFree, got %v", err)
	}

	ctx2 := NewContext()
	defer ctx2.Close()
	if err := Catch(func() { ctx2.TypeAttr(AttrByVal, nil) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("nil 类型应 panic ErrInvalidArg, got %v", err)
	}
	ctx3 := NewContext()
	defer ctx3.Close()
	if err := Catch(func() { ctx2.TypeAttr(AttrByVal, ctx3.Int(32)) }); err == nil || err.Reason != ErrCrossContext {
		t.Fatalf("跨 Context 类型应 panic ErrCrossContext, got %v", err)
	}
}

// TestAttributeWrongKind 语义契约负向：访问器用于错误种类的属性 / 非法 kind 构造。
func TestAttributeWrongKind(t *testing.T) {
	requireDebug(t)

	ctx := NewContext()
	defer ctx.Close()
	enum := ctx.EnumAttr(AttrNoInline, 0)
	str := ctx.StringAttr("k", "v")
	typ := ctx.TypeAttr(AttrByVal, ctx.Int(32))

	tests := []struct {
		name string
		fn   func()
		want ErrKind
	}{
		{"EnumKind(str)", func() { str.EnumKind() }, ErrTypeMismatch},
		{"EnumValue(str)", func() { str.EnumValue() }, ErrTypeMismatch},
		{"TypeValue(enum)", func() { enum.TypeValue() }, ErrTypeMismatch},
		{"StringKind(enum)", func() { enum.StringKind() }, ErrTypeMismatch},
		{"StringValue(enum)", func() { enum.StringValue() }, ErrTypeMismatch},
		{"StringKind(typ)", func() { typ.StringKind() }, ErrTypeMismatch},
		{"EnumAttr(0)", func() { ctx.EnumAttr(0, 0) }, ErrInvalidArg},
		{"EnumAttrName(未知)", func() { ctx.EnumAttrName("definitely-not-an-attribute", 0) }, ErrNotFound},
		{"TypeAttr(0)", func() { ctx.TypeAttr(0, ctx.Int(32)) }, ErrInvalidArg},
		{"MustAttributeKind(未知)", func() { MustAttributeKind("definitely-not-an-attribute") }, ErrNotFound},
	}
	for _, tc := range tests {
		err := Catch(tc.fn)
		if err == nil || err.Reason != tc.want {
			t.Errorf("%s: want %v, got %v", tc.name, tc.want, err)
		}
	}
}
