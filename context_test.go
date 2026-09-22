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
