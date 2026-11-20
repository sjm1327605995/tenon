package ui

import "testing"

func TestHitCacheReusesAndInvalidates(t *testing.T) {
	tree := Div(Style(Width(100), Height(100)), OnHover(func(bool) {}), Text("hit"))
	h := MountDefault(tree)
	defer h.Close()
	old := input
	input = &harnessInput{x: 20, y: 20}
	defer func() { input = old }()

	h.g.updateHover()
	if got := h.g.hitTests; got != 1 {
		t.Fatalf("initial hit tests=%d, want 1", got)
	}
	h.g.hasLayoutAnim = true
	h.g.updateHover()
	if got := h.g.hitTests; got != 1 {
		t.Fatalf("unchanged tree forced another hit test: %d", got)
	}

	input = &harnessInput{x: 30, y: 20}
	h.g.updateHover()
	if got := h.g.hitTests; got != 2 {
		t.Fatalf("cursor movement did not refresh hit test: %d", got)
	}
	h.g.layout()
	h.g.updateHover()
	if got := h.g.hitTests; got != 3 {
		t.Fatalf("layout change did not invalidate hit cache: %d", got)
	}
}
