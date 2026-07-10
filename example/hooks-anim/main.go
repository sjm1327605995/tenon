package main

import (
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func App(_ struct{}) *ui.Node {
	open, setOpen := ui.UseState(false)

	target := float32(0)
	if open {
		target = 150
	}
	h := ui.UseTween(target, 280, ui.EaseInOut) // 面板高度平滑过渡
	op := h / 150                               // 内容透明度随高度淡入淡出

	label := "展开 ▼"
	if open {
		label = "收起 ▲"
	}

	return ui.Div(
		ui.Style(ui.Column, ui.Width(440), ui.Height(380), ui.Bg(ui.Hex("#0f172a")),
			ui.Padding(28), ui.Gap(18), ui.ItemsCenter),

		ui.Text("UseTween 折叠动画", ui.FontSize(24), ui.TextColor(ui.White)),

		ui.Button(
			ui.Style(ui.PaddingXY(22, 11), ui.Bg(ui.Hex("#38bdf8")), ui.Radius(8),
				ui.ItemsCenter, ui.JustifyCenter),
			ui.OnClick(func() { setOpen(!open) }),
			ui.Text(label, ui.FontSize(15), ui.TextColor(ui.Hex("#0f172a"))),
		),

		// 折叠面板：高度用补间驱动，Clip 裁掉溢出，Opacity 让内容淡入
		ui.Div(
			ui.Style(ui.Width(380), ui.Height(h), ui.Bg(ui.Hex("#1e293b")), ui.Radius(12),
				ui.Border(1, ui.Hex("#334155")), ui.Clip, ui.Padding(18), ui.Opacity(op),
				ui.Column, ui.Gap(8)),
			ui.Text("平滑过渡的折叠内容", ui.FontSize(16), ui.TextColor(ui.White)),
			ui.Text("高度用 EaseInOut 缓动，透明度随之淡入淡出。",
				ui.FontSize(14), ui.TextColor(ui.Hex("#94a3b8"))),
			ui.Text("动画期间引擎每帧只重渲染这个组件。",
				ui.FontSize(14), ui.TextColor(ui.Hex("#94a3b8"))),
		),
	)
}

func main() {
	ui.Run(ui.Use(App, struct{}{}))
}
