package shadcn

import (
	"strings"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// ComboboxOption 是一个可选项：Value 用于回调/选中判断，Label 用于显示与搜索。
type ComboboxOption struct {
	Value string
	Label string
}

type ComboboxProps struct {
	Options           []ComboboxOption
	Value             string // 当前选中值（受控）
	OnChange          func(string)
	Placeholder       string  // 未选中时触发按钮的占位文字
	SearchPlaceholder string  // 搜索框占位文字
	Empty             string  // 无匹配时的提示
	Width             float32 // 触发按钮最小宽度（0 用默认）
}

// Combobox 是可搜索下拉：点击展开，输入即时按标签过滤选项，选中回填并打勾。
// 键盘：搜索框内按 ↓ 进入列表，↑↓ 在选项间移动（ArrowNav），Enter 选中，Esc 关闭。
func Combobox(p ComboboxProps) *ui.Node { return ui.Use(combobox, p) }

func combobox(p ComboboxProps) *ui.Node {
	th := ui.UseTheme()
	open, setOpen := ui.UseState(false)
	query, setQuery := ui.UseState("")
	ref, rect := ui.UseMeasure()
	ui.UseEscape(open, func() { setOpen(false) })
	// 关闭时清空搜索，下次打开是干净的
	ui.UseEffect(func() ui.Cleanup {
		if !open {
			setQuery("")
		}
		return nil
	}, open)

	// 触发按钮文字：选中项的 Label，否则占位符
	label, labelColor := p.Placeholder, th.MutedForeground
	if label == "" {
		label = "Select…"
	}
	for _, o := range p.Options {
		if o.Value == p.Value {
			label, labelColor = o.Label, th.Foreground
			break
		}
	}

	minW := p.Width
	if minW <= 0 {
		minW = 220
	}
	trigger := ui.Div(ref,
		ui.Style(ui.Row, ui.ItemsCenter, ui.JustifyBetween, ui.Gap(8), ui.Height(38),
			ui.PaddingXY(12, 0), ui.Radius(th.Radius), ui.Border(1, th.Input),
			ui.Bg(th.Background), ui.MinWidth(minW)),
		ui.OnClick(func() { setOpen(!open) }),
		ui.Text(label, ui.FontSize(14), ui.TextColor(labelColor)),
		ui.Icon(ui.IconChevronDown, 16, ui.TextColor(th.MutedForeground)),
	)

	searchPh := p.SearchPlaceholder
	if searchPh == "" {
		searchPh = "Search…"
	}
	empty := p.Empty
	if empty == "" {
		empty = "No results."
	}

	// 过滤（大小写不敏感的子串匹配）后生成选项行
	rows := []*ui.Node{ui.Style(ui.Column)}
	matches := 0
	for _, o := range p.Options {
		if query != "" && !strings.Contains(strings.ToLower(o.Label), strings.ToLower(query)) {
			continue
		}
		matches++
		opt := o
		rows = append(rows, comboRow(opt.Label, opt.Value == p.Value, func() {
			if p.OnChange != nil {
				p.OnChange(opt.Value)
			}
			setOpen(false)
		}))
	}
	if matches == 0 {
		rows = append(rows, ui.Div(ui.Style(ui.PaddingXY(10, 8)),
			ui.Text(empty, ui.FontSize(14), ui.TextColor(th.MutedForeground))))
	}

	// 面板内容：整块设为纵向 ArrowNav 组（搜索框 ↓ 进列表、行间 ↑↓ 移动）
	content := []*ui.Node{
		ui.ArrowNav(ui.NavVertical),
		ui.Input(
			ui.Style(ui.Height(36), ui.PaddingXY(8, 0), ui.Bg(th.Popover), ui.TextColor(th.Foreground),
				ui.FontSize(14), ui.ItemsCenter),
			ui.Value(query), ui.OnChange(setQuery), ui.Placeholder(searchPh),
		),
		ui.Div(ui.Style(ui.Height(1), ui.Bg(th.Border))),
		ui.ScrollView(append([]*ui.Node{ui.Style(ui.Column, ui.MaxHeight(240), ui.PaddingXY(0, 4))}, rows...)...),
	}

	return ui.Fragment(
		trigger,
		ui.If(open, floatPanel(th, rect, func() { setOpen(false) },
			[]ui.StyleOpt{ui.Column, ui.Padding(4), ui.Gap(4), ui.MinWidth(max(rect.W, minW))}, content...)),
	)
}

// comboRow 是一行选项：左侧对号列（选中显示 ✓）+ 标签，悬停高亮。
type comboRowProps struct {
	label    string
	selected bool
	onClick  func()
}

func comboRow(label string, selected bool, onClick func()) *ui.Node {
	return ui.Use(comboRowImpl, comboRowProps{label: label, selected: selected, onClick: onClick})
}

func comboRowImpl(p comboRowProps) *ui.Node {
	th := ui.UseTheme()
	hovered, _, ia := ui.UseInteraction()
	st := []ui.StyleOpt{ui.Row, ui.ItemsCenter, ui.Gap(8), ui.PaddingXY(8, 7), ui.Radius(th.Radius - 2)}
	if hovered {
		st = append(st, ui.Bg(th.Accent), ui.TextColor(th.AccentForeground))
	}
	return ui.Div(ui.Style(st...), ui.OnClick(p.onClick), ia,
		ui.Div(ui.Style(ui.Width(16), ui.ItemsCenter),
			ui.If(p.selected, ui.Icon(ui.IconCheck, 14))),
		ui.Text(p.label, ui.FontSize(14)),
	)
}
