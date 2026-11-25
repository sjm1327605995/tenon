package main

import (
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

func main() {
	text := value.New("Edit אבג and أبجد")
	ui.WindowTitle("Tenon — BiDi text")
	ui.WindowSize(760, 420)
	ui.Run(ui.VStack(18,
		ui.Text("BiDi text / 双向文字").FontSize(24),
		ui.Text("English אבג Arabic أبجد and English"),
		ui.RichText(
			ui.Text("RichText: "),
			ui.Text("אבג", ui.TextColor(ui.Hex("#dc2626"))),
			ui.Text(" + "),
			ui.Text("أبجد", ui.Bold, ui.TextColor(ui.Hex("#2563eb"))),
			ui.Text(" + English"),
		),
		ui.Text("普通文本和 RichText 使用视觉 Run；输入框保留逻辑编辑顺序。"),
		shadcn.TextField(text).Placeholder("Mixed-direction input"),
	).Width(680).Padding(28))
}
