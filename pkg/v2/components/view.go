package components

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// View is the basic container Element.
type View struct {
	core.BaseElement
	backgroundColor color.Color
	borderColor     color.Color
	shadowColor     color.Color
	shadowBlur      float32
	shadowOffsetX   float32
	shadowOffsetY   float32
	borderRadius    core.BorderRadius
}

// NewView creates a View.
func NewView() *View {
	v := &View{}
	v.Init(v)
	return v
}

// ElementType returns type identifier.
func (v *View) ElementType() string { return "View" }

// IsNative returns true as View is a native rendering component.
func (v *View) IsNative() bool { return true }

// SyncFrom 同步新 View 的属性到当前 Element（声明式重建）。
func (v *View) SyncFrom(src core.Element) {
	other, ok := src.(*View)
	if !ok {
		return
	}
	sb := &SyncBuilder{}
	syncColor(sb, &v.backgroundColor, other.backgroundColor)
	syncColor(sb, &v.borderColor, other.borderColor)
	if !colorsEqual(v.shadowColor, other.shadowColor) || v.shadowBlur != other.shadowBlur ||
		v.shadowOffsetX != other.shadowOffsetX || v.shadowOffsetY != other.shadowOffsetY {
		v.shadowColor = other.shadowColor
		v.shadowBlur = other.shadowBlur
		v.shadowOffsetX = other.shadowOffsetX
		v.shadowOffsetY = other.shadowOffsetY
		sb.NeedDraw = true
	}
	syncField(sb, &v.borderRadius, other.borderRadius)
	sb.MarkDraw(v)
}

