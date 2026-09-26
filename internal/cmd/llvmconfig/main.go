// llvmconfig 从本机 llvm-config 生成 internal/binding/cgo.go 的 #cgo flags。
//
// 用法：
//
//	go run ./internal/cmd/llvmconfig          # 用探测到的 llvm-config 重写 cgo.go（机器专用版）
//	go run ./internal/cmd/llvmconfig --check  # 只报告探测与查询结果，不写文件
//
// 探测顺序：$LLVM_CONFIG > $LLVM_PREFIX/bin/llvm-config > PATH 中的 llvm-config-22
// > PATH 中的 llvm-config。生成器不做 LLVM 大版本校验：所选 llvm-config 报告的
// includedir/libdir/libs 原样采用，版本是否合适由使用者自行判断。
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// modulePath 是本仓库的 module 路径，用于从 cwd 向上定位 go.mod。
const modulePath = "github.com/kkkunny/go-llvm"

// cgoFile 是相对模块根目录的目标文件。
const cgoFile = "internal/binding/cgo.go"

// CFLAGS/CXXFLAGS 的公共部分：与 A 层候选版保持一致。
// CXXFLAGS 不采用 llvm-config --cxxflags 原文：其中可能带 -fno-exceptions，
// 而 C++ shim 必须开启异常（见 Core.cpp 的 extern "C" 边界约定）。
const (
	commonCFlags   = "-D_GNU_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D__STDC_LIMIT_MACROS"
	commonCXXFlags = "-std=c++17 -fexceptions -D_GNU_SOURCE -D_GLIBCXX_USE_CXX11_ABI=1 -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D__STDC_LIMIT_MACROS"
)

// llvmConfig 保存一次 llvm-config 查询的全部结果。
type llvmConfig struct {
	binary     string // llvm-config 可执行文件路径
	version    string // --version
	includeDir string // --includedir
	libDir     string // --libdir
	libs       string // --libs，例如 -lLLVM-22
	systemLibs string // --system-libs，可能为空
}

// options 汇总 run 的注入点，便于测试替换 cwd/env/命令执行与输出。
type options struct {
	check    bool   // 只报告，不写文件
	cwd      string // 启动目录，用于向上查找 module 根
	getenv   func(string) string
	lookPath func(string) (string, error)
	runCmd   func(name string, args ...string) (string, error)
	stdout   io.Writer
}

func main() {
	check := flag.Bool("check", false, "只报告 llvm-config 探测与查询结果，不重写 cgo.go")
	flag.Parse()
	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "llvmconfig: 不接受位置参数：%v\n", flag.Args())
		os.Exit(2)
	}
	cwd, err := os.Getwd()
	if err != nil {
		fatal(err)
	}
	err = run(options{
		check:    *check,
		cwd:      cwd,
		getenv:   os.Getenv,
		lookPath: exec.LookPath,
		runCmd:   execCommand,
		stdout:   os.Stdout,
	})
	if err != nil {
		fatal(err)
	}
}

// fatal 按 CLI 约定把错误写到 stderr 并以退出码 1 结束。
func fatal(err error) {
	fmt.Fprintf(os.Stderr, "llvmconfig: %v\n", err)
	os.Exit(1)
}

// run 执行完整的「发现 → 查询 → 报告或重写」流程；--check 在查询后直接报告。
func run(opts options) error {
	binary, err := discoverLLVMConfig(opts.getenv, opts.lookPath)
	if err != nil {
		return err
	}
	cfg, err := queryLLVMConfig(opts.runCmd, binary)
	if err != nil {
		return err
	}
	if opts.check {
		reportCheck(opts.stdout, cfg)
		return nil
	}
	root, err := findModuleRoot(opts.cwd)
	if err != nil {
		return err
	}
	target := filepath.Join(root, filepath.FromSlash(cgoFile))
	if err := os.WriteFile(target, renderCgo(cfg), 0o644); err != nil {
		return fmt.Errorf("写入 %s: %w", target, err)
	}
	fmt.Fprintf(opts.stdout, "wrote %s (llvm-config: %s, version: %s)\n", target, cfg.binary, cfg.version)
	return nil
}

// discoverLLVMConfig 按固定优先级选择 llvm-config 可执行文件。
// 显式设置了 $LLVM_CONFIG / $LLVM_PREFIX 但不可用时直接报错，避免悄悄退回 PATH 上
// 可能版本不同的工具链；自动探测不校验大版本。
func discoverLLVMConfig(getenv func(string) string, lookPath func(string) (string, error)) (string, error) {
	if p := getenv("LLVM_CONFIG"); p != "" {
		found, err := lookPath(p)
		if err != nil {
			return "", fmt.Errorf("$LLVM_CONFIG=%s 不可执行：%w", p, err)
		}
		return found, nil
	}
	if prefix := getenv("LLVM_PREFIX"); prefix != "" {
		candidate := filepath.Join(prefix, "bin", "llvm-config")
		found, err := lookPath(candidate)
		if err != nil {
			return "", fmt.Errorf("$LLVM_PREFIX=%s 下没有可执行的 bin/llvm-config：%w", prefix, err)
		}
		return found, nil
	}
	for _, name := range []string{"llvm-config-22", "llvm-config"} {
		if found, err := lookPath(name); err == nil {
			return found, nil
		}
	}
	return "", errors.New("找不到 llvm-config：已尝试 $LLVM_CONFIG、$LLVM_PREFIX/bin/llvm-config、" +
		"PATH 中的 llvm-config-22 与 llvm-config。请安装 LLVM 22 开发包（Debian/Ubuntu: llvm-22-dev；" +
		"Arch: llvm；macOS: brew install llvm@22），或用 LLVM_CONFIG=/path/to/llvm-config 明确指定；" +
		"无法安装时也可用 CGO_CFLAGS/CGO_CXXFLAGS/CGO_LDFLAGS 手动覆盖编译链接 flags")
}

