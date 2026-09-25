package binding

/*
#include "llvm-c/Orc.h"
#include "llvm-c/LLJIT.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type (
	// LLVMOrcLLJITRef A reference to an orc::LLJIT instance.
	LLVMOrcLLJITRef struct{ c C.LLVMOrcLLJITRef }

	// LLVMOrcLLJITBuilderRef A reference to an orc::LLJITBuilder instance.
	LLVMOrcLLJITBuilderRef struct{ c C.LLVMOrcLLJITBuilderRef }

	// LLVMOrcJITDylibRef A reference to an orc::JITDylib instance.
	LLVMOrcJITDylibRef struct{ c C.LLVMOrcJITDylibRef }

	// LLVMOrcExecutionSessionRef A reference to an orc::ExecutionSession instance.
	LLVMOrcExecutionSessionRef struct{ c C.LLVMOrcExecutionSessionRef }

	// LLVMOrcSymbolStringPoolEntryRef A reference to an orc::SymbolStringPool table entry.
	LLVMOrcSymbolStringPoolEntryRef struct {
		c C.LLVMOrcSymbolStringPoolEntryRef
	}

	// LLVMOrcResourceTrackerRef A reference to an orc::ResourceTracker instance.
	LLVMOrcResourceTrackerRef struct{ c C.LLVMOrcResourceTrackerRef }

	// LLVMOrcThreadSafeModuleRef A reference to an orc::ThreadSafeModule instance.
	LLVMOrcThreadSafeModuleRef struct{ c C.LLVMOrcThreadSafeModuleRef }

	// LLVMOrcThreadSafeContextRef A reference to an orc::ThreadSafeContext instance.
	LLVMOrcThreadSafeContextRef struct{ c C.LLVMOrcThreadSafeContextRef }

	// LLVMOrcJITTargetMachineBuilderRef A reference to an orc::JITTargetMachineBuilder instance.
	LLVMOrcJITTargetMachineBuilderRef struct {
		c C.LLVMOrcJITTargetMachineBuilderRef
	}

	// LLVMOrcMaterializationUnitRef A reference to an orc::MaterializationUnit instance.
	LLVMOrcMaterializationUnitRef struct {
		c C.LLVMOrcMaterializationUnitRef
	}
)

func (ref LLVMOrcLLJITRef) IsNil() bool                   { return ref.c == nil }
func (ref LLVMOrcLLJITBuilderRef) IsNil() bool            { return ref.c == nil }
func (ref LLVMOrcJITDylibRef) IsNil() bool                { return ref.c == nil }
func (ref LLVMOrcExecutionSessionRef) IsNil() bool        { return ref.c == nil }
func (ref LLVMOrcSymbolStringPoolEntryRef) IsNil() bool   { return ref.c == nil }
func (ref LLVMOrcResourceTrackerRef) IsNil() bool         { return ref.c == nil }
func (ref LLVMOrcThreadSafeModuleRef) IsNil() bool        { return ref.c == nil }
func (ref LLVMOrcThreadSafeContextRef) IsNil() bool       { return ref.c == nil }
func (ref LLVMOrcJITTargetMachineBuilderRef) IsNil() bool { return ref.c == nil }
func (ref LLVMOrcMaterializationUnitRef) IsNil() bool     { return ref.c == nil }

// LLVMJITSymbolGenericFlags Represents generic linkage flags for a symbol definition.
type LLVMJITSymbolGenericFlags uint8

const (
	LLVMJITSymbolGenericFlagsNone                           LLVMJITSymbolGenericFlags = C.LLVMJITSymbolGenericFlagsNone
	LLVMJITSymbolGenericFlagsExported                       LLVMJITSymbolGenericFlags = C.LLVMJITSymbolGenericFlagsExported
	LLVMJITSymbolGenericFlagsWeak                           LLVMJITSymbolGenericFlags = C.LLVMJITSymbolGenericFlagsWeak
	LLVMJITSymbolGenericFlagsCallable                       LLVMJITSymbolGenericFlags = C.LLVMJITSymbolGenericFlagsCallable
	LLVMJITSymbolGenericFlagsMaterializationSideEffectsOnly LLVMJITSymbolGenericFlags = C.LLVMJITSymbolGenericFlagsMaterializationSideEffectsOnly
)

// LLVMJITSymbolFlags Represents the linkage flags for a symbol definition.
type LLVMJITSymbolFlags struct {
	GenericFlags LLVMJITSymbolGenericFlags
	TargetFlags  uint8
}

// LLVMJITEvaluatedSymbol Represents an evaluated symbol address and flags.
type LLVMJITEvaluatedSymbol struct {
	Address uint64
	Flags   LLVMJITSymbolFlags
}

// LLVMOrcCSymbolMapPair Represents a pair of a symbol name and an evaluated symbol.
type LLVMOrcCSymbolMapPair struct {
	Name LLVMOrcSymbolStringPoolEntryRef
	Sym  LLVMJITEvaluatedSymbol
}

// LLVMOrcCreateLLJITBuilder Create an LLVMOrcLLJITBuilder.
// The client owns the resulting LLJITBuilder and should dispose of it using
// LLVMOrcDisposeLLJITBuilder once they are done with it.
func LLVMOrcCreateLLJITBuilder() LLVMOrcLLJITBuilderRef {
	return LLVMOrcLLJITBuilderRef{c: C.LLVMOrcCreateLLJITBuilder()}
}

// LLVMOrcDisposeLLJITBuilder Dispose of an LLVMOrcLLJITBuilderRef.
func LLVMOrcDisposeLLJITBuilder(b LLVMOrcLLJITBuilderRef) {
	C.LLVMOrcDisposeLLJITBuilder(b.c)
}

// LLVMOrcLLJITBuilderSetJITTargetMachineBuilder Set the JITTargetMachineBuilder
// to be used when constructing the LLJIT instance. Takes ownership of the JTMB argument.
func LLVMOrcLLJITBuilderSetJITTargetMachineBuilder(b LLVMOrcLLJITBuilderRef, jtmb LLVMOrcJITTargetMachineBuilderRef) {
	C.LLVMOrcLLJITBuilderSetJITTargetMachineBuilder(b.c, jtmb.c)
}

// LLVMOrcCreateLLJIT Create an LLJIT instance from an LLJITBuilder.
// Takes ownership of the Builder argument even if the function returns an error.
func LLVMOrcCreateLLJIT(builder LLVMOrcLLJITBuilderRef) (LLVMOrcLLJITRef, error) {
	var out LLVMOrcLLJITRef
	err := orcError2Error(C.LLVMOrcCreateLLJIT(&out.c, builder.c))
	if err != nil {
		return LLVMOrcLLJITRef{}, err
	}
	return out, nil
}

// LLVMOrcDisposeLLJIT Dispose of an LLJIT instance.
func LLVMOrcDisposeLLJIT(j LLVMOrcLLJITRef) error {
	return orcError2Error(C.LLVMOrcDisposeLLJIT(j.c))
}

// LLVMOrcLLJITGetExecutionSession Returns a non-owning reference to the LLJIT instance's execution session.
func LLVMOrcLLJITGetExecutionSession(j LLVMOrcLLJITRef) LLVMOrcExecutionSessionRef {
	return LLVMOrcExecutionSessionRef{c: C.LLVMOrcLLJITGetExecutionSession(j.c)}
}

// LLVMOrcLLJITGetMainJITDylib Returns the main JITDylib for this LLJIT instance.
func LLVMOrcLLJITGetMainJITDylib(j LLVMOrcLLJITRef) LLVMOrcJITDylibRef {
	return LLVMOrcJITDylibRef{c: C.LLVMOrcLLJITGetMainJITDylib(j.c)}
}

// LLVMOrcLLJITGetTripleString Returns the target triple string for this LLJIT instance.
func LLVMOrcLLJITGetTripleString(j LLVMOrcLLJITRef) string {
	return C.GoString(C.LLVMOrcLLJITGetTripleString(j.c))
}

// LLVMOrcLLJITGetDataLayoutStr Get the LLJIT instance's default data layout string.
func LLVMOrcLLJITGetDataLayoutStr(j LLVMOrcLLJITRef) string {
	return C.GoString(C.LLVMOrcLLJITGetDataLayoutStr(j.c))
}

// LLVMOrcLLJITMangleAndIntern Mangle the given symbol name and intern it in the execution session.
func LLVMOrcLLJITMangleAndIntern(j LLVMOrcLLJITRef, name string) LLVMOrcSymbolStringPoolEntryRef {
	return string2CString(name, func(name *C.char) LLVMOrcSymbolStringPoolEntryRef {
		return LLVMOrcSymbolStringPoolEntryRef{c: C.LLVMOrcLLJITMangleAndIntern(j.c, name)}
	})
}

// LLVMOrcLLJITAddLLVMIRModule Add an IR module to the given JITDylib.
// Transfers ownership of the TSM argument to the LLJIT instance.
func LLVMOrcLLJITAddLLVMIRModule(j LLVMOrcLLJITRef, jd LLVMOrcJITDylibRef, tsm LLVMOrcThreadSafeModuleRef) error {
	return orcError2Error(C.LLVMOrcLLJITAddLLVMIRModule(j.c, jd.c, tsm.c))
}

// LLVMOrcLLJITAddObjectFile Add a buffer representing an object file to the given JITDylib.
// Transfers ownership of the buffer to the LLJIT instance.
func LLVMOrcLLJITAddObjectFile(j LLVMOrcLLJITRef, jd LLVMOrcJITDylibRef, buf LLVMMemoryBufferRef) error {
	return orcError2Error(C.LLVMOrcLLJITAddObjectFile(j.c, jd.c, buf.c))
}

// LLVMOrcLLJITLookup Look up the given symbol in the main JITDylib of the given LLJIT instance.
func LLVMOrcLLJITLookup(j LLVMOrcLLJITRef, name string) (unsafe.Pointer, error) {
	var addr unsafe.Pointer
	err := string2CString(name, func(name *C.char) error {
		return orcError2Error(C.LLVMOrcLLJITLookup(j.c, (*C.LLVMOrcExecutorAddress)(unsafe.Pointer(&addr)), name))
	})
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// LLVMOrcJITTargetMachineBuilderDetectHost Create a JITTargetMachineBuilder by detecting the host.
func LLVMOrcJITTargetMachineBuilderDetectHost() (LLVMOrcJITTargetMachineBuilderRef, error) {
	var out LLVMOrcJITTargetMachineBuilderRef
	err := orcError2Error(C.LLVMOrcJITTargetMachineBuilderDetectHost(&out.c))
	if err != nil {
		return LLVMOrcJITTargetMachineBuilderRef{}, err
	}
	return out, nil
}

// LLVMOrcJITTargetMachineBuilderCreateFromTargetMachine Create a JITTargetMachineBuilder from the given TargetMachine.
func LLVMOrcJITTargetMachineBuilderCreateFromTargetMachine(tm LLVMTargetMachineRef) LLVMOrcJITTargetMachineBuilderRef {
	return LLVMOrcJITTargetMachineBuilderRef{c: C.LLVMOrcJITTargetMachineBuilderCreateFromTargetMachine(tm.c)}
}

// LLVMOrcDisposeJITTargetMachineBuilder Dispose of a JITTargetMachineBuilder.
func LLVMOrcDisposeJITTargetMachineBuilder(jtmb LLVMOrcJITTargetMachineBuilderRef) {
	C.LLVMOrcDisposeJITTargetMachineBuilder(jtmb.c)
}

// LLVMOrcCreateNewThreadSafeContext Create a ThreadSafeContextRef containing a new LLVMContext.
func LLVMOrcCreateNewThreadSafeContext() LLVMOrcThreadSafeContextRef {
	return LLVMOrcThreadSafeContextRef{c: C.LLVMOrcCreateNewThreadSafeContext()}
}

// LLVMOrcCreateNewThreadSafeContextFromLLVMContext Create a ThreadSafeContextRef
// from a given LLVMContext, which must not be associated with any existing
// ThreadSafeContext. The underlying ThreadSafeContext takes ownership of the
// LLVMContext object.
func LLVMOrcCreateNewThreadSafeContextFromLLVMContext(ctx LLVMContextRef) LLVMOrcThreadSafeContextRef {
	return LLVMOrcThreadSafeContextRef{c: C.LLVMOrcCreateNewThreadSafeContextFromLLVMContext(ctx.c)}
}

// LLVMOrcDisposeThreadSafeContext Dispose of a ThreadSafeContext.
func LLVMOrcDisposeThreadSafeContext(tsctx LLVMOrcThreadSafeContextRef) {
	C.LLVMOrcDisposeThreadSafeContext(tsctx.c)
}

// LLVMOrcCreateNewThreadSafeModule Create a ThreadSafeModule wrapper around the
// given LLVM module. This takes ownership of the M argument.
func LLVMOrcCreateNewThreadSafeModule(m LLVMModuleRef, tsctx LLVMOrcThreadSafeContextRef) LLVMOrcThreadSafeModuleRef {
	return LLVMOrcThreadSafeModuleRef{c: C.LLVMOrcCreateNewThreadSafeModule(m.c, tsctx.c)}
}

// LLVMOrcDisposeThreadSafeModule Dispose of a ThreadSafeModule.
func LLVMOrcDisposeThreadSafeModule(tsm LLVMOrcThreadSafeModuleRef) {
	C.LLVMOrcDisposeThreadSafeModule(tsm.c)
}

// LLVMOrcExecutionSessionIntern Intern the given name in the execution session's symbol string pool.
func LLVMOrcExecutionSessionIntern(es LLVMOrcExecutionSessionRef, name string) LLVMOrcSymbolStringPoolEntryRef {
	return string2CString(name, func(name *C.char) LLVMOrcSymbolStringPoolEntryRef {
		return LLVMOrcSymbolStringPoolEntryRef{c: C.LLVMOrcExecutionSessionIntern(es.c, name)}
	})
}

// LLVMOrcReleaseSymbolStringPoolEntry Reduces the ref-count for of a SymbolStringPool entry.
func LLVMOrcReleaseSymbolStringPoolEntry(s LLVMOrcSymbolStringPoolEntryRef) {
	C.LLVMOrcReleaseSymbolStringPoolEntry(s.c)
}

// LLVMOrcAbsoluteSymbols Create a MaterializationUnit to define the given
// symbols as pointing to the corresponding raw addresses.
// Takes ownership of the elements of the Syms array.
func LLVMOrcAbsoluteSymbols(syms []LLVMOrcCSymbolMapPair) LLVMOrcMaterializationUnitRef {
	if len(syms) == 0 {
		return LLVMOrcMaterializationUnitRef{c: C.LLVMOrcAbsoluteSymbols(nil, 0)}
	}
	csyms := make([]C.LLVMOrcCSymbolMapPair, len(syms))
	for i, sym := range syms {
		csyms[i] = C.LLVMOrcCSymbolMapPair{
			Name: sym.Name.c,
			Sym: C.LLVMJITEvaluatedSymbol{
				Address: C.LLVMOrcExecutorAddress(sym.Sym.Address),
				Flags: C.LLVMJITSymbolFlags{
					GenericFlags: C.uint8_t(sym.Sym.Flags.GenericFlags),
					TargetFlags:  C.uint8_t(sym.Sym.Flags.TargetFlags),
				},
			},
		}
	}
	return LLVMOrcMaterializationUnitRef{c: C.LLVMOrcAbsoluteSymbols(&csyms[0], C.size_t(len(csyms)))}
}

// LLVMOrcJITDylibDefine Add the given MaterializationUnit to the given JITDylib.
// On success the JITDylib takes ownership of MU.
func LLVMOrcJITDylibDefine(jd LLVMOrcJITDylibRef, mu LLVMOrcMaterializationUnitRef) error {
	return orcError2Error(C.LLVMOrcJITDylibDefine(jd.c, mu.c))
}

// LLVMOrcDisposeMaterializationUnit Dispose of a MaterializationUnit.
func LLVMOrcDisposeMaterializationUnit(mu LLVMOrcMaterializationUnitRef) {
	C.LLVMOrcDisposeMaterializationUnit(mu.c)
}

// LLVMOrcJITDylibCreateResourceTracker Return a reference to a newly created resource tracker.
func LLVMOrcJITDylibCreateResourceTracker(jd LLVMOrcJITDylibRef) LLVMOrcResourceTrackerRef {
	return LLVMOrcResourceTrackerRef{c: C.LLVMOrcJITDylibCreateResourceTracker(jd.c)}
}

// LLVMOrcReleaseResourceTracker Reduces the ref-count of a ResourceTracker.
func LLVMOrcReleaseResourceTracker(rt LLVMOrcResourceTrackerRef) {
	C.LLVMOrcReleaseResourceTracker(rt.c)
}

// orcError2Error converts a non-null LLVMErrorRef to a Go error, consuming it.
func orcError2Error(err C.LLVMErrorRef) error {
	if err == nil {
		return nil
	}
	msg := LLVMGetErrorMessage(LLVMErrorRef{c: err})
	return &orcError{msg: msg}
}

type orcError struct{ msg string }

func (e *orcError) Error() string { return e.msg }

// LLVMOrcJITDylibClear Calls remove on all trackers associated with this JITDylib.
func LLVMOrcJITDylibClear(jd LLVMOrcJITDylibRef) error {
	return orcError2Error(C.LLVMOrcJITDylibClear(jd.c))
}

// LLVMOrcJITDylibGetDefaultResourceTracker Return the default resource tracker for the JITDylib.
// 注意：LLVM 22 的实现未按文档增加引用计数（缺 Retain），C 侧句柄的 release 会破坏
// JITDylib 自身持有的引用，故公开层不暴露该句柄；卸载全部符号请用 LLVMOrcJITDylibClear。
func LLVMOrcJITDylibGetDefaultResourceTracker(jd LLVMOrcJITDylibRef) LLVMOrcResourceTrackerRef {
	return LLVMOrcResourceTrackerRef{c: C.LLVMOrcJITDylibGetDefaultResourceTracker(jd.c)}
}

// LLVMOrcResourceTrackerRemove Remove all symbols tracked by the resource tracker (unload).
func LLVMOrcResourceTrackerRemove(rt LLVMOrcResourceTrackerRef) error {
	return orcError2Error(C.LLVMOrcResourceTrackerRemove(rt.c))
}

// LLVMOrcResourceTrackerTransferTo Transfer ownership of tracked symbols to another tracker.
func LLVMOrcResourceTrackerTransferTo(src, dst LLVMOrcResourceTrackerRef) {
	C.LLVMOrcResourceTrackerTransferTo(src.c, dst.c)
}

// LLVMOrcLLJITAddLLVMIRModuleWithRT Add an IR module to the LLJIT under the given resource tracker.
func LLVMOrcLLJITAddLLVMIRModuleWithRT(j LLVMOrcLLJITRef, rt LLVMOrcResourceTrackerRef, tsm LLVMOrcThreadSafeModuleRef) error {
	return orcError2Error(C.LLVMOrcLLJITAddLLVMIRModuleWithRT(j.c, rt.c, tsm.c))
}

// LLVMOrcLLJITAddObjectFileWithRT Add an object file to the LLJIT under the given resource tracker.
func LLVMOrcLLJITAddObjectFileWithRT(j LLVMOrcLLJITRef, rt LLVMOrcResourceTrackerRef, buf LLVMMemoryBufferRef) error {
	return orcError2Error(C.LLVMOrcLLJITAddObjectFileWithRT(j.c, rt.c, buf.c))
}
