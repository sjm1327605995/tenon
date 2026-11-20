package ui

import "testing"

func TestChildBitmapGrowsAndReusesWords(t *testing.T) {
	var bitmap childBitmap
	for _, index := range []uint32{0, 63, 64, 255, 1000} {
		bitmap.Set(index)
		if !bitmap.Has(index) {
			t.Fatalf("index %d is not set", index)
		}
	}
	if len(bitmap.words) != 16 {
		t.Fatalf("bitmap words=%d, want power-of-two capacity 16", len(bitmap.words))
	}
	bitmap.Delete(64)
	if bitmap.Has(64) {
		t.Fatal("index 64 was not deleted")
	}
	bitmap.Reset()
	if bitmap.Any() || len(bitmap.words) != 16 {
		t.Fatal("Reset did not retain cleared storage")
	}
}

func TestHoverBitmapRetainsCurrentState(t *testing.T) {
	enters := 0
	h := MountDefault(Div(Style(Width(200), Height(100)), OnHover(func(in bool) {
		if in {
			enters++
		}
	})))
	defer h.Close()
	old := input
	input = &harnessInput{x: 100, y: 50}
	defer func() { input = old }()
	h.g.updateHover()
	rn := h.g.rootRN
	input = &harnessInput{x: 101, y: 50}
	h.g.updateHover()
	if rn.hover.self != [2]bool{true, true} {
		t.Fatalf("hover state did not retain current and previous bits: %v", rn.hover.self)
	}
	if enters != 1 {
		t.Fatalf("enters=%d, want 1 after moving inside", enters)
	}
}

func TestRenderNodeLevelsAndLevelLookup(t *testing.T) {
	h := MountDefault(Div(Style(Width(100), Height(100)), Div(Text("leaf"))))
	defer h.Close()
	root := h.g.rootRN
	child := root.children[0]
	grandchild := child.children[0]
	if root.level != 0 || child.level != 1 || grandchild.level != 2 || h.g.maxLevel != 2 {
		t.Fatalf("unexpected levels: root=%d child=%d grandchild=%d max=%d",
			root.level, child.level, grandchild.level, h.g.maxLevel)
	}
	var atLevelOne []*renderNode
	root.forEachAtLevel(1, func(n *renderNode) { atLevelOne = append(atLevelOne, n) })
	if len(atLevelOne) != 1 || atLevelOne[0] != child {
		t.Fatal("level lookup did not return the direct child")
	}
}
