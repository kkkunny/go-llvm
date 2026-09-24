package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestAttributeKinds(t *testing.T) {
	if llvm.AttributeKindForName("noinline") == 0 {
		t.Fatal("noinline kind should exist")
	}
	if llvm.AttributeKindForName("definitely-not-an-attribute") != 0 {
		t.Fatal("unknown kind should be 0")
	}
	if err := llvm.Catch(func() { llvm.MustAttributeKind("definitely-not-an-attribute") }); err == nil || err.Reason != llvm.ErrNotFound {
		t.Fatalf("unknown kind should panic ErrNotFound, got %v", err)
	}

	ctx := llvm.NewContext()
	defer ctx.Close()
	a := ctx.EnumAttr(llvm.AttrNoInline, 0)
	if !a.IsEnum() || a.IsString() || a.IsType() {
		t.Fatal("noinline should be an enum attribute")
	}
	if a.EnumKind() != llvm.AttrNoInline || a.EnumValue() != 0 || a.String() != "noinline" {
		t.Fatalf("noinline attr = %d/%d/%s", a.EnumKind(), a.EnumValue(), a.String())
	}
	al := ctx.AlignAttr(8)
	if al.String() != "align 8" {
		t.Fatalf("align attr = %s", al.String())
	}
	s := ctx.StringAttr("key", "val")
	if !s.IsString() || s.StringKind() != "key" || s.StringValue() != "val" || s.String() != `"key"="val"` {
		t.Fatalf("string attr = %s/%s/%s", s.StringKind(), s.StringValue(), s.String())
	}
	byval := ctx.ByValAttr(ctx.Int(32))
	if !byval.IsType() || byval.String() != "byval(i32)" {
		t.Fatalf("byval attr = %s", byval.String())
	}

	if err := llvm.Catch(func() { ctx.EnumAttr(0, 0) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("zero kind should panic ErrInvalidArg, got %v", err)
	}
}

func TestFunctionAttrs(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "attrs")
	defer m.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{ptr}, false))
	fn.AddAttr(llvm.AttrFunction, ctx.EnumAttr(llvm.AttrNoInline, 0))
	fn.AddAttr(llvm.AttrFunction, ctx.StringAttr("my-attr", "v1"))
	fn.AddAttr(llvm.AttrParam(0), ctx.AlignAttr(8))
	fn.AddAttr(llvm.AttrReturn, ctx.EnumAttr(llvm.AttrNoUndef, 0))

	if got, ok := fn.EnumAttr(llvm.AttrFunction, llvm.AttrNoInline); !ok || got.String() != "noinline" {
		t.Fatalf("function noinline = %v %v", got, ok)
	}
	if got, ok := fn.StringAttr(llvm.AttrFunction, "my-attr"); !ok || got.StringValue() != "v1" {
		t.Fatalf("function string attr = %v %v", got, ok)
	}
	if got, ok := fn.EnumAttr(llvm.AttrParam(0), llvm.AttrAlign); !ok || got.EnumValue() != 8 {
		t.Fatalf("param align = %v %v", got, ok)
	}
	if fn.AttrCount(llvm.AttrFunction) != 2 {
		t.Fatalf("function attr count = %d, want 2", fn.AttrCount(llvm.AttrFunction))
	}
	if fn.Param(0).Index() != 0 || fn.Param(0).AttrCount() != 1 {
		t.Fatalf("param index/count = %d/%d", fn.Param(0).Index(), fn.Param(0).AttrCount())
	}
	if got := m.String(); !strings.Contains(got, "attributes #0 = {") || !strings.Contains(got, `"my-attr"="v1"`) {
		t.Fatalf("function attr group missing:\n%s", got)
	}

	if fn.CallConv() != llvm.CallConvC {
		t.Fatalf("default call conv = %v", fn.CallConv())
	}
	fn.SetCallConv(llvm.CallConvFast)
	if fn.CallConv() != llvm.CallConvFast {
		t.Fatalf("call conv after set = %v", fn.CallConv())
	}

	fn.RemoveStringAttr(llvm.AttrFunction, "my-attr")
	if _, ok := fn.StringAttr(llvm.AttrFunction, "my-attr"); ok {
		t.Fatal("string attr should be removed")
	}
	fn.RemoveEnumAttr(llvm.AttrFunction, llvm.AttrNoInline)
	if _, ok := fn.EnumAttr(llvm.AttrFunction, llvm.AttrNoInline); ok {
		t.Fatal("enum attr should be removed")
	}

	entry := fn.NewBlock("entry")
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)
	b.Ret(ctx.ConstInt(i32, 0).Value)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := m.String()
	for _, want := range []string{
		"define fastcc noundef i32 @f(ptr align 8 %0)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q:\n%s", want, got)
		}
	}

	// 跨 Context 属性
	ctx2 := llvm.NewContext()
	defer ctx2.Close()
	foreign := ctx2.EnumAttr(llvm.AttrNoInline, 0)
	if err := llvm.Catch(func() { fn.AddAttr(llvm.AttrFunction, foreign) }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("foreign attr should panic ErrCrossContext, got %v", err)
	}
}

