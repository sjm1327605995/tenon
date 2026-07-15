package shadcn

import (
	"fmt"
	"testing"
	"time"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// Phase-3: Menubar + DatePicker 内联渲染。
func TestMenubarAndDatePicker(t *testing.T) {
	h := ui.MountDefault(ui.Use(func(_ struct{}) *ui.Node {
		return ui.Div(
			Menubar([]MenubarMenu{{Label: "File", Items: []MenuItem{{Label: "New"}}}, {Label: "Edit"}}),
			DatePicker(DatePickerProps{Value: time.Date(2026, 7, 14, 0, 0, 0, 0, time.UTC)}),
		)
	}, struct{}{}))
	for _, s := range []string{"File", "Edit", "2026-07-14"} {
		if !h.Root().ByText(s).Exists() {
			t.Fatalf("missing %q; texts=%v", s, h.Root().Texts())
		}
	}
}

// Phase-3: AlertDialog 打开后显示内容，点击操作按钮回调。
func TestAlertDialogOpens(t *testing.T) {
	acted := 0
	h := ui.MountDefault(AlertDialog(AlertDialogProps{
		Open: true, Title: "T-title", Description: "D-desc",
		ActionLabel: "Go", OnAction: func() { acted++ }}))
	h.Step(220) // 过渡挂载
	ovs := h.Overlays()
	if len(ovs) == 0 {
		t.Fatal("alert dialog overlay did not open")
	}
	fired := false
	for _, ov := range ovs {
		if ov.ByText("T-title").Exists() && ov.ByText("D-desc").Exists() {
			fired = ov.ByText("Go").Click()
		}
	}
	if !fired || acted != 1 {
		t.Fatalf("alert action click didn't fire (fired=%v acted=%d)", fired, acted)
	}
}

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

// Phase-3: ContextMenu 右键弹出，选中回调并关闭。
func TestContextMenuRightClick(t *testing.T) {
	sel := ""
	h := ui.MountDefault(ui.Use(func(_ struct{}) *ui.Node {
		return ContextMenu(
			ui.Div(ui.Style(ui.Width(200), ui.Height(100)), ui.Text("area")),
			[]MenuItem{
				{Label: "Copy", OnSelect: func() { sel = "Copy" }},
				{Label: "Delete", OnSelect: func() { sel = "Delete" }},
			},
		)
	}, struct{}{}))
	if len(h.Overlays()) != 0 {
		t.Fatal("menu should be closed initially")
	}
	b := h.Root().ByText("area").Bounds()
	if !h.RightClickAt(b.X+b.W/2, b.Y+b.H/2) {
		t.Fatal("right-click found no onContextMenu handler")
	}
	ovs := h.Overlays()
	if len(ovs) == 0 {
		t.Fatal("context menu did not open on right-click")
	}
	if !ovs[0].ByText("Delete").Click() {
		t.Fatal("menu item click found no handler")
	}
	if sel != "Delete" {
		t.Fatalf("OnSelect = %q want Delete", sel)
	}
	if len(h.Overlays()) != 0 {
		t.Fatal("menu should close after select")
	}
}

// Phase-4: DataTable 分页/搜索基础渲染 + 比较函数。
func TestDataTablePaging(t *testing.T) {
	rows := []map[string]string{
		{"name": "Alice", "amt": "100"}, {"name": "Bob", "amt": "50"}, {"name": "Carol", "amt": "200"},
	}
	cols := []DataColumn{{Key: "name", Header: "Name", Sortable: true}, {Key: "amt", Header: "Amt"}}
	h := ui.MountDefault(DataTable(DataTableProps{Columns: cols, Rows: rows, Search: true, PageSize: 2}))
	if !h.Root().ByText("Alice").Exists() || !h.Root().ByText("Bob").Exists() {
		t.Fatalf("page-1 rows missing; texts=%v", h.Root().Texts())
	}
	if h.Root().ByText("Carol").Exists() {
		t.Fatal("Carol should be on page 2, not page 1")
	}
}

func TestLessValue(t *testing.T) {
	if !lessValue("50", "100") { // 数值比较
		t.Fatal("50 < 100 (numeric) failed")
	}
	if lessValue("100", "50") {
		t.Fatal("100 < 50 should be false")
	}
	if !lessValue("apple", "banana") { // 字符串比较
		t.Fatal("apple < banana (string) failed")
	}
}

// Phase-4: Sidebar 渲染分组菜单项。
func TestSidebarRenders(t *testing.T) {
	h := ui.MountDefault(Sidebar(SidebarProps{
		Header: H4("App"),
		Groups: []SidebarGroup{{Label: "Nav", Items: []SidebarItem{{Label: "Home", Active: true}, {Label: "Settings"}}}},
	}))
	for _, s := range []string{"App", "Nav", "Home", "Settings"} {
		if !h.Root().ByText(s).Exists() {
			t.Fatalf("missing %q; texts=%v", s, h.Root().Texts())
		}
	}
}

// Phase-4: Line/Area/Pie 通过 Vector 画出描边/填充路径。
func TestChartsDrawPaths(t *testing.T) {
	h := ui.MountDefault(ui.Use(func(_ struct{}) *ui.Node {
		return ui.Div(
			LineChart(LineChartProps{Data: []float32{1, 3, 2, 5}, Width: 200, Height: 100, Area: true}),
			PieChart(PieChartProps{Size: 120, Slices: []PieSlice{
				{Value: 1, Color: ui.Hex("#ff0000")}, {Value: 1, Color: ui.Hex("#00ff00")}}}),
		)
	}, struct{}{}))
	var stroke, fill int
	for _, op := range h.Paint() {
		switch op.Kind {
		case "strokepath":
			stroke++
		case "path":
			fill++
		}
	}
	if stroke < 1 || fill < 3 { // line 描边×1；area 填充×1 + pie 填充×2
		t.Fatalf("charts paths: stroke=%d fill=%d (want >=1 / >=3)", stroke, fill)
	}
}

// Form：无效提交显示错误，修正后提交回调。
func TestFormValidation(t *testing.T) {
	var submitted map[string]string
	app := func(_ struct{}) *ui.Node {
		return Form(FormProps{
			SubmitLabel: "Save",
			OnSubmit:    func(v map[string]string) { submitted = v },
			Fields: []FormField{
				{Name: "email", Label: "Email", Validate: func(v string) string {
					if v == "" {
						return "必填"
					}
					return ""
				}, Control: func(v string, oc func(string)) *ui.Node {
					return Input(InputProps{Value: v, OnChange: oc})
				}},
			},
		})
	}
	h := ui.MountDefault(ui.Use(app, struct{}{}))
	// 空提交 -> 错误显示，OnSubmit 不触发
	h.Root().ByText("Save").Click()
	if !h.Root().ByText("必填").Exists() {
		t.Fatalf("expected validation error; texts=%v", h.Root().Texts())
	}
	if submitted != nil {
		t.Fatal("OnSubmit should not fire with invalid form")
	}
	// 填入值 -> 再提交 -> 回调
	h.Root().ByKind("input").Type("a@b.com")
	h.Root().ByText("Save").Click()
	if submitted == nil || submitted["email"] != "a@b.com" {
		t.Fatalf("OnSubmit not called with valid form: %v", submitted)
	}
}
