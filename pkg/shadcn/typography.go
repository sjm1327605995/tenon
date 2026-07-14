package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

// Typography：与 shadcn/ui 文档一致的排版辅助。颜色默认继承容器前景色，
// muted 类用主题的 MutedForeground。

func H1(s string) *ui.Node { return ui.Text(s, ui.FontSize(36), ui.Bold) }
func H2(s string) *ui.Node { return ui.Text(s, ui.FontSize(30), ui.Semibold) }
func H3(s string) *ui.Node { return ui.Text(s, ui.FontSize(24), ui.Semibold) }
func H4(s string) *ui.Node { return ui.Text(s, ui.FontSize(20), ui.Semibold) }

func P(s string) *ui.Node     { return ui.Text(s, ui.FontSize(16)) }
func Large(s string) *ui.Node { return ui.Text(s, ui.FontSize(18), ui.Semibold) }
func Small(s string) *ui.Node { return ui.Text(s, ui.FontSize(14), ui.Medium) }

// Lead 是引导段落（较大、次要色）。
func Lead(s string) *ui.Node { return ui.Use(leadText, s) }

func leadText(s string) *ui.Node {
	th := ui.UseTheme()
	return ui.Text(s, ui.FontSize(20), ui.TextColor(th.MutedForeground))
}

// Muted 是次要说明文字。
func Muted(s string) *ui.Node { return ui.Use(mutedText, s) }

// InlineCode 是行内代码片段（次要背景 + 圆角）。
func InlineCode(s string) *ui.Node { return ui.Use(inlineCode, s) }

func inlineCode(s string) *ui.Node {
	th := ui.UseTheme()
	return ui.Div(
		ui.Style(ui.Row, ui.ItemsCenter, ui.PaddingXY(6, 2), ui.Radius(radiusSm(th)), ui.Bg(th.Muted)),
		ui.Text(s, ui.FontSize(14), ui.TextColor(th.Foreground)),
	)
}

// Blockquote 是引用块（左侧竖线 + 斜体次要色）。
func Blockquote(s string) *ui.Node { return ui.Use(blockquote, s) }

func blockquote(s string) *ui.Node {
	th := ui.UseTheme()
	return ui.Div(
		ui.Style(ui.Row, ui.PaddingXY(14, 6)),
		ui.Div(ui.Style(ui.Width(3), ui.Radius(2), ui.Bg(th.Border))),
		ui.Div(ui.Style(ui.PaddingXY(12, 0)),
			ui.Text(s, ui.FontSize(16), ui.Italic, ui.TextColor(th.MutedForeground))),
	)
}
