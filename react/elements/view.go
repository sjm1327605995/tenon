// Package elements implements the React-like element system for the Tenon framework.
// Elements are the building blocks of the UI that users directly interact with.
package elements

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/sjm1327605995/tenon/react/api"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/yoga"
	"image"
	"image/color"
)

// View represents a container element that can hold child elements and apply styles.
// It extends the components.View with additional element-specific functionality.
//
// Fields:
//   - View: The underlying component view that handles layout and basic rendering
//   - Children: The list of child elements contained within this view
type View struct {
	ElementBase
	Background  color.NRGBA // Background is the background color with alpha channel
	BorderColor color.NRGBA // BorderColor is the border color with alpha channel
}

func (v *View) Paint(ctx layout.Context) {

	//x, y := int(node.LayoutLeft()), int(node.LayoutTop())

	w, h := int(v.Node.LayoutWidth()), int(v.Node.LayoutHeight())
	size := image.Pt(w, h)

	// TODO: Currently all border sizes are consistent
	borderWidth := v.Node.LayoutBorder(yoga.EdgeLeft)

	whalf := (int(borderWidth) + 1) / 2
	// Draw background
	if v.Background.A > 0 {
		bodySize := size
		if borderWidth > 0 {
			bodySize.X -= whalf
			bodySize.Y -= whalf
		}
		paint.FillShape(ctx.Ops, v.Background, clip.Outline{
			Path: clip.RRect{
				Rect: image.Rectangle{Min: image.Pt(whalf, whalf), Max: bodySize},
				//SE:   b.radiusSE,
				//SW:   b.radiusSW,
				//NW:   v.radiusNW,
				//NE:   v.radiusNE,
			}.Path(ctx.Ops),
		}.Op())
	}
	// Draw border
	if borderWidth > 0 {
		paint.FillShape(ctx.Ops, v.BorderColor,
			clip.Stroke{
				Path: clip.RRect{
					Rect: image.Rect(whalf, whalf, size.X-whalf, size.Y-whalf),
					//SE:   v.radiusSE,
					//SW:   v.radiusSW,
					//NW:   v.radiusNW,
					//NE:   v.radiusNE,
				}.Path(ctx.Ops),
				Width: borderWidth,
			}.Op(),
		)
	}
}

// SetStyle applies the given style to this view element.
// This method delegates to the style's Apply method to apply all style properties.
//
// Parameters:
//   - style: The style configuration to apply to this view
func (v *View) SetStyle(style *styles.Style) {
	style.Apply(v)
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

// SetExtendedStyle applies extended style properties to this view element.
// This method handles specialized style types like BackgroundColor and BorderColor.
//
// Parameters:
//   - extendedStyle: The extended style to apply to this view
func (v *View) SetExtendedStyle(extendedStyle styles.IExtendedStyle) {
	switch e := extendedStyle.(type) {
	case styles.BackgroundColor:
		v.Background = e.Color
	case styles.BorderColor:
		v.BorderColor = e.Color
	}
}

// NewView creates a new View element with default properties.
//
// Returns:
//   - A pointer to a newly created View element
func NewView() *View {
	return &View{
		ElementBase: ElementBase{
			Node: yoga.NewNode(),
		},
	}
}
