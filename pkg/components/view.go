package components

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// View 是最基础的容器宿主组件。
type View struct {
	core.BaseHost
}

// NewView 创建一个新的 View。
func NewView() *View {
	v := &View{}
	v.Init(v)
	return v
}

// Draw 绘制 View 的背景、阴影和边框。
func (v *View) Draw(screen *ebiten.Image) {
	el := v.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := v.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	v.drawShadow(screen, el, bounds)

	if el.BackgroundColor != nil {
		v.drawBackground(screen, el, bounds)
	}

	if el.BorderColor != nil {
		v.drawBorder(screen, el, bounds)
	}
}

func (v *View) drawShadow(screen *ebiten.Image, el *core.Element, b core.LayoutBounds) {
	if el.ShadowColor == nil || el.ShadowBlur <= 0 {
		return
	}
	sw := b.Width + el.ShadowBlur*2
	sh := b.Height + el.ShadowBlur*2
	sx := b.X - el.ShadowBlur + el.ShadowOffsetX
	sy := b.Y - el.ShadowBlur + el.ShadowOffsetY
	vector.FillRect(screen, sx, sy, sw, sh, el.ShadowColor, false)
}

func (v *View) drawBackground(screen *ebiten.Image, el *core.Element, b core.LayoutBounds) {
	if hasRadius(el.BorderRadius) {
		v.drawRoundedRectFill(screen, b.X, b.Y, b.Width, b.Height, el.BorderRadius, el.BackgroundColor)
	} else {
		vector.FillRect(screen, b.X, b.Y, b.Width, b.Height, el.BackgroundColor, false)
	}
}

func (v *View) drawBorder(screen *ebiten.Image, el *core.Element, b core.LayoutBounds) {
	yogaNode := el.Yoga
	if yogaNode == nil {
		return
	}
	borderTop := yogaNode.StyleGetBorder(yoga.EdgeTop)
	borderRight := yogaNode.StyleGetBorder(yoga.EdgeRight)
	borderBottom := yogaNode.StyleGetBorder(yoga.EdgeBottom)
	borderLeft := yogaNode.StyleGetBorder(yoga.EdgeLeft)

	if hasRadius(el.BorderRadius) {
		maxBorder := max(borderTop, max(borderRight, max(borderBottom, borderLeft)))
		v.drawRoundedRectStroke(screen, b.X, b.Y, b.Width, b.Height, el.BorderRadius, maxBorder, el.BorderColor)
		return
	}

	if borderTop > 0 {
		vector.FillRect(screen, b.X, b.Y, b.Width, borderTop, el.BorderColor, false)
	}
	if borderRight > 0 {
		vector.FillRect(screen, b.X+b.Width-borderRight, b.Y, borderRight, b.Height, el.BorderColor, false)
	}
	if borderBottom > 0 {
		vector.FillRect(screen, b.X, b.Y+b.Height-borderBottom, b.Width, borderBottom, el.BorderColor, false)
	}
	if borderLeft > 0 {
		vector.FillRect(screen, b.X, b.Y, borderLeft, b.Height, el.BorderColor, false)
	}
}

func (v *View) drawRoundedRectFill(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, clr color.Color) {
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, r)

	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.FillPath(screen, &path, &vector.FillOptions{}, op)
}

func (v *View) drawRoundedRectStroke(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, stroke float32, clr color.Color) {
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, r)

	strokeOp := &vector.StrokeOptions{Width: stroke, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.StrokePath(screen, &path, strokeOp, op)
}

func buildRoundedRectPath(path *vector.Path, x, y, w, h float32, r core.BorderRadius) {
	path.MoveTo(x+r.TopLeft, y)
	path.LineTo(x+w-r.TopRight, y)
	if r.TopRight > 0 {
		path.Arc(x+w-r.TopRight, y+r.TopRight, r.TopRight, -math.Pi/2, 0, vector.Clockwise)
	} else {
		path.LineTo(x+w, y)
	}
	path.LineTo(x+w, y+h-r.BottomRight)
	if r.BottomRight > 0 {
		path.Arc(x+w-r.BottomRight, y+h-r.BottomRight, r.BottomRight, 0, math.Pi/2, vector.Clockwise)
	} else {
		path.LineTo(x+w, y+h)
	}
	path.LineTo(x+r.BottomLeft, y+h)
	if r.BottomLeft > 0 {
		path.Arc(x+r.BottomLeft, y+h-r.BottomLeft, r.BottomLeft, math.Pi/2, math.Pi, vector.Clockwise)
	} else {
		path.LineTo(x, y+h)
	}
	path.LineTo(x, y+r.TopLeft)
	if r.TopLeft > 0 {
		path.Arc(x+r.TopLeft, y+r.TopLeft, r.TopLeft, math.Pi, 3*math.Pi/2, vector.Clockwise)
	} else {
		path.LineTo(x, y)
	}
	path.Close()
}

