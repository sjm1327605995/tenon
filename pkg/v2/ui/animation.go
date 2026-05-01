package ui

import "time"

// AnimationStatus 表示动画控制器的当前状态。
type AnimationStatus int

const (
	AnimationDismissed AnimationStatus = iota
	AnimationForward
	AnimationReverse
	AnimationCompleted
)

// AnimationController 是动画时间线控制器。
type AnimationController struct {
	Duration   time.Duration
	LowerBound float64 // 默认 0
	UpperBound float64 // 默认 1
	Value      float64 // 当前值
	Status     AnimationStatus
	running    bool // 内部运行标志
}

// Forward 开始从当前值向 UpperBound 推进。
func (a *AnimationController) Forward() {
	a.Status = AnimationForward
	a.running = true
}

// Reverse 开始从当前值向 LowerBound 推进。
func (a *AnimationController) Reverse() {
	a.Status = AnimationReverse
	a.running = true
}

// Stop 停止动画推进。
func (a *AnimationController) Stop() {
	a.running = false
}

// Tick 根据 dt 推进动画。返回 true 表示值或状态有变化。
func (a *AnimationController) Tick(dt time.Duration) bool {
	if !a.running {
		return false
	}

	if a.Duration <= 0 {
		oldValue := a.Value
		if a.Status == AnimationForward {
			a.Value = a.UpperBound
			a.Status = AnimationCompleted
		} else if a.Status == AnimationReverse {
			a.Value = a.LowerBound
			a.Status = AnimationDismissed
		}
		a.running = false
		return a.Value != oldValue
	}

	oldValue := a.Value
	delta := (a.UpperBound - a.LowerBound) * float64(dt) / float64(a.Duration)

	if a.Status == AnimationForward {
		a.Value += delta
		if a.Value >= a.UpperBound {
			a.Value = a.UpperBound
			a.Status = AnimationCompleted
			a.running = false
			return true
		}
	} else if a.Status == AnimationReverse {
		a.Value -= delta
		if a.Value <= a.LowerBound {
			a.Value = a.LowerBound
			a.Status = AnimationDismissed
			a.running = false
			return true
		}
	} else {
		return false
	}

	return a.Value != oldValue
}

// Tween 是值区间插值器。
type Tween[T any] struct {
	Begin T
	End   T
}

// Evaluate 根据 progress 对 Begin/End 进行插值。
// 当前仅支持 float64 和 float32。
func (t *Tween[T]) Evaluate(progress float64) T {
	var zero T
	switch any(zero).(type) {
	case float64:
		b := any(t.Begin).(float64)
		e := any(t.End).(float64)
		return any(b + (e-b)*progress).(T)
	case float32:
		b := any(t.Begin).(float32)
		e := any(t.End).(float32)
		return any(b + (e-b)*float32(progress)).(T)
	default:
		return zero
	}
}

// Curve 定义缓动曲线。
type Curve interface {
	Transform(t float64) float64 // 输入 [0,1]，输出 [0,1]
}

// LinearCurve 是线性曲线。
type LinearCurve struct{}

func (LinearCurve) Transform(t float64) float64 {
	return t
}

// EaseInOutCurve 是 ease-in-out-cubic 曲线。
type EaseInOutCurve struct{}

func (EaseInOutCurve) Transform(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	f := -2*t + 2
	return 1 - f*f*f/2
}

// Animation 绑定 Controller + Tween + Curve。
type Animation[T any] struct {
	Controller *AnimationController
	Tween      *Tween[T]
	Curve      Curve
}

// Value 读取当前动画值（Controller.Value 经 Curve 变换后插值）。
func (a *Animation[T]) Value() T {
	progress := a.Curve.Transform(a.Controller.Value)
	return a.Tween.Evaluate(progress)
}
