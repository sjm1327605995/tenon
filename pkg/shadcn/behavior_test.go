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

// Phase-2 表单组件构造并渲染。
func TestFormComponentsRender(t *testing.T) {
	clicks := 0
	h := ui.MountDefault(ui.Use(func(_ struct{}) *ui.Node {
		return ui.Div(
			Field(FieldProps{Label: "Email", Description: "desc-here"},
				Input(InputProps{Placeholder: "x@y.com"})),
			Field(FieldProps{Label: "Pass", Error: "err-here"},
				Input(InputProps{Value: "v"})),
			InputGroup(InputGroupProps{Leading: ui.Icon(ui.IconSearch, 16), Placeholder: "search"}),
			ButtonGroup([]ButtonGroupItem{
				{Label: "Day", Active: true}, {Label: "Week", OnClick: func() { clicks++ }},
			}),
		)
	}, struct{}{}))
	for _, s := range []string{"Email", "desc-here", "Pass", "err-here", "Day", "Week"} {
		if !h.Root().ByText(s).Exists() {
			t.Fatalf("missing %q; texts=%v", s, h.Root().Texts())
		}
	}
	if !h.Root().ByText("Week").Click() || clicks != 1 {
		t.Fatalf("ButtonGroup item click did not fire (clicks=%d)", clicks)
	}
}

// Item + InputOTP：渲染与交互。
func TestItemAndOTP(t *testing.T) {
	clicked := 0
	h := ui.MountDefault(ui.Use(func(_ struct{}) *ui.Node {
		return ui.Div(
			Item(ItemProps{Title: "T1", Description: "D1", OnClick: func() { clicked++ }}),
			InputOTP(InputOTPProps{Length: 6, Value: "42x8"}), // 非数字被过滤 -> 428
		)
	}, struct{}{}))
	for _, s := range []string{"T1", "D1", "4", "2", "8"} {
		if !h.Root().ByText(s).Exists() {
			t.Fatalf("missing %q; texts=%v", s, h.Root().Texts())
		}
	}
	if h.Root().ByText("x").Exists() {
		t.Fatal("OTP should have filtered non-digit 'x'")
	}
	if !h.Root().ByText("T1").Click() || clicked != 1 {
		t.Fatalf("Item click didn't fire (clicked=%d)", clicked)
	}
}

func TestDigitsOnly(t *testing.T) {
	if got := digitsOnly("a1b2c3!"); got != "123" {
		t.Fatalf("digitsOnly = %q want 123", got)
	}
}

// Accordion：标题与内容都渲染，点击标题触发切换（不 panic）。
func TestAccordionToggle(t *testing.T) {
	h := ui.MountDefault(Accordion([]AccordionItemData{
		{Title: "A", Content: []*ui.Node{ui.Text("body-A")}},
		{Title: "B", Content: []*ui.Node{ui.Text("body-B")}},
	}))
	for _, s := range []string{"A", "B", "body-A", "body-B"} {
		if !h.Root().ByText(s).Exists() {
			t.Fatalf("missing %q; texts=%v", s, h.Root().Texts())
		}
	}
	if !h.Root().ByText("B").Click() {
		t.Fatal("clicking accordion header fired no toggle handler")
	}
	h.Step(250) // 推进滑动/旋转动画，确认无 panic
}

// Phase-1 基础组件都能构造并渲染出文本。
func TestPrimitivesRender(t *testing.T) {
	h := ui.MountDefault(ui.Use(func(_ struct{}) *ui.Node {
		return ui.Div(
			Kbd("K"),
			KbdGroup("Ctrl", "K"),
			Spinner(SpinnerProps{}),
			Empty(EmptyProps{Title: "空", Description: "无数据"}),
			H1("标题"), H2("H2"), P("正文"), Lead("引导"), Muted("次要"),
			InlineCode("go"), Blockquote("引用"), Large("大"), Small("小"),
		)
	}, struct{}{}))
	for _, s := range []string{"K", "Ctrl", "空", "无数据", "标题", "H2", "正文", "引导", "次要", "go", "引用", "大", "小"} {
		if !h.Root().ByText(s).Exists() {
			t.Fatalf("missing rendered text %q; texts=%v", s, h.Root().Texts())
		}
	}
}

// Combobox：打开 -> 输入过滤 -> 选中回填并关闭。
func TestComboboxFilterAndSelect(t *testing.T) {
	var selected string
	app := func(_ struct{}) *ui.Node {
		val, setVal := ui.UseState("")
		selected = val
		return Combobox(ComboboxProps{
			Options: []ComboboxOption{
				{Value: "go", Label: "Go"},
				{Value: "rs", Label: "Rust"},
				{Value: "py", Label: "Python"},
			},
			Value:       val,
			OnChange:    func(v string) { setVal(v) },
			Placeholder: "Pick language",
		})
	}
	h := ui.MountDefault(ui.Use(app, struct{}{}))

	if !h.Root().ByText("Pick language").Exists() {
		t.Fatal("placeholder not shown initially")
	}
	if len(h.Overlays()) != 0 {
		t.Fatal("panel should be closed initially")
	}

	h.Root().ByText("Pick language").Click() // 打开
	if len(h.Overlays()) == 0 {
		t.Fatal("panel did not open on trigger click")
	}
	ov := h.Overlays()[0]
	for _, lbl := range []string{"Go", "Rust", "Python"} {
		if !ov.ByText(lbl).Exists() {
			t.Fatalf("option %q missing before filter (texts=%v)", lbl, ov.Texts())
		}
	}

	ov.ByKind("input").Type("ru") // 过滤
	ov = h.Overlays()[0]
	if ov.ByText("Go").Exists() || ov.ByText("Python").Exists() {
		t.Fatalf("filter did not hide non-matches: %v", ov.Texts())
	}
	if !ov.ByText("Rust").Exists() {
		t.Fatalf("filter hid the match: %v", ov.Texts())
	}

	ov.ByText("Rust").Click() // 选中
	if selected != "rs" {
		t.Fatalf("OnChange value = %q want rs", selected)
	}
	if len(h.Overlays()) != 0 {
		t.Fatal("panel should close after select")
	}
	if !h.Root().ByText("Rust").Exists() {
		t.Fatalf("trigger should show selected label; texts=%v", h.Root().Texts())
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
