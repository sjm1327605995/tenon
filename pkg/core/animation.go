package core

import "time"

// Animation 是动画接口，所有动画类型都应实现此接口。
type Animation interface {
	// Update 使用 deltaTime（秒）更新动画状态，返回 true 表示动画仍在运行。
	Update(deltaTime float32) bool
	// Stop 立即停止动画。
	Stop()
}

// Tween 是一个补间动画，在指定时长内从 0 过渡到 1。
type Tween struct {
	duration   time.Duration
	elapsed    time.Duration
	easing     EasingFunction
	onUpdate   func(progress float32)
	onComplete func()
	running    bool
}

// NewTween 创建一个新的 Tween 动画。
func NewTween(duration time.Duration, easing EasingFunction) *Tween {
	if easing == nil {
		easing = EaseLinear
	}
	return &Tween{
		duration: duration,
		easing:   easing,
	}
}

// OnUpdate 设置每帧更新回调，参数 progress 范围 [0, 1]。
func (t *Tween) OnUpdate(fn func(progress float32)) *Tween {
	t.onUpdate = fn
	return t
}

// OnComplete 设置动画完成回调。
func (t *Tween) OnComplete(fn func()) *Tween {
	t.onComplete = fn
	return t
}

// Start 启动动画。
func (t *Tween) Start() {
	t.elapsed = 0
	t.running = true
}

// Update 更新动画状态。
func (t *Tween) Update(deltaTime float32) bool {
	if !t.running {
		return false
	}
	t.elapsed += time.Duration(deltaTime * float32(time.Second))
	if t.elapsed >= t.duration {
		t.elapsed = t.duration
		t.running = false
		progress := float32(1)
		if t.onUpdate != nil {
			t.onUpdate(t.easing(progress))
		}
		if t.onComplete != nil {
			t.onComplete()
		}
		return false
	}
	progress := float32(t.elapsed) / float32(t.duration)
	if t.onUpdate != nil {
		t.onUpdate(t.easing(progress))
	}
	return true
}

// Stop 停止动画。
func (t *Tween) Stop() {
	t.running = false
}

// IsRunning 返回动画是否正在运行。
func (t *Tween) IsRunning() bool {
	return t.running
}
