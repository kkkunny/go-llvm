package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

// TestAttrFunctionAndParamRoles 覆盖函数/参数两类角色的属性全量访问器。
func TestAttrFunctionAndParamRoles(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "attrs")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn.AddAttr(llvm.AttrFunction, ctx.EnumAttr(llvm.AttrNoInline, 0))
	fn.AddAttr(llvm.AttrFunction, ctx.StringAttr("f-key", "f-val"))

	// Function.Attrs 读取全部属性（枚举 + 字符串）
	fa := fn.Attrs(llvm.AttrFunction)
	if len(fa) != 2 {
		t.Fatalf("function attrs = %v, want 2", fa)
	}

	// 参数属性：添加、全量读取、按种类读取、移除
	p := fn.Param(0)
	if p.Index() != 0 {
		t.Fatalf("param index = %d, want 0", p.Index())
	}
	p.AddAttr(ctx.AlignAttr(8))
	p.AddAttr(ctx.StringAttr("p-key", "p-val"))
	if pa := p.Attrs(); len(pa) != 2 {
		t.Fatalf("param attrs = %v, want 2", pa)
	}
	if got, ok := p.EnumAttr(llvm.AttrAlign); !ok || got.EnumValue() != 8 {
		t.Fatalf("param align = %v %v", got, ok)
	}
	if got, ok := p.StringAttr("p-key"); !ok || got.StringValue() != "p-val" {
		t.Fatalf("param string attr = %v %v", got, ok)
	}
	if p.AttrCount() != 2 {
		t.Fatalf("param attr count = %d, want 2", p.AttrCount())
	}
	p.RemoveEnumAttr(llvm.AttrAlign)
	p.RemoveStringAttr("p-key")
	if p.AttrCount() != 0 {
		t.Fatalf("param attr count after remove = %d, want 0", p.AttrCount())
	}
	if _, ok := p.EnumAttr(llvm.AttrAlign); ok {
		t.Fatal("param align attr should be removed")
	}
	if _, ok := p.StringAttr("p-key"); ok {
		t.Fatal("param string attr should be removed")
	}

	entry := fn.NewBlock("entry")
	b := NewBuilderAt(entry)
	defer b.Close()
	ret := b.Ret(ctx.ConstInt(i32, 0).Value)
	_ = ret

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	out := m.String()
	if !strings.Contains(out, "attributes #0 = { noinline") {
		t.Fatalf("IR missing function attrs:\n%s", out)
	}

	// 跨 Context 的属性（崩溃类地板：始终校验）
	ctx2 := llvm.NewContext()
	defer ctx2.Close()
	foreign := ctx2.EnumAttr(llvm.AttrNoInline, 0)
	if err := llvm.Catch(func() { fn.AddAttr(llvm.AttrFunction, foreign) }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("function foreign attr should panic ErrCrossContext, got %v", err)
	}
	if err := llvm.Catch(func() { p.AddAttr(foreign) }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("param foreign attr should panic ErrCrossContext, got %v", err)
	}
}

// TestAttrInvokeRoles 覆盖 Invoke 调用点属性全量访问器；Call 已由既有测试覆盖，这里补 Attrs/StringAttr。
func TestAttrInvokeRoles(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "attrs")
	defer m.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	pers := m.NewFunction("pers", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	g := m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn.SetPersonality(pers)

	entry := fn.NewBlock("entry")
	cont := fn.NewBlock("cont")
	lpad := fn.NewBlock("lpad")

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	iv := b.Invoke[llvm.IntT](g, []llvm.AnyValue{ctx.ConstInt(i32, 1).Value}, cont, lpad, "iv")
	iv.AddAttr(llvm.AttrFunction, ctx.EnumAttr(llvm.AttrNoInline, 0))
	iv.AddAttr(llvm.AttrFunction, ctx.StringAttr("iv-key", "iv-val"))
	if iv.AttrCount(llvm.AttrFunction) != 2 {
		t.Fatalf("invoke attr count = %d, want 2", iv.AttrCount(llvm.AttrFunction))
	}
	if ia := iv.Attrs(llvm.AttrFunction); len(ia) != 2 {
		t.Fatalf("invoke attrs = %v, want 2", ia)
	}
	if got, ok := iv.EnumAttr(llvm.AttrFunction, llvm.AttrNoInline); !ok || !got.IsEnum() {
		t.Fatalf("invoke noinline = %v %v", got, ok)
	}
	if got, ok := iv.StringAttr(llvm.AttrFunction, "iv-key"); !ok || got.StringValue() != "iv-val" {
		t.Fatalf("invoke string attr = %v %v", got, ok)
	}
	iv.RemoveEnumAttr(llvm.AttrFunction, llvm.AttrNoInline)
	iv.RemoveStringAttr(llvm.AttrFunction, "iv-key")
	if iv.AttrCount(llvm.AttrFunction) != 0 {
		t.Fatalf("invoke attr count after remove = %d, want 0", iv.AttrCount(llvm.AttrFunction))
	}

	// Call 角色的 Attrs/StringAttr 补缺（放在 invoke 的正常出口块）
	b.MoveToEnd(cont)
	call := b.Call[llvm.IntT](g, []llvm.AnyValue{ctx.ConstInt(i32, 1).Value}, "c")
	call.AddAttr(llvm.AttrFunction, ctx.EnumAttr(llvm.AttrNoInline, 0))
	call.AddAttr(llvm.AttrFunction, ctx.StringAttr("c-key", "c-val"))
	if ca := call.Attrs(llvm.AttrFunction); len(ca) != 2 {
		t.Fatalf("call attrs = %v, want 2", ca)
	}
	if got, ok := call.StringAttr(llvm.AttrFunction, "c-key"); !ok || got.StringValue() != "c-val" {
		t.Fatalf("call string attr = %v %v", got, ok)
	}

	b.Ret(iv)
	b.MoveToEnd(lpad)
	lp := b.LandingPad(ctx.Struct([]llvm.AnyType{ptr, i32}, false), "lp")
	lp.SetCleanup(true)
	b.Resume(lp)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	out := m.String()
	if !strings.Contains(out, `"c-key"="c-val"`) {
		t.Fatalf("IR missing call site string attr:\n%s", out)
	}

	// 跨 Context 的调用点属性（崩溃类地板：始终校验）
	ctx2 := llvm.NewContext()
	defer ctx2.Close()
	foreign := ctx2.EnumAttr(llvm.AttrNoInline, 0)
	if err := llvm.Catch(func() { call.AddAttr(llvm.AttrFunction, foreign) }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("call foreign attr should panic ErrCrossContext, got %v", err)
	}
	if err := llvm.Catch(func() { iv.AddAttr(llvm.AttrFunction, foreign) }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("invoke foreign attr should panic ErrCrossContext, got %v", err)
	}
}
