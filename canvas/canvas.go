package canvas

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/font"
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/widget"
)

// EbitenCanvas implements widget.Canvas using ebiten.
type EbitenCanvas struct {
	target         *ebiten.Image
	clipStack      []geometry.Rect
	transformStack []geometry.Point
	screenOrigin   geometry.Point
}

// New creates a new EbitenCanvas wrapping the given ebiten.Image.
func New(img *ebiten.Image) *EbitenCanvas {
	return &EbitenCanvas{
		target:         img,
		clipStack:      []geometry.Rect{},
		transformStack: []geometry.Point{{X: 0, Y: 0}},
	}
}

func (c *EbitenCanvas) currentOffset() geometry.Point {
	return c.transformStack[len(c.transformStack)-1]
}

func (c *EbitenCanvas) toScreen(r geometry.Rect) geometry.Rect {
	off := c.currentOffset()
	return geometry.Rect{
		Min: geometry.Pt(r.Min.X+off.X, r.Min.Y+off.Y),
		Max: geometry.Pt(r.Max.X+off.X, r.Max.Y+off.Y),
	}
}

func (c *EbitenCanvas) toScreenPt(p geometry.Point) geometry.Point {
	off := c.currentOffset()
	return geometry.Pt(p.X+off.X, p.Y+off.Y)
}

func toEbitenColor(c widget.Color) color.Color {
	r, g, b, a := c.RGBA8()
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

// Clear fills the entire canvas with the given color.
func (c *EbitenCanvas) Clear(col widget.Color) {
	c.target.Fill(toEbitenColor(col))
}

// DrawRect fills a rectangle with the given color.
func (c *EbitenCanvas) DrawRect(r geometry.Rect, col widget.Color) {
	sr := c.toScreen(r)
	vector.DrawFilledRect(c.target,
		float32(sr.Min.X), float32(sr.Min.Y),
		float32(sr.Width()), float32(sr.Height()),
		toEbitenColor(col), true)
}

// StrokeRect draws the outline of a rectangle.
func (c *EbitenCanvas) StrokeRect(r geometry.Rect, col widget.Color, strokeWidth float32) {
	sr := c.toScreen(r)
	vector.StrokeRect(c.target,
		float32(sr.Min.X), float32(sr.Min.Y),
		float32(sr.Width()), float32(sr.Height()),
		strokeWidth, toEbitenColor(col), true)
}

// DrawRoundRect fills a rounded rectangle with the given color.
func (c *EbitenCanvas) DrawRoundRect(r geometry.Rect, col widget.Color, radius float32) {
	sr := c.toScreen(r)
	path := roundedRectPath(sr, radius)
	drawPath(c.target, &path, toEbitenColor(col))
}

// StrokeRoundRect draws the outline of a rounded rectangle.
func (c *EbitenCanvas) StrokeRoundRect(r geometry.Rect, col widget.Color, radius float32, strokeWidth float32) {
	sr := c.toScreen(r)
	path := roundedRectPath(sr, radius)
	strokePath(c.target, &path, toEbitenColor(col), strokeWidth)
}

// DrawCircle fills a circle with the given color.
func (c *EbitenCanvas) DrawCircle(center geometry.Point, radius float32, col widget.Color) {
	sc := c.toScreenPt(center)
	vector.DrawFilledCircle(c.target,
		float32(sc.X), float32(sc.Y),
		radius, toEbitenColor(col), true)
}

// StrokeCircle draws the outline of a circle.
func (c *EbitenCanvas) StrokeCircle(center geometry.Point, radius float32, col widget.Color, strokeWidth float32) {
	sc := c.toScreenPt(center)
	vector.StrokeCircle(c.target,
		float32(sc.X), float32(sc.Y),
		radius, strokeWidth, toEbitenColor(col), true)
}

// StrokeArc draws a circular arc outline.
func (c *EbitenCanvas) StrokeArc(center geometry.Point, radius float32, startAngle, sweepAngle float64,
	col widget.Color, strokeWidth float32) {
	sc := c.toScreenPt(center)
	path := vector.Path{}
	path.MoveTo(float32(sc.X)+radius*cosf32(startAngle), float32(sc.Y)+radius*sinf32(startAngle))
	path.Arc(float32(sc.X), float32(sc.Y), radius, float32(startAngle), float32(startAngle+sweepAngle), vector.Clockwise)
	strokePath(c.target, &path, toEbitenColor(col), strokeWidth)
}

// DrawLine draws a line between two points.
func (c *EbitenCanvas) DrawLine(from, to geometry.Point, col widget.Color, strokeWidth float32) {
	sf := c.toScreenPt(from)
	st := c.toScreenPt(to)
	vector.StrokeLine(c.target,
		float32(sf.X), float32(sf.Y),
		float32(st.X), float32(st.Y),
		strokeWidth, toEbitenColor(col), true)
}



// DrawImage draws an image at the specified position.
func (c *EbitenCanvas) DrawImage(img image.Image, at geometry.Point) {
	sat := c.toScreenPt(at)
	eimg := ebiten.NewImageFromImage(img)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(sat.X), float64(sat.Y))
	c.target.DrawImage(eimg, op)
}

