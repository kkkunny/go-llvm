package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestModuleFunctionGolden(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "test")
	defer m.Close()

	i32 := ctx.Int(32)
	m.NewFunction("add", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))

	want := "declare i32 @add(i32, i32)"
	if got := m.String(); !strings.Contains(got, want) {
		t.Fatalf("module output missing %q:\n%s", want, got)
	}
	if got := m.Source(); got != "test" {
		t.Fatalf("Source() = %q", got)
	}
}

func TestModuleLookup(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "lookup")
	defer m.Close()

	i32 := ctx.Int(32)
	m.NewFunction("foo", ctx.Fn(i32, nil, false))

	got, ok := m.GetFunction("foo")
	if !ok || got.Name() != "foo" {
		t.Fatalf("GetFunction(foo) = %v, %v", got.Name(), ok)
	}
	if _, ok := m.GetFunction("bar"); ok {
		t.Fatal("GetFunction(bar) should not exist")
	}

	g := m.NewGlobal("g", i32)
	g.SetInitializer(ctx.ConstInt(i32, 7))
	gg, ok := m.GetGlobal("g")
	if !ok {
		t.Fatalf("GetGlobal(g) = %v", ok)
	}
	init, hasInit := gg.Initializer()
	if !hasInit || init.String() != "i32 7" {
		t.Fatalf("Initializer() = %v, %v", init, hasInit)
	}
	if _, ok := m.GetGlobal("nope"); ok {
		t.Fatal("GetGlobal(nope) should not exist")
	}

	m.DelGlobal(gg)
	if _, ok := m.GetGlobal("g"); ok {
		t.Fatal("global should be deleted")
	}
}

func TestFunctionShape(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "fn")
	defer m.Close()

	i32 := ctx.Int(32)
	f64 := ctx.Float(llvm.FloatDouble)
	fn := m.NewFunction("mix", ctx.Fn(f64, []llvm.AnyType{i32, f64}, true))

	if fn.Name() != "mix" {
		t.Fatalf("Name() = %q", fn.Name())
	}
	if fn.CountParams() != 2 {
		t.Fatalf("CountParams() = %d", fn.CountParams())
	}
	if !fn.Signature().IsVarArg() {
		t.Fatal("should be vararg")
	}
	if got := fn.Signature().Return().String(); got != "double" {
		t.Fatalf("return = %q", got)
	}
	if got := fn.Param(0).Type().String(); got != "i32" {
		t.Fatalf("param0 type = %q", got)
	}
	if got := fn.ParamAs[llvm.FloatT](1).Type().String(); got != "double" {
		t.Fatalf("param1 type = %q", got)
	}
	if !fn.OnlyDecl() {
		t.Fatal("function without blocks should be declaration only")
	}

	blk := fn.NewBlock("entry")
	if fn.OnlyDecl() {
		t.Fatal("function with block should not be declaration only")
	}
	if got := blk.Name(); got != "entry" {
		t.Fatalf("block name = %q", got)
	}
	if blk.Belong().Name() != "mix" {
		t.Fatalf("Belong() = %q", blk.Belong().Name())
	}
	if blocks := fn.Blocks(); len(blocks) != 1 {
		t.Fatalf("Blocks() = %d", len(blocks))
	}
	if entry, ok := fn.EntryBlock(); !ok || entry.Name() != "entry" {
		t.Fatalf("EntryBlock() = %v, %v", entry, ok)
	}
	if !blk.Empty() {
		t.Fatal("new block should be empty")
	}
}

func TestFunctionLinkage(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "linkage")
	defer m.Close()

	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))
	if fn.Linkage() != llvm.LinkageExternal {
		t.Fatalf("default linkage = %v", fn.Linkage())
	}
	fn.SetLinkage(llvm.LinkageInternal)
	if fn.Linkage() != llvm.LinkageInternal {
		t.Fatalf("linkage = %v", fn.Linkage())
	}
}

func TestModuleVerify(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "verify")
	defer m.Close()

	if err := m.Verify(); err != nil {
		t.Fatalf("empty module should verify: %v", err)
	}

	m.NewFunction("decl", ctx.Fn(ctx.Void(), nil, false))
	if err := m.Verify(); err != nil {
		t.Fatalf("declaration-only module should verify: %v", err)
	}
}

func TestModuleCloseLifecycle(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "life")
	i32 := ctx.Int(32)
	g := m.NewGlobal("g", i32)
	v := g

	if !v.Alive() {
		t.Fatal("value should be alive while module open")
	}
	if err := m.Close(); err != nil {
		t.Fatal(err)
	}
	if v.Alive() {
		t.Fatal("value should be dead after module close")
	}
	if err := m.Close(); err == nil || err.(*llvm.Error).Reason != llvm.ErrClosed {
		t.Fatalf("want ErrClosed, got %v", err)
	}
}

