package g3n

import (
	"math"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

type labelPanel struct {
	*gui.Label
	color    math32.Color4
	uniModel gls.Uniform
	clipData
}

func (p *labelPanel) SetColor4(c *math32.Color4) *gui.Panel {
	p.color = *c
	p.Label.SetColor4(c)
	return p.Label.GetPanel()
}

func (p *labelPanel) Color4() math32.Color4 { return p.color }

func (p *labelPanel) RenderSetup(gs *gls.GLS, info *core.RenderInfo) {
	p.applyClip(gs)
	renderPanelTransform(gs, p.Label, p.Label.GetPanel(), &p.uniModel, info)
}

type clippedPanel interface {
	setClips([]ui.Rect)
	transformClips(ui.Layer)
}

type painter struct {
	root          *core.Node
	gs            *gls.GLS
	w, h          int
	layers        []int
	clips         []ui.Rect
	projection    ui.Layer
	hasProjection bool
}

func newPainter(root *core.Node, gs *gls.GLS, w, h int) *painter {
	return &painter{root: root, gs: gs, w: w, h: h}
}

func (p *painter) add(panel layerPanel) {
	if clipped, ok := panel.(clippedPanel); ok {
		clipped.setClips(p.clips)
	}
	p.root.Add(panel)
}

func color4(c ui.Color) *math32.Color4 {
	const scale = 1.0 / 255.0
	return &math32.Color4{R: float32(c.R) * scale, G: float32(c.G) * scale, B: float32(c.B) * scale, A: float32(c.A) * scale}
}

func (p *painter) panel(x, y, w, h float32, c ui.Color) *transformedPanel {
	panel := gui.NewPanel(w, h)
	panel.SetPosition(x, y)
	panel.SetColor4(color4(c))
	wrapped := withTransforms(panel)
	p.add(wrapped)
	return wrapped
}

func (p *painter) FillRect(x, y, w, h, r float32, c ui.Color) {
	if w <= 0 || h <= 0 || c.A == 0 {
		return
	}
	if r == 0 {
		p.panel(x, y, w, h, c)
		return
	}
	p.add(gradientPanel(x, y, w, h, r, c, c, 0))
}

func (p *painter) FillGradient(x, y, w, h, r float32, from, to ui.Color, angle float32) {
	if w > 0 && h > 0 {
		p.add(gradientPanel(x, y, w, h, r, from, to, angle))
	}
}

func (p *painter) StrokeRect(x, y, w, h, r float32, width float32, c ui.Color) {
	if w <= 0 || h <= 0 || width <= 0 || c.A == 0 {
		return
	}
	width = min(width, minf(w, h)/2)
	p.add(strokePanel(x, y, w, h, r, width, c))
}

func (p *painter) Shadow(x, y, w, h, r float32, offX, offY, blur, spread float32, c ui.Color) {
	steps := max(1, min(8, int(math.Ceil(float64(blur)/3))))
	for i := steps - 1; i >= 0; i-- {
		t := (float32(i) + 0.5) / float32(steps)
		expand := spread + blur*t
		layer := c
		layer.A = uint8(float32(c.A) * (1 - t) * 2 / float32(steps))
		p.FillRect(x+offX-expand, y+offY-expand, w+expand*2, h+expand*2, r+expand, layer)
	}
}

func (p *painter) Line(x0, y0, x1, y1 float32, c ui.Color) {
	distance := math.Hypot(float64(x1-x0), float64(y1-y0))
	steps := max(1, int(distance/4))
	for i := 0; i <= steps; i++ {
		t := float32(i) / float32(steps)
		p.FillRect(x0+(x1-x0)*t-0.5, y0+(y1-y0)*t-0.5, 1, 1, 0, c)
	}
}

func (p *painter) DrawText(s string, face ui.FontFace, c ui.Color, x, y float32, _, _ bool) {
	f, ok := face.(*font)
	if !ok || f == nil || s == "" {
		return
	}
	label := &labelPanel{Label: gui.NewLabelWithFont(s, f.font)}
	label.uniModel.Init("ModelMatrix")
	label.SetFontSize(float64(f.px))
	label.SetColor4(color4(c))
	label.SetPosition(x, y)
	p.add(label)
}

func (p *painter) DrawImage(img ui.Bitmap, dst ui.Rect, opacity float32) {
	b, ok := img.(*bitmap)
	if !ok || b == nil {
		return
	}
	opacity = max(0, min(1, opacity))
	alpha := uint8(opacity*255 + 0.5)
	panel := gui.NewImageFromTex(b.textureForOpacity(alpha).Incref())
	panel.SetPosition(dst.X, dst.Y)
	panel.SetSize(dst.W, dst.H)
	panel.SetColor4(color4(ui.Transparent))
	p.add(withTransforms(panel))
}

func (p *painter) FillPath(ui.VecPath, float32, float32, ui.Color)            {}
func (p *painter) StrokePath(ui.VecPath, float32, float32, float32, ui.Color) {}

func (p *painter) PushClip(r ui.Rect, _ float32) {
	p.clips = append(p.clips, r)
}

func (p *painter) PopClip() {
	if len(p.clips) > 0 {
		p.clips = p.clips[:len(p.clips)-1]
	}
}

func (p *painter) BeginLayer() { p.layers = append(p.layers, len(p.root.Children())) }

func (p *painter) EndLayer(t ui.LayerTransform) {
	if len(p.layers) == 0 {
		return
	}
	start := p.layers[len(p.layers)-1]
	p.layers = p.layers[:len(p.layers)-1]
	p.applyLayer(start, t.Snapshot())
}

func (p *painter) BeginProjected(t ui.LayerTransform) {
	p.projection = t.Snapshot()
	p.hasProjection = true
	p.BeginLayer()
}
func (p *painter) BeginProjectedElement(ui.Rect) {}
func (p *painter) EndProjectedElement()          {}
func (p *painter) EndProjected() {
	if !p.hasProjection {
		return
	}
	start := p.layers[len(p.layers)-1]
	p.layers = p.layers[:len(p.layers)-1]
	p.applyLayer(start, p.projection)
	p.hasProjection = false
}

func (p *painter) ProjectedImage(img ui.Bitmap, quad [4]ui.Point, opacity float32, _ string, _ *ui.Rect, _ float32) {
	minX, minY := quad[0].X, quad[0].Y
	maxX, maxY := minX, minY
	for _, point := range quad[1:] {
		minX, minY = min(minX, point.X), min(minY, point.Y)
		maxX, maxY = max(maxX, point.X), max(maxY, point.Y)
	}
	p.DrawImage(img, ui.Rect{X: minX, Y: minY, W: maxX - minX, H: maxY - minY}, opacity)
}

var _ ui.Painter = (*painter)(nil)
