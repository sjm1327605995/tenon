package ui

import (
	"image"
	"testing"
)

func TestInputOpacityFadesTextAndBorder(t *testing.T) {
	for _, multiline := range []bool{false, true} {
		name := "single-line"
		if multiline {
			name = "multiline"
		}
		t.Run(name, func(t *testing.T) {
			field := func(opacity float32) *Node {
				attrs := []*Node{Value("Opacity"), Style(Width(180), Height(50), PaddingXY(8, 8),
					Bg(White), Border(1, Black), TextColor(Black), FontSize(16), Opacity(opacity))}
				if multiline {
					attrs = append(attrs, Multiline())
				}
				return Input(attrs...)
			}
			h := Mount(Div(Style(Row, Fill, Bg(White)), field(1), field(0.5)), 360, 60)
			img := renderOps(t, 360, 60, func(p *gioPainter) { paint(p, h.g.rootRN) })
			darkest := func(r image.Rectangle) uint32 {
				minR := uint32(255)
				for y := r.Min.Y; y < r.Max.Y; y++ {
					for x := r.Min.X; x < r.Max.X; x++ {
						r, _, _, _ := img.At(x, y).RGBA()
						minR = min(minR, r>>8)
					}
				}
				return minR
			}
			for _, region := range []image.Rectangle{image.Rect(8, 5, 165, 45), image.Rect(20, 0, 160, 2)} {
				opaque := darkest(region)
				faded := darkest(region.Add(image.Pt(180, 0)))
				if opaque > 245 || faded < opaque+10 || faded > 250 {
					t.Fatalf("opacity did not fade the whole input at %+v: opaque=%d faded=%d", region, opaque, faded)
				}
			}
		})
	}
}
