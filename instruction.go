package llvm

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/kkkunny/go-llvm/internal/binding"
)

type Instruction interface {
	binding() binding.LLVMValueRef
	Belong() Block
}

func lookupInstruction(ref binding.LLVMValueRef) Instruction {
	switch binding.LLVMGetInstructionOpcode(ref) {
	case binding.LLVMRet:
		return Return(ref)
	case binding.LLVMBr:
		return Br(ref)
	case binding.LLVMUnreachable:
		return Unreachable(ref)
	case binding.LLVMAdd, binding.LLVMFAdd:
		return Add(ref)
	case binding.LLVMSub, binding.LLVMFSub:
		return Sub(ref)
	case binding.LLVMMul, binding.LLVMFMul:
		return Mul(ref)
	case binding.LLVMSDiv, binding.LLVMUDiv, binding.LLVMFDiv:
		return Div(ref)
	case binding.LLVMSRem, binding.LLVMURem, binding.LLVMFRem:
		return Rem(ref)
	case binding.LLVMShl:
		return Shl(ref)
	case binding.LLVMLShr, binding.LLVMAShr:
		return Shr(ref)
	case binding.LLVMAnd:
		return And(ref)
	case binding.LLVMOr:
		return Or(ref)
	case binding.LLVMXor:
		return Xor(ref)
	case binding.LLVMFNeg:
		return Neg(ref)
	case binding.LLVMAlloca:
		return Alloca(ref)
	case binding.LLVMLoad:
		return Load(ref)
	case binding.LLVMStore:
		return Store(ref)
	case binding.LLVMGetElementPtr:
		return GetElementPtr(ref)
	case binding.LLVMTrunc, binding.LLVMFPTrunc:
		return Trunc(ref)
	case binding.LLVMZExt, binding.LLVMSExt, binding.LLVMFPExt:
		return Expand(ref)
	case binding.LLVMFPToUI, binding.LLVMFPToSI:
		return FloatToInt(ref)
	case binding.LLVMUIToFP, binding.LLVMSIToFP:
		return IntToFloat(ref)
	case binding.LLVMPtrToInt:
		return PtrToInt(ref)
	case binding.LLVMIntToPtr:
		return IntToPtr(ref)
	case binding.LLVMBitCast:
		return BitCast(ref)
	case binding.LLVMICmp:
		return IntCmp(ref)
	case binding.LLVMFCmp:
		return IntCmp(ref)
	case binding.LLVMPHI:
		return PHI(ref)
	case binding.LLVMCall:
		return Call(ref)
	case binding.LLVMInvoke:
		return Invoke(ref)
	case binding.LLVMExtractElement:
		return ExtractElement(ref)
	case binding.LLVMExtractValue:
		return ExtractValue(ref)
	case binding.LLVMResume:
		return Resume(ref)
	case binding.LLVMLandingPad:
		return LandingPad(ref)
	case binding.LLVMCleanupRet:
		return CleanupRet(ref)
	case binding.LLVMCatchRet:
		return CatchRet(ref)
	case binding.LLVMCatchPad:
		return CatchPad(ref)
	case binding.LLVMCleanupPad:
		return CleanupPad(ref)
	case binding.LLVMCatchSwitch:
		return CatchSwitch(ref)
	case binding.LLVMSwitch:
		return Switch(ref)
	case binding.LLVMSelect:
		return Select(ref)
	default:
		panic(fmt.Errorf("unknown enum value `%d`", binding.LLVMGetInstructionOpcode(ref)))
	}
}

type Terminator interface {
	Instruction
	terminator()
}

func lookupTerminator(ref binding.LLVMValueRef) Terminator {
	switch binding.LLVMGetInstructionOpcode(ref) {
	case binding.LLVMRet:
		return Return(ref)
	case binding.LLVMBr:
		return Br(ref)
	case binding.LLVMUnreachable:
		return Unreachable(ref)
	default:
		panic(fmt.Errorf("unknown enum value `%d`", binding.LLVMGetInstructionOpcode(ref)))
	}
}

type Return binding.LLVMValueRef

func (b Builder) CreateRet(v *Value) Return {
	if v == nil {
		return Return(binding.LLVMBuildRetVoid(b.binding()))
	} else {
		return Return(binding.LLVMBuildRet(b.binding(), (*v).binding()))
	}
}

