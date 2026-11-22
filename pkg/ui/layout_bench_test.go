package ui

import "testing"

func layoutBenchTree(rows int) *Node {
	children := make([]*Node, 0, rows)
	for i := 0; i < rows; i++ {
		children = append(children, Div(
			Style(Row, Gap(8)),
			Text("label"),
			Div(Style(Width(40), Height(20)), Text("value")),
		))
	}
	return Div(append([]*Node{Style(Column, Gap(4))}, children...)...)
}

func BenchmarkRelink(b *testing.B) {
	h := MountDefault(layoutBenchTree(200))
	defer h.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		relink(h.g.rootFiber)
	}
}

func BenchmarkLayoutCleanTree(b *testing.B) {
	h := MountDefault(layoutBenchTree(200))
	defer h.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.g.layout()
	}
}

func BenchmarkResolveInherited(b *testing.B) {
	h := MountDefault(layoutBenchTree(200))
	defer h.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resolveInherited(h.g.rootRN, inhText{})
	}
}

func BenchmarkResolveInheritedDirtySubtree(b *testing.B) {
	h := MountDefault(layoutBenchTree(200))
	defer h.Close()
	branch := h.g.rootRN.children[0]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		branch.markInheritedDirty()
		resolveInherited(h.g.rootRN, inhText{})
	}
}
