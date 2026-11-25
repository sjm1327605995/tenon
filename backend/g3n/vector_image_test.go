package g3n

import (
	"testing"

	"github.com/g3n/engine/core"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

func TestRasterizeVectorFillAndStroke(t *testing.T) {
	fill := rasterizeVector("M0 0H10V10H0Z", 2, true, 0, ui.White)
	if fill == nil || fill.bitmap.source == nil {
		t.Fatal("filled vector was not rasterized")
	}
	bounds := fill.bitmap.source.Bounds()
	if alpha := fill.bitmap.source.NRGBAAt(bounds.Dx()/2, bounds.Dy()/2).A; alpha != 255 {
		t.Fatalf("filled vector center alpha = %d, want 255", alpha)
	}
	if alpha := fill.bitmap.source.NRGBAAt(0, 0).A; alpha != 0 {
		t.Fatalf("filled vector padding alpha = %d, want 0", alpha)
	}

	stroke := rasterizeVector("M2 2L18 18", 1, false, 2, ui.White)
	if stroke == nil || stroke.bitmap.source == nil {
		t.Fatal("stroked vector was not rasterized")
	}
	if alpha := stroke.bitmap.source.NRGBAAt(9, 9).A; alpha == 0 {
		t.Fatal("stroked vector has no ink along its diagonal")
	}
}

func TestVectorImageCacheIsScopedByAppearance(t *testing.T) {
	path := &vectorPath{svg: "M0 0H10V10H0Z", scale: 1}
	first := path.image(true, 0, ui.White)
	second := path.image(true, 0, ui.White)
	if first == nil || first != second {
		t.Fatal("identical vector appearance did not reuse its cached image")
	}
	if other := path.image(true, 0, ui.Black); other == first {
		t.Fatal("different vector colors shared the same rasterized image")
	}
	if stroked := path.image(false, 2, ui.White); stroked == first {
		t.Fatal("fill and stroke appearances shared the same rasterized image")
	}
}

func TestPainterAddsVectorImages(t *testing.T) {
	root := core.NewNode()
	painter := newPainter(root, nil, 100, 100)
	path := &vectorPath{svg: "M0 0H10V10H0Z", scale: 1}
	painter.FillPath(path, 5, 5, ui.White)
	painter.StrokePath(path, 5, 5, 2, ui.Black)
	if len(root.Children()) != 2 {
		t.Fatalf("vector paint produced %d panels, want 2", len(root.Children()))
	}
}

func TestBackendDeclaresVerifiedCapabilities(t *testing.T) {
	capabilities := (&Backend{}).Capabilities()
	if !capabilities.SVGPaths || !capabilities.LayerTransforms {
		t.Fatal("G3N verified SVG or layer-transform support is not declared")
	}
	if capabilities.RoundedClipping || capabilities.ExactProjection || capabilities.Accessibility {
		t.Fatal("G3N declared an advanced capability that is not yet implemented")
	}
}