func (i Return) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Return) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (Return) terminator() {}

type Br binding.LLVMValueRef

func (b Builder) CreateBr(block Block) Br {
	return Br(binding.LLVMBuildBr(b.binding(), block.binding()))
}

func (b Builder) CreateCondBr(cond Value, tblock, fblock Block) Br {
	return Br(binding.LLVMBuildCondBr(b.binding(), cond.binding(), tblock.binding(), fblock.binding()))
}

func (i Br) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Br) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (Br) terminator() {}

type Unreachable binding.LLVMValueRef

func (b Builder) CreateUnreachable() Unreachable {
	return Unreachable(binding.LLVMBuildUnreachable(b.binding()))
}

func (i Unreachable) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Unreachable) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (Unreachable) terminator() {}

type Add binding.LLVMValueRef

func (b Builder) CreateSAdd(name string, l, r Value) Add {
	return Add(binding.LLVMBuildNSWAdd(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateUAdd(name string, l, r Value) Add {
	return Add(binding.LLVMBuildNUWAdd(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateFAdd(name string, l, r Value) Add {
	return Add(binding.LLVMBuildFAdd(b.binding(), l.binding(), r.binding(), name))
}

func (i Add) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Add) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Add) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Add) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Sub binding.LLVMValueRef

func (b Builder) CreateSSub(name string, l, r Value) Sub {
	return Sub(binding.LLVMBuildNSWSub(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateUSub(name string, l, r Value) Sub {
	return Sub(binding.LLVMBuildNUWSub(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateFSub(name string, l, r Value) Sub {
	return Sub(binding.LLVMBuildFSub(b.binding(), l.binding(), r.binding(), name))
}

func (i Sub) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Sub) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Sub) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Sub) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Mul binding.LLVMValueRef

func (b Builder) CreateSMul(name string, l, r Value) Mul {
	return Mul(binding.LLVMBuildNSWMul(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateUMul(name string, l, r Value) Mul {
	return Mul(binding.LLVMBuildNUWMul(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateFMul(name string, l, r Value) Mul {
	return Mul(binding.LLVMBuildFMul(b.binding(), l.binding(), r.binding(), name))
}

func (i Mul) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Mul) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Mul) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Mul) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Div binding.LLVMValueRef

func (b Builder) CreateSDiv(name string, l, r Value) Div {
	return Div(binding.LLVMBuildSDiv(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateUDiv(name string, l, r Value) Div {
	return Div(binding.LLVMBuildUDiv(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateFDiv(name string, l, r Value) Div {
	return Div(binding.LLVMBuildFDiv(b.binding(), l.binding(), r.binding(), name))
}

func (i Div) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Div) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Div) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Div) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Rem binding.LLVMValueRef

func (b Builder) CreateSRem(name string, l, r Value) Rem {
	return Rem(binding.LLVMBuildSRem(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateURem(name string, l, r Value) Rem {
	return Rem(binding.LLVMBuildURem(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateFRem(name string, l, r Value) Rem {
	return Rem(binding.LLVMBuildFRem(b.binding(), l.binding(), r.binding(), name))
}

func (i Rem) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Rem) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Rem) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Rem) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Shl binding.LLVMValueRef

func (b Builder) CreateShl(name string, l, r Value) Shl {
	return Shl(binding.LLVMBuildShl(b.binding(), l.binding(), r.binding(), name))
}

func (i Shl) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Shl) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Shl) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Shl) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Shr binding.LLVMValueRef

func (b Builder) CreateLShr(name string, l, r Value) Shr {
	return Shr(binding.LLVMBuildLShr(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateAShr(name string, l, r Value) Shr {
	return Shr(binding.LLVMBuildAShr(b.binding(), l.binding(), r.binding(), name))
}

func (i Shr) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Shr) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Shr) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Shr) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type And binding.LLVMValueRef

func (b Builder) CreateAnd(name string, l, r Value) And {
	return And(binding.LLVMBuildAnd(b.binding(), l.binding(), r.binding(), name))
}

func (i And) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i And) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v And) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v And) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Or binding.LLVMValueRef

func (b Builder) CreateOr(name string, l, r Value) Or {
	return Or(binding.LLVMBuildOr(b.binding(), l.binding(), r.binding(), name))
}

func (i Or) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Or) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Or) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Or) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Xor binding.LLVMValueRef

func (b Builder) CreateXor(name string, l, r Value) Xor {
	return Xor(binding.LLVMBuildXor(b.binding(), l.binding(), r.binding(), name))
}

func (i Xor) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Xor) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Xor) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Xor) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Neg binding.LLVMValueRef

func (b Builder) CreateSNeg(name string, v Value) Neg {
	return Neg(binding.LLVMBuildNSWNeg(b.binding(), v.binding(), name))
}

func (b Builder) CreateUNeg(name string, v Value) Neg {
	return Neg(binding.LLVMBuildNUWNeg(b.binding(), v.binding(), name))
}

func (b Builder) CreateFNeg(name string, v Value) Neg {
	return Neg(binding.LLVMBuildFNeg(b.binding(), v.binding(), name))
}

func (i Neg) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Neg) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Neg) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Neg) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Not binding.LLVMValueRef

