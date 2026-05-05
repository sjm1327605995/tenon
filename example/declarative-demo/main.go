package main

import (
	"fmt"

	. "github.com/sjm1327605995/tenon/pkg/v2/declarative"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/yoga"
)

// 声明式 API 示例 — 类 React + SwiftUI 风格
func main() {
	count := 0

	SetTheme(ui.DefaultLightTheme())

	Run(func() ui.Widget {
		return VStack(
			// 标题
			Text("Tenon Declarative API").
				FontSize(28).
				Color(Black),

			Text("React + SwiftUI style UI in Go").
				FontSize(16).
				Color(Gray),

			// 计数器卡片
			Card(
				VStack(
					Text(fmt.Sprintf("Count: %d", count)).
						FontSize(24),

					HStack(
						Button("+1").
							Style(ButtonPrimary).
							OnClick(func() { count++ }),

						Button("Reset").
							Style(ButtonOutline).
							OnClick(func() { count = 0 }),
					).Gap(8),
				).Gap(12),
			),

			// 输入框
			Input("Type something...").
				OnChange(func(s string) {
					fmt.Println("Input:", s)
				}),

			// 水平布局
			HStack(
				Container(Text("A")).Background(Red).CornerRadius(8).Padding(16),
				Container(Text("B")).Background(Green).CornerRadius(8).Padding(16),
				Container(Text("C")).Background(Blue).CornerRadius(8).Padding(16),
			).Gap(8).Justify(yoga.JustifySpaceBetween),

			Spacer(),
			Text("Built with Tenon v2").FontSize(12).Color(Gray),
		).
			Gap(16).
			Padding(24).
			Align(yoga.AlignStretch)
	}, 800, 600)
}
