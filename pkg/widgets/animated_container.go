package widgets

import (
	"image/color"
	"time"

	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
)

// AnimatedContainer 是一个隐式动画容器。
// 当 Width/Height/BackgroundColor/BorderRadius 等属性变化时，自动启动过渡动画。
type AnimatedContainer struct {
	engine.BaseWidget
	Child           engine.Widget
	Duration        time.Duration
	Curve           engine.Curve
	Width           float32
	Height          float32
	BackgroundColor *render.Color
	BorderRadius    float32
	BorderColor     *render.Color
	BorderWidth     float32
	Padding         engine.EdgeInsets
	Margin          engine.EdgeInsets
}

// NewAnimatedContainer 创建隐式动画容器（必须后续通过 Builder/State 驱动属性变化）。
func NewAnimatedContainer() AnimatedContainer {
	return AnimatedContainer{Duration: 300 * time.Millisecond}
}

func (a AnimatedContainer) WithChild(child engine.Widget) AnimatedContainer {
	a.Child = child
	return a
}

func (a AnimatedContainer) WithDuration(d time.Duration) AnimatedContainer {
	a.Duration = d
	return a
}

func (a AnimatedContainer) WithCurve(c engine.Curve) AnimatedContainer {
	a.Curve = c
	return a
}

func (a AnimatedContainer) WithSize(w, h float32) AnimatedContainer {
	a.Width = w
	a.Height = h
	return a
}

func (a AnimatedContainer) WithBackground(c render.Color) AnimatedContainer {
	a.BackgroundColor = &c
	return a
}

func (a AnimatedContainer) WithRadius(r float32) AnimatedContainer {
	a.BorderRadius = r
	return a
}

func (a AnimatedContainer) WithBorder(c render.Color, w float32) AnimatedContainer {
	a.BorderColor = &c
	a.BorderWidth = w
	return a
}

func (a AnimatedContainer) WithPadding(p engine.EdgeInsets) AnimatedContainer {
	a.Padding = p
	return a
}

func (a AnimatedContainer) WithMargin(m engine.EdgeInsets) AnimatedContainer {
	a.Margin = m
	return a
}

func (a AnimatedContainer) CreateElement() engine.Element {
	return engine.NewStatefulElement(a)
}

func (a AnimatedContainer) CreateState() engine.State {
	return &animatedContainerState{}
}

// animatedContainerState 维护动画的当前值和控制器。
type animatedContainerState struct {
	engine.BaseStateOf[AnimatedContainer]

	ctrl *engine.AnimationController

	// 当前动画值
	width        float32
	height       float32
	bgColor      *render.Color
	borderRadius float32
	borderColor  *render.Color
	borderWidth  float32
	padding      engine.EdgeInsets
	margin       engine.EdgeInsets

	// 动画起始值
	startWidth        float32
	startHeight       float32
	startBgColor      *render.Color
	startBorderRadius float32
	startBorderColor  *render.Color
	startBorderWidth  float32
	startPadding      engine.EdgeInsets
	startMargin       engine.EdgeInsets

	// 动画目标值
	targetWidth        float32
	targetHeight       float32
	targetBgColor      *render.Color
	targetBorderRadius float32
	targetBorderColor  *render.Color
	targetBorderWidth  float32
	targetPadding      engine.EdgeInsets
	targetMargin       engine.EdgeInsets

	// 是否有属性正在动画中
	animating bool
}

func (s *animatedContainerState) InitState() {
	w := s.Widget()
	s.ctrl = &engine.AnimationController{Duration: w.Duration}
	if engine.DefaultEngine() != nil {
		engine.DefaultEngine().RegisterAnimation(s.ctrl)
	}
	s.ctrl.AddListener(s.onAnimationTick)

	// 初始当前值等于目标值
	s.width = w.Width
	s.height = w.Height
	s.bgColor = w.BackgroundColor
	s.borderRadius = w.BorderRadius
	s.borderColor = w.BorderColor
	s.borderWidth = w.BorderWidth
	s.padding = w.Padding
	s.margin = w.Margin
}

func (s *animatedContainerState) Dispose() {
	if engine.DefaultEngine() != nil {
		engine.DefaultEngine().UnregisterAnimation(s.ctrl)
	}
	s.ctrl.Stop()
	s.ctrl.RemoveListener(s.onAnimationTick)
}

