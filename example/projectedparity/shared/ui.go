package shared

import (
	"image"
	"image/color"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func Root() *ui.Node {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.SetNRGBA(0, 0, color.NRGBA{R: 239, G: 68, B: 68, A: 255})
	img.SetNRGBA(1, 0, color.NRGBA{R: 34, G: 197, B: 94, A: 255})
	img.SetNRGBA(0, 1, color.NRGBA{R: 59, G: 130, B: 246, A: 255})
	img.SetNRGBA(1, 1, color.NRGBA{R: 250, G: 204, B: 21, A: 255})
	return ui.Div(
		ui.Style(ui.Fill, ui.Column, ui.ItemsCenter, ui.JustifyCenter, ui.Gap(20), ui.Bg(ui.Hex("#0f172a"))),
		ui.Div(
			ui.Style(ui.Width(240), ui.Height(180), ui.Perspective(700), ui.RotateX(-10), ui.RotateY(30)),
			ui.Img(ui.SrcImage("projection-colors", img), ui.Style(ui.Width(240), ui.Height(180))),
		),
		ui.Text("Shared homography", ui.FontSize(16), ui.TextColor(ui.Hex("#cbd5e1"))),
	)
}
