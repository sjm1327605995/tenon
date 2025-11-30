// Package components provides core component implementations for the React framework.
// This package contains fundamental components for building user interfaces, such as View.
package components

import (
	"image/color"

	"github.com/sjm1327605995/tenon/react/yoga"
)

// View is the most basic UI container component, similar to a div in HTML.
// It can contain other components and supports Flexbox layout, background colors, borders, and other style properties.
type View struct {
	yoga        *yoga.Node  // Yoga layout node for handling Flexbox layout
	Background  color.NRGBA // Background is the background color with alpha channel
	BorderColor color.NRGBA // BorderColor is the border color with alpha channel
}

// Yoga returns the View component's Yoga layout node.
// This method enables the View to implement the api.StyleElement interface, allowing styles to be applied.
func (v *View) Yoga() *yoga.Node {
	return v.yoga
}

// NewView creates and returns a new View component instance.
// It initializes the necessary Yoga layout node with default background and border colors.
func NewView() *View {
	return &View{
		yoga: yoga.NewNode(),
	}
}
