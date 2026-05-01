package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/yoga"
)

// RenderSlider is a value slider.
type RenderSlider struct {
	RenderBox
	MinValue    float32
	MaxValue    float32
	Value       float32
	TrackColor  color.Color
	FillColor   color.Color
	ThumbColor  color.Color
	ThumbRadius float32
	TrackHeight float32
	OnChange    func(value float32)
}

func NewRenderSlider() *RenderSlider {
	r := &RenderSlider{
		MinValue:    0,
		MaxValue:    100,
		Value:       0,
		ThumbRadius: 8,
		TrackHeight: 4,
	}
	r.RenderBox.Init(r)
	r.yoga = yoga.NewNode()
	r.yoga.StyleSetHeight(24)
	r.StyleSetMinWidth(200)
	return r
}

func (r *RenderSlider) SetValue(v float32) {
	if v < r.MinValue {
		v = r.MinValue
	}
	if v > r.MaxValue {
		v = r.MaxValue
	}
	if r.Value != v {
		r.Value = v
		r.MarkNeedsPaint()
	}
}

func (r *RenderSlider) Paint(screen *ebiten.Image, offset Offset) {
	bounds := r.bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	trackY := offset.Y + bounds.Y + bounds.Height/2
	trackTop := trackY - r.TrackHeight/2

	// Track background
	if r.TrackColor != nil {
		vector.FillRect(screen, offset.X+bounds.X, trackTop, bounds.Width, r.TrackHeight, r.TrackColor, false)
	}

	// Fill
	if r.FillColor != nil && r.MaxValue > r.MinValue {
		ratio := (r.Value - r.MinValue) / (r.MaxValue - r.MinValue)
		fillW := bounds.Width * ratio
		if fillW > 0 {
			vector.FillRect(screen, offset.X+bounds.X, trackTop, fillW, r.TrackHeight, r.FillColor, false)
		}
	}

	// Thumb
	thumbX := offset.X + bounds.X
	if r.MaxValue > r.MinValue {
		ratio := (r.Value - r.MinValue) / (r.MaxValue - r.MinValue)
		thumbX = offset.X + bounds.X + bounds.Width*ratio
	}
	if r.ThumbColor != nil {
		r.drawFilledCircle(screen, thumbX, trackY, r.ThumbRadius, r.ThumbColor)
		r.drawStrokedCircle(screen, thumbX, trackY, r.ThumbRadius, 1.5, color.RGBA{R: 229, G: 229, B: 229, A: 255})
	}
}

func (r *RenderSlider) getAbsoluteX() float32 {
	x := r.bounds.X
	parent := r.GetParent()
	for parent != nil {
		x += parent.GetBounds().X
		parent = parent.GetParent()
	}
	return x
}

func (r *RenderSlider) HandleDrag(mx, my float32) {
	absX := r.getAbsoluteX()
	relX := mx - absX
	if r.bounds.Width > 0 {
		ratio := relX / r.bounds.Width
		if ratio < 0 {
			ratio = 0
		}
		if ratio > 1 {
			ratio = 1
		}
		newValue := r.MinValue + (r.MaxValue-r.MinValue)*ratio
		if newValue != r.Value {
			r.Value = newValue
			r.MarkNeedsPaint()
			if r.OnChange != nil {
				r.OnChange(newValue)
			}
		}
	}
}

func (r *RenderSlider) drawFilledCircle(screen *ebiten.Image, cx, cy, radius float32, clr color.Color) {
	if clr == nil {
		return
	}
	path := &vector.Path{}
	path.MoveTo(cx+radius, cy)
	path.Arc(cx, cy, radius, 0, 3.14159265*2, vector.Clockwise)
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.FillPath(screen, path, &vector.FillOptions{}, op)
}

func (r *RenderSlider) drawStrokedCircle(screen *ebiten.Image, cx, cy, radius, strokeWidth float32, clr color.Color) {
	if clr == nil {
		return
	}
	path := &vector.Path{}
	path.MoveTo(cx+radius, cy)
	path.Arc(cx, cy, radius, 0, 3.14159265*2, vector.Clockwise)
	strokeOp := &vector.StrokeOptions{Width: strokeWidth, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.StrokePath(screen, path, strokeOp, op)
}
