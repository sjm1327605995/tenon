package main

import (
	"log"

	"github.com/sjm1327605995/tenon/app"
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/primitives"
	"github.com/sjm1327605995/tenon/widget"
)

func main() {
	root := primitives.NewColumn(
		primitives.NewBox(
			primitives.NewText("Hello, Tenon + Yoga + Ebiten!"),
		),
		primitives.NewRow(
			primitives.NewBox(
				primitives.NewText("Left"),
			),
			primitives.NewBox(
				primitives.NewText("Right"),
			),
		),
	)

	// Style the root box
	rootBox := root.Children()[0].(*primitives.Box)
	rootBox.Background = widget.ColorLightGray
	rootBox.Padding = geometry.Insets{Left: 20, Top: 20, Right: 20, Bottom: 20}

	a := app.New(root)
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
