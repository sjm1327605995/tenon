package ui

import (
	"fmt"
	"testing"
)

func syntaxFirstView() *Node {
	count, set := UseState(1)
	return Button(Text(fmt.Sprint(count))).OnClick(func() { set(count + 1) })
}

func syntaxSecondView() *Node {
	count, _ := UseState(99)
	return Text(fmt.Sprint(count))
}

func TestViewUsesOriginalFunctionIdentity(t *testing.T) {
	h := MountDefault(View(func() *Node {
		second, set := UseState(false)
		view := View(syntaxFirstView)
		if second {
			view = View(syntaxSecondView)
		}
		return VStack(0, view, Button(Text("Switch")).OnClick(func() { set(!second) }))
	}))
	h.Root().ByText("1").Click()
	if !h.Root().ByText("2").Exists() {
		t.Fatal("View did not preserve state during its local rerender")
	}
	h.Root().ByText("Switch").Click()
	if !h.Root().ByText("99").Exists() {
		t.Fatal("distinct View functions shared hook state")
	}
}

func TestForPreservesKeyedRowStateDuringReorder(t *testing.T) {
	row := func(id string) *Node {
		count, set := UseState(0)
		return Button(Text(fmt.Sprintf("%s:%d", id, count))).OnClick(func() { set(count + 1) })
	}
	h := MountDefault(View(func() *Node {
		reversed, set := UseState(false)
		items := []string{"a", "b"}
		if reversed {
			items = []string{"b", "a"}
		}
		return VStack(0, Text("Header"), For(items, func(id string) *Node {
			return Use(row, id).Key(id)
		}), Button(Text("Reverse")).OnClick(func() { set(!reversed) }))
	}))
	h.Root().ByText("a:0").Click()
	h.Root().ByText("Reverse").Click()
	texts := h.Root().Texts()
	if len(texts) != 4 || texts[0] != "Header" || texts[1] != "b:0" || texts[2] != "a:1" {
		t.Fatalf("For lost keyed state or order: %v", texts)
	}
	if h.Root().Count() != 4 {
		t.Fatal("For inserted an extra layout container")
	}
}

func TestForHandlesEmptyAndFilteredResults(t *testing.T) {
	h := MountDefault(VStack(0,
		For([]int(nil), func(int) *Node { t.Fatal("empty list called renderer"); return nil }),
		For([]int{0, 1, 2}, func(v int) *Node {
			if v == 1 {
				return Text("One")
			}
			return nil
		})))
	if h.Root().Count() != 1 || !h.Root().ByText("One").Exists() {
		t.Fatalf("For did not ignore empty/filtered results: %v", h.Root().Texts())
	}
}
