package components

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

func hasRadius(r core.BorderRadius) bool {
	return r.TopLeft > 0 || r.TopRight > 0 || r.BottomRight > 0 || r.BottomLeft > 0
}

func drawRoundedRectFill(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, clr color.Color) {
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, r)
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.FillPath(screen, &path, &vector.FillOptions{}, op)
}

func drawRoundedRectStroke(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, stroke float32, clr color.Color) {
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, r)
	strokeOp := &vector.StrokeOptions{Width: stroke, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.StrokePath(screen, &path, strokeOp, op)
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// colorsEqual 比较两个 color.Color 是否相等。
func colorsEqual(a, b color.Color) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	ra, ga, ba, aa := a.RGBA()
	rb, gb, bb, ab := b.RGBA()
	return ra == rb && ga == gb && ba == bb && aa == ab
}

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

// buildRoundedRectPath 构建圆角矩形路径。
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
