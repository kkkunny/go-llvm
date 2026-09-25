package ir

import (
	"iter"
	"reflect"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
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
func (f Function) CountParams() uint {
	return uint(binding.LLVMCountParams(f.Ref()))
}

// Param 第 i 个参数（擦除种类）
func (f Function) Param(i uint) Param {
	const op = "ir.Function.Param"
	if checks.Debug && i >= f.CountParams() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "parameter index %d out of range", i)
	}
	ref := binding.LLVMGetParam(f.Ref(), uint32(i))
	return Param{Value: llvm.NewValue[llvm.DynT](f.Context(), f.Lifetime(), ref)}
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

// Blocks 全部基本块；需要零分配遍历时用 AllBlocks
func (f Function) Blocks() []Block {
	refs := binding.LLVMGetBasicBlocks(f.Ref())
	blocks := make([]Block, len(refs))
	for i, ref := range refs {
		blocks[i] = wrapBlock(f.Context(), f.Lifetime(), ref)
	}
	return blocks
}

// AllBlocks 惰性遍历全部基本块（range 友好，无切片分配）
func (f Function) AllBlocks() iter.Seq[Block] {
	return func(yield func(Block) bool) {
		for ref := binding.LLVMGetFirstBasicBlock(f.Ref()); !ref.IsNil(); ref = binding.LLVMGetNextBasicBlock(ref) {
			if !yield(wrapBlock(f.Context(), f.Lifetime(), ref)) {
				return
			}
		}
	}
}

// AllParams 惰性遍历全部参数（range 友好，无切片分配）
func (f Function) AllParams() iter.Seq[Param] {
	return func(yield func(Param) bool) {
		n := f.CountParams()
		for i := uint(0); i < n; i++ {
			if !yield(f.Param(i)) {
				return
			}
		}
	}
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
func (f Function) OnlyDecl() bool {
	return binding.LLVMIsDeclaration(f.Ref())
}

// Verify 校验函数
func (f Function) Verify() bool {
	return !binding.LLVMVerifyFunction(f.Ref(), binding.LLVMReturnStatusAction)
}

// Linkage 链接类型
func (f Function) Linkage() llvm.Linkage {
	return llvm.Linkage(binding.LLVMGetLinkage(f.Ref()))
}

// SetLinkage 设置链接类型
func (f Function) SetLinkage(l llvm.Linkage) {
	binding.LLVMSetLinkage(f.Ref(), binding.LLVMLinkage(l))
}

// SetPersonality 设置 personality 函数（EH 展开用）
func (f Function) SetPersonality(pers llvm.ValueRef[llvm.FnT]) {
	const op = "ir.Function.SetPersonality"
	pv := pers.AsValue()
	if pv.IsNil() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil personality")
	}
	f.Context().CheckValues(op, pv)
	binding.LLVMSetPersonalityFn(f.Ref(), pv.Ref())
}

// Personality personality 函数（未设置时返回 false）
func (f Function) Personality() (llvm.Value[llvm.FnT], bool) {
	ref := binding.LLVMGetPersonalityFn(f.Ref())
	if ref.IsNil() {
		return llvm.Value[llvm.FnT]{}, false
	}
	return llvm.NewValue[llvm.FnT](f.Context(), f.Lifetime(), ref), true
}

// Section 段名
func (f Function) Section() string {
	return binding.LLVMGetSection(f.Ref())
}

// SetSection 设置段名
func (f Function) SetSection(s string) {
	binding.LLVMSetSection(f.Ref(), s)
}

// GC GC 策略名
func (f Function) GC() string {
	return binding.LLVMGetGC(f.Ref())
}

// SetGC 设置 GC 策略名
func (f Function) SetGC(name string) {
	binding.LLVMSetGC(f.Ref(), name)
}

// PrefixData 前缀数据（无则空句柄）
func (f Function) PrefixData() llvm.Value[llvm.DynT] {
	return llvm.ValueOf(f.Context(), f.Lifetime(), binding.LLVMGetPrefixData(f.Ref()))
}

// SetPrefixData 设置前缀数据
func (f Function) SetPrefixData(v llvm.AnyValue) {
	const op = "ir.Function.SetPrefixData"
	f.Context().CheckValues(op, v)
	binding.LLVMSetPrefixData(f.Ref(), v.Ref())
}

// PrologueData 序言数据（无则空句柄）
func (f Function) PrologueData() llvm.Value[llvm.DynT] {
	return llvm.ValueOf(f.Context(), f.Lifetime(), binding.LLVMGetPrologueData(f.Ref()))
}

// SetPrologueData 设置序言数据
func (f Function) SetPrologueData(v llvm.AnyValue) {
	const op = "ir.Function.SetPrologueData"
	f.Context().CheckValues(op, v)
	binding.LLVMSetPrologueData(f.Ref(), v.Ref())
}

// Param 函数参数角色（内嵌 Value[DynT]）
type Param struct {
	llvm.Value[llvm.DynT]
}

// SetAlign 设置参数对齐
func (p Param) SetAlign(n uint32) {
	const op = "ir.Param.SetAlign"
	preAlign(op, n)
	binding.LLVMSetAlignment(p.Ref(), n)
}

// Belong 参数所属函数
func (p Param) Belong() Function {
	return Function{Value: llvm.NewValue[llvm.FnT](p.Context(), p.Lifetime(), binding.LLVMGetParamParent(p.Ref()))}
}

// GoFunc Go 签名绑定的函数句柄
type GoFunc[F any] struct {
	fn     Function
	goType reflect.Type
}

// Function 返回函数角色
func (f GoFunc[F]) Function() Function { return f.fn }

// GoType 返回绑定的 Go 函数类型
func (f GoFunc[F]) GoType() reflect.Type { return f.goType }
