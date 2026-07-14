package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type ButtonGroupItem struct {
	Label   string
	OnClick func()
	Active  bool // 选中态（分段控件）
}

// ButtonGroup 是相连的分段按钮组：外圈统一圆角边框，内部以竖线分隔，选中项高亮。
//
//	shadcn.ButtonGroup([]shadcn.ButtonGroupItem{
//	    {Label: "Day", Active: true}, {Label: "Week"}, {Label: "Month"},
//	})
func ButtonGroup(items []ButtonGroupItem) *ui.Node { return ui.Use(buttonGroup, items) }

func buttonGroup(items []ButtonGroupItem) *ui.Node {
	th := ui.UseTheme()
	kids := []*ui.Node{ui.Style(ui.Row, ui.Height(36), // 默认 align-items:stretch，项撑满行高
		ui.Border(1, th.Border), ui.Radius(radiusMd(th)), ui.Clip)}
	for i, it := range items {
		if i > 0 { // 项间竖直分隔线
			kids = append(kids, ui.Div(ui.Style(ui.Width(1), ui.HeightPct(100), ui.Bg(th.Border))))
		}
		kids = append(kids, ui.Use(bgItem, it))
	}
	return ui.Div(kids...)
}

func bgItem(it ButtonGroupItem) *ui.Node {
	th := ui.UseTheme()
	hovered, _, ia := ui.UseInteraction()

	st := []ui.StyleOpt{ui.Row, ui.ItemsCenter, ui.JustifyCenter, ui.Grow(1), ui.PaddingXY(14, 0)}
	fg := th.Foreground
	switch {
	case it.Active:
		st = append(st, ui.Bg(th.Primary))
		fg = th.PrimaryForeground
	case hovered:
		st = append(st, ui.Bg(th.Accent))
		fg = th.AccentForeground
	}
	return ui.Div(ui.Style(st...), ui.OnClick(it.OnClick), ia,
		ui.Text(it.Label, ui.FontSize(14), ui.Medium, ui.TextColor(fg)))
}
