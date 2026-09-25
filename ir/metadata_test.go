package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestMetadata(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "md")
	defer m.Close()

	s := ctx.MDString("hello")
	if !s.IsString() || s.IsNode() || s.StringValue() != "hello" {
		t.Fatalf("MDString = %v/%v/%s", s.IsString(), s.IsNode(), s.StringValue())
	}
	n := ctx.MDNode(s, ctx.ValueAsMetadata(ctx.ConstInt(ctx.Int(32), 7).Value))
	if !n.IsNode() || n.IsString() || len(n.Operands()) != 2 {
		t.Fatalf("MDNode = %v/%v/%d", n.IsNode(), n.IsString(), len(n.Operands()))
	}
	if got := n.Operands()[0]; got.StringValue() != "hello" {
		t.Fatalf("operand 0 = %s", got)
	}

	m.AddNamedMetadataOperand("my.md", n)
	if m.NamedMetadataCount("my.md") != 1 {
		t.Fatalf("named md count = %d", m.NamedMetadataCount("my.md"))
	}
	ops := m.NamedMetadataOperands("my.md")
	if len(ops) != 1 || !ops[0].IsNode() || len(ops[0].Operands()) != 2 {
		t.Fatalf("named md operands = %v", ops)
	}
	if m.NamedMetadataCount("no.such.md") != 0 {
		t.Fatal("missing named md should have 0 operands")
	}

	m.AddModuleFlag(llvm.ModuleFlagOverride, "my.flag", ctx.MDString("v"))
	flag, ok := m.ModuleFlag("my.flag")
	if !ok || !flag.IsString() || flag.StringValue() != "v" {
		t.Fatalf("module flag = %v %v", flag, ok)
	}
	if _, ok := m.ModuleFlag("no.such.flag"); ok {
		t.Fatal("missing module flag should be absent")
	}

	got := m.String()
	for _, want := range []string{
		"!my.md = !{!0}",
		`!0 = !{!"hello", i32 7}`,
		"!llvm.module.flags = !{!1}",
		`!"my.flag"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q:\n%s", want, got)
		}
	}

	// 非 MDString/MDNode 的操作数访问
	if err := llvm.Catch(func() { _ = ctx.MDString("x").Operands() }); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("Operands on MDString should panic ErrTypeMismatch, got %v", err)
	}
	ctx2 := llvm.NewContext()
	defer ctx2.Close()
	if err := llvm.Catch(func() { ctx.MDNode(ctx2.MDString("x")) }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("cross-context MDNode should panic ErrCrossContext, got %v", err)
	}
}

func TestInstMetadata(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "md")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	entry := fn.NewBlock("entry")

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	arg := fn.ParamAs[llvm.IntT](0)
	add := b.Add(arg, arg, "x")
	AttachMetadata(add, "my.kind", ctx.MDNode(ctx.MDString("tag")))
	md, ok := InstMetadata(add, "my.kind")
	if !ok || !md.IsNode() || len(md.Operands()) != 1 || md.Operands()[0].StringValue() != "tag" {
		t.Fatalf("inst metadata = %v %v", md, ok)
	}
	if _, ok := InstMetadata(add, "other.kind"); ok {
		t.Fatal("missing inst metadata should be absent")
	}
	b.Ret(add)

	if err := llvm.Catch(func() { AttachMetadata(add, "bad.kind", ctx.MDString("x")) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("non-MDNode attachment should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { AttachMetadata(nil, "my.kind", ctx.MDNode(ctx.MDString("x"))) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil instruction attachment should panic ErrInvalidArg, got %v", err)
	}
	ctx2 := llvm.NewContext()
	defer ctx2.Close()
	if err := llvm.Catch(func() {
		AttachMetadata(add, "my.kind", ctx2.MDNode(ctx2.MDString("x")))
	}); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("cross-context metadata should panic ErrCrossContext, got %v", err)
	}

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := m.String()
	for _, want := range []string{
		"!my.kind !0",
		`!0 = !{!"tag"}`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q:\n%s", want, got)
		}
	}
}

func TestInlineAsm(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "asm")
	defer m.Close()

	sig := ctx.Fn(ctx.Void(), nil, false)
	asm := ctx.InlineAsm(sig, "nop", "", true, false, llvm.InlineAsmATT, false)
	fn := m.NewFunction("f", sig)
	entry := fn.NewBlock("entry")

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)
	b.Call[llvm.VoidT](asm, nil, "")
	b.RetVoid()

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	if got := m.String(); !strings.Contains(got, `call void asm sideeffect "nop", ""()`) {
		t.Fatalf("missing inline asm call:\n%s", got)
	}
}

func TestBlockAddress(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "ba")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, nil, false))
	entry := fn.NewBlock("entry")
	next := fn.NewBlock("next")

	ba := BlockAddress(fn, next)
	if got := BlockAddressFunction(ba); got.Name() != "f" {
		t.Fatalf("blockaddress function = %s", got.Name())
	}
	if got := BlockAddressBlock(ba); got != next {
		t.Fatalf("blockaddress block = %s", got.Name())
	}

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)
	b.Br(next)
	b.MoveToEnd(next)
	b.Ret(ctx.ConstInt(i32, 0).Value)

	// 让常量被引用，否则打印时会被省略
	g := m.NewGlobal("g", ctx.Ptr(0))
	g.SetInitializer(ba)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	if got := m.String(); !strings.Contains(got, "blockaddress(@f, %next)") {
		t.Fatalf("missing blockaddress:\n%s", got)
	}

	// 跨 Context 的函数/块（崩溃类地板：始终校验）
	ctx2 := llvm.NewContext()
	defer ctx2.Close()
	m2 := NewModule(ctx2, "ba2")
	defer m2.Close()
	i32b := ctx2.Int(32)
	fn2 := m2.NewFunction("f2", ctx2.Fn(i32b, nil, false))
	blk2 := fn2.NewBlock("entry")
	if err := llvm.Catch(func() { BlockAddress(fn2, next) }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("foreign block should panic ErrCrossContext, got %v", err)
	}
	if err := llvm.Catch(func() { BlockAddress(fn, blk2) }); err == nil || err.Reason != llvm.ErrCrossContext {
		t.Fatalf("foreign function should panic ErrCrossContext, got %v", err)
	}
}

func TestComdat(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "comdat")
	defer m.Close()

	i32 := ctx.Int(32)
	c := m.GetOrInsertComdat("mycomdat")
	if c.Name() != "mycomdat" || c.SelectionKind() != llvm.ComdatAny {
		t.Fatalf("comdat = %s %v", c.Name(), c.SelectionKind())
	}
	c.SetSelectionKind(llvm.ComdatNoDeduplicate)
	if c.SelectionKind() != llvm.ComdatNoDeduplicate {
		t.Fatalf("selection kind = %v", c.SelectionKind())
	}

	g := m.NewGlobal("g", i32)
	if _, ok := g.Comdat(); ok {
		t.Fatal("global should have no comdat initially")
	}
	g.SetComdat(c)
	got, ok := g.Comdat()
	if !ok || got.SelectionKind() != llvm.ComdatNoDeduplicate {
		t.Fatalf("global comdat = %v %v", got.SelectionKind(), ok)
	}

	// 同名 GetOrInsert 返回同一 comdat
	if m.GetOrInsertComdat("mycomdat").SelectionKind() != llvm.ComdatNoDeduplicate {
		t.Fatal("GetOrInsertComdat should return the existing comdat")
	}

	if got := m.String(); !strings.Contains(got, "$mycomdat = comdat nodeduplicate") || !strings.Contains(got, "comdat($mycomdat)") {
		t.Fatalf("missing comdat IR:\n%s", got)
	}

	if err := llvm.Catch(func() { m.GetOrInsertComdat("") }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("empty comdat name should panic ErrInvalidArg, got %v", err)
	}

	// 访问器：Ref/Context/Lifetime 与零值 Check
	if c.Ref().IsNil() {
		t.Fatal("comdat ref should not be nil")
	}
	if c.Context() != ctx || c.Lifetime() == nil || !c.Lifetime().Alive() {
		t.Fatalf("comdat accessors = %v %v", c.Context(), c.Lifetime())
	}
	if err := llvm.Catch(func() { Comdat{}.Name() }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("zero comdat should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { Comdat{}.SetSelectionKind(llvm.ComdatAny) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("zero comdat SetSelectionKind should panic ErrInvalidArg, got %v", err)
	}
}
