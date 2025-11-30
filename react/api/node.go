// Package api provides core API interfaces and functionality for the React framework.
package api

// Node interface is the core representation of a React component, inheriting features from both Component and Element interfaces.
// Node represents a renderable component instance, containing rendering logic and DOM structure information.
type Node interface {
	Component
	Element
}

// Component interface defines the basic behavior of React components.
// Any type that implements this interface can be used as a React component.
type Component interface {
	// Render returns the rendering result of the component, which is a Node instance.
	// This method is the core of a React component, determining its final UI representation.
	Render() Node
}
