package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
)

// Slider creates a value slider.
func Slider(min, max, value float32, onChange func(float32)) engine.Widget {
	return sliderWidget{MinValue: min, MaxValue: max, Value: value, OnChange: onChange}
}

// internal widget type
type sliderWidget struct {
	engine.BaseWidget
	MinValue float32
	MaxValue float32
	Value    float32
	OnChange func(value float32)
}

func (s sliderWidget) CreateElement() engine.Element {
	return engine.NewRenderObjectElement(s)
}

// CreateRenderObject implements RenderObjectFactory.
func (s sliderWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	theme := engine.GetTheme()
	r := render.NewRenderSlider()
	r.MinValue = s.MinValue
	r.MaxValue = s.MaxValue
	r.SetValue(s.Value)
	r.TrackColor = theme.SliderTrackColor
	r.FillColor = theme.SliderFillColor
	r.ThumbColor = theme.SliderThumbColor
	r.OnChange = s.OnChange
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (s sliderWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	old := oldWidget.(sliderWidget)
	r := ro.(*render.RenderSlider)
	r.SetMinValue(s.MinValue)
	r.SetMaxValue(s.MaxValue)
	r.SetValue(s.Value)
	r.OnChange = s.OnChange
	if old.MinValue != s.MinValue || old.MaxValue != s.MaxValue || old.Value != s.Value {
		// 确保布局刷新
		r.MarkNeedsLayout()
	}
}
