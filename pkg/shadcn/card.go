package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type nodesProps struct{ children []*ui.Node }

// Card 是带边框的卡片容器。
func Card(children ...*ui.Node) *ui.Node { return ui.Use(card, nodesProps{children}) }

func card(p nodesProps) *ui.Node {
	th := ui.UseTheme()
	// shadcn v4: rounded-xl border bg-card shadow-sm flex-col gap-6 p-6。
	// 内边距挂在容器上（而非各分区），因此 Card(任意内容) 与 Card(CardHeader(...), ...)
	// 都自带正确内边距，无需开发者记得配合分区组件。
	base := ui.Style(ui.Column, ui.Gap(24), ui.Padding(24),
		ui.Bg(th.Card), ui.TextColor(th.CardForeground),
		ui.Border(1, th.Border), ui.Radius(radiusXl(th)), shadowSm())
	return ui.Div(append([]*ui.Node{base}, p.children...)...)
}

// CardHeader / CardContent / CardFooter 是可选的布局分区（纵向分组，间距由 Card 的 gap-6 提供，
// 水平内边距由 Card 容器统一提供）。
func CardHeader(children ...*ui.Node) *ui.Node {
	return ui.Div(append([]*ui.Node{ui.Style(ui.Column, ui.Gap(6))}, children...)...)
}

func CardContent(children ...*ui.Node) *ui.Node {
	return ui.Div(append([]*ui.Node{ui.Style(ui.Column)}, children...)...)
}

func CardFooter(children ...*ui.Node) *ui.Node {
	return ui.Div(append([]*ui.Node{ui.Style(ui.Row, ui.Gap(8), ui.ItemsCenter)}, children...)...)
}

// CardTitle 继承 Card 前景色（text-base font-semibold）。
func CardTitle(text string) *ui.Node { return ui.Text(text, ui.FontSize(16), ui.Semibold) }

// CardDescription 使用弱化前景色。
func CardDescription(text string) *ui.Node { return ui.Use(mutedText, text) }

func mutedText(text string) *ui.Node {
	th := ui.UseTheme()
	return ui.Text(text, ui.FontSize(14), ui.TextColor(th.MutedForeground))
}
