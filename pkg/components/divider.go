package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Divider 是分割线组件。
type Divider struct {
	core.BaseHost
	thickness float32
	clr       color.Color
}

// NewDivider 创建一个分割线。
func NewDivider() *Divider {
	d := &Divider{
		thickness: 1,
		clr:       color.RGBA{R: 222, G: 226, B: 230, A: 255},
	}
	d.Init(d)
	d.GetElement().Yoga.StyleSetHeight(1)
	return d
}

// Draw 绘制分割线。
func (d *Divider) Draw(screen *ebiten.Image) {
	el := d.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := d.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	centerY := bounds.Y + bounds.Height/2 - d.thickness/2
	vector.FillRect(screen, bounds.X, centerY, bounds.Width, d.thickness, d.clr, false)
}

// ==================== 链式 API ====================

func (d *Divider) SetColor(clr color.Color) *Divider {
	d.clr = clr
	return d
}
func (d *Divider) SetThickness(thickness float32) *Divider {
	d.thickness = thickness
	return d
}
func (d *Divider) SetMargin(edge yoga.Edge, value float32) *Divider {
	d.GetElement().Yoga.StyleSetMargin(edge, value)
	return d
}
func (d *Divider) SetWidth(width float32) *Divider {
	d.GetElement().Yoga.StyleSetWidth(width)
	return d
}
func (d *Divider) SetFlexGrow(grow float32) *Divider {
	d.GetElement().Yoga.StyleSetFlexGrow(grow)
	return d
}

// SyncFrom 同步分割线属性。
func (d *Divider) SyncFrom(other core.Host) {
	if o, ok := other.(*Divider); ok {
		d.thickness = o.thickness
		d.clr = o.clr
	}
}
