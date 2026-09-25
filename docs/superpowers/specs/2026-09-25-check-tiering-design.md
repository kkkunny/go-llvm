# go-llvm 校验分层与 debug/release 双模式 —— 设计（已实施）

日期：2026-09-25
状态：已实施（第一期 M1-M3 核心 + 第二期 M4 核心；偏差见 §9）

---

## 1. 背景与问题

原实现为把异常暴露在 Go 侧（复用 panic/recover 与堆栈），在每个 API 入口做全量预检。
三个问题：

1. **效率**：`BenchmarkBuilderBinop` 366 ns/op 中，4 次 cgo 有 3 次纯为校验
   （`LLVMGetInsertBlock` + `LLVMTypeOf` ×2），校验让每条指令贵 2-3 倍。
2. **工作量/美观**：非测试代码显式校验调用 **342 处**（26 个文件）。
3. **原则性**：类型不一致、GEP 越界等语义契约，业界普遍交还 LLVM assert + `Verify`。

性能剖析结论：纯 Go 校验（判空/Lifetime）仅 ~3 ns，真正的成本是**校验路径上的 cgo
往返**；此外指令名参数（`const char *Name`）此前的 `C.CString` + `C.free` 每次额外
贡献 2 次 cgo 穿越，是最大的单项开销。

## 2. 调研：inkwell 的分层

1. 编译期强类型 + Rust 生命周期/所有权扛大头；
2. 廉价结构校验常开，`build_*` 返回 `Result<_, BuilderError>`（`UnsetPosition`/
   `NotSameType`/对齐/越界/原子序）；
3. 跨 Context 校验挂 `debug_assertions`（`check_val_context`）；
4. 查不了的标 `unsafe`（JIT 调用等）；深层 IR 正确性信任 LLVM，`module.verify()` 兜底
   （issue #661 明确讨论过"打印前自动 verify" vs 性能，选择了不自动查）。

关键差异：Rust FFI ~2 ns，inkwell 扛得起每次调用查两次类型；cgo ~90 ns/次，Go 绑定
必须把 cgo 从校验路径中消灭，并借用 build 模式区分层级。

## 3. 决策

| # | 决策 | 结果 |
|---|---|---|
| D-1 | 3+1' 修正版：release 保留纯 Go 崩溃类地板，语义类校验仅 debug | ✅ |
| D-2 | 建立 debug/release 双模式机制并启用 A-D 增强 | ✅ |
| D-3 | 零 cgo 校验重构（类型缓存/插入点跟踪/名字参数快路径） | ✅ |
| D-4 | 咽喉点收敛 | ⚠️ 部分（见 §9） |
| D-5 | 失败语义不变（程序员错误 panic `*llvm.Error`） | ✅ |
| D-6 | 开关用 build tag `llvm_release` + `const checks.Debug` | ✅ |
| D-7 | 默认 debug 全量，`-tags=llvm_release` 进入信任模式 | ✅ |

## 4. 三层模型（已落地）

```
第 1 层 编译期种类安全      Value[T]/Type[T]/ValueRef —— 任何模式
第 2 层 崩溃类地板（纯 Go） 判空/已释放/跨 Context/已关闭/插入点失效 —— 任何模式
第 3 层 语义契约与调试增强  类型一致/对齐/越界/原子序/实参/goroutine/Verify/... —— 仅 checks.Debug
```

- `internal/checks/checks.go`：`const Debug`，由 `debug.go`（`!llvm_release`）/
  `release.go`（`llvm_release`）两个载体文件选择；release 下第 3 层被编译器
  死代码消除。
- 负向测试：`ir/debug_test.go` 的 `requireDebug(t)` 在信任构建下跳过语义契约类误用测试；
  崩溃类地板测试两种构建都必须通过。

## 5. 实施内容

### 5.1 模式机制（M1）

- 新增 `internal/checks`；所有语义校验以 `if checks.Debug { ... }` 包裹：
  `preSameType`、`preAlign`、`checkKind`（Call/CallIndirect/Invoke/ExtractValue/
  ExtractElement/InsertElement）、`checkCallArgs` 的参数个数/类型部分、
  `elementTypeAt`（ExtractValue/InsertValue）、`preOrdering{Fence,RMW,CmpXchg,Load,Store}`、
  CondBr/Select 的 i1 校验、Switch case 类型、各类索引越界、`Function.Param` 越界。
- 无条件保留：插入点检查、`len(indices)==0`（防 Go 侧越界）、跨 Context 指针比较等。

### 5.2 零 cgo 校验（M2）

- `Value[T].ty`：调试层构造时缓存 `LLVMTypeOf`（`newValue` 统一入口；release 不查询、
  零额外开销）；`Value.Type()`/`As()`/`preSameType`/`checkCallArgs`/`typeString`
  在调试层全部零 cgo。
- `prePosition` 删除 `LLVMGetInsertBlock` cgo，改为 Go 侧插入点跟踪
  （`b.inserted == nil || !b.inserted.life.Alive()`）；公开 API 无块删除入口，
  检测能力不弱于原实现。
- `string2CString`：空名字走静态空 C 串，消除每次调用的 `CString`/`free` 两次
  cgo 穿越与 C 堆分配（实测 LLVM 22 对 `LLVMBuildAdd(..., NULL)` 会段错误，
  故不能用 NULL）。
- `Builder.valueRefs`/`gep` 的 scratch 复用：release 复用；debug 每次新分配（B4，
  把"只在紧随调用内消费"的约定违规变为不可能）。

