package components

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawFilledCirclePath 使用 vector.Path 绘制抗锯齿填充圆形。
func drawFilledCirclePath(screen *ebiten.Image, cx, cy, radius float32, clr color.Color) {
	var path vector.Path
	path.MoveTo(cx+radius, cy)
	path.Arc(cx, cy, radius, 0, math.Pi, vector.Clockwise)
	path.Arc(cx, cy, radius, math.Pi, 2*math.Pi, vector.Clockwise)
	path.Close()

	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.FillPath(screen, &path, &vector.FillOptions{}, op)
}

// strokeCirclePath 使用 vector.Path 绘制抗锯齿描边圆形。
func strokeCirclePath(screen *ebiten.Image, cx, cy, radius, strokeWidth float32, clr color.Color) {
	var path vector.Path
	path.MoveTo(cx+radius, cy)
	path.Arc(cx, cy, radius, 0, math.Pi, vector.Clockwise)
	path.Arc(cx, cy, radius, math.Pi, 2*math.Pi, vector.Clockwise)
	path.Close()

	strokeOp := &vector.StrokeOptions{Width: strokeWidth, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.StrokePath(screen, &path, strokeOp, op)
}
