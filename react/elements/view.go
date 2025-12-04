package elements

import (
	"github.com/sjm1327605995/tenon/react/api"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/core"
)

// View is a wrapper for a VNode to provide fluent-style methods.
type View struct {
	*core.VNode
}

// NewView creates a new VNode of type "View" and wraps it in a View struct.
func NewView() *View {
	vnode := core.NewVNode("View", nil)
	return &View{vnode}
}

// Style applies the given style to the View's VNode.
func (v *View) Style(style *styles.Style) *View {
	v.Props["style"] = style
	return v
}

// Child adds child components to the View's VNode.
// In the VDOM world, children are other VNodes.
func (v *View) Child(children ...api.Component) *View {
	for _, child := range children {
		// The Render method of a component should now return a VNode.
		// We'll need to adjust the Component interface and implementations.
		v.Children = append(v.Children, child.Render())
	}
	return v
}

// Render implements the api.Component interface.
// It returns the underlying VNode.
func (v *View) Render() *core.VNode {
	return v.VNode
}
