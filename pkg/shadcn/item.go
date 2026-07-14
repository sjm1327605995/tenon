package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type ItemProps struct {
	Media       *ui.Node // 左侧图标/头像，可为 nil
	Title       string
	Description string
	Trailing    *ui.Node // 右侧操作/元信息，可为 nil
	OnClick     func()   // 非 nil 时整行可点、悬停高亮
}

// Item 是列表项基元：媒体 + 标题/描述 + 尾部操作。用于设置项、列表、菜单等。
func Item(p ItemProps) *ui.Node { return ui.Use(item, p) }

func item(p ItemProps) *ui.Node {
	th := ui.UseTheme()
	hovered, _, ia := ui.UseInteraction()

	st := []ui.StyleOpt{ui.Row, ui.ItemsCenter, ui.Gap(12), ui.PaddingXY(12, 10), ui.Radius(radiusMd(th))}
	if p.OnClick != nil && hovered {
		st = append(st, ui.Bg(th.Accent))
	}
	kids := []*ui.Node{ui.Style(st...)}
	if p.OnClick != nil {
		kids = append(kids, ui.OnClick(p.OnClick), ia)
	}
	if p.Media != nil {
		kids = append(kids, p.Media)
	}
	col := []*ui.Node{ui.Style(ui.Column, ui.Gap(2), ui.Grow(1))}
	if p.Title != "" {
		col = append(col, ui.Text(p.Title, ui.FontSize(14), ui.Medium, ui.TextColor(th.Foreground)))
	}
	if p.Description != "" {
		col = append(col, ui.Text(p.Description, ui.FontSize(13), ui.TextColor(th.MutedForeground)))
	}
	kids = append(kids, ui.Div(col...))
	if p.Trailing != nil {
		kids = append(kids, p.Trailing)
	}
	return ui.Div(kids...)
}
