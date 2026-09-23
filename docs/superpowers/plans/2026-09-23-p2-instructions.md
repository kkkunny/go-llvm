# P2-1/P2-2/P2-4 指令补全（EH / 原子 / 向量）实现计划

> **面向 AI 代理的工作者：** 必需子技能：使用 superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans 逐任务实现此计划。步骤使用复选框（`- [ ]`）语法来跟踪进度。

**目标：** 为 `ir` 包补齐 P2 指令面——EH 全套（invoke/landingpad/funclet）、原子指令与内存序（fence/atomicrmw/cmpxchg、volatile）、向量指令（extractelement/insertelement/shufflevector）与 `ConstVector` 常量。

**架构：** 沿用既定分层：`internal/binding` 只做 LLVM-C 1:1 绑定（枚举 `= C.Name`），root `llvm` 提供公共枚举（`AtomicOrdering`/`RMWOp`）与常量，`ir` 提供 Builder 方法与指令角色（内嵌 `Value[T]`）。预检在进 cgo 前完成（种类/类型/内存序合法性/块归属），失败 `panic(*llvm.Error)`。EH 的 personality 经 `Function.SetPersonality` 设置，landingpad 的 `PersFn` 参数传零值（LLVM 22 已忽略）。

**技术栈：** Go 1.27、cgo、LLVM 22.1.8（`llvm-c/Core.h`）、标准 `testing`、golden IR 测试。

---

## 已实证的 LLVM 22 行为（探针结论，勿凭记忆改动）

以下均为本机 LLVM 22.1.8 实测结论，实现与测试以此为准：

1. `LLVMBuildInvoke2` 的第二个参数是**函数类型**（不是返回类型）；传错类型直接段错误。现有绑定 `LLVMBuildInvoke(builder, ty, fn, args, then, catch, name)` 中的 `ty` 即函数类型，与 `LLVMBuildCall` 的 `sig.Ref()` 用法一致。
2. `LLVMBuildLandingPad(B, Ty, PersFn=NULL, NumClauses=0, Name)` 安全可用；personality 由 `LLVMSetPersonalityFn(fn, pers)` 挂在函数上（IR 形如 `define ... personality ptr @pers`）。
3. `invoke` 的 `LLVMGetSuccessor(inst, 0)` = 正常出口块，`(inst, 1)` = unwind 出口块；`LLVMGetNumArgOperands`/`LLVMGetCalledValue`/`LLVMGetOperand` 对 invoke 均可用（LLVM 侧同为 CallBase），`Invoke` 角色可安全复用 `Call[T]` 的实参操作。
4. `catchpad` 的 `LLVMGetNumOperands` = 实参个数 + 1（末位操作数是 parent pad）；`LLVMGetArgOperand(pad, i)` 直接按下标取实参（i=0 即第一个实参）；`LLVMGetParentCatchSwitch(catchpad)` 可取父 catchswitch。
5. `fence` 返回 void 值——**命名 void 值会被 verifier 拒绝**（"Instruction has a name, but provides a void value!"）。因此 `Fence`/`Resume`/`CatchRet`/`CleanupRet` 等 void 指令的 Builder 方法**不带 name 参数**（与既有 `Store`/`Ret` 一致）。
6. `atomicrmw`/`cmpxchg` 的 C 构建 API 无 name 参数；构建后 `LLVMSetValueName` 命名可行（IR 打印 `%cx = cmpxchg ...`）。二者可带 name（返回值非 void）。
7. EH 结构过 verifier 的合法形态（golden 测试必须照此构建）：
   - EH pad（catchswitch/catchpad/cleanuppad）**不能位于 entry 块**；
   - EH pad **必须由 unwind 边进入**（invoke 的 unwind label，或 catchswitch 的 handler）；
   - `catchpad` 的目标块经 `catchret` 进入；`cleanuppad` 同理。
   - 已验证通过 verifier 的三例见任务 5 的 golden 构建代码。
8. 返回种类：`fence`→void、`atomicrmw`→操作数同种类（如 i32）、`cmpxchg`→`{T, i1}` 字面量结构体（StructT）、`landingpad`→所给类型（StructT）、`invoke`→函数返回种类。

---

## 文件结构

| 文件 | 动作 | 职责 |
|---|---|---|
| `internal/binding/Core.go` | 增补 | EH 子句/handler/funclet/personality 绑定、原子绑定、vector 绑定、`LLVMAtomicOrdering`/`LLVMAtomicRMWBinOp` 枚举（`= C.Name`） |
| `atomic.go` | 新建 | root 公共枚举 `AtomicOrdering`/`RMWOp`（转发 binding 值） |
| `ir/function.go` | 增补 | `Function.SetPersonality`/`Personality` |
| `ir/builder_eh.go` | 新建 | `Invoke`/`LandingPad`/`CatchSwitch`/`FuncletPad` 角色 + `Invoke`/`InvokeIndirect`/`LandingPad`/`Resume`/`CatchSwitch`/`CatchPad`/`CleanupPad`/`CatchRet`/`CleanupRet` |
| `ir/builder_atomic.go` | 新建 | `Fence`/`AtomicRMW[T]`/`CmpXchg` 角色 + `Fence`/`AtomicRMW`/`CmpXchg` + 内存序预检 |
| `ir/builder_mem.go` | 增补 | `Load[T]`/`Store` 角色加 `SetVolatile`/`IsVolatile`/`SetOrdering`/`Ordering` |
| `ir/builder_vec.go` | 新建 | `ExtractElement`/`InsertElement`/`ShuffleVector` |
| `constant.go` | 增补 | `Context.ConstVector` |
| `ir/builder_call.go` | 增补 | `checkCallArgs` 共享助手（Invoke 复用 call 的实参校验） |
| `ir/zz_refcheck_test.go` | 增补 | 新角色的 `ValueRef`/`AnyValue` 断言 |
| `ir/builder_eh_test.go` | 新建 | EH 构建/角色/预检测试 |
| `ir/builder_atomic_test.go` | 新建 | 原子构建/角色/预检测试 |
| `ir/builder_vec_test.go` | 新建 | 向量构建/预检测试 |
| `constant_test.go` | 增补 | `TestConstVector` |
| `ir/golden_test.go` + `ir/testdata/golden/*.ll` | 增补 | `TestGoldenEH`/`TestGoldenAtomic`/`TestGoldenVec` + 对应 golden 文件 |
| `AGENTS.md` | 增补 | 「Builder return types」角色清单更新 |

---

## 任务 1：binding EH 增补（子句/handler/funclet/personality）

**文件：**
- 修改：`internal/binding/Core.go`（EH 构建函数区域之后，约 `LLVMAddCase` 附近）
- 修改：`ir/function.go`（`SetLinkage` 之后）
- 测试：`ir/builder_eh_test.go`（本任务只测 personality 往返）

- [x] **步骤 1：编写失败的测试**

创建 `ir/builder_eh_test.go`：

```go
package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestFunctionPersonality(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "eh")
	defer m.Close()

	i32 := ctx.Int(32)
	pers := m.NewFunction("pers", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))

	fn.SetPersonality(pers)
	got, ok := fn.Personality()
	if !ok {
		t.Fatal("personality should be set")
	}
	if got.Name() != "pers" {
		t.Fatalf("personality = %s, want pers", got.Name())
	}
	if !strings.Contains(m.String(), "personality ptr @pers") {
		t.Fatalf("module should contain personality:\n%s", m.String())
	}

	if err := llvm.Catch(func() { fn.SetPersonality(llvm.Value[llvm.FnT]{}) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("nil personality should panic ErrInvalidArg, got %v", err)
	}
}
```

- [x] **步骤 2：运行测试验证失败**

运行：`go test ./ir -run TestFunctionPersonality -v`
预期：编译失败，报错 `fn.SetPersonality undefined`（type Function has no field or method SetPersonality）。

- [x] **步骤 3：编写最少实现代码**

`internal/binding/Core.go` 在 `LLVMAddCase` 之后插入（英文注释风格与周围一致）：

```go
// LLVMGetNumClauses Get the number of clauses on the landingpad instruction.
func LLVMGetNumClauses(landingPad LLVMValueRef) uint32 {
	return uint32(C.LLVMGetNumClauses(landingPad.c))
}

// LLVMGetClause Get the value of the clause at index Idx on the landingpad instruction.
func LLVMGetClause(landingPad LLVMValueRef, idx uint32) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetClause(landingPad.c, C.unsigned(idx))}
}

// LLVMAddClause Add a catch or filter clause to the landingpad instruction.
func LLVMAddClause(landingPad, clauseVal LLVMValueRef) {
	C.LLVMAddClause(landingPad.c, clauseVal.c)
}

// LLVMIsCleanup Get the 'cleanup' flag in the landingpad instruction.
func LLVMIsCleanup(landingPad LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsCleanup(landingPad.c))
}

// LLVMSetCleanup Set the 'cleanup' flag in the landingpad instruction.
func LLVMSetCleanup(landingPad LLVMValueRef, val bool) {
	C.LLVMSetCleanup(landingPad.c, bool2LLVMBool(val))
}

// LLVMAddHandler Add a destination to the catchswitch instruction.
func LLVMAddHandler(catchSwitch LLVMValueRef, dest LLVMBasicBlockRef) {
	C.LLVMAddHandler(catchSwitch.c, dest.c)
}

// LLVMGetNumHandlers Get the number of handlers on the catchswitch instruction.
func LLVMGetNumHandlers(catchSwitch LLVMValueRef) uint32 {
	return uint32(C.LLVMGetNumHandlers(catchSwitch.c))
}

// LLVMGetHandlers Obtain the basic blocks acting as handlers for a catchswitch instruction.
func LLVMGetHandlers(catchSwitch LLVMValueRef) []LLVMBasicBlockRef {
	n := int(LLVMGetNumHandlers(catchSwitch))
	if n == 0 {
		return nil
	}
	handlers := make([]C.LLVMBasicBlockRef, n)
	C.LLVMGetHandlers(catchSwitch.c, &handlers[0])
	refs := make([]LLVMBasicBlockRef, n)
	for i, h := range handlers {
		refs[i] = LLVMBasicBlockRef{c: h}
	}
	return refs
}

// LLVMGetArgOperand Get a funcletpad argument at the given index.
func LLVMGetArgOperand(funclet LLVMValueRef, i uint32) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetArgOperand(funclet.c, C.unsigned(i))}
}

// LLVMSetArgOperand Set a funcletpad argument at the given index.
func LLVMSetArgOperand(funclet LLVMValueRef, i uint32, value LLVMValueRef) {
	C.LLVMSetArgOperand(funclet.c, C.unsigned(i), value.c)
}

// LLVMGetParentCatchSwitch Get the parent catchswitch instruction of a catchpad instruction.
func LLVMGetParentCatchSwitch(catchPad LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetParentCatchSwitch(catchPad.c)}
}

// LLVMGetPersonalityFn Get the personality function attached to the function.
func LLVMGetPersonalityFn(fn LLVMValueRef) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMGetPersonalityFn(fn.c)}
}

// LLVMSetPersonalityFn Set the personality function attached to the function.
func LLVMSetPersonalityFn(fn, persFn LLVMValueRef) {
	C.LLVMSetPersonalityFn(fn.c, persFn.c)
}
```

`ir/function.go` 在 `SetLinkage` 之后插入：

```go
// SetPersonality 设置 personality 函数（EH 展开用）
func (f Function) SetPersonality(pers llvm.ValueRef[llvm.FnT]) {
	const op = "ir.Function.SetPersonality"
	f.Check(op)
	pv := pers.AsValue()
	if pv.IsNil() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil personality")
	}
	f.Context().CheckValues(op, pv)
	binding.LLVMSetPersonalityFn(f.Ref(), pv.Ref())
}

// Personality personality 函数（未设置时返回 false）
func (f Function) Personality() (llvm.Value[llvm.FnT], bool) {
	f.Check("ir.Function.Personality")
	ref := binding.LLVMGetPersonalityFn(f.Ref())
	if ref.IsNil() {
		return llvm.Value[llvm.FnT]{}, false
	}
	return llvm.NewValue[llvm.FnT](f.Context(), f.Lifetime(), ref), true
}
```

- [x] **步骤 4：运行测试验证通过**

运行：`go test ./ir -run TestFunctionPersonality -v && go build ./... && go vet ./...`
预期：PASS，build/vet 干净。

- [x] **步骤 5：Commit**

```bash
git add internal/binding/Core.go ir/function.go ir/builder_eh_test.go
git commit -m "feat(ir): personality 绑定与 Function.SetPersonality/Personality"
```

---

## 任务 2：Invoke/InvokeIndirect + Invoke 角色

