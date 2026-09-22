# go-llvm 重构设计：对标 inkwell 的 Go 泛型 LLVM 封装

- 日期：2026-09-22
- 状态：待用户审查
- 目标：将现有 `github.com/kkkunny/go-llvm` 重构为类型安全、高度封装、Go 惯用的 LLVM 封装库，参考 `github.com/TheDan64/inkwell`

## 1. 背景

现有库为双层结构：`internal/binding`（LLVM-C 的 1:1 cgo 封装 + 少量 C++ shim）与 root `llvm` 包（手写封装，约 3000 行）。主要问题：

| 维度 | 现状 | 问题 |
|---|---|---|
| 类型安全 | 单一 `Value` 接口 + 类型断言 | `CreateSAdd` 传浮点值会触发 LLVM 内部崩溃，编译期零防护 |
| 错误处理 | 少量 `error`；未知枚举 `panic`；fatal handler 逻辑写反（`error.go:20`） | C++ 异常/fatal error 直接 abort 进程 |
| 生命周期 | 到处暴露 `Free()` | use-after-free / double-free 无检测 |
| API 覆盖 | IR 构建较全 | 缺 IR 解析、bitcode 读写、vector、原子指令、属性、metadata、Object |
| JIT | legacy MCJIT + `MapFunctionToGo` | LLVM 遗留 API；仅 amd64；反射细节裸露给用户 |
| cgo 边界 | `execution.c/h` 位于 root 包 | 违反 AGENTS.md「all cgo lives in internal/binding」 |

## 2. 已确认决策

| 决策点 | 结论 |
|---|---|
| API 兼容性 | 完全重构，不保留旧 API；`internal/binding` 保持 1:1 定位 |
| 类型安全机制 | 全泛型 `Value[T]` / `Type[T]`，允许使用 Go 1.27 泛型方法（方法自带类型参数） |
| 类型参数粒度 | **仅种类级**：T ∈ {VoidT, IntT, FloatT, PtrT, StructT, ArrayT, VecT, FnT, LabelT, MetaT, TokenT, DynT}；位宽、元素类型、函数签名运行时预检 |
| 错误处理 | 混合式：运行时可失败 → `error`；程序员错误 → `panic(*llvm.Error)`（可 recover）；C++ 异常经 shim 捕获转 Go error |
| 资源管理 | 分级所有权：资源对象实现 `io.Closer`；`Value[T]`/`Type[T]` 由 Context/Module 托管，内部检测 double-close 与 use-after-free |
| JIT 后端 | ORC LLJIT 优先；legacy MCJIT/Interpreter 删除 |
| Go 类型映射 | 纳入首批：`TypeOf[T]`、`ConstOf[T]`、`NewFunc[F]`、JIT `Func[F]` 返回真实 Go 函数值 |
| 架构 | **多子包**（见 §3） |
| Go 版本 | `go 1.27`（泛型方法），本机已验证 go1.27.1 |

## 3. 架构：多子包

### 3.1 包布局与依赖方向

```
用户代码
   │
   ├── llvm             核心词汇：Kind / Type[T] / Value[T] / 常量 / Context / DataLayout /
   │                    MemoryBuffer / Error / 生命周期 / Go 类型映射 / 版本常量
   ├── llvm/ir          Module / Function / Block / Global / Builder / 指令角色 /
   │                    Verify / Print / ParseIR / Bitcode 读写 / Link
   ├── llvm/target      Target / TargetMachine / CodeModel / RelocMode / OptLevel / emit OBJ·ASM
   ├── llvm/jit         LLJIT / Func[F] / MapFunc[F] / MapSymbol / RunMain
   ├── llvm/pass        PassBuilderOptions / RunPasses / AutoOpt
   ├── llvm/object      Object 文件读取（P3）
   └── llvm/debug       DIBuilder 调试信息（P3）

internal/binding       1:1 cgo（LLVM-C + C++ shim + bridge.c），全部 cgo 唯一归属地
   ↓
系统 libLLVM
```

