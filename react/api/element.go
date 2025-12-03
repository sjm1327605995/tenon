// Package api provides core API interfaces and functionality for the React framework.
package api

import (
	"gioui.org/layout"
	"github.com/sjm1327605995/tenon/react/api/styles"
)

// Element is the fundamental interface for all renderable elements, providing necessary methods for third-party renderers.
// Types implementing this interface can be processed and displayed by the React rendering system.
type Element interface {
	styles.StyleElement
	// GetChildrenCount returns the number of child elements.
	GetChildrenCount() int
	// GetChildAt returns the child element at the specified index.
	// If the index is out of bounds, it may return nil or raise an error.
	GetChildAt(index int) Element

	Paint(ctx layout.Context)
	// GetChildren returns a list of all child elements of the element.
	GetChildren() []Element
	// SetExtendedStyle sets extended styles, with specific logic implemented by each element.
	// Extended styles provide custom style properties beyond standard styles.
	SetExtendedStyle(style styles.IExtendedStyle)
	// SetStyle sets the basic style of the element.
	SetStyle(style *styles.Style)
}
