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

func (r *RenderBox) Paint(screen *ebiten.Image, offset Offset) {
	bounds := r.bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	x := offset.X + bounds.X
	y := offset.Y + bounds.Y
	w := bounds.Width
	h := bounds.Height

	// 有变换时：先画到临时 Image，再整体变换贴回
	if !r.transform.IsIdentity() {
		margin := r.ShadowBlur + max(r.ShadowOffsetX, r.ShadowOffsetY)
		tmpW := int(w + margin*2)
		tmpH := int(h + margin*2)
		if tmpW <= 0 {
			tmpW = 1
		}
		if tmpH <= 0 {
			tmpH = 1
		}
		tmp := ebiten.NewImage(tmpW, tmpH)
		defer tmp.Dispose()

		localBounds := Bounds{X: margin, Y: margin, Width: w, Height: h}
		r.paintContent(tmp, localBounds)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Concat(BuildTransformGeoM(Bounds{X: x, Y: y, Width: w, Height: h}, r.transform))
		ApplyColorScaleAlpha(&op.ColorScale, r.transform.Alpha)
		screen.DrawImage(tmp, op)
		return
	}

	// 无变换时：直接绘制到 screen
	r.paintContent(screen, Bounds{X: x, Y: y, Width: w, Height: h})
}

// paintContent 绘制背景、阴影、边框等内容（不含变换和裁剪）。
func (r *RenderBox) paintContent(screen *ebiten.Image, b Bounds) {
	x, y, w, h := b.X, b.Y, b.Width, b.Height

	// 绘制阴影
	if r.ShadowColor != nil {
		DrawShadow(screen, x, y, w, h, r.BorderRadius, r.ShadowColor, r.ShadowBlur, r.ShadowOffsetX, r.ShadowOffsetY)
	}

	// 绘制背景
	if r.BackgroundColor != nil {
		if !r.BorderRadius.IsZero() {
			DrawRoundedRectFill(screen, x, y, w, h, r.BorderRadius, r.BackgroundColor)
		} else {
			DrawRect(screen, x, y, w, h, r.BackgroundColor)
		}
	}

	// 绘制边框（子元素 focus 时使用 FocusedBorderColor）
	borderColor := r.BorderColor
	if r.hasFocusedChild() && r.FocusedBorderColor != nil {
		borderColor = r.FocusedBorderColor
	}
	if borderColor != nil && r.BorderWidth > 0 {
		if !r.BorderRadius.IsZero() {
			DrawRoundedRectStroke(screen, x, y, w, h, r.BorderRadius, r.BorderWidth, borderColor)
		} else {
			DrawBorder(screen, x, y, w, h, r.BorderWidth, borderColor)
		}
	}
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
func (r *RenderBox) StyleSetMargin(edge yoga.Edge, v float32)  { r.yoga.StyleSetMargin(edge, v); r.MarkNeedsLayout() }
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
