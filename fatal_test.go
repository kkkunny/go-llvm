package llvm

import "testing"

// TestFatalHandlers 只验证 fatal 回调的安装与复位，不触发 fatal 路径——
// LLVM 的 report_fatal_error 会退出进程，无法通过 recover 阻止。
func TestFatalHandlers(t *testing.T) {
	SetFatalErrorHandler(func(msg string) {})
	// 重复安装/复位应保持幂等
	SetFatalErrorHandler(func(msg string) {})
	ResetFatalErrorHandler()
	ResetFatalErrorHandler()
}

// 注意：EnablePrettyStackTrace 不在测试中调用。该 API 会安装 LLVM 的崩溃信号处理器
// （上游文档：“Enables dumping a pretty stack trace when the program crashes”），
// 它会抢占 Go 运行时用于栈增长的 SIGSEGV，导致此后任意测试触发
// “PLEASE submit a bug report to LLVM … fatal error: runtime: split stack overflow”。
// 该行为已在本地实测复现（PROBE=enable 后跑 ExampleTry 必崩），属上游 API 与 Go 运行时的
// 已知冲突，故不放入进程内测试；详见 task-10 报告。
