package widgets

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/yoga"
)

// ContainerWidget 是一个全能的装饰性盒子，可设置背景、圆角、边框、阴影、尺寸等。
type ContainerWidget struct {
	ui.BaseWidget
	child            ui.Widget
	backgroundColor  *render.Color
	backgroundImage  *ebiten.Image
	backgroundSlice  render.BorderSlice
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
	zIndex           int
	rotation         float32
	scaleX           float32
	scaleY           float32
	alpha            float32
	focusable        bool
}

// Container 创建容器 Widget，包裹一个子 Widget。
func Container(child ui.Widget) ContainerWidget {
	return ContainerWidget{child: child}
}

func (c ContainerWidget) Background(color render.Color) ContainerWidget {
	c.backgroundColor = &color
	return c
}

// BackgroundImage 设置 NinePatch 或普通拉伸背景图片。
// slice 全为 0 时使用普通拉伸，否则使用 9-slice 缩放。
func (c ContainerWidget) BackgroundImage(img *ebiten.Image, slice render.BorderSlice) ContainerWidget {
	c.backgroundImage = img
	c.backgroundSlice = slice
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

func (c ContainerWidget) Pad(insets ui.EdgeInsets) ContainerWidget { return c.Padding(insets) }
func (c ContainerWidget) Padding(insets ui.EdgeInsets) ContainerWidget {
	c.padding = insets
	return c
}

func (c ContainerWidget) Marginf(insets ui.EdgeInsets) ContainerWidget { return c.Margin(insets) }
func (c ContainerWidget) Margin(insets ui.EdgeInsets) ContainerWidget {
	c.margin = insets
	return c
}

func (c ContainerWidget) OnTap(fn func()) ContainerWidget {
	c.onClick = fn
	return c
}

func (c ContainerWidget) ZIndex(v int) ContainerWidget {
	c.zIndex = v
	return c
}

func (c ContainerWidget) Rotate(deg float32) ContainerWidget {
	c.rotation = deg
	return c
}

func (c ContainerWidget) Scale(x, y float32) ContainerWidget {
	c.scaleX = x
	c.scaleY = y
	return c
}

func (c ContainerWidget) Opacity(a float32) ContainerWidget {
	if a < 0 {
		a = 0
	}
	if a > 1 {
		a = 1
	}
	c.alpha = a
	return c
}

func (c ContainerWidget) Focusable(v bool) ContainerWidget {
	c.focusable = v
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
	return ui.NewSingleChildRenderObjectElement(c)
}

// CreateRenderObject implements RenderObjectFactory.
func (c ContainerWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	r := render.NewRenderBox()
	applyContainerProps(r, ContainerWidget{}, c)
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (c ContainerWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
	old := oldWidget.(ContainerWidget)
	applyContainerProps(ro.(*render.RenderBox), old, c)
}

// GetChildWidget implements SingleChildProvider.
func (c ContainerWidget) GetChildWidget() ui.Widget {
	return c.child
}

func applyContainerProps(r *render.RenderBox, old, w ContainerWidget) {
	if !render.ColorPtrEquals(old.backgroundColor, w.backgroundColor) {
		r.SetBackgroundColor(w.backgroundColor)
	}
	if old.backgroundImage != w.backgroundImage || old.backgroundSlice != w.backgroundSlice {
		r.SetBackgroundImage(w.backgroundImage, w.backgroundSlice)
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
	r.SetOnClick(w.onClick)
	if old.focusable != w.focusable {
		r.SetFocusable(w.focusable)
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

	if old.zIndex != w.zIndex {
		r.SetZIndex(w.zIndex)
	}

	if old.rotation != w.rotation || old.scaleX != w.scaleX || old.scaleY != w.scaleY || old.alpha != w.alpha {
		t := r.GetTransform()
		t.Rotation = w.rotation
		if w.scaleX != 0 || w.scaleY != 0 {
			t.ScaleX = w.scaleX
			t.ScaleY = w.scaleY
		}
		t.Alpha = w.alpha
		if t.OriginX == 0 && t.OriginY == 0 {
			t.OriginX = 0.5
			t.OriginY = 0.5
		}
		r.SetTransform(t)
	}
}
