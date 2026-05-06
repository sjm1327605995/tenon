package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/yoga"
)

// RenderProgressBar is a linear progress indicator.
type RenderProgressBar struct {
	RenderBox
	Progress     float32 // 0.0 ~ 1.0
	TrackColor   color.Color
	FillColor    color.Color
	BarRadius    float32
}

func NewRenderProgressBar() *RenderProgressBar {
	r := &RenderProgressBar{}
	r.RenderBox.Init(r)
	r.yoga = yoga.NewNode()
	r.BarRadius = 4
	r.StyleSetHeight(8)
	r.StyleSetWidthPercent(100)
	return r
}

func (r *RenderProgressBar) SetProgress(v float32) {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	if r.Progress != v {
		r.Progress = v
		r.MarkNeedsPaint()
	}
}

func (r *RenderProgressBar) Paint(screen *ebiten.Image, offset Offset) {
	bounds := r.bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}
	x := offset.X + bounds.X
	y := offset.Y + bounds.Y
	w := bounds.Width
	h := bounds.Height

	br := UniformBorderRadius(r.BarRadius)
	br = clampRadius(br, w, h)

	// Track
	DrawRoundedRectFill(screen, x, y, w, h, br, r.TrackColor)

	// Fill
	if r.FillColor != nil && r.Progress > 0 {
		fillW := w * r.Progress
		if fillW > w {
			fillW = w
		}
		DrawRoundedRectFill(screen, x, y, fillW, h, br, r.FillColor)
	}
}
