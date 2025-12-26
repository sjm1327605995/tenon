package main

import (
	"image/color"

	"github.com/sjm1327605995/tenon/core/gio"
	"github.com/sjm1327605995/tenon/core/ui"
	"github.com/sjm1327605995/tenon/yoga"

	"gioui.org/unit"
)

// 用户组件创建函数类型
type ComponentCreator func() ui.UI

// 用户定义的组件示例
func createUserComponent() ui.UI {
	return ui.View(
		ui.View().Background(color.NRGBA{G: 255, A: 255}).Height(ui.Px(100)).Width(ui.Px(100)),
		ui.View().Background(color.NRGBA{R: 255, A: 255}).Height(ui.Px(100)).Width(ui.Px(100)).
			Border(yoga.EdgeAll, 3).BorderRadius(20),
	).Background(color.NRGBA{B: 255, A: 128}).JustifyContent(yoga.JustifyCenter).AlignItems(yoga.AlignCenter)
}

func main() {
	// 创建用户组件
	component := createUserComponent()

	// 使用新的core/gio库运行应用，用户不再需要关心Gio层面的内容
	gio.RunApp(
		gio.AppConfig{
			Title:  "Component Demo",
			Width:  unit.Dp(800),
			Height: unit.Dp(600),
		},
		component,
	)
}
