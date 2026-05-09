package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/yoga"
)

// RenderBox 是所有需要 Yoga 布局的 RenderObject 的基类。
// 提供基本的矩形边界、背景色、圆角、边框、阴影、变换、裁剪等绘制能力。
type RenderBox struct {
	BaseRenderObject

	BackgroundColor    *Color
	BorderRadius       BorderRadius
	BorderColor        *Color
	BorderWidth        float32
	FocusedBorderColor *Color
	ShadowColor        *Color
	ShadowBlur         float32
	ShadowOffsetX      float32
	ShadowOffsetY      float32
	clipChildren       bool
	focusable          bool

	// BackgroundImage 是背景图片；与 BackgroundColor 互斥，优先绘制图片。
	BackgroundImage *ebiten.Image
	// BackgroundSlice 是 9-slice 切图参数；非零值时使用 NinePatch 绘制，否则普通拉伸。
	BackgroundSlice BorderSlice
}

func NewRenderBox() *RenderBox {
	r := &RenderBox{}
	r.BaseRenderObject.Init(r)
	r.yoga = yoga.NewNode()
	r.transform = DefaultTransform()
	return r
}

func (r *RenderBox) PerformLayout() {
	// RenderBox 默认不处理特殊布局，子类（如 RenderFlex）覆盖
}

// Focus/Blur/IsFocused 让 RenderBox 实现 Focusabler 接口，
// 使 Container 等组件也能参与焦点管理。
func (r *RenderBox) Focus() {
	r.focused = true
	r.MarkNeedsPaint()
}

func (r *RenderBox) Blur() {
	r.focused = false
	r.MarkNeedsPaint()
}

func (r *RenderBox) IsFocused() bool {
	return r.focused
}

func (r *RenderBox) SetFocusable(v bool) {
	r.focusable = v
}

func (r *RenderBox) IsFocusable() bool {
	return r.focusable
}

func (r *RenderBox) Paint(screen *ebiten.Image, offset Offset) {
	bounds := r.bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	x := offset.X + bounds.X
	y := offset.Y + bounds.Y

	// 变换由 Engine.drawRenderObjectWithTransform 统一处理，
	// Paint 只负责绘制内容（背景、阴影、边框）。
	r.paintContent(screen, Bounds{X: x, Y: y, Width: bounds.Width, Height: bounds.Height})
}

// paintContent 绘制背景、阴影、边框等内容（不含变换和裁剪）。
func (r *RenderBox) paintContent(screen *ebiten.Image, b Bounds) {
	x, y, w, h := b.X, b.Y, b.Width, b.Height

	var shadowColor color.Color
	if r.ShadowColor != nil {
		shadowColor = r.ShadowColor
	}
	PaintBoxShadow(screen, x, y, w, h, r.BorderRadius, shadowColor, r.ShadowBlur, r.ShadowOffsetX, r.ShadowOffsetY)

	var bgColor color.Color
	if r.BackgroundColor != nil {
		bgColor = r.BackgroundColor
	}
	PaintBackground(screen, x, y, w, h, r.BorderRadius, bgColor)

	// 绘制背景图片（优先于背景色）
	if r.BackgroundImage != nil {
		if r.BackgroundSlice.Top > 0 || r.BackgroundSlice.Right > 0 || r.BackgroundSlice.Bottom > 0 || r.BackgroundSlice.Left > 0 {
			ninePatch := &RenderNinePatch{Source: r.BackgroundImage, Slice: r.BackgroundSlice}
			ninePatch.Paint(screen, Offset{X: x, Y: y})
		} else {
			op := &ebiten.DrawImageOptions{}
			sw := float64(w) / float64(r.BackgroundImage.Bounds().Dx())
			sh := float64(h) / float64(r.BackgroundImage.Bounds().Dy())
			op.GeoM.Scale(sw, sh)
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(r.BackgroundImage, op)
		}
	}

	// 绘制边框（子元素 focus 时使用 FocusedBorderColor）
	var borderColor color.Color
	if r.BorderColor != nil {
		borderColor = r.BorderColor
	}
	if r.hasFocusedChild() && r.FocusedBorderColor != nil {
		borderColor = r.FocusedBorderColor
	}
	PaintBorder(screen, x, y, w, h, r.BorderRadius, r.BorderWidth, borderColor)
}

// hasFocusedChild 检查直接子元素是否有处于 focus 状态的。
func (r *RenderBox) hasFocusedChild() bool {
	for _, child := range r.children {
		if f, ok := child.(interface{ IsFocused() bool }); ok {
			if f.IsFocused() {
				return true
			}
		}
	}
	return false
}

// SetBackgroundColor 设置背景色。
func (r *RenderBox) SetBackgroundColor(c *Color) {
	if r.BackgroundColor != nil && r.BackgroundColor.Equals(c) {
		return
	}
	r.BackgroundColor = c
	r.MarkNeedsPaint()
}

func (r *RenderBox) SetBackgroundImage(img *ebiten.Image, slice BorderSlice) {
	if r.BackgroundImage == img && r.BackgroundSlice == slice {
		return
	}
	r.BackgroundImage = img
	r.BackgroundSlice = slice
	r.MarkNeedsPaint()
}

func (r *RenderBox) SetBorderRadius(v float32) {
	br := UniformBorderRadius(v)
	if r.BorderRadius == br {
		return
	}
	r.BorderRadius = br
	r.MarkNeedsPaint()
}

