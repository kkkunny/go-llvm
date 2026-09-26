package ir

import (
	"runtime"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/checks"
)

// goroutineCheckInterval 单 goroutine 契约采样间隔：每 N 次调用核对一次 owner，
// 避免每次 runtime.Stack 的开销；定位/关闭等关键路径全量核对。
const goroutineCheckInterval = 1024

// curGoroutineID 解析当前 goroutine id（仅调试层使用；解析失败返回 0，此时跳过检查）。
// Go 未暴露官方 API，这里解析 runtime.Stack 头部，属调试层可接受的妥协。
func curGoroutineID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	const prefix = "goroutine "
	if n <= len(prefix) || string(buf[:len(prefix)]) != prefix {
		return 0
	}
	id := uint64(0)
	for _, c := range buf[len(prefix):n] {
		if c < '0' || c > '9' {
			break
		}
		id = id*10 + uint64(c-'0')
	}
	return id
}

// checkOwner 调试层单 goroutine 契约执法（B2）：首次调用记录 owner goroutine，
// 之后按采样（always=false）或全量（always=true）核对，跨 goroutine 使用 panic。
// 该检查补的是 race detector 看不见 C 侧状态的盲区。
func checkOwner(op string, owner *uint64, ops *uint64, what string, always bool) {
	if !checks.Debug {
		return
	}
	if *owner == 0 {
		*owner = curGoroutineID()
		return
	}
	*ops++
	// 关键：先判断采样间隔再取 goroutine id（runtime.Stack ~µs 级，不能每次调用都付）
	if !always && *ops%goroutineCheckInterval != 0 {
		return
	}
	if id := curGoroutineID(); id != *owner {
		llvm.Panicf(llvm.ErrInvalidArg, op,
			"%s used from another goroutine (owner=%d, current=%d); handles are not goroutine-safe", what, *owner, id)
	}
}
