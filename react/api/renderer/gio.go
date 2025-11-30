// Package renderer provides the Gio renderer implementation for the React framework.
// The Gio renderer is responsible for rendering React components to the Gio graphics library.
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

// Gio is a React renderer implementation based on the Gio graphics library.
// It implements the api.Renderer interface and is responsible for rendering React component trees to Gio windows.
type Gio struct {
	window  *app.Window    // Gio应用窗口
	element api.Element    // 根元素
	ctx     layout.Context // 布局上下文
	size    image.Point    // 窗口大小
}

// DrawView draws the View component to the Gio rendering context.
// It handles the rendering of visual properties such as background color and borders.
// Currently, border width is consistent for all edges (TODO: support different widths for different edges)
func (g *Gio) DrawView(view *components.View) {
	node := view.Yoga()
	//x, y := int(node.LayoutLeft()), int(node.LayoutTop())

	w, h := int(node.LayoutWidth()), int(node.LayoutHeight())
	size := image.Pt(w, h)

	// TODO: Currently all border sizes are consistent
	borderWidth := node.LayoutBorder(yoga.EdgeLeft)

	whalf := (int(borderWidth) + 1) / 2
	// Draw background
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
	// Draw border
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

// NewGio creates and returns a new Gio renderer instance.
// It initializes the necessary window resources.
func NewGio() *Gio {
	return &Gio{
		window: new(app.Window),
	}
}

// SetElement sets the root element to be rendered by the renderer.
// The element parameter is the root node of the React component tree.
func (g *Gio) SetElement(element api.Element) {
	g.element = element
}

// Run starts the main loop of the renderer, processing window events and performing rendering.
// It returns an error if there are issues during startup or operation.
func (g *Gio) Run() error {
	var ops = new(op.Ops)
	for {
		switch e := g.window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.ConfigEvent:
			// Handle window size change
			if !e.Config.Size.Eq(g.size) && e.Config.Size.X > 0 {
				g.element.SetStyle(styles.NewStyle().
					Width(float32(e.Config.Size.X)).Height(float32(e.Config.Size.Y)))
			}

		case app.FrameEvent:
			// Create new layout context for rendering
			g.ctx = app.NewContext(ops, e)
			// Calculate layout
			g.element.Yoga().CalculateLayout(yoga.Undefined, yoga.Undefined, yoga.DirectionInherit)
			// Start rendering
			g.Draw(g.ctx, g.element)
			// Pass rendering operations to GPU
			e.Frame(g.ctx.Ops)
		}
	}
}

// Draw recursively draws the element and its children.
// The ctx parameter is the Gio layout context, and element is the element to be drawn.
// It returns layout dimension information.
func (g *Gio) Draw(ctx layout.Context, element api.Element) layout.Dimensions {
	node := element.Yoga()
	x, y := int(node.LayoutLeft()), int(node.LayoutTop())
	size := image.Pt(x, y)
	// Set offset
	defer op.Offset(image.Pt(x, y)).Push(ctx.Ops).Pop()
	// Set constraints
	w, h := int(node.LayoutWidth()), int(node.LayoutHeight())
	ctx.Constraints.Max = image.Pt(w, h)
	// Call the element's rendering method
	element.Rendering(g)
	// Recursively render child elements
	children := element.GetChildren()
	for _, child := range children {
		g.Draw(ctx, child)
	}
	return layout.Dimensions{Size: size}
}
