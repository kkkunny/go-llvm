package jit

import (
	"testing"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/checks"
	"github.com/kkkunny/go-llvm/ir"
	"github.com/kkkunny/go-llvm/target"
)

// requireDebugJIT 语义类（签名核对）负向测试仅在调试构建下运行；
// -tags=llvm_release 时 B3 核对被编译消除，此类测试跳过。
func requireDebugJIT(t *testing.T) {
	t.Helper()
	if !checks.Debug {
		t.Skip("JIT signature checks are compiled out in llvm_release builds")
	}
}

// sigModule 定义恒等函数 f(i32) -> i32 与声明 g(i32) -> i32，并交给 JIT
func sigModule(t *testing.T) (*LLJIT, func()) {
	t.Helper()
	target.InitNative()
	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "sigcheck")
	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(fn.ParamAs[llvm.IntT](0))
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	j, err := NewLLJIT()
	if err != nil {
		t.Fatal(err)
	}
	if err := j.AddIRModule(m); err != nil {
		t.Fatal(err)
	}
	return j, func() { _ = j.Close() }
}

func TestFuncSignatureCheck(t *testing.T) {
	requireDebugJIT(t)
	j, cleanup := sigModule(t)
	defer cleanup()

	// 匹配签名正常可用
	f, err := j.Func[func(int32) int32]("f")
	if err != nil {
		t.Fatal(err)
	}
	if got := f(7); got != 7 {
		t.Fatalf("f(7) = %d, want 7", got)
	}

	// 不匹配签名在注册期 panic ErrTypeMismatch，而非调用期 UB
	if err := llvm.Catch(func() { _, _ = j.Func[func(int64) int32]("f") }); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("mismatched signature should panic ErrTypeMismatch, got %v", err)
	}

	// 未记录符号（由 MapSymbol 引入）跳过核对，交由 Lookup 报错
	if _, err := j.Func[func(int64) int32]("no_such_symbol"); err == nil {
		t.Fatal("unknown symbol should fail")
	} else if e, ok := err.(*llvm.Error); !ok || e.Reason != llvm.ErrNotFound {
		t.Fatalf("unknown symbol should fail with ErrNotFound, got %v", err)
	}
}

func TestMapFuncSignatureCheck(t *testing.T) {
	requireDebugJIT(t)
	j, cleanup := sigModule(t)
	defer cleanup()

	// 与声明 g(i32)->i32 同签名：允许（包装体把声明升级为定义）
	if err := j.MapFunc[func(int32) int32]("g", func(x int32) int32 { return x + 1 }); err != nil {
		t.Fatalf("matching MapFunc should succeed, got %v", err)
	}

	// 不同签名：注册期 panic ErrTypeMismatch
	if err := llvm.Catch(func() {
		_ = j.MapFunc[func(float64) float64]("g", func(x float64) float64 { return x })
	}); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("mismatched MapFunc should panic ErrTypeMismatch, got %v", err)
	}
}
