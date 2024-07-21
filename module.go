package llvm

import (
	"errors"
	"strings"

	"github.com/kkkunny/go-llvm/internal/binding"
)

type Module struct {
	binding.LLVMModuleRef
	target *Target
}

func newModule(ref binding.LLVMModuleRef) Module {
	return Module{LLVMModuleRef: ref}
}

func (ctx Context) NewModule(name string) Module {
	return newModule(binding.LLVMModuleCreateWithNameInContext(name, ctx.binding()))
}

func (m Module) binding() binding.LLVMModuleRef {
	return m.LLVMModuleRef
}

func (m Module) Free() {
	binding.LLVMDisposeModule(m.binding())
}

func (m Module) String() string {
	return binding.LLVMPrintModuleToString(m.binding())
}

func (m Module) Clone() Module {
	return Module{
		LLVMModuleRef: binding.LLVMCloneModule(m.binding()),
		target:        m.target,
	}
}

func (m Module) Context() Context {
	return Context(binding.LLVMGetModuleContext(m.binding()))
}

func (m Module) GetSource() string {
	return binding.LLVMGetSourceFileName(m.binding())
}

func (m Module) SetSource(source string) {
	binding.LLVMSetSourceFileName(m.binding(), source)
}

// RunPasses https://llvm.org/docs/Passes.html
func (m Module) RunPasses(target *Target, option PassOption, pass ...string) error {
	return binding.LLVMRunPasses(m.binding(), strings.Join(pass, ","), target.machine, option.binding())
}

func (m Module) Verify() error {
	msg, fail := binding.LLVMVerifyModule(m.binding(), binding.LLVMReturnStatusAction)
	if fail {
		return errors.New(msg)
	}
	return nil
}

func (m Module) AddConstructor(prior uint16, f Function) {
	name := "llvm.global_ctors"
	ctx := m.Context()
	ft := ctx.FunctionType(false, ctx.VoidType())
	st := ctx.StructType(false, ctx.IntegerType(32), ctx.PointerType(ft), ctx.PointerType(ctx.IntegerType(8)))
	stv := ctx.ConstNamedStruct(st, ctx.ConstInteger(st.GetElem(0).(IntegerType), int64(prior)), f, ctx.ConstNull(st.GetElem(2)))

	var elems []Constant
	ctors, ok := m.GetGlobal(name)
	if ok {
		v, ok := ctors.GetInitializer()
		if ok {
			av := v.(ConstArray)
			avt := av.Type().(ArrayType)
			for i := uint32(0); i < avt.Capacity(); i++ {
				elems = append(elems, av.GetElem(uint(i)))
			}
		}
		m.DelGlobal(ctors)
	}
	elems = append(elems, stv)
	ctors = m.NewGlobal(name, ctx.ArrayType(st, uint32(len(elems))), ctx.ConstArray(st, elems...))
	ctors.SetLinkage(AppendingLinkage)
}

func (m Module) AddDestructor(prior uint16, f Function) {
	name := "llvm.global_dtors"
	ctx := m.Context()
	ft := ctx.FunctionType(false, ctx.VoidType())
	st := ctx.StructType(false, ctx.IntegerType(32), ctx.PointerType(ft), ctx.PointerType(ctx.IntegerType(8)))
	stv := ctx.ConstNamedStruct(st, ctx.ConstInteger(st.GetElem(0).(IntegerType), int64(prior)), f, ctx.ConstNull(st.GetElem(2)))

	var elems []Constant
	dtors, ok := m.GetGlobal(name)
	if ok {
		v, ok := dtors.GetInitializer()
		if ok {
			av := v.(ConstArray)
			avt := av.Type().(ArrayType)
			for i := uint32(0); i < avt.Capacity(); i++ {
				elems = append(elems, av.GetElem(uint(i)))
			}
		}
		m.DelGlobal(dtors)
	}
	elems = append(elems, stv)
	dtors = m.NewGlobal(name, ctx.ArrayType(st, uint32(len(elems))), ctx.ConstArray(st, elems...))
	dtors.SetLinkage(AppendingLinkage)
}
