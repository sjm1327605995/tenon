package engine

import "time"

// AnimationGroup 并行运行多个 AnimationController。
type AnimationGroup struct {
	controllers []*AnimationController
}

func NewAnimationGroup(controllers ...*AnimationController) *AnimationGroup {
	return &AnimationGroup{controllers: controllers}
}

func (g *AnimationGroup) Forward() {
	for _, c := range g.controllers {
		c.Forward()
	}
}

func (g *AnimationGroup) Reverse() {
	for _, c := range g.controllers {
		c.Reverse()
	}
}

func (g *AnimationGroup) Stop() {
	for _, c := range g.controllers {
		c.Stop()
	}
}

func (g *AnimationGroup) Add(c *AnimationController) {
	g.controllers = append(g.controllers, c)
}

// AnimationSequence 顺序运行多个 AnimationController。
type AnimationSequence struct {
	controllers []*AnimationController
	index       int
	listener    func()
}

func NewAnimationSequence(controllers ...*AnimationController) *AnimationSequence {
	return &AnimationSequence{controllers: controllers}
}

func (s *AnimationSequence) Forward() {
	s.index = 0
	if len(s.controllers) == 0 {
		return
	}
	s.controllers[0].Forward()
	s.listener = s.onTick
	s.controllers[0].AddListener(s.listener)
}

func (s *AnimationSequence) Reverse() {
	// 反向：从最后一个倒序执行
	s.index = len(s.controllers) - 1
	if s.index < 0 {
		return
	}
	s.controllers[s.index].Reverse()
	s.listener = s.onTickReverse
	s.controllers[s.index].AddListener(s.listener)
}

func (s *AnimationSequence) Stop() {
	for _, c := range s.controllers {
		c.Stop()
	}
	if s.index >= 0 && s.index < len(s.controllers) && s.listener != nil {
		s.controllers[s.index].RemoveListener(s.listener)
		s.listener = nil
	}
}

func (s *AnimationSequence) Add(c *AnimationController) {
	s.controllers = append(s.controllers, c)
}

func (s *AnimationSequence) onTick() {
	if s.index < 0 || s.index >= len(s.controllers) {
		return
	}
	c := s.controllers[s.index]
	if c.Status == AnimationCompleted {
		c.RemoveListener(s.listener)
		s.index++
		if s.index < len(s.controllers) {
			next := s.controllers[s.index]
			next.Forward()
			next.AddListener(s.listener)
		}
	}
}

func (s *AnimationSequence) onTickReverse() {
	if s.index < 0 || s.index >= len(s.controllers) {
		return
	}
	c := s.controllers[s.index]
	if c.Status == AnimationDismissed {
		c.RemoveListener(s.listener)
		s.index--
		if s.index >= 0 {
			prev := s.controllers[s.index]
			prev.Reverse()
			prev.AddListener(s.listener)
		}
	}
}

// CurvedAnimation 将 Curve 直接封装进 Controller 的监听器中，
// 提供便捷的 Value() 方法读取当前缓动后的值。
type CurvedAnimation struct {
	Controller *AnimationController
	Curve      Curve
}

func NewCurvedAnimation(ctrl *AnimationController, curve Curve) *CurvedAnimation {
	return &CurvedAnimation{Controller: ctrl, Curve: curve}
}

func (c *CurvedAnimation) Value() float64 {
	return c.Curve.Transform(c.Controller.Value)
}

// DelayedAnimation 包装一个 Controller，在延迟时间后才开始推进。
type DelayedAnimation struct {
	Controller *AnimationController
	Delay      time.Duration
	elapsed    time.Duration
	started    bool
}

func (d *DelayedAnimation) Tick(dt time.Duration) bool {
	if d.started {
		return d.Controller.Tick(dt)
	}
	d.elapsed += dt
	if d.elapsed >= d.Delay {
		d.started = true
		overflow := d.elapsed - d.Delay
		return d.Controller.Tick(overflow)
	}
	return false
}

func (d *DelayedAnimation) Forward() {
	d.Controller.Forward()
}

func (d *DelayedAnimation) Reverse() {
	d.Controller.Reverse()
}

func (d *DelayedAnimation) Stop() {
	d.Controller.Stop()
	d.started = false
	d.elapsed = 0
}
