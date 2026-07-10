package main

import (
	"fmt"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// Counter 是一个函数组件，用 hooks 管理本地状态；无需手动声明重渲染。
func Counter(_ struct{}) *ui.Node {
	count, setCount := ui.UseState(0)

	ui.UseEffect(func() ui.Cleanup {
		fmt.Printf("count -> %d\n", count)
		return nil
	}, count)

	btn := func(label string, bg, fg ui.Color, onClick func()) *ui.Node {
		return ui.Button(
			ui.Style(ui.Width(56), ui.Height(44), ui.Bg(bg), ui.Radius(8),
				ui.ItemsCenter, ui.JustifyCenter),
			ui.OnClick(onClick),
			ui.Text(label, ui.FontSize(24), ui.TextColor(fg)),
		)
	}

	// 外层用 Fill 填满视口并居中，窗口缩放时卡片始终居中（自适应）
	return ui.Div(
		ui.Style(ui.Fill, ui.ItemsCenter, ui.JustifyCenter, ui.Bg(ui.Hex("#eef2f7"))),

		ui.Div(
			ui.Style(ui.Column, ui.Gap(24), ui.Padding(32), ui.ItemsCenter,
				ui.Width(600), ui.Height(400), ui.Bg(ui.White), ui.Radius(16),
				ui.Border(1, ui.LightGray)),

			ui.Text("Tenon Hooks Counter", ui.FontSize(28), ui.TextColor(ui.Hex("#1f2937"))),

			ui.Div(
				ui.Style(ui.Row, ui.Gap(20), ui.ItemsCenter),
				btn("-", ui.LightGray, ui.DarkGray, func() { setCount(count - 1) }),
				ui.Text(fmt.Sprintf("%d", count), ui.FontSize(40), ui.TextColor(ui.Blue)),
				btn("+", ui.Blue, ui.White, func() { setCount(count + 1) }),
			),

			ui.Text("缩放窗口：卡片始终居中；仅点击时该组件重渲染", ui.FontSize(14), ui.TextColor(ui.Gray)),
		),
	)
}

func main() {
	ui.Run(ui.Use(Counter, struct{}{}))
}