// queryLLVMConfig 依次执行 llvm-config 的查询参数并整理结果。
func queryLLVMConfig(runCmd func(name string, args ...string) (string, error), binary string) (llvmConfig, error) {
	cfg := llvmConfig{binary: binary}
	queries := []struct {
		arg string
		dst *string
	}{
		{"--version", &cfg.version},
		{"--includedir", &cfg.includeDir},
		{"--libdir", &cfg.libDir},
		{"--libs", &cfg.libs},
		{"--system-libs", &cfg.systemLibs},
	}
	for _, q := range queries {
		out, err := runCmd(binary, q.arg)
		if err != nil {
			return llvmConfig{}, fmt.Errorf("%s %s: %w", binary, q.arg, err)
		}
		*q.dst = strings.TrimSpace(out)
	}
	if cfg.includeDir == "" || cfg.libDir == "" {
		return llvmConfig{}, fmt.Errorf("%s 未返回 includedir/libdir，无法生成 flags", binary)
	}
	if cfg.libs == "" {
		return llvmConfig{}, fmt.Errorf("%s --libs 为空，无法生成链接 flags", binary)
	}
	return cfg, nil
}

// reportCheck 打印 --check 报告：只搬运 llvm-config 的自报信息，不判断版本适用性。
func reportCheck(w io.Writer, cfg llvmConfig) {
	fmt.Fprintf(w, "llvm-config: %s\n", cfg.binary)
	fmt.Fprintf(w, "version: %s\n", cfg.version)
	fmt.Fprintf(w, "includedir: %s\n", cfg.includeDir)
	fmt.Fprintf(w, "libdir: %s\n", cfg.libDir)
	fmt.Fprintf(w, "libs: %s\n", cfg.libs)
	fmt.Fprintf(w, "system-libs: %s\n", orNone(cfg.systemLibs))
}

// orNone 把空字符串显示为 (none)，避免报告里出现看不清的空值。
func orNone(s string) string {
	if s == "" {
		return "(none)"
	}
	return s
}

// findModuleRoot 从 start 目录向上查找 module 为本仓库的 go.mod，返回模块根目录。
func findModuleRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if err == nil && goModModule(data) == modulePath {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("在 %s 及其父目录中找不到 module %s 的 go.mod", start, modulePath)
		}
		dir = parent
	}
}

// goModModule 返回 go.mod 中 module 指令的路径，没有则返回空串。
func goModModule(data []byte) string {
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "module" {
			return fields[1]
		}
	}
	return ""
}

// renderCgo 渲染机器专用版 cgo.go。
// CFLAGS 与 CXXFLAGS 对称携带 -I<includedir>：C++ shim 同样包含 LLVM 头文件。
// //go:generate 保留在生成结果里，保证 cgo.go 被重写后 `go generate` 仍可再次运行。
func renderCgo(cfg llvmConfig) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "// Code generated by internal/cmd/llvmconfig from %s (%s); DO NOT EDIT.\n", cfg.binary, cfg.version)
	b.WriteString("//\n")
	b.WriteString("// This machine-specific version replaces the portable per-OS candidate list that is\n")
	b.WriteString("// checked in by default. Re-run `go generate ./internal/binding` or `make config`\n")
	b.WriteString("// after upgrading LLVM; do not commit it unless the whole team shares the layout.\n")
	b.WriteString("//\n")
	b.WriteString("//go:generate go run ../cmd/llvmconfig\n")
	b.WriteString("\n")
	b.WriteString("package binding\n")
	b.WriteString("\n")
	b.WriteString("/*\n")
	fmt.Fprintf(&b, "#cgo CFLAGS: %s\n", commonCFlags)
	fmt.Fprintf(&b, "#cgo CFLAGS: -I%s\n", cfg.includeDir)
	fmt.Fprintf(&b, "#cgo CXXFLAGS: %s\n", commonCXXFlags)
	fmt.Fprintf(&b, "#cgo CXXFLAGS: -I%s\n", cfg.includeDir)
	fmt.Fprintf(&b, "#cgo LDFLAGS: %s\n", linkFlags(cfg))
	b.WriteString("*/\n")
	b.WriteString("import \"C\"\n")
	return []byte(b.String())
}

// linkFlags 拼出 LDFLAGS：-L<libdir> 后接 --libs 与 --system-libs（后者可为空）。
func linkFlags(cfg llvmConfig) string {
	flags := []string{"-L" + cfg.libDir, cfg.libs}
	if cfg.systemLibs != "" {
		flags = append(flags, cfg.systemLibs)
	}
	return strings.Join(flags, " ")
}

// execCommand 运行外部命令并返回其标准输出；失败时附带 stderr 便于定位。
func execCommand(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && len(exitErr.Stderr) > 0 {
			return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(string(exitErr.Stderr)))
		}
		return "", err
	}
	return string(out), nil
}
