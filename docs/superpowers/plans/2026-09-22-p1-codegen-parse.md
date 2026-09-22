# P1-1 目标与代码生成 + P1-2 序列化与解析 实现计划

> **面向 AI 代理的工作者：** 必需子技能：superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans 逐任务实现此计划。步骤使用复选框（`- [ ]`）语法跟踪进度。

**目标：** 落地 `llvm/target`（Target/TargetMachine/emit OBJ·ASM）、root `llvm.DataLayout` 与 `llvm.MemoryBuffer`，并在 `llvm/ir` 完成 IR 文本/bitcode 的解析与写出。

**架构：** `DataLayout`/`MemoryBuffer` 作为多子包共用资源放 root（P0.5 决策）；`llvm/target` 依赖 `llvm/ir`（emit 必须见 `*ir.Module`）；`internal/binding` 继续 1:1 映射 LLVM-C（新增 BitReader/BitWriter/IRReader 与 Core 的 MemoryBuffer 族）。所有 `DataLayout` 均为 owned 句柄（`Module.DataLayout()` 经布局字符串重建副本），避免 `LLVMGetModuleDataLayout` 借用引用的双重释放。

**技术栈：** Go 1.27、cgo、LLVM 22（C API，无新 shim）、标准 `testing`。

**规格：** `docs/superpowers/specs/2026-09-22-go-llvm-inkwell-style-redesign-design.md` §4.4/§4.7、§5 的 P1-1/P1-2。

---

## 文件结构

```
internal/binding/Types.go          修改：LLVMMemoryBufferRef
internal/binding/Core.go           修改：MemoryBuffer 族 + LLVMPrintModuleToFile
internal/binding/BitReader.go      创建：LLVMParseBitcodeInContext
internal/binding/BitWriter.go      创建：LLVMWriteBitcodeToFile / LLVMWriteBitcodeToMemoryBuffer
internal/binding/IRReader.go       创建：LLVMParseIRInContext（包装非 deprecated 的 ...2）
internal/binding/TargetMachine.go  修改：LLVMTargetMachineEmitToMemoryBuffer
error.go                           修改：ErrKind 新增 ErrParse
memorybuffer.go                    创建：MemoryBuffer（ReadFile/NewMemoryBuffer/Bytes/Len/Alive/Close）
datalayout.go                      创建：DataLayout + ByteOrder（全套查询）
ir/parse.go                        创建：ParseIR / ParseBitcode
ir/module.go                       修改：newModule 铸造助手；DataLayout/WriteToFile/WriteBitcode/Bitcode
target/target.go                   创建：Arch/Init/InitAll/InitNative/Target/宿主查询
target/machine.go                  创建：OptLevel/RelocMode/CodeModel/FileType/TargetMachine/emit/SetTo
*_test.go                          各包单元测试
```

**关键决策（实现时不可偏离）：**

1. `ErrKind` 新增 `ErrParse`（IR/bitcode 解析失败），与 `ErrVerify` 区分。
2. 解析 API：`LLVMParseIRInContext2`（非 deprecated，buffer 不被消费）；`LLVMParseBitcodeInContext`（v1 带 `OutMessage`，仅注释级 deprecated、无弃用属性）。两者 buffer 均归调用方 `Close`。
3. `DataLayout` 一律 owned（`NewDataLayout(string)` / `DataLayoutOf(ref)` / `TargetMachine.DataLayout()`），实现 `Close()`；二次 Close → `ErrClosed`。
4. `llvm/target` 的 `SetTo(*ir.Module)` 承接旧 `Module.SetTarget`：写 triple + data layout。
5. 旧 `IsLinux/IsDarwin/IsWindows` 不迁移（YAGNI）；旧 `dataLayout.GetSizeOfType` 实际返回 bit 却注释 byte，新 API 更名 `SizeOfTypeInBits` 修正。
6. `Arch` 表驱动初始化：NVPTX 无 AsmPrinter/AsmParser/Disassembler，XCore 无 AsmPrinter/AsmParser（binding 中确实未声明），表中以 nil 表示跳过。

---

## 核心定义（任务 1-6 共用，实现时逐字采用）

### error.go（新增 ErrParse）

在 `ErrVerify` 之后插入一行：

```go
	ErrVerify                      // IR 验证失败
	ErrParse                       // IR/bitcode 解析失败（ir.ParseIR/ParseBitcode）
	ErrUnsupported                 // 映射遇到不支持的类型
```

### internal/binding/Types.go（新增句柄）

在 `LLVMDiagnosticInfoRef` 之后加入：

```go
	LLVMMemoryBufferRef struct{ c C.LLVMMemoryBufferRef }
```

并在 `IsNil` 列表中加入：

```go
func (ref LLVMMemoryBufferRef) IsNil() bool { return ref.c == nil }
```

### internal/binding/Core.go（MemoryBuffer 族 + PrintModuleToFile）

在 `LLVMPrintModuleToString` 之后加入：

```go
// LLVMPrintModuleToFile Print a module to a file.
func LLVMPrintModuleToFile(m LLVMModuleRef, filename string) error {
	return string2CString(filename, func(filename *C.char) error {
		return llvmError2Error(func(errstr **C.char) C.LLVMBool {
			return C.LLVMPrintModuleToFile(m.c, filename, errstr)
		})
	})
}
```

在文件末尾加入：

```go
// LLVMCreateMemoryBufferWithContentsOfFile Read a file into a memory buffer.
func LLVMCreateMemoryBufferWithContentsOfFile(path string) (LLVMMemoryBufferRef, error) {
	var buf LLVMMemoryBufferRef
	err := string2CString(path, func(cpath *C.char) error {
		return llvmError2Error(func(errstr **C.char) C.LLVMBool {
			return C.LLVMCreateMemoryBufferWithContentsOfFile(cpath, &buf.c, errstr)
		})
	})
	if err != nil {
		return LLVMMemoryBufferRef{}, err
	}
	return buf, nil
}

// LLVMCreateMemoryBufferWithMemoryRangeCopy Create a memory buffer from a memory range, copying the data.
func LLVMCreateMemoryBufferWithMemoryRangeCopy(data []byte, name string) LLVMMemoryBufferRef {
	return string2CString(name, func(name *C.char) LLVMMemoryBufferRef {
		var ptr *C.char
		if len(data) > 0 {
			ptr = (*C.char)(unsafe.Pointer(&data[0]))
		}
		return LLVMMemoryBufferRef{c: C.LLVMCreateMemoryBufferWithMemoryRangeCopy(ptr, C.size_t(len(data)), name)}
	})
}

// LLVMGetBufferStart Get the start of the buffer.
func LLVMGetBufferStart(buf LLVMMemoryBufferRef) []byte {
	ptr := C.LLVMGetBufferStart(buf.c)
	if ptr == nil {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(ptr)), int(C.LLVMGetBufferSize(buf.c)))
}

// LLVMGetBufferSize Get the size of the buffer.
func LLVMGetBufferSize(buf LLVMMemoryBufferRef) uint64 {
	return uint64(C.LLVMGetBufferSize(buf.c))
}

// LLVMDisposeMemoryBuffer Dispose of a memory buffer.
func LLVMDisposeMemoryBuffer(buf LLVMMemoryBufferRef) {
	C.LLVMDisposeMemoryBuffer(buf.c)
}
```

### internal/binding/BitReader.go（新建）

```go
package binding

/*
#include "llvm-c/BitReader.h"
*/
import "C"

// LLVMParseBitcodeInContext Read bitcode from a memory buffer and convert it to an LLVM module.
// The memory buffer is not consumed; the caller keeps ownership.
func LLVMParseBitcodeInContext(c LLVMContextRef, mem LLVMMemoryBufferRef) (LLVMModuleRef, error) {
	var out LLVMModuleRef
	err := llvmError2Error(func(errstr **C.char) C.LLVMBool {
		return C.LLVMParseBitcodeInContext(c.c, mem.c, &out.c, errstr)
	})
	if err != nil {
		return LLVMModuleRef{}, err
	}
	return out, nil
}
```

### internal/binding/BitWriter.go（新建）

