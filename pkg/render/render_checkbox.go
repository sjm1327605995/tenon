package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/yoga"
)

// RenderCheckbox is a checkable box with optional label.
type RenderCheckbox struct {
	RenderBox
	Checked     bool
	BoxSize     float32
	BorderColor color.Color
	FillColor   color.Color
	CheckColor  color.Color
	OnChange    func(checked bool)
}

func NewRenderCheckbox() *RenderCheckbox {
	r := &RenderCheckbox{BoxSize: 18}
	r.RenderBox.Init(r)
	r.yoga = yoga.NewNode()
	r.StyleSetFlexDirection(yoga.FlexDirectionRow)
	r.StyleSetAlignItems(yoga.AlignCenter)
	// 为 label 留出空间，避免 checkbox 与文字重叠
	r.StyleSetPadding(yoga.EdgeLeft, r.BoxSize+8)
	r.StyleSetMinHeight(r.BoxSize)
	return r
}

func (r *RenderCheckbox) SetChecked(v bool) {
	if r.Checked != v {
		r.Checked = v
		r.MarkNeedsPaint()
	}
}

func (r *RenderCheckbox) Paint(screen *ebiten.Image, offset Offset) {
	bounds := r.bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	boxY := offset.Y + bounds.Y + (bounds.Height-r.BoxSize)/2
	boxX := offset.X + bounds.X

	if r.Checked {
		DrawRoundedRectFill(screen, boxX, boxY, r.BoxSize, r.BoxSize, UniformBorderRadius(3), r.FillColor)
		r.drawCheckmark(screen, boxX, boxY, r.BoxSize)
	} else {
		DrawRoundedRectFill(screen, boxX, boxY, r.BoxSize, r.BoxSize, UniformBorderRadius(3), color.White)
		DrawRoundedRectStroke(screen, boxX, boxY, r.BoxSize, r.BoxSize, UniformBorderRadius(3), 1.5, r.BorderColor)
	}

	// 调用父类 Paint 绘制子节点（label）
	r.RenderBox.Paint(screen, offset)
}

func (r *RenderCheckbox) drawCheckmark(screen *ebiten.Image, x, y, size float32) {
	padding := size * 0.2
	x1 := x + padding
	y1 := y + size/2
	x2 := x + size*0.4
	y2 := y + size - padding
	x3 := x + size - padding
	y3 := y + padding

	path := &vector.Path{}
	path.MoveTo(float32(x1), float32(y1))
	path.LineTo(float32(x2), float32(y2))
	path.LineTo(float32(x3), float32(y3))

	strokeOp := &vector.StrokeOptions{Width: 2, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(r.CheckColor)
	op.AntiAlias = true
	vector.StrokePath(screen, path, strokeOp, op)
}
