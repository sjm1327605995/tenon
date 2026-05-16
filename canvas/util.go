package canvas

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/geometry"
)

// roundedRectPath creates a vector path for a rounded rectangle.
func roundedRectPath(r geometry.Rect, radius float32) vector.Path {
	path := vector.Path{}
	x, y, w, h := r.Min.X, r.Min.Y, r.Width(), r.Height()
	if radius < 0 {
		radius = 0
	}
	if radius > w/2 {
		radius = w / 2
	}
	if radius > h/2 {
		radius = h / 2
	}

	path.MoveTo(x+radius, y)
	path.LineTo(x+w-radius, y)
	path.Arc(x+w-radius, y+radius, radius, -math.Pi/2, 0, vector.Clockwise)
	path.LineTo(x+w, y+h-radius)
	path.Arc(x+w-radius, y+h-radius, radius, 0, math.Pi/2, vector.Clockwise)
	path.LineTo(x+radius, y+h)
	path.Arc(x+radius, y+h-radius, radius, math.Pi/2, math.Pi, vector.Clockwise)
	path.LineTo(x, y+radius)
	path.Arc(x+radius, y+radius, radius, math.Pi, math.Pi*1.5, vector.Clockwise)
	path.Close()
	return path
}

func cosf32(a float64) float32 { return float32(math.Cos(a)) }
func sinf32(a float64) float32 { return float32(math.Sin(a)) }
