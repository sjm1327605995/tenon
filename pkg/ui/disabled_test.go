package ui

import (
	"testing"

	"gioui.org/io/key"
	"gioui.org/io/pointer"
)

func TestDisabledInputRejectsKeyboardIMEAndPaste(t *testing.T) {
	changes, submits := 0, 0
	h := MountDefault(Input(Disabled(true), Value("locked"),
		OnChange(func(string) { changes++ }), OnSubmit(func(string) { submits++ })))
	rn := h.Root().rn
	// 即使后端还持有禁用前的焦点/异步编辑，落地时也必须重新检查禁用态。
	h.g.focusedFiber = rn.owner
	rn.caretPos, rn.selAnchor = len(rn.value), 0
	old := input
	hi := &harnessInput{chars: []rune("x")}
	hi.just[keyEnter], hi.just[keyBackspace] = true, true
	input = hi
	t.Cleanup(func() { input = old })
	h.g.handleKeyboardNav()
	h.g.editFocusedInput()
	pasteIntoFocused(h.g, "paste")
	var ime gioIMEState
	ime.applyEdits(h.g, &gioInput{edits: []gioEdit{{rng: key.Range{Start: 0, End: 6}, text: "替换"}}})
	if rn.value != "locked" || changes != 0 || submits != 0 {
		t.Fatalf("disabled input accepted an edit: value=%q changes=%d submits=%d", rn.value, changes, submits)
	}
	gioIn.setCursor(rn.bounds.X+1, rn.bounds.Y+1)
	if got := gioCursor(h.g); got != pointer.CursorDefault {
		t.Fatalf("disabled input cursor=%v, want default", got)
	}
}

func TestDisablingFocusedInputClearsInteractionAndCanReenable(t *testing.T) {
	var setDisabled func(bool)
	changes := 0
	h := MountDefault(Use(func(_ struct{}) *Node {
		disabled, set := UseState(false)
		setDisabled = set
		value, setValue := UseState("hello")
		return Input(Disabled(disabled), Value(value), OnChange(func(s string) { changes++; setValue(s) }),
			Style(Width(180), Height(36)))
	}, struct{}{}))
	q := h.Root().ByKind("input").Focus()
	h.g.inputSelecting, h.g.imeComposing = true, true
	q.rn.composeLo, q.rn.composeHi = 0, 2
	q.rn.selAnchor = 1
	setDisabled(true)
	h.Step(0)
	if h.Focused().Exists() || h.g.inputSelecting || h.g.imeComposing || q.rn.composeHi != 0 || q.rn.selAnchor != q.rn.caretPos {
		t.Fatal("disabling input retained focus, selection capture, or composition")
	}
	setDisabled(false)
	h.Step(0)
	q = h.Root().ByKind("input")
	b := q.Bounds()
	if !h.MouseDown(b.X+b.W-4, b.Y+b.H/2).Exists() || !q.Type("!") || q.Value() != "hello!" || changes != 1 {
		t.Fatalf("reenabled input did not resume editing: value=%q changes=%d", q.Value(), changes)
	}
}

func TestDisabledSubtreeBlocksActivationButAllowsOuterTooltip(t *testing.T) {
	activated, hovered := 0, false
	h := MountDefault(Div(Style(Column), OnHover(func(on bool) { hovered = on }),
		Div(Disabled(true), Style(Width(180), Height(50)),
			Button(Style(Width(180), Height(50)), Text("Blocked"),
				OnClick(func() { activated++ }), OnPress(func(bool) { activated++ }),
				OnDrag(func(float32, float32) { activated++ }), OnContextMenu(func(float32, float32) { activated++ }))),
		Input(Placeholder("enabled"))))
	old := input
	hi := &harnessInput{x: 20, y: 20}
	hi.btnJust[btnLeft], hi.btnJust[btnRight], hi.btnHeld[btnLeft] = true, true, true
	input = hi
	t.Cleanup(func() { input = old })
	h.g.handleInput()
	if activated != 0 || h.Focused().Exists() || !hovered {
		t.Fatalf("disabled subtree: activations=%d focused=%v outer-hover=%v", activated, h.Focused().Exists(), hovered)
	}
	if got := h.Tab().Placeholder(); got != "enabled" {
		t.Fatalf("Tab entered disabled subtree: %q", got)
	}
}

func TestDisablingPressedControlReleasesCapture(t *testing.T) {
	var setDisabled func(bool)
	released := 0
	h := MountDefault(Use(func(_ struct{}) *Node {
		disabled, set := UseState(false)
		setDisabled = set
		return Button(Disabled(disabled), Style(Width(100), Height(36)),
			OnClick(func() {}), OnPress(func(down bool) {
				if !down {
					released++
				}
			}),
			OnDrag(func(float32, float32) {}))
	}, struct{}{}))
	rn := h.Root().rn
	h.g.pressedNode, h.g.dragging = rn, rn
	setDisabled(true)
	h.Step(0)
	if h.g.pressedNode != nil || h.g.dragging != nil || released != 1 {
		t.Fatalf("disabled control retained capture: pressed=%v dragging=%v release=%d", h.g.pressedNode != nil, h.g.dragging != nil, released)
	}
}

func TestDisabledBoundInputStillReceivesExternalUpdates(t *testing.T) {
	sig := NewSignal("initial")
	h := MountDefault(Input(Disabled(true), BindInput(sig)))
	q := h.Root().ByKind("input")
	if q.Type("user edit") || sig.Get() != "initial" {
		t.Fatal("disabled input wrote into its signal")
	}
	sig.Set("external update")
	h.Step(0)
	if q.Value() != "external update" {
		t.Fatalf("disabled input failed to display an external update: %q", q.Value())
	}
}