**文件：**
- 修改：`internal/binding/Core.go`（已有的 `LLVMBuildInvoke`/`LLVMBuildLandingPad` 等 EH 构建绑定不动；本任务纯 ir 层）
- 新建：`ir/builder_eh.go`
- 修改：`ir/builder_call.go`（抽出 `checkCallArgs`）
- 修改：`ir/zz_refcheck_test.go`
- 测试：`ir/builder_eh_test.go`

- [x] **步骤 1：编写失败的测试**

`ir/builder_eh_test.go` 追加（需用合法 EH CFG——见探针结论 7，unwind 目标先用 `Unreachable` 占位、任务 3 换 landingpad）：

```go
func buildEHModule(t *testing.T) (*llvm.Context, *Module, *Builder) {
	t.Helper()
	ctx := llvm.NewContext()
	m := NewModule(ctx, "eh")
	b := NewBuilder(ctx)
	return ctx, m, b
}

func TestBuilderInvoke(t *testing.T) {
	ctx, m, b := buildEHModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	g := m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	entry := fn.NewBlock("entry")
	cont := fn.NewBlock("cont")
	lpad := fn.NewBlock("lpad")

	b.MoveToEnd(entry)
	iv := b.Invoke[llvm.IntT](g, []llvm.AnyValue{fn.ParamAs[llvm.IntT](0)}, cont, lpad, "v")

	if iv.NormalBlock() != cont {
		t.Fatalf("normal block = %s, want cont", iv.NormalBlock().Name())
	}
	if iv.UnwindBlock() != lpad {
		t.Fatalf("unwind block = %s, want lpad", iv.UnwindBlock().Name())
	}
	if iv.ArgCount() != 1 {
		t.Fatalf("arg count = %d, want 1", iv.ArgCount())
	}
	if called, ok := iv.CalledFunction(); !ok || called.Name() != "g" {
		t.Fatalf("called function = %v %v, want g", called, ok)
	}

	b.MoveToEnd(cont)
	b.Ret(iv)

	b.MoveToEnd(lpad)
	b.Unreachable() // 任务 3 替换为 landingpad+resume

	got := m.String()
	if !strings.Contains(got, "invoke i32 @g(i32 %0)") {
		t.Fatalf("missing invoke:\n%s", got)
	}
	if !strings.Contains(got, "to label %cont unwind label %lpad") {
		t.Fatalf("missing successors:\n%s", got)
	}
}

func TestBuilderInvokePrecheck(t *testing.T) {
	ctx, m, b := buildEHModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	g := m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn := m.NewFunction("f", ctx.Fn(i32, nil, false))
	entry := fn.NewBlock("entry")
	cont := fn.NewBlock("cont")
	lpad := fn.NewBlock("lpad")
	b.MoveToEnd(entry)

	// 未定位
	b2 := NewBuilder(ctx)
	defer b2.Close()
	if err := llvm.Catch(func() {
		b2.Invoke[llvm.IntT](g, []llvm.AnyValue{fn.ParamAs[llvm.IntT](0)}, cont, lpad, "")
	}); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("unpositioned invoke should panic ErrInvalidArg, got %v", err)
	}

	// 返回种类不符
	if err := llvm.Catch(func() {
		b.Invoke[llvm.FloatT](g, []llvm.AnyValue{fn.ParamAs[llvm.IntT](0)}, cont, lpad, "")
	}); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("wrong return kind should panic ErrTypeMismatch, got %v", err)
	}

	// 实参个数不符
	if err := llvm.Catch(func() {
		b.Invoke[llvm.IntT](g, nil, cont, lpad, "")
	}); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("wrong arg count should panic ErrTypeMismatch, got %v", err)
	}
}
```

- [x] **步骤 2：运行测试验证失败**

运行：`go test ./ir -run TestBuilderInvoke -v`
预期：编译失败，报错 `b.Invoke undefined`（type *Builder has no method Invoke）。

- [x] **步骤 3：编写最少实现代码**

`ir/builder_call.go`：把 `call` 方法里的实参校验段抽出为共享助手（`call` 改为调用它，行为不变）：

```go
// checkCallArgs 调用类指令（call/invoke）公共实参预检：个数/类型与签名匹配
func (b *Builder) checkCallArgs(op string, sig llvm.FnType, args []llvm.AnyValue) {
	for _, a := range args {
		b.checkVal(op, coreAny(a))
	}
	n := binding.LLVMCountParamTypes(sig.Ref())
	got := uint(len(args))
	if got < uint(n) {
		llvm.Panicf(llvm.ErrTypeMismatch, op, "expect at least %d arguments, got %d", n, got)
	}
	if got > uint(n) && !sig.IsVarArg() {
		llvm.Panicf(llvm.ErrTypeMismatch, op, "expect %d arguments, got %d", n, got)
	}
	if n > 0 {
		params := binding.LLVMGetParamTypes(sig.Ref())
		for i, p := range params {
			if !p.Equal(binding.LLVMTypeOf(args[i].Ref())) {
				llvm.Panicf(llvm.ErrTypeMismatch, op, "argument %d type %s does not match parameter type %s",
					i, typeString(b.ctx, args[i].Ref()), typeRefString(b.ctx, p))
			}
		}
	}
}
```

`call` 相应变为：

```go
// call 调用公共路径：实参预检后发指令（callee 已由调用方校验）
func (b *Builder) call(op string, callee binding.LLVMValueRef, sig llvm.FnType, args []llvm.AnyValue, name string) binding.LLVMValueRef {
	b.pre(op)
	b.checkCallArgs(op, sig, args)
	return binding.LLVMBuildCall(b.ref, sig.Ref(), callee, b.valueRefs(args), name)
}
```

新建 `ir/builder_eh.go`：

```go
package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Invoke invoke 指令角色（内嵌 Value[T]）；实参操作同 Call（LLVM 侧同为 CallBase），
// 另有正常/异常出口块访问
type Invoke[T llvm.Kind] struct {
	llvm.Value[T]
}

// NormalBlock 正常出口块
func (c Invoke[T]) NormalBlock() Block {
	c.Check("ir.Invoke.NormalBlock")
	return wrapBlock(c.Context(), c.Lifetime(), binding.LLVMGetSuccessor(c.Ref(), 0))
}

// UnwindBlock 异常出口块
func (c Invoke[T]) UnwindBlock() Block {
	c.Check("ir.Invoke.UnwindBlock")
	return wrapBlock(c.Context(), c.Lifetime(), binding.LLVMGetSuccessor(c.Ref(), 1))
}

// ArgCount 实参个数
func (c Invoke[T]) ArgCount() uint32 {
	c.Check("ir.Invoke.ArgCount")
	return binding.LLVMGetNumArgOperands(c.Ref())
}

// Arg 第 i 个实参（擦除种类）
func (c Invoke[T]) Arg(i uint32) llvm.Value[llvm.DynT] {
	const op = "ir.Invoke.Arg"
	c.Check(op)
	if i >= c.ArgCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "argument index %d out of range", i)
	}
	return llvm.ValueOf(c.Context(), c.Lifetime(), binding.LLVMGetOperand(c.Ref(), i))
}

// SetArg 替换第 i 个实参
func (c Invoke[T]) SetArg(i uint32, v llvm.AnyValue) {
	const op = "ir.Invoke.SetArg"
	c.Check(op)
	if i >= c.ArgCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "argument index %d out of range", i)
	}
	c.Context().CheckValues(op, v)
	binding.LLVMSetOperand(c.Ref(), i, v.Ref())
}

// CalledFunction 被调用函数（非间接 invoke 时）
func (c Invoke[T]) CalledFunction() (llvm.Value[llvm.FnT], bool) {
	c.Check("ir.Invoke.CalledFunction")
	ref := binding.LLVMGetCalledValue(c.Ref())
	if ref.IsNil() {
		return llvm.Value[llvm.FnT]{}, false
	}
	return llvm.NewValue[llvm.FnT](c.Context(), c.Lifetime(), ref), true
}

// ===== EH 构建方法 =====

// Invoke 插入 invoke 调用：then 为正常出口、unwind 为异常出口；
// 返回种类 U 在调用前与函数返回类型比对，不符 panic
func (b *Builder) Invoke[U llvm.Kind](fn llvm.ValueRef[llvm.FnT], args []llvm.AnyValue, then, unwind Block, name string) Invoke[U] {
	const op = "ir.Builder.Invoke"
	fv := fn.AsValue()
	b.pre(op, core(fv))
	b.preBlockOwn(op, then)
	b.preBlockOwn(op, unwind)
	sig := llvm.AsFnType(llvm.TypeOfRef(b.ctx, binding.LLVMGetFunctionType(fv.Ref())))
	checkKind[U](op, b.ctx, binding.LLVMGetReturnType(sig.Ref()))
	b.checkCallArgs(op, sig, args)
	ref := binding.LLVMBuildInvoke(b.ref, sig.Ref(), fv.Ref(), b.valueRefs(args), then.ref, unwind.ref, name)
	return Invoke[U]{Value: llvm.NewValue[U](b.ctx, b.inserted.life, ref)}
}

// InvokeIndirect 通过函数指针 invoke（不透明指针 + 签名）；返回种类 U 与签名返回类型比对
func (b *Builder) InvokeIndirect[U llvm.Kind](fnPtr llvm.ValueRef[llvm.PtrT], sig llvm.FnType, args []llvm.AnyValue, then, unwind Block, name string) Invoke[U] {
	const op = "ir.Builder.InvokeIndirect"
	pv := fnPtr.AsValue()
	b.pre(op, core(pv))
	b.preBlockOwn(op, then)
	b.preBlockOwn(op, unwind)
	if sig.Context() != b.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "signature belongs to another context")
	}
	checkKind[U](op, b.ctx, binding.LLVMGetReturnType(sig.Ref()))
	b.checkCallArgs(op, sig, args)
	ref := binding.LLVMBuildInvoke(b.ref, sig.Ref(), pv.Ref(), b.valueRefs(args), then.ref, unwind.ref, name)
	return Invoke[U]{Value: llvm.NewValue[U](b.ctx, b.inserted.life, ref)}
}
```

`ir/zz_refcheck_test.go` 追加断言：

```go
	_ llvm.ValueRef[llvm.IntT] = Invoke[llvm.IntT]{}
	_ llvm.AnyValue            = Invoke[llvm.IntT]{}
```

- [x] **步骤 4：运行测试验证通过**

运行：`go test ./ir -run 'TestBuilderInvoke|TestFunctionPersonality' -v && go build ./... && go vet ./...`
预期：PASS。注意 `TestBuilderInvoke` 里 lpad 块暂为 `Unreachable`（无 landingpad 时 invoke 的 unwind 边语义不完整，但 `Module.Verify` 在任务 3 前不加入断言；本任务测试不调 `Verify`）。

- [x] **步骤 5：Commit**

```bash
git add ir/builder_call.go ir/builder_eh.go ir/zz_refcheck_test.go ir/builder_eh_test.go
git commit -m "feat(ir): Invoke/InvokeIndirect 与 Invoke 角色"
```

---

## 任务 3：LandingPad/Resume + personality 联通

**文件：**
- 修改：`ir/builder_eh.go`
- 修改：`ir/zz_refcheck_test.go`
- 测试：`ir/builder_eh_test.go`

- [x] **步骤 1：编写失败的测试**

`ir/builder_eh_test.go` 追加：

```go
func TestBuilderLandingPadResume(t *testing.T) {
	ctx, m, b := buildEHModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	pers := m.NewFunction("pers", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	g := m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	ti := m.NewGlobal("ti", ptr)

	fn := m.NewFunction("f", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	fn.SetPersonality(pers)
	entry := fn.NewBlock("entry")
	cont := fn.NewBlock("cont")
	lpad := fn.NewBlock("lpad")

	b.MoveToEnd(entry)
	iv := b.Invoke[llvm.IntT](g, []llvm.AnyValue{fn.ParamAs[llvm.IntT](0)}, cont, lpad, "v")

	b.MoveToEnd(cont)
	b.Ret(iv)

	b.MoveToEnd(lpad)
	padTy := ctx.Struct([]llvm.AnyType{ptr, i32}, false)
	lp := b.LandingPad(padTy, "lp")
	lp.AddClause(ti)
	lp.AddClause(ti)
	lp.SetCleanup(true)
	if lp.ClauseCount() != 2 {
		t.Fatalf("clause count = %d, want 2", lp.ClauseCount())
	}
	if got := lp.Clause(0); got.Name() != "ti" {
		t.Fatalf("clause 0 = %s, want ti", got.Name())
	}
	if !lp.IsCleanup() {
		t.Fatal("cleanup should be set")
	}
	lp.SetCleanup(false)
	if lp.IsCleanup() {
		t.Fatal("cleanup should be cleared")
	}
	b.Resume(lp)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := m.String()
	if !strings.Contains(got, "landingpad { ptr, i32 }") {
		t.Fatalf("missing landingpad:\n%s", got)
	}
	if !strings.Contains(got, "catch ptr @ti") {
		t.Fatalf("missing catch clause:\n%s", got)
	}
	if !strings.Contains(got, "resume { ptr, i32 } %lp") {
		t.Fatalf("missing resume:\n%s", got)
	}
}

func TestBuilderLandingPadPrecheck(t *testing.T) {
	ctx, m, b := buildEHModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))
	entry := fn.NewBlock("entry")
	b.MoveToEnd(entry)

	padTy := ctx.Struct([]llvm.AnyType{ptr, i32}, false)
	lp := b.LandingPad(padTy, "lp")

	// 子句必须是常量
	if err := llvm.Catch(func() { lp.AddClause(fn.Param(0)) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("non-constant clause should panic ErrInvalidArg, got %v", err)
	}
	// 子句下标越界
	if err := llvm.Catch(func() { lp.Clause(9) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("clause out of range should panic ErrInvalidArg, got %v", err)
	}
	b.Unreachable()
}
```

