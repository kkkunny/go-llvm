package ir

import (
	"path/filepath"
	"testing"

	"github.com/kkkunny/go-llvm"
)

// TestUseAfterModuleClose 模块关闭后，各角色方法必须 panic ErrUseAfterFree，而不是把悬垂句柄交给 LLVM
func TestUseAfterModuleClose(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "dead")
	i32 := ctx.Int(32)
	one := ctx.ConstInt(i32, 1, false)

	fnTy := ctx.Fn(i32, []llvm.AnyType{i32}, false)
	fn := m.NewFunction("f", fnTy)
	blk := fn.NewBlock("entry")
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(blk)
	alloca := b.Alloca(i32, "a")
	load := b.Load(alloca, i32, "v")
	b.Ret(load.Value)
	g := m.NewGlobal("g", i32)

	fn2 := m.NewFunction("f2", fnTy)
	blk2 := fn2.NewBlock("entry")
	b.MoveToEnd(blk2)
	phi := b.PHI(i32, "p")
	call := b.Call[llvm.IntT](fn.Value, []llvm.AnyValue{one}, "c")
	sw := b.Switch(one, blk2)
	sw.AddCase(one, blk)

	if err := m.Close(); err != nil {
		t.Fatal(err)
	}

	deadPath := filepath.Join(t.TempDir(), "dead.ll")
	cases := []struct {
		name string
		call func()
	}{
		{"Module.String", func() { _ = m.String() }},
		{"Module.Verify", func() { _ = m.Verify() }},
		{"Module.NewGlobal", func() { m.NewGlobal("x", i32) }},
		{"Module.NewFunction", func() { m.NewFunction("x", fnTy) }},
		{"Module.GetGlobal", func() { m.GetGlobal("g") }},
		{"Module.Clone", func() { m.Clone() }},
		{"Module.DataLayout", func() { _ = m.DataLayout() }},
		{"Module.Bitcode", func() { m.Bitcode() }},
		{"Module.WriteToFile", func() { _ = m.WriteToFile(deadPath) }},
		{"Block.Insts", func() { _ = blk.Insts() }},
		{"Block.Name", func() { _ = blk.Name() }},
		{"Block.Next", func() { blk.Next() }},
		{"Block.Belong", func() { blk.Belong() }},
		{"Function.Param", func() { _ = fn.Param(0) }},
		{"Function.NewBlock", func() { fn.NewBlock("x") }},
		{"Function.Linkage", func() { _ = fn.Linkage() }},
		{"Function.Verify", func() { _ = fn.Verify() }},
		{"Global.SetInitializer", func() { g.SetInitializer(one) }},
		{"Global.SetConstant", func() { g.SetConstant(true) }},
		{"Global.SetAlign", func() { g.SetAlign(4) }},
		{"Global.ValueType", func() { _ = g.ValueType() }},
		{"Alloca.Align", func() { _ = alloca.Align() }},
		{"Load.Align", func() { _ = load.Align() }},
		{"Call.ArgCount", func() { _ = call.ArgCount() }},
		{"Call.SetArg", func() { call.SetArg(0, one) }},
		{"Call.CalledFunction", func() { call.CalledFunction() }},
		{"Phi.Count", func() { _ = phi.Count() }},
		{"Phi.AddIncoming", func() {
			phi.AddIncoming(Incoming[llvm.IntT]{Value: load.Value, Block: blk})
		}},
		{"Switch.Count", func() { _ = sw.Count() }},
		{"Switch.AddCase", func() { sw.AddCase(one, blk) }},
		{"Switch.DefaultBlock", func() { sw.DefaultBlock() }},
		{"Switch.CaseValue", func() { _ = sw.CaseValue(0) }},
		{"Param.SetAlign", func() { fn.Param(0).SetAlign(4) }},
		{"Param.Belong", func() { fn.Param(0).Belong() }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := llvm.Catch(tc.call)
			if err == nil || err.Reason != llvm.ErrUseAfterFree {
				t.Fatalf("want ErrUseAfterFree, got %v", err)
			}
		})
	}
}

// TestNilModuleCheck nil 模块句柄必须 panic ErrInvalidArg
func TestNilModuleCheck(t *testing.T) {
	var m *Module
	if err := llvm.Catch(func() { _ = m.String() }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("want ErrInvalidArg, got %v", err)
	}
}
