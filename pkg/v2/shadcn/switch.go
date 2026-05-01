package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// Switch creates a toggle switch.
func Switch(checked bool, onChange func(bool)) ui.Widget {
	return switchWidget{Checked: checked, OnChange: onChange}
}

// internal widget type
type switchWidget struct {
	ui.BaseWidget
	Checked  bool
	OnChange func(checked bool)
}

func (s switchWidget) CreateElement() ui.Element {
	e := &switchElement{}
	e.RenderObjectElement.BaseElement.Init(e, s)
	return e
}

type switchElement struct {
	ui.RenderObjectElement
}

func (e *switchElement) CreateRenderObject() render.RenderObject {
	w := e.GetWidget().(switchWidget)
	theme := ui.GetTheme()
	r := render.NewRenderSwitch()
	r.Checked = w.Checked
	r.OffColor = theme.SwitchOffColor
	r.OnColor = theme.SwitchOnColor
	r.ThumbColor = theme.SwitchThumbColor
	r.OnChange = w.OnChange
	r.SetOnClick(func() {
		r.Checked = !r.Checked
		if r.OnChange != nil {
			r.OnChange(r.Checked)
		}
	})
	return r
}

func (e *switchElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(switchWidget)
	r := e.GetRenderObject().(*render.RenderSwitch)
	r.SetChecked(w.Checked)
	r.OnChange = w.OnChange
}
