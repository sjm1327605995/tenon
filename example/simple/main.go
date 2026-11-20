// Command simple demonstrates the compact, pure-Go API.
// Run with: go run ./example/simple
package main

import (
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func App() *ui.Node {
	name, setName := ui.UseState("")
	saved, setSaved := ui.UseState("尚未保存")
	save := func() {
		if name == "" {
			setSaved("请先输入姓名")
			return
		}
		setSaved("已保存：" + name)
	}
	return ui.ThemeProvider(ui.LightTheme,
		ui.Center(
			ui.Style(ui.Fill, ui.Bg(ui.Hex("#f4f4f5"))),
			ui.VStack(14,
				ui.Text("用 Go，舒服地写 UI").FontSize(22).Style(ui.Semibold),
				shadcn.Label("姓名"),
				shadcn.TextField(name, setName).Placeholder("输入姓名，回车即可保存").OnSubmit(func(string) { save() }),
				shadcn.TextButton("保存", save),
				ui.Text(saved).FontSize(14).TextColor(ui.DarkGray),
				shadcn.Separator(shadcn.SeparatorProps{}),
				ui.For([]string{"原生 Go 自动补全", "无需手写空 props", "样式与布局树分开阅读"}, func(tip string) *ui.Node {
					return ui.Text("• " + tip).FontSize(13).Key(tip)
				}),
			).Width(380).Padding(24).Radius(12).Bg(ui.White),
		))
}

func main() {
	ui.WindowTitle("Tenon — Simple API")
	ui.WindowSize(520, 460)
	ui.Run(ui.View(App))
}
