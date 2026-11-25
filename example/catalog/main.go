package main

import (
	"fmt"
	"time"

	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

func catalogSection(name string, node *ui.Node) *ui.Node {
	return shadcn.Card(
		shadcn.CardHeader(shadcn.CardTitle(name)),
		shadcn.CardContent(node),
	)
}

func catalogInline(children ...*ui.Node) *ui.Node {
	return ui.Div(append([]*ui.Node{ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(8))}, children...)...)
}

type catalogApp struct {
	ui.Component
	text, message value.Bind[string]
	status        value.Bind[string]
	checked       bool
	switched      bool
	radio         string
	slider        float32
	activeTab     int
	pressed       bool
	toggleGroup   string
	page          int
	date          time.Time
	collapsible   bool
	dialog        bool
	sheet         bool
	command       bool
	selectValue   string
	comboboxValue string
}

func (a *catalogApp) Render() ui.NodeSource {
	if a.status == nil {
		a.text, a.message = value.New(""), value.New("")
		a.status = value.New("Ready")
		a.radio, a.toggleGroup, a.page = "A", "A", 1
	}
	sections := []*ui.Node{
		ui.VStack(6,
			ui.Text("Component catalog", ui.FontSize(28), ui.Bold),
			ui.Text("Themeable controls with keyboard, pointer and reactive state.").FontSize(14).TextColor(ui.Hex("#64748b")),
		),
		catalogSection("Interaction status", ui.Text(a.status)),
		catalogSection("Button", catalogInline(shadcn.Button(shadcn.ButtonProps{OnClick: func() { a.status.Set("Button clicked") }}, ui.Text("Continue")))),
		catalogSection("Badge", catalogInline(shadcn.Badge(shadcn.BadgeProps{}, ui.Text("New")))),
		catalogSection("Card", shadcn.Card(shadcn.CardContent(ui.Text("Content")))),
		catalogSection("Input", shadcn.Input(shadcn.InputProps{Placeholder: "Email", Bind: a.text})),
		catalogSection("Textarea", shadcn.Textarea(shadcn.TextareaProps{Placeholder: "Message", Rows: 3, Bind: a.message})),
		catalogSection("Label", shadcn.Label("Label")),
		catalogSection("Checkbox", shadcn.Checkbox(shadcn.CheckboxProps{Checked: a.checked, OnChange: func(v bool) {
			a.Context().Update(func() { a.checked = v })
			a.status.Set("Checkbox changed")
		}})),
		catalogSection("Switch", shadcn.Switch(shadcn.SwitchProps{Checked: a.switched, OnChange: func(v bool) {
			a.Context().Update(func() { a.switched = v })
			a.status.Set("Switch changed")
		}})),
		catalogSection("RadioGroup", shadcn.RadioGroup(shadcn.RadioGroupProps{Value: a.radio, Options: []string{"A", "B"}, OnChange: func(v string) {
			a.Context().Update(func() { a.radio = v })
			a.status.Set("Radio selection: " + v)
		}})),
		catalogSection("Slider", shadcn.Slider(shadcn.SliderProps{Value: a.slider, Min: 0, Max: 10, OnChange: func(v float32) {
			a.Context().Update(func() { a.slider = v })
			a.status.Set("Slider changed")
		}})),
		catalogSection("Progress", shadcn.Progress(0.6)),
		catalogSection("Separator", shadcn.Separator(shadcn.SeparatorProps{})),
		catalogSection("Skeleton", shadcn.Skeleton(240, 18)),
		catalogSection("Avatar", catalogInline(shadcn.Avatar("TN", 44))),
		catalogSection("Alert", shadcn.Alert(shadcn.AlertProps{}, shadcn.AlertTitle("Ready"), shadcn.AlertDescription("Component rendered"))),
		catalogSection("Tabs", shadcn.Tabs(shadcn.TabsProps{Tabs: []string{"A", "B"}, Active: a.activeTab, OnChange: func(v int) {
			a.Context().Update(func() { a.activeTab = v })
			a.status.Set("Tab changed")
		}})),
		catalogSection("Toggle", catalogInline(shadcn.Toggle(shadcn.ToggleProps{Pressed: a.pressed, OnChange: func(v bool) {
			a.Context().Update(func() { a.pressed = v })
			a.status.Set("Toggle changed")
		}}, ui.Text("Bold")))),
		catalogSection("Dialog", ui.VStack(8,
			shadcn.Button(shadcn.ButtonProps{OnClick: func() {
				a.Context().Update(func() { a.dialog = true })
			}}, ui.Text("Open dialog")),
			shadcn.Dialog(shadcn.DialogProps{Open: a.dialog, OnClose: func() {
				a.Context().Update(func() { a.dialog = false })
			}}, shadcn.DialogTitle("Dialog"), ui.Text("Dialog content")),
		)),
		catalogSection("Popover", catalogInline(shadcn.Popover(ui.Text("Open popover"), ui.Text("Popover content")))),
		catalogSection("Tooltip", catalogInline(shadcn.Tooltip("Helpful hint", ui.Text("Hover for tooltip")))),
		catalogSection("DropdownMenu", catalogInline(shadcn.DropdownMenu(ui.Text("Open menu"), []shadcn.MenuItem{{Label: "Item", OnSelect: func() { a.status.Set("Dropdown item selected") }}}))),
		catalogSection("Select", shadcn.Select(shadcn.SelectProps{Options: []string{"A", "B"}, Value: a.selectValue, Placeholder: "Select", OnChange: func(v string) {
			a.Context().Update(func() { a.selectValue = v })
			a.status.Set("Selected: " + v)
		}})),
		catalogSection("Combobox", shadcn.Combobox(shadcn.ComboboxProps{Options: []shadcn.ComboboxOption{{Value: "a", Label: "Apple"}}, Value: a.comboboxValue, Placeholder: "Search", OnChange: func(v string) {
			a.Context().Update(func() { a.comboboxValue = v })
			a.status.Set("Combobox selection: " + v)
		}})),
		catalogSection("Table", shadcn.Table(shadcn.TableRow(shadcn.TableHead("Name"), shadcn.TableHead("Value")), shadcn.TableRow(shadcn.TableCell(ui.Text("A")), shadcn.TableCell(ui.Text("1"))))),
		catalogSection("Accordion", shadcn.Accordion([]shadcn.AccordionItemData{{Title: "Section", Content: []*ui.Node{ui.Text("Content")}}})),
		catalogSection("Toaster", shadcn.Toaster()),
		catalogSection("Sheet", ui.VStack(8,
			shadcn.Button(shadcn.ButtonProps{OnClick: func() {
				a.Context().Update(func() { a.sheet = true })
			}}, ui.Text("Open sheet")),
			shadcn.Sheet(shadcn.SheetProps{Open: a.sheet, OnClose: func() {
				a.Context().Update(func() { a.sheet = false })
			}}, shadcn.SheetTitle("Sheet"), ui.Text("Sheet content")),
		)),
		catalogSection("Command", ui.VStack(8,
			shadcn.Button(shadcn.ButtonProps{OnClick: func() {
				a.Context().Update(func() { a.command = true })
			}}, ui.Text("Open command")),
			shadcn.Command(shadcn.CommandProps{Open: a.command, OnClose: func() {
				a.Context().Update(func() { a.command = false })
			}, Items: []shadcn.CommandItem{{Label: "Run command", OnSelect: func() { a.status.Set("Command executed") }}}}),
		)),
		catalogSection("Calendar", shadcn.Calendar(shadcn.CalendarProps{Value: a.date, OnChange: func(v time.Time) {
			a.Context().Update(func() { a.date = v })
			a.status.Set("Calendar selection: " + v.Format("2006-01-02"))
		}})),
		catalogSection("NavigationMenu", shadcn.NavigationMenu([]shadcn.NavItem{{Label: "Home", OnSelect: func() { a.status.Set("Navigation: Home") }}})),
		catalogSection("Menubar", shadcn.Menubar([]shadcn.MenubarMenu{{Label: "File", Items: []shadcn.MenuItem{{Label: "New", OnSelect: func() { a.status.Set("Menu: New") }}}}})),
		catalogSection("Breadcrumb", shadcn.Breadcrumb(
			[]string{"Home", "Components"},
			func(index int) { a.status.Set(fmt.Sprintf("Breadcrumb: %d", index)) },
		)),
		catalogSection("Pagination", shadcn.Pagination(shadcn.PaginationProps{Page: a.page, Total: 3, OnChange: func(v int) {
			a.Context().Update(func() { a.page = v })
			a.status.Set("Page changed")
		}})),
		catalogSection("DatePicker", shadcn.DatePicker(shadcn.DatePickerProps{Value: a.date, Placeholder: "Pick a date", OnChange: func(v time.Time) {
			a.Context().Update(func() { a.date = v })
			a.status.Set("Date selected: " + v.Format("2006-01-02"))
		}})),
		catalogSection("Collapsible", shadcn.Collapsible(a.collapsible, func() {
			a.Context().Update(func() { a.collapsible = !a.collapsible })
		}, ui.Text("Toggle"), ui.Text("Content"))),
		catalogSection("ToggleGroup", shadcn.ToggleGroup(shadcn.ToggleGroupProps{Value: a.toggleGroup, Options: []string{"A", "B"}, OnChange: func(v string) {
			a.Context().Update(func() { a.toggleGroup = v })
			a.status.Set("Toggle group: " + v)
		}})),
		catalogSection("AspectRatio", shadcn.AspectRatio(160, 1.5, ui.Text("16:9 content"))),
		catalogSection("ScrollArea", shadcn.ScrollArea(100, ui.Text("Scrollable content"))),
		catalogSection("HoverCard", shadcn.HoverCard(ui.Text("Hover card"), ui.Text("Details"))),
		catalogSection("Carousel", shadcn.Carousel(shadcn.CarouselProps{Slides: []*ui.Node{ui.Text("First"), ui.Text("Second")}, Width: 200, Height: 100})),
		catalogSection("Resizable", shadcn.Resizable(shadcn.ResizableProps{Width: 300, Height: 100, Left: ui.Text("Left"), Right: ui.Text("Right")})),
		catalogSection("BarChart", shadcn.BarChart(shadcn.BarChartProps{Data: []float32{1, 2, 3}, Labels: []string{"A", "B", "C"}, Width: 240, Height: 120})),
	}
	return ui.ScrollView(ui.Style(ui.Fill, ui.ItemsCenter), ui.VStack(16, sections...).Width(760).Padding(24))
}

func main() {
	ui.WindowTitle("Tenon — component catalog")
	ui.WindowSize(1000, 800)
	ui.Run(ui.Render(&catalogApp{}))
}
