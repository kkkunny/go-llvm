package llvm

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/kkkunny/go-llvm/internal/binding"
)

type Instruction interface {
	binding() binding.LLVMValueRef
	Belong() Block
	Next() (Instruction, bool)
	Prev() (Instruction, bool)
	RemoveFromBlock()
}

type genericInst[T any] binding.LLVMValueRef

func newGenericInst[T any](ref binding.LLVMValueRef) genericInst[T] {
	return genericInst[T](ref)
}

func (i genericInst[T]) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i genericInst[T]) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (i genericInst[T]) Next() (Instruction, bool) {
	ref := binding.LLVMGetNextInstruction(i.binding())
	if ref.IsNil() {
		return nil, false
	}
	return lookupInstruction(ref), true
}

func (i genericInst[T]) Prev() (Instruction, bool) {
	ref := binding.LLVMGetPreviousInstruction(i.binding())
	if ref.IsNil() {
		return nil, false
	}
	return lookupInstruction(ref), true
}

func (i genericInst[T]) RemoveFromBlock() {
	binding.LLVMInstructionRemoveFromParent(i.binding())
}

func lookupInstruction(ref binding.LLVMValueRef) Instruction {
	if ref.IsNil() {
		return nil
	}

	switch binding.LLVMGetInstructionOpcode(ref) {
	case binding.LLVMRet, binding.LLVMBr, binding.LLVMUnreachable, binding.LLVMSwitch:
		return lookupTerminator(ref)
	case binding.LLVMAdd, binding.LLVMFAdd:
		return newValueInst[_Add](ref)
	case binding.LLVMSub, binding.LLVMFSub:
		return newValueInst[_Sub](ref)
	case binding.LLVMMul, binding.LLVMFMul:
		return newValueInst[_Mul](ref)
	case binding.LLVMSDiv, binding.LLVMUDiv, binding.LLVMFDiv:
		return newValueInst[_Div](ref)
	case binding.LLVMSRem, binding.LLVMURem, binding.LLVMFRem:
		return newValueInst[_Rem](ref)
	case binding.LLVMShl:
		return newValueInst[_Shl](ref)
	case binding.LLVMLShr, binding.LLVMAShr:
		return newValueInst[_Shr](ref)
	case binding.LLVMAnd:
		return newValueInst[_And](ref)
	case binding.LLVMOr:
		return newValueInst[_Or](ref)
	case binding.LLVMXor:
		return newValueInst[_Xor](ref)
	case binding.LLVMFNeg:
		return newValueInst[_Neg](ref)
	case binding.LLVMAlloca:
		return newValueInstWithAlign[_Alloca](ref)
	case binding.LLVMLoad:
		return newValueInstWithAlign[_Load](ref)
	case binding.LLVMStore:
		return newGenericInstWithAlign[_Store](ref)
	case binding.LLVMGetElementPtr:
		return newValueInst[_GetElementPtr](ref)
	case binding.LLVMTrunc, binding.LLVMFPTrunc:
		return newValueInst[_Trunc](ref)
	case binding.LLVMZExt, binding.LLVMSExt, binding.LLVMFPExt:
		return newValueInst[_Expand](ref)
	case binding.LLVMFPToUI, binding.LLVMFPToSI:
		return newValueInst[_FloatToInt](ref)
	case binding.LLVMUIToFP, binding.LLVMSIToFP:
		return newValueInst[_IntToFloat](ref)
	case binding.LLVMPtrToInt:
		return newValueInst[_PtrToInt](ref)
	case binding.LLVMIntToPtr:
		return newValueInst[_IntToPtr](ref)
	case binding.LLVMBitCast:
		return newValueInst[_BitCast](ref)
	case binding.LLVMICmp:
		return newIntCmpInst(ref)
	case binding.LLVMFCmp:
		return newFloatCmpInst(ref)
	case binding.LLVMPHI:
		return newPhiInst(ref)
	case binding.LLVMCall:
		return newValueInst[_Call](ref)
	case binding.LLVMInvoke:
		return newValueInst[_Invoke](ref)
	case binding.LLVMExtractElement:
		return newValueInst[_ExtractElement](ref)
	case binding.LLVMExtractValue:
		return newValueInst[_ExtractValue](ref)
	case binding.LLVMResume:
		return newValueInst[_Resume](ref)
	case binding.LLVMLandingPad:
		return newValueInst[_LandingPad](ref)
	case binding.LLVMCleanupRet:
		return newValueInst[_CleanupRet](ref)
	case binding.LLVMCatchRet:
		return newValueInst[_CatchRet](ref)
	case binding.LLVMCatchPad:
		return newValueInst[_CatchPad](ref)
	case binding.LLVMCleanupPad:
		return newValueInst[_CleanupPad](ref)
	case binding.LLVMCatchSwitch:
		return newValueInst[_CatchSwitch](ref)
	case binding.LLVMSelect:
		return newValueInst[_Select](ref)
	default:
		panic(fmt.Errorf("unknown enum value `%d`", binding.LLVMGetInstructionOpcode(ref)))
	}
}

