package ui

import "testing"

func TestPortalScrollUpdatesPaintAndHitBounds(t *testing.T) {
	clicked := false
	h := Mount(Div(Style(Fill), Portal(
		ScrollView(Style(Absolute, Left(20), Top(20), Width(160), Height(100)),
			Div(Style(Height(180))),
			Button(Style(Height(30)), OnClick(func() { clicked = true }), Text("Last item")),
		))), 320, 240)
	before := h.Overlays()[0].ByText("Last item").Bounds()
	h.Overlays()[0].ByKind("scroll").ScrollBy(1000)
	last := h.Overlays()[0].ByText("Last item")
	if last.Bounds().Y >= before.Y || last.Bounds().Y+last.Bounds().H > 120 {
		t.Fatalf("portal content bounds did not follow scroll: before=%+v after=%+v", before, last.Bounds())
	}
	if !last.Click() || !clicked {
		t.Fatal("scrolled portal item could not be clicked at its visible position")
	}
	for _, op := range h.Paint() {
		if op.Kind == "text" && op.Text == "Last item" {
			if op.Y0 != last.Bounds().Y || op.X0 != last.Bounds().X {
				t.Fatalf("paint position (%v,%v) differs from hit bounds %+v", op.X0, op.Y0, last.Bounds())
			}
			return
		}
	}
	t.Fatal("scrolled item was not painted")
}
