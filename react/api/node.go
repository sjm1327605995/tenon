package api

import "github.com/sjm1327605995/tenon/react/core"

// Component is the primary interface for all Tenon components.
// Any struct that implements this interface can be used as a component.
type Component interface {
	// Render returns the Virtual DOM representation of the component.
	// This method is the core of a component, defining its UI structure.
	Render() *core.VNode
}

// Node is an alias for VNode, representing any node in the virtual DOM tree.
// It can be a component's output or a simple element.
type Node = *core.VNode
