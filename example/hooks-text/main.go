package main

import (
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// 演示文本自动折行 + 样式继承（根容器设 color/font-size，段落继承）。
func App(_ struct{}) *ui.Node {
	return ui.Div(
		ui.Style(ui.Column, ui.Width(600), ui.Height(460), ui.Bg(ui.Hex("#0f172a")),
			ui.Padding(36), ui.Gap(20),
			ui.TextColor(ui.Hex("#e2e8f0")), ui.FontSize(16)), // 下传给后代文本

		ui.Text("文本折行与样式继承", ui.FontSize(26), ui.TextColor(ui.White)), // 覆盖继承

		ui.Div(
			ui.Style(ui.Column, ui.Width(400), ui.Bg(ui.Hex("#1e293b")), ui.Radius(12),
				ui.Padding(22), ui.Gap(12)),

			ui.Text("这是一段较长的中文段落，用来演示自动折行：当宽度不足以容纳整行时，文字会在合适的位置断开并换到下一行，颜色和字号都从父容器继承而来。"),

			ui.Text("English paragraphs wrap at word boundaries when the available width is exceeded, flowing naturally onto multiple lines.",
				ui.TextColor(ui.Hex("#94a3b8"))), // 只覆盖颜色，字号仍继承
		),
	)
}

func main() {
	ui.Run(ui.Use(App, struct{}{}))
}
