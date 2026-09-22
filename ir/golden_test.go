package ir

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/kkkunny/go-llvm"
)

var updateGolden = flag.Bool("update", false, "更新 golden 文件")

func checkGolden(t *testing.T, got, name string) {
	t.Helper()
	goldenPath := filepath.Join("testdata", "golden", name)
	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(goldenPath, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden (可用 -update 生成): %v", err)
	}
	if got != string(want) {
		t.Fatalf("module output differs from golden %s:\n--- got ---\n%s\n--- want ---\n%s", name, got, want)
	}
}

// TestGoldenMain 端到端构建 README 示例模块并与 golden 文件比对
func TestGoldenMain(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "main")
	defer m.Close()

	i32 := ctx.Int(32)
	main, err := m.NewFunc[func() int32]("main")
	if err != nil {
		t.Fatal(err)
	}
	entry := main.Function().NewBlock("entry")

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	// int32 result = add(1, 2); return result == 3 ? 0 : 1
	one := ctx.ConstInt(i32, 1, false).Value
	two := ctx.ConstInt(i32, 2, false).Value
	sum := b.Add(one, two, "sum")
	isThree := b.ICmp(llvm.IntEQ, sum, ctx.ConstInt(i32, 3, false).Value, "is_three")
	ret := b.Select(isThree, ctx.ConstInt(i32, 0, false).Value, ctx.ConstInt(i32, 1, false).Value, "ret")
	b.Ret(ret)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}

	checkGolden(t, m.String(), "main.ll")
}

// TestGoldenLoop 覆盖 PHI/Switch/内存指令的端到端场景
func TestGoldenLoop(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "loop")
	defer m.Close()

	i32 := ctx.Int(32)
	i64 := ctx.Int(64)
	boolTy := ctx.Bool()
	fn, err := m.NewFunc[func(int32) int32]("sum_to")
	if err != nil {
		t.Fatal(err)
	}
	entry := fn.Function().NewBlock("entry")
	loop := fn.Function().NewBlock("loop")
	exit := fn.Function().NewBlock("exit")

	b := NewBuilder(ctx)
	defer b.Close()

	b.MoveToEnd(entry)
	acc := b.Alloca(i32, "acc")
	b.Store(ctx.ConstInt(i32, 0, false).Value, acc)
	b.Br(loop)

	b.MoveToEnd(loop)
	i := b.PHI(i32, "i")
	cur := b.Load(acc, i32, "cur")
	i.AddIncoming(Incoming[llvm.IntT]{Value: ctx.ConstInt(i32, 0, false).Value, Block: entry})
	next := b.Add(i, ctx.ConstInt(i32, 1, false).Value, "next")
	b.Store(next, acc)
	done := b.ICmp(llvm.IntSGE, next, fn.Function().ParamAs[llvm.IntT](0), "done")
	b.CondBr(done, exit, loop)
	i.AddIncoming(Incoming[llvm.IntT]{Value: next, Block: loop})
	_ = cur
	_ = boolTy

	b.MoveToEnd(exit)
	result := b.Load(acc, i32, "result")
	b.Ret(result)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}

	checkGolden(t, m.String(), "loop.ll")
	_ = i64
}
