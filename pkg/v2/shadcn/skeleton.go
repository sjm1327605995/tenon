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
	return ui.NewRenderObjectElement(s)
}

// CreateRenderObject implements RenderObjectFactory.
func (s skeletonWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	theme := ui.GetTheme()
	r := render.NewRenderBox()
	r.SetBackgroundColor(newColor(theme.MutedColor))
	r.SetBorderRadius(theme.BorderRadius / 2)
	r.StyleSetWidth(s.width)
	r.StyleSetHeight(s.height)
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (s skeletonWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
	old := oldWidget.(skeletonWidget)
	r := ro.(*render.RenderBox)
	r.SetBackgroundColor(newColor(ui.GetTheme().MutedColor))
	if old.width != s.width {
		r.StyleSetWidth(s.width)
	}
	if old.height != s.height {
		r.StyleSetHeight(s.height)
	}
}
