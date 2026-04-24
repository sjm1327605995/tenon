package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Image 是图片宿主组件。
type Image struct {
	core.BaseHost
	Source    string
	ebitenImg *ebiten.Image
}

// NewImage 创建一个图片组件。
func NewImage() *Image {
	i := &Image{}
	i.Init(i)
	return i
}

// Draw 绘制图片。
func (i *Image) Draw(screen *ebiten.Image) {
	el := i.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := i.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	if i.ebitenImg != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(bounds.X), float64(bounds.Y))
		screen.DrawImage(i.ebitenImg, op)
	} else {
		placeholder := color.RGBA{R: 224, G: 224, B: 224, A: 255}
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, placeholder, false)
	}
}

// SetSource 设置图片路径（当前仅做标记，未实现加载逻辑）。
func (i *Image) SetSource(source string) *Image {
	i.Source = source
	return i
}

// SetEbitenImage 直接设置已加载的 Ebiten 图片。
func (i *Image) SetEbitenImage(img *ebiten.Image) *Image {
	i.ebitenImg = img
	return i
}

// SetWidth 设置图片宽度。
func (i *Image) SetWidth(width float32) *Image {
	i.GetElement().Yoga.StyleSetWidth(width)
	return i
}

// SetHeight 设置图片高度。
func (i *Image) SetHeight(height float32) *Image {
	i.GetElement().Yoga.StyleSetHeight(height)
	return i
}

// SetMargin 设置外边距。
func (i *Image) SetMargin(edge yoga.Edge, value float32) *Image {
	i.GetElement().Yoga.StyleSetMargin(edge, value)
	return i
}

// SyncFrom 同步图片属性。
func (i *Image) SyncFrom(other core.Host) {
	if o, ok := other.(*Image); ok {
		i.Source = o.Source
		i.ebitenImg = o.ebitenImg
	}
}
