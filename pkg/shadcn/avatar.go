package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/pkg/widgets"
)

// Avatar creates a circular user avatar with fallback text.
func Avatar(fallback string, size float32) engine.Widget {
	if size <= 0 {
		size = 40
	}
	return avatarWidget{fallback: fallback, size: size}
}

// internal widget type
type avatarWidget struct {
	engine.BaseWidget
	fallback string
	size     float32
}

func (a avatarWidget) CreateElement() engine.Element {
	return engine.NewSingleChildRenderObjectElement(a)
}

// CreateRenderObject implements RenderObjectFactory.
func (a avatarWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	r := render.NewRenderBox()
	r.StyleSetWidth(a.size)
	r.StyleSetHeight(a.size)
	r.SetBorderRadius(999) // circle (clamped to half size)
	r.StyleSetJustifyContent(engine.JustifyCenter)
	r.StyleSetAlignItems(engine.AlignCenter)
	r.SetBackgroundColor(newColor(engine.GetTheme().SecondaryColor))
	r.SetBorderColor(newColor(engine.GetTheme().BorderColor))
	r.SetBorderWidth(1)
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (a avatarWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	old := oldWidget.(avatarWidget)
	r := ro.(*render.RenderBox)
	if old.size != a.size {
		r.StyleSetWidth(a.size)
		r.StyleSetHeight(a.size)
	}
	// Theme-derived colors: always update in case theme changed
	r.SetBackgroundColor(newColor(engine.GetTheme().SecondaryColor))
	r.SetBorderColor(newColor(engine.GetTheme().BorderColor))
}

// GetChildWidget implements SingleChildProvider.
func (a avatarWidget) GetChildWidget() engine.Widget {
	return widgets.Text(a.fallback).FontSize(14).Color(newColor(engine.GetTheme().SecondaryForegroundColor))
}
