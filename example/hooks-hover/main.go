package main

import (
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

type CardProps struct {
	Label string
	Color ui.Color
}

// Card：悬停时用弹性缩放（UseTween + Scale）放大，展示 hover + transform + 动画联动。
func Card(p CardProps) *ui.Node {
	hovered, setHovered := ui.UseState(false)

	target := float32(1.0)
	if hovered {
		target = 1.1
	}
	s := ui.UseTween(target, 160, ui.EaseOut)

	return ui.Div(
		ui.Style(ui.Width(130), ui.Height(150), ui.Bg(p.Color), ui.Radius(16),
			ui.Scale(s), ui.ItemsCenter, ui.JustifyCenter),
		ui.OnHover(setHovered),
		ui.Text(p.Label, ui.FontSize(18), ui.TextColor(ui.White)),
	)
}

func App(_ struct{}) *ui.Node {
	return ui.Div(
		ui.Style(ui.Column, ui.Width(760), ui.Height(320), ui.Bg(ui.Hex("#0f172a")),
			ui.Padding(36), ui.Gap(24), ui.ItemsCenter, ui.JustifyCenter),
		ui.Text("鼠标悬停卡片（hover + 弹性缩放）", ui.FontSize(22), ui.TextColor(ui.White)),
		ui.Div(
			ui.Style(ui.Row, ui.Gap(24), ui.ItemsCenter),
			ui.Use(Card, CardProps{"Scale", ui.Hex("#ef4444")}),
			ui.Use(Card, CardProps{"Hover", ui.Hex("#22c55e")}),
			ui.Use(Card, CardProps{"Spring", ui.Hex("#3b82f6")}),
			ui.Use(Card, CardProps{"Tween", ui.Hex("#a855f7")}),
		),
	)
}

func main() {
	ui.Run(ui.Use(App, struct{}{}))
}