依赖方向严格单向：`llvm` ← `llvm/ir` ← `llvm/target` ← `llvm/jit`，且 `llvm/ir` ← `llvm/pass`。
root `llvm` 包不 import 任何子包。**`llvm/target` 依赖 `llvm/ir`**：emit OBJ/ASM 与 `SetTarget`
便捷方法都必须拿到 `*ir.Module`（`LLVMTargetMachineEmitToFile` 的入参是模块）；`llvm/ir`
保持对 target/JIT 零依赖，`DataLayout`/`MemoryBuffer` 作为多子包共用的资源词汇放 root
（P0.5 决策，替代早期“ir 依赖 target”的方案）。

### 3.2 关键约束与解法

1. **公共 API 不得用泛型别名转发**：已验证泛型类型别名（`type Value[T] = core.Value[T]`）在 `go doc`/pkg.go.dev 中不展示方法，损害可发现性。因此共享泛型核心直接定义在 root `llvm` 包，子包直接 import 使用，不做别名转发。
2. **跨包生命周期令牌**：root 的 `Value[T]` 需要感知所属 Module 是否已关闭，但 Module 定义在 `llvm/ir`（root 不能 import 它）。解法：root 定义不透明生命周期令牌（如 `llvm.Lifetime`，含存活标志），`ir.Module` 持有一份并传给每个 `Value` 构造；`Value` 操作前检查 `ctx` 与 `lifetime` 存活状态。反向依赖被消除。
3. **cgo 归属**：桥接代码（固定签名 C 通道 + `//export` Go 回调）全部放 `internal/binding`；`internal/binding` 通过注册表回调 root 包（root 在 init 或首次使用时注册 handler），避免 internal → root 的 import。
- `DataLayout` 放 root：`ir.Module` 与 `target.TargetMachine` 都要用，且需要 `Type[T]` 参与查询；放 root 可避免 `ir` ↔ `target` 循环依赖。`MemoryBuffer` 同理放 root（ir 的 ParseIR/bitcode、jit 的 AddObjectFile、P3 的 object 都要用）。`SetTarget` 这类同时见 target 与 ir 的便捷方法放 `llvm/target`（target → ir 方向）。

### 3.3 各包职责边界

| 包 | 职责 | 明确不做 |
|---|---|---|
| `llvm` | 类型系统、值句柄、常量、Context、错误、生命周期、Go 类型映射、DataLayout、MemoryBuffer、版本 | 不涉及 Module/Builder/目标机器 |
| `llvm/ir` | 模块与 IR 构建、验证、打印、解析、bitcode、链接 | 不涉及目标机器与执行 |
| `llvm/target` | 目标注册与查询、目标机器、代码生成产出（`EmitToFile(m *ir.Module)`） | 不涉及执行 |
| `llvm/jit` | ORC LLJIT、符号查找与定义、Go 互调桥 | 不涉及 AOT 代码生成细节 |
| `llvm/pass` | 优化管线执行 | 不定义 pass 本身 |

## 4. 核心设计

### 4.1 类型系统（种类级泛型）

```go
// 封闭种类标记集（约束接口，不允许用户扩展）
type Kind interface{ kind() }
type VoidT, IntT, FloatT, PtrT, StructT, ArrayT, VecT, FnT, LabelT, MetaT, TokenT, DynT struct{}

type Type[T Kind]  struct { ref binding.LLVMTypeRef;  ctx *Context }
type Value[T Kind] struct { ref binding.LLVMValueRef; ctx *Context; life *Lifetime }

// 常用别名（非泛型，方法直接可见）
type IntValue = Value[IntT]; type FloatValue = Value[FloatT]; type PtrValue = Value[PtrT]
type DynValue = Value[DynT]   // 类型擦除逃生舱：IR 解析、指令遍历、动态查找
```

