package ir

import (
	"testing"

	"github.com/kkkunny/go-llvm"
)

// TestValueHelperFloorChecks 崩溃类地板：nil/已释放值进入包级辅助 API 必须 panic ErrInvalidArg，
// 而不是把悬垂句柄交给 LLVM。两种构建模式都必须成立。
func TestValueHelperFloorChecks(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "floor")
	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()

	a := fn.ParamAs[llvm.IntT](0)
	add := b.Add(a, a, "x")
	b.Ret(add)
	// 关闭模块使全部值句柄失效
	if err := m.Close(); err != nil {
		t.Fatal(err)
	}

	dead := add
	cases := map[string]func(){
		"OperandCount": func() { OperandCount(dead) },
		"OperandAt":    func() { OperandAt(dead, 0) },
		"SetOperand":   func() { SetOperand(dead, 0, dead) },
		"Operands": func() {
			for range Operands(dead) {
			}
		},
		"Uses": func() {
			for range Uses(dead) {
			}
		},
		"ReplaceAllUsesOld": func() {
			ReplaceAllUses(dead, dead)
		},
		"OpOf":            func() { OpOf(dead) },
		"CanFastMath":     func() { CanFastMath(dead) },
		"FastMathOf":      func() { FastMathOf(dead) },
		"SetFastMath":     func() { SetFastMath(dead, FastMathNone) },
		"GEPNoWrapOf":     func() { GEPNoWrapOf(dead) },
		"SetGEPNoWrap":    func() { SetGEPNoWrap(dead, NoWrapNone) },
		"SetNSW":          func() { SetNSW(dead, true) },
		"SetNUW":          func() { SetNUW(dead, true) },
		"SetExact":        func() { SetExact(dead, true) },
		"SetNNeg":         func() { SetNNeg(dead, true) },
		"SetInBounds":     func() { SetInBounds(dead, true) },
		"SetTailCall":     func() { SetTailCall(dead, true) },
		"IsTailCall":      func() { IsTailCall(dead) },
		"SetTailCallKind": func() { SetTailCallKind(dead, TailCallNone) },
		"SetParamAlign":   func() { SetParamAlign(dead, 0, 4) },
		"SyncScopeOf":     func() { SyncScopeOf(dead) },
		"SetSyncScope":    func() { SetSyncScope(dead, 0) },
		"SuccessorCount":  func() { SuccessorCount(dead) },
		"Successor":       func() { Successor(dead, 0) },
		"SetSuccessor":    func() { SetSuccessor(dead, 0, blk) },
		"IsConditional":   func() { IsConditional(dead) },
		"Condition":       func() { Condition(dead) },
		"SetCondition":    func() { SetCondition(dead, dead) },
		"AttachMetadata":  func() { AttachMetadata(dead, "k", llvm.Metadata{}) },
		"InstMetadata":    func() { InstMetadata(dead, "k") },
		"BlockAddressFunction": func() {
			BlockAddressFunction(dead)
		},
		"BlockAddressBlock": func() { BlockAddressBlock(dead) },
	}
	for name, call := range cases {
		t.Run(name, func(t *testing.T) {
			err := llvm.Catch(call)
			if err == nil || err.Reason != llvm.ErrInvalidArg {
				t.Fatalf("want ErrInvalidArg, got %v", err)
			}
		})
	}

	// 容器句柄（Function/Block）走 Ref 地板：模块关闭后为 ErrUseAfterFree
	if err := llvm.Catch(func() { BlockAddress(fn, blk) }); err == nil || err.Reason != llvm.ErrUseAfterFree {
		t.Fatalf("BlockAddress on dead function should panic ErrUseAfterFree, got %v", err)
	}

	// nil 接口同样必须被地板拦下
	nilCases := map[string]func(){
		"OperandCount": func() { OperandCount(nil) },
		"OperandAt":    func() { OperandAt(nil, 0) },
		"SetOperand":   func() { SetOperand(nil, 0, nil) },
		"Uses": func() {
			for range Uses(nil) {
			}
		},
		"ReplaceAllUses": func() { ReplaceAllUses(nil, nil) },
		"OpOf":           func() { OpOf(nil) },
		"CanFastMath":    func() { CanFastMath(nil) },
		"SuccessorCount": func() { SuccessorCount(nil) },
		"IsConditional":  func() { IsConditional(nil) },
		"Condition":      func() { Condition(nil) },
		"SetCondition":   func() { SetCondition(nil, nil) },
	}
	for name, call := range nilCases {
		t.Run("nil/"+name, func(t *testing.T) {
			err := llvm.Catch(call)
			if err == nil || err.Reason != llvm.ErrInvalidArg {
				t.Fatalf("want ErrInvalidArg, got %v", err)
			}
		})
	}
}

// TestValueHelperCrossContextChecks 跨 Context 的替换/操作数/条件必须 panic ErrCrossContext。
func TestValueHelperCrossContextChecks(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "ctx1")
	defer m.Close()
	ctx2 := llvm.NewContext()
	defer ctx2.Close()
	m2 := NewModule(ctx2, "ctx2")
	defer m2.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	blk := fn.NewBlock("entry")
	b := NewBuilderAt(blk)
	defer b.Close()
	add := b.Add(fn.ParamAs[llvm.IntT](0), fn.ParamAs[llvm.IntT](0), "x")
	br := b.CondBr(b.ICmp(llvm.IntEQ, add, add, "c"), blk, blk)

	i32b := ctx2.Int(32)
	foreign := ctx2.ConstInt(i32b, 1)
	fn2 := m2.NewFunction("f2", ctx2.Fn(i32b, nil, false))
	blk2 := fn2.NewBlock("entry")

	cases := map[string]func(){
		"SetOperand":     func() { SetOperand(add, 0, foreign) },
		"ReplaceAllUses": func() { ReplaceAllUses(add, foreign) },
		"SetSuccessor":   func() { SetSuccessor(br, 0, blk2) },
		"SetCondition":   func() { SetCondition(br, foreign) },
	}
	for name, call := range cases {
		t.Run(name, func(t *testing.T) {
			err := llvm.Catch(call)
			if err == nil || err.Reason != llvm.ErrCrossContext {
				t.Fatalf("want ErrCrossContext, got %v", err)
			}
		})
	}
}
