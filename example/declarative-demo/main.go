package main

import (
	"fmt"

	. "github.com/sjm1327605995/tenon/pkg/v2/declarative"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/yoga"
)

// appState 持有应用状态（闭包外持久化）
type appState struct {
	count int
	input string
}

func main() {
	SetTheme(ui.DefaultLightTheme())

	state := &appState{}

	Run(func() ui.Widget {
		return VStack(
			// 标题
			Text("Tenon Declarative API").
				FontSize(28).
				Color(Black),

			Text("React + SwiftUI style UI in Go").
				FontSize(16).
				Color(Gray),

			// 计数器（用 StatefulBuilder 包裹，setState 触发局部重建）
			CounterCard(state),

			// 输入框
			Input("Type something...").
				OnChange(func(s string) {
					state.input = s
					fmt.Println("Input:", s)
				}),

			// 水平布局
			HStack(
				Container(Text("A").Color(White)).Background(Red).CornerRadius(8).Padding(16),
				Container(Text("B").Color(White)).Background(Green).CornerRadius(8).Padding(16),
				Container(Text("C").Color(White)).Background(Blue).CornerRadius(8).Padding(16),
			).Gap(8).Justify(yoga.JustifySpaceBetween),

			Spacer(),
			Text("Built with Tenon v2").FontSize(12).Color(Gray),
		).
			Gap(16).
			Padding(24).
			Align(yoga.AlignStretch)
	}, 800, 600)
}

// CounterCard 用 StatefulBuilder 实现有状态的计数器卡片
func CounterCard(state *appState) ui.Widget {
	return ui.NewStatefulBuilder(func(ctx ui.BuildContext, setState func(func())) ui.Widget {
		return Card(
			VStack(
				Text(fmt.Sprintf("Count: %d", state.count)).
					FontSize(24).
					Color(Black),

				HStack(
					Button("+1").
						Style(ButtonPrimary).
						OnClick(func() {
							state.count++
							setState(nil) // 触发重建
						}),

					Button("Reset").
						Style(ButtonOutline).
						OnClick(func() {
							state.count = 0
							setState(nil)
						}),
				).Gap(8),
			).Gap(12),
		)
	})
}
