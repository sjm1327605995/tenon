package ui

import (
	"fmt"
	"testing"
)

func TestKeyboardFocusScrollsIntoView(t *testing.T) {
	for _, portal := range []bool{false, true} {
		t.Run(fmt.Sprintf("portal=%v", portal), func(t *testing.T) {
			list := ScrollView(Style(Width(180), Height(90)), ArrowNav(NavVertical),
				Button(Style(Height(30)), OnClick(func() {}), Text("First")),
				Div(Style(Height(200))),
				Button(Style(Height(30)), OnClick(func() {}), Text("Last")))
			if portal {
				list = Portal(list)
			}
			h := Mount(Div(Style(Fill), list), 320, 240)
			root := func() *Query {
				if portal {
					return h.Overlays()[0]
				}
				return h.Root()
			}
			for _, arrow := range []bool{false, true} {
				root().Find(func(q *Query) bool { return q.Clickable() && q.AllText() == "First" }).Focus()
				if arrow {
					h.Arrow(NavVertical, true)
				} else {
					h.Tab()
				}
				b, viewport := h.Focused().Bounds(), root().ByKind("scroll").Bounds()
				if h.Focused().AllText() != "Last" || b.Y < viewport.Y || b.Y+b.H > viewport.Y+viewport.H+0.1 {
					t.Fatalf("keyboard focus remains outside scroll viewport (arrow=%v): focus=%+v viewport=%+v", arrow, b, viewport)
				}
				if arrow {
					h.Arrow(NavVertical, true) // wrap to the first item
				} else {
					h.ShiftTab()
				}
				b = h.Focused().Bounds()
				if h.Focused().AllText() != "First" || b.Y < viewport.Y || b.Y+b.H > viewport.Y+viewport.H+0.1 {
					t.Fatalf("returning focus did not scroll back: %+v", b)
				}
			}
		})
	}
}

func TestKeyboardFocusRevealsTargetInNestedScrollViews(t *testing.T) {
	h := MountDefault(ScrollView(Style(Width(180), Height(100)),
		Div(Style(Height(100))),
		ScrollView(Style(Height(80)),
			Div(Style(Height(160))),
			Button(Style(Height(30)), OnClick(func() {}), Text("Target")))))
	h.Tab()
	b := h.Focused().Bounds()
	for _, scroll := range h.Root().FindAll(func(q *Query) bool { return q.Kind() == "scroll" }) {
		v := scroll.Bounds()
		if b.Y < v.Y || b.Y+b.H > v.Y+v.H+0.1 {
			t.Fatalf("nested scroll did not reveal focused target: target=%+v viewport=%+v", b, v)
		}
	}
}

func TestTabIntoInputCollapsesSelectionAtEnd(t *testing.T) {
	h := MountDefault(Use(func(_ struct{}) *Node {
		value, setValue := UseState("hello")
		return Input(Value(value), OnChange(setValue))
	}, struct{}{}))
	h.pressKey(keyTab)
	q := h.Focused()
	if q.rn.caretPos != 5 || q.rn.selAnchor != 5 {
		t.Fatalf("Tab unexpectedly selected the input text: caret=%d anchor=%d", q.rn.caretPos, q.rn.selAnchor)
	}
	old := input
	input = &harnessInput{chars: []rune("!")}
	t.Cleanup(func() { input = old })
	h.g.editFocusedInput()
	h.settle()
	if got := h.Root().ByKind("input").Value(); got != "hello!" {
		t.Fatalf("typing after Tab replaced existing text: %q", got)
	}
}

func TestMouseDownFocusesClickableAncestor(t *testing.T) {
	clicks := 0
	h := MountDefault(Button(OnClick(func() { clicks++ }), Text("Save")))
	b := h.Root().ByText("Save").Bounds()
	focus := h.MouseDown(b.X+b.W/2, b.Y+b.H/2)
	if !focus.Clickable() || focus.AllText() != "Save" || clicks != 1 {
		t.Fatal("clicking button text did not focus and activate its button")
	}
	if !h.Enter() || clicks != 2 {
		t.Fatal("clicked button could not subsequently activate through the keyboard")
	}
}
