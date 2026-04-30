package native

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/internal/core"
)

func HasRadius(r core.BorderRadius) bool {
	return r.TopLeft > 0 || r.TopRight > 0 || r.BottomRight > 0 || r.BottomLeft > 0
}

func DrawRoundedRectFill(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, clr color.Color) {
	var path vector.Path
	BuildRoundedRectPath(&path, x, y, w, h, r)
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.FillPath(screen, &path, &vector.FillOptions{}, op)
}

func DrawRoundedRectStroke(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, stroke float32, clr color.Color) {
	var path vector.Path
	BuildRoundedRectPath(&path, x, y, w, h, r)
	strokeOp := &vector.StrokeOptions{Width: stroke, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.StrokePath(screen, &path, strokeOp, op)
}

func DrawFilledCirclePath(screen *ebiten.Image, cx, cy, radius float32, clr color.Color) {
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

func StrokeCirclePath(screen *ebiten.Image, cx, cy, radius, strokeWidth float32, clr color.Color) {
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

func BuildRoundedRectPath(path *vector.Path, x, y, w, h float32, r core.BorderRadius) {
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
