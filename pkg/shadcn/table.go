package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

// Table 是简单表格容器（等宽列：单元格用 Grow 平分）。
// 用法：Table(TableRow(TableHead("名称"), TableHead("邮箱")), TableRow(TableCell(...), TableCell(...)))
func Table(children ...*ui.Node) *ui.Node { return ui.Use(table, nodesProps{children}) }

func table(p nodesProps) *ui.Node {
	th := ui.UseTheme()
	base := ui.Style(ui.Column, ui.Border(1, th.Border), ui.Radius(th.Radius), ui.Clip)
	return ui.Div(append([]*ui.Node{base}, p.children...)...)
}

// TableRow 是一行，底部带分隔线。
func TableRow(children ...*ui.Node) *ui.Node { return ui.Use(tableRow, nodesProps{children}) }

func tableRow(p nodesProps) *ui.Node {
	th := ui.UseTheme()
	inner := ui.Div(append([]*ui.Node{ui.Style(ui.Row, ui.ItemsCenter)}, p.children...)...)
	return ui.Div(ui.Style(ui.Column),
		inner,
		ui.Div(ui.Style(ui.Height(1), ui.Bg(th.Border))),
	)
}

// TableHead 是表头单元格（弱化文字）。
func TableHead(text string) *ui.Node { return ui.Use(tableHead, text) }

func tableHead(text string) *ui.Node {
	th := ui.UseTheme()
	return ui.Div(ui.Style(ui.Grow(1), ui.PaddingXY(14, 10)),
		ui.Text(text, ui.FontSize(13), ui.TextColor(th.MutedForeground)))
}

// TableCell 是数据单元格。
func TableCell(children ...*ui.Node) *ui.Node {
	base := ui.Style(ui.Grow(1), ui.PaddingXY(14, 10))
	return ui.Div(append([]*ui.Node{base}, children...)...)
}
