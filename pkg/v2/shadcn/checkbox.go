package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// Checkbox creates a checkable box with optional label.
func Checkbox(label string, checked bool, onChange func(bool)) ui.Widget {
	return checkboxWidget{Label: label, Checked: checked, OnChange: onChange}
}

// internal widget type
type checkboxWidget struct {
	ui.BaseWidget
	Label    string
	Checked  bool
	OnChange func(checked bool)
}

func (c checkboxWidget) CreateElement() ui.Element {
	return ui.NewSingleChildRenderObjectElement(c)
}

// CreateRenderObject implements RenderObjectFactory.
func (c checkboxWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	theme := ui.GetTheme()
	r := render.NewRenderCheckbox()
	r.Checked = c.Checked
	r.BorderColor = theme.CheckboxBorderColor
	r.FillColor = theme.CheckboxFillColor
	r.CheckColor = theme.CheckboxCheckColor
	r.OnChange = c.OnChange
	r.SetOnClick(func() {
		r.Checked = !r.Checked
		if r.OnChange != nil {
			r.OnChange(r.Checked)
		}
	})
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (c checkboxWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
	old := oldWidget.(checkboxWidget)
	r := ro.(*render.RenderCheckbox)
	r.SetChecked(c.Checked)
	r.OnChange = c.OnChange
	if old.Label != c.Label {
		// Label 变化由 SingleChildProvider 自动处理子元素 diff
	}
}

// GetChildWidget implements SingleChildProvider.
func (c checkboxWidget) GetChildWidget() ui.Widget {
	if c.Label == "" {
		return nil
	}
	return widgets.Text(c.Label).FontSize(14).Color(ui.GetTheme().TextColor)
}
