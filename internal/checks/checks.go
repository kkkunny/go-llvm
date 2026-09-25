// Package checks exposes the library-wide switch for the validation layers.
//
// Validation is layered:
//
//  1. Compile-time kind safety (the Kind generics) — always on.
//  2. Crash-class floor (nil, use-after-free, cross-context, closed) — pure Go, always on.
//  3. Debug layer (semantic contracts and extra diagnostics) — only when [Debug] is true.
//
// [Debug] is chosen by build tag: it is true by default and false when building with
// -tags=llvm_release, which deletes layer 3 at compile time (constant-condition dead-code
// elimination). In release builds semantic misuse falls back to LLVM asserts and module
// verification, matching the release semantics of C, Rust and inkwell.
package checks

// Debug 是否为调试构建（默认）；-tags=llvm_release 时为 false。
// 所有"仅调试层"的校验与增强必须以 `if checks.Debug { ... }` 包裹。
const Debug = debugBuild
