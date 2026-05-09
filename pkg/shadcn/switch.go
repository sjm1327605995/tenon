package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
)

// Switch creates a toggle switch.
func Switch(checked bool, onChange func(bool)) engine.Widget {
	return switchWidget{Checked: checked, OnChange: onChange}
}

// internal widget type
type switchWidget struct {
	engine.BaseWidget
	Checked  bool
	OnChange func(checked bool)
}

func (s switchWidget) CreateElement() engine.Element {
	return engine.NewRenderObjectElement(s)
}

// CreateRenderObject implements RenderObjectFactory.
func (s switchWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	theme := engine.GetTheme()
	r := render.NewRenderSwitch()
	r.Checked = s.Checked
	r.OffColor = theme.SwitchOffColor
	r.OnColor = theme.SwitchOnColor
	r.ThumbColor = theme.SwitchThumbColor
	r.OnChange = s.OnChange
	r.SetOnClick(func() {
		r.Checked = !r.Checked
		if r.OnChange != nil {
			r.OnChange(r.Checked)
		}
	})
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (s switchWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	r := ro.(*render.RenderSwitch)
	r.SetChecked(s.Checked)
	r.OnChange = s.OnChange
}
