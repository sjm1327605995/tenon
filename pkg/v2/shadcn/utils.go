package shadcn

import (
	"image/color"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
)

// newColor converts color.Color to *render.Color.
func newColor(c color.Color) *render.Color {
	if c == nil {
		return nil
	}
	r, g, b, a := c.RGBA()
	return render.NewColor(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
}

// colorToRender converts color.Color to render.Color value type.
func colorToRender(c color.Color) render.Color {
	r, g, b, a := c.RGBA()
	return render.Color{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}
