package llvm

import "github.com/kkkunny/go-llvm/internal/errs"

// ErrKind 错误类别
type ErrKind = errs.Kind

const (
	ErrTypeMismatch    = errs.ErrTypeMismatch    // 操作数/种类不匹配
	ErrCrossContext    = errs.ErrCrossContext    // 跨 Context 混用
	ErrUseAfterFree    = errs.ErrUseAfterFree    // 使用已释放资源
	ErrClosed          = errs.ErrClosed          // 重复 Close
	ErrNotFound        = errs.ErrNotFound        // 查找失败
	ErrInvalidArg      = errs.ErrInvalidArg      // 参数非法
	ErrVerify          = errs.ErrVerify          // IR 验证失败
	ErrParse           = errs.ErrParse           // IR/bitcode 解析失败（ir.ParseIR/ParseBitcode）
	ErrUnsupported     = errs.ErrUnsupported     // 映射遇到不支持的类型
	ErrCodeGen         = errs.ErrCodeGen         // 目标代码生成失败（target 包）
	ErrJIT             = errs.ErrJIT             // JIT 构造/符号解析失败（jit 包）
	ErrIO              = errs.ErrIO              // 文件/内存缓冲读写失败（ir/target/jit）
	ErrLink            = errs.ErrLink            // 模块链接失败（ir.Module.Link）
	ErrPass            = errs.ErrPass            // 优化管线执行失败（pass 包）
	ErrInternal        = errs.ErrInternal        // 兜底
	ErrVersionMismatch = errs.ErrVersionMismatch // 运行时 LLVM 库与编译期头文件版本不一致
)

// Error Go 化的 LLVM 错误；程序员错误经 panic 抛出，可在调用方 recover 后用 errors.As 还原为 *Error。
//
// 错误类别的构造、包装与收敛机制位于 internal/errs，不属于公开 API；
// 本包仅以别名形式保留 [Error] 与 [ErrKind] 供调用方检查错误。
type Error = errs.Error
