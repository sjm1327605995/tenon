package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type SidebarItem struct {
	Label   string
	Icon    *ui.Node
	Active  bool
	OnClick func()
}

type SidebarGroup struct {
	Label string
	Items []SidebarItem
}

type SidebarProps struct {
	Header    *ui.Node // 顶部（logo/标题），可为 nil
	Groups    []SidebarGroup
	Footer    *ui.Node // 底部（用户/设置），可为 nil
	Collapsed bool     // 折叠为图标条
	Width     float32  // 展开宽度（默认 240）
}

// Sidebar 是可折叠的应用侧边栏：顶部 + 分组菜单 + 底部；折叠时只显示图标。
// 放在填满窗口的 Row 里，与主内容并列。
func Sidebar(p SidebarProps) *ui.Node { return ui.Use(sidebar, p) }

func sidebar(p SidebarProps) *ui.Node {
	th := ui.UseTheme()
	w := p.Width
	if w <= 0 {
		w = 240
	}
	if p.Collapsed {
		w = 60
	}
	kids := []*ui.Node{ui.Style(ui.Column, ui.Width(w), ui.HeightPct(100), ui.Gap(4),
		ui.Padding(8), ui.Bg(th.Card), ui.TextColor(th.CardForeground), ui.Clip)}

	if p.Header != nil {
		kids = append(kids, ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.PaddingXY(8, 10)), p.Header))
	}

	group := []*ui.Node{ui.Style(ui.Column, ui.Gap(2), ui.Grow(1))}
	for _, g := range p.Groups {
		if g.Label != "" && !p.Collapsed {
			group = append(group, ui.Div(ui.Style(ui.PaddingXY(8, 6)),
				ui.Text(g.Label, ui.FontSize(12), ui.Medium, ui.TextColor(th.MutedForeground))))
		}
		for _, it := range g.Items {
			group = append(group, ui.Use(sidebarItem, sidebarItemProps{item: it, collapsed: p.Collapsed}))
		}
	}
	kids = append(kids, ui.Div(group...))

	if p.Footer != nil {
		kids = append(kids, ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.PaddingXY(8, 8)), p.Footer))
	}
	// 面板 + 右侧 1px 分隔线（区分主内容；引擎无单边 border）
	return ui.Div(ui.Style(ui.Row, ui.HeightPct(100)),
		ui.Div(kids...),
		ui.Div(ui.Style(ui.Width(1), ui.HeightPct(100), ui.Bg(th.Border))),
	)
}

type sidebarItemProps struct {
	item      SidebarItem
	collapsed bool
}

func sidebarItem(p sidebarItemProps) *ui.Node {
	th := ui.UseTheme()
	hovered, _, ia := ui.UseInteraction()
	it := p.item

	st := []ui.StyleOpt{ui.Row, ui.ItemsCenter, ui.Gap(10), ui.Height(36), ui.PaddingXY(8, 0), ui.Radius(radiusMd(th))}
	fg := th.CardForeground
	switch {
	case it.Active:
		st = append(st, ui.Bg(th.Accent))
		fg = th.AccentForeground
	case hovered:
		st = append(st, ui.Bg(over(th.Accent, th.Card, 0.5)))
	}
	if p.collapsed {
		st = append(st, ui.JustifyCenter)
	}
	kids := []*ui.Node{ui.Style(st...), ui.OnClick(it.OnClick), ia}
	if it.Icon != nil {
		kids = append(kids, ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.TextColor(fg)), it.Icon))
	}
	if !p.collapsed {
		kids = append(kids, ui.Text(it.Label, ui.FontSize(14), ui.Medium, ui.TextColor(fg)))
	}
	return ui.Div(kids...)
}
