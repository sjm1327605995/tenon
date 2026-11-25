package ui

// SVGPathSink receives backend-neutral SVG path geometry after applying scale.
// Backends can rasterize the same commands instead of maintaining separate parsers.
type SVGPathSink interface {
	MoveTo(x, y float32)
	LineTo(x, y float32)
	CubeTo(x1, y1, x2, y2, x, y float32)
	QuadTo(x1, y1, x, y float32)
	ClosePath()
}

// ParseSVGPath parses an SVG path's d attribute and emits scaled geometry.
// Elliptical arcs are flattened; all other commands retain their curve types.
func ParseSVGPath(d string, scale float32, sink SVGPathSink) {
	if sink == nil {
		return
	}
	parseSVGInto(d, scale, exportedSVGPathSink{sink})
}

type exportedSVGPathSink struct{ SVGPathSink }

func (s exportedSVGPathSink) moveTo(x, y float32) { s.MoveTo(x, y) }
func (s exportedSVGPathSink) lineTo(x, y float32) { s.LineTo(x, y) }
func (s exportedSVGPathSink) cubeTo(x1, y1, x2, y2, x, y float32) {
	s.CubeTo(x1, y1, x2, y2, x, y)
}
func (s exportedSVGPathSink) quadTo(x1, y1, x, y float32) { s.QuadTo(x1, y1, x, y) }
func (s exportedSVGPathSink) closePath()                  { s.ClosePath() }