```go
package binding

/*
#include "llvm-c/BitWriter.h"
*/
import "C"
import "errors"

// LLVMWriteBitcodeToFile Write a module to the specified path.
func LLVMWriteBitcodeToFile(m LLVMModuleRef, path string) error {
	return string2CString(path, func(path *C.char) error {
		if C.LLVMWriteBitcodeToFile(m.c, path) != 0 {
			return errors.New("failed to write bitcode file")
		}
		return nil
	})
}

// LLVMWriteBitcodeToMemoryBuffer Write a module to a memory buffer owned by the caller.
func LLVMWriteBitcodeToMemoryBuffer(m LLVMModuleRef) LLVMMemoryBufferRef {
	return LLVMMemoryBufferRef{c: C.LLVMWriteBitcodeToMemoryBuffer(m.c)}
}
```

### internal/binding/IRReader.go（新建）

```go
package binding

/*
#include "llvm-c/IRReader.h"
*/
import "C"

// LLVMParseIRInContext Read LLVM IR from a memory buffer and convert it into an LLVM module.
// Wraps the non-deprecated LLVMParseIRInContext2; the memory buffer is not consumed.
func LLVMParseIRInContext(c LLVMContextRef, mem LLVMMemoryBufferRef) (LLVMModuleRef, error) {
	var out LLVMModuleRef
	err := llvmError2Error(func(errstr **C.char) C.LLVMBool {
		return C.LLVMParseIRInContext2(c.c, mem.c, &out.c, errstr)
	})
	if err != nil {
		return LLVMModuleRef{}, err
	}
	return out, nil
}
```

### internal/binding/TargetMachine.go（新增 EmitToMemoryBuffer）

在 `LLVMTargetMachineEmitToFile` 之后加入：

```go
// LLVMTargetMachineEmitToMemoryBuffer Emits an asm or object file for the given module to a memory buffer.
func LLVMTargetMachineEmitToMemoryBuffer(t LLVMTargetMachineRef, m LLVMModuleRef, codegen LLVMCodeGenFileType) (LLVMMemoryBufferRef, error) {
	var out LLVMMemoryBufferRef
	err := llvmError2Error(func(errstr **C.char) C.LLVMBool {
		return C.LLVMTargetMachineEmitToMemoryBuffer(t.c, m.c, C.LLVMCodeGenFileType(codegen), errstr, &out.c)
	})
	if err != nil {
		return LLVMMemoryBufferRef{}, err
	}
	return out, nil
}
```

### memorybuffer.go（新建，root）

```go
package llvm

import (
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
	return &MemoryBuffer{ref: ref}
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

// Close 释放缓冲；二次调用返回 ErrClosed
func (b *MemoryBuffer) Close() error {
	if b.closed {
		return &Error{Reason: ErrClosed, Op: "llvm.MemoryBuffer.Close", Msg: "memory buffer already closed"}
	}
	b.closed = true
	binding.LLVMDisposeMemoryBuffer(b.ref)
	return nil
}
```

### datalayout.go（新建，root）

```go
package llvm

import (
	"github.com/kkkunny/go-llvm/internal/binding"
)

// ByteOrder 目标字节序
type ByteOrder int32

const (
	LittleEndian ByteOrder = ByteOrder(binding.LLVMLittleEndian)
	BigEndian    ByteOrder = ByteOrder(binding.LLVMBigEndian)
)

// DataLayout 目标数据布局查询句柄；由 NewDataLayout 或 TargetMachine 创建，用毕 Close
type DataLayout struct {
	ref    binding.LLVMTargetDataRef
	closed bool
}

// NewDataLayout 由布局字符串创建数据布局（拥有句柄）
func NewDataLayout(layout string) *DataLayout {
	ref := binding.LLVMCreateTargetData(layout)
	if ref.IsNil() {
		errPanic(ErrInvalidArg, "llvm.NewDataLayout", "invalid data layout string %q", layout)
	}
	return &DataLayout{ref: ref}
}

// DataLayoutOf 由底层句柄构建数据布局（供 llvm/target 桥接使用，接管所有权）
func DataLayoutOf(ref binding.LLVMTargetDataRef) *DataLayout {
	return &DataLayout{ref: ref}
}

// check 前置校验：句柄非空且未释放
func (d *DataLayout) check(op string) {
	if d.ref.IsNil() {
		errPanic(ErrInvalidArg, op, "nil data layout")
	}
	if d.closed {
		errPanic(ErrUseAfterFree, op, "data layout is closed")
	}
}

// checkType 前置校验类型参数
func (d *DataLayout) checkType(op string, t AnyType) {
	d.check(op)
	if t == nil || t.IsNil() {
		errPanic(ErrInvalidArg, op, "nil type")
	}
}

// checkValue 前置校验值参数
func (d *DataLayout) checkValue(op string, v AnyValue) {
	d.check(op)
	if v == nil || v.IsNil() {
		errPanic(ErrInvalidArg, op, "nil value")
	}
}

// Close 释放句柄；二次调用返回 ErrClosed
func (d *DataLayout) Close() error {
	if d.closed {
		return &Error{Reason: ErrClosed, Op: "llvm.DataLayout.Close", Msg: "data layout already closed"}
	}
	d.closed = true
	binding.LLVMDisposeTargetData(d.ref)
	return nil
}

// String 布局字符串
func (d *DataLayout) String() string {
	d.check("llvm.DataLayout.String")
	return binding.LLVMCopyStringRepOfTargetData(d.ref)
}

// ByteOrder 字节序
func (d *DataLayout) ByteOrder() ByteOrder {
	d.check("llvm.DataLayout.ByteOrder")
	return ByteOrder(binding.LLVMByteOrder(d.ref))
}

// PointerSize 指针字节数
func (d *DataLayout) PointerSize() uint32 {
	d.check("llvm.DataLayout.PointerSize")
	return binding.LLVMPointerSize(d.ref)
}

// PointerSizeForAS 指定地址空间的指针字节数
func (d *DataLayout) PointerSizeForAS(addrspace uint32) uint32 {
	d.check("llvm.DataLayout.PointerSizeForAS")
	return binding.LLVMPointerSizeForAS(d.ref, addrspace)
}

// IntPtrType 指针等宽整数类型
func (d *DataLayout) IntPtrType(ctx *Context) IntType {
	d.check("llvm.DataLayout.IntPtrType")
	ctx.checkAlive("llvm.DataLayout.IntPtrType")
	return IntType{Type[IntT]{ref: binding.LLVMIntPtrTypeInContext(ctx.ref, d.ref), ctx: ctx}}
}

// IntPtrTypeForAS 指定地址空间指针等宽整数类型
func (d *DataLayout) IntPtrTypeForAS(ctx *Context, addrspace uint32) IntType {
	d.check("llvm.DataLayout.IntPtrTypeForAS")
	ctx.checkAlive("llvm.DataLayout.IntPtrTypeForAS")
	return IntType{Type[IntT]{ref: binding.LLVMIntPtrTypeForASInContext(ctx.ref, d.ref, addrspace), ctx: ctx}}
}

// SizeOfTypeInBits 类型位宽（bit）
func (d *DataLayout) SizeOfTypeInBits(t AnyType) uint64 {
	d.checkType("llvm.DataLayout.SizeOfTypeInBits", t)
	return binding.LLVMSizeOfTypeInBits(d.ref, t.Ref())
}

// StoreSizeOfType 实际存储大小（byte）
func (d *DataLayout) StoreSizeOfType(t AnyType) uint64 {
	d.checkType("llvm.DataLayout.StoreSizeOfType", t)
	return binding.LLVMStoreSizeOfType(d.ref, t.Ref())
}

// ABISizeOfType ABI 大小（byte）
func (d *DataLayout) ABISizeOfType(t AnyType) uint64 {
	d.checkType("llvm.DataLayout.ABISizeOfType", t)
	return binding.LLVMABISizeOfType(d.ref, t.Ref())
}

// ABIAlignOfType ABI 对齐（byte）
func (d *DataLayout) ABIAlignOfType(t AnyType) uint32 {
	d.checkType("llvm.DataLayout.ABIAlignOfType", t)
	return binding.LLVMABIAlignmentOfType(d.ref, t.Ref())
}

// CallFrameAlignOfType 调用栈帧对齐（byte）
func (d *DataLayout) CallFrameAlignOfType(t AnyType) uint32 {
	d.checkType("llvm.DataLayout.CallFrameAlignOfType", t)
	return binding.LLVMCallFrameAlignmentOfType(d.ref, t.Ref())
}

// PrefAlignOfType 编译器推荐对齐（byte）
func (d *DataLayout) PrefAlignOfType(t AnyType) uint32 {
	d.checkType("llvm.DataLayout.PrefAlignOfType", t)
	return binding.LLVMPreferredAlignmentOfType(d.ref, t.Ref())
}

// PrefAlignOfGlobal 全局变量推荐对齐（byte）
func (d *DataLayout) PrefAlignOfGlobal(g AnyValue) uint32 {
	d.checkValue("llvm.DataLayout.PrefAlignOfGlobal", g)
	return binding.LLVMPreferredAlignmentOfGlobal(d.ref, g.Ref())
}

// ElementAtOffset 包含指定字节偏移的结构体元素下标
func (d *DataLayout) ElementAtOffset(st AnyType, offset uint64) uint32 {
	d.checkType("llvm.DataLayout.ElementAtOffset", st)
	return binding.LLVMElementAtOffset(d.ref, st.Ref(), offset)
}

// OffsetOfElement 指定结构体元素的字节偏移
func (d *DataLayout) OffsetOfElement(st AnyType, i uint32) uint64 {
	d.checkType("llvm.DataLayout.OffsetOfElement", st)
	return binding.LLVMOffsetOfElement(d.ref, st.Ref(), i)
}
```

