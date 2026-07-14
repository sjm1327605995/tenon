package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

// Kbd 渲染一个键帽（用于展示键盘快捷键，如 ⌘ / K / Esc）。多个键用 KbdGroup 组合。
func Kbd(text string) *ui.Node { return ui.Use(kbd, text) }

func kbd(text string) *ui.Node {
	th := ui.UseTheme()
	return ui.Div(
		ui.Style(ui.Row, ui.ItemsCenter, ui.JustifyCenter, ui.Height(20), ui.MinWidth(20),
			ui.PaddingXY(5, 0), ui.Radius(radiusSm(th)), ui.Bg(th.Muted), ui.Border(1, th.Border)),
		ui.Text(text, ui.FontSize(11), ui.TextColor(th.MutedForeground), ui.Medium),
	)
}

// KbdGroup 把多个键帽横向排列（如 ⌘ + K）。
func KbdGroup(keys ...string) *ui.Node {
	kids := []*ui.Node{ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(4))}
	for _, k := range keys {
		kids = append(kids, Kbd(k))
	}
	return ui.Div(kids...)
}
