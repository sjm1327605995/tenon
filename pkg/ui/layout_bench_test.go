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

func BenchmarkResolveInherited(b *testing.B) {
	h := MountDefault(layoutBenchTree(200))
	defer h.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resolveInherited(h.g.rootRN, inhText{})
	}
}
