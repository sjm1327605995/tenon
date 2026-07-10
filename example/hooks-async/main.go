package main

import (
	"fmt"
	"time"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// 演示后台 goroutine（模拟网络请求）完成后用 ui.Post 安全回到 UI 线程更新状态。
func App(_ struct{}) *ui.Node {
	loading, setLoading := ui.UseState(false)
	result, setResult := ui.UseState("")
	count, setCount := ui.UseState(0)

	load := func() {
		if loading {
			return
		}
		setLoading(true)
		setResult("")
		go func() {
			time.Sleep(1200 * time.Millisecond) // 模拟耗时请求（在后台 goroutine）
			data := fmt.Sprintf("请求 #%d 完成 @ %s", count+1, time.Now().Format("15:04:05"))
			ui.Post(func() { // 回到渲染线程，setState 安全
				setLoading(false)
				setResult(data)
				setCount(count + 1)
			})
		}()
	}

	// 挂载时自动触发一次，便于直接观察异步流程
	ui.UseEffect(func() ui.Cleanup { load(); return nil })

	return ui.Div(
		ui.Style(ui.Fill, ui.ItemsCenter, ui.JustifyCenter, ui.Bg(ui.Hex("#0f172a")),
			ui.Column, ui.Gap(22), ui.TextColor(ui.White)),

		ui.Text("后台任务 + ui.Post 安全更新", ui.FontSize(22)),

		ui.Button(
			ui.Style(ui.PaddingXY(24, 12), ui.Radius(8), ui.ItemsCenter, ui.JustifyCenter,
				ui.Bg(ui.Hex("#3b82f6")), ui.StyleIf(loading, ui.Opacity(0.5))),
			ui.OnClick(load),
			ui.Text("发起请求", ui.FontSize(15), ui.TextColor(ui.White)),
		),

		ui.If(loading, ui.Div(ui.Style(ui.Row, ui.Gap(10), ui.ItemsCenter),
			ui.Spinner(24, ui.Hex("#38bdf8")),
			ui.Text("加载中…", ui.FontSize(14), ui.TextColor(ui.Hex("#94a3b8"))),
		)),

		ui.If(result != "", ui.Text(result, ui.FontSize(15), ui.TextColor(ui.Hex("#4ade80")))),
	)
}

func main() {
	ui.Run(ui.Use(App, struct{}{}))
}
