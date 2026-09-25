// Package jit executes IR with LLVM's ORC-based LLJIT engine.
//
// [NewLLJIT] creates a JIT for the host target. Modules are added with [LLJIT.AddIRModule]
// (or [ResourceTracker.AddIRModule] to track them for later removal) and symbols are
// resolved lazily through [LLJIT.Lookup], [LLJIT.Func] or [LLJIT.RunMain].
//
// # Ownership
//
// [LLJIT] is an independent ownership root: it is not registered with
// [github.com/kkkunny/go-llvm.Context] and must be closed explicitly. [LLJIT.AddIRModule]
// transfers ownership of the module and of its entire context to the JIT, invalidating both
// Go handles immediately; the JIT releases them on [LLJIT.Close]. [LLJIT.AddObjectFile]
// consumes a [github.com/kkkunny/go-llvm.MemoryBuffer] the same way.
// [ResourceTracker.Remove] unloads its tracked symbols without closing the JIT, and
// [LLJIT.ClearSymbols] unloads every symbol definition.
//
// # Concurrency
//
// [LLJIT.Func], [LLJIT.MapFunc], [LLJIT.Lookup], [LLJIT.MapSymbol] and the resource tracker
// methods may be called from multiple goroutines; the adapter cache and the bridge registry
// are lock-protected. Calls that conflict with [LLJIT.Close] must be serialized by the
// caller.
//
// # Go interop
//
// [LLJIT.MapFunc] registers a host Go function as a JIT symbol: it generates an IR wrapper
// with the real signature that packs arguments into slots and dispatches them through the
// callGoChannel symbol. In the other direction, [LLJIT.Func] wraps a JIT symbol as a real Go
// function value via reflect.MakeFunc, at the cost of reflection and one cgo round trip per
// call. [LLJIT.AddProcessSymbols] exposes host process symbols (libc, libm, ...) to extern
// declarations, and [LLJIT.MapSymbol] defines a single host address as a JIT symbol.
package jit
