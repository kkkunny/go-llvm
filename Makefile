# 重写 internal/binding/cgo.go 的 #cgo 编译链接 flags，用于本机/非标准前缀的 LLVM 布局。
#
# 用法：
#   make config                                   # 使用生成器探测到的 llvm-config
#   LLVM_CONFIG=/path/to/llvm-config make config  # 指定工具链
#
# 探测与生成逻辑在 internal/cmd/llvmconfig（不做 LLVM 大版本校验）：
#   $LLVM_CONFIG > $LLVM_PREFIX/bin/llvm-config > PATH llvm-config-23 > PATH llvm-config。
# 生成结果是机器专用版；提交前应确认写回可移植候选版，或先 review diff。
#
# 版本 tag 契约（见 AGENTS.md "Version tags"）：不带 tag 构建 master 最新支持的
# LLVM；带 VERSION_TAG（llvmNN）精确锁定该大版本。MIN/MAX/VERSION_TAG 必须与
# README 支持矩阵、CI 矩阵一起升级。
MIN_SUPPORT_MAJOR_VERSION = 23
MAX_SUPPORT_MAJOR_VERSION = 23
VERSION_TAG = llvm23

.PHONY: test test-tag test-release vet bench bench-release config

# 默认（调试）构建：三层校验全开（崩溃类地板 + 语义契约 + 调试增强）
test:
	go test ./...
# 版本 tag 契约构建：显式锁定当前支持的大版本（目前与默认等价）
test-tag:
	go test -tags=$(VERSION_TAG) ./...
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
