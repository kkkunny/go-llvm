# 重写 internal/binding/cgo.go 的 #cgo 编译链接 flags，用于升级 LLVM 时对齐新工具链。
#
# 用法：
#   make config                    # 探测 MIN/MAX 范围内可用的 llvm-config
#   make config MAJOR_VERSION=23   # 指定工具链
#
# 生成规则：按仓库约定排列多候选目录——/usr/lib/llvm-NN 优先，然后 /usr、/usr/local、
# /usr/lib64；链接沿用 unversioned -lLLVM。非标准前缀请在消费者侧用 CGO_* 环境变量覆盖。
MIN_SUPPORT_MAJOR_VERSION = 22
MAX_SUPPORT_MAJOR_VERSION = 22
LLVM_CONFIG_BIN =

define find_llvm_config_bin
    $(eval TMP := $(shell which llvm-config$(1) 2>/dev/null))
    $(if $(TMP),$(eval LLVM_CONFIG_BIN := llvm-config$(1)))
endef

ifdef MAJOR_VERSION
	ifneq ($(shell which llvm-config$(MAJOR_VERSION) 2>/dev/null),)
		LLVM_CONFIG_BIN = llvm-config$(MAJOR_VERSION)
	else
		LLVM_CONFIG_BIN = llvm-config
	endif
else
	_ := $(foreach i,$(shell seq $(MIN_SUPPORT_MAJOR_VERSION) $(MAX_SUPPORT_MAJOR_VERSION)),$(call find_llvm_config_bin,$(i)))
	ifeq ($(LLVM_CONFIG_BIN),)
		LLVM_CONFIG_BIN = llvm-config
	endif
endif

# 探测到的 llvm-config 实际大版本，决定 /usr/lib/llvm-NN 候选目录
LLVM_MAJOR = $(shell $(LLVM_CONFIG_BIN) --version 2>/dev/null | cut -d. -f1)
LLVM_PREFIX = /usr/lib/llvm-$(LLVM_MAJOR)
CGO_FILE = internal/binding/cgo.go

.PHONY: config test test-release vet bench bench-release

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
	@test -n "$(LLVM_MAJOR)" || { echo "llvm-config not found: $(LLVM_CONFIG_BIN)" >&2; exit 1; }
	@printf 'package binding\n\n/*\n' > $(CGO_FILE)
	@printf '#cgo CFLAGS: -D_GNU_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D__STDC_LIMIT_MACROS\n' >> $(CGO_FILE)
	@printf '#cgo CFLAGS: -I%s/include -I/usr/include -I/usr/local/include\n' '$(LLVM_PREFIX)' >> $(CGO_FILE)
	@printf '#cgo CXXFLAGS: -std=c++17 -fexceptions -D_GNU_SOURCE -D_GLIBCXX_USE_CXX11_ABI=1 -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D__STDC_LIMIT_MACROS\n' >> $(CGO_FILE)
	@printf '#cgo CXXFLAGS: -I%s/include -I/usr/include -I/usr/local/include\n' '$(LLVM_PREFIX)' >> $(CGO_FILE)
	@printf '#cgo LDFLAGS: -L%s/lib -L/usr/lib64 -L/usr/lib -L/usr/local/lib -lLLVM\n' '$(LLVM_PREFIX)' >> $(CGO_FILE)
	@printf '*/\nimport "C"\n' >> $(CGO_FILE)
	@echo "wrote $(CGO_FILE) with $(LLVM_CONFIG_BIN) ($(LLVM_MAJOR))"
