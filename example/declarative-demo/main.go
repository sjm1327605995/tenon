package main

import (
	"fmt"

	. "github.com/sjm1327605995/tenon/pkg/v2/declarative"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/yoga"
)

type appState struct {
	count int
}

func main() {
	SetTheme(ui.DefaultLightTheme())
	state := &appState{}

	Run(func() ui.Widget {
		return VStack(
			Text("Tenon Declarative API").FontSize(28).Color(Black),
			Text("React + SwiftUI style UI in Go").FontSize(16).Color(Gray),

			// 计数器
			ui.NewStatefulBuilder(func(ctx ui.BuildContext, setState func(func())) ui.Widget {
				return Card(
					VStack(
						Text(fmt.Sprintf("Count: %d", state.count)).FontSize(24).Color(Black),
						HStack(
							Button("+1").Style(ButtonPrimary).OnClick(func() {
								state.count++
								setState(nil)
							}),
							Button("Reset").Style(ButtonOutline).OnClick(func() {
								state.count = 0
								setState(nil)
							}),
						).Gap(8),
					).Gap(12),
				)
			}),

			Input("Type something...").OnChange(func(s string) {
				fmt.Println("Input:", s)
			}),

			HStack(
				Container(Text("A").Color(White)).Background(Red).CornerRadius(8).Padding(16),
				Container(Text("B").Color(White)).Background(Green).CornerRadius(8).Padding(16),
				Container(Text("C").Color(White)).Background(Blue).CornerRadius(8).Padding(16),
			).Gap(8).Justify(yoga.JustifySpaceBetween),

			Spacer(),
			Text("Built with Tenon v2").FontSize(12).Color(Gray),
		).Gap(16).Padding(24).Align(yoga.AlignStretch)
	}, 800, 600)
}
