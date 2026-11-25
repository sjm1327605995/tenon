package ui

import (
	"image"
	"image/color"
	"math"
)

func newRadialGradientImage(width, height int, from, to Color) *image.NRGBA {
	if width <= 0 || height <= 0 {
		return image.NewNRGBA(image.Rect(0, 0, 0, 0))
	}
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	centerX, centerY := float64(width)/2, float64(height)/2
	maxDistance := math.Hypot(centerX, centerY)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			distance := math.Hypot(float64(x)-centerX, float64(y)-centerY)
			t := min(distance/maxDistance, 1)
			img.SetNRGBA(x, y, color.NRGBA{
				R: mixGradientChannel(from.R, to.R, t),
				G: mixGradientChannel(from.G, to.G, t),
				B: mixGradientChannel(from.B, to.B, t),
				A: mixGradientChannel(from.A, to.A, t),
			})
		}
	}
	return img
}

func mixGradientChannel(from, to uint8, t float64) uint8 {
	return uint8(math.Round(float64(from) + (float64(to)-float64(from))*t))
}

func (rn *renderNode) paintRadialGradient(painter painter, opacity float32) {
	bounds := rn.bounds
	width, height := int(bounds.W), int(bounds.H)
	if width <= 0 || height <= 0 {
		return
	}
	if rn.radialImg == nil || rn.radialW != width || rn.radialH != height || rn.radialFrom != rn.gradFrom || rn.radialTo != rn.gradTo {
		rn.radialImg = backendNewBitmap(newRadialGradientImage(width, height, rn.gradFrom, rn.gradTo))
		rn.radialW, rn.radialH = width, height
		rn.radialFrom, rn.radialTo = rn.gradFrom, rn.gradTo
	}
	if rn.radialImg == nil {
		return
	}
	painter.PushClip(bounds, rn.radius)
	painter.DrawImage(rn.radialImg, bounds, opacity)
	painter.PopClip()
}
