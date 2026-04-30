package components

import (
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// Resizable is a container with a draggable resize handle.
type Resizable struct {
	core.BaseElement
	Content   *native.View
	handle    *native.View
	direction string // "horizontal" | "vertical"
	minSize   float32
	maxSize   float32
}

// NewResizable creates a resizable panel.
func NewResizable(Content core.Element, direction string) *Resizable {
	r := &Resizable{
		direction: direction,
		minSize:   100,
		maxSize:   800,
	}
	r.Init(r)

	if c, ok := Content.(*native.View); ok {
		r.Content = c
	} else {
		r.Content = native.NewView()
		r.Content.Add(Content)
	}

	r.handle = native.NewView()
	r.handle.SetBackgroundColor(core.GetTheme().BorderColor)
	if direction == "horizontal" {
		r.SetFlexDirection(yoga.FlexDirectionRow)
		r.Content.SetWidth(200)
		r.handle.SetWidth(4)
		r.handle.SetHeightPercent(100)
	} else {
		r.SetFlexDirection(yoga.FlexDirectionColumn)
		r.Content.SetHeight(200)
		r.handle.SetHeight(4)
		r.handle.SetWidthPercent(100)
	}

	r.Add(r.Content, r.handle)
	return r
}

// ElementType returns type identifier.
func (r *Resizable) ElementType() string { return "Resizable" }

// Panel returns the Content panel.
func (r *Resizable) Panel() *native.View { return r.Content }
