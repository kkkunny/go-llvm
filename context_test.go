package llvm

import "testing"

type closeRecorder struct{ closed *bool }

func (c *closeRecorder) Close() error {
	*c.closed = true
	return nil
}

func TestContextCloseCascade(t *testing.T) {
	ctx := NewContext()
	closed := false
	ctx.Own(&closeRecorder{closed: &closed})
	if err := ctx.Close(); err != nil {
		t.Fatal(err)
	}
	if !closed {
		t.Fatal("Own 的子资源未被级联关闭")
	}
	if ctx.Alive() {
		t.Fatal("Close 后 Context 应不存活")
	}
}

func TestContextCloseCascadeReverseOrder(t *testing.T) {
	ctx := NewContext()
	var order []int
	ctx.Own(closeFunc(func() error { order = append(order, 1); return nil }))
	ctx.Own(closeFunc(func() error { order = append(order, 2); return nil }))
	_ = ctx.Close()
	if len(order) != 2 || order[0] != 2 || order[1] != 1 {
		t.Fatalf("期望逆序级联关闭，got %v", order)
	}
}

func TestContextDoubleClose(t *testing.T) {
	ctx := NewContext()
	_ = ctx.Close()
	err := ctx.Close()
	if err == nil || err.(*Error).Reason != ErrClosed {
		t.Fatalf("want ErrClosed, got %v", err)
	}
}

func TestContextOwnUnown(t *testing.T) {
	ctx := NewContext()
	closed := false
	unown := ctx.Own(&closeRecorder{closed: &closed})
	unown()
	unown()
	if err := ctx.Close(); err != nil {
		t.Fatal(err)
	}
	if closed {
		t.Fatal("注销后的子资源不应被级联关闭")
	}
}

func TestContextOwnUnownSelective(t *testing.T) {
	ctx := NewContext()
	var order []int
	ctx.Own(closeFunc(func() error { order = append(order, 1); return nil }))
	unown := ctx.Own(closeFunc(func() error { order = append(order, 2); return nil }))
	ctx.Own(closeFunc(func() error { order = append(order, 3); return nil }))
	unown()
	_ = ctx.Close()
	if len(order) != 2 || order[0] != 3 || order[1] != 1 {
		t.Fatalf("期望注销中间项后逆序关闭 [3,1]，got %v", order)
	}
}

func TestContextDisown(t *testing.T) {
	ctx := NewContext()
	closed := false
	ctx.Own(&closeRecorder{closed: &closed})
	life := ctx.Lifetime()
	if life == nil || life != ctx.Lifetime() || !life.Alive() {
		t.Fatal("Lifetime 应返回同一存活的令牌")
	}

	ctx.Disown()
	if !closed {
		t.Fatal("Disown 应级联关闭未移交的子资源")
	}
	if ctx.Alive() || life.Alive() {
		t.Fatal("Disown 后 Context 与令牌都应失效")
	}
	ctx.Disown() // 已移交后幂等
	if err := ctx.Close(); err == nil || err.(*Error).Reason != ErrClosed {
		t.Fatalf("Disown 后 Close 应为 ErrClosed, got %v", err)
	}

	// 已关闭的 Context 上 Disown 为 no-op
	closedCtx := NewContext()
	_ = closedCtx.Close()
	closedCtx.Disown()
}

func TestContextSyncScopeID(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	// 同一名称在同一 Context 内应映射到同一 ID（LLVM-C 合约）
	system := ctx.SyncScopeID("system")
	if got := ctx.SyncScopeID("system"); got != system {
		t.Fatalf("同名 scope 应返回同一 ID: %d vs %d", got, system)
	}
	custom := ctx.SyncScopeID("go-llvm-test-scope")
	if got := ctx.SyncScopeID("go-llvm-test-scope"); got != custom {
		t.Fatalf("同名 scope 应返回同一 ID: %d vs %d", got, custom)
	}
	if custom == system {
		t.Fatalf("不同 scope 名不应共用 ID: %d", custom)
	}
}

func TestTypeAliveCheck(t *testing.T) {
	ctx := NewContext()
	i32 := ctx.Int(32)
	if i32.String() != "i32" {
		t.Fatalf("String() = %q", i32.String())
	}
	_ = ctx.Close()

	err := Catch(func() { _ = i32.String() })
	if err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("want ErrUseAfterFree, got %v", err)
	}
	err = Catch(func() { _ = i32.Bits() })
	if err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("want ErrUseAfterFree from Bits, got %v", err)
	}
}

func TestContextConstructAfterClose(t *testing.T) {
	ctx := NewContext()
	_ = ctx.Close()
	err := Catch(func() { _ = ctx.Int(32) })
	if err == nil || err.Reason != ErrUseAfterFree {
		t.Fatalf("want ErrUseAfterFree, got %v", err)
	}
}

func TestLifetime(t *testing.T) {
	l := NewLifetime()
	if !l.Alive() {
		t.Fatal("新令牌应存活")
	}
	l.Kill()
	if l.Alive() {
		t.Fatal("Kill 后应不存活")
	}
}

type closeFunc func() error

func (f closeFunc) Close() error { return f() }
