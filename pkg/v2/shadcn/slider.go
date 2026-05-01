package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// Slider creates a value slider.
func Slider(min, max, value float32, onChange func(float32)) ui.Widget {
	return sliderWidget{MinValue: min, MaxValue: max, Value: value, OnChange: onChange}
}

// internal widget type
type sliderWidget struct {
	ui.BaseWidget
	MinValue float32
	MaxValue float32
	Value    float32
	OnChange func(value float32)
}

func (s sliderWidget) CreateElement() ui.Element {
	e := &sliderElement{}
	e.RenderObjectElement.BaseElement.Init(e, s)
	return e
}

type sliderElement struct {
	ui.RenderObjectElement
}

func (e *sliderElement) CreateRenderObject() render.RenderObject {
	w := e.GetWidget().(sliderWidget)
	theme := ui.GetTheme()
	r := render.NewRenderSlider()
	r.MinValue = w.MinValue
	r.MaxValue = w.MaxValue
	r.SetValue(w.Value)
	r.TrackColor = theme.SliderTrackColor
	r.FillColor = theme.SliderFillColor
	r.ThumbColor = theme.SliderThumbColor
	r.OnChange = w.OnChange
	return r
}

func (e *sliderElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(sliderWidget)
	r := e.GetRenderObject().(*render.RenderSlider)
	r.MinValue = w.MinValue
	r.MaxValue = w.MaxValue
	r.SetValue(w.Value)
	r.OnChange = w.OnChange
}
