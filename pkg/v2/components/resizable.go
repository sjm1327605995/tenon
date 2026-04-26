package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Resizable is a container with a draggable resize handle.
type Resizable struct {
	core.BaseElement
	content   *View
	handle    *View
	direction string // "horizontal" | "vertical"
	minSize   float32
	maxSize   float32
}

// NewResizable creates a resizable panel.
func NewResizable(content core.Element, direction string) *Resizable {
	r := &Resizable{
		direction: direction,
		minSize:   100,
		maxSize:   800,
	}
	r.Init(r)

	if c, ok := content.(*View); ok {
		r.content = c
	} else {
		r.content = NewView()
		r.content.Add(content)
	}

	r.handle = NewView()
	r.handle.SetBackgroundColor(core.GetTheme().BorderColor)
	if direction == "horizontal" {
		r.SetFlexDirection(yoga.FlexDirectionRow)
		r.content.SetWidth(200)
		r.handle.SetWidth(4)
		r.handle.SetHeightPercent(100)
	} else {
		r.SetFlexDirection(yoga.FlexDirectionColumn)
		r.content.SetHeight(200)
		r.handle.SetHeight(4)
		r.handle.SetWidthPercent(100)
	}

	r.Add(r.content, r.handle)
	return r
}

// ElementType returns type identifier.
func (r *Resizable) ElementType() string { return "Resizable" }

// Panel returns the content panel.
func (r *Resizable) Panel() *View { return r.content }
