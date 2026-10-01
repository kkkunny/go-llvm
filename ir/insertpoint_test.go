package ir

import (
	"testing"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/errs"
)

// instNames 返回块内指令的名字（未命名指令为空字符串），按出现顺序。
func instNames(blk Block) []string {
	var names []string
	for inst := range blk.AllInsts() {
		names = append(names, inst.Name())
	}
	return names
}

func sameNames(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

// TestInsertPointRestoreAtEnd 保存"块末尾"插入点、嵌套生成后恢复，新指令落回原块末尾。
func TestInsertPointRestoreAtEnd(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "ip-end")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()

	slot := b.Alloca(i32, "slot")
	b.Load(slot, i32, "a")
	p := b.SaveInsertPoint()

	// 模拟嵌套函数生成：切到另一个函数再恢复
	other := m.NewFunction("g", ctx.Fn(i32, nil, false))
	b.MoveToEnd(other.NewBlock("entry"))
	b.Ret(ctx.ConstInt(i32, 0))
	b.RestoreInsertPoint(p)

	b.Load(slot, i32, "c")
	b.RetVoid()

	if got := instNames(blk); !sameNames(got, []string{"slot", "a", "c", ""}) {
		t.Fatalf("entry insts = %v, want [slot a c ret-void]\n%s", got, m.String())
	}
	if got := instNames(other.Blocks()[0]); !sameNames(got, []string{""}) {
		t.Fatalf("nested function should keep its own instructions, got %v", got)
	}
}

// TestInsertPointRestoreBeforeInst 保存"某指令之前"插入点，恢复后新指令插在该指令之前。
func TestInsertPointRestoreBeforeInst(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "ip-before")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()

	slot := b.Alloca(i32, "slot")
	a := b.Load(slot, i32, "a")
	b.MoveBefore(a)
	p := b.SaveInsertPoint()

	other := m.NewFunction("g", ctx.Fn(i32, nil, false))
	b.MoveToEnd(other.NewBlock("entry"))
	b.Ret(ctx.ConstInt(i32, 0))
	b.RestoreInsertPoint(p)

	b.Load(slot, i32, "c")
	b.MoveToEnd(blk)
	b.RetVoid()

	if got := instNames(blk); !sameNames(got, []string{"slot", "c", "a", ""}) {
		t.Fatalf("entry insts = %v, want [slot c a ret-void]\n%s", got, m.String())
	}
}

// TestInsertPointUnpositioned 未定位的构建器可保存/恢复"未定位"状态。
func TestInsertPointUnpositioned(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "ip-none")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, nil, false))
	blk := fn.NewBlock("entry")

	b := NewBuilder(ctx)
	defer b.Close()

	p := b.SaveInsertPoint()
	b.MoveToEnd(blk)
	if _, ok := b.CurrentBlock(); !ok {
		t.Fatal("builder should be positioned")
	}
	b.RestoreInsertPoint(p)
	if _, ok := b.CurrentBlock(); ok {
		t.Fatal("restored unpositioned point should leave builder unpositioned")
	}
	if err := errs.Catch(func() { b.RetVoid() }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("unpositioned builder should panic ErrInvalidArg, got %v", err)
	}
}

// TestInsertPointValidation 快照中的指令已不在其块内时 panic，
// 而不是把非法位置交给 LLVM。
func TestInsertPointValidation(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "ip-valid")
	defer m.Close()
	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()
	slot := b.Alloca(i32, "slot")
	a := b.Load(slot, i32, "a")
	b.MoveBefore(a)
	p := b.SaveInsertPoint()

	// 把快照的块改成另一个块：指令不再属于该块，恢复必须被拦截
	blk2 := fn.NewBlock("second")
	p.blk = blk2
	if err := errs.Catch(func() { b.RestoreInsertPoint(p) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("instruction outside saved block should panic ErrInvalidArg, got %v", err)
	}
}
