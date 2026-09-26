package llvm

import "testing"

// TestKindMarkers 覆盖各 Kind 标记的 kind() 方法，并验证 kindOf/sameKind 的判定。
func TestKindMarkers(t *testing.T) {
	kinds := []Kind{
		VoidT{}, IntT{}, FloatT{}, PtrT{}, StructT{}, ArrayT{},
		VecT{}, FnT{}, LabelT{}, MetaT{}, TokenT{}, DynT{},
	}
	for _, k := range kinds {
		k.kind() // 覆盖各标记的空实现（编译期约束的运行时体现）
	}
	if !sameKind(kindOf[IntT](), IntT{}) {
		t.Fatal("kindOf[IntT] 应等于 IntT{}")
	}
	if sameKind(kindOf[IntT](), FloatT{}) {
		t.Fatal("IntT 与 FloatT 不应相等")
	}
	if _, ok := kindOf[PtrT]().(PtrT); !ok {
		t.Fatal("kindOf[PtrT] 应返回 PtrT")
	}
}
