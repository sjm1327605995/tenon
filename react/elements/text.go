package elements

import (
	"github.com/sjm1327605995/tenon/react/api"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/components"
)

type Text struct {
	*components.Text
}

// SetStyle applies the given style to this view element.
// This method delegates to the style's Apply method to apply all style properties.
//
// Parameters:
//   - style: The style configuration to apply to this view
func (v *Text) SetStyle(style *styles.Style) {
	style.Apply(v)
}

func (v *Text) Rendering(renderer api.Renderer) {
	renderer.DrawText(v.Text)
}

func (v *Text) GetChildrenCount() int {
	return 0
}

// GetChildAt returns the child element at the specified index.
// Returns nil if the index is out of bounds.
//
// Parameters:
//   - index: The index of the child element to retrieve
//
// Returns:
//   - The child element at the specified index, or nil if the index is invalid
func (v *Text) GetChildAt(index int) api.Element {
	return nil
}

// Render implements the api.Component interface and returns the view itself as a renderable node.
//
// Returns:
//   - The view as a renderable api.Node
func (v *Text) Render() api.Node {
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
func (v *Text) Style(option *styles.Style) *Text {
	v.SetStyle(option)
	return v
}

func (v *Text) GetChildren() []api.Element {
	return nil
}

// SetExtendedStyle applies extended style properties to this view element.
// This method handles specialized style types like BackgroundColor and BorderColor.
//
// Parameters:
//   - extendedStyle: The extended style to apply to this view
func (v *Text) SetExtendedStyle(extendedStyle styles.IExtendedStyle) {
	//switch e := extendedStyle.(type) {
	//
	//}
}

func NewText() *Text {
	return &Text{
		Text: components.NewText(),
	}
}
