package shadcn

import (
	"fmt"
	"testing"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// These tests drive real interactions through ui.Mount (the headless harness),
// asserting behavior — clicks fire, controlled state flips the rendered output —
// rather than merely that constructors return non-nil (see shadcn_test.go).

func TestButtonClickFires(t *testing.T) {
	clicks := 0
	h := ui.MountDefault(Button(ButtonProps{OnClick: func() { clicks++ }}, ui.Text("Go")))

	if !h.Root().ByText("Go").Click() {
		t.Fatal("Button click found no handler")
	}
	if clicks != 1 {
		t.Fatalf("clicks = %d, want 1", clicks)
	}
}

// TestButtonSizes locks button heights to the shadcn spec (h-9/h-8/h-10 = 36/32/40),
// verified through real layout via the harness.
func TestButtonSizes(t *testing.T) {
	cases := []struct {
		size Size
		w, h float32
	}{
		{SizeDefault, 0, 36},
		{SizeSm, 0, 32},
		{SizeLg, 0, 40},
		{SizeIcon, 36, 36},
	}
	for _, c := range cases {
		h := ui.MountDefault(Button(ButtonProps{Size: c.size}, ui.Text("x")))
		b := h.Root().Bounds()
		if b.H != c.h {
			t.Errorf("size %d height = %v, want %v", c.size, b.H, c.h)
		}
		if c.w != 0 && b.W != c.w {
			t.Errorf("size %d width = %v, want %v", c.size, b.W, c.w)
		}
	}
}

// TestButtonPaintsBackgroundAndLabel uses the recording painter (ui.Harness.Paint)
// to assert a Button actually renders a filled background + its label text —
// a headless golden-style paint check, beyond construction/layout.
func TestButtonPaintsBackgroundAndLabel(t *testing.T) {
	h := ui.MountDefault(Button(ButtonProps{}, ui.Text("Save")))
	ops := h.Paint()

	var rects, labels int
	for _, op := range ops {
		switch op.Kind {
		case "rect":
			if op.Rect.W > 0 && op.Rect.H > 0 {
				rects++
			}
		case "text":
			if op.Text == "Save" {
				labels++
			}
		}
	}
	if rects == 0 || labels != 1 {
		t.Fatalf("button paint: rects=%d label=%d; want >=1 / 1\n%+v", rects, labels, ops)
	}
}

// Tabs 用左右方向键在标签间移动焦点（ArrowNav）。
func TestTabsArrowNav(t *testing.T) {
	h := ui.MountDefault(Tabs(TabsProps{Tabs: []string{"A", "B", "C"}, Active: 0, OnChange: func(int) {}}))
	label := func(s string) *ui.Query {
		return h.Root().Find(func(q *ui.Query) bool { return q.Clickable() && q.AllText() == s })
	}
	label("A").Focus()
	if got := h.Arrow(ui.NavHorizontal, true).AllText(); got != "B" {
		t.Fatalf("Right -> %q want B", got)
	}
	if got := h.Arrow(ui.NavHorizontal, false).AllText(); got != "A" {
		t.Fatalf("Left -> %q want A", got)
	}
	if got := h.Arrow(ui.NavHorizontal, false).AllText(); got != "C" { // 环形回绕
		t.Fatalf("Left wrap -> %q want C", got)
	}
}

// 打开的 Dialog 应把键盘焦点困在对话框内（不 Tab 到背景按钮）。
func TestDialogTrapsFocus(t *testing.T) {
	app := func(_ struct{}) *ui.Node {
		noop := func() {}
		return ui.Div(
			Button(ButtonProps{OnClick: noop}, ui.Text("bg")),
			Dialog(DialogProps{Open: true},
				Button(ButtonProps{OnClick: noop}, ui.Text("ok")),
				Button(ButtonProps{OnClick: noop}, ui.Text("cancel")),
			),
		)
	}
	h := ui.MountDefault(ui.Use(app, struct{}{}))
	h.Step(200) // 让进场过渡完成，浮层挂载

	seen := map[string]bool{}
	for i := 0; i < 4; i++ {
		seen[h.Tab().AllText()] = true
	}
	if seen["bg"] {
		t.Fatalf("focus escaped dialog to background: %v", seen)
	}
	if !seen["ok"] || !seen["cancel"] {
		t.Fatalf("dialog buttons not reachable: %v", seen)
	}
}

func TestCheckboxTogglesControlledState(t *testing.T) {
	app := func(_ struct{}) *ui.Node {
		checked, set := ui.UseState(false)
		status := "off"
		if checked {
			status = "on"
		}
		return ui.Div(
			Checkbox(CheckboxProps{Checked: checked, OnChange: func(v bool) { set(v) }}),
			ui.Text(status),
		)
	}
	h := ui.MountDefault(ui.Use(app, struct{}{}))

	if !h.Root().ByText("off").Exists() {
		t.Fatalf("initial status texts = %v, want off", h.Root().Texts())
	}

	box := h.Root().Clickables()
	if len(box) == 0 {
		t.Fatal("checkbox is not clickable")
	}
	box[0].Click()

	if !h.Root().ByText("on").Exists() {
		t.Fatalf("after click texts = %v, want on", h.Root().Texts())
	}
}

func TestPopoverOpensAndEscapeCloses(t *testing.T) {
	h := ui.MountDefault(ui.Div(Popover(ui.Text("open"), ui.Text("content"))))

	if n := len(h.Overlays()); n != 0 {
		t.Fatalf("overlays before open = %d, want 0", n)
	}

	// Click the trigger — bubbles to the Popover's toggle.
	if !h.Root().ByText("open").Click() {
		t.Fatal("Popover trigger click found no handler")
	}
	overlays := h.Overlays()
	if len(overlays) != 1 {
		t.Fatalf("overlays after open = %d, want 1", len(overlays))
	}
	if !overlays[0].ByText("content").Exists() {
		t.Fatalf("popover content not shown; overlay texts = %v", overlays[0].Texts())
	}

	// Esc runs the topmost UseEscape handler, closing the popover.
	h.Escape()
	if n := len(h.Overlays()); n != 0 {
		t.Fatalf("overlays after Escape = %d, want 0", n)
	}
}

func TestTabsSwitchActive(t *testing.T) {
	app := func(_ struct{}) *ui.Node {
		active, set := ui.UseState(0)
		return ui.Div(
			Tabs(TabsProps{Tabs: []string{"A", "B", "C"}, Active: active, OnChange: func(i int) { set(i) }}),
			ui.Text(fmt.Sprintf("active=%d", active)),
		)
	}
	h := ui.MountDefault(ui.Use(app, struct{}{}))

	if !h.Root().ByText("active=0").Exists() {
		t.Fatalf("initial texts = %v, want active=0", h.Root().Texts())
	}

	// Click the "C" tab label; bubbling reaches the tab's onClick.
	if !h.Root().ByText("C").Click() {
		t.Fatal("tab C click found no handler")
	}
	if !h.Root().ByText("active=2").Exists() {
		t.Fatalf("after clicking C texts = %v, want active=2", h.Root().Texts())
	}
}
