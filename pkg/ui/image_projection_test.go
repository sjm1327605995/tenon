package ui

import (
	"image"
	"testing"
)

func TestWarpImageToQuadUsesSharedHomography(t *testing.T) {
	source := image.NewRGBA(image.Rect(0, 0, 2, 2))
	source.Set(0, 0, Red)
	source.Set(1, 0, Red)
	source.Set(0, 1, Red)
	source.Set(1, 1, Red)
	warped, box := WarpImageToQuad(source, [4]Point{{X: 10, Y: 20}, {X: 110, Y: 20}, {X: 10, Y: 120}, {X: 110, Y: 120}})
	if warped == nil || box.W != 100 || box.H != 100 {
		t.Fatalf("warped image or bounds = %v, %+v", warped, box)
	}
	_, _, _, alpha := warped.At(50, 50).RGBA()
	if alpha != 0xffff {
		t.Fatalf("warped center alpha = %x, want opaque", alpha)
	}
}
