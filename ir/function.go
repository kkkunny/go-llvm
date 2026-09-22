package ir

import (
	"reflect"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Function 函数角色（内嵌 Value[FnT]，自动实现 llvm.ValueRef/AnyValue）
type Function struct {
	llvm.Value[llvm.FnT]
}

// Signature 函数类型
func (f Function) Signature() llvm.FnType {
	return llvm.AsFnType(llvm.TypeOfRef(f.Context(), binding.LLVMGetFunctionType(f.Ref())))
}

// CountParams 参数个数
func (f Function) CountParams() uint { return uint(binding.LLVMCountParams(f.Ref())) }

// Param 第 i 个参数（擦除种类）
func (f Function) Param(i uint) Param {
	ref := binding.LLVMGetParam(f.Ref(), uint32(i))
	return Param{Value: wrapValue[llvm.DynT](f.Context(), f.Lifetime(), ref)}
}

// ParamAs 第 i 个参数（泛型方法；种类不符 panic）
func (f Function) ParamAs[U llvm.Kind](i uint) llvm.Value[U] {
	return f.Param(i).MustAs[U]()
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
	ref := binding.LLVMAppendBasicBlockInContext(f.Context().Ref(), f.Ref(), name)
	return wrapBlock(f.Context(), f.Lifetime(), ref)
}

// Blocks 全部基本块
func (f Function) Blocks() []Block {
	refs := binding.LLVMGetBasicBlocks(f.Ref())
	blocks := make([]Block, len(refs))
	for i, ref := range refs {
		blocks[i] = wrapBlock(f.Context(), f.Lifetime(), ref)
	}
	return blocks
}

// EntryBlock 入口块
func (f Function) EntryBlock() (Block, bool) {
	ref := binding.LLVMGetEntryBasicBlock(f.Ref())
	if ref.IsNil() {
		return Block{}, false
	}
	return wrapBlock(f.Context(), f.Lifetime(), ref), true
}

// OnlyDecl 是否只有声明（无基本块）
func (f Function) OnlyDecl() bool { return binding.LLVMIsDeclaration(f.Ref()) }

// Verify 校验函数
func (f Function) Verify() bool {
	return !binding.LLVMVerifyFunction(f.Ref(), binding.LLVMReturnStatusAction)
}

// Linkage 链接类型
func (f Function) Linkage() llvm.Linkage { return llvm.Linkage(binding.LLVMGetLinkage(f.Ref())) }

// SetLinkage 设置链接类型
func (f Function) SetLinkage(l llvm.Linkage) {
	binding.LLVMSetLinkage(f.Ref(), binding.LLVMLinkage(l))
}

// Param 函数参数角色（内嵌 Value[DynT]）
type Param struct {
	llvm.Value[llvm.DynT]
}

// SetAlign 设置参数对齐
func (p Param) SetAlign(n uint32) { binding.LLVMSetAlignment(p.Ref(), n) }

// Belong 参数所属函数
func (p Param) Belong() Function {
	return Function{Value: wrapValue[llvm.FnT](p.Context(), p.Lifetime(), binding.LLVMGetParamParent(p.Ref()))}
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