### ir/parse.go（新建）

```go
package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// ParseIR 从内存缓冲解析 LLVM IR 文本；缓冲不被消费，仍由调用方 Close。
// 解析失败返回 ErrParse 与完整诊断。
func ParseIR(ctx *llvm.Context, buf *llvm.MemoryBuffer) (*Module, error) {
	const op = "ir.ParseIR"
	if !ctx.Alive() {
		errPanic(llvm.ErrUseAfterFree, op, "context is closed")
	}
	if buf == nil {
		errPanic(llvm.ErrInvalidArg, op, "nil memory buffer")
	}
	if !buf.Alive() {
		errPanic(llvm.ErrUseAfterFree, op, "memory buffer is closed")
	}
	ref, err := binding.LLVMParseIRInContext(ctx.Ref(), buf.Ref())
	if err != nil {
		return nil, llvm.WrapError(llvm.ErrParse, op, err)
	}
	return newModule(ctx, ref), nil
}

// ParseBitcode 从内存缓冲解析 bitcode；缓冲不被消费，仍由调用方 Close。
// 解析失败返回 ErrParse 与完整诊断。
func ParseBitcode(ctx *llvm.Context, buf *llvm.MemoryBuffer) (*Module, error) {
	const op = "ir.ParseBitcode"
	if !ctx.Alive() {
		errPanic(llvm.ErrUseAfterFree, op, "context is closed")
	}
	if buf == nil {
		errPanic(llvm.ErrInvalidArg, op, "nil memory buffer")
	}
	if !buf.Alive() {
		errPanic(llvm.ErrUseAfterFree, op, "memory buffer is closed")
	}
	ref, err := binding.LLVMParseBitcodeInContext(ctx.Ref(), buf.Ref())
	if err != nil {
		return nil, llvm.WrapError(llvm.ErrParse, op, err)
	}
	return newModule(ctx, ref), nil
}
```

### ir/module.go（铸造助手 + 序列化 + DataLayout）

`NewModule`/`Clone` 改为共用铸造助手：

```go
// newModule 由底层句柄铸造模块并登记到 Context 生命周期
func newModule(ctx *llvm.Context, ref binding.LLVMModuleRef) *Module {
	m := &Module{ref: ref, ctx: ctx, life: llvm.NewLifetime()}
	m.unown = ctx.Own(m)
	return m
}

// NewModule 创建模块并登记到 Context 生命周期
func NewModule(ctx *llvm.Context, name string) *Module {
	if !ctx.Alive() {
		errPanic(llvm.ErrUseAfterFree, "ir.NewModule", "context is closed")
	}
	return newModule(ctx, binding.LLVMModuleCreateWithNameInContext(name, ctx.Ref()))
}
```

`Clone` 的构造部分改为：

```go
func (m *Module) Clone() *Module {
	return newModule(m.ctx, binding.LLVMCloneModule(m.ref))
}
```

新增方法：

```go
// DataLayout 模块数据布局（owned 副本，用毕 Close）
func (m *Module) DataLayout() *llvm.DataLayout {
	return llvm.NewDataLayout(binding.LLVMGetDataLayoutStr(m.ref))
}

// WriteToFile 将模块 IR 文本写入文件；失败返回 ErrIO
func (m *Module) WriteToFile(path string) error {
	if err := binding.LLVMPrintModuleToFile(m.ref, path); err != nil {
		return llvm.WrapError(llvm.ErrIO, "ir.Module.WriteToFile", err)
	}
	return nil
}

// WriteBitcode 将模块 bitcode 写入文件；失败返回 ErrIO
func (m *Module) WriteBitcode(path string) error {
	if err := binding.LLVMWriteBitcodeToFile(m.ref, path); err != nil {
		return llvm.WrapError(llvm.ErrIO, "ir.Module.WriteBitcode", err)
	}
	return nil
}

// Bitcode 将模块序列化为 bitcode 内存缓冲
func (m *Module) Bitcode() *llvm.MemoryBuffer {
	return llvm.MemoryBufferOf(binding.LLVMWriteBitcodeToMemoryBuffer(m.ref))
}
```

### target/target.go（新建）