type Terminator interface {
	Instruction
	terminator()
}

type terminatorInst[T any] struct{ genericInst[T] }

func newTerminatorInst[T any](ref binding.LLVMValueRef) terminatorInst[T] {
	return terminatorInst[T]{newGenericInst[T](ref)}
}

func (i terminatorInst[T]) terminator() {}

func lookupTerminator(ref binding.LLVMValueRef) Terminator {
	if ref.IsNil() {
		return nil
	}

	switch binding.LLVMGetInstructionOpcode(ref) {
	case binding.LLVMRet:
		return newTerminatorInst[_Return](ref)
	case binding.LLVMBr:
		return newTerminatorInst[_Br](ref)
	case binding.LLVMUnreachable:
		return newTerminatorInst[_Unreachable](ref)
	case binding.LLVMSwitch:
		return newTerminatorInst[_Switch](ref)
	default:
		panic(fmt.Errorf("unknown enum value `%d`", binding.LLVMGetInstructionOpcode(ref)))
	}
}

type valueInst[T any] struct{ genericInst[T] }

func newValueInst[T any](ref binding.LLVMValueRef) valueInst[T] {
	return valueInst[T]{newGenericInst[T](ref)}
}

func (i valueInst[T]) Type() Type {
	return lookupType(binding.LLVMTypeOf(i.binding()))
}

func (i valueInst[T]) String() string {
	return binding.LLVMPrintValueToString(i.binding())
}

type valueInstWithAlign[T any] struct{ valueInst[T] }

func newValueInstWithAlign[T any](ref binding.LLVMValueRef) valueInstWithAlign[T] {
	return valueInstWithAlign[T]{newValueInst[T](ref)}
}

func (i valueInstWithAlign[T]) SetAlign(align uint32) {
	binding.LLVMSetAlignment(i.binding(), align)
}

func (i valueInstWithAlign[T]) GetAlign() uint32 {
	return binding.LLVMGetAlignment(i.binding())
}

type genericInstWithAlign[T any] struct{ genericInst[T] }

func newGenericInstWithAlign[T any](ref binding.LLVMValueRef) genericInstWithAlign[T] {
	return genericInstWithAlign[T]{newGenericInst[T](ref)}
}

func (i genericInstWithAlign[T]) SetAlign(align uint32) {
	binding.LLVMSetAlignment(i.binding(), align)
}

func (i genericInstWithAlign[T]) GetAlign() uint32 {
	return binding.LLVMGetAlignment(i.binding())
}

type intCmpInst struct{ valueInst[_IntCmp] }

func newIntCmpInst(ref binding.LLVMValueRef) intCmpInst {
	return intCmpInst{newValueInst[_IntCmp](ref)}
}

func (i intCmpInst) GetOperator() IntPredicate {
	return IntPredicate(binding.LLVMGetICmpPredicate(i.binding()))
}

type floatCmpInst struct{ valueInst[_FloatCmp] }

func newFloatCmpInst(ref binding.LLVMValueRef) floatCmpInst {
	return floatCmpInst{newValueInst[_FloatCmp](ref)}
}

