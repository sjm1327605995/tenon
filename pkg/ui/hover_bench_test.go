package ui

import "testing"

func BenchmarkUpdateHover(b *testing.B) {
	tree := Div(Style(Width(40), Height(40)), OnHover(func(bool) {}))
	tree = Div(Style(Width(60), Height(60)), OnHover(func(bool) {}), tree)
	tree = Div(Style(Width(80), Height(80)), OnHover(func(bool) {}), tree)
	tree = Div(Style(Width(100), Height(100)), OnHover(func(bool) {}), tree)
	h := MountDefault(tree)
	defer h.Close()
	h.MoveTo(20, 20)
	h.g.hasLayoutAnim = true
	old := input
	input = &harnessInput{x: 20, y: 20}
	defer func() { input = old }()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.g.updateHover()
	}
}
