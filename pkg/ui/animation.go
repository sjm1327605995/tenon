package ui

// Easing 是缓动函数，输入/输出均为 0..1。
type Easing func(float32) float32

func Linear(t float32) float32 { return t }

func EaseOut(t float32) float32 {
	u := 1 - t
	return 1 - u*u*u
}

func EaseIn(t float32) float32 { return t * t * t }

func EaseInOut(t float32) float32 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	u := -2*t + 2
	return 1 - u*u*u/2
}

type tweenHook struct {
	cur, from, target float32
	elapsed, duration float32
	ease              Easing
	fiber             *Fiber
	active            bool
}

// UseTween 把一个数值平滑过渡到 target；target 变化时自动从当前值开始动画，
// 动画期间引擎每帧重渲染该组件。durationMs 为毫秒，ease 可为 nil（默认 EaseOut）。
func UseTween(target, durationMs float32, ease Easing) float32 {
	f := currentFiber
	_, raw := nextHook(f, func() any {
		return &tweenHook{cur: target, target: target, from: target}
	})
	h := raw.(*tweenHook)
	h.fiber = f
	if ease == nil {
		ease = EaseOut
	}
	h.ease = ease
	if durationMs <= 0 {
		durationMs = 1
	}
	h.duration = durationMs

	if h.target != target {
		h.from = h.cur
		h.target = target
		h.elapsed = 0
		if !h.active {
			h.active = true
			if activeGame != nil {
				activeGame.anims = append(activeGame.anims, h)
			}
		}
	}
	return h.cur
}

type loopHook struct {
	elapsed float32 // 秒
	fiber   *Fiber
	active  bool
}

// UseElapsed 返回自挂载以来的秒数，并让引擎每帧重渲染该组件（用于持续动画，如加载指示器）。
func UseElapsed() float32 {
	f := currentFiber
	_, raw := nextHook(f, func() any { return &loopHook{fiber: f} })
	h := raw.(*loopHook)
	h.fiber = f
	if !h.active {
		h.active = true
		if activeGame != nil {
			activeGame.loops = append(activeGame.loops, h)
		}
	}
	return h.elapsed
}

type transitionHook struct {
	tween   tweenHook
	mounted bool
}

func boolf(b bool) float32 {
	if b {
		return 1
	}
	return 0
}

// UseTransition 管理进出场：返回 (mounted, progress)。
// visible 变 true 时立即 mounted 且 progress 0→1（进场）；
// 变 false 时 progress 1→0（退场），退场结束前保持 mounted，之后才卸载。
// 典型用法：If(mounted, Node(Opacity(progress), Scale(...)))。
func UseTransition(visible bool, durationMs float32) (bool, float32) {
	f := currentFiber
	_, raw := nextHook(f, func() any {
		v := boolf(visible)
		return &transitionHook{
			tween:   tweenHook{cur: v, from: v, target: v},
			mounted: visible,
		}
	})
	h := raw.(*transitionHook)
	tw := &h.tween
	tw.fiber = f
	tw.ease = EaseOut
	if durationMs <= 0 {
		durationMs = 1
	}
	tw.duration = durationMs

	if visible {
		h.mounted = true
	}
	target := boolf(visible)
	if tw.target != target {
		tw.from = tw.cur
		tw.target = target
		tw.elapsed = 0
		if !tw.active {
			tw.active = true
			if activeGame != nil {
				activeGame.anims = append(activeGame.anims, tw)
			}
		}
	}
	if !visible && !tw.active && tw.cur == 0 {
		h.mounted = false
	}
	return h.mounted, tw.cur
}

// advance 推进动画一帧，返回是否仍在进行。
func (h *tweenHook) advance(dt float32) bool {
	h.elapsed += dt
	t := h.elapsed / h.duration
	if t >= 1 {
		h.cur = h.target
		h.active = false
		return false
	}
	h.cur = h.from + (h.target-h.from)*h.ease(t)
	return true
}
