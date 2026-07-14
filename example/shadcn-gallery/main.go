package main

import (
	"fmt"
	"time"

	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func section(title string, body ...*ui.Node) *ui.Node {
	return shadcn.Card(
		shadcn.CardHeader(shadcn.CardTitle(title)),
		shadcn.CardContent(body...),
		ui.Div(ui.Style(ui.Height(24))),
	)
}

func rowN(children ...*ui.Node) *ui.Node {
	return ui.Div(append([]*ui.Node{ui.Style(ui.Row, ui.Gap(12), ui.ItemsCenter)}, children...)...)
}

func App(_ struct{}) *ui.Node {
	dark, setDark := ui.UseState(false)
	check, setCheck := ui.UseState(true)
	sw, setSw := ui.UseState(false)
	radio, setRadio := ui.UseState("A")
	vol, setVol := ui.UseState(float32(50))
	text, setText := ui.UseState("")
	tab, setTab := ui.UseState(0)
	toggle, setToggle := ui.UseState(false)
	open, setOpen := ui.UseState(false)
	sel, setSel := ui.UseState("")
	combo, setCombo := ui.UseState("")
	note, setNote := ui.UseState("")
	date, setDate := ui.UseState(time.Time{})
	cmdOpen, setCmdOpen := ui.UseState(false)
	sheetOpen, setSheetOpen := ui.UseState(false)
	page, setPage := ui.UseState(1)
	align, setAlign := ui.UseState("居中")
	collapse, setCollapse := ui.UseState(false)

	theme := ui.LightTheme
	if dark {
		theme = ui.DarkTheme
	}

	return ui.ThemeProvider(theme,
		ui.Div(
			ui.Style(ui.Column, ui.Width(720), ui.Height(760), ui.Bg(theme.Background),
				ui.TextColor(theme.Foreground), ui.Padding(24), ui.Gap(16)),

			// 顶栏（渐变横幅）
			ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.JustifyBetween, ui.PaddingXY(16, 12), ui.Radius(12),
				ui.LinearGradient(ui.Hex("#6366f1"), ui.Hex("#ec4899"), 60), ui.TextColor(ui.White)),
				ui.Text("shadcn/ui 组件库", ui.FontSize(24)),
				rowN(
					shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeSecondary}, ui.Text("Tenon")),
					shadcn.Switch(shadcn.SwitchProps{Checked: dark, OnChange: setDark}),
					shadcn.Label("暗色"),
				),
			),

			ui.ScrollView(
				ui.Style(ui.Column, ui.Gap(16), ui.Grow(1)),

				section("Button",
					rowN(
						shadcn.Button(shadcn.ButtonProps{}, ui.Text("Default")),
						shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Secondary}, ui.Text("Secondary")),
						shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Destructive}, ui.Text("Destructive")),
						shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline}, ui.Text("Outline")),
						shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Ghost}, ui.Text("Ghost")),
					),
				),

				section("Badge",
					rowN(
						shadcn.Badge(shadcn.BadgeProps{}, ui.Text("Default")),
						shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeSecondary}, ui.Text("Secondary")),
						shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeDestructive}, ui.Text("Destructive")),
						shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeOutline}, ui.Text("Outline")),
					),
				),

				section("Form",
					rowN(
						shadcn.Checkbox(shadcn.CheckboxProps{Checked: check, OnChange: setCheck}),
						shadcn.Label("接受条款"),
						shadcn.Switch(shadcn.SwitchProps{Checked: sw, OnChange: setSw}),
						shadcn.Label("通知"),
					),
					ui.Div(ui.Style(ui.Height(10))),
					shadcn.RadioGroup(shadcn.RadioGroupProps{Value: radio, Options: []string{"A", "B", "C"}, OnChange: setRadio}),
					ui.Div(ui.Style(ui.Height(10))),
					shadcn.Input(shadcn.InputProps{Value: text, OnChange: setText, Placeholder: "输入点什么…"}),
				),

				section("Slider & Progress",
					rowN(
						shadcn.Slider(shadcn.SliderProps{Value: vol, Min: 0, Max: 100, OnChange: setVol}),
						ui.Text(fmt.Sprintf("%.0f", vol), ui.FontSize(14)),
					),
					ui.Div(ui.Style(ui.Height(10))),
					shadcn.Progress(vol/100),
				),

				section("Tabs & Toggle",
					shadcn.Tabs(shadcn.TabsProps{Tabs: []string{"账户", "密码", "通知"}, Active: tab, OnChange: setTab}),
					ui.Div(ui.Style(ui.Height(10))),
					shadcn.Toggle(shadcn.ToggleProps{Pressed: toggle, OnChange: setToggle}, ui.Text("加粗")),
				),

				section("Alert",
					shadcn.Alert(shadcn.AlertProps{},
						shadcn.AlertTitle("提示"),
						shadcn.AlertDescription("这是一条默认样式的提示信息。"),
					),
					ui.Div(ui.Style(ui.Height(10))),
					shadcn.Alert(shadcn.AlertProps{Variant: shadcn.AlertDestructive},
						shadcn.AlertTitle("错误"),
						shadcn.AlertDescription("出问题了，请稍后重试。"),
					),
				),

				section("Display",
					rowN(
						shadcn.Avatar("孙", 44),
						shadcn.Avatar("AI", 44),
						shadcn.Skeleton(120, 16),
						shadcn.Skeleton(80, 16),
					),
					ui.Div(ui.Style(ui.Height(12))),
					shadcn.Separator(shadcn.SeparatorProps{}),
				),

				section("Overlays（锚定浮层）",
					rowN(
						shadcn.Select(shadcn.SelectProps{Value: sel, Options: []string{"苹果", "香蕉", "橙子"},
							OnChange: setSel, Placeholder: "选择水果"}),
						shadcn.Combobox(shadcn.ComboboxProps{Value: combo, OnChange: setCombo,
							Placeholder: "搜索框架…", SearchPlaceholder: "输入过滤…",
							Options: []shadcn.ComboboxOption{
								{Value: "go", Label: "Go"}, {Value: "rust", Label: "Rust"},
								{Value: "ts", Label: "TypeScript"}, {Value: "py", Label: "Python"},
								{Value: "zig", Label: "Zig"},
							}}),
						shadcn.Popover(shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline}, ui.Text("Popover")),
							ui.Text("这是一个锚定在按钮下方的浮层。", ui.FontSize(13)),
						),
						shadcn.DropdownMenu(shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Secondary}, ui.Text("菜单")),
							[]shadcn.MenuItem{
								{Label: "个人资料", OnSelect: func() {}},
								{Label: "设置", OnSelect: func() {}},
								{Label: "退出", OnSelect: func() {}},
							}),
						shadcn.Tooltip("悬停提示", shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeOutline}, ui.Text("Tooltip"))),
					),
				),

				section("Textarea",
					shadcn.Textarea(shadcn.TextareaProps{Value: note, OnChange: setNote, Placeholder: "多行文本，Enter 换行…", Rows: 3}),
				),

				section("Table",
					shadcn.Table(
						shadcn.TableRow(shadcn.TableHead("名称"), shadcn.TableHead("角色"), shadcn.TableHead("状态")),
						shadcn.TableRow(shadcn.TableCell(ui.Text("孙江萌")), shadcn.TableCell(ui.Text("管理员")), shadcn.TableCell(shadcn.Badge(shadcn.BadgeProps{}, ui.Text("在线")))),
						shadcn.TableRow(shadcn.TableCell(ui.Text("AI 助手")), shadcn.TableCell(ui.Text("成员")), shadcn.TableCell(shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeSecondary}, ui.Text("离线")))),
					),
				),

				section("Accordion",
					shadcn.Accordion([]shadcn.AccordionItemData{
						{Title: "这是什么？", Content: []*ui.Node{ui.Text("一个基于 tenon/pkg/ui 的 shadcn 风格组件库。", ui.FontSize(13))}},
						{Title: "支持主题吗？", Content: []*ui.Node{ui.Text("支持，通过 ThemeProvider 切换明暗。", ui.FontSize(13))}},
						{Title: "带动画吗？", Content: []*ui.Node{ui.Text("展开/收起有高度过渡动画。", ui.FontSize(13))}},
					}),
				),

				section("Breadcrumb / Pagination / ToggleGroup",
					shadcn.Breadcrumb([]string{"首页", "组件", "按钮"}, func(int) {}),
					ui.Div(ui.Style(ui.Height(10))),
					rowN(
						shadcn.Pagination(shadcn.PaginationProps{Page: page, Total: 5, OnChange: setPage}),
						shadcn.ToggleGroup(shadcn.ToggleGroupProps{Value: align, Options: []string{"左", "居中", "右"}, OnChange: setAlign}),
					),
				),

				section("Collapsible",
					shadcn.Collapsible(collapse, func() { setCollapse(!collapse) },
						rowN(ui.Text("高级选项", ui.FontSize(15)), ui.Text("▾", ui.FontSize(12))),
						ui.Text("这里是折叠内容，点击标题展开或收起。", ui.FontSize(13)),
					),
				),

				section("NavigationMenu",
					shadcn.NavigationMenu([]shadcn.NavItem{
						{Label: "首页", OnSelect: func() {}},
						{Label: "产品", Items: []shadcn.MenuItem{{Label: "概览"}, {Label: "定价"}, {Label: "文档"}}},
						{Label: "关于", OnSelect: func() {}},
					}),
				),

				section("Carousel",
					shadcn.Carousel(shadcn.CarouselProps{Width: 320, Height: 130, Slides: []*ui.Node{
						ui.Div(ui.Style(ui.Width(320), ui.Height(130), ui.Bg(ui.Hex("#ef4444")), ui.ItemsCenter, ui.JustifyCenter), ui.Text("Slide 1", ui.FontSize(20), ui.TextColor(ui.White))),
						ui.Div(ui.Style(ui.Width(320), ui.Height(130), ui.Bg(ui.Hex("#22c55e")), ui.ItemsCenter, ui.JustifyCenter), ui.Text("Slide 2", ui.FontSize(20), ui.TextColor(ui.White))),
						ui.Div(ui.Style(ui.Width(320), ui.Height(130), ui.Bg(ui.Hex("#3b82f6")), ui.ItemsCenter, ui.JustifyCenter), ui.Text("Slide 3", ui.FontSize(20), ui.TextColor(ui.White))),
					}}),
				),

				section("Resizable",
					shadcn.Resizable(shadcn.ResizableProps{Width: 420, Height: 120,
						Left:  ui.Div(ui.Style(ui.Grow(1), ui.ItemsCenter, ui.JustifyCenter), ui.Text("左栏", ui.FontSize(14))),
						Right: ui.Div(ui.Style(ui.Grow(1), ui.ItemsCenter, ui.JustifyCenter), ui.Text("右栏（拖动中间分隔条）", ui.FontSize(13))),
					}),
				),

				section("BarChart & HoverCard",
					shadcn.BarChart(shadcn.BarChartProps{Width: 320, Height: 140,
						Data: []float32{40, 75, 55, 90, 30, 65}, Labels: []string{"一", "二", "三", "四", "五", "六"}}),
					ui.Div(ui.Style(ui.Height(10))),
					shadcn.HoverCard(
						shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeOutline}, ui.Text("悬停查看")),
						ui.Text("HoverCard", ui.FontSize(15)),
						ui.Text("悬停触发的信息卡片，移入卡片本身也保持展开。", ui.FontSize(13)),
					),
				),

				section("Calendar",
					shadcn.Calendar(shadcn.CalendarProps{Value: date, OnChange: setDate}),
					ui.Div(ui.Style(ui.Height(8))),
					ui.Text(fmt.Sprintf("已选：%s", func() string {
						if date.IsZero() {
							return "无"
						}
						return date.Format("2006-01-02")
					}()), ui.FontSize(13)),
				),

				section("Dialog / Sheet / Command / Toast",
					rowN(
						shadcn.Button(shadcn.ButtonProps{OnClick: func() { setOpen(true) }}, ui.Text("对话框")),
						shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Secondary, OnClick: func() { setSheetOpen(true) }}, ui.Text("抽屉")),
						shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline, OnClick: func() { setCmdOpen(true) }}, ui.Text("命令面板")),
						shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Ghost, OnClick: func() {
							shadcn.Toast("已保存", "你的更改已成功保存。")
						}}, ui.Text("通知")),
					),
				),
			),

			shadcn.Dialog(shadcn.DialogProps{Open: open, OnClose: func() { setOpen(false) }},
				shadcn.DialogTitle("确认操作"),
				shadcn.DialogDescription("这是一个通过 Portal 渲染的模态对话框。"),
				ui.Div(ui.Style(ui.Row, ui.Gap(8), ui.JustifyEnd),
					shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline, OnClick: func() { setOpen(false) }}, ui.Text("取消")),
					shadcn.Button(shadcn.ButtonProps{OnClick: func() { setOpen(false) }}, ui.Text("确认")),
				),
			),

			shadcn.Sheet(shadcn.SheetProps{Open: sheetOpen, OnClose: func() { setSheetOpen(false) }},
				shadcn.SheetTitle("设置"),
				shadcn.SheetDescription("从右侧滑入的抽屉面板。"),
				shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline, OnClick: func() { setSheetOpen(false) }}, ui.Text("关闭")),
			),

			shadcn.Command(shadcn.CommandProps{Open: cmdOpen, OnClose: func() { setCmdOpen(false) },
				Items: []shadcn.CommandItem{
					{Label: "新建文件", OnSelect: func() { shadcn.Toast("新建文件", "") }},
					{Label: "打开设置", OnSelect: func() { setSheetOpen(true) }},
					{Label: "切换主题", OnSelect: func() { setDark(!dark) }},
					{Label: "退出", OnSelect: func() {}},
				}}),

			shadcn.Toaster(),
		),
	)
}

func main() {
	ui.Run(ui.Use(App, struct{}{}))
}
