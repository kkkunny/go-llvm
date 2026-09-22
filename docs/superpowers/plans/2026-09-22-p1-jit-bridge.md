# P1-3 ~ P1-6 ORC LLJIT + Go 互调桥 实现计划（as-built）

> **状态：已完成。** 本文件按实现结果记录（步骤全部勾选），执行顺序与提交序列见文末。
> P1-1/P1-2（目标代码生成 + 解析/序列化）见 `2026-09-22-p1-codegen-parse.md`。

**目标：** 落地 `llvm/jit`（ORC LLJIT）、`internal/binding` ORC/桥绑定、Go↔native 互调桥（无汇编、全平台可移植），移除 legacy 执行引擎。

**架构：** `llvm/jit` 依赖 `llvm/target` → `llvm/ir` → `llvm`；ORC 句柄与桥通道全部封装在 `internal/binding`；按签名的 IR 适配器/包装体由 P0 Builder 生成，经 JIT 自身编译。

**技术栈：** Go 1.27、cgo、LLVM 22.1.8（`llvm-c/Orc.h`、`llvm-c/LLJIT.h`）、标准 `testing`。

---

## 文件结构（as-built）

```
internal/binding/Orc.go         ORC 句柄/枚举/结构 + LLJIT + 绝对符号 + ThreadSafeModule + orcError2Error
internal/binding/bridge.h/.c    固定签名通道：llvmBridgeCall / llvmBridgeGoChannel / llvmBridgeGoChannelAddr
internal/binding/bridge.go      BridgeCall / SetBridgeDispatch / CStringArray / //export goLLVMBridgeDispatch
context.go                      + Context.Disown（移交 JIT 时级联关闭子资源、不释放底层 Context）
memorybuffer.go                 + MemoryBuffer.Disown（ORC 消费缓冲时使 Go 句柄失效）
gotype.go                       + FnSignatureOfGo（reflect.Type 版签名映射）
ir/builder_call.go              + CallIndirect[U]（不透明函数指针 + 签名调用）
jit/lljit.go                    LLJIT：NewLLJIT/Close/AddIRModule/AddObjectFile/Lookup/MapSymbol/Triple/DataLayoutStr
jit/adapter.go                  槽位类别表、签名校验、按签名缓存、IR 适配器生成（unpackSlot/packSlot）
jit/bridge.go                   注册表、dispatchToGo、boxValue/unboxValue（reflect 派发）
jit/highlevel.go                Func[F] / MapFunc[F] / ensureGoChannel / compileWrapper / RunMain
jit/*_test.go                   LLJIT、桥、端到端（fib、Go 回调、OBJ 往返、IR 往返）
```

## 关键决策（as-built）

1. **所有权移交**：`LLVMOrcCreateNewThreadSafeModule` 接管 Module；ThreadSafeContext 接管 LLVMContext。`AddIRModule` 立即 `Module.Disown()` + `Context.Disown()`（Go 侧句柄全部失效），底层由 LLJIT 释放；`AddObjectFile` 消费 `MemoryBuffer`（`Disown`）。`Context.Disown` 会先级联关闭未移交的子资源（Builder 等），避免泄漏。
2. **桥 ABI**：适配器固定 `uint64_t adapter(void *fn, uint64_t *slots)`；native→Go 通道固定 `uint64_t callGoChannel(int64_t idx, uint64_t *slots)`；槽位 `[16]uint64`（每槽 8 字节，最多 15 参数），返回值即通道返回值。C→Go 走 `//export goLLVMBridgeDispatch` + `bridge.c` trampoline；Go→native 走 `bridge.c` 中的函数指针调用。
3. **适配器缓存**：按 `reflect.Type` 缓存于 `*LLJIT`；Go→native 适配器编译进 JIT 并 `Lookup` 缓存地址；native→Go 包装体以真实 native 签名定义符号 `name`，`callGoChannel` 经 `MapSymbol` 一次性挂接。
4. **首版桥类型**：bool/int32/uint32/int64/uint64/int/uint/float32/float64/指针（含 `unsafe.Pointer`、`*T`）；int8/int16 系、聚合、多返回值、变参 → `ErrUnsupported`（ABI 扩展属性归 P2-3）。
5. **错误归类**：ORC `LLVMErrorRef` 经 `LLVMGetErrorMessage` 消费后归入 `ErrJIT`；符号查找失败归 `ErrNotFound`；桥类型不支持归 `ErrUnsupported`。

## 任务分解（全部完成）

### 任务 8：ORC LLJIT 与符号定义绑定

