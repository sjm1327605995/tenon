package main

import (
	"fmt"
	"math/rand"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

var palette = []struct {
	Label string
	Color ui.Color
}{
	{"红", ui.Hex("#ef4444")},
	{"橙", ui.Hex("#f97316")},
	{"绿", ui.Hex("#22c55e")},
	{"蓝", ui.Hex("#3b82f6")},
	{"紫", ui.Hex("#a855f7")},
}

// 列表重排布局动画：打乱顺序时，带 Animated 的项从旧位置平滑滑到新位置（FLIP）。
func App(_ struct{}) *ui.Node {
	order, setOrder := ui.UseState([]int{0, 1, 2, 3, 4})

	shuffle := func() {
		n := append([]int{}, order...)
		rand.Shuffle(len(n), func(i, j int) { n[i], n[j] = n[j], n[i] })
		setOrder(n)
	}

	rows := []*ui.Node{ui.Style(ui.Column, ui.Gap(12))}
	for _, idx := range order {
		it := palette[idx]
		rows = append(rows, ui.Keyed(fmt.Sprintf("%d", idx),
			ui.Div(
				ui.Style(ui.Width(360), ui.Height(54), ui.Bg(it.Color), ui.Radius(10),
					ui.Animated, ui.ItemsCenter, ui.PaddingXY(18, 0)),
				ui.Text(it.Label+" 号色块", ui.FontSize(17), ui.TextColor(ui.White)),
			),
		))
	}

	return ui.Div(
		ui.Style(ui.Column, ui.Width(460), ui.Height(560), ui.Bg(ui.Hex("#0f172a")),
			ui.Padding(30), ui.Gap(18), ui.ItemsCenter, ui.TextColor(ui.White)),

		ui.Text("列表重排动画", ui.FontSize(22)),
		ui.Button(
			ui.Style(ui.PaddingXY(22, 11), ui.Bg(ui.Hex("#3b82f6")), ui.Radius(8),
				ui.ItemsCenter, ui.JustifyCenter),
			ui.OnClick(shuffle),
			ui.Text("打乱顺序", ui.FontSize(15)),
		),
		ui.Div(rows...),
	)
}

func main() {
	ui.Run(ui.Use(App, struct{}{}))
}
