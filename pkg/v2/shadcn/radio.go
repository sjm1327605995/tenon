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
	return ui.NewSingleChildRenderObjectElement(r)
}

// CreateRenderObject implements RenderObjectFactory.
func (r radioWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	theme := ui.GetTheme()
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
func (r radioWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
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
func (r radioWidget) GetChildWidget() ui.Widget {
	if r.Label == "" {
		return nil
	}
	return widgets.Text(r.Label).FontSize(14).Color(ui.GetTheme().TextColor)
}
