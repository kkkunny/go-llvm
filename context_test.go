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
