package llvm

import (
	"strings"
	"testing"
)

func TestMetadataMDString(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	s := ctx.MDString("hello")
	if s.IsNil() || s.Context() != ctx {
		t.Fatalf("MDString 句柄异常: nil=%v ctx=%p", s.IsNil(), s.Context())
	}
	cast := MetadataOf(ctx, s.Ref())
	if cast.Ref() != s.Ref() || cast.Context() != ctx {
		t.Fatal("MetadataOf 应保留底层句柄与所属 Context")
	}
	cast.Check("llvm.Test.Metadata.Check")

	if !s.IsString() || s.IsNode() || s.IsValueAsMetadata() {
		t.Fatalf("MDString 判定错误: string=%v node=%v value=%v",
			s.IsString(), s.IsNode(), s.IsValueAsMetadata())
	}
	if got := s.StringValue(); got != "hello" {
		t.Fatalf("StringValue() = %q", got)
	}
	if got := s.String(); got != `!"hello"` {
		t.Fatalf("MDString String() = %q", got)
	}
	// Value 形式：metadata 类型
	if got := s.Value().Type().String(); got != "metadata" {
		t.Fatalf("Metadata.Value().Type() = %q", got)
	}
}

func TestMetadataMDNode(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	n := ctx.MDNode(ctx.MDString("hello"), ctx.MDString("world"))
	if !n.IsNode() || n.IsString() || n.IsValueAsMetadata() {
		t.Fatalf("MDNode 判定错误: string=%v node=%v value=%v",
			n.IsString(), n.IsNode(), n.IsValueAsMetadata())
	}
	ops := n.Operands()
	if len(ops) != 2 {
		t.Fatalf("Operands() 长度 = %d, want 2", len(ops))
	}
	if got := ops[0].StringValue(); got != "hello" {
		t.Fatalf("ops[0] = %q", got)
	}
	if got := ops[1].StringValue(); got != "world" {
		t.Fatalf("ops[1] = %q", got)
	}
	if got := n.String(); !strings.Contains(got, `!{!"hello", !"world"}`) {
		t.Fatalf("MDNode String() = %q", got)
	}

	// 空节点返回 nil 操作数
	if got := ctx.MDNode().Operands(); got != nil {
		t.Fatalf("空 MDNode Operands() = %v, want nil", got)
	}
}

func TestMetadataValueAsMetadata(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	vm := ctx.ValueAsMetadata(ctx.ConstSInt(ctx.Int(32), 7))
	if !vm.IsValueAsMetadata() || vm.IsString() {
		t.Fatalf("ValueAsMetadata 判定错误: string=%v value=%v",
			vm.IsString(), vm.IsValueAsMetadata())
	}
	// LLVM-C 的 LLVMIsAMDNode 对 ValueAsMetadata 同样返回非 nil：单常量 MDNode 会被
	// 规范化为 ConstantAsMetadata，LLVM-C 再反向重建为操作数即被包装值的节点
	// （上游注明 LLVMIsAMDNode “a bit of a lier”，真节点判定需再排除 ValueAsMetadata）。
	if !vm.IsNode() {
		t.Fatal("LLVM 22 的 LLVMIsAMDNode 应将 ValueAsMetadata 视为节点")
	}
	if ops := vm.Operands(); len(ops) != 1 || ops[0].String() != "i32 7" {
		t.Fatalf("ValueAsMetadata 操作数 = %v", ops)
	}
	if got := vm.String(); got != "i32 7" {
		t.Fatalf("ValueAsMetadata String() = %q", got)
	}
}

// TestMetadataCrashFloor 崩溃类地板负向：nil 句柄 / 已释放 Context / 跨 Context 元素。
func TestMetadataCrashFloor(t *testing.T) {
	if err := Catch(func() { Metadata{}.IsString() }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("nil 元数据应 panic ErrInvalidArg, got %v", err)
	}

	ctx := NewContext()
	ref := ctx.MDString("x").Ref()
	if err := Catch(func() { (Metadata{ref: ref}).IsString() }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("ctx 为 nil 的元数据应 panic ErrInvalidArg, got %v", err)
	}

	m := ctx.MDString("x")
	_ = ctx.Close()
	if err := Catch(func() { _ = m.String() }); err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("已释放 Context 上的元数据应 panic ErrUseAfterFree, got %v", err)
	}

	ctx2 := NewContext()
	defer ctx2.Close()
	if err := Catch(func() { ctx2.MDNode(ctx2.MDString("a"), Metadata{}) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("nil 元素应 panic ErrInvalidArg, got %v", err)
	}
	ctx3 := NewContext()
	defer ctx3.Close()
	if err := Catch(func() { ctx2.MDNode(ctx3.MDString("other")) }); err == nil || err.Reason != ErrCrossContext {
		t.Fatalf("跨 Context 元素应 panic ErrCrossContext, got %v", err)
	}
	if err := Catch(func() { ctx2.ValueAsMetadata(nil) }); err == nil || err.Reason != ErrInvalidArg {
		t.Fatalf("nil 值包装应 panic ErrInvalidArg, got %v", err)
	}
}

// TestMetadataWrongKind 语义契约负向：访问器用于错误种类的元数据。
func TestMetadataWrongKind(t *testing.T) {
	requireDebug(t)

	ctx := NewContext()
	defer ctx.Close()
	s := ctx.MDString("hello")
	n := ctx.MDNode(s)

	if err := Catch(func() { n.StringValue() }); err == nil || err.Reason != ErrTypeMismatch {
		t.Fatalf("MDNode 取字符串应 panic ErrTypeMismatch, got %v", err)
	}
	if err := Catch(func() { s.Operands() }); err == nil || err.Reason != ErrTypeMismatch {
		t.Fatalf("MDString 取操作数应 panic ErrTypeMismatch, got %v", err)
	}
}

func TestModuleFlagBehaviorValues(t *testing.T) {
	// 与 llvm-c/Core.h 的 LLVMModuleFlagBehavior 枚举一一对应
	want := []struct {
		name string
		got  ModuleFlagBehavior
		val  int
	}{
		{"Error", ModuleFlagError, 0},
		{"Warning", ModuleFlagWarning, 1},
		{"Require", ModuleFlagRequire, 2},
		{"Override", ModuleFlagOverride, 3},
		{"Append", ModuleFlagAppend, 4},
		{"AppendUnique", ModuleFlagAppendUnique, 5},
	}
	for _, c := range want {
		if int(c.got) != c.val {
			t.Errorf("ModuleFlag%s = %d, want %d", c.name, int(c.got), c.val)
		}
	}
}