func TestModuleDisown(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "disown")
	i32 := ctx.Int(32)
	g := m.NewGlobal("g", i32)
	v := g

	m.Disown()
	if v.Alive() {
		t.Fatal("value should be dead immediately after disown")
	}
	if err := llvm.Catch(func() { _ = m.String() }); err == nil || err.Reason != llvm.ErrUseAfterFree {
		t.Fatalf("use after disown should panic ErrUseAfterFree, got %v", err)
	}
	// 注销后 Context.Close 不应触碰该模块
	if err := m.Close(); err == nil || err.(*llvm.Error).Reason != llvm.ErrClosed {
		t.Fatalf("Close after Disown should return ErrClosed, got %v", err)
	}
}

func TestBelongAfterModuleClose(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "belong")
	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	blk := fn.NewBlock("entry")

	belong := blk.Belong()
	param := belong.Param(0)
	if got := param.Type().String(); got != "i32" {
		t.Fatalf("param type = %q", got)
	}
	if got := belong.NewBlock("entry2"); got.Name() != "entry2" {
		t.Fatalf("NewBlock name = %q", got.Name())
	}
	if blocks := belong.Blocks(); len(blocks) != 2 {
		t.Fatalf("Blocks() = %d", len(blocks))
	}
}

func TestModuleCloseWithContext(t *testing.T) {
	ctx := llvm.NewContext()
	m := NewModule(ctx, "ctx-close")
	if err := ctx.Close(); err != nil {
		t.Fatal(err)
	}
	if err := m.Close(); err == nil || err.(*llvm.Error).Reason != llvm.ErrClosed {
		t.Fatalf("module should be closed by context cascade, got %v", err)
	}
}

func TestModuleClone(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "orig")
	defer m.Close()

	i32 := ctx.Int(32)
	m.NewFunction("f", ctx.Fn(i32, nil, false))

	clone := m.Clone()
	defer clone.Close()
	if clone.String() != m.String() {
		t.Fatalf("clone mismatch:\n%s\n---\n%s", clone.String(), m.String())
	}
	if _, ok := clone.GetFunction("f"); !ok {
		t.Fatal("clone should contain function f")
	}
}

func TestGlobal(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "global")
	defer m.Close()

	i32 := ctx.Int(32)
	g := m.NewGlobal("g", i32)
	if got := g.Type().String(); got != "ptr" {
		t.Fatalf("global value type = %q", got)
	}
	if g.IsConstant() {
		t.Fatal("new global should not be constant")
	}
	g.SetConstant(true)
	if !g.IsConstant() {
		t.Fatal("global should be constant")
	}
	g.SetAlign(8)
	if g.Align() != 8 {
		t.Fatalf("align = %d", g.Align())
	}
	if _, ok := g.Initializer(); ok {
		t.Fatal("global without initializer should return false")
	}
	g.SetInitializer(ctx.ConstInt(i32, 3))
	if init, ok := g.Initializer(); !ok || init.String() != "i32 3" {
		t.Fatalf("Initializer() = %v, %v", init, ok)
	}

	c := m.NewGlobalConst("c", ctx.ConstInt(i32, 5))
	if !c.IsConstant() {
		t.Fatal("NewGlobalConst should be constant")
	}
	ci, ok := c.Initializer()
	if !ok || ci.String() != "i32 5" {
		t.Fatalf("constant initializer = %v, %v", ci, ok)
	}
}

func TestModuleSourceAndTriple(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "src")
	defer m.Close()

	m.SetSource("file.c")
	if got := m.Source(); got != "file.c" {
		t.Fatalf("Source() = %q", got)
	}
	m.SetTargetTriple("x86_64-unknown-linux-gnu")
	if got := m.TargetTriple(); got != "x86_64-unknown-linux-gnu" {
		t.Fatalf("TargetTriple() = %q", got)
	}
	m.SetDataLayout("e-m:e-p:64:64-i64:64-n8:16:32:64-S128")
	if !strings.Contains(m.String(), "target datalayout") {
		t.Fatalf("module output missing datalayout:\n%s", m.String())
	}
}

func TestModuleIteration(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "t")
	defer m.Close()

	i32 := ctx.Int(32)
	m.NewFunction("f1", ctx.Fn(i32, nil, false))
	m.NewFunction("f2", ctx.Fn(i32, nil, false))
	m.NewGlobal("g1", i32)
	m.NewGlobal("g2", i32)

	fns := map[string]bool{}
	for fn := range m.AllFunctions() {
		fns[fn.Name()] = true
	}
	if !fns["f1"] || !fns["f2"] {
		t.Fatalf("AllFunctions = %v", fns)
	}

	gs := map[string]bool{}
	for g := range m.AllGlobals() {
		gs[g.Name()] = true
	}
	if !gs["g1"] || !gs["g2"] {
		t.Fatalf("AllGlobals = %v", gs)
	}
}
