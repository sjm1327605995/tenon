// Command struct demonstrates UI declarations written as Go struct literals.
// Run with: go run ./example/struct
package main

import (
	d "github.com/sjm1327605995/tenon/pkg/decl"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func Profile() d.View {
	name, setName := ui.UseState("")
	saved, setSaved := ui.UseState("尚未保存")
	disabled, setDisabled := ui.UseState(false)
	save := func() { setSaved("已保存：" + name) }

	return d.Column{
		Width: 380, Gap: 14, Padding: 24, Background: ui.White, Radius: 12,
		Children: []d.View{
			d.Text{Value: "用结构体声明界面", FontSize: 22, Bold: true},
			d.Text{Value: "姓名", FontSize: 14},
			d.Input{
				Value: name, OnChange: setName, Placeholder: "输入姓名",
				Disabled: disabled, OnSubmit: func(string) { save() },
			},
			d.Row{Gap: 8, Children: []d.View{
				d.Switch{Checked: disabled, OnChange: setDisabled},
				d.Text{Value: "禁用输入与保存", FontSize: 14},
			}},
			d.Button{Text: "保存", OnClick: save, Disabled: disabled},
			d.Text{Value: saved, FontSize: 14, Color: ui.DarkGray},
			ui.Text("也可以混用已有 ui.Node").FontSize(12),
		},
	}
}

func main() {
	ui.WindowTitle("Tenon — Struct UI")
	ui.WindowSize(520, 430)
	d.Run(func() d.View {
		return d.ThemeProvider{
			Theme: ui.LightTheme,
			Children: []d.View{
				d.Box{
					WidthPct: 100, HeightPct: 100, Background: ui.Hex("#f4f4f5"),
					Style:    []ui.StyleOpt{ui.ItemsCenter, ui.JustifyCenter},
					Children: []d.View{d.Func(Profile)},
				},
			},
		}
	})
}
