package llvm

import (
	"runtime"

	"github.com/kkkunny/go-llvm/internal/binding"
)

// MemoryBuffer 内存缓冲；IR 解析、bitcode 读写、目标代码产出共用，用毕 Close
type MemoryBuffer struct {
	ref    binding.LLVMMemoryBufferRef
	closed bool
}

// ReadFile 读取文件到内存缓冲；失败返回 ErrIO
func ReadFile(path string) (*MemoryBuffer, error) {
	ref, err := binding.LLVMCreateMemoryBufferWithContentsOfFile(path)
	if err != nil {
		return nil, WrapError(ErrIO, "llvm.ReadFile", err)
	}
	return MemoryBufferOf(ref), nil
}

// NewMemoryBuffer 由字节切片创建内存缓冲（拷贝数据）
func NewMemoryBuffer(data []byte, name string) *MemoryBuffer {
	return MemoryBufferOf(binding.LLVMCreateMemoryBufferWithMemoryRangeCopy(data, name))
}

// MemoryBufferOf 由底层句柄构建内存缓冲（供 llvm/* 子包桥接使用，接管所有权）
func MemoryBufferOf(ref binding.LLVMMemoryBufferRef) *MemoryBuffer {
	b := &MemoryBuffer{ref: ref}
	runtime.SetFinalizer(b, (*MemoryBuffer).finalize)
	return b
}

// finalize GC 兜底：忘记 Close 时释放底层缓冲
func (b *MemoryBuffer) finalize() {
	if b.closed {
		return
	}
	b.closed = true
	binding.LLVMDisposeMemoryBuffer(b.ref)
}

// Ref 返回底层句柄（供 llvm/* 子包桥接使用）
func (b *MemoryBuffer) Ref() binding.LLVMMemoryBufferRef { return b.ref }

// Alive 缓冲是否可用（未释放）
func (b *MemoryBuffer) Alive() bool { return !b.ref.IsNil() && !b.closed }

// check 前置校验：句柄非空且未释放
func (b *MemoryBuffer) check(op string) {
	if b.ref.IsNil() {
		errPanic(ErrInvalidArg, op, "nil memory buffer")
	}
	if b.closed {
		errPanic(ErrUseAfterFree, op, "memory buffer is closed")
	}
}

// Bytes 缓冲内容视图（零拷贝）；Close 后失效
func (b *MemoryBuffer) Bytes() []byte {
	b.check("llvm.MemoryBuffer.Bytes")
	return binding.LLVMGetBufferStart(b.ref)
}

// Len 缓冲字节数
func (b *MemoryBuffer) Len() int {
	b.check("llvm.MemoryBuffer.Len")
	return int(binding.LLVMGetBufferSize(b.ref))
}

// Disown 移交底层缓冲给外部接管方（如 JIT）：此后 Go 侧句柄失效，Close 返回 ErrClosed；
// 底层缓冲由接管方释放。
func (b *MemoryBuffer) Disown() {
	if b.closed {
		return
	}
	b.closed = true
	runtime.SetFinalizer(b, nil) // 接管方负责释放，禁止 finalizer 二次释放
}

// Close 释放缓冲；二次调用返回 ErrClosed
func (b *MemoryBuffer) Close() error {
	if b.closed {
		return &Error{Reason: ErrClosed, Op: "llvm.MemoryBuffer.Close", Msg: "memory buffer already closed"}
	}
	b.closed = true
	runtime.SetFinalizer(b, nil)
	binding.LLVMDisposeMemoryBuffer(b.ref)
	return nil
}
