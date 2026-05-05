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
	e := &avatarElement{}
	e.SingleChildRenderObjectElement.RenderObjectElement.BaseElement.Init(e, a)
	return e
}

type avatarElement struct {
	ui.SingleChildRenderObjectElement
	ro *render.RenderBox
}

func (e *avatarElement) CreateRenderObject() render.RenderObject {
	r := render.NewRenderBox()
	w := e.GetWidget().(avatarWidget)
	r.StyleSetWidth(w.size)
	r.StyleSetHeight(w.size)
	r.SetBorderRadius(999) // circle (clamped to half size)
	r.StyleSetJustifyContent(ui.JustifyCenter)
	r.StyleSetAlignItems(ui.AlignCenter)
	r.SetBackgroundColor(newColor(ui.GetTheme().SecondaryColor))
	r.SetBorderColor(newColor(ui.GetTheme().BorderColor))
	r.SetBorderWidth(1)
	return r
}

func (e *avatarElement) UpdateRenderObject(oldWidget ui.Widget) {
	r := e.ro
	w := e.GetWidget().(avatarWidget)
	old := oldWidget.(avatarWidget)

	if old.size != w.size {
		r.StyleSetWidth(w.size)
		r.StyleSetHeight(w.size)
	}
	// Theme-derived colors: always update in case theme changed
	r.SetBackgroundColor(newColor(ui.GetTheme().SecondaryColor))
	r.SetBorderColor(newColor(ui.GetTheme().BorderColor))
}

func (e *avatarElement) Mount(parent ui.Element, slot int) {
	e.ro = e.CreateRenderObject().(*render.RenderBox)
	e.RenderObject = e.ro
	e.SingleChildRenderObjectElement.Mount(parent, slot)
	w := e.GetWidget().(avatarWidget)
	child := widgets.Text(w.fallback).FontSize(14).Color(newColor(ui.GetTheme().SecondaryForegroundColor))
	e.Child = ui.UpdateChild(e, nil, child)
}

func (e *avatarElement) UpdateChild(oldWidget ui.Widget) {
	w := e.GetWidget().(avatarWidget)
	child := widgets.Text(w.fallback).FontSize(14).Color(newColor(ui.GetTheme().SecondaryForegroundColor))
	e.Child = ui.UpdateChild(e, e.Child, child)
}
