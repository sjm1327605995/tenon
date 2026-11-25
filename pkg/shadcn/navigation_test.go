package shadcn

import (
	"fmt"
	"testing"
	"time"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func TestDropdownMenuArrowNavigationAndEnterSelect(t *testing.T) {
	selected := ""
	h := ui.MountDefault(DropdownMenu(ui.Text("Menu"), []MenuItem{
		{Label: "A"}, {Label: "B", OnSelect: func() { selected = "B" }},
	}))
	defer h.Close()
	if !h.Root().ByText("Menu").Click() || len(h.Overlays()) != 1 {
		t.Fatal("dropdown did not open")
	}
	if got := h.Tab().AllText(); got != "Menu" {
		t.Fatalf("first Tab focus=%q, want Menu", got)
	}
	if got := h.Tab().AllText(); got != "A" {
		t.Fatalf("second Tab focus=%q, want A", got)
	}
	if got := h.Arrow(ui.NavVertical, true).AllText(); got != "B" {
		t.Fatalf("arrow focus=%q, want B", got)
	}
	if !h.Enter() || selected != "B" {
		t.Fatalf("Enter selection=%q, want B", selected)
	}
	if len(h.Overlays()) != 0 {
		t.Fatal("selected dropdown item did not close the menu")
	}
}

func TestCommandFiltersAndSelects(t *testing.T) {
	selected := ""
	app := func(_ struct{}) *ui.Node {
		open, setOpen := ui.UseState(true)
		return commandTestApp(open, func() { setOpen(false) }, func(label string) { selected = label })
	}
	h := ui.MountDefault(ui.Use(app, struct{}{}))
	defer h.Close()
	h.Step(200)
	if len(h.Overlays()) != 1 {
		t.Fatal("command palette did not open")
	}
	h.Overlays()[0].ByPlaceholder("输入命令搜索…").SetValue("av")
	if !h.Overlays()[0].ByText("Avocado").Exists() || h.Overlays()[0].ByText("Banana").Exists() {
		t.Fatal("command search did not filter items case-insensitively")
	}
	h.Overlays()[0].ByText("Avocado").Click()
	if selected != "Avocado" {
		t.Fatalf("selected=%q, want Avocado", selected)
	}
	h.Step(200)
	if len(h.Overlays()) != 0 {
		t.Fatal("selected command did not close the palette")
	}
}

func commandTestApp(open bool, close func(), selectItem func(string)) *ui.Node {
	return Command(CommandProps{
		Open: open, OnClose: close,
		Items: []CommandItem{
			{Label: "Apple"}, {Label: "Banana"}, {Label: "Avocado", OnSelect: func() { selectItem("Avocado") }},
		},
	})
}

func TestMenubarSelectsAndCloses(t *testing.T) {
	selected := ""
	h := ui.MountDefault(Menubar([]MenubarMenu{{
		Label: "File", Items: []MenuItem{{Label: "New", OnSelect: func() { selected = "New" }}},
	}}))
	defer h.Close()
	if !h.Root().ByText("File").Click() || len(h.Overlays()) != 1 {
		t.Fatal("menubar did not open")
	}
	h.Overlays()[0].ByText("New").Click()
	if selected != "New" || len(h.Overlays()) != 0 {
		t.Fatalf("menubar selection=%q overlays=%d", selected, len(h.Overlays()))
	}
}

func TestPaginationNavigationAndBoundary(t *testing.T) {
	app := func(_ struct{}) *ui.Node {
		page, setPage := ui.UseState(1)
		return ui.Div(
			Pagination(PaginationProps{Page: page, Total: 3, OnChange: func(value int) { setPage(value) }}),
			ui.Text(fmt.Sprintf("page=%d", page)),
		)
	}
	h := ui.MountDefault(ui.Use(app, struct{}{}))
	defer h.Close()
	if !h.Root().ByText("›").Click() || !h.Root().ByText("page=2").Exists() {
		t.Fatal("pagination next did not advance the page")
	}
	if !h.Root().ByText("‹").Click() || !h.Root().ByText("page=1").Exists() {
		t.Fatal("pagination previous did not return to the first page")
	}
	if h.Root().ByText("‹").Click() {
		t.Fatal("previous control was active at the first page")
	}
}

func TestDatePickerSelectsAndCloses(t *testing.T) {
	var selected time.Time
	app := func(_ struct{}) *ui.Node {
		value, setValue := ui.UseState(time.Date(2026, 7, 14, 0, 0, 0, 0, time.UTC))
		return ui.Div(DatePicker(DatePickerProps{Value: value, OnChange: func(next time.Time) {
			selected = next
			setValue(next)
		}}))
	}
	h := ui.MountDefault(ui.Use(app, struct{}{}))
	defer h.Close()
	if !h.Root().ByText("2026-07-14").Click() || len(h.Overlays()) != 1 {
		t.Fatal("date picker did not open")
	}
	if !h.Overlays()[0].ByText("15").Click() {
		t.Fatal("calendar day 15 was not selectable")
	}
	if selected.Day() != 15 || len(h.Overlays()) != 0 || !h.Root().ByText("2026-07-15").Exists() {
		t.Fatalf("date selection=%v overlays=%d texts=%v", selected, len(h.Overlays()), h.Root().Texts())
	}
}

func TestCalendarNavigationAndDateSelection(t *testing.T) {
	selected := time.Time{}
	h := ui.MountDefault(Calendar(CalendarProps{
		Value:    time.Date(2026, 7, 14, 0, 0, 0, 0, time.UTC),
		OnChange: func(value time.Time) { selected = value },
	}))
	defer h.Close()
	if !h.Root().ByText("›").Click() || !h.Root().ByText("2026 年 8 月").Exists() {
		t.Fatal("calendar did not advance to August")
	}
	if !h.Root().ByText("‹").Click() || !h.Root().ByText("2026 年 7 月").Exists() {
		t.Fatal("calendar did not return to July")
	}
	h.Root().ByText("12").Click()
	if selected.Day() != 12 || selected.Month() != time.July {
		t.Fatalf("selected date=%v, want July 12", selected)
	}
}

func TestNavigationMenuLinksAndDropdown(t *testing.T) {
	selected := ""
	h := ui.MountDefault(NavigationMenu([]NavItem{
		{Label: "Home", OnSelect: func() { selected = "Home" }},
		{Label: "Products", Items: []MenuItem{{Label: "Docs", OnSelect: func() { selected = "Docs" }}}},
	}))
	defer h.Close()
	h.Root().ByText("Home").Click()
	if selected != "Home" {
		t.Fatalf("direct navigation selection=%q, want Home", selected)
	}
	if !h.Root().ByText("Products").Click() || len(h.Overlays()) != 1 {
		t.Fatal("navigation dropdown did not open")
	}
	h.Overlays()[0].ByText("Docs").Click()
	if selected != "Docs" || len(h.Overlays()) != 0 {
		t.Fatalf("dropdown selection=%q overlays=%d", selected, len(h.Overlays()))
	}
}

func TestBreadcrumbReportsSelectedIndex(t *testing.T) {
	selected := -1
	h := ui.MountDefault(Breadcrumb([]string{"Home", "Library", "Current"}, func(index int) { selected = index }))
	defer h.Close()
	if !h.Root().ByText("Library").Click() || selected != 1 {
		t.Fatalf("breadcrumb selected index=%d, want 1", selected)
	}
}
