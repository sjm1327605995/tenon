package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type contextMenuProps struct {
	trigger *ui.Node
	items   []MenuItem
}

// ContextMenu 在 trigger 区域右键时，于光标处弹出菜单（复用 MenuItem / menuRow）。
// 点空白处或 Esc 关闭，↑↓ 方向键在项间移动。
func ContextMenu(trigger *ui.Node, items []MenuItem) *ui.Node {
	return ui.Use(contextMenu, contextMenuProps{trigger: trigger, items: items})
}

func contextMenu(p contextMenuProps) *ui.Node {
	th := ui.UseTheme()
	open, setOpen := ui.UseState(false)
	pos, setPos := ui.UseState([2]float32{})
	ui.UseEscape(open, func() { setOpen(false) })

	wrapped := ui.Div(
		ui.OnContextMenu(func(x, y float32) { setPos([2]float32{x, y}); setOpen(true) }),
		p.trigger,
	)
	if !open {
		return ui.Fragment(wrapped)
	}

	rows := []*ui.Node{ui.ArrowNav(ui.NavVertical)}
	for _, it := range p.items {
		item := it
		rows = append(rows, menuRow(item.Label, func() {
			if item.OnSelect != nil {
				item.OnSelect()
			}
			setOpen(false)
		}))
	}
	panel := append([]*ui.Node{
		ui.Style(ui.Absolute, ui.Left(pos[0]), ui.Top(pos[1]), ui.Column, ui.Padding(4), ui.MinWidth(160),
			ui.Bg(th.Popover), ui.TextColor(th.PopoverForeground), ui.Border(1, th.Border),
			ui.Radius(th.Radius), shadowMd()),
		ui.OnClick(func() {}), // 吞掉点击，避免冒泡到遮罩关闭
	}, rows...)

	menu := ui.Portal(
		ui.Div(ui.Style(ui.Grow(1)), ui.OnClick(func() { setOpen(false) }), // 点空白关闭
			ui.Div(panel...)),
	)
	return ui.Fragment(wrapped, menu)
}
