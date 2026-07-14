package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type InputGroupProps struct {
	Value       string
	OnChange    func(string)
	Placeholder string
	Leading     *ui.Node // 前置附加物（图标/文字/按钮），可为 nil
	Trailing    *ui.Node // 后置附加物，可为 nil
}

// InputGroup 是带前/后置附加物的输入框：外层是统一的边框容器，内部输入框无边框。
//
//	shadcn.InputGroup(shadcn.InputGroupProps{
//	    Leading:     ui.Icon(ui.IconSearch, 16),
//	    Placeholder: "搜索…",
//	    Trailing:    shadcn.Kbd("⌘K"),
//	})
func InputGroup(p InputGroupProps) *ui.Node { return ui.Use(inputGroup, p) }

func inputGroup(p InputGroupProps) *ui.Node {
	th := ui.UseTheme()
	kids := []*ui.Node{ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(8), ui.Height(36),
		ui.PaddingXY(10, 0), ui.Radius(radiusMd(th)), ui.Border(1, th.Input),
		ui.Bg(th.Background), shadowXs())}

	if p.Leading != nil {
		kids = append(kids, ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.TextColor(th.MutedForeground)), p.Leading))
	}
	// 无边框、透明、撑满的输入框
	kids = append(kids, ui.Input(
		ui.Style(ui.Grow(1), ui.ItemsCenter, ui.Bg(ui.Color{}), ui.TextColor(th.Foreground), ui.FontSize(14)),
		ui.Value(p.Value), ui.OnChange(p.OnChange), ui.Placeholder(p.Placeholder),
	))
	if p.Trailing != nil {
		kids = append(kids, ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.TextColor(th.MutedForeground)), p.Trailing))
	}
	return ui.Div(kids...)
}
