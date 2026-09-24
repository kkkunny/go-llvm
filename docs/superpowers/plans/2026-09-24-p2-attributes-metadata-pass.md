# P2-3/P2-5/P2-6 补全实现计划（as-built）

> 本文件为事后回填的执行记录；实现与验证均已完成，提交见 git log。
> 设计依据：`docs/superpowers/specs/2026-09-22-go-llvm-inkwell-style-redesign-design.md` §5 的 P2-3/P2-5/P2-6。

**目标：** 补齐属性/调用约定、metadata/module flag/comdat/inline asm/blockaddress、诊断 handler/`Module.Link`/ctors-dtors/`llvm/pass`。

**验证命令：** `go build ./... && go vet ./... && go test ./... -count=1`

---

## 任务 A：P2-3 属性系统与调用约定

- [x] `internal/binding/Attribute.go`：`LLVMCallConv` 枚举、枚举/字符串/类型属性构造与查询、按位置（返回/函数/参数）增删查、调用点专用 `LLVM*CallSite*Attribute` 系列、函数/指令调用约定
- [x] root `attribute.go`：`Attribute`/`AttributeKind`/`AttrIndex`/`CallConv`、常用 kind 表（`AttributeKindForName` 动态解析；LLVM 22 已无 `argmemonly` 等旧名，`memory` 取代）、`EnumAttr`/`StringAttr`/`TypeAttr` 与 `AlignAttr`/`ByValAttr`/`SRetAttr` 糖
- [x] `ir/attrs.go`：`Function`/`Param`/`Call`/`Invoke` 四类角色的属性增删查与调用约定；内部按“函数 vs 调用点”分流（LLVM-C 两套 API 不能混用）
- [x] 测试 `ir/builder_attrs_test.go`：kind 解析、往返、打印、跨 Context 预检、调用点/Invoke 调用约定

## 任务 B：P2-5 metadata / module flag / comdat / inline asm / blockaddress

- [x] `internal/binding/Metadata.go`：MDString/MDNode/值互转、命名元数据、模块 flag、MDKind、指令元数据、blockaddress、inline asm；`internal/binding/Comdat.go`；`Types.go` 增加 metadata/comdat/named-md ref
- [x] root `metadata.go`：`Metadata`（`IsString`/`IsNode`/`Operands`/`StringValue`）、`ModuleFlagBehavior`；`inlineasm.go`：`InlineAsmDialect`/`Context.InlineAsm`；`comdat.go`：`ComdatSelectionKind`
- [x] `ir/metadata.go`：命名元数据、模块 flag、指令元数据附着（要求 MDNode）、`BlockAddress` 及反查；`ir/comdat.go`：`Comdat` 角色 + `Global.Comdat/SetComdat`
- [x] `Builder.Call`/`Invoke` 支持 InlineAsm 目标（新增 `callSig`：Function → `LLVMGetFunctionType`；InlineAsm → `LLVMGetInlineAsmFunctionType`）
- [x] 测试 `ir/metadata_test.go`：元数据往返/打印、指令附着、内联汇编调用、blockaddress、comdat

## 任务 C：P2-6 诊断 / Link / ctors-dtors / llvm/pass

- [x] `ErrorHandling.{h,c,go}`：诊断 trampoline + Go 注册表；root `diagnostic.go`：`DiagnosticSeverity`/`Context.SetDiagnosticHandler`/`ClearDiagnosticHandler`；`Core.cpp` 增加 `LLVMGoEmitError` shim 供测试触发（LLVM 默认 handler 对 error 会 `exit(1)`，测试只查回调指针）
- [x] `ir.Module.Link`：同 Context 前置校验；源模块被 LLVM 消费 → `Disown` 后 Go 侧句柄立即失效；失败归 `ErrLink`
- [x] `ir.AppendCtor/AppendDtor`：读取既有 `{i32, ptr, ptr}` 数组 + 删除重建 appending 全局
- [x] `pass` 包：`Level`、`Option`（VerifyEach/DebugLogging/LoopVectorization/SLPVectorization/LoopUnrolling/LoopInterleaving/ForgetAllSCEVInLoopUnroll/InlinerThreshold）、`RunPasses`/`RunPassesOnFunction`/`AutoOpt`；失败归 `ErrPass`
- [x] 测试：`diagnostic_test.go`、`ir/link_ctors_test.go`、`pass/pass_test.go`（O2 常量折叠、单函数管线、非法管线/级别、全选项）

## 备注（实现中确认的 LLVM 行为）

- LLVM-C 的 `LLVMAddAttributeAtIndex` 只接受 Function；调用点必须走 `LLVMAddCallSiteAttribute` 等专用 API。
- `align` 属性在 verifier 中只允许指针（及向量）类型。
- InlineAsm 值的 `LLVMTypeOf` 是不透明 `ptr`，签名须经 `LLVMGetInlineAsmFunctionType`。
- `LLVMSetMetadata` 要求 MDNode，裸 MDString 附着会产生非法/垃圾结果。
- 跨 Context 链接会产出 “Function context does not match Module context” 的非法模块，`Link` 予以拒绝。
- LLVM 22 枚举属性名称：`argmemonly`/`inaccessiblememonly`/`immutable` 已移除，`memory` 存在。