- [x] `internal/binding/Orc.go`：句柄（LLJIT/JITDylib/ExecutionSession/SymbolStringPoolEntry/ResourceTracker/ThreadSafeModule/ThreadSafeContext/JITTargetMachineBuilder/MaterializationUnit）；`LLVMJITSymbolGenericFlags` 等枚举 `= C.Name`；`LLVMOrcCSymbolMapPair`/`LLVMJITEvaluatedSymbol` 结构映射；`LLVMOrcCreateLLJITBuilder`/`CreateLLJIT`/`DisposeLLJIT`/`LLJITGetExecutionSession`/`LLJITGetMainJITDylib`/`LLJITGetTripleString`/`LLJITGetDataLayoutStr`/`LLJITMangleAndIntern`/`LLJITAddLLVMIRModule`/`LLJITAddObjectFile`/`LLJITLookup`/`JITTargetMachineBuilderDetectHost`/`CreateNewThreadSafeContext(FromLLVMContext)`/`CreateNewThreadSafeModule`/`ExecutionSessionIntern`/`AbsoluteSymbols`/`JITDylibDefine`/`DisposeMaterializationUnit`/`orcError2Error`
- [x] `go build ./... && go vet ./...`

### 任务 9：llvm/jit LLJIT 封装

- [x] `Context.Disown` / `MemoryBuffer.Disown`（移交语义）
- [x] `jit/lljit.go`：`NewLLJIT`（DetectHost，ErrJIT）、`Close`（二次 ErrClosed）、`AddIRModule`（移交 Module+Context）、`AddObjectFile`（移交缓冲）、`Lookup`（ErrNotFound）、`MapSymbol`（MangleAndIntern + AbsoluteSymbols + Define）、`Triple`/`DataLayoutStr`
- [x] `jit/lljit_test.go`：Add+Lookup、MapSymbol 别名调用、AddObjectFile、Close/移交后校验

### 任务 10：Go 互调桥固定签名通道

- [x] `internal/binding/bridge.h`/`bridge.c`/`bridge.go`：`llvmBridgeCall`、`llvmBridgeGoChannel`、`llvmBridgeGoChannelAddr`、`//export goLLVMBridgeDispatch`、`SetBridgeDispatch`、`CStringArray`
- [x] `go build ./... && go vet ./...`

### 任务 11：按签名 IR 适配器与反射派发

- [x] `gotype.go` + `FnSignatureOfGo`；`ir/builder_call.go` + `CallIndirect[U]`
- [x] `jit/adapter.go`：槽位类别表（`slotI1/slotI32/slotI64/slotF32/slotF64/slotPtr`）、`checkBridgeFunc`、`adapterFor` 缓存、`compileAdapter`（解箱→间接调用→装箱）、`unpackSlot`/`packSlot`
- [x] `jit/bridge.go`：`registerGoFunc`/`dispatchToGo`/`boxValue`/`unboxValue`
- [x] `go test ./jit`

### 任务 12：JIT 高层 API

- [x] `jit/highlevel.go`：`(*LLJIT).Func[F]`（reflect.MakeFunc 真函数值）、`(*LLJIT).MapFunc[F]`（真实签名包装体 + `callGoChannel`）、`ensureGoChannel`、`compileWrapper`、`(*LLJIT).RunMain`
- [x] `jit/bridge_test.go`：Func 成功/ErrNotFound/ErrUnsupported、MapFunc 回调、float 签名、RunMain exit code

### 任务 13：移除 legacy 执行引擎与 gg 依赖

- [x] 删除 `internal/binding/ExecutionEngine.go`（MCJIT/Interpreter/GenericValue 全套）
- [x] `go mod tidy` 并删除 `go.sum`：`bytedance/gg` 依赖归零（go.mod 仅剩 module/go 两行）
- [x] `go build ./... && go vet ./... && go test ./...`

### 任务 14：端到端测试与文档收尾

- [x] `jit/e2e_test.go`：递归 fib（PHI+自引用调用）、Go 回调（MapFunc 经 JIT 调用宿主函数）、OBJ 往返（`Emit(ObjectFile)` → `AddObjectFile` → `Func`）、IR 往返（打印→ParseIR→JIT）
- [x] `README.md`：`llvm/jit` 去 P1 标记 + JIT execution 示例
- [x] `AGENTS.md`：jit 包边界、桥机制（`bridge.c` 通道）、ORC 所有权移交约定
- [x] 全量验收：`go build ./... && go vet ./... && go test ./...`（含 `-race`）全绿

## 提交序列（as-built）

```
feat: ORC LLJIT与符号定义绑定
feat: llvm/jit LLJIT封装与所有权移交
feat: Go互调桥固定签名通道与按签名IR适配器
refactor: 移除legacy执行引擎与gg依赖
test: JIT端到端测试与P1文档收尾
```

## 验收对照（规格 §5 P1）

| 验收项 | 覆盖 |
|---|---|
| JIT 跑通 fib | `TestE2EFib`、`TestE2EIRRoundTrip` |
| Go 回调 | `TestE2EGoCallback`、`TestLLJITMapFunc` |
| `Func[F]` 拿 Go 函数值 | `TestLLJITFunc`、`TestLLJITFuncFloatsAndPointers` |
| OBJ/ASM 产出正确 | `target/machine_test.go`、`TestE2EObjectRoundTrip` |
| go.mod 归零 | T13（无 require、无 go.sum） |