```go
package target

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Arch 目标架构
type Arch uint8

const (
	AArch64 Arch = iota
	AMDGPU
	ARM
	AVR
	BPF
	Hexagon
	Lanai
	LoongArch
	Mips
	MSP430
	NVPTX
	PowerPC
	RISCV
	Sparc
	SystemZ
	VE
	WebAssembly
	X86
	XCore
)

// archInit 单架构初始化函数表；nil 表示该组件在 binding 中不可用
type archInit struct {
	info         func()
	target       func()
	mc           func()
	asmPrinter   func()
	asmParser    func()
	disassembler func()
}

var archInits = map[Arch]archInit{
	AArch64:     {binding.LLVMInitializeAArch64TargetInfo, binding.LLVMInitializeAArch64Target, binding.LLVMInitializeAArch64TargetMC, binding.LLVMInitializeAArch64AsmPrinter, binding.LLVMInitializeAArch64AsmParser, binding.LLVMInitializeAArch64Disassembler},
	AMDGPU:      {binding.LLVMInitializeAMDGPUTargetInfo, binding.LLVMInitializeAMDGPUTarget, binding.LLVMInitializeAMDGPUTargetMC, binding.LLVMInitializeAMDGPUAsmPrinter, binding.LLVMInitializeAMDGPUAsmParser, binding.LLVMInitializeAMDGPUDisassembler},
	ARM:         {binding.LLVMInitializeARMTargetInfo, binding.LLVMInitializeARMTarget, binding.LLVMInitializeARMTargetMC, binding.LLVMInitializeARMAsmPrinter, binding.LLVMInitializeARMAsmParser, binding.LLVMInitializeARMDisassembler},
	AVR:         {binding.LLVMInitializeAVRTargetInfo, binding.LLVMInitializeAVRTarget, binding.LLVMInitializeAVRTargetMC, binding.LLVMInitializeAVRAsmPrinter, binding.LLVMInitializeAVRAsmParser, binding.LLVMInitializeAVRDisassembler},
	BPF:         {binding.LLVMInitializeBPFTargetInfo, binding.LLVMInitializeBPFTarget, binding.LLVMInitializeBPFTargetMC, binding.LLVMInitializeBPFAsmPrinter, binding.LLVMInitializeBPFAsmParser, binding.LLVMInitializeBPFDisassembler},
	Hexagon:     {binding.LLVMInitializeHexagonTargetInfo, binding.LLVMInitializeHexagonTarget, binding.LLVMInitializeHexagonTargetMC, binding.LLVMInitializeHexagonAsmPrinter, binding.LLVMInitializeHexagonAsmParser, binding.LLVMInitializeHexagonDisassembler},
	Lanai:       {binding.LLVMInitializeLanaiTargetInfo, binding.LLVMInitializeLanaiTarget, binding.LLVMInitializeLanaiTargetMC, binding.LLVMInitializeLanaiAsmPrinter, binding.LLVMInitializeLanaiAsmParser, binding.LLVMInitializeLanaiDisassembler},
	LoongArch:   {binding.LLVMInitializeLoongArchTargetInfo, binding.LLVMInitializeLoongArchTarget, binding.LLVMInitializeLoongArchTargetMC, binding.LLVMInitializeLoongArchAsmPrinter, binding.LLVMInitializeLoongArchAsmParser, binding.LLVMInitializeLoongArchDisassembler},
	Mips:        {binding.LLVMInitializeMipsTargetInfo, binding.LLVMInitializeMipsTarget, binding.LLVMInitializeMipsTargetMC, binding.LLVMInitializeMipsAsmPrinter, binding.LLVMInitializeMipsAsmParser, binding.LLVMInitializeMipsDisassembler},
	MSP430:      {binding.LLVMInitializeMSP430TargetInfo, binding.LLVMInitializeMSP430Target, binding.LLVMInitializeMSP430TargetMC, binding.LLVMInitializeMSP430AsmPrinter, binding.LLVMInitializeMSP430AsmParser, binding.LLVMInitializeMSP430Disassembler},
	NVPTX:       {binding.LLVMInitializeNVPTXTargetInfo, binding.LLVMInitializeNVPTXTarget, binding.LLVMInitializeNVPTXTargetMC, nil, nil, nil},
	PowerPC:     {binding.LLVMInitializePowerPCTargetInfo, binding.LLVMInitializePowerPCTarget, binding.LLVMInitializePowerPCTargetMC, binding.LLVMInitializePowerPCAsmPrinter, binding.LLVMInitializePowerPCAsmParser, binding.LLVMInitializePowerPCDisassembler},
	RISCV:       {binding.LLVMInitializeRISCVTargetInfo, binding.LLVMInitializeRISCVTarget, binding.LLVMInitializeRISCVTargetMC, binding.LLVMInitializeRISCVAsmPrinter, binding.LLVMInitializeRISCVAsmParser, binding.LLVMInitializeRISCVDisassembler},
	Sparc:       {binding.LLVMInitializeSparcTargetInfo, binding.LLVMInitializeSparcTarget, binding.LLVMInitializeSparcTargetMC, binding.LLVMInitializeSparcAsmPrinter, binding.LLVMInitializeSparcAsmParser, binding.LLVMInitializeSparcDisassembler},
	SystemZ:     {binding.LLVMInitializeSystemZTargetInfo, binding.LLVMInitializeSystemZTarget, binding.LLVMInitializeSystemZTargetMC, binding.LLVMInitializeSystemZAsmPrinter, binding.LLVMInitializeSystemZAsmParser, binding.LLVMInitializeSystemZDisassembler},
	VE:          {binding.LLVMInitializeVETargetInfo, binding.LLVMInitializeVETarget, binding.LLVMInitializeVETargetMC, binding.LLVMInitializeVEAsmPrinter, binding.LLVMInitializeVEAsmParser, binding.LLVMInitializeVEDisassembler},
	WebAssembly: {binding.LLVMInitializeWebAssemblyTargetInfo, binding.LLVMInitializeWebAssemblyTarget, binding.LLVMInitializeWebAssemblyTargetMC, binding.LLVMInitializeWebAssemblyAsmPrinter, binding.LLVMInitializeWebAssemblyAsmParser, binding.LLVMInitializeWebAssemblyDisassembler},
	X86:         {binding.LLVMInitializeX86TargetInfo, binding.LLVMInitializeX86Target, binding.LLVMInitializeX86TargetMC, binding.LLVMInitializeX86AsmPrinter, binding.LLVMInitializeX86AsmParser, binding.LLVMInitializeX86Disassembler},
	XCore:       {binding.LLVMInitializeXCoreTargetInfo, binding.LLVMInitializeXCoreTarget, binding.LLVMInitializeXCoreTargetMC, nil, nil, binding.LLVMInitializeXCoreDisassembler},
}

// InitAll 初始化全部目标（infos/targets/MCs/asm printers/asm parsers/disassemblers）
func InitAll() {
	binding.LLVMInitializeAllTargetInfos()
	binding.LLVMInitializeAllTargets()
	binding.LLVMInitializeAllTargetMCs()
	binding.LLVMInitializeAllAsmPrinters()
	binding.LLVMInitializeAllAsmParsers()
	binding.LLVMInitializeAllDisassemblers()
}

// InitNative 初始化宿主目标；失败返回 ErrCodeGen
func InitNative() error {
	for _, f := range []struct {
		op string
		fn func() error
	}{
		{"target.InitNative", binding.LLVMInitializeNativeTarget},
		{"target.InitNativeAsmPrinter", binding.LLVMInitializeNativeAsmPrinter},
		{"target.InitNativeAsmParser", binding.LLVMInitializeNativeAsmParser},
		{"target.InitNativeDisassembler", binding.LLVMInitializeNativeDisassembler},
	} {
		if err := f.fn(); err != nil {
			return llvm.WrapError(llvm.ErrCodeGen, f.op, err)
		}
	}
	return nil
}

// Init 初始化指定架构的全部可用组件；未知架构 panic ErrInvalidArg
func Init(arch Arch) {
	init, ok := archInits[arch]
	if !ok {
		llvm.Panicf(llvm.ErrInvalidArg, "target.Init", "unknown arch %d", arch)
	}
	init.info()
	init.target()
	init.mc()
	if init.asmPrinter != nil {
		init.asmPrinter()
	}
	if init.asmParser != nil {
		init.asmParser()
	}
	if init.disassembler != nil {
		init.disassembler()
	}
}

// Target 目标描述（初始化后全局唯一，无需释放）
type Target struct{ ref binding.LLVMTargetRef }

// Ref 返回底层句柄（供 llvm/target 内部桥接使用）
func (t Target) Ref() binding.LLVMTargetRef { return t.ref }

// FromName 按目标名查找（如 "x86-64"）
func FromName(name string) (Target, bool) {
	ref := binding.LLVMGetTargetFromName(name)
	if ref.IsNil() {
		return Target{}, false
	}
	return Target{ref: ref}, true
}

// FromTriple 按三元组查找；未知三元组返回 ErrNotFound
func FromTriple(triple string) (Target, error) {
	ref, err := binding.LLVMGetTargetFromTriple(triple)
	if err != nil {
		return Target{}, llvm.WrapError(llvm.ErrNotFound, "target.FromTriple", err)
	}
	return Target{ref: ref}, nil
}

// NativeTarget 宿主目标
func NativeTarget() (Target, error) { return FromTriple(DefaultTriple()) }

// Name 目标名
func (t Target) Name() string { return binding.LLVMGetTargetName(t.ref) }

// Description 目标描述
func (t Target) Description() string { return binding.LLVMGetTargetDescription(t.ref) }

// HasJIT 是否支持 JIT
func (t Target) HasJIT() bool { return binding.LLVMTargetHasJIT(t.ref) }

// HasTargetMachine 是否支持目标机器
func (t Target) HasTargetMachine() bool { return binding.LLVMTargetHasTargetMachine(t.ref) }

// HasAsmBackend 是否有汇编后端
func (t Target) HasAsmBackend() bool { return binding.LLVMTargetHasAsmBackend(t.ref) }

// DefaultTriple 宿主三元组
func DefaultTriple() string { return binding.LLVMGetDefaultTargetTriple() }

// NormalizeTriple 规范化三元组
func NormalizeTriple(triple string) string { return binding.LLVMNormalizeTargetTriple(triple) }

// HostCPUName 宿主 CPU 名
func HostCPUName() string { return binding.LLVMGetHostCPUName() }

// HostCPUFeatures 宿主 CPU 特性串
func HostCPUFeatures() string { return binding.LLVMGetHostCPUFeatures() }
```

