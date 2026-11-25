package shadcn

import (
	"fmt"
	"math"
	"testing"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func assertInViewport(t *testing.T, r ui.Rect, w, h float32) {
	t.Helper()
	if r.W <= 0 || r.H <= 0 || r.X < 0 || r.Y < 0 || r.X+r.W > w+0.1 || r.Y+r.H > h+0.1 {
		t.Fatalf("floating content outside %vx%v viewport: %+v", w, h, r)
	}
}

func floatingPanelQuery(h *ui.Harness) *ui.Query {
	return h.Overlays()[0].Child(0).Child(0)
}

func TestPopoverStaysInViewportAfterResize(t *testing.T) {
	h := ui.Mount(ui.Div(ui.Style(ui.Fill),
		ui.Div(ui.Style(ui.Absolute, ui.Right(4), ui.Bottom(4)),
			Popover(ui.Text("Open panel"), ui.Text("A floating panel with useful content")))), 360, 240)
	h.Root().ByText("Open panel").Click()
	if len(h.Overlays()) != 1 {
		t.Fatal("popover did not open")
	}
	assertInViewport(t, floatingPanelQuery(h).Bounds(), 360, 240)
	h.Resize(240, 160)
	assertInViewport(t, floatingPanelQuery(h).Bounds(), 240, 160)
	h.Escape()
	if len(h.Overlays()) != 0 {
		t.Fatal("Escape did not close repositioned popover")
	}
}

func TestPopoverWithFixedAnchorStaysInViewportAfterResize(t *testing.T) {
	h := ui.Mount(ui.Div(ui.Style(ui.Fill),
		ui.Div(ui.Style(ui.Absolute, ui.Left(16), ui.Top(64)),
			Popover(ui.Text("Fixed anchor"), ui.Div(ui.Style(ui.Height(140)), ui.Text("Content"))))), 360, 300)
	h.Root().ByText("Fixed anchor").Click()
	h.Resize(240, 150)
	assertInViewport(t, floatingPanelQuery(h).Bounds(), 240, 150)
}

func TestSelectUsesTransformedAnchorAndWorksWithoutCallback(t *testing.T) {
	h := ui.Mount(ui.Div(ui.Style(ui.Fill),
		ui.Div(ui.Style(ui.Absolute, ui.Left(100), ui.Top(80), ui.Width(200), ui.Height(100), ui.ScaleXY(1.5, 1.25)),
			Select(SelectProps{Options: []string{"A", "B"}, Placeholder: "Choose", Style: []ui.StyleOpt{ui.Width(120), ui.Height(30)}}))), 400, 300)
	if !h.ClickAt(140, 86.25) {
		t.Fatal("scaled select trigger did not open")
	}
	if len(h.Overlays()) != 1 {
		t.Fatal("select did not open")
	}
	panel := floatingPanelQuery(h).Bounds()
	if math.Abs(float64(panel.X-50)) > 0.01 || math.Abs(float64(panel.Y-111)) > 0.01 {
		t.Fatalf("select panel is not aligned to transformed trigger: %+v", panel)
	}
	if !h.Overlays()[0].ByText("B").Click() {
		t.Fatal("select option did not fire")
	}
	if len(h.Overlays()) != 0 || !h.Root().ByText("B").Exists() {
		t.Fatal("uncontrolled select did not display the chosen value and close")
	}
}

func TestSelectItemsSeparateLabelFromValue(t *testing.T) {
	selected := ""
	h := ui.Mount(Select(SelectProps{
		Value:    "5",
		Items:    []SelectOption{{Label: "Master Rule 2020", Value: "5"}, {Label: "Master Rule 2017", Value: "4"}},
		OnChange: func(value string) { selected = value },
	}), 400, 300)
	defer h.Close()
	if !h.Root().ByText("Master Rule 2020").Click() || !h.Overlays()[0].ByText("Master Rule 2017").Click() {
		t.Fatal("label/value select interaction failed")
	}
	if selected != "4" || !h.Root().ByText("Master Rule 2017").Exists() {
		t.Fatalf("selected value=%q texts=%v", selected, h.Root().Texts())
	}
}

func TestSelectLongListRemainsScrollableAndSelectable(t *testing.T) {
	options := make([]string, 30)
	for i := range options {
		options[i] = fmt.Sprintf("Option %02d", i)
	}
	var selected string
	h := ui.Mount(ui.Div(ui.Style(ui.Fill),
		ui.Div(ui.Style(ui.Absolute, ui.Right(4), ui.Top(60)),
			Select(SelectProps{Options: options, Placeholder: "Choose", OnChange: func(v string) { selected = v }}))), 180, 180)
	h.Root().ByText("Choose").Click()
	assertInViewport(t, floatingPanelQuery(h).Bounds(), 180, 180)
	scroll := h.Overlays()[0].ByKind("scroll")
	if !scroll.Exists() || !scroll.ScrollBy(10000) {
		t.Fatal("long select list has no scrollable viewport")
	}
	if !h.Overlays()[0].ByText(options[29]).Click() || selected != options[29] {
		t.Fatalf("last option not selectable after scrolling: %q", selected)
	}
	if len(h.Overlays()) != 0 {
		t.Fatal("select did not close after choosing the last option")
	}
}

func TestTooltipFlipsBelowTopEdgeAndWraps(t *testing.T) {
	h := ui.Mount(ui.Div(ui.Style(ui.Fill),
		ui.Div(ui.Style(ui.Absolute, ui.Right(4), ui.Top(0)),
			Tooltip("A longer tooltip that must remain readable near the edge", ui.Text("Help")))), 220, 160)
	b := h.Root().ByText("Help").Bounds()
	h.MoveTo(b.X+b.W/2, b.Y+b.H/2)
	if len(h.Overlays()) != 1 {
		t.Fatal("tooltip did not open on hover")
	}
	panel := h.Overlays()[0].Child(0)
	assertInViewport(t, panel.Bounds(), 220, 160)
	if panel.Bounds().Y < b.Y+b.H {
		t.Fatalf("top-edge tooltip did not flip below trigger: %+v", panel.Bounds())
	}
	h.MoveTo(10, 150)
	if len(h.Overlays()) != 0 {
		t.Fatal("tooltip stayed open after leaving trigger")
	}
}

func TestHoverCardSurvivesPointerCrossingGap(t *testing.T) {
	clicked := 0
	h := ui.Mount(ui.Div(ui.Style(ui.Fill),
		ui.Div(ui.Style(ui.Absolute, ui.Left(24), ui.Top(32)),
			HoverCard(ui.Text("Profile"),
				Button(ButtonProps{OnClick: func() { clicked++ }}, ui.Text("View profile"))))), 360, 240)
	b := h.Root().ByText("Profile").Bounds()
	h.MoveTo(b.X+b.W/2, b.Y+b.H/2)
	h.Step(160)
	if len(h.Overlays()) != 1 {
		t.Fatal("hover card did not open")
	}
	h.MoveTo(b.X+b.W/2, b.Y+b.H+1)
	h.Step(40)
	if len(h.Overlays()) != 1 {
		t.Fatal("hover card closed while pointer was crossing the trigger/panel gap")
	}
	button := h.Overlays()[0].ByText("View profile")
	r := button.Bounds()
	h.MoveTo(r.X+r.W/2, r.Y+r.H/2)
	h.Step(160)
	if len(h.Overlays()) != 1 || !h.Overlays()[0].ByText("View profile").Click() || clicked != 1 {
		t.Fatal("hover card content did not stay interactive")
	}
	h.MoveTo(350, 230)
	h.Step(160)
	if len(h.Overlays()) != 0 {
		t.Fatal("hover card did not close after leaving both trigger and panel")
	}
}

func TestContextMenuStaysInViewport(t *testing.T) {
	selected := false
	h := ui.Mount(ContextMenu(ui.Div(ui.Style(ui.Fill), ui.Text("Surface")), []MenuItem{
		{Label: "First"}, {Label: "Last", OnSelect: func() { selected = true }},
	}), 240, 160)
	if !h.RightClickAt(238, 158) || len(h.Overlays()) != 1 {
		t.Fatal("context menu did not open at bottom-right corner")
	}
	assertInViewport(t, floatingPanelQuery(h).Bounds(), 240, 160)
	if !h.Overlays()[0].ByText("Last").Click() || !selected || len(h.Overlays()) != 0 {
		t.Fatal("repositioned context menu could not select and close")
	}
}

func TestSheetEscapeClosesAndPanelClickDoesNot(t *testing.T) {
	app := func(_ struct{}) *ui.Node {
		open, setOpen := ui.UseState(true)
		return ui.Div(Sheet(SheetProps{Open: open, OnClose: func() { setOpen(false) }}, ui.Text("Details")))
	}
	h := ui.MountDefault(ui.Use(app, struct{}{}))
	defer h.Close()
	h.Step(250)
	if len(h.Overlays()) != 1 {
		t.Fatal("sheet did not open")
	}
	h.Overlays()[0].ByText("Details").Click()
	if len(h.Overlays()) != 1 {
		t.Fatal("clicking sheet content closed the drawer")
	}
	h.Escape()
	h.Step(250)
	if len(h.Overlays()) != 0 {
		t.Fatal("Escape did not close the sheet")
	}
}

func TestCollapsibleToggleChangesContentHeight(t *testing.T) {
	app := func(_ struct{}) *ui.Node {
		open, setOpen := ui.UseState(false)
		return Collapsible(open, func() { setOpen(!open) }, ui.Text("Toggle"), ui.Text("Content"))
	}
	h := ui.MountDefault(ui.Use(app, struct{}{}))
	defer h.Close()
	if got := h.Root().Child(1).Bounds().H; got != 0 {
		t.Fatalf("initial content height=%v, want 0", got)
	}
	h.Root().ByText("Toggle").Click()
	h.Step(220)
	if got := h.Root().Child(1).Bounds().H; got <= 0 {
		t.Fatalf("expanded content height=%v, inner=%+v, texts=%v", got, h.Root().Child(1).Child(0).Bounds(), h.Root().Texts())
	}
	h.Root().ByText("Toggle").Click()
	h.Step(220)
	if got := h.Root().Child(1).Bounds().H; got != 0 {
		t.Fatalf("collapsed content height=%v, want 0", got)
	}
}

func TestTooltipHoverShowsAndHides(t *testing.T) {
	h := ui.MountDefault(Tooltip("Help text", ui.Text("Trigger")))
	defer h.Close()
	if len(h.Overlays()) != 0 {
		t.Fatal("tooltip was visible before hover")
	}
	trigger := h.Root().ByText("Trigger")
	trigger.Hover(true)
	if len(h.Overlays()) != 1 || !h.Overlays()[0].ByText("Help text").Exists() {
		t.Fatal("hover did not show tooltip")
	}
	trigger.Hover(false)
	if len(h.Overlays()) != 0 {
		t.Fatal("hover leave did not hide tooltip")
	}
}
