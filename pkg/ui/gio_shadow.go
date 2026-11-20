package ui

import (
	"fmt"
	"image"
	"math"
)

type shadowSpec struct {
	width, height int
	alpha         []float32
}

func shadowCacheKey(w, h, r, blur, spread float32, c Color) string {
	return fmt.Sprintf("shadow:%g:%g:%g:%g:%g:%02x%02x%02x%02x", w, h, r, blur, spread, c.R, c.G, c.B, c.A)
}

func roundedRectCoverage(x, y, w, h, r float32) float32 {
	r = clampRadius(w, h, r)
	cx, cy := w/2, h/2
	qx := absf(x-cx) - (cx - r)
	qy := absf(y-cy) - (cy - r)
	if qx <= 0 && qy <= 0 {
		return 1
	}
	var distance float32
	if qx > 0 && qy > 0 {
		distance = float32(math.Hypot(float64(qx), float64(qy))) - r
	} else {
		distance = maxf(qx, qy) - r
	}
	return minf(1, maxf(0, 0.5-distance))
}

func newShadowSpec(width, height int, bw, bh, br, blur float32, c Color) *shadowSpec {
	spec := &shadowSpec{width: width, height: height, alpha: make([]float32, width*height)}
	const layers = 5
	layerAlpha := float32(c.A) / 255 / layers
	for i := layers; i >= 0; i-- {
		g := blur * float32(i) / layers
		x, y := blur-g, blur-g
		w, h := bw+g*2, bh+g*2
		r := br + g
		for py := 0; py < height; py++ {
			for px := 0; px < width; px++ {
				coverage := roundedRectCoverage(float32(px)+0.5-x, float32(py)+0.5-y, w, h, r)
				if coverage == 0 {
					continue
				}
				idx := py*width + px
				spec.alpha[idx] += coverage * layerAlpha * (1 - spec.alpha[idx])
			}
		}
	}
	return spec
}

func (s *shadowSpec) image(c Color) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, s.width, s.height))
	for i, alpha := range s.alpha {
		img.Pix[i*4] = c.R
		img.Pix[i*4+1] = c.G
		img.Pix[i*4+2] = c.B
		img.Pix[i*4+3] = uint8(alpha*255 + 0.5)
	}
	return img
}

func cachedShadow(w, h, r, blur, spread float32, c Color) (bitmap, Rect, bool) {
	bx, by := -spread, -spread
	bw, bh := w+spread*2, h+spread*2
	br := maxf(0, r+spread)
	x0 := float32(math.Floor(float64(bx - blur)))
	y0 := float32(math.Floor(float64(by - blur)))
	x1 := float32(math.Ceil(float64(bx + bw + blur)))
	y1 := float32(math.Ceil(float64(by + bh + blur)))
	width, height := int(x1-x0), int(y1-y0)
	if width <= 0 || height <= 0 || width*height > 16<<20 {
		return nil, Rect{}, false
	}

	key := shadowCacheKey(w, h, r, blur, spread, c)
	if cached, ok := lookupImage(key); ok {
		return cached, Rect{X: x0, Y: y0, W: float32(width), H: float32(height)}, true
	}
	spec := newShadowSpec(width, height, bw, bh, br, blur, c)
	bmp := backendNewBitmap(spec.image(c))
	storeImage(key, bmp)
	return bmp, Rect{X: x0, Y: y0, W: float32(width), H: float32(height)}, true
}