### target/machine.go（新建）

```go
package target

import (
	"errors"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/ir"
)

// OptLevel 代码生成优化级别
type OptLevel int32

const (
	OptNone       OptLevel = OptLevel(binding.LLVMCodeGenLevelNone)
	OptLess       OptLevel = OptLevel(binding.LLVMCodeGenLevelLess)
	OptDefault    OptLevel = OptLevel(binding.LLVMCodeGenLevelDefault)
	OptAggressive OptLevel = OptLevel(binding.LLVMCodeGenLevelAggressive)
)

// RelocMode 重定位模式
type RelocMode int32

const (
	RelocDefault      RelocMode = RelocMode(binding.LLVMRelocDefault)
	RelocStatic       RelocMode = RelocMode(binding.LLVMRelocStatic)
	RelocPIC          RelocMode = RelocMode(binding.LLVMRelocPIC)
	RelocDynamicNoPic RelocMode = RelocMode(binding.LLVMRelocDynamicNoPic)
	RelocROPI         RelocMode = RelocMode(binding.LLVMRelocROPI)
	RelocRWPI         RelocMode = RelocMode(binding.LLVMRelocRWPI)
	RelocROPI_RWPI    RelocMode = RelocMode(binding.LLVMRelocROPI_RWPI)
)

// CodeModel 代码模型
type CodeModel int32

const (
	CodeModelDefault    CodeModel = CodeModel(binding.LLVMCodeModelDefault)
	CodeModelJITDefault CodeModel = CodeModel(binding.LLVMCodeModelJITDefault)
	CodeModelTiny       CodeModel = CodeModel(binding.LLVMCodeModelTiny)
	CodeModelSmall      CodeModel = CodeModel(binding.LLVMCodeModelSmall)
	CodeModelKernel     CodeModel = CodeModel(binding.LLVMCodeModelKernel)
	CodeModelMedium     CodeModel = CodeModel(binding.LLVMCodeModelMedium)
	CodeModelLarge      CodeModel = CodeModel(binding.LLVMCodeModelLarge)
)

// FileType 产出文件类型
type FileType int32

const (
	AsmFile    FileType = FileType(binding.LLVMAssemblyFile)
	ObjectFile FileType = FileType(binding.LLVMObjectFile)
)

// TargetMachine 目标机器；独立所有权根（不经 Context.Own），用毕 Close
type TargetMachine struct {
	ref    binding.LLVMTargetMachineRef
	closed bool
}

// NewTargetMachine 创建目标机器；失败返回 ErrCodeGen
func NewTargetMachine(t Target, triple, cpu, features string, opt OptLevel, reloc RelocMode, cm CodeModel) (*TargetMachine, error) {
	ref := binding.LLVMCreateTargetMachine(t.ref, triple, cpu, features, binding.LLVMCodeGenOptLevel(opt), binding.LLVMRelocMode(reloc), binding.LLVMCodeModel(cm))
	if ref.IsNil() {
		return nil, llvm.WrapError(llvm.ErrCodeGen, "target.NewTargetMachine", errors.New("LLVMCreateTargetMachine returned null"))
	}
	return &TargetMachine{ref: ref}, nil
}

// check 前置校验：未释放
func (m *TargetMachine) check(op string) {
	if m.closed {
		llvm.Panicf(llvm.ErrUseAfterFree, op, "target machine is closed")
	}
}

// checkModule 前置校验模块可用
func checkModule(op string, mod *ir.Module) {
	if mod == nil || mod.Ref().IsNil() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil module")
	}
	if !mod.Lifetime().Alive() {
		llvm.Panicf(llvm.ErrUseAfterFree, op, "module is closed")
	}
}

// Close 释放目标机器；二次调用返回 ErrClosed
func (m *TargetMachine) Close() error {
	if m.closed {
		return &llvm.Error{Reason: llvm.ErrClosed, Op: "target.TargetMachine.Close", Msg: "target machine already closed"}
	}
	m.closed = true
	binding.LLVMDisposeTargetMachine(m.ref)
	return nil
}

// Triple 目标三元组
func (m *TargetMachine) Triple() string {
	m.check("target.TargetMachine.Triple")
	return binding.LLVMGetTargetMachineTriple(m.ref)
}

// CPU 目标 CPU
func (m *TargetMachine) CPU() string {
	m.check("target.TargetMachine.CPU")
	return binding.LLVMGetTargetMachineCPU(m.ref)
}

// Features 目标特性串
func (m *TargetMachine) Features() string {
	m.check("target.TargetMachine.Features")
	return binding.LLVMGetTargetMachineFeatureString(m.ref)
}

// Target 目标机器对应的目标
func (m *TargetMachine) Target() Target {
	m.check("target.TargetMachine.Target")
	return Target{ref: binding.LLVMGetTargetMachineTarget(m.ref)}
}

// DataLayout 目标数据布局（owned，用毕 Close）
func (m *TargetMachine) DataLayout() *llvm.DataLayout {
	m.check("target.TargetMachine.DataLayout")
	return llvm.DataLayoutOf(binding.LLVMCreateTargetDataLayout(m.ref))
}

// SetAsmVerbosity 设置汇编输出详细程度
func (m *TargetMachine) SetAsmVerbosity(verbose bool) {
	m.check("target.TargetMachine.SetAsmVerbosity")
	binding.LLVMSetTargetMachineAsmVerbosity(m.ref, verbose)
}

// SetTo 把目标三元组与数据布局写入模块（旧 Module.SetTarget 的替代）
func (m *TargetMachine) SetTo(mod *ir.Module) {
	const op = "target.TargetMachine.SetTo"
	m.check(op)
	checkModule(op, mod)
	mod.SetTargetTriple(m.Triple())
	dl := m.DataLayout()
	defer dl.Close()
	mod.SetDataLayout(dl.String())
}

// EmitToFile 将模块编译为汇编/目标文件；失败返回 ErrCodeGen
func (m *TargetMachine) EmitToFile(mod *ir.Module, path string, ft FileType) error {
	const op = "target.TargetMachine.EmitToFile"
	m.check(op)
	checkModule(op, mod)
	if err := binding.LLVMTargetMachineEmitToFile(m.ref, mod.Ref(), path, binding.LLVMCodeGenFileType(ft)); err != nil {
		return llvm.WrapError(llvm.ErrCodeGen, op, err)
	}
	return nil
}

// Emit 将模块编译为汇编/目标代码内存缓冲
func (m *TargetMachine) Emit(mod *ir.Module, ft FileType) (*llvm.MemoryBuffer, error) {
	const op = "target.TargetMachine.Emit"
	m.check(op)
	checkModule(op, mod)
	buf, err := binding.LLVMTargetMachineEmitToMemoryBuffer(m.ref, mod.Ref(), binding.LLVMCodeGenFileType(ft))
	if err != nil {
		return nil, llvm.WrapError(llvm.ErrCodeGen, op, err)
	}
	return llvm.MemoryBufferOf(buf), nil
}
```

---

## 任务分解

### 任务 1：binding 补齐（内存缓冲/解析/bitcode/emit）

**文件：**
- 修改：`internal/binding/Types.go`、`internal/binding/Core.go`、`internal/binding/TargetMachine.go`
- 创建：`internal/binding/BitReader.go`、`internal/binding/BitWriter.go`、`internal/binding/IRReader.go`

binding 层无测试（1:1 cgo，由上层任务测试覆盖）。

- [ ] **步骤 1：Types.go 增加 LLVMMemoryBufferRef**——按「核心定义」逐字添加类型与 `IsNil`
- [ ] **步骤 2：Core.go 增加 MemoryBuffer 族与 LLVMPrintModuleToFile**——按「核心定义」逐字添加
- [ ] **步骤 3：创建 BitReader.go / BitWriter.go / IRReader.go**——按「核心定义」逐字添加
- [ ] **步骤 4：TargetMachine.go 增加 LLVMTargetMachineEmitToMemoryBuffer**——按「核心定义」逐字添加
- [ ] **步骤 5：编译确认**

运行：`go build ./... && go vet ./...`
预期：无输出（成功）

