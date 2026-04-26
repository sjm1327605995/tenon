package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

// AspectRatio maintains a specific width-to-height ratio for its content.
type AspectRatio struct {
	core.BaseElement
	ratio float32
	child core.Element
}

// NewAspectRatio creates an aspect ratio container.
func NewAspectRatio(ratio float32) *AspectRatio {
	a := &AspectRatio{ratio: ratio}
	a.Init(a)
	return a
}

// ElementType returns type identifier.
func (a *AspectRatio) ElementType() string { return "AspectRatio" }

// SetChild sets the content element.
func (a *AspectRatio) SetChild(child core.Element) *AspectRatio {
	if a.child != nil {
		a.RemoveChild(a.child)
	}
	a.child = child
	if child != nil {
		a.AppendChild(child)
	}
	return a
}

// Update adjusts height based on width to maintain ratio.
func (a *AspectRatio) Update() error {
	bounds := a.GetBounds()
	if bounds.Width > 0 && a.ratio > 0 {
		a.SetHeight(bounds.Width / a.ratio)
	}
	return nil
}