// PushClip pushes a clipping rectangle onto the clip stack.
func (c *EbitenCanvas) PushClip(r geometry.Rect) {
	c.clipStack = append(c.clipStack, r)
}

// PopClip removes the most recently pushed clipping region.
func (c *EbitenCanvas) PopClip() {
	if len(c.clipStack) > 0 {
		c.clipStack = c.clipStack[:len(c.clipStack)-1]
	}
}

// PushTransform pushes a translation transform onto the transform stack.
func (c *EbitenCanvas) PushTransform(offset geometry.Point) {
	current := c.currentOffset()
	c.transformStack = append(c.transformStack, geometry.Pt(current.X+offset.X, current.Y+offset.Y))
}

// PopTransform removes the most recently pushed transform.
func (c *EbitenCanvas) PopTransform() {
	if len(c.transformStack) > 1 {
		c.transformStack = c.transformStack[:len(c.transformStack)-1]
	}
}

// TransformOffset returns the current cumulative transform offset.
func (c *EbitenCanvas) TransformOffset() geometry.Point {
	return c.currentOffset()
}

// ScreenOriginBase returns the screen-space base offset for this canvas.
func (c *EbitenCanvas) ScreenOriginBase() geometry.Point {
	return c.screenOrigin
}

// ClipBounds returns the current clip rectangle.
func (c *EbitenCanvas) ClipBounds() geometry.Rect {
	if len(c.clipStack) == 0 {
		return geometry.Rect{Min: geometry.Pt(0, 0), Max: geometry.Pt(1e9, 1e9)}
	}
	return c.clipStack[len(c.clipStack)-1]
}

// DrawText draws text within the given bounding rectangle.
func (c *EbitenCanvas) DrawText(s string, bounds geometry.Rect, fontSize float32, col widget.Color, bold bool, align widget.TextAlign) {
	sr := c.toScreen(bounds)
	if sr.IsEmpty() {
		return
	}

	face := text.NewGoXFace(font.Face(fontSize))
	w, h := text.Measure(s, face, 0)

	x := float64(sr.Min.X)
	y := float64(sr.Min.Y)
	switch align {
	case widget.TextAlignCenter:
		x += (float64(sr.Width()) - w) / 2
	case widget.TextAlignRight:
		x += float64(sr.Width()) - w
	}
	y += (float64(sr.Height()) - h) / 2

	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	r8, g8, b8, a8 := col.RGBA8()
	op.ColorScale.Scale(float32(r8)/255, float32(g8)/255, float32(b8)/255, float32(a8)/255)
	text.Draw(c.target, s, face, op)
}

// MeasureText returns the width in pixels of the given text string.
func (c *EbitenCanvas) MeasureText(s string, fontSize float32, bold bool) float32 {
	face := text.NewGoXFace(font.Face(fontSize))
	w, _ := text.Measure(s, face, 0)
	return float32(w)
}
