package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
)

// Skeleton creates a placeholder loading block.
func Skeleton(w, h float32) engine.Widget {
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
	engine.BaseWidget
	width  float32
	height float32
}

func (s skeletonWidget) CreateElement() engine.Element {
	return engine.NewRenderObjectElement(s)
}

// CreateRenderObject implements RenderObjectFactory.
func (s skeletonWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	theme := engine.GetTheme()
	r := render.NewRenderBox()
	r.SetBackgroundColor(newColor(theme.MutedColor))
	r.SetBorderRadius(theme.BorderRadius / 2)
	r.StyleSetWidth(s.width)
	r.StyleSetHeight(s.height)
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (s skeletonWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	old := oldWidget.(skeletonWidget)
	r := ro.(*render.RenderBox)
	r.SetBackgroundColor(newColor(engine.GetTheme().MutedColor))
	if old.width != s.width {
		r.StyleSetWidth(s.width)
	}
	if old.height != s.height {
		r.StyleSetHeight(s.height)
	}
}
