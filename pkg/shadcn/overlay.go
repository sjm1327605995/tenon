package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

// floatPanel 在锚点附近渲染浮层：自动翻转、贴边，内容过长时可滚动。
// 点击浮层外区域触发 onClose。
func floatPanel(th ui.Theme, anchor ui.Rect, onClose func(), minWidth float32, extra []ui.StyleOpt, content ...*ui.Node) *ui.Node {
	return ui.Use(floatPanelC, floatProps{th: th, anchor: anchor, onClose: onClose, minWidth: minWidth, extra: extra, content: content})
}

type floatProps struct {
	th       ui.Theme
	anchor   ui.Rect
	onClose  func()
	minWidth float32
	extra    []ui.StyleOpt
	content  []*ui.Node
}

// useFloatingPosition 用实测尺寸定位，首帧测量前隐藏，避免浮层闪到错误的位置。
// 所有坐标均为逻辑像素；窄窗口内最小宽度也必须服从视口上限。
func useFloatingPosition(anchor ui.Rect, gap float32, above bool, opacity float32) (*ui.Node, []ui.StyleOpt) {
	ref, rect := ui.UseMeasure()
	vp := ui.Viewport()
	marginX, marginY := min(float32(8), vp.W/2), min(float32(8), vp.H/2)
	left := clampf(anchor.X, marginX, max(marginX, vp.W-marginX-rect.W))
	below, over := anchor.Y+anchor.H+gap, anchor.Y-rect.H-gap
	top, alternate := below, over
	if above {
		top, alternate = over, below
	}
	fits := func(y float32) bool { return y >= marginY && y+rect.H <= vp.H-marginY }
	if !fits(top) && fits(alternate) {
		top = alternate
	}
	top = clampf(top, marginY, max(marginY, vp.H-marginY-rect.H))
	if rect.W == 0 || rect.H == 0 {
		opacity = 0
	}
	return ref, []ui.StyleOpt{
		ui.Absolute, ui.Left(left), ui.Top(top), ui.Opacity(opacity),
		ui.MaxWidth(max(float32(0), vp.W-2*marginX)),
		ui.MaxHeight(max(float32(0), vp.H-2*marginY)),
	}
}

func floatPanelC(p floatProps) *ui.Node {
	enter := ui.UseEnterTransition(100)
	pref, position := useFloatingPosition(p.anchor, 6, false, enter)
	vp := ui.Viewport()
	panel := append([]ui.StyleOpt{
		ui.Bg(p.th.Popover), ui.TextColor(p.th.PopoverForeground),
		ui.Border(1, p.th.Border), ui.Radius(p.th.Radius), shadowMd(),
	}, p.extra...)
	panel = append(panel, ui.MinWidth(min(p.minWidth, max(float32(0), vp.W-16))))
	panel = append(panel, position...)

	return ui.Portal(
		ui.Div(
			ui.Style(ui.Grow(1)),
			ui.OnClick(p.onClose),
			ui.ScrollView(append([]*ui.Node{ui.Style(panel...), pref, ui.OnClick(func() {})}, p.content...)...),
		),
	)
}

// ---- 菜单行（悬停高亮的可点击项，供 Select / DropdownMenu 复用）----

type menuRowProps struct {
	label   string
	onClick func()
}

func menuRow(label string, onClick func()) *ui.Node {
	return ui.Use(menuRowImpl, menuRowProps{label: label, onClick: onClick})
}

func menuRowImpl(p menuRowProps) *ui.Node {
	th := ui.UseTheme()
	hovered, _, ia := ui.UseInteraction()
	st := []ui.StyleOpt{ui.Row, ui.ItemsCenter, ui.PaddingXY(10, 7), ui.Radius(th.Radius - 2)}
	if hovered {
		st = append(st, ui.Bg(th.Accent), ui.TextColor(th.AccentForeground))
	}
	return ui.Div(ui.Style(st...), ui.OnClick(p.onClick), ia,
		ui.Text(p.label, ui.FontSize(14)))
}

// ---- Popover ----

type popoverProps struct {
	trigger *ui.Node
	content []*ui.Node
}

// Popover 点击 trigger 展开一个锚定在其下方的浮层。
func Popover(trigger *ui.Node, content ...*ui.Node) *ui.Node {
	return ui.Use(popover, popoverProps{trigger: trigger, content: content})
}

func popover(p popoverProps) *ui.Node {
	th := ui.UseTheme()
	open, setOpen := ui.UseState(false)
	ref, rect := ui.UseMeasure()
	ui.UseEscape(open, func() { setOpen(false) })
	return ui.Fragment(
		ui.Div(ref, ui.OnClick(func() { setOpen(!open) }), p.trigger),
		ui.If(open, floatPanel(th, rect, func() { setOpen(false) },
			rect.W, []ui.StyleOpt{ui.Column, ui.Gap(8), ui.Padding(14)}, p.content...)),
	)
}

