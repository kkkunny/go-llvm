package jit

import (
	"reflect"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
	"github.com/kkkunny/go-llvm/internal/checks"
	"github.com/kkkunny/go-llvm/ir"
)

// ===== 调试层 JIT 签名核对（B3）=====
//
// Func[F]/MapFunc[F] 按 Go 签名构造调用桥；签名与 JIT 符号的真实 LLVM 类型不匹配
// 属于未定义行为（对标 inkwell 在 JIT 调用处标 unsafe 的契约）。调试构建在
// AddIRModule 时记录模块内函数签名（name → 规范化类型文本），在 Func/MapFunc
// 注册时核对；release 构建不记录、不核对，零成本。

// recordModuleSigs 记录模块内函数签名；必须在模块被 Disown 前调用。
func (j *LLJIT) recordModuleSigs(m *ir.Module) {
	j.sigMu.Lock()
	defer j.sigMu.Unlock()
	for ref := binding.LLVMGetFirstFunction(m.Ref()); !ref.IsNil(); ref = binding.LLVMGetNextFunction(ref) {
		name := binding.LLVMGetValueName(ref)
		if name == "" {
			continue
		}
		// 注意：LLVM 22 不透明指针下 LLVMTypeOf(函数值) 是 ptr，须用 LLVMGetFunctionType
		ty := binding.LLVMGetFunctionType(ref)
		if ty.IsNil() {
			continue
		}
		if j.symbolSigs == nil {
			j.symbolSigs = make(map[string]string)
		}
		j.symbolSigs[name] = binding.LLVMPrintTypeToString(ty)
	}
}

// checkSymbolSig 核对 Go 签名与已知 JIT 符号类型；不匹配 panic ErrTypeMismatch。
// 未知符号（MapSymbol/对象文件引入）跳过；不可过桥的签名留给 adapterFor 报错。
func (j *LLJIT) checkSymbolSig(op, name string, ft reflect.Type) {
	if !checks.Debug {
		return
	}
	j.sigMu.Lock()
	want, ok := j.symbolSigs[name]
	j.sigMu.Unlock()
	if !ok {
		return
	}
	got := j.goSigString(ft)
	if got == "" {
		return
	}
	if got != want {
		llvm.Panicf(llvm.ErrTypeMismatch, op,
			"Go signature %s maps to LLVM type %q, but JIT symbol %q has type %q; calling it would be undefined behavior",
			ft, got, name, want)
	}
}

// goSigString Go 函数签名对应的 LLVM 类型文本（按 reflect.Type 缓存，仅调试构建调用）
func (j *LLJIT) goSigString(ft reflect.Type) string {
	j.sigMu.Lock()
	if s, ok := j.goSigs[ft]; ok {
		j.sigMu.Unlock()
		return s
	}
	j.sigMu.Unlock()

	// 用临时 Context 做 Go→LLVM 签名映射；不支持和类型直接返回空串（交由既有错误路径）
	ctx := llvm.NewContext()
	defer ctx.Close()
	sig, _, err := llvm.FnSignatureOfGo(ctx, ft)
	if err != nil {
		return ""
	}
	s := sig.String()

	j.sigMu.Lock()
	if j.goSigs == nil {
		j.goSigs = make(map[reflect.Type]string)
	}
	j.goSigs[ft] = s
	j.sigMu.Unlock()
	return s
}
