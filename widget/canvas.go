package widget

import (
	"image"

	"github.com/sjm1327605995/tenon/geometry"
)

// TextAlign specifies horizontal text alignment within bounds.
type TextAlign uint8

const (
	// TextAlignLeft aligns text to the left edge (default).
	TextAlignLeft TextAlign = iota

	// TextAlignCenter centers text horizontally.
	TextAlignCenter

	// TextAlignRight aligns text to the right edge.
	TextAlignRight
)

// textAlignNames maps each TextAlign to its human-readable name.
var textAlignNames = [...]string{
	TextAlignLeft:   "Left",
	TextAlignCenter: "Center",
	TextAlignRight:  "Right",
}

// unknownStr is the string representation for unknown/unrecognized values.
const unknownStr = "Unknown"

// String returns a human-readable name for the text alignment.
func (a TextAlign) String() string {
	if int(a) < len(textAlignNames) {
		return textAlignNames[a]
	}
	return unknownStr
}

// Float64 returns the alignment as a float64 value for rendering backends.
// Left=0, Center=0.5, Right=1.
func (a TextAlign) Float64() float64 {
	switch a {
	case TextAlignCenter:
		return 0.5
	case TextAlignRight:
		return 1.0
	default:
		return 0.0
	}
}

// Canvas provides drawing operations for widgets.
type Canvas interface {
	// Clear fills the entire canvas with the given color.
	Clear(color Color)

	// DrawRect fills a rectangle with the given color.
	DrawRect(r geometry.Rect, color Color)

	// StrokeRect draws the outline of a rectangle.
	StrokeRect(r geometry.Rect, color Color, strokeWidth float32)

	// DrawRoundRect fills a rounded rectangle with the given color.
	DrawRoundRect(r geometry.Rect, color Color, radius float32)

	// StrokeRoundRect draws the outline of a rounded rectangle.
	StrokeRoundRect(r geometry.Rect, color Color, radius float32, strokeWidth float32)

	// DrawCircle fills a circle with the given color.
	DrawCircle(center geometry.Point, radius float32, color Color)

	// StrokeCircle draws the outline of a circle.
	StrokeCircle(center geometry.Point, radius float32, color Color, strokeWidth float32)

	// StrokeArc draws a circular arc outline from startAngle with the given sweep.
	StrokeArc(center geometry.Point, radius float32, startAngle, sweepAngle float64,
		color Color, strokeWidth float32)

	// DrawLine draws a line between two points.
	DrawLine(from, to geometry.Point, color Color, strokeWidth float32)

	// DrawText draws text within the given bounding rectangle.
	DrawText(text string, bounds geometry.Rect, fontSize float32, color Color, bold bool, align TextAlign)

	// MeasureText returns the width in pixels of the given text string.
	MeasureText(text string, fontSize float32, bold bool) float32

	// DrawImage draws an image at the specified position.
	DrawImage(img image.Image, at geometry.Point)

	// PushClip pushes a clipping rectangle onto the clip stack.
	PushClip(r geometry.Rect)

	// PopClip removes the most recently pushed clipping region.
	PopClip()

	// PushTransform pushes a translation transform onto the transform stack.
	PushTransform(offset geometry.Point)

	// PopTransform removes the most recently pushed transform.
	PopTransform()

	// TransformOffset returns the current cumulative transform offset.
	TransformOffset() geometry.Point

	// ScreenOriginBase returns the screen-space base offset for this canvas.
	ScreenOriginBase() geometry.Point

	// ClipBounds returns the current clip rectangle.
	ClipBounds() geometry.Rect
}

// LineCap specifies how the endpoints of stroked arcs and lines are drawn.
type LineCap uint8

const (
	LineCapButt   LineCap = iota // flat end, stops exactly at endpoint
	LineCapRound                 // semicircle extending past endpoint
	LineCapSquare                // half-square extending past endpoint
)

// TextStyle specifies font properties for styled text rendering.
type TextStyle struct {
	FontFamily string
	FontSize   float32
	Bold       bool
	Italic     bool
	Color      Color
	Align      TextAlign
}

// StyledTextDrawer is an optional interface for canvases that support
// rendering text with custom fonts from the font registry.
type StyledTextDrawer interface {
	DrawStyledText(text string, bounds geometry.Rect, style TextStyle)
	MeasureStyledText(text string, style TextStyle) float32
}

