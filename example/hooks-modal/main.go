package main

import (
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// 模态对话框：用 Portal 渲染到顶层浮层，脱离父级裁剪；点击遮罩或按钮关闭。
func App(_ struct{}) *ui.Node {
	open, setOpen := ui.UseState(false)
	mounted, p := ui.UseTransition(open, 200) // 进出场：遮罩淡入淡出、卡片缩放

	return ui.Div(
		ui.Style(ui.Width(680), ui.Height(440), ui.Bg(ui.Hex("#0f172a")),
			ui.ItemsCenter, ui.JustifyCenter),

		ui.Button(
			ui.Style(ui.PaddingXY(24, 12), ui.Bg(ui.Hex("#3b82f6")), ui.Radius(8),
				ui.ItemsCenter, ui.JustifyCenter),
			ui.OnClick(func() { setOpen(true) }),
			ui.Text("打开对话框", ui.FontSize(16), ui.TextColor(ui.White)),
		),

		ui.If(mounted, ui.Portal(
			// 遮罩：铺满全屏、居中内容，点击关闭；透明度随进出场
			ui.Div(
				ui.Style(ui.Grow(1), ui.ItemsCenter, ui.JustifyCenter,
					ui.Bg(ui.Color{R: 0, G: 0, B: 0, A: 150}), ui.Opacity(p)),
				ui.OnClick(func() { setOpen(false) }),

				// 对话框卡片：吞掉点击，避免冒泡到遮罩；随进出场缩放
				ui.Div(
					ui.Style(ui.Width(340), ui.Height(210), ui.Bg(ui.White), ui.Radius(16),
						ui.Padding(24), ui.Column, ui.Gap(16), ui.ItemsCenter, ui.JustifyCenter,
						ui.Scale(0.9+0.1*p)),
					ui.OnClick(func() {}),
					ui.Text("这是一个模态对话框", ui.FontSize(19), ui.TextColor(ui.Hex("#111827"))),
					ui.Text("Portal 渲染到顶层，点击遮罩或按钮关闭。",
						ui.FontSize(14), ui.TextColor(ui.Hex("#6b7280"))),
					ui.Button(
						ui.Style(ui.PaddingXY(22, 10), ui.Bg(ui.Hex("#3b82f6")), ui.Radius(8),
							ui.ItemsCenter, ui.JustifyCenter),
						ui.OnClick(func() { setOpen(false) }),
						ui.Text("关闭", ui.FontSize(15), ui.TextColor(ui.White)),
					),
				),
			),
		)),
	)
}

func main() {
	ui.Run(ui.Use(App, struct{}{}))
}
