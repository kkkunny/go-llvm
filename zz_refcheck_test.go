package llvm

var (
	_ TypeRef[IntT]    = Type[IntT]{}
	_ TypeRef[IntT]    = IntType{}
	_ TypeRef[FloatT]  = FloatType{}
	_ TypeRef[PtrT]    = PtrType{}
	_ TypeRef[StructT] = StructType{}
	_ TypeRef[ArrayT]  = ArrayType{}
	_ TypeRef[VecT]    = VecType{}
	_ TypeRef[FnT]     = FnType{}
	_ TypeRef[VoidT]   = VoidType{}
	_ TypeRef[DynT]    = Type[DynT]{}
	_ ValueRef[IntT]   = Value[IntT]{}
	_ ValueRef[IntT]   = IntConst{}
	_ ValueRef[FloatT] = FloatConst{}
	_ ValueRef[DynT]   = Value[DynT]{}
	_ AnyType          = IntType{}
	_ AnyType          = StructType{}
	_ AnyType          = FnType{}
	_ AnyType          = Type[DynT]{}
	_ AnyValue         = IntConst{}
	_ AnyValue         = FloatConst{}
	_ AnyValue         = Value[DynT]{}
)