- 依赖类型：`func (v Value[T]) Type() Type[T]`。
- Go 1.27 泛型方法（仅具体结构体可用）：
  - `func (v Value[T]) As[U Kind]() (Value[U], error)`、`MustAs[U Kind]() Value[U]`
  - `func (b Builder) Load[U Kind](p PtrValue, t Type[U], name string) Value[U]`
  - `func (b Builder) Select[T Kind](cond IntValue, x, y Value[T], name string) Value[T]`
  - `func (m *Module) NewFunc[F any](name string) (Func[F], error)`
- 角色包装承载指令专属操作，Builder 返回角色（均嵌入 `Value[T]`）：`Phi[T]`、`Switch`、`Call[T]`、`Alloca`、`Load[T]`、`Store`、`GEP` 等。
- 未知 kind 降级为 `Value[DynT]`，不 panic。
- 混合集合用仅含非泛型方法的小接口 `AnyValue` / `AnyType`；配套 `As[T]` / `AsDyn` 包级泛型函数。

### 4.2 Builder 与预检

```go
b.Add/Sub/Mul(l, r IntValue, name string) IntValue          // 另有 AddNSW / AddNUW 等变体
b.SDiv/UDiv/FDiv/Shl/LShr/AShr/And/Or/Xor/...               // IntT 约束；F* 系约束 FloatT
b.ICmp(op IntPred, l, r AnyValue, name string) IntValue     // 返回 i1
b.Alloca(t AnyType, name string) PtrValue
b.Store(v AnyValue, p PtrValue)
b.GEP(elem AnyType, p PtrValue, idx []IntValue, name string) PtrValue
b.ZExt/SExt/Trunc(v IntValue, to Type[IntT], name string) IntValue   // 目标种类从 to 推导
b.BitCast[U Kind](v AnyValue, to Type[U], name string) Value[U]
b.Call[U Kind](fn Value[FnT], args []AnyValue, name string) Value[U] // 返回种类调用前预检
b.Ret(v AnyValue) / b.Br / b.CondBr(IntValue, ...) / b.Switch / b.Unreachable
b.PHI[T Kind](t Type[T], name string) Phi[T] / b.Select[T Kind](...)
```

统一预检框架（进 cgo 前执行，失败 `panic(*Error)` 并带操作与定位信息）：

1. Builder 已定位（`CurrentBlock` 非空）；
2. 所有值与类型属于同一 Context，且生命周期令牌存活；
3. 语义约束：操作数种类一致、call 实参个数与函数类型匹配、对齐为 2 的幂、名字非空等。

命名约定：去掉旧 `Create` 前缀（`b.Add`），root 及子包注释沿用中文。

### 4.3 错误模型

```go
type Error struct { Reason ErrKind; Op, Msg string }   // 实现 error
// ErrTypeMismatch / ErrCrossContext / ErrUseAfterFree / ErrClosed / ErrNotFound /
// ErrInvalidArg / ErrVerify / ErrUnsupported / ErrCodeGen / ErrJIT / ErrIO / ErrInternal
func Catch(fn func()) *Error      // recover(*Error) → error
func Must[T any](v T, err error) T
func WrapError(reason ErrKind, op string, err error) *Error  // binding 裸 error → *Error
```

- 运行时可失败操作 → `error`；程序员错误 → `panic(*Error)`（可 recover 转 error）。
- C++ 异常：shim 以 `-fexceptions` 编译，入口 `try/catch (...)` → 错误码/消息 → Go error。
- LLVM `Error`/`Expected`（ORC 等）经 `LLVMErrorRef` → Go error。
- fatal error（`report_fatal_error`）：语义上退出进程，跨 cgo 回调不可 recover。防线为：前置校验 + shim 异常捕获消灭绝大多数 abort 路径；handler 记录诊断、执行用户 hook 后退出；文档如实说明。
- `Module.Verify()` 经 diagnostic handler 收集结构化错误（P2-6）。

### 4.4 生命周期

