package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// SetFatalErrorHandler 安装 fatal error 回调。
//
// 注意：LLVM 的 report_fatal_error 语义是打印诊断后退出进程，回调只会在退出前被调用，
// 无法通过 recover 阻止退出。库的主要防线是前置校验与 C++ shim 的异常捕获，
// 该回调用于在退出前记录现场。
func SetFatalErrorHandler(handler func(string)) {
	binding.LLVMInstallFatalErrorHandlerGo(handler)
}

// ResetFatalErrorHandler 恢复 LLVM 默认的 fatal error 行为
func ResetFatalErrorHandler() {
	binding.LLVMResetFatalErrorHandler()
}

// EnablePrettyStackTrace 启用 LLVM 内建堆栈跟踪。
//
// 警告：该 API 会安装 LLVM 崩溃信号处理器，与 Go runtime 用于栈增长的 SIGSEGV 冲突；
// 实测调用后进程会在栈增长时崩溃（runtime: split stack overflow），普通 Go 程序与
// 测试中不要调用。
func EnablePrettyStackTrace() {
	binding.LLVMEnablePrettyStackTrace()
}