func (b Builder) CreateNot(name string, v Value) Not {
	return Not(binding.LLVMBuildNot(b.binding(), v.binding(), name))
}

func (i Not) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Not) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Not) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Not) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Alloca binding.LLVMValueRef

func (b Builder) CreateAlloca(name string, t Type) Alloca {
	return Alloca(binding.LLVMBuildAlloca(b.binding(), t.binding(), name))
}

func (b Builder) CreateAllocaWithSize(name string, t Type, size Value) Alloca {
	return Alloca(binding.LLVMBuildArrayAlloca(b.binding(), t.binding(), size.binding(), name))
}

func (i Alloca) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Alloca) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Alloca) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Alloca) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

func (v Alloca) SetAlign(align uint32) {
	binding.LLVMSetAlignment(v.binding(), uint32(align))
}

func (v Alloca) GetAlign() uint32 {
	return binding.LLVMGetAlignment(v.binding())
}

type Load binding.LLVMValueRef

func (b Builder) CreateLoad(name string, t Type, p Value) Load {
	return Load(binding.LLVMBuildLoad(b.binding(), t.binding(), p.binding(), name))
}

func (i Load) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Load) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Load) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Load) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

func (v Load) SetAlign(align uint32) {
	binding.LLVMSetAlignment(v.binding(), uint32(align))
}

func (v Load) GetAlign() uint32 {
	return binding.LLVMGetAlignment(v.binding())
}

type Store binding.LLVMValueRef

func (b Builder) CreateStore(from, to Value) Store {
	return Store(binding.LLVMBuildStore(b.binding(), from.binding(), to.binding()))
}

func (i Store) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Store) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (i Store) SetAlign(align uint32) {
	binding.LLVMSetAlignment(i.binding(), uint32(align))
}

func (i Store) GetAlign() uint32 {
	return binding.LLVMGetAlignment(i.binding())
}

type GetElementPtr binding.LLVMValueRef

