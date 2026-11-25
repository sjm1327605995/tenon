package main

import (
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

func main() {
	ui.WindowTitle("Tenon multiwindow")
	ui.WindowSize(360, 240)
	first := value.New("Window A")
	second := value.New("Window B")
	root := func(title string, binding value.Bind[string]) *ui.Node {
		return ui.VStack(14,
			ui.Text(title).FontSize(22),
			shadcn.TextField(binding).Placeholder("Independent input"),
			ui.Text("每个窗口拥有独立输入、IME、剪贴板和生命周期。"),
		).Padding(24)
	}
	if err := ui.RunWindows(root("Window A", first), root("Window B", second)); err != nil {
		panic(err)
	}
}
