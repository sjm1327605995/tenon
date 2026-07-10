package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type nodesProps struct{ children []*ui.Node }

// Card 是带边框的卡片容器。
func Card(children ...*ui.Node) *ui.Node { return ui.Use(card, nodesProps{children}) }

func card(p nodesProps) *ui.Node {
	th := ui.UseTheme()
	base := ui.Style(ui.Column, ui.Bg(th.Card), ui.TextColor(th.CardForeground),
		ui.Border(1, th.Border), ui.Radius(th.Radius+4))
	return ui.Div(append([]*ui.Node{base}, p.children...)...)
}

// CardHeader / CardContent / CardFooter 是布局分区。
func CardHeader(children ...*ui.Node) *ui.Node {
	return ui.Div(append([]*ui.Node{ui.Style(ui.Column, ui.Gap(6), ui.Padding(24))}, children...)...)
}

func CardContent(children ...*ui.Node) *ui.Node {
	return ui.Div(append([]*ui.Node{ui.Style(ui.Column, ui.Gap(8), ui.PaddingXY(24, 0))}, children...)...)
}

func CardFooter(children ...*ui.Node) *ui.Node {
	return ui.Div(append([]*ui.Node{ui.Style(ui.Row, ui.Gap(8), ui.ItemsCenter, ui.Padding(24))}, children...)...)
}

// CardTitle 继承 Card 前景色，字号加大。
func CardTitle(text string) *ui.Node { return ui.Text(text, ui.FontSize(18)) }

// CardDescription 使用弱化前景色。
func CardDescription(text string) *ui.Node { return ui.Use(mutedText, text) }

func mutedText(text string) *ui.Node {
	th := ui.UseTheme()
	return ui.Text(text, ui.FontSize(14), ui.TextColor(th.MutedForeground))
}
