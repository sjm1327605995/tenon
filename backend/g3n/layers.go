package g3n

import (
	"math"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

type layerPanel interface {
	gui.IPanel
	SetSize(width, height float32)
	SetRotationZ(float32)
	SetColor4(*math32.Color4) *gui.Panel
	Color4() math32.Color4
}

func (p *painter) applyLayer(start int, t ui.Layer) {
	children := p.root.Children()
	for i := start; i < len(children); i++ {
		panel, ok := children[i].(layerPanel)
		if !ok {
			continue
		}
		position := panel.Position()
		w, h := panel.Width(), panel.Height()
		cx, cy := position.X+w/2, position.Y+h/2
		dx, dy := cx-t.CX, cy-t.CY
		scaleY := t.ScaleY
		if scaleY == 0 {
			scaleY = t.Scale
		}
		sx, sy := dx*t.Scale, dy*scaleY
		rad := float64(t.Rotate) * math.Pi / 180
		cos, sin := float32(math.Cos(rad)), float32(math.Sin(rad))
		rx, ry := sx*cos-sy*sin, sx*sin+sy*cos
		cx, cy = t.CX+rx+t.TX, t.CY+ry+t.TY
		panel.SetPosition(cx-w*t.Scale/2, cy-h*scaleY/2)
		panel.SetSize(w*t.Scale, h*scaleY)
		panel.SetRotationZ(float32(-rad))
		color := panel.Color4()
		color.A *= t.Opacity
		panel.SetColor4(&color)
		if clipped, ok := panel.(clippedPanel); ok {
			clipped.transformClips(t)
		}
	}
}
