// Package react implements the core React-like UI framework for the Tenon project.
// This package provides the main functionality for creating, rendering, and managing UI components.
package react

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/sjm1327605995/tenon/react/api"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/elements"
	"github.com/sjm1327605995/tenon/react/event"
	"github.com/sjm1327605995/tenon/react/yoga"
	"image"
)

// ReactDOM is the main entry point for rendering React components.
// It manages the root component, renderer, styling, and event handling.
//
// Fields:
//   - root: The root component of the UI tree
//   - renderer: The renderer implementation used to draw components
//   - style: The default style applied to the root element
//   - event: Channel for handling UI events

type ReactDOM struct {
	window  *app.Window
	root    api.Component
	style   *styles.Style
	element api.Element
	size    image.Point
	event   chan event.Event
}

// Render renders the provided components into the ReactDOM.
// It creates a root view, applies the default style, adds the provided components as children,
// and starts the rendering process.
//
// Parameters:
//   - children: The components to render as children of the root view
//
// Returns:
//   - An error if there was an issue during rendering
func (r *ReactDOM) Render(children ...api.Component) error {
	element := elements.NewView().
		Style(r.style).
		Child(children...)
	r.element = element
	r.root = element
	r.root.Render()
	return r.run()
}

// NewReactDOM creates a new ReactDOM instance with default settings.
// It initializes the renderer with a new Gio renderer, sets default dimensions of 800x600,
// and creates an event channel.
//
// Returns:
//   - A pointer to a newly created ReactDOM instance
func NewReactDOM() *ReactDOM {
	return &ReactDOM{
		style: styles.NewStyle().Width(800).Height(600),
		event: make(chan event.Event, 10),
	}
}
func (r *ReactDOM) run() error {
	r.window = new(app.Window)
	var ops = new(op.Ops)
	for {
		switch e := r.window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.ConfigEvent:
			// Handle window size change
			if !e.Config.Size.Eq(r.size) && e.Config.Size.X > 0 {
				r.element.SetStyle(styles.NewStyle().
					Width(float32(e.Config.Size.X)).Height(float32(e.Config.Size.Y)))
				r.size = e.Config.Size
			}

		case app.FrameEvent:
			if e.Metric != elements.Metric {
				elements.Metric = e.Metric
			}
			// Create new layout context for rendering
			ctx := app.NewContext(ops, e)
			// Calculate layout
			r.element.Yoga().CalculateLayout(float32(e.Size.X), float32(e.Size.Y), yoga.DirectionInherit)
			// Start rendering
			r.draw(ctx, r.element)
			// Pass rendering operations to GPU
			e.Frame(ctx.Ops)
		}
	}
}
func (r *ReactDOM) draw(ctx layout.Context, element api.Element) layout.Dimensions {
	node := element.Yoga()
	x, y := int(node.LayoutLeft()), int(node.LayoutTop())
	size := image.Pt(x, y)
	// Set offset
	defer op.Offset(image.Pt(x, y)).Push(ctx.Ops).Pop()
	// Set constraints
	w, h := int(node.LayoutWidth()), int(node.LayoutHeight())
	ctx.Constraints.Max = image.Pt(w, h)
	// Call the element's rendering method
	element.Paint(ctx)
	// Recursively render child elements
	children := element.GetChildren()
	for _, child := range children {
		r.draw(ctx, child)
	}
	return layout.Dimensions{Size: size}
}
