package ui

import (
	"image"
	"image/color"
)

func grayscaleImage(source image.Image, amount float32) image.Image {
	if amount <= 0 || source == nil {
		return source
	}
	amount = min(amount, 1)
	bounds := source.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width == 0 || height == 0 {
		return source
	}
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := imageNRGBAAt(source, bounds.Min.X+x, bounds.Min.Y+y)
			luminance := 0.2126*float32(pixel.R) + 0.7152*float32(pixel.G) + 0.0722*float32(pixel.B)
			img.SetNRGBA(x, y, color.NRGBA{
				R: grayscaleChannel(pixel.R, luminance, amount),
				G: grayscaleChannel(pixel.G, luminance, amount),
				B: grayscaleChannel(pixel.B, luminance, amount),
				A: pixel.A,
			})
		}
	}
	return img
}

func grayscaleChannel(channel uint8, luminance, amount float32) uint8 {
	value := float32(channel)*(1-amount) + luminance*amount
	return uint8(value + 0.5)
}
