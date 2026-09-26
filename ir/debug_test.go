package ir

import (
	"testing"

	"github.com/kkkunny/go-llvm/internal/checks"
)

// requireDebug 语义契约类负向测试仅在调试构建（默认）下执行。
// -tags=llvm_release 时第 3 层校验被编译消除，此类测试会向 LLVM 传入非法 IR，
// 必须在信任构建下跳过；崩溃类地板（判空/已释放/跨 Context）的负向测试不受影响。
func requireDebug(t *testing.T) {
	t.Helper()
	if !checks.Debug {
		t.Skip("semantic contract checks are compiled out in llvm_release builds")
	}
}
