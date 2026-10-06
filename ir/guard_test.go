package ir

import (
	"testing"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/errs"
)

// TestConstantOperandGuards 覆盖 LLVM unwrap<Constant> 类入口的无条件校验。
func TestConstantOperandGuards(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "guards")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()
	inst := b.Add(fn.ParamAs[llvm.IntT](0), i32.Const(2), "sum").Dyn()

	// 初始化器必须是常量
	g := m.NewGlobal("g", i32)
	if err := errs.Catch(func() { g.SetInitializer(inst) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("SetInitializer(instruction) should panic ErrInvalidArg, got %v", err)
	}
	// 别名/IFunc 目标必须是常量
	if err := errs.Catch(func() { m.NewAlias("a", i32, inst) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("NewAlias(instruction) should panic ErrInvalidArg, got %v", err)
	}
	if err := errs.Catch(func() { m.NewIFunc("if", i32, inst) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("NewIFunc(instruction) should panic ErrInvalidArg, got %v", err)
	}
	// personality / prefix / prologue 必须是常量
	if err := errs.Catch(func() { fn.SetPrefixData(inst) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("SetPrefixData(instruction) should panic ErrInvalidArg, got %v", err)
	}
	if err := errs.Catch(func() { fn.SetPrologueData(inst) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("SetPrologueData(instruction) should panic ErrInvalidArg, got %v", err)
	}
	// switch case 必须是常量
	cont := fn.NewBlock("cont")
	sw := b.Switch(i32.Const(0), cont)
	if err := errs.Catch(func() { sw.AddCase(inst.MustAs[llvm.IntT](), cont) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("AddCase(instruction) should panic ErrInvalidArg, got %v", err)
	}
}

// TestBuilderFloorGuards 覆盖 Builder 的崩溃类地板。
func TestBuilderFloorGuards(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "floor")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, nil, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	b.Ret(i32.Const(0))

	// 零值 Builder 不能触碰底层句柄
	var zero Builder
	if err := errs.Catch(func() { zero.CurrentBlock() }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("zero builder CurrentBlock should panic ErrInvalidArg, got %v", err)
	}
	// Close 后 CurrentBlock 必须可恢复地 panic，而不是 UAF
	if err := b.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if err := errs.Catch(func() { b.CurrentBlock() }); err == nil || err.Reason != llvm.ErrClosed {
		t.Fatalf("closed builder CurrentBlock should panic ErrClosed, got %v", err)
	}

	// MoveBefore 只接受指令
	b2 := NewBuilder(ctx)
	defer b2.Close()
	if err := errs.Catch(func() { b2.MoveBefore(i32.Const(1).Value.Dyn()) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("MoveBefore(constant) should panic ErrInvalidArg, got %v", err)
	}
}

// TestCmpXchgFailureOrdering unordered 不是合法的失败序（上游 isValidFailureOrdering 排除）。
func TestCmpXchgFailureOrdering(t *testing.T) {
	requireDebug(t)
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "atomic")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()
	ptr := b.Alloca(i32, "p").Value
	if err := errs.Catch(func() {
		b.CmpXchg(ptr, i32.Const(0), i32.Const(1), llvm.AtomicMonotonic, llvm.AtomicUnordered, false, "x")
	}); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("unordered failure ordering should panic ErrInvalidArg, got %v", err)
	}
	b.RetVoid()
}

// TestMetadataGuards 覆盖 MDNode 位置的值包装拒绝。
func TestMetadataGuards(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "md")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()
	inst := b.RetVoid()

	// ValueAsMetadata（如 i32 3）不能直接附着为指令元数据
	vam := ctx.ValueAsMetadata(i32.Const(3).Value)
	if err := errs.Catch(func() { AttachMetadata(inst, "prof", vam) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("AttachMetadata(ValueAsMetadata) should panic ErrInvalidArg, got %v", err)
	}
	// MDString 不能作为命名元数据操作数
	if err := errs.Catch(func() { m.AddNamedMetadataOperand("llvm.ident", ctx.MDString("x")) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("AddNamedMetadataOperand(MDString) should panic ErrInvalidArg, got %v", err)
	}
}

// TestBlockAddressCrossFunction blockaddress 的块必须属于给定函数（release 下否则静默错误）。
func TestBlockAddressCrossFunction(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "ba")
	defer m.Close()

	voidFn := ctx.Fn(ctx.Void(), nil, false)
	f1 := m.NewFunction("f1", voidFn)
	f2 := m.NewFunction("f2", voidFn)
	b1 := NewBuilderAt(f1.NewBlock("entry"))
	b1.RetVoid()
	b2 := NewBuilderAt(f2.NewBlock("entry"))
	b2.RetVoid()

	blkOfF2, _ := f2.EntryBlock()
	if err := errs.Catch(func() { BlockAddress(f1, blkOfF2) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("BlockAddress with foreign block should panic ErrInvalidArg, got %v", err)
	}
}

// TestCalledFunctionIndirect 间接调用的 CalledFunction 必须返回 false。
func TestCalledFunctionIndirect(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "indirect")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{ctx.Ptr(0)}, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()
	callee := fn.ParamAs[llvm.PtrT](0)
	call := b.CallIndirect[llvm.IntT](callee, ctx.Fn(i32, nil, false), nil, "c")
	if _, ok := call.CalledFunction(); ok {
		t.Fatalf("indirect call CalledFunction should return false")
	}
	b.Ret(call)
}
