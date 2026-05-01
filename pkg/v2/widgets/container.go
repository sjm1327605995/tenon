package widgets

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/yoga"
)

// ContainerWidget 是一个全能的装饰性盒子，可设置背景、圆角、边框、阴影、尺寸等。
type ContainerWidget struct {
	ui.BaseWidget
	child            ui.Widget
	backgroundColor  *render.Color
	borderRadius     float32
	borderColor      *render.Color
	borderWidth      float32
	shadowColor      *render.Color
	shadowBlur       float32
	shadowOffsetX    float32
	shadowOffsetY    float32
	width            float32
	height           float32
	widthPercent     float32
	heightPercent    float32
	flexGrow         float32
	flexShrink       float32
	flexBasis        float32
	padding          ui.EdgeInsets
	margin           ui.EdgeInsets
	justify          yoga.Justify
	alignItemsVal    yoga.Align
	onClick          func()
}

// Container 创建容器 Widget，包裹一个子 Widget。
func Container(child ui.Widget) ContainerWidget {
	return ContainerWidget{child: child}
}

func (c ContainerWidget) Background(color render.Color) ContainerWidget {
	c.backgroundColor = &color
	return c
}

func (c ContainerWidget) Radius(v float32) ContainerWidget {
	c.borderRadius = v
	return c
}

func (c ContainerWidget) Border(color render.Color, width float32) ContainerWidget {
	c.borderColor = &color
	c.borderWidth = width
	return c
}

func (c ContainerWidget) Shadow(color render.Color, blur, offsetX, offsetY float32) ContainerWidget {
	c.shadowColor = &color
	c.shadowBlur = blur
	c.shadowOffsetX = offsetX
	c.shadowOffsetY = offsetY
	return c
}

func (c ContainerWidget) W(v float32) ContainerWidget {
	c.width = v
	return c
}

func (c ContainerWidget) H(v float32) ContainerWidget {
	c.height = v
	return c
}

func (c ContainerWidget) WPct(v float32) ContainerWidget {
	c.widthPercent = v
	return c
}

func (c ContainerWidget) HPct(v float32) ContainerWidget {
	c.heightPercent = v
	return c
}

func (c ContainerWidget) Grow(v float32) ContainerWidget {
	c.flexGrow = v
	return c
}

func (c ContainerWidget) Shrink(v float32) ContainerWidget {
	c.flexShrink = v
	return c
}

func (c ContainerWidget) Basis(v float32) ContainerWidget {
	c.flexBasis = v
	return c
}

func (c ContainerWidget) Pad(insets ui.EdgeInsets) ContainerWidget {
	c.padding = insets
	return c
}

func (c ContainerWidget) Marginf(insets ui.EdgeInsets) ContainerWidget {
	c.margin = insets
	return c
}

func (c ContainerWidget) OnTap(fn func()) ContainerWidget {
	c.onClick = fn
	return c
}

func (c ContainerWidget) JustifyContent(v yoga.Justify) ContainerWidget {
	c.justify = v
	return c
}

func (c ContainerWidget) AlignItems(v yoga.Align) ContainerWidget {
	c.alignItemsVal = v
	return c
}

func (c ContainerWidget) CreateElement() ui.Element {
	e := &ContainerElement{}
	e.SingleChildRenderObjectElement.RenderObjectElement.BaseElement.Init(e, c)
	return e
}

// ContainerElement 是 ContainerWidget 对应的 Element。
type ContainerElement struct {
	ui.SingleChildRenderObjectElement
}

func (e *ContainerElement) CreateRenderObject() render.RenderObject {
	r := render.NewRenderBox()
	applyContainerProps(r, ContainerWidget{}, e.GetWidget().(ContainerWidget))
	return r
}

func (e *ContainerElement) UpdateRenderObject(oldWidget ui.Widget) {
	r := e.GetRenderObject().(*render.RenderBox)
	old := oldWidget.(ContainerWidget)
	w := e.GetWidget().(ContainerWidget)
	applyContainerProps(r, old, w)
}

func (e *ContainerElement) UpdateChild(oldWidget ui.Widget) {
	w := e.GetWidget().(ContainerWidget)
	e.Child = ui.UpdateChild(e, e.Child, w.child)
}

func (e *ContainerElement) Mount(parent ui.Element, slot int) {
	e.RenderObject = e.CreateRenderObject()
	e.SingleChildRenderObjectElement.Mount(parent, slot)
	w := e.GetWidget().(ContainerWidget)
	if w.child != nil {
		e.Child = ui.UpdateChild(e, nil, w.child)
	}
}

func applyContainerProps(r *render.RenderBox, old, w ContainerWidget) {
	if !render.ColorPtrEquals(old.backgroundColor, w.backgroundColor) {
		r.SetBackgroundColor(w.backgroundColor)
	}
	if old.borderRadius != w.borderRadius {
		r.SetBorderRadius(w.borderRadius)
	}
	if !render.ColorPtrEquals(old.borderColor, w.borderColor) {
		r.SetBorderColor(w.borderColor)
	}
	if old.borderWidth != w.borderWidth {
		r.SetBorderWidth(w.borderWidth)
	}
	if !render.ColorPtrEquals(old.shadowColor, w.shadowColor) || old.shadowBlur != w.shadowBlur ||
		old.shadowOffsetX != w.shadowOffsetX || old.shadowOffsetY != w.shadowOffsetY {
		r.SetShadow(w.shadowColor, w.shadowBlur, w.shadowOffsetX, w.shadowOffsetY)
	}
	if (old.onClick == nil) != (w.onClick == nil) {
		r.SetOnClick(w.onClick)
	}

	if old.width != w.width {
		if w.width > 0 {
			r.StyleSetWidth(w.width)
		} else {
			r.StyleSetWidthAuto()
		}
	}
	if old.widthPercent != w.widthPercent {
		if w.widthPercent > 0 {
			r.StyleSetWidthPercent(w.widthPercent)
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
	if old.heightPercent != w.heightPercent {
		if w.heightPercent > 0 {
			r.StyleSetHeightPercent(w.heightPercent)
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
	if old.flexBasis != w.flexBasis {
		if w.flexBasis > 0 {
			r.StyleSetFlexBasis(w.flexBasis)
		} else {
			r.StyleSetFlexBasisAuto()
		}
	}
	if old.padding != w.padding {
		r.StyleSetPadding(ui.EdgeTop, w.padding.Top)
		r.StyleSetPadding(ui.EdgeRight, w.padding.Right)
		r.StyleSetPadding(ui.EdgeBottom, w.padding.Bottom)
		r.StyleSetPadding(ui.EdgeLeft, w.padding.Left)
	}
	if old.margin != w.margin {
		r.StyleSetMargin(ui.EdgeTop, w.margin.Top)
		r.StyleSetMargin(ui.EdgeRight, w.margin.Right)
		r.StyleSetMargin(ui.EdgeBottom, w.margin.Bottom)
		r.StyleSetMargin(ui.EdgeLeft, w.margin.Left)
	}
	if old.justify != w.justify {
		r.StyleSetJustifyContent(w.justify)
	}
	if old.alignItemsVal != w.alignItemsVal {
		r.StyleSetAlignItems(w.alignItemsVal)
	}
}
