package main

import (
	"fmt"
	"image/color"

	"github.com/sjm1327605995/tenon/core/gio"
	"github.com/sjm1327605995/tenon/core/ui"
	"github.com/sjm1327605995/tenon/yoga"

	"gioui.org/unit"
)

// CounterProps 计数器组件属性
type CounterProps struct {
	InitialCount int
}

// Counter 计数器组件，使用React-like hooks
func Counter(props CounterProps) ui.UI {
	// 使用类似React的useState hook
	count, setCount := ui.UseState(props.InitialCount)

	return ui.View(
		// 显示当前计数
		ui.Text().Content(fmt.Sprintf("Count: %d", count)),

		// 增加按钮
		ui.View(
			ui.Text().Content("+"),
		).Background(color.NRGBA{G: 255, A: 255}).Height(ui.Px(50)).Width(ui.Px(100)).OnClick(func() {
			setCount(count + 1)
			fmt.Println(count)
		}),
		// 减少按钮
		ui.View(
			ui.Text().Content("-"),
		).Background(color.NRGBA{R: 255, A: 255}).Height(ui.Px(50)).Width(ui.Px(100)).OnClick(func() {
			setCount(count - 1)
		}),
	).
		FlexDirection(yoga.FlexDirectionColumn).
		JustifyContent(yoga.JustifyCenter).
		AlignItems(yoga.AlignCenter).
		Background(color.NRGBA{B: 255, A: 128})
}

// 计数器组件创建函数，返回ui.UI接口
func createCounterComponent() ui.UI {
	// 使用新的函数组件创建方式
	return ui.NewFunctionComponent(Counter, CounterProps{InitialCount: 0})
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
