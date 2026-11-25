package ui

import (
	"image"
	"image/color"
	"strconv"
)

func (rn *renderNode) imageKey() string {
	key := rn.imgSrc
	if rn.imgBlur > 0 {
		key += "\x00blur=" + strconv.FormatFloat(float64(rn.imgBlur), 'f', -1, 32)
	}
	if rn.imgGray > 0 {
		key += "\x00gray=" + strconv.FormatFloat(float64(rn.imgGray), 'f', -1, 32)
	}
	return key
}

func processImage(source image.Image, blur, grayscale float32) image.Image {
	return grayscaleImage(blurredImage(source, blur), grayscale)
}

func blurredImage(source image.Image, radius float32) image.Image {
	if radius <= 0 || source == nil {
		return source
	}
	bounds := source.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width == 0 || height == 0 {
		return source
	}
	pixels := make([]color.NRGBA, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixels[y*width+x] = imageNRGBAAt(source, bounds.Min.X+x, bounds.Min.Y+y)
		}
	}
	horizontal := make([]color.NRGBA, len(pixels))
	output := pixels
	radiusPixels := int(radius)
	if radiusPixels < 1 {
		radiusPixels = 1
	}
	for y := 0; y < height; y++ {
		blurImageAxis(pixels, horizontal, y*width, width, 1, radiusPixels)
	}
	for x := 0; x < width; x++ {
		blurImageAxis(horizontal, output, x, height, width, radiusPixels)
	}
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for i, pixel := range output {
		img.SetNRGBA(i%width, i/width, pixel)
	}
	return img
}

func blurImageAxis(source, target []color.NRGBA, start, count, stride, radius int) {
	if count == 0 {
		return
	}
	left := 0
	right := min(count-1, radius)
	var red, green, blue, alpha int
	for i := left; i <= right; i++ {
		pixel := source[start+i*stride]
		red += int(pixel.R)
		green += int(pixel.G)
		blue += int(pixel.B)
		alpha += int(pixel.A)
	}
	for i := 0; i < count; i++ {
		target[start+i*stride] = color.NRGBA{
			R: uint8(red / (right - left + 1)),
			G: uint8(green / (right - left + 1)),
			B: uint8(blue / (right - left + 1)),
			A: uint8(alpha / (right - left + 1)),
		}
		nextLeft := max(0, i+1-radius)
		if nextLeft > left {
			pixel := source[start+left*stride]
			red -= int(pixel.R)
			green -= int(pixel.G)
			blue -= int(pixel.B)
			alpha -= int(pixel.A)
			left = nextLeft
		}
		nextRight := min(count-1, i+1+radius)
		if nextRight > right {
			right = nextRight
			pixel := source[start+right*stride]
			red += int(pixel.R)
			green += int(pixel.G)
			blue += int(pixel.B)
			alpha += int(pixel.A)
		}
	}
}