func (s *animatedContainerState) DidUpdateWidget(old engine.Widget) {
	oldW := old.(AnimatedContainer)
	newW := s.Widget()

	changed := false

	if oldW.Width != newW.Width {
		s.startWidth = s.width
		s.targetWidth = newW.Width
		changed = true
	}
	if oldW.Height != newW.Height {
		s.startHeight = s.height
		s.targetHeight = newW.Height
		changed = true
	}
	if !render.ColorPtrEquals(oldW.BackgroundColor, newW.BackgroundColor) {
		s.startBgColor = cloneColor(s.bgColor)
		s.targetBgColor = newW.BackgroundColor
		changed = true
	}
	if oldW.BorderRadius != newW.BorderRadius {
		s.startBorderRadius = s.borderRadius
		s.targetBorderRadius = newW.BorderRadius
		changed = true
	}
	if !render.ColorPtrEquals(oldW.BorderColor, newW.BorderColor) {
		s.startBorderColor = cloneColor(s.borderColor)
		s.targetBorderColor = newW.BorderColor
		changed = true
	}
	if oldW.BorderWidth != newW.BorderWidth {
		s.startBorderWidth = s.borderWidth
		s.targetBorderWidth = newW.BorderWidth
		changed = true
	}
	if oldW.Padding != newW.Padding {
		s.startPadding = s.padding
		s.targetPadding = newW.Padding
		changed = true
	}
	if oldW.Margin != newW.Margin {
		s.startMargin = s.margin
		s.targetMargin = newW.Margin
		changed = true
	}

	if changed {
		s.ctrl.Value = 0
		s.ctrl.Status = engine.AnimationDismissed
		s.animating = true
		s.ctrl.Forward()
	}
}

func (s *animatedContainerState) onAnimationTick() {
	w := s.Widget()
	c := w.Curve
	if c == nil {
		c = engine.LinearCurve{}
	}
	progress := c.Transform(s.ctrl.Value)

	s.width = lerpFloat32(s.startWidth, s.targetWidth, progress)
	s.height = lerpFloat32(s.startHeight, s.targetHeight, progress)
	s.bgColor = lerpColor(s.startBgColor, s.targetBgColor, progress)
	s.borderRadius = lerpFloat32(s.startBorderRadius, s.targetBorderRadius, progress)
	s.borderColor = lerpColor(s.startBorderColor, s.targetBorderColor, progress)
	s.borderWidth = lerpFloat32(s.startBorderWidth, s.targetBorderWidth, progress)
	s.padding = lerpEdgeInsets(s.startPadding, s.targetPadding, progress)
	s.margin = lerpEdgeInsets(s.startMargin, s.targetMargin, progress)

	if s.ctrl.Status == engine.AnimationCompleted || s.ctrl.Status == engine.AnimationDismissed {
		s.animating = false
	}

	s.SetState(nil)
}

func (s *animatedContainerState) Build(ctx engine.BuildContext) engine.Widget {
	w := s.Widget()
	c := Container(w.Child)
	if s.width > 0 {
		c = c.W(s.width)
	}
	if s.height > 0 {
		c = c.H(s.height)
	}
	if s.bgColor != nil {
		c = c.Background(*s.bgColor)
	}
	if s.borderRadius > 0 {
		c = c.Radius(s.borderRadius)
	}
	if s.borderColor != nil && s.borderWidth > 0 {
		c = c.Border(*s.borderColor, s.borderWidth)
	}
	c = c.Pad(s.padding).Marginf(s.margin)
	return c
}

// ---- helpers ----

func lerpFloat32(a, b float32, t float64) float32 {
	return a + (b-a)*float32(t)
}

func lerpColor(a, b *render.Color, t float64) *render.Color {
	if a == nil && b == nil {
		return nil
	}
	if a == nil {
		ac := color.RGBA{R: 0, G: 0, B: 0, A: 0}
		a = render.NewColorFrom(ac)
	}
	if b == nil {
		bc := color.RGBA{R: 0, G: 0, B: 0, A: 0}
		b = render.NewColorFrom(bc)
	}
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	// RGBA() returns 0-65535, divide by 257 to get 0-255
	return render.NewColorFrom(color.RGBA{
		R: uint8(lerpFloat32(float32(ar)/257, float32(br)/257, t)),
		G: uint8(lerpFloat32(float32(ag)/257, float32(bg)/257, t)),
		B: uint8(lerpFloat32(float32(ab)/257, float32(bb)/257, t)),
		A: uint8(lerpFloat32(float32(aa)/257, float32(ba)/257, t)),
	})
}

func lerpEdgeInsets(a, b engine.EdgeInsets, t float64) engine.EdgeInsets {
	return engine.EdgeInsets{
		Top:    lerpFloat32(a.Top, b.Top, t),
		Right:  lerpFloat32(a.Right, b.Right, t),
		Bottom: lerpFloat32(a.Bottom, b.Bottom, t),
		Left:   lerpFloat32(a.Left, b.Left, t),
	}
}

func cloneColor(c *render.Color) *render.Color {
	if c == nil {
		return nil
	}
	cp := *c
	return &cp
}
