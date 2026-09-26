package llvm

import "fmt"

// ErrKind 错误类别
type ErrKind int

const (
	ErrTypeMismatch    ErrKind = iota // 操作数/种类不匹配
	ErrCrossContext                   // 跨 Context 混用
	ErrUseAfterFree                   // 使用已释放资源
	ErrClosed                         // 重复 Close
	ErrNotFound                       // 查找失败
	ErrInvalidArg                     // 参数非法
	ErrVerify                         // IR 验证失败
	ErrParse                          // IR/bitcode 解析失败（ir.ParseIR/ParseBitcode）
	ErrUnsupported                    // 映射遇到不支持的类型
	ErrCodeGen                        // 目标代码生成失败（target 包）
	ErrJIT                            // JIT 构造/符号解析失败（jit 包）
	ErrIO                             // 文件/内存缓冲读写失败（ir/target/jit）
	ErrLink                           // 模块链接失败（ir.Module.Link）
	ErrPass                           // 优化管线执行失败（pass 包）
	ErrInternal                       // 兜底
	ErrVersionMismatch                // 运行时 LLVM 库与编译期头文件版本不一致
)

// Error Go 化的 LLVM 错误；程序员错误经 panic 抛出，可用 Catch 收敛为 error
type Error struct {
	Reason ErrKind
	Op     string // 如 "llvm.Builder.Add"
	Msg    string
	cause  error // 底层错误链（WrapError 保留），支持 errors.Is/As 追溯
}

// Error 实现 error 接口，返回 "Op: Msg" 形式的错误文本（如 "llvm.Builder.Add: ..."）。
// 文本不含 Reason 字段，也不展开底层 cause；底层错误经 [Error.Unwrap] 追溯。
func (e *Error) Error() string { return fmt.Sprintf("%s: %s", e.Op, e.Msg) }

// Unwrap 返回底层错误（若有）
func (e *Error) Unwrap() error { return e.cause }

// errPanic 构造并 panic（内部用）
func errPanic(reason ErrKind, op, format string, args ...any) {
	panic(&Error{Reason: reason, Op: op, Msg: fmt.Sprintf(format, args...)})
}

// Panicf 构造并 panic *Error（供 llvm/* 子包统一错误类型）
func Panicf(reason ErrKind, op, format string, args ...any) {
	errPanic(reason, op, format, args...)
}

// Catch 执行 fn 并将其 panic 的 *Error 收敛为 error 返回；非 *Error panic 原样重抛；无 panic 返回 nil
func Catch(fn func()) (err *Error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(*Error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()
	fn()
	return nil
}

// Try 执行 fn 并返回其结果；fn panic 的 *Error 收敛为 error 返回（非 *Error 原样重抛）
func Try[T any](fn func() T) (v T, err *Error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(*Error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()
	v = fn()
	return v, nil
}

// Must err 非 nil 时 panic *Error（非 *Error 会经 WrapError 归一，保证 Catch 可收敛），否则返回 v
func Must[T any](v T, err error) T {
	if err != nil {
		if e, ok := err.(*Error); ok {
			panic(e)
		}
		panic(WrapError(ErrInternal, "llvm.Must", err))
	}
	return v
}

// WrapError 把底层错误（如 internal/binding 返回的 error）包装为 *Error。
// err 为 nil 返回 nil；已是 *Error 则原样返回。上层包（ir/target/jit/pass）统一用它归类错误。
func WrapError(reason ErrKind, op string, err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return &Error{Reason: reason, Op: op, Msg: err.Error(), cause: err}
}
