// Package api provides core API interfaces and functionality for the React framework.
package api

import (
	"github.com/sjm1327605995/tenon/react/components"
)

// Renderer is the interface that third-party renderers need to implement, defining methods for rendering React components to specific platforms.
// Types implementing this interface can serve as rendering engines for React applications, responsible for actual graphics drawing.
type Renderer interface {
	// DrawView draws the view component to the rendering target.
	// This method is the core of the rendering process, responsible for actual graphics drawing.
	DrawView(view *components.View)
	// SetElement sets the root element to be rendered.
	// During rendering, the renderer will traverse the entire component tree starting from this root element.
	SetElement(element Element)
	// Run starts the main loop of the renderer, beginning to process rendering and events.
	// It returns an error if there are issues during startup.
	Run() error
}
