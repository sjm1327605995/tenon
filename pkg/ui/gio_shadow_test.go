package ui

import "testing"

func TestShadowCacheReusesBitmap(t *testing.T) {
	resetImageCache(t)
	color := Color{R: 10, G: 20, B: 30, A: 120}
	first, _, ok := cachedShadow(64, 28, 6, 8, 0, color)
	if !ok {
		t.Fatal("shadow bitmap was not generated")
	}
	second, _, ok := cachedShadow(64, 28, 6, 8, 0, color)
	if !ok || first != second {
		t.Fatal("identical shadow did not reuse its cached bitmap")
	}
	other, _, _ := cachedShadow(64, 28, 6, 8, 0, Color{A: 120})
	if other == first {
		t.Fatal("shadow cache key ignored its color")
	}
}

func TestGioCachedShadowRenders(t *testing.T) {
	resetImageCache(t)
	img := renderOps(t, 100, 100, func(p *gioPainter) {
		p.FillRect(0, 0, 100, 100, 0, Color{255, 255, 255, 255})
		p.Shadow(30, 30, 30, 20, 4, 0, 0, 8, 0, Color{A: 255})
	})
	if n := darkPixels(img, 20, 30, 30, 60); n == 0 {
		t.Fatal("cached shadow did not render outside the source rectangle")
	}
	entries, _, _ := ImageCacheStats()
	if entries != 1 {
		t.Fatalf("shadow cache entries=%d, want 1", entries)
	}
}
