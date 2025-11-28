package renderer

import (
	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/sjm1327605995/tenon/react/common"
	"github.com/sjm1327605995/tenon/react/common/component"
	"github.com/sjm1327605995/tenon/react/yoga"
	"image"
	"image/color"
)

type Gio struct {
	window  *app.Window
	element common.Element
	size    image.Point
	Ops     *op.Ops
}

func (g *Gio) DrawView(view *component.View) {
	node := view.Yoga()
	x, y := int(node.LayoutLeft()), int(node.LayoutTop())
	w, h := int(node.LayoutWidth()), int(node.LayoutHeight())
	defer op.Offset(image.Pt(x, y)).Push(g.Ops).Pop()
	drawRedRect(g.Ops, image.Pt(w, h))

}

func drawRedRect(ops *op.Ops, point image.Point) {
	defer clip.Rect{Max: point}.Push(ops).Pop()
	paint.ColorOp{Color: color.NRGBA{R: 0x80, A: 0xFF}}.Add(ops)
	paint.PaintOp{}.Add(ops)
}
func NewGio() *Gio {
	return &Gio{
		Ops:    &op.Ops{},
		window: new(app.Window),
	}
}
func (g *Gio) SetElement(element common.Element) {
	g.element = element
}
func (g *Gio) Run() error {

	for {
		switch e := g.window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.ConfigEvent:
			if !e.Config.Size.Eq(g.size) && e.Config.Size.X > 0 {

			}

		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(g.Ops, e)
			g.element.Yoga().CalculateLayout(yoga.Undefined, yoga.Undefined, yoga.DirectionInherit)
			g.Draw(g.element)
			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
func (g *Gio) Draw(element common.Element) {
	element.Rendering(g)
	children := element.GetChildren()
	for _, child := range children {
		g.Draw(child)
	}
}