func (b Builder) CreateGEP(name string, t Type, from Value, indices ...Value) GetElementPtr {
	var vs []binding.LLVMValueRef
	if len(indices) > 0 {
		vs = lo.Map(indices, func(item Value, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return GetElementPtr(binding.LLVMBuildGEP(b.binding(), t.binding(), from.binding(), vs, name))
}

func (b Builder) CreateInBoundsGEP(name string, t Type, from Value, indices ...Value) GetElementPtr {
	var vs []binding.LLVMValueRef
	if len(indices) > 0 {
		vs = lo.Map(indices, func(item Value, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return GetElementPtr(binding.LLVMBuildInBoundsGEP(b.binding(), t.binding(), from.binding(), vs, name))
}

func (b Builder) CreateStructGEP(name string, t Type, from Value, i uint) GetElementPtr {
	return GetElementPtr(binding.LLVMBuildStructGEP(b.binding(), t.binding(), from.binding(), i, name))
}

func (i GetElementPtr) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i GetElementPtr) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v GetElementPtr) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v GetElementPtr) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Trunc binding.LLVMValueRef

func (b Builder) CreateTrunc(name string, from Value, to IntegerType) Trunc {
	return Trunc(binding.LLVMBuildTrunc(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateFPTrunc(name string, from Value, to FloatType) Trunc {
	return Trunc(binding.LLVMBuildFPTrunc(b.binding(), from.binding(), to.binding(), name))
}

func (i Trunc) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (v Trunc) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Trunc) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

func (i Trunc) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

type Expand binding.LLVMValueRef

func (b Builder) CreateZExt(name string, from Value, to IntegerType) Expand {
	return Expand(binding.LLVMBuildZExt(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateSExt(name string, from Value, to IntegerType) Expand {
	return Expand(binding.LLVMBuildSExt(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateFPExt(name string, from Value, to FloatType) Expand {
	return Expand(binding.LLVMBuildFPExt(b.binding(), from.binding(), to.binding(), name))
}

func (i Expand) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Expand) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Expand) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Expand) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type FloatToInt binding.LLVMValueRef

func (b Builder) CreateFPToUI(name string, from Value, to IntegerType) FloatToInt {
	return FloatToInt(binding.LLVMBuildFPToUI(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateFPToSI(name string, from Value, to IntegerType) FloatToInt {
	return FloatToInt(binding.LLVMBuildFPToSI(b.binding(), from.binding(), to.binding(), name))
}

func (i FloatToInt) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i FloatToInt) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v FloatToInt) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v FloatToInt) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type IntToFloat binding.LLVMValueRef

func (b Builder) CreateUIToFP(name string, from Value, to FloatType) IntToFloat {
	return IntToFloat(binding.LLVMBuildUIToFP(b.binding(), from.binding(), to.binding(), name))
}

func (b Builder) CreateSIToFP(name string, from Value, to FloatType) IntToFloat {
	return IntToFloat(binding.LLVMBuildSIToFP(b.binding(), from.binding(), to.binding(), name))
}

func (i IntToFloat) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i IntToFloat) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v IntToFloat) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v IntToFloat) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type PtrToInt binding.LLVMValueRef

func (b Builder) CreatePtrToInt(name string, from Value, to IntegerType) PtrToInt {
	return PtrToInt(binding.LLVMBuildPtrToInt(b.binding(), from.binding(), to.binding(), name))
}

func (i PtrToInt) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i PtrToInt) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v PtrToInt) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v PtrToInt) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type IntToPtr binding.LLVMValueRef

func (b Builder) CreateIntToPtr(name string, from Value, to PointerType) IntToPtr {
	return IntToPtr(binding.LLVMBuildIntToPtr(b.binding(), from.binding(), to.binding(), name))
}

func (i IntToPtr) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i IntToPtr) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v IntToPtr) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v IntToPtr) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type BitCast binding.LLVMValueRef

func (b Builder) CreateBitCast(name string, from Value, to Type) BitCast {
	return BitCast(binding.LLVMBuildBitCast(b.binding(), from.binding(), to.binding(), name))
}

func (i BitCast) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i BitCast) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v BitCast) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v BitCast) String() string {
	return binding.LLVMPrintValueToString(v.binding())
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

type IntCmp binding.LLVMValueRef

func (b Builder) CreateIntCmp(name string, op IntPredicate, l, r Value) IntCmp {
	return IntCmp(binding.LLVMBuildICmp(b.binding(), binding.LLVMIntPredicate(op), l.binding(), r.binding(), name))
}

func (i IntCmp) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i IntCmp) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v IntCmp) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v IntCmp) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

func (c IntCmp) GetOperator() IntPredicate {
	return IntPredicate(binding.LLVMGetICmpPredicate(c.binding()))
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

type FloatCmp binding.LLVMValueRef

func (b Builder) CreateFloatCmp(name string, op FloatPredicate, l, r Value) FloatCmp {
	return FloatCmp(binding.LLVMBuildFCmp(b.binding(), binding.LLVMRealPredicate(op), l.binding(), r.binding(), name))
}

func (i FloatCmp) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i FloatCmp) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v FloatCmp) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v FloatCmp) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

func (c FloatCmp) GetOperator() FloatPredicate {
	return FloatPredicate(binding.LLVMGetFCmpPredicate(c.binding()))
}

