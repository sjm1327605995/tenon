package ui

import (
	"testing"
)

func TestRadialGradientImagePixels(t *testing.T) {
	from, to := Hex("#ff0000"), Hex("#0000ff")
	img := newRadialGradientImage(20, 20, from, to)
	if got := img.NRGBAAt(10, 10); got.R != 255 || got.B != 0 {
		t.Fatalf("center pixel=%v, want red", got)
	}
	if got := img.NRGBAAt(0, 0); got.R != 0 || got.B != 255 {
		t.Fatalf("corner pixel=%v, want blue", got)
	}
}

func TestRadialGradientPaintUsesClippedCachedBitmap(t *testing.T) {
	rn := &renderNode{
		kind: rnBox, bounds: Rect{X: 10, Y: 20, W: 20, H: 20}, radius: 4,
		hasGradient: true, gradRadial: true, gradFrom: Hex("#ff0000"), gradTo: Hex("#0000ff"),
	}
	painter := &recordPainter{}
	paintSelf(painter, rn)
	if rn.radialImg == nil || len(painter.ops) != 3 || painter.ops[0].Kind != "clip" || painter.ops[1].Kind != "image" || painter.ops[2].Kind != "unclip" {
		t.Fatalf("radial paint ops=%+v", painter.ops)
	}
	first := rn.radialImg
	paintSelf(painter, rn)
	if rn.radialImg != first {
		t.Fatal("unchanged radial gradient rebuilt its bitmap")
	}
	rn.bounds.W = 40
	paintSelf(painter, rn)
	if rn.radialImg == first {
		t.Fatal("resized radial gradient retained the wrong bitmap")
	}
}
