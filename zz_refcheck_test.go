package llvm

var (
	_ TypeRef[IntT]  = Type[IntT]{}
	_ TypeRef[IntT]  = IntType{}
	_ TypeRef[DynT]  = Type[DynT]{}
	_ ValueRef[IntT] = Value[IntT]{}
	_ ValueRef[IntT] = IntConst{}
	_ ValueRef[DynT] = Value[DynT]{}
	_ AnyType        = IntType{}
	_ AnyValue       = IntConst{}
)
