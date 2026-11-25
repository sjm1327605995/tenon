package ui

import (
	"image"
	"image/color"
	"testing"
)

var benchmarkFilteredImage image.Image

func imageFilterBenchmarkSource() *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, 512, 512))
	for y := 0; y < 512; y++ {
		for x := 0; x < 512; x++ {
			img.SetNRGBA(x, y, color.NRGBA{
				R: uint8(x), G: uint8(y), B: uint8(x + y), A: 255,
			})
		}
	}
	return img
}

func BenchmarkImageBlur(b *testing.B) {
	source := imageFilterBenchmarkSource()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkFilteredImage = blurredImage(source, 8)
	}
}

func BenchmarkImageGrayscale(b *testing.B) {
	source := imageFilterBenchmarkSource()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkFilteredImage = grayscaleImage(source, 0.65)
	}
}

func BenchmarkImageGrayscaleYCbCr(b *testing.B) {
	source := image.NewYCbCr(image.Rect(0, 0, 512, 512), image.YCbCrSubsampleRatio444)
	for i := range source.Y {
		source.Y[i] = uint8(i)
	}
	for i := range source.Cb {
		source.Cb[i], source.Cr[i] = uint8(i*3), uint8(i*5)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkFilteredImage = grayscaleImage(source, 0.65)
	}
}

func BenchmarkImageGrayscaleGray(b *testing.B) {
	source := image.NewGray(image.Rect(0, 0, 512, 512))
	for i := range source.Pix {
		source.Pix[i] = uint8(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkFilteredImage = grayscaleImage(source, 0.65)
	}
}
