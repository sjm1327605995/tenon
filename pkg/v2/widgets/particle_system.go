package widgets

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// ParticleSystemWidget 配置并显示粒子效果。
// 例子：
//   ParticleSystem().
//       W(200).H(200).
//       Emitter(render.PointEmitter{X:100,Y:100}).
//       Behavior(render.ParticleBehavior{MinSpeed:50,MaxSpeed:100,...}).
//       Burst(50).
//       Loop(true).
//       AutoStart(true)
type ParticleSystemWidget struct {
	ui.BaseWidget
	Width      float32
	Height     float32
	Emitter    render.ParticleEmitter
	Behavior   render.ParticleBehavior
	BurstCount int
	Loop       bool
	LoopDelay  float32 // 秒
	AutoStart  bool
}

// ParticleSystem 创建粒子系统 Widget。
func ParticleSystem() ParticleSystemWidget {
	return ParticleSystemWidget{
		Width:      200,
		Height:     200,
		Emitter:    render.PointEmitter{},
		Behavior:   DefaultParticleBehavior(),
		BurstCount: 30,
		Loop:       false,
		AutoStart:  true,
	}
}

// DefaultParticleBehavior 返回默认粒子行为（白色小圆点，向上飘散）。
func DefaultParticleBehavior() render.ParticleBehavior {
	return render.ParticleBehavior{
		MinSpeed: 20, MaxSpeed: 50,
		MinAngle: -0.5, MaxAngle: 0.5,
		MinSize: 2, MaxSize: 5,
		MinLife: 0.5, MaxLife: 1.5,
		AX: 0, AY: 0,
		MinRotSpeed: 0, MaxRotSpeed: 0,
		Colors: []*render.Color{
			render.NewColor(255, 255, 255, 255),
		},
		SizeOverLife: true,
	}
}

func (p ParticleSystemWidget) W(v float32) ParticleSystemWidget {
	p.Width = v
	return p
}

func (p ParticleSystemWidget) H(v float32) ParticleSystemWidget {
	p.Height = v
	return p
}

func (p ParticleSystemWidget) WithEmitter(e render.ParticleEmitter) ParticleSystemWidget {
	p.Emitter = e
	return p
}

func (p ParticleSystemWidget) WithBehavior(b render.ParticleBehavior) ParticleSystemWidget {
	p.Behavior = b
	return p
}

func (p ParticleSystemWidget) WithBurst(n int) ParticleSystemWidget {
	p.BurstCount = n
	return p
}

func (p ParticleSystemWidget) WithLoop(v bool) ParticleSystemWidget {
	p.Loop = v
	return p
}

func (p ParticleSystemWidget) WithLoopDelay(v float32) ParticleSystemWidget {
	p.LoopDelay = v
	return p
}

func (p ParticleSystemWidget) WithAutoStart(v bool) ParticleSystemWidget {
	p.AutoStart = v
	return p
}

func (p ParticleSystemWidget) CreateElement() ui.Element {
	e := &ParticleSystemElement{}
	e.RenderObjectElement.BaseElement.Init(e, p)
	return e
}

// ParticleSystemElement 管理 ParticleSystemRenderObject。
type ParticleSystemElement struct {
	ui.RenderObjectElement
}

func (e *ParticleSystemElement) CreateRenderObject() render.RenderObject {
	w := e.GetWidget().(ParticleSystemWidget)
	ro := render.NewParticleSystemRenderObject(w.Emitter, w.Behavior)
	ro.SetBounds(render.Bounds{Width: w.Width, Height: w.Height})
	ro.SetZIndex(100) // 默认在上层
	if w.AutoStart {
		ro.Start()
	}
	return ro
}

func (e *ParticleSystemElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(ParticleSystemWidget)
	old := oldWidget.(ParticleSystemWidget)
	ro := e.GetRenderObject().(*render.ParticleSystemRenderObject)
	if old.Width != w.Width || old.Height != w.Height {
		ro.SetBounds(render.Bounds{Width: w.Width, Height: w.Height})
	}
	if old.Emitter != w.Emitter {
		ro.SetEmitter(w.Emitter)
	}
	// Behavior 包含 slice，无法直接比较，简单重新设置
	ro.SetBehavior(w.Behavior)
	ro.SetLoop(w.Loop)
	ro.SetLoopDelay(w.LoopDelay)
	ro.SetBurstCount(w.BurstCount)
	if w.AutoStart && !old.AutoStart {
		ro.Start()
	}
}
