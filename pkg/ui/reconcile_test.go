package ui

import "testing"

func TestUnkeyedReconcileSkipsNilAndReusesFibers(t *testing.T) {
	parent := &Fiber{}
	old := mountList(parent, []*Node{nil, Div(Text("a")), nil, Div(Text("b"))})
	first, second := old[0], old[1]

	old = reconcileList(parent, old, []*Node{Div(Text("a2")), nil, Div(Text("b2")), Div(Text("c"))})
	if len(old) != 3 || old[0] != first || old[1] != second {
		t.Fatal("unkeyed reconcile did not preserve positional fiber identity after nil")
	}
	old = reconcileList(parent, old, []*Node{nil, Div(Text("a3"))})
	if len(old) != 1 || old[0] != first {
		t.Fatal("unkeyed reconcile did not preserve the remaining fiber")
	}
	if !second.unmounted {
		t.Fatal("removed unkeyed fiber was not unmounted")
	}
}
