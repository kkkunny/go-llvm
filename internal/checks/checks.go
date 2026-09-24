// Package checks 提供全库统一的校验模式开关。
//
// go-llvm 的校验分三层：
//  1. 编译期种类安全（Kind 泛型，任何模式下都在）；
//  2. 崩溃类地板（判空/已释放/跨 Context/已关闭，纯 Go，任何模式下都在）；
//  3. 调试层（语义契约与诊断增强，仅 Debug 为 true 时启用）。
//
// Debug 由构建约束选择：默认 true（安全默认）；使用 -tags=llvm_release 构建时
// 为 false，第 3 层在编译期被整体消除（const 条件死代码消除），此时语义误用
// 交由 LLVM assert + Verify 兜底，与 C/Rust/inkwell 的 release 语义一致。
package checks

// Debug 是否为调试构建（默认）；-tags=llvm_release 时为 false。
// 所有"仅调试层"的校验与增强必须以 `if checks.Debug { ... }` 包裹。
const Debug = debugBuild
