package core

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/sjm1327605995/tenon/react/yoga"
	"image"
)

// Engine manages the application routing and rendering
type Engine struct {
	Size         image.Point
	refreshAll   bool
	window       *app.Window
	routes       []Route
	currentRoute *Route
	currentPage  Component
	currentNode  Node
	updateChan   chan func(ctx layout.Context) // Channel for pending updates
}

// NewEngine creates a new engine instance
func NewEngine() *Engine {
	return &Engine{
		routes:      make([]Route, 0),
		currentNode: newEmptyNode(),
		updateChan:  make(chan func(ctx layout.Context), 1),
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

			select {
			// Handle window events

			// Handle pending updates
			case updateFn := <-e.updateChan:
				// Execute the update function
				updateFn(gtx)
				// Request a new frame to render the changes
				e.window.Invalidate()
			default:
				//if !e.Size.Eq(evt.Size) {
				e.Refresh(gtx, e.currentNode)
				e.currentNode.Yoga().
					CalculateLayout(float32(gtx.Constraints.Max.X), float32(gtx.Constraints.Max.Y), yoga.DirectionInherit)
				e.Size = evt.Size
				//}
				if gio := e.currentNode.Gio(); gio != nil {
					gio.Layout(gtx)
				}
			}

			// Render the frame
			evt.Frame(gtx.Ops)
		}

	}
}
func (e *Engine) AddRoute(r Route) {
	e.routes = append(e.routes, r)
	// Set the current page to the first route added
	if e.currentPage == nil {
		e.currentPage = r.ComponentFn()
	}
}
func (e *Engine) Refresh(gtx layout.Context, n Node) {
	n.Update(gtx)
	children := n.Children()
	for i := range children {
		e.Refresh(gtx, children[i])
	}
}

// Update enqueues an update to be processed in the main event loop
func (e *Engine) Update(updateFn func()) {
	// Wrap the update function to ensure re-rendering after update
	wrappedUpdate := func(ctx layout.Context) {
		// Execute the provided update function
		updateFn()
		// Re-render the current page to get the new node tree
		newNode := e.currentPage.Render()
		// Reconcile the new node with the old one
		e.currentNode = e.reconcile(e.currentNode, newNode)
	}
	// Enqueue the wrapped update function
	e.updateChan <- wrappedUpdate
}

// reconcile compares the old and new nodes and returns the reconciled node
func (e *Engine) reconcile(oldNode, newNode Node) Node {
	// If either node is nil, return the non-nil one
	if oldNode == nil {
		return newNode
	}
	if newNode == nil {
		return oldNode
	}

	// For now, we'll do a simple comparison and replace if different
	// In a more complete implementation, we'd compare node types and properties
	// This is a basic version - in React, this would be more sophisticated

	// Recursively reconcile children
	oldChildren := oldNode.Children()
	newChildren := newNode.Children()
	maxLen := len(oldChildren)
	if len(newChildren) > maxLen {
		maxLen = len(newChildren)
	}

	var reconciledChildren []Node
	for i := 0; i < maxLen; i++ {
		var oldChild, newChild Node
		if i < len(oldChildren) {
			oldChild = oldChildren[i]
		}
		if i < len(newChildren) {
			newChild = newChildren[i]
		}

		reconciledChild := e.reconcile(oldChild, newChild)
		if reconciledChild != nil {
			reconciledChildren = append(reconciledChildren, reconciledChild)
		}
	}

	// Check if the current component's shouldComponentUpdate method exists and returns true
	// If it does, we should use the new node
	if e.currentPage != nil {
		if e.currentPage.ShouldComponentUpdate() {
			// Use the new node since component should update
			return newNode
		}
	}

	// For now, always return the new node
	// In a real implementation, we'd compare and update instead of replacing
	return newNode
}
