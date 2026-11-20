// Command component demonstrates embedding Component and obtaining Context from
// its inherited accessor. Run with: go run ./example/component
package main

import (
	"fmt"

	d "github.com/sjm1327605995/tenon/pkg/decl"
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

type Profile struct {
	d.Component
	Title string
	Name  value.Bind[string]
	Saved value.Bind[string]
}

func (p *Profile) save() {
	if p.Name.Get() == "" {
		p.Saved.Set("请先输入姓名")
	} else {
		p.Saved.Set("已保存：" + p.Name.Get())
	}
}

func (p *Profile) Render() d.View {
	ctx := p.Context()
	th, viewport := ctx.Theme(), ctx.Viewport()
	return ui.Center(ui.Style(ui.Fill, ui.Bg(th.Muted)),
		ui.VStack(14,
			ui.Text(p.Title).FontSize(22).Style(ui.Bold),
			shadcn.TextField(p.Name).Placeholder("姓名").OnSubmit(func(string) { p.save() }),
			shadcn.TextButton("保存", p.save),
			ui.Text(p.Saved).FontSize(14).TextColor(th.MutedForeground),
			ui.Text(value.Map(p.Name, func(name string) string { return "你好，" + name })).FontSize(14),
			ui.Text(fmt.Sprintf("当前窗口 %.0f × %.0f", viewport.W, viewport.H)).FontSize(12).TextColor(th.MutedForeground),
		).Width(380).Padding(24).Radius(12).Bg(th.Card))
}

func main() {
	ui.WindowTitle("Tenon — Component Context")
	ui.WindowSize(520, 400)
	d.Run(&Profile{Title: "Component + Value / Bind", Name: value.New(""), Saved: value.New("尚未保存")})
}

var _ d.IComponent = (*Profile)(nil)
