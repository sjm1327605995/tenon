package main

import (
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// 复刻 shadcn/ui Accordion 文档效果：高度滑动 + chevron 旋转 + 明暗主题切换。
func App(_ struct{}) *ui.Node {
	dark, setDark := ui.UseState(true)
	theme := ui.LightTheme
	if dark {
		theme = ui.DarkTheme
	}

	items := []shadcn.AccordionItemData{
		{Title: "Is it accessible?", Content: []*ui.Node{
			ui.Text("Yes. It adheres to the WAI-ARIA design pattern.")}},
		{Title: "Is it styled?", Content: []*ui.Node{
			ui.Text("Yes. It comes with default styles that match the other components’ aesthetic.")}},
		{Title: "Is it animated?", Content: []*ui.Node{
			ui.Text("Yes. It’s animated by default with a smooth height slide, but you can disable it if you prefer.")}},
	}

	return ui.ThemeProvider(theme,
		ui.Div(
			ui.Style(ui.Column, ui.Fill, ui.ItemsCenter, ui.JustifyCenter, ui.Gap(20),
				ui.Bg(theme.Background), ui.TextColor(theme.Foreground)),

			ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(10)),
				shadcn.Label("暗色主题"),
				shadcn.Switch(shadcn.SwitchProps{Checked: dark, OnChange: setDark}),
			),

			// 预览卡片（对应 shadcn 文档里的组件预览框）
			ui.Div(
				ui.Style(ui.Column, ui.Width(460), ui.Bg(theme.Card), ui.TextColor(theme.CardForeground),
					ui.Border(1, theme.Border), ui.Radius(theme.Radius+4), ui.PaddingXY(20, 8)),
				shadcn.Accordion(items),
			),
		),
	)
}

func main() { ui.Run(ui.Use(App, struct{}{})) }
