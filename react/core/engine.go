package core

import (
	"gioui.org/app"
	"gioui.org/op"
	"github.com/sjm1327605995/tenon/react/yoga"
)

// Engine manages the application routing and rendering
type Engine struct {
	window       *app.Window
	routes       []Route
	currentRoute *Route
	currentPage  Component
	currentNode  Node
}

// NewEngine creates a new engine instance
func NewEngine() *Engine {
	return &Engine{
		routes:      make([]Route, 0),
		currentNode: newEmptyNode(),
	}
}

// Navigate navigates to the specified path
func (e *Engine) Navigate(path string) bool {
	for _, route := range e.routes {
		if route.Path == path {
			// Unmount current page if exists
			if e.currentPage != nil {
				e.currentPage.ComponentWillUnmount()
			}

			// Create new component instance
			e.currentRoute = &route
			e.currentPage = route.ComponentFn()

			// Call component lifecycle methods in order
			e.currentPage.Constructor()
			// Note: GetDerivedStateFromProps, ShouldComponentUpdate would be called here in a real React implementation

			// Render the component to get DOM node
			e.currentNode = e.currentPage.Render()

			//// Calculate layout on the component's root node
			//componentNode.Yoga().CalculateLayout(0, 0, yoga.DirectionLTR)

			// Call componentDidMount after rendering is complete
			e.currentPage.ComponentDidMount()

			return true
		}
	}
	return false
}

// Run starts the Gio application with the specified window size
func (e *Engine) Run() error {
	// Create window if not already created
	if e.window == nil {
		e.window = new(app.Window)
	}
	var ops op.Ops
	for {
		switch evt := e.window.Event().(type) {
		case app.DestroyEvent:
			return evt.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, evt)

			// Generate drawing commands by laying out the current node
			if e.currentNode != nil {
				e.currentNode.Yoga().CalculateLayout(float32(gtx.Constraints.Max.X), float32(gtx.Constraints.Max.Y), yoga.DirectionInherit)
				e.currentNode.Update(gtx)
				e.currentNode.Gio().Layout(gtx)
			}

			// Render the frame
			evt.Frame(gtx.Ops)
		}
	}
}
func (e *Engine) AddRoute(r Route) {
	e.routes = append(e.routes, r)
}
