package main

import (
	"fmt"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// 演示列表虚拟化：10 万行，只渲染视口附近的十几行，滚动流畅。
// 开 ui.ShowStats（F12 切换）可看到重绘/布局开销不随行数增长。
func App(_ struct{}) *ui.Node {
	const n = 100000
	return ui.Div(
		ui.Style(ui.Column, ui.Fill, ui.ItemsCenter, ui.JustifyCenter, ui.Bg(ui.Hex("#0f172a")),
			ui.Gap(12), ui.Padding(24), ui.TextColor(ui.Hex("#e2e8f0"))),

		ui.Text(fmt.Sprintf("VirtualList · %d 行 · 只渲染可视窗口", n), ui.FontSize(18), ui.Bold),

		ui.Div(ui.Style(ui.Width(420), ui.Height(440), ui.Bg(ui.Hex("#1e293b")),
			ui.Radius(12), ui.Border(1, ui.Hex("#334155")), ui.Padding(6), ui.Clip),
			ui.VirtualList(ui.VirtualListProps{
				Count: n, ItemHeight: 40, Height: 428,
				Render: func(i int) *ui.Node {
					bg := ui.Hex("#1e293b")
					if i%2 == 1 {
						bg = ui.Hex("#243449")
					}
					return ui.Div(
						ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(12), ui.Fill,
							ui.PaddingXY(14, 0), ui.Radius(8), ui.Bg(bg)),
						ui.Text(fmt.Sprintf("#%05d", i), ui.FontSize(13), ui.TextColor(ui.Hex("#94a3b8"))),
						ui.Text(fmt.Sprintf("列表项 %d", i), ui.FontSize(15)),
					)
				},
			}),
		),
	)
}

func main() {
	ui.ShowStats = true
	ui.Run(ui.Use(App, struct{}{}))
}
