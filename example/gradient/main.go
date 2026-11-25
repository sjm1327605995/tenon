package main

import (
	"image"
	"image/color"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func main() {
	sample := sampleImage()
	ui.WindowTitle("Tenon — gradients and image filters")
	ui.WindowSize(520, 610)
	ui.Run(ui.VStack(16,
		ui.Text("Radial gradient").FontSize(22),
		ui.Div(ui.Style(ui.Width(240), ui.Height(140), ui.Radius(20),
			ui.RadialGradient(ui.Hex("#fbbf24"), ui.Hex("#7c3aed")))),
		ui.Text("Linear gradient"),
		ui.Div(ui.Style(ui.Width(240), ui.Height(64), ui.Radius(12),
			ui.LinearGradient(ui.Hex("#22c55e"), ui.Hex("#2563eb"), 25))),
		ui.Text("Image filters").FontSize(22),
		ui.HStack(16,
			ui.VStack(6,
				ui.Text("Original"),
				ui.Img(ui.SrcImage("gradient:sample", sample), ui.Style(ui.Width(140), ui.Height(120))),
			),
			ui.VStack(6,
				ui.Text("Blur radius 5"),
				ui.Img(ui.SrcImage("gradient:sample", sample),
					ui.Style(ui.Width(140), ui.Height(120), ui.ImageBlur(5))),
			),
			ui.VStack(6,
				ui.Text("Grayscale"),
				ui.Img(ui.SrcImage("gradient:sample", sample),
					ui.Style(ui.Width(140), ui.Height(120), ui.ImageGrayscale(1))),
			),
		),
	).Padding(28))
}

func sampleImage() image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, 180, 120))
	palette := []color.NRGBA{
		{R: 15, G: 23, B: 42, A: 255},
		{R: 248, G: 250, B: 252, A: 255},
		{R: 249, G: 115, B: 22, A: 255},
		{R: 34, G: 197, B: 94, A: 255},
	}
	for y := 0; y < 120; y++ {
		for x := 0; x < 180; x++ {
			img.SetNRGBA(x, y, palette[(x/15+y/15)%len(palette)])
		}
	}
	return img
}
