package components

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
)

// LoadingSpinner 是一个旋转的 loading 指示器。
type LoadingSpinner struct {
	core.BaseHost
	angle float32
	clr   color.Color
	size  float32
}

// NewLoadingSpinner 创建一个 loading 旋转器。
func NewLoadingSpinner() *LoadingSpinner {
	ls := &LoadingSpinner{
		clr:  color.RGBA{R: 255, G: 255, B: 255, A: 255},
		size: 14,
	}
	ls.Init(ls)
	ls.GetElement().Yoga.StyleSetWidth(ls.size)
	ls.GetElement().Yoga.StyleSetHeight(ls.size)
	return ls
}

// SetColor 设置旋转器颜色。
func (ls *LoadingSpinner) SetColor(clr color.Color) *LoadingSpinner {
	ls.clr = clr
	return ls
}

// SetSize 设置旋转器大小。
func (ls *LoadingSpinner) SetSize(size float32) *LoadingSpinner {
	ls.size = size
	ls.GetElement().Yoga.StyleSetWidth(size)
	ls.GetElement().Yoga.StyleSetHeight(size)
	return ls
}

// Update 更新旋转角度。
func (ls *LoadingSpinner) Update() error {
	ls.angle += 0.15
	if ls.angle > 2*math.Pi {
		ls.angle -= 2 * math.Pi
	}
	return nil
}

// Draw 绘制旋转圆弧。
func (ls *LoadingSpinner) Draw(screen *ebiten.Image) {
	bounds := ls.GetLayoutBounds()
	cx := bounds.X + bounds.Width/2
	cy := bounds.Y + bounds.Height/2
	r := ls.size / 2
	stroke := float32(2)

	var path vector.Path
	startAngle := ls.angle
	endAngle := ls.angle + 1.2
	path.MoveTo(cx+r*float32(math.Cos(float64(startAngle))), cy+r*float32(math.Sin(float64(startAngle))))
	path.Arc(cx, cy, r, startAngle, endAngle, vector.Clockwise)
	strokeOp := &vector.StrokeOptions{Width: stroke, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(ls.clr)
	vector.StrokePath(screen, &path, strokeOp, op)
}