- [x] **步骤 2：运行测试验证失败**

运行：`go test ./ir -run 'TestBuilderLandingPad|TestBuilderInvoke' -v`
预期：编译失败，报错 `b.LandingPad undefined`。同时把任务 2 测试里 lpad 的 `b.Unreachable()` 替换为 `b.LandingPad(...)+b.Resume(...)` 后 `Verify` 断言——一并在此步骤改为：

```go
	b.MoveToEnd(lpad)
	padTy := ctx.Struct([]llvm.AnyType{ctx.Ptr(0), i32}, false)
	lp := b.LandingPad(padTy, "lp")
	b.Resume(lp)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
```

- [x] **步骤 3：编写最少实现代码**

`ir/builder_eh.go` 追加（角色定义放文件顶部角色区，构建方法放构建区）：

```go
// LandingPad landingpad 指令角色（内嵌 Value[T]，T 通常为 StructT）
type LandingPad[T llvm.Kind] struct {
	llvm.Value[T]
}

// AddClause 追加 catch/filter 子句（catch：类型信息全局；filter：常量数组）
func (l LandingPad[T]) AddClause(v llvm.AnyValue) {
	const op = "ir.LandingPad.AddClause"
	l.Check(op)
	l.Context().CheckValues(op, v)
	if !v.IsConstant() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "clause must be a constant")
	}
	binding.LLVMAddClause(l.Ref(), v.Ref())
}

// ClauseCount 子句数量
func (l LandingPad[T]) ClauseCount() uint32 {
	l.Check("ir.LandingPad.ClauseCount")
	return binding.LLVMGetNumClauses(l.Ref())
}

// Clause 第 i 条子句（擦除种类）
func (l LandingPad[T]) Clause(i uint32) llvm.Value[llvm.DynT] {
	const op = "ir.LandingPad.Clause"
	l.Check(op)
	if i >= l.ClauseCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "clause index %d out of range", i)
	}
	return llvm.ValueOf(l.Context(), l.Lifetime(), binding.LLVMGetClause(l.Ref(), i))
}

// SetCleanup 设置 cleanup 标志
func (l LandingPad[T]) SetCleanup(v bool) {
	l.Check("ir.LandingPad.SetCleanup")
	binding.LLVMSetCleanup(l.Ref(), v)
}

// IsCleanup 是否 cleanup
func (l LandingPad[T]) IsCleanup() bool {
	l.Check("ir.LandingPad.IsCleanup")
	return binding.LLVMIsCleanup(l.Ref())
}
```

构建方法：

```go
// LandingPad 插入 landingpad（t 须为首类聚合类型）；personality 经 Function.SetPersonality 设置。
// 注：LLVM 22 的 PersFn 构建参数已废弃，传零值
func (b *Builder) LandingPad[T llvm.Kind](t llvm.TypeRef[T], name string) LandingPad[T] {
	const op = "ir.Builder.LandingPad"
	tt := t.AsType()
	b.pre(op)
	if tt.Context() != b.ctx {
		llvm.Panicf(llvm.ErrCrossContext, op, "type belongs to another context")
	}
	ref := binding.LLVMBuildLandingPad(b.ref, tt.Ref(), binding.LLVMValueRef{}, 0, name)
	return LandingPad[T]{Value: llvm.NewValue[T](b.ctx, b.inserted.life, ref)}
}

// Resume 以 landingpad 值恢复异常传播（void 指令，无 name）
func (b *Builder) Resume(exn llvm.AnyValue) llvm.Value[llvm.VoidT] {
	const op = "ir.Builder.Resume"
	b.pre(op, coreAny(exn))
	ref := binding.LLVMBuildResume(b.ref, exn.Ref())
	return llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)
}
```

`ir/zz_refcheck_test.go` 追加：

```go
	_ llvm.ValueRef[llvm.StructT] = LandingPad[llvm.StructT]{}
	_ llvm.AnyValue               = LandingPad[llvm.StructT]{}
```

- [x] **步骤 4：运行测试验证通过**

运行：`go test ./ir -v && go build ./... && go vet ./...`
预期：全绿（含任务 2 测试的 Verify 断言）。

- [x] **步骤 5：Commit**

```bash
git add ir/builder_eh.go ir/zz_refcheck_test.go ir/builder_eh_test.go
git commit -m "feat(ir): LandingPad/Resume 与 landingpad 子句操作"
```

---

## 任务 4：funclet 全套（CatchSwitch/CatchPad/CleanupPad/CatchRet/CleanupRet）

**文件：**
- 修改：`ir/builder_eh.go`
- 修改：`ir/zz_refcheck_test.go`
- 测试：`ir/builder_eh_test.go`

- [x] **步骤 1：编写失败的测试**

`ir/builder_eh_test.go` 追加（合法 funclet CFG——探针结论 7 的过验形态）：

```go
func TestBuilderFunclets(t *testing.T) {
	ctx, m, b := buildEHModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	pers := m.NewFunction("pers", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	h := m.NewFunction("h", ctx.Fn(ctx.Void(), nil, false))
	ti := m.NewGlobal("ti", ptr)

	// ---- catchswitch -> catchpad -> catchret ----
	fn := m.NewFunction("funclets", ctx.Fn(ctx.Void(), nil, false))
	fn.SetPersonality(pers)
	entry := fn.NewBlock("entry")
	cont := fn.NewBlock("cont")
	dispatch := fn.NewBlock("dispatch")
	hnd := fn.NewBlock("hnd")
	done := fn.NewBlock("done")

	b.MoveToEnd(entry)
	b.Invoke[llvm.VoidT](h, nil, cont, dispatch, "") // EH pad 需 unwind 边进入

	b.MoveToEnd(cont)
	b.RetVoid()

	b.MoveToEnd(dispatch) // EH pad 不得在 entry 块
	cs := b.CatchSwitch(nil, Block{}, "cs") // parent=nil 即 within none；unwindTo 零块即 unwind to caller
	cs.AddHandler(hnd)
	if cs.HandlerCount() != 1 {
		t.Fatalf("handler count = %d, want 1", cs.HandlerCount())
	}
	if cs.HandlerAt(0) != hnd {
		t.Fatalf("handler 0 = %s, want hnd", cs.HandlerAt(0).Name())
	}

	b.MoveToEnd(hnd)
	cp := b.CatchPad(cs, []llvm.AnyValue{ti}, "cp")
	if got := cp.ParentCatchSwitch(); got.Name() != "cs" {
		t.Fatalf("parent catchswitch = %s, want cs", got.Name())
	}
	if cp.ArgCount() != 1 {
		t.Fatalf("catchpad arg count = %d, want 1", cp.ArgCount())
	}
	if got := cp.Arg(0); got.Name() != "ti" {
		t.Fatalf("catchpad arg 0 = %s, want ti", got.Name())
	}
	b.CatchRet(cp, done)

	b.MoveToEnd(done)
	b.RetVoid()

	// ---- cleanuppad -> cleanupret ----
	kn := m.NewFunction("cleanup", ctx.Fn(ctx.Void(), nil, false))
	kn.SetPersonality(pers)
	e2 := kn.NewBlock("entry")
	c2 := kn.NewBlock("cont")
	clean := kn.NewBlock("clean")

	b.MoveToEnd(e2)
	b.Invoke[llvm.VoidT](h, nil, c2, clean, "")
	b.MoveToEnd(c2)
	b.RetVoid()

	b.MoveToEnd(clean)
	clp := b.CleanupPad(nil, nil, "clp")
	if clp.ArgCount() != 0 {
		t.Fatalf("cleanuppad arg count = %d, want 0", clp.ArgCount())
	}
	b.CleanupRet(clp, Block{}) // unwind to caller

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := m.String()
	if !strings.Contains(got, "catchswitch within none [label %hnd] unwind to caller") {
		t.Fatalf("missing catchswitch:\n%s", got)
	}
	if !strings.Contains(got, "catchpad within %cs [ptr @ti]") {
		t.Fatalf("missing catchpad:\n%s", got)
	}
	if !strings.Contains(got, "catchret from %cp to label %done") {
		t.Fatalf("missing catchret:\n%s", got)
	}
	if !strings.Contains(got, "cleanuppad within none []") {
		t.Fatalf("missing cleanuppad:\n%s", got)
	}
	if !strings.Contains(got, "cleanupret from %clp unwind to caller") {
		t.Fatalf("missing cleanupret:\n%s", got)
	}
}

func TestBuilderFuncletPrecheck(t *testing.T) {
	ctx, m, b := buildEHModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))
	entry := fn.NewBlock("entry")
	b.MoveToEnd(entry)

	cs := b.CatchSwitch(nil, Block{}, "cs")

	// handler 下标越界
	if err := llvm.Catch(func() { cs.HandlerAt(3) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("handler out of range should panic ErrInvalidArg, got %v", err)
	}
	// funclet pad 实参下标越界
	pad := b.CatchPad(cs, nil, "cp")
	if err := llvm.Catch(func() { pad.Arg(0) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("pad arg out of range should panic ErrInvalidArg, got %v", err)
	}
	b.Unreachable()
}
```

- [x] **步骤 2：运行测试验证失败**

运行：`go test ./ir -run 'TestBuilderFunclet' -v`
预期：编译失败，报错 `b.CatchSwitch undefined`。

- [x] **步骤 3：编写最少实现代码**

`ir/builder_eh.go` 追加角色：

```go
// CatchSwitch catchswitch 指令角色（内嵌 Value[TokenT]）
type CatchSwitch struct {
	llvm.Value[llvm.TokenT]
}

// AddHandler 追加处理器入口块（须与 catchswitch 同函数）
func (s CatchSwitch) AddHandler(blk Block) {
	const op = "ir.CatchSwitch.AddHandler"
	s.Check(op)
	blk.Check(op)
	binding.LLVMAddHandler(s.Ref(), blk.ref)
}

// HandlerCount 处理器数量
func (s CatchSwitch) HandlerCount() uint32 {
	s.Check("ir.CatchSwitch.HandlerCount")
	return binding.LLVMGetNumHandlers(s.Ref())
}

// HandlerAt 第 i 个处理器入口块
func (s CatchSwitch) HandlerAt(i uint32) Block {
	const op = "ir.CatchSwitch.HandlerAt"
	s.Check(op)
	if i >= s.HandlerCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "handler index %d out of range", i)
	}
	return wrapBlock(s.Context(), s.Lifetime(), binding.LLVMGetHandlers(s.Ref())[i])
}

// FuncletPad funclet pad 指令角色（catchpad/cleanuppad，内嵌 Value[TokenT]）
type FuncletPad struct {
	llvm.Value[llvm.TokenT]
}

// ArgCount 实参个数（LLVMGetNumOperands 含末位 parent pad，故减 1）
func (p FuncletPad) ArgCount() uint32 {
	p.Check("ir.FuncletPad.ArgCount")
	return uint32(binding.LLVMGetNumOperands(p.Ref())) - 1
}

// Arg 第 i 个实参（擦除种类）
func (p FuncletPad) Arg(i uint32) llvm.Value[llvm.DynT] {
	const op = "ir.FuncletPad.Arg"
	p.Check(op)
	if i >= p.ArgCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "argument index %d out of range", i)
	}
	return llvm.ValueOf(p.Context(), p.Lifetime(), binding.LLVMGetArgOperand(p.Ref(), i))
}

// SetArg 替换第 i 个实参
func (p FuncletPad) SetArg(i uint32, v llvm.AnyValue) {
	const op = "ir.FuncletPad.SetArg"
	p.Check(op)
	if i >= p.ArgCount() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "argument index %d out of range", i)
	}
	p.Context().CheckValues(op, v)
	binding.LLVMSetArgOperand(p.Ref(), i, v.Ref())
}

// ParentCatchSwitch 所属 catchswitch（仅 catchpad 有意义）
func (p FuncletPad) ParentCatchSwitch() CatchSwitch {
	p.Check("ir.FuncletPad.ParentCatchSwitch")
	return CatchSwitch{Value: llvm.NewValue[llvm.TokenT](p.Context(), p.Lifetime(), binding.LLVMGetParentCatchSwitch(p.Ref()))}
}
```

