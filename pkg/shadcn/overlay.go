package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

// floatPanel 在锚点矩形附近渲染浮层（Portal）：默认在下方展开，空间不足时贴边翻转到上方。
// 点击浮层外区域触发 onClose。
func floatPanel(th ui.Theme, anchor ui.Rect, onClose func(), extra []ui.StyleOpt, content ...*ui.Node) *ui.Node {
	return ui.Use(floatPanelC, floatProps{th: th, anchor: anchor, onClose: onClose, extra: extra, content: content})
}

type floatProps struct {
	th      ui.Theme
	anchor  ui.Rect
	onClose func()
	extra   []ui.StyleOpt
	content []*ui.Node
}

func floatPanelC(p floatProps) *ui.Node {
	pref, prect := ui.UseMeasure() // 测量浮层自身高度以决定翻转
	vp := ui.Viewport()

	top := p.anchor.Y + p.anchor.H + 6
	if prect.H > 0 && top+prect.H > vp.H && p.anchor.Y-prect.H-6 >= 0 {
		top = p.anchor.Y - prect.H - 6 // 下方放不下 -> 翻到上方
	}
	op := float32(1)
	if prect.H == 0 {
		op = 0 // 首帧测量前隐藏，避免位置跳动
	}

	panel := append([]ui.StyleOpt{
		ui.Absolute, ui.Left(p.anchor.X), ui.Top(top), ui.Opacity(op),
		ui.Bg(p.th.Popover), ui.TextColor(p.th.PopoverForeground),
		ui.Border(1, p.th.Border), ui.Radius(p.th.Radius),
	}, p.extra...)

	return ui.Portal(
		ui.Div(
			ui.Style(ui.Grow(1)),
			ui.OnClick(p.onClose),
			ui.Div(append([]*ui.Node{ui.Style(panel...), pref, ui.OnClick(func() {})}, p.content...)...),
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
			[]ui.StyleOpt{ui.Column, ui.Gap(8), ui.Padding(14), ui.MinWidth(rect.W)}, p.content...)),
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
		ui.If(hovered, ui.Portal(
			ui.Div(
				ui.Style(ui.Absolute, ui.Left(rect.X), ui.Top(rect.Y-32),
					ui.Bg(th.Foreground), ui.Radius(6), ui.PaddingXY(10, 5)),
				ui.Text(p.text, ui.FontSize(12), ui.TextColor(th.Background)),
			),
		)),
	)
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
	return ui.Fragment(
		ui.Div(ref, ui.OnHover(setTHover), p.trigger),
		ui.If(open, ui.Portal(
			ui.Div(append([]*ui.Node{
				ui.Style(ui.Absolute, ui.Left(rect.X), ui.Top(rect.Y+rect.H+2),
					ui.Column, ui.Gap(6), ui.Padding(14), ui.MinWidth(240),
					ui.Bg(th.Popover), ui.TextColor(th.PopoverForeground),
					ui.Border(1, th.Border), ui.Radius(th.Radius)),
				ui.OnHover(setCHover),
			}, p.content...)...),
		)),
	)
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
	return ui.Fragment(
		ui.Div(ref, ui.OnClick(func() { setOpen(!open) }), p.trigger),
		ui.If(open, floatPanel(th, rect, func() { setOpen(false) },
			[]ui.StyleOpt{ui.Column, ui.Padding(4), ui.MinWidth(max(rect.W, 160))}, rows...)),
	)
}
