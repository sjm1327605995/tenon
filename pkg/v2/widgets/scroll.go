package widgets

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// ScrollWidget 是滚动容器，可包裹一个子 Widget，当内容超出视口时支持滚动。
type ScrollWidget struct {
	ui.BaseWidget
	child      ui.Widget
	width      float32
	height     float32
	maxHeight  float32
	flexGrow   float32
	flexShrink float32
}

// Scroll 创建滚动容器。
func Scroll(child ui.Widget) ScrollWidget {
	return ScrollWidget{child: child}
}

func (s ScrollWidget) W(v float32) ScrollWidget {
	s.width = v
	return s
}

func (s ScrollWidget) H(v float32) ScrollWidget {
	s.height = v
	return s
}

func (s ScrollWidget) MaxH(v float32) ScrollWidget { return s.MaxHeight(v) }
func (s ScrollWidget) MaxHeight(v float32) ScrollWidget {
	s.maxHeight = v
	return s
}

func (s ScrollWidget) Grow(v float32) ScrollWidget {
	s.flexGrow = v
	return s
}

func (s ScrollWidget) Shrink(v float32) ScrollWidget {
	s.flexShrink = v
	return s
}

func (s ScrollWidget) CreateElement() ui.Element {
	return ui.NewSingleChildRenderObjectElement(s)
}

// CreateRenderObject implements RenderObjectFactory.
func (s ScrollWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	r := render.NewRenderScroll()
	applyScrollProps(r, ScrollWidget{}, s)
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (s ScrollWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
	old := oldWidget.(ScrollWidget)
	applyScrollProps(ro.(*render.RenderScroll), old, s)
}

// GetChildWidget implements SingleChildProvider.
func (s ScrollWidget) GetChildWidget() ui.Widget {
	return s.child
}

func applyScrollProps(r *render.RenderScroll, old, w ScrollWidget) {
	if old.width != w.width {
		if w.width > 0 {
			r.StyleSetWidth(w.width)
		} else {
			r.StyleSetWidthAuto()
		}
	}
	if old.height != w.height {
		if w.height > 0 {
			r.StyleSetHeight(w.height)
		} else {
			r.StyleSetHeightAuto()
		}
	}
	if old.maxHeight != w.maxHeight && w.maxHeight > 0 {
		r.StyleSetMaxHeight(w.maxHeight)
	}
	if old.flexGrow != w.flexGrow {
		r.StyleSetFlexGrow(w.flexGrow)
	}
	if old.flexShrink != w.flexShrink {
		r.StyleSetFlexShrink(w.flexShrink)
	}
}
