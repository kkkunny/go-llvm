package llvm

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm/internal/binding"
)

// TestLinkedLLVMVersion 冒烟校验运行时链接的 LLVM 库大版本与编译期头文件一致。
// 该比对是 NewContext 调试构建版本校验（checkLinkedVersion）的门槛；其余测试都会
// 经过 NewContext，因此同时覆盖了校验通过的正常路径。
func TestLinkedLLVMVersion(t *testing.T) {
	major, minor, patch := binding.LLVMGetVersion()
	if want := uint32(binding.LLVM_VERSION_MAJOR); major != want {
		t.Fatalf("运行时 LLVM 库为 %d.%d.%d，编译期头文件为 %s（major %d）",
			major, minor, patch, binding.LLVM_VERSION_STRING, want)
	}
}

// TestCheckVersionMatchMismatch 覆盖版本错配路径：运行时 major 与头文件不一致时
// panic [ErrVersionMismatch]，消息包含两侧版本信息。校验只在调试构建启用，故用 requireDebug 门控。
func TestCheckVersionMatchMismatch(t *testing.T) {
	requireDebug(t)
	err := Catch(func() {
		checkVersionMatch(func() (uint32, uint32, uint32) { return 21, 1, 0 }, 22, "22.1.8")
	})
	if err == nil {
		t.Fatal("版本错配时应 panic")
	}
	if err.Reason != ErrVersionMismatch || err.Op != "llvm.NewContext" {
		t.Fatalf("want ErrVersionMismatch/llvm.NewContext, got %+v", err)
	}
	if !strings.Contains(err.Msg, "21.1.0") || !strings.Contains(err.Msg, "22.1.8") {
		t.Fatalf("消息应包含运行时版本与头文件版本：%s", err.Msg)
	}
}
