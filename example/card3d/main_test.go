package main

import (
	"testing"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func TestCardContentFitsWithinCard(t *testing.T) {
	h := ui.Mount(ui.Use(App, struct{}{}), 720, 760)
	defer h.Close()
	photo := h.Root().ByKind("image")
	third := h.Root().ByText("冰岛南岸")
	if !photo.Exists() || !third.Exists() {
		t.Fatal("card image or third content region is missing")
	}
	photoBounds := photo.Bounds()
	thirdBounds := third.Bounds()
	if thirdBounds.W == 0 || thirdBounds.H == 0 {
		t.Fatalf("third content has zero bounds: %+v", thirdBounds)
	}
	cardTop, cardBottom := photoBounds.Y, photoBounds.Y+400
	if thirdBounds.Y < cardTop || thirdBounds.Y+thirdBounds.H > cardBottom {
		t.Fatalf("third content is outside the card: third=%+v card vertical range=%v..%v", thirdBounds, cardTop, cardBottom)
	}
}

func TestCardStartsProjected(t *testing.T) {
	h := ui.Mount(ui.Use(App, struct{}{}), 720, 760)
	defer h.Close()
	photo := h.Root().ByKind("image")
	if !photo.Exists() {
		t.Fatal("card image is missing")
	}
	layout, visual := photo.Bounds(), photo.VisualBounds()
	if visual == layout {
		t.Fatalf("initial image is not projected: layout and visual bounds are both %+v", layout)
	}
}
