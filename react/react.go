// Package react implements the core React-like UI framework for the Tenon project.
// This package provides the main functionality for creating, rendering, and managing UI components.
package react

import (
	"github.com/sjm1327605995/tenon/react/api"
	"github.com/sjm1327605995/tenon/react/api/renderer"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/elements"
	"github.com/sjm1327605995/tenon/react/event"
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
	root     api.Component
	renderer api.Renderer
	style    *styles.Style
	event    chan event.Event
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
func (h *ReactDOM) Render(children ...api.Component) error {
	element := elements.NewView().
		Style(h.style).
		Child(children...)
	h.root = element
	h.root.Render()
	h.renderer.SetElement(element)
	return h.renderer.Run()
}

// NewReactDOM creates a new ReactDOM instance with default settings.
// It initializes the renderer with a new Gio renderer, sets default dimensions of 800x600,
// and creates an event channel.
//
// Returns:
//   - A pointer to a newly created ReactDOM instance
func NewReactDOM() *ReactDOM {
	return &ReactDOM{
		renderer: renderer.NewGio(),
		style:    styles.NewStyle().Width(800).Height(600),
		event:    make(chan event.Event, 10),
	}
}
