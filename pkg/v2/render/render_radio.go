package render

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/yoga"
)

// RenderRadio is a single-select radio button.
type RenderRadio struct {
	RenderBox
	Selected    bool
	BoxSize     float32
	BorderColor color.Color
	FillColor   color.Color
	InnerColor  color.Color
	OnChange    func(selected bool)
}

func NewRenderRadio() *RenderRadio {
	r := &RenderRadio{BoxSize: 18}
	r.RenderBox.Init(r)
	r.yoga = yoga.NewNode()
	r.yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
	r.yoga.StyleSetAlignItems(yoga.AlignCenter)
	// 为 label 留出空间，避免 radio circle 与文字重叠
	r.StyleSetPadding(yoga.EdgeLeft, r.BoxSize+8)
	r.StyleSetMinHeight(r.BoxSize)
	return r
}

func (r *RenderRadio) SetSelected(v bool) {
	if r.Selected != v {
		r.Selected = v
		r.MarkNeedsPaint()
	}
}

func (r *RenderRadio) Paint(screen *ebiten.Image, offset Offset) {
	bounds := r.bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	centerY := offset.Y + bounds.Y + bounds.Height/2
	centerX := offset.X + bounds.X + r.BoxSize/2

	if r.Selected {
		r.drawFilledCircle(screen, centerX, centerY, r.BoxSize/2, r.FillColor)
		r.drawFilledCircle(screen, centerX, centerY, r.BoxSize/4, r.InnerColor)
	} else {
		r.drawFilledCircle(screen, centerX, centerY, r.BoxSize/2, color.White)
		r.drawStrokedCircle(screen, centerX, centerY, r.BoxSize/2, 1.5, r.BorderColor)
	}

	// 调用父类 Paint 绘制子节点（label）
	r.RenderBox.Paint(screen, offset)
}

func (r *RenderRadio) drawFilledCircle(screen *ebiten.Image, cx, cy, radius float32, clr color.Color) {
	if clr == nil {
		return
	}
	path := &vector.Path{}
	path.MoveTo(float32(cx+radius), float32(cy))
	path.Arc(float32(cx), float32(cy), float32(radius), 0, float32(math.Pi*2), vector.Clockwise)
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.FillPath(screen, path, &vector.FillOptions{}, op)
}

func (r *RenderRadio) drawStrokedCircle(screen *ebiten.Image, cx, cy, radius, strokeWidth float32, clr color.Color) {
	if clr == nil {
		return
	}
	path := &vector.Path{}
	path.MoveTo(float32(cx+radius), float32(cy))
	path.Arc(float32(cx), float32(cy), float32(radius), 0, float32(math.Pi*2), vector.Clockwise)
	strokeOp := &vector.StrokeOptions{Width: strokeWidth, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.StrokePath(screen, path, strokeOp, op)
}
