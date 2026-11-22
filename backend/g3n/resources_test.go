package g3n

import (
	"image"
	"testing"

	"github.com/g3n/engine/window"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

func TestBuiltinFontMeasuresCJK(t *testing.T) {
	backend := &Backend{fonts: make(map[fontKey]*font)}
	face := backend.NewFont(16, 400, false)
	if face == nil {
		t.Fatal("built-in font was not created")
	}
	if got := face.Measure("中文 Tenon", 20); got <= 0 {
		t.Fatalf("CJK text width = %v, want positive", got)
	}
	ascent, _, _ := face.Metrics()
	if ascent <= 0 {
		t.Fatalf("font ascent = %v, want positive", ascent)
	}
}

func TestInputKeyMapping(t *testing.T) {
	for key, want := range map[window.Key]ui.InputKey{
		window.KeyEnter:       ui.KeyEnter,
		window.KeyKPEnter:     ui.KeyEnter,
		window.KeyLeftShift:   ui.KeyShift,
		window.KeyRightShift:  ui.KeyShift,
		window.KeyLeftControl: ui.KeyCtrl,
		window.KeyBackspace:   ui.KeyBackspace,
	} {
		got := inputKeys(key)
		if len(got) != 1 || got[0] != want {
			t.Fatalf("inputKeys(%v) = %v, want %v", key, got, want)
		}
	}
}

func TestBitmapOpacityKeepsStraightColor(t *testing.T) {
	source := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	source.Pix = []uint8{20, 40, 60, 128}
	got := bitmapRGBA(source, 128)
	if got.Pix[0] != 20 || got.Pix[1] != 40 || got.Pix[2] != 60 || got.Pix[3] != 64 {
		t.Fatalf("opacity pixel = %v, want [20 40 60 64]", got.Pix)
	}
}

func TestColor4Conversion(t *testing.T) {
	got := color4(ui.Color{R: 255, G: 128, B: 0, A: 64})
	if got.R != 1 || got.G != float32(128)/255 || got.B != 0 || got.A != float32(64)/255 {
		t.Fatalf("color4 = %+v", got)
	}
}
