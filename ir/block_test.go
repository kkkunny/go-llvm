package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

// TestBlockStructure 覆盖基本块改名、首末指令、前后块遍历与空块判定。
func TestBlockStructure(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "blocks")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	entry := fn.NewBlock("entry")
	mid := fn.NewBlock("mid")
	exit := fn.NewBlock("exit")
	empty := fn.NewBlock("empty")

	b := NewBuilderAt(entry)
	defer b.Close()
	a := fn.ParamAs[llvm.IntT](0)
	b.Br(mid)

	b.MoveToEnd(mid)
	x := b.Add(a, a, "x")
	b.Br(exit)

	b.MoveToEnd(exit)
	b.Ret(x)

	// SetName 后 Name 必须同步；改名的块仍可正常遍历
	mid.SetName("renamed")
	if got := mid.Name(); got != "renamed" {
		t.Fatalf("block name after SetName = %q, want renamed", got)
	}
	if got := mid.Belong(); got.Name() != fn.Name() {
		t.Fatalf("block belong = %s, want %s", got.Name(), fn.Name())
	}

	// FirstInst / LastInst
	first, ok := mid.FirstInst()
	if !ok || !strings.Contains(first.String(), "add i32 %0, %0") {
		t.Fatalf("FirstInst = %v %v", first, ok)
	}
	last, ok := mid.LastInst()
	if !ok || !strings.Contains(last.String(), "br label %exit") {
		t.Fatalf("LastInst = %v %v", last, ok)
	}
	if count := len(mid.Insts()); count != 2 {
		t.Fatalf("mid instructions = %d, want 2", count)
	}

	// 空块：FirstInst/LastInst 返回 false，Empty 为真
	if _, ok := empty.FirstInst(); ok {
		t.Fatal("empty block should have no first instruction")
	}
	if _, ok := empty.LastInst(); ok {
		t.Fatal("empty block should have no last instruction")
	}
	if !empty.Empty() {
		t.Fatal("empty block should report Empty")
	}
	if entry.Empty() {
		t.Fatal("entry block should not be empty")
	}

	// Next/Prev 链：entry -> renamed -> exit -> empty
	next, ok := entry.Next()
	if !ok || next.Name() != "renamed" {
		t.Fatalf("entry.Next = %s %v", next.Name(), ok)
	}
	prev, ok := exit.Prev()
	if !ok || prev.Name() != "renamed" {
		t.Fatalf("exit.Prev = %s %v", prev.Name(), ok)
	}
	if back, ok := next.Prev(); !ok || back.Name() != "entry" {
		t.Fatalf("renamed.Prev = %s %v", back.Name(), ok)
	}
	if _, ok := entry.Prev(); ok {
		t.Fatal("first block should have no previous block")
	}
	if _, ok := empty.Next(); ok {
		t.Fatal("last block should have no next block")
	}
	// 给空块补终结指令，使模块可校验（空块断言已在上方完成）
	b.MoveToEnd(empty)
	b.Unreachable()

	// 打印文本使用改名后的标签
	out := m.String()
	if !strings.Contains(out, "br label %renamed") || !strings.Contains(out, "renamed:") {
		t.Fatalf("IR missing renamed block:\n%s", out)
	}
	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
}

// TestBlockTerminator 覆盖终结指令查询：Terminator/IsTerminating 必须区分
// "已有终结指令"与"仅含普通指令"（含 void 调用，LastInst 类型无法区分的情形），空块返回 false。
func TestBlockTerminator(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "terminator")
	defer m.Close()

	void := ctx.Void()
	i32 := ctx.Int(32)
	nop := m.NewFunction("nop", ctx.Fn(void, nil, false))
	fn := m.NewFunction("f", ctx.Fn(void, []llvm.AnyType{i32}, false))
	entry := fn.NewBlock("entry")
	mid := fn.NewBlock("mid")
	exit := fn.NewBlock("exit")

	b := NewBuilderAt(entry)
	defer b.Close()
	a := fn.ParamAs[llvm.IntT](0)

	// 未终结：只有 add 与 void 调用（后者类型为 void，不能靠 LastInst 类型判断）
	b.Add(a, a, "x")
	b.Call[llvm.VoidT](nop, nil, "")
	if term, ok := entry.Terminator(); ok {
		t.Fatalf("block with only non-terminator instructions has terminator %s", term.String())
	}
	if entry.IsTerminating() {
		t.Fatal("block with only non-terminator instructions should not be terminating")
	}

	// 终结后：Terminator 必须与刚插入的 br 一致，OpOf 分类正确
	br := b.Br(mid)
	term, ok := entry.Terminator()
	if !ok {
		t.Fatal("terminated block should report a terminator")
	}
	if !term.Ref().Equal(br.Ref()) {
		t.Fatalf("Terminator = %s, want br %s", term.String(), br.String())
	}
	if op, ok := OpOf(term); !ok || op != OpBr {
		t.Fatalf("terminator opcode = %v %v, want br", op, ok)
	}
	if !entry.IsTerminating() {
		t.Fatal("terminated block should report IsTerminating")
	}

	// 空块：尚未终结，Terminator 返回 false
	if term, ok := mid.Terminator(); ok {
		t.Fatalf("empty block has terminator %s", term.String())
	}
	if mid.IsTerminating() {
		t.Fatal("empty block should not be terminating")
	}

	// ret 终结的块
	b.MoveToEnd(mid)
	ret := b.RetVoid()
	term, ok = mid.Terminator()
	if !ok {
		t.Fatal("ret block should report a terminator")
	}
	if !term.Ref().Equal(ret.Ref()) {
		t.Fatalf("mid Terminator = %s, want ret", term.String())
	}
	if !mid.IsTerminating() {
		t.Fatal("ret block should report IsTerminating")
	}

	// 补上 exit 的终结指令，模块必须通过校验
	b.MoveToEnd(exit)
	b.Unreachable()
	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
}
