package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// Avatar creates a circular user avatar with fallback text.
func Avatar(fallback string, size float32) ui.Widget {
	if size <= 0 {
		size = 40
	}
	return avatarWidget{fallback: fallback, size: size}
}

// internal widget type
type avatarWidget struct {
	ui.BaseWidget
	fallback string
	size     float32
}

func (a avatarWidget) CreateElement() ui.Element {
	return ui.NewSingleChildRenderObjectElement(a)
}

// CreateRenderObject implements RenderObjectFactory.
func (a avatarWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	r := render.NewRenderBox()
	r.StyleSetWidth(a.size)
	r.StyleSetHeight(a.size)
	r.SetBorderRadius(999) // circle (clamped to half size)
	r.StyleSetJustifyContent(ui.JustifyCenter)
	r.StyleSetAlignItems(ui.AlignCenter)
	r.SetBackgroundColor(newColor(ui.GetTheme().SecondaryColor))
	r.SetBorderColor(newColor(ui.GetTheme().BorderColor))
	r.SetBorderWidth(1)
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (a avatarWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
	old := oldWidget.(avatarWidget)
	r := ro.(*render.RenderBox)
	if old.size != a.size {
		r.StyleSetWidth(a.size)
		r.StyleSetHeight(a.size)
	}
	// Theme-derived colors: always update in case theme changed
	r.SetBackgroundColor(newColor(ui.GetTheme().SecondaryColor))
	r.SetBorderColor(newColor(ui.GetTheme().BorderColor))
}

// GetChildWidget implements SingleChildProvider.
func (a avatarWidget) GetChildWidget() ui.Widget {
	return widgets.Text(a.fallback).FontSize(14).Color(newColor(ui.GetTheme().SecondaryForegroundColor))
}
