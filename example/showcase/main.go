package main

import (
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func App(_ struct{}) *ui.Node {
	otp, setOTP := ui.UseState("428")
	th := ui.DarkTheme
	return ui.ThemeProvider(th,
		ui.Div(
			ui.Style(ui.Column, ui.Fill, ui.ItemsCenter, ui.JustifyCenter, ui.Gap(24),
				ui.Bg(th.Background), ui.TextColor(th.Foreground)),

			ui.Div(ui.Style(ui.Column, ui.Width(380), ui.Gap(20),
				ui.Bg(th.Card), ui.Border(1, th.Border), ui.Radius(th.Radius+4), ui.Padding(24)),

				shadcn.H3("Item 列表项"),
				ui.Div(ui.Style(ui.Column, ui.Gap(2)),
					shadcn.Item(shadcn.ItemProps{
						Media: ui.Icon(ui.IconSearch, 18),
						Title: "搜索", Description: "全局搜索命令与文件",
						Trailing: shadcn.KbdGroup("Ctrl", "K"), OnClick: func() {}}),
					shadcn.Item(shadcn.ItemProps{
						Media: ui.Icon(ui.IconTrash, 18),
						Title: "回收站", Description: "30 天后自动清空",
						Trailing: shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeSecondary}, ui.Text("12")),
						OnClick:  func() {}}),
				),

				shadcn.H3("Input OTP 验证码"),
				shadcn.InputOTP(shadcn.InputOTPProps{Length: 6, Value: otp, OnChange: setOTP}),
			),
		),
	)
}

func main() { ui.Run(ui.Use(App, struct{}{})) }
