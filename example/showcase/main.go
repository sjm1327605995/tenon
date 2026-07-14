package main

import (
	"os"
	"time"

	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func App(scene string) *ui.Node {
	th := ui.DarkTheme
	menus := []shadcn.MenubarMenu{
		{Label: "File", Items: []shadcn.MenuItem{{Label: "New"}, {Label: "Open"}, {Label: "Save"}}},
		{Label: "Edit", Items: []shadcn.MenuItem{{Label: "Undo"}, {Label: "Redo"}}},
		{Label: "View", Items: []shadcn.MenuItem{{Label: "Zoom In"}, {Label: "Zoom Out"}}},
	}
	base := ui.Div(ui.Style(ui.Column, ui.Fill, ui.ItemsCenter, ui.JustifyCenter, ui.Gap(24),
		ui.Bg(th.Background), ui.TextColor(th.Foreground)),
		shadcn.H3("Menubar + DatePicker"),
		shadcn.Menubar(menus),
		shadcn.DatePicker(shadcn.DatePickerProps{Value: time.Date(2026, 7, 14, 0, 0, 0, 0, time.Local)}),
	)
	kids := []*ui.Node{ui.Style(ui.Fill), base}
	switch scene {
	case "alert":
		kids = append(kids, shadcn.AlertDialog(shadcn.AlertDialogProps{
			Open: true, Title: "确认删除？",
			Description: "此操作不可撤销，将永久删除该项目及其所有数据。",
			ActionLabel: "删除", Destructive: true}))
	case "drawer":
		kids = append(kids, shadcn.Drawer(shadcn.DrawerProps{Open: true, Height: 280},
			shadcn.H3("分享到"),
			shadcn.Muted("从底部滑入的抽屉内容。"),
			shadcn.Button(shadcn.ButtonProps{}, ui.Text("复制链接"))))
	}
	return ui.ThemeProvider(th, ui.Div(kids...))
}

func main() { ui.Run(ui.Use(App, os.Getenv("SCENE"))) }
