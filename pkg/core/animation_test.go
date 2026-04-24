package core

import (
	"testing"
	"time"
)

func TestEaseLinear(t *testing.T) {
	if EaseLinear(0) != 0 {
		t.Error("EaseLinear(0) should be 0")
	}
	if EaseLinear(1) != 1 {
		t.Error("EaseLinear(1) should be 1")
	}
	if EaseLinear(0.5) != 0.5 {
		t.Error("EaseLinear(0.5) should be 0.5")
	}
}

func TestEaseInOutQuad(t *testing.T) {
	if EaseInOutQuad(0) != 0 {
		t.Error("EaseInOutQuad(0) should be 0")
	}
	if EaseInOutQuad(1) != 1 {
		t.Error("EaseInOutQuad(1) should be 1")
	}
	mid := EaseInOutQuad(0.5)
	if mid < 0.4 || mid > 0.6 {
		t.Errorf("EaseInOutQuad(0.5) should be near 0.5, got %f", mid)
	}
}

func TestEaseOutQuad(t *testing.T) {
	if EaseOutQuad(0) != 0 {
		t.Error("EaseOutQuad(0) should be 0")
	}
	if EaseOutQuad(1) != 1 {
		t.Error("EaseOutQuad(1) should be 1")
	}
}

func TestEaseOutBounce(t *testing.T) {
	if EaseOutBounce(0) != 0 {
		t.Error("EaseOutBounce(0) should be 0")
	}
	if EaseOutBounce(1) != 1 {
		t.Error("EaseOutBounce(1) should be 1")
	}
}

func TestEaseOutElastic(t *testing.T) {
	if EaseOutElastic(0) != 0 {
		t.Error("EaseOutElastic(0) should be 0")
	}
	if EaseOutElastic(1) != 1 {
		t.Error("EaseOutElastic(1) should be 1")
	}
}

func TestLerpFloat32(t *testing.T) {
	if LerpFloat32(0, 100, 0.5) != 50 {
		t.Error("LerpFloat32(0, 100, 0.5) should be 50")
	}
	if LerpFloat32(10, 20, 0) != 10 {
		t.Error("LerpFloat32(10, 20, 0) should be 10")
	}
	if LerpFloat32(10, 20, 1) != 20 {
		t.Error("LerpFloat32(10, 20, 1) should be 20")
	}
}

func TestTweenBasic(t *testing.T) {
	var progress float32 = -1
	completed := false

	tw := NewTween(100*time.Millisecond, EaseLinear).
		OnUpdate(func(p float32) {
			progress = p
		}).
		OnComplete(func() {
			completed = true
		})

	tw.Start()
	if !tw.IsRunning() {
		t.Fatal("Tween should be running after Start")
	}

	// 第 1 帧：50ms
	stillRunning := tw.Update(0.05)
	if !stillRunning {
		t.Error("Tween should still be running at 50ms")
	}
	if progress < 0.4 || progress > 0.6 {
		t.Errorf("Progress should be near 0.5 at 50ms, got %f", progress)
	}

	// 第 2 帧：再 60ms（总计 110ms，超过 duration）
	stillRunning = tw.Update(0.06)
	if stillRunning {
		t.Error("Tween should have finished")
	}
	if !completed {
		t.Error("OnComplete should have been called")
	}
	if progress != 1 {
		t.Errorf("Final progress should be 1, got %f", progress)
	}
}

func TestTweenStop(t *testing.T) {
	tw := NewTween(1*time.Second, EaseLinear).OnUpdate(func(p float32) {})
	tw.Start()
	if !tw.IsRunning() {
		t.Fatal("Tween should be running")
	}
	tw.Stop()
	if tw.IsRunning() {
		t.Error("Tween should be stopped")
	}
}

func TestEngineAddAnimation(t *testing.T) {
	root := newTestHost("root")
	engine := NewEngine(root, 100, 100)

	tw := NewTween(100*time.Millisecond, EaseLinear).OnUpdate(func(p float32) {})
	tw.Start()
	engine.AddAnimation(tw)

	if len(engine.animations) != 1 {
		t.Fatalf("expected 1 animation, got %d", len(engine.animations))
	}

	// 更新多帧使动画完成（updateAnimations 会限制单帧最大 deltaTime）
	engine.updateAnimations(0.05)
	engine.updateAnimations(0.05)
	engine.updateAnimations(0.05)
	if len(engine.animations) != 0 {
		t.Fatalf("expected 0 animations after completion, got %d", len(engine.animations))
	}
}
