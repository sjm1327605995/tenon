package ui

import "testing"

func TestScaleXYUsesIndependentAxes(t *testing.T) {
	h := Mount(Div(Style(Width(100), Height(100), ScaleXY(2, 0.5))), 200, 200)
	defer h.Close()
	rn := h.Root().rn
	if !rn.hasTransform() {
		t.Fatal("ScaleXY did not create a transformed layer")
	}
	transform := layerOf(rn, nil)
	snapshot := transform.Snapshot()
	if snapshot.Scale != 2 || snapshot.ScaleY != 0.5 {
		t.Fatalf("scale = (%v,%v), want (2,0.5)", snapshot.Scale, snapshot.ScaleY)
	}
	x, y := rn.invTransform(75, 37.5, nil)
	if x != 62.5 || y != 25 {
		t.Fatalf("inverse point = (%v,%v), want (62.5,25)", x, y)
	}
}
