package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// Skeleton creates a placeholder loading block.
func Skeleton(w, h float32) ui.Widget {
	if w <= 0 {
		w = 200
	}
	if h <= 0 {
		h = 20
	}
	return skeletonWidget{width: w, height: h}
}

// internal widget type
type skeletonWidget struct {
	ui.BaseWidget
	width  float32
	height float32
}

func (s skeletonWidget) CreateElement() ui.Element {
	e := &skeletonElement{}
	e.RenderObjectElement.BaseElement.Init(e, s)
	return e
}

type skeletonElement struct {
	ui.RenderObjectElement
	ro *render.RenderBox
}

func (e *skeletonElement) CreateRenderObject() render.RenderObject {
	w := e.GetWidget().(skeletonWidget)
	theme := ui.GetTheme()
	r := render.NewRenderBox()
	r.SetBackgroundColor(newColor(theme.MutedColor))
	r.SetBorderRadius(theme.BorderRadius / 2)
	r.StyleSetWidth(w.width)
	r.StyleSetHeight(w.height)
	return r
}

func (e *skeletonElement) Mount(parent ui.Element, slot int) {
	e.ro = e.CreateRenderObject().(*render.RenderBox)
	e.RenderObject = e.ro
	e.RenderObjectElement.Mount(parent, slot)
}

func (e *skeletonElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(skeletonWidget)
	r := e.ro
	old := oldWidget.(skeletonWidget)

	// Background color is theme-derived; always update in case theme changed
	r.SetBackgroundColor(newColor(ui.GetTheme().MutedColor))
	if old.width != w.width {
		r.StyleSetWidth(w.width)
	}
	if old.height != w.height {
		r.StyleSetHeight(w.height)
	}
}
