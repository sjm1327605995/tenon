package core

// AnimationEngine 负责动画管理。
type AnimationEngine struct {
	engine *Engine
}

func newAnimationEngine(e *Engine) *AnimationEngine {
	return &AnimationEngine{engine: e}
}

func (a *AnimationEngine) AddAnimation(anim Animation) {
	a.engine.animations = append(a.engine.animations, anim)
}

func (a *AnimationEngine) updateAnimations(deltaTime float32) {
	const maxDelta = 1.0 / 30.0
	if deltaTime > maxDelta {
		deltaTime = maxDelta
	}
	alive := a.engine.animations[:0]
	for _, anim := range a.engine.animations {
		if anim.Update(deltaTime) {
			alive = append(alive, anim)
		}
	}
	a.engine.animations = alive
}
