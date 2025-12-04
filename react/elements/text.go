package elements

import (
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/core"
)

// Text is a wrapper for a VNode to provide fluent-style methods.
type Text struct {
	*core.VNode
}

// NewText creates a new VNode of type "Text" and wraps it.
func NewText() *Text {
	vnode := core.NewVNode("Text", nil)
	return &Text{vnode}
}

// Style applies the given style to the Text's VNode.
func (t *Text) Style(style *styles.Style) *Text {
	t.Props["style"] = style
	return t
}

// Content sets the text content in the VNode's props.
func (t *Text) Content(content string) *Text {
	t.Props["content"] = content
	return t
}

// Render implements the api.Component interface.
// It returns the underlying VNode.
func (t *Text) Render() *core.VNode {
	return t.VNode
}