构建方法：

```go
// CatchSwitch 插入 catchswitch（终结指令）；parent 为 nil 表示 within none，
// unwindTo 为零块表示 unwind to caller。处理器用 CatchSwitch.AddHandler 追加
func (b *Builder) CatchSwitch(parent llvm.ValueRef[llvm.TokenT], unwindTo Block, name string) CatchSwitch {
	const op = "ir.Builder.CatchSwitch"
	b.pre(op)
	parentRef := b.preParentPad(op, parent)
	unwindRef := b.preOptBlock(op, unwindTo)
	ref := binding.LLVMBuildCatchSwitch(b.ref, parentRef, unwindRef, 0, name)
	return CatchSwitch{Value: llvm.NewValue[llvm.TokenT](b.ctx, b.inserted.life, ref)}
}

// CatchPad 插入 catchpad（须为块首指令）；parent 为 nil 表示 within none
func (b *Builder) CatchPad(parent llvm.ValueRef[llvm.TokenT], args []llvm.AnyValue, name string) FuncletPad {
	const op = "ir.Builder.CatchPad"
	b.pre(op)
	parentRef := b.preParentPad(op, parent)
	for _, a := range args {
		b.checkVal(op, coreAny(a))
	}
	ref := binding.LLVMBuildCatchPad(b.ref, parentRef, b.valueRefs(args), name)
	return FuncletPad{Value: llvm.NewValue[llvm.TokenT](b.ctx, b.inserted.life, ref)}
}

// CleanupPad 插入 cleanuppad（须为块首指令）；parent 为 nil 表示 within none
func (b *Builder) CleanupPad(parent llvm.ValueRef[llvm.TokenT], args []llvm.AnyValue, name string) FuncletPad {
	const op = "ir.Builder.CleanupPad"
	b.pre(op)
	parentRef := b.preParentPad(op, parent)
	for _, a := range args {
		b.checkVal(op, coreAny(a))
	}
	ref := binding.LLVMBuildCleanupPad(b.ref, parentRef, b.valueRefs(args), name)
	return FuncletPad{Value: llvm.NewValue[llvm.TokenT](b.ctx, b.inserted.life, ref)}
}

// CatchRet 从 catchpad 转移到目标块（void 终结指令，无 name）
func (b *Builder) CatchRet(pad FuncletPad, to Block) llvm.Value[llvm.VoidT] {
	const op = "ir.Builder.CatchRet"
	b.pre(op, core(pad.Value))
	b.preBlockOwn(op, to)
	ref := binding.LLVMBuildCatchRet(b.ref, pad.Ref(), to.ref)
	return llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)
}

// CleanupRet 从 cleanuppad 转移；unwindTo 为零块表示 unwind to caller（void 终结指令，无 name）
func (b *Builder) CleanupRet(pad FuncletPad, unwindTo Block) llvm.Value[llvm.VoidT] {
	const op = "ir.Builder.CleanupRet"
	b.pre(op, core(pad.Value))
	unwindRef := b.preOptBlock(op, unwindTo)
	ref := binding.LLVMBuildCleanupRet(b.ref, pad.Ref(), unwindRef)
	return llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)
}

// preParentPad 预检可选父 pad（nil = within none）；返回底层句柄（零值 = none）
func (b *Builder) preParentPad(op string, parent llvm.ValueRef[llvm.TokenT]) binding.LLVMValueRef {
	if parent == nil {
		return binding.LLVMValueRef{}
	}
	pv := parent.AsValue()
	b.checkVal(op, core(pv))
	return pv.Ref()
}

// preOptBlock 预检可选目标块（零块 = to/unwind to caller）；返回底层句柄（零值 = caller）
func (b *Builder) preOptBlock(op string, blk Block) binding.LLVMBasicBlockRef {
	if blk.ref.IsNil() {
		return binding.LLVMBasicBlockRef{}
	}
	b.preBlockOwn(op, blk)
	return blk.ref
}
```

`ir/zz_refcheck_test.go` 追加：

```go
	_ llvm.ValueRef[llvm.TokenT] = CatchSwitch{}
	_ llvm.ValueRef[llvm.TokenT] = FuncletPad{}
	_ llvm.AnyValue              = CatchSwitch{}
	_ llvm.AnyValue              = FuncletPad{}
```

- [x] **步骤 4：运行测试验证通过**

运行：`go test ./ir -v && go build ./... && go vet ./...`
预期：全绿，`TestBuilderFunclets` 的 `Verify` 通过（合法 funclet CFG）。

- [x] **步骤 5：Commit**

```bash
git add ir/builder_eh.go ir/zz_refcheck_test.go ir/builder_eh_test.go
git commit -m "feat(ir): funclet 全套（catchswitch/catchpad/catchret/cleanuppad/cleanupret）"
```

---

## 任务 5：EH golden 端到端

**文件：**
- 修改：`ir/golden_test.go`
- 新建：`ir/testdata/golden/eh.ll`（由 `-update` 生成）

- [x] **步骤 1：编写失败的测试**

`ir/golden_test.go` 追加：

```go
// TestGoldenEH 覆盖 Itanium 式（invoke/landingpad/resume）与 funclet 式（catchswitch/catchpad/catchret/cleanuppad/cleanupret）
func TestGoldenEH(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "eh")
	defer m.Close()

	i32 := ctx.Int(32)
	ptr := ctx.Ptr(0)
	pers := m.NewFunction("pers", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	g := m.NewFunction("g", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	h := m.NewFunction("h", ctx.Fn(ctx.Void(), nil, false))
	ti := m.NewGlobal("ti", ptr)

	b := NewBuilder(ctx)
	defer b.Close()

	// ---- itanium ----
	f1 := m.NewFunction("itanium", ctx.Fn(i32, []llvm.AnyType{i32}, false))
	f1.SetPersonality(pers)
	e1 := f1.NewBlock("entry")
	c1 := f1.NewBlock("cont")
	l1 := f1.NewBlock("lpad")

	b.MoveToEnd(e1)
	v := b.Invoke[llvm.IntT](g, []llvm.AnyValue{f1.ParamAs[llvm.IntT](0)}, c1, l1, "v")
	b.MoveToEnd(c1)
	b.Ret(v)
	b.MoveToEnd(l1)
	lp := b.LandingPad(ctx.Struct([]llvm.AnyType{ptr, i32}, false), "lp")
	lp.AddClause(ti)
	lp.SetCleanup(true)
	b.Resume(lp)

	// ---- funclets ----
	f2 := m.NewFunction("funclets", ctx.Fn(ctx.Void(), nil, false))
	f2.SetPersonality(pers)
	e2 := f2.NewBlock("entry")
	c2 := f2.NewBlock("cont")
	d2 := f2.NewBlock("dispatch")
	hn := f2.NewBlock("hnd")
	dn := f2.NewBlock("done")

	b.MoveToEnd(e2)
	b.Invoke[llvm.VoidT](h, nil, c2, d2, "")
	b.MoveToEnd(c2)
	b.RetVoid()
	b.MoveToEnd(d2)
	cs := b.CatchSwitch(nil, Block{}, "cs")
	cs.AddHandler(hn)
	b.MoveToEnd(hn)
	cp := b.CatchPad(cs, []llvm.AnyValue{ti}, "cp")
	b.CatchRet(cp, dn)
	b.MoveToEnd(dn)
	b.RetVoid()

	// ---- cleanup ----
	f3 := m.NewFunction("cleanup", ctx.Fn(ctx.Void(), nil, false))
	f3.SetPersonality(pers)
	e3 := f3.NewBlock("entry")
	c3 := f3.NewBlock("cont")
	k3 := f3.NewBlock("clean")

	b.MoveToEnd(e3)
	b.Invoke[llvm.VoidT](h, nil, c3, k3, "")
	b.MoveToEnd(c3)
	b.RetVoid()
	b.MoveToEnd(k3)
	clp := b.CleanupPad(nil, nil, "clp")
	b.CleanupRet(clp, Block{})

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	checkGolden(t, m.String(), "eh.ll")
}
```

- [x] **步骤 2：运行测试验证失败**

运行：`go test ./ir -run TestGoldenEH -v`
预期：FAIL，报错 `read golden (可用 -update 生成): open testdata/golden/eh.ll: no such file or directory`。

- [x] **步骤 3：生成 golden 并核对**

运行：`go test ./ir -run TestGoldenEH -update && cat ir/testdata/golden/eh.ll`
核对生成文件包含以下形态（顺序/空白以生成为准）：

```
define i32 @itanium(i32 %0) personality ptr @pers {
entry:
  %v = invoke i32 @g(i32 %0)
          to label %cont unwind label %lpad
cont:
  ret i32 %v
lpad:
  %lp = landingpad { ptr, i32 }
          cleanup
          catch ptr @ti
  resume { ptr, i32 } %lp
define void @funclets() personality ptr @pers {
  %cs = catchswitch within none [label %hnd] unwind to caller
  %cp = catchpad within %cs [ptr @ti]
  catchret from %cp to label %done
define void @cleanup() personality ptr @pers {
  %clp = cleanuppad within none []
  cleanupret from %clp unwind to caller
```

若差异仅为空白/注释排布，接受生成结果；若语句缺失或结构不同，回头检查构建代码。

- [x] **步骤 4：运行测试验证通过**

运行：`go test ./ir -run TestGolden -v`
预期：`TestGoldenMain`/`TestGoldenLoop`/`TestGoldenEH` 全部 PASS。

- [x] **步骤 5：Commit**

```bash
git add ir/golden_test.go ir/testdata/golden/eh.ll
git commit -m "test(ir): EH 端到端 golden（itanium/funclet/cleanup）"
```

---

## 任务 6：原子枚举 + Fence/AtomicRMW/CmpXchg

**文件：**
- 修改：`internal/binding/Core.go`（枚举定义 + 原子绑定，紧跟 `LLVMBuildFCmp` 附近）
- 新建：`atomic.go`（root 公共枚举）
- 新建：`ir/builder_atomic.go`
- 修改：`ir/zz_refcheck_test.go`
- 测试：`ir/builder_atomic_test.go`

- [x] **步骤 1：编写失败的测试**

新建 `ir/builder_atomic_test.go`：

