package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type MenubarMenu struct {
	Label string
	Items []MenuItem
}

// Menubar 是应用菜单栏：一排菜单标题，每个点击展开下拉菜单（复用 DropdownMenu）。
func Menubar(menus []MenubarMenu) *ui.Node { return ui.Use(menubar, menus) }

func menubar(menus []MenubarMenu) *ui.Node {
	th := ui.UseTheme()
	kids := []*ui.Node{ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(2), ui.Padding(3),
		ui.Border(1, th.Border), ui.Radius(radiusMd(th)), ui.Bg(th.Background))}
	for _, m := range menus {
		kids = append(kids, DropdownMenu(ui.Use(menubarTrigger, m.Label), m.Items))
	}
	return ui.Div(kids...)
}

func menubarTrigger(label string) *ui.Node {
	th := ui.UseTheme()
	hovered, _, ia := ui.UseInteraction()
	st := []ui.StyleOpt{ui.Row, ui.ItemsCenter, ui.PaddingXY(10, 5), ui.Radius(radiusSm(th))}
	if hovered {
		st = append(st, ui.Bg(th.Accent))
	}
	return ui.Div(ui.Style(st...), ia, ui.Text(label, ui.FontSize(14), ui.Medium))
}