- `*Context`、`*Module`、`*Builder`、`*LLJIT`、`*TargetMachine`、`*MemoryBuffer` 实现 `io.Closer`；二次 `Close` 返回 `ErrClosed`。
- `Value[T]` / `Type[T]` / `Block` / `Func[F]` 不暴露 Free；持有 `ctx` + `life` 令牌，操作前检查存活 → use-after-free 即时 `panic(*Error)`。
- `Context` 为所有权根：`Own` 返回注销函数（资源自行关闭/移交时解除登记），`Close` 逆序级联释放全部；`Module.Close()` 使其下值失效。
- **所有权转移**：`Module.Disown()` 从 Context 解除登记并返回 `release`；JIT 接管模块（ThreadSafeModule）时调用，由接管方在释放底层模块时执行 `release`，避免与 `Context.Close` 双重释放。`TargetMachine`/`MemoryBuffer`/`LLJIT` 不隶属任何 `Context`，是各自独立的所有权根，不经 `Own` 登记。

### 4.5 Go 类型映射（首批）

```go
func TypeOf[T any](ctx *Context) (Type[DynT], error)          // bool→i1, intN/uintN→iN, float64→double,
                                                              // *T/unsafe.Pointer→ptr, 标量字段 struct→literal struct
func ConstOf[T any](ctx *Context, v T) (DynValue, error)
func (m *Module) NewFunc[F any](name string) (Func[F], error) // Go 签名→FunctionType
```

首版仅支持标量/指针/bool；按值 struct、string、slice、map、chan、func、interface → 明确返回 `ErrUnsupported`（提前报错）。

### 4.6 JIT 与 Go 互调桥（无汇编、全平台可移植）

```go
func (e *LLJIT) Func[F any](name string) (F, error)      // reflect.MakeFunc 返回真实 Go 函数值
func (e *LLJIT) MapFunc[F any](name string, f F) error   // 宿主 Go 函数注册为 JIT 符号（替代 MapFunctionToGo）
func (e *LLJIT) MapSymbol(name string, p unsafe.Pointer) error
```

机制（泛化现有 `MapFunctionToGo` 的「IR 装箱 + 固定签名通道」思路）：

- **Go→native**（`Func[F]` 调用）：`reflect.MakeFunc` 闭包装箱实参 → cgo 调固定签名 `callBridge`（C）→ 按签名缓存的 **IR 适配器**（用 P0 的 Builder 生成，间接调用目标并还原实参）。
- **native→Go**（`MapFunc`）：IR 包装体装箱实参 → 固定签名 `callGoChannel(idx, slots)`（C → `//export` → reflect 调用用户 Go 函数）；经 ORC `DefineAbsoluteSymbols` 挂接。
- 适配器按签名缓存，首次调用编译一次；无需按架构手写汇编。
- 低层逃生舱：仍提供 `unsafe.Pointer` 直取符号，供热路径直接按 C ABI 调用。

### 4.7 目标、代码生成与优化

- `Target` 初始化（native / 按 arch）、`TargetMachine(triple, cpu, features, opt, reloc, codeModel)`、emit OBJ/ASM 到文件与 `[]byte`。
- `DataLayout` 全套查询（现 `target.go` 迁移）。
- `RunPasses(pipeline string) error`、`AutoOpt(level)`（shim 异常捕获）。
- bitcode 读写、`ParseIR`、`MemoryBuffer`。

## 5. 大步骤分期

> 只拆大步骤；细分步骤在执行各阶段时再拆分。

### P0 — 类型系统骨架 + 高频 IR 构建
验收：golden IR 单测通过；README 示例以新 API 端到端跑通；`go build ./...` + `go vet ./...` 干净。

