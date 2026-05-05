package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// Radio creates a radio button with optional label.
func Radio(label string, selected bool, onChange func(bool)) ui.Widget {
	return radioWidget{Label: label, Selected: selected, OnChange: onChange}
}

// internal widget type
type radioWidget struct {
	ui.BaseWidget
	Label    string
	Selected bool
	OnChange func(selected bool)
}

func (r radioWidget) CreateElement() ui.Element {
	e := &radioElement{}
	e.SingleChildRenderObjectElement.RenderObjectElement.BaseElement.Init(e, r)
	return e
}

type radioElement struct {
	ui.SingleChildRenderObjectElement
	ro *render.RenderRadio
}

func (e *radioElement) CreateRenderObject() render.RenderObject {
	w := e.GetWidget().(radioWidget)
	theme := ui.GetTheme()
	r := render.NewRenderRadio()
	r.Selected = w.Selected
	r.BorderColor = theme.RadioBorderColor
	r.FillColor = theme.RadioFillColor
	r.InnerColor = theme.RadioInnerColor
	r.OnChange = w.OnChange
	r.SetOnClick(func() {
		if r.OnChange != nil {
			r.OnChange(!r.Selected)
		}
	})
	return r
}

func (e *radioElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(radioWidget)
	r := e.ro
	r.SetSelected(w.Selected)
	r.OnChange = w.OnChange
}

func (e *radioElement) Mount(parent ui.Element, slot int) {
	e.ro = e.CreateRenderObject().(*render.RenderRadio)
	e.RenderObject = e.ro
	e.SingleChildRenderObjectElement.Mount(parent, slot)
	w := e.GetWidget().(radioWidget)
	if w.Label != "" {
		child := widgets.Text(w.Label).FontSize(14).Color(ui.GetTheme().TextColor)
		e.Child = ui.UpdateChild(e, nil, child)
	}
}

func (e *radioElement) UpdateChild(oldWidget ui.Widget) {
	w := e.GetWidget().(radioWidget)
	if w.Label != "" {
		child := widgets.Text(w.Label).FontSize(14).Color(ui.GetTheme().TextColor)
		e.Child = ui.UpdateChild(e, e.Child, child)
	} else {
		e.Child = ui.UpdateChild(e, e.Child, nil)
	}
}
