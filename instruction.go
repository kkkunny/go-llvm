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
	case binding.LLVMShl, binding.LLVMLShr, binding.LLVMAShr:
		return Shl(ref)
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
	case binding.LLVMCall:
		return Call(ref)
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

func (b Builder) CreateLShr(name string, l, r Value) Shl {
	return Shl(binding.LLVMBuildLShr(b.binding(), l.binding(), r.binding(), name))
}

func (b Builder) CreateAShr(name string, l, r Value) Shl {
	return Shl(binding.LLVMBuildAShr(b.binding(), l.binding(), r.binding(), name))
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

func (b Builder) CreateAllocaWithSize(t Type, size Value, name string) Alloca {
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
