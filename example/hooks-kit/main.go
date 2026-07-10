package main

import (
	"fmt"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func Row(label string, control *ui.Node) *ui.Node {
	return ui.Div(
		ui.Style(ui.Row, ui.ItemsCenter, ui.JustifyBetween, ui.Gap(16)),
		ui.Text(label, ui.FontSize(15)),
		control,
	)
}

func App(_ struct{}) *ui.Node {
	agree, setAgree := ui.UseState(false)
	notify, setNotify := ui.UseState(true)
	choice, setChoice := ui.UseState(0)
	volume, setVolume := ui.UseState(float32(40))
	tab, setTab := ui.UseState(0)

	return ui.Div(
		ui.Style(ui.Column, ui.Width(560), ui.Height(680), ui.Bg(ui.Hex("#f1f5f9")),
			ui.Padding(28), ui.Gap(16), ui.TextColor(ui.Hex("#0f172a"))),

		ui.Text("基础组件套件", ui.FontSize(24)),

		ui.Tabs(ui.TabsProps{Tabs: []string{"控件", "展示"}, Active: tab, OnChange: setTab}),

		ui.If(tab == 0, ui.Card(
			Row("Checkbox", ui.Checkbox(agree, setAgree)),
			ui.Divider(),
			Row("Switch", ui.Switch(notify, setNotify)),
			ui.Divider(),
			Row("Radio", ui.Div(
				ui.Style(ui.Row, ui.Gap(16)),
				ui.Radio(choice == 0, func() { setChoice(0) }),
				ui.Radio(choice == 1, func() { setChoice(1) }),
				ui.Radio(choice == 2, func() { setChoice(2) }),
			)),
			ui.Divider(),
			Row(fmt.Sprintf("Slider (%.0f)", volume), ui.Slider(volume, 0, 100, setVolume)),
			ui.Divider(),
			Row("Progress", ui.ProgressBar(volume/100)),
		)),

		ui.If(tab == 1, ui.Card(
			Row("Badge", ui.Div(ui.Style(ui.Row, ui.Gap(8)),
				ui.Badge("New", ui.Hex("#3b82f6")),
				ui.Badge("Beta", ui.Hex("#a855f7")),
				ui.Badge("Done", ui.Hex("#22c55e")),
			)),
			ui.Divider(),
			Row("Avatar", ui.Div(ui.Style(ui.Row, ui.Gap(8)),
				ui.Avatar("孙", 40),
				ui.Avatar("AI", 40),
			)),
			ui.Divider(),
			Row("Spinner", ui.Spinner(28, ui.Hex("#3b82f6"))),
		)),
	)
}

func main() {
	ui.Run(ui.Use(App, struct{}{}))
}