- [ ] **步骤 6：Commit**

```bash
git add internal/binding/
git commit -m "feat: binding补齐内存缓冲/解析/bitcode/emit绑定"
```

### 任务 2：root MemoryBuffer

**文件：**
- 创建：`memorybuffer.go`、`memorybuffer_test.go`

- [ ] **步骤 1：编写失败测试**

```go
package llvm

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMemoryBufferRoundTrip(t *testing.T) {
	buf := NewMemoryBuffer([]byte("hello"), "test")
	defer buf.Close()

	if !buf.Alive() {
		t.Fatal("new buffer should be alive")
	}
	if got := string(buf.Bytes()); got != "hello" {
		t.Fatalf("Bytes() = %q", got)
	}
	if buf.Len() != 5 {
		t.Fatalf("Len() = %d", buf.Len())
	}
	if err := buf.Close(); err == nil || err.(*Error).Reason != ErrClosed {
		t.Fatalf("double close should return ErrClosed, got %v", err)
	}
}

func TestReadFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "in.txt")
	if err := os.WriteFile(path, []byte("abc"), 0o644); err != nil {
		t.Fatal(err)
	}
	buf, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	defer buf.Close()
	if got := string(buf.Bytes()); got != "abc" {
		t.Fatalf("Bytes() = %q", got)
	}

	if _, err := ReadFile(filepath.Join(t.TempDir(), "missing")); err == nil || err.(*Error).Reason != ErrIO {
		t.Fatalf("missing file should return ErrIO, got %v", err)
	}
}

func TestMemoryBufferAfterClose(t *testing.T) {
	buf := NewMemoryBuffer([]byte("x"), "test")
	_ = buf.Close()
	if buf.Alive() {
		t.Fatal("closed buffer should not be alive")
	}
	if err := Catch(func() { buf.Bytes() }); err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("Bytes after close should panic ErrUseAfterFree, got %v", err)
	}
}
```

- [ ] **步骤 2：运行测试验证失败**

运行：`go test . -run 'TestMemoryBuffer|TestReadFile' -v`
预期：FAIL，`undefined: NewMemoryBuffer` / `undefined: ReadFile`

- [ ] **步骤 3：实现 memorybuffer.go**——按「核心定义」逐字实现

- [ ] **步骤 4：运行测试验证通过**

运行：`go test . -run 'TestMemoryBuffer|TestReadFile' -v && go vet ./...`
预期：PASS

- [ ] **步骤 5：Commit**

```bash
git add memorybuffer.go memorybuffer_test.go
git commit -m "feat: MemoryBuffer所有权封装"
```

### 任务 3：root DataLayout

**文件：**
- 创建：`datalayout.go`、`datalayout_test.go`

- [ ] **步骤 1：编写失败测试**

```go
package llvm

import (
	"strings"
	"testing"
)

func TestDataLayoutQueries(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()
	dl := NewDataLayout("e-p:64:64-i64:64-n8:16:32:64")
	defer dl.Close()

	if got := dl.ByteOrder(); got != LittleEndian {
		t.Fatalf("ByteOrder() = %v", got)
	}
	if got := dl.PointerSize(); got != 8 {
		t.Fatalf("PointerSize() = %d", got)
	}
	if got := dl.PointerSizeForAS(0); got != 8 {
		t.Fatalf("PointerSizeForAS(0) = %d", got)
	}

	i64 := ctx.Int(64)
	if got := dl.SizeOfTypeInBits(i64); got != 64 {
		t.Fatalf("SizeOfTypeInBits(i64) = %d", got)
	}
	if got := dl.StoreSizeOfType(i64); got != 8 {
		t.Fatalf("StoreSizeOfType(i64) = %d", got)
	}
	if got := dl.ABISizeOfType(i64); got != 8 {
		t.Fatalf("ABISizeOfType(i64) = %d", got)
	}
	if got := dl.ABIAlignOfType(i64); got != 8 {
		t.Fatalf("ABIAlignOfType(i64) = %d", got)
	}
	if got := dl.PrefAlignOfType(i64); got != 8 {
		t.Fatalf("PrefAlignOfType(i64) = %d", got)
	}
	if got := dl.CallFrameAlignOfType(i64); got == 0 {
		t.Fatal("CallFrameAlignOfType(i64) should be non-zero")
	}
	if !dl.IntPtrType(ctx).Equal(ctx.Int(64)) {
		t.Fatalf("IntPtrType() = %s", dl.IntPtrType(ctx))
	}
	if !dl.IntPtrTypeForAS(ctx, 0).Equal(ctx.Int(64)) {
		t.Fatalf("IntPtrTypeForAS() = %s", dl.IntPtrTypeForAS(ctx, 0))
	}
	if got := dl.String(); !strings.Contains(got, "p:64:64") || !strings.Contains(got, "i64:64") {
		t.Fatalf("String() = %q", got)
	}
}

func TestDataLayoutStructQueries(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()
	dl := NewDataLayout("e-p:64:64-i64:64-n8:16:32:64")
	defer dl.Close()

	st := ctx.Struct([]AnyType{ctx.Int(8), ctx.Int(64)}, false)
	if got := dl.OffsetOfElement(st, 1); got != 8 {
		t.Fatalf("OffsetOfElement(st, 1) = %d", got)
	}
	if got := dl.ElementAtOffset(st, 8); got != 1 {
		t.Fatalf("ElementAtOffset(st, 8) = %d", got)
	}
	if got := dl.SizeOfTypeInBits(st); got != 128 {
		t.Fatalf("SizeOfTypeInBits(st) = %d", got)
	}
	if got := dl.ABISizeOfType(st); got != 16 {
		t.Fatalf("ABISizeOfType(st) = %d", got)
	}
}

func TestDataLayoutClose(t *testing.T) {
	dl := NewDataLayout("e-p:64:64")
	_ = dl.Close()
	if err := dl.Close(); err == nil || err.(*Error).Reason != ErrClosed {
		t.Fatalf("double close should return ErrClosed, got %v", err)
	}
	if err := Catch(func() { dl.PointerSize() }); err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("use after close should panic ErrUseAfterFree, got %v", err)
	}
}
```

- [ ] **步骤 2：运行测试验证失败**

运行：`go test . -run TestDataLayout -v`
预期：FAIL，`undefined: NewDataLayout`

- [ ] **步骤 3：实现 datalayout.go**——按「核心定义」逐字实现

- [ ] **步骤 4：运行测试验证通过**

运行：`go test . -run TestDataLayout -v && go vet ./...`
预期：PASS

- [ ] **步骤 5：Commit**

```bash
git add datalayout.go datalayout_test.go
git commit -m "feat: DataLayout全套查询与ByteOrder"
```

### 任务 4：ir 解析与序列化

**文件：**
- 修改：`error.go`、`ir/module.go`
- 创建：`ir/parse.go`、`ir/parse_test.go`

- [ ] **步骤 1：编写失败测试**

