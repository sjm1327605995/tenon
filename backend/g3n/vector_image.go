package g3n

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/sjm1327605995/tenon/pkg/ui"
	"golang.org/x/image/vector"
)

type vectorImageKey struct {
	fill   bool
	stroke float32
	color  ui.Color
}

type vectorImageCache struct {
	sync.Mutex
	images map[vectorImageKey]*vectorImage
}

type vectorImage struct {
	bitmap *bitmap
	x, y   float32
}

func (p *vectorPath) image(fill bool, stroke float32, c ui.Color) *vectorImage {
	key := vectorImageKey{fill: fill, stroke: stroke, color: c}
	p.cache.Lock()
	defer p.cache.Unlock()
	if cached := p.cache.images[key]; cached != nil {
		return cached
	}
	created := rasterizeVector(p.svg, p.scale, fill, stroke, c)
	if created == nil {
		return nil
	}
	if p.cache.images == nil {
		p.cache.images = make(map[vectorImageKey]*vectorImage)
	}
	p.cache.images[key] = created
	return created
}

type vectorBounds struct {
	minX, minY, maxX, maxY float32
	valid                  bool
}

func (b *vectorBounds) add(x, y float32) {
	if !b.valid {
		b.minX, b.minY, b.maxX, b.maxY, b.valid = x, y, x, y, true
		return
	}
	b.minX, b.minY = min(b.minX, x), min(b.minY, y)
	b.maxX, b.maxY = max(b.maxX, x), max(b.maxY, y)
}

func (b *vectorBounds) MoveTo(x, y float32) { b.add(x, y) }
func (b *vectorBounds) LineTo(x, y float32) { b.add(x, y) }
func (b *vectorBounds) CubeTo(x1, y1, x2, y2, x, y float32) {
	b.add(x1, y1)
	b.add(x2, y2)
	b.add(x, y)
}
func (b *vectorBounds) QuadTo(x1, y1, x, y float32) {
	b.add(x1, y1)
	b.add(x, y)
}
func (*vectorBounds) ClosePath() {}

type offsetVectorSink struct {
	ui.SVGPathSink
	dx, dy float32
}

func (s offsetVectorSink) MoveTo(x, y float32) { s.SVGPathSink.MoveTo(x+s.dx, y+s.dy) }
func (s offsetVectorSink) LineTo(x, y float32) { s.SVGPathSink.LineTo(x+s.dx, y+s.dy) }
func (s offsetVectorSink) CubeTo(x1, y1, x2, y2, x, y float32) {
	s.SVGPathSink.CubeTo(x1+s.dx, y1+s.dy, x2+s.dx, y2+s.dy, x+s.dx, y+s.dy)
}
func (s offsetVectorSink) QuadTo(x1, y1, x, y float32) {
	s.SVGPathSink.QuadTo(x1+s.dx, y1+s.dy, x+s.dx, y+s.dy)
}
func (s offsetVectorSink) ClosePath() { s.SVGPathSink.ClosePath() }

func rasterizeVector(d string, scale float32, fill bool, stroke float32, c ui.Color) *vectorImage {
	bounds := &vectorBounds{}
	ui.ParseSVGPath(d, scale, bounds)
	if !bounds.valid {
		return nil
	}
	padding := float32(1)
	if !fill {
		padding += stroke / 2
	}
	minX := float32(math.Floor(float64(bounds.minX - padding)))
	minY := float32(math.Floor(float64(bounds.minY - padding)))
	maxX := float32(math.Ceil(float64(bounds.maxX + padding)))
	maxY := float32(math.Ceil(float64(bounds.maxY + padding)))
	width, height := int(maxX-minX), int(maxY-minY)
	if width <= 0 || height <= 0 {
		return nil
	}
	raster := vector.NewRasterizer(width, height)
	sink := offsetVectorSink{dx: -minX, dy: -minY}
	if fill {
		sink.SVGPathSink = raster
		ui.ParseSVGPath(d, scale, sink)
	} else {
		sink.SVGPathSink = &vectorStrokeSink{raster: raster, radius: stroke / 2}
		ui.ParseSVGPath(d, scale, sink)
	}
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	raster.Draw(img, img.Bounds(), image.NewUniform(color.NRGBA{R: c.R, G: c.G, B: c.B, A: c.A}), image.Point{})
	return &vectorImage{bitmap: newBitmap(img), x: minX, y: minY}
}

type vectorStrokeSink struct {
	raster   *vector.Rasterizer
	radius   float32
	x, y     float32
	startX   float32
	startY   float32
	hasPoint bool
}

func (s *vectorStrokeSink) MoveTo(x, y float32) {
	s.x, s.y = x, y
	s.startX, s.startY = x, y
	s.hasPoint = true
}
func (s *vectorStrokeSink) LineTo(x, y float32) {
	if !s.hasPoint {
		s.MoveTo(x, y)
		return
	}
	s.segment(s.x, s.y, x, y)
	s.x, s.y = x, y
}
func (s *vectorStrokeSink) CubeTo(x1, y1, x2, y2, x, y float32) {
	if !s.hasPoint {
		s.MoveTo(x1, y1)
	}
	const steps = 16
	x0, y0 := s.x, s.y
	for i := 1; i <= steps; i++ {
		t := float32(i) / steps
		u := 1 - t
		px := u*u*u*x0 + 3*u*u*t*x1 + 3*u*t*t*x2 + t*t*t*x
		py := u*u*u*y0 + 3*u*u*t*y1 + 3*u*t*t*y2 + t*t*t*y
		s.LineTo(px, py)
	}
}
func (s *vectorStrokeSink) QuadTo(x1, y1, x, y float32) {
	if !s.hasPoint {
		s.MoveTo(x1, y1)
	}
	const steps = 12
	x0, y0 := s.x, s.y
	for i := 1; i <= steps; i++ {
		t := float32(i) / steps
		u := 1 - t
		px := u*u*x0 + 2*u*t*x1 + t*t*x
		py := u*u*y0 + 2*u*t*y1 + t*t*y
		s.LineTo(px, py)
	}
}
func (s *vectorStrokeSink) ClosePath() {
	if s.hasPoint {
		s.LineTo(s.startX, s.startY)
		s.hasPoint = false
	}
}

func (s *vectorStrokeSink) segment(x0, y0, x1, y1 float32) {
	dx, dy := x1-x0, y1-y0
	length := float32(math.Hypot(float64(dx), float64(dy)))
	if length == 0 {
		s.point(x0, y0)
		return
	}
	nx, ny := -dy/length*s.radius, dx/length*s.radius
	s.raster.MoveTo(x0+nx, y0+ny)
	s.raster.LineTo(x1+nx, y1+ny)
	s.raster.LineTo(x1-nx, y1-ny)
	s.raster.LineTo(x0-nx, y0-ny)
	s.raster.ClosePath()
	s.point(x0, y0)
	s.point(x1, y1)
}

func (s *vectorStrokeSink) point(x, y float32) {
	if s.radius <= 0 {
		return
	}
	const steps = 12
	s.raster.MoveTo(x+s.radius, y)
	for i := 1; i <= steps; i++ {
		angle := 2 * math.Pi * float64(i) / steps
		s.raster.LineTo(x+float32(math.Cos(angle))*s.radius, y+float32(math.Sin(angle))*s.radius)
	}
	s.raster.ClosePath()
}