func TestCallSiteAttrs(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "attrs")
	defer m.Close()

	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))
	g := m.NewFunction("g", ctx.Fn(ctx.Void(), nil, false))
	entry := fn.NewBlock("entry")

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	call := b.Call[llvm.VoidT](g, nil, "")
	call.AddAttr(llvm.AttrFunction, ctx.EnumAttr(llvm.AttrNoInline, 0))
	call.AddAttr(llvm.AttrFunction, ctx.StringAttr("site", "yes"))
	call.SetCallConv(llvm.CallConvCold)

	if call.CallConv() != llvm.CallConvCold {
		t.Fatalf("call conv = %v", call.CallConv())
	}
	if got, ok := call.EnumAttr(llvm.AttrFunction, llvm.AttrNoInline); !ok || !got.IsEnum() {
		t.Fatalf("call noinline = %v %v", got, ok)
	}
	if call.AttrCount(llvm.AttrFunction) != 2 {
		t.Fatalf("call attr count = %d, want 2", call.AttrCount(llvm.AttrFunction))
	}

	b.RetVoid()
	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := m.String()
	if !strings.Contains(got, "call coldcc void @g()") {
		t.Fatalf("missing cold call:\n%s", got)
	}

	call.RemoveEnumAttr(llvm.AttrFunction, llvm.AttrNoInline)
	call.RemoveStringAttr(llvm.AttrFunction, "site")
	if call.AttrCount(llvm.AttrFunction) != 0 {
		t.Fatalf("call attr count after remove = %d, want 0", call.AttrCount(llvm.AttrFunction))
	}
}

func TestInvokeAttrs(t *testing.T) {
	ctx, m, b := buildEHModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	pers := m.NewFunction("pers", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	h := m.NewFunction("h", ctx.Fn(ctx.Void(), nil, false))
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))
	fn.SetPersonality(pers)
	entry := fn.NewBlock("entry")
	cont := fn.NewBlock("cont")
	lpad := fn.NewBlock("lpad")

	b.MoveToEnd(entry)
	iv := b.Invoke[llvm.VoidT](h, nil, cont, lpad, "")
	iv.AddAttr(llvm.AttrFunction, ctx.EnumAttr(llvm.AttrNoInline, 0))
	if iv.CallConv() != llvm.CallConvC {
		t.Fatalf("invoke call conv = %v", iv.CallConv())
	}
	iv.SetCallConv(llvm.CallConvFast)
	if iv.CallConv() != llvm.CallConvFast {
		t.Fatalf("invoke call conv after set = %v", iv.CallConv())
	}
	if got, ok := iv.EnumAttr(llvm.AttrFunction, llvm.AttrNoInline); !ok || !got.IsEnum() {
		t.Fatalf("invoke noinline = %v %v", got, ok)
	}

	b.MoveToEnd(cont)
	b.RetVoid()
	b.MoveToEnd(lpad)
	lp := b.LandingPad(ctx.Struct([]llvm.AnyType{ptr, i32}, false), "lp")
	lp.SetCleanup(true)
	b.Resume(lp)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	if got := m.String(); !strings.Contains(got, "invoke fastcc void @h()") {
		t.Fatalf("missing fastcc invoke:\n%s", got)
	}
}
