package shadcn

import (
	"testing"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// TestConstruct 确认所有组件构造器都能返回非空节点（编译期签名 + 运行期不 panic）。
// 端到端渲染由 example/accordion 的实机运行覆盖。
func TestConstruct(t *testing.T) {
	nodes := map[string]*ui.Node{
		"Button":       Button(ButtonProps{Variant: Outline, OnClick: func() {}}, ui.Text("x")),
		"Badge":        Badge(BadgeProps{Variant: BadgeDestructive}, ui.Text("x")),
		"Card":         Card(CardHeader(CardTitle("t")), CardContent(CardDescription("d"))),
		"Input":        Input(InputProps{Placeholder: "p"}),
		"Label":        Label("l"),
		"Checkbox":     Checkbox(CheckboxProps{Checked: true}),
		"Switch":       Switch(SwitchProps{Checked: true}),
		"RadioGroup":   RadioGroup(RadioGroupProps{Value: "A", Options: []string{"A", "B"}}),
		"Slider":       Slider(SliderProps{Value: 5, Min: 0, Max: 10}),
		"Progress":     Progress(0.5),
		"Separator":    Separator(SeparatorProps{}),
		"Skeleton":     Skeleton(100, 16),
		"Avatar":       Avatar("AB", 40),
		"Alert":        Alert(AlertProps{Variant: AlertDestructive}, AlertTitle("t"), AlertDescription("d")),
		"Tabs":         Tabs(TabsProps{Tabs: []string{"a", "b"}, Active: 0}),
		"Toggle":       Toggle(ToggleProps{Pressed: true}, ui.Text("B")),
		"Dialog":       Dialog(DialogProps{Open: true}, DialogTitle("t")),
		"Textarea":     Textarea(TextareaProps{Placeholder: "p"}),
		"Popover":      Popover(ui.Text("open"), ui.Text("content")),
		"Tooltip":      Tooltip("hi", ui.Text("hover")),
		"DropdownMenu": DropdownMenu(ui.Text("menu"), []MenuItem{{Label: "A", OnSelect: func() {}}}),
		"Select":       Select(SelectProps{Options: []string{"A", "B"}, Placeholder: "pick"}),
		"Table":        Table(TableRow(TableHead("a"), TableHead("b")), TableRow(TableCell(ui.Text("1")), TableCell(ui.Text("2")))),
		"Accordion":    Accordion([]AccordionItemData{{Title: "t", Content: []*ui.Node{ui.Text("c")}}}),
		"Toaster":      Toaster(),
		"Sheet":        Sheet(SheetProps{Open: true}, ui.Text("x")),
		"Command":      Command(CommandProps{Open: true, Items: []CommandItem{{Label: "a"}}}),
		"Calendar":     Calendar(CalendarProps{}),
		"NavMenu":      NavigationMenu([]NavItem{{Label: "a", Items: []MenuItem{{Label: "b"}}}}),
		"Breadcrumb":   Breadcrumb([]string{"a", "b"}, func(int) {}),
		"Pagination":   Pagination(PaginationProps{Page: 1, Total: 3}),
		"AspectRatio":  AspectRatio(160, 1.5, ui.Text("x")),
		"ScrollArea":   ScrollArea(120, ui.Text("x")),
		"Collapsible":  Collapsible(true, func() {}, ui.Text("t"), ui.Text("c")),
		"ToggleGroup":  ToggleGroup(ToggleGroupProps{Value: "a", Options: []string{"a", "b"}}),
		"HoverCard":    HoverCard(ui.Text("t"), ui.Text("c")),
		"Carousel":     Carousel(CarouselProps{Slides: []*ui.Node{ui.Text("1"), ui.Text("2")}, Width: 200, Height: 100}),
		"Resizable":    Resizable(ResizableProps{Width: 300, Height: 100, Left: ui.Text("L"), Right: ui.Text("R")}),
		"BarChart":     BarChart(BarChartProps{Data: []float32{1, 2, 3}, Labels: []string{"a", "b", "c"}, Width: 200, Height: 120}),
	}
	for name, n := range nodes {
		if n == nil {
			t.Fatalf("%s returned nil", name)
		}
	}
}