```go
package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func buildAtomicModule(t *testing.T) (*llvm.Context, *Module, *Builder) {
	t.Helper()
	ctx := llvm.NewContext()
	m := NewModule(ctx, "atomic")
	b := NewBuilder(ctx)
	return ctx, m, b
}

func TestBuilderAtomics(t *testing.T) {
	ctx, m, b := buildAtomicModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{ctx.Ptr(0)}, false))
	entry := fn.NewBlock("entry")
	b.MoveToEnd(entry)
	p := fn.ParamAs[llvm.PtrT](0)
	one := ctx.ConstInt(i32, 1).Value

	fe := b.Fence(llvm.AtomicSequentiallyConsistent, false)
	if fe.Ordering() != llvm.AtomicSequentiallyConsistent {
		t.Fatalf("fence ordering = %v", fe.Ordering())
	}
	fe.SetOrdering(llvm.AtomicAcquire)
	if fe.Ordering() != llvm.AtomicAcquire {
		t.Fatalf("fence ordering after set = %v", fe.Ordering())
	}

	rm := b.AtomicRMW(llvm.RMWAdd, p, one, llvm.AtomicMonotonic, false, "old")
	if rm.Ordering() != llvm.AtomicMonotonic || rm.Op() != llvm.RMWAdd {
		t.Fatalf("rmw = %v %v", rm.Ordering(), rm.Op())
	}
	rm.SetVolatile(true)
	rm.SetAlign(4)
	if !rm.IsVolatile() || rm.Align() != 4 {
		t.Fatalf("rmw volatile/align = %v %d", rm.IsVolatile(), rm.Align())
	}

	cx := b.CmpXchg(p, one, ctx.ConstInt(i32, 2).Value, llvm.AtomicAcquire, llvm.AtomicMonotonic, true, "cx")
	if cx.SuccessOrdering() != llvm.AtomicAcquire || cx.FailureOrdering() != llvm.AtomicMonotonic {
		t.Fatalf("cmpxchg orderings = %v %v", cx.SuccessOrdering(), cx.FailureOrdering())
	}
	if !cx.IsWeak() {
		t.Fatal("cmpxchg should be weak")
	}
	cx.SetWeak(false)
	cx.SetAlign(4)

	b.RetVoid()
	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := m.String()
	for _, want := range []string{
		"fence acquire",
		"atomicrmw volatile add i32 %0, i32 1 monotonic, align 4",
		"cmpxchg i32 %0, i32 1, i32 2 acquire monotonic, align 4",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q:\n%s", want, got)
		}
	}
}

func TestBuilderAtomicPrecheck(t *testing.T) {
	ctx, m, b := buildAtomicModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{ctx.Ptr(0)}, false))
	entry := fn.NewBlock("entry")
	b.MoveToEnd(entry)
	p := fn.ParamAs[llvm.PtrT](0)
	one := ctx.ConstInt(i32, 1).Value

	// fence 不允许 monotonic/not_atomic
	if err := llvm.Catch(func() { b.Fence(llvm.AtomicMonotonic, false) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("fence monotonic should panic ErrInvalidArg, got %v", err)
	}
	// atomicrmw 不允许 unordered/not_atomic
	if err := llvm.Catch(func() {
		b.AtomicRMW(llvm.RMWAdd, p, one, llvm.AtomicNotAtomic, false, "")
	}); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("rmw not_atomic should panic ErrInvalidArg, got %v", err)
	}
	// cmpxchg 失败序不得为 release/acq_rel
	if err := llvm.Catch(func() {
		b.CmpXchg(p, one, one, llvm.AtomicAcquire, llvm.AtomicRelease, false, "")
	}); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("cmpxchg failure release should panic ErrInvalidArg, got %v", err)
	}
	// cmpxchg 失败序不得强于成功序
	if err := llvm.Catch(func() {
		b.CmpXchg(p, one, one, llvm.AtomicMonotonic, llvm.AtomicAcquire, false, "")
	}); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("cmpxchg failure stronger should panic ErrInvalidArg, got %v", err)
	}
	// cmp/new 类型必须一致
	if err := llvm.Catch(func() {
		b.CmpXchg(p, one, ctx.ConstFloat(ctx.Float(llvm.FloatDouble), 1).Value, llvm.AtomicAcquire, llvm.AtomicMonotonic, false, "")
	}); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("cmpxchg type mismatch should panic ErrTypeMismatch, got %v", err)
	}
	b.RetVoid()
}
```

- [x] **步骤 2：运行测试验证失败**

运行：`go test ./ir -run TestBuilderAtomic -v`
预期：编译失败，报错 `b.Fence undefined` / `llvm.AtomicSequentiallyConsistent undefined`。

- [x] **步骤 3：编写绑定、枚举与实现**

`internal/binding/Core.go` 枚举区（`LLVMDiagnosticSeverity` 之后）追加：

```go
// LLVMAtomicOrdering is the ordering of a fence/atomic load-store instruction.
type LLVMAtomicOrdering int32

const (
	LLVMAtomicOrderingNotAtomic              LLVMAtomicOrdering = C.LLVMAtomicOrderingNotAtomic
	LLVMAtomicOrderingUnordered              LLVMAtomicOrdering = C.LLVMAtomicOrderingUnordered
	LLVMAtomicOrderingMonotonic              LLVMAtomicOrdering = C.LLVMAtomicOrderingMonotonic
	LLVMAtomicOrderingAcquire                LLVMAtomicOrdering = C.LLVMAtomicOrderingAcquire
	LLVMAtomicOrderingRelease                LLVMAtomicOrdering = C.LLVMAtomicOrderingRelease
	LLVMAtomicOrderingAcquireRelease         LLVMAtomicOrdering = C.LLVMAtomicOrderingAcquireRelease
	LLVMAtomicOrderingSequentiallyConsistent LLVMAtomicOrdering = C.LLVMAtomicOrderingSequentiallyConsistent
)

// LLVMAtomicRMWBinOp is the operation of an atomicrmw instruction.
type LLVMAtomicRMWBinOp int32

const (
	LLVMAtomicRMWBinOpXchg     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpXchg
	LLVMAtomicRMWBinOpAdd      LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpAdd
	LLVMAtomicRMWBinOpSub      LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpSub
	LLVMAtomicRMWBinOpAnd      LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpAnd
	LLVMAtomicRMWBinOpNand     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpNand
	LLVMAtomicRMWBinOpOr       LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpOr
	LLVMAtomicRMWBinOpXor      LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpXor
	LLVMAtomicRMWBinOpMax      LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpMax
	LLVMAtomicRMWBinOpMin      LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpMin
	LLVMAtomicRMWBinOpUMax     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpUMax
	LLVMAtomicRMWBinOpUMin     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpUMin
	LLVMAtomicRMWBinOpFAdd     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpFAdd
	LLVMAtomicRMWBinOpFSub     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpFSub
	LLVMAtomicRMWBinOpFMax     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpFMax
	LLVMAtomicRMWBinOpFMin     LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpFMin
	LLVMAtomicRMWBinOpUIncWrap LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpUIncWrap
	LLVMAtomicRMWBinOpUDecWrap LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpUDecWrap
	LLVMAtomicRMWBinOpUSubCond LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpUSubCond
	LLVMAtomicRMWBinOpUSubSat  LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpUSubSat
	LLVMAtomicRMWBinOpFMaximum LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpFMaximum
	LLVMAtomicRMWBinOpFMinimum LLVMAtomicRMWBinOp = C.LLVMAtomicRMWBinOpFMinimum
)
```

`internal/binding/Core.go` 原子绑定（`LLVMBuildFCmp` 之后）：

```go
// LLVMGetVolatile Get the volatile flag of the memory access instruction.
func LLVMGetVolatile(inst LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMGetVolatile(inst.c))
}

// LLVMSetVolatile Set the volatile flag of the memory access instruction.
func LLVMSetVolatile(inst LLVMValueRef, isVolatile bool) {
	C.LLVMSetVolatile(inst.c, bool2LLVMBool(isVolatile))
}

// LLVMGetWeak Get the weak flag of the cmpxchg instruction.
func LLVMGetWeak(cmpXchgInst LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMGetWeak(cmpXchgInst.c))
}

// LLVMSetWeak Set the weak flag of the cmpxchg instruction.
func LLVMSetWeak(cmpXchgInst LLVMValueRef, isWeak bool) {
	C.LLVMSetWeak(cmpXchgInst.c, bool2LLVMBool(isWeak))
}

// LLVMGetOrdering Get the ordering of the memory access instruction.
func LLVMGetOrdering(inst LLVMValueRef) LLVMAtomicOrdering {
	return LLVMAtomicOrdering(C.LLVMGetOrdering(inst.c))
}

// LLVMSetOrdering Set the ordering of the memory access instruction.
func LLVMSetOrdering(inst LLVMValueRef, ordering LLVMAtomicOrdering) {
	C.LLVMSetOrdering(inst.c, C.LLVMAtomicOrdering(ordering))
}

// LLVMGetAtomicRMWBinOp Get the operation of the atomicrmw instruction.
func LLVMGetAtomicRMWBinOp(inst LLVMValueRef) LLVMAtomicRMWBinOp {
	return LLVMAtomicRMWBinOp(C.LLVMGetAtomicRMWBinOp(inst.c))
}

// LLVMSetAtomicRMWBinOp Set the operation of the atomicrmw instruction.
func LLVMSetAtomicRMWBinOp(inst LLVMValueRef, binOp LLVMAtomicRMWBinOp) {
	C.LLVMSetAtomicRMWBinOp(inst.c, C.LLVMAtomicRMWBinOp(binOp))
}

// LLVMGetCmpXchgSuccessOrdering Get the success ordering of the cmpxchg instruction.
func LLVMGetCmpXchgSuccessOrdering(cmpXchgInst LLVMValueRef) LLVMAtomicOrdering {
	return LLVMAtomicOrdering(C.LLVMGetCmpXchgSuccessOrdering(cmpXchgInst.c))
}

// LLVMSetCmpXchgSuccessOrdering Set the success ordering of the cmpxchg instruction.
func LLVMSetCmpXchgSuccessOrdering(cmpXchgInst LLVMValueRef, ordering LLVMAtomicOrdering) {
	C.LLVMSetCmpXchgSuccessOrdering(cmpXchgInst.c, C.LLVMAtomicOrdering(ordering))
}

// LLVMGetCmpXchgFailureOrdering Get the failure ordering of the cmpxchg instruction.
func LLVMGetCmpXchgFailureOrdering(cmpXchgInst LLVMValueRef) LLVMAtomicOrdering {
	return LLVMAtomicOrdering(C.LLVMGetCmpXchgFailureOrdering(cmpXchgInst.c))
}

// LLVMSetCmpXchgFailureOrdering Set the failure ordering of the cmpxchg instruction.
func LLVMSetCmpXchgFailureOrdering(cmpXchgInst LLVMValueRef, ordering LLVMAtomicOrdering) {
	C.LLVMSetCmpXchgFailureOrdering(cmpXchgInst.c, C.LLVMAtomicOrdering(ordering))
}

// LLVMIsAtomicSingleThread Get the singlethread flag of the atomic instruction.
func LLVMIsAtomicSingleThread(inst LLVMValueRef) bool {
	return llvmBool2bool(C.LLVMIsAtomicSingleThread(inst.c))
}

// LLVMSetAtomicSingleThread Set the singlethread flag of the atomic instruction.
func LLVMSetAtomicSingleThread(inst LLVMValueRef, singleThread bool) {
	C.LLVMSetAtomicSingleThread(inst.c, bool2LLVMBool(singleThread))
}

// LLVMBuildFence Create a fence instruction. Note: fence is void-valued and must not be named.
func LLVMBuildFence(builder LLVMBuilderRef, ordering LLVMAtomicOrdering, singleThread bool, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildFence(builder.c, C.LLVMAtomicOrdering(ordering), bool2LLVMBool(singleThread), name)}
	})
}

// LLVMBuildAtomicRMW Create an atomicrmw instruction. The C API takes no name;
// apply LLVMSetValueName afterwards if needed.
func LLVMBuildAtomicRMW(builder LLVMBuilderRef, op LLVMAtomicRMWBinOp, ptr, val LLVMValueRef, ordering LLVMAtomicOrdering, singleThread bool) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildAtomicRMW(builder.c, C.LLVMAtomicRMWBinOp(op), ptr.c, val.c, C.LLVMAtomicOrdering(ordering), bool2LLVMBool(singleThread))}
}

// LLVMBuildAtomicCmpXchg Create an atomic cmpxchg instruction. The C API takes no name;
// apply LLVMSetValueName afterwards if needed.
func LLVMBuildAtomicCmpXchg(builder LLVMBuilderRef, ptr, cmp, new LLVMValueRef, successOrdering, failureOrdering LLVMAtomicOrdering, singleThread bool) LLVMValueRef {
	return LLVMValueRef{c: C.LLVMBuildAtomicCmpXchg(builder.c, ptr.c, cmp.c, new.c, C.LLVMAtomicOrdering(successOrdering), C.LLVMAtomicOrdering(failureOrdering), bool2LLVMBool(singleThread))}
}
```

新建 `atomic.go`（root 包，公共枚举）：

