package llvm

// Kind 值/类型种类标记（封闭约束，不允许用户扩展）
type Kind interface{ kind() }

type (
	// VoidT 空类型
	VoidT struct{}
	// IntT 整数类型
	IntT struct{}
	// FloatT 浮点类型
	FloatT struct{}
	// PtrT 指针类型
	PtrT struct{}
	// StructT 结构体类型
	StructT struct{}
	// ArrayT 数组类型
	ArrayT struct{}
	// VecT 向量类型
	VecT struct{}
	// FnT 函数类型
	FnT struct{}
	// LabelT 基本块类型
	LabelT struct{}
	// MetaT 元数据类型
	MetaT struct{}
	// TokenT token 类型
	TokenT struct{}
	// DynT 类型擦除（未知/动态来源）
	DynT struct{}
)

func (VoidT) kind()   {}
func (IntT) kind()    {}
func (FloatT) kind()  {}
func (PtrT) kind()    {}
func (StructT) kind() {}
func (ArrayT) kind()  {}
func (VecT) kind()    {}
func (FnT) kind()     {}
func (LabelT) kind()  {}
func (MetaT) kind()   {}
func (TokenT) kind()  {}
func (DynT) kind()    {}

// kindOf 返回类型参数 T 对应的种类标记（T 的约束即标记集合，直接返回实例）
func kindOf[T Kind]() Kind {
	var t T
	return t
}
