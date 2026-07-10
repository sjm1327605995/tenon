package main

import (
	"fmt"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// 可拖拽卡片：OnDrag 逐帧位移 -> UseState -> TranslateXY，命中测试跟随 transform。
func App(_ struct{}) *ui.Node {
	x, setX := ui.UseState(float32(0))
	y, setY := ui.UseState(float32(0))
	hovered, setHovered := ui.UseState(false)

	bg := ui.Hex("#3b82f6")
	if hovered {
		bg = ui.Hex("#2563eb")
	}

	return ui.Div(
		ui.Style(ui.Width(640), ui.Height(420), ui.Bg(ui.Hex("#0f172a")),
			ui.ItemsCenter, ui.JustifyCenter),

		ui.Div(
			ui.Style(ui.Width(150), ui.Height(150), ui.Bg(bg), ui.Radius(18),
				ui.TranslateXY(x, y), ui.ItemsCenter, ui.JustifyCenter, ui.Column, ui.Gap(6)),
			ui.OnHover(setHovered),
			ui.OnDrag(func(dx, dy float32) { setX(x + dx); setY(y + dy) }),
			ui.Text("拖动我", ui.FontSize(20), ui.TextColor(ui.White)),
			ui.Text(fmt.Sprintf("(%.0f, %.0f)", x, y), ui.FontSize(13),
				ui.TextColor(ui.Hex("#dbeafe"))),
		),
	)
}

func main() {
	ui.Run(ui.Use(App, struct{}{}))
}
