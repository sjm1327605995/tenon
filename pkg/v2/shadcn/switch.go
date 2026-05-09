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
	return ui.NewRenderObjectElement(s)
}

// CreateRenderObject implements RenderObjectFactory.
func (s switchWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	theme := ui.GetTheme()
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
func (s switchWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
	r := ro.(*render.RenderSwitch)
	r.SetChecked(s.Checked)
	r.OnChange = s.OnChange
}