```go
package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// AtomicOrdering 原子内存序
type AtomicOrdering binding.LLVMAtomicOrdering

const (
	AtomicNotAtomic              = AtomicOrdering(binding.LLVMAtomicOrderingNotAtomic)
	AtomicUnordered              = AtomicOrdering(binding.LLVMAtomicOrderingUnordered)
	AtomicMonotonic              = AtomicOrdering(binding.LLVMAtomicOrderingMonotonic)
	AtomicAcquire                = AtomicOrdering(binding.LLVMAtomicOrderingAcquire)
	AtomicRelease                = AtomicOrdering(binding.LLVMAtomicOrderingRelease)
	AtomicAcquireRelease         = AtomicOrdering(binding.LLVMAtomicOrderingAcquireRelease)
	AtomicSequentiallyConsistent = AtomicOrdering(binding.LLVMAtomicOrderingSequentiallyConsistent)
)

// RMWOp 原子读改写操作
type RMWOp binding.LLVMAtomicRMWBinOp

const (
	RMWXchg     = RMWOp(binding.LLVMAtomicRMWBinOpXchg)
	RMWAdd      = RMWOp(binding.LLVMAtomicRMWBinOpAdd)
	RMWSub      = RMWOp(binding.LLVMAtomicRMWBinOpSub)
	RMWAnd      = RMWOp(binding.LLVMAtomicRMWBinOpAnd)
	RMNand      = RMWOp(binding.LLVMAtomicRMWBinOpNand)
	RMWOr       = RMWOp(binding.LLVMAtomicRMWBinOpOr)
	RMWXor      = RMWOp(binding.LLVMAtomicRMWBinOpXor)
	RMWMax      = RMWOp(binding.LLVMAtomicRMWBinOpMax)
	RMWMin      = RMWOp(binding.LLVMAtomicRMWBinOpMin)
	RMWUMax     = RMWOp(binding.LLVMAtomicRMWBinOpUMax)
	RMWUMin     = RMWOp(binding.LLVMAtomicRMWBinOpUMin)
	RMWFAdd     = RMWOp(binding.LLVMAtomicRMWBinOpFAdd)
	RMWFSub     = RMWOp(binding.LLVMAtomicRMWBinOpFSub)
	RMWFMax     = RMWOp(binding.LLVMAtomicRMWBinOpFMax)
	RMWFMin     = RMWOp(binding.LLVMAtomicRMWBinOpFMin)
	RMWUIncWrap = RMWOp(binding.LLVMAtomicRMWBinOpUIncWrap)
	RMWUDecWrap = RMWOp(binding.LLVMAtomicRMWBinOpUDecWrap)
	RMWUSubCond = RMWOp(binding.LLVMAtomicRMWBinOpUSubCond)
	RMWUSubSat  = RMWOp(binding.LLVMAtomicRMWBinOpUSubSat)
	RMWFMaximum = RMWOp(binding.LLVMAtomicRMWBinOpFMaximum)
	RMWFMinimum = RMWOp(binding.LLVMAtomicRMWBinOpFMinimum)
)
```

新建 `ir/builder_atomic.go`：

```go
package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// ===== 内存序预检 =====

// preOrderingFence fence 合法内存序：acquire/release/acq_rel/seq_cst
func preOrderingFence(op string, o llvm.AtomicOrdering) {
	switch o {
	case llvm.AtomicAcquire, llvm.AtomicRelease, llvm.AtomicAcquireRelease, llvm.AtomicSequentiallyConsistent:
	default:
		llvm.Panicf(llvm.ErrInvalidArg, op, "invalid fence ordering %d", o)
	}
}

// preOrderingRMW atomicrmw/成功比较交换合法内存序：monotonic 起（not_atomic/unordered 非法）
func preOrderingRMW(op string, o llvm.AtomicOrdering) {
	switch o {
	case llvm.AtomicMonotonic, llvm.AtomicAcquire, llvm.AtomicRelease, llvm.AtomicAcquireRelease, llvm.AtomicSequentiallyConsistent:
	default:
		llvm.Panicf(llvm.ErrInvalidArg, op, "invalid atomic ordering %d", o)
	}
}

// orderingRank 内存序强度序（用于失败序不得强于成功序）
func orderingRank(o llvm.AtomicOrdering) int {
	switch o {
	case llvm.AtomicNotAtomic:
		return 0
	case llvm.AtomicUnordered:
		return 1
	case llvm.AtomicMonotonic:
		return 2
	case llvm.AtomicAcquire, llvm.AtomicRelease:
		return 3
	case llvm.AtomicAcquireRelease:
		return 4
	default:
		return 5 // seq_cst
	}
}

// preOrderingCmpXchg 比较交换内存序预检：失败序不得为 release/acq_rel 且不得强于成功序
func preOrderingCmpXchg(op string, success, failure llvm.AtomicOrdering) {
	preOrderingRMW(op, success)
	switch failure {
	case llvm.AtomicUnordered, llvm.AtomicMonotonic, llvm.AtomicAcquire, llvm.AtomicSequentiallyConsistent:
	default:
		llvm.Panicf(llvm.ErrInvalidArg, op, "invalid cmpxchg failure ordering %d", failure)
	}
	if orderingRank(failure) > orderingRank(success) {
		llvm.Panicf(llvm.ErrInvalidArg, op, "failure ordering %d is stronger than success ordering %d", failure, success)
	}
}

// ===== 指令角色 =====

// Fence fence 指令角色（内嵌 Value[VoidT]；void 值指令不可命名）
type Fence struct {
	llvm.Value[llvm.VoidT]
}

// Ordering 内存序
func (f Fence) Ordering() llvm.AtomicOrdering {
	f.Check("ir.Fence.Ordering")
	return llvm.AtomicOrdering(binding.LLVMGetOrdering(f.Ref()))
}

// SetOrdering 设置内存序
func (f Fence) SetOrdering(o llvm.AtomicOrdering) {
	const op = "ir.Fence.SetOrdering"
	f.Check(op)
	preOrderingFence(op, o)
	binding.LLVMSetOrdering(f.Ref(), binding.LLVMAtomicOrdering(o))
}

// IsSingleThread 是否 syncscope("singlethread")
func (f Fence) IsSingleThread() bool {
	f.Check("ir.Fence.IsSingleThread")
	return binding.LLVMIsAtomicSingleThread(f.Ref())
}

// SetSingleThread 设置 syncscope("singlethread")
func (f Fence) SetSingleThread(v bool) {
	f.Check("ir.Fence.SetSingleThread")
	binding.LLVMSetAtomicSingleThread(f.Ref(), v)
}

// AtomicRMW atomicrmw 指令角色（内嵌 Value[T]，T 为操作数种类，返回旧值）
type AtomicRMW[T llvm.Kind] struct {
	llvm.Value[T]
}

// Op 读改写操作
func (r AtomicRMW[T]) Op() llvm.RMWOp {
	r.Check("ir.AtomicRMW.Op")
	return llvm.RMWOp(binding.LLVMGetAtomicRMWBinOp(r.Ref()))
}

// SetOp 设置读改写操作
func (r AtomicRMW[T]) SetOp(op llvm.RMWOp) {
	r.Check("ir.AtomicRMW.SetOp")
	binding.LLVMSetAtomicRMWBinOp(r.Ref(), binding.LLVMAtomicRMWBinOp(op))
}

// Ordering 内存序
func (r AtomicRMW[T]) Ordering() llvm.AtomicOrdering {
	r.Check("ir.AtomicRMW.Ordering")
	return llvm.AtomicOrdering(binding.LLVMGetOrdering(r.Ref()))
}

// SetOrdering 设置内存序
func (r AtomicRMW[T]) SetOrdering(o llvm.AtomicOrdering) {
	const op = "ir.AtomicRMW.SetOrdering"
	r.Check(op)
	preOrderingRMW(op, o)
	binding.LLVMSetOrdering(r.Ref(), binding.LLVMAtomicOrdering(o))
}

// IsVolatile 是否 volatile
func (r AtomicRMW[T]) IsVolatile() bool {
	r.Check("ir.AtomicRMW.IsVolatile")
	return binding.LLVMGetVolatile(r.Ref())
}

// SetVolatile 设置 volatile
func (r AtomicRMW[T]) SetVolatile(v bool) {
	r.Check("ir.AtomicRMW.SetVolatile")
	binding.LLVMSetVolatile(r.Ref(), v)
}

// Align 对齐字节数
func (r AtomicRMW[T]) Align() uint32 {
	r.Check("ir.AtomicRMW.Align")
	return binding.LLVMGetAlignment(r.Ref())
}

// SetAlign 设置对齐字节数
func (r AtomicRMW[T]) SetAlign(n uint32) {
	const op = "ir.AtomicRMW.SetAlign"
	r.Check(op)
	preAlign(op, n)
	binding.LLVMSetAlignment(r.Ref(), n)
}

// CmpXchg cmpxchg 指令角色（内嵌 Value[StructT]，结果为 {旧值, i1 成功标志}）
type CmpXchg struct {
	llvm.Value[llvm.StructT]
}

// SuccessOrdering 成功序
func (x CmpXchg) SuccessOrdering() llvm.AtomicOrdering {
	x.Check("ir.CmpXchg.SuccessOrdering")
	return llvm.AtomicOrdering(binding.LLVMGetCmpXchgSuccessOrdering(x.Ref()))
}

// SetSuccessOrdering 设置成功序
func (x CmpXchg) SetSuccessOrdering(o llvm.AtomicOrdering) {
	const op = "ir.CmpXchg.SetSuccessOrdering"
	x.Check(op)
	preOrderingRMW(op, o)
	binding.LLVMSetCmpXchgSuccessOrdering(x.Ref(), binding.LLVMAtomicOrdering(o))
}

// FailureOrdering 失败序
func (x CmpXchg) FailureOrdering() llvm.AtomicOrdering {
	x.Check("ir.CmpXchg.FailureOrdering")
	return llvm.AtomicOrdering(binding.LLVMGetCmpXchgFailureOrdering(x.Ref()))
}

// SetFailureOrdering 设置失败序（不得强于成功序）
func (x CmpXchg) SetFailureOrdering(o llvm.AtomicOrdering) {
	const op = "ir.CmpXchg.SetFailureOrdering"
	x.Check(op)
	preOrderingCmpXchg(op, x.SuccessOrdering(), o)
	binding.LLVMSetCmpXchgFailureOrdering(x.Ref(), binding.LLVMAtomicOrdering(o))
}

// IsWeak 是否 weak
func (x CmpXchg) IsWeak() bool {
	x.Check("ir.CmpXchg.IsWeak")
	return binding.LLVMGetWeak(x.Ref())
}

// SetWeak 设置 weak
func (x CmpXchg) SetWeak(v bool) {
	x.Check("ir.CmpXchg.SetWeak")
	binding.LLVMSetWeak(x.Ref(), v)
}

// IsVolatile 是否 volatile
func (x CmpXchg) IsVolatile() bool {
	x.Check("ir.CmpXchg.IsVolatile")
	return binding.LLVMGetVolatile(x.Ref())
}

// SetVolatile 设置 volatile
func (x CmpXchg) SetVolatile(v bool) {
	x.Check("ir.CmpXchg.SetVolatile")
	binding.LLVMSetVolatile(x.Ref(), v)
}

// Align 对齐字节数
func (x CmpXchg) Align() uint32 {
	x.Check("ir.CmpXchg.Align")
	return binding.LLVMGetAlignment(x.Ref())
}

// SetAlign 设置对齐字节数
func (x CmpXchg) SetAlign(n uint32) {
	const op = "ir.CmpXchg.SetAlign"
	x.Check(op)
	preAlign(op, n)
	binding.LLVMSetAlignment(x.Ref(), n)
}

// ===== 构建方法 =====

// Fence 插入 fence（void 值指令不可命名，故无 name 参数）
func (b *Builder) Fence(order llvm.AtomicOrdering, singleThread bool) Fence {
	const op = "ir.Builder.Fence"
	b.pre(op)
	preOrderingFence(op, order)
	ref := binding.LLVMBuildFence(b.ref, binding.LLVMAtomicOrdering(order), singleThread, "")
	return Fence{Value: llvm.NewValue[llvm.VoidT](b.ctx, b.inserted.life, ref)}
}

// AtomicRMW 插入 atomicrmw，返回旧值（种类 T 与 val 一致）；
// C 构建 API 无 name，构建后经 SetValueName 命名
func (b *Builder) AtomicRMW[T llvm.Kind](op llvm.RMWOp, ptr llvm.ValueRef[llvm.PtrT], val llvm.ValueRef[T], order llvm.AtomicOrdering, singleThread bool, name string) AtomicRMW[T] {
	const opr = "ir.Builder.AtomicRMW"
	pv, vv := ptr.AsValue(), val.AsValue()
	b.pre(opr, core(pv), core(vv))
	preOrderingRMW(opr, order)
	ref := binding.LLVMBuildAtomicRMW(b.ref, binding.LLVMAtomicRMWBinOp(op), pv.Ref(), vv.Ref(), binding.LLVMAtomicOrdering(order), singleThread)
	binding.LLVMSetValueName(ref, name)
	return AtomicRMW[T]{Value: llvm.NewValue[T](b.ctx, b.inserted.life, ref)}
}

// CmpXchg 插入原子比较交换，返回 {旧值, i1 成功标志}；cmp/new 须同类型。
// C 构建 API 无 name，构建后经 SetValueName 命名
func (b *Builder) CmpXchg(ptr llvm.ValueRef[llvm.PtrT], cmp, new llvm.AnyValue, success, failure llvm.AtomicOrdering, singleThread bool, name string) CmpXchg {
	const op = "ir.Builder.CmpXchg"
	pv := ptr.AsValue()
	b.pre(op, core(pv), coreAny(cmp), coreAny(new))
	b.preSameType(op, coreAny(cmp), coreAny(new))
	preOrderingCmpXchg(op, success, failure)
	ref := binding.LLVMBuildAtomicCmpXchg(b.ref, pv.Ref(), cmp.Ref(), new.Ref(),
		binding.LLVMAtomicOrdering(success), binding.LLVMAtomicOrdering(failure), singleThread)
	binding.LLVMSetValueName(ref, name)
	return CmpXchg{Value: llvm.NewValue[llvm.StructT](b.ctx, b.inserted.life, ref)}
}
```

