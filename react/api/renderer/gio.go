package renderer

import (
	"image"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/sjm1327605995/tenon/react/api"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/components"
	"github.com/sjm1327605995/tenon/react/yoga"
)

type Gio struct {
	window  *app.Window
	element api.Element
	ctx     layout.Context
	size    image.Point
}

func (g *Gio) DrawView(view *components.View) {
	node := view.Yoga()
	//x, y := int(node.LayoutLeft()), int(node.LayoutTop())

	w, h := int(node.LayoutWidth()), int(node.LayoutHeight())
	size := image.Pt(w, h)

	//TODO 现在所有边框大小一致
	borderWidth := node.LayoutBorder(yoga.EdgeLeft)

	whalf := (int(borderWidth) + 1) / 2
	if view.Background.A > 0 {
		bodySize := size
		if borderWidth > 0 {
			bodySize.X -= whalf
			bodySize.Y -= whalf
		}
		paint.FillShape(g.ctx.Ops, view.Background, clip.Outline{
			Path: clip.RRect{
				Rect: image.Rectangle{Min: image.Pt(whalf, whalf), Max: bodySize},
				//SE:   b.radiusSE,
				//SW:   b.radiusSW,
				//NW:   v.radiusNW,
				//NE:   v.radiusNE,
			}.Path(g.ctx.Ops),
		}.Op())
	}
	if borderWidth > 0 {
		paint.FillShape(g.ctx.Ops, view.BorderColor,
			clip.Stroke{
				Path: clip.RRect{
					Rect: image.Rect(whalf, whalf, size.X-whalf, size.Y-whalf),
					//SE:   v.radiusSE,
					//SW:   v.radiusSW,
					//NW:   v.radiusNW,
					//NE:   v.radiusNE,
				}.Path(g.ctx.Ops),
				Width: borderWidth,
			}.Op(),
		)
	}
}

func NewGio() *Gio {
	return &Gio{
		window: new(app.Window),
	}
}
func (g *Gio) SetElement(element api.Element) {
	g.element = element
}
func (g *Gio) Run() error {
	var ops = new(op.Ops)
	for {
		switch e := g.window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.ConfigEvent:
			if !e.Config.Size.Eq(g.size) && e.Config.Size.X > 0 {
				g.element.SetStyle(styles.NewStyle().
					Width(float32(e.Config.Size.X)).Height(float32(e.Config.Size.Y)))
			}

		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			g.ctx = app.NewContext(ops, e)
			g.element.Yoga().CalculateLayout(yoga.Undefined, yoga.Undefined, yoga.DirectionInherit)
			g.Draw(g.ctx, g.element)
			// Pass the drawing operations to the GPU.
			e.Frame(g.ctx.Ops)
		}
	}
}
func (g *Gio) Draw(ctx layout.Context, element api.Element) layout.Dimensions {
	node := element.Yoga()
	x, y := int(node.LayoutLeft()), int(node.LayoutTop())
	size := image.Pt(x, y)
	defer op.Offset(image.Pt(x, y)).Push(ctx.Ops).Pop()
	w, h := int(node.LayoutWidth()), int(node.LayoutHeight())
	ctx.Constraints.Max = image.Pt(w, h)
	element.Rendering(g)
	children := element.GetChildren()
	for _, child := range children {
		g.Draw(ctx, child)

	}
	return layout.Dimensions{Size: size}
}