```go
package ir

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

// stripModuleID 去掉 `; ModuleID = ...` 行：解析往返时该行会变为缓冲名或被 bitcode 丢弃
func stripModuleID(s string) string {
	var b strings.Builder
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(line, "; ModuleID =") {
			continue
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}

func buildRetModule(t *testing.T, name string) (*llvm.Context, *Module) {
	t.Helper()
	ctx := llvm.NewContext()
	m := NewModule(ctx, name)
	i32 := ctx.Int(32)
	fn := m.NewFunction("main", ctx.Fn(i32, nil, false))
	b := NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(ctx.ConstInt(i32, 0, false))
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	if err := m.Verify(); err != nil {
		t.Fatal(err)
	}
	return ctx, m
}

func TestParseIRRoundTrip(t *testing.T) {
	ctx, m := buildRetModule(t, "roundtrip")
	defer ctx.Close()
	defer m.Close()

	text := m.String()
	buf := llvm.NewMemoryBuffer([]byte(text), "roundtrip.ll")
	defer buf.Close()
	parsed, err := ParseIR(ctx, buf)
	if err != nil {
		t.Fatal(err)
	}
	defer parsed.Close()

	// 解析后 ModuleID 变为缓冲名，其余 IR 必须逐字一致
	if stripModuleID(parsed.String()) != stripModuleID(text) {
		t.Fatalf("round trip mismatch:\n%s\n---\n%s", parsed.String(), text)
	}
	if !buf.Alive() {
		t.Fatal("ParseIR must not consume the memory buffer")
	}
}

func TestParseIRInvalid(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	buf := llvm.NewMemoryBuffer([]byte("this is not IR"), "bad.ll")
	defer buf.Close()

	_, err := ParseIR(ctx, buf)
	if err == nil || err.(*llvm.Error).Reason != llvm.ErrParse {
		t.Fatalf("invalid IR should return ErrParse, got %v", err)
	}
}

func TestBitcodeRoundTrip(t *testing.T) {
	ctx, m := buildRetModule(t, "bitcode")
	defer ctx.Close()
	defer m.Close()

	path := filepath.Join(t.TempDir(), "m.bc")
	if err := m.WriteBitcode(path); err != nil {
		t.Fatal(err)
	}
	buf, err := llvm.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	defer buf.Close()
	parsed, err := ParseBitcode(ctx, buf)
	if err != nil {
		t.Fatal(err)
	}
	defer parsed.Close()
	if stripModuleID(parsed.String()) != stripModuleID(m.String()) {
		t.Fatalf("bitcode round trip mismatch:\n%s\n---\n%s", parsed.String(), m.String())
	}

	bc := m.Bitcode()
	defer bc.Close()
	parsed2, err := ParseBitcode(ctx, bc)
	if err != nil {
		t.Fatal(err)
	}
	defer parsed2.Close()
	if stripModuleID(parsed2.String()) != stripModuleID(m.String()) {
		t.Fatalf("memory bitcode round trip mismatch")
	}

	bad := llvm.NewMemoryBuffer([]byte("junk"), "bad.bc")
	defer bad.Close()
	if _, err := ParseBitcode(ctx, bad); err == nil || err.(*llvm.Error).Reason != llvm.ErrParse {
		t.Fatalf("invalid bitcode should return ErrParse, got %v", err)
	}
}

func TestModuleWriteToFile(t *testing.T) {
	ctx, m := buildRetModule(t, "writefile")
	defer ctx.Close()
	defer m.Close()

	path := filepath.Join(t.TempDir(), "m.ll")
	if err := m.WriteToFile(path); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "define i32 @main()") {
		t.Fatalf("written IR:\n%s", data)
	}

	if err := m.WriteToFile(filepath.Join(t.TempDir(), "no/such/dir/m.ll")); err == nil || err.(*llvm.Error).Reason != llvm.ErrIO {
		t.Fatalf("bad path should return ErrIO, got %v", err)
	}
}

func TestModuleDataLayout(t *testing.T) {
	ctx, m := buildRetModule(t, "dl")
	defer ctx.Close()
	defer m.Close()

	dl := m.DataLayout()
	defer dl.Close()
	if dl.PointerSize() == 0 {
		t.Fatal("module data layout pointer size should be non-zero")
	}

	g := m.NewGlobal("g", ctx.Int(32))
	g.SetAlign(4)
	if got := dl.PrefAlignOfGlobal(g); got != 4 {
		t.Fatalf("PrefAlignOfGlobal(g) = %d", got)
	}
}
```

- [ ] **步骤 2：运行测试验证失败**

运行：`go test ./ir -run 'TestParseIR|TestBitcode|TestModuleWriteToFile|TestModuleDataLayout' -v`
预期：FAIL，`undefined: ParseIR` / `undefined: ParseBitcode`

- [ ] **步骤 3：error.go 增加 ErrParse**——按「核心定义」插入
- [ ] **步骤 4：ir/module.go 重构 newModule 并新增方法**——按「核心定义」逐字实现（`NewModule`/`Clone` 改为委托 `newModule`，新增 `DataLayout`/`WriteToFile`/`WriteBitcode`/`Bitcode`）
- [ ] **步骤 5：创建 ir/parse.go**——按「核心定义」逐字实现

- [ ] **步骤 6：运行测试验证通过**

运行：`go test ./ir -v && go vet ./...`
预期：全部 PASS（含既有测试）

- [ ] **步骤 7：Commit**

```bash
git add error.go ir/module.go ir/parse.go ir/parse_test.go
git commit -m "feat: IR文本与bitcode解析/序列化"
```

### 任务 5：llvm/target 目标初始化与查询

**文件：**
- 创建：`target/target.go`、`target/target_test.go`

- [ ] **步骤 1：编写失败测试**

```go
package target

import (
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestInitAndHostQueries(t *testing.T) {
	InitNative()

	native, err := NativeTarget()
	if err != nil {
		t.Fatal(err)
	}
	if !native.HasTargetMachine() {
		t.Fatalf("native target %q should have a target machine", native.Name())
	}
	if !native.HasAsmBackend() {
		t.Fatalf("native target %q should have an asm backend", native.Name())
	}
	if native.Name() == "" || native.Description() == "" {
		t.Fatal("native target name/description should be non-empty")
	}
	if triple := DefaultTriple(); triple == "" || NormalizeTriple(triple) == "" {
		t.Fatal("host triple should be non-empty")
	}
	if HostCPUName() == "" {
		t.Fatal("host cpu name should be non-empty")
	}
}

func TestInitAllAndLookup(t *testing.T) {
	InitAll()
	Init(X86)

	if got, ok := FromName("x86-64"); !ok || got.Name() != "x86-64" {
		t.Fatalf("FromName(x86-64) = %q, %v", got.Name(), ok)
	}
	if _, ok := FromName("definitely-not-a-target"); ok {
		t.Fatal("unknown target name should not resolve")
	}
	if _, err := FromTriple("x86_64-unknown-linux-gnu"); err != nil {
		t.Fatalf("FromTriple(x86_64-unknown-linux-gnu) = %v", err)
	}
	if _, err := FromTriple("not-a-triple"); err == nil || err.(*llvm.Error).Reason != llvm.ErrNotFound {
		t.Fatalf("invalid triple should return ErrNotFound, got %v", err)
	}
}

func TestInitUnknownArch(t *testing.T) {
	if err := llvm.Catch(func() { Init(Arch(255)) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("unknown arch should panic ErrInvalidArg, got %v", err)
	}
}
```

- [ ] **步骤 2：运行测试验证失败**

运行：`go test ./target -v`
预期：FAIL，`undefined: InitNative`（包不存在）

- [ ] **步骤 3：实现 target/target.go**——按「核心定义」逐字实现

- [ ] **步骤 4：运行测试验证通过**

运行：`go test ./target -v && go vet ./...`
预期：PASS

- [ ] **步骤 5：Commit**

```bash
git add target/target.go target/target_test.go
git commit -m "feat: 目标注册/初始化与宿主查询"
```

### 任务 6：TargetMachine 与代码生成

**文件：**
- 创建：`target/machine.go`、`target/machine_test.go`

- [ ] **步骤 1：编写失败测试**

