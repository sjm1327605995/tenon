package g3n

import (
	"image"
	"image/color"
	"testing"

	"github.com/sjm1327605995/tenon/pkg/ui"
)

func TestFontMetricsMatchSyntheticStyleThresholds(t *testing.T) {
	backend := &Backend{fonts: make(map[fontKey]*font)}
	for _, test := range []struct {
		weight  int
		italic  bool
		bold    bool
		oblique bool
	}{
		{weight: 400}, {weight: 500}, {weight: 600, bold: true}, {weight: 700, bold: true},
		{weight: 400, italic: true, oblique: true}, {weight: 700, italic: true, bold: true, oblique: true},
	} {
		face := backend.NewFont(24, test.weight, test.italic).(*font)
		_, bold, italic := face.Metrics()
		if bold != test.bold || italic != test.oblique {
			t.Errorf("weight=%d italic=%v -> faux=%v/%v, want %v/%v", test.weight, test.italic, bold, italic, test.bold, test.oblique)
		}
	}
}

func TestStyledTextChangesPixels(t *testing.T) {
	backend := &Backend{fonts: make(map[fontKey]*font)}
	face := backend.NewFont(28, 400, false).(*font)
	regular := styledTextImage(face, "Hamburg", false, false)
	if alpha := regular.NRGBAAt(0, 0).A; alpha != 0 {
		t.Fatalf("text background alpha = %d, want transparent", alpha)
	}
	bold := styledTextImage(face, "Hamburg", true, false)
	italic := styledTextImage(face, "Hamburg", false, true)
	regularInk, boldInk := textInk(regular), textInk(bold)
	if boldInk <= regularInk {
		t.Fatalf("bold ink = %d, regular ink = %d, want more bold pixels", boldInk, regularInk)
	}
	row := firstTextInkRow(regular)
	if row < 0 {
		t.Fatal("regular text has no ink")
	}
	if leftTextInk(regular, row) >= leftTextInk(italic, row) {
		t.Fatalf("italic did not shear upper glyphs rightward: regular x=%d italic x=%d ink=%d firstRow=%d", leftTextInk(regular, row), leftTextInk(italic, row), textInk(italic), firstTextInkRow(italic))
	}
}

func textInk(img *image.NRGBA) int {
	count := 0
	for i := 3; i < len(img.Pix); i += 4 {
		if img.Pix[i] > 16 {
			count++
		}
	}
	return count
}

func firstTextInkRow(img *image.NRGBA) int {
	for y := 0; y < img.Bounds().Dy(); y++ {
		if leftTextInk(img, y) >= 0 {
			return y
		}
	}
	return -1
}

func leftTextInk(img *image.NRGBA, y int) int {
	for x := 0; x < img.Bounds().Dx(); x++ {
		if img.NRGBAAt(x, y).A > 16 {
			return x
		}
	}
	return -1
}

func TestTintTextImageUsesGlyphAlpha(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	img.SetNRGBA(0, 0, color.NRGBA{A: 128})
	img.SetNRGBA(1, 0, color.NRGBA{})
	tintTextImage(img, ui.Color{R: 255, A: 128})
	if got := img.NRGBAAt(0, 0); got.R != 255 || got.A != 64 {
		t.Fatalf("tinted glyph = %+v, want red with alpha 64", got)
	}
	if got := img.NRGBAAt(1, 0); got.A != 0 {
		t.Fatalf("transparent background alpha = %d, want 0", got.A)
	}
}
