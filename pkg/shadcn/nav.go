package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type NavItem struct {
	Label    string
	OnSelect func()     // 无子项时的点击回调
	Items    []MenuItem // 有子项时展开下拉
}

type navProps struct{ items []NavItem }

// NavigationMenu 是横向导航条；带子项的项点击展开下拉菜单。
func NavigationMenu(items []NavItem) *ui.Node { return ui.Use(navMenu, navProps{items: items}) }

func navMenu(p navProps) *ui.Node {
	th := ui.UseTheme()
	kids := []*ui.Node{ui.Style(ui.Row, ui.Gap(4), ui.Padding(4), ui.Radius(th.Radius), ui.Bg(th.Muted))}
	for _, it := range p.items {
		item := it
		label := ui.Div(
			ui.Style(ui.PaddingXY(14, 7), ui.Radius(th.Radius-2), ui.ItemsCenter, ui.JustifyCenter),
			ui.Text(item.Label, ui.FontSize(14)),
		)
		if len(item.Items) > 0 {
			kids = append(kids, DropdownMenu(label, item.Items))
		} else {
			onSel := item.OnSelect
			kids = append(kids, ui.Div(
				ui.Style(ui.PaddingXY(14, 7), ui.Radius(th.Radius-2), ui.ItemsCenter, ui.JustifyCenter),
				ui.OnClick(onSel),
				ui.Text(item.Label, ui.FontSize(14)),
			))
		}
	}
	return ui.Div(kids...)
}
