package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

// Slider is a draggable value slider.
type Slider struct {
	core.BaseElement
	minValue    float32
	maxValue    float32
	value       float32
	onChange    func(value float32)
	trackColor  color.Color
	fillColor   color.Color
	thumbColor  color.Color
	thumbRadius float32
	trackHeight float32
	dragging    bool
}

// NewSlider creates a slider.
func NewSlider(min, max float32) *Slider {
	theme := core.GetTheme()
	s := &Slider{
		minValue:    min,
		maxValue:    max,
		value:       min,
		trackColor:  theme.SliderTrackColor,
		fillColor:   theme.SliderFillColor,
		thumbColor:  theme.SliderThumbColor,
		thumbRadius: 8,
		trackHeight: 4,
	}
	s.Init(s)
	s.SetFlag(core.FlagFocusable)
	s.SetHeight(24)
	return s
}

// ElementType returns type identifier.
func (s *Slider) ElementType() string { return "Slider" }

// Draw renders the track, fill and thumb.
func (s *Slider) Draw(screen *ebiten.Image) {
	if !s.IsVisible() {
		return
	}
	bounds := s.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	trackY := bounds.Y + bounds.Height/2
	trackTop := trackY - s.trackHeight/2

	// Track background
	vector.FillRect(screen, bounds.X, trackTop, bounds.Width, s.trackHeight, s.trackColor, false)

	// Fill
	if s.maxValue > s.minValue {
		ratio := (s.value - s.minValue) / (s.maxValue - s.minValue)
		fillW := bounds.Width * ratio
		if fillW > 0 {
			vector.FillRect(screen, bounds.X, trackTop, fillW, s.trackHeight, s.fillColor, false)
		}
	}

	// Thumb
	thumbX := bounds.X
	if s.maxValue > s.minValue {
		ratio := (s.value - s.minValue) / (s.maxValue - s.minValue)
		thumbX = bounds.X + bounds.Width*ratio
	}
	drawFilledCirclePath(screen, thumbX, trackY, s.thumbRadius, s.thumbColor)
	strokeCirclePath(screen, thumbX, trackY, s.thumbRadius, 1.5, core.GetTheme().BorderColor)
}

// HandleEvent processes drag and click events.
func (s *Slider) HandleEvent(e *core.Event) bool {
	bounds := s.GetBounds()

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
		s.Mark(core.FlagNeedDraw)
		if s.onChange != nil {
			s.onChange(s.value)
		}
	}
}

// SetValue sets the slider value.
func (s *Slider) SetValue(value float32) *Slider {
	if value < s.minValue {
		value = s.minValue
	}
	if value > s.maxValue {
		value = s.maxValue
	}
	s.value = value
	s.Mark(core.FlagNeedDraw)
	return s
}

// SetOnChange sets the change callback.
func (s *Slider) SetOnChange(fn func(value float32)) *Slider {
	s.onChange = fn
	return s
}

// SetTrackColor sets the track background color.
func (s *Slider) SetTrackColor(clr color.Color) *Slider {
	s.trackColor = clr
	s.Mark(core.FlagNeedDraw)
	return s
}

// SetFillColor sets the fill color.
func (s *Slider) SetFillColor(clr color.Color) *Slider {
	s.fillColor = clr
	s.Mark(core.FlagNeedDraw)
	return s
}

// SetThumbColor sets the thumb color.
func (s *Slider) SetThumbColor(clr color.Color) *Slider {
	s.thumbColor = clr
	s.Mark(core.FlagNeedDraw)
	return s
}

// SetThumbRadius sets the thumb radius.
func (s *Slider) SetThumbRadius(radius float32) *Slider {
	s.thumbRadius = radius
	s.Mark(core.FlagNeedDraw)
	return s
}

// SetTrackHeight sets the track height.
func (s *Slider) SetTrackHeight(height float32) *Slider {
	s.trackHeight = height
	s.Mark(core.FlagNeedDraw)
	return s
}
