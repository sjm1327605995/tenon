package g3n

import (
	"math"

	"github.com/g3n/engine/gls"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

type clipData struct {
	clips []ui.Rect
}

func (c *clipData) setClips(clips []ui.Rect) {
	c.clips = append(c.clips[:0], clips...)
}

func (c *clipData) transformClips(t ui.Layer) {
	scaleY := t.ScaleY
	if scaleY == 0 {
		scaleY = t.Scale
	}
	for i, clip := range c.clips {
		x0, y0 := clip.X, clip.Y
		x1, y1 := x0+clip.W, y0+clip.H
		points := [4][2]float32{{x0, y0}, {x1, y0}, {x0, y1}, {x1, y1}}
		minX, minY := float32(1e20), float32(1e20)
		maxX, maxY := float32(-1e20), float32(-1e20)
		for _, point := range points {
			x, y := transformLayerPoint(point[0], point[1], t, scaleY)
			minX, minY = min(minX, x), min(minY, y)
			maxX, maxY = max(maxX, x), max(maxY, y)
		}
		c.clips[i] = ui.Rect{X: minX, Y: minY, W: maxX - minX, H: maxY - minY}
	}
}

func transformLayerPoint(x, y float32, t ui.Layer, scaleY float32) (float32, float32) {
	dx, dy := (x-t.CX)*t.Scale, (y-t.CY)*scaleY
	cos, sin := float32(math.Cos(float64(t.Rotate)*math.Pi/180)), float32(math.Sin(float64(t.Rotate)*math.Pi/180))
	return t.CX + dx*cos - dy*sin + t.TX, t.CY + dx*sin + dy*cos + t.TY
}

func (c *clipData) applyClip(gs *gls.GLS) {
	if len(c.clips) == 0 {
		gs.Disable(gls.SCISSOR_TEST)
		return
	}
	_, _, viewportW, viewportH := gs.GetViewport()
	left, top := float32(0), float32(0)
	right, bottom := float32(viewportW), float32(viewportH)
	for _, clip := range c.clips {
		left, top = max(left, clip.X), max(top, clip.Y)
		right, bottom = min(right, clip.X+clip.W), min(bottom, clip.Y+clip.H)
	}
	if right < left {
		right = left
	}
	if bottom < top {
		bottom = top
	}
	gs.Enable(gls.SCISSOR_TEST)
	gs.Scissor(int32(left), int32(float32(viewportH)-bottom), uint32(right-left), uint32(bottom-top))
}
