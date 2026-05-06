package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/yoga"
)

// RenderSwitch is a toggle switch.
type RenderSwitch struct {
	RenderBox
	Checked    bool
	TrackWidth float32
	TrackHeight float32
	OffColor   color.Color
	OnColor    color.Color
	ThumbColor color.Color
	OnChange   func(checked bool)
}

func NewRenderSwitch() *RenderSwitch {
	r := &RenderSwitch{TrackWidth: 44, TrackHeight: 24}
	r.RenderBox.Init(r)
	r.yoga = yoga.NewNode()
	r.yoga.StyleSetWidth(r.TrackWidth)
	r.yoga.StyleSetHeight(r.TrackHeight)
	return r
}

func (r *RenderSwitch) SetChecked(v bool) {
	if r.Checked != v {
		r.Checked = v
		r.MarkNeedsPaint()
	}
}

func (r *RenderSwitch) Paint(screen *ebiten.Image, offset Offset) {
	bounds := r.bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}
	x := offset.X + bounds.X
	y := offset.Y + bounds.Y
	w := bounds.Width
	h := bounds.Height

	trackColor := r.OffColor
	if r.Checked {
		trackColor = r.OnColor
	}

	// Track (rounded rect)
	radius := h / 2
	if trackColor != nil {
		DrawRoundedRectFill(screen, x, y, w, h, UniformBorderRadius(radius), trackColor)
	}

	// Thumb (circle)
	thumbRadius := h/2 - 2
	thumbY := y + h/2
	leftX := x + thumbRadius + 2
	rightX := x + w - thumbRadius - 2
	var thumbX float32
	if r.Checked {
		thumbX = rightX
	} else {
		thumbX = leftX
	}

	if r.ThumbColor != nil {
		DrawFilledCirclePath(screen, thumbX, thumbY, thumbRadius, r.ThumbColor)
	}
}
