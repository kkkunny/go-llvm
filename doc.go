// Package llvm provides Go bindings to a system-installed LLVM.
//
// This package holds the core vocabulary of the library: kinds, types, values, constants,
// contexts, errors and lifetimes. IR construction, code generation, execution and
// optimization live in the sub-packages.
//
// # Package layout
//
//   - llvm: the core vocabulary — [Kind], [Type] and [Value], constants, [Context], error
//     values and lifetimes.
//   - [github.com/kkkunny/go-llvm/ir]: IR construction and inspection —
//     [github.com/kkkunny/go-llvm/ir.Module], functions, blocks, builder and instructions.
//   - [github.com/kkkunny/go-llvm/target]: target machines and code generation.
//   - [github.com/kkkunny/go-llvm/jit]: ORC LLJIT execution engine.
//   - [github.com/kkkunny/go-llvm/pass]: optimization pipelines.
//
// # Type safety
//
// Types and values are expressed with kind-level generics, [Type] and [Value]: the type
// parameter distinguishes categories such as integer, float, pointer, struct or function,
// never bit widths or element types. Category-specific operations live on role wrappers that
// embed the generic view, such as [IntType.Bits] or
// [github.com/kkkunny/go-llvm/ir.Alloca.SetAlign]. [TypeRef] and [ValueRef] are implemented
// by both the generic and the role form, so generic APIs accept either.
//
// # Error handling
//
// Recoverable runtime failures return an error. Programmer errors panic with an [Error];
// recover such panics with [Catch], or use [Try] when a value is produced.
package llvm
