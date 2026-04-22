package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Slider 是滑块组件。
type Slider struct {
	core.BaseHost
	minValue     float32
	maxValue     float32
	value        float32
	onChange     func(value float32)
	trackColor   color.Color
	fillColor    color.Color
	thumbColor   color.Color
	thumbRadius  float32
	trackHeight  float32
	dragging     bool
}

// NewSlider 创建一个滑块。
func NewSlider(min, max float32) *Slider {
	s := &Slider{
		minValue:    min,
		maxValue:    max,
		value:       min,
		trackColor:  color.RGBA{R: 224, G: 224, B: 224, A: 255},
		fillColor:   color.RGBA{R: 0, G: 123, B: 255, A: 255},
		thumbColor:  color.White,
		thumbRadius: 8,
		trackHeight: 4,
	}
	s.Init(s)
	s.SetFocusable(true)
	s.GetElement().Yoga.StyleSetHeight(24)
	return s
}

// Draw 绘制滑块轨道、填充和滑块按钮。
func (s *Slider) Draw(screen *ebiten.Image) {
	el := s.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := s.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	trackY := bounds.Y + bounds.Height/2
	trackTop := trackY - s.trackHeight/2

	// 轨道背景
	vector.FillRect(screen, bounds.X, trackTop, bounds.Width, s.trackHeight, s.trackColor, false)

	// 填充部分
	if s.maxValue > s.minValue {
		ratio := (s.value - s.minValue) / (s.maxValue - s.minValue)
		fillW := bounds.Width * ratio
		if fillW > 0 {
			vector.FillRect(screen, bounds.X, trackTop, fillW, s.trackHeight, s.fillColor, false)
		}
	}

	// 滑块按钮
	thumbX := bounds.X
	if s.maxValue > s.minValue {
		ratio := (s.value - s.minValue) / (s.maxValue - s.minValue)
		thumbX = bounds.X + bounds.Width*ratio
	}
	drawFilledCirclePath(screen, thumbX, trackY, s.thumbRadius, s.thumbColor)
	// 滑块边框
	strokeCirclePath(screen, thumbX, trackY, s.thumbRadius, 1.5, color.RGBA{R: 200, G: 200, B: 200, A: 255})
}

// HandleEvent 处理拖拽和点击事件。
func (s *Slider) HandleEvent(e *core.Event) bool {
	bounds := s.GetLayoutBounds()

	switch e.Type {
	case core.EventMouseDown:
		s.dragging = true
		s.updateValueFromX(e.X, bounds)
		return true
	case core.EventMouseUp:
		s.dragging = false
	case core.EventMouseMove:
		if s.dragging {
			s.updateValueFromX(e.X, bounds)
			return true
		}
	case core.EventClick:
		s.updateValueFromX(e.X, bounds)
		return true
	}
	return false
}

func (s *Slider) updateValueFromX(x float32, bounds core.LayoutBounds) {
	if bounds.Width <= 0 {
		return
	}
	ratio := (x - bounds.X) / bounds.Width
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	newValue := s.minValue + ratio*(s.maxValue-s.minValue)
	if newValue != s.value {
		s.value = newValue
		if s.onChange != nil {
			s.onChange(s.value)
		}
	}
}

// ==================== 链式 API ====================

func (s *Slider) SetValue(value float32) *Slider {
	if value < s.minValue {
		value = s.minValue
	}
	if value > s.maxValue {
		value = s.maxValue
	}
	s.value = value
	return s
}
func (s *Slider) SetOnChange(fn func(value float32)) *Slider {
	s.onChange = fn
	return s
}
func (s *Slider) SetTrackColor(clr color.Color) *Slider {
	s.trackColor = clr
	return s
}
func (s *Slider) SetFillColor(clr color.Color) *Slider {
	s.fillColor = clr
	return s
}
func (s *Slider) SetThumbColor(clr color.Color) *Slider {
	s.thumbColor = clr
	return s
}
func (s *Slider) SetThumbRadius(radius float32) *Slider {
	s.thumbRadius = radius
	return s
}
func (s *Slider) SetTrackHeight(height float32) *Slider {
	s.trackHeight = height
	return s
}
func (s *Slider) SetWidth(width float32) *Slider {
	s.GetElement().Yoga.StyleSetWidth(width)
	return s
}
func (s *Slider) SetMargin(edge yoga.Edge, value float32) *Slider {
	s.GetElement().Yoga.StyleSetMargin(edge, value)
	return s
}

// SyncFrom 同步滑块属性（保留 dragging 等交互状态）。
func (s *Slider) SyncFrom(other core.Host) {
	if o, ok := other.(*Slider); ok {
		s.minValue = o.minValue
		s.maxValue = o.maxValue
		s.value = o.value
		s.onChange = o.onChange
		s.trackColor = o.trackColor
		s.fillColor = o.fillColor
		s.thumbColor = o.thumbColor
		s.thumbRadius = o.thumbRadius
		s.trackHeight = o.trackHeight
	}
}
