package ui

import (
	"image"
	"image/color"
)

func imageNRGBAAt(source image.Image, x, y int) color.NRGBA {
	switch img := source.(type) {
	case *image.NRGBA:
		return img.NRGBAAt(x, y)
	case *image.NRGBA64:
		pixel := img.NRGBA64At(x, y)
		return color.NRGBA{R: uint8(pixel.R >> 8), G: uint8(pixel.G >> 8), B: uint8(pixel.B >> 8), A: uint8(pixel.A >> 8)}
	case *image.RGBA:
		return premultipliedNRGBA(img.RGBAAt(x, y).RGBA())
	case *image.RGBA64:
		return premultipliedNRGBA(img.RGBA64At(x, y).RGBA())
	case *image.YCbCr:
		return premultipliedNRGBA(img.YCbCrAt(x, y).RGBA())
	default:
		return premultipliedNRGBA(source.At(x, y).RGBA())
	}
}

func premultipliedNRGBA(red, green, blue, alpha uint32) color.NRGBA {
	if alpha == 0 {
		return color.NRGBA{}
	}
	if alpha != 0xffff {
		red = red * 0xffff / alpha
		green = green * 0xffff / alpha
		blue = blue * 0xffff / alpha
	}
	return color.NRGBA{R: uint8(red >> 8), G: uint8(green >> 8), B: uint8(blue >> 8), A: uint8(alpha >> 8)}
}
