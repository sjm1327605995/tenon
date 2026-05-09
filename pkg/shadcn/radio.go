package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/pkg/widgets"
)

// Radio creates a radio button with optional label.
func Radio(label string, selected bool, onChange func(bool)) engine.Widget {
	return radioWidget{Label: label, Selected: selected, OnChange: onChange}
}

// internal widget type
type radioWidget struct {
	engine.BaseWidget
	Label    string
	Selected bool
	OnChange func(selected bool)
}

func (r radioWidget) CreateElement() engine.Element {
	return engine.NewSingleChildRenderObjectElement(r)
}

// CreateRenderObject implements RenderObjectFactory.
func (r radioWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	theme := engine.GetTheme()
	ro := render.NewRenderRadio()
	ro.Selected = r.Selected
	ro.BorderColor = theme.RadioBorderColor
	ro.FillColor = theme.RadioFillColor
	ro.InnerColor = theme.RadioInnerColor
	ro.OnChange = r.OnChange
	ro.SetOnClick(func() {
		if ro.OnChange != nil {
			ro.OnChange(!ro.Selected)
		}
	})
	return ro
}

// UpdateRenderObject implements RenderObjectUpdater.
func (r radioWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	old := oldWidget.(radioWidget)
	radioRO := ro.(*render.RenderRadio)
	radioRO.SetSelected(r.Selected)
	radioRO.OnChange = r.OnChange
	if old.Label != r.Label {
		// Label 变化由 SingleChildProvider 自动处理子元素 diff
		// RenderRadio 本身无需额外操作
	}
}

// GetChildWidget implements SingleChildProvider.
func (r radioWidget) GetChildWidget() engine.Widget {
	if r.Label == "" {
		return nil
	}
	return widgets.Text(r.Label).FontSize(14).Color(engine.GetTheme().TextColor)
}
