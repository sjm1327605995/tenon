package main

import (
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// 演示文本自动折行 + 样式继承 + 富文本混排 + 输入法组字（IME）。
func App(_ struct{}) *ui.Node {
	name, setName := ui.UseState("")
	note, setNote := ui.UseState("")
	return ui.Div(
		ui.Style(ui.Column, ui.Width(600), ui.Height(620), ui.Bg(ui.Hex("#0f172a")),
			ui.Padding(36), ui.Gap(18),
			ui.TextColor(ui.Hex("#e2e8f0")), ui.FontSize(16)), // 下传给后代文本

		ui.Text("文本折行 · 样式继承 · 富文本 · 输入法", ui.FontSize(24), ui.Bold, ui.TextColor(ui.White)),

		// 字重展示（合成粗体）
		ui.Div(ui.Style(ui.Row, ui.Gap(16), ui.ItemsCenter),
			ui.Text("常规 Regular", ui.FontSize(18)),
			ui.Text("中等 Medium", ui.FontSize(18), ui.Medium),
			ui.Text("半粗 Semibold", ui.FontSize(18), ui.Semibold),
			ui.Text("粗体 Bold", ui.FontSize(18), ui.Bold),
			ui.Text("斜体 Italic", ui.FontSize(18), ui.Italic),
		),

		// 富文本：一个 Text 内多段不同颜色/字重/字号，统一折行、基线对齐
		ui.RichText(
			ui.Text("富文本 "),
			ui.Text("RichText", ui.Bold, ui.TextColor(ui.Hex("#38bdf8"))),
			ui.Text(" 可以在同一段落里混排 "),
			ui.Text("大字号", ui.FontSize(26), ui.TextColor(ui.Hex("#f472b6"))),
			ui.Text("、"),
			ui.Text("粗斜体", ui.Bold, ui.Italic, ui.TextColor(ui.Hex("#fbbf24"))),
			ui.Text(" 与普通文字，并随宽度自动折行。"),
		),

		ui.Div(
			ui.Style(ui.Column, ui.Width(400), ui.Bg(ui.Hex("#1e293b")), ui.Radius(12),
				ui.Padding(22), ui.Gap(12)),
			ui.Text("这是一段较长的中文段落，用来演示自动折行：当宽度不足以容纳整行时，文字会在合适的位置断开并换到下一行，颜色和字号都从父容器继承而来。"),
			ui.Text("English paragraphs wrap at word boundaries when the available width is exceeded, flowing naturally onto multiple lines.",
				ui.TextColor(ui.Hex("#94a3b8"))),
		),

		// 输入法组字：中文输入时预编辑串带下划线显示，回车/选词后提交
		ui.Text("试试中文输入（IME 组字预编辑）：", ui.FontSize(14), ui.TextColor(ui.Hex("#94a3b8"))),
		ui.Input(
			ui.Style(ui.Width(400), ui.Height(38), ui.Bg(ui.White), ui.Radius(8),
				ui.PaddingXY(10, 6), ui.TextColor(ui.Hex("#0f172a"))),
			ui.Placeholder("单行：输入姓名…"),
			ui.Value(name), ui.OnChange(setName),
		),
		ui.Input(
			ui.Style(ui.Width(400), ui.Height(120), ui.Bg(ui.White), ui.Radius(8),
				ui.Padding(10), ui.TextColor(ui.Hex("#0f172a"))),
			ui.Multiline(),
			ui.Placeholder("多行：写点备注，支持折行与选区…"),
			ui.Value(note), ui.OnChange(setNote),
		),
	)
}

func main() {
	ui.ShowStats = true // 启动即显示性能 HUD（F12 切换）：观察空闲时重绘降到 ~0/s
	ui.Run(ui.Use(App, struct{}{}))
}
