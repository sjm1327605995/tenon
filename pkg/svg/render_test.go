package svg

import (
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func TestRenderSearchIcon(t *testing.T) {
	d := "M909.6 854.5L649.9 594.8C690.2 542.7 714 479.2 714 410.5 714 234.3 571.7 92 395.5 92S77 234.3 77 410.5 219.3 729 395.5 729c68.7 0 132.2-23.8 184.3-64.1l259.7 259.7c7.8 7.8 20.5 7.8 28.3 0l42.8-42.8c7.8-7.8 7.8-20.5 0-28.3zM395.5 643c-128.2 0-232.5-104.3-232.5-232.5S267.3 178 395.5 178 628 282.3 628 410.5 523.7 643 395.5 643z"

	minX, minY, maxX, maxY, err := ParsePathBounds(d)
	if err != nil {
		t.Fatal(err)
	}
	viewW := maxX - minX
	viewH := maxY - minY
	scale := float32(64) / max(viewW, viewH)

	path, err := ParsePathScaledAndShifted(d, scale, -minX*scale, -minY*scale)
	if err != nil {
		t.Fatal(err)
	}

	imgW := int(viewW * scale)
	imgH := int(viewH * scale)
	img := ebiten.NewImage(imgW, imgH)
	img.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255})

	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	op.AntiAlias = true

	vector.FillPath(img, path, &vector.FillOptions{}, op)

	f, err := os.Create("search_icon_test.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	t.Logf("Written search_icon_test.png (%dx%d)", imgW, imgH)
}
