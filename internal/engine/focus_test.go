package engine

import (
	"testing"
)

func TestFocusManagerRegister(t *testing.T) {
	fm := NewFocusManager()
	node := &FocusNode{CanFocus: true, Focusable: true}
	fm.Register(node)

	if len(fm.nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(fm.nodes))
	}
}

func TestFocusManagerFocus(t *testing.T) {
	fm := NewFocusManager()
	var focused, blurred bool
	node := &FocusNode{
		CanFocus:  true,
		Focusable: true,
		OnFocus:   func() { focused = true },
		OnBlur:    func() { blurred = true },
	}
	fm.Register(node)

	fm.Focus(node)
	if !focused {
		t.Error("expected OnFocus called")
	}
	if fm.GetFocused() != node {
		t.Error("expected node to be focused")
	}

	fm.Unfocus()
	if !blurred {
		t.Error("expected OnBlur called")
	}
	if fm.GetFocused() != nil {
		t.Error("expected no focused node")
	}
}

func TestFocusManagerNextFocus(t *testing.T) {
	fm := NewFocusManager()
	var focusOrder []int
	n1 := &FocusNode{CanFocus: true, Focusable: true, OnFocus: func() { focusOrder = append(focusOrder, 1) }}
	n2 := &FocusNode{CanFocus: true, Focusable: true, OnFocus: func() { focusOrder = append(focusOrder, 2) }}
	n3 := &FocusNode{CanFocus: true, Focusable: true, OnFocus: func() { focusOrder = append(focusOrder, 3) }}
	fm.Register(n1)
	fm.Register(n2)
	fm.Register(n3)

	fm.NextFocus() // → n1
	fm.NextFocus() // → n2
	fm.NextFocus() // → n3
	fm.NextFocus() // → n1 (wrap)

	if len(focusOrder) != 4 {
		t.Fatalf("expected 4 focus events, got %d", len(focusOrder))
	}
	if focusOrder[0] != 1 || focusOrder[1] != 2 || focusOrder[2] != 3 || focusOrder[3] != 1 {
		t.Errorf("unexpected focus order: %v", focusOrder)
	}
}

func TestFocusManagerPreviousFocus(t *testing.T) {
	fm := NewFocusManager()
	var focusOrder []int
	n1 := &FocusNode{CanFocus: true, Focusable: true, OnFocus: func() { focusOrder = append(focusOrder, 1) }}
	n2 := &FocusNode{CanFocus: true, Focusable: true, OnFocus: func() { focusOrder = append(focusOrder, 2) }}
	fm.Register(n1)
	fm.Register(n2)

	fm.PreviousFocus() // → n2 (wrap from end)
	fm.PreviousFocus() // → n1

	if len(focusOrder) != 2 {
		t.Fatalf("expected 2 focus events, got %d", len(focusOrder))
	}
	if focusOrder[0] != 2 || focusOrder[1] != 1 {
		t.Errorf("unexpected focus order: %v", focusOrder)
	}
}

func TestFocusManagerSkipNonFocusable(t *testing.T) {
	fm := NewFocusManager()
	var focusOrder []int
	n1 := &FocusNode{CanFocus: true, Focusable: true, OnFocus: func() { focusOrder = append(focusOrder, 1) }}
	n2 := &FocusNode{CanFocus: true, Focusable: false, OnFocus: func() { focusOrder = append(focusOrder, 2) }} // 不可 Tab 到达
	n3 := &FocusNode{CanFocus: true, Focusable: true, OnFocus: func() { focusOrder = append(focusOrder, 3) }}
	fm.Register(n1)
	fm.Register(n2)
	fm.Register(n3)

	fm.NextFocus() // → n1
	fm.NextFocus() // → n3 (skip n2)
	fm.NextFocus() // → n1 (wrap)

	if len(focusOrder) != 3 {
		t.Fatalf("expected 3 focus events, got %d", len(focusOrder))
	}
	if focusOrder[0] != 1 || focusOrder[1] != 3 || focusOrder[2] != 1 {
		t.Errorf("unexpected focus order: %v", focusOrder)
	}
}

func TestFocusManagerUnregister(t *testing.T) {
	fm := NewFocusManager()
	n1 := &FocusNode{CanFocus: true, Focusable: true}
	n2 := &FocusNode{CanFocus: true, Focusable: true}
	fm.Register(n1)
	fm.Register(n2)

	fm.Focus(n1)
	fm.Unregister(n1)

	if fm.GetFocused() != nil {
		t.Error("expected unfocused after unregister")
	}

	fm.NextFocus()
	if fm.GetFocused() != n2 {
		t.Error("expected n2 to be focused")
	}
}

func TestFocusManagerClear(t *testing.T) {
	fm := NewFocusManager()
	fm.Register(&FocusNode{CanFocus: true, Focusable: true})
	fm.Register(&FocusNode{CanFocus: true, Focusable: true})
	fm.NextFocus()

	fm.Clear()
	if fm.GetFocused() != nil {
		t.Error("expected no focused node after clear")
	}
	if len(fm.nodes) != 0 {
		t.Errorf("expected 0 nodes after clear, got %d", len(fm.nodes))
	}
}

func TestFocusManagerNoNodes(t *testing.T) {
	fm := NewFocusManager()
	fm.NextFocus()    // 不应 panic
	fm.PreviousFocus() // 不应 panic
	if fm.GetFocused() != nil {
		t.Error("expected nil focused with no nodes")
	}
}
