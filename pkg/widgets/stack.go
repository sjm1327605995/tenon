package widgets

import (
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/yoga"
)

// StackWidget 是层叠布局容器，支持绝对定位子元素。
type StackWidget struct {
	engine.BaseWidget
	children   []engine.Widget
	zIndex     int
	width      float32
	height     float32
	background *render.Color
	borderRadius float32
}

// Stack 创建层叠布局容器。
func Stack(children ...engine.Widget) StackWidget {
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

func (s StackWidget) CreateElement() engine.Element {
	return engine.NewMultiChildRenderObjectElement(s)
}

// CreateRenderObject implements RenderObjectFactory.
func (s StackWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	r := render.NewRenderStack()
	r.SetZIndex(s.zIndex)
	if s.background != nil {
		r.SetBackgroundColor(s.background)
	}
	if s.borderRadius > 0 {
		r.SetBorderRadius(s.borderRadius)
	}
	if s.width > 0 {
		r.StyleSetWidth(s.width)
	}
	if s.height > 0 {
		r.StyleSetHeight(s.height)
	}
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (s StackWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	r := ro.(*render.RenderStack)
	old := oldWidget.(StackWidget)
	if old.zIndex != s.zIndex {
		r.SetZIndex(s.zIndex)
	}
	if !render.ColorPtrEquals(old.background, s.background) {
		r.SetBackgroundColor(s.background)
	}
	if old.borderRadius != s.borderRadius {
		r.SetBorderRadius(s.borderRadius)
	}
	if old.width != s.width {
		if s.width > 0 {
			r.StyleSetWidth(s.width)
		} else {
			r.StyleSetWidthAuto()
		}
	}
	if old.height != s.height {
		if s.height > 0 {
			r.StyleSetHeight(s.height)
		} else {
			r.StyleSetHeightAuto()
		}
	}
}

// GetChildrenWidgets implements MultiChildProvider.
func (s StackWidget) GetChildrenWidgets() []engine.Widget {
	return s.children
}

// PositionedWidget 用于在 Stack 中对子元素进行绝对定位。
type PositionedWidget struct {
	engine.BaseWidget
	child  engine.Widget
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
func Positioned(child engine.Widget) PositionedWidget {
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

func (p PositionedWidget) CreateElement() engine.Element {
	return engine.NewSingleChildRenderObjectElement(p)
}

// CreateRenderObject implements RenderObjectFactory.
func (p PositionedWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	r := render.NewRenderBox()
	r.StyleSetPositionType(yoga.PositionTypeAbsolute)
	if p.center {
		// 清除所有方向约束，居中将由 syncBounds 手动计算
		r.StyleSetPosition(yoga.EdgeLeft, yoga.Undefined)
		r.StyleSetPosition(yoga.EdgeRight, yoga.Undefined)
		r.StyleSetPosition(yoga.EdgeTop, yoga.Undefined)
		r.StyleSetPosition(yoga.EdgeBottom, yoga.Undefined)
	} else {
		if p.left >= 0 {
			r.StyleSetPosition(yoga.EdgeLeft, p.left)
		}
		if p.top >= 0 {
			r.StyleSetPosition(yoga.EdgeTop, p.top)
		}
		if p.right >= 0 {
			r.StyleSetPosition(yoga.EdgeRight, p.right)
		}
		if p.bottom >= 0 {
			r.StyleSetPosition(yoga.EdgeBottom, p.bottom)
		}
	}
	if p.width > 0 {
		r.StyleSetWidth(p.width)
	}
	if p.height > 0 {
		r.StyleSetHeight(p.height)
	}
	r.SetZIndex(p.zIndex)
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (p PositionedWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	r := ro.(*render.RenderBox)
	old := oldWidget.(PositionedWidget)

	if old.center != p.center {
		if p.center {
			// 清除所有方向约束，居中将由 syncBounds 手动计算
			r.StyleSetPosition(yoga.EdgeLeft, yoga.Undefined)
			r.StyleSetPosition(yoga.EdgeRight, yoga.Undefined)
			r.StyleSetPosition(yoga.EdgeTop, yoga.Undefined)
			r.StyleSetPosition(yoga.EdgeBottom, yoga.Undefined)
		} else {
			// 取消 center 时重置为未设置
			r.StyleSetPosition(yoga.EdgeLeft, yoga.Undefined)
			r.StyleSetPosition(yoga.EdgeRight, yoga.Undefined)
			r.StyleSetPosition(yoga.EdgeTop, yoga.Undefined)
			r.StyleSetPosition(yoga.EdgeBottom, yoga.Undefined)
		}
	}

	if !p.center {
		if old.left != p.left {
			if p.left >= 0 {
				r.StyleSetPosition(yoga.EdgeLeft, p.left)
			} else {
				r.StyleSetPosition(yoga.EdgeLeft, yoga.Undefined)
			}
		}
		if old.top != p.top {
			if p.top >= 0 {
				r.StyleSetPosition(yoga.EdgeTop, p.top)
			} else {
				r.StyleSetPosition(yoga.EdgeTop, yoga.Undefined)
			}
		}
		if old.right != p.right {
			if p.right >= 0 {
				r.StyleSetPosition(yoga.EdgeRight, p.right)
			} else {
				r.StyleSetPosition(yoga.EdgeRight, yoga.Undefined)
			}
		}
		if old.bottom != p.bottom {
			if p.bottom >= 0 {
				r.StyleSetPosition(yoga.EdgeBottom, p.bottom)
			} else {
				r.StyleSetPosition(yoga.EdgeBottom, yoga.Undefined)
			}
		}
	}
	if old.width != p.width {
		if p.width > 0 {
			r.StyleSetWidth(p.width)
		} else {
			r.StyleSetWidthAuto()
		}
	}
	if old.height != p.height {
		if p.height > 0 {
			r.StyleSetHeight(p.height)
		} else {
			r.StyleSetHeightAuto()
		}
	}
	if old.zIndex != p.zIndex {
		r.SetZIndex(p.zIndex)
	}
}

// GetChildWidget implements SingleChildProvider.
func (p PositionedWidget) GetChildWidget() engine.Widget {
	return p.child
}