`ir/zz_refcheck_test.go` 追加：

```go
	_ llvm.ValueRef[llvm.VoidT]   = Fence{}
	_ llvm.ValueRef[llvm.IntT]    = AtomicRMW[llvm.IntT]{}
	_ llvm.ValueRef[llvm.StructT] = CmpXchg{}
	_ llvm.AnyValue               = Fence{}
	_ llvm.AnyValue               = AtomicRMW[llvm.IntT]{}
	_ llvm.AnyValue               = CmpXchg{}
```

- [x] **步骤 4：运行测试验证通过**

运行：`go test ./ir -run TestBuilderAtomic -v && go build ./... && go vet ./...`
预期：PASS。若 `fence acquire` 的打印含 `syncscope(...)` 或参数序与断言不同，以生成 IR 为准微调断言字符串（语义不变）。

- [x] **步骤 5：Commit**

```bash
git add internal/binding/Core.go atomic.go ir/builder_atomic.go ir/zz_refcheck_test.go ir/builder_atomic_test.go
git commit -m "feat(ir): fence/atomicrmw/cmpxchg 原子指令与内存序预检"
```

---

## 任务 7：Load/Store volatile 与原子内存序

**文件：**
- 修改：`ir/builder_mem.go`（`Load[T]`/`Store` 角色扩 volatile/ordering）
- 测试：`ir/builder_atomic_test.go` 追加用例

- [x] **步骤 1：编写失败的测试**

`ir/builder_atomic_test.go` 追加：

```go
func TestBuilderLoadStoreAtomic(t *testing.T) {
	ctx, m, b := buildAtomicModule(t)
	defer ctx.Close()
	defer m.Close()
	defer b.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), []llvm.AnyType{ctx.Ptr(0)}, false))
	entry := fn.NewBlock("entry")
	b.MoveToEnd(entry)
	p := fn.ParamAs[llvm.PtrT](0)

	ld := b.Load(p, i32, "v")
	ld.SetVolatile(true)
	ld.SetOrdering(llvm.AtomicAcquire)
	ld.SetAlign(4)
	if !ld.IsVolatile() || ld.Ordering() != llvm.AtomicAcquire || ld.Align() != 4 {
		t.Fatalf("load flags = %v %v %d", ld.IsVolatile(), ld.Ordering(), ld.Align())
	}

	st := b.Store(ctx.ConstInt(i32, 1).Value, p)
	st.SetVolatile(true)
	st.SetOrdering(llvm.AtomicRelease)
	if !st.IsVolatile() || st.Ordering() != llvm.AtomicRelease {
		t.Fatalf("store flags = %v %v", st.IsVolatile(), st.Ordering())
	}

	b.RetVoid()
	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := m.String()
	if !strings.Contains(got, "load atomic volatile i32, ptr %0 acquire, align 4") {
		t.Fatalf("missing atomic load:\n%s", got)
	}
	if !strings.Contains(got, "store atomic volatile i32 1, ptr %0 release") {
		t.Fatalf("missing atomic store:\n%s", got)
	}

	// load 不得设 release/acq_rel；store 不得设 acquire/acq_rel
	if err := llvm.Catch(func() { ld.SetOrdering(llvm.AtomicRelease) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("load release should panic ErrInvalidArg, got %v", err)
	}
	if err := llvm.Catch(func() { st.SetOrdering(llvm.AtomicAcquire) }); err == nil || err.Reason != llvm.ErrInvalidArg {
		t.Fatalf("store acquire should panic ErrInvalidArg, got %v", err)
	}
}
```

- [x] **步骤 2：运行测试验证失败**

运行：`go test ./ir -run TestBuilderLoadStoreAtomic -v`
预期：编译失败，报错 `ld.SetVolatile undefined`。

- [x] **步骤 3：编写最少实现代码**

`ir/builder_atomic.go` 预检区追加两个助手：

```go
// preOrderingLoad load 合法内存序：unordered/monotonic/acquire/seq_cst（及 not_atomic 复位）
func preOrderingLoad(op string, o llvm.AtomicOrdering) {
	switch o {
	case llvm.AtomicNotAtomic, llvm.AtomicUnordered, llvm.AtomicMonotonic, llvm.AtomicAcquire, llvm.AtomicSequentiallyConsistent:
	default:
		llvm.Panicf(llvm.ErrInvalidArg, op, "invalid load ordering %d", o)
	}
}

// preOrderingStore store 合法内存序：unordered/monotonic/release/seq_cst（及 not_atomic 复位）
func preOrderingStore(op string, o llvm.AtomicOrdering) {
	switch o {
	case llvm.AtomicNotAtomic, llvm.AtomicUnordered, llvm.AtomicMonotonic, llvm.AtomicRelease, llvm.AtomicSequentiallyConsistent:
	default:
		llvm.Panicf(llvm.ErrInvalidArg, op, "invalid store ordering %d", o)
	}
}
```

`ir/builder_mem.go` 的 `Load[T]` 角色（`Align` 之后）追加：

```go
// SetVolatile 设置 volatile 访问
func (l Load[T]) SetVolatile(v bool) {
	l.Check("ir.Load.SetVolatile")
	binding.LLVMSetVolatile(l.Ref(), v)
}

// IsVolatile 是否 volatile 访问
func (l Load[T]) IsVolatile() bool {
	l.Check("ir.Load.IsVolatile")
	return binding.LLVMGetVolatile(l.Ref())
}

// SetOrdering 设置原子内存序
func (l Load[T]) SetOrdering(o llvm.AtomicOrdering) {
	const op = "ir.Load.SetOrdering"
	l.Check(op)
	preOrderingLoad(op, o)
	binding.LLVMSetOrdering(l.Ref(), binding.LLVMAtomicOrdering(o))
}

// Ordering 原子内存序（非原子访问为 AtomicNotAtomic）
func (l Load[T]) Ordering() llvm.AtomicOrdering {
	l.Check("ir.Load.Ordering")
	return llvm.AtomicOrdering(binding.LLVMGetOrdering(l.Ref()))
}
```

`Store` 角色同构追加（操作名分别为 `ir.Store.SetVolatile`/`ir.Store.IsVolatile`/`ir.Store.SetOrdering`/`ir.Store.Ordering`，SetOrdering 用 `preOrderingStore`）：

```go
// SetVolatile 设置 volatile 访问
func (s Store) SetVolatile(v bool) {
	s.Check("ir.Store.SetVolatile")
	binding.LLVMSetVolatile(s.Ref(), v)
}

// IsVolatile 是否 volatile 访问
func (s Store) IsVolatile() bool {
	s.Check("ir.Store.IsVolatile")
	return binding.LLVMGetVolatile(s.Ref())
}

// SetOrdering 设置原子内存序
func (s Store) SetOrdering(o llvm.AtomicOrdering) {
	const op = "ir.Store.SetOrdering"
	s.Check(op)
	preOrderingStore(op, o)
	binding.LLVMSetOrdering(s.Ref(), binding.LLVMAtomicOrdering(o))
}

// Ordering 原子内存序（非原子访问为 AtomicNotAtomic）
func (s Store) Ordering() llvm.AtomicOrdering {
	s.Check("ir.Store.Ordering")
	return llvm.AtomicOrdering(binding.LLVMGetOrdering(s.Ref()))
}
```

- [x] **步骤 4：运行测试验证通过**

运行：`go test ./ir -v && go build ./... && go vet ./...`
预期：全绿。IR 打印若为 `load atomic volatile ...` 之外的词序，按实际输出微调断言。

- [x] **步骤 5：Commit**

```bash
git add ir/builder_mem.go ir/builder_atomic_test.go
git commit -m "feat(ir): Load/Store volatile 与原子内存序"
```

---

## 任务 8：vector 指令 + ConstVector

**文件：**
- 修改：`internal/binding/Core.go`（`LLVMBuildInsertElement`/`LLVMBuildShuffleVector`/`LLVMConstVector`）
- 新建：`ir/builder_vec.go`
- 修改：`constant.go`（`Context.ConstVector`）
- 测试：`ir/builder_vec_test.go`、`constant_test.go`

- [x] **步骤 1：编写失败的测试**

新建 `ir/builder_vec_test.go`：

```go
package ir

import (
	"strings"
	"testing"

	"github.com/kkkunny/go-llvm"
)

func TestBuilderVectorInsts(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "vec")
	defer m.Close()

	i32 := ctx.Int(32)
	vty := ctx.Vec(i32, 4)
	fn := m.NewFunction("f", ctx.Fn(i32, nil, false))
	entry := fn.NewBlock("entry")

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	zero := ctx.ConstZero(vty)
	one := ctx.ConstInt(i32, 1).Value
	idx := ctx.ConstInt(i32, 2).Value

	ins := b.InsertElement(zero, one, idx, "ins")
	ex := b.ExtractElement[llvm.IntT](ins, idx, "ex")
	sh := b.ShuffleVector(ins, ins, llvm.ValueOf(ctx, fn.Lifetime(), ctx.ConstVector(i32,
		ctx.ConstInt(i32, 3).Value, ctx.ConstInt(i32, 2).Value,
		ctx.ConstInt(i32, 1).Value, ctx.ConstInt(i32, 0).Value).Ref()), "sh")

	b.Ret(ex)
	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := m.String()
	for _, want := range []string{
		"%ins = insertelement <4 x i32> zeroinitializer, i32 1, i32 2",
		"%ex = extractelement <4 x i32> %ins, i32 2",
		"shufflevector <4 x i32> %ins, <4 x i32> %ins, <4 x i32> <i32 3, i32 2, i32 1, i32 0>",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q:\n%s", want, got)
		}
	}
	_ = sh
}

func TestBuilderVectorPrecheck(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "vec-pre")
	defer m.Close()

	i32 := ctx.Int(32)
	vty := ctx.Vec(i32, 2)
	fn := m.NewFunction("f", ctx.Fn(ctx.Void(), nil, false))
	entry := fn.NewBlock("entry")
	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	zero := ctx.ConstZero(vty)
	idx := ctx.ConstInt(i32, 0).Value

	// 元素类型不符
	if err := llvm.Catch(func() {
		b.InsertElement(zero, ctx.ConstFloat(ctx.Float(llvm.FloatDouble), 1).Value, idx, "")
	}); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("elem type mismatch should panic ErrTypeMismatch, got %v", err)
	}
	// ExtractElement 种类不符
	if err := llvm.Catch(func() {
		b.ExtractElement[llvm.FloatT](zero, idx, "")
	}); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("kind mismatch should panic ErrTypeMismatch, got %v", err)
	}
	// ShuffleVector v1/v2 类型不符
	if err := llvm.Catch(func() {
		b.ShuffleVector(zero, ctx.ConstZero(ctx.Vec(ctx.Int(64), 2)), zero, "")
	}); err == nil || err.Reason != llvm.ErrTypeMismatch {
		t.Fatalf("shuffle type mismatch should panic ErrTypeMismatch, got %v", err)
	}
	b.RetVoid()
}
```

`constant_test.go` 追加：

```go
func TestConstVector(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	i32 := ctx.Int(32)
	v := ctx.ConstVector(i32, ctx.ConstInt(i32, 1).Value, ctx.ConstInt(i32, 2).Value)
	if got := v.String(); !strings.Contains(got, "<i32 1, i32 2>") {
		t.Fatalf("const vector = %s", got)
	}

	// 元素类型不符
	if err := Catch(func() {
		ctx.ConstVector(i32, ctx.ConstFloat(ctx.Float(FloatDouble), 1).Value)
	}); err == nil || err.Reason != ErrTypeMismatch {
		t.Fatalf("elem type mismatch should panic ErrTypeMismatch, got %v", err)
	}
}
```

注：`TestBuilderVectorInsts` 中 shufflevector 的 mask 用 `llvm.ValueOf` 转换是多余的（`ConstVector` 返回 `Value[VecT]` 本身即 `ValueRef[VecT]`），最终测试代码直接传常量：

```go
	mask := ctx.ConstVector(i32,
		ctx.ConstInt(i32, 3).Value, ctx.ConstInt(i32, 2).Value,
		ctx.ConstInt(i32, 1).Value, ctx.ConstInt(i32, 0).Value)
	sh := b.ShuffleVector(ins, ins, mask, "sh")
```

- [x] **步骤 2：运行测试验证失败**

