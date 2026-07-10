package ui

import (
	"fmt"
	"testing"
)

// counterApp is a minimal stateful component used to exercise the harness.
func counterApp(_ struct{}) *Node {
	count, set := UseState(0)
	return Div(Style(Row, Gap(8)),
		Button(Style(Width(40), Height(24)), OnClick(func() { set(count - 1) }), Text("-")),
		Text(fmt.Sprintf("%d", count), FontSize(16)),
		Button(Style(Width(40), Height(24)), OnClick(func() { set(count + 1) }), Text("+")),
	)
}

func TestHarnessClickDrivesState(t *testing.T) {
	h := Mount(Use(counterApp, struct{}{}), 300, 100)

	if got := h.Root().ByText("0"); !got.Exists() {
		t.Fatalf("initial texts = %v, want a 0", h.Root().Texts())
	}

	// Click "+" twice, "-" once -> 1.
	h.Root().ByText("+").Click()
	h.Root().ByText("+").Click()
	h.Root().ByText("-").Click()

	if got := h.Root().ByText("1"); !got.Exists() {
		t.Fatalf("after +,+,- texts = %v, want a 1", h.Root().Texts())
	}
}

func TestHarnessClickAtHitsThroughOverlap(t *testing.T) {
	fired := ""
	app := Div(Style(Width(200), Height(100)),
		Button(Style(Width(200), Height(100)), OnClick(func() { fired = "bg" }),
			Button(Style(Absolute, Left(10), Top(10), Width(50), Height(30)),
				OnClick(func() { fired = "fg" }), Text("x")),
		),
	)
	h := Mount(app, 200, 100)

	// A point inside the inner button should hit the inner handler, not the outer.
	if !h.ClickAt(20, 20) {
		t.Fatal("ClickAt found no handler")
	}
	if fired != "fg" {
		t.Fatalf("clicked handler = %q, want fg (nearest onClick wins)", fired)
	}

	// A point only over the outer button hits the outer handler.
	fired = ""
	h.ClickAt(150, 80)
	if fired != "bg" {
		t.Fatalf("clicked handler = %q, want bg", fired)
	}
}

func TestHarnessInputTyping(t *testing.T) {
	var current string
	app := func(_ struct{}) *Node {
		val, set := UseState("")
		current = val
		return Input(Style(Width(200), Height(30), FontSize(16)),
			Value(val), OnChange(func(s string) { set(s) }))
	}
	h := Mount(Use(app, struct{}{}), 240, 60)

	in := h.Root().ByKind("input")
	if !in.Exists() {
		t.Fatal("no input node found")
	}
	in.Focus().Type("hi")
	// Re-query: controlled input re-renders on change.
	if v := h.Root().ByKind("input").Value(); v != "hi" {
		t.Fatalf("input value = %q, want hi", v)
	}
	if current != "hi" {
		t.Fatalf("state after typing = %q, want hi", current)
	}
}

func TestHarnessHover(t *testing.T) {
	h := Mount(Div(Style(Width(50), Height(50)),
		OnHover(func(bool) {})), 100, 100)
	// The engine wires a hover handler; toggling must not panic and must settle.
	q := h.Root()
	if q.rn.onHover == nil {
		t.Fatal("onHover not wired")
	}
	q.Hover(true)
	q.Hover(false)
}

func TestHarnessKeyboardFocus(t *testing.T) {
	activated := ""
	app := Div(Style(Column),
		Button(Style(Width(50), Height(20)), OnClick(func() { activated = "a" }), Text("a")),
		Button(Style(Width(50), Height(20)), OnClick(func() { activated = "b" }), Text("b")),
	)
	h := Mount(app, 200, 200)

	if h.Focused().Exists() {
		t.Fatal("nothing should be focused before Tab")
	}
	if got := h.Tab().AllText(); got != "a" { // first focusable
		t.Fatalf("Tab #1 focused %q, want a", got)
	}
	if !h.Focused().IsFocused() {
		t.Fatal("Focused() node should report IsFocused")
	}
	if got := h.Tab().AllText(); got != "b" {
		t.Fatalf("Tab #2 focused %q, want b", got)
	}
	if got := h.Tab().AllText(); got != "a" { // wraps around
		t.Fatalf("Tab #3 focused %q, want a (wrap)", got)
	}
	if got := h.ShiftTab().AllText(); got != "b" { // backward wraps
		t.Fatalf("ShiftTab focused %q, want b", got)
	}

	// Enter activates the focused (button "b").
	if !h.Enter() {
		t.Fatal("Enter did not activate focused button")
	}
	if activated != "b" {
		t.Fatalf("activated = %q, want b", activated)
	}
}

func TestHarnessEscape(t *testing.T) {
	fired := 0
	app := func(_ struct{}) *Node {
		UseEscape(true, func() { fired++ })
		return Text("x")
	}
	h := Mount(Use(app, struct{}{}), 100, 100)

	h.Escape()
	if fired != 1 {
		t.Fatalf("escape handler fired=%d, want 1", fired)
	}
}

func TestHarnessScroll(t *testing.T) {
	kids := []*Node{Style(Height(100), Width(120), Column)}
	for i := 0; i < 5; i++ {
		kids = append(kids, Div(Style(Height(40), Width(100)), Text(fmt.Sprintf("row%d", i))))
	}
	h := Mount(ScrollView(kids...), 140, 100)

	sc := h.Root()
	if sc.Kind() != "scroll" {
		t.Fatalf("root kind = %s, want scroll", sc.Kind())
	}
	y0 := sc.Child(0).Bounds().Y

	if !sc.ScrollBy(50) {
		t.Fatal("ScrollBy found no scroll container")
	}
	if got := h.Root().Child(0).Bounds().Y; got != y0-50 {
		t.Fatalf("after ScrollBy(50) row0.Y=%v, want %v", got, y0-50)
	}

	// Over-scroll clamps to content height (200) - viewport (100) = 100.
	h.Root().ScrollBy(999)
	if got := h.Root().Child(0).Bounds().Y; got != y0-100 {
		t.Fatalf("over-scroll row0.Y=%v, want clamped %v", got, y0-100)
	}
}

func TestHarnessInputEditing(t *testing.T) {
	var cur string
	app := func(_ struct{}) *Node {
		v, set := UseState("hello")
		cur = v
		return Input(Style(Width(200), Height(30), FontSize(16)),
			Value(v), OnChange(func(s string) { set(s) }))
	}
	h := Mount(Use(app, struct{}{}), 240, 60)

	h.Root().ByKind("input").Focus().Backspace(2) // "hello" -> "hel"
	if v := h.Root().ByKind("input").Value(); v != "hel" {
		t.Fatalf("after Backspace(2) value=%q, want hel", v)
	}
	h.Root().ByKind("input").Clear()
	if cur != "" {
		t.Fatalf("after Clear state=%q, want empty", cur)
	}
}

func TestHarnessMissesAreSafe(t *testing.T) {
	h := Mount(Div(Text("only")), 100, 100)
	miss := h.Root().ByText("nope")
	if miss.Exists() {
		t.Fatal("ByText should miss")
	}
	// Read + action methods on a missing node must be no-ops, not panics.
	if miss.Text() != "" || miss.Click() || miss.Type("x") || miss.Bounds() != (Rect{}) {
		t.Fatal("methods on a missing Query should be safe no-ops")
	}
}
