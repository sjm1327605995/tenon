package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// DrawerSide defines which side the drawer slides from.
type DrawerSide int

const (
	DrawerLeft DrawerSide = iota
	DrawerRight
	DrawerTop
	DrawerBottom
)

// Drawer is a sliding panel that appears from an edge of the screen.
type Drawer struct {
	core.BaseElement
	panel     *native.View
	isOpen    bool
	side      DrawerSide
	width     float32
	height    float32
	maskColor color.Color
	onClose   func()
}

// NewDrawer creates a drawer.
func NewDrawer(side DrawerSide) *Drawer {
	d := &Drawer{
		side:      side,
		width:     300,
		height:    300,
		maskColor: color.RGBA{R: 0, G: 0, B: 0, A: 120},
	}
	d.Init(d)
	d.SetVisible(false)
	d.SetPositionType(yoga.PositionTypeAbsolute)
	d.SetPosition(yoga.EdgeLeft, 0)
	d.SetPosition(yoga.EdgeTop, 0)
	d.SetWidthPercent(100)
	d.SetHeightPercent(100)

	theme := core.GetTheme()
	d.panel = native.NewView()
	d.panel.SetBackgroundColor(theme.CardColor)
	d.panel.SetShadow(theme.ShadowColor, 16, 0, 4)
	d.panel.SetPadding(yoga.EdgeAll, 24)
	d.panel.SetFlexDirection(yoga.FlexDirectionColumn)
	d.Add(d.panel)
	return d
}

// ElementType returns type identifier.
func (d *Drawer) ElementType() string { return "Drawer" }

// Draw renders the backdrop mask.
func (d *Drawer) Draw(screen *ebiten.Image) {
	bounds := d.GetBounds()
	vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, d.maskColor, false)
}

// Open shows the drawer.
func (d *Drawer) Open() *Drawer {
	d.isOpen = true
	d.SetVisible(true)
	d.applyPanelLayout()
	d.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	if eng := d.GetEngine(); eng != nil {
		eng.AddOverlay(d)
	}
	return d
}

// Close hides the drawer.
func (d *Drawer) Close() *Drawer {
	d.isOpen = false
	d.SetVisible(false)
	if eng := d.GetEngine(); eng != nil {
		eng.RemoveOverlay(d)
	}
	if d.onClose != nil {
		d.onClose()
	}
	return d
}

func (d *Drawer) applyPanelLayout() {
	switch d.side {
	case DrawerLeft:
		d.panel.SetWidth(d.width)
		d.panel.SetHeightPercent(100)
	case DrawerRight:
		d.panel.SetWidth(d.width)
		d.panel.SetHeightPercent(100)
		d.panel.SetPositionType(yoga.PositionTypeAbsolute)
		d.panel.SetPosition(yoga.EdgeRight, 0)
		d.panel.SetPosition(yoga.EdgeTop, 0)
	case DrawerTop:
		d.panel.SetWidthPercent(100)
		d.panel.SetHeight(d.height)
	case DrawerBottom:
		d.panel.SetWidthPercent(100)
		d.panel.SetHeight(d.height)
		d.panel.SetPositionType(yoga.PositionTypeAbsolute)
		d.panel.SetPosition(yoga.EdgeLeft, 0)
		d.panel.SetPosition(yoga.EdgeBottom, 0)
	}
}

// HandleEvent processes mask click to close.
func (d *Drawer) HandleEvent(e *core.Event) bool {
	if !d.isOpen {
		return false
	}
	switch e.Type {
	case core.EventClick:
		pb := d.panel.GetBounds()
		if e.X < pb.X || e.X >= pb.X+pb.Width || e.Y < pb.Y || e.Y >= pb.Y+pb.Height {
			d.Close()
			return true
		}
	}
	return false
}

// Panel returns the inner Content container.
func (d *Drawer) Panel() *native.View { return d.panel }

// SetSize sets the drawer panel size.
func (d *Drawer) SetSize(size float32) *Drawer {
	if d.side == DrawerLeft || d.side == DrawerRight {
		d.width = size
	} else {
		d.height = size
	}
	d.applyPanelLayout()
	return d
}

// SetOnClose sets the close callback.
func (d *Drawer) SetOnClose(fn func()) *Drawer {
	d.onClose = fn
	return d
}