### 5.3 调试层增强（M4）

| 项 | 实现 |
|---|---|
| A1 panic 现场 | `Builder` 维护最近 8 条 op（debug）；`b.panicf` 在消息中附加 `recent builder ops` |
| A3 诊断回调默认安装 | `NewContext` 默认安装 Go 回调（error/warning → 日志，不 panic、不退出）；`SetLogger(nil)` 静默 |
| B2 单 goroutine 契约 | `Builder`/`Module` 记录 owner goroutine id，采样（每 1024 次）+ 定位/关闭/移动插入点全量核对，跨 goroutine panic |
| B4 scratch 约定 | 见 5.2 |
| B5 边界 Verify | debug 下 `Module.Link`（双方）、`target.EmitToFile`、`jit.AddIRModule` 先 Verify，失败 panic `ErrVerify` |
| C1 资源审计 | debug 下 `Context.Close` 级联前报告未显式关闭的子资源数 |
| D2 benchmark 矩阵 | `make bench` / `make bench-release` |
| D3 golden 双模式一致 | 两种构建都跑 `TestGolden`（输出必须逐字节一致） |

### 5.4 工程化

- Makefile：`test` / `test-release` / `vet` / `bench` / `bench-release`。
- README：三层校验说明 + 信任构建用法；AGENTS.md：新增 *Checks and build modes* 约定、
  布局表加入 `internal/checks`、验证步骤加入 release 矩阵、性能约定更新。

## 6. 实测结果

| Benchmark | 改造前 | debug（默认） | release（`llvm_release`） |
|---|---|---|---|
| `BuilderBinop` | 366 ns/op | ~305 ns/op | ~255 ns/op |
| `BuilderCall` | 830 ns/op | ~730 ns/op | ~520 ns/op |

- release 模式为**真实地板**（语义校验编译期消除）；debug 与 release 差距约 20%（Binop）/40%（Call，剩余 cgo 为 `LLVMGetReturnType`/参数类型查询与调试层 scratch 分配）。
- 剖析显示剩余成本为：单次 cgo 穿越（~90 ns）+ LLVM C++ 指令分配 + API 边界的接口装箱（`ValueRef` 参数导致 2 allocs/op）。
  原设计估计"366 → ~150 ns"过于乐观：cgo 与装箱是 API 形状的固有成本，除非改变 `Value[T]` 的表示或方法签名。
- 指令名空串快路径（`string2CString` 静态空缓冲）单独贡献约 50-70 ns/op，是绑定层最大的单点收益。

## 7. 测试策略

- 两种构建全绿：`go test ./...` 与 `go test -tags=llvm_release ./...`。
- 语义契约负向测试挂 `requireDebug(t)`；地板负向测试（判空/UAF/跨 Context/双击关闭）两模式都跑。
- golden IR 双模式逐字节一致。

## 8. 明确不做（YAGNI）

- 双层 `Unchecked*` API（模式机制替代）。
- 运行时热切换（将来可把 `const checks.Debug` 换成 `var`，架构兼容）。
- 自动在 `Module.String()` 打印前 Verify：与"调试期导出构造中模块"的常规工作流冲突
  （inkwell 讨论同样结论），Verify 放在代码生成/链接边界。

## 9. 与初版设计的偏差

1. **D-4 咽喉点收敛已完成（后续补充）**：`Value.Ref()`/`Type.Ref()` 内建崩溃类地板校验
   （快路径保持可内联，慢路径 `checkFloor`），基于 `Value[T]`/`Type[T]` 的角色方法序言
   `Check` 全部删除（ir 包字段不可达，天然只能经 `Ref()` 触达句柄）；根包 `type.go`/
   `constant.go` 的 `.ref` 直读改为 `Ref()`。显式 `.Check(` 调用点 239 → 87，剩余集中在
   自持句柄的容器类型（`Block`/`Module`/`Comdat`/`Metadata`/`Attribute`）与参数校验
   （Argument `Attribute`），保留其显式 `Check`。`RawRef()` 仅供 `ir.preVal` 等预检内部
   使用，避免 `Ref()` 先 panic 丢失 builder 现场（op + recent ops）。实测两种构建
   benchmark 无回退。
2. **FnType 参数缓存未做**：调试层用 `LLVMGetParamTypes` 单次调用 + `Value.ty` 已把
   实参类型校验降到 1 次 cgo；`FnType` 全量缓存需要改所有构造路径，收益不足。
3. **A2 打印前 Verify 移除**：见 §8；B5 已覆盖真正的代码生成边界。
4. **B3 JIT 签名核对已完成（后续补充）**：调试构建在 `AddIRModule` 时记录模块内函数
   签名（`name → LLVM 类型文本`，迭代用新增的 `LLVMGetFirstFunction/NextFunction`
   binding；注意 LLVM 22 不透明指针下须用 `LLVMGetFunctionType` 而非 `LLVMTypeOf`），
   `Func[F]`/`MapFunc[F]` 注册时对照 `FnSignatureOfGo` 的映射文本，不匹配 panic
   `ErrTypeMismatch`；未知符号（`MapSymbol`/对象文件引入）跳过，不可过桥签名交由既有
   错误路径。release 构建不记录、不核对，零成本。
5. **异步/块级 Lifetime**：公开 API 无块删除路径，`prePosition` 用模块级令牌已足够；
   若未来暴露 `EraseBlock`，需补块专属令牌。