func (i floatCmpInst) GetOperator() FloatPredicate {
	return FloatPredicate(binding.LLVMGetFCmpPredicate(i.binding()))
}

type phiInst struct{ valueInst[_PHI] }

func newPhiInst(ref binding.LLVMValueRef) phiInst {
	return phiInst{newValueInst[_PHI](ref)}
}

func (i phiInst) AddIncomings(incomings ...struct {
	Value Value
	Block Block
}) {
	values := lo.Map(incomings, func(item struct {
		Value Value
		Block Block
	}, index int) binding.LLVMValueRef {
		return item.Value.binding()
	})
	blocks := lo.Map(incomings, func(item struct {
		Value Value
		Block Block
	}, index int) binding.LLVMBasicBlockRef {
		return item.Block.binding()
	})
	binding.LLVMAddIncoming(i.binding(), values, blocks)
}

func (i phiInst) CountIncomings() uint32 {
	return binding.LLVMCountIncoming(i.binding())
}

func (i phiInst) GetIncoming(index uint32) (Value, Block) {
	return lookupValue(binding.LLVMGetIncomingValue(i.binding(), index)), Block(binding.LLVMGetIncomingBlock(i.binding(), index))
}

func (i phiInst) Incomings() (res []struct {
	Value Value
	Block Block
}) {
	for j := uint32(0); j < i.CountIncomings(); j++ {
		v, b := i.GetIncoming(j)
		res = append(res, struct {
			Value Value
			Block Block
		}{Value: v, Block: b})
	}
	return res
}

type (
	_Return         struct{}
	_Br             struct{}
	_Unreachable    struct{}
	_Add            struct{}
	_Sub            struct{}
	_Mul            struct{}
	_Div            struct{}
	_Rem            struct{}
	_Shl            struct{}
	_Shr            struct{}
	_And            struct{}
	_Or             struct{}
	_Xor            struct{}
	_Neg            struct{}
	_Not            struct{}
	_Alloca         struct{}
	_Load           struct{}
	_Store          struct{}
	_GetElementPtr  struct{}
	_Trunc          struct{}
	_Expand         struct{}
	_FloatToInt     struct{}
	_IntToFloat     struct{}
	_PtrToInt       struct{}
	_IntToPtr       struct{}
	_BitCast        struct{}
	_IntCmp         struct{}
	_FloatCmp       struct{}
	_Call           struct{}
	_PHI            struct{}
	_ExtractElement struct{}
	_ExtractValue   struct{}
	_Invoke         struct{}
	_Resume         struct{}
	_LandingPad     struct{}
	_CleanupRet     struct{}
	_CatchRet       struct{}
	_CatchPad       struct{}
	_CleanupPad     struct{}
	_CatchSwitch    struct{}
	_Switch         struct{}
	_Select         struct{}
)

type (
	Return         = terminatorInst[_Return]
	Br             = terminatorInst[_Br]
	Unreachable    = terminatorInst[_Unreachable]
	Add            = valueInst[_Add]
	Sub            = valueInst[_Sub]
	Mul            = valueInst[_Mul]
	Div            = valueInst[_Div]
	Rem            = valueInst[_Rem]
	Shl            = valueInst[_Shl]
	Shr            = valueInst[_Shr]
	And            = valueInst[_And]
	Or             = valueInst[_Or]
	Xor            = valueInst[_Xor]
	Neg            = valueInst[_Neg]
	Not            = valueInst[_Not]
	Alloca         = valueInstWithAlign[_Alloca]
	Load           = valueInstWithAlign[_Load]
	Store          = genericInstWithAlign[_Store]
	GetElementPtr  = valueInst[_GetElementPtr]
	Trunc          = valueInst[_Trunc]
	Expand         = valueInst[_Expand]
	FloatToInt     = valueInst[_FloatToInt]
	IntToFloat     = valueInst[_IntToFloat]
	PtrToInt       = valueInst[_PtrToInt]
	IntToPtr       = valueInst[_IntToPtr]
	BitCast        = valueInst[_BitCast]
	IntCmp         = intCmpInst
	FloatCmp       = floatCmpInst
	Call           = valueInst[_Call]
	PHI            = phiInst
	ExtractElement = valueInst[_ExtractElement]
	ExtractValue   = valueInst[_ExtractValue]
	Invoke         = valueInst[_Invoke]
	Resume         = valueInst[_Resume]
	LandingPad     = valueInst[_LandingPad]
	CleanupRet     = valueInst[_CleanupRet]
	CatchRet       = valueInst[_CatchRet]
	CatchPad       = valueInst[_CatchPad]
	CleanupPad     = valueInst[_CleanupPad]
	CatchSwitch    = valueInst[_CatchSwitch]
	Switch         = terminatorInst[_Switch]
	Select         = valueInst[_Select]
)

