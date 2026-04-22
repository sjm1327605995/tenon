package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Switch 是开关组件。
type Switch struct {
	core.BaseHost
	checked       bool
	onChange      func(checked bool)
	trackWidth    float32
	trackHeight   float32
	offColor      color.Color
	onColor       color.Color
	thumbColor    color.Color
}

// NewSwitch 创建一个开关。
func NewSwitch() *Switch {
	s := &Switch{
		checked:     false,
		trackWidth:  44,
		trackHeight: 24,
		offColor:    color.RGBA{R: 204, G: 204, B: 204, A: 255},
		onColor:     color.RGBA{R: 0, G: 123, B: 255, A: 255},
		thumbColor:  color.White,
	}
	s.Init(s)
	s.SetFocusable(true)
	s.GetElement().Yoga.StyleSetWidth(s.trackWidth)
	s.GetElement().Yoga.StyleSetHeight(s.trackHeight)
	return s
}

// Draw 绘制开关轨道和按钮。
func (s *Switch) Draw(screen *ebiten.Image) {
	el := s.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := s.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	trackColor := s.offColor
	if s.checked {
		trackColor = s.onColor
	}

	// 绘制圆角轨道
	s.drawRoundedRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, bounds.Height/2, trackColor)

	// 绘制滑块
	thumbRadius := bounds.Height/2 - 2
	thumbY := bounds.Y + bounds.Height/2
	thumbX := bounds.X + thumbRadius + 2
	if s.checked {
		thumbX = bounds.X + bounds.Width - thumbRadius - 2
	}
	drawFilledCirclePath(screen, thumbX, thumbY, thumbRadius, s.thumbColor)
}

func (s *Switch) drawRoundedRect(screen *ebiten.Image, x, y, w, h, r float32, clr color.Color) {
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, core.BorderRadius{TopLeft: r, TopRight: r, BottomRight: r, BottomLeft: r})
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.FillPath(screen, &path, &vector.FillOptions{}, op)
}

// HandleEvent 处理点击事件。
func (s *Switch) HandleEvent(e *core.Event) bool {
	switch e.Type {
	case core.EventClick:
		s.checked = !s.checked
		if s.onChange != nil {
			s.onChange(s.checked)
		}
		return true
	}
	return false
}

// ==================== 链式 API ====================

func (s *Switch) SetChecked(checked bool) *Switch {
	s.checked = checked
	return s
}
func (s *Switch) SetOnChange(fn func(checked bool)) *Switch {
	s.onChange = fn
	return s
}
func (s *Switch) SetOffColor(clr color.Color) *Switch {
	s.offColor = clr
	return s
}
func (s *Switch) SetOnColor(clr color.Color) *Switch {
	s.onColor = clr
	return s
}
func (s *Switch) SetThumbColor(clr color.Color) *Switch {
	s.thumbColor = clr
	return s
}
func (s *Switch) SetTrackSize(width, height float32) *Switch {
	s.trackWidth = width
	s.trackHeight = height
	s.GetElement().Yoga.StyleSetWidth(width)
	s.GetElement().Yoga.StyleSetHeight(height)
	return s
}
func (s *Switch) SetMargin(edge yoga.Edge, value float32) *Switch {
	s.GetElement().Yoga.StyleSetMargin(edge, value)
	return s
}

// SyncFrom 同步开关属性。
func (s *Switch) SyncFrom(other core.Host) {
	if o, ok := other.(*Switch); ok {
		s.checked = o.checked
		s.onChange = o.onChange
		s.trackWidth = o.trackWidth
		s.trackHeight = o.trackHeight
		s.offColor = o.offColor
		s.onColor = o.onColor
		s.thumbColor = o.thumbColor
	}
}
