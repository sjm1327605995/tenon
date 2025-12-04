package elements

import (
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/core"
)

// Image is a wrapper for a VNode to provide fluent-style methods.
type Image struct {
	*core.VNode
}

// NewImage creates a new VNode of type "Image" and wraps it.
func NewImage() *Image {
	vnode := core.NewVNode("Image", nil)
	return &Image{vnode}
}

// Style applies the given style to the Image's VNode.
func (i *Image) Style(style *styles.Style) *Image {
	i.Props["style"] = style
	return i
}

// Source sets the image source path in the VNode's props.
func (i *Image) Source(path string) *Image {
	i.Props["source"] = path
	return i
}

// Render implements the api.Component interface.
// It returns the underlying VNode.
func (i *Image) Render() *core.VNode {
	return i.VNode
}
