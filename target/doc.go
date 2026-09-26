// Package target registers LLVM targets and lowers IR to assembly or object code.
//
// # Initialization
//
// Code generation requires an initialized target: [InitNative] registers the host backend,
// [Init] registers one architecture and [InitAll] registers every available backend.
// [NativeTarget] returns the host [Target] and [DefaultTriple] its triple; [FromTriple] and
// [FromName] look targets up by triple or name.
//
// # Code generation
//
// [NewTargetMachine] builds a [TargetMachine] from a target, triple, CPU and features.
// [TargetMachine.EmitToFile] writes assembly or an object file ([AsmFile], [ObjectFile]) to
// disk, and [TargetMachine.Emit] returns the same as a
// [github.com/kkkunny/go-llvm.MemoryBuffer]. Before emitting, call [TargetMachine.ApplyTo]
// to write the machine's triple and data layout into the
// [github.com/kkkunny/go-llvm/ir.Module].
//
// [TargetMachine] is an independent ownership root (not registered with
// [github.com/kkkunny/go-llvm.Context]) and must be released with [TargetMachine.Close].
package target