func (b Builder) doBeforeCreateInst() context.Context {
	return nil
}

func (b Builder) doAfterCreateInst(_ context.Context, _ Instruction) {}

func (b Builder) CreateRet(v Value) (inst Return) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	if v == nil {
		return newTerminatorInst[_Return](binding.LLVMBuildRetVoid(b.binding()))
	} else {
		return newTerminatorInst[_Return](binding.LLVMBuildRet(b.binding(), v.binding()))
	}
}

func (b Builder) CreateBr(block Block) (inst Br) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newTerminatorInst[_Br](binding.LLVMBuildBr(b.binding(), block.binding()))
}

func (b Builder) CreateCondBr(cond Value, tblock, fblock Block) (inst Br) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newTerminatorInst[_Br](binding.LLVMBuildCondBr(b.binding(), cond.binding(), tblock.binding(), fblock.binding()))
}

func (b Builder) CreateUnreachable() (inst Unreachable) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newTerminatorInst[_Unreachable](binding.LLVMBuildUnreachable(b.binding()))
}

func (b Builder) CreateSAdd(name string, l, r Value) (inst Add) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Add](binding.LLVMBuildNSWAdd(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateUAdd(name string, l, r Value) (inst Add) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Add](binding.LLVMBuildNUWAdd(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateFAdd(name string, l, r Value) (inst Add) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Add](binding.LLVMBuildFAdd(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateSSub(name string, l, r Value) (inst Sub) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Sub](binding.LLVMBuildNSWSub(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateUSub(name string, l, r Value) (inst Sub) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Sub](binding.LLVMBuildNUWSub(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateFSub(name string, l, r Value) (inst Sub) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Sub](binding.LLVMBuildFSub(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateSMul(name string, l, r Value) (inst Mul) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Mul](binding.LLVMBuildNSWMul(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateUMul(name string, l, r Value) (inst Mul) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Mul](binding.LLVMBuildNUWMul(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateFMul(name string, l, r Value) (inst Mul) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Mul](binding.LLVMBuildFMul(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateSDiv(name string, l, r Value) (inst Div) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Div](binding.LLVMBuildSDiv(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateUDiv(name string, l, r Value) (inst Div) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Div](binding.LLVMBuildUDiv(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateFDiv(name string, l, r Value) (inst Div) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Div](binding.LLVMBuildFDiv(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateSRem(name string, l, r Value) (inst Rem) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Rem](binding.LLVMBuildSRem(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateURem(name string, l, r Value) (inst Rem) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Rem](binding.LLVMBuildURem(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateFRem(name string, l, r Value) (inst Rem) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Rem](binding.LLVMBuildFRem(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateShl(name string, l, r Value) (inst Shl) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Shl](binding.LLVMBuildShl(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateLShr(name string, l, r Value) (inst Shr) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Shr](binding.LLVMBuildLShr(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateAShr(name string, l, r Value) (inst Shr) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Shr](binding.LLVMBuildAShr(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateAnd(name string, l, r Value) (inst And) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_And](binding.LLVMBuildAnd(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateOr(name string, l, r Value) (inst Or) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Or](binding.LLVMBuildOr(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateXor(name string, l, r Value) (inst Xor) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Xor](binding.LLVMBuildXor(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateSNeg(name string, v Value) (inst Neg) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Neg](binding.LLVMBuildNSWNeg(b.binding(), v.binding(), name))
}

func (b Builder) CreateUNeg(name string, v Value) (inst Neg) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Neg](binding.LLVMBuildNUWNeg(b.binding(), v.binding(), name))
}

func (b Builder) CreateFNeg(name string, v Value) (inst Neg) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Neg](binding.LLVMBuildFNeg(b.binding(), v.binding(), name))
}

func (b Builder) CreateNot(name string, v Value) (inst Not) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Not](binding.LLVMBuildNot(b.binding(), v.binding(), name))
}

func (b Builder) CreateAlloca(name string, t Type) (inst Alloca) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInstWithAlign[_Alloca](binding.LLVMBuildAlloca(b.binding(), t.binding(), name))
}

func (b Builder) CreateAllocaWithSize(name string, t Type, size Value) (inst Alloca) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInstWithAlign[_Alloca](binding.LLVMBuildArrayAlloca(b.binding(), t.binding(), size.binding(), name))
}

func (b Builder) CreateLoad(name string, t Type, p Value) (inst Load) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInstWithAlign[_Load](binding.LLVMBuildLoad(b.binding(), t.binding(), p.binding(), name))
}

func (b Builder) CreateStore(from, to Value) (inst Store) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newGenericInstWithAlign[_Store](binding.LLVMBuildStore(b.binding(), from.binding(), to.binding()))
}

func (b Builder) CreateGEP(name string, t Type, from Value, indices ...Value) (inst GetElementPtr) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	var vs []binding.LLVMValueRef
	if len(indices) > 0 {
		vs = lo.Map(indices, func(item Value, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return newValueInst[_GetElementPtr](binding.LLVMBuildGEP(b.binding(), t.binding(), from.binding(), vs, name))
}

func (b Builder) CreateInBoundsGEP(name string, t Type, from Value, indices ...Value) (inst GetElementPtr) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	var vs []binding.LLVMValueRef
	if len(indices) > 0 {
		vs = lo.Map(indices, func(item Value, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return newValueInst[_GetElementPtr](binding.LLVMBuildInBoundsGEP(b.binding(), t.binding(), from.binding(), vs, name))
}

func (b Builder) CreateStructGEP(name string, t Type, from Value, i uint) (inst GetElementPtr) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_GetElementPtr](binding.LLVMBuildStructGEP(b.binding(), t.binding(), from.binding(), i, name))
}

func (b Builder) CreateTrunc(name string, from Value, to IntegerType) (inst Trunc) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Trunc](binding.LLVMBuildTrunc(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateFPTrunc(name string, from Value, to FloatType) (inst Trunc) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Trunc](binding.LLVMBuildFPTrunc(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateZExt(name string, from Value, to IntegerType) (inst Expand) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Expand](binding.LLVMBuildZExt(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateSExt(name string, from Value, to IntegerType) (inst Expand) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Expand](binding.LLVMBuildSExt(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateFPExt(name string, from Value, to FloatType) (inst Expand) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Expand](binding.LLVMBuildFPExt(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateFPToUI(name string, from Value, to IntegerType) (inst FloatToInt) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_FloatToInt](binding.LLVMBuildFPToUI(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateFPToSI(name string, from Value, to IntegerType) (inst FloatToInt) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_FloatToInt](binding.LLVMBuildFPToSI(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateUIToFP(name string, from Value, to FloatType) (inst IntToFloat) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_IntToFloat](binding.LLVMBuildUIToFP(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateSIToFP(name string, from Value, to FloatType) (inst IntToFloat) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_IntToFloat](binding.LLVMBuildSIToFP(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreatePtrToInt(name string, from Value, to IntegerType) (inst PtrToInt) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_PtrToInt](binding.LLVMBuildPtrToInt(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateIntToPtr(name string, from Value, to PointerType) (inst IntToPtr) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_IntToPtr](binding.LLVMBuildIntToPtr(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateBitCast(name string, from Value, to Type) (inst BitCast) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_BitCast](binding.LLVMBuildBitCast(b.binding(), from.binding(), to.binding(), name))
}

type IntPredicate binding.LLVMIntPredicate

const (
	IntEQ  = IntPredicate(binding.LLVMIntEQ)
	IntNE  = IntPredicate(binding.LLVMIntNE)
	IntUGT = IntPredicate(binding.LLVMIntUGT)
	IntUGE = IntPredicate(binding.LLVMIntUGE)
	IntULT = IntPredicate(binding.LLVMIntULT)
	IntULE = IntPredicate(binding.LLVMIntULE)
	IntSGT = IntPredicate(binding.LLVMIntSGT)
	IntSGE = IntPredicate(binding.LLVMIntSGE)
	IntSLT = IntPredicate(binding.LLVMIntSLT)
	IntSLE = IntPredicate(binding.LLVMIntSLE)
)

func (b Builder) CreateIntCmp(name string, op IntPredicate, l, r Value) (inst IntCmp) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newIntCmpInst(binding.LLVMBuildICmp(b.binding(), binding.LLVMIntPredicate(op), l.binding(), r.binding(), name))
}

type FloatPredicate binding.LLVMRealPredicate

const (
	FloatOEQ = FloatPredicate(binding.LLVMRealOEQ)
	FloatOGT = FloatPredicate(binding.LLVMRealOGT)
	FloatOGE = FloatPredicate(binding.LLVMRealOGE)
	FloatOLT = FloatPredicate(binding.LLVMRealOLT)
	FloatOLE = FloatPredicate(binding.LLVMRealOLE)
	FloatONE = FloatPredicate(binding.LLVMRealONE)
	FloatORD = FloatPredicate(binding.LLVMRealORD)
	FloatUNO = FloatPredicate(binding.LLVMRealUNO)
	FloatUEQ = FloatPredicate(binding.LLVMRealUEQ)
	FloatUGT = FloatPredicate(binding.LLVMRealUGT)
	FloatUGE = FloatPredicate(binding.LLVMRealUGE)
	FloatULT = FloatPredicate(binding.LLVMRealULT)
	FloatULE = FloatPredicate(binding.LLVMRealULE)
	FloatUNE = FloatPredicate(binding.LLVMRealUNE)
)

func (b Builder) CreateFloatCmp(name string, op FloatPredicate, l, r Value) (inst FloatCmp) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newFloatCmpInst(binding.LLVMBuildFCmp(b.binding(), binding.LLVMRealPredicate(op), l.binding(), r.binding(), name))
}

func (b Builder) CreateCall(name string, ft Type, fn Value, args ...Value) (inst Call) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	var as []binding.LLVMValueRef
	if len(args) > 0 {
		as = lo.Map(args, func(item Value, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return newValueInst[_Call](binding.LLVMBuildCall(b.binding(), ft.binding(), fn.binding(), as, name))
}

func (b Builder) CreateMalloc(name string, t Type) (inst Call) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Call](binding.LLVMBuildMalloc(b.binding(), t.binding(), name))
}

func (b Builder) CreateMallocWithSize(t Type, size Value, name string) (inst Call) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Call](binding.LLVMBuildArrayMalloc(b.binding(), t.binding(), size.binding(), name))
}

func (b Builder) CreateFree(p Value) (inst Call) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Call](binding.LLVMBuildFree(b.binding(), p.binding()))
}

func (b Builder) CreateMemSet(ptr, value, length Value, align uint32) (inst Call) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Call](binding.LLVMBuildMemSet(b.binding(), ptr.binding(), value.binding(), length.binding(), align))
}

func (b Builder) CreateMemCpy(dst Value, dstAlign32 uint32, src Value, srcAlign uint32, size Value) (inst Call) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Call](binding.LLVMBuildMemCpy(b.binding(), dst.binding(), dstAlign32, src.binding(), srcAlign, size.binding()))
}

func (b Builder) CreateMemMove(dst Value, dstAlign32 uint32, src Value, srcAlign uint32, size Value) (inst Call) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Call](binding.LLVMBuildMemMove(b.binding(), dst.binding(), dstAlign32, src.binding(), srcAlign, size.binding()))
}

func (b Builder) CreatePHI(name string, t Type, incomings ...struct {
	Value Value
	Block Block
}) (inst PHI) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	phi := newPhiInst(binding.LLVMBuildPhi(b.binding(), t.binding(), name))
	phi.AddIncomings(incomings...)
	return phi
}

func (b Builder) CreateExtractElement(name string, vec, index Value) (inst ExtractElement) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_ExtractElement](binding.LLVMBuildExtractElement(b.binding(), vec.binding(), index.binding(), name))
}

func (b Builder) CreateExtractValue(name string, vec Value, index uint) (inst ExtractValue) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_ExtractValue](binding.LLVMBuildExtractValue(b.binding(), vec.binding(), uint32(index), name))
}

func (b Builder) CreateInvoke(name string, ft Type, then, catch Block, fn Value, args ...Value) (inst Invoke) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	var as []binding.LLVMValueRef
	if len(args) > 0 {
		as = lo.Map(args, func(item Value, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return newValueInst[_Invoke](binding.LLVMBuildInvoke(b.binding(), ft.binding(), fn.binding(), as, then.binding(), catch.binding(), name))
}

func (b Builder) CreateResume(v *Value) (inst Resume) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	if v == nil {
		ctx := b.CurrentBlock().Belong().FunctionType().Context()
		return newValueInst[_Resume](binding.LLVMBuildResume(b.binding(), ctx.ConstNull(ctx.VoidType()).binding()))
	}
	return newValueInst[_Resume](binding.LLVMBuildResume(b.binding(), (*v).binding()))
}

func (b Builder) CreateLandingPad(name string, t Type, f Value, numClauses uint32) (inst LandingPad) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_LandingPad](binding.LLVMBuildLandingPad(b.binding(), t.binding(), f.binding(), numClauses, name))
}

func (b Builder) CreateCleanupRet(catchPad Value, bb Block) (inst CleanupRet) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_CleanupRet](binding.LLVMBuildCleanupRet(b.binding(), catchPad.binding(), bb.binding()))
}

func (b Builder) CreateCatchRet(catchPad Value, bb Block) (inst CatchRet) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_CatchRet](binding.LLVMBuildCatchRet(b.binding(), catchPad.binding(), bb.binding()))
}

