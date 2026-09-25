package ir

import (
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestIntrinsic(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "t")
	defer m.Close()

	ins, ok := FindIntrinsic("llvm.trap")
	if !ok {
		t.Fatalf("llvm.trap not found")
	}
	if ins.IsOverloaded() {
		t.Fatalf("llvm.trap should not be overloaded")
	}
	if ins.Name() != "llvm.trap" {
		t.Fatalf("Name = %s", ins.Name())
	}
	fn, ok := ins.Declaration(m, nil)
	if !ok {
		t.Fatalf("declaration failed")
	}
	if fn.Name() != "llvm.trap" {
		t.Fatalf("fn name = %s", fn.Name())
	}
	if !fn.OnlyDecl() {
		t.Fatalf("declaration should be a decl")
	}

	// 重载 intrinsic：需要参数类型
	ov, ok := FindIntrinsic("llvm.sadd.with.overflow.i32")
	if !ok {
		t.Fatalf("sadd.with.overflow.i32 not found")
	}
	i32 := ctx.Int(32)
	fn2, ok := ov.Declaration(m, []llvm.AnyType{i32})
	if !ok {
		t.Fatalf("overloaded declaration failed")
	}
	if fn2.CountParams() != 2 {
		t.Fatalf("params = %d, want 2", fn2.CountParams())
	}

	// 重载但缺参数：应返回 false（LLVM 对空参数的重载查询会崩，前置拦截）
	if _, ok := ov.Declaration(m, nil); ok {
		t.Fatalf("overloaded intrinsic without params should fail")
	}

	if _, ok := FindIntrinsic("llvm.not.an.intrinsic"); ok {
		t.Fatalf("unknown intrinsic should not be found")
	}
}
