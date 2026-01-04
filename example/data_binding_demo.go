package main

import (
	"fmt"
	"image/color"

	"github.com/sjm1327605995/tenon/core/gio"
	"github.com/sjm1327605995/tenon/core/ui"
	"github.com/sjm1327605995/tenon/yoga"

	"gioui.org/unit"
)

// 计数器组件创建函数，返回ui.UI接口
func createCounterComponent() ui.UI {
	// 使用ui.UseState创建状态
	count, _ := ui.UseState(0)

	// 创建一个无状态组件，包装计数器UI
	return ui.NewStatelessComponent(func() ui.UI {
		return ui.View(
			// 显示当前计数
			ui.Text().Content(fmt.Sprintf("Count: %d", count.Value())),

			// 增加按钮
			ui.View(
				ui.Text().Content("+"),
			).Background(color.NRGBA{G: 255, A: 255}).Height(ui.Px(50)).Width(ui.Px(100)),
			ui.View(
				ui.Text().Content("-"),
			).Background(color.NRGBA{R: 255, A: 255}).Height(ui.Px(50)).Width(ui.Px(100)),
		).
			FlexDirection(yoga.FlexDirectionColumn).
			JustifyContent(yoga.JustifyCenter).
			AlignItems(yoga.AlignCenter).
			Background(color.NRGBA{B: 255, A: 128})
	})
}

func main() {
	// 创建用户组件
	component := createCounterComponent()

	// 使用新的core/gio库运行应用，用户不再需要关心Gio层面的内容
	gio.RunApp(
		gio.AppConfig{
			Title:  "Data Binding Demo",
			Width:  unit.Dp(800),
			Height: unit.Dp(600),
		},
		component,
	)
}
