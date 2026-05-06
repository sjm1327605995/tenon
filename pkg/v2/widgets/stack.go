package widgets

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/yoga"
)

// StackWidget 是层叠布局容器，支持绝对定位子元素。
type StackWidget struct {
	ui.BaseWidget
	children   []ui.Widget
	zIndex     int
	width      float32
	height     float32
	background *render.Color
	borderRadius float32
}

// Stack 创建层叠布局容器。
func Stack(children ...ui.Widget) StackWidget {
	return StackWidget{children: children}
}

func (s StackWidget) Z(v int) StackWidget {
	s.zIndex = v
	return s
}

func (s StackWidget) W(v float32) StackWidget {
	s.width = v
	return s
}

func (s StackWidget) H(v float32) StackWidget {
	s.height = v
	return s
}

func (s StackWidget) Background(c render.Color) StackWidget {
	s.background = &c
	return s
}

func (s StackWidget) Radius(v float32) StackWidget {
	s.borderRadius = v
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
	if w.background != nil {
		r.SetBackgroundColor(w.background)
	}
	if w.borderRadius > 0 {
		r.SetBorderRadius(w.borderRadius)
	}
	if w.width > 0 {
		r.StyleSetWidth(w.width)
	}
	if w.height > 0 {
		r.StyleSetHeight(w.height)
	}
	return r
}

func (e *StackElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(StackWidget)
	old := oldWidget.(StackWidget)
	if old.zIndex != w.zIndex {
		e.ro.SetZIndex(w.zIndex)
	}
	if !render.ColorPtrEquals(old.background, w.background) {
		e.ro.SetBackgroundColor(w.background)
	}
	if old.borderRadius != w.borderRadius {
		e.ro.SetBorderRadius(w.borderRadius)
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
	center bool
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

// Center 将子元素在 Stack 中水平和垂直居中。
// 子元素需要具有固定尺寸（或通过内容确定尺寸），否则会被拉伸填满。
func (p PositionedWidget) Center() PositionedWidget {
	p.center = true
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
	if w.center {
		// 清除所有方向约束，居中将由 syncBounds 手动计算
		r.StyleSetPosition(yoga.EdgeLeft, yoga.Undefined)
		r.StyleSetPosition(yoga.EdgeRight, yoga.Undefined)
		r.StyleSetPosition(yoga.EdgeTop, yoga.Undefined)
		r.StyleSetPosition(yoga.EdgeBottom, yoga.Undefined)
	} else {
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

	if old.center != w.center {
		if w.center {
			// 清除所有方向约束，居中将由 syncBounds 手动计算
			e.ro.StyleSetPosition(yoga.EdgeLeft, yoga.Undefined)
			e.ro.StyleSetPosition(yoga.EdgeRight, yoga.Undefined)
			e.ro.StyleSetPosition(yoga.EdgeTop, yoga.Undefined)
			e.ro.StyleSetPosition(yoga.EdgeBottom, yoga.Undefined)
		} else {
			// 取消 center 时重置为未设置
			e.ro.StyleSetPosition(yoga.EdgeLeft, yoga.Undefined)
			e.ro.StyleSetPosition(yoga.EdgeRight, yoga.Undefined)
			e.ro.StyleSetPosition(yoga.EdgeTop, yoga.Undefined)
			e.ro.StyleSetPosition(yoga.EdgeBottom, yoga.Undefined)
		}
	}

	if !w.center {
		if old.left != w.left {
			if w.left >= 0 {
				e.ro.StyleSetPosition(yoga.EdgeLeft, w.left)
			} else {
				e.ro.StyleSetPosition(yoga.EdgeLeft, yoga.Undefined)
			}
		}
		if old.top != w.top {
			if w.top >= 0 {
				e.ro.StyleSetPosition(yoga.EdgeTop, w.top)
			} else {
				e.ro.StyleSetPosition(yoga.EdgeTop, yoga.Undefined)
			}
		}
		if old.right != w.right {
			if w.right >= 0 {
				e.ro.StyleSetPosition(yoga.EdgeRight, w.right)
			} else {
				e.ro.StyleSetPosition(yoga.EdgeRight, yoga.Undefined)
			}
		}
		if old.bottom != w.bottom {
			if w.bottom >= 0 {
				e.ro.StyleSetPosition(yoga.EdgeBottom, w.bottom)
			} else {
				e.ro.StyleSetPosition(yoga.EdgeBottom, yoga.Undefined)
			}
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
