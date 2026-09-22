package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func arithModule(t *testing.T) (*llvm.Context, *Module, *Builder, Function) {
	t.Helper()
	ctx := llvm.NewContext()
	m := NewModule(ctx, "arith")
	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32, i32}, false))
	b := NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	return ctx, m, b, fn
}

func TestBuilderIntArithGolden(t *testing.T) {
	ctx, m, b, fn := arithModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	a := fn.ParamAs[llvm.IntT](0)
	c := fn.ParamAs[llvm.IntT](1)

	b.Add(a, c, "add")
	b.AddNSW(a, c, "addnsw")
	b.AddNUW(a, c, "addnuw")
	b.Sub(a, c, "sub")
	b.Mul(a, c, "mul")
	b.SDiv(a, c, "sdiv")
	b.UDiv(a, c, "udiv")
	b.SRem(a, c, "srem")
	b.URem(a, c, "urem")
	b.Shl(a, c, "shl")
	b.LShr(a, c, "lshr")
	b.AShr(a, c, "ashr")
	b.And(a, c, "and")
	b.Or(a, c, "or")
	b.Xor(a, c, "xor")
	b.Neg(a, "neg")
	b.Not(a, "not")
	b.RetVoid()

	got := m.String()
	for _, want := range []string{
		"%add = add i32 %0, %1",
		"%addnsw = add nsw i32 %0, %1",
		"%addnuw = add nuw i32 %0, %1",
		"%sub = sub i32 %0, %1",
		"%mul = mul i32 %0, %1",
		"%sdiv = sdiv i32 %0, %1",
		"%udiv = udiv i32 %0, %1",
		"%srem = srem i32 %0, %1",
		"%urem = urem i32 %0, %1",
		"%shl = shl i32 %0, %1",
		"%lshr = lshr i32 %0, %1",
		"%ashr = ashr i32 %0, %1",
		"%and = and i32 %0, %1",
		"%or = or i32 %0, %1",
		"%xor = xor i32 %0, %1",
		"%neg = sub i32 0, %0",
		"%not = xor i32 %0, -1",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("module output missing %q:\n%s", want, got)
		}
	}
}

func TestBuilderFloatArithGolden(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "farith")
	defer m.Close()

	f64 := ctx.Float(llvm.FloatDouble)
	fn := m.NewFunction("f", ctx.Fn(f64, []llvm.AnyType{f64, f64}, false))
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(fn.NewBlock("entry"))

	a := fn.ParamAs[llvm.FloatT](0)
	c := fn.ParamAs[llvm.FloatT](1)
	b.FAdd(a, c, "fadd")
	b.FSub(a, c, "fsub")
	b.FMul(a, c, "fmul")
	b.FDiv(a, c, "fdiv")
	b.FRem(a, c, "frem")
	b.FNeg(a, "fneg")
	b.RetVoid()

	got := m.String()
	for _, want := range []string{
		"%fadd = fadd double %0, %1",
		"%fsub = fsub double %0, %1",
		"%fmul = fmul double %0, %1",
		"%fdiv = fdiv double %0, %1",
		"%frem = frem double %0, %1",
		"%fneg = fneg double %0",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("module output missing %q:\n%s", want, got)
		}
	}
}

func TestBuilderCmpAndSelect(t *testing.T) {
	ctx, m, b, fn := arithModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	a := fn.ParamAs[llvm.IntT](0)
	c := fn.ParamAs[llvm.IntT](1)

	lt := b.ICmp(llvm.IntSLT, a, c, "lt")
	eq := b.ICmp(llvm.IntEQ, a, c, "eq")
	sel := b.Select(lt, a, c, "sel")
	b.Ret(sel)
	_ = eq

	got := m.String()
	for _, want := range []string{
		"%lt = icmp slt i32 %0, %1",
		"%eq = icmp eq i32 %0, %1",
		"%sel = select i1 %lt, i32 %0, i32 %1",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("module output missing %q:\n%s", want, got)
		}
	}

	f64 := ctx.Float(llvm.FloatDouble)
	ffn := m.NewFunction("g", ctx.Fn(ctx.Bool(), []llvm.AnyType{f64, f64}, false))
	b.MoveToEnd(ffn.NewBlock("entry"))
	fc := b.FCmp(llvm.FloatOLT, ffn.ParamAs[llvm.FloatT](0), ffn.ParamAs[llvm.FloatT](1), "flt")
	b.Ret(fc)
	if got := m.String(); !strings.Contains(got, "%flt = fcmp olt double %0, %1") {
		t.Fatalf("module output missing fcmp:\n%s", got)
	}
}

func TestBuilderTypeMismatch(t *testing.T) {
	ctx, m, b, fn := arithModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := fn.ParamAs[llvm.IntT](0)
	i64 := ctx.ConstInt(ctx.Int(64), 1, false).Value

	err := llvm.Catch(func() { b.Add(i32, i64, "bad") })
	if err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("mixed widths should panic ErrTypeMismatch, got %v", err)
	}

	err = llvm.Catch(func() { b.Select(ctx.ConstInt(ctx.Int(32), 1, false), i32, i32, "bad") })
	if err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("non-i1 select condition should panic, got %v", err)
	}
}
