package widgets

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/yoga"
)

// StackWidget 是层叠布局容器，支持绝对定位子元素。
type StackWidget struct {
	ui.BaseWidget
	children []ui.Widget
	zIndex   int
}

// Stack 创建层叠布局容器。
func Stack(children ...ui.Widget) StackWidget {
	return StackWidget{children: children}
}

func (s StackWidget) Z(v int) StackWidget {
	s.zIndex = v
	return s
}

func (s StackWidget) CreateElement() ui.Element {
	e := &StackElement{}
	e.MultiChildRenderObjectElement.RenderObjectElement.BaseElement.Init(e, s)
	return e
}

// StackElement 是 StackWidget 对应的 Element。
type StackElement struct {
	ui.MultiChildRenderObjectElement
	ro *render.RenderStack
}

func (e *StackElement) CreateRenderObject() render.RenderObject {
	r := render.NewRenderStack()
	w := e.GetWidget().(StackWidget)
	r.SetZIndex(w.zIndex)
	return r
}

func (e *StackElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(StackWidget)
	old := oldWidget.(StackWidget)
	if old.zIndex != w.zIndex {
		e.ro.SetZIndex(w.zIndex)
	}
}

func (e *StackElement) UpdateChildren(oldWidget ui.Widget) {
	w := e.GetWidget().(StackWidget)
	newWidgets := make([]ui.Widget, len(w.children))
	for i, child := range w.children {
		newWidgets[i] = child
	}
	e.Children = ui.UpdateChildren(e, e.Children, newWidgets)
}

func (e *StackElement) Mount(parent ui.Element, slot int) {
	e.ro = e.CreateRenderObject().(*render.RenderStack)
	e.RenderObject = e.ro
	e.MultiChildRenderObjectElement.Mount(parent, slot)
	w := e.GetWidget().(StackWidget)
	newWidgets := make([]ui.Widget, len(w.children))
	for i, child := range w.children {
		newWidgets[i] = child
	}
	e.Children = ui.UpdateChildren(e, nil, newWidgets)
}

// PositionedWidget 用于在 Stack 中对子元素进行绝对定位。
type PositionedWidget struct {
	ui.BaseWidget
	child  ui.Widget
	left   float32
	top    float32
	right  float32
	bottom float32
	width  float32
	height float32
	zIndex int
}

// Positioned 创建绝对定位包装器。
func Positioned(child ui.Widget) PositionedWidget {
	return PositionedWidget{child: child, left: -1, top: -1, right: -1, bottom: -1}
}

func (p PositionedWidget) L(v float32) PositionedWidget {
	p.left = v
	return p
}

func (p PositionedWidget) T(v float32) PositionedWidget {
	p.top = v
	return p
}

func (p PositionedWidget) R(v float32) PositionedWidget {
	p.right = v
	return p
}

func (p PositionedWidget) B(v float32) PositionedWidget {
	p.bottom = v
	return p
}

func (p PositionedWidget) W(v float32) PositionedWidget {
	p.width = v
	return p
}

func (p PositionedWidget) H(v float32) PositionedWidget {
	p.height = v
	return p
}

func (p PositionedWidget) Z(v int) PositionedWidget {
	p.zIndex = v
	return p
}

func (p PositionedWidget) CreateElement() ui.Element {
	e := &PositionedElement{}
	e.SingleChildRenderObjectElement.RenderObjectElement.BaseElement.Init(e, p)
	return e
}

// PositionedElement 是 PositionedWidget 对应的 Element。
type PositionedElement struct {
	ui.SingleChildRenderObjectElement
	ro *render.RenderBox
}

func (e *PositionedElement) Mount(parent ui.Element, slot int) {
	e.ro = e.CreateRenderObject().(*render.RenderBox)
	e.RenderObject = e.ro
	e.SingleChildRenderObjectElement.Mount(parent, slot)
	p := e.GetWidget().(PositionedWidget)
	if p.child != nil {
		e.Child = ui.UpdateChild(e, nil, p.child)
	}
}

func (e *PositionedElement) CreateRenderObject() render.RenderObject {
	w := e.GetWidget().(PositionedWidget)
	r := render.NewRenderBox()
	r.StyleSetPositionType(yoga.PositionTypeAbsolute)
	if w.left >= 0 {
		r.StyleSetPosition(yoga.EdgeLeft, w.left)
	}
	if w.top >= 0 {
		r.StyleSetPosition(yoga.EdgeTop, w.top)
	}
	if w.right >= 0 {
		r.StyleSetPosition(yoga.EdgeRight, w.right)
	}
	if w.bottom >= 0 {
		r.StyleSetPosition(yoga.EdgeBottom, w.bottom)
	}
	if w.width > 0 {
		r.StyleSetWidth(w.width)
	}
	if w.height > 0 {
		r.StyleSetHeight(w.height)
	}
	r.SetZIndex(w.zIndex)
	return r
}

func (e *PositionedElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(PositionedWidget)
	old := oldWidget.(PositionedWidget)

	if old.left != w.left {
		if w.left >= 0 {
		e.ro.StyleSetPosition(yoga.EdgeLeft, w.left)
		} else {
		e.ro.StyleSetPosition(yoga.EdgeLeft, 0)
		}
	}
	if old.top != w.top {
		if w.top >= 0 {
		e.ro.StyleSetPosition(yoga.EdgeTop, w.top)
		} else {
		e.ro.StyleSetPosition(yoga.EdgeTop, 0)
		}
	}
	if old.right != w.right {
		if w.right >= 0 {
		e.ro.StyleSetPosition(yoga.EdgeRight, w.right)
		} else {
		e.ro.StyleSetPosition(yoga.EdgeRight, 0)
		}
	}
	if old.bottom != w.bottom {
		if w.bottom >= 0 {
		e.ro.StyleSetPosition(yoga.EdgeBottom, w.bottom)
		} else {
		e.ro.StyleSetPosition(yoga.EdgeBottom, 0)
		}
	}
	if old.width != w.width {
		if w.width > 0 {
		e.ro.StyleSetWidth(w.width)
		} else {
		e.ro.StyleSetWidthAuto()
		}
	}
	if old.height != w.height {
		if w.height > 0 {
		e.ro.StyleSetHeight(w.height)
		} else {
		e.ro.StyleSetHeightAuto()
		}
	}
	if old.zIndex != w.zIndex {
	e.ro.SetZIndex(w.zIndex)
	}
}

func (e *PositionedElement) UpdateChild(oldWidget ui.Widget) {
	p := e.GetWidget().(PositionedWidget)
	e.Child = ui.UpdateChild(e, e.Child, p.child)
}
