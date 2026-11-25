package ui

import (
	"image"
	"image/color"
	"math"
	"testing"
)

func TestGrayscaleImageUsesLuminanceAndPreservesAlpha(t *testing.T) {
	source := image.NewNRGBA(image.Rect(4, 7, 5, 8))
	source.SetNRGBA(4, 7, color.NRGBA{R: 200, G: 100, B: 50, A: 128})

	if got := grayscaleImage(source, 0); got != source {
		t.Fatal("zero grayscale allocated a replacement image")
	}
	result := grayscaleImage(source, 1)
	pixel := color.NRGBAModel.Convert(result.At(0, 0)).(color.NRGBA)
	if pixel.R != 118 || pixel.G != 118 || pixel.B != 118 || pixel.A != 128 {
		t.Fatalf("grayscale pixel = %#v, want luminance 118 and alpha 128", pixel)
	}
	original := source.NRGBAAt(4, 7)
	if original.R != 200 || original.G != 100 || original.B != 50 {
		t.Fatal("grayscale modified the source image")
	}
}

func TestImageGrayscaleUsesDistinctCacheEntries(t *testing.T) {
	resetImageCache(t)
	img := solidImage(2, 2, color.RGBA{R: 200, G: 100, B: 50, A: 255})
	h := Mount(Div(
		Img(SrcImage("mem:gray-cache", img)),
		Img(SrcImage("mem:gray-cache", img), Style(ImageGrayscale(0.5))),
		Img(SrcImage("mem:gray-cache", img), Style(ImageGrayscale(1))),
	), 100, 100)
	images := h.Root().FindAll(func(q *Query) bool { return q.rn.kind == rnImage })
	if len(images) != 3 {
		t.Fatalf("found %d image nodes, want 3", len(images))
	}
	if images[0].rn.img == images[1].rn.img || images[1].rn.img == images[2].rn.img {
		t.Fatal("original and grayscale variants share a bitmap")
	}
	if n, _, _ := ImageCacheStats(); n != 3 {
		t.Fatalf("cache entries = %d, want original plus two grayscale amounts", n)
	}
}

func TestImageGrayscaleChangeReloadsBitmap(t *testing.T) {
	resetImageCache(t)
	img := solidImage(2, 2, color.RGBA{R: 200, G: 100, B: 50, A: 255})
	app := func(_ struct{}) *Node {
		grayscale, setGrayscale := UseState(false)
		amount := float32(0)
		if grayscale {
			amount = 1
		}
		return Div(
			Img(SrcImage("mem:gray-change", img), Style(ImageGrayscale(amount))),
			Button(Text("toggle"), OnClick(func() { setGrayscale(true) })),
		)
	}
	h := Mount(Use(app, struct{}{}), 100, 100)
	original := h.Root().ByKind("image").rn.img
	bounds := h.Root().ByText("toggle").Bounds()
	h.ClickAt(bounds.X+bounds.W/2, bounds.Y+bounds.H/2)

	changed := h.Root().ByKind("image").rn.img
	if changed == nil || changed == original {
		t.Fatal("changing grayscale amount did not replace the bitmap")
	}
	pixel := color.NRGBAModel.Convert(changed.(*gioImage).src.At(0, 0)).(color.NRGBA)
	if pixel.R != 118 || pixel.G != 118 || pixel.B != 118 {
		t.Fatalf("pixel after enabling grayscale = %#v, want luminance 118", pixel)
	}
	if _, ok := lookupImage("mem:gray-change"); !ok {
		t.Fatal("original cache entry was not retained")
	}
	if _, ok := lookupImage("mem:gray-change\x00gray=1"); !ok {
		t.Fatal("grayscale cache entry is missing")
	}
}

func TestImageFilterKeysIncludeCombinedParameters(t *testing.T) {
	rn := renderNode{imgSrc: "image", imgBlur: 2, imgGray: 0.5}
	if got, want := rn.imageKey(), "image\x00blur=2\x00gray=0.5"; got != want {
		t.Fatalf("combined image key = %q, want %q", got, want)
	}

	style := newStyleProps()
	ImageBlur(float32(math.NaN()))(&style)
	ImageGrayscale(float32(math.NaN()))(&style)
	if style.imgBlur != 0 || style.imgGray != 0 {
		t.Fatalf("NaN filters = blur %v, grayscale %v, want disabled", style.imgBlur, style.imgGray)
	}
	ImageGrayscale(2)(&style)
	if style.imgGray != 1 {
		t.Fatalf("grayscale amount = %v, want clamped 1", style.imgGray)
	}
}
