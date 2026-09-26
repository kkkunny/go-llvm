# 重写 internal/binding/cgo.go 的 #cgo 编译链接 flags，用于本机/非标准前缀的 LLVM 布局。
#
# 用法：
#   make config                                   # 使用生成器探测到的 llvm-config
#   LLVM_CONFIG=/path/to/llvm-config make config  # 指定工具链
#
# 探测与生成逻辑在 internal/cmd/llvmconfig（不做 LLVM 大版本校验）：
#   $LLVM_CONFIG > $LLVM_PREFIX/bin/llvm-config > PATH llvm-config-22 > PATH llvm-config。
# 生成结果是机器专用版；提交前应确认写回可移植候选版，或先 review diff。
#
# MIN/MAX_SUPPORT_MAJOR_VERSION 仅作为版本标记，供 AGENTS.md / README 支持矩阵引用。
MIN_SUPPORT_MAJOR_VERSION = 22
MAX_SUPPORT_MAJOR_VERSION = 22

.PHONY: test test-release vet bench bench-release config

# 默认（调试）构建：三层校验全开（崩溃类地板 + 语义契约 + 调试增强）
test:
	go test ./...
# 信任构建：语义契约与调试增强在编译期消除，仅保留纯 Go 崩溃类地板
test-release:
	go test -tags=llvm_release ./...
vet:
	go vet ./...
bench:
	go test -run '^$$' -bench . -benchmem ./...
bench-release:
	go test -tags=llvm_release -run '^$$' -bench . -benchmem ./...
config:
	go run ./internal/cmd/llvmconfig
