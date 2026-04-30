package native

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/internal/core"
)

// Divider is a horizontal or vertical dividing line.
type Divider struct {
	core.BaseElement
	thickness float32
	clr       color.Color
}

// NewDivider creates a divider.
func NewDivider() *Divider {
	theme := core.GetTheme()
	d := &Divider{
		thickness: 1,
		clr:       theme.DividerColor,
	}
	d.Init(d)
	d.SetHeight(1)
	return d
}

// ElementType returns type identifier.
func (d *Divider) ElementType() string { return "Divider" }

// Draw renders the divider.
func (d *Divider) Draw(screen *ebiten.Image) {
	bounds := d.GetBounds()
	centerY := bounds.Y + bounds.Height/2 - d.thickness/2
	vector.FillRect(screen, bounds.X, centerY, bounds.Width, d.thickness, d.clr, false)
}

// SetColor sets the divider color.
func (d *Divider) SetColor(clr color.Color) *Divider {
	d.clr = clr
	d.Mark(core.FlagNeedDraw)
	return d
}

// SetThickness sets the line thickness.
func (d *Divider) SetThickness(thickness float32) *Divider {
	d.thickness = thickness
	d.Mark(core.FlagNeedDraw)
	return d
}


