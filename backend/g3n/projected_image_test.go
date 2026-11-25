package g3n

import (
	"image"
	"testing"

	"github.com/g3n/engine/core"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

func TestProjectedImageUsesWarpedBitmapAndCache(t *testing.T) {
	clear(projectedImages)
	source := newBitmap(image.NewNRGBA(image.Rect(0, 0, 2, 2)))
	root := core.NewNode()
	painter := newPainter(root, nil, 200, 200)
	quad := [4]ui.Point{{X: 10, Y: 20}, {X: 110, Y: 20}, {X: 10, Y: 120}, {X: 110, Y: 120}}
	painter.ProjectedImage(source, quad, 1, "source", nil, 0)
	if len(root.Children()) != 1 {
		t.Fatalf("projected image produced %d panels, want 1", len(root.Children()))
	}
	entry := projectedImages[projectedImageKey("source", quad)]
	if entry.bitmap == nil || entry.bitmap.source == nil || entry.box.W != 100 || entry.box.H != 100 {
		t.Fatalf("projected cache entry = %+v", entry)
	}
	painter.ProjectedImage(source, quad, 1, "source", nil, 0)
	if cached := projectedImages[projectedImageKey("source", quad)]; cached.bitmap != entry.bitmap {
		t.Fatal("identical projection did not reuse its cached bitmap")
	}
}

func TestProjectedImageCacheSeparatesQuads(t *testing.T) {
	clear(projectedImages)
	first := [4]ui.Point{{X: 0, Y: 0}, {X: 100, Y: 0}, {X: 0, Y: 100}, {X: 100, Y: 100}}
	second := [4]ui.Point{{X: 10, Y: 0}, {X: 90, Y: 20}, {X: 0, Y: 90}, {X: 100, Y: 100}}
	if projectedImageKey("source", first) == projectedImageKey("source", second) {
		t.Fatal("different projection quads share a cache key")
	}
}