// Draw renders background, shadow and border.
func (v *View) Draw(screen *ebiten.Image) {
	bounds := v.GetBounds()

	t := v.GetTransform()
	if !t.IsIdentity() {
		// 需要包含 shadow 的边距
		margin := v.shadowBlur + max(v.shadowOffsetX, v.shadowOffsetY)
		if margin < 0 {
			margin = 0
		}
		mw := int(bounds.Width + margin*2)
		mh := int(bounds.Height + margin*2)
		if mw <= 0 {
			mw = 1
		}
		if mh <= 0 {
			mh = 1
		}
		tmp := ebiten.NewImage(mw, mh)
		defer tmp.Dispose()
		localBounds := core.LayoutBounds{
			X:      margin,
			Y:      margin,
			Width:  bounds.Width,
			Height: bounds.Height,
		}
		v.drawShadow(tmp, localBounds)
		if v.backgroundColor != nil {
			v.drawBackground(tmp, localBounds)
		}
		if v.borderColor != nil {
			v.drawBorder(tmp, localBounds)
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Concat(core.BuildTransformGeoM(bounds, t))
		core.ApplyColorScaleAlpha(&op.ColorScale, t.Alpha)
		screen.DrawImage(tmp, op)
		return
	}

	v.drawShadow(screen, bounds)

	if v.backgroundColor != nil {
		v.drawBackground(screen, bounds)
	}

	if v.borderColor != nil {
		v.drawBorder(screen, bounds)
	}
}

func (v *View) drawShadow(screen *ebiten.Image, b core.LayoutBounds) {
	if v.shadowColor == nil || v.shadowBlur <= 0 {
		return
	}
	sw := b.Width + v.shadowBlur*2
	sh := b.Height + v.shadowBlur*2
	sx := b.X - v.shadowBlur + v.shadowOffsetX
	sy := b.Y - v.shadowBlur + v.shadowOffsetY
	vector.FillRect(screen, sx, sy, sw, sh, v.shadowColor, false)
}

func (v *View) drawBackground(screen *ebiten.Image, b core.LayoutBounds) {
	if hasRadius(v.borderRadius) {
		drawRoundedRectFill(screen, b.X, b.Y, b.Width, b.Height, v.borderRadius, v.backgroundColor)
	} else {
		vector.FillRect(screen, b.X, b.Y, b.Width, b.Height, v.backgroundColor, false)
	}
}

func (v *View) drawBorder(screen *ebiten.Image, b core.LayoutBounds) {
	if hasRadius(v.borderRadius) {
		bt := v.GetYoga().StyleGetBorder(yoga.EdgeTop)
		br := v.GetYoga().StyleGetBorder(yoga.EdgeRight)
		bb := v.GetYoga().StyleGetBorder(yoga.EdgeBottom)
		bl := v.GetYoga().StyleGetBorder(yoga.EdgeLeft)
		maxBorder := max(bt, max(br, max(bb, bl)))
		drawRoundedRectStroke(screen, b.X, b.Y, b.Width, b.Height, v.borderRadius, maxBorder, v.borderColor)
		return
	}

	bt := v.GetYoga().StyleGetBorder(yoga.EdgeTop)
	br := v.GetYoga().StyleGetBorder(yoga.EdgeRight)
	bb := v.GetYoga().StyleGetBorder(yoga.EdgeBottom)
	bl := v.GetYoga().StyleGetBorder(yoga.EdgeLeft)

	if bt > 0 {
		vector.FillRect(screen, b.X, b.Y, b.Width, bt, v.borderColor, false)
	}
	if br > 0 {
		vector.FillRect(screen, b.X+b.Width-br, b.Y, br, b.Height, v.borderColor, false)
	}
	if bb > 0 {
		vector.FillRect(screen, b.X, b.Y+b.Height-bb, b.Width, bb, v.borderColor, false)
	}
	if bl > 0 {
		vector.FillRect(screen, b.X, b.Y, bl, b.Height, v.borderColor, false)
	}
}

// Chain API

func (v *View) SetBackgroundColor(c color.Color) *View {
	v.backgroundColor = c
	v.Mark(core.FlagNeedDraw)
	return v
}

func (v *View) SetBorderColor(c color.Color) *View {
	v.borderColor = c
	v.Mark(core.FlagNeedDraw)
	return v
}

func (v *View) SetBorderRadius(radius float32) *View {
	v.borderRadius = core.BorderRadius{
		TopLeft: radius, TopRight: radius,
		BottomRight: radius, BottomLeft: radius,
	}
	v.Mark(core.FlagNeedDraw)
	return v
}

func (v *View) SetBorderRadius4(tl, tr, br, bl float32) *View {
	v.borderRadius = core.BorderRadius{TopLeft: tl, TopRight: tr, BottomRight: br, BottomLeft: bl}
	v.Mark(core.FlagNeedDraw)
	return v
}

func (v *View) SetShadow(c color.Color, blur, offsetX, offsetY float32) *View {
	v.shadowColor = c
	v.shadowBlur = blur
	v.shadowOffsetX = offsetX
	v.shadowOffsetY = offsetY
	v.Mark(core.FlagNeedDraw)
	return v
}

// Utility functions

func (v *View) DebugProps() map[string]interface{} {
	props := make(map[string]interface{})
	if v.backgroundColor != nil {
		props["backgroundColor"] = colorToCSS(v.backgroundColor)
	}
	if v.borderColor != nil {
		props["borderColor"] = colorToCSS(v.borderColor)
	}
	if v.borderRadius.TopLeft != 0 || v.borderRadius.TopRight != 0 || v.borderRadius.BottomRight != 0 || v.borderRadius.BottomLeft != 0 {
		props["borderRadius"] = map[string]float32{
			"topLeft":     v.borderRadius.TopLeft,
			"topRight":    v.borderRadius.TopRight,
			"bottomRight": v.borderRadius.BottomRight,
			"bottomLeft":  v.borderRadius.BottomLeft,
		}
	}
	if v.shadowColor != nil {
		props["shadow"] = map[string]interface{}{
			"color":   colorToCSS(v.shadowColor),
			"blur":    v.shadowBlur,
			"offsetX": v.shadowOffsetX,
			"offsetY": v.shadowOffsetY,
		}
	}
	return props
}

func colorToCSS(c color.Color) string {
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("rgba(%d,%d,%d,%.2f)", r>>8, g>>8, b>>8, float64(a>>8)/255.0)
}
