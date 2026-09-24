//go:build llvm_release

package checks

// debugBuild 信任构建（-tags=llvm_release）：关闭调试层，仅保留崩溃类地板。
const debugBuild = false