// ---- Tooltip ----

type tooltipProps struct {
	text    string
	trigger *ui.Node
}

// Tooltip 悬停 trigger 时在其上方显示提示文本。
func Tooltip(text string, trigger *ui.Node) *ui.Node {
	return ui.Use(tooltip, tooltipProps{text: text, trigger: trigger})
}

func tooltip(p tooltipProps) *ui.Node {
	th := ui.UseTheme()
	hovered, _, ia := ui.UseInteraction()
	ref, rect := ui.UseMeasure()
	return ui.Fragment(
		ui.Div(ref, ia, p.trigger),
		ui.If(hovered, ui.Use(tooltipPanel, tooltipPanelProps{th: th, anchor: rect, text: p.text})),
	)
}

type tooltipPanelProps struct {
	th     ui.Theme
	anchor ui.Rect
	text   string
}

func tooltipPanel(p tooltipPanelProps) *ui.Node {
	ref, position := useFloatingPosition(p.anchor, 6, true, 1)
	style := append([]ui.StyleOpt{ui.Bg(p.th.Foreground), ui.Radius(6), ui.PaddingXY(10, 5), ui.Clip}, position...)
	return ui.Portal(ui.Div(ref, ui.Style(style...),
		ui.Text(p.text, ui.FontSize(12), ui.TextColor(p.th.Background))))
}

// ---- HoverCard ----

type hoverCardProps struct {
	trigger *ui.Node
	content []*ui.Node
}

// HoverCard 悬停 trigger 时在其下方展示一张卡片（悬停卡片本身也保持展开）。
func HoverCard(trigger *ui.Node, content ...*ui.Node) *ui.Node {
	return ui.Use(hoverCard, hoverCardProps{trigger: trigger, content: content})
}

func hoverCard(p hoverCardProps) *ui.Node {
	th := ui.UseTheme()
	tHover, setTHover := ui.UseState(false)
	cHover, setCHover := ui.UseState(false)
	ref, rect := ui.UseMeasure()
	open := tHover || cHover
	// 离开触发器后保留短暂退场，让指针有时间跨过触发器与卡片之间的空隙。
	mounted, progress := ui.UseTransition(open, 120)
	return ui.Fragment(
		ui.Div(ref, ui.OnHover(setTHover), p.trigger),
		ui.If(mounted, ui.Use(hoverCardPanel, hoverCardPanelProps{
			th: th, anchor: rect, progress: progress, onHover: setCHover, content: p.content,
		})),
	)
}

type hoverCardPanelProps struct {
	th       ui.Theme
	anchor   ui.Rect
	progress float32
	onHover  func(bool)
	content  []*ui.Node
}

func hoverCardPanel(p hoverCardPanelProps) *ui.Node {
	ref, position := useFloatingPosition(p.anchor, 6, false, p.progress)
	style := []ui.StyleOpt{
		ui.Column, ui.Gap(6), ui.Padding(14), ui.MinWidth(min(float32(240), max(float32(0), ui.Viewport().W-16))),
		ui.Bg(p.th.Popover), ui.TextColor(p.th.PopoverForeground),
		ui.Border(1, p.th.Border), ui.Radius(p.th.Radius), shadowMd(),
	}
	style = append(style, position...)
	return ui.Portal(ui.ScrollView(append([]*ui.Node{ref, ui.Style(style...), ui.OnHover(p.onHover)}, p.content...)...))
}

// ---- DropdownMenu ----

type MenuItem struct {
	Label    string
	OnSelect func()
}

type dropdownProps struct {
	trigger *ui.Node
	items   []MenuItem
}

// DropdownMenu 点击 trigger 展开菜单项列表。
func DropdownMenu(trigger *ui.Node, items []MenuItem) *ui.Node {
	return ui.Use(dropdown, dropdownProps{trigger: trigger, items: items})
}

func dropdown(p dropdownProps) *ui.Node {
	th := ui.UseTheme()
	open, setOpen := ui.UseState(false)
	ref, rect := ui.UseMeasure()
	ui.UseEscape(open, func() { setOpen(false) })
	rows := make([]*ui.Node, 0, len(p.items))
	for _, it := range p.items {
		item := it
		rows = append(rows, menuRow(item.Label, func() {
			if item.OnSelect != nil {
				item.OnSelect()
			}
			setOpen(false)
		}))
	}
	// 上下方向键在菜单项间移动焦点（ArrowNav 属性会落到浮层面板的 Div 上）
	content := append([]*ui.Node{ui.ArrowNav(ui.NavVertical)}, rows...)
	return ui.Fragment(
		ui.Div(ref, ui.OnClick(func() { setOpen(!open) }), p.trigger),
		ui.If(open, floatPanel(th, rect, func() { setOpen(false) },
			max(rect.W, 160), []ui.StyleOpt{ui.Column, ui.Padding(4)}, content...)),
	)
}
