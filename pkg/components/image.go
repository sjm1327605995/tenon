package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
)

type Image struct {
	core.BaseComponent
	Source    string
	ebitenImg *ebiten.Image
}

func NewImage() *Image {
	i := &Image{
		BaseComponent: core.NewBaseComponent(),
	}
	i.Init(i)
	return i
}

func (i *Image) Draw(screen *ebiten.Image) {
	element := i.Render()
	if element == nil || !element.Visible {
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
		placeholderColor := color.RGBA{R: 224, G: 224, B: 224, A: 255}
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, placeholderColor, false)
	}
}

func (i *Image) Update() error {
	return nil
}

func (i *Image) DrawOverlay(screen *ebiten.Image) {
}

func (i *Image) HandleInput() bool {
	return false
}

func (i *Image) SetSource(source string) *Image {
	i.Source = source
	return i
}

func (i *Image) SetEbitenImage(img *ebiten.Image) *Image {
	i.ebitenImg = img
	return i
}
