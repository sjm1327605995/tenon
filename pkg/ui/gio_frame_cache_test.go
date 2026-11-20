package ui

import (
	"image"
	"testing"
	"time"

	"gioui.org/gpu/headless"
	"gioui.org/op"
)

func TestGioFrameCacheReusesStaticFrame(t *testing.T) {
	h := MountDefault(Text("static"))
	defer h.Close()
	now := time.UnixMilli(1000)
	version := h.g.framePaintVersion(now)
	var ops op.Ops

	h.g.frameCache.Add(h.g, &ops, version)
	if h.g.frameCache.builds != 1 {
		t.Fatal("initial frame was not recorded")
	}
	ops.Reset()
	h.g.frameCache.Add(h.g, &ops, version)
	if h.g.frameCache.builds != 1 {
		t.Fatal("static frame was recorded twice")
	}
	h.g.visualVersion++
	ops.Reset()
	h.g.frameCache.Add(h.g, &ops, h.g.framePaintVersion(now))
	if h.g.frameCache.builds != 2 {
		t.Fatal("visual version change did not invalidate frame cache")
	}
}

func TestGioFrameCachePixelsMatch(t *testing.T) {
	h := MountDefault(Div(Style(Width(80), Height(40), Bg(White), Radius(8))))
	defer h.Close()
	win, err := headless.NewWindow(100, 60)
	if err != nil {
		t.Skipf("headless GPU unavailable: %v", err)
	}
	defer win.Release()
	capture := func() *image.RGBA {
		var ops op.Ops
		h.g.frameCache.Add(h.g, &ops, 1)
		if err := win.Frame(&ops); err != nil {
			t.Fatal(err)
		}
		img := image.NewRGBA(image.Rect(0, 0, 100, 60))
		if err := win.Screenshot(img); err != nil {
			t.Fatal(err)
		}
		return img
	}

	first := capture()
	second := capture()
	maxDelta := 0
	for i := range first.Pix {
		delta := int(first.Pix[i]) - int(second.Pix[i])
		if delta < 0 {
			delta = -delta
		}
		if delta > maxDelta {
			maxDelta = delta
		}
	}
	if maxDelta != 0 {
		t.Fatalf("cached frame differs from the original by up to %d levels", maxDelta)
	}
}

func TestCaretBlinkChangesFrameVersion(t *testing.T) {
	h := MountDefault(Input(Value("x")))
	defer h.Close()
	h.Root().ByKind("input").Focus()
	if h.g.framePaintVersion(time.UnixMilli(0)) == h.g.framePaintVersion(time.UnixMilli(500)) {
		t.Fatal("caret blink phase did not change the frame version")
	}
}
