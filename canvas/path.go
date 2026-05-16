package canvas

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawPath fills a vector path with the given color.
func drawPath(dst *ebiten.Image, path *vector.Path, c color.Color) {
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	if len(vs) == 0 {
		return
	}
	// Apply color to vertices
	r, g, b, a := c.RGBA()
	rf := float32(r) / 0xffff
	gf := float32(g) / 0xffff
	bf := float32(b) / 0xffff
	af := float32(a) / 0xffff
	for i := range vs {
		vs[i].ColorR = rf
		vs[i].ColorG = gf
		vs[i].ColorB = bf
		vs[i].ColorA = af
	}
	op := &ebiten.DrawTrianglesOptions{}
	dst.DrawTriangles(vs, is, whitePixel, op)
}

// strokePath strokes a vector path with the given color and width.
func strokePath(dst *ebiten.Image, path *vector.Path, c color.Color, width float32) {
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width: width,
	})
	if len(vs) == 0 {
		return
	}
	r, g, b, a := c.RGBA()
	rf := float32(r) / 0xffff
	gf := float32(g) / 0xffff
	bf := float32(b) / 0xffff
	af := float32(a) / 0xffff
	for i := range vs {
		vs[i].ColorR = rf
		vs[i].ColorG = gf
		vs[i].ColorB = bf
		vs[i].ColorA = af
	}
	op := &ebiten.DrawTrianglesOptions{}
	dst.DrawTriangles(vs, is, whitePixel, op)
}

// whitePixel is a 1x1 white image used as a texture for colored triangles.
var whitePixel = ebiten.NewImage(1, 1)

func init() {
	whitePixel.Fill(color.White)
}
