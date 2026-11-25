package ui

import (
	"image"
	"image/color"
	"testing"
)

func TestImageNRGBAAtMatchesStandardConversion(t *testing.T) {
	bounds := image.Rect(2, 3, 3, 4)
	nrgba := image.NewNRGBA(bounds)
	nrgba.SetNRGBA(2, 3, color.NRGBA{R: 200, G: 100, B: 50, A: 128})
	nrgba64 := image.NewNRGBA64(bounds)
	nrgba64.SetNRGBA64(2, 3, color.NRGBA64{R: 0xab12, G: 0x5634, B: 0x1278, A: 0x8080})
	rgba := image.NewRGBA(bounds)
	rgba.SetRGBA(2, 3, color.RGBA{R: 100, G: 50, B: 25, A: 128})
	rgba64 := image.NewRGBA64(bounds)
	rgba64.SetRGBA64(2, 3, color.RGBA64{R: 0x8000, G: 0x4000, B: 0x2000, A: 0x8000})
	ycbcr := image.NewYCbCr(bounds, image.YCbCrSubsampleRatio444)
	ycbcr.Y[ycbcr.YOffset(2, 3)] = 120
	ycbcr.Cb[ycbcr.COffset(2, 3)] = 200
	ycbcr.Cr[ycbcr.COffset(2, 3)] = 80

	for _, test := range []struct {
		name  string
		img   image.Image
		color color.Color
	}{
		{name: "nrgba", img: nrgba, color: nrgba.NRGBAAt(2, 3)},
		{name: "nrgba64", img: nrgba64, color: nrgba64.NRGBA64At(2, 3)},
		{name: "rgba", img: rgba, color: rgba.RGBAAt(2, 3)},
		{name: "rgba64", img: rgba64, color: rgba64.RGBA64At(2, 3)},
		{name: "ycbcr", img: ycbcr, color: color.YCbCr{Y: 120, Cb: 200, Cr: 80}},
		{name: "uniform", img: image.NewUniform(color.RGBA{R: 100, G: 50, B: 25, A: 128}), color: color.RGBA{R: 100, G: 50, B: 25, A: 128}},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := imageNRGBAAt(test.img, 2, 3)
			want := color.NRGBAModel.Convert(test.color).(color.NRGBA)
			if got != want {
				t.Fatalf("converted pixel = %#v, want %#v", got, want)
			}
		})
	}
}