```go
package target

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
)

func newMachineModule(t *testing.T) (*llvm.Context, *ir.Module, *TargetMachine) {
	t.Helper()
	InitNative()
	native, err := NativeTarget()
	if err != nil {
		t.Fatal(err)
	}
	tm, err := NewTargetMachine(native, DefaultTriple(), HostCPUName(), HostCPUFeatures(), OptDefault, RelocPIC, CodeModelDefault)
	if err != nil {
		t.Fatal(err)
	}

	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "codegen")
	i32 := ctx.Int(32)
	fn := m.NewFunction("main", ctx.Fn(i32, nil, false))
	b := ir.NewBuilder(ctx)
	b.MoveToEnd(fn.NewBlock("entry"))
	b.Ret(ctx.ConstInt(i32, 0, false))
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	tm.SetTo(m)
	if err := m.Verify(); err != nil {
		t.Fatal(err)
	}
	return ctx, m, tm
}

func TestTargetMachineSetTo(t *testing.T) {
	ctx, m, tm := newMachineModule(t)
	defer ctx.Close()
	defer m.Close()
	defer tm.Close()

	if tm.Triple() != m.TargetTriple() {
		t.Fatalf("SetTo triple: module %q, machine %q", m.TargetTriple(), tm.Triple())
	}
	if tm.Target().Name() == "" {
		t.Fatal("Target() should be non-nil")
	}
	mdl := m.DataLayout()
	defer mdl.Close()
	tmdl := tm.DataLayout()
	defer tmdl.Close()
	if mdl.String() != tmdl.String() {
		t.Fatalf("SetTo data layout: module %q, machine %q", mdl.String(), tmdl.String())
	}
}

func TestTargetMachineEmit(t *testing.T) {
	ctx, m, tm := newMachineModule(t)
	defer ctx.Close()
	defer m.Close()
	defer tm.Close()

	asm, err := tm.Emit(m, AsmFile)
	if err != nil {
		t.Fatal(err)
	}
	defer asm.Close()
	if !strings.Contains(string(asm.Bytes()), "main") {
		t.Fatalf("asm output missing main:\n%s", asm.Bytes())
	}

	obj, err := tm.Emit(m, ObjectFile)
	if err != nil {
		t.Fatal(err)
	}
	defer obj.Close()
	if obj.Len() == 0 {
		t.Fatal("object output should be non-empty")
	}

	path := filepath.Join(t.TempDir(), "m.s")
	if err := tm.EmitToFile(m, path, AsmFile); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "main") {
		t.Fatalf("asm file missing main:\n%s", data)
	}
}

func TestTargetMachineChecks(t *testing.T) {
	ctx, m, tm := newMachineModule(t)
	defer ctx.Close()
	defer m.Close()

	closed := ir.NewModule(ctx, "closed")
	if err := closed.Close(); err != nil {
		t.Fatal(err)
	}
	if err := llvm.Catch(func() { tm.EmitToFile(closed, "x.s", AsmFile) }); err == nil || err.Reason != llvm.ErrUseAfterFree {
		t.Fatalf("emit closed module should panic ErrUseAfterFree, got %v", err)
	}

	if err := tm.Close(); err != nil {
		t.Fatal(err)
	}
	if err := tm.Close(); err == nil || err.(*llvm.Error).Reason != llvm.ErrClosed {
		t.Fatalf("double close should return ErrClosed, got %v", err)
	}
	if err := llvm.Catch(func() { tm.Triple() }); err == nil || err.Reason != llvm.ErrUseAfterFree {
		t.Fatalf("use after close should panic ErrUseAfterFree, got %v", err)
	}
}
```

- [ ] **步骤 2：运行测试验证失败**

运行：`go test ./target -run TestTargetMachine -v`
预期：FAIL，`undefined: NewTargetMachine`

- [ ] **步骤 3：实现 target/machine.go**——按「核心定义」逐字实现

- [ ] **步骤 4：运行测试验证通过**

运行：`go test ./target -v && go vet ./...`
预期：PASS

- [ ] **步骤 5：Commit**

```bash
git add target/machine.go target/machine_test.go
git commit -m "feat: TargetMachine与OBJ/ASM代码生成"
```

### 任务 7：文档同步与全量验收

**文件：**
- 修改：`README.md`、`AGENTS.md`

- [ ] **步骤 1：README 更新**——包表去掉 `llvm/target` 的 `(P1)` 标记（改为 `Target machines and code generation`），并在 Usage 末尾追加代码生成示例：

````markdown
### Code generation

```go
target.InitNative()
native, _ := target.NativeTarget()
tm, _ := target.NewTargetMachine(native, target.DefaultTriple(), target.HostCPUName(), target.HostCPUFeatures(),
	target.OptDefault, target.RelocPIC, target.CodeModelDefault)
defer tm.Close()

tm.SetTo(module)                                  // 写入 triple + data layout
_ = tm.EmitToFile(module, "main.o", target.ObjectFile)
asm, _ := tm.Emit(module, target.AsmFile)          // 或产出到内存缓冲
defer asm.Close()
```
````

- [ ] **步骤 2：AGENTS.md 更新**——布局表 `llvm/target (P1)` 改为 `llvm/target`；错误约定句补 `ErrParse`：

```markdown
- **Errors**: recoverable runtime failures return `error`; programmer errors `panic(*llvm.Error)` with `Reason`/`Op`/`Msg`, recoverable via `llvm.Catch`. Data-driven unsupported cases (`TypeOf[string]`, variadic Go funcs) return `ErrUnsupported`. `internal/binding` returns plain `error` (it must not import the root package); sub-packages wrap it with `llvm.WrapError(reason, op, err)` and use `ErrCodeGen`/`ErrJIT`/`ErrIO`/`ErrParse` for target/JIT/IO/parse failures.
```

- [ ] **步骤 3：全量验收**

运行：`go build ./... && go vet ./... && go test ./...`
预期：全绿（`ok github.com/kkkunny/go-llvm`、`ok .../ir`、`ok .../target`）

- [ ] **步骤 4：Commit**

```bash
git add README.md AGENTS.md
git commit -m "docs: 同步P1-1/P1-2包职责与错误约定"
```

---

## 迁移映射（旧 → 新）

| 旧 API（重构前 root 包） | 去向 |
|---|---|
| `Arch` / `InitializeTarget(arch)` / `InitializeAllTargets` 等 | `target.Arch` / `target.Init` / `target.InitAll` / `target.InitNative` |
| `NativeTarget` / `NewTargetFromTriple` | `target.NativeTarget` / `target.FromTriple`（新增 `FromName`） |
| `targetInfo.Name/Description/HasJIT/HasTargetMachine/HasAsmBackend` | `target.Target` 同名方法 |
| `machine` / `newMachine` | `target.TargetMachine` / `target.NewTargetMachine`（枚举改 `OptLevel`/`RelocMode`/`CodeModel`） |
| `machine.Free` / `dataLayout.Free` | `Close() error`（二次 `ErrClosed`） |
| `dataLayout` 全套查询 | `llvm.DataLayout`（`SizeOfType`→`SizeOfTypeInBits` 修正 bit 语义；新增 `PointerSizeForAS`/`IntPtrTypeForAS`/`ElementAtOffset`） |
| `Module.SetTarget` / `GetTarget` | `TargetMachine.SetTo(m)` / `target.FromTriple(m.TargetTriple())` |
| `Target.WriteASMToFile` / `WriteOBJToFile` | `TargetMachine.EmitToFile(m, path, AsmFile/ObjectFile)` 与 `Emit`（内存） |
| `Target.IsLinux/IsDarwin/IsWindows` | 删除（YAGNI，可自行 `strings.Contains(triple, ...)`） |
| `Module.getDataLayout`（借用引用） | `Module.DataLayout()`（owned 副本，避免 Dispose 借用句柄） |
| （无） | 新增 `ir.ParseIR` / `ir.ParseBitcode` / `Module.WriteToFile` / `Module.WriteBitcode` / `Module.Bitcode` / `llvm.ReadFile` / `llvm.NewMemoryBuffer` |

## 自检结果

1. **规格覆盖**：§4.7 的 Target 初始化/目标机器/emit 文件与内存/DataLayout 全套查询 → 任务 3/5/6；§5 P1-2 的 bitcode 读写、`ParseIR`、`MemoryBuffer` → 任务 1/2/4；§3.1 的 `llvm/target → llvm/ir` 依赖方向与 `EmitToFile(m *ir.Module)` → 任务 6 签名一致；§4.4 的独立所有权根（TargetMachine/MemoryBuffer 不经 `Own`）→ 任务 2/6 实现与测试覆盖。
2. **占位符**：无 TODO/待定；所有新增文件均有逐字代码，所有任务均有可运行测试与命令。
3. **类型一致性**：`MemoryBufferOf`/`DataLayoutOf` 桥接命名与既有 `ValueOf`/`TypeOfRef` 一致；`newModule` 同时服务 `NewModule`/`Clone`/`ParseIR`/`ParseBitcode`；`ErrParse` 在 error.go 定义、ir 使用；`checkModule` 在 target 包内单点定义、`SetTo`/`EmitToFile`/`Emit` 共用；binding 函数名与 LLVM-C 头文件（BitReader.h/BitWriter.h/IRReader.h/Core.h/TargetMachine.h）逐一对应。
