package ir

import (
	"github.com/kkkunny/go-llvm"
)

var (
	_ llvm.ValueRef[llvm.PtrT]    = Alloca{}
	_ llvm.ValueRef[llvm.IntT]    = Load[llvm.IntT]{}
	_ llvm.ValueRef[llvm.VoidT]   = Store{}
	_ llvm.ValueRef[llvm.IntT]    = Call[llvm.IntT]{}
	_ llvm.ValueRef[llvm.IntT]    = Phi[llvm.IntT]{}
	_ llvm.ValueRef[llvm.IntT]    = Switch{}
	_ llvm.ValueRef[llvm.PtrT]    = Global{}
	_ llvm.ValueRef[llvm.FnT]     = Function{}
	_ llvm.ValueRef[llvm.DynT]    = Param{}
	_ llvm.ValueRef[llvm.IntT]    = Invoke[llvm.IntT]{}
	_ llvm.ValueRef[llvm.StructT] = LandingPad[llvm.StructT]{}
	_ llvm.AnyValue               = Alloca{}
	_ llvm.AnyValue               = Load[llvm.IntT]{}
	_ llvm.AnyValue               = Store{}
	_ llvm.AnyValue               = Call[llvm.IntT]{}
	_ llvm.AnyValue               = Phi[llvm.IntT]{}
	_ llvm.AnyValue               = Switch{}
	_ llvm.AnyValue               = Global{}
	_ llvm.AnyValue               = Function{}
	_ llvm.AnyValue               = Param{}
	_ llvm.AnyValue               = Invoke[llvm.IntT]{}
	_ llvm.AnyValue               = LandingPad[llvm.StructT]{}
)
