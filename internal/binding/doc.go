// Package binding is the only cgo layer in the library. It mirrors LLVM-C one-to-one:
// wrappers keep the exact C API names and signatures, and enum constants are declared as the
// local header names so new members are picked up automatically.
//
// APIs that LLVM-C does not expose are filled in by a small C++ shim (Core.h/Core.cpp).
// C++ exceptions must be caught at the extern "C" boundary and reported as errors, never
// propagated into cgo.
package binding
