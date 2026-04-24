package components

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/pkg/svg"
)

// SVGIcon 渲染 SVG path 数据的图标组件。
type SVGIcon struct {
	core.BaseHost
	pathData string
	path     *vector.Path
	clr      color.Color
	size     float32
	viewW    float32
	viewH    float32
	drawW    float32
	drawH    float32
}

// NewSVGIcon 从 SVG path 数据创建图标。
func NewSVGIcon(pathData string) *SVGIcon {
	si := &SVGIcon{
		pathData: pathData,
		clr:      color.RGBA{R: 255, G: 255, B: 255, A: 255},
		size:     14,
	}
	si.Init(si)
	si.GetElement().Yoga.StyleSetWidth(si.size)
	si.GetElement().Yoga.StyleSetHeight(si.size)
	si.parse()
	return si
}

// SetColor 设置图标颜色。
func (si *SVGIcon) SetColor(clr color.Color) *SVGIcon {
	si.clr = clr
	return si
}

// SetSize 设置图标大小。
func (si *SVGIcon) SetSize(size float32) *SVGIcon {
	si.size = size
	si.GetElement().Yoga.StyleSetWidth(size)
	si.GetElement().Yoga.StyleSetHeight(size)
	si.parse()
	return si
}

func (si *SVGIcon) parse() {
	if si.pathData == "" {
		return
	}
	// 第一遍：获取 bounds
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
	// 第二遍：按 size 缩放，并把 minX/minY 平移到 0
	scale := si.size / max(si.viewW, si.viewH)
	si.drawW = si.viewW * scale
	si.drawH = si.viewH * scale
	p, err := svg.ParsePathScaledAndShifted(si.pathData, scale, -minX*scale, -minY*scale)
	if err != nil {
		return
	}
	si.path = p
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// Draw 绘制 SVG 图标。
func (si *SVGIcon) Draw(screen *ebiten.Image) {
	if si.path == nil {
		return
	}
	bounds := si.GetLayoutBounds()

	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(si.clr)
	op.AntiAlias = true

	// path 已经缩放并平移到以 (0,0) 为起点
	// 计算在组件内的居中偏移
	offsetX := bounds.X + (bounds.Width-si.drawW)/2
	offsetY := bounds.Y + (bounds.Height-si.drawH)/2

	// 由于 vector.FillPath 没有 transform 支持，
	// 我们创建一个临时 Image 在上面绘制，然后用 DrawImage 偏移到正确位置。
	// 使用 Ceil 避免浮点截断导致 path 边缘被切掉
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

	// 将偏移量取整，避免亚像素渲染导致模糊/错位
	drawOp := &ebiten.DrawImageOptions{}
	drawOp.GeoM.Translate(math.Round(float64(offsetX)), math.Round(float64(offsetY)))
	screen.DrawImage(img, drawOp)
}

// SyncFrom 同步图标属性。
func (si *SVGIcon) SyncFrom(other core.Host) {
	if o, ok := other.(*SVGIcon); ok {
		si.pathData = o.pathData
		si.clr = o.clr
		si.size = o.size
		si.parse()
		si.GetElement().Yoga.StyleSetWidth(si.size)
		si.GetElement().Yoga.StyleSetHeight(si.size)
	}
}