func (r *RenderBox) SetBorderRadiusDetailed(br BorderRadius) {
	if r.BorderRadius == br {
		return
	}
	r.BorderRadius = br
	r.MarkNeedsPaint()
}

func (r *RenderBox) SetBorderColor(c *Color) {
	if r.BorderColor != nil && r.BorderColor.Equals(c) {
		return
	}
	r.BorderColor = c
	r.MarkNeedsPaint()
}

func (r *RenderBox) SetBorderWidth(v float32) {
	if r.BorderWidth == v {
		return
	}
	r.BorderWidth = v
	r.MarkNeedsPaint()
}

func (r *RenderBox) SetFocusedBorderColor(c *Color) {
	if r.FocusedBorderColor != nil && r.FocusedBorderColor.Equals(c) {
		return
	}
	r.FocusedBorderColor = c
	r.MarkNeedsPaint()
}

func (r *RenderBox) SetShadow(color *Color, blur, offsetX, offsetY float32) {
	if ColorPtrEquals(r.ShadowColor, color) && r.ShadowBlur == blur &&
		r.ShadowOffsetX == offsetX && r.ShadowOffsetY == offsetY {
		return
	}
	r.ShadowColor = color
	r.ShadowBlur = blur
	r.ShadowOffsetX = offsetX
	r.ShadowOffsetY = offsetY
	r.MarkNeedsPaint()
}

func (r *RenderBox) SetClipChildren(v bool) {
	if r.clipChildren == v {
		return
	}
	r.clipChildren = v
	r.MarkNeedsPaint()
}

func (r *RenderBox) ClipChildren() bool {
	return r.clipChildren
}

// 布局属性代理到 Yoga 节点
func (r *RenderBox) StyleSetWidth(v float32)            { r.yoga.StyleSetWidth(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetWidthPercent(v float32)     { r.yoga.StyleSetWidthPercent(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetWidthAuto()                 { r.yoga.StyleSetWidthAuto(); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetHeight(v float32)           { r.yoga.StyleSetHeight(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetHeightPercent(v float32)    { r.yoga.StyleSetHeightPercent(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetHeightAuto()                { r.yoga.StyleSetHeightAuto(); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetMinWidth(v float32)         { r.yoga.StyleSetMinWidth(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetMinHeight(v float32)        { r.yoga.StyleSetMinHeight(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetMaxWidth(v float32)         { r.yoga.StyleSetMaxWidth(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetMaxHeight(v float32)        { r.yoga.StyleSetMaxHeight(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetFlexDirection(v yoga.FlexDirection) { r.yoga.StyleSetFlexDirection(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetJustifyContent(v yoga.Justify)     { r.yoga.StyleSetJustifyContent(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetAlignItems(v yoga.Align)           { r.yoga.StyleSetAlignItems(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetAlignSelf(v yoga.Align)            { r.yoga.StyleSetAlignSelf(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetFlexGrow(v float32)                { r.yoga.StyleSetFlexGrow(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetFlexShrink(v float32)              { r.yoga.StyleSetFlexShrink(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetFlexBasis(v float32)               { r.yoga.StyleSetFlexBasis(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetFlexBasisAuto()                    { r.yoga.StyleSetFlexBasisAuto(); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetFlexWrap(v yoga.Wrap)              { r.yoga.StyleSetFlexWrap(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetPadding(edge yoga.Edge, v float32) { r.yoga.StyleSetPadding(edge, v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetMargin(edge yoga.Edge, v float32)      { r.yoga.StyleSetMargin(edge, v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetMarginAuto(edge yoga.Edge)              { r.yoga.StyleSetMarginAuto(edge); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetBorder(edge yoga.Edge, v float32)  { r.yoga.StyleSetBorder(edge, v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetPosition(edge yoga.Edge, v float32){ r.yoga.StyleSetPosition(edge, v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetPositionType(v yoga.PositionType)  { r.yoga.StyleSetPositionType(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetAspectRatio(v float32)             { r.yoga.StyleSetAspectRatio(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetDisplay(v yoga.Display)             { r.yoga.StyleSetDisplay(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetOverflow(v yoga.Overflow)           { r.yoga.StyleSetOverflow(v); r.MarkNeedsLayout() }
func (r *RenderBox) StyleSetGap(gutter yoga.Gutter, v float32)  { r.yoga.StyleSetGap(gutter, v); r.MarkNeedsLayout() }

// Color 是可空的 RGBA 颜色类型，实现 color.Color 接口。
type Color struct {
	R, G, B, A uint8
}

func NewColor(r, g, b, a uint8) *Color {
	return &Color{R: r, G: g, B: b, A: a}
}

func NewColorFrom(c color.Color) *Color {
	r, g, b, a := c.RGBA()
	return &Color{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

func (c *Color) RGBA() (r, g, b, a uint32) {
	if c == nil {
		return 0, 0, 0, 0
	}
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(c.A)
	a |= a << 8
	return
}

func (c *Color) Equals(other *Color) bool {
	if c == nil || other == nil {
		return c == other
	}
	return c.R == other.R && c.G == other.G && c.B == other.B && c.A == other.A
}

// ColorPtrEquals 比较两个 *Color 指针指向的颜色是否相等（处理 nil）。
func ColorPtrEquals(a, b *Color) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equals(b)
}
