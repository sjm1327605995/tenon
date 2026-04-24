package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Image displays an image.
type Image struct {
	core.BaseElement
	src      *ebiten.Image
	drawOpts *ebiten.DrawImageOptions
}

// NewImage creates an image element.
func NewImage() *Image {
	img := &Image{}
	img.Init(img)
	img.GetYoga().SetMeasureFunc(img.measure)
	return img
}

func (img *Image) measure(node *yoga.Node, width float32, widthMode yoga.MeasureMode, height float32, heightMode yoga.MeasureMode) yoga.Size {
	if img.src == nil {
		return yoga.Size{Width: 0, Height: 0}
	}
	b := img.src.Bounds()
	return yoga.Size{Width: float32(b.Dx()), Height: float32(b.Dy())}
}

// ElementType returns type identifier.
func (img *Image) ElementType() string { return "Image" }

// Draw renders the image.
func (img *Image) Draw(screen *ebiten.Image) {
	if !img.IsVisible() || img.src == nil {
		return
	}
	bounds := img.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	op := &ebiten.DrawImageOptions{}
	if img.drawOpts != nil {
		*op = *img.drawOpts
	}

	// Scale to fit bounds
	sw := float64(bounds.Width) / float64(img.src.Bounds().Dx())
	sh := float64(bounds.Height) / float64(img.src.Bounds().Dy())
	op.GeoM.Scale(sw, sh)

	// Apply transform (concatenates on the left: Transform * Scale)
	t := img.GetTransform()
	op.GeoM.Concat(core.BuildTransformGeoM(bounds, t))
	core.ApplyColorScaleAlpha(&op.ColorScale, t.Alpha)

	screen.DrawImage(img.src, op)
}

// SetEbitenImage sets the image source.
func (img *Image) SetEbitenImage(src *ebiten.Image) *Image {
	img.src = src
	img.Mark(core.FlagNeedDraw)
	return img
}

// SetSource loads image from file path.
func (img *Image) SetSource(path string) *Image {
	if loaded, _, err := ebitenutil.NewImageFromFile(path); err == nil {
		img.src = loaded
		img.Mark(core.FlagNeedMeasure | core.FlagNeedDraw)
	}
	return img
}

// SetDrawOptions sets custom draw options.
func (img *Image) SetDrawOptions(opts *ebiten.DrawImageOptions) *Image {
	img.drawOpts = opts
	return img
}

// Measure implements yoga measure for image aspect ratio.
func (img *Image) Measure() (width, height float32) {
	if img.src == nil {
		return 0, 0
	}
	b := img.src.Bounds()
	return float32(b.Dx()), float32(b.Dy())
}
