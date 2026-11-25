package main

import (
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

type Profile struct {
	ui.Component
	Name value.Bind[string]
}

func (p *Profile) Render() ui.NodeSource {
	theme := p.Context().Theme()
	return ui.VStack(12,
		ui.Text("个人资料").FontSize(22),
		shadcn.TextField(p.Name).Placeholder("姓名"),
		ui.Text(p.Name),
		shadcn.TextButton("清空", func() { p.Name.Set("") }),
	).Width(360).Padding(24).Bg(theme.Card).TextColor(theme.CardForeground)
}

func main() {
	ui.WindowTitle("Tenon — 个人资料")
	ui.WindowSize(480, 320)
	ui.Run(ui.Render(&Profile{Name: value.New("Tenon")}))
}
