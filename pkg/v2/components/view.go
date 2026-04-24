package components

import (
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

// Draw renders background, shadow and border.
func (v *View) Draw(screen *ebiten.Image) {
	if !v.IsVisible() {
		return
	}
	bounds := v.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

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
