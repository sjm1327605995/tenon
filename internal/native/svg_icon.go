package native

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/svg"
	"github.com/sjm1327605995/tenon/internal/core"
)

// SVGIcon renders SVG path data as an icon.
type SVGIcon struct {
	core.BaseElement
	pathData string
	path     *vector.Path
	clr      color.Color
	size     float32
	viewW    float32
	viewH    float32
	drawW    float32
	drawH    float32
}

// NewSVGIcon creates an icon from SVG path data.
func NewSVGIcon(pathData string) *SVGIcon {
	si := &SVGIcon{
		pathData: pathData,
		clr:      color.RGBA{R: 255, G: 255, B: 255, A: 255},
		size:     14,
	}
	si.Init(si)
	si.SetWidth(si.size)
	si.SetHeight(si.size)
	si.parse()
	return si
}

// ElementType returns type identifier.
func (si *SVGIcon) ElementType() string { return "SVGIcon" }

// IsNative returns true as SVGIcon is a native rendering component.
func (si *SVGIcon) IsNative() bool { return true }

// SetColor sets the icon color.
func (si *SVGIcon) SetColor(clr color.Color) *SVGIcon {
	si.clr = clr
	si.Mark(core.FlagNeedDraw)
	return si
}

// SetSize sets the icon size.
func (si *SVGIcon) SetSize(size float32) *SVGIcon {
	si.size = size
	si.SetWidth(size)
	si.SetHeight(size)
	si.parse()
	si.Mark(core.FlagNeedDraw)
	return si
}

func (si *SVGIcon) parse() {
	if si.pathData == "" {
		return
	}
	minX, minY, maxX, maxY, err := svg.ParsePathBounds(si.pathData)
	if err != nil {
		return
	}
	si.viewW = maxX - minX
	si.viewH = maxY - minY
	if si.viewW <= 0 {
		si.viewW = 1
	}
	if si.viewH <= 0 {
		si.viewH = 1
	}
	scale := si.size / max(si.viewW, si.viewH)
	si.drawW = si.viewW * scale
	si.drawH = si.viewH * scale
	p, err := svg.ParsePathScaledAndShifted(si.pathData, scale, -minX*scale, -minY*scale)
	if err != nil {
		return
	}
	si.path = p
}

// Draw renders the SVG icon.
func (si *SVGIcon) Draw(screen *ebiten.Image) {
	if si.path == nil {
		return
	}
	bounds := si.GetBounds()

	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(si.clr)
	op.AntiAlias = true

	offsetX := bounds.X + (bounds.Width-si.drawW)/2
	offsetY := bounds.Y + (bounds.Height-si.drawH)/2

	imgW := int(math.Ceil(float64(si.drawW)))
	imgH := int(math.Ceil(float64(si.drawH)))
	if imgW < 1 {
		imgW = 1
	}
	if imgH < 1 {
		imgH = 1
	}

	img := ebiten.NewImage(imgW, imgH)
	vector.FillPath(img, si.path, &vector.FillOptions{}, op)

	drawOp := &ebiten.DrawImageOptions{}
	drawOp.GeoM.Translate(math.Round(float64(offsetX)), math.Round(float64(offsetY)))
	screen.DrawImage(img, drawOp)
}