func (b Builder) CreateCatchPad(name string, parentPad Value, arg ...Value) (inst CatchPad) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	var args []binding.LLVMValueRef
	if len(arg) > 0 {
		args = lo.Map(arg, func(item Value, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return newValueInst[_CatchPad](binding.LLVMBuildCatchPad(b.binding(), parentPad.binding(), args, name))
}

func (b Builder) CreateCleanupPad(name string, parentPad Value, arg ...Value) (inst CleanupPad) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	var args []binding.LLVMValueRef
	if len(arg) > 0 {
		args = lo.Map(arg, func(item Value, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return newValueInst[_CleanupPad](binding.LLVMBuildCleanupPad(b.binding(), parentPad.binding(), args, name))
}

func (b Builder) CreateCatchSwitch(name string, parentPad Value, unwindBB Block, numHandlers uint32) (inst CatchSwitch) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_CatchSwitch](binding.LLVMBuildCatchSwitch(b.binding(), parentPad.binding(), unwindBB.binding(), numHandlers, name))
}

func (b Builder) CreateSwitch(v Value, defaultBlock Block, conds ...struct {
	Value Value
	Block Block
}) (inst Switch) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	instObj := binding.LLVMBuildSwitch(b.binding(), v.binding(), defaultBlock.binding(), uint32(len(conds)))
	for _, cond := range conds {
		binding.LLVMAddCase(instObj, cond.Value.binding(), cond.Block.binding())
	}
	return newTerminatorInst[_Switch](instObj)
}

func (b Builder) CreateSelect(name string, cond Value, trueValue, falseValue Value) (inst Select) {
	defer func(ctx context.Context) { b.doAfterCreateInst(ctx, inst) }(b.doBeforeCreateInst())

	return newValueInst[_Select](binding.LLVMBuildSelect(b.binding(), cond.binding(), trueValue.binding(), falseValue.binding(), name))
}
