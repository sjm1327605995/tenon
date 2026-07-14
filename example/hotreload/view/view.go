package view

import (
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// View 是被热更新解释的预览组件。改这里、保存，窗口即时更新（无需重启/重编译）。
// 限制：解释器不支持泛型，所以这里不能用 ui.UseState/ui.Use，是一个无状态视图；
// 但可以随意使用 pkg/ui 的非泛型 API 和 pkg/shadcn 的全部组件。
func View() *ui.Node {
	return ui.Div(
		ui.Style(ui.Column, ui.Fill, ui.ItemsCenter, ui.JustifyCenter,
			ui.Gap(16), ui.Padding(32), ui.Bg(ui.Hex("#0b1020")), ui.TextColor(ui.White)),

		shadcn.H1("Hot Reload 🔥"),
		shadcn.Muted("编辑 view/view.go 保存，这里即时刷新"),

		ui.Div(ui.Style(ui.Row, ui.Gap(10)),
			shadcn.Button(shadcn.ButtonProps{}, ui.Text("Primary")),
			shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline}, ui.Text("Outline")),
			shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Destructive}, ui.Text("Delete")),
		),

		ui.Div(ui.Style(ui.Row, ui.Gap(8), ui.ItemsCenter),
			shadcn.Badge(shadcn.BadgeProps{}, ui.Text("New")),
			shadcn.KbdGroup("Ctrl", "S"),
			shadcn.Spinner(shadcn.SpinnerProps{Size: 18}),
		),
	)
}
