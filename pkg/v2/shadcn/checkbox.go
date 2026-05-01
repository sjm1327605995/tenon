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
	e := &checkboxElement{}
	e.SingleChildRenderObjectElement.RenderObjectElement.BaseElement.Init(e, c)
	return e
}

type checkboxElement struct {
	ui.SingleChildRenderObjectElement
}

func (e *checkboxElement) CreateRenderObject() render.RenderObject {
	w := e.GetWidget().(checkboxWidget)
	theme := ui.GetTheme()
	r := render.NewRenderCheckbox()
	r.Checked = w.Checked
	r.BorderColor = theme.CheckboxBorderColor
	r.FillColor = theme.CheckboxFillColor
	r.CheckColor = theme.CheckboxCheckColor
	r.OnChange = w.OnChange
	r.SetOnClick(func() {
		r.Checked = !r.Checked
		if r.OnChange != nil {
			r.OnChange(r.Checked)
		}
	})
	return r
}

func (e *checkboxElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(checkboxWidget)
	r := e.GetRenderObject().(*render.RenderCheckbox)
	r.SetChecked(w.Checked)
	r.OnChange = w.OnChange
}

func (e *checkboxElement) Mount(parent ui.Element, slot int) {
	e.RenderObject = e.CreateRenderObject()
	e.SingleChildRenderObjectElement.Mount(parent, slot)
	w := e.GetWidget().(checkboxWidget)
	if w.Label != "" {
		child := widgets.Text(w.Label).FontSize(14).Color(ui.GetTheme().TextColor)
		e.Child = ui.UpdateChild(e, nil, child)
	}
}

func (e *checkboxElement) UpdateChild(oldWidget ui.Widget) {
	w := e.GetWidget().(checkboxWidget)
	if w.Label != "" {
		child := widgets.Text(w.Label).FontSize(14).Color(ui.GetTheme().TextColor)
		e.Child = ui.UpdateChild(e, e.Child, child)
	} else {
		e.Child = ui.UpdateChild(e, e.Child, nil)
	}
}