type Call binding.LLVMValueRef

func (b Builder) CreateCall(name string, ft Type, fn Value, args ...Value) Call {
	var as []binding.LLVMValueRef
	if len(args) > 0 {
		as = lo.Map(args, func(item Value, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return Call(binding.LLVMBuildCall(b.binding(), ft.binding(), fn.binding(), as, name))
}

func (b Builder) CreateMalloc(name string, t Type) Call {
	return Call(binding.LLVMBuildMalloc(b.binding(), t.binding(), name))
}

func (b Builder) CreateMallocWithSize(t Type, size Value, name string) Call {
	return Call(binding.LLVMBuildArrayMalloc(b.binding(), t.binding(), size.binding(), name))
}

func (b Builder) CreateFree(p Value) Call {
	return Call(binding.LLVMBuildFree(b.binding(), p.binding()))
}

func (b Builder) CreateMemSet(ptr, value, length Value, align uint32) Call {
	return Call(binding.LLVMBuildMemSet(b.binding(), ptr.binding(), value.binding(), length.binding(), align))
}

func (b Builder) CreateMemCpy(dst Value, dstAlign32 uint32, src Value, srcAlign uint32, size Value) Call {
	return Call(binding.LLVMBuildMemCpy(b.binding(), dst.binding(), dstAlign32, src.binding(), srcAlign, size.binding()))
}

func (b Builder) CreateMemMove(dst Value, dstAlign32 uint32, src Value, srcAlign uint32, size Value) Call {
	return Call(binding.LLVMBuildMemMove(b.binding(), dst.binding(), dstAlign32, src.binding(), srcAlign, size.binding()))
}

func (i Call) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Call) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Call) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Call) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type PHI binding.LLVMValueRef

func (b Builder) CreatePHI(name string, t Type, incomings ...struct {
	Value Value
	Block Block
}) PHI {
	phi := PHI(binding.LLVMBuildPhi(b.binding(), t.binding(), name))
	phi.AddIncomings(incomings...)
	return phi
}

func (i PHI) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i PHI) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v PHI) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v PHI) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

func (phi PHI) AddIncomings(incomings ...struct {
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
	binding.LLVMAddIncoming(phi.binding(), values, blocks)
}

func (phi PHI) CountIncomings() uint {
	return uint(binding.LLVMCountIncoming(phi.binding()))
}

func (phi PHI) GetIncoming(i uint) (Value, Block) {
	return lookupValue(binding.LLVMGetIncomingValue(phi.binding(), uint32(i))), Block(binding.LLVMGetIncomingBlock(phi.binding(), uint32(i)))
}

func (phi PHI) Incomings() (res []struct {
	Value Value
	Block Block
}) {
	for i := uint(0); i < phi.CountIncomings(); i++ {
		v, b := phi.GetIncoming(i)
		res = append(res, struct {
			Value Value
			Block Block
		}{Value: v, Block: b})
	}
	return res
}

type ExtractElement binding.LLVMValueRef

func (b Builder) CreateExtractElement(name string, vec, index Value) ExtractElement {
	return ExtractElement(binding.LLVMBuildExtractElement(b.binding(), vec.binding(), index.binding(), name))
}

func (i ExtractElement) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i ExtractElement) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v ExtractElement) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v ExtractElement) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type ExtractValue binding.LLVMValueRef

func (b Builder) CreateExtractValue(name string, vec Value, index uint) ExtractValue {
	return ExtractValue(binding.LLVMBuildExtractValue(b.binding(), vec.binding(), uint32(index), name))
}

func (i ExtractValue) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i ExtractValue) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v ExtractValue) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v ExtractValue) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Invoke binding.LLVMValueRef

func (b Builder) CreateInvoke(name string, ft Type, then, catch Block, fn Value, args ...Value) Invoke {
	var as []binding.LLVMValueRef
	if len(args) > 0 {
		as = lo.Map(args, func(item Value, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return Invoke(binding.LLVMBuildInvoke(b.binding(), ft.binding(), fn.binding(), as, then.binding(), catch.binding(), name))
}

func (i Invoke) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Invoke) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Invoke) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Invoke) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Resume binding.LLVMValueRef

