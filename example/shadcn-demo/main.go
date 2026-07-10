package main

import (
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func App(_ struct{}) *ui.Node {
	dark, setDark := ui.UseState(false)
	theme := ui.LightTheme
	if dark {
		theme = ui.DarkTheme
	}

	row := func(children ...*ui.Node) *ui.Node {
		return ui.Div(append([]*ui.Node{ui.Style(ui.Row, ui.Gap(12), ui.ItemsCenter)}, children...)...)
	}

	// ThemeProvider 下：shadcn 组件按主题取色；也能和基础 ui 组件混用
	return ui.ThemeProvider(theme,
		ui.Div(
			ui.Style(ui.Column, ui.Width(680), ui.Height(420), ui.Bg(theme.Background),
				ui.Padding(36), ui.Gap(20), ui.TextColor(theme.Foreground)),

			row(
				ui.Text("shadcn 风格按钮", ui.FontSize(24)),
				// 基础 ui 组件（Switch）与 shadcn 混用
				ui.Switch(dark, setDark),
				ui.Text("暗色", ui.FontSize(13), ui.TextColor(theme.MutedForeground)),
			),

			row(
				shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Default}, ui.Text("Default")),
				shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Secondary}, ui.Text("Secondary")),
				shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Destructive}, ui.Text("Destructive")),
			),
			row(
				shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline}, ui.Text("Outline")),
				shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Ghost}, ui.Text("Ghost")),
				shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Link}, ui.Text("Link")),
			),
			row(
				shadcn.Button(shadcn.ButtonProps{Size: shadcn.SizeSm}, ui.Text("Small")),
				shadcn.Button(shadcn.ButtonProps{Size: shadcn.SizeLg}, ui.Text("Large")),
				shadcn.Button(shadcn.ButtonProps{Disabled: true}, ui.Text("Disabled")),
			),
		),
	)
}

func main() {
	ui.Run(ui.Use(App, struct{}{}))
}
