package components

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

// LoadingSpinner is a rotating loading indicator.
type LoadingSpinner struct {
	core.BaseElement
	angle  float32
	clr    color.Color
	size   float32
	stroke float32
}

// NewLoadingSpinner creates a loading spinner.
func NewLoadingSpinner() *LoadingSpinner {
	ls := &LoadingSpinner{
		clr:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
		size:   14,
		stroke: 2,
	}
	ls.Init(ls)
	ls.SetWidth(ls.size)
	ls.SetHeight(ls.size)
	return ls
}

// ElementType returns type identifier.
func (ls *LoadingSpinner) ElementType() string { return "LoadingSpinner" }

// SetColor sets the spinner color.
func (ls *LoadingSpinner) SetColor(clr color.Color) *LoadingSpinner {
	ls.clr = clr
	ls.Mark(core.FlagNeedDraw)
	return ls
}

// SetSize sets the spinner size.
func (ls *LoadingSpinner) SetSize(size float32) *LoadingSpinner {
	ls.size = size
	ls.SetWidth(size)
	ls.SetHeight(size)
	ls.Mark(core.FlagNeedDraw)
	return ls
}

// SetStroke sets the stroke width.
func (ls *LoadingSpinner) SetStroke(stroke float32) *LoadingSpinner {
	ls.stroke = stroke
	ls.Mark(core.FlagNeedDraw)
	return ls
}

// Update rotates the spinner.
func (ls *LoadingSpinner) Update() error {
	ls.angle += 0.15
	if ls.angle > 2*math.Pi {
		ls.angle -= 2 * math.Pi
	}
	return nil
}

// Draw renders the rotating arc.
func (ls *LoadingSpinner) Draw(screen *ebiten.Image) {
	if !ls.IsVisible() {
		return
	}
	bounds := ls.GetBounds()
	cx := bounds.X + bounds.Width/2
	cy := bounds.Y + bounds.Height/2
	r := ls.size / 2

	var path vector.Path
	startAngle := ls.angle
	endAngle := ls.angle + 1.2
	path.MoveTo(cx+r*float32(math.Cos(float64(startAngle))), cy+r*float32(math.Sin(float64(startAngle))))
	path.Arc(cx, cy, r, startAngle, endAngle, vector.Clockwise)
	strokeOp := &vector.StrokeOptions{Width: ls.stroke, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(ls.clr)
	vector.StrokePath(screen, &path, strokeOp, op)
}
