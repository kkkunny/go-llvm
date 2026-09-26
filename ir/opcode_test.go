package ir

import (
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestOpOf(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "t")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()

	a := fn.ParamAs[llvm.IntT](0)
	c := fn.ParamAs[llvm.IntT](1)
	add := b.Add(a, c, "x")
	b.Ret(add)

	// 非指令值（函数参数）：返回 false
	if _, ok := OpOf(a); ok {
		t.Fatalf("param should not be an instruction")
	}

	got := map[string]bool{}
	for inst := range blk.AllInsts() {
		op, ok := OpOf(inst)
		if !ok {
			t.Fatalf("inst %v is not recognized as instruction", inst)
		}
		got[op.String()] = true
	}
	if !got["add"] || !got["ret"] {
		t.Fatalf("opcodes = %v, want add & ret", got)
	}

	// 未知操作码退化为 op<N>，不 panic
	if got := Op(9999).String(); got != "op9999" {
		t.Fatalf("unknown opcode String = %q, want op9999", got)
	}
}

// TestOpIsTerminator 覆盖终结操作码分类：静态集合必须与 LLVM 的
// LLVMIsATerminatorInst 判定一致（后者经 SuccessorCount 的常开拒绝语义暴露）。
func TestOpIsTerminator(t *testing.T) {
	for _, op := range []Op{
		OpRet, OpBr, OpSwitch, OpIndirectBr, OpInvoke, OpUnreachable,
		OpResume, OpCleanupRet, OpCatchRet, OpCatchSwitch, OpCallBr,
	} {
		if !op.IsTerminator() {
			t.Fatalf("%v should be a terminator", op)
		}
	}
	for _, op := range []Op{
		OpAdd, OpLoad, OpStore, OpCall, OpPHI, OpFence, OpCatchPad,
		OpCleanupPad, OpLandingPad, OpFreeze, OpUserOp1,
	} {
		if op.IsTerminator() {
			t.Fatalf("%v should not be a terminator", op)
		}
	}

	// 交叉验证：真实指令上 Op.IsTerminator 与 LLVM 判定一致
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "terminator-ops")
	defer m.Close()

	void := ctx.Void()
	i32 := ctx.Int(32)
	nop := m.NewFunction("nop", ctx.Fn(void, nil, false))
	fn := m.NewFunction("f", ctx.Fn(void, []llvm.AnyType{i32}, false))
	entry := fn.NewBlock("entry")
	mid := fn.NewBlock("mid")
	sw := fn.NewBlock("sw")
	exit := fn.NewBlock("exit")

	b := NewBuilderAt(entry)
	defer b.Close()
	a := fn.ParamAs[llvm.IntT](0)

	check := func(name string, inst llvm.AnyValue, want bool) {
		t.Helper()
		op, ok := OpOf(inst)
		if !ok {
			t.Fatalf("%s is not recognized as an instruction", name)
		}
		if got := op.IsTerminator(); got != want {
			t.Fatalf("%s IsTerminator = %v, want %v", name, got, want)
		}
		// SuccessorCount 对非终结指令 panic，以 LLVM 的判定作真值交叉验证
		llvmTerm := llvm.Catch(func() { SuccessorCount(inst) }) == nil
		if llvmTerm != want {
			t.Fatalf("%s LLVM terminator = %v, want %v", name, llvmTerm, want)
		}
	}

	check("add", b.Add(a, a, "x"), false)
	check("call", b.Call[llvm.VoidT](nop, nil, ""), false)
	check("br", b.Br(mid), true)

	b.MoveToEnd(mid)
	check("ret", b.RetVoid(), true)

	b.MoveToEnd(sw)
	check("switch", b.Switch(a, exit), true)

	b.MoveToEnd(exit)
	check("unreachable", b.Unreachable(), true)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
}
