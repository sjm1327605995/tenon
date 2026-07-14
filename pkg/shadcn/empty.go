package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type EmptyProps struct {
	Icon        *ui.Node // 可选：顶部图标/插画节点
	Title       string
	Description string
}

// Empty 是空状态占位：居中的图标 + 标题 + 描述 + 可选操作按钮。
//
//	shadcn.Empty(shadcn.EmptyProps{
//	    Icon: ui.Icon(ui.IconSearch, 28), Title: "没有结果", Description: "换个关键词试试。",
//	}, shadcn.Button(shadcn.ButtonProps{}, ui.Text("清除筛选")))
func Empty(p EmptyProps, actions ...*ui.Node) *ui.Node {
	return ui.Use(empty, emptyProps{p: p, actions: actions})
}

type emptyProps struct {
	p       EmptyProps
	actions []*ui.Node
}

func empty(ep emptyProps) *ui.Node {
	th := ui.UseTheme()
	p := ep.p
	kids := []*ui.Node{ui.Style(ui.Column, ui.ItemsCenter, ui.JustifyCenter, ui.Gap(8),
		ui.Padding(32), ui.TextColor(th.Foreground))}

	if p.Icon != nil {
		kids = append(kids, ui.Div(
			ui.Style(ui.Width(48), ui.Height(48), ui.ItemsCenter, ui.JustifyCenter,
				ui.Radius(radiusLg(th)), ui.Bg(th.Muted), ui.TextColor(th.MutedForeground)),
			p.Icon))
	}
	if p.Title != "" {
		kids = append(kids, ui.Text(p.Title, ui.FontSize(16), ui.Semibold))
	}
	if p.Description != "" {
		kids = append(kids, ui.Text(p.Description, ui.FontSize(14), ui.TextColor(th.MutedForeground)))
	}
	if len(ep.actions) > 0 {
		row := append([]*ui.Node{ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(8))}, ep.actions...)
		kids = append(kids, ui.Div(row...))
	}
	return ui.Div(kids...)
}
