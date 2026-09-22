package ir

import (
	"reflect"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Function 函数角色
type Function struct {
	v   llvm.Value[llvm.FnT]
	mod *Module
}

// Value 返回底层泛型值
func (f Function) Value() llvm.Value[llvm.FnT] { return f.v }

// AsValue 实现 llvm.ValueRef[FnT]
func (f Function) AsValue() llvm.Value[llvm.FnT] { return f.v }

// Name 函数名
func (f Function) Name() string { return f.v.Name() }

// SetName 设置函数名
func (f Function) SetName(name string) { f.v.SetName(name) }

// Signature 函数类型
func (f Function) Signature() llvm.FnType {
	return llvm.AsFnType(llvm.TypeOfRef(f.v.Context(), binding.LLVMGetFunctionType(f.v.Ref())))
}

// CountParams 参数个数
func (f Function) CountParams() uint { return uint(binding.LLVMCountParams(f.v.Ref())) }

// Param 第 i 个参数（擦除种类）
func (f Function) Param(i uint) Param {
	ref := binding.LLVMGetParam(f.v.Ref(), uint32(i))
	return Param{v: llvm.NewValue[llvm.DynT](f.v.Context(), f.mod.life, ref)}
}

// ParamAs 第 i 个参数（泛型方法；种类不符 panic）
func (f Function) ParamAs[U llvm.Kind](i uint) llvm.Value[U] {
	return f.Param(i).Value().MustAs[U]()
}

// Params 全部参数
func (f Function) Params() []Param {
	n := f.CountParams()
	params := make([]Param, n)
	for i := range params {
		params[i] = f.Param(uint(i))
	}
	return params
}

// NewBlock 追加基本块
func (f Function) NewBlock(name string) Block {
	ref := binding.LLVMAppendBasicBlockInContext(f.v.Context().Ref(), f.v.Ref(), name)
	return Block{ref: ref, ctx: f.v.Context(), life: f.mod.life}
}

// Blocks 全部基本块
func (f Function) Blocks() []Block {
	refs := binding.LLVMGetBasicBlocks(f.v.Ref())
	blocks := make([]Block, len(refs))
	for i, ref := range refs {
		blocks[i] = Block{ref: ref, ctx: f.v.Context(), life: f.mod.life}
	}
	return blocks
}

// EntryBlock 入口块
func (f Function) EntryBlock() (Block, bool) {
	ref := binding.LLVMGetEntryBasicBlock(f.v.Ref())
	if ref.IsNil() {
		return Block{}, false
	}
	return Block{ref: ref, ctx: f.v.Context(), life: f.mod.life}, true
}

// OnlyDecl 是否只有声明（无基本块）
func (f Function) OnlyDecl() bool { return binding.LLVMIsDeclaration(f.v.Ref()) }

// Verify 校验函数
func (f Function) Verify() bool {
	return !binding.LLVMVerifyFunction(f.v.Ref(), binding.LLVMReturnStatusAction)
}

// Linkage 链接类型
func (f Function) Linkage() llvm.Linkage { return llvm.Linkage(binding.LLVMGetLinkage(f.v.Ref())) }

// SetLinkage 设置链接类型
func (f Function) SetLinkage(l llvm.Linkage) {
	binding.LLVMSetLinkage(f.v.Ref(), binding.LLVMLinkage(l))
}

// Param 函数参数角色
type Param struct {
	v llvm.Value[llvm.DynT]
}

// Value 返回底层擦除值
func (p Param) Value() llvm.Value[llvm.DynT] { return p.v }

// Type 参数类型
func (p Param) Type() llvm.AnyType { return p.v.Type() }

// SetAlign 设置参数对齐
func (p Param) SetAlign(n uint32) { binding.LLVMSetAlignment(p.v.Ref(), n) }

// Belong 参数所属函数
func (p Param) Belong() Function {
	return Function{v: llvm.NewValue[llvm.FnT](p.v.Context(), p.v.Lifetime(), binding.LLVMGetParamParent(p.v.Ref()))}
}

// Func Go 签名绑定的函数句柄
type Func[F any] struct {
	fn     Function
	goType reflect.Type
}

// Function 返回函数角色
func (f Func[F]) Function() Function { return f.fn }

// GoType 返回绑定的 Go 函数类型
func (f Func[F]) GoType() reflect.Type { return f.goType }
