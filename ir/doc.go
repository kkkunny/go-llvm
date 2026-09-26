// Package ir constructs, parses and prints LLVM IR.
//
// [Module] is the unit of construction: it owns functions, globals, aliases, comdats and
// metadata. [Function] values contain basic blocks, blocks contain instructions, and
// [Builder] appends instructions at a tracked insert point. Modules round-trip through IR
// text or bitcode: [Module.String], [Module.WriteToFile], [ParseIR] and [ParseBitcode].
//
// # Ownership
//
// A [Module] is rooted in a [github.com/kkkunny/go-llvm.Context]: [NewModule] registers it
// with the context, closing the module invalidates every value created from it, and closing
// the context cascades to its modules in reverse creation order. Values, types, blocks and
// functions carry the module's lifetime token and panic when used after it ends.
// [Module.Disown] hands a module to an external owner (such as the JIT) and immediately
// invalidates the Go handle.
//
// # Verification
//
// Call [Module.Verify] before handing a module to [github.com/kkkunny/go-llvm/target] or
// [github.com/kkkunny/go-llvm/jit]; debug builds verify automatically at those boundaries.
// [Module.String] deliberately does not verify — it is meant for inspecting modules under
// construction, where the IR may not be valid yet.
//
// # Roles and kind safety
//
// Instructions and values are returned as kind-level generics
// ([github.com/kkkunny/go-llvm.Value]); operations specific to a category live on role
// wrappers such as [Alloca], [Call], [Phi] and [Switch], which embed the generic value and
// promote its methods. [github.com/kkkunny/go-llvm.ValueRef] and
// [github.com/kkkunny/go-llvm.TypeRef] are implemented by both forms, so generic APIs accept
// either.
package ir
