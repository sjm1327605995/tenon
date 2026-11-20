package ui

import "testing"

func BenchmarkReconcileListUnkeyed(b *testing.B) {
	parent := &Fiber{}
	nodes := make([]*Node, 200)
	for i := range nodes {
		nodes[i] = Div(Style(Width(20), Height(20)))
	}
	old := mountList(parent, nodes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		old = reconcileList(parent, old, nodes)
	}
}

func BenchmarkReconcileListUnkeyedWithText(b *testing.B) {
	parent := &Fiber{}
	nodes := make([]*Node, 200)
	for i := range nodes {
		nodes[i] = Div(Style(Width(20), Height(20)), Text("unchanged"))
	}
	old := mountList(parent, nodes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		old = reconcileList(parent, old, nodes)
	}
}
