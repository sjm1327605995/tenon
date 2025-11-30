// Package elements implements the React-like element system for the Tenon framework.
// Elements are the building blocks of the UI that users directly interact with.
package elements

import (
	"github.com/sjm1327605995/tenon/react/api"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/components"
)

// View represents a container element that can hold child elements and apply styles.
// It extends the components.View with additional element-specific functionality.
//
// Fields:
//   - View: The underlying component view that handles layout and basic rendering
//   - Children: The list of child elements contained within this view
type View struct {
	*components.View
	Children []api.Element
}

// SetStyle applies the given style to this view element.
// This method delegates to the style's Apply method to apply all style properties.
//
// Parameters:
//   - style: The style configuration to apply to this view
func (v *View) SetStyle(style *styles.Style) {
	style.Apply(v)
}

// Rendering implements the rendering functionality using the provided renderer.
// It delegates the actual drawing to the renderer's DrawView method.
//
// Parameters:
//   - renderer: The renderer to use for drawing this view
func (v *View) Rendering(renderer api.Renderer) {
	renderer.DrawView(v.View)
}

// GetChildrenCount returns the number of child elements in this view.
// This is determined by querying the Yoga layout engine for child nodes.
//
// Returns:
//   - The number of child elements in the view
func (v *View) GetChildrenCount() int {
	return len(v.Yoga().GetChildren())
}

// GetChildAt returns the child element at the specified index.
// Returns nil if the index is out of bounds.
//
// Parameters:
//   - index: The index of the child element to retrieve
//
// Returns:
//   - The child element at the specified index, or nil if the index is invalid
func (v *View) GetChildAt(index int) api.Element {
	if index < 0 || index >= len(v.Children) {
		return nil
	}
	return v.Children[index]
}

// Render implements the api.Component interface and returns the view itself as a renderable node.
//
// Returns:
//   - The view as a renderable api.Node
func (v *View) Render() api.Node {
	return v
}

// Style applies the given style to this view and returns the view itself for method chaining.
// This is a convenience method that wraps SetStyle with a return value for fluent API usage.
//
// Parameters:
//   - option: The style configuration to apply to this view
//
// Returns:
//   - The view itself, allowing for method chaining
func (v *View) Style(option *styles.Style) *View {
	v.SetStyle(option)
	return v
}

// Child adds child components to this view and returns the view itself for method chaining.
// Each component is rendered and added to both the Yoga layout tree and the view's children list.
//
// Parameters:
//   - nodes: The components to add as children to this view
//
// Returns:
//   - The view itself, allowing for method chaining
func (v *View) Child(nodes ...api.Component) *View {
	for i := range nodes {
		element := nodes[i].Render()
		v.Yoga().InsertChild(element.Yoga(), uint32(i))
		v.Children = append(v.Children, element)
	}
	return v
}

// GetChildren returns all child elements of this view.
//
// Returns:
//   - A slice of all child elements contained within this view
func (v *View) GetChildren() []api.Element {
	return v.Children
}

// GetView returns the underlying component.View instance, used by the style system to set properties.
// This provides access to the core view implementation for systems that need direct access.
//
// Returns:
//   - The underlying component.View instance
func (v *View) GetView() *components.View {
	return v.View
}

// SetExtendedStyle applies extended style properties to this view element.
// This method handles specialized style types like BackgroundColor and BorderColor.
//
// Parameters:
//   - extendedStyle: The extended style to apply to this view
func (v *View) SetExtendedStyle(extendedStyle styles.IExtendedStyle) {
	switch e := extendedStyle.(type) {
	case styles.BackgroundColor:
		v.View.Background = e.Color
	case styles.BorderColor:
		v.View.BorderColor = e.Color
	}
}

// NewView creates a new View element with default properties.
//
// Returns:
//   - A pointer to a newly created View element
func NewView() *View {
	return &View{
		View: components.NewView(),
	}
}
