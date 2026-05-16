package main

import (
	"log"

	"github.com/sjm1327605995/tenon/app"
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/primitives"
	"github.com/sjm1327605995/tenon/theme"
	"github.com/sjm1327605995/tenon/widget"
)

func main() {
	// Build a simple demo UI
	root := buildUI()

	// Create app and run
	a := app.New(root)
	_ = a

	// Set theme via context (widgets can access it during draw)
	root.Layout(nil, geometry.Constraints{
		MaxWidth:  800,
		MaxHeight: 600,
	})
	if tp, ok := root.(interface{ SetThemeProvider(widget.ThemeProvider) }); ok {
		tp.SetThemeProvider(theme.DefaultLight())
	}

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}

func buildUI() widget.Widget {
	title := primitives.NewText("Tenon UI Framework")
	title.FontSize = 32
	title.Bold = true
	title.Color = widget.Hex(0x6200EE)

	subtitle := primitives.NewText("Yoga Layout + Ebiten Render")
	subtitle.FontSize = 16
	subtitle.Color = widget.ColorDarkGray

	box1 := primitives.NewBox(primitives.NewText("Box 1"))
	box1.Background = widget.ColorLightGray
	box1.Padding = geometry.UniformInsets(20)

	box2 := primitives.NewBox(primitives.NewText("Box 2"))
	box2.Background = widget.Hex(0xE8DEF8)
	box2.Padding = geometry.UniformInsets(20)

	row := primitives.NewRow(box1, box2)
	row.Gap = 16

	col := primitives.NewColumn(
		primitives.NewBox(title),
		primitives.NewBox(subtitle),
		row,
	)
	col.Gap = 16
	col.Padding = geometry.UniformInsets(24)

	return col
}