func hasRadius(r core.BorderRadius) bool {
	return r.TopLeft > 0 || r.TopRight > 0 || r.BottomRight > 0 || r.BottomLeft > 0
}

// ==================== 链式 API ====================

func (v *View) SetWidth(width float32) *View {
	v.GetElement().Yoga.StyleSetWidth(width)
	return v
}
func (v *View) SetHeight(height float32) *View {
	v.GetElement().Yoga.StyleSetHeight(height)
	return v
}
func (v *View) SetMinWidth(width float32) *View {
	v.GetElement().Yoga.StyleSetMinWidth(width)
	return v
}
func (v *View) SetMinHeight(height float32) *View {
	v.GetElement().Yoga.StyleSetMinHeight(height)
	return v
}
func (v *View) SetMaxWidth(width float32) *View {
	v.GetElement().Yoga.StyleSetMaxWidth(width)
	return v
}
func (v *View) SetMaxHeight(height float32) *View {
	v.GetElement().Yoga.StyleSetMaxHeight(height)
	return v
}
func (v *View) SetFlexDirection(dir yoga.FlexDirection) *View {
	v.GetElement().Yoga.StyleSetFlexDirection(dir)
	return v
}
func (v *View) SetJustifyContent(justify yoga.Justify) *View {
	v.GetElement().Yoga.StyleSetJustifyContent(justify)
	return v
}
func (v *View) SetAlignItems(align yoga.Align) *View {
	v.GetElement().Yoga.StyleSetAlignItems(align)
	return v
}
func (v *View) SetAlignSelf(align yoga.Align) *View {
	v.GetElement().Yoga.StyleSetAlignSelf(align)
	return v
}
func (v *View) SetFlexGrow(grow float32) *View {
	v.GetElement().Yoga.StyleSetFlexGrow(grow)
	return v
}
func (v *View) SetFlexShrink(shrink float32) *View {
	v.GetElement().Yoga.StyleSetFlexShrink(shrink)
	return v
}
func (v *View) SetFlexBasis(basis float32) *View {
	v.GetElement().Yoga.StyleSetFlexBasis(basis)
	return v
}
func (v *View) SetFlexWrap(wrap yoga.Wrap) *View {
	v.GetElement().Yoga.StyleSetFlexWrap(wrap)
	return v
}
func (v *View) SetPadding(edge yoga.Edge, value float32) *View {
	v.GetElement().Yoga.StyleSetPadding(edge, value)
	return v
}
func (v *View) SetMargin(edge yoga.Edge, value float32) *View {
	v.GetElement().Yoga.StyleSetMargin(edge, value)
	return v
}
func (v *View) SetBorder(edge yoga.Edge, value float32) *View {
	v.GetElement().Yoga.StyleSetBorder(edge, value)
	return v
}
func (v *View) SetPosition(edge yoga.Edge, value float32) *View {
	v.GetElement().Yoga.StyleSetPosition(edge, value)
	return v
}
func (v *View) SetPositionType(pt yoga.PositionType) *View {
	v.GetElement().Yoga.StyleSetPositionType(pt)
	return v
}
func (v *View) SetAspectRatio(ratio float32) *View {
	v.GetElement().Yoga.StyleSetAspectRatio(ratio)
	return v
}
func (v *View) SetDisplay(display yoga.Display) *View {
	v.GetElement().Yoga.StyleSetDisplay(display)
	return v
}
func (v *View) SetOverflow(overflow yoga.Overflow) *View {
	v.GetElement().Yoga.StyleSetOverflow(overflow)
	return v
}
func (v *View) SetBackgroundColor(clr color.Color) *View {
	v.GetElement().BackgroundColor = clr
	return v
}
func (v *View) SetBorderColor(clr color.Color) *View {
	v.GetElement().BorderColor = clr
	return v
}
func (v *View) SetBorderRadius(radius float32) *View {
	v.GetElement().BorderRadius = core.BorderRadius{
		TopLeft: radius, TopRight: radius,
		BottomRight: radius, BottomLeft: radius,
	}
	return v
}
func (v *View) SetBorderRadius4(tl, tr, br, bl float32) *View {
	v.GetElement().BorderRadius = core.BorderRadius{
		TopLeft: tl, TopRight: tr,
		BottomRight: br, BottomLeft: bl,
	}
	return v
}
func (v *View) SetShadow(clr color.Color, blur, offsetX, offsetY float32) *View {
	el := v.GetElement()
	el.ShadowColor = clr
	el.ShadowBlur = blur
	el.ShadowOffsetX = offsetX
	el.ShadowOffsetY = offsetY
	return v
}
func (v *View) SetVisible(visible bool) *View {
	v.GetElement().Visible = visible
	return v
}
func (v *View) SetPointerEvents(pe core.PointerEvents) *View {
	v.GetElement().PointerEvents = pe
	return v
}
func (v *View) SetKey(key string) *View {
	v.GetElement().Key = key
	return v
}
func (v *View) Add(children ...core.Component) *View {
	for _, c := range children {
		v.AddChild(c)
	}
	return v
}
