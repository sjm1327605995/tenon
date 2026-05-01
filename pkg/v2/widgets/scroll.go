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

func (s ScrollWidget) Grow(v float32) ScrollWidget {
	s.flexGrow = v
	return s
}

func (s ScrollWidget) Shrink(v float32) ScrollWidget {
	s.flexShrink = v
	return s
}

func (s ScrollWidget) CreateElement() ui.Element {
	e := &ScrollElement{}
	e.SingleChildRenderObjectElement.RenderObjectElement.BaseElement.Init(e, s)
	return e
}

// ScrollElement 是 ScrollWidget 对应的 Element。
type ScrollElement struct {
	ui.SingleChildRenderObjectElement
}

func (e *ScrollElement) CreateRenderObject() render.RenderObject {
	r := render.NewRenderScroll()
	w := e.GetWidget().(ScrollWidget)
	applyScrollProps(r, ScrollWidget{}, w)
	return r
}

func (e *ScrollElement) UpdateRenderObject(oldWidget ui.Widget) {
	r := e.GetRenderObject().(*render.RenderScroll)
	old := oldWidget.(ScrollWidget)
	w := e.GetWidget().(ScrollWidget)
	applyScrollProps(r, old, w)
}

func (e *ScrollElement) UpdateChild(oldWidget ui.Widget) {
	w := e.GetWidget().(ScrollWidget)
	e.Child = ui.UpdateChild(e, e.Child, w.child)
}

func (e *ScrollElement) Mount(parent ui.Element, slot int) {
	e.RenderObject = e.CreateRenderObject()
	e.SingleChildRenderObjectElement.Mount(parent, slot)
	w := e.GetWidget().(ScrollWidget)
	if w.child != nil {
		e.Child = ui.UpdateChild(e, nil, w.child)
	}
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
	if old.flexGrow != w.flexGrow {
		r.StyleSetFlexGrow(w.flexGrow)
	}
	if old.flexShrink != w.flexShrink {
		r.StyleSetFlexShrink(w.flexShrink)
	}
}