func (b Builder) CreateResume(v *Value) Resume {
	if v == nil {
		ctx := b.CurrentBlock().Belong().FunctionType().Context()
		return Resume(binding.LLVMBuildResume(b.binding(), ctx.ConstNull(ctx.VoidType()).binding()))
	}
	return Resume(binding.LLVMBuildResume(b.binding(), (*v).binding()))
}

func (i Resume) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Resume) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Resume) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Resume) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type LandingPad binding.LLVMValueRef

func (b Builder) CreateLandingPad(name string, t Type, f Value, numClauses uint32) LandingPad {
	return LandingPad(binding.LLVMBuildLandingPad(b.binding(), t.binding(), f.binding(), numClauses, name))
}

func (i LandingPad) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i LandingPad) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v LandingPad) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v LandingPad) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type CleanupRet binding.LLVMValueRef

func (b Builder) CreateCleanupRet(catchPad Value, bb Block) CleanupRet {
	return CleanupRet(binding.LLVMBuildCleanupRet(b.binding(), catchPad.binding(), bb.binding()))
}

func (i CleanupRet) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i CleanupRet) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v CleanupRet) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v CleanupRet) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type CatchRet binding.LLVMValueRef

func (b Builder) CreateCatchRet(catchPad Value, bb Block) CatchRet {
	return CatchRet(binding.LLVMBuildCatchRet(b.binding(), catchPad.binding(), bb.binding()))
}

func (i CatchRet) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i CatchRet) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v CatchRet) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v CatchRet) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type CatchPad binding.LLVMValueRef

func (b Builder) CreateCatchPad(name string, parentPad Value, arg ...Value) CatchPad {
	var args []binding.LLVMValueRef
	if len(arg) > 0 {
		args = lo.Map(arg, func(item Value, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return CatchPad(binding.LLVMBuildCatchPad(b.binding(), parentPad.binding(), args, name))
}

func (i CatchPad) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i CatchPad) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v CatchPad) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v CatchPad) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type CleanupPad binding.LLVMValueRef

func (b Builder) CreateCleanupPad(name string, parentPad Value, arg ...Value) CleanupPad {
	var args []binding.LLVMValueRef
	if len(arg) > 0 {
		args = lo.Map(arg, func(item Value, index int) binding.LLVMValueRef {
			return item.binding()
		})
	}
	return CleanupPad(binding.LLVMBuildCleanupPad(b.binding(), parentPad.binding(), args, name))
}

func (i CleanupPad) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i CleanupPad) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v CleanupPad) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v CleanupPad) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type CatchSwitch binding.LLVMValueRef

func (b Builder) CreateCatchSwitch(name string, parentPad Value, unwindBB Block, numHandlers uint32) CatchSwitch {
	return CatchSwitch(binding.LLVMBuildCatchSwitch(b.binding(), parentPad.binding(), unwindBB.binding(), numHandlers, name))
}

func (i CatchSwitch) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i CatchSwitch) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v CatchSwitch) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v CatchSwitch) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}

type Switch binding.LLVMValueRef

func (b Builder) CreateSwitch(v Value, defaultBlock Block, conds ...struct {
	Value Value
	Block Block
}) Switch {
	inst := binding.LLVMBuildSwitch(b.binding(), v.binding(), defaultBlock.binding(), uint32(len(conds)))
	for _, cond := range conds{
		binding.LLVMAddCase(inst, cond.Value.binding(), cond.Block.binding())
	}
	return Switch(inst)
}

func (i Switch) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Switch) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (Switch) terminator() {}

type Select binding.LLVMValueRef

func (b Builder) CreateSelect(name string, cond Value, trueValue, falseValue Value) Select {
	return Select(binding.LLVMBuildSelect(b.binding(), cond.binding(), trueValue.binding(), falseValue.binding(), name))
}

func (i Select) binding() binding.LLVMValueRef {
	return binding.LLVMValueRef(i)
}

func (i Select) Belong() Block {
	return Block(binding.LLVMGetInstructionParent(i.binding()))
}

func (v Select) Type() Type {
	return lookupType(binding.LLVMTypeOf(v.binding()))
}

func (v Select) String() string {
	return binding.LLVMPrintValueToString(v.binding())
}
