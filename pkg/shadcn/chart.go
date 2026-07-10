package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type BarChartProps struct {
	Data          []float32
	Labels        []string
	Width, Height float32
}

// BarChart 是简单柱状图。
func BarChart(p BarChartProps) *ui.Node { return ui.Use(barChart, p) }

func barChart(p BarChartProps) *ui.Node {
	th := ui.UseTheme()
	maxV := float32(1)
	for _, v := range p.Data {
		if v > maxV {
			maxV = v
		}
	}
	bars := []*ui.Node{ui.Style(ui.Row, ui.Gap(10), ui.Width(p.Width), ui.Height(p.Height))}
	for i, v := range p.Data {
		h := v / maxV * (p.Height - 26)
		label := ""
		if i < len(p.Labels) {
			label = p.Labels[i]
		}
		bars = append(bars, ui.Div(
			ui.Style(ui.Column, ui.Grow(1), ui.JustifyEnd, ui.ItemsCenter, ui.Gap(4)),
			ui.Div(ui.Style(ui.Width(26), ui.Height(h), ui.Radius(4), ui.Bg(th.Primary))),
			ui.Text(label, ui.FontSize(11), ui.TextColor(th.MutedForeground)),
		))
	}
	return ui.Div(bars...)
}