运行：`go test ./ir -run TestBuilderVector -v && go test . -run TestConstVector -v`
预期：编译失败，报错 `b.InsertElement undefined` / `ctx.ConstVector undefined`。

- [x] **步骤 3：编写绑定与实现**

`internal/binding/Core.go`（`LLVMBuildExtractValue` 之后）追加：

```go
// LLVMBuildInsertElement Insert a value into a vector element.
func LLVMBuildInsertElement(builder LLVMBuilderRef, vecVal, eltVal, index LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildInsertElement(builder.c, vecVal.c, eltVal.c, index.c, name)}
	})
}

// LLVMBuildShuffleVector Create a shufflevector instruction.
func LLVMBuildShuffleVector(builder LLVMBuilderRef, v1, v2, mask LLVMValueRef, name string) LLVMValueRef {
	return string2CString(name, func(name *C.char) LLVMValueRef {
		return LLVMValueRef{c: C.LLVMBuildShuffleVector(builder.c, v1.c, v2.c, mask.c, name)}
	})
}
```

`internal/binding/Core.go`（`LLVMConstNamedStruct` 附近）追加：

```go
// LLVMConstVector Create a ConstantVector from values.
func LLVMConstVector(scalarConstantVals []LLVMValueRef) LLVMValueRef {
	ptr, length := slice2Ptr[LLVMValueRef, C.LLVMValueRef](scalarConstantVals)
	return LLVMValueRef{c: C.LLVMConstVector(ptr, length)}
}
```

新建 `ir/builder_vec.go`：

```go
package ir

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// ExtractElement 提取向量第 idx 个元素；结果种类 U 与元素类型比对，不符 panic
func (b *Builder) ExtractElement[U llvm.Kind](vec llvm.ValueRef[llvm.VecT], idx llvm.ValueRef[llvm.IntT], name string) llvm.Value[U] {
	const op = "ir.Builder.ExtractElement"
	vv, iv := vec.AsValue(), idx.AsValue()
	b.pre(op, core(vv), core(iv))
	elemTy := binding.LLVMGetElementType(binding.LLVMTypeOf(vv.Ref()))
	checkKind[U](op, b.ctx, elemTy)
	ref := binding.LLVMBuildExtractElement(b.ref, vv.Ref(), iv.Ref(), name)
	return llvm.NewValue[U](b.ctx, b.inserted.life, ref)
}

// InsertElement 将元素插入向量第 idx 位；elem 须与向量元素类型一致
func (b *Builder) InsertElement(vec llvm.ValueRef[llvm.VecT], elem llvm.AnyValue, idx llvm.ValueRef[llvm.IntT], name string) llvm.Value[llvm.VecT] {
	const op = "ir.Builder.InsertElement"
	vv, iv := vec.AsValue(), idx.AsValue()
	b.pre(op, core(vv), coreAny(elem), core(iv))
	elemTy := binding.LLVMGetElementType(binding.LLVMTypeOf(vv.Ref()))
	if ref := binding.LLVMTypeOf(elem.Ref()); !elemTy.Equal(ref) {
		llvm.Panicf(llvm.ErrTypeMismatch, op, "element type %s does not match vector element type %s",
			typeString(b.ctx, elem.Ref()), typeRefString(b.ctx, elemTy))
	}
	ref := binding.LLVMBuildInsertElement(b.ref, vv.Ref(), elem.Ref(), iv.Ref(), name)
	return llvm.NewValue[llvm.VecT](b.ctx, b.inserted.life, ref)
}

// ShuffleVector 按 mask 重排两个同型向量；mask 为 <N x i32> 常量向量
func (b *Builder) ShuffleVector(v1, v2 llvm.ValueRef[llvm.VecT], mask llvm.ValueRef[llvm.VecT], name string) llvm.Value[llvm.VecT] {
	const op = "ir.Builder.ShuffleVector"
	a, c, mv := v1.AsValue(), v2.AsValue(), mask.AsValue()
	b.preSameType(op, core(a), core(c))
	b.checkVal(op, core(mv))
	maskElem := binding.LLVMGetElementType(binding.LLVMTypeOf(mv.Ref()))
	i32Ref := b.ctx.Int(32).Ref()
	if binding.LLVMGetTypeKind(maskElem) != binding.LLVMIntegerTypeKind || !maskElem.Equal(i32Ref) {
		llvm.Panicf(llvm.ErrTypeMismatch, op, "mask element must be i32, got %s", typeRefString(b.ctx, maskElem))
	}
	ref := binding.LLVMBuildShuffleVector(b.ref, a.Ref(), c.Ref(), mv.Ref(), name)
	return llvm.NewValue[llvm.VecT](b.ctx, b.inserted.life, ref)
}
```

`constant.go`（`ConstArray` 之后）追加：

```go
// ConstVector 构造向量常量；元素类型/归属不符则 panic
func (ctx *Context) ConstVector(elem AnyType, elems ...AnyValue) Value[VecT] {
	const op = "llvm.Context.ConstVector"
	ctx.CheckType(op, elem)
	ctx.CheckValues(op, elems...)
	elemRef := elem.Ref()
	for i, e := range elems {
		if ref := binding.LLVMTypeOf(e.Ref()); !elemRef.Equal(ref) {
			errPanic(ErrTypeMismatch, op,
				"element %d type %s does not match vector element type %s", i, TypeOfRef(ctx, ref), elem)
		}
	}
	ref := binding.LLVMConstVector(AnyValuesToRefs(elems))
	return Value[VecT]{ref: ref, ctx: ctx, life: ctx.life}
}
```

- [x] **步骤 4：运行测试验证通过**

运行：`go test ./ir -run TestBuilderVector -v && go test . -run TestConstVector -v && go build ./... && go vet ./...`
预期：PASS。IR 打印（`insertelement`/`shufflevector` 常量折行等）以实际输出为准微调断言。

- [x] **步骤 5：Commit**

```bash
git add internal/binding/Core.go ir/builder_vec.go constant.go ir/builder_vec_test.go constant_test.go
git commit -m "feat(ir): vector 指令（extractelement/insertelement/shufflevector）与 ConstVector"
```

---

## 任务 9：原子/向量 golden + 文档同步 + 全量回归

**文件：**
- 修改：`ir/golden_test.go`
- 新建：`ir/testdata/golden/atomic.ll`、`ir/testdata/golden/vec.ll`
- 修改：`AGENTS.md`（「Builder return types」行）
- 修改：`docs/superpowers/specs/2026-09-22-go-llvm-inkwell-style-redesign-design.md`（§9 as-built 表）

- [x] **步骤 1：编写失败的测试**

`ir/golden_test.go` 追加：

```go
// TestGoldenAtomic 覆盖 fence/atomicrmw/cmpxchg 与 load/store 原子形态
func TestGoldenAtomic(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "atomic")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("atomics", ctx.Fn(ctx.Void(), []llvm.AnyType{ctx.Ptr(0)}, false))
	entry := fn.NewBlock("entry")

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	p := fn.ParamAs[llvm.PtrT](0)
	one := ctx.ConstInt(i32, 1).Value

	b.Fence(llvm.AtomicSequentiallyConsistent, false)
	rm := b.AtomicRMW(llvm.RMWAdd, p, one, llvm.AtomicMonotonic, false, "old")
	rm.SetAlign(4)
	b.Store(one, p).SetOrdering(llvm.AtomicRelease)
	ld := b.Load(p, i32, "cur")
	ld.SetOrdering(llvm.AtomicAcquire)
	b.CmpXchg(p, ld, ctx.ConstInt(i32, 2).Value, llvm.AtomicAcquire, llvm.AtomicMonotonic, true, "cx")
	b.RetVoid()

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	checkGolden(t, m.String(), "atomic.ll")
}

// TestGoldenVec 覆盖 insertelement/extractelement/shufflevector 与 ConstVector
func TestGoldenVec(t *testing.T) {
	ctx := llvm.NewContext()
	defer ctx.Close()
	m := NewModule(ctx, "vec")
	defer m.Close()

	i32 := ctx.Int(32)
	fn := m.NewFunction("vec_demo", ctx.Fn(i32, nil, false))
	entry := fn.NewBlock("entry")

	b := NewBuilder(ctx)
	defer b.Close()
	b.MoveToEnd(entry)

	zero := ctx.ConstZero(ctx.Vec(i32, 4))
	ins := b.InsertElement(zero, ctx.ConstInt(i32, 42).Value, ctx.ConstInt(i32, 1).Value, "ins")
	mask := ctx.ConstVector(i32,
		ctx.ConstInt(i32, 0).Value, ctx.ConstInt(i32, 0).Value,
		ctx.ConstInt(i32, 0).Value, ctx.ConstInt(i32, 0).Value)
	sh := b.ShuffleVector(ins, ins, mask, "sh")
	ex := b.ExtractElement[llvm.IntT](sh, ctx.ConstInt(i32, 0).Value, "ex")
	b.Ret(ex)

	if err := m.Verify(); err != nil {
		t.Fatalf("verify: %v", err)
	}
	checkGolden(t, m.String(), "vec.ll")
}
```

- [x] **步骤 2：运行测试验证失败**

运行：`go test ./ir -run 'TestGoldenAtomic|TestGoldenVec' -v`
预期：FAIL，报错 `read golden (可用 -update 生成): open testdata/golden/atomic.ll: no such file or directory`。

- [x] **步骤 3：生成 golden 并核对**

运行：`go test ./ir -run 'TestGoldenAtomic|TestGoldenVec' -update`
核对关键语句存在（`atomic.ll`：`fence seq_cst`、`atomicrmw ... add i32 %0, i32 1 monotonic, align 4`、`store atomic i32 1, ptr %0 release`、`load atomic i32, ptr %0 acquire`、`cmpxchg weak ...`；`vec.ll`：`insertelement`、`shufflevector`、`extractelement`、`<i32 42, ...>`）。语句齐全即接受生成结果。

- [x] **步骤 4：全量回归**

运行：`go build ./... && go vet ./... && go test ./...`
预期：全部 PASS、build/vet 干净。

- [x] **步骤 5：文档同步**

`AGENTS.md` 「Core conventions」的「Builder return types」条目更新为：

```markdown
- **Builder return types**: return a role wrapper only when the instruction has role-specific operations (`Alloca`/`Load`/`Store`/`Call`/`Invoke`/`Phi`/`Switch`/`LandingPad`/`CatchSwitch`/`FuncletPad`/`Fence`/`AtomicRMW`/`CmpXchg`); everything else returns the plain `Value[T]`.
```

设计稿 §9 as-built 表追加一行：

```markdown
| §5 P2-1/P2-2/P2-4 指令面 | EH：`Invoke`/`LandingPad`/`CatchSwitch`/`FuncletPad` 角色 + `Invoke(Indirect)`/`LandingPad`/`Resume`/`CatchSwitch`/`CatchPad`/`CleanupPad`/`CatchRet`/`CleanupRet`；personality 走 `Function.SetPersonality`（`LLVMBuildLandingPad` 的 PersFn 参数传零值）。原子：`Fence`/`AtomicRMW[T]`/`CmpXchg` + Load/Store 的 volatile/ordering，内存序组合在预检层校验。vector：`ExtractElement`/`InsertElement`/`ShuffleVector` + `Context.ConstVector`。void 值指令（fence/resume/catchret/cleanupret）不提供 name 参数（LLVM 命名 void 值会被 verifier 拒绝） |
```

- [x] **步骤 6：Commit**

```bash
git add ir/golden_test.go ir/testdata/golden/atomic.ll ir/testdata/golden/vec.ll AGENTS.md docs/superpowers/specs/2026-09-22-go-llvm-inkwell-style-redesign-design.md
git commit -m "test(ir): 原子/vector golden 基线与文档同步"
git commit -m "test(ir): 原子/vector golden 基线与文档同步"
```

（上方 commit 行重复为笔误，只执行一次。）

---

## 验收对照（spec §5 P2-1/P2-2/P2-4）

- **P2-1 EH 全套**：invoke/landingpad/catchpad/catchret/cleanuppad/cleanupret/resume/catchswitch/catchswitch handler → 任务 1-5；golden `eh.ll`。
- **P2-2 原子与内存序**：fence/atomicrmw/cmpxchg、volatile、对齐属性 → 任务 6-7、9（`SetAlign` 见 AtomicRMW/CmpXchg/Load/Store）。
- **P2-4 vector 与聚合**：extractelement/insertelement/shufflevector（insertvalue/extractvalue 已有）→ 任务 8-9；`ConstVector` 补齐常量面。
- 每任务含预检测试（ErrInvalidArg/ErrTypeMismatch/ErrCrossContext）与 golden 或 IR 断言；`go build/vet/test` 每任务收尾必过。
