package g3n

import (
	"image"
	"image/color"
	"math"

	"github.com/g3n/engine/gui"
	"github.com/sjm1327605995/tenon/pkg/ui"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

func styledTextPanel(f *font, s string, c ui.Color, x, y float32, bold, italic bool) *transformedPanel {
	img := styledTextImage(f, s, bold, italic)
	tintTextImage(img, c)
	panel := gui.NewImageFromTex(newBitmap(img).texture.Incref())
	panel.SetPosition(x, y)
	panel.SetColor4(color4(ui.Transparent))
	return withTransforms(panel)
}

func tintTextImage(img *image.NRGBA, c ui.Color) {
	for i := 0; i < len(img.Pix); i += 4 {
		alpha := uint16(img.Pix[i+3])
		img.Pix[i] = c.R
		img.Pix[i+1] = c.G
		img.Pix[i+2] = c.B
		img.Pix[i+3] = uint8((alpha*uint16(c.A) + 127) / 255)
	}
}

func styledTextImage(f *font, s string, bold, italic bool) *image.NRGBA {
	source := f.font.DrawText(s)
	img := image.NewNRGBA(source.Bounds())
	xdraw.Draw(img, img.Bounds(), source, source.Bounds().Min, xdraw.Src)
	if bold {
		mask := image.NewUniform(color.Alpha{A: 140})
		for _, offset := range []image.Point{{X: 1}, {X: -1}, {Y: 1}, {Y: -1}, {X: 1, Y: 1}, {X: -1, Y: -1}} {
			xdraw.DrawMask(img, source.Bounds().Add(offset), source, image.Point{}, mask, image.Point{}, xdraw.Over)
		}
	}
	if !italic {
		return img
	}
	shear := math.Tan(0.21)
	baseline := float64(f.font.Metrics().Ascent.Ceil())
	padding := int(math.Ceil(shear*float64(img.Bounds().Dy()))) + 1
	output := image.NewNRGBA(image.Rect(0, 0, img.Bounds().Dx()+padding*2, img.Bounds().Dy()))
	matrix := f64.Aff3{1, -shear, float64(padding) + shear*baseline, 0, 1, 0}
	xdraw.ApproxBiLinear.Transform(output, matrix, img, img.Bounds(), xdraw.Over, nil)
	return output
}