// Color represents an RGBA color with float32 components.
type Color struct {
	R, G, B, A float32
}

// RGBA creates a Color from red, green, blue, and alpha components.
func RGBA(r, g, b, a float32) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// RGB creates an opaque Color from red, green, and blue components.
func RGB(r, g, b float32) Color {
	return Color{R: r, G: g, B: b, A: 1}
}

// RGBA8 creates a Color from 8-bit RGBA values (0-255).
func RGBA8(r, g, b, a uint8) Color {
	return Color{
		R: float32(r) / colorMax8,
		G: float32(g) / colorMax8,
		B: float32(b) / colorMax8,
		A: float32(a) / colorMax8,
	}
}

// RGB8 creates an opaque Color from 8-bit RGB values (0-255).
func RGB8(r, g, b uint8) Color {
	return Color{
		R: float32(r) / colorMax8,
		G: float32(g) / colorMax8,
		B: float32(b) / colorMax8,
		A: 1,
	}
}

// Hex creates a Color from a 24-bit hex value (0xRRGGBB).
func Hex(hex uint32) Color {
	return Color{
		R: float32((hex>>16)&0xFF) / colorMax8,
		G: float32((hex>>8)&0xFF) / colorMax8,
		B: float32(hex&0xFF) / colorMax8,
		A: 1,
	}
}

// HexA creates a Color from a 32-bit hex value with alpha (0xRRGGBBAA).
func HexA(hex uint32) Color {
	return Color{
		R: float32((hex>>24)&0xFF) / colorMax8,
		G: float32((hex>>16)&0xFF) / colorMax8,
		B: float32((hex>>8)&0xFF) / colorMax8,
		A: float32(hex&0xFF) / colorMax8,
	}
}

// colorMax8 is the maximum value for 8-bit color components.
const colorMax8 = 255.0

// WithAlpha returns a copy of the color with the given alpha value.
func (c Color) WithAlpha(a float32) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: a}
}

// Lerp returns a color linearly interpolated between c and other.
func (c Color) Lerp(other Color, t float32) Color {
	return Color{
		R: c.R + (other.R-c.R)*t,
		G: c.G + (other.G-c.G)*t,
		B: c.B + (other.B-c.B)*t,
		A: c.A + (other.A-c.A)*t,
	}
}

// IsOpaque returns true if the color is fully opaque (alpha == 1).
func (c Color) IsOpaque() bool {
	return c.A >= 1.0
}

// IsTransparent returns true if the color is fully transparent (alpha == 0).
func (c Color) IsTransparent() bool {
	return c.A <= 0.0
}

// RGBA8 returns the color as 8-bit RGBA components (0-255).
func (c Color) RGBA8() (r, g, b, a uint8) {
	return uint8(clamp01(c.R) * colorMax8),
		uint8(clamp01(c.G) * colorMax8),
		uint8(clamp01(c.B) * colorMax8),
		uint8(clamp01(c.A) * colorMax8)
}

// clamp01 clamps a float32 value to the range [0, 1].
func clamp01(v float32) float32 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// Common color constants.
var (
	ColorTransparent = Color{R: 0, G: 0, B: 0, A: 0}
	ColorBlack       = Color{R: 0, G: 0, B: 0, A: 1}
	ColorWhite       = Color{R: 1, G: 1, B: 1, A: 1}
	ColorRed         = Color{R: 1, G: 0, B: 0, A: 1}
	ColorGreen       = Color{R: 0, G: 1, B: 0, A: 1}
	ColorBlue        = Color{R: 0, G: 0, B: 1, A: 1}
	ColorYellow      = Color{R: 1, G: 1, B: 0, A: 1}
	ColorCyan        = Color{R: 0, G: 1, B: 1, A: 1}
	ColorMagenta     = Color{R: 1, G: 0, B: 1, A: 1}
	ColorGray        = Color{R: 0.5, G: 0.5, B: 0.5, A: 1}
	ColorLightGray   = Color{R: 0.75, G: 0.75, B: 0.75, A: 1}
	ColorDarkGray    = Color{R: 0.25, G: 0.25, B: 0.25, A: 1}
)