- P0-1 基建：`go.mod` → 1.27；`llvm` 包骨架、`Error`/`Catch`/`Must`、生命周期令牌；修复 fatal handler 逻辑
- P0-2 类型系统核心：`Kind` / `Type[T]` / `Value[T]` / `As` / lookup 分发 / `Dyn` 降级
- P0-3 上下文与常量：`Context`、常量全套（Int/Float/String/Null/Zero/Array/Struct/GEP/表达式高频）
- P0-4 IR 构建：`llvm/ir` 的 `Module`/`Function`/`Block`/`Global`/`Builder` + 高频指令 + 统一预检框架
- P0-5 Go 类型映射声明侧：`TypeOf` / `ConstOf` / `NewFunc[F]`
- P0-6 验证与测试基建：`Verify`（结构化错误）/`Print`/`SetSource`、golden 测试、CI 等价命令

### P1 — 代码生成 + 解析 + ORC JIT
验收：JIT 跑通 fib、Go 回调、`Func[F]` 拿 Go 函数值；OBJ/ASM 产出正确。

- P1-1 目标与 DataLayout：`llvm/target` 迁移重写
- P1-2 序列化与解析：bitcode 读写、`ParseIR`、`MemoryBuffer`
- P1-3 ORC LLJIT 绑定与封装：`internal/binding` ORC 全套 + `llvm/jit` 基础
- P1-4 Go 互调桥：`bridge.c` 固定签名通道 + 按签名 IR 适配器生成
- P1-5 JIT 高层 API：`Func[F]` / `MapFunc[F]` / `MapSymbol` / `RunMain`
- P1-6 移除 legacy 执行引擎（MCJIT/Interpreter/GenericValue，含 `internal/binding/ExecutionEngine.go`）+ 端到端测试；随后移除仅供其使用的 `bytedance/gg` 依赖（go.mod 归零）

### P2 — 中频补全
验收：各指令/属性有 golden 测试；旧 API 能力无回退。

- P2-1 异常处理指令全套（invoke/landingpad/catchpad/cleanupret/resume…）
- P2-2 原子与内存序（fence/atomicrmw/cmpxchg、volatile、对齐属性）
- P2-3 属性系统与调用约定全集（函数/参数/调用点/返回值/属性组）
- P2-4 vector 与聚合指令补全（shufflevector/insertelement/insertvalue 等）
- P2-5 metadata / module flag / comdat / inline asm / blockaddress
- P2-6 诊断 handler 对外暴露、`Module.Link`、ctors/dtors、`llvm/pass` 完善

### P3 — 冷门功能
- P3-1 `llvm/debug`：DIBuilder 调试信息
- P3-2 `llvm/object`：Object 文件读取（sections/symbols/relocs）
- P3-3 Disassembler
- P3-4 ORC 高级（ResourceTracker/generators）与 Analysis/Transforms 杂项

## 6. 迁移与文档

- 旧 root 包 API 全部删除；`execution.c/h`、`execution_amd64.go` 机制并入 `internal/binding` 桥。
- README 更新为新 API 示例与包布局说明；AGENTS.md 增补泛型约定、子包边界、桥机制。
- 版本线维护流程（LLVM 升级）保持 README 现有约定。

## 7. 风险

| 风险 | 应对 |
|---|---|
| `report_fatal_error` 跨 cgo 不可 recover | 前置校验 + shim 异常捕获为主；handler 记录诊断后退出；文档如实说明 |
| 反射桥每次调用有装箱开销 | 适配器按签名缓存；提供 `unsafe.Pointer` 低层入口给热路径 |
| Go→C ABI 首版仅标量/指针 | 明确 `ErrUnsupported`，后续按需扩展（按值聚合、string 约定） |
| 多子包引入的导入仪式感 | root 包承载核心词汇，用户主要 import `llvm` + `llvm/ir`；后端按需引入 |
| ORC C API 随 LLVM 版本演进 | 绑定按需增量；沿用 AGENTS.md 的枚举 `= C.Name` 与头文件核对流程 |

## 8. 待确认项

1. Builder 方法去掉 `Create` 前缀（`b.Add` 而非 `b.CreateSAdd`）——默认采纳，如有异议请提出。
2. legacy MCJIT/Interpreter 直接删除、不做兼容层——默认采纳。
